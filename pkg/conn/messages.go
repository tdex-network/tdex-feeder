package conn

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/tdex-network/tdex-feeder/config"
	"github.com/tdex-network/tdex-feeder/pkg/markets"
	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
)

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

func CreatePingMessage() RequestMessage {
	m := RequestMessage{Event: "ping"}
	return m
}

func CreateSubscribeToMarketMessage(market string) RequestMessage {
	s := Subscription{Name: "ticker"}
	m := RequestMessage{"subscribe", []string{market}, &s, 0}
	return m
}

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

// GetMessages keeps a loop that gets the messages from the remote host
// and calls a function to handle the received messages when necessary.
func GetMessages(done chan string, cSocket *websocket.Conn, clientgRPC pboperator.OperatorClient, marketsConfigs []config.Market) {
	defer close(done)
	numberOfMarkets := len(marketsConfigs)
	marketsInfos := make(markets.MarketsInformations, numberOfMarkets)

	for i, marketConfig := range marketsConfigs {
		marketsInfos[i] = markets.DefaultMarketInfo(marketConfig)
		defer marketsInfos[i].GetInterval().Stop()
	}

	for {
		_, message, err := cSocket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		//log.Println(string(message))
		handleMessages(message, marketsInfos, clientgRPC)
	}
}

func handleMessages(message []byte, marketsInfos markets.MarketsInformations, clientgRPC pboperator.OperatorClient) {
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
		go UpdateMarketPricegRPC(marketsInfos, clientgRPC)
	}
}
