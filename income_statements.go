package polygon

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type FinancialsOptionTimeframe string

const (
	FinancialsOptionTimeframeAnnual               FinancialsOptionTimeframe = "annual"
	FinancialsOptionTimeframeQuarterly            FinancialsOptionTimeframe = "quarterly"
	FinancialsOptionTimeframeTrailingTwelveMonths FinancialsOptionTimeframe = "trailing_twelve_months"
)

// IncomeStatementsResponse top-level response for the new income statements endpoint.
// This endpoint replaces the deprecated Financials endpoint.
type IncomeStatementsResponse struct {
	Results   []IncomeStatement `json:"results"`
	Status    string            `json:"status"`
	RequestID string            `json:"request_id"`
	NextURL   string            `json:"next_url"`
}

// IncomeStatement represents a single income statement period (annual, quarterly, or trailing twelve months)
// Only a subset of fields are modeled explicitly; additional numeric fields can be added as needed.
// The Polygon docs list these numeric metrics as numbers; we use float64. Optional fields default to zero if missing.
type IncomeStatement struct {
	BasicEarningsPerShare                       float64  `json:"basic_earnings_per_share"`
	BasicSharesOutstanding                      float64  `json:"basic_shares_outstanding"`
	CIK                                         string   `json:"cik"`
	ConsolidatedNetIncomeLoss                   float64  `json:"consolidated_net_income_loss"`
	CostOfRevenue                               float64  `json:"cost_of_revenue"`
	DepreciationDepletionAmortization           float64  `json:"depreciation_depletion_amortization"`
	DilutedEarningsPerShare                     float64  `json:"diluted_earnings_per_share"`
	DilutedSharesOutstanding                    float64  `json:"diluted_shares_outstanding"`
	DiscontinuedOperations                      float64  `json:"discontinued_operations"`
	EBITDA                                      float64  `json:"ebitda"`
	EquityInAffiliates                          float64  `json:"equity_in_affiliates"`
	ExtraordinaryItems                          float64  `json:"extraordinary_items"`
	FilingDate                                  string   `json:"filing_date"`
	FiscalQuarter                               int      `json:"fiscal_quarter"`
	FiscalYear                                  int      `json:"fiscal_year"`
	GrossProfit                                 float64  `json:"gross_profit"`
	IncomeBeforeIncomeTaxes                     float64  `json:"income_before_income_taxes"`
	IncomeTaxes                                 float64  `json:"income_taxes"`
	InterestExpense                             float64  `json:"interest_expense"`
	InterestIncome                              float64  `json:"interest_income"`
	NetIncomeLossAttributableCommonShareholders float64  `json:"net_income_loss_attributable_common_shareholders"`
	NoncontrollingInterest                      float64  `json:"noncontrolling_interest"`
	OperatingIncome                             float64  `json:"operating_income"`
	OtherIncomeExpense                          float64  `json:"other_income_expense"`
	OtherOperatingExpenses                      float64  `json:"other_operating_expenses"`
	PeriodEnd                                   string   `json:"period_end"`
	PreferredStockDividendsDeclared             float64  `json:"preferred_stock_dividends_declared"`
	ResearchDevelopment                         float64  `json:"research_development"`
	Revenue                                     float64  `json:"revenue"`
	SellingGeneralAdministrative                float64  `json:"selling_general_administrative"`
	Tickers                                     []string `json:"tickers"`
	Timeframe                                   string   `json:"timeframe"` // quarterly | annual | trailing_twelve_months
	TotalOperatingExpenses                      float64  `json:"total_operating_expenses"`
	TotalOtherIncomeExpense                     float64  `json:"total_other_income_expense"`
}

// IncomeStatementsOption holds query params for the income statements endpoint.
type IncomeStatementsOption struct {
	CIK           string                    `url:"cik,omitempty"`
	Tickers       string                    `url:"tickers,omitempty"`         // matches API expecting a value contained in array
	PeriodEnd     string                    `url:"period_end,omitempty"`      // YYYY-MM-DD
	FilingDate    string                    `url:"filing_date,omitempty"`     // YYYY-MM-DD
	FilingDateGTE string                    `url:"filing_date.gte,omitempty"` // Query by the date the financial statement was filed (greater than or equal to) in YYYY-MM-DD format.
	FilingDateLTE string                    `url:"filing_date.lte,omitempty"` // Query by the date the financial statement was filed (less than or equal to) in YYYY-MM-DD format.
	FiscalYear    string                    `url:"fiscal_year,omitempty"`
	FiscalQuarter string                    `url:"fiscal_quarter,omitempty"`
	Timeframe     FinancialsOptionTimeframe `url:"timeframe,omitempty"` // quarterly, annual, trailing_twelve_months
	Limit         uint                      `url:"limit,omitempty"`     // default 100, max 50000
	Sort          string                    `url:"sort,omitempty"`      // e.g. period_end.desc
}

var ErrIncomeStatementsNoResults = errors.New("no income statements results")

// IncomeStatements retrieves income statements data. This replaces the deprecated Financials endpoint.
func (c Client) IncomeStatements(ctx context.Context, opt *IncomeStatementsOption) (resp IncomeStatementsResponse, err error) {
	c = c.UseFinancialsV1Endpoints()
	if opt == nil {
		opt = new(IncomeStatementsOption)
	}
	endpoint, err := c.endpointWithOpts("/income-statements", opt)
	if err != nil {
		return
	}
	if err = c.GetJSON(ctx, endpoint, &resp); err != nil {
		err = fmt.Errorf("get json: %w", err)
		return
	}
	if strings.ToUpper(resp.Status) != "OK" {
		err = fmt.Errorf("%v: %w", resp.Status, ErrIncomeStatementsNoResults)
		return
	}
	if len(resp.Results) == 0 {
		err = ErrIncomeStatementsNoResults
		return
	}
	return
}

// GetRevenue helper returns the revenue for a statement period.
func (is IncomeStatement) GetRevenue() float64 { return is.Revenue }

// GetEarningsPerShare helper returns the diluted earnings per share for a statement period.
func (is IncomeStatement) GetEarningsPerShare() float64 {
	return is.DilutedEarningsPerShare
}
