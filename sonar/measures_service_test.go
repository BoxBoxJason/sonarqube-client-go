package sonargo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// MeasuresService Test Suite
// -----------------------------------------------------------------------------

// TestMeasuresService_Component tests the Component method.
func TestMeasuresService_Component(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/measures/component") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("component") != "my-project" {
			t.Errorf("unexpected component: %s", r.URL.Query().Get("component"))
		}
		if r.URL.Query().Get("metricKeys") != "coverage,bugs,vulnerabilities" {
			t.Errorf("unexpected metricKeys: %s", r.URL.Query().Get("metricKeys"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
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
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &MeasuresComponentOption{
		Component:  "my-project",
		MetricKeys: []string{"coverage", "bugs", "vulnerabilities"},
	}

	result, resp, err := client.Measures.Component(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.Component.Key != "my-project" {
		t.Errorf("unexpected component key: %s", result.Component.Key)
	}
	if len(result.Component.Measures) != 3 {
		t.Errorf("expected 3 measures, got %d", len(result.Component.Measures))
	}
}

// TestMeasuresService_Component_ValidationError tests validation for Component.
func TestMeasuresService_Component_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Component
	opt := &MeasuresComponentOption{
		MetricKeys: []string{"coverage"},
	}
	_, _, err := client.Measures.Component(opt)
	if err == nil {
		t.Error("expected validation error for missing Component")
	}

	// Test missing MetricKeys
	opt = &MeasuresComponentOption{
		Component: "my-project",
	}
	_, _, err = client.Measures.Component(opt)
	if err == nil {
		t.Error("expected validation error for missing MetricKeys")
	}
}

// TestMeasuresService_ComponentTree tests the ComponentTree method.
func TestMeasuresService_ComponentTree(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/measures/component_tree") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
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
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &MeasuresComponentTreeOption{
		Component:  "my-project",
		MetricKeys: []string{"coverage"},
		Strategy:   "all",
	}

	result, resp, err := client.Measures.ComponentTree(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.BaseComponent.Key != "my-project" {
		t.Errorf("unexpected base component key: %s", result.BaseComponent.Key)
	}
	if len(result.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(result.Components))
	}
}

// TestMeasuresService_ComponentTree_ValidationError tests validation for ComponentTree.
func TestMeasuresService_ComponentTree_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Component
	opt := &MeasuresComponentTreeOption{
		MetricKeys: []string{"coverage"},
	}
	_, _, err := client.Measures.ComponentTree(opt)
	if err == nil {
		t.Error("expected validation error for missing Component")
	}

	// Test missing MetricKeys
	opt = &MeasuresComponentTreeOption{
		Component: "my-project",
	}
	_, _, err = client.Measures.ComponentTree(opt)
	if err == nil {
		t.Error("expected validation error for missing MetricKeys")
	}

	// Test invalid Strategy
	opt = &MeasuresComponentTreeOption{
		Component:  "my-project",
		MetricKeys: []string{"coverage"},
		Strategy:   "invalid",
	}
	_, _, err = client.Measures.ComponentTree(opt)
	if err == nil {
		t.Error("expected validation error for invalid Strategy")
	}

	// Test invalid MetricSortFilter
	opt = &MeasuresComponentTreeOption{
		Component:        "my-project",
		MetricKeys:       []string{"coverage"},
		MetricSortFilter: "invalid",
	}
	_, _, err = client.Measures.ComponentTree(opt)
	if err == nil {
		t.Error("expected validation error for invalid MetricSortFilter")
	}
}

// TestMeasuresService_Search tests the Search method.
func TestMeasuresService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/measures/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"measures": [
				{"component": "project1", "metric": "coverage", "value": "80.0"},
				{"component": "project2", "metric": "coverage", "value": "75.0"}
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &MeasuresSearchOption{
		MetricKeys:  []string{"coverage"},
		ProjectKeys: []string{"project1", "project2"},
	}

	result, resp, err := client.Measures.Search(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Measures) != 2 {
		t.Errorf("expected 2 measures, got %d", len(result.Measures))
	}
}

// TestMeasuresService_Search_ValidationError tests validation for Search.
func TestMeasuresService_Search_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing MetricKeys
	opt := &MeasuresSearchOption{
		ProjectKeys: []string{"project1"},
	}
	_, _, err := client.Measures.Search(opt)
	if err == nil {
		t.Error("expected validation error for missing MetricKeys")
	}

	// Test missing ProjectKeys
	opt = &MeasuresSearchOption{
		MetricKeys: []string{"coverage"},
	}
	_, _, err = client.Measures.Search(opt)
	if err == nil {
		t.Error("expected validation error for missing ProjectKeys")
	}
}

// TestMeasuresService_SearchHistory tests the SearchHistory method.
func TestMeasuresService_SearchHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/measures/search_history") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
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
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &MeasuresSearchHistoryOption{
		Component: "my-project",
		Metrics:   []string{"coverage"},
		From:      "2024-01-01",
		To:        "2024-03-01",
	}

	result, resp, err := client.Measures.SearchHistory(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Measures) != 1 {
		t.Errorf("expected 1 measure, got %d", len(result.Measures))
	}
	if len(result.Measures[0].History) != 3 {
		t.Errorf("expected 3 history entries, got %d", len(result.Measures[0].History))
	}
}

// TestMeasuresService_SearchHistory_ValidationError tests validation for SearchHistory.
func TestMeasuresService_SearchHistory_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Component
	opt := &MeasuresSearchHistoryOption{
		Metrics: []string{"coverage"},
	}
	_, _, err := client.Measures.SearchHistory(opt)
	if err == nil {
		t.Error("expected validation error for missing Component")
	}

	// Test missing Metrics
	opt = &MeasuresSearchHistoryOption{
		Component: "my-project",
	}
	_, _, err = client.Measures.SearchHistory(opt)
	if err == nil {
		t.Error("expected validation error for missing Metrics")
	}
}
