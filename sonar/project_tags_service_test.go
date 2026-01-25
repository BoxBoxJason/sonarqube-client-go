package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectTags_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_tags/search" {
			t.Errorf("expected path /api/project_tags/search, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"tags": ["security", "performance", "bug"]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.ProjectTags.Search(nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(result.Tags))
	}

	if result.Tags[0] != "security" {
		t.Errorf("expected first tag to be 'security', got %s", result.Tags[0])
	}
}

func TestProjectTags_Search_WithPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("p")
		if page != "2" {
			t.Errorf("expected page '2', got %s", page)
		}

		pageSize := r.URL.Query().Get("ps")
		if pageSize != "10" {
			t.Errorf("expected pageSize '10', got %s", pageSize)
		}

		query := r.URL.Query().Get("q")
		if query != "sec" {
			t.Errorf("expected query 'sec', got %s", query)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tags": ["security"]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectTagsSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 10,
		},
		Query: "sec",
	}

	result, resp, err := client.ProjectTags.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result.Tags))
	}
}

func TestProjectTags_Set(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_tags/set" {
			t.Errorf("expected path /api/project_tags/set, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		tags := r.URL.Query().Get("tags")
		if tags != "security,performance" {
			t.Errorf("expected tags 'security,performance', got %s", tags)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectTagsSetOption{
		Project: "my-project",
		Tags:    []string{"security", "performance"},
	}

	resp, err := client.ProjectTags.Set(opt)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectTags_Set_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectTags.Set(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, err = client.ProjectTags.Set(&ProjectTagsSetOption{
		Tags: []string{"tag1"},
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectTags_ValidateSearchOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid.
	err := client.ProjectTags.ValidateSearchOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid.
	err = client.ProjectTags.ValidateSearchOpt(&ProjectTagsSearchOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}
}

func TestProjectTags_ValidateSetOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.ProjectTags.ValidateSetOpt(&ProjectTagsSetOption{
		Project: "my-project",
		Tags:    []string{"tag1", "tag2"},
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.ProjectTags.ValidateSetOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail.
	err = client.ProjectTags.ValidateSetOpt(&ProjectTagsSetOption{
		Tags: []string{"tag1"},
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}
