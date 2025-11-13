package polygon

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Financials top-level response
type Financials struct {
	Results   []FinancialResult `json:"results"`
	Status    string            `json:"status"`
	RequestID string            `json:"request_id"`
	NextURL   string            `json:"next_url"`
}

// FinancialResult result (one filing/period)
type FinancialResult struct {
	StartDate    string              `json:"start_date"`    // e.g., "2009-06-28"
	EndDate      string              `json:"end_date"`      // e.g., "2009-09-26"
	Timeframe    string              `json:"timeframe"`     // "quarterly" | "annual"
	FiscalPeriod string              `json:"fiscal_period"` // e.g., "Q4"
	FiscalYear   string              `json:"fiscal_year"`   // e.g., "2009"
	CIK          string              `json:"cik"`
	SIC          string              `json:"sic"`
	Tickers      []string            `json:"tickers"`
	CompanyName  string              `json:"company_name"`
	Financials   FinancialStatements `json:"financials"`
}

// FinancialStatements groups of statements. Each statement is a map of metric key -> Metric.
type FinancialStatements struct {
	IncomeStatement     map[string]Metric `json:"income_statement"`
	BalanceSheet        map[string]Metric `json:"balance_sheet"`
	ComprehensiveIncome map[string]Metric `json:"comprehensive_income"`
	CashFlowStatement   map[string]Metric `json:"cash_flow_statement"`
}

// Metric a numeric metric entry (used across all statements)
type Metric struct {
	Value float64 `json:"value"` // Polygon numbers fit well in float64
	Unit  string  `json:"unit"`
	Label string  `json:"label"`
	Order int     `json:"order"`
}

// GetRevenue returns the "revenues" metric value if available.
func (fr *FinancialResult) GetRevenue() float64 {
	if fr.Financials.IncomeStatement == nil {
		return 0
	}
	if rev, ok := fr.Financials.IncomeStatement["revenues"]; ok {
		return rev.Value
	}
	return 0
}

// GetConsolidatedNetIncomeLoss returns the "net income / loss" value if available
func (fr *FinancialResult) GetConsolidatedNetIncomeLoss() float64 {
	if fr.Financials.IncomeStatement == nil {
		return 0
	}
	if rev, ok := fr.Financials.IncomeStatement["consolidated_net_income_loss"]; ok {
		return rev.Value
	}
	return 0
}

type FinancialsOption struct {
	Ticker             string                    `url:"ticker"`                          // Query by company ticker.
	CIK                string                    `url:"cik,omitempty"`                   // Query by central index key (CIK) Number
	CompanyName        string                    `url:"company_name,omitempty"`          // Query by company name.
	SIC                string                    `url:"sic,omitempty"`                   // Query by standard industrial classification (SIC)
	FilingDate         string                    `url:"filing_date,omitempty"`           // Query by the date the financial statement was filed (YYYY-MM-DD).
	FilingDateGTE      string                    `url:"filing_date.gte,omitempty"`       // Query by the date the financial statement was filed (greater than or equal to) in YYYY-MM-DD format.
	FilingDateLTE      string                    `url:"filing_date.lte,omitempty"`       // Query by the date the financial statement was filed (less than or equal to) in YYYY-MM-DD format.
	PeriodOfReportDate string                    `url:"period_of_report_date,omitempty"` // The period of report for the filing with financials data in YYYY-MM-DD format.
	Timeframe          FinancialsOptionTimeframe `url:"timeframe,omitempty"`             // Query by timeframe. Annual financials originate from 10-K filings, and quarterly financials originate from 10-Q filings
	IncludeSources     bool                      `url:"include_sources,omitempty"`       // Whether or not to include the `xpath` and `formula` attributes for each financial data point
	Order              Order                     `url:"order,omitempty"`                 // Order results by `filing_date` or `period_of_report_date`. Default is ascending (`asc`). Use `desc` for descending.
	Limit              uint                      `url:"limit,omitempty"`                 // Limit the number of results returned. Default is 10, max is 1000.
	Sort               string                    `url:"sort,omitempty"`                  // Sort field used for ordering
}

var ErrFinancialsNoResults = errors.New("no financials results")

// Financials Retrieve historical financial data for a specified stock ticker
//
// Deprecated: This API is deprecated and will be removed in a future version.
func (c Client) Financials(ctx context.Context, ticker string, opt *FinancialsOption) (resp Financials, err error) {
	c = c.UseVXEndpoints()
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	if opt == nil {
		opt = new(FinancialsOption)
	}
	opt.Ticker = ticker
	endpoint, err := c.endpointWithOpts("/reference/financials", opt)
	if err != nil {
		return
	}
	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		err = fmt.Errorf("get json: %w", err)
		return
	}
	if resp.Status != "OK" {
		err = fmt.Errorf("%v: %w", resp.Status, ErrFinancialsNoResults)
		return
	}
	if len(resp.Results) == 0 {
		err = ErrFinancialsNoResults
		return
	}
	return resp, err
}
