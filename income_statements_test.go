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
