package markets

import (
	"time"

	"github.com/tdex-network/tdex-feeder/config"
)

type MarketsInformations []*marketInfo

type marketInfo struct {
	config   config.Market
	price    float64
	interval *time.Ticker
}

func defaultMarketInfo(market config.Market) *marketInfo {
	var marketInfo marketInfo
	marketInfo.config = market
	marketInfo.price = 0.00
	marketInfo.interval = time.NewTicker(time.Second * time.Duration(market.Interval))
	return &marketInfo
}
