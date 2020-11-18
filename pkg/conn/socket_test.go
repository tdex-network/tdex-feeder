package conn

import (
	"testing"
)

const (
	kraken = "ws.kraken.com"
)

func TestConnectToSocket(t *testing.T) {
	_, err := ConnectToSocket(kraken)
	if err != nil {
		t.Fatal(err)
	}
}
