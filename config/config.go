package config

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
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

	err := validate()
	if err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	// this skip the check for default config file (avoid make test fail)
	vip.SetDefault(ConfigFilePathKey, "./config.json")
}

// GetConfigPath return the path of the config.json file
func GetConfigPath() string {
	return vip.GetString(ConfigFilePathKey)
}

func validate() error {
	return validatePath(vip.GetString(ConfigFilePathKey))
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
