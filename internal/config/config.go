package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

const (
	ConfigKey   = "CONFIG_PATH"
	LogLevelKey = "LOG_LEVEL"
)

var (
	defaultConfig = "."
)

func init() {
	viper.SetEnvPrefix("FEEDER")
	viper.AutomaticEnv()
	viper.SetDefault(ConfigKey, defaultConfig)
	viper.SetDefault(LogLevelKey, 4)

	configPath := viper.GetString(ConfigKey)
	if configPath != defaultConfig {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath(configPath)
	}
}

type Config struct {
	PriceFeeder      string              `mapstructure:"price_feeder"`
	Interval         int                 `mapstructure:"interval"`
	Targets          []Target            `mapstructure:"targets"`
	WellKnownMarkets map[string][]Market `mapstructure:"well_known_markets"`
}

func GetConfigPath() string {
	return viper.ConfigFileUsed()
}

func NewConfigFromFile() (Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	cfg := Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.PriceFeeder == "" {
		return fmt.Errorf("price_feeder must not be nil")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be a positive value")
	}
	if len(c.Targets) <= 0 {
		return fmt.Errorf("targets must not be empty")
	}
	for _, t := range c.Targets {
		if err := t.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c Config) MergeWellKnownMarkets(
	priceFeeder string, markets []ports.Market,
) error {
	mkts := make([]Market, 0, len(markets))
	for _, m := range markets {
		mkt := Market{m.BaseAsset(), m.QuoteAsset(), m.Ticker()}
		mkts = append(mkts, mkt)
	}

	if c.WellKnownMarkets == nil {
		c.WellKnownMarkets = make(map[string][]Market)
	}
	c.WellKnownMarkets[priceFeeder] = mkts

	raw := make(map[string][]map[string]string)
	for feeder, markets := range c.WellKnownMarkets {
		raw[feeder] = make([]map[string]string, 0, len(markets))
		for _, m := range markets {
			raw[feeder] = append(raw[feeder], m.RawMap())
		}
	}

	viper.Set("well_known_markets", raw)
	return viper.WriteConfig()
}
