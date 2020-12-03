package ports

import (
	"errors"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
	log "github.com/sirupsen/logrus"
)

type KrakenWebSocket interface {
	Connect(address string, tickersToSubscribe []string) error
	StartListen() (chan TickerWithPrice, error)
	Close() error
}

type krakenWebSocket struct {
	krakenWS *ws.Client
	tickerWithPriceChan chan TickerWithPrice
}

func NewKrakenWebSocket() KrakenWebSocket {
	return &krakenWebSocket{
		krakenWS: ws.New(),
		tickerWithPriceChan: make(chan TickerWithPrice),
	}
}

func (socket *krakenWebSocket) Connect(address string, tickersToSubscribe []string) error {
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

func (socket *krakenWebSocket) Close() error {
	socket.krakenWS.Close()
	socket.krakenWS = nil
	return nil
}

func (socket *krakenWebSocket) StartListen() (chan TickerWithPrice, error) {
	if socket.krakenWS == nil {
		return nil, errors.New("Socket not connected")
	}

	go func() {
		for obj := range socket.krakenWS.Listen() {
			switch obj := obj.(type) {
			case error:
				log.Debug("Channel closed: ", obj)
			case ws.DataUpdate:
				tickerUpdate, ok := obj.Data.(ws.TickerUpdate)
				if ok {
					result := TickerWithPrice{
						Ticker: tickerUpdate.Pair,
						Price: tickerUpdate.Close.Today.(float64),
					}
					socket.tickerWithPriceChan <- result
				}
			}
		}
	}()

	return socket.tickerWithPriceChan, nil
}