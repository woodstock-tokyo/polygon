package polygon

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEPSStatus    = errors.New("eps: unexpected status")
	ErrEPSNoResults = errors.New("eps: no results")
	ErrPERStatus    = errors.New("per: unexpected status")
	ErrPERNoResults = errors.New("per: no results")
)

type EPSResponse struct {
	Results []struct {
		StartDate   string   `json:"start_date"`
		EndDate     string   `json:"end_date"`
		Timeframe   string   `json:"timeframe"`
		FiscalYear  string   `json:"fiscal_year"`
		Tickers     []string `json:"tickers"`
		CompanyName string   `json:"company_name"`
		Financials  struct {
			IncomeStatement struct {
				DilutedEarningsPerShare struct {
					Value float64 `json:"value"`
					Unit  string  `json:"unit"`
					Label string  `json:"label"`
					Order int     `json:"order"`
				} `json:"diluted_earnings_per_share"`
			} `json:"income_statement"`
		} `json:"financials"`
	} `json:"results"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
}

// EarningsPerShare get stock EPS from Polygon
func (c Client) EarningsPerShareFromPolygon(ctx context.Context, ticker string) (resp EPSResponse, err error) {
	c = c.UseVXEndpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))

	endpoint, err := c.endpointWithOpts(fmt.Sprintf("/reference/financials?ticker=%s&timeframe=ttm&limit=1", ticker), nil)
	if err != nil {
		return
	}

	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		err = fmt.Errorf("get json: %w", err)
		return
	}

	if resp.Status != "OK" {
		err = fmt.Errorf("%v: %w", resp.Status, ErrRSIStatus)
		return
	}

	if len(resp.Results) == 0 {
		err = ErrEPSNoResults
		return
	}

	return resp, err
}

// EarningsPerShare get stock EPS for a given ticker
func (c Client) EarningsPerShare(ctx context.Context, ticker string) (float64, error) {
	c = c.UseVXEndpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))

	resp, err := c.EarningsPerShareFromPolygon(ctx, ticker)
	if err != nil {
		return 0, err
	}

	if len(resp.Results) == 0 {
		return 0, ErrRSINoResults
	}

	return resp.Results[0].Financials.IncomeStatement.DilutedEarningsPerShare.Value, err
}
