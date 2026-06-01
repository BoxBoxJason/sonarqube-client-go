package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// FileGraph
// -----------------------------------------------------------------------------

func TestArchitectureService_FileGraph(t *testing.T) {
	response := ArchitectureFileGraph{
		Nodes: []ArchitectureFileNode{
			{Id: "1", Name: "main.go", Path: "src/main.go"},
			{Id: "2", Name: "utils.go", Path: "src/utils.go"},
		},
		Edges: []ArchitectureFileEdge{
			{Source: "1", Target: "2"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/file-graph", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		ProjectKey: "my-project",
		BranchKey:  "main",
		Source:     "src/main.go",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Nodes, 2)
	assert.Len(t, result.Edges, 1)
}

func TestArchitectureService_FileGraph_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Architecture.FileGraph(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		BranchKey: "main",
		Source:    "src/main.go",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Architecture.FileGraph(context.Background(), &ArchitectureFileGraphOptions{
		ProjectKey: "my-project",
		Source:     "src/main.go",
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
