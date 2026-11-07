package api

import (
	"encoding/json"
	"net/http"

	"github.com/sankevichsima-lab/Currency-Tracker/internal/cache"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
)

// Handler holds dependencies for API endpoints.
// Handler хранит зависимости для обработчиков API.
type Handler struct {
	db    *storage.Storage
	cache *cache.Cache
}

// NewHandler initializes a new Handler instance.
// NewHandler создает новый экземпляр Handler.
func NewHandler(db *storage.Storage, cache *cache.Cache) *Handler {
	return &Handler{
		db:    db,
		cache: cache,
	}
}

// GetPriceHandler returns the latest cached price for a ticker.
// GetPriceHandler возвращает последнюю актуальную цену из кеша по тикеру.
func (h *Handler) GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ticker parameter is required",
		})
		return
	}

	// Read latest price directly from RAM cache.
	// Читаем последнюю цену напрямую из оперативной памяти.
	item, ok := h.cache.Get(ticker)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "price not found in cache",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// GetHistoryHandler returns price history from the database for a ticker.
// GetHistoryHandler возвращает историю цен монеты из базы данных.
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

	// Fetch historical records from SQLite.
	// Запрашиваем историю котировок из SQLite.
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
