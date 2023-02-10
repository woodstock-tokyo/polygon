package polygon

import (
	"context"
	"testing"
)

func TestDividend(t *testing.T) {
	client := NewClient(token, edgeID, edgeIPAddress)
	_, err := client.LastestDiviend(context.Background(), "AAPL", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
