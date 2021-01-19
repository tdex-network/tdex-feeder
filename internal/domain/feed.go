package domain

// MarketPrice represents a new price associated with a given market
type MarketPrice struct {
	Market Market
	Price  Price
}

// Feed represents a source of MarketPrice data
type Feed interface {
	AddMarketPrice(marketPrice MarketPrice)
	GetMarketPriceChan() <-chan MarketPrice
}

type feed struct {
	marketPriceChan chan MarketPrice
}

// NewFeed creates a Feed (i.e an empty channel)
func NewFeed() Feed {
	return &feed{
		marketPriceChan: make(chan MarketPrice),
	}
}

// AddMarketPrice send a new marketPrice value inside the Feed's channel.
func (f feed) AddMarketPrice(marketPrice MarketPrice) {
	f.marketPriceChan <- marketPrice
}

func (f feed) GetMarketPriceChan() <-chan MarketPrice {
	return f.marketPriceChan
}
