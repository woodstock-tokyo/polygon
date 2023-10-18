package polygon

import (
	"context"
	"testing"
)

func TestStockSplits(t *testing.T) {
	client := NewClient(token, WithEdge(edgeID, edgeIPAddress))
	_, err := client.StockSplits(context.Background(), nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
