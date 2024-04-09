package polygon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

type TestWebsocketClient struct {
	conn *websocket.Conn
}

func (c *TestWebsocketClient) Dial(urlStr string, reqHeader http.Header) {
	c.conn, _ = websocket.Dial(urlStr, "", "http://localhost")
}

func (c *TestWebsocketClient) WriteMessage(messageType int, data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *TestWebsocketClient) ReadMessage() (messageType int, message []byte, err error) {
	var msg = make([]byte, 512)
	_, err = c.conn.Read(msg)
	return 1, msg, err
}

// NOTE: it takes more than a minute
func TestSubscribeAggregatesPerMinute(t *testing.T) {
	aggregateChan := make(chan Aggregate)
	errChan := make(chan error)
	websocketClient := &TestWebsocketClient{}

	go func() {
		client := NewClient(token)
		err := client.SubscribeAggregates(websocketClient, []string{"AAPL", "NVDA"}, EventTypeAM)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		for {
			_, msg, err := websocketClient.ReadMessage()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			msg = bytes.Trim(msg, "\x00")
			msgString := string(msg[:])
			// if it's a message object, just print it
			if strings.Contains(msgString, "message") {
				fmt.Printf("%s\n", msgString)
				continue
			}

			var stocks []Aggregate
			err = json.Unmarshal(msg, &stocks)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			for _, stock := range stocks {
				aggregateChan <- stock
			}
		}
	}()

	select {
	case aggregate := <-aggregateChan:
		fmt.Printf("%+v", aggregate)
	case err := <-errChan:
		t.Errorf("unexpected error: %v", err)
	}
}
