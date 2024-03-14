package polygon

import (
	"context"
	"fmt"
)

// OpenClose Get the open, close and afterhours prices of a stock symbol on a certain date.
type OpenClose struct {
	Status     string  `json:"status,omitempty"`
	From       string  `json:"from,omitempty"`
	Symbol     string  `json:"symbol"`
	Open       float64 `json:"open"`
	Close      float64 `json:"close"`
	High       float64 `json:"high,omitempty"`
	Low        float64 `json:"low,omitempty"`
	Volume     float64 `json:"volume,omitempty"`
	AfterHours float64 `json:"afterHours,omitempty"`
	PreMarket  float64 `json:"preMarket,omitempty"`
}

// OpenCloseOption prev close option
type OpenCloseOption struct {
	Adjusted bool `url:"adjusted,omitempty"`
}

// StockOpenClose Get the open, close and afterhours prices of a stock symbol on a certain date.
func (c Client) StockOpenClose(ctx context.Context, ticker string, date string, opt *OpenCloseOption) (OpenClose, error) {
	c = c.UseV1Endpoints()
	p := OpenClose{}

	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/open-close/%s/%s", ticker, date), opt)
	if err != nil {
		return p, err
	}
	err = c.GetJSON(ctx, endpoint, &p)
	return p, err
}

// CryptoOpenClose Get the open, close prices of a crypto pair on a certain date.
func (c Client) CryptoOpenClose(ctx context.Context, from, to string, date string, opt *OpenCloseOption) (OpenClose, error) {
	c = c.UseV1Endpoints()
	p := OpenClose{}
	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/open-close/crypto/%s/%s/%s", from, to, date), opt)
	if err != nil {
		return p, err
	}
	err = c.GetJSON(ctx, endpoint, &p)
	return p, err
}
