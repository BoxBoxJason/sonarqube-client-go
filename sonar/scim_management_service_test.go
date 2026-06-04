package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Disable
// -----------------------------------------------------------------------------

func TestScimManagementService_Disable(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/scim_management/disable", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.ScimManagement.Disable(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Enable
// -----------------------------------------------------------------------------

func TestScimManagementService_Enable(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/scim_management/enable", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.ScimManagement.Enable(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Status
// -----------------------------------------------------------------------------

func TestScimManagementService_Status(t *testing.T) {
	response := ScimManagementStatus{
		Enabled: true,
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/scim_management/status", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.ScimManagement.Status(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.True(t, result.Enabled)
}
