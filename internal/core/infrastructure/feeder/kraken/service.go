package krakenfeeder

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

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
	conn        *websocket.Conn
	writeTicker *time.Ticker
	lock        *sync.RWMutex
	chLock      *sync.Mutex

	marketByTicker map[string]ports.Market
	lastPriceFeed  ports.PriceFeed
	feedChan       chan ports.PriceFeed
	quitChan       chan struct{}
}

func NewKrakenPriceFeeder(args ...interface{}) (ports.PriceFeeder, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of args")
	}

	interval, ok := args[0].(int)
	if !ok {
		return nil, fmt.Errorf("unknown interval arg type")
	}

	markets, ok := args[1].([]ports.Market)
	if !ok {
		return nil, fmt.Errorf("unknown marktes arg type")
	}

	writeTicker := time.NewTicker(time.Duration(interval) * time.Millisecond)
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
		writeTicker:    writeTicker,
		lock:           &sync.RWMutex{},
		chLock:         &sync.Mutex{},
		marketByTicker: mktByTicker,
		feedChan:       make(chan ports.PriceFeed),
		quitChan:       make(chan struct{}, 1),
	}, nil
}

func (s *service) Start() error {
	mustReconnect, err := s.start()
	for mustReconnect {
		log.WithError(err).Warn("connection dropped unexpectedly. Trying to reconnect...")

		tickers := make([]string, 0, len(s.marketByTicker))
		for ticker := range s.marketByTicker {
			tickers = append(tickers, ticker)
		}

		var conn *websocket.Conn
		conn, err = connectAndSubscribe(tickers)
		if err != nil {
			return err
		}
		s.conn = conn

		log.Debug("connection and subscriptions re-established. Restarting...")
		mustReconnect, err = s.start()
	}

	return err
}

func (s *service) Stop() {
	s.quitChan <- struct{}{}
}

func (s *service) FeedChan() chan ports.PriceFeed {
	return s.feedChan
}

func (s *service) start() (mustReconnect bool, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			mustReconnect = true
		}
	}()

	go func() {
		for range s.writeTicker.C {
			s.writeToFeedChan()
		}
	}()

	for {
		select {
		case <-s.quitChan:
			s.writeTicker.Stop()
			s.closeChannels()
			err = s.conn.Close()
			return false, err
		default:
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					panic(err)
				}
			}

			priceFeed := s.parseFeed(message)
			if priceFeed == nil {
				continue
			}

			s.writePriceFeed(priceFeed)
		}
	}
}

func (s *service) readPriceFeed() ports.PriceFeed {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.lastPriceFeed
}

func (s *service) writePriceFeed(priceFeed ports.PriceFeed) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastPriceFeed = priceFeed
}

func (s *service) writeToFeedChan() {
	s.chLock.Lock()
	defer s.chLock.Unlock()

	priceFeed := s.readPriceFeed()
	if priceFeed != nil {
		s.feedChan <- priceFeed
	}
}

func (s *service) closeChannels() {
	s.chLock.Lock()
	defer s.chLock.Unlock()

	close(s.feedChan)
	close(s.quitChan)
}

func (s *service) parseFeed(msg []byte) ports.PriceFeed {
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

	mkt, ok := s.marketByTicker[ticker]
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
			basePrice:  basePrice.StringFixed(8),
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
