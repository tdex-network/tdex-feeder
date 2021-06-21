package ports

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/macaroon.v2"

	"github.com/tdex-network/tdex-feeder/internal/domain"
	pboperator "github.com/tdex-network/tdex-protobuf/generated/go/operator"
	"github.com/tdex-network/tdex-protobuf/generated/go/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 200)
)

// TdexDaemonPriceUpdater is a grpc client using to call UpdateMarketPrice RPC of tdex daemon
type TdexDaemonPriceUpdater interface {
	UpdateMarketPrice(ctx context.Context, marketPrice domain.MarketPrice) error
}

// NewTdexDaemonPriceUpdater uses the operatorInterfaceEndpoint to create a gRPC client.
func NewTdexDaemonPriceUpdater(
	operatorInterfaceEndpoint, macaroonsPath, tlsCertPath string,
) (TdexDaemonPriceUpdater, error) {
	if macOK, certOK := macaroonsPath != "", tlsCertPath != ""; macOK != certOK {
		return nil, fmt.Errorf("both macaroons filepath and TLS cert path must be defined")
	}

	connGrpc, err := connectToGRPC(operatorInterfaceEndpoint, macaroonsPath, tlsCertPath)
	if err != nil {
		return nil, err
	}

	// TOSO: Add health check,
	return &tdexDaemonPriceUpdater{
		clientGRPC: pboperator.NewOperatorClient(connGrpc),
	}, nil
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
			QuotePrice: marketPrice.Price.QuotePrice,
		},
	}

	_, err := updater.clientGRPC.UpdateMarketPrice(ctx, &args)
	if err != nil {
		return err
	}

	return nil
}

// ConnectTogRPC dials and returns a new client connection to a remote host
func connectToGRPC(daemonEndpoint, macPath, certPath string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithDefaultCallOptions(maxMsgRecvSize)}

	if len(macPath) <= 0 {
		opts = append(opts, grpc.WithInsecure())
	} else {
		tlsCreds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, fmt.Errorf("could not read TLS certificate:  %s", err)
		}

		macBytes, err := ioutil.ReadFile(macPath)
		if err != nil {
			return nil, fmt.Errorf("could not read macaroon %s: %s", macPath, err)
		}
		mac := &macaroon.Macaroon{}
		err = mac.UnmarshalBinary(macBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse macaroon %s: %s", macPath, err)
		}
		macCreds := NewMacaroonCredential(mac)
		opts = append(opts, grpc.WithPerRPCCredentials(macCreds))
		opts = append(opts, grpc.WithTransportCredentials(tlsCreds))
	}

	conn, err := grpc.Dial(daemonEndpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to RPC server: %v", err)
	}

	log.Println("Connected to gRPC:", daemonEndpoint)
	return conn, nil
}
