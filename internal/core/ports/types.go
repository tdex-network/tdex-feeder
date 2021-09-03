package ports

import (
	"crypto/sha256"
	"encoding/hex"
)

type Market interface {
	BaseAsset() string
	QuoteAsset() string
	Ticker() string
}

type Price interface {
	BasePrice() string
	QuotePrice() string
}

type PriceFeed interface {
	GetMarket() Market
	GetPrice() Price
}

func MarketKey(mkt Market) string {
	key := mkt.BaseAsset() + mkt.QuoteAsset()
	keyBytes := sha256.Sum256([]byte(key))
	return hex.EncodeToString(keyBytes[:])
}
