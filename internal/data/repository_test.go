package data

import (
	"os"
	"path/filepath"
	"testing"
)

// fixtureCSV é um dataset pequeno e determinístico usado nos testes unitários.
// Último ano disponível = 2022.
//   total CO₂ (2022) = 300 + 600 + 50 = 950
//   média renováveis (2022) = (40 + 46 + 80) / 3 = 55,333... -> 55,33
//   países = 3
const fixtureCSV = `country,year,co2_total_mt,co2_per_capita_t,renewable_share_pct
Atlantis,2020,100.0,1.0,50.0
Atlantis,2021,200.0,2.0,60.0
Atlantis,2022,300.0,3.0,40.0
Brazil,2020,500.0,2.5,42.0
Brazil,2021,400.0,2.0,44.0
Brazil,2022,600.0,3.0,46.0
Zedland,2022,50.0,0.5,80.0
`

// writeFixture grava um CSV temporário e devolve o caminho.
func writeFixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fixture.csv")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("gravar fixture: %v", err)
	}
	return path
}

// loadFixture carrega o repositório a partir do fixture padrão.
func loadFixture(t *testing.T) *Repository {
	t.Helper()
	repo, err := LoadFromCSV(writeFixture(t, fixtureCSV))
	if err != nil {
		t.Fatalf("LoadFromCSV: %v", err)
	}
	return repo
}

func TestLoadFromCSV_OK(t *testing.T) {
	repo := loadFixture(t)
	if got := len(repo.all); got != 7 {
		t.Errorf("registros carregados = %d, esperado 7", got)
	}
	if got := len(repo.Countries()); got != 3 {
		t.Errorf("países = %d, esperado 3", got)
	}
}

func TestLoadFromCSV_FileNotFound(t *testing.T) {
	if _, err := LoadFromCSV(filepath.Join(t.TempDir(), "nao-existe.csv")); err == nil {
		t.Fatal("esperado erro ao abrir arquivo inexistente, got nil")
	}
}

func TestLoadFromCSV_BadHeader(t *testing.T) {
	bad := "pais,ano,co2,a,b\nBrazil,2022,1,1,1\n"
	if _, err := LoadFromCSV(writeFixture(t, bad)); err == nil {
		t.Fatal("esperado erro de cabeçalho inválido, got nil")
	}
}

func TestLoadFromCSV_BadRow(t *testing.T) {
	bad := "country,year,co2_total_mt,co2_per_capita_t,renewable_share_pct\nBrazil,XXXX,1,1,1\n"
	if _, err := LoadFromCSV(writeFixture(t, bad)); err == nil {
		t.Fatal("esperado erro de parsing de linha (year inválido), got nil")
	}
}

func TestParseRow_OK(t *testing.T) {
	ind, err := parseRow([]string{"Brazil", "2022", "600.0", "3.0", "46.0"})
	if err != nil {
		t.Fatalf("parseRow: %v", err)
	}
	if ind.Country != "Brazil" || ind.Year != 2022 || ind.CO2TotalMt != 600.0 ||
		ind.CO2PerCapitaT != 3.0 || ind.RenewableShare != 46.0 {
		t.Errorf("parseRow devolveu %+v", ind)
	}
}

func TestParseRow_TrimsSpaces(t *testing.T) {
	ind, err := parseRow([]string{" Brazil ", " 2022 ", " 600.0 ", " 3.0 ", " 46.0 "})
	if err != nil {
		t.Fatalf("parseRow: %v", err)
	}
	if ind.Country != "Brazil" || ind.Year != 2022 {
		t.Errorf("espaços não removidos: %+v", ind)
	}
}

func TestCountries_Sorted(t *testing.T) {
	repo := loadFixture(t)
	got := repo.Countries()
	want := []string{"Atlantis", "Brazil", "Zedland"}
	if len(got) != len(want) {
		t.Fatalf("Countries() = %v, esperado %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Countries()[%d] = %q, esperado %q", i, got[i], want[i])
		}
	}
}

func TestYearsAndLatest(t *testing.T) {
	repo := loadFixture(t)
	years := repo.Years()
	if len(years) != 3 || years[0] != 2020 || years[2] != 2022 {
		t.Errorf("Years() = %v, esperado [2020 2021 2022]", years)
	}
	if repo.LatestYear() != 2022 {
		t.Errorf("LatestYear() = %d, esperado 2022", repo.LatestYear())
	}
}

func TestSummary(t *testing.T) {
	s := loadFixture(t).Summary()
	if s.LatestYear != 2022 {
		t.Errorf("LatestYear = %d, esperado 2022", s.LatestYear)
	}
	if s.TotalCountries != 3 {
		t.Errorf("TotalCountries = %d, esperado 3", s.TotalCountries)
	}
	if s.GlobalCO2Mt != 950.0 {
		t.Errorf("GlobalCO2Mt = %v, esperado 950", s.GlobalCO2Mt)
	}
	if s.AvgRenewableShare != 55.33 {
		t.Errorf("AvgRenewableShare = %v, esperado 55.33", s.AvgRenewableShare)
	}
}

func TestTop(t *testing.T) {
	repo := loadFixture(t)
	top := repo.Top(2)
	if len(top) != 2 {
		t.Fatalf("Top(2) devolveu %d itens, esperado 2", len(top))
	}
	if top[0].Country != "Brazil" || top[0].CO2TotalMt != 600.0 {
		t.Errorf("Top[0] = %+v, esperado Brazil/600", top[0])
	}
	if top[1].Country != "Atlantis" || top[1].CO2TotalMt != 300.0 {
		t.Errorf("Top[1] = %+v, esperado Atlantis/300", top[1])
	}
}

func TestTop_DefaultsAndClamp(t *testing.T) {
	repo := loadFixture(t)
	// n <= 0 deve assumir 5 (UC05 A1). Há apenas 3 países no último ano,
	// então o resultado é truncado para 3.
	if got := len(repo.Top(0)); got != 3 {
		t.Errorf("Top(0) devolveu %d, esperado 3 (default 5 truncado)", got)
	}
	if got := len(repo.Top(-1)); got != 3 {
		t.Errorf("Top(-1) devolveu %d, esperado 3", got)
	}
	if got := len(repo.Top(100)); got != 3 {
		t.Errorf("Top(100) devolveu %d, esperado 3 (clamp ao total)", got)
	}
}

func TestSeriesFor_Found(t *testing.T) {
	repo := loadFixture(t)
	s, ok := repo.SeriesFor("Brazil", 0, 0)
	if !ok {
		t.Fatal("SeriesFor(Brazil) ok=false, esperado true")
	}
	if len(s.Points) != 3 {
		t.Errorf("pontos = %d, esperado 3", len(s.Points))
	}
	if s.Points[0].Year != 2020 || s.Points[2].Year != 2022 {
		t.Errorf("série fora de ordem: %+v", s.Points)
	}
}

func TestSeriesFor_NotFound(t *testing.T) {
	repo := loadFixture(t)
	if _, ok := repo.SeriesFor("Narnia", 0, 0); ok {
		t.Error("SeriesFor(Narnia) ok=true, esperado false")
	}
}

func TestSeriesFor_YearRange(t *testing.T) {
	repo := loadFixture(t)
	s, _ := repo.SeriesFor("Brazil", 2021, 2022)
	if len(s.Points) != 2 {
		t.Fatalf("intervalo [2021,2022] devolveu %d pontos, esperado 2", len(s.Points))
	}
	for _, p := range s.Points {
		if p.Year < 2021 || p.Year > 2022 {
			t.Errorf("ano %d fora do intervalo", p.Year)
		}
	}

	// from apenas
	if s, _ := repo.SeriesFor("Brazil", 2022, 0); len(s.Points) != 1 || s.Points[0].Year != 2022 {
		t.Errorf("from=2022 devolveu %+v", s.Points)
	}
	// to apenas
	if s, _ := repo.SeriesFor("Brazil", 0, 2020); len(s.Points) != 1 || s.Points[0].Year != 2020 {
		t.Errorf("to=2020 devolveu %+v", s.Points)
	}
}

func TestSeriesFor_DoesNotMutateRepository(t *testing.T) {
	repo := loadFixture(t)
	// Recorta a série e altera a cópia retornada; o repositório não pode mudar.
	s, _ := repo.SeriesFor("Brazil", 0, 0)
	if len(s.Points) > 0 {
		s.Points[0].CO2TotalMt = -999
	}
	again, _ := repo.SeriesFor("Brazil", 0, 0)
	if again.Points[0].CO2TotalMt == -999 {
		t.Error("mutação na série retornada afetou o repositório")
	}
}

func TestSeriesForMany_IgnoresUnknown(t *testing.T) {
	repo := loadFixture(t)
	out := repo.SeriesForMany([]string{"Brazil", "Narnia", "Atlantis"}, 0, 0)
	if len(out) != 2 {
		t.Fatalf("SeriesForMany devolveu %d séries, esperado 2 (Narnia ignorado)", len(out))
	}
}

func TestEmptyRepository(t *testing.T) {
	onlyHeader := "country,year,co2_total_mt,co2_per_capita_t,renewable_share_pct\n"
	repo, err := LoadFromCSV(writeFixture(t, onlyHeader))
	if err != nil {
		t.Fatalf("LoadFromCSV (vazio): %v", err)
	}
	if repo.LatestYear() != 0 {
		t.Errorf("LatestYear vazio = %d, esperado 0", repo.LatestYear())
	}
	s := repo.Summary()
	if s.TotalCountries != 0 || s.GlobalCO2Mt != 0 {
		t.Errorf("Summary vazio = %+v, esperado zeros", s)
	}
	if len(repo.Top(5)) != 0 {
		t.Error("Top em repositório vazio deveria ser vazio")
	}
}
