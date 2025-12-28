package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMeasures_Component(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &MeasuresComponentOption{}
	_, resp, err := client.Measures.Component(opt)
	if err != nil {
		t.Fatalf("Component failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMeasures_ComponentTree(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &MeasuresComponentTreeOption{}
	_, resp, err := client.Measures.ComponentTree(opt)
	if err != nil {
		t.Fatalf("ComponentTree failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMeasures_Search(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &MeasuresSearchOption{}
	_, resp, err := client.Measures.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMeasures_SearchHistory(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &MeasuresSearchHistoryOption{}
	_, resp, err := client.Measures.SearchHistory(opt)
	if err != nil {
		t.Fatalf("SearchHistory failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
