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
	"github.com/tdex-network/tdex-feeder/pkg/marketinfos"

	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
)

const (
	defaultConfigPath = "./config/config.json"
)

func main() {

	confFlag := flag.String("conf", defaultConfigPath, "Configuration File Path")
	flag.Parse()

	conf := config.LoadConfig(*confFlag)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	cSocket := conn.ConnectToSocket(conf.Kraken_ws_endpoint)
	defer cSocket.Close()

	// Set up a connection to the gRPC server.
	conngRPC := conn.ConnectTogRPC(conf.Daemon_endpoint)
	defer conngRPC.Close()
	clientgRPC := pboperator.NewOperatorClient(conngRPC)

	numberOfMarkets := len(conf.Markets)
	marketsInfos := make([]*marketinfos.MarketInfo, numberOfMarkets)
	for i, marketConfig := range conf.Markets {
		marketsInfos[i] = marketinfos.DefaultMarketInfo(marketConfig)
		defer marketsInfos[i].GetInterval().Stop()
		m := conn.CreateSubscribeToMarketMessage(marketConfig.Kraken_ticker)
		conn.SendRequestMessage(cSocket, m)
	}

	done := make(chan string)
	go conn.GetMessages(done, cSocket, marketsInfos)
	go conn.UpdateMarketPricegRPC(marketsInfos, clientgRPC)

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
