package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectLinks_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_links/create" {
			t.Errorf("expected path /api/project_links/create, got %s", r.URL.Path)
		}

		name := r.URL.Query().Get("name")
		if name != "Homepage" {
			t.Errorf("expected name 'Homepage', got %s", name)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		url := r.URL.Query().Get("url")
		if url != "https://example.com" {
			t.Errorf("expected url 'https://example.com', got %s", url)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"link": {
				"id": "1",
				"name": "Homepage",
				"type": "homepage",
				"url": "https://example.com"
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectLinksCreateOption{
		Name:       "Homepage",
		ProjectKey: "my-project",
		URL:        "https://example.com",
	}

	result, resp, err := client.ProjectLinks.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Link.ID != "1" {
		t.Errorf("expected link ID '1', got %s", result.Link.ID)
	}

	if result.Link.Name != "Homepage" {
		t.Errorf("expected link Name 'Homepage', got %s", result.Link.Name)
	}
}

func TestProjectLinks_Create_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectLinks.Create(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Name should fail validation.
	_, _, err = client.ProjectLinks.Create(&ProjectLinksCreateOption{
		ProjectKey: "my-project",
		URL:        "https://example.com",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Missing URL should fail validation.
	_, _, err = client.ProjectLinks.Create(&ProjectLinksCreateOption{
		Name:       "Homepage",
		ProjectKey: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Missing ProjectID and ProjectKey should fail validation.
	_, _, err = client.ProjectLinks.Create(&ProjectLinksCreateOption{
		Name: "Homepage",
		URL:  "https://example.com",
	})
	if err == nil {
		t.Error("expected error for missing ProjectID and ProjectKey")
	}
}

func TestProjectLinks_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_links/delete" {
			t.Errorf("expected path /api/project_links/delete, got %s", r.URL.Path)
		}

		id := r.URL.Query().Get("id")
		if id != "1" {
			t.Errorf("expected id '1', got %s", id)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectLinksDeleteOption{
		ID: "1",
	}

	resp, err := client.ProjectLinks.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectLinks_Delete_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectLinks.Delete(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing ID should fail validation.
	_, err = client.ProjectLinks.Delete(&ProjectLinksDeleteOption{})
	if err == nil {
		t.Error("expected error for missing ID")
	}
}

func TestProjectLinks_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_links/search" {
			t.Errorf("expected path /api/project_links/search, got %s", r.URL.Path)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"links": [
				{"id": "1", "name": "Homepage", "type": "homepage", "url": "https://example.com"},
				{"id": "2", "name": "CI", "type": "ci", "url": "https://ci.example.com"}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectLinksSearchOption{
		ProjectKey: "my-project",
	}

	result, resp, err := client.ProjectLinks.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Links) != 2 {
		t.Errorf("expected 2 links, got %d", len(result.Links))
	}

	if result.Links[0].ID != "1" {
		t.Errorf("expected first link ID '1', got %s", result.Links[0].ID)
	}
}

func TestProjectLinks_Search_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectLinks.Search(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing ProjectID and ProjectKey should fail validation.
	_, _, err = client.ProjectLinks.Search(&ProjectLinksSearchOption{})
	if err == nil {
		t.Error("expected error for missing ProjectID and ProjectKey")
	}
}

func TestProjectLinks_ValidateCreateOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option with ProjectKey should pass.
	err := client.ProjectLinks.ValidateCreateOpt(&ProjectLinksCreateOption{
		Name:       "Homepage",
		ProjectKey: "my-project",
		URL:        "https://example.com",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Valid option with ProjectID should pass.
	err = client.ProjectLinks.ValidateCreateOpt(&ProjectLinksCreateOption{
		Name:      "Homepage",
		ProjectID: "project-id",
		URL:       "https://example.com",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxLinkNameLength+1; i++ {
		longName += "a"
	}
	err = client.ProjectLinks.ValidateCreateOpt(&ProjectLinksCreateOption{
		Name:       longName,
		ProjectKey: "my-project",
		URL:        "https://example.com",
	})
	if err == nil {
		t.Error("expected error for name exceeding max length")
	}
}
