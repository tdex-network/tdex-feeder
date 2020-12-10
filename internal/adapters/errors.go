package adapters

import "errors"

var (
	ErrDaemonEndpointIsEmpty       = errors.New("daemon endpoint is empty")
	ErrKrakenEndpointIsEmpty       = errors.New("kraken websocket endpoint is empty")
	ErrNeedAtLeastOneMarketToFeed  = errors.New("need at least 1 market to feed")
	ErrKrakenTickerIsEmpty         = errors.New("krakenTicker should not be an empty string")
	ErrIntervalIsNotPositiveNumber = errors.New("interval must be greater (or equal) than 0")
)

type ErrInvalidAssetHash struct {
	asset string
}

func (e ErrInvalidAssetHash) Error() string {
	return "the string '" + e.asset + "' is an invalid asset string."
}
