package polygon

import (
	"context"
	"fmt"
	"testing"
)

func TestSummary(t *testing.T) {
	client := NewClient(token, edgeID, edgeIPAddress)

	assets := []Asset{
		{Ticker: "AAPL", AssetType: "stock"},
		{Ticker: "EUR/USD", AssetType: "forex"},
		{Ticker: "BTC", AssetType: "crypto"},
		{Ticker: "SPY250321C00380000", AssetType: "option"},
	}

	s, err := client.Summary(context.Background(), assets)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	fmt.Println(s)
}
