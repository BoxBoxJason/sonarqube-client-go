package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityAlertsService_Search(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/security-alerts/alerts", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.SecurityAlerts.Search(context.Background(), &SecurityAlertsSearchOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestSecurityAlertsService_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.SecurityAlerts.Search(context.Background(), nil)
	require.Error(t, err)
}

func TestSecurityAlertsService_Get(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/security-alerts/alerts/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.SecurityAlerts.Get(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestSecurityAlertsService_Get_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.SecurityAlerts.Get(context.Background(), "")
	require.Error(t, err)
}
