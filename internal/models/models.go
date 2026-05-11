package models

// Indicator é o registro de um indicador climático para um par (país, ano).
type Indicator struct {
	Country         string  `json:"country"`
	Year            int     `json:"year"`
	CO2TotalMt      float64 `json:"co2_total_mt"`
	CO2PerCapitaT   float64 `json:"co2_per_capita_t"`
	RenewableShare  float64 `json:"renewable_share_pct"`
}

// Series representa a série histórica de um país.
type Series struct {
	Country string      `json:"country"`
	Points  []Indicator `json:"points"`
}

// Summary é a visão geral exibida na home do dashboard.
type Summary struct {
	LatestYear       int     `json:"latest_year"`
	TotalCountries   int     `json:"total_countries"`
	GlobalCO2Mt      float64 `json:"global_co2_mt"`
	AvgRenewableShare float64 `json:"avg_renewable_share_pct"`
	YearsCovered     []int   `json:"years_covered"`
}

// TopEmitter é um item da lista de maiores emissores.
type TopEmitter struct {
	Country    string  `json:"country"`
	Year       int     `json:"year"`
	CO2TotalMt float64 `json:"co2_total_mt"`
}
