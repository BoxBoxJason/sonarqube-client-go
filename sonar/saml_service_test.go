package sonar

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

func TestSamlService_Validation(t *testing.T) {
	data := []byte(`<html><body>SAML Validation Result</body></html>`)
	// Real SAML assertions are base64-encoded blobs typically several KB long
	// once base64-encoded. Use an oversized value here to prove it survives
	// the round trip: were it (incorrectly) sent as a URL query parameter, a
	// value this large would be rejected or truncated by most HTTP servers.
	samlResponse := strings.Repeat("a", 10*1024)

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "unexpected HTTP method")
		assert.Equal(t, "/saml/validation", r.URL.Path, "unexpected URL path")
		assert.Empty(t, r.URL.RawQuery, "SAMLResponse must not be sent as a URL query parameter")
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"), "unexpected Content-Type")

		require.NoError(t, r.ParseForm())
		assert.Equal(t, samlResponse, r.FormValue("SAMLResponse"), "SAMLResponse must be transmitted in the request body")

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})
	client := newTestClient(t, server.URL)

	result, resp, err := client.Saml.Validation(context.Background(), &SamlValidationOptions{
		SAMLResponse: samlResponse,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestSamlService_Validation_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.Saml.Validation(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.Saml.Validation(context.Background(), &SamlValidationOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// ValidationInit
// -----------------------------------------------------------------------------

func TestSamlService_ValidationInit(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodGet, "/saml/validation_init", http.StatusOK))
	client := newTestClient(t, server.URL)

	resp, err := client.Saml.ValidationInit(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
