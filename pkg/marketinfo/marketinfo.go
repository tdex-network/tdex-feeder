package marketinfo

import (
	"time"

	"github.com/tdex-network/tdex-feeder/config"
)

// MarketInfo stores the informations necessary for
// handling different market pair prices in real-time.
type MarketInfo struct {
	config config.Market
	price  float64
	ticker *time.Ticker
}

// InitialMarketInfo returns a pointer to a MarketInfo struct
// with the default configurations.
func InitialMarketInfo(market config.Market) *MarketInfo {
	return &MarketInfo{
		config: market,
		price:  0.00,
		ticker: time.NewTicker(time.Second * time.Duration(market.Interval)),
	}
}

func (marketInfo *MarketInfo) GetConfig() config.Market {
	return marketInfo.config
}

func (marketInfo *MarketInfo) GetPrice() float64 {
	return marketInfo.price
}

func (marketInfo *MarketInfo) SetPrice(value float64) {
	marketInfo.price = value
}

func (marketInfo *MarketInfo) GetTicker() *time.Ticker {
	return marketInfo.ticker
}
