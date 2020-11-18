package conn

import "testing"

const (
	endpoint = "localhost:9000"
)

func TestConnectTogRPC(t *testing.T) {
	_, err := ConnectToSocket(endpoint)
	if err != nil {
		t.Fatal(err)
	}
}
