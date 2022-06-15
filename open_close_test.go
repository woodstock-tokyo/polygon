package polygon

import (
	"context"
	"testing"
)

func TestOpenClose(t *testing.T) {
	client := NewClient(token)

	opt := &OpenCloseOption{
		Adjusted: true,
	}

	_, err := client.OpenClose(context.Background(), "AAPL", "2022-06-01", opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
