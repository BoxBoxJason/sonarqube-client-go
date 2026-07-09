package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupportService_Info(t *testing.T) {
	response := &SupportInfo{
		System:   map[string]any{"OS": "Linux"},
		Database: map[string]any{"Name": "PostgreSQL"},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/support/info", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Support.Info(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, response, result)
}

func TestSupportService_Info_ErrorResponse(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/support/info", http.StatusForbidden, `{"errors":[{"msg":"Insufficient privileges"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.URL)

	result, resp, err := client.Support.Info(context.Background())
	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Nil(t, result)
}

// TestSupportService_Info_LicenseNotFound mirrors the response observed when
// live-verifying this endpoint against a real, unlicensed SonarQube 2025.2
// Enterprise instance: GET /api/support/info returns HTTP 400 with
// {"errors":[{"msg":"License not found"}]} even with full admin permissions,
// because the endpoint requires an actually installed commercial license to
// succeed (see .github/reviews/pr-265.md).
func TestSupportService_Info_LicenseNotFound(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/support/info", http.StatusBadRequest, `{"errors":[{"msg":"License not found"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.URL)

	result, resp, err := client.Support.Info(context.Background())
	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Nil(t, result)
}
