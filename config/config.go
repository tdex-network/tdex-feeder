package config

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	// ConfigFilePathKey is the location of the config.json file.
	ConfigFilePathKey = "CONFIG_PATH"
	// LogLevelKey ...
	LogLevelKey = "LOG_LEVEL"
)

var vip *viper.Viper

func init() {
	vip = viper.New()
	vip.SetEnvPrefix("FEEDER")
	vip.AutomaticEnv()

	vip.SetDefault(LogLevelKey, 4)
	vip.SetDefault(ConfigFilePathKey, "./config.json")

	validate()
}

func GetConfigPath() string {
	return vip.GetString(ConfigFilePathKey)
}

func validate() {
	if err := validatePath(vip.GetString(ConfigFilePathKey)); err != nil {
		log.Fatal(err)
	}
}

func validatePath(path string) error {
	if path != "" {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return errors.New("not a file")
		}
	}

	return nil
}
