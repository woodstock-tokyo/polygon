package polygon

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRSIStatus    = errors.New("rsi: unexpected status")
	ErrRSINoResults = errors.New("rsi: no results")
)

// RSIOption options for fetching RSI
type RSIOption struct {
	Timespan string `url:"timespan"`
	Adjusted bool   `json:"adjusted"`
	Window   uint   `json:"window"`
	Limit    uint   `url:"limit,omitempty"`
	Order    Order  `url:"order,omitempty"`
	Sort     string `url:"sort,omitempty"`
}

type RSIResponse struct {
	Results struct {
		Underlying struct {
			URL string `json:"url"`
		} `json:"underlying"`
		Values []struct {
			Timestamp int64   `json:"timestamp"`
			Value     float64 `json:"value"`
		} `json:"values"`
	} `json:"results"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	NextURL   string `json:"next_url"`
}

// RSIStockOneDay gives a one day RSI for a stock
func (c Client) RSIStockOneDay(ctx context.Context, ticker string) (float64, error) {
	c = c.UseV1Endpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	opt := &RSIOption{
		Timespan: "day",
		Adjusted: true,
		Window:   2,
		Limit:    1,
		Order:    Decend,
	}
	endpoint, err := c.endpointWithOpts("/indicators/rsi/"+ticker, opt)
	if err != nil {
		return 0, err
	}

	var resp RSIResponse
	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		return 0, fmt.Errorf("get json: %w", err)
	}
	if resp.Status != "OK" {
		return 0, fmt.Errorf("%v: %w", resp.Status, ErrRSIStatus)
	}
	if len(resp.Results.Values) == 0 {
		return 0, ErrRSINoResults
	}

	return resp.Results.Values[0].Value, err
}
