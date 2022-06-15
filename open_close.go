package polygon

import (
	"context"
	"fmt"
)

// OpenClose Get the open, close and afterhours prices of a stock symbol on a certain date.
type OpenClose struct {
	Status     string  `json:"status"`
	From       string  `json:"from"`
	Symbol     string  `json:"symbol"`
	Open       float64 `json:"open"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Volume     float64 `json:"volume"`
	AfterHours float64 `json:"afterHours"`
	PreMarket  float64 `json:"preMarket"`
}

// Valid check whether open close is valid or not
func (o OpenClose) Valid() bool {
	return o.Status == "OK" || o.Status == "DELAYED"
}

// OpenCloseOption prev close option
type OpenCloseOption struct {
	Adjusted bool `url:"adjusted,omitempty"`
}

// News retrieves the given number of news articles for the given stock symbol.
func (c Client) OpenClose(ctx context.Context, ticker string, date string, opt *OpenCloseOption) (OpenClose, error) {
	c = c.UseV1Endpoints()
	p := OpenClose{}
	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/open-close/%s/%s", ticker, date), opt)
	if err != nil {
		return p, err
	}
	err = c.GetJSON(ctx, endpoint, &p)
	return p, err
}
