package ports

import (
	"errors"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
	log "github.com/sirupsen/logrus"
)

// KrakenWebSocket is the interface to manage kraken web socket streams
type KrakenWebSocket interface {
	Connect(tickersToSubscribe []string) error
	Start() (chan TickerWithPrice, error)
	Close() error
}

type krakenWebSocket struct {
	krakenWS            *ws.Client
	tickerWithPriceChan chan TickerWithPrice
	isListening         bool
}

// NewKrakenWebSocket is a factory function for KrakenWebSocket interface
func NewKrakenWebSocket() KrakenWebSocket {
	return &krakenWebSocket{
		krakenWS:            ws.New(),
		tickerWithPriceChan: make(chan TickerWithPrice),
		isListening:         false,
	}
}

// Connect method will connect to the websocket kraken server, ping it and subscribe to tickers threads.
func (socket *krakenWebSocket) Connect(tickersToSubscribe []string) error {
	// connect to server
	err := socket.krakenWS.Connect()
	if err != nil {
		return err
	}
	// test if the server is alive
	err = socket.krakenWS.Ping()
	if err != nil {
		return err
	}

	// subscribe to tickers
	err = socket.krakenWS.SubscribeTicker(tickersToSubscribe)
	if err != nil {
		return err
	}

	return nil
}

func (socket *krakenWebSocket) Start() (chan TickerWithPrice, error) {
	if socket.krakenWS == nil {
		return nil, errors.New("Socket not connected")
	}

	if socket.isListening {
		return nil, errors.New("socket is already listening")
	}

	go func() {
		for obj := range socket.krakenWS.Listen() {
			log.Warn(obj)
			switch obj := obj.(type) {
			case error:
				log.Warn("Channel closed: ", obj)
				socket.Close()
				break
			case ws.DataUpdate:
				tickerUpdate, ok := obj.Data.(ws.TickerUpdate)
				if ok {
					result := TickerWithPrice{
						Ticker: tickerUpdate.Pair,
						Price:  tickerUpdate.Close.Today.(float64),
					}
					socket.tickerWithPriceChan <- result
				}
			default:
			}
		}
	}()

	return socket.tickerWithPriceChan, nil
}

func (socket *krakenWebSocket) Close() error {
	if !socket.krakenWS.IsConnected() {
		socket.krakenWS.Close()
		socket.krakenWS = nil
	}
	return nil
}
