package grpcdaemon

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/tdex-network/tdex-daemon/pkg/macaroons"
	"github.com/tdex-network/tdex-daemon/pkg/tdexdconnect"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"

	pb "github.com/tdex-network/tdex-daemon/api-spec/protobuf/gen/tdex-daemon/v1"
	pbtypes "github.com/tdex-network/tdex-daemon/api-spec/protobuf/gen/tdex/v1"
)

var (
	maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 200)
)

type service struct {
	rpcAddress string

	unlockerClient pb.WalletUnlockerServiceClient
	operatorClient pb.OperatorServiceClient
}

func NewGRPCClient(
	addr, macaroonsPath, tlsCertPath string,
) (ports.TdexClient, error) {
	unlockerConn, err := createGRPCConnFromFile(addr, macaroonsPath, tlsCertPath)
	if err != nil {
		return nil, err
	}

	operatorConn, err := createGRPCConnFromFile(addr, macaroonsPath, tlsCertPath)
	if err != nil {
		return nil, err
	}

	unlockerClient := pb.NewWalletUnlockerServiceClient(unlockerConn)
	operatorClient := pb.NewOperatorServiceClient(operatorConn)

	return &service{
		rpcAddress:     addr,
		unlockerClient: unlockerClient,
		operatorClient: operatorClient,
	}, nil
}

func NewGRPCClientFromURL(url string) (ports.TdexClient, error) {
	addr, tlsCert, macaroon, err := tdexdconnect.Decode(url)
	if err != nil {
		return nil, err
	}

	unlockerConn, err := createGRPCConn(addr, macaroon, tlsCert)
	if err != nil {
		return nil, err
	}

	operatorConn, err := createGRPCConn(addr, macaroon, tlsCert)
	if err != nil {
		return nil, err
	}

	unlockerClient := pb.NewWalletUnlockerServiceClient(unlockerConn)
	operatorClient := pb.NewOperatorServiceClient(operatorConn)

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
		context.Background(), &pb.IsReadyRequest{},
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
		context.Background(), &pb.UpdateMarketPriceRequest{
			Market: &pbtypes.Market{
				BaseAsset:  mkt.BaseAsset(),
				QuoteAsset: mkt.QuoteAsset(),
			},
			Price: &pbtypes.Price{
				BasePrice:  basePrice,
				QuotePrice: quotePrice,
			},
		},
	)
	return err
}

func (s *service) ListMarkets() ([]ports.Market, error) {
	res, err := s.operatorClient.ListMarkets(
		context.Background(), &pb.ListMarketsRequest{},
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

func createGRPCConn(
	daemonEndpoint string, macBytes, certBytes []byte,
) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithDefaultCallOptions(maxMsgRecvSize)}

	if len(macBytes) <= 0 {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// TLS credentials
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse TLS certificate: %s", err)
		}
		cp := x509.NewCertPool()
		cp.AddCert(cert)
		tlsCreds := credentials.NewClientTLSFromCert(cp, "")

		// macaroons credentials
		mac := &macaroon.Macaroon{}
		if err := mac.UnmarshalBinary(macBytes); err != nil {
			return nil, fmt.Errorf("could not parse macaroon: %s", err)
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

func createGRPCConnFromFile(
	daemonEndpoint, macPath, certPath string,
) (*grpc.ClientConn, error) {
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
