package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// AddComment Tests
// =============================================================================

func TestHotspots_AddComment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/add_comment" {
			t.Errorf("expected path /api/hotspots/add_comment, got %s", r.URL.Path)
		}

		hotspot := r.URL.Query().Get("hotspot")
		if hotspot != "hotspot123" {
			t.Errorf("expected hotspot 'hotspot123', got %s", hotspot)
		}

		comment := r.URL.Query().Get("comment")
		if comment != "This is a comment" {
			t.Errorf("expected comment 'This is a comment', got %s", comment)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Hotspots.AddComment(&HotspotsAddCommentOption{
		Hotspot: "hotspot123",
		Comment: "This is a comment",
	})
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestHotspots_AddComment_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsAddCommentOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsAddCommentOption{Comment: "This is a comment"},
		},
		{
			name: "missing comment",
			opt:  &HotspotsAddCommentOption{Hotspot: "hotspot123"},
		},
		{
			name: "comment too long",
			opt: &HotspotsAddCommentOption{
				Hotspot: "hotspot123",
				Comment: string(make([]byte, MaxHotspotCommentLength+1)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.AddComment(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Assign Tests
// =============================================================================

func TestHotspots_Assign(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/assign" {
			t.Errorf("expected path /api/hotspots/assign, got %s", r.URL.Path)
		}

		hotspot := r.URL.Query().Get("hotspot")
		if hotspot != "hotspot123" {
			t.Errorf("expected hotspot 'hotspot123', got %s", hotspot)
		}

		assignee := r.URL.Query().Get("assignee")
		if assignee != "john.doe" {
			t.Errorf("expected assignee 'john.doe', got %s", assignee)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Hotspots.Assign(&HotspotsAssignOption{
		Hotspot:  "hotspot123",
		Assignee: "john.doe",
	})
	if err != nil {
		t.Fatalf("Assign failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestHotspots_Assign_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsAssignOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsAssignOption{Assignee: "john.doe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.Assign(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// ChangeStatus Tests
// =============================================================================

func TestHotspots_ChangeStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/change_status" {
			t.Errorf("expected path /api/hotspots/change_status, got %s", r.URL.Path)
		}

		hotspot := r.URL.Query().Get("hotspot")
		if hotspot != "hotspot123" {
			t.Errorf("expected hotspot 'hotspot123', got %s", hotspot)
		}

		status := r.URL.Query().Get("status")
		if status != "REVIEWED" {
			t.Errorf("expected status 'REVIEWED', got %s", status)
		}

		resolution := r.URL.Query().Get("resolution")
		if resolution != "SAFE" {
			t.Errorf("expected resolution 'SAFE', got %s", resolution)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Hotspots.ChangeStatus(&HotspotsChangeStatusOption{
		Hotspot:    "hotspot123",
		Status:     "REVIEWED",
		Resolution: "SAFE",
	})
	if err != nil {
		t.Fatalf("ChangeStatus failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestHotspots_ChangeStatus_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsChangeStatusOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsChangeStatusOption{Status: "REVIEWED"},
		},
		{
			name: "missing status",
			opt:  &HotspotsChangeStatusOption{Hotspot: "hotspot123"},
		},
		{
			name: "invalid status",
			opt:  &HotspotsChangeStatusOption{Hotspot: "hotspot123", Status: "INVALID"},
		},
		{
			name: "invalid resolution",
			opt:  &HotspotsChangeStatusOption{Hotspot: "hotspot123", Status: "REVIEWED", Resolution: "INVALID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.ChangeStatus(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// DeleteComment Tests
// =============================================================================

func TestHotspots_DeleteComment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/delete_comment" {
			t.Errorf("expected path /api/hotspots/delete_comment, got %s", r.URL.Path)
		}

		comment := r.URL.Query().Get("comment")
		if comment != "comment123" {
			t.Errorf("expected comment 'comment123', got %s", comment)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Hotspots.DeleteComment(&HotspotsDeleteCommentOption{
		Comment: "comment123",
	})
	if err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestHotspots_DeleteComment_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsDeleteCommentOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing comment",
			opt:  &HotspotsDeleteCommentOption{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.DeleteComment(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// EditComment Tests
// =============================================================================

func TestHotspots_EditComment(t *testing.T) {
	expectedResponse := &HotspotsEditComment{
		CreatedAt: "2024-01-01T12:00:00+0000",
		HTMLText:  "<p>Updated comment</p>",
		Key:       "comment123",
		Login:     "john.doe",
		Markdown:  "Updated comment",
		Updatable: true,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/edit_comment" {
			t.Errorf("expected path /api/hotspots/edit_comment, got %s", r.URL.Path)
		}

		comment := r.URL.Query().Get("comment")
		if comment != "comment123" {
			t.Errorf("expected comment 'comment123', got %s", comment)
		}

		text := r.URL.Query().Get("text")
		if text != "Updated comment" {
			t.Errorf("expected text 'Updated comment', got %s", text)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Hotspots.EditComment(&HotspotsEditCommentOption{
		Comment: "comment123",
		Text:    "Updated comment",
	})
	if err != nil {
		t.Fatalf("EditComment failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Key != expectedResponse.Key {
		t.Errorf("expected key '%s', got %s", expectedResponse.Key, result.Key)
	}

	if result.Markdown != expectedResponse.Markdown {
		t.Errorf("expected markdown '%s', got %s", expectedResponse.Markdown, result.Markdown)
	}
}

func TestHotspots_EditComment_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsEditCommentOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing comment",
			opt:  &HotspotsEditCommentOption{Text: "Updated comment"},
		},
		{
			name: "missing text",
			opt:  &HotspotsEditCommentOption{Comment: "comment123"},
		},
		{
			name: "text too long",
			opt: &HotspotsEditCommentOption{
				Comment: "comment123",
				Text:    string(make([]byte, MaxHotspotCommentLength+1)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.EditComment(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// List Tests
// =============================================================================

func TestHotspots_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/list" {
			t.Errorf("expected path /api/hotspots/list, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [
				{"key": "project:src/main.go", "name": "main.go", "qualifier": "FIL"}
			],
			"hotspots": [
				{
					"key": "hotspot123",
					"component": "project:src/main.go",
					"status": "TO_REVIEW",
					"vulnerabilityProbability": "HIGH"
				}
			],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Hotspots.List(&HotspotsListOption{
		Project: "my-project",
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Hotspots) != 1 {
		t.Fatalf("expected 1 hotspot, got %d", len(result.Hotspots))
	}

	if result.Hotspots[0].Key != "hotspot123" {
		t.Errorf("expected key 'hotspot123', got %s", result.Hotspots[0].Key)
	}
}

func TestHotspots_List_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsListOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project",
			opt:  &HotspotsListOption{},
		},
		{
			name: "invalid status",
			opt:  &HotspotsListOption{Project: "my-project", Status: "INVALID"},
		},
		{
			name: "invalid resolution",
			opt:  &HotspotsListOption{Project: "my-project", Resolution: "INVALID"},
		},
		{
			name: "page size too large",
			opt:  &HotspotsListOption{Project: "my-project", PaginationArgs: PaginationArgs{PageSize: MaxHotspotListPageSize + 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.List(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Pull Tests
// =============================================================================

func TestHotspots_Pull(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/pull" {
			t.Errorf("expected path /api/hotspots/pull, got %s", r.URL.Path)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		branchName := r.URL.Query().Get("branchName")
		if branchName != "main" {
			t.Errorf("expected branchName 'main', got %s", branchName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Hotspots.Pull(&HotspotsPullOption{
		ProjectKey: "my-project",
		BranchName: "main",
	})
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestHotspots_Pull_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsPullOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project key",
			opt:  &HotspotsPullOption{BranchName: "main"},
		},
		{
			name: "missing branch name",
			opt:  &HotspotsPullOption{ProjectKey: "my-project"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.Pull(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Search Tests
// =============================================================================

func TestHotspots_Search_WithProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/search" {
			t.Errorf("expected path /api/hotspots/search, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [
				{"key": "project:src/main.go", "name": "main.go", "qualifier": "FIL"}
			],
			"hotspots": [
				{
					"key": "hotspot123",
					"component": "project:src/main.go",
					"status": "TO_REVIEW",
					"vulnerabilityProbability": "HIGH"
				}
			],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Hotspots.Search(&HotspotsSearchOption{
		Project: "my-project",
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Hotspots) != 1 {
		t.Fatalf("expected 1 hotspot, got %d", len(result.Hotspots))
	}

	if result.Hotspots[0].Key != "hotspot123" {
		t.Errorf("expected key 'hotspot123', got %s", result.Hotspots[0].Key)
	}
}

func TestHotspots_Search_WithHotspots(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		hotspots := r.URL.Query().Get("hotspots")
		if hotspots != "hotspot1,hotspot2" {
			t.Errorf("expected hotspots 'hotspot1,hotspot2', got %s", hotspots)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"components": [], "hotspots": [], "paging": {"pageIndex": 1, "pageSize": 100, "total": 0}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, resp, err := client.Hotspots.Search(&HotspotsSearchOption{
		Hotspots: []string{"hotspot1", "hotspot2"},
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestHotspots_Search_WithFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		status := r.URL.Query().Get("status")
		if status != "REVIEWED" {
			t.Errorf("expected status 'REVIEWED', got %s", status)
		}

		resolution := r.URL.Query().Get("resolution")
		if resolution != "SAFE" {
			t.Errorf("expected resolution 'SAFE', got %s", resolution)
		}

		inNewCodePeriod := r.URL.Query().Get("inNewCodePeriod")
		if inNewCodePeriod != "true" {
			t.Errorf("expected inNewCodePeriod 'true', got %s", inNewCodePeriod)
		}

		onlyMine := r.URL.Query().Get("onlyMine")
		if onlyMine != "true" {
			t.Errorf("expected onlyMine 'true', got %s", onlyMine)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"components": [], "hotspots": [], "paging": {"pageIndex": 1, "pageSize": 100, "total": 0}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, resp, err := client.Hotspots.Search(&HotspotsSearchOption{
		Project:         "my-project",
		Status:          "REVIEWED",
		Resolution:      "SAFE",
		InNewCodePeriod: true,
		OnlyMine:        true,
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestHotspots_Search_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsSearchOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project and hotspots",
			opt:  &HotspotsSearchOption{},
		},
		{
			name: "invalid status",
			opt:  &HotspotsSearchOption{Project: "my-project", Status: "INVALID"},
		},
		{
			name: "invalid resolution",
			opt:  &HotspotsSearchOption{Project: "my-project", Resolution: "INVALID"},
		},
		{
			name: "invalid owasp asvs level",
			opt:  &HotspotsSearchOption{Project: "my-project", OwaspAsvsLevel: "5"},
		},
		{
			name: "invalid owasp top 10",
			opt:  &HotspotsSearchOption{Project: "my-project", OwaspTop10: []string{"a1", "invalid"}},
		},
		{
			name: "invalid sans top 25",
			opt:  &HotspotsSearchOption{Project: "my-project", SansTop25: []string{"insecure-interaction", "invalid"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.Search(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Show Tests
// =============================================================================

func TestHotspots_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/hotspots/show" {
			t.Errorf("expected path /api/hotspots/show, got %s", r.URL.Path)
		}

		hotspot := r.URL.Query().Get("hotspot")
		if hotspot != "hotspot123" {
			t.Errorf("expected hotspot 'hotspot123', got %s", hotspot)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"key": "hotspot123",
			"status": "TO_REVIEW",
			"canChangeStatus": true,
			"message": "Security issue found",
			"component": {
				"key": "project:src/main.go",
				"name": "main.go"
			},
			"project": {
				"key": "my-project",
				"name": "My Project"
			},
			"rule": {
				"key": "java:S2092",
				"name": "Cookies should be secure"
			},
			"users": [
				{"login": "john.doe", "name": "John Doe", "active": true}
			],
			"changelog": [
				{
					"user": "john.doe",
					"creationDate": "2024-01-01T12:00:00+0000",
					"diffs": [
						{"key": "status", "oldValue": "TO_REVIEW", "newValue": "REVIEWED"}
					]
				}
			],
			"comment": [
				{
					"key": "comment123",
					"login": "john.doe",
					"markdown": "This is safe"
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Hotspots.Show(&HotspotsShowOption{
		Hotspot: "hotspot123",
	})
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Key != "hotspot123" {
		t.Errorf("expected key 'hotspot123', got %s", result.Key)
	}

	if result.Status != "TO_REVIEW" {
		t.Errorf("expected status 'TO_REVIEW', got %s", result.Status)
	}

	if !result.CanChangeStatus {
		t.Error("expected canChangeStatus to be true")
	}

	if len(result.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(result.Users))
	}

	if result.Users[0].Login != "john.doe" {
		t.Errorf("expected user login 'john.doe', got %s", result.Users[0].Login)
	}

	if len(result.Changelog) != 1 {
		t.Fatalf("expected 1 changelog entry, got %d", len(result.Changelog))
	}

	if result.Changelog[0].User != "john.doe" {
		t.Errorf("expected changelog user 'john.doe', got %s", result.Changelog[0].User)
	}

	if len(result.Changelog[0].Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result.Changelog[0].Diffs))
	}

	if result.Changelog[0].Diffs[0].Key != "status" {
		t.Errorf("expected diff key 'status', got %s", result.Changelog[0].Diffs[0].Key)
	}
}

func TestHotspots_Show_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name string
		opt  *HotspotsShowOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsShowOption{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.Show(tt.opt)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Validation Function Tests
// =============================================================================

func TestHotspots_ValidateAddCommentOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *HotspotsAddCommentOption
		wantErr bool
	}{
		{
			name: "valid option",
			opt: &HotspotsAddCommentOption{
				Comment: "Valid comment",
				Hotspot: "hotspot123",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name: "missing comment",
			opt: &HotspotsAddCommentOption{
				Hotspot: "hotspot123",
			},
			wantErr: true,
		},
		{
			name: "missing hotspot",
			opt: &HotspotsAddCommentOption{
				Comment: "Valid comment",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateAddCommentOpt(tt.opt)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestHotspots_ValidateChangeStatusOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *HotspotsChangeStatusOption
		wantErr bool
	}{
		{
			name: "valid TO_REVIEW status",
			opt: &HotspotsChangeStatusOption{
				Hotspot: "hotspot123",
				Status:  "TO_REVIEW",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with SAFE resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "SAFE",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with FIXED resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "FIXED",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with ACKNOWLEDGED resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "ACKNOWLEDGED",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name: "invalid status",
			opt: &HotspotsChangeStatusOption{
				Hotspot: "hotspot123",
				Status:  "INVALID_STATUS",
			},
			wantErr: true,
		},
		{
			name: "invalid resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "INVALID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateChangeStatusOpt(tt.opt)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestHotspots_ValidateSearchOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *HotspotsSearchOption
		wantErr bool
	}{
		{
			name: "valid with project",
			opt: &HotspotsSearchOption{
				Project: "my-project",
			},
			wantErr: false,
		},
		{
			name: "valid with hotspots",
			opt: &HotspotsSearchOption{
				Hotspots: []string{"hotspot1", "hotspot2"},
			},
			wantErr: false,
		},
		{
			name: "valid with all OWASP filters",
			opt: &HotspotsSearchOption{
				Project:        "my-project",
				OwaspTop10:     []string{"a1", "a2"},
				OwaspTop102021: []string{"a3", "a4"},
				OwaspAsvsLevel: "2",
			},
			wantErr: false,
		},
		{
			name: "valid with SANS filter",
			opt: &HotspotsSearchOption{
				Project:   "my-project",
				SansTop25: []string{"insecure-interaction", "porous-defenses"},
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing project and hotspots",
			opt:     &HotspotsSearchOption{},
			wantErr: true,
		},
		{
			name: "invalid OwaspAsvsLevel",
			opt: &HotspotsSearchOption{
				Project:        "my-project",
				OwaspAsvsLevel: "4",
			},
			wantErr: true,
		},
		{
			name: "invalid OwaspTop10 value",
			opt: &HotspotsSearchOption{
				Project:    "my-project",
				OwaspTop10: []string{"a11"},
			},
			wantErr: true,
		},
		{
			name: "invalid SansTop25 value",
			opt: &HotspotsSearchOption{
				Project:   "my-project",
				SansTop25: []string{"invalid-category"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateSearchOpt(tt.opt)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestHotspots_ValidateListOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *HotspotsListOption
		wantErr bool
	}{
		{
			name: "valid basic option",
			opt: &HotspotsListOption{
				Project: "my-project",
			},
			wantErr: false,
		},
		{
			name: "valid with all optional params",
			opt: &HotspotsListOption{
				Project:         "my-project",
				Branch:          "main",
				InNewCodePeriod: true,
				Status:          "TO_REVIEW",
				Resolution:      "SAFE",
				PaginationArgs:  PaginationArgs{PageSize: 100},
			},
			wantErr: false,
		},
		{
			name: "valid max page size",
			opt: &HotspotsListOption{
				Project:        "my-project",
				PaginationArgs: PaginationArgs{PageSize: MaxHotspotListPageSize},
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing project",
			opt:     &HotspotsListOption{},
			wantErr: true,
		},
		{
			name: "page size exceeds max",
			opt: &HotspotsListOption{
				Project:        "my-project",
				PaginationArgs: PaginationArgs{PageSize: MaxHotspotListPageSize + 1},
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			opt: &HotspotsListOption{
				Project: "my-project",
				Status:  "CLOSED",
			},
			wantErr: true,
		},
		{
			name: "invalid resolution",
			opt: &HotspotsListOption{
				Project:    "my-project",
				Resolution: "WONTFIX",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateListOpt(tt.opt)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
