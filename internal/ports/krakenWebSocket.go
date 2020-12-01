package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type KrakenWebSocket interface {
	Connect(address string, tickersToSubscribe []string) error
	Read() (*TickerWithPrice, error)
	Close() error
}

type krakenWebSocket struct {
	connSocket *websocket.Conn
}

func NewKrakenWebSocket() KrakenWebSocket {
	return &krakenWebSocket{
		connSocket: nil,
	}
}

func (socket *krakenWebSocket) Connect(address string, tickersToSubscribe []string) error {
	conn, err := connectToSocket(address)
	if err != nil {
		return err
	}

	for _, ticker := range tickersToSubscribe {
		msg := createSubscribeToMarketMessage(ticker)
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			return err
		}
	}

	socket.connSocket = conn

	return nil
}

func (socket *krakenWebSocket) Close() error {
	err := socket.connSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	socket.connSocket = nil
	return err
}

func (socket *krakenWebSocket) Read() (*TickerWithPrice, error) {
	if socket.connSocket == nil {
		return nil, errors.New("Socket not connected")
	}

	var msgAsJson []interface{}

	_, message, err := socket.connSocket.ReadMessage()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(message), &msgAsJson)	
	if err != nil {
		return nil, err
	}

	if len(msgAsJson) < 4 {
		return nil, errors.New("Invalid message" + fmt.Sprint(message))
	}

	pricesJson := msgAsJson[1].(map[string]interface{})
	priceAsk := pricesJson["c"].([]interface{})
	price, _ := strconv.ParseFloat(priceAsk[0].(string), 64)

	ticker := msgAsJson[3].(string)

	return &TickerWithPrice{
		Ticker: ticker,
		Price: price,
	}, nil
}

// ConnectToSocket dials and returns a new client connection to a remote host
func connectToSocket(address string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: address, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return c, err
	}
	log.Info("Connected to socket:", u.String())
	return c, nil
}

// CreateSubscribeToMarketMessage gets a string with a market pair and returns
// a RequestMessage struct with instructions to subscrive to that market pair ticker.
func createSubscribeToMarketMessage(ticker string) requestMessage {
	s := subscription{Name: "ticker"}
	return requestMessage{"subscribe", []string{ticker}, &s, 0}
}