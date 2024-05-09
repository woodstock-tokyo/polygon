package polygon

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func Test_SMA(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	sma, err := c.SimpleMovingAverage(ctx, "AAPL", &SMAOption{
		Timespan: "day",
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

func Test_SMAErrorNotFound(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	_, err := c.SimpleMovingAverage(ctx, "NOT_A_SYMBOL", &SMAOption{
		Timespan: "day",
		Window:   50,
		Limit:    1,
	})

	if !errors.Is(err, ErrSMANoResults) {
		t.Fatal("unexpected error:", err)
	}
}
