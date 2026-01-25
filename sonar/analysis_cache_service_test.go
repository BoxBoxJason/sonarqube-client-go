package sonargo

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalysisCache_Clear(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/analysis_cache/clear" {
			t.Errorf("expected path /api/analysis_cache/clear, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.AnalysisCache.Clear(nil)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAnalysisCache_Clear_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "feature" {
			t.Errorf("expected branch 'feature', got %s", branch)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AnalysisCacheClearOption{
		Project: "my-project",
		Branch:  "feature",
	}

	resp, err := client.AnalysisCache.Clear(opt)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAnalysisCache_Get(t *testing.T) {
	// Mock binary data
	mockData := []byte("mock cached data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/analysis_cache/get" {
			t.Errorf("expected path /api/analysis_cache/get, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AnalysisCacheGetOption{
		Project: "my-project",
	}

	resp, err := client.AnalysisCache.Get(opt)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	if !bytes.Equal(body, mockData) {
		t.Errorf("expected body %q, got %q", mockData, body)
	}
}

func TestAnalysisCache_Get_WithOptions(t *testing.T) {
	// Mock binary data
	mockData := []byte("mock cached data with options")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "main" {
			t.Errorf("expected branch 'main', got %s", branch)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AnalysisCacheGetOption{
		Project: "my-project",
		Branch:  "main",
	}

	resp, err := client.AnalysisCache.Get(opt)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	if !bytes.Equal(body, mockData) {
		t.Errorf("expected body %q, got %q", mockData, body)
	}
}

func TestAnalysisCache_ValidateClearOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid.
	err := client.AnalysisCache.ValidateClearOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid.
	err = client.AnalysisCache.ValidateClearOpt(&AnalysisCacheClearOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}
}

func TestAnalysisCache_ValidateGetOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be invalid.
	err := client.AnalysisCache.ValidateGetOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Empty option should be invalid.
	err = client.AnalysisCache.ValidateGetOpt(&AnalysisCacheGetOption{})
	if err == nil {
		t.Error("expected error for empty option")
	}

	// Option with Project should be valid.
	err = client.AnalysisCache.ValidateGetOpt(&AnalysisCacheGetOption{Project: "my-project"})
	if err != nil {
		t.Errorf("expected nil error for valid option, got %v", err)
	}
}
