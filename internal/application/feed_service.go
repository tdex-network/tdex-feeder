package application

import (
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
	listening bool
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
		listening: false,
		tickersToMarketMap: tickersToMarketMap,
	}, nil
}

func (f *krakenFeedService) GetFeed() domain.Feed {
	return f.feed
}

func (f *krakenFeedService) Start() {
	log.Println("Start listening kraken service")
	for f.listening {
		tickerWithPrice, err := f.krakenWebSocket.Read()
		if err != nil {
			log.Debug("Read message error: ", err)
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

func (f *krakenFeedService) Stop() {
	f.listening = false
	f.krakenWebSocket.Close()
}