package polygon

import (
	"context"
	"testing"
)

func TestEPS(t *testing.T) {
	client := NewClient(token)

	_, err := client.EarningsPerShare(context.Background(), "AAPL")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
