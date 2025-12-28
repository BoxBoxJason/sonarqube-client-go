package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDuplications_Show(t *testing.T) {
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
	opt := &DuplicationsShowOption{}
	_, resp, err := client.Duplications.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
