package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/tdex-network/tdex-feeder/config"
	"github.com/tdex-network/tdex-feeder/internal/adapters"
)

const (
	containerName    = "tdexd-feeder-test"
	daemonEndpoint   = "127.0.0.1:9000"
	krakenWsEndpoint = "ws.kraken.com"
	// nigiriUrl = "https://nigiri.network/liquid/api"
	nigiriUrl = "http://localhost:3001"
	password  = "vulpemsecret"
)

func TestFeeder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	runDaemonAndInitConfigFile(t)
	t.Cleanup(stopAndDeleteContainer)

	go main()
	time.Sleep(30 * time.Second)
	os.Exit(0)
}

func runDaemonAndInitConfigFile(t *testing.T) {
	usdt := runDaemonAndCreateMarket(t)

	configJson := adapters.ConfigJson{
		DaemonEndpoint:   daemonEndpoint,
		KrakenWsEndpoint: krakenWsEndpoint,
		Markets: []adapters.MarketJson{
			adapters.MarketJson{
				KrakenTicker: "LTC/USDT",
				BaseAsset:    "5ac9f65c0efcc4775e0baec4ec03abdde22473cd3cf33c0419ca290e0751b225",
				QuoteAsset:   usdt,
				Interval:     500,
			},
		},
	}

	bytes, err := json.Marshal(configJson)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile(config.GetConfigPath(), bytes, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
}

func runDaemonAndCreateMarket(t *testing.T) string {
	_, err := execute(
		"docker", "run", "--name", containerName,
		"-p", "9000:9000",
		"-d",
		"-v", "tdexd:/.tdex-daemon",
		"-e", "TDEX_NETWORK=regtest",
		"-e", "TDEX_EXPLORER_ENDPOINT="+nigiriUrl,
		"-e", "TDEX_FEE_ACCOUNT_BALANCE_TRESHOLD=1000",
		"-e", "TDEX_BASE_ASSET=5ac9f65c0efcc4775e0baec4ec03abdde22473cd3cf33c0419ca290e0751b225",
		"-e", "TDEX_LOG_LEVEL=5",
		"--network=host",
		"tdexd:latest",
	)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)

	_, err = runCLICommand("config", "init")
	if err != nil {
		t.Error(err)
	}

	// init the wallet
	seed, err := runCLICommand("genseed")
	if err != nil {
		t.Error(err)
	}

	_, err = runCLICommand("init", "--seed", seed, "--password", password)
	if err != nil {
		t.Error(err)
	}

	_, err = runCLICommand("unlock", "--password", password)
	if err != nil {
		t.Error(err)
	}

	depositMarketJson, err := runCLICommand("depositmarket", "--base_asset", "", "--quote_asset", "")
	if err != nil {
		t.Error(err)
	}

	var depositMarketResult map[string]interface{}

	err = json.Unmarshal([]byte(depositMarketJson), &depositMarketResult)
	if err != nil {
		t.Error(t, err)
	}

	address := depositMarketResult["address"].(string)
	usdt := fundMarketAddress(t, address)

	return usdt
}

func stopAndDeleteContainer() {
	_, err := execute("docker", "stop", containerName)
	if err != nil {
		panic(err)
	}

	_, err = execute("docker", "container", "rm", containerName)
	if err != nil {
		panic(err)
	}
}

func fundMarketAddress(t *testing.T, address string) string {
	_, err := faucet(address)
	if err != nil {
		t.Error(err)
	}

	_, shitcoin, err := mint(address, 100)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(3 * time.Second)
	return shitcoin
}

func mint(address string, amount int) (string, string, error) {
	url := fmt.Sprintf("%s/mint", nigiriUrl)
	payload := map[string]interface{}{"address": address, "quantity": amount}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", errors.New("Internal server error")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	respBody := map[string]interface{}{}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return "", "", err
	}

	if respBody["asset"].(string) == "" {
		return mint(address, amount)
	}

	return respBody["txId"].(string), respBody["asset"].(string), nil
}

func faucet(address string) (string, error) {
	url := fmt.Sprintf("%s/faucet", nigiriUrl)
	payload := map[string]string{"address": address}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := map[string]string{}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return "", err
	}

	return respBody["txId"], nil
}

func execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	result := out.String()
	if err == nil {
		return result, nil
	}

	return out.String(), errors.New(fmt.Sprint(err) + ": " + stderr.String())
}

func runCLICommand(cliCommand string, args ...string) (string, error) {
	commandArgs := []string{"exec", containerName, "tdex", cliCommand}
	commandArgs = append(commandArgs, args...)
	output, err := execute("docker", commandArgs...)
	return output, err
}
