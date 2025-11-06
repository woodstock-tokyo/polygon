package polygon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDividend(t *testing.T) {
	client := NewClient(token)
	dividend, err := client.LastestDiviend(context.Background(), "AAPL", &DividendOption{Limit: 1})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("dividend: %+v", dividend)

	_, err = client.LastestDiviend(context.Background(), "COIN", &DividendOption{Limit: 1})
	assert.Error(t, err)
}
