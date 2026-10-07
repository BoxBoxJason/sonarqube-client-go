package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryService_GetIssueDensityHistory(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/history/issue-density-history", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetIssueDensityHistory(context.Background(), &HistoryGetIssueDensityHistoryOptions{EntityId: "test", EntityType: "PORTFOLIO", StartDate: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestHistoryService_GetIssueDensityHistory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.History.GetIssueDensityHistory(context.Background(), nil)
	require.Error(t, err)
}

func TestHistoryService_GetIssueResolutionHistory(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/history/issue-resolution-history", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetIssueResolutionHistory(context.Background(), &HistoryGetIssueResolutionHistoryOptions{EntityId: "test", EntityType: "PORTFOLIO", Statistic: "MTTR", StartDate: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestHistoryService_GetIssueResolutionHistory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.History.GetIssueResolutionHistory(context.Background(), nil)
	require.Error(t, err)
}

func TestHistoryService_GetProjectIssueResolution(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/history/project-issue-resolution", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetProjectIssueResolution(context.Background(), &HistoryGetProjectIssueResolutionOptions{Statistic: "MTTR"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestHistoryService_GetProjectIssueResolution_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.History.GetProjectIssueResolution(context.Background(), nil)
	require.Error(t, err)
}

func TestHistoryService_GetProjectScaResolution(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/history/project-sca-resolution", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetProjectScaResolution(context.Background(), &HistoryGetProjectScaResolutionOptions{Statistic: "SCA_MTTR"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestHistoryService_GetProjectScaResolution_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.History.GetProjectScaResolution(context.Background(), nil)
	require.Error(t, err)
}

func TestHistoryService_GetScaResolutionHistory(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/history/sca-resolution-history", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetScaResolutionHistory(context.Background(), &HistoryGetScaResolutionHistoryOptions{EntityId: "test", EntityType: "PORTFOLIO", Statistic: "SCA_MTTR", StartDate: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestHistoryService_GetScaResolutionHistory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.History.GetScaResolutionHistory(context.Background(), nil)
	require.Error(t, err)
}
