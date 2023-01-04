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
