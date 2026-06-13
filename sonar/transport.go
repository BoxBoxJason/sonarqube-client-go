package sonar

import (
	"crypto/tls"
	"net/http"
	"time"
)

// TransportConfig holds configuration for the SDK-managed HTTP transport.
// It is ignored when WithHTTPClient is also used.
type TransportConfig struct {
	// TLSClientConfig customizes the TLS settings used for HTTPS connections.
	//
	// Security warning: setting TLSClientConfig.InsecureSkipVerify to true
	// disables server certificate verification, which exposes the connection to
	// man-in-the-middle attacks and can leak the auth token and source-code
	// metadata exchanged with SonarQube. Only use it against a trusted local or
	// development instance, never in production.
	TLSClientConfig *tls.Config
	// MaxIdleConns controls the maximum number of idle (keep-alive) connections
	// across all hosts. Zero means no limit.
	MaxIdleConns int
	// IdleConnTimeout is the maximum time an idle connection is kept before
	// closing. Zero means no limit.
	IdleConnTimeout time.Duration
	// TLSHandshakeTimeout is the maximum time to wait for a TLS handshake.
	// Zero means no timeout.
	TLSHandshakeTimeout time.Duration
	// DisableCompression disables transparent gzip request/response compression.
	DisableCompression bool
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
