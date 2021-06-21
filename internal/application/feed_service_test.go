package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tdex-network/tdex-feeder/internal/domain"
	"github.com/tdex-network/tdex-feeder/pkg/feeder"
	"github.com/tdex-network/tdex-feeder/pkg/testutils"
)

const (
	baseAsset        = "5ac9f65c0efcc4775e0baec4ec03abdde22473cd3cf33c0419ca290e0751b225"
	quoteAsset       = "a64b14f3de72bc602d0786e6f034720a879a6b9339d59b09ddd49e1783ed227a"
	krakenTicker     = "XBT/USDT"
	krakenWsEndpoint = "ws.kraken.com"
)

func TestKrakenFeedService(t *testing.T) {
	tickerMap := make(map[string]domain.Market)
	tickerMap[krakenTicker] = domain.Market{
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
	}

	svc, err := NewKrakenFeedService(krakenWsEndpoint, tickerMap)
	if err != nil {
		t.Fatal(err)
	}
	err = svc.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer svc.Stop()

	feed := svc.GetFeed()
	target := &testutils.MockTarget{MarketPrices: []domain.MarketPrice{}}
	feeder := feeder.NewFeeder([]domain.Feed{feed}, []domain.Target{target})
	go func() {
		err := feeder.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(10 * time.Second)
	feeder.Stop()

	assert.Equal(t, true, len(target.MarketPrices) > 0)
}
