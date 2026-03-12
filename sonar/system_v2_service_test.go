package sonar

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// GetMigrationsStatus
// =============================================================================

func TestSystemV2_GetMigrationsStatus(t *testing.T) {
	response := SystemDbMigrationsStatusV2{
		Status:         "MIGRATION_RUNNING",
		CompletedSteps: 5,
		TotalSteps:     10,
		StartedAt:      "2024-01-01T00:00:00+0000",
		Message:        "Migration in progress",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/system/migrations-status", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.System.GetMigrationsStatus()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "MIGRATION_RUNNING", result.Status)
	assert.Equal(t, int32(5), result.CompletedSteps)
	assert.Equal(t, int32(10), result.TotalSteps)
}

// =============================================================================
// CheckLiveness
// =============================================================================

func TestSystemV2_CheckLiveness(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodGet, "/v2/system/liveness", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.System.CheckLiveness(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestSystemV2_CheckLiveness_WithPasscode(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2/system/liveness", r.URL.Path)
		assert.Equal(t, "my-passcode", r.Header.Get("X-Sonar-Passcode"))
		w.WriteHeader(http.StatusNoContent)
	})
	client := newTestClient(t, server.url())

	resp, err := client.V2.System.CheckLiveness(&SystemPasscodeOptionV2{Passcode: "my-passcode"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// =============================================================================
// GetHealth
// =============================================================================

func TestSystemV2_GetHealth(t *testing.T) {
	response := SystemHealthV2{
		Status: "GREEN",
		Causes: []string{},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/system/health", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.System.GetHealth(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GREEN", result.Status)
}

func TestSystemV2_GetHealth_WithPasscode(t *testing.T) {
	response := SystemHealthV2{
		Status: "YELLOW",
		Causes: []string{"Elasticsearch is YELLOW"},
	}
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2/system/health", r.URL.Path)
		assert.Equal(t, "secret-passcode", r.Header.Get("X-Sonar-Passcode"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	})
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.System.GetHealth(&SystemPasscodeOptionV2{Passcode: "secret-passcode"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "YELLOW", result.Status)
	assert.Len(t, result.Causes, 1)
}
