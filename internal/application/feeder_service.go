package application

import (
	log "github.com/sirupsen/logrus"
	"github.com/tdex-network/tdex-feeder/internal/domain"
)

type FeederService interface {
	Start() error
	Stop() error
}

type feederService struct {
	tdexFeeder domain.TdexFeeder
	krakenService FeedService
}

type NewFeederServiceArgs struct {
	OperatorEndpoint string
	KrakenWSaddress string
	TickerToMarket map[string]domain.Market
}

func NewFeederService(args NewFeederServiceArgs) FeederService {
	target := NewTdexDaemonTarget(args.OperatorEndpoint)
	krakenFeedService, err := NewKrakenFeedService(args.KrakenWSaddress, args.TickerToMarket)
	if err != nil {
		log.Fatal(err)
	}

	feeder := domain.NewTdexFeeder([]domain.Feed{krakenFeedService.GetFeed()}, []domain.Target{target})

	return &feederService{
		tdexFeeder: feeder,
		krakenService: krakenFeedService,
	}
}

func (feeder *feederService) Start() error {
	go feeder.krakenService.Start()
	err := feeder.tdexFeeder.Start()
	return err
}

func (feeder *feederService) Stop() error {
	feeder.krakenService.Stop()
	feeder.tdexFeeder.Stop()
	return nil	
}