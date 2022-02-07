package polygon

import (
	"context"
	"testing"
)

func TestMarket(t *testing.T) {
	client := NewClient(token)

	_, err := client.MarketStatus(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
