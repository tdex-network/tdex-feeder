// Copyright (c) 2020 The VulpemVentures developers

// Feeder allows to connect an external price feed to the TDex Daemon to determine the current market price.
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tdex-network/tdex-feeder/internal/adapters"
	"github.com/tdex-network/tdex-feeder/internal/application"
)

const (
	envConfigPathKey = "TDEX_FEEDER_CONFIG_PATH"
	defaultConfigPath = "./config.json"
)

func main() {
	// Interrupt Notification.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// retrieve feeder service from config file
	envConfigPath := os.Getenv(envConfigPathKey)
	if envConfigPath == "" {
		envConfigPath = defaultConfigPath
	}
	feeder := configFileToFeederService(envConfigPath)


	go func ()  {
		err := feeder.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// check for interupt
	for range interrupt {
		log.Println("Shutting down the feeder")
		time.Sleep(time.Second)
		err := feeder.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func configFileToFeederService(configFilePath string) application.FeederService {
	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	configBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	config := &adapters.Config{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatal(err)
	}

	feeder := config.ToFeederService()
	return feeder
}

