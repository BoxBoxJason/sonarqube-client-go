package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

func TestProjectPullRequestsService_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_pull_requests/delete", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.ProjectPullRequests.Delete(context.Background(), &ProjectPullRequestsDeleteOptions{
		Project:     "my-project",
		PullRequest: "123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectPullRequestsService_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.ProjectPullRequests.Delete(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.ProjectPullRequests.Delete(context.Background(), &ProjectPullRequestsDeleteOptions{PullRequest: "123"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.ProjectPullRequests.Delete(context.Background(), &ProjectPullRequestsDeleteOptions{Project: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// List
// -----------------------------------------------------------------------------

func TestProjectPullRequestsService_List(t *testing.T) {
	response := ProjectPullRequestsList{
		PullRequests: []ProjectPullRequest{
			{
				Key:      "123",
				Title:    "My Pull Request",
				Branch:   "feature/my-feature",
				Base:     "main",
				IsOrphan: false,
				URL:      "https://github.com/org/repo/pull/123",
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_pull_requests/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.ProjectPullRequests.List(context.Background(), &ProjectPullRequestsListOptions{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.PullRequests, 1)
	assert.Equal(t, "123", result.PullRequests[0].Key)
	assert.Equal(t, "My Pull Request", result.PullRequests[0].Title)
}

func TestProjectPullRequestsService_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.ProjectPullRequests.List(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.ProjectPullRequests.List(context.Background(), &ProjectPullRequestsListOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
