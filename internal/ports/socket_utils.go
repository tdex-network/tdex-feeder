package ports

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	krakenWebSocketURL = "ws.kraken.com"
)

type Subscription struct {
	Name     string `json:"name"`
	Interval int    `json:"interval,omitempty"`
	Token    string `json:"token,omitempty"`
	Depth    int    `json:"depth,omitempty"`
	Snapshop bool   `json:"snapshot,omitempty"`
}

// RequestMessage is the data structure used to create
// jsons in order subscribe to market updates on Kraken
type RequestMessage struct {
	Event        string        `json:"event"`
	Pair         []string      `json:"pair,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`
	Reqid        int           `json:"reqid,omitempty"`
}

// CreatePingMessage returns a RequestMessage struct
// with a ping Event.
func createPingMessage() RequestMessage {
	return RequestMessage{Event: "ping"}
}

// CreateSubscribeToMarketMessage gets a string with a market pair and returns
// a RequestMessage struct with instructions to subscrive to that market pair ticker.
func createSubscribeToMarketMessage(marketpairs []string) RequestMessage {
	s := Subscription{Name: "ticker"}
	return RequestMessage{"subscribe", marketpairs, &s, 0}
}

// SendRequestMessage gets a socket connection and a RequestMessage struct,
// marshalls the struct and sends the message using the socket.
func sendRequestMessage(c *websocket.Conn, m RequestMessage) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(b))
	if err != nil {
		return err
	}
	return nil
}

func toTickerWithPrice(message []byte) (*TickerWithPrice, error) {
	var result []interface{}
	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 4 {
		pair := result[3].(string)
		price, err := strconv.ParseFloat(result[1].(map[string]interface{})["c"].([]interface{})[0].(string), 64)
		if err != nil {
			return nil, err
		}

		return &TickerWithPrice{
			Ticker: pair,
			Price:  price,
		}, nil
	}

	return nil, errors.New("message is not a subscribe response")
}
