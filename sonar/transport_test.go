package sonar

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTransportConfig_SetsCustomTransport(t *testing.T) {
	cfg := TransportConfig{
		MaxIdleConns:        42,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  true,
	}

	client, err := NewClient(nil, WithTransportConfig(cfg))
	require.NoError(t, err)
	require.NotNil(t, client)

	assert.NotSame(t, http.DefaultClient, client.httpClient, "expected a dedicated http.Client, not the default")
	require.NotNil(t, client.httpClient.Transport, "expected a non-nil transport")

	transport, ok := client.httpClient.Transport.(*http.Transport)
	require.True(t, ok, "expected transport to be *http.Transport")

	assert.Equal(t, 42, transport.MaxIdleConns)
	assert.Equal(t, 30*time.Second, transport.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, transport.TLSHandshakeTimeout)
	assert.True(t, transport.DisableCompression)
}

func TestWithTransportConfig_IgnoredWhenHTTPClientProvided(t *testing.T) {
	customTransport := &http.Transport{MaxIdleConns: 99}     //nolint:exhaustruct
	customClient := &http.Client{Transport: customTransport} //nolint:exhaustruct

	cfg := TransportConfig{
		MaxIdleConns: 1,
	}

	client, err := NewClient(nil, WithHTTPClient(customClient), WithTransportConfig(cfg))
	require.NoError(t, err)
	require.NotNil(t, client)

	assert.Same(t, customClient, client.httpClient, "expected the provided http.Client to be used")

	transport, ok := client.httpClient.Transport.(*http.Transport)
	require.True(t, ok, "expected transport to be *http.Transport")
	assert.Equal(t, 99, transport.MaxIdleConns, "expected the custom client's transport to be preserved")
}

func TestWithTransportConfig_DefaultUnchanged(t *testing.T) {
	client, err := NewClient(nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// The default client is a dedicated *http.Client (not the shared
	// http.DefaultClient) with the default timeout and no explicit transport
	// (it falls back to http.DefaultTransport at request time).
	assert.NotSame(t, http.DefaultClient, client.httpClient, "expected a dedicated http.Client, not the shared default")
	assert.Equal(t, defaultHTTPTimeout, client.httpClient.Timeout, "expected the default timeout to be applied")
	assert.Nil(t, client.httpClient.Transport, "expected no explicit transport when no options are provided")
}

func TestWithTransportConfig_TLSConfig(t *testing.T) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec // intentional for testing
		MinVersion:         tls.VersionTLS12,
	}

	cfg := TransportConfig{
		TLSClientConfig: tlsCfg,
	}

	client, err := NewClient(nil, WithTransportConfig(cfg))
	require.NoError(t, err)
	require.NotNil(t, client)

	require.NotNil(t, client.httpClient.Transport, "expected a non-nil transport")

	transport, ok := client.httpClient.Transport.(*http.Transport)
	require.True(t, ok, "expected transport to be *http.Transport")

	require.NotNil(t, transport.TLSClientConfig, "expected TLSClientConfig to be set")
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify) //nolint:gosec
	assert.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
}
