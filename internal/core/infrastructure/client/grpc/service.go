package grpcdaemon

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/tdex-network/tdex-daemon/pkg/macaroons"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"

	pboperator "github.com/tdex-network/tdex-daemon/api-spec/protobuf/gen/operator"
	pbunlocker "github.com/tdex-network/tdex-daemon/api-spec/protobuf/gen/walletunlocker"
	pbtypes "github.com/tdex-network/tdex-protobuf/generated/go/types"
)

var (
	maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 200)
)

type service struct {
	rpcAddress string

	unlockerClient pbunlocker.WalletUnlockerClient
	operatorClient pboperator.OperatorClient
}

func NewGRPCClient(
	addr, macaroonsPath, tlsCertPath string,
) (ports.TdexClient, error) {
	unlockerConn, err := createGRPCConn(addr, macaroonsPath, tlsCertPath)
	if err != nil {
		return nil, err
	}

	operatorConn, err := createGRPCConn(addr, macaroonsPath, tlsCertPath)
	if err != nil {
		return nil, err
	}

	unlockerClient := pbunlocker.NewWalletUnlockerClient(unlockerConn)
	operatorClient := pboperator.NewOperatorClient(operatorConn)

	return &service{
		rpcAddress:     addr,
		unlockerClient: unlockerClient,
		operatorClient: operatorClient,
	}, nil
}

func (s *service) RPCAddress() string {
	return s.rpcAddress
}

func (s *service) IsReady() (bool, error) {
	res, err := s.unlockerClient.IsReady(
		context.Background(), &pbunlocker.IsReadyRequest{},
	)
	if err != nil {
		return false, err
	}
	return res.GetInitialized() && res.GetUnlocked(), nil
}

func (s *service) UpdateMarketPrice(
	mkt ports.Market, price ports.Price,
) error {
	basePrice, _ := strconv.ParseFloat(price.BasePrice(), 32)
	quotePrice, _ := strconv.ParseFloat(price.QuotePrice(), 32)
	_, err := s.operatorClient.UpdateMarketPrice(
		context.Background(), &pboperator.UpdateMarketPriceRequest{
			Market: &pbtypes.Market{
				BaseAsset:  mkt.BaseAsset(),
				QuoteAsset: mkt.QuoteAsset(),
			},
			Price: &pbtypes.Price{
				BasePrice:  float32(basePrice),
				QuotePrice: float32(quotePrice),
			},
		},
	)
	return err
}

func (s *service) ListMarkets() ([]ports.Market, error) {
	res, err := s.operatorClient.ListMarkets(
		context.Background(), &pboperator.ListMarketsRequest{},
	)
	if err != nil {
		return nil, err
	}

	mkts := res.GetMarkets()
	markets := make([]ports.Market, 0, len(mkts))
	for _, mkt := range mkts {
		markets = append(markets, market{
			baseAsset:  mkt.GetMarket().GetBaseAsset(),
			quoteAsset: mkt.GetMarket().GetQuoteAsset(),
		})
	}
	return markets, nil
}

func createGRPCConn(daemonEndpoint, macPath, certPath string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithDefaultCallOptions(maxMsgRecvSize)}

	if len(macPath) <= 0 {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// TLS credentials
		tlsCreds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, fmt.Errorf("could not read TLS certificate:  %s", err)
		}
		// macaroons credentials
		macBytes, err := ioutil.ReadFile(macPath)
		if err != nil {
			return nil, fmt.Errorf("could not read macaroon %s: %s", macPath, err)
		}
		mac := &macaroon.Macaroon{}
		if err = mac.UnmarshalBinary(macBytes); err != nil {
			return nil, fmt.Errorf("could not parse macaroon %s: %s", macPath, err)
		}
		macCreds := macaroons.NewMacaroonCredential(mac, true)

		opts = append(opts, grpc.WithPerRPCCredentials(macCreds))
		opts = append(opts, grpc.WithTransportCredentials(tlsCreds))
	}

	conn, err := grpc.Dial(daemonEndpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to RPC server: %v", err)
	}

	return conn, nil
}
