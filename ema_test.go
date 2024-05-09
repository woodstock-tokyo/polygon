package polygon

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func Test_EMA(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	sma, err := c.ExponentialMovingAverage(ctx, "AAPL", &EMAOption{
		Timespan: Day,
		Window:   50,
		Limit:    1,
	})

	if err != nil {
		t.Fatal(fmt.Errorf("get sma: %w", err))
	}

	if len(sma.Results.Values) == 0 {
		t.Error("unexpected sma:", sma)
	}
}

func Test_EMAErrorNotFound(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	_, err := c.ExponentialMovingAverage(ctx, "NOT_A_SYMBOL", &EMAOption{
		Timespan: Day,
		Window:   50,
		Limit:    1,
	})

	if !errors.Is(err, ErrEMANoResults) {
		t.Fatal("unexpected error:", err)
	}
}
