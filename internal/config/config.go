package config

import (
	"github.com/spf13/viper"
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

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
}

func GetConfigPath() string {
	configPath := viper.GetString(ConfigKey)
	return configPath
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
