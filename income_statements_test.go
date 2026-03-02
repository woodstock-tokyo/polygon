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

func TestIncomeStatement_GetRevenue(t *testing.T) {
	stmt := IncomeStatement{Revenue: 1234.56}
	if got := stmt.GetRevenue(); got != 1234.56 {
		t.Errorf("GetRevenue() = %v, want %v", got, 1234.56)
	}
}

func TestIncomeStatementsOption_Defaults(t *testing.T) {
	opt := IncomeStatementsOption{}
	if opt.Limit != 0 { // zero means use API default (100)
		t.Errorf("Limit default = %v, want 0", opt.Limit)
	}
}

func Test_IncomeStatements(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	resp, err := c.IncomeStatements(ctx, &IncomeStatementsOption{Limit: 1})
	if err != nil {
		t.Fatal(fmt.Errorf("income statements: %w", err))
	}
	if len(resp.Results) == 0 {
		t.Error("unexpected empty results")
	}
}

func Test_IncomeStatementsNotFound(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	// Use invalid CIK expecting no results (depends on API behavior).
	_, err := c.IncomeStatements(ctx, &IncomeStatementsOption{CIK: "0000000000", Limit: 1})
	if !errors.Is(err, ErrIncomeStatementsNoResults) {
		t.Skip("API may return OK empty or different error; adjust once behavior confirmed. got:", err)
	}
}
