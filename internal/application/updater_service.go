package application

import (
	"context"

	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/internal/ports"
)

// Implements the domain.Target interface
type TdexDaemonTarget struct {
	Endpoint string
	priceUpdater ports.TdexDaemonPriceUpdater
}

func NewTdexDaemonTarget(tdexDaemonOperatorInterfaceEnpoint string) domain.Target {
	return &TdexDaemonTarget{
		Endpoint: tdexDaemonOperatorInterfaceEnpoint,
		priceUpdater: ports.NewTdexDaemonPriceUpdater(context.Background(), tdexDaemonOperatorInterfaceEnpoint),
	}
}

func (daemon *TdexDaemonTarget) Push(marketPrice domain.MarketPrice) error {
	return daemon.priceUpdater.UpdateMarketPrice(context.Background(), marketPrice)
}