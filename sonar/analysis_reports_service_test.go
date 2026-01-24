package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalysisReports_IsQueueEmpty_True(t *testing.T) {
	// Create mock server that returns "true"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/analysis_reports/is_queue_empty" {
			t.Errorf("expected path /api/analysis_reports/is_queue_empty, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("true"))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.AnalysisReports.IsQueueEmpty()
	if err != nil {
		t.Fatalf("IsQueueEmpty failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.IsEmpty {
		t.Error("expected IsEmpty to be true")
	}
}

func TestAnalysisReports_IsQueueEmpty_False(t *testing.T) {
	// Create mock server that returns "false"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("false"))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.AnalysisReports.IsQueueEmpty()
	if err != nil {
		t.Fatalf("IsQueueEmpty failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.IsEmpty {
		t.Error("expected IsEmpty to be false")
	}
}
