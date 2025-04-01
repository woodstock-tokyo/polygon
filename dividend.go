package polygon

import (
	"context"
	"fmt"
)

// Dividend Get a list of historical cash dividends, including the ticker symbol, declaration date, ex-dividend date, record date, pay date, frequency, and amount.
type Dividend struct {
	Results   []DividendResult `json:"results"`
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	NextURL   string           `json:"next_url"`
}

// DividendResult dividend result item
type DividendResult struct {
	CashAmount      float64 `json:"cash_amount"`
	Currency        string  `json:"currency"`
	DeclarationDate string  `json:"declaration_date"`
	DividendType    string  `json:"dividend_type"`
	ExDividendDate  string  `json:"ex_dividend_date"`
	Frequency       int     `json:"frequency"`
	PayDate         string  `json:"pay_date"`
	RecordDate      string  `json:"record_date"`
	Ticker          string  `json:"ticker"`
}

type DividendType string

const (
	CD DividendType = "CD"
	SC DividendType = "SC"
	LT DividendType = "LT"
	ST DividendType = "ST"
)

type DividendOption struct {
	Ticker string `url:"ticker"`
	// Query by the number of times per year the dividend is paid out. Possible values are 0 (one-time), 1 (annually), 2 (bi-annually), 4 (quarterly), and 12 (monthly).
	Frequency int `url:"frequency,omitempty"`
	// Query by the type of dividend. Dividends that have been paid and/or are expected to be paid on consistent schedules are denoted as CD.
	// Special Cash dividends that have been paid that are infrequent or unusual, and/or can not be expected to occur in the future are denoted as SC.
	DividendType      DividendType `url:"dividend_type,omitempty"`
	Order             Order        `url:"order,omitempty"`
	Limit             uint         `url:"limit,omitempty"`
	Sort              string       `url:"sort,omitempty"`
	ExDividendDateGTE string       `url:"ex_dividend_date.gte,omitempty"`
}

// Dividend Get a list of historical cash dividends, including the ticker symbol, declaration date, ex-dividend date, record date, pay date, frequency, and amount.
func (c Client) Dividend(ctx context.Context, ticker string, opt *DividendOption) (Dividend, error) {
	c = c.UseV3Endpoints()
	d := Dividend{}

	if opt == nil {
		opt = new(DividendOption)
	}

	opt.Ticker = ticker
	endpoint, err := c.endpointWithOpts("/reference/dividends", opt)
	if err != nil {
		return d, err
	}
	err = c.GetJSON(ctx, endpoint, &d)
	return d, err
}

// LastestDiviend retrieves the latest dividend for a given ticker
func (c Client) LastestDiviend(ctx context.Context, ticker string, opt *DividendOption) (DividendResult, error) {
	if opt == nil {
		opt = &DividendOption{}
	}

	opt.Limit = 1
	opt.Order = Descend
	// opt.Sort = "declaration_date" // to avoid Polygon bug , temporary remove sort key

	d, err := c.Dividend(ctx, ticker, opt)
	if err != nil {
		return DividendResult{}, err
	}

	if len(d.Results) == 0 {
		return DividendResult{}, fmt.Errorf("no dividend found for %s", ticker)
	}

	if d.Status != "OK" {
		return DividendResult{}, fmt.Errorf("status is not OK: %s", d.Status)
	}

	return d.Results[0], nil
}
