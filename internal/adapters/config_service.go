package adapters

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/tdex-network/tdex-feeder/internal/application"
	"github.com/tdex-network/tdex-feeder/internal/domain"
)

type MarketJson struct {
	BaseAsset    string `json:"base_asset"`
	QuoteAsset   string `json:"quote_asset"`
	KrakenTicker string `json:"kraken_ticker"`
	Interval     int    `json:"interval"`
}

type ConfigJson struct {
	DaemonEndpoint   string       `json:"daemon_endpoint"`
	KrakenWsEndpoint string       `json:"kraken_ws_endpoint"`
	Markets          []MarketJson `json:"markets"`
}

type Config struct {
	daemonEndpoint  string
	krakenWSaddress string
	markets         map[string]domain.Market
	marketIntervals map[domain.Market]time.Duration
}

func (config *Config) ToFeederService() application.FeederService {
	feederSvc := application.NewFeederService(application.NewFeederServiceArgs{
		KrakenWSaddress:  config.krakenWSaddress,
		OperatorEndpoint: config.daemonEndpoint,
		TickerToMarket:   config.markets,
		MarketToInterval: config.marketIntervals,
	})

	return feederSvc
}

func (config *Config) UnmarshalJSON(data []byte) error {
	jsonConfig := &ConfigJson{}
	err := json.Unmarshal(data, jsonConfig)
	if err != nil {
		return err
	}

	err = jsonConfig.validate()
	if err != nil {
		return err
	}

	config.daemonEndpoint = jsonConfig.DaemonEndpoint
	config.krakenWSaddress = jsonConfig.KrakenWsEndpoint

	configTickerToMarketMap := make(map[string]domain.Market)
	marketIntervalsMap := make(map[domain.Market]time.Duration)

	for _, marketJson := range jsonConfig.Markets {
		market := domain.Market{
			BaseAsset:  marketJson.BaseAsset,
			QuoteAsset: marketJson.QuoteAsset,
		}

		configTickerToMarketMap[marketJson.KrakenTicker] = market
		marketIntervalsMap[market] = time.Duration(marketJson.Interval) * time.Millisecond
	}

	config.markets = configTickerToMarketMap
	config.marketIntervals = marketIntervalsMap

	return nil
}

func (configJson ConfigJson) validate() error {
	if configJson.DaemonEndpoint == "" {
		return ErrDaemonEndpointIsEmpty
	}

	if configJson.KrakenWsEndpoint == "" {
		return ErrKrakenEndpointIsEmpty
	}

	if len(configJson.Markets) == 0 {
		return ErrNeedAtLeastOneMarketToFeed
	}

	for _, marketJson := range configJson.Markets {
		if marketJson.KrakenTicker == "" {
			return ErrKrakenTickerIsEmpty
		}

		err := validateAssetString(marketJson.BaseAsset)
		if err != nil {
			return err
		}

		err = validateAssetString(marketJson.QuoteAsset)
		if err != nil {
			return err
		}

		if marketJson.Interval < 0 {
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
		return ErrInvalidAssetHash{asset: asset}.Error()
	}

	return nil
}
