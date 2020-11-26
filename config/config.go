package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	defaultDaemonEndpoint   = "localhost:9000"
	defaultKrakenWsEndpoint = "ws.kraken.com"
)

// Config defines the struct for the configuration JSON file
type Config struct {
	DaemonEndpoint   string   `json:"daemon_endpoint,required"`
	DaemonMacaroon   string   `json:"daemon_macaroon"`
	KrakenWsEndpoint string   `json:"kraken_ws_endpoint,required"`
	Markets          []Market `json:"markets,required"`
}

// DefaultConfig returns the datastructure needed
// for a default connection.
func defaultConfig() Config {
	return Config{
		DaemonEndpoint:   defaultDaemonEndpoint,
		KrakenWsEndpoint: defaultKrakenWsEndpoint,
		Markets:          nil,
	}
}

// LoadConfigFromFile reads a file with the intended running behaviour
// and returns a Config struct with the respective configurations.
func loadConfigFromFile(filePath string) (Config, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	var config Config

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return Config{}, err
	}
	err = checkConfigParsing(config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// checkConfigParsing checks if all the required fields
// were correctly loaded into the Config struct.
func checkConfigParsing(config Config) error {
	fields := reflect.ValueOf(config)
	for i := 0; i < fields.NumField(); i++ {
		tags := fields.Type().Field(i).Tag
		if strings.Contains(string(tags), "required") && fields.Field(i).IsZero() {
			return errors.New("Config required field is missing: " + string(tags))
		}
	}
	for _, market := range config.Markets {
		fields := reflect.ValueOf(market)
		for i := 0; i < fields.NumField(); i++ {
			tags := fields.Type().Field(i).Tag
			if strings.Contains(string(tags), "required") && fields.Field(i).IsZero() {
				return errors.New("Config required field is missing: " + string(tags))
			}
		}
	}
	return nil
}

// LoadConfig handles the default behaviour for loading
// config.json files. In case the file is not found,
// it loads the default config.
func LoadConfig(filePath string) (Config, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Debug("File not found: %s. Loading default config.\n", filePath)
		return defaultConfig(), nil
	}
	return loadConfigFromFile(filePath)
}
