package polygon

import (
	"fmt"
	"strings"
)

// CryptoEventTypeEnum event type enum
type CryptoEventTypeEnum string

const (
	// CryptoEventTypeOther others
	CryptoEventTypeOther CryptoEventTypeEnum = ""
	// CryptoEventTypeXA minute aggregates
	CryptoEventTypeXA CryptoEventTypeEnum = "XA"
	// CryptoEventTypeXAS second aggregates
	CryptoEventTypeXAS CryptoEventTypeEnum = "XAS"
)

type CryptoAggregate struct {
	Event                 CryptoEventTypeEnum // The event type.
	Pair                  string              // The crypto pair.
	TickOpen              float64             // The opening tick price for this aggregate window.
	TickClose             float64             // The closing tick price for this aggregate window.
	TickHigh              float64             // The highest tick price for this aggregate window.
	TickLow               float64             // The lowest tick price for this aggregate window.
	TickVWAP              float64             // The volume of trades during this aggregate window.
	VWAP                  float64             // Today's volume weighted average price.
	AverageTradeSize      float64             // The average trade size for this aggregate window.
	StartTimestamp        int64               // The timestamp of the starting tick for this aggregate window in Unix Milliseconds.
	EndTimestamp          int64               // The timestamp of the ending tick for this aggregate window in Unix Milliseconds.
	Performance           float64             // performance from last market close
	PerformancePercentage float64             // performance percentage from last market close
}

func (c Client) SubscribeCryptoAggregates(client WebSocketClient, pairs []string, eventType CryptoEventTypeEnum) (err error) {
	// connect
	client.Dial(fmt.Sprintf("%s/crypto", c.websocketBaseURL), nil)
	// auth
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", c.token)))
	if err != nil {
		return
	}

	channel := resolveCryptoChannel(pairs, eventType)
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channel)))
	if err != nil {
		return
	}

	return
}

func resolveCryptoChannel(pairs []string, eventType CryptoEventTypeEnum) string {
	var sb strings.Builder
	for i, pair := range pairs {
		sb.WriteString(fmt.Sprintf("%s.%s", eventType, pair))
		if i < len(pairs)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}
