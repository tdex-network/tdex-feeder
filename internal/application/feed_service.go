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
	feed               domain.Feed
	krakenWebSocket    ports.KrakenWebSocket
	stopChan           chan bool
	tickersToMarketMap map[string]domain.Market
}

func NewKrakenFeedService(
	address string,
	tickersToMarketMap map[string]domain.Market,
) (FeedService, error) {
	newFeed := domain.NewFeed()

	tickersToSubscribe := make([]string, 0)
	for k := range tickersToMarketMap {
		tickersToSubscribe = append(tickersToSubscribe, k)
	}

	krakenSocket := ports.NewKrakenWebSocket()
	err := krakenSocket.Connect(address, tickersToSubscribe)
	if err != nil {
		return nil, err
	}

	return &krakenFeedService{
		krakenWebSocket:    krakenSocket,
		feed:               newFeed,
		stopChan:           make(chan bool),
		tickersToMarketMap: tickersToMarketMap,
	}, nil
}

// Start is the main function of krakenFeedService
// when start, the services is listening for new data from kraken server
func (f *krakenFeedService) Start() {
	listening := true
	log.Info("Kraken web socket feed is listening")
	tickerWithPriceChan, err := f.krakenWebSocket.StartListen()
	if err != nil {
		log.Fatal(err)
	}
	for listening {
		select {
		case <-f.stopChan:
			listening = false
			err := f.krakenWebSocket.Close()
			if err != nil {
				log.Fatal(err)
			}

			log.Info("Kraken web socket feed stopped")
			break
		case tickerWithPrice := <-tickerWithPriceChan:
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
	}
}

// Stop just send data to the stopChan in order to stop listening from kraken web socket
func (f *krakenFeedService) Stop() {
	f.stopChan <- true
}

// GetFeed is a getter function for kraken's feed member
func (f *krakenFeedService) GetFeed() domain.Feed {
	return f.feed
}
