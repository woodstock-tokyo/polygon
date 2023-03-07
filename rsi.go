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

// LatestRelativeStrengthIndex get latest stock RSI by day for a given ticker
func (c Client) LatestRelativeStrengthIndex(ctx context.Context, ticker string) (float64, error) {
	if c.launchPad {
		return 0.0, fmt.Errorf("launchpad does not support this endpoint")
	}

	c = c.UseV1Endpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	opt := &RSIOption{
		Timespan: "day",
		Adjusted: true,
		Window:   2,
		Limit:    1,
		Order:    Descend,
	}

	resp, err := c.RelativeStrengthIndex(ctx, ticker, opt)
	if err != nil {
		return 0, err
	}

	if len(resp.Results.Values) == 0 {
		return 0, ErrRSINoResults
	}

	return resp.Results.Values[0].Value, err
}

// RelativeStrengthIndex get stock RSI for a given ticker
func (c Client) RelativeStrengthIndex(ctx context.Context, ticker string, opt *RSIOption) (resp RSIResponse, err error) {
	if c.launchPad {
		return RSIResponse{}, fmt.Errorf("launchpad does not support this endpoint")
	}

	c = c.UseV1Endpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	endpoint, err := c.endpointWithOpts("/indicators/rsi/"+ticker, opt)
	if err != nil {
		return
	}

	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		err = fmt.Errorf("get json: %w", err)
		return
	}

	if resp.Status != "OK" {
		err = fmt.Errorf("%v: %w", resp.Status, ErrRSIStatus)
		return
	}

	if len(resp.Results.Values) == 0 {
		err = ErrRSINoResults
		return
	}

	return resp, err
}
