package config

import (
	"encoding/hex"
	"fmt"

	"github.com/tdex-network/tdex-daemon/pkg/tdexdconnect"
	"github.com/tdex-network/tdex-feeder/internal/core/ports"
)

type Target struct {
	RPCAddress      string `mapstructure:"rpc_address"`
	TLSCertPath     string `mapstructure:"tls_cert_path,omitempty"`
	MacaroonsPath   string `mapstructure:"macaroons_path,omitempty"`
	TdexdconnectURL string `mapstructure:"tdexdconnect_url,omitempty"`
}

func (t Target) validate() error {
	if url := t.TdexdconnectURL; url != "" {
		_, _, _, err := tdexdconnect.Decode(url)
		return err
	}

	if t.RPCAddress == "" {
		return fmt.Errorf("target rpc address must not be nil")
	}
	macPathOk := t.MacaroonsPath != ""
	certPathOk := t.TLSCertPath != ""
	if macPathOk != certPathOk {
		return fmt.Errorf(
			"macaroon and tls cert paths must be both either set or unset for a target",
		)
	}
	return nil
}

type Market struct {
	CBaseAsset  string   `mapstructure:"base_asset"`
	CQuoteAsset string   `mapstructure:"quote_asset"`
	CTicker     string   `mapstructure:"ticker"`
	CTargets    []Target `mapstructure:"targets"`
}

func (m Market) BaseAsset() string {
	return m.CBaseAsset
}

func (m Market) QuoteAsset() string {
	return m.CQuoteAsset
}

func (m Market) Ticker() string {
	return m.CTicker
}

func (m Market) validate() error {
	if m.BaseAsset() == "" {
		return fmt.Errorf("market base asset must not be nil")
	}
	ba, err := hex.DecodeString(m.BaseAsset())
	if err != nil {
		return fmt.Errorf("market base asset must be an hex string")
	}
	if len(ba) != 32 {
		return fmt.Errorf("market base asset must be a 64-chars hex string")
	}
	if m.QuoteAsset() == "" {
		return fmt.Errorf("market quote asset must not be nil")
	}
	qa, err := hex.DecodeString(m.QuoteAsset())
	if err != nil {
		return fmt.Errorf("market quote asset must be an hex string")
	}
	if len(qa) != 32 {
		return fmt.Errorf("market quote asset must be a 64-chars hex string")
	}
	if m.Ticker() == "" {
		return fmt.Errorf("market ticker must not be nil")
	}
	if len(m.CTargets) <= 0 {
		return fmt.Errorf("market must have at least one target")
	}
	for _, t := range m.CTargets {
		if err := t.validate(); err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	PriceFeeder string   `mapstructure:"price_feeder"`
	Interval    int      `mapstructure:"interval"`
	Markets     []Market `mapstructure:"markets"`
}

func (c Config) Validate() error {
	if c.PriceFeeder == "" {
		return fmt.Errorf("price_feeder must not be nil")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be a positive value")
	}
	if len(c.Markets) <= 0 {
		return fmt.Errorf("markets must not be empty")
	}
	for _, mkt := range c.Markets {
		if err := mkt.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c Config) PortableMarkets() []ports.Market {
	markets := make([]ports.Market, 0, len(c.Markets))
	for _, mkt := range c.Markets {
		markets = append(markets, mkt)
	}
	return markets
}
