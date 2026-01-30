package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGithubProvisioning_Check(t *testing.T) {
	response := `{
		"application": {
			"autoProvisioning": {"status": "success"},
			"jit": {"status": "success"}
		},
		"installations": [
			{
				"autoProvisioning": {"status": "success"},
				"jit": {"status": "success"},
				"organization": "my-org"
			}
		]
	}`
	handler := mockHandler(t, http.MethodPost, "/github_provisioning/check", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.GithubProvisioning.Check()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "success", result.Application.AutoProvisioning.Status)
	assert.Equal(t, "success", result.Application.Jit.Status)
	require.Len(t, result.Installations, 1)
	assert.Equal(t, "my-org", result.Installations[0].Organization)
}

func TestGithubProvisioning_Check_WithError(t *testing.T) {
	response := `{
		"application": {
			"autoProvisioning": {"status": "failed", "errorMessage": "Invalid token"},
			"jit": {"status": "success"}
		},
		"installations": []
	}`
	handler := mockHandler(t, http.MethodPost, "/github_provisioning/check", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.GithubProvisioning.Check()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "failed", result.Application.AutoProvisioning.Status)
	assert.Equal(t, "Invalid token", result.Application.AutoProvisioning.ErrorMessage)
}

func TestGithubProvisioning_Check_EmptyResponse(t *testing.T) {
	handler := mockHandler(t, http.MethodPost, "/github_provisioning/check", http.StatusOK, `{}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.GithubProvisioning.Check()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}
