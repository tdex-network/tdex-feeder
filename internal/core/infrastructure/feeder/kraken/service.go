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

	marketByTicker      map[string]ports.Market
	latestFeedsByTicker map[string]ports.PriceFeed
	feedChan            chan ports.PriceFeed
	quitChan            chan struct{}
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
		conn:                conn,
		writeTicker:         writeTicker,
		lock:                &sync.RWMutex{},
		chLock:              &sync.Mutex{},
		marketByTicker:      mktByTicker,
		latestFeedsByTicker: make(map[string]ports.PriceFeed),
		feedChan:            make(chan ports.PriceFeed),
		quitChan:            make(chan struct{}, 1),
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
			// Referred to:
			//
			// https://support.kraken.com/hc/en-us/articles/360044504011-WebSocket-API-unexpected-disconnections-from-market-data-feeds
			//
			// Sometimes it can happen that the line below panics instead of
			// returning an UnexpectedCloseError. Because of this it's
			// mandatory here to recover a potential panic to signal that the
			// connection must be re-established.
			// Even in case the line below returns an UnexpectedCloseError,
			// this is used to panic so the deferred recover function is reused
			// to still signal the need for a reconnection with kraken websocket.
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

			s.writePriceFeed(priceFeed.GetMarket().Ticker(), priceFeed)
		}
	}
}

func (s *service) readPriceFeeds() []ports.PriceFeed {
	s.lock.RLock()
	defer s.lock.RUnlock()

	feeds := make([]ports.PriceFeed, 0, len(s.latestFeedsByTicker))
	for _, priceFeed := range s.latestFeedsByTicker {
		feeds = append(feeds, priceFeed)
	}
	return feeds
}

func (s *service) writePriceFeed(mktTicker string, priceFeed ports.PriceFeed) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.latestFeedsByTicker[mktTicker] = priceFeed
}

func (s *service) writeToFeedChan() {
	s.chLock.Lock()
	defer s.chLock.Unlock()

	priceFeeds := s.readPriceFeeds()
	for _, priceFeed := range priceFeeds {
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
