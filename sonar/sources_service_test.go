package sonargo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// SourcesService Test Suite
// -----------------------------------------------------------------------------

// TestSourcesService_Index tests the Index method.
func TestSourcesService_Index(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/sources/index") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("key") != "my-project:src/main.go" {
			t.Errorf("unexpected key: %s", r.URL.Query().Get("key"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"sources":{"1":"package main","2":"","3":"import \"fmt\""}}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SourcesIndexOption{
		Key:  "my-project:src/main.go",
		From: 1,
		To:   10,
	}

	result, resp, err := client.Sources.Index(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.Sources["1"] != "package main" {
		t.Errorf("unexpected source line 1: %s", result.Sources["1"])
	}
}

// TestSourcesService_Index_ValidationError tests validation for Index.
func TestSourcesService_Index_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Key
	opt := &SourcesIndexOption{}
	_, _, err := client.Sources.Index(opt)
	if err == nil {
		t.Error("expected validation error for missing Key")
	}
}

// TestSourcesService_IssueSnippets tests the IssueSnippets method.
func TestSourcesService_IssueSnippets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/sources/issue_snippets") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"my-project:src/main.go": {
				"component": {
					"key": "my-project:src/main.go",
					"name": "main.go",
					"qualifier": "FIL"
				},
				"sources": [
					{"line": 10, "code": "func main() {"}
				]
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SourcesIssueSnippetsOption{
		IssueKey: "AX1234567890",
	}

	result, resp, err := client.Sources.IssueSnippets(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	snippet, ok := (*result)["my-project:src/main.go"]
	if !ok {
		t.Error("expected snippet for my-project:src/main.go")
	}
	if snippet.Component.Key != "my-project:src/main.go" {
		t.Errorf("unexpected component key: %s", snippet.Component.Key)
	}
}

// TestSourcesService_IssueSnippets_ValidationError tests validation for IssueSnippets.
func TestSourcesService_IssueSnippets_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing IssueKey
	opt := &SourcesIssueSnippetsOption{}
	_, _, err := client.Sources.IssueSnippets(opt)
	if err == nil {
		t.Error("expected validation error for missing IssueKey")
	}
}

// TestSourcesService_Lines tests the Lines method.
func TestSourcesService_Lines(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/sources/lines") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"sources": [
				{
					"line": 1,
					"code": "<span class=\"k\">package</span> main",
					"scmAuthor": "john.doe@example.com",
					"scmDate": "2024-01-15T10:30:00+0000",
					"scmRevision": "abc123",
					"duplicated": false
				}
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SourcesLinesOption{
		Key:    "my-project:src/main.go",
		Branch: "main",
		From:   1,
		To:     100,
	}

	result, resp, err := client.Sources.Lines(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Sources) != 1 {
		t.Errorf("expected 1 source line, got %d", len(result.Sources))
	}
	if result.Sources[0].Line != 1 {
		t.Errorf("unexpected line number: %d", result.Sources[0].Line)
	}
	if result.Sources[0].SCMAuthor != "john.doe@example.com" {
		t.Errorf("unexpected author: %s", result.Sources[0].SCMAuthor)
	}
}

// TestSourcesService_Lines_ValidationError tests validation for Lines.
func TestSourcesService_Lines_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Key
	opt := &SourcesLinesOption{
		Branch: "main",
	}
	_, _, err := client.Sources.Lines(opt)
	if err == nil {
		t.Error("expected validation error for missing Key")
	}
}

// TestSourcesService_Raw tests the Raw method.
func TestSourcesService_Raw(t *testing.T) {
	expectedContent := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/sources/raw") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SourcesRawOption{
		Key: "my-project:src/main.go",
	}

	result, resp, err := client.Sources.Raw(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result != expectedContent {
		t.Errorf("unexpected raw content: %s", result)
	}
}

// TestSourcesService_Raw_ValidationError tests validation for Raw.
func TestSourcesService_Raw_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Key
	opt := &SourcesRawOption{}
	_, _, err := client.Sources.Raw(opt)
	if err == nil {
		t.Error("expected validation error for missing Key")
	}
}

// TestSourcesService_Scm tests the Scm method.
func TestSourcesService_Scm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/sources/scm") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"scm": [
				[1, "john.doe@example.com", "2024-01-15T10:30:00+0000", "abc123"],
				[2, "jane.smith@example.com", "2024-01-14T09:00:00+0000", "def456"]
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SourcesScmOption{
		Key:           "my-project:src/main.go",
		CommitsByLine: true,
	}

	result, resp, err := client.Sources.Scm(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Scm) != 2 {
		t.Errorf("expected 2 SCM entries, got %d", len(result.Scm))
	}
}

// TestSourcesService_Scm_ValidationError tests validation for Scm.
func TestSourcesService_Scm_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Key
	opt := &SourcesScmOption{
		CommitsByLine: true,
	}
	_, _, err := client.Sources.Scm(opt)
	if err == nil {
		t.Error("expected validation error for missing Key")
	}
}

// TestSourcesService_Show tests the Show method.
func TestSourcesService_Show(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/sources/show") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"sources": [
				[1, "package main"],
				[2, ""],
				[3, "import \"fmt\""]
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SourcesShowOption{
		Key:  "my-project:src/main.go",
		From: 1,
		To:   10,
	}

	result, resp, err := client.Sources.Show(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Sources) != 3 {
		t.Errorf("expected 3 source lines, got %d", len(result.Sources))
	}
}

// TestSourcesService_Show_ValidationError tests validation for Show.
func TestSourcesService_Show_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Key
	opt := &SourcesShowOption{}
	_, _, err := client.Sources.Show(opt)
	if err == nil {
		t.Error("expected validation error for missing Key")
	}
}
