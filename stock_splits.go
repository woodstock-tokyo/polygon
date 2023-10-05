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
	ExecutionDate string `json:"execution_date"`
	SplitFrom     int    `json:"split_from"`
	SplitTo       int    `json:"split_to"`
	Ticker        string `json:"ticker"`
}

type StockSplitsOption struct {
	Ticker        string `url:"ticker,omitempty"`
	ExecutionDate string `url:"execution_date,omitempty"`
	Order         Order  `url:"order,omitempty"`
	Limit         uint   `url:"limit,omitempty"`
	Sort          string `url:"sort,omitempty"`
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
