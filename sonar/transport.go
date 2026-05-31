package sonar

import (
	"crypto/tls"
	"net/http"
	"time"
)

// TransportConfig holds configuration for the SDK-managed HTTP transport.
// It is ignored when WithHTTPClient is also used.
type TransportConfig struct {
	TLSClientConfig     *tls.Config
	MaxIdleConns        int
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
	DisableCompression  bool
}

// buildTransport creates an *http.Transport from cfg.
func buildTransport(cfg TransportConfig) *http.Transport {
	//nolint:exhaustruct // only configure the fields the caller specified
	return &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConns,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		TLSHandshakeTimeout: cfg.TLSHandshakeTimeout,
		TLSClientConfig:     cfg.TLSClientConfig,
		DisableCompression:  cfg.DisableCompression,
	}
}
