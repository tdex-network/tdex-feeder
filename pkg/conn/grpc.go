package conn

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tdex-network/tdex-feeder/pkg/marketinfo"
	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
	pbtypes "github.com/tdex-network/tdex-protobuf/generated/go/types"
	"google.golang.org/grpc"
)

// ConnectTogRPC dials and returns a new client connection to a remote host
func ConnectTogRPC(daemon_endpoint string) *grpc.ClientConn {
	conn, err := grpc.Dial(daemon_endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to gRPC: %v", err)
	}
	return conn
}

func UpdateMarketPricegRPC(marketsInfos []*marketinfo.MarketInfo, clientgRPC pboperator.OperatorClient) {
	for {
		for _, marketsInfo := range marketsInfos {
			select {
			case <-marketsInfo.GetInterval().C:
				if marketsInfo.GetPrice() == 0.00 {
					log.Println("Can't send gRPC request with no price")
				} else {
					log.Println("Sending gRPC request:", marketsInfo.GetConfig().Kraken_ticker, marketsInfo.GetPrice())
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					r, err := clientgRPC.UpdateMarketPrice(ctx, &pboperator.UpdateMarketPriceRequest{
						Market: &pbtypes.Market{BaseAsset: marketsInfo.GetConfig().Base_asset, QuoteAsset: marketsInfo.GetConfig().Quote_asset},
						Price:  &pbtypes.Price{BasePrice: 1 / float32(marketsInfo.GetPrice()), QuotePrice: float32(marketsInfo.GetPrice())}})
					if err != nil {
						log.Println(err)
					} else {
						log.Println(r)
					}
				}

			}
		}
	}
}
