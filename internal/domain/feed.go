package domain

import (
	"sync"

	"github.com/google/uuid"
)

type MarketPrice struct {
	Market Market
	Price Price
}

type Feed interface {
	AddMarketPrice(marketPrice MarketPrice)
	getMarketPriceChan() <-chan MarketPrice	
}

type feed struct {
	id string
	marketPriceChan chan MarketPrice
}

func NewFeed() (Feed, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &feed{
		id: uuid.String(),
		marketPriceChan: make(chan MarketPrice),
	}, nil
}

func (f feed) AddMarketPrice(marketPrice MarketPrice) {
	f.marketPriceChan <- marketPrice
}

func (f feed) getMarketPriceChan() <-chan MarketPrice {
	return f.marketPriceChan
}

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