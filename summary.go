package polygon

import (
	"context"
	"fmt"
	"strings"
)

// Summary Get everything needed to visualize the tick-by-tick movement of a list of tickers.
type Summary struct {
	RequestID string          `json:"request_id"`
	Results   []SummaryResult `json:"results"`
	Status    string          `json:"status"`
}

// SummaryResult summary result item
type SummaryResult struct {
	Branding struct {
		IconURL string `json:"icon_url"`
		LogoURL string `json:"logo_url"`
	} `json:"branding"`
	MarketStatus MarketStatus   `json:"market_status"`
	Name         string         `json:"name"`
	Price        float64        `json:"price"`
	Session      SummarySession `json:"session"`
	Ticker       string         `json:"ticker"`
	Type         string         `json:"type"`
}

// SummarySession summary session item
type SummarySession struct {
	Change                    float64 `json:"change"`
	ChangePercent             float64 `json:"change_percent"`
	Close                     float64 `json:"close"`
	EarlyTradingChange        float64 `json:"early_trading_change"`
	EarlyTradingChangePercent float64 `json:"early_trading_change_percent"`
	High                      float64 `json:"high"`
	LateTradingChange         float64 `json:"late_trading_change"`
	LateTradingChangePercent  float64 `json:"late_trading_change_percent"`
	Low                       float64 `json:"low"`
	Open                      float64 `json:"open"`
	PreviousClose             float64 `json:"previous_close"`
	Volume                    float64 `json:"volume"`
}

// SummaryOption summary option
type SummaryOption struct {
	TickerAnyOf string `url:"ticker.any_of,omitempty"`
}

// Asset
type Asset struct {
	Ticker    string
	AssetType string
}

// Resolve resolve ticker format for sumary end point
func (a Asset) resolveTicker() string {
	switch a.AssetType {
	case "stock":
		return a.Ticker
	case "option":
		return fmt.Sprintf("O:%s", strings.ToUpper(a.Ticker))
	case "forex":
		return fmt.Sprintf("C:%s", strings.ToUpper(strings.ReplaceAll(a.Ticker, "/", "")))
	case "crypto":
		return fmt.Sprintf("X:%sUSD", strings.ToUpper(a.Ticker))

	}

	return a.Ticker
}

// Summary Get everything needed to visualize the tick-by-tick movement of a list of tickers.
func (c Client) Summary(ctx context.Context, assets []Asset) (Summary, error) {

	opt := SummaryOption{}

	for _, asset := range assets {
		if opt.TickerAnyOf == "" {
			opt.TickerAnyOf = asset.resolveTicker()
			continue
		}

		opt.TickerAnyOf += fmt.Sprintf(",%s", asset.resolveTicker())
	}

	c = c.UseV1Endpoints()
	s := Summary{}
	endpoint, err := c.endpointWithOpts("/summaries", opt)
	if err != nil {
		return s, err
	}
	err = c.GetJSON(ctx, endpoint, &s)
	return s, err
}
