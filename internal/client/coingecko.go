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
	Bitcoin CoinPrice `json:"bitcoin"`
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

func (c *CoinGeckoClient) GetBitcoinPrice() (float64, error) {
	url := fmt.Sprintf("%s/simple/price?ids=bitcoin&vs_currencies=usd", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("problem with getting data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status: %d", resp.StatusCode)
	}

	var result CoinGeckoResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("problem with decoding: %w", err)
	}
	return result.Bitcoin.USD, nil
}
