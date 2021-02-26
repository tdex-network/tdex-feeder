package adapters

import "errors"

var (
	// ErrDaemonEndpointIsEmpty is returned if the config contains an empty tdex-daemon endpoint
	ErrDaemonEndpointIsEmpty = errors.New("daemon endpoint is empty")
	// ErrKrakenEndpointIsEmpty is returned if the config contains an empty kraken WS endpoint
	ErrKrakenEndpointIsEmpty = errors.New("kraken websocket endpoint is empty")
	// ErrNeedAtLeastOneMarketToFeed is returned if the config does not contain market to feed
	ErrNeedAtLeastOneMarketToFeed = errors.New("need at least 1 market to feed")
	// ErrKrakenTickerIsEmpty is returned if the ticker specified in config is an empty string
	ErrKrakenTickerIsEmpty = errors.New("krakenTicker should not be an empty string")
	// ErrIntervalIsNotPositiveNumber is returned if the interval is < 0
	ErrIntervalIsNotPositiveNumber = errors.New("interval must be greater (or equal) than 0")
)

// ErrInvalidAssetHash is returned if the given string `asset` is not a valid asset hash string
type ErrInvalidAssetHash struct {
	asset string
}

var _ error = &ErrInvalidAssetHash{}

func (e ErrInvalidAssetHash) Error() string {
	return "the string '" + e.asset + "' is an invalid asset string."
}
