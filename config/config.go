package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
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
	Daemon_endpoint    string   `json:"daemon_endpoint,required"`
	Daemon_macaroon    string   `json:"daemon_macaroon"`
	Kraken_ws_endpoint string   `json:"kraken_ws_endpoint,required"`
	Markets            []Market `json:"markets,required"`
}

type Market struct {
	Base_asset    string `json:"base_asset,required"`
	Quote_asset   string `json:"quote_asset,required"`
	Kraken_ticker string `json:"kraken_ticker,required"`
	Interval      int    `json:"interval,required"`
}

// DefaultConfig returns the datastructure needed
// for a default connection.
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

// LoadConfigFromFile reads a file with the intended running behaviour
// and returns a Config struct with the respective configurations.
func LoadConfigFromFile(filePath string) Config {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()

	var config Config

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(byteValue, &config)

	checkConfigParsing(config)

	return config
}

// checkConfigParsing checks if all the required fields
// were correctly loaded into the Config struct.
func checkConfigParsing(config Config) {
	fields := reflect.ValueOf(config)
	for i := 0; i < fields.NumField(); i++ {
		tags := fields.Type().Field(i).Tag
		if strings.Contains(string(tags), "required") && fields.Field(i).IsZero() {
			log.Println(errors.New("Config required field is missing: " + string(tags)))
		}
	}
	for _, market := range config.Markets {
		fields := reflect.ValueOf(market)
		for i := 0; i < fields.NumField(); i++ {
			tags := fields.Type().Field(i).Tag
			if strings.Contains(string(tags), "required") && fields.Field(i).IsZero() {
				log.Println(errors.New("Config required field is missing: " + string(tags)))
			}
		}
	}
}

// LoadConfig handles the default behaviour for loading
// config.json files. In case the file is not found,
// it loads the default config.
func LoadConfig(filePath string) Config {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Println("File not found. Loading default config.")
		return DefaultConfig()
	}
	return LoadConfigFromFile(filePath)
}
