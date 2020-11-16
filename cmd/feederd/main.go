package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tdex-network/tdex-feeder/config"
	"github.com/tdex-network/tdex-feeder/pkg/conn"
	"github.com/tdex-network/tdex-feeder/pkg/marketinfo"

	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
)

const (
	defaultConfigPath = "./config/config.json"
)

func main() {
	// Checks for command line flags for Config Path
	confFlag := flag.String("conf", defaultConfigPath, "Configuration File Path")
	flag.Parse()

	// Loads Config File
	conf := config.LoadConfig(*confFlag)

	// Interrupt Notification
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Dials the connection the the Socket
	cSocket := conn.ConnectToSocket(conf.Kraken_ws_endpoint)
	defer cSocket.Close()

	// Set up the connection to the gRPC server.
	conngRPC := conn.ConnectTogRPC(conf.Daemon_endpoint)
	defer conngRPC.Close()
	clientgRPC := pboperator.NewOperatorClient(conngRPC)

	// Loads Config Markets infos into Data Structure and Subscribes to Messages from this Markets
	numberOfMarkets := len(conf.Markets)
	marketsInfos := make([]*marketinfo.MarketInfo, numberOfMarkets)
	for i, marketConfig := range conf.Markets {
		marketsInfos[i] = marketinfo.DefaultMarketInfo(marketConfig)
		defer marketsInfos[i].GetInterval().Stop()
		m := conn.CreateSubscribeToMarketMessage(marketConfig.Kraken_ticker)
		conn.SendRequestMessage(cSocket, m)
	}

	// Gets messages from subscriptions
	done := make(chan string)
	go conn.GetMessages(done, cSocket, marketsInfos)

	// Periodically sends gRPC request to update price
	go conn.UpdateMarketPricegRPC(marketsInfos, clientgRPC)

	// Loop to keep cycle alive. Aways Interrupt to close the connection.
	for {
		select {
		case <-interrupt:
			log.Println("Shutting down Feeder")
			err := cSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
