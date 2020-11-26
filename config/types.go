package config

type Market struct {
	BaseAsset    string `json:"base_asset,required"`
	QuoteAsset   string `json:"quote_asset,required"`
	KrakenTicker string `json:"kraken_ticker,required"`
	Interval     int    `json:"interval,required"`
}
