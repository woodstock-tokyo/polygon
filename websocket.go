package polygon

import (
	"bytes"
	"encoding/json"
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
	ReadMessage() (messageType int, message []byte, err error)
}

//////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// polygon websocket //////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

// EventTypeEnum event type enum
type EventTypeEnum string

const (
	// EventTypeOther others
	EventTypeOther EventTypeEnum = ""
	// EventTypeAM minute aggregates
	EventTypeAM EventTypeEnum = "AM"
	// EventTypeA second aggregates
	EventTypeA EventTypeEnum = "A"
)

type Aggregate struct {
	Event             EventTypeEnum `json:"ev"`
	Symbol            string        `json:"sym"`
	TickVolume        int64         `json:"v"`
	AccumulatedVolume int64         `json:"av"`
	Open              float64       `json:"op"`
	TickVWAP          float64       `json:"vw"`
	TickOpen          float64       `json:"o"`
	TickClose         float64       `json:"c"`
	TickHigh          float64       `json:"h"`
	TickLow           float64       `json:"l"`
	VWAP              float64       `json:"a"`
	AverageTradeSize  float64       `json:"z"`
	StartTimestamp    int64         `json:"s"`
	EndTimestamp      int64         `json:"e"`
	OTC               *bool         `json:"otc"`
}

func (c Client) SubscribeAggregatesPerMinute(client WebSocketClient, pipe chan Aggregate, errorPipe chan error, symbols []string) (err error) {
	// connect
	client.Dial(fmt.Sprintf("%s/stocks", c.websocketBaseURL), nil)

	// auth
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", c.token)))
	if err != nil {
		return
	}

	// subscribe
	// https://polygon.io/docs/stocks/ws_stocks_am
	channel := resolveStockChannel(symbols, EventTypeAM)
	err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channel)))
	if err != nil {
		return
	}

	for {
		_, msg, err := client.ReadMessage()
		if err != nil {
			return err
		}
		msg = bytes.Trim(msg, "\x00")
		msgString := string(msg[:])
		if strings.Contains(msgString, "message") {
			fmt.Printf("%s\n", msgString)
			continue
		}

		// it's slice in the case of multiple items
		// https://polygon.io/docs/stocks/ws_getting-started
		var stocks []Aggregate
		err = json.Unmarshal(msg, &stocks)
		if err != nil {
			errorPipe <- err
			continue
		}

		for _, stock := range stocks {
			pipe <- stock
		}
	}
}

func (c Client) SubscribeAggregatesPerSecond(client WebSocketClient, aggregateChan chan Aggregate, errorChan chan error, symbols []string) (err error) {
	// connect
	client.Dial(fmt.Sprintf("%s/stocks", c.websocketBaseURL), nil)
	// auth
	if err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", c.token))); err != nil {
		return
	}

	// subscribe
	// https://polygon.io/docs/stocks/ws_stocks_a
	channel := resolveStockChannel(symbols, EventTypeA)
	fmt.Println("Subscribed to channel", channel)
	if err = client.WriteMessage(TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channel))); err != nil {
		return
	}

	for {
		_, msg, err := client.ReadMessage()
		if err != nil {
			errorChan <- err
			continue
		}

		// it's slice in the case of multiple items
		// https://polygon.io/docs/stocks/ws_getting-started
		var stocks []Aggregate
		err = json.Unmarshal(msg, &stocks)
		if err != nil {
			errorChan <- err
			continue
		}

		for _, stock := range stocks {
			aggregateChan <- stock
		}
	}
}

func resolveStockChannel(symbols []string, eventType EventTypeEnum) string {
	var sb strings.Builder
	for i, symbol := range symbols {
		sb.WriteString(fmt.Sprintf("%s.%s", eventType, symbol))
		if i < len(symbols)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}
