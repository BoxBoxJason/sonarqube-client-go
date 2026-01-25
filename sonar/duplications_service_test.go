package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDuplications_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/duplications/show" {
			t.Errorf("expected path /api/duplications/show, got %s", r.URL.Path)
		}

		key := r.URL.Query().Get("key")
		if key != "com.example:MyFile.java" {
			t.Errorf("expected key 'com.example:MyFile.java', got %s", key)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"duplications": [
				{
					"blocks": [
						{"_ref": "1", "from": 10, "size": 5},
						{"_ref": "2", "from": 20, "size": 5}
					]
				}
			],
			"files": {
				"1": {"key": "com.example:MyFile.java", "name": "MyFile.java", "projectName": "My Project"},
				"2": {"key": "com.example:OtherFile.java", "name": "OtherFile.java", "projectName": "My Project"}
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DuplicationsShowOption{
		Key: "com.example:MyFile.java",
	}

	result, resp, err := client.Duplications.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Duplications) != 1 {
		t.Errorf("expected 1 duplication group, got %d", len(result.Duplications))
	}

	if len(result.Duplications[0].Blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(result.Duplications[0].Blocks))
	}

	if result.Duplications[0].Blocks[0].Ref != "1" {
		t.Errorf("expected first block ref '1', got %s", result.Duplications[0].Blocks[0].Ref)
	}

	if result.Duplications[0].Blocks[0].From != 10 {
		t.Errorf("expected first block from 10, got %d", result.Duplications[0].Blocks[0].From)
	}

	if len(result.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(result.Files))
	}

	if file, ok := result.Files["1"]; !ok {
		t.Error("expected file with key '1'")
	} else if file.Name != "MyFile.java" {
		t.Errorf("expected file name 'MyFile.java', got %s", file.Name)
	}
}

func TestDuplications_Show_WithBranch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		branch := r.URL.Query().Get("branch")
		if branch != "feature" {
			t.Errorf("expected branch 'feature', got %s", branch)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"duplications": [], "files": {}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DuplicationsShowOption{
		Key:    "com.example:MyFile.java",
		Branch: "feature",
	}

	_, resp, err := client.Duplications.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDuplications_Show_WithPullRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pullRequest := r.URL.Query().Get("pullRequest")
		if pullRequest != "123" {
			t.Errorf("expected pullRequest '123', got %s", pullRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"duplications": [], "files": {}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DuplicationsShowOption{
		Key:         "com.example:MyFile.java",
		PullRequest: "123",
	}

	_, resp, err := client.Duplications.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDuplications_Show_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Duplications.Show(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Key should fail validation.
	_, _, err = client.Duplications.Show(&DuplicationsShowOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}
}

func TestDuplications_ValidateShowOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.Duplications.ValidateShowOpt(&DuplicationsShowOption{
		Key: "com.example:MyFile.java",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.Duplications.ValidateShowOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Key should fail.
	err = client.Duplications.ValidateShowOpt(&DuplicationsShowOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}
}
