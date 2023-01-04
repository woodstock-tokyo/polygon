package polygon

import (
	"context"
	"fmt"
	"time"
)

// Aggregation Get aggregate bars for a ticker over a given date range in custom time window sizes
type Aggregation struct {
	Ticker       string              `json:"ticker"`
	QueryCount   int                 `json:"queryCount"`
	ResultsCount int                 `json:"resultsCount"`
	Adjusted     bool                `json:"adjusted"`
	Results      []AggregationResult `json:"results"`
	Status       string              `json:"status"`
	RequestID    string              `json:"request_id"`
	Count        int                 `json:"count"`
}

// AggregationResult aggregation result item
type AggregationResult struct {
	Open                   float64 `json:"o"`
	Close                  float64 `json:"c"`
	High                   float64 `json:"h"`
	Low                    float64 `json:"l"`
	TransactionNumber      int     `json:"n"`
	Volume                 float64 `json:"v"`
	VolumeWeightedAvgPrice float64 `json:"vw"`
	Timestamp              int64   `json:"t"`
}

// AggregationResultSortFunc default sort function
func AggregationResultSortFunc() func(AggregationResult, AggregationResult) bool {
	return func(r1, r2 AggregationResult) bool {
		return r1.Timestamp < r2.Timestamp
	}
}

// Valid check whether aggregation is valid or not
func (a Aggregation) Valid() bool {
	return (a.Status == "OK" || a.Status == "DELAYED") && a.Count > 0
}

// Valid check whether aggregation is valid or not
func (ar AggregationResult) Time() time.Time {
	return time.UnixMilli(ar.Timestamp)
}

// AggregationOption aggregation option
type AggregationOption struct {
	Adjusted bool  `url:"adjusted,omitempty"`
	Sort     Order `url:"sort,omitempty"`
	Limit    int   `url:"limit,omitempty"`
}

// Aggregation Get aggregate bars for a ticker over a given date range in custom time window sizes
func (c Client) Aggregation(ctx context.Context, ticker string, multiplier int, timespan Timespan, from, to time.Time, opt *AggregationOption) (Aggregation, error) {
	a := Aggregation{}
	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/aggs/ticker/%s/range/%d/%s/%s/%s", ticker, multiplier, timespan, ttoa(from), ttoa(to)), opt)
	if err != nil {
		return a, err
	}
	err = c.GetJSON(ctx, endpoint, &a)
	return a, err
}
