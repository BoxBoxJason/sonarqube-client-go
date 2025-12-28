package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIssues_AddComment(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesAddCommentOption{}
	_, resp, err := client.Issues.AddComment(opt)
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_AnticipatedTransitions(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &IssuesAnticipatedTransitionsOption{}
	resp, err := client.Issues.AnticipatedTransitions(opt)
	if err != nil {
		t.Fatalf("AnticipatedTransitions failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestIssues_Assign(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesAssignOption{}
	_, resp, err := client.Issues.Assign(opt)
	if err != nil {
		t.Fatalf("Assign failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_Authors(t *testing.T) {
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
	opt := &IssuesAuthorsOption{}
	_, resp, err := client.Issues.Authors(opt)
	if err != nil {
		t.Fatalf("Authors failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_BulkChange(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesBulkChangeOption{}
	_, resp, err := client.Issues.BulkChange(opt)
	if err != nil {
		t.Fatalf("BulkChange failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_Changelog(t *testing.T) {
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
	opt := &IssuesChangelogOption{}
	_, resp, err := client.Issues.Changelog(opt)
	if err != nil {
		t.Fatalf("Changelog failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_ComponentTags(t *testing.T) {
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
	opt := &IssuesComponentTagsOption{}
	_, resp, err := client.Issues.ComponentTags(opt)
	if err != nil {
		t.Fatalf("ComponentTags failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_DeleteComment(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesDeleteCommentOption{}
	_, resp, err := client.Issues.DeleteComment(opt)
	if err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_DoTransition(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesDoTransitionOption{}
	_, resp, err := client.Issues.DoTransition(opt)
	if err != nil {
		t.Fatalf("DoTransition failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_EditComment(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesEditCommentOption{}
	_, resp, err := client.Issues.EditComment(opt)
	if err != nil {
		t.Fatalf("EditComment failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_List(t *testing.T) {
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
	opt := &IssuesListOption{}
	_, resp, err := client.Issues.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_Pull(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("[]"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &IssuesPullOption{}
	_, resp, err := client.Issues.Pull(opt)
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_PullTaint(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("[]"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &IssuesPullTaintOption{}
	_, resp, err := client.Issues.PullTaint(opt)
	if err != nil {
		t.Fatalf("PullTaint failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_Reindex(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &IssuesReindexOption{}
	resp, err := client.Issues.Reindex(opt)
	if err != nil {
		t.Fatalf("Reindex failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestIssues_Search(t *testing.T) {
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
	opt := &IssuesSearchOption{}
	_, resp, err := client.Issues.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_SetSeverity(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesSetSeverityOption{}
	_, resp, err := client.Issues.SetSeverity(opt)
	if err != nil {
		t.Fatalf("SetSeverity failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_SetTags(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesSetTagsOption{}
	_, resp, err := client.Issues.SetTags(opt)
	if err != nil {
		t.Fatalf("SetTags failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_SetType(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
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
	opt := &IssuesSetTypeOption{}
	_, resp, err := client.Issues.SetType(opt)
	if err != nil {
		t.Fatalf("SetType failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIssues_Tags(t *testing.T) {
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
	opt := &IssuesTagsOption{}
	_, resp, err := client.Issues.Tags(opt)
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
