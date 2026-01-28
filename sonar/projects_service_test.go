package sonargo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// ProjectsService Test Suite
// -----------------------------------------------------------------------------

// TestProjectsService_BulkDelete tests the BulkDelete method.
func TestProjectsService_BulkDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/bulk_delete") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("projects") != "project1,project2" {
			t.Errorf("unexpected projects: %s", r.FormValue("projects"))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsBulkDeleteOption{
		Projects: []string{"project1", "project2"},
	}

	resp, err := client.Projects.BulkDelete(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestProjectsService_BulkDelete_ValidationError tests validation for BulkDelete.
func TestProjectsService_BulkDelete_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test no filter provided
	opt := &ProjectsBulkDeleteOption{}
	_, err := client.Projects.BulkDelete(opt)
	if err == nil {
		t.Error("expected validation error for no filter")
	}

	// Test invalid qualifier
	opt = &ProjectsBulkDeleteOption{
		Query:      "test",
		Qualifiers: []string{"INVALID"},
	}
	_, err = client.Projects.BulkDelete(opt)
	if err == nil {
		t.Error("expected validation error for invalid qualifier")
	}
}

// TestProjectsService_Create tests the Create method.
func TestProjectsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/create") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("name") != "My Project" {
			t.Errorf("unexpected name: %s", r.FormValue("name"))
		}
		if r.FormValue("project") != "my-project" {
			t.Errorf("unexpected project: %s", r.FormValue("project"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"project": {
				"key": "my-project",
				"name": "My Project",
				"qualifier": "TRK",
				"visibility": "private"
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsCreateOption{
		Name:       "My Project",
		Project:    "my-project",
		Visibility: "private",
	}

	result, resp, err := client.Projects.Create(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.Project.Key != "my-project" {
		t.Errorf("unexpected project key: %s", result.Project.Key)
	}
	if result.Project.Visibility != "private" {
		t.Errorf("unexpected visibility: %s", result.Project.Visibility)
	}
}

// TestProjectsService_Create_ValidationError tests validation for Create.
func TestProjectsService_Create_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Name
	opt := &ProjectsCreateOption{
		Project: "my-project",
	}
	_, _, err := client.Projects.Create(opt)
	if err == nil {
		t.Error("expected validation error for missing Name")
	}

	// Test missing Project
	opt = &ProjectsCreateOption{
		Name: "My Project",
	}
	_, _, err = client.Projects.Create(opt)
	if err == nil {
		t.Error("expected validation error for missing Project")
	}

	// Test Name too long
	opt = &ProjectsCreateOption{
		Name:    strings.Repeat("a", MaxProjectNameLength+1),
		Project: "my-project",
	}
	_, _, err = client.Projects.Create(opt)
	if err == nil {
		t.Error("expected validation error for Name too long")
	}

	// Test Project key too long
	opt = &ProjectsCreateOption{
		Name:    "My Project",
		Project: strings.Repeat("a", MaxProjectKeyLength+1),
	}
	_, _, err = client.Projects.Create(opt)
	if err == nil {
		t.Error("expected validation error for Project too long")
	}

	// Test invalid visibility
	opt = &ProjectsCreateOption{
		Name:       "My Project",
		Project:    "my-project",
		Visibility: "invalid",
	}
	_, _, err = client.Projects.Create(opt)
	if err == nil {
		t.Error("expected validation error for invalid Visibility")
	}
}

// TestProjectsService_Delete tests the Delete method.
func TestProjectsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/delete") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsDeleteOption{
		Project: "my-project",
	}

	resp, err := client.Projects.Delete(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestProjectsService_Delete_ValidationError tests validation for Delete.
func TestProjectsService_Delete_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Project
	opt := &ProjectsDeleteOption{}
	_, err := client.Projects.Delete(opt)
	if err == nil {
		t.Error("expected validation error for missing Project")
	}
}

// TestProjectsService_Search tests the Search method.
func TestProjectsService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 2},
			"components": [
				{
					"key": "project1",
					"name": "Project One",
					"qualifier": "TRK",
					"visibility": "public",
					"lastAnalysisDate": "2024-01-15T10:30:00+0000"
				},
				{
					"key": "project2",
					"name": "Project Two",
					"qualifier": "TRK",
					"visibility": "private",
					"managed": true
				}
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsSearchOption{
		Query:      "project",
		Qualifiers: []string{"TRK"},
	}

	result, resp, err := client.Projects.Search(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(result.Components))
	}
	if result.Components[0].Key != "project1" {
		t.Errorf("unexpected project key: %s", result.Components[0].Key)
	}
}

// TestProjectsService_Search_ValidationError tests validation for Search.
func TestProjectsService_Search_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test invalid qualifier
	opt := &ProjectsSearchOption{
		Qualifiers: []string{"INVALID"},
	}
	_, _, err := client.Projects.Search(opt)
	if err == nil {
		t.Error("expected validation error for invalid qualifier")
	}
}

// TestProjectsService_SearchMyProjects tests the SearchMyProjects method.
func TestProjectsService_SearchMyProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/search_my_projects") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1},
			"projects": [
				{
					"key": "my-project",
					"name": "My Project",
					"description": "A test project",
					"qualityGate": "OK"
				}
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsSearchMyProjectsOption{}

	result, resp, err := client.Projects.SearchMyProjects(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(result.Projects))
	}
}

// TestProjectsService_SearchMyScannableProjects tests the SearchMyScannableProjects method.
func TestProjectsService_SearchMyScannableProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/search_my_scannable_projects") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"projects": [
				{"key": "project1", "name": "Project One"},
				{"key": "project2", "name": "Project Two"}
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Projects.SearchMyScannableProjects(&ProjectsSearchMyScannableProjectsOption{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(result.Projects))
	}
}

// TestProjectsService_UpdateDefaultVisibility tests the UpdateDefaultVisibility method.
func TestProjectsService_UpdateDefaultVisibility(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/update_default_visibility") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsUpdateDefaultVisibilityOption{
		ProjectVisibility: "private",
	}

	resp, err := client.Projects.UpdateDefaultVisibility(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestProjectsService_UpdateDefaultVisibility_ValidationError tests validation for UpdateDefaultVisibility.
func TestProjectsService_UpdateDefaultVisibility_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing ProjectVisibility
	opt := &ProjectsUpdateDefaultVisibilityOption{}
	_, err := client.Projects.UpdateDefaultVisibility(opt)
	if err == nil {
		t.Error("expected validation error for missing ProjectVisibility")
	}

	// Test invalid visibility
	opt = &ProjectsUpdateDefaultVisibilityOption{
		ProjectVisibility: "invalid",
	}
	_, err = client.Projects.UpdateDefaultVisibility(opt)
	if err == nil {
		t.Error("expected validation error for invalid ProjectVisibility")
	}
}

// TestProjectsService_UpdateKey tests the UpdateKey method.
func TestProjectsService_UpdateKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/update_key") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsUpdateKeyOption{
		From: "old-project-key",
		To:   "new-project-key",
	}

	resp, err := client.Projects.UpdateKey(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestProjectsService_UpdateKey_ValidationError tests validation for UpdateKey.
func TestProjectsService_UpdateKey_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing From
	opt := &ProjectsUpdateKeyOption{
		To: "new-key",
	}
	_, err := client.Projects.UpdateKey(opt)
	if err == nil {
		t.Error("expected validation error for missing From")
	}

	// Test missing To
	opt = &ProjectsUpdateKeyOption{
		From: "old-key",
	}
	_, err = client.Projects.UpdateKey(opt)
	if err == nil {
		t.Error("expected validation error for missing To")
	}
}

// TestProjectsService_UpdateVisibility tests the UpdateVisibility method.
func TestProjectsService_UpdateVisibility(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/update_visibility") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &ProjectsUpdateVisibilityOption{
		Project:    "my-project",
		Visibility: "public",
	}

	resp, err := client.Projects.UpdateVisibility(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestProjectsService_UpdateVisibility_ValidationError tests validation for UpdateVisibility.
func TestProjectsService_UpdateVisibility_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Project
	opt := &ProjectsUpdateVisibilityOption{
		Visibility: "public",
	}
	_, err := client.Projects.UpdateVisibility(opt)
	if err == nil {
		t.Error("expected validation error for missing Project")
	}

	// Test missing Visibility
	opt = &ProjectsUpdateVisibilityOption{
		Project: "my-project",
	}
	_, err = client.Projects.UpdateVisibility(opt)
	if err == nil {
		t.Error("expected validation error for missing Visibility")
	}

	// Test invalid visibility
	opt = &ProjectsUpdateVisibilityOption{
		Project:    "my-project",
		Visibility: "invalid",
	}
	_, err = client.Projects.UpdateVisibility(opt)
	if err == nil {
		t.Error("expected validation error for invalid Visibility")
	}
}
