package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

func TestSamlService_Validation(t *testing.T) {
	data := []byte(`<html><body>SAML Validation Result</body></html>`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodPost, "/saml/validation", http.StatusOK, "text/html", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Saml.Validation(context.Background(), &SamlValidationOptions{
		SAMLResponse: "base64encodedsamlresponse",
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
