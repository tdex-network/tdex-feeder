package feeder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/pkg/testutils"
)

func TestFeeder(t *testing.T) {
	feed := domain.NewFeed()
	feedBis := domain.NewFeed()

	target := &testutils.MockTarget{
		MarketPrices: make([]domain.MarketPrice, 0),
	}

	feeder := NewFeeder(
		[]domain.Feed{feed, feedBis},
		[]domain.Target{target},
	)

	marketPrice := domain.MarketPrice{
		Market: domain.Market{
			BaseAsset:  "1111",
			QuoteAsset: "0000",
		},
		Price: domain.Price{
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

	assert.Equal(t, 15, len(target.MarketPrices))
}
