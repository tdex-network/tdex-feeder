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
	t.Run("should connect to a valid kraken web socket address", func(t *testing.T) {
		_, err := createAndConnect()
		assert.Nil(t, err)
	})
}

func TestRead(t *testing.T) {
	t.Run("should read marketInfo from websocket kraken stream", func(t *testing.T) {
		ws, err := createAndConnect()
		if err != nil {
			t.Error(err)
		}

		tickerWithPrice, err := ws.Read()
		assert.Nil(t, err)
		assert.NotNil(t, tickerWithPrice)
	})
} 