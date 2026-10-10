package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Price struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type Client struct {
	URL  string // full endpoint, e.g. https://api.binance.com/api/v3/ticker/price
	HTTP *http.Client
}

// NewClient returns a Client with a safe default HTTP client (never nil, always a timeout).
func NewClient(url string) *Client {
	return &Client{
		URL:  url,
		HTTP: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetPrice(ctx context.Context, symbol string) (*Price, error) {
	u := c.URL + "?symbol=" + symbol

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance: unexpected status %d for %s", resp.StatusCode, symbol)
	}

	var price Price
	err = json.NewDecoder(resp.Body).Decode(&price)
	if err != nil {
		return nil, err
	}
	return &price, nil
}
