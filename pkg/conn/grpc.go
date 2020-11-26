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

const (
	timeout = 3
)

// ConnectTogRPC dials and returns a new client connection to a remote host
func ConnectTogRPC(daemonEndpoint string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, daemonEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return conn, err
	}
	log.Println("Connected to gRPC:", daemonEndpoint)
	return conn, nil
}

// UpdateMarketPricegRPC calls the tdex daemon UpdateMarketPrice rpc endpoint to update a defined market
func UpdateMarketPricegRPC(marketInfo marketinfo.MarketInfo, clientgRPC pboperator.OperatorClient) {

	if marketInfo.Price == 0.00 {
		log.Println("Can't send gRPC request with no price")
		return
	}

	log.Printf("%s %g for market %s-%s", marketInfo.Config.KrakenTicker, marketInfo.Price, marketInfo.Config.BaseAsset[:4], marketInfo.Config.QuoteAsset[:4])

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := clientgRPC.UpdateMarketPrice(ctx, &pboperator.UpdateMarketPriceRequest{
		Market: &pbtypes.Market{
			BaseAsset:  marketInfo.Config.BaseAsset,
			QuoteAsset: marketInfo.Config.QuoteAsset,
		},
		Price: &pbtypes.Price{
			BasePrice:  1 / float32(marketInfo.Price),
			QuotePrice: float32(marketInfo.Price),
		},
	})

	if err != nil {
		log.Println(err)
		return
	}
}
