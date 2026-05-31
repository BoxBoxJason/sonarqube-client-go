package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupportService_Info(t *testing.T) {
	data := []byte(`{"System":{"OS":"Linux"},"Database":{"Name":"PostgreSQL"}}`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/support/info", http.StatusOK, "application/json", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Support.Info(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}
