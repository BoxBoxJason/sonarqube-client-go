package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaService_ListReachabilityPaths(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/reachability/list-paths", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Sca.ListReachabilityPaths(context.Background(), &ScaListReachabilityPathsOptions{LanguageKey: "test", ProjectKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestScaService_ListReachabilityPaths_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Sca.ListReachabilityPaths(context.Background(), nil)
	require.Error(t, err)
}
