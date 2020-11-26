// Copyright (c) 2020 The VulpemVentures developers

// Feeder allows to connect an external price feed to the TDex Daemon to determine the current market price.
package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/gorilla/websocket"
	"github.com/tdex-network/tdex-feeder/config"
	"github.com/tdex-network/tdex-feeder/pkg/conn"
	"github.com/tdex-network/tdex-feeder/pkg/marketinfo"

	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
)

const (
	defaultConfigPath = "./config.json"
)

func main() {
	interrupt, cSocket, marketsInfos, conngRPC := setup()
	infiniteLoops(interrupt, cSocket, marketsInfos, conngRPC)
}

func setup() (chan os.Signal, *websocket.Conn, []marketinfo.MarketInfo, *grpc.ClientConn) {
	conf := checkFlags()

	// Interrupt Notification.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Dials the connection the the Socket.
	cSocket, err := conn.ConnectToSocket(conf.KrakenWsEndpoint)
	if err != nil {
		log.Fatal("Socket Connection Error: ", err)
	}
	marketsInfos := loadMarkets(conf, cSocket)
	if len(marketsInfos) == 0 {
		log.Warn("list of market to feed is empty")
	}

	// Set up the connection to the gRPC server.
	conngRPC, err := conn.ConnectTogRPC(conf.DaemonEndpoint)
	if err != nil {
		log.Fatal("gRPC Connection Error: ", err)
	}
	return interrupt, cSocket, marketsInfos, conngRPC
}

// Checks for command line flags for Config Path and Debug mode.
// Loads flags as required.
func checkFlags() config.Config {
	confFlag := flag.String("conf", defaultConfigPath, "Configuration File Path")
	debugFlag := flag.Bool("debug", false, "Log Debug Informations")
	flag.Parse()
	if *debugFlag == true {
		log.SetLevel(log.DebugLevel)
	}
	// Loads Config File.
	conf, err := config.LoadConfig(*confFlag)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

// Loads Config Markets infos into Data Structure and Subscribes to
// Messages from this Markets.
func loadMarkets(conf config.Config, cSocket *websocket.Conn) []marketinfo.MarketInfo {
	numberOfMarkets := len(conf.Markets)
	marketsInfos := make([]marketinfo.MarketInfo, numberOfMarkets)
	for i, marketConfig := range conf.Markets {
		marketsInfos[i] = marketinfo.InitialMarketInfo(marketConfig)
		m := conn.CreateSubscribeToMarketMessage(marketConfig.KrakenTicker)
		err := conn.SendRequestMessage(cSocket, m)
		if err != nil {
			log.Fatal("Couldn't send request message: ", err)
		}
	}
	return marketsInfos
}

func infiniteLoops(interrupt chan os.Signal, cSocket *websocket.Conn, marketsInfos []marketinfo.MarketInfo, conngRPC *grpc.ClientConn) {
	defer cSocket.Close()
	defer conngRPC.Close()
	clientgRPC := pboperator.NewOperatorClient(conngRPC)
	done := make(chan string)
	// Handles Messages from subscriptions. Will periodically call the
	// gRPC UpdateMarketPrice with the price info from the messages.
	go conn.HandleMessages(done, cSocket, marketsInfos, clientgRPC)
	checkInterrupt(interrupt, cSocket, done)
}

// Loop to keep cycle alive. Waits Interrupt to close the connection.
func checkInterrupt(interrupt chan os.Signal, cSocket *websocket.Conn, done chan string) {
	for {
		for range interrupt {
			log.Println("Shutting down Feeder")
			err := cSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Fatal("write close:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
