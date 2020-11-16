package marketinfo

import (
	"time"

	"github.com/tdex-network/tdex-feeder/config"
)

type MarketInfo struct {
	config   config.Market
	price    float64
	interval *time.Ticker
}

func DefaultMarketInfo(market config.Market) *MarketInfo {
	var marketInfo MarketInfo
	marketInfo.SetConfig(market)
	marketInfo.SetPrice(0.00)
	marketInfo.SetInterval(market.Interval)
	return &marketInfo
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
