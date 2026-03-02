// Copyright (c) 2026 Woodstock K.K.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
	aggregateChan := make(chan StockAggregate)
	errChan := make(chan error)
	websocketClient := &TestWebsocketClient{}

	go func() {
		client := NewClient(token)
		err := client.SubscribeStockAggregates(websocketClient, []string{"AAPL", "NVDA"}, StockEventTypeAM)
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

			var stocks []StockAggregate
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
