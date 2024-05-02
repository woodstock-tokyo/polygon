package polygon

import (
	"fmt"
	"strings"
)

// ForexEventTypeEnum event type enum
type ForexEventTypeEnum string

const (
	// ForexEventTypeOther others
	ForexEventTypeOther ForexEventTypeEnum = ""
	// ForexEventTypeCA minute aggregates
	ForexEventTypeCA ForexEventTypeEnum = "CA"
	// ForexEventTypeCAS second aggregates
	ForexEventTypeCAS ForexEventTypeEnum = "CAS"
)

type ForexAggregate struct {
	Event                 ForexEventTypeEnum // The event type.
	Pair                  string             // The current pair.
	Open                  float64            // Today's official opening price.
	TickVWAP              float64            // The tick's volume weighted average price.
	TickOpen              float64            // The opening tick price for this aggregate window.
	TickClose             float64            // The closing tick price for this aggregate window.
	TickHigh              float64            // The highest tick price for this aggregate window.
	TickLow               float64            // The lowest tick price for this aggregate window.
	StartTimestamp        int64              // The timestamp of the starting tick for this aggregate window in Unix Milliseconds.
	EndTimestamp          int64              // The timestamp of the ending tick for this aggregate window in Unix Milliseconds.
	Performance           float64            // performance from last market close
	PerformancePercentage float64            // performance percentage from last market close

}

func (c Client) SubscribeForexAggregates(client WebSocketClient, pairs []string, eventType ForexEventTypeEnum) (err error) {
	// connect
	client.Dial(fmt.Sprintf("%s/forex", c.websocketBaseURL), nil)
	// auth
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", c.token)))
	if err != nil {
		return
	}

	channel := resolveForexChannel(pairs, eventType)
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channel)))
	if err != nil {
		return
	}

	return
}

func resolveForexChannel(pairs []string, eventType ForexEventTypeEnum) string {
	var sb strings.Builder
	for i, pair := range pairs {
		sb.WriteString(fmt.Sprintf("%s.%s", eventType, pair))
		if i < len(pairs)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}
