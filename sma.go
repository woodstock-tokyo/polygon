package polygon

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrSMAStatus    = errors.New("sma: unexpected status")
	ErrSMANoResults = errors.New("sma: no results")
)

// SMAOption options for fetching SMA
type SMAOption struct {
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

type SMAResponse struct {
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

// SimpleMovingAverage get stock SMA for a given ticker
func (c Client) SimpleMovingAverage(ctx context.Context, ticker string, opt *SMAOption) (resp SMAResponse, err error) {
	c = c.UseV1Endpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	endpoint, err := c.endpointWithOpts("/indicators/sma/"+ticker, opt)
	if err != nil {
		return
	}

	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		err = fmt.Errorf("get json: %w", err)
		return
	}

	if resp.Status != "OK" {
		err = fmt.Errorf("%v: %w", resp.Status, ErrSMAStatus)
		return
	}

	if len(resp.Results.Values) == 0 {
		err = ErrSMANoResults
		return
	}

	return resp, err
}
