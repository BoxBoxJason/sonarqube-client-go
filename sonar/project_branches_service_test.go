package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectBranches_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_branches/delete" {
			t.Errorf("expected path /api/project_branches/delete, got %s", r.URL.Path)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "feature-1" {
			t.Errorf("expected branch 'feature-1', got %s", branch)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBranchesDeleteOption{
		Branch:  "feature-1",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectBranches_Delete_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectBranches.Delete(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.Delete(&ProjectBranchesDeleteOption{
		Project: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing Branch")
	}

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.Delete(&ProjectBranchesDeleteOption{
		Branch: "feature-1",
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBranches_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_branches/list" {
			t.Errorf("expected path /api/project_branches/list, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"branches": [
				{
					"name": "main",
					"isMain": true,
					"type": "LONG",
					"status": {"qualityGateStatus": "OK"},
					"analysisDate": "2024-01-01T00:00:00+0000",
					"excludedFromPurge": true
				},
				{
					"name": "feature-1",
					"isMain": false,
					"type": "BRANCH",
					"status": {"qualityGateStatus": "ERROR"},
					"excludedFromPurge": false
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBranchesListOption{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBranches.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(result.Branches))
	}

	if result.Branches[0].Name != "main" {
		t.Errorf("expected first branch name 'main', got %s", result.Branches[0].Name)
	}

	if !result.Branches[0].IsMain {
		t.Error("expected first branch to be main")
	}

	if result.Branches[0].Status.QualityGateStatus != "OK" {
		t.Errorf("expected first branch status 'OK', got %s", result.Branches[0].Status.QualityGateStatus)
	}

	if result.Branches[1].ExcludedFromPurge {
		t.Error("expected second branch to not be excluded from purge")
	}
}

func TestProjectBranches_List_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectBranches.List(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, _, err = client.ProjectBranches.List(&ProjectBranchesListOption{})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBranches_Rename(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_branches/rename" {
			t.Errorf("expected path /api/project_branches/rename, got %s", r.URL.Path)
		}

		name := r.URL.Query().Get("name")
		if name != "main" {
			t.Errorf("expected name 'main', got %s", name)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBranchesRenameOption{
		Name:    "main",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.Rename(opt)
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectBranches_Rename_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectBranches.Rename(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Name should fail validation.
	_, err = client.ProjectBranches.Rename(&ProjectBranchesRenameOption{
		Project: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.Rename(&ProjectBranchesRenameOption{
		Name: "main",
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxBranchNameLength+1; i++ {
		longName += "a"
	}
	_, err = client.ProjectBranches.Rename(&ProjectBranchesRenameOption{
		Name:    longName,
		Project: "my-project",
	})
	if err == nil {
		t.Error("expected error for name exceeding max length")
	}
}

func TestProjectBranches_SetAutomaticDeletionProtection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_branches/set_automatic_deletion_protection" {
			t.Errorf("expected path /api/project_branches/set_automatic_deletion_protection, got %s", r.URL.Path)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "feature-1" {
			t.Errorf("expected branch 'feature-1', got %s", branch)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		value := r.URL.Query().Get("value")
		if value != "true" {
			t.Errorf("expected value 'true', got %s", value)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   true,
	}

	resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(opt)
	if err != nil {
		t.Fatalf("SetAutomaticDeletionProtection failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectBranches_SetAutomaticDeletionProtection_False(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		if value != "false" {
			t.Errorf("expected value 'false', got %s", value)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   false,
	}

	resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(opt)
	if err != nil {
		t.Fatalf("SetAutomaticDeletionProtection failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectBranches_SetAutomaticDeletionProtection_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectBranches.SetAutomaticDeletionProtection(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.SetAutomaticDeletionProtection(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Project: "my-project",
		Value:   true,
	})
	if err == nil {
		t.Error("expected error for missing Branch")
	}

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.SetAutomaticDeletionProtection(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch: "feature-1",
		Value:  true,
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBranches_SetMain(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_branches/set_main" {
			t.Errorf("expected path /api/project_branches/set_main, got %s", r.URL.Path)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "main" {
			t.Errorf("expected branch 'main', got %s", branch)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBranchesSetMainOption{
		Branch:  "main",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.SetMain(opt)
	if err != nil {
		t.Fatalf("SetMain failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectBranches_SetMain_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectBranches.SetMain(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.SetMain(&ProjectBranchesSetMainOption{
		Project: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing Branch")
	}

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.SetMain(&ProjectBranchesSetMainOption{
		Branch: "main",
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBranches_ValidateDeleteOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.ProjectBranches.ValidateDeleteOpt(&ProjectBranchesDeleteOption{
		Branch:  "feature-1",
		Project: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestProjectBranches_ValidateListOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.ProjectBranches.ValidateListOpt(&ProjectBranchesListOption{
		Project: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestProjectBranches_ValidateRenameOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.ProjectBranches.ValidateRenameOpt(&ProjectBranchesRenameOption{
		Name:    "main",
		Project: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestProjectBranches_ValidateSetAutomaticDeletionProtectionOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option with true should pass.
	err := client.ProjectBranches.ValidateSetAutomaticDeletionProtectionOpt(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   true,
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Valid option with false should pass.
	err = client.ProjectBranches.ValidateSetAutomaticDeletionProtectionOpt(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   false,
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestProjectBranches_ValidateSetMainOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.ProjectBranches.ValidateSetMainOpt(&ProjectBranchesSetMainOption{
		Branch:  "main",
		Project: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
