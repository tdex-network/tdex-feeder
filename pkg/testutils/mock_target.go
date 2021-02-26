package testutils

import "github.com/tdex-network/tdex-feeder/internal/domain"

// MockTarget simulates a target interface
type MockTarget struct {
	MarketPrices []domain.MarketPrice
}

// Push is the target method
func (t *MockTarget) Push(marketPrice domain.MarketPrice) {
	t.MarketPrices = append(t.MarketPrices, marketPrice)
}
