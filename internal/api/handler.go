package api

import (
	"encoding/json"
	"net/http"

	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
)

type Handler struct {
	db *storage.Storage
}

func NewHandler(db *storage.Storage) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ticker parameter is required",
		})
		return
	}

	history, err := h.db.GetPriceHistory(ticker)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to fetch history",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
