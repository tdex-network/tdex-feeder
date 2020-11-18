package conn

import (
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// ConnectToSocket dials and returns a new client connection to a remote host
func ConnectToSocket(address string) *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: address, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Debug("Connected to ", u.String())
	return c
}
