package polygon

import (
	"context"
	"fmt"
)

// TickerDetail Get a single ticker supported by Polygon.io. This response will have detailed information about the ticker and the company behind it
type TickerDetail struct {
	Results   TickerDetailResult `json:"results"`
	Status    string             `json:"status"`
	RequestID string             `json:"request_id"`
}

type TickerType string

const (
	TickerTypeCommonStock TickerType = "CS"
)

// TickerDetailResult ticker detail result item
type TickerDetailResult struct {
	Ticker          string     `json:"ticker"`
	Name            string     `json:"name"`
	Market          string     `json:"market"`
	Locale          string     `json:"locale"`
	PrimaryExchange string     `json:"primary_exchange"`
	Type            TickerType `json:"type"`
	Active          bool       `json:"active"`
	CurrencyName    string     `json:"currency_name"`
	MarketCap       float64    `json:"market_cap"`
	PhoneNumber     string     `json:"phone_number"`
	Address         struct {
		Address1   string `json:"address1"`
		City       string `json:"city"`
		State      string `json:"state"`
		PostalCode string `json:"postal_code"`
	} `json:"address"`
	Description    string `json:"description"`
	Sector         string `json:"sic_description"`
	HomePageURL    string `json:"homepage_url"`
	TotalEmployees int    `json:"total_employees"`
	ListDate       string `json:"list_date"`
	Branding       struct {
		LogoURL string `json:"logo_url"`
		IconURL string `json:"icon_url"`
	} `json:"branding"`
	SIC string `json:"sic_code"`
}

type TickerDetailOption struct {
	Date string `url:"date,omitempty"`
}

func (c Client) TickerDetail(ctx context.Context, ticker string, opt *TickerDetailOption) (TickerDetail, error) {
	c = c.UseV3Endpoints()
	d := TickerDetail{}

	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/reference/tickers/%s", ticker), opt)
	if err != nil {
		return d, err
	}
	err = c.GetJSON(ctx, endpoint, &d)
	return d, err
}
