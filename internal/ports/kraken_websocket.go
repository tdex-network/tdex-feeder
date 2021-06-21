package ports

import (
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// KrakenWebSocket is the interface to manage kraken web socket streams
type KrakenWebSocket interface {
	Connect() error
	Start() (chan TickerWithPrice, error)
	Stop() error
}

type krakenWebSocket struct {
	krakenWebSocketConn *websocket.Conn
	tickerWithPriceChan chan TickerWithPrice
	tickersToSubscribe  []string
	quitChan            chan bool
}

// NewKrakenWebSocket is a factory function for KrakenWebSocket interface
func NewKrakenWebSocket(tickersToSubscribe []string) KrakenWebSocket {
	return &krakenWebSocket{
		krakenWebSocketConn: nil,
		tickerWithPriceChan: make(chan TickerWithPrice),
		tickersToSubscribe:  tickersToSubscribe,
		quitChan:            make(chan bool, 1),
	}
}

// Connect method will connect to the websocket kraken server, ping it and subscribe to tickers threads.
func (socket *krakenWebSocket) Connect() error {
	// connect to server
	url := url.URL{Scheme: "wss", Host: krakenWebSocketURL, Path: "/"}
	websocketConn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}

	socket.krakenWebSocketConn = websocketConn
	subscribeMsg := createSubscribeToMarketMessage(socket.tickersToSubscribe)
	return sendRequestMessage(socket.krakenWebSocketConn, subscribeMsg)
}

func (socket *krakenWebSocket) Start() (chan TickerWithPrice, error) {
	go socket.listen()
	return socket.tickerWithPriceChan, nil
}

func (socket *krakenWebSocket) listen() {
	for {
		select {
		case <-socket.quitChan:
			return
		default:
			_, message, err := socket.krakenWebSocketConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Warn("error: ", err)
				}
				break
			}

			tickerWithPrice, err := toTickerWithPrice(message)
			if err != nil {
				continue
			}

			socket.tickerWithPriceChan <- *tickerWithPrice
		}
	}
}

func (socket *krakenWebSocket) Stop() error {
	socket.quitChan <- true
	err := socket.krakenWebSocketConn.Close()
	if err != nil {
		return err
	}
	close(socket.tickerWithPriceChan)
	return nil
}
