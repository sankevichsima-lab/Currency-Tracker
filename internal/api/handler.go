package api

import (
	"encoding/json"
	"net/http"

	"github.com/sankevichsima-lab/Currency-Tracker/internal/auth"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/cache"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
)

// Handler holds dependencies for API endpoints.
// Handler хранит зависимости для обработчиков API.
type Handler struct {
	db        *storage.Storage
	cache     *cache.Cache
	jwtSecret string
}

// LoginRequest represents login credentials payload.
// LoginRequest представляет тело запроса на вход.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents token response payload.
// LoginResponse представляет тело ответа с выданным токеном.
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler issues a JWT token for valid credentials.
// LoginHandler выдает JWT-токен для пользователя.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(req.Username, h.jwtSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// NewHandler initializes a new Handler instance.
// NewHandler создает новый экземпляр Handler.
func NewHandler(db *storage.Storage, cache *cache.Cache, jwtSecret string) *Handler {
	return &Handler{
		db:        db,
		cache:     cache,
		jwtSecret: jwtSecret,
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
