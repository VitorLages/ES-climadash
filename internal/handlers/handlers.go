package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"climadash/internal/data"
)

// Handlers agrupa os endpoints HTTP.
type Handlers struct {
	repo *data.Repository
}

func New(repo *data.Repository) *Handlers {
	return &Handlers{repo: repo}
}

// Register registra todas as rotas no mux fornecido.
func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", h.health)
	mux.HandleFunc("/api/countries", h.countries)
	mux.HandleFunc("/api/countries/", h.countryDetail)
	mux.HandleFunc("/api/years", h.years)
	mux.HandleFunc("/api/summary", h.summary)
	mux.HandleFunc("/api/emissions", h.emissions)
	mux.HandleFunc("/api/top", h.top)
}

func (h *Handlers) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) countries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.repo.Countries())
}

func (h *Handlers) countryDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	country := strings.TrimPrefix(r.URL.Path, "/api/countries/")
	country = strings.TrimSpace(country)
	if country == "" {
		writeError(w, http.StatusBadRequest, "país é obrigatório")
		return
	}
	if dec, err := url.PathUnescape(country); err == nil {
		country = dec
	}
	from, to, err := parseYearRange(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	series, ok := h.repo.SeriesFor(country, from, to)
	if !ok {
		writeError(w, http.StatusNotFound, "país não encontrado")
		return
	}
	writeJSON(w, http.StatusOK, series)
}

func (h *Handlers) years(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.repo.Years())
}

func (h *Handlers) summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	writeJSON(w, http.StatusOK, h.repo.Summary())
}

func (h *Handlers) emissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	raw := strings.TrimSpace(r.URL.Query().Get("countries"))
	if raw == "" {
		writeError(w, http.StatusBadRequest, "parâmetro 'countries' é obrigatório")
		return
	}
	parts := strings.Split(raw, ",")
	countries := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		countries = append(countries, p)
	}
	if len(countries) > 5 {
		writeError(w, http.StatusBadRequest, "máximo de 5 países por consulta")
		return
	}
	from, to, err := parseYearRange(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, h.repo.SeriesForMany(countries, from, to))
}

func (h *Handlers) top(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	n := 5
	if raw := r.URL.Query().Get("n"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			n = v
		}
	}
	writeJSON(w, http.StatusOK, h.repo.Top(n))
}

func parseYearRange(r *http.Request) (from, to int, err error) {
	if raw := r.URL.Query().Get("from"); raw != "" {
		v, e := strconv.Atoi(raw)
		if e != nil {
			return 0, 0, errBadParam("from")
		}
		from = v
	}
	if raw := r.URL.Query().Get("to"); raw != "" {
		v, e := strconv.Atoi(raw)
		if e != nil {
			return 0, 0, errBadParam("to")
		}
		to = v
	}
	if from != 0 && to != 0 && from > to {
		return 0, 0, errBadRange()
	}
	return from, to, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type paramError struct{ msg string }

func (e *paramError) Error() string { return e.msg }

func errBadParam(name string) error { return &paramError{msg: "parâmetro inválido: " + name} }
func errBadRange() error            { return &paramError{msg: "intervalo inválido: from > to"} }
