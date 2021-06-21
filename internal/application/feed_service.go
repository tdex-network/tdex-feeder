package application

import (
	log "github.com/sirupsen/logrus"

	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/internal/ports"
)

// FeedService is the interface wrapping krakenWS and transform it into a domain.Feed
type FeedService interface {
	Start() error
	Stop() error
	GetFeed() domain.Feed
}

type krakenFeedService struct {
	feed               domain.Feed
	krakenWebSocket    ports.KrakenWebSocket
	tickersToMarketMap map[string]domain.Market
}

// NewKrakenFeedService is the factory function for FeedService
func NewKrakenFeedService(
	address string,
	tickersToMarketMap map[string]domain.Market,
) (FeedService, error) {
	newFeed := domain.NewFeed()

	tickersToSubscribe := make([]string, 0)
	for k := range tickersToMarketMap {
		tickersToSubscribe = append(tickersToSubscribe, k)
	}

	krakenSocket := ports.NewKrakenWebSocket(tickersToSubscribe)
	err := krakenSocket.Connect()
	if err != nil {
		return nil, err
	}

	return &krakenFeedService{
		krakenWebSocket:    krakenSocket,
		feed:               newFeed,
		tickersToMarketMap: tickersToMarketMap,
	}, nil
}

// Start is the main function of krakenFeedService
// when start, the services is listening for new data from kraken server
func (f *krakenFeedService) Start() error {
	log.Info("Kraken web socket feed is listening")
	tickerWithPriceChan, err := f.krakenWebSocket.Start()
	if err != nil {
		return err
	}

	go func() {
		for tickerWithPrice := range tickerWithPriceChan {
			log.Debug("Kraken web socket receive message = " + string(tickerWithPrice.Ticker))

			market, ok := f.tickersToMarketMap[tickerWithPrice.Ticker]
			if !ok {
				log.Debug("Market not found for ticker: ", tickerWithPrice.Ticker)
				continue
			}

			f.feed.AddMarketPrice(domain.MarketPrice{
				Market: market,
				Price: domain.Price{
					BasePrice:  1 / float32(tickerWithPrice.Price),
					QuotePrice: float32(tickerWithPrice.Price),
				},
			})

		}
	}()

	return nil
}

// Stop closes the connection with the kraken websocket
func (f *krakenFeedService) Stop() error {
	return f.krakenWebSocket.Stop()
}

// GetFeed is a getter function for kraken's feed member
func (f *krakenFeedService) GetFeed() domain.Feed {
	return f.feed
}
