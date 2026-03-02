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

type Market struct {
	Status     MarketStatus `json:"market"`
	EarlyHours bool         `json:"earlyHours"`
	AfterHours bool         `json:"afterHours"`
	ServerTime string       `json:"serverTime"`
	Exchanges  struct {
		NYSE   MarketStatus `json:"nyse"`
		Nasdaq MarketStatus `json:"nasdaq"`
		OTC    MarketStatus `json:"otc"`
	} `json:"exchanges"`
	Currencies struct {
		FX     MarketStatus `json:"fx"`
		Crypto MarketStatus `json:"crypto"`
	} `json:"currencies"`
}

func (m Market) String(market ...string) string {
	_market := "stock"
	if len(market) == 1 {
		_market = market[0]
	}

	switch _market {
	case "stock":
		if m.Status == Open {
			return string(Open)
		}

		if m.Status == Closed {
			return string(Closed)
		}

		if m.Status == Overnight {
			return string(Overnight)
		}

		if m.EarlyHours {
			return string(EarlyHours)
		}

		if m.AfterHours {
			return string(AfterHours)
		}

		return ""

	case "crypto":
		return string(m.Currencies.Crypto)

	case "forex":
		return string(m.Currencies.FX)

	default:
		return ""
	}
}

type MarketOption struct{}

// MarketStatus Get the current trading status of the exchanges and overall financial markets.
func (c Client) MarketStatus(ctx context.Context) (Market, error) {
	c = c.UseV1Endpoints()

	m := Market{}
	endpoint, err := c.endpointWithOpts("/marketstatus/now", new(MarketOption))
	if err != nil {
		return m, err
	}
	err = c.GetJSON(ctx, endpoint, &m)
	return m, err
}
