package transport

import (
	"encoding/json"
	"net/http"

	"github.com/Helltale/no-fines/internal/domain"
	"github.com/Helltale/no-fines/internal/service"
)

type ExchangeHandler struct {
	exchangeService *service.ExchangeService
}

func NewExchangeHandler(exchangeService *service.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{exchangeService: exchangeService}
}

func (h *ExchangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/exchange/rate":
		h.GetBestRate(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *ExchangeHandler) GetBestRate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BaseCurrency  string `json:"base_currency"`
		QuoteCurrency string `json:"quote_currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rate, err := h.exchangeService.GetBestRate(r.Context(), domain.CurrencyPair{
		BaseCurrency:  req.BaseCurrency,
		QuoteCurrency: req.QuoteCurrency,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"rate": rate})
}
