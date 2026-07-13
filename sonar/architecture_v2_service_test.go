package sonar

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// FileGraph
// -----------------------------------------------------------------------------

func TestArchitectureService_FileGraph(t *testing.T) {
	// The raw payload intentionally contains characters (quotes, newlines) that
	// only round-trip correctly if the response is treated as a JSON-encoded
	// string and unescaped, rather than copied as opaque raw bytes.
	rawPayload := "{\"nodes\":[{\"id\":\"1\"}],\"edges\":[]}\nwith a \"quote\""

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "unexpected HTTP method")
		assert.Equal(t, "/v2/architecture/file-graph", r.URL.Path, "unexpected URL path")
		// The API spec declares this endpoint's 200 response as
		// "application/json" only (unlike genuinely text/plain endpoints such
		// as AnalysisService.GetVersion). Live verification against a real
		// SonarQube 2025.2 Enterprise instance showed that V2 endpoints
		// strictly enforce their declared content type via Spring content
		// negotiation: sending "Accept: text/plain" against a JSON-only V2
		// endpoint returns 406 Not Acceptable instead of the payload. This
		// assertion locks in that the client must request application/json.
		assert.Equal(t, "application/json", r.Header.Get("Accept"), "FileGraph must not request text/plain")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(rawPayload))
	})
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		ProjectKey: "my-project",
		BranchKey:  "main",
		Source:     "java",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, rawPayload, *result)
}

func TestArchitectureService_FileGraph_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Architecture.FileGraph(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		BranchKey: "main",
		Source:    "java",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		ProjectKey: "my-project",
		Source:     "java",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		ProjectKey: "my-project",
		BranchKey:  "main",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// SearchGraphs
// -----------------------------------------------------------------------------

func TestArchitectureService_SearchGraphs(t *testing.T) {
	response := map[string]any{
		"graphs": []map[string]any{
			{
				"id":             "graph-1",
				"branchId":       "branch-1",
				"type":           "file_graph",
				"ecosystem":      "java",
				"perspectiveKey": "default",
				"graphVersion":   "1",
			},
		},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/graphs", http.StatusOK,
		map[string]string{"projectKey": "my-project", "branchKey": "main"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.SearchGraphs(context.Background(), &ArchitectureSearchGraphsOptions{
		ProjectKey: "my-project",
		BranchKey:  "main",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "graph-1", result[0]["id"])
	assert.Equal(t, "file_graph", result[0]["type"])
}

func TestArchitectureService_SearchGraphs_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Architecture.SearchGraphs(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.SearchGraphs(context.Background(), &ArchitectureSearchGraphsOptions{
		BranchKey: "main",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.SearchGraphs(context.Background(), &ArchitectureSearchGraphsOptions{
		ProjectKey: "my-project",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// GetGraph
// -----------------------------------------------------------------------------

func TestArchitectureService_GetGraph(t *testing.T) {
	// As with FileGraph, this payload intentionally contains quotes/newlines that
	// only round-trip correctly if the response is treated as a JSON-encoded
	// string and unescaped, rather than copied as opaque raw bytes.
	rawPayload := "{\"nodes\":[{\"id\":\"1\"}],\"edges\":[]}\nwith a \"quote\""

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "unexpected HTTP method")
		assert.Equal(t, "/v2/architecture/graphs/graph-1", r.URL.Path, "unexpected URL path")
		assert.Equal(t, "application/json", r.Header.Get("Accept"), "GetGraph must not request text/plain")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(rawPayload))
	})
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetGraph(context.Background(), "graph-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, rawPayload, *result)
}

func TestArchitectureService_GetGraph_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Architecture.GetGraph(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
