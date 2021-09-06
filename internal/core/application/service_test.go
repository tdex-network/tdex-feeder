package application_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tdex-network/tdex-feeder/internal/core/application"
	krakenfeeder "github.com/tdex-network/tdex-feeder/internal/core/infrastructure/feeder/kraken"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

func TestService(t *testing.T) {
	appSvc, err := newTestService()
	require.NoError(t, err)

	go func() {
		time.Sleep(10 * time.Second)
		appSvc.Stop()
	}()

	err = appSvc.Start()
	require.NoError(t, err)
}

func newTestService() (application.Service, error) {
	interval := 2000 // 2s interval
	tickers := []string{"XBT/USDT", "XBT/EUR"}
	markets := mockedMarkets(tickers)

	priceFeeder, err := krakenfeeder.NewKrakenPriceFeeder(interval, markets)
	if err != nil {
		return nil, err
	}

	targetsByMarket := make(application.IndexedTargetsByMarket)
	for _, mkt := range markets {
		mktKey := ports.MarketKey(mkt)
		targets := make(map[string]ports.TdexClient)
		daemons := mockedDaemons(2)
		for _, d := range daemons {
			targets[d.RPCAddress()] = d
		}
		targetsByMarket[mktKey] = targets
	}

	return application.NewService(priceFeeder, targetsByMarket), nil
}

func mockedMarkets(tickers []string) []ports.Market {
	markets := make([]ports.Market, 0, len(tickers))
	for _, ticker := range tickers {
		markets = append(markets, &mockMarket{
			baseAsset:  randomHex(32),
			quoteAsset: randomHex(32),
			ticker:     ticker,
		})
	}
	return markets
}

func mockedDaemons(num int) []ports.TdexClient {
	daemons := make([]ports.TdexClient, 0, num)
	for i := 0; i < num; i++ {
		mockedDaemon := &mockDaemon{}
		mockedDaemon.On("RPCAddress").Return(randomAddr())
		mockedDaemon.On("IsReady").Return(true, nil)
		mockedDaemon.On("UpdateMarketPrice", mock.Anything, mock.Anything).Return(nil)
		daemons = append(daemons, mockedDaemon)
	}
	return daemons
}

type mockMarket struct {
	baseAsset  string
	quoteAsset string
	ticker     string
}

func (m *mockMarket) BaseAsset() string {
	return m.baseAsset
}

func (m *mockMarket) QuoteAsset() string {
	return m.quoteAsset
}

func (m *mockMarket) Ticker() string {
	return m.ticker
}

type mockDaemon struct {
	mock.Mock
}

func (m *mockDaemon) RPCAddress() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m mockDaemon) IsReady() (bool, error) {
	args := m.Called()

	var res bool
	if a := args.Get(0); a != nil {
		res = a.(bool)
	}
	return res, args.Error(1)
}

func (m mockDaemon) UpdateMarketPrice(mkt ports.Market, prc ports.Price) error {
	args := m.Called(mkt, prc)
	return args.Error(0)
}

func randomAddr() string {
	return fmt.Sprintf("localhost:%d", randomIntInRange(1024, 9999))
}

func randomHex(len int) string {
	return hex.EncodeToString(randomBytes(len))
}

func randomBytes(len int) []byte {
	b := make([]byte, len)
	rand.Read(b)
	return b
}

func randomIntInRange(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(int(n.Int64())) + min
}
