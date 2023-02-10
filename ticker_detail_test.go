package polygon

import (
	"context"
	"testing"
)

func TestTickerDetail(t *testing.T) {
	client := NewClient(token, edgeID, edgeIPAddress)
	_, err := client.TickerDetail(context.Background(), "AAPL", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
