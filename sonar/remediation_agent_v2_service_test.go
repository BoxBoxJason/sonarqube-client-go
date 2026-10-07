package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemediationAgentService_DeleteScheduledAgentConfig(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/remediation-agent/scheduled-agent-configs/test-id", http.StatusOK))
	client := newTestClient(t, server.url())

	resp, err := client.V2.RemediationAgent.DeleteScheduledAgentConfig(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemediationAgentService_DeleteScheduledAgentConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.RemediationAgent.DeleteScheduledAgentConfig(context.Background(), "")
	require.Error(t, err)
}

func TestRemediationAgentService_GetBacklogJob(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/remediation-agent/backlog-jobs/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.GetBacklogJob(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_GetBacklogJob_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.GetBacklogJob(context.Background(), "")
	require.Error(t, err)
}

func TestRemediationAgentService_ListJobs(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/remediation-agent/jobs", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.ListJobs(context.Background(), &RemediationAgentListJobsOptions{ProjectKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_ListJobs_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.ListJobs(context.Background(), nil)
	require.Error(t, err)
}

func TestRemediationAgentService_GetPullRequestJob(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/remediation-agent/pull-request-jobs/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.GetPullRequestJob(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_GetPullRequestJob_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.GetPullRequestJob(context.Background(), "")
	require.Error(t, err)
}

func TestRemediationAgentService_SearchScheduledAgentConfigs(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/remediation-agent/scheduled-agent-configs", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.SearchScheduledAgentConfigs(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_GetScheduledAgentConfig(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/remediation-agent/scheduled-agent-configs/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.GetScheduledAgentConfig(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_GetScheduledAgentConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.GetScheduledAgentConfig(context.Background(), "")
	require.Error(t, err)
}

func TestRemediationAgentService_GetEffectiveScheduledAgentConfig(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/remediation-agent/scheduled-agent-configs/test-id/effective", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.GetEffectiveScheduledAgentConfig(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_GetEffectiveScheduledAgentConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.GetEffectiveScheduledAgentConfig(context.Background(), "")
	require.Error(t, err)
}

func TestRemediationAgentService_CreateBacklogJob(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/remediation-agent/backlog-jobs", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.CreateBacklogJob(context.Background(), &RemediationAgentCreateBacklogJobOptions{ProjectKey: "test", IssueKeys: []string{"test"}})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_CreateBacklogJob_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.CreateBacklogJob(context.Background(), nil)
	require.Error(t, err)
}

func TestRemediationAgentService_CreatePullRequestJob(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/remediation-agent/pull-request-jobs", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.CreatePullRequestJob(context.Background(), &RemediationAgentCreatePullRequestJobOptions{ProjectKey: "test", PullRequestKey: "test", IssueKeys: []string{"test"}})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_CreatePullRequestJob_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.CreatePullRequestJob(context.Background(), nil)
	require.Error(t, err)
}

func TestRemediationAgentService_CreateScheduledAgentConfig(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/remediation-agent/scheduled-agent-configs", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.CreateScheduledAgentConfig(context.Background(), &RemediationAgentCreateScheduledAgentConfigOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_CreateScheduledAgentConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.CreateScheduledAgentConfig(context.Background(), nil)
	require.Error(t, err)
}

func TestRemediationAgentService_UpdateScheduledAgentConfig(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPut, "/v2/remediation-agent/scheduled-agent-configs/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.RemediationAgent.UpdateScheduledAgentConfig(context.Background(), "test-id", &RemediationAgentUpdateScheduledAgentConfigOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestRemediationAgentService_UpdateScheduledAgentConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.RemediationAgent.UpdateScheduledAgentConfig(context.Background(), "", nil)
	require.Error(t, err)
}
