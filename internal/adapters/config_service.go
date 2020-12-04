package adapters

import (
	"encoding/json"
	"errors"
	"regexp"

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
}

func (config *Config) ToFeederService() application.FeederService {
	feederSvc := application.NewFeederService(application.NewFeederServiceArgs{
		KrakenWSaddress:  config.krakenWSaddress,
		OperatorEndpoint: config.daemonEndpoint,
		TickerToMarket:   config.markets,
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

	for _, marketJson := range jsonConfig.Markets {
		configTickerToMarketMap[marketJson.KrakenTicker] = domain.Market{
			BaseAsset:  marketJson.BaseAsset,
			QuoteAsset: marketJson.QuoteAsset,
		}
	}

	config.markets = configTickerToMarketMap

	return nil
}

func (configJson ConfigJson) validate() error {
	if configJson.DaemonEndpoint == "" {
		return errors.New("daemon endpoint is empty")
	}

	if configJson.KrakenWsEndpoint == "" {
		return errors.New("kraken websocket endpoint is empty")
	}

	if len(configJson.Markets) == 0 {
		return errors.New("need at least 1 market to feed")
	}

	for _, marketJson := range configJson.Markets {
		if marketJson.KrakenTicker == "" {
			return errors.New("krakenTicker is empty")
		}

		err := validateAssetString(marketJson.BaseAsset)
		if err != nil {
			return err
		}

		err = validateAssetString(marketJson.QuoteAsset)
		if err != nil {
			return err
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
		return errors.New(asset + " is an invalid asset string.")
	}

	return nil
}
