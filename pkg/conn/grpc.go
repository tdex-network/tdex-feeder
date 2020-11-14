package conn

import (
	"context"
	"log"
	"time"

	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
	pbtypes "github.com/tdex-network/tdex-protobuf/generated/go/types"
	"google.golang.org/grpc"
)

// ConnectTogRPC dials and returns a new client connection to a remote host
func ConnectTogRPC(daemon_endpoint string) *grpc.ClientConn {
	conn, err := grpc.Dial(daemon_endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func UpdateMarketPricegRPC(marketsInfos marketsInformations, clientgRPC pboperator.OperatorClient) {
	for _, marketsInfo := range marketsInfos {
		select {
		case <-marketsInfo.interval.C:
			if marketsInfo.price == 0.00 {
				log.Println("Can't send gRPC request with no price")
			} else {
				log.Println("Sending gRPC request:", marketsInfo.config.Kraken_ticker, marketsInfo.price)
				// Contact the server and print out its response.
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := clientgRPC.UpdateMarketPrice(ctx, &pboperator.UpdateMarketPriceRequest{
					Market: &pbtypes.Market{BaseAsset: marketsInfo.config.Base_asset, QuoteAsset: marketsInfo.config.Quote_asset},
					Price:  &pbtypes.Price{BasePrice: 1 / float32(marketsInfo.price), QuotePrice: float32(marketsInfo.price)}})
				if err != nil {
					log.Println(err)
				}
				log.Println(r)
			}

		}
	}
}
