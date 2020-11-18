package marketinfo

import (
	"time"

	"github.com/tdex-network/tdex-feeder/config"
)

// MarketInfo stores the informations necessary for
// handling different market pair prices in real-time.
type MarketInfo struct {
	config   config.Market
	price    float64
	interval *time.Ticker
}

// DefaultMarketInfo returns a pointer to a MarketInfo struct
// with the default configurations.
func DefaultMarketInfo(market config.Market) *MarketInfo {
	return &MarketInfo{
		config:   market,
		price:    0.00,
		interval: time.NewTicker(time.Second * time.Duration(market.Interval)),
	}
}

func (marketInfo *MarketInfo) GetConfig() config.Market {
	return marketInfo.config
}

func (marketInfo *MarketInfo) SetConfig(market config.Market) {
	marketInfo.config = market
}

func (marketInfo *MarketInfo) GetPrice() float64 {
	return marketInfo.price
}

func (marketInfo *MarketInfo) SetPrice(value float64) {
	marketInfo.price = value
}

func (marketInfo *MarketInfo) GetInterval() *time.Ticker {
	return marketInfo.interval
}

func (marketInfo *MarketInfo) SetInterval(interval int) {
	marketInfo.interval = time.NewTicker(time.Second * time.Duration(interval))
}
