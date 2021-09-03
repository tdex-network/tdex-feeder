package application

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

type Service interface {
	Start() error
	Stop()
}

type IndexedTargetsByMarket map[string]map[string]ports.TdexDaemon

type service struct {
	priceFeeder      ports.PriceFeeder
	targetsByAddress map[string]ports.TdexDaemon
	targetsByMarket  IndexedTargetsByMarket

	lock *sync.Mutex
}

func NewService(
	priceFeeder ports.PriceFeeder, targetsByMarket IndexedTargetsByMarket,
) Service {
	targetsByAddress := make(map[string]ports.TdexDaemon)
	for _, targets := range targetsByMarket {
		for _, t := range targets {
			targetsByAddress[t.RPCAddress()] = t
		}
	}

	return &service{
		priceFeeder:      priceFeeder,
		targetsByMarket:  targetsByMarket,
		targetsByAddress: targetsByAddress,
		lock:             &sync.Mutex{},
	}
}

func (s *service) Start() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	errCount := 0
	for _, t := range s.targetsByAddress {
		ok, err := t.IsReady()
		if err != nil {
			log.WithError(err).Warnf("cannot connect to daemon %s", t.RPCAddress())
			errCount++
			continue
		}
		if !ok {
			log.Warnf("daemon %s is not ready", t.RPCAddress())
			errCount++
		}
	}

	if errCount == len(s.targetsByAddress) {
		return fmt.Errorf(
			"none of the provided targets are either reachable or ready to be fed " +
				"with prices",
		)
	}

	if errCount > 0 {
		log.Warn(
			"some targets are neither reachable or ready to be fed with prices. " +
				"Make sure they are up and reachable to receive price feeds",
		)
	}

	var err error
	go func() {
		if _err := s.priceFeeder.Start(); _err != nil {
			err = _err
		}
	}()

	for priceFeed := range s.priceFeeder.FeedChan() {
		mkt := priceFeed.GetMarket()
		price := priceFeed.GetPrice()
		marketKey := ports.MarketKey(mkt)
		targets := s.targetsByMarket[marketKey]
		log.Infof(
			"received price feed: market %s, price %s/%s",
			mkt.Ticker(), price.BasePrice(), price.QuotePrice(),
		)

		for _, target := range targets {
			t := target
			go func() {
				if isReady, _ := t.IsReady(); isReady {
					if err := t.UpdateMarketPrice(mkt, price); err != nil {
						log.WithError(err).Warnf(
							"error while updating target %s", t.RPCAddress(),
						)
						return
					}
					log.Infof(
						"updated market price (%s) for target %s",
						mkt.Ticker(), t.RPCAddress(),
					)
				}
			}()
		}
	}

	return err
}

func (s *service) Stop() {
	log.Info("stopping price feeder...")
	s.priceFeeder.Stop()
	log.Info("done")
}
