package polygon

import (
	"context"
	"testing"
)

func TestMarket(t *testing.T) {
	client := NewClient(token, edgeID, edgeIPAddress)

	_, err := client.MarketStatus(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
