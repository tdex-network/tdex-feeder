package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/tdex-network/tdex-feeder/internal/config"
	"github.com/tdex-network/tdex-feeder/internal/core/application"
	grpcclient "github.com/tdex-network/tdex-feeder/internal/core/infrastructure/client/grpc"
	krakenfeeder "github.com/tdex-network/tdex-feeder/internal/core/infrastructure/feeder/kraken"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

type indexedPriceFeeders map[string]func(args ...interface{}) (ports.PriceFeeder, error)

func (i indexedPriceFeeders) supported() []string {
	keys := make([]string, 0, len(i))
	for k := range i {
		keys = append(keys, k)
	}
	return keys
}

var (
	priceFeeders = indexedPriceFeeders{
		"kraken": krakenfeeder.NewKrakenPriceFeeder,
	}
)

func main() {
	cfg, err := config.NewConfigFromFile()
	if err != nil {
		log.WithError(err).Fatalf(
			"error while reading config from file %s", config.GetConfigPath(),
		)
	}

	priceFeederFactory, ok := priceFeeders[cfg.PriceFeeder]
	if !ok {
		log.Fatalf(
			"price feeder must be one of: '%s'", priceFeeders.supported(),
		)
	}

	priceFeeder, err := priceFeederFactory(cfg.Interval)
	if err != nil {
		log.WithError(err).Fatal("error while initializing price feeder")
	}

	wellKnownMarkets := priceFeeder.WellKnownMarkets()
	knownMarketsExistInConfig := len(cfg.WellKnownMarkets) > 0 && len(cfg.WellKnownMarkets[cfg.PriceFeeder]) > 0
	if knownMarketsExistInConfig {
		for _, mkt := range cfg.WellKnownMarkets[cfg.PriceFeeder] {
			wellKnownMarkets = append(wellKnownMarkets, mkt)
		}
	}

	indexedTargets := make(application.IndexedTargetsByMarket)
	marketsByKey := make(map[string]ports.Market)
	for _, t := range cfg.Targets {
		var target ports.TdexClient
		if t.TdexdconnectURL != "" {
			target, err = grpcclient.NewGRPCClientFromURL(t.TdexdconnectURL)
			if err != nil {
				log.WithError(err).Fatalf(
					"error while connecting with target with url %s", t.TdexdconnectURL,
				)
			}
		} else {
			target, err = grpcclient.NewGRPCClient(t.RPCAddress, t.MacaroonsPath, t.TLSCertPath)
			if err != nil {
				log.WithError(err).Fatalf(
					"error while connecting with target %s", t.RPCAddress,
				)
			}
		}
		markets, err := target.ListMarkets()
		if err != nil {
			log.WithError(err).Fatalf(
				"failed to list markets for target %s", t.RPCAddress,
			)
		}

		for _, m := range markets {
			for _, wellKnownMarket := range wellKnownMarkets {
				if m.BaseAsset() == wellKnownMarket.BaseAsset() && m.QuoteAsset() == wellKnownMarket.QuoteAsset() {
					mktKey := ports.MarketKey(wellKnownMarket)
					if indexedTargets[mktKey] == nil {
						indexedTargets[mktKey] = make(map[string]ports.TdexClient)
					}
					indexedTargets[mktKey][t.RPCAddress] = target
					marketsByKey[mktKey] = wellKnownMarket
				}
			}
		}
	}

	markets := make([]ports.Market, 0, len(marketsByKey))
	for _, mkt := range marketsByKey {
		markets = append(markets, mkt)
	}

	if err := priceFeeder.SubscribeMarkets(markets); err != nil {
		log.WithError(err).Fatalf(
			"failed to subscribe price feeder %s to markets", cfg.PriceFeeder,
		)
	}

	if !knownMarketsExistInConfig {
		log.Info("writing price feeder's known markets to config")

		if err := cfg.MergeWellKnownMarkets(
			cfg.PriceFeeder, wellKnownMarkets,
		); err != nil {
			log.WithError(err).Fatal(
				"failed to write well known markets to config file",
			)
		}
	}

	app := application.NewService(priceFeeder, indexedTargets)

	defer app.Stop()

	log.Info("starting service")
	go func() {
		if err := app.Start(); err != nil {
			log.WithError(err).Fatal("service exited with error")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Info("shutting down")
}
