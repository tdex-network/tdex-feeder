package application

import "github.com/tdex-network/tdex-feeder/internal/domain"

type FeederService interface {
	Start() error
	Stop() error
}

type feederService struct {
	tdexFeeder domain.TdexFeeder
	krakenService krakenFeedService
}

func (feeder *feederService) Start() error {

}

func (feeder *feederService) Stop() error {
	
}