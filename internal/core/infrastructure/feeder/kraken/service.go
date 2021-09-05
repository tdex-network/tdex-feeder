package krakenfeeder

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

const (
	// KrakenWebSocketURL is the base url to open a connection with kraken.
	// This can be tweaked if in the future it might change, even if unlikely.
	KrakenWebSocketURL = "ws.kraken.com"
)

type service struct {
	conn *websocket.Conn

	marketByTicker map[string]ports.Market
	feedChan       chan ports.PriceFeed
	quitChan       chan struct{}
}

func NewKrakenPriceFeeder(args ...interface{}) (ports.PriceFeeder, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of args")
	}

	markets, ok := args[0].([]ports.Market)
	if !ok {
		return nil, fmt.Errorf("unknown args type")
	}

	mktTickers := make([]string, 0, len(markets))
	mktByTicker := make(map[string]ports.Market)
	for _, mkt := range markets {
		mktTickers = append(mktTickers, mkt.Ticker())
		mktByTicker[mkt.Ticker()] = mkt
	}

	conn, err := connectAndSubscribe(mktTickers)
	if err != nil {
		return nil, err
	}
	conn.CloseHandler()

	return &service{
		conn:           conn,
		marketByTicker: mktByTicker,
		feedChan:       make(chan ports.PriceFeed),
		quitChan:       make(chan struct{}, 1),
	}, nil
}

func (k *service) Start() error {
	defer func() {
		close(k.feedChan)
		close(k.quitChan)
	}()

	mustReconnect, err := k.start()
	for mustReconnect {
		log.WithError(err).Warn("connection dropped unexpectedly. Trying to reconnect...")

		tickers := make([]string, 0, len(k.marketByTicker))
		for ticker := range k.marketByTicker {
			tickers = append(tickers, ticker)
		}

		var conn *websocket.Conn
		conn, err = connectAndSubscribe(tickers)
		if err != nil {
			return err
		}
		k.conn = conn

		log.Debug("connection and subscriptions re-established. Restarting...")
		mustReconnect, err = k.start()
	}

	return err
}

func (k *service) Stop() {
	k.quitChan <- struct{}{}
}

func (k *service) FeedChan() chan ports.PriceFeed {
	return k.feedChan
}

func (k *service) start() (mustReconnect bool, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			mustReconnect = true
		}
	}()

	for {
		select {
		case <-k.quitChan:
			err = k.conn.Close()
			return false, err
		default:
			_, message, err := k.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					panic(err)
				}
			}

			priceFeed := k.parseFeed(message)
			if priceFeed == nil {
				continue
			}

			k.feedChan <- priceFeed
		}
	}
}

func (k *service) parseFeed(msg []byte) ports.PriceFeed {
	var i []interface{}
	if err := json.Unmarshal(msg, &i); err != nil {
		return nil
	}
	if len(i) != 4 {
		return nil
	}

	ticker, ok := i[3].(string)
	if !ok {
		return nil
	}

	mkt, ok := k.marketByTicker[ticker]
	if !ok {
		return nil
	}

	ii, ok := i[1].(map[string]interface{})
	if !ok {
		return nil
	}

	iii, ok := ii["c"].([]interface{})
	if !ok {
		return nil
	}

	if len(iii) < 1 {
		return nil
	}
	priceStr, ok := iii[0].(string)
	if !ok {
		return nil
	}

	quotePrice, err := decimal.NewFromString(priceStr)
	if err != nil {
		return nil
	}
	basePrice := decimal.NewFromInt(1).Div(quotePrice)

	return &priceFeed{
		market: mkt,
		price: &price{
			basePrice:  basePrice.String(),
			quotePrice: quotePrice.String(),
		},
	}
}

func connectAndSubscribe(mktTickers []string) (*websocket.Conn, error) {
	url := fmt.Sprintf("wss://%s", KrakenWebSocketURL)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	msg := map[string]interface{}{
		"event": "subscribe",
		"pair":  mktTickers,
		"subscription": map[string]string{
			"name": "ticker",
		},
	}

	buf, _ := json.Marshal(msg)
	if err := conn.WriteMessage(websocket.TextMessage, buf); err != nil {
		return nil, fmt.Errorf("cannot subscribe to given markets: %s", err)
	}

	return conn, nil
}
