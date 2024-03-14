package polygon

import (
	"context"
	"testing"
)

func TestStockSplits(t *testing.T) {
	client := NewClient(token)
	_, err := client.StockSplits(context.Background(), nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
