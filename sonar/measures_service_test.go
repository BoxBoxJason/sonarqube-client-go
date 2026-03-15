package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// MeasuresService Test Suite
// -----------------------------------------------------------------------------

// TestMeasuresService_Component tests the Component method.
func TestMeasuresService_Component(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/measures/component", http.StatusOK, `{
		"component": {
			"key": "my-project",
			"name": "My Project",
			"qualifier": "TRK",
			"measures": [
				{"metric": "coverage", "value": "85.5"},
				{"metric": "bugs", "value": "12"},
				{"metric": "vulnerabilities", "value": "3"}
			]
		},
		"metrics": [
			{"key": "coverage", "name": "Coverage", "type": "PERCENT"}
		]
	}`))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &MeasuresComponentOptions{
		Component:  "my-project",
		MetricKeys: []string{"coverage", "bugs", "vulnerabilities"},
	}

	result, resp, err := client.Measures.Component(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my-project", result.Component.Key)
	assert.Len(t, result.Component.Measures, 3)
}

// TestMeasuresService_Component_ValidationError tests validation for Component.
func TestMeasuresService_Component_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Component
	opt := &MeasuresComponentOptions{
		MetricKeys: []string{"coverage"},
	}
	_, _, err := client.Measures.Component(opt)
	assert.Error(t, err)

	// Test missing MetricKeys
	opt = &MeasuresComponentOptions{
		Component: "my-project",
	}
	_, _, err = client.Measures.Component(opt)
	assert.Error(t, err)
}

// TestMeasuresService_ComponentTree tests the ComponentTree method.
func TestMeasuresService_ComponentTree(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/measures/component_tree", http.StatusOK, `{
		"paging": {"pageIndex": 1, "pageSize": 100, "total": 2},
		"baseComponent": {
			"key": "my-project",
			"name": "My Project",
			"qualifier": "TRK"
		},
		"components": [
			{
				"key": "my-project:src",
				"name": "src",
				"qualifier": "DIR",
				"measures": [{"metric": "coverage", "value": "90.0"}]
			},
			{
				"key": "my-project:src/main.go",
				"name": "main.go",
				"qualifier": "FIL",
				"measures": [{"metric": "coverage", "value": "85.0"}]
			}
		],
		"metrics": [{"key": "coverage", "name": "Coverage", "type": "PERCENT"}]
	}`))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &MeasuresComponentTreeOptions{
		Component:  "my-project",
		MetricKeys: []string{"coverage"},
		Strategy:   MeasureStrategyAll,
	}

	result, resp, err := client.Measures.ComponentTree(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my-project", result.BaseComponent.Key)
	assert.Len(t, result.Components, 2)
}

// TestMeasuresService_ComponentTree_ValidationError tests validation for ComponentTree.
func TestMeasuresService_ComponentTree_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Component
	opt := &MeasuresComponentTreeOptions{
		MetricKeys: []string{"coverage"},
	}
	_, _, err := client.Measures.ComponentTree(opt)
	assert.Error(t, err)

	// Test missing MetricKeys
	opt = &MeasuresComponentTreeOptions{
		Component: "my-project",
	}
	_, _, err = client.Measures.ComponentTree(opt)
	assert.Error(t, err)

	// Test invalid Strategy
	opt = &MeasuresComponentTreeOptions{
		Component:  "my-project",
		MetricKeys: []string{"coverage"},
		Strategy:   "invalid",
	}
	_, _, err = client.Measures.ComponentTree(opt)
	assert.Error(t, err)

	// Test invalid MetricSortFilter
	opt = &MeasuresComponentTreeOptions{
		Component:        "my-project",
		MetricKeys:       []string{"coverage"},
		MetricSortFilter: "invalid",
	}
	_, _, err = client.Measures.ComponentTree(opt)
	assert.Error(t, err)
}

// TestMeasuresService_Search tests the Search method.
func TestMeasuresService_Search(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/measures/search", http.StatusOK, `{
		"measures": [
			{"component": "project1", "metric": "coverage", "value": "80.0"},
			{"component": "project2", "metric": "coverage", "value": "75.0"}
		]
	}`))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &MeasuresSearchOptions{
		MetricKeys:  []string{"coverage"},
		ProjectKeys: []string{"project1", "project2"},
	}

	result, resp, err := client.Measures.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Measures, 2)
}

// TestMeasuresService_Search_ValidationError tests validation for Search.
func TestMeasuresService_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing MetricKeys
	opt := &MeasuresSearchOptions{
		ProjectKeys: []string{"project1"},
	}
	_, _, err := client.Measures.Search(opt)
	assert.Error(t, err)

	// Test missing ProjectKeys
	opt = &MeasuresSearchOptions{
		MetricKeys: []string{"coverage"},
	}
	_, _, err = client.Measures.Search(opt)
	assert.Error(t, err)
}

// TestMeasuresService_SearchHistory tests the SearchHistory method.
func TestMeasuresService_SearchHistory(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/measures/search_history", http.StatusOK, `{
		"paging": {"pageIndex": 1, "pageSize": 100, "total": 3},
		"measures": [
			{
				"metric": "coverage",
				"history": [
					{"date": "2024-01-01T00:00:00+0000", "value": "75.0"},
					{"date": "2024-01-15T00:00:00+0000", "value": "80.0"},
					{"date": "2024-02-01T00:00:00+0000", "value": "85.0"}
				]
			}
		]
	}`))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &MeasuresSearchHistoryOptions{
		Component: "my-project",
		Metrics:   []string{"coverage"},
		From:      "2024-01-01",
		To:        "2024-03-01",
	}

	result, resp, err := client.Measures.SearchHistory(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Measures, 1)
	assert.Len(t, result.Measures[0].History, 3)
}

// TestMeasuresService_SearchHistory_ValidationError tests validation for SearchHistory.
func TestMeasuresService_SearchHistory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Component
	opt := &MeasuresSearchHistoryOptions{
		Metrics: []string{"coverage"},
	}
	_, _, err := client.Measures.SearchHistory(opt)
	assert.Error(t, err)

	// Test missing Metrics
	opt = &MeasuresSearchHistoryOptions{
		Component: "my-project",
	}
	_, _, err = client.Measures.SearchHistory(opt)
	assert.Error(t, err)
}
