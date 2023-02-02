package polygon

import (
	"context"
	"testing"
)

func TestPrevClose(t *testing.T) {
	client := NewClient(token, edgeID, edgeIPAddress)

	opt := &PrevCloseOption{
		Adjusted: true,
	}

	_, err := client.PrevClose(context.Background(), "AAPL", opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = client.PrevClose(context.Background(), "X:BTCUSD", opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
