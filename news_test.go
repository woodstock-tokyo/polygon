package polygon

import (
	"context"
	"testing"
)

func TestNews(t *testing.T) {
	client := NewClient(token)

	opt := &NewsOption{
		Published_GreaterThan: "2020-02-01",
	}

	_, err := client.News(context.Background(), "AAPL", opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
