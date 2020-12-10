package application

import "github.com/tdex-network/tdex-feeder/internal/domain"

type mockTarget struct {
	marketPrices []domain.MarketPrice
}

func (t *mockTarget) Push(marketPrice domain.MarketPrice) {
	t.marketPrices = append(t.marketPrices, marketPrice)
}
