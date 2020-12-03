package ports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	address = "ws.kraken.com"
	ticker = "LTC/USDT"
)

func createAndConnect() (KrakenWebSocket, error) {
	krakenWS := NewKrakenWebSocket()
	err := krakenWS.Connect(address, []string{ticker})
	return krakenWS, err
}

func TestConnectToKrakenWebSocket(t *testing.T) {
	_, err := createAndConnect()
	assert.Nil(t, err)
}

func TestListen(t *testing.T) {
	ws, err := createAndConnect()
	if err != nil {
		t.Error(err)
	}

	tickerWithPriceChan, err := ws.StartListen()
	assert.Nil(t, err)

	nextTickerWithPrice := <- tickerWithPriceChan
	assert.NotNil(t, nextTickerWithPrice)
} 