package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wireHandler checks the path and query exactly as they are received.
func wireHandler(t *testing.T, path, query string) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, path, r.URL.EscapedPath())
		assert.Equal(t, query, r.URL.RawQuery)
		w.WriteHeader(http.StatusNoContent)
	}
}

func TestPathIdentifiersAreEscapedOnce(t *testing.T) {
	tests := map[string]struct {
		call func(client *Client) (*http.Response, error)
		path string
	}{
		"PermissionMappingRoleWithSpace": {
			call: func(client *Client) (*http.Response, error) {
				return client.V2.DopTranslation.DeleteGithubPermissionMapping(context.Background(), "Security Reviewer")
			},
			path: "/v2/dop-translation/github-permission-mappings/Security%20Reviewer",
		},
		"ScheduledAgentConfigWildcard": {
			call: func(client *Client) (*http.Response, error) {
				return client.V2.RemediationAgent.DeleteScheduledAgentConfig(context.Background(), "*")
			},
			path: "/v2/remediation-agent/scheduled-agent-configs/%2A",
		},
		"IdentifierWithSlashAndPercent": {
			call: func(client *Client) (*http.Response, error) {
				return client.V2.DopTranslation.DeleteGitlabPermissionMapping(context.Background(), "a/b%20c")
			},
			path: "/v2/dop-translation/gitlab-permission-mappings/a%2Fb%2520c",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server := newTestServer(t, wireHandler(t, test.path, ""))
			client := newTestClient(t, server.URL)

			_, err := test.call(client)
			require.NoError(t, err)
		})
	}
}
