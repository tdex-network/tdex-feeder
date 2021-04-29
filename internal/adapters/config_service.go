package adapters

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/tdex-network/tdex-feeder/internal/application"
	"github.com/tdex-network/tdex-feeder/internal/domain"
)

// MarketJSON is the struct describing the shape of market specs in config JSON file
type MarketJSON struct {
	BaseAsset    string `json:"base_asset"`
	QuoteAsset   string `json:"quote_asset"`
	KrakenTicker string `json:"kraken_ticker"`
	Interval     int    `json:"interval"`
}

// ConfigJSON is the struct describing the shape of config JSON file
type ConfigJSON struct {
	DaemonEndpoint string       `json:"daemon_endpoint"`
	Markets        []MarketJSON `json:"markets"`
}

// Config is the config of the application retreived from config JSON file
type Config struct {
	daemonEndpoint  string
	markets         map[string]domain.Market
	marketIntervals map[domain.Market]time.Duration
}

// ToFeederService transforms a Config into FeederService
func (config *Config) ToFeederService() application.FeederService {
	feederSvc := application.NewFeederService(application.NewFeederServiceArgs{
		OperatorEndpoint: config.daemonEndpoint,
		TickerToMarket:   config.markets,
		MarketToInterval: config.marketIntervals,
	})

	return feederSvc
}

// UnmarshalJSON ...
func (config *Config) UnmarshalJSON(data []byte) error {
	jsonConfig := &ConfigJSON{}
	err := json.Unmarshal(data, jsonConfig)
	if err != nil {
		return err
	}

	err = jsonConfig.validate()
	if err != nil {
		return err
	}

	config.daemonEndpoint = jsonConfig.DaemonEndpoint

	configTickerToMarketMap := make(map[string]domain.Market)
	marketIntervalsMap := make(map[domain.Market]time.Duration)

	for _, marketJSON := range jsonConfig.Markets {
		market := domain.Market{
			BaseAsset:  marketJSON.BaseAsset,
			QuoteAsset: marketJSON.QuoteAsset,
		}

		configTickerToMarketMap[marketJSON.KrakenTicker] = market
		marketIntervalsMap[market] = time.Duration(marketJSON.Interval) * time.Millisecond
	}

	config.markets = configTickerToMarketMap
	config.marketIntervals = marketIntervalsMap

	return nil
}

func (configJson ConfigJSON) validate() error {
	if configJson.DaemonEndpoint == "" {
		return ErrDaemonEndpointIsEmpty
	}

	if len(configJson.Markets) == 0 {
		return ErrNeedAtLeastOneMarketToFeed
	}

	for _, marketJSON := range configJson.Markets {
		if marketJSON.KrakenTicker == "" {
			return ErrKrakenTickerIsEmpty
		}

		err := validateAssetString(marketJSON.BaseAsset)
		if err != nil {
			return err
		}

		err = validateAssetString(marketJSON.QuoteAsset)
		if err != nil {
			return err
		}

		if marketJSON.Interval < 0 {
			return ErrIntervalIsNotPositiveNumber
		}
	}

	return nil
}

func validateAssetString(asset string) error {
	const regularExpression = `[0-9A-Fa-f]{64}`

	matched, err := regexp.Match(regularExpression, []byte(asset))
	if err != nil {
		return err
	}

	if !matched {
		return ErrInvalidAssetHash{asset: asset}
	}

	return nil
}
