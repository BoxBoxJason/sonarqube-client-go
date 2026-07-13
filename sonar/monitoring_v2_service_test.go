package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// GetActiveAlerts
// -----------------------------------------------------------------------------

func TestMonitoringV2_GetActiveAlerts(t *testing.T) {
	response := MonitoringActiveAlerts{
		Alerts: []MonitoringAlert{
			{
				Key:         "EXAMPLE_COUNTER_REACHED_FIVE",
				Message:     "The example counter reached five",
				ActiveSince: "2024-01-01T12:00:00Z",
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/monitoring/alerts", http.StatusOK, response))
	client := newTestClient(t, server.URL)
	svc := &MonitoringServiceV2{client: client}

	result, resp, err := svc.GetActiveAlerts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Alerts, 1)
	assert.Equal(t, "EXAMPLE_COUNTER_REACHED_FIVE", result.Alerts[0].Key)
	assert.Equal(t, "The example counter reached five", result.Alerts[0].Message)
	assert.Equal(t, "2024-01-01T12:00:00Z", result.Alerts[0].ActiveSince)
}

func TestMonitoringV2_GetActiveAlerts_Empty(t *testing.T) {
	response := MonitoringActiveAlerts{}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/monitoring/alerts", http.StatusOK, response))
	client := newTestClient(t, server.URL)
	svc := &MonitoringServiceV2{client: client}

	result, resp, err := svc.GetActiveAlerts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, result.Alerts)
}
