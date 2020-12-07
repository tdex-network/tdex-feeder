package ports

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/tdex-network/tdex-feeder/internal/domain"
	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
	"github.com/tdex-network/tdex-protobuf/generated/go/types"
	"google.golang.org/grpc"
)

type TdexDaemonPriceUpdater interface {
	UpdateMarketPrice(ctx context.Context, marketPrice domain.MarketPrice) error
}

// NewTdexDaemonPriceUpdater uses the operatorInterfaceEndpoint to create a gRPC client.
func NewTdexDaemonPriceUpdater(ctx context.Context, operatorInterfaceEndpoint string) TdexDaemonPriceUpdater {
	connGrpc, err := connectToGRPC(ctx, operatorInterfaceEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	operatorClient := pboperator.NewOperatorClient(connGrpc)

	return &tdexDaemonPriceUpdater{
		clientGRPC: operatorClient,
	}
}

type tdexDaemonPriceUpdater struct {
	clientGRPC pboperator.OperatorClient
}

// UpdateMarketPrice gets a marketPrice and sends updateMarketPrice request through gRPC client.
func (updater *tdexDaemonPriceUpdater) UpdateMarketPrice(ctx context.Context, marketPrice domain.MarketPrice) error {
	if marketPrice.Price.BasePrice == 0.00 {
		return errors.New("Base price is 0.00")
	}

	if marketPrice.Price.BasePrice == 0.00 {
		return errors.New("Quote price is 0.00")
	}

	args := pboperator.UpdateMarketPriceRequest{
		Market: &types.Market{
			BaseAsset:  marketPrice.Market.BaseAsset,
			QuoteAsset: marketPrice.Market.QuoteAsset,
		},
		Price: &types.Price{
			BasePrice:  marketPrice.Price.BasePrice,
			QuotePrice: marketPrice.Price.BasePrice,
		},
	}

	_, err := updater.clientGRPC.UpdateMarketPrice(ctx, &args)
	if err != nil {
		return err
	}

	return nil
}

// ConnectTogRPC dials and returns a new client connection to a remote host
func connectToGRPC(ctx context.Context, daemonEndpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, daemonEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return conn, err
	}
	log.Println("Connected to gRPC:", daemonEndpoint)
	return conn, nil
}
