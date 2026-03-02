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
	"errors"
	"fmt"
	"testing"
)

func TestFinancialResult_GetRevenue(t *testing.T) {
	fr := FinancialResult{
		Financials: FinancialStatements{
			IncomeStatement: map[string]Metric{
				"revenues": {Value: 123.45},
			},
		},
	}
	if got := fr.GetRevenue(); got != 123.45 {
		t.Errorf("GetRevenue() = %v, want %v", got, 123.45)
	}

	fr = FinancialResult{
		Financials: FinancialStatements{
			IncomeStatement: map[string]Metric{},
		},
	}
	if got := fr.GetRevenue(); got != 0 {
		t.Errorf("GetRevenue() = %v, want %v", got, 0)
	}
}

func TestFinancialsOption_Defaults(t *testing.T) {
	opt := FinancialsOption{}
	if opt.Ticker != "" {
		t.Errorf("Ticker default = %v, want empty string", opt.Ticker)
	}
	if opt.Limit != 0 {
		t.Errorf("Limit default = %v, want 0", opt.Limit)
	}
}

// MockClient for testing Financials method
type MockClient struct{}

func (c MockClient) UseVXEndpoints() Client {
	return Client{}
}

func (c MockClient) GetJSON(ctx context.Context, endpoint string, v interface{}) error {
	f := v.(*Financials)
	*f = Financials{
		Results: []FinancialResult{{
			Financials: FinancialStatements{
				IncomeStatement: map[string]Metric{"revenues": {Value: 999.99}},
			},
		}},
	}
	return nil
}

func Test_Financials(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	f, err := c.Financials(ctx, "AAPL", &FinancialsOption{
		Limit: 1,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("get financials: %w", err))
	}
	if len(f.Results) == 0 {
		t.Error("unexpected financials results:", f)
	}
}

func Test_FinancialsErrorNotFound(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	_, err := c.Financials(ctx, "NOT_A_SYMBOL", &FinancialsOption{
		Limit: 1,
	})
	if !errors.Is(err, ErrFinancialsNoResults) {
		t.Fatal("unexpected error:", err)
	}
}
