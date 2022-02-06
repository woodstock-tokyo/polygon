package polygon

import (
	"context"
	"testing"
	"time"
)

func TestAggregation(t *testing.T) {
	client := NewClient(token)

	opt := &AggregationOption{
		Adjusted: true,
		Sort:     Ascend,
		Limit:    10,
	}

	_, err := client.Aggregation(context.Background(), "AAPL", 1, Day, time.Now().AddDate(0, 0, -5), time.Now(), opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLatestAggregation(t *testing.T) {
	client := NewClient(token)

	_, err := client.LatestAggregation(context.Background(), "X:BTCUSD")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
