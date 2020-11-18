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
func ConnectTogRPC(daemon_endpoint string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, daemon_endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func UpdateMarketPricegRPC(marketsInfos []*marketinfo.MarketInfo, clientgRPC pboperator.OperatorClient) {
	for {
		for _, marketInfo := range marketsInfos {
			select {
			case <-marketInfo.GetInterval().C:
				if marketInfo.GetPrice() == 0.00 {
					log.Println("Can't send gRPC request with no price")
				}
				if marketInfo.GetPrice() != 0.00 {
					log.Println("Sending gRPC request:", marketInfo.GetConfig().Kraken_ticker, marketInfo.GetPrice())
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					r, err := clientgRPC.UpdateMarketPrice(ctx, &pboperator.UpdateMarketPriceRequest{
						Market: &pbtypes.Market{BaseAsset: marketInfo.GetConfig().Base_asset, QuoteAsset: marketInfo.GetConfig().Quote_asset},
						Price:  &pbtypes.Price{BasePrice: 1 / float32(marketInfo.GetPrice()), QuotePrice: float32(marketInfo.GetPrice())}})
					if err != nil {
						log.Println(err)
					}
					if err == nil {
						log.Println(r)
					}
				}
			}
		}
	}
}
