package ports

type PriceFeeder interface {
	Start() error
	Stop()

	FeedChan() chan PriceFeed
}

type TdexClient interface {
	RPCAddress() string

	IsReady() (bool, error)
	UpdateMarketPrice(market Market, price Price) error
}
