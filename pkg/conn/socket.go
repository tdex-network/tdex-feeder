package conn

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// ConnectToSocket dials and returns a new client connection to a remote host
func ConnectToSocket(address string) *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: address, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Printf("Connected to %s", u.String())
	return c
}
