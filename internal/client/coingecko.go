package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CoinPrice struct {
	USD float64 `json:"usd"`
}

type CoinGeckoResponse struct {
	Bitcoin  CoinPrice `json:"bitcoin"`
	Ethereum CoinPrice `json:"ethereum"`
}

type CoinGeckoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCoinGeckoClient(baseURL string, timeOut time.Duration) *CoinGeckoClient {
	client := &http.Client{
		Timeout: timeOut,
	}
	return &CoinGeckoClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

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

	var result map[string]CoinPrice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("problem with decoding: %w", err)
	}

	priceData, ok := result[coin]
	if !ok {
		return 0, fmt.Errorf("coin %s not found in response", coin)
	}

	return priceData.USD, nil
}
