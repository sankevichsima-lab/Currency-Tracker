package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetBitcoinPrice_Success tests successful price fetching using a mock HTTP server.
// TestGetBitcoinPrice_Success проверяет успешное получение цены через тестовый HTTP-сервер.
func TestGetBitcoinPrice_Success(t *testing.T) {
	// Start a local test server returning a mock JSON response.
	// Запускаем локальный тестовый сервер, возвращающий мок-ответ в формате JSON.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"bitcoin": {"usd": 50000.0}}`))
	}))
	defer server.Close()

	client := NewCoinGeckoClient(server.URL, 2*time.Second)
	price, err := client.GetPrice("bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price != 50000.0 {
		t.Errorf("expected price 50000.0, got %v", price)
	}
}

// TestGetBitcoinPrice_ServerError verifies client behavior when the external server returns an error.
// TestGetBitcoinPrice_ServerError проверяет поведение клиента, если внешний сервер возвращает ошибку.
func TestGetBitcoinPrice_ServerError(t *testing.T) {
	// Start a test server that simulates an internal server error.
	// Запускаем тестовый сервер, имитирующий внутреннюю ошибку 500.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewCoinGeckoClient(server.URL, 2*time.Second)
	_, err := client.GetPrice("bitcoin")
	if err == nil {
		t.Error("expected error from server, got nil")
	}
}
