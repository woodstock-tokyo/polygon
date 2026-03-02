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
	err = c.GetJSONWithRetries(ctx, endpoint, &p)
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
	err = c.GetJSONWithRetries(ctx, endpoint, &p)
	return p, err
}
