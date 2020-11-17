package conn

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/tdex-network/tdex-feeder/pkg/marketinfo"
)

// RequestMessage is the data structure used to create
// jsons in order subscribe to market updates on Kraken
type RequestMessage struct {
	Event        string        `json:"event"`
	Pair         []string      `json:"pair,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`
	Reqid        int           `json:"reqid,omitempty"`
}

type Subscription struct {
	Name     string `json:"name"`
	Interval int    `json:"interval,omitempty"`
	Token    string `json:"token,omitempty"`
	Depth    int    `json:"depth,omitempty"`
	Snapshop bool   `json:"snapshot,omitempty"`
}

// CreatePingMessage returns a RequestMessage struct
// with a ping Event.
func CreatePingMessage() RequestMessage {
	return RequestMessage{Event: "ping"}
}

// CreateSubscribeToMarketMessage gets a string with a market pair and returns
// a RequestMessage struct with instructions to subscrive to that market pair ticker.
func CreateSubscribeToMarketMessage(marketpair string) RequestMessage {
	s := Subscription{Name: "ticker"}
	return RequestMessage{"subscribe", []string{marketpair}, &s, 0}
}

// SendRequestMessage gets a socket connection and a RequestMessage struct,
// marshalls the struct and sends the message using the socket.
func SendRequestMessage(c *websocket.Conn, m RequestMessage) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(b))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

// GetMessages keeps a loop that gets the data from the remote host
// and calls a function to handle the received json.
func GetMessages(done chan string, cSocket *websocket.Conn, marketsInfos []*marketinfo.MarketInfo) {
	defer close(done)
	for {
		_, message, err := cSocket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Println(string(message))
		handleMessages(message, marketsInfos)
	}
}

// handleMessages gets a json with a market pair updated price and
// a list of marketInfo structs and sets the price for that market pair
// on the respective struct.
func handleMessages(message []byte, marketsInfos []*marketinfo.MarketInfo) {
	var result []interface{}
	json.Unmarshal([]byte(message), &result)
	if len(result) == 4 {
		pricesJson := result[1].(map[string]interface{})
		priceAsk := pricesJson["c"].([]interface{})
		price, err := strconv.ParseFloat(priceAsk[0].(string), 64)
		if err == nil {
			for i, marketsInfo := range marketsInfos {
				if marketsInfo.GetConfig().Kraken_ticker == result[3] {
					marketsInfos[i].SetPrice(price)
				}
			}
		}
	}
}
