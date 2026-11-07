package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetBitcoinPrice_Success(t *testing.T) {
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

func TestGetBitcoinPrice_ServerError(t *testing.T) {
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
