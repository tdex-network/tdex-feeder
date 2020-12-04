package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeeder(t *testing.T) {
	feed := NewFeed()
	feedBis := NewFeed()

	target := &mockTarget{
		marketPrices: make([]MarketPrice, 0),
	}

	feeder := NewTdexFeeder([]Feed{feed, feedBis}, []Target{target})

	marketPrice := MarketPrice{
		Market: Market{
			BaseAsset:  "1111",
			QuoteAsset: "0000",
		},
		Price: Price{
			BasePrice:  0.2,
			QuotePrice: 1,
		},
	}

	go func() {
		err := feeder.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)
	assert.Equal(t, true, feeder.IsRunning())

	go func() {
		for i := 0; i < 5; i++ {
			feedBis.AddMarketPrice(marketPrice)
		}
	}()

	for i := 0; i < 10; i++ {
		feed.AddMarketPrice(marketPrice)
	}

	time.Sleep(500 * time.Millisecond)
	feeder.Stop()

	assert.Equal(t, 15, len(target.marketPrices))
}

type mockTarget struct {
	marketPrices []MarketPrice
}

func (t *mockTarget) Push(marketPrice MarketPrice) {
	t.marketPrices = append(t.marketPrices, marketPrice)
}
