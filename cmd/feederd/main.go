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

	priceFeeder, err := priceFeederFactory(cfg.Interval, cfg.PortableMarkets())
	if err != nil {
		log.WithError(err).Fatal("error while initializing price feeder")
	}

	indexedTargets := make(application.IndexedTargetsByMarket)
	for _, mkt := range cfg.Markets {
		targets := make(map[string]ports.TdexClient)
		for _, t := range mkt.CTargets {
			target, err := grpcclient.NewGRPCClient(t.RPCAddress, t.MacaroonsPath, t.TLSCertPath)
			if err != nil {
				log.WithError(err).Fatalf(
					"error while connecting with target %s", t.RPCAddress,
				)
			}
			targets[target.RPCAddress()] = target
		}
		mktKey := ports.MarketKey(mkt)
		indexedTargets[mktKey] = targets
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
