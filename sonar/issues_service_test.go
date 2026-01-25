package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIssues_AddComment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/issues/add_comment" {
			t.Errorf("expected path /api/issues/add_comment, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"issue": {"key": "AU-Tpxb--iU5OvuD2FLy", "message": "Test issue"},
			"components": [],
			"rules": [],
			"users": []
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesAddCommentOption{
		Issue: "AU-Tpxb--iU5OvuD2FLy",
		Text:  "This is a comment",
	}
	result, resp, err := client.Issues.AddComment(opt)
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Key != "AU-Tpxb--iU5OvuD2FLy" {
		t.Errorf("expected issue key AU-Tpxb--iU5OvuD2FLy, got %s", result.Issue.Key)
	}
}

func TestIssues_AddComment_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test nil option
	_, _, err := client.Issues.AddComment(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Issue
	_, _, err = client.Issues.AddComment(&IssuesAddCommentOption{Text: "test"})
	if err == nil {
		t.Error("expected error for missing Issue")
	}

	// Test missing Text
	_, _, err = client.Issues.AddComment(&IssuesAddCommentOption{Issue: "key"})
	if err == nil {
		t.Error("expected error for missing Text")
	}
}

func TestIssues_Assign(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"issue": {"key": "test-key", "assignee": "admin"},
			"components": [],
			"rules": [],
			"users": []
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesAssignOption{
		Issue:    "test-key",
		Assignee: "admin",
	}
	result, resp, err := client.Issues.Assign(opt)
	if err != nil {
		t.Fatalf("Assign failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Assignee != "admin" {
		t.Errorf("expected assignee admin, got %s", result.Issue.Assignee)
	}
}

func TestIssues_Authors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"authors": ["john@example.com", "jane@example.com"]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesAuthorsOption{
		Project:  "my-project",
		PageSize: 50,
	}
	result, resp, err := client.Issues.Authors(opt)
	if err != nil {
		t.Fatalf("Authors failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(result.Authors))
	}
}

func TestIssues_Authors_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test page size out of range
	_, _, err := client.Issues.Authors(&IssuesAuthorsOption{PageSize: 150})
	if err == nil {
		t.Error("expected error for page size > 100")
	}
}

func TestIssues_BulkChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"total": 10, "success": 8, "ignored": 1, "failures": 1}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesBulkChangeOption{
		Issues:      []string{"issue1", "issue2"},
		SetSeverity: "MAJOR",
		AddTags:     []string{"tag1", "tag2"},
	}
	result, resp, err := client.Issues.BulkChange(opt)
	if err != nil {
		t.Fatalf("BulkChange failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Success != 8 {
		t.Errorf("expected 8 successes, got %d", result.Success)
	}
}

func TestIssues_BulkChange_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test nil option
	_, _, err := client.Issues.BulkChange(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test empty issues
	_, _, err = client.Issues.BulkChange(&IssuesBulkChangeOption{})
	if err == nil {
		t.Error("expected error for empty issues")
	}

	// Test invalid severity
	_, _, err = client.Issues.BulkChange(&IssuesBulkChangeOption{
		Issues:      []string{"issue1"},
		SetSeverity: "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid severity")
	}

	// Test invalid type
	_, _, err = client.Issues.BulkChange(&IssuesBulkChangeOption{
		Issues:  []string{"issue1"},
		SetType: "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}

	// Test invalid transition
	_, _, err = client.Issues.BulkChange(&IssuesBulkChangeOption{
		Issues:       []string{"issue1"},
		DoTransition: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}

func TestIssues_Changelog(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"changelog": [
				{
					"user": "admin",
					"userName": "Admin User",
					"creationDate": "2023-01-01T00:00:00+0000",
					"diffs": [
						{"key": "severity", "oldValue": "MAJOR", "newValue": "CRITICAL"}
					]
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesChangelogOption{Issue: "test-key"}
	result, resp, err := client.Issues.Changelog(opt)
	if err != nil {
		t.Fatalf("Changelog failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Changelog) != 1 {
		t.Errorf("expected 1 changelog entry, got %d", len(result.Changelog))
	}
}

func TestIssues_ComponentTags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tags": [{"key": "security", "value": 10}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesComponentTagsOption{ComponentUuid: "uuid-123"}
	result, resp, err := client.Issues.ComponentTags(opt)
	if err != nil {
		t.Fatalf("ComponentTags failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result.Tags))
	}
}

func TestIssues_DeleteComment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issue": {"key": "test-key"}, "components": [], "rules": [], "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesDeleteCommentOption{Comment: "comment-key"}
	result, resp, err := client.Issues.DeleteComment(opt)
	if err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Key != "test-key" {
		t.Error("unexpected issue key")
	}
}

func TestIssues_DoTransition(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issue": {"key": "test-key"}, "components": [], "rules": [], "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesDoTransitionOption{
		Issue:      "test-key",
		Transition: "confirm",
	}
	result, resp, err := client.Issues.DoTransition(opt)
	if err != nil {
		t.Fatalf("DoTransition failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Key != "test-key" {
		t.Error("unexpected issue key")
	}
}

func TestIssues_DoTransition_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test invalid transition
	_, _, err := client.Issues.DoTransition(&IssuesDoTransitionOption{
		Issue:      "test-key",
		Transition: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}

func TestIssues_EditComment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issue": {"key": "test-key"}, "components": [], "rules": [], "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesEditCommentOption{
		Comment: "comment-key",
		Text:    "Updated comment",
	}
	result, resp, err := client.Issues.EditComment(opt)
	if err != nil {
		t.Fatalf("EditComment failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Key != "test-key" {
		t.Error("unexpected issue key")
	}
}

func TestIssues_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"issues": [{"key": "issue-1"}],
			"components": [],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesListOption{
		Project: "my-project",
	}
	result, resp, err := client.Issues.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(result.Issues))
	}
}

func TestIssues_List_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test nil option
	_, _, err := client.Issues.List(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing project and component
	_, _, err = client.Issues.List(&IssuesListOption{})
	if err == nil {
		t.Error("expected error for missing project and component")
	}

	// Test invalid type
	_, _, err = client.Issues.List(&IssuesListOption{
		Project: "my-project",
		Types:   []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestIssues_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"issues": [{"key": "issue-1", "message": "Test issue"}],
			"components": [],
			"rules": [],
			"users": [],
			"facets": [],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesSearchOption{
		Projects:         []string{"my-project"},
		Severities:       []string{"BLOCKER", "CRITICAL"},
		Types:            []string{"BUG"},
		ImpactSeverities: []string{"HIGH"},
	}
	result, resp, err := client.Issues.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(result.Issues))
	}
}

func TestIssues_Search_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test invalid severity
	_, _, err := client.Issues.Search(&IssuesSearchOption{
		Severities: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid severity")
	}

	// Test invalid type
	_, _, err = client.Issues.Search(&IssuesSearchOption{
		Types: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}

	// Test invalid impact severity
	_, _, err = client.Issues.Search(&IssuesSearchOption{
		ImpactSeverities: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid impact severity")
	}

	// Test invalid impact software quality
	_, _, err = client.Issues.Search(&IssuesSearchOption{
		ImpactSoftwareQualities: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid impact software quality")
	}

	// Test invalid clean code attribute category
	_, _, err = client.Issues.Search(&IssuesSearchOption{
		CleanCodeAttributeCategories: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid clean code attribute category")
	}
}

func TestIssues_SetSeverity(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issue": {"key": "test-key", "severity": "BLOCKER"}, "components": [], "rules": [], "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesSetSeverityOption{
		Issue:    "test-key",
		Severity: "BLOCKER",
	}
	result, resp, err := client.Issues.SetSeverity(opt)
	if err != nil {
		t.Fatalf("SetSeverity failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Severity != "BLOCKER" {
		t.Errorf("expected severity BLOCKER, got %s", result.Issue.Severity)
	}
}

func TestIssues_SetSeverity_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test invalid severity
	_, _, err := client.Issues.SetSeverity(&IssuesSetSeverityOption{
		Issue:    "test-key",
		Severity: "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid severity")
	}
}

func TestIssues_SetTags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issue": {"key": "test-key", "tags": ["security", "cwe"]}, "components": [], "rules": [], "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesSetTagsOption{
		Issue: "test-key",
		Tags:  []string{"security", "cwe"},
	}
	result, resp, err := client.Issues.SetTags(opt)
	if err != nil {
		t.Fatalf("SetTags failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Issue.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Issue.Tags))
	}
}

func TestIssues_SetType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"issue": {"key": "test-key", "type": "BUG"}, "components": [], "rules": [], "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesSetTypeOption{
		Issue: "test-key",
		Type:  "BUG",
	}
	result, resp, err := client.Issues.SetType(opt)
	if err != nil {
		t.Fatalf("SetType failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Issue.Type != "BUG" {
		t.Errorf("expected type BUG, got %s", result.Issue.Type)
	}
}

func TestIssues_SetType_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing issue
	_, _, err := client.Issues.SetType(&IssuesSetTypeOption{Type: "BUG"})
	if err == nil {
		t.Error("expected error for missing issue")
	}

	// Test missing type
	_, _, err = client.Issues.SetType(&IssuesSetTypeOption{Issue: "test-key"})
	if err == nil {
		t.Error("expected error for missing type")
	}

	// Test invalid type
	_, _, err = client.Issues.SetType(&IssuesSetTypeOption{
		Issue: "test-key",
		Type:  "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestIssues_Tags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tags": ["security", "cwe", "java"]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesTagsOption{
		Project:  "my-project",
		PageSize: 100,
	}
	result, resp, err := client.Issues.Tags(opt)
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(result.Tags))
	}
}

func TestIssues_Pull(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesPullOption{
		ProjectKey: "my-project",
		BranchName: "main",
	}
	result, resp, err := client.Issues.Pull(opt)
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestIssues_Pull_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test nil option
	_, _, err := client.Issues.Pull(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing project key
	_, _, err = client.Issues.Pull(&IssuesPullOption{})
	if err == nil {
		t.Error("expected error for missing project key")
	}

	// Test invalid language
	_, _, err = client.Issues.Pull(&IssuesPullOption{
		ProjectKey: "my-project",
		Languages:  []string{"invalid-language"},
	})
	if err == nil {
		t.Error("expected error for invalid language")
	}
}

func TestIssues_PullTaint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesPullTaintOption{
		ProjectKey: "my-project",
		BranchName: "main",
	}
	result, resp, err := client.Issues.PullTaint(opt)
	if err != nil {
		t.Fatalf("PullTaint failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestIssues_Reindex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesReindexOption{Project: "my-project"}
	resp, err := client.Issues.Reindex(opt)
	if err != nil {
		t.Fatalf("Reindex failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestIssues_AnticipatedTransitions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &IssuesAnticipatedTransitionsOption{ProjectKey: "my-project"}
	resp, err := client.Issues.AnticipatedTransitions(opt)
	if err != nil {
		t.Fatalf("AnticipatedTransitions failed: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", resp.StatusCode)
	}
}

func TestIssues_AnticipatedTransitions_Validation(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test nil option
	_, err := client.Issues.AnticipatedTransitions(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing project key
	_, err = client.Issues.AnticipatedTransitions(&IssuesAnticipatedTransitionsOption{})
	if err == nil {
		t.Error("expected error for missing project key")
	}
}
