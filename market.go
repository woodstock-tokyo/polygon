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
