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
