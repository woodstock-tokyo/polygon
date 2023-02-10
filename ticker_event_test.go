package polygon

import (
	"context"
	"testing"
)

func TestTickerEvent(t *testing.T) {
	client := NewClient(token, edgeID, edgeIPAddress)
	_, err := client.TickerEvent(context.Background(), "AAPL", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
