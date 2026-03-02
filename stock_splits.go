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
)

// StockSplits Get a list of historical stock splits
type StockSplits struct {
	Results []StockSplitsResult `json:"results"`
	Status  string              `json:"status"`
	NextURL string              `json:"next_url"`
}

// StockSplitsResult stock splits result item
type StockSplitsResult struct {
	ExecutionDate string  `json:"execution_date"`
	SplitFrom     float64 `json:"split_from"`
	SplitTo       float64 `json:"split_to"`
	Ticker        string  `json:"ticker"`
}

type StockSplitsOption struct {
	Ticker           string `url:"ticker,omitempty"`
	ExecutionDate    string `url:"execution_date,omitempty"`
	ExecutionDateGTE string `url:"execution_date.gte,omitempty"`
	Order            Order  `url:"order,omitempty"`
	Limit            uint   `url:"limit,omitempty"`
	Sort             string `url:"sort,omitempty"`
}

// StockSplits Get a list of historical stock splits
func (c Client) StockSplits(ctx context.Context, opt *StockSplitsOption) (StockSplits, error) {
	c = c.UseV3Endpoints()
	d := StockSplits{}

	if opt == nil {
		opt = new(StockSplitsOption)
	}

	endpoint, err := c.endpointWithOpts("/reference/splits", opt)
	if err != nil {
		return d, err
	}
	err = c.GetJSON(ctx, endpoint, &d)
	return d, err
}
