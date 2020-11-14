package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	defaultDaemon_endpoint    = "localhost:9000"
	defaultKraken_ws_endpoint = "ws.kraken.com"
	defaultBase_asset         = "lbtc"
	defaultQuote_asset        = "usd"
	defaultKraken_ticker      = "XBT/USD"
	defaultInterval           = 30
)

type Config struct {
	Daemon_endpoint    string   `json:"daemon_endpoint"`
	Daemon_macaroon    string   `json:"daemon_macaroon"`
	Kraken_ws_endpoint string   `json:"kraken_ws_endpoint"`
	Markets            []Market `json:"markets"`
}

type Market struct {
	Base_asset    string `json:"base_asset"`
	Quote_asset   string `json:"quote_asset"`
	Kraken_ticker string `json:"kraken_ticker"`
	Interval      int    `json:"interval"`
}

func DefaultConfig() Config {
	return Config{
		Daemon_endpoint:    defaultDaemon_endpoint,
		Kraken_ws_endpoint: defaultKraken_ws_endpoint,
		Markets: []Market{
			Market{
				Base_asset:    defaultBase_asset,
				Quote_asset:   defaultQuote_asset,
				Kraken_ticker: defaultKraken_ticker,
				Interval:      defaultInterval,
			},
		},
	}
}

func LoadConfigFromFile(filePath string) Config {
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	var config Config

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	return config
}

func LoadConfig(filePath string) Config {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Println("File not found. Loading default config.")
		return DefaultConfig()
	}
	return LoadConfigFromFile(filePath)
}
