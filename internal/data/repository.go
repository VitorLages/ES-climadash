package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"climadash/internal/models"
)

// Repository é um repositório em memória dos indicadores carregados do CSV.
type Repository struct {
	all    []models.Indicator
	byCtry map[string][]models.Indicator
	years  []int
}

// LoadFromCSV carrega o dataset a partir de um arquivo CSV no caminho informado.
func LoadFromCSV(path string) (*Repository, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("abrir CSV %q: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 5
	r.TrimLeadingSpace = true

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("ler cabeçalho: %w", err)
	}
	expected := []string{"country", "year", "co2_total_mt", "co2_per_capita_t", "renewable_share_pct"}
	for i, col := range expected {
		if strings.TrimSpace(strings.ToLower(header[i])) != col {
			return nil, fmt.Errorf("cabeçalho inesperado: coluna %d = %q, esperado %q", i, header[i], col)
		}
	}

	var all []models.Indicator
	for line := 2; ; line++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("linha %d: %w", line, err)
		}
		ind, err := parseRow(row)
		if err != nil {
			return nil, fmt.Errorf("linha %d: %w", line, err)
		}
		all = append(all, ind)
	}

	repo := &Repository{
		byCtry: make(map[string][]models.Indicator),
	}
	yearSet := make(map[int]struct{})
	for _, ind := range all {
		repo.byCtry[ind.Country] = append(repo.byCtry[ind.Country], ind)
		yearSet[ind.Year] = struct{}{}
	}
	for _, series := range repo.byCtry {
		sort.Slice(series, func(i, j int) bool { return series[i].Year < series[j].Year })
	}
	for y := range yearSet {
		repo.years = append(repo.years, y)
	}
	sort.Ints(repo.years)
	repo.all = all

	return repo, nil
}

func parseRow(row []string) (models.Indicator, error) {
	year, err := strconv.Atoi(strings.TrimSpace(row[1]))
	if err != nil {
		return models.Indicator{}, fmt.Errorf("year inválido %q: %w", row[1], err)
	}
	total, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return models.Indicator{}, fmt.Errorf("co2_total_mt inválido %q: %w", row[2], err)
	}
	perCap, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil {
		return models.Indicator{}, fmt.Errorf("co2_per_capita_t inválido %q: %w", row[3], err)
	}
	ren, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil {
		return models.Indicator{}, fmt.Errorf("renewable_share_pct inválido %q: %w", row[4], err)
	}
	return models.Indicator{
		Country:        strings.TrimSpace(row[0]),
		Year:           year,
		CO2TotalMt:     total,
		CO2PerCapitaT:  perCap,
		RenewableShare: ren,
	}, nil
}

// Countries retorna a lista de países ordenada alfabeticamente.
func (r *Repository) Countries() []string {
	out := make([]string, 0, len(r.byCtry))
	for c := range r.byCtry {
		out = append(out, c)
	}
	sort.Strings(out)
	return out
}

// Years retorna os anos disponíveis (já ordenados).
func (r *Repository) Years() []int {
	out := make([]int, len(r.years))
	copy(out, r.years)
	return out
}

// LatestYear retorna o último ano com dados.
func (r *Repository) LatestYear() int {
	if len(r.years) == 0 {
		return 0
	}
	return r.years[len(r.years)-1]
}

// SeriesFor retorna a série histórica de um país, opcionalmente recortada por [from,to].
// Retorna (Series{}, false) se o país não existir.
func (r *Repository) SeriesFor(country string, from, to int) (models.Series, bool) {
	points, ok := r.byCtry[country]
	if !ok {
		return models.Series{}, false
	}
	filtered := filterByYear(points, from, to)
	out := make([]models.Indicator, len(filtered))
	copy(out, filtered)
	return models.Series{Country: country, Points: out}, true
}

// SeriesForMany retorna séries de vários países, ignorando os inexistentes.
func (r *Repository) SeriesForMany(countries []string, from, to int) []models.Series {
	out := make([]models.Series, 0, len(countries))
	for _, c := range countries {
		s, ok := r.SeriesFor(c, from, to)
		if !ok {
			continue
		}
		out = append(out, s)
	}
	return out
}

// Summary devolve a visão geral.
func (r *Repository) Summary() models.Summary {
	latest := r.LatestYear()
	var total float64
	var sumRen float64
	var count int
	for _, ind := range r.all {
		if ind.Year == latest {
			total += ind.CO2TotalMt
			sumRen += ind.RenewableShare
			count++
		}
	}
	avg := 0.0
	if count > 0 {
		avg = sumRen / float64(count)
	}
	return models.Summary{
		LatestYear:        latest,
		TotalCountries:    len(r.byCtry),
		GlobalCO2Mt:       round2(total),
		AvgRenewableShare: round2(avg),
		YearsCovered:      r.Years(),
	}
}

// Top retorna os N maiores emissores no último ano disponível.
func (r *Repository) Top(n int) []models.TopEmitter {
	if n <= 0 {
		n = 5
	}
	latest := r.LatestYear()
	var list []models.TopEmitter
	for _, ind := range r.all {
		if ind.Year != latest {
			continue
		}
		list = append(list, models.TopEmitter{
			Country:    ind.Country,
			Year:       ind.Year,
			CO2TotalMt: ind.CO2TotalMt,
		})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CO2TotalMt > list[j].CO2TotalMt })
	if n > len(list) {
		n = len(list)
	}
	return list[:n]
}

func filterByYear(in []models.Indicator, from, to int) []models.Indicator {
	if from == 0 && to == 0 {
		return in
	}
	out := in[:0:0]
	for _, ind := range in {
		if from != 0 && ind.Year < from {
			continue
		}
		if to != 0 && ind.Year > to {
			continue
		}
		out = append(out, ind)
	}
	return out
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
