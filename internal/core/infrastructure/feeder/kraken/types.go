package krakenfeeder

import (
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

type price struct {
	basePrice  string
	quotePrice string
}

func (p *price) BasePrice() string {
	return p.basePrice
}

func (p *price) QuotePrice() string {
	return p.quotePrice
}

type priceFeed struct {
	market ports.Market
	price  *price
}

func (p *priceFeed) GetMarket() ports.Market {
	return p.market
}

func (p *priceFeed) GetPrice() ports.Price {
	return p.price
}

type market struct {
	baseAsset  string
	quoteAsset string
	ticker     string
}

func (m market) BaseAsset() string {
	return m.baseAsset
}

func (m market) QuoteAsset() string {
	return m.quoteAsset
}

func (m market) Ticker() string {
	return m.ticker
}
