package polygon

import (
	"context"
	"fmt"
	"testing"
)

const testToken = "use.your.own.token"

func TestNews(t *testing.T) {
	client := NewClient(testToken)

	opt := &NewsOption{
		Published_GreaterThan: "2020-02-01",
	}

	news, err := client.News(context.Background(), "AAPL", opt)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	fmt.Println(news.Results)
}
