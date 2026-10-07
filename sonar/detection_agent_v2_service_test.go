package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectionAgentService_GetInstanceConfiguration(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/detection-agent/instance-config", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DetectionAgent.GetInstanceConfiguration(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDetectionAgentService_ListJobs(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/detection-agent/jobs", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DetectionAgent.ListJobs(context.Background(), &DetectionAgentListJobsOptions{ProjectId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDetectionAgentService_ListJobs_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DetectionAgent.ListJobs(context.Background(), nil)
	require.Error(t, err)
}

func TestDetectionAgentService_GetCreditCost(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/detection-agent/jobs/credit-cost", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DetectionAgent.GetCreditCost(context.Background(), &DetectionAgentGetCreditCostOptions{BranchId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDetectionAgentService_GetCreditCost_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DetectionAgent.GetCreditCost(context.Background(), nil)
	require.Error(t, err)
}

func TestDetectionAgentService_GetProjectConfiguration(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/detection-agent/project-configs/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DetectionAgent.GetProjectConfiguration(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDetectionAgentService_GetProjectConfiguration_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DetectionAgent.GetProjectConfiguration(context.Background(), "")
	require.Error(t, err)
}

func TestDetectionAgentService_UpsertProjectConfiguration(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/detection-agent/project-configs/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DetectionAgent.UpsertProjectConfiguration(context.Background(), "test-id", &DetectionAgentUpsertProjectConfigurationOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDetectionAgentService_UpsertProjectConfiguration_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DetectionAgent.UpsertProjectConfiguration(context.Background(), "", nil)
	require.Error(t, err)
}

func TestDetectionAgentService_CreateJob(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/detection-agent/jobs", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DetectionAgent.CreateJob(context.Background(), &DetectionAgentCreateJobOptions{ProjectId: "test", BranchId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDetectionAgentService_CreateJob_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DetectionAgent.CreateJob(context.Background(), nil)
	require.Error(t, err)
}

func TestDetectionAgentService_SubmitJobResult(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/v2/detection-agent/jobs/results", http.StatusOK))
	client := newTestClient(t, server.url())

	resp, err := client.V2.DetectionAgent.SubmitJobResult(context.Background(), "test-value", &DetectionAgentSubmitJobResultOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDetectionAgentService_SubmitJobResult_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.DetectionAgent.SubmitJobResult(context.Background(), "", nil)
	require.Error(t, err)
}
