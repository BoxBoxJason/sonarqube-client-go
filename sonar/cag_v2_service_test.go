package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCagService_GetEntitlement(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/cag-entitlement/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.GetEntitlement(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_GetEntitlement_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Cag.GetEntitlement(context.Background(), "")
	require.Error(t, err)
}

func TestCagService_GetUsageStats(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/cag-usage-stats/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.GetUsageStats(context.Background(), "test-id", &CagGetUsageStatsOptions{Period: "LAST_DAY"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_GetUsageStats_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Cag.GetUsageStats(context.Background(), "", nil)
	require.Error(t, err)
}

func TestCagService_GetGuideMetrics(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/impact/guide-metrics", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.GetGuideMetrics(context.Background(), &CagGetGuideMetricsOptions{OrganizationId: "test", From: "test", To: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_GetGuideMetrics_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Cag.GetGuideMetrics(context.Background(), nil)
	require.Error(t, err)
}

func TestCagService_GetProjectActivity(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/impact/project-activity", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.GetProjectActivity(context.Background(), &CagGetProjectActivityOptions{OrganizationId: "test", From: "test", To: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_GetProjectActivity_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Cag.GetProjectActivity(context.Background(), nil)
	require.Error(t, err)
}

func TestCagService_GetProjectImpact(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/impact/projects", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.GetProjectImpact(context.Background(), &CagGetProjectImpactOptions{OrganizationId: "test", From: "test", To: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_GetProjectImpact_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Cag.GetProjectImpact(context.Background(), nil)
	require.Error(t, err)
}

func TestCagService_GetVerifyMetrics(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/impact/verify-metrics", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.GetVerifyMetrics(context.Background(), &CagGetVerifyMetricsOptions{OrganizationId: "test", From: "test", To: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_GetVerifyMetrics_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Cag.GetVerifyMetrics(context.Background(), nil)
	require.Error(t, err)
}

func TestCagService_Ping(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/cag/ping", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Cag.Ping(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestCagService_RecordImpactEvent(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/v2/cag/cag-impact-events", http.StatusAccepted))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Cag.RecordImpactEvent(context.Background(), &CagRecordImpactEventOptions{OrganizationId: "test", ProjectId: "test", InvocationId: "test", CagInstanceId: "test", EventType: "NAVIGATION", EventVersion: 1, Success: true, Transport: "MCP", Payload: nil})
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestCagService_RecordImpactEvent_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Cag.RecordImpactEvent(context.Background(), nil)
	require.Error(t, err)
}

func TestCagService_RecordUsageEvent(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/v2/cag/cag-usage-events", http.StatusAccepted))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Cag.RecordUsageEvent(context.Background(), &CagRecordUsageEventOptions{OrganizationId: "test", ProjectId: "test", InvocationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestCagService_RecordUsageEvent_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Cag.RecordUsageEvent(context.Background(), nil)
	require.Error(t, err)
}
