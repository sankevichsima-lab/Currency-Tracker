package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CoinPrice stores the price in USD.
// CoinPrice хранит цену монеты в долларах США.
type CoinPrice struct {
	USD float64 `json:"usd"`
}

// CoinGeckoResponse maps the API response for predefined coins.
// CoinGeckoResponse описывает структуру ответа API для конкретных монет.
type CoinGeckoResponse struct {
	Bitcoin  CoinPrice `json:"bitcoin"`
	Ethereum CoinPrice `json:"ethereum"`
}

// CoinGeckoClient handles HTTP communication with CoinGecko API.
// CoinGeckoClient отвечает за HTTP-запросы к API CoinGecko.
type CoinGeckoClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCoinGeckoClient creates a new API client with a request timeout.
// NewCoinGeckoClient создает новый клиент API с заданным таймаутом запросов.
func NewCoinGeckoClient(baseURL string, timeOut time.Duration) *CoinGeckoClient {
	client := &http.Client{
		Timeout: timeOut,
	}
	return &CoinGeckoClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

// GetPrice fetches the current USD price for any supported cryptocurrency.
// GetPrice получает текущую цену в USD для любой поддерживаемой криптовалюты.
func (c *CoinGeckoClient) GetPrice(coin string) (float64, error) {
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.baseURL, coin)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("problem with getting data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status: %d", resp.StatusCode)
	}

	// Decode into a map to handle dynamic coin tickers in JSON keys.
	// Декодируем в map, чтобы парсить произвольные имена монет из ключей JSON.
	var result map[string]CoinPrice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("problem with decoding: %w", err)
	}

	// Check if the requested coin is present in the response map.
	// Проверяем наличие запрашиваемой монеты в полученной карте.
	priceData, ok := result[coin]
	if !ok {
		return 0, fmt.Errorf("coin %s not found in response", coin)
	}

	return priceData.USD, nil
}
