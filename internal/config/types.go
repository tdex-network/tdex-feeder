package config

import (
	"fmt"
)

type Target struct {
	RPCAddress    string `mapstructure:"rpc_address"`
	TLSCertPath   string `mapstructure:"tls_cert_path,omitempty"`
	MacaroonsPath string `mapstructure:"macaroons_path,omitempty"`
}

func (t Target) validate() error {
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
	CBaseAsset  string `mapstructure:"base_asset"`
	CQuoteAsset string `mapstructure:"quote_asset"`
	CTicker     string `mapstructure:"ticker"`
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

func (m Market) RawMap() map[string]string {
	return map[string]string{
		"base_asset":  m.CBaseAsset,
		"quote_asset": m.CQuoteAsset,
		"ticker":      m.CTicker,
	}
}
