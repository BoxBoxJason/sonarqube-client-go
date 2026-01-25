package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetrics_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/metrics/search" {
			t.Errorf("expected path /api/metrics/search, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"metrics": [
				{
					"id": "1",
					"key": "coverage",
					"name": "Coverage",
					"description": "Code coverage percentage",
					"domain": "Coverage",
					"type": "PERCENT",
					"direction": 1,
					"qualitative": true,
					"hidden": false,
					"custom": false
				},
				{
					"id": "2",
					"key": "bugs",
					"name": "Bugs",
					"description": "Number of bugs",
					"domain": "Reliability",
					"type": "INT",
					"direction": -1,
					"qualitative": true,
					"hidden": false,
					"custom": false
				}
			],
			"paging": {
				"pageIndex": 1,
				"pageSize": 100,
				"total": 2
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Metrics.Search(nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(result.Metrics))
	}

	if result.Metrics[0].Key != "coverage" {
		t.Errorf("expected first metric key 'coverage', got %s", result.Metrics[0].Key)
	}

	if result.Metrics[0].Name != "Coverage" {
		t.Errorf("expected first metric name 'Coverage', got %s", result.Metrics[0].Name)
	}

	if result.Metrics[0].Direction != 1 {
		t.Errorf("expected first metric direction 1, got %d", result.Metrics[0].Direction)
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected paging total 2, got %d", result.Paging.Total)
	}
}

func TestMetrics_Search_WithPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("p")
		if page != "2" {
			t.Errorf("expected page '2', got %s", page)
		}

		pageSize := r.URL.Query().Get("ps")
		if pageSize != "50" {
			t.Errorf("expected pageSize '50', got %s", pageSize)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"metrics": [], "paging": {"pageIndex": 2, "pageSize": 50, "total": 0}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &MetricsSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	}

	_, resp, err := client.Metrics.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMetrics_Types(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/metrics/types" {
			t.Errorf("expected path /api/metrics/types, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"types": ["INT", "FLOAT", "PERCENT", "BOOL", "STRING", "MILLISEC", "RATING", "DATA", "DISTRIB", "LEVEL", "WORK_DUR"]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Metrics.Types()
	if err != nil {
		t.Fatalf("Types failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Types) != 11 {
		t.Errorf("expected 11 types, got %d", len(result.Types))
	}

	if result.Types[0] != "INT" {
		t.Errorf("expected first type 'INT', got %s", result.Types[0])
	}
}

func TestMetrics_ValidateSearchOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid.
	err := client.Metrics.ValidateSearchOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid.
	err = client.Metrics.ValidateSearchOpt(&MetricsSearchOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}

	// Valid pagination should be valid.
	err = client.Metrics.ValidateSearchOpt(&MetricsSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     1,
			PageSize: 100,
		},
	})
	if err != nil {
		t.Errorf("expected nil error for valid pagination, got %v", err)
	}
}
