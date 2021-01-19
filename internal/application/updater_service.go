package application

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/internal/ports"
)

// TdexDaemonTarget implements the domain.Target interface and manage interval for each market
type TdexDaemonTarget struct {
	Endpoint           string
	priceUpdater       ports.TdexDaemonPriceUpdater
	priceUpdaterLocker sync.Locker
	marketsToUpdate    map[domain.Market]domain.Price
	closeChan          chan bool
}

// NewTdexDaemonTarget configure a tdexDaemonUpdater using the endpoint
// and start goroutines depending of the configured intervals for each market.
func NewTdexDaemonTarget(
	tdexDaemonOperatorInterfaceEnpoint string,
	marketToIntervalMap map[domain.Market]time.Duration,
) domain.Target {
	now := time.Now()
	mapLastSent := make(map[domain.Market]time.Time)

	for market := range marketToIntervalMap {
		mapLastSent[market] = now
	}

	tdexTarget := &TdexDaemonTarget{
		Endpoint:           tdexDaemonOperatorInterfaceEnpoint,
		priceUpdater:       ports.NewTdexDaemonPriceUpdater(context.Background(), tdexDaemonOperatorInterfaceEnpoint),
		priceUpdaterLocker: &sync.Mutex{},
		closeChan:          make(chan bool, 1),
		marketsToUpdate:    make(map[domain.Market]domain.Price),
	}

	for market, interval := range marketToIntervalMap {
		go func(duration time.Duration, market domain.Market) {
			for {
				select {
				case <-tdexTarget.closeChan:
					log.Info("Stop the tdex updater")
					break
				case <-time.After(duration):
					tdexTarget.updatePrice(market)
					continue
				}
			}
		}(interval, market)
	}

	return tdexTarget
}

// Push is a method of the Target interface
// The tdexDaemonTarget stores the marketPrice in a local cache.
func (daemon *TdexDaemonTarget) Push(marketPrice domain.MarketPrice) {
	daemon.priceUpdaterLocker.Lock()
	defer daemon.priceUpdaterLocker.Unlock()

	daemon.marketsToUpdate[marketPrice.Market] = marketPrice.Price
}

// Stop is used to stop all the goroutines launched in NewTdexDaemonTarget
func (daemon *TdexDaemonTarget) Stop() {
	daemon.closeChan <- true
}

func (daemon *TdexDaemonTarget) updatePrice(market domain.Market) {
	daemon.priceUpdaterLocker.Lock()
	defer daemon.priceUpdaterLocker.Unlock()

	price, ok := daemon.marketsToUpdate[market]
	if ok {
		err := daemon.priceUpdater.UpdateMarketPrice(context.Background(), domain.MarketPrice{Market: market, Price: price})
		if err != nil {
			log.Error("error updatePrice: ", err)
			return
		}
		delete(daemon.marketsToUpdate, market)
	}
}
