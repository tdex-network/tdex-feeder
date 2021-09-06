module github.com/tdex-network/tdex-feeder

go 1.15

require (
	github.com/gorilla/websocket v1.4.2
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tdex-network/tdex-daemon v0.5.4-0.20210901130622-99c797e9450a
	github.com/tdex-network/tdex-daemon/pkg/macaroons v0.0.0-20210813140257-70d50a8b72a4
	github.com/tdex-network/tdex-protobuf v0.0.0-20210507104156-d509331cccdb
	google.golang.org/grpc v1.40.0
	gopkg.in/macaroon.v2 v2.1.0
)
