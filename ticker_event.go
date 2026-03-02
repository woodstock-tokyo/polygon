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
)

// TickerEvent Get a timeline of events for the entity associated with the given ticker, CUSIP, or Composite FIGI.
type TickerEvent struct {
	Results   TickerEventResult `json:"results"`
	Status    string            `json:"status"`
	RequestID string            `json:"request_id"`
}

// TickerEventResult ticker event result item
type TickerEventResult struct {
	Name   string            `json:"name"`
	FIGI   string            `json:"figi"`
	CIK    string            `json:"cik"`
	Events []TickerEventItem `json:"events"`
}

type TickerEventItem struct {
	Type         string                `json:"type"`
	Date         string                `json:"date"`
	TickerChange TickerChangeEventItem `json:"ticker_change"`
}

type TickerChangeEventItem struct {
	Ticker string `json:"ticker"`
}

type TickerEventOption struct {
	// A comma-separated list of the types of event to include. Currently ticker_change is the only supported event_type. Leave blank to return all supported event_types.
	Types string `url:"types,omitempty"`
}

// TickerEvent Get a timeline of events for the entity associated with the given ticker, CUSIP, or Composite FIGI.
func (c Client) TickerEvent(ctx context.Context, ticker string, opt *TickerEventOption) (TickerEvent, error) {
	c = c.UseVXEndpoints()
	e := TickerEvent{}

	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/reference/tickers/%s/events", ticker), opt)
	if err != nil {
		return e, err
	}
	err = c.GetJSON(ctx, endpoint, &e)
	return e, err
}
