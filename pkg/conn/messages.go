package conn

import (
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/tdex-network/tdex-feeder/pkg/marketinfo"
	"github.com/tdex-network/tdex-protobuf/generated/go/operator"
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
func SendRequestMessage(c *websocket.Conn, m RequestMessage) error {
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

// HandleMessages is responsible for the perpetual loop of receiving messages
// from subscriptions, retrieving the price from them and send the gRPC request
// to update the market price in the predeterminated interval.
func HandleMessages(done chan string, cSocket *websocket.Conn, marketsInfos []marketinfo.MarketInfo, clientgRPC operator.OperatorClient) {
	defer close(done)
	for {
		_, message, err := cSocket.ReadMessage()
		if err != nil {
			log.Debug("read:", err)
			return
		}
		log.Debug(string(message))
		marketsInfos = retrievePriceFromMessage(message, marketsInfos)
		marketsInfos = checkInterval(marketsInfos, clientgRPC)
	}
}

// checkInterval handles the gRPC calls for UpdateMarketPrice
// at a predeterminated inteval for each market.
func checkInterval(marketsInfos []marketinfo.MarketInfo, clientgRPC operator.OperatorClient) []marketinfo.MarketInfo {
	for i, marketInfo := range marketsInfos {
		if time.Since(marketInfo.LastSent).Round(time.Second) == time.Duration(marketInfo.Config.Interval*int(math.Pow10(9))) {
			UpdateMarketPricegRPC(marketInfo, clientgRPC)
			marketInfo.LastSent = time.Now()
			marketsInfos[i] = marketInfo
		}
	}
	return marketsInfos
}

// retrievePriceFromMessage gets a message from a subscription and retrieves the
// price information, updating the price of the specific market.
func retrievePriceFromMessage(message []byte, marketsInfos []marketinfo.MarketInfo) []marketinfo.MarketInfo {
	var result []interface{}
	json.Unmarshal([]byte(message), &result)
	if len(result) == 4 {
		pricesJson := result[1].(map[string]interface{})
		priceAsk := pricesJson["c"].([]interface{})
		price, _ := strconv.ParseFloat(priceAsk[0].(string), 64)
		for i, marketInfo := range marketsInfos {
			if marketInfo.Config.KrakenTicker == result[3] {
				marketInfo.Price = price
				marketsInfos[i] = marketInfo
			}
		}
	}
	return marketsInfos
}
