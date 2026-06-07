package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"climadash/internal/data"
)

const fixtureCSV = `country,year,co2_total_mt,co2_per_capita_t,renewable_share_pct
Atlantis,2020,100.0,1.0,50.0
Atlantis,2021,200.0,2.0,60.0
Atlantis,2022,300.0,3.0,40.0
Brazil,2020,500.0,2.5,42.0
Brazil,2021,400.0,2.0,44.0
Brazil,2022,600.0,3.0,46.0
Zedland,2022,50.0,0.5,80.0
`

// newTestServer monta um servidor httptest com o fixture padrão.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fixture.csv")
	if err := os.WriteFile(path, []byte(fixtureCSV), 0o600); err != nil {
		t.Fatalf("gravar fixture: %v", err)
	}
	repo, err := data.LoadFromCSV(path)
	if err != nil {
		t.Fatalf("LoadFromCSV: %v", err)
	}
	mux := http.NewServeMux()
	New(repo).Register(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// getJSON faz GET e decodifica o corpo JSON em out; devolve o status.
func getJSON(t *testing.T, url string, out any) int {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("decode %s: %v", url, err)
		}
	}
	return resp.StatusCode
}

func TestHealth(t *testing.T) {
	srv := newTestServer(t)
	var body map[string]string
	if code := getJSON(t, srv.URL+"/api/health", &body); code != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", code)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, esperado ok", body["status"])
	}
}

func TestCountriesEndpoint(t *testing.T) {
	srv := newTestServer(t)
	var list []string
	code := getJSON(t, srv.URL+"/api/countries", &list)
	if code != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", code)
	}
	if len(list) != 3 || list[0] != "Atlantis" {
		t.Errorf("countries = %v, esperado [Atlantis Brazil Zedland]", list)
	}
}

func TestYearsEndpoint(t *testing.T) {
	srv := newTestServer(t)
	var years []int
	getJSON(t, srv.URL+"/api/years", &years)
	if len(years) != 3 || years[2] != 2022 {
		t.Errorf("years = %v, esperado [2020 2021 2022]", years)
	}
}

func TestSummaryEndpoint(t *testing.T) {
	srv := newTestServer(t)
	var s struct {
		LatestYear        int     `json:"latest_year"`
		TotalCountries    int     `json:"total_countries"`
		GlobalCO2Mt       float64 `json:"global_co2_mt"`
		AvgRenewableShare float64 `json:"avg_renewable_share_pct"`
	}
	getJSON(t, srv.URL+"/api/summary", &s)
	if s.LatestYear != 2022 || s.TotalCountries != 3 || s.GlobalCO2Mt != 950 || s.AvgRenewableShare != 55.33 {
		t.Errorf("summary = %+v", s)
	}
}

func TestCountryDetail_OK(t *testing.T) {
	srv := newTestServer(t)
	var s struct {
		Country string `json:"country"`
		Points  []struct {
			Year int `json:"year"`
		} `json:"points"`
	}
	code := getJSON(t, srv.URL+"/api/countries/Brazil", &s)
	if code != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", code)
	}
	if s.Country != "Brazil" || len(s.Points) != 3 {
		t.Errorf("detalhe = %+v", s)
	}
}

func TestCountryDetail_NotFound(t *testing.T) {
	srv := newTestServer(t)
	if code := getJSON(t, srv.URL+"/api/countries/Narnia", nil); code != http.StatusNotFound {
		t.Errorf("status = %d, esperado 404", code)
	}
}

func TestCountryDetail_YearRange(t *testing.T) {
	srv := newTestServer(t)
	var s struct {
		Points []struct {
			Year int `json:"year"`
		} `json:"points"`
	}
	getJSON(t, srv.URL+"/api/countries/Brazil?from=2021&to=2022", &s)
	if len(s.Points) != 2 {
		t.Errorf("intervalo devolveu %d pontos, esperado 2", len(s.Points))
	}
}

func TestCountryDetail_InvalidRange(t *testing.T) {
	srv := newTestServer(t)
	if code := getJSON(t, srv.URL+"/api/countries/Brazil?from=2022&to=2020", nil); code != http.StatusBadRequest {
		t.Errorf("from>to: status = %d, esperado 400", code)
	}
}

func TestEmissions_OK(t *testing.T) {
	srv := newTestServer(t)
	var series []struct {
		Country string `json:"country"`
	}
	code := getJSON(t, srv.URL+"/api/emissions?countries=Brazil,Atlantis", &series)
	if code != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", code)
	}
	if len(series) != 2 {
		t.Errorf("séries = %d, esperado 2", len(series))
	}
}

func TestEmissions_MissingParam(t *testing.T) {
	srv := newTestServer(t)
	if code := getJSON(t, srv.URL+"/api/emissions", nil); code != http.StatusBadRequest {
		t.Errorf("sem 'countries': status = %d, esperado 400", code)
	}
}

func TestEmissions_TooMany(t *testing.T) {
	srv := newTestServer(t)
	url := srv.URL + "/api/emissions?countries=a,b,c,d,e,f"
	if code := getJSON(t, url, nil); code != http.StatusBadRequest {
		t.Errorf("6 países: status = %d, esperado 400", code)
	}
}

func TestEmissions_IgnoresUnknown(t *testing.T) {
	srv := newTestServer(t)
	var series []struct {
		Country string `json:"country"`
	}
	getJSON(t, srv.URL+"/api/emissions?countries=Brazil,Narnia", &series)
	if len(series) != 1 || series[0].Country != "Brazil" {
		t.Errorf("séries = %+v, esperado apenas Brazil", series)
	}
}

func TestTop_OK(t *testing.T) {
	srv := newTestServer(t)
	var top []struct {
		Country    string  `json:"country"`
		CO2TotalMt float64 `json:"co2_total_mt"`
	}
	getJSON(t, srv.URL+"/api/top?n=2", &top)
	if len(top) != 2 {
		t.Fatalf("top = %d, esperado 2", len(top))
	}
	if top[0].Country != "Brazil" {
		t.Errorf("primeiro = %q, esperado Brazil", top[0].Country)
	}
}

func TestTop_InvalidNDefaultsTo5(t *testing.T) {
	srv := newTestServer(t)
	var top []struct {
		Country string `json:"country"`
	}
	// n inválido -> default 5 (UC05 A1); há 3 países -> 3 itens.
	getJSON(t, srv.URL+"/api/top?n=abc", &top)
	if len(top) != 3 {
		t.Errorf("n=abc devolveu %d, esperado 3 (default 5 truncado)", len(top))
	}
}

func TestMethodNotAllowed(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Post(srv.URL+"/api/summary", "application/json", nil)
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("POST /api/summary: status = %d, esperado 405", resp.StatusCode)
	}
}
