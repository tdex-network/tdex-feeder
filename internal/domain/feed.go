package domain

import "sync"

type MarketPrice struct {
	Market Market
	Price Price
}

type Feed interface {
	getId() string
	getMarketPriceFeed() <-chan MarketPrice
}

func merge(feeds ...Feed) <-chan MarketPrice {
	mergedChan := make(chan MarketPrice)
	var wg sync.WaitGroup

	wg.Add(len(feeds))
	for _, feed := range feeds {
		c := feed.getMarketPriceFeed()
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