package polygon

import (
	"context"
	"testing"
)

func TestPreevClose(t *testing.T) {
	client := NewClient(token)

	opt := &PrevCloseOption{
		Adjusted: true,
	}

	_, err := client.PrevClose(context.Background(), "AAPL", opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
