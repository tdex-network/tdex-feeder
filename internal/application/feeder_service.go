package application

import (
	"time"

	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/pkg/feeder"
)

// FeederService is a tdex-configured feeder. It takes data from KrakenWS endpoint and update a tdex daemon
type FeederService interface {
	Start() error
	Stop() error
}

type feederService struct {
	tdexFeeder    feeder.Service
	krakenService FeedService
	target        *TdexDaemonTarget
}

// NewFeederServiceArgs is a wrapper for NewFeederService arguments
type NewFeederServiceArgs struct {
	OperatorEndpoint string
	MacaroonsPath    string
	TLSCertPath      string
	MarketToInterval map[domain.Market]time.Duration
	KrakenWSaddress  string
	TickerToMarket   map[string]domain.Market
}

// NewFeederService is the factory function for the FeederService
func NewFeederService(args NewFeederServiceArgs) (FeederService, error) {
	target, err := NewTdexDaemonTarget(
		args.OperatorEndpoint, args.MacaroonsPath, args.TLSCertPath,
		args.MarketToInterval,
	)
	if err != nil {
		return nil, err
	}

	krakenFeedService, err := NewKrakenFeedService(args.KrakenWSaddress, args.TickerToMarket)
	if err != nil {
		return nil, err
	}

	feeder := feeder.NewFeeder(
		[]domain.Feed{krakenFeedService.GetFeed()},
		[]domain.Target{target},
	)

	return &feederService{
		tdexFeeder:    feeder,
		krakenService: krakenFeedService,
		target:        target.(*TdexDaemonTarget),
	}, nil
}

func (feeder *feederService) Start() error {
	err := feeder.krakenService.Start()
	if err != nil {
		return err
	}

	return feeder.tdexFeeder.Start()
}

func (feeder *feederService) Stop() error {
	err := feeder.krakenService.Stop()
	if err != nil {
		return err
	}
	feeder.target.Stop()
	feeder.tdexFeeder.Stop()
	return nil
}
