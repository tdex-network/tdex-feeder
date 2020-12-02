package application

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/internal/ports"
)

type FeedService interface {
	Start() 
	Stop()
	GetFeed() domain.Feed
}

type krakenFeedService struct {
	feed domain.Feed
	krakenWebSocket ports.KrakenWebSocket
	stopChan chan bool
	tickersToMarketMap map[string]domain.Market
}

func NewKrakenFeedService(
	address string, 
	tickersToMarketMap map[string]domain.Market,
) (FeedService, error) {
	newFeed, err := domain.NewFeed()
	if err != nil {
		return nil, err
	}

	tickersToSubscribe := make([]string, 0)
	for k := range tickersToMarketMap {
		tickersToSubscribe = append(tickersToSubscribe, k)
	}

	krakenSocket := ports.NewKrakenWebSocket()
	err = krakenSocket.Connect(address, tickersToSubscribe)
	if err != nil {
		return nil, err
	}
	
	return &krakenFeedService{
		krakenWebSocket: krakenSocket,
		feed: newFeed,
		stopChan: make(chan bool),
		tickersToMarketMap: tickersToMarketMap,
	}, nil
}

func (f *krakenFeedService) GetFeed() domain.Feed {
	return f.feed
}

func (f *krakenFeedService) Start() {
	listening := true
	log.Println("Start listening kraken service")
	for listening {
		select {
		case <-f.stopChan:
			listening = false
			err := f.krakenWebSocket.Close()
			if err != nil {
				log.Fatal(err)
			}

			log.Info("Feed service stopped")
			break;
		case <-time.After(500 * time.Millisecond):
			log.Info("Read socket interval")
			tickerWithPrice, err := f.krakenWebSocket.Read()
			if (tickerWithPrice != nil) {
				log.Info("msg =" + string(tickerWithPrice.Ticker))
			}
			if err != nil {
				log.Debug("Read message error: ", err)
				continue
			}

			market, ok := f.tickersToMarketMap[tickerWithPrice.Ticker]
			if !ok {
				log.Debug("Market not found for ticker: ", tickerWithPrice.Ticker)
			}

			f.feed.AddMarketPrice(domain.MarketPrice{
				Market: market,
				Price: domain.Price{
					BasePrice: 1 / float32(tickerWithPrice.Price),
					QuotePrice: float32(tickerWithPrice.Price),
				},
			})
		}
	}
}

func (f *krakenFeedService) Stop() {
	f.stopChan <- true
}