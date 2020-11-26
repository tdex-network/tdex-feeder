package marketinfo

import (
	"time"

	"github.com/tdex-network/tdex-feeder/config"
)

// MarketInfo stores the informations necessary for
// handling different market pair prices in real-time.
type MarketInfo struct {
	Config   config.Market
	LastSent time.Time
	Price    float64
}

// InitialMarketInfo returns a pointer to a MarketInfo struct
// with the default configurations.
func InitialMarketInfo(market config.Market) MarketInfo {
	return MarketInfo{
		Config:   market,
		LastSent: time.Now(),
		Price:    0,
	}
}
