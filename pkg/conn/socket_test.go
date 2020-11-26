package conn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	kraken = "ws.kraken.com"
)

func TestConnectToSocket(t *testing.T) {
	_, err := ConnectToSocket(kraken)
	assert.Nil(t, err)
}
