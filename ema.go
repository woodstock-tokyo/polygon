package polygon

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEMAStatus    = errors.New("ema: unexpected status")
	ErrEMANoResults = errors.New("ema: no results")
)

// EMAOption options for fetching EMA
type EMAOption struct {
	Timespan                    Timespan `url:"timespan"`
	Timestamp                   uint     `url:"timestamp,omitempty"`
	TimestampGreaterThan        uint     `url:"timestamp.gt,omitempty"`
	TimestampLessThan           uint     `url:"timestamp.lt,omitempty"`
	TimestampGreaterThanOrEqual uint     `url:"timestamp.gte,omitempty"`
	TimestampLessThanOrEqual    uint     `url:"timestamp.lte,omitempty"`
	Adjusted                    bool     `url:"adjusted,omitempty"`
	Window                      uint     `url:"window,omitempty"`
	Limit                       uint     `url:"limit,omitempty"`
	Order                       Order    `url:"order,omitempty"`
	Sort                        string   `url:"sort,omitempty"`
}

type EMAResponse struct {
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

// ExponentialMovingAverage get stock EMA for a given ticker
func (c Client) ExponentialMovingAverage(ctx context.Context, ticker string, opt *EMAOption) (resp EMAResponse, err error) {
	c = c.UseV1Endpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	endpoint, err := c.endpointWithOpts("/indicators/ema/"+ticker, opt)
	if err != nil {
		return
	}

	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		err = fmt.Errorf("get json: %w", err)
		return
	}

	if resp.Status != "OK" {
		err = fmt.Errorf("%v: %w", resp.Status, ErrEMAStatus)
		return
	}

	if len(resp.Results.Values) == 0 {
		err = ErrEMANoResults
		return
	}

	return resp, err
}
