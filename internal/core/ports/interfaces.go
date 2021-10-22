package ports

type PriceFeeder interface {
	WellKnownMarkets() []Market
	SubscribeMarkets([]Market) error

	Start() error
	Stop()

	FeedChan() chan PriceFeed
}

type TdexClient interface {
	RPCAddress() string

	IsReady() (bool, error)
	ListMarkets() ([]Market, error)
	UpdateMarketPrice(market Market, price Price) error
}
