package ports

// TickerWithPrice is a struct using to represent ticker to subscribe in kraken web socket feed
type TickerWithPrice struct {
	Ticker string
	Price  float64
}
