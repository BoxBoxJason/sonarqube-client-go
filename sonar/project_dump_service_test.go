package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectDump_Export(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_dump/export" {
			t.Errorf("expected path /api/project_dump/export, got %s", r.URL.Path)
		}

		key := r.URL.Query().Get("key")
		if key != "my-project" {
			t.Errorf("expected key 'my-project', got %s", key)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"projectId": "proj-123",
			"projectKey": "my-project",
			"projectName": "My Project",
			"taskId": "task-456"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectDumpExportOption{
		Key: "my-project",
	}

	result, resp, err := client.ProjectDump.Export(opt)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.ProjectID != "proj-123" {
		t.Errorf("expected projectId 'proj-123', got %s", result.ProjectID)
	}

	if result.ProjectKey != "my-project" {
		t.Errorf("expected projectKey 'my-project', got %s", result.ProjectKey)
	}

	if result.TaskID != "task-456" {
		t.Errorf("expected taskId 'task-456', got %s", result.TaskID)
	}
}

func TestProjectDump_Export_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectDump.Export(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Key should fail validation.
	_, _, err = client.ProjectDump.Export(&ProjectDumpExportOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}
}

func TestProjectDump_Status(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_dump/status" {
			t.Errorf("expected path /api/project_dump/status, got %s", r.URL.Path)
		}

		key := r.URL.Query().Get("key")
		if key != "my-project" {
			t.Errorf("expected key 'my-project', got %s", key)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"canBeExported": true,
			"canBeImported": false,
			"dumpToImport": "",
			"exportedDump": "/path/to/dump.zip"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectDumpStatusOption{
		Key: "my-project",
	}

	result, resp, err := client.ProjectDump.Status(opt)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.CanBeExported {
		t.Error("expected canBeExported to be true")
	}

	if result.CanBeImported {
		t.Error("expected canBeImported to be false")
	}

	if result.ExportedDump != "/path/to/dump.zip" {
		t.Errorf("expected exportedDump '/path/to/dump.zip', got %s", result.ExportedDump)
	}
}

func TestProjectDump_Status_WithID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id != "proj-123" {
			t.Errorf("expected id 'proj-123', got %s", id)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"canBeExported": true, "canBeImported": true}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectDumpStatusOption{
		ID: "proj-123",
	}

	result, resp, err := client.ProjectDump.Status(opt)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestProjectDump_Status_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectDump.Status(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing both ID and Key should fail validation.
	_, _, err = client.ProjectDump.Status(&ProjectDumpStatusOption{})
	if err == nil {
		t.Error("expected error for missing ID and Key")
	}
}

func TestProjectDump_ValidateExportOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.ProjectDump.ValidateExportOpt(&ProjectDumpExportOption{
		Key: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.ProjectDump.ValidateExportOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Key should fail.
	err = client.ProjectDump.ValidateExportOpt(&ProjectDumpExportOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}
}

func TestProjectDump_ValidateStatusOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option with Key should pass.
	err := client.ProjectDump.ValidateStatusOpt(&ProjectDumpStatusOption{
		Key: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Valid option with ID should pass.
	err = client.ProjectDump.ValidateStatusOpt(&ProjectDumpStatusOption{
		ID: "proj-123",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.ProjectDump.ValidateStatusOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing both ID and Key should fail.
	err = client.ProjectDump.ValidateStatusOpt(&ProjectDumpStatusOption{})
	if err == nil {
		t.Error("expected error for missing ID and Key")
	}
}
