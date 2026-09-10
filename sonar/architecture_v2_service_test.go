package sonar

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

// -----------------------------------------------------------------------------
// Directives
// -----------------------------------------------------------------------------

func TestArchitectureService_ListDirectives(t *testing.T) {
	response := map[string]any{
		"directives": []map[string]any{
			{
				"id":        "d1",
				"projectId": "p1",
				"origin":    "manual",
				"directiveData": map[string]any{
					"fromNodeKey":    "a",
					"toNodeKey":      "b",
					"graphType":      "file_graph",
					"graphEcoSystem": "java",
					"type":           "REMOVE_EDGE",
				},
			},
		},
		"page": map[string]any{"pageIndex": 1, "pageSize": 50, "total": 1},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/directives", http.StatusOK,
		map[string]string{"projectId": "p1"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.ListDirectives(context.Background(), &ArchitectureListDirectivesOptions{ProjectId: "p1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Directives, 1)
	assert.Equal(t, "d1", result.Directives[0].Id)
	assert.Equal(t, "REMOVE_EDGE", result.Directives[0].DirectiveData.Type)
	assert.Equal(t, int32(1), result.Page.Total)
}

func TestArchitectureService_ListDirectives_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListDirectives(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.ListDirectives(context.Background(), &ArchitectureListDirectivesOptions{})
	assert.Error(t, err)
}

func TestArchitectureService_CreateDirective(t *testing.T) {
	response := map[string]any{"id": "d2", "projectId": "p1"}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/directives", http.StatusCreated, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.CreateDirective(context.Background(), &ArchitectureCreateDirectiveOptions{
		ProjectId: "p1",
		DirectiveData: ArchitectureDirectiveData{
			FromNodeKey:    "a",
			GraphType:      "file_graph",
			GraphEcoSystem: "java",
			Type:           "REMOVE_EDGE",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "d2", result.Id)
}

func TestArchitectureService_CreateDirective_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreateDirective(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.CreateDirective(context.Background(), &ArchitectureCreateDirectiveOptions{ProjectId: "p1"})
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.CreateDirective(context.Background(), &ArchitectureCreateDirectiveOptions{
		ProjectId:     "p1",
		DirectiveData: ArchitectureDirectiveData{FromNodeKey: "a", GraphType: "bogus", GraphEcoSystem: "java", Type: "REMOVE_EDGE"},
	})
	assert.Error(t, err)
}

func TestArchitectureService_GetDirective(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/directives/d1", http.StatusOK,
		map[string]any{"id": "d1", "projectId": "p1"}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetDirective(context.Background(), "d1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "d1", result.Id)
}

func TestArchitectureService_GetDirective_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetDirective(context.Background(), "")
	assert.Error(t, err)
}

func TestArchitectureService_DeleteDirective(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/directives/d1", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Architecture.DeleteDirective(context.Background(), "d1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeleteDirective_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeleteDirective(context.Background(), "")
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Intended architecture
// -----------------------------------------------------------------------------

func TestArchitectureService_GetIntendedArchitecture(t *testing.T) {
	response := []map[string]any{
		{
			"id":       "i1",
			"language": "java",
			"type":     "namespace",
			"constraints": []map[string]any{
				{"from": "a", "to": "b"},
			},
			"groups": []map[string]any{
				{"label": "core", "patterns": []string{"com.core.*"}},
			},
		},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/intended", http.StatusOK,
		map[string]string{"projectId": "p1"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetIntendedArchitecture(context.Background(), &ArchitectureIntendedOptions{ProjectId: "p1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "i1", result[0].Id)
	require.Len(t, result[0].Constraints, 1)
	assert.Equal(t, "b", result[0].Constraints[0].To)
}

func TestArchitectureService_GetIntendedArchitecture_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetIntendedArchitecture(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.GetIntendedArchitecture(context.Background(), &ArchitectureIntendedOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Models
// -----------------------------------------------------------------------------

func TestArchitectureService_ListModels(t *testing.T) {
	response := []map[string]any{
		{
			"id":        "m1",
			"projectId": "p1",
			"model": map[string]any{
				"perspectives": []map[string]any{
					{"qualifiers": "file", "language": "java"},
				},
			},
		},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/models", http.StatusOK,
		map[string]string{"projectId": "p1"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.ListModels(context.Background(), &ArchitectureModelsListOptions{ProjectId: "p1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "m1", result[0].Id)
	require.Len(t, result[0].Model.Perspectives, 1)
	assert.Equal(t, "java", result[0].Model.Perspectives[0].Language)
}

func TestArchitectureService_ListModels_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListModels(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.ListModels(context.Background(), &ArchitectureModelsListOptions{})
	assert.Error(t, err)
}

func TestArchitectureService_CreateModel(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/models", http.StatusCreated,
		map[string]any{"id": "m2", "projectId": "p1"}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.CreateModel(context.Background(), &ArchitectureCreateModelOptions{
		ProjectId: "p1",
		Model: &ArchitectureModelData{
			Perspectives: []ArchitectureModelPerspective{{Qualifiers: "file", Language: "java"}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "m2", result.Id)
}

func TestArchitectureService_CreateModel_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreateModel(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.CreateModel(context.Background(), &ArchitectureCreateModelOptions{})
	assert.Error(t, err)
}

func TestArchitectureService_GetModel(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/models/m1", http.StatusOK,
		map[string]any{"id": "m1", "projectId": "p1"}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetModel(context.Background(), "m1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "m1", result.Id)
}

func TestArchitectureService_GetModel_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetModel(context.Background(), "")
	assert.Error(t, err)
}

func TestArchitectureService_UpdateModel(t *testing.T) {
	server := newTestServer(t, mockPatchHandler(t, "/v2/architecture/models/m1", http.StatusOK,
		map[string]any{"data": "x"}, map[string]any{"id": "m1", "data": "x"}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.UpdateModel(context.Background(), "m1", &ArchitectureUpdateModelOptions{Data: "x"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "m1", result.Id)
}

func TestArchitectureService_UpdateModel_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.UpdateModel(context.Background(), "", &ArchitectureUpdateModelOptions{})
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.UpdateModel(context.Background(), "m1", nil)
	assert.Error(t, err)
}

func TestArchitectureService_DeleteModel(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/models/m1", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Architecture.DeleteModel(context.Background(), "m1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeleteModel_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeleteModel(context.Background(), "")
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Project configurations
// -----------------------------------------------------------------------------

func TestArchitectureService_GetProjectConfigurations(t *testing.T) {
	response := []map[string]any{
		{
			"projectId":            "p1",
			"intendedArchitecture": "def",
			"boundaryDescriptors":  []string{"b1"},
			"directives":           []string{"d1"},
		},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/project-configurations", http.StatusOK,
		map[string]string{"projectKey": "p1"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetProjectConfigurations(context.Background(), &ArchitectureProjectConfigurationsOptions{ProjectKey: "p1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "p1", result[0].ProjectId)
	assert.Equal(t, []string{"b1"}, result[0].BoundaryDescriptors)
}

func TestArchitectureService_GetPrivateProjectConfigurations(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/private/architecture/project-configurations", http.StatusOK,
		[]map[string]any{{"projectId": "p1"}}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetPrivateProjectConfigurations(context.Background(), &ArchitectureProjectConfigurationsOptions{ProjectId: "p1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "p1", result[0].ProjectId)
}

// -----------------------------------------------------------------------------
// Snapshots
// -----------------------------------------------------------------------------

func TestArchitectureService_ListSnapshots(t *testing.T) {
	response := map[string]any{
		"snapshots": []map[string]any{
			{
				"id":             "s1",
				"projectId":      "p1",
				"organizationId": "o1",
				"branchId":       "b1",
				"aggregatedMeasures": map[string]any{
					"architectureCoverage": "0.9",
					"numGraphNodes":        "10",
				},
				"measures": []map[string]any{
					{"ecosystem": "java", "values": map[string]any{"numGraphNodes": "10"}},
				},
			},
		},
		"page": map[string]any{"pageIndex": 1, "pageSize": 50, "total": 1},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/snapshots", http.StatusOK,
		map[string]string{"organizationId": "o1", "branchId": "b1"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.ListSnapshots(context.Background(), &ArchitectureListSnapshotsOptions{
		OrganizationId: "o1",
		BranchId:       "b1",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Snapshots, 1)
	assert.Equal(t, "s1", result.Snapshots[0].Id)
	assert.Equal(t, "0.9", result.Snapshots[0].AggregatedMeasures.ArchitectureCoverage)
	require.Len(t, result.Snapshots[0].Measures, 1)
	assert.Equal(t, "java", result.Snapshots[0].Measures[0].Ecosystem)
}

func TestArchitectureService_ListSnapshots_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListSnapshots(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.ListSnapshots(context.Background(), &ArchitectureListSnapshotsOptions{OrganizationId: "o1"})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Structure
// -----------------------------------------------------------------------------

func TestArchitectureService_GetStructure(t *testing.T) {
	response := []map[string]any{
		{
			"id":        "g1",
			"nodeCount": 2,
			"edgeCount": 1,
			"nodes": []map[string]any{
				{
					"id":   1,
					"name": "root",
					"kind": "package",
					"edges": []map[string]any{
						{"toId": 2, "weight": 5},
					},
					"nodes": []map[string]any{
						{"id": 2, "name": "child", "kind": "class"},
					},
				},
			},
		},
	}

	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/architecture/structure", http.StatusOK,
		map[string]string{"organizationId": "o1", "branchId": "b1", "ecosystem": "java"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.GetStructure(context.Background(), &ArchitectureStructureOptions{
		OrganizationId: "o1",
		BranchId:       "b1",
		Ecosystem:      "java",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "g1", result[0].Id)
	require.Len(t, result[0].Nodes, 1)
	assert.Equal(t, "root", result[0].Nodes[0].Name)
	require.Len(t, result[0].Nodes[0].Edges, 1)
	assert.Equal(t, 2, result[0].Nodes[0].Edges[0].ToId)
	require.Len(t, result[0].Nodes[0].Nodes, 1)
	assert.Equal(t, "child", result[0].Nodes[0].Nodes[0].Name)
}

func TestArchitectureService_GetStructure_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetStructure(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.GetStructure(context.Background(), &ArchitectureStructureOptions{
		OrganizationId: "o1", BranchId: "b1", Ecosystem: "cobol",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Analysis and scanner data
// -----------------------------------------------------------------------------

func TestArchitectureService_StoreAnalysis(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v2/architecture/analysis", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode("analysis-123"))
	})
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.StoreAnalysis(context.Background(), &ArchitectureStoreAnalysisOptions{CeTaskId: "ce-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "analysis-123", *result)
}

func TestArchitectureService_StoreAnalysis_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.StoreAnalysis(context.Background(), nil)
	assert.Error(t, err)
}

func TestArchitectureService_InitiateScannerDataUpload(t *testing.T) {
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/architecture/scanner-data", http.StatusCreated,
		map[string]any{"analysisId": "a1"}, map[string]any{"id": "u1", "uploadUrl": "https://example/upload"}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Architecture.InitiateScannerDataUpload(context.Background(), &ArchitectureScannerDataOptions{AnalysisId: "a1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "u1", result.Id)
	assert.Equal(t, "https://example/upload", result.UploadUrl)
}

func TestArchitectureService_InitiateScannerDataUpload_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.InitiateScannerDataUpload(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.V2.Architecture.InitiateScannerDataUpload(context.Background(), &ArchitectureScannerDataOptions{})
	assert.Error(t, err)
}

func TestArchitectureService_UploadScannerData(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v2/architecture/scanner-data-upload", r.URL.Path)
		assert.Equal(t, "ce-1", r.URL.Query().Get("ceTaskId"))
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		file, header, ferr := r.FormFile("payload")
		require.NoError(t, ferr)

		defer func() { _ = file.Close() }()

		data, ferr := io.ReadAll(file)
		require.NoError(t, ferr)
		assert.Equal(t, "report-bytes", string(data))
		assert.Equal(t, "payload", header.Filename)

		w.WriteHeader(http.StatusCreated)
	})
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Architecture.UploadScannerData(context.Background(), &ArchitectureUploadScannerDataOptions{
		CeTaskId: "ce-1",
		Payload:  strings.NewReader("report-bytes"),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestArchitectureService_UploadScannerData_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.UploadScannerData(context.Background(), nil)
	assert.Error(t, err)

	_, err = client.V2.Architecture.UploadScannerData(context.Background(), &ArchitectureUploadScannerDataOptions{CeTaskId: "ce-1"})
	assert.Error(t, err)

	_, err = client.V2.Architecture.UploadScannerData(context.Background(), &ArchitectureUploadScannerDataOptions{Payload: strings.NewReader("x")})
	assert.Error(t, err)
}
