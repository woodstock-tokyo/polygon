package polygon

import (
	"context"
	"fmt"
	"time"
)

// PrevClose Get the previous day's open, high, low, and close (OHLC) for the specified stock ticker.
type PrevClose struct {
	Ticker       string            `json:"ticker"`
	QueryCount   int               `json:"queryCount"`
	ResultsCount int               `json:"resultsCount"`
	Adjusted     bool              `json:"adjusted"`
	Results      []PrevCloseResult `json:"results"`
	Status       string            `json:"status"`
	RequestID    string            `json:"request_id"`
	Count        int               `json:"count"`
}

// PrevCloseResult prev close result item
type PrevCloseResult struct {
	Ticker                 string  `json:"T"`
	Open                   float64 `json:"o"`
	Close                  float64 `json:"c"`
	High                   float64 `json:"h"`
	Low                    float64 `json:"l"`
	TransactionNumber      int     `json:"n"`
	Volume                 float64 `json:"v"`
	VolumeWeightedAvgPrice float64 `json:"vw"`
	Timestamp              int64   `json:"t"`
}

// Valid check whether prev close is valid or not
func (p PrevClose) Valid() bool {
	return (p.Status == "OK" || p.Status == "DELAYED") && p.Count > 0
}

// Time check whether prev close is valid or not
func (pr PrevCloseResult) Time() time.Time {
	return time.UnixMilli(pr.Timestamp)
}

// PrevCloseOption prev close option
type PrevCloseOption struct {
	Adjusted bool `url:"adjusted,omitempty"`
}

// PrevClose Get the previous day's open, high, low, and close (OHLC) for the specified stock ticker.
func (c Client) PrevClose(ctx context.Context, ticker string, opt *PrevCloseOption) (PrevClose, error) {
	p := PrevClose{}

	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/aggs/ticker/%s/prev", ticker), opt)
	if err != nil {
		return p, err
	}
	err = c.GetJSONWithRetries(ctx, endpoint, &p)
	return p, err
}
