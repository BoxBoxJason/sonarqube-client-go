package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatures_List(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/features/list", http.StatusOK, `["branch", "pr-decoration"]`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Features.List()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, "branch", (*result)[0])
	assert.Equal(t, "pr-decoration", (*result)[1])
}

func TestFeatures_List_Empty(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/features/list", http.StatusOK, `[]`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Features.List()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Empty(t, *result)
}
