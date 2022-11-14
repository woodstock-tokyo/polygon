package polygon

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func Test_RSI(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	rsi, err := c.RSIStockOneDay(ctx, "AAPL")
	if err != nil {
		t.Fatal(fmt.Errorf("get rsi: %w", err))
	}

	if rsi == 0 {
		t.Error("unexpected rsi:", rsi)
	}
}

func Test_RSIErrorNotFound(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	_, err := c.RSIStockOneDay(ctx, "NOT_A_SYMBOL")

	if !errors.Is(err, ErrRSINoResults) {
		t.Fatal("unexpected error:", err)
	}
}
