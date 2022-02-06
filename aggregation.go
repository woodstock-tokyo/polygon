package polygon

import (
	"context"
	"fmt"
	"time"
)

// Aggregation Get aggregate bars for a ticker over a given date range in custom time window sizes
type Aggregation struct {
	Ticker       string `json:"ticker"`
	QueryCount   int    `json:"queryCount"`
	ResultsCount int    `json:"resultsCount"`
	Adjusted     bool   `json:"adjusted"`
	Results      []struct {
		Open                   float64 `json:"o"`
		Close                  float64 `json:"c"`
		High                   float64 `json:"h"`
		Low                    float64 `json:"l"`
		TransactionNumber      int     `json:"n"`
		Volume                 float64 `json:"v"`
		VolumeWeightedAvgPrice float64 `json:"vw"`
	} `json:"results"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Count     int    `json:"count"`
}

// Valid check whether aggregation is valid or not
func (a Aggregation) Valid() bool {
	return (a.Status == "OK" || a.Status == "DELAYED") && a.Count > 0
}

// AggregationOption aggregation option
type AggregationOption struct {
	Adjusted bool  `url:"adjusted,omitempty"`
	Sort     Order `url:"sort,omitempty"`
	Limit    int   `url:"limit,omitempty"`
}

// News retrieves the given number of news articles for the given stock symbol.
func (c Client) Aggregation(ctx context.Context, ticker string, multiplier int, timespan Timespan, from, to time.Time, opt *AggregationOption) (Aggregation, error) {
	a := Aggregation{}
	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/aggs/ticker/%s/range/%d/%s/%s/%s", ticker, multiplier, timespan, ttoa(from), ttoa(to)), opt)
	if err != nil {
		return a, err
	}
	err = c.GetJSON(ctx, endpoint, &a)
	return a, err
}

// LatestAggregation get lastest aggregation
// however it will return empty result for stocks out of market time
// for cryptos, it always returns the valid aggregation data
func (c Client) LatestAggregation(ctx context.Context, ticker string) (Aggregation, error) {
	opt := AggregationOption{
		Adjusted: true,
		Sort:     Decend,
		Limit:    1,
	}

	est, _ := time.LoadLocation("America/New_York")
	targetDate := time.Now().In(est)
	return c.Aggregation(ctx, ticker, 1, Minute, targetDate, targetDate, &opt)
}
