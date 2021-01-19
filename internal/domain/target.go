package domain

// Target interface is used to push new prices fetched from feeds
type Target interface {
	Push(marketPrice MarketPrice)
}
