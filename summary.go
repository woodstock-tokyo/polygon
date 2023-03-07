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
	Options      SummaryOptions `json:"options"`
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

// SummaryOptions summary options item
type SummaryOptions struct {
	ContractType      string `json:"contract_type"`
	ExerciseStyle     string `json:"exercise_style"`
	ExpirationDate    string `json:"expiration_date"`
	SharesPerContract uint   `json:"shares_per_contract"`
	StrikePrice       uint   `json:"strike_price"`
	UnderlyingTicker  string `json:"underlying_ticker"`
}

// SummaryOption summary option
type SummaryOption struct {
	TickerAnyOf string `url:"ticker.any_of,omitempty"`
}

// Asset
type SummaryAsset struct {
	Ticker    string
	AssetType string
}

// resolveTicker resolve ticker format for sumary end point
func (sa SummaryAsset) resolveTicker() string {
	switch sa.AssetType {
	case "stock":
		return sa.Ticker
	case "option":
		return fmt.Sprintf("O:%s", strings.ToUpper(sa.Ticker))
	case "forex":
		return fmt.Sprintf("C:%s", strings.ToUpper(strings.ReplaceAll(sa.Ticker, "/", "")))
	case "crypto":
		return fmt.Sprintf("X:%sUSD", strings.ToUpper(sa.Ticker))
	}

	return sa.Ticker
}

// Summary Get everything needed to visualize the tick-by-tick movement of a list of tickers.
func (c Client) Summary(ctx context.Context, assets []SummaryAsset) (Summary, error) {
	s := Summary{}
	c = c.UseV1Endpoints()

	if !c.launchPad {
		return s, fmt.Errorf("only launchpad supports this endpoint")
	}

	opt := SummaryOption{}
	for _, asset := range assets {
		if opt.TickerAnyOf == "" {
			opt.TickerAnyOf = asset.resolveTicker()
			continue
		}

		opt.TickerAnyOf += fmt.Sprintf(",%s", asset.resolveTicker())
	}

	endpoint, err := c.endpointWithOpts("/summaries", opt)
	if err != nil {
		return s, err
	}

	err = c.GetJSON(ctx, endpoint, &s)
	return s, err
}
