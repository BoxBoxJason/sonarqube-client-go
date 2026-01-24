package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Version(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		// Return plain text version
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("9.9"))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	v, resp, err := client.Server.Version()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if v == nil || *v != "9.9" {
		t.Errorf("expected version '9.9', got '%v'", v)
	}
}

func TestServer_Version_ErrorResponse(t *testing.T) {
	// Create mock server that returns an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"errors":["internal error"]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method - expect an error
	v, resp, err := client.Server.Version()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Response should be present and indicate failure
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

	if v != nil {
		t.Errorf("expected nil version on error, got %v", v)
	}
}
