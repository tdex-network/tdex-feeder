package grpcdaemon

type market struct {
	baseAsset  string
	quoteAsset string
}

func (m market) BaseAsset() string {
	return m.baseAsset
}

func (m market) QuoteAsset() string {
	return m.quoteAsset
}

func (m market) Ticker() string {
	return ""
}
