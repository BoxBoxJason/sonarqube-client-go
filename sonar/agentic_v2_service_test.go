package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgenticService_SearchJobs(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/agentic/jobs", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Agentic.SearchJobs(context.Background(), &AgenticSearchJobsOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestAgenticService_SearchJobs_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Agentic.SearchJobs(context.Background(), nil)
	require.Error(t, err)
}

func TestAgenticService_DownloadJobLogs(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/agentic/jobs/logs", http.StatusOK, "zip-bytes"))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Agentic.DownloadJobLogs(context.Background(), &AgenticDownloadJobLogsOptions{JobId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestAgenticService_DownloadJobLogs_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Agentic.DownloadJobLogs(context.Background(), nil)
	require.Error(t, err)
}
