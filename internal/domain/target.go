package domain

type Target interface {
	Push(marketPrice MarketPrice) error
}