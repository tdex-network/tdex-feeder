package domain

import (
	"sync"
)

type MarketPrice struct {
	Market Market
	Price  Price
}

// Feed represents a source of MarketPrice data
type Feed interface {
	AddMarketPrice(marketPrice MarketPrice)
	getMarketPriceChan() <-chan MarketPrice
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

func (f feed) getMarketPriceChan() <-chan MarketPrice {
	return f.marketPriceChan
}

// merge gathers several feeds into a unique channel
func merge(feeds ...Feed) <-chan MarketPrice {
	mergedChan := make(chan MarketPrice)
	var wg sync.WaitGroup

	wg.Add(len(feeds))
	for _, feed := range feeds {
		c := feed.getMarketPriceChan()
		go func(c <-chan MarketPrice) {
			for marketPrice := range c {
				mergedChan <- marketPrice
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(mergedChan)
	}()

	return mergedChan
}
