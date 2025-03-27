package polygon

import (
	"fmt"
	"net/http"
	"strings"
)

//////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// websocket interface ////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

type MessageType uint

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

type WebSocketClient interface {
	Dial(urlStr string, reqHeader http.Header)
	WriteMessage(messageType int, data []byte) error
}

// ////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////// polygon websocket //////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////
// StockEventTypeEnum event type enum
type StockEventTypeEnum string

const (
	// StockEventTypeOther others
	StockEventTypeOther StockEventTypeEnum = ""
	// StockEventTypeAM minute aggregates
	StockEventTypeAM StockEventTypeEnum = "AM"
	// StockEventTypeA second aggregates
	StockEventTypeA StockEventTypeEnum = "A"
)

type StockAggregate struct {
	Event             StockEventTypeEnum `json:"ev"`
	Symbol            string             `json:"sym"`
	TickVolume        float64            `json:"v"`
	AccumulatedVolume int64              `json:"av"`
	Open              float64            `json:"op"`
	TickVWAP          float64            `json:"vw"`
	TickOpen          float64            `json:"o"`
	TickClose         float64            `json:"c"`
	TickHigh          float64            `json:"h"`
	TickLow           float64            `json:"l"`
	VWAP              float64            `json:"a"`
	AverageTradeSize  float64            `json:"z"`
	StartTimestamp    int64              `json:"s"`
	EndTimestamp      int64              `json:"e"`
	OTC               *bool              `json:"otc"`
}

func (s *StockAggregate) ValidOpen() float64 {
	if s.Open == 0 {
		return s.TickOpen
	}
	return s.Open
}

func (c Client) SubscribeStockAggregates(client WebSocketClient, symbols []string, eventType StockEventTypeEnum) (err error) {
	// connect
	client.Dial(fmt.Sprintf("%s/stocks", c.websocketBaseURL), nil)
	// auth
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", c.token)))
	if err != nil {
		return
	}

	// subscribe
	// https://polygon.io/docs/stocks/ws_stocks_am
	channel := resolveStockChannel(symbols, eventType)
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channel)))
	if err != nil {
		return
	}

	return
}

func resolveStockChannel(symbols []string, eventType StockEventTypeEnum) string {
	var sb strings.Builder
	for i, symbol := range symbols {
		sb.WriteString(fmt.Sprintf("%s.%s", eventType, symbol))
		if i < len(symbols)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}
