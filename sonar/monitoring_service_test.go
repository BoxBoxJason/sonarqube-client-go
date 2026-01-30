package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitoring_Metrics(t *testing.T) {
	metricsContent := `# HELP sonarqube_health Health check status
# TYPE sonarqube_health gauge
sonarqube_health 1`

	handler := mockBinaryHandler(t, http.MethodGet, "/monitoring/metrics", http.StatusOK, "text/plain", []byte(metricsContent))
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Monitoring.Metrics()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)
	assert.Contains(t, *result, "sonarqube_health")
}
