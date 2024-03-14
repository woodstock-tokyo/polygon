package polygon

import (
	"context"
	"testing"
)

func TestSummary(t *testing.T) {
	client := NewClient(token)
	assets := []SummaryAsset{
		{Ticker: "AAPL", AssetType: "stock"},
		{Ticker: "EUR/USD", AssetType: "forex"},
		{Ticker: "BTC", AssetType: "crypto"},
		{Ticker: "SPY250321C00380000", AssetType: "option"},
	}

	_, err := client.Summary(context.Background(), assets)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
