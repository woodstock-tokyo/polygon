// Copyright (c) 2026 Woodstock K.K.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
func AggregationResultSortFunc() func(AggregationResult, AggregationResult) int {
	return func(r1, r2 AggregationResult) int {
		if r1.Timestamp < r2.Timestamp {
			return -1
		}

		if r1.Timestamp > r2.Timestamp {
			return 1
		}

		return 0
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
