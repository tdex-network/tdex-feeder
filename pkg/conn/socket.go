package conn

import (
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// ConnectToSocket dials and returns a new client connection to a remote host
func ConnectToSocket(address string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: address, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return c, err
	}
	log.Println("Connected to socket:", u.String())
	return c, nil
}
