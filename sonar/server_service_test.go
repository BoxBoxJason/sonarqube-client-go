package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Version(t *testing.T) {
	handler := mockBinaryHandler(t, http.MethodGet, "/server/version", http.StatusOK, "text/plain", []byte("9.9"))
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	version, resp, err := client.Server.Version()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, version)
	assert.Equal(t, "9.9", *version)
}

func TestServer_Version_ErrorResponse(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/server/version", http.StatusInternalServerError, `{"errors":[{"msg":"internal error"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	version, resp, err := client.Server.Version()
	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Nil(t, version)
}
