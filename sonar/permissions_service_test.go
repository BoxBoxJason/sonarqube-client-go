package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// -----------------------------------------------------------------------------
// AddGroup Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/add_group" {
			t.Errorf("expected path /api/permissions/add_group, got %s", r.URL.Path)
		}

		groupName := r.URL.Query().Get("groupName")
		if groupName != "developers" {
			t.Errorf("expected groupName 'developers', got %s", groupName)
		}

		permission := r.URL.Query().Get("permission")
		if permission != "admin" {
			t.Errorf("expected permission 'admin', got %s", permission)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsAddGroupOption{
		GroupName:  "developers",
		Permission: "admin",
	}

	resp, err := client.Permissions.AddGroup(opt)
	if err != nil {
		t.Fatalf("AddGroup failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_AddGroup_WithProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		permission := r.URL.Query().Get("permission")
		if permission != "user" {
			t.Errorf("expected permission 'user', got %s", permission)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsAddGroupOption{
		GroupName:  "developers",
		Permission: "user",
		ProjectKey: "my-project",
	}

	resp, err := client.Permissions.AddGroup(opt)
	if err != nil {
		t.Fatalf("AddGroup failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_AddGroup_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.AddGroup(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing GroupName should fail validation.
	_, err = client.Permissions.AddGroup(&PermissionsAddGroupOption{
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing GroupName")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddGroup(&PermissionsAddGroupOption{
		GroupName: "developers",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddGroup(&PermissionsAddGroupOption{
		GroupName:  "developers",
		Permission: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Permission")
	}
}

// -----------------------------------------------------------------------------
// AddGroupToTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddGroupToTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/add_group_to_template" {
			t.Errorf("expected path /api/permissions/add_group_to_template, got %s", r.URL.Path)
		}

		groupName := r.URL.Query().Get("groupName")
		if groupName != "developers" {
			t.Errorf("expected groupName 'developers', got %s", groupName)
		}

		permission := r.URL.Query().Get("permission")
		if permission != "admin" {
			t.Errorf("expected permission 'admin', got %s", permission)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsAddGroupToTemplateOption{
		GroupName:    "developers",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddGroupToTemplate(opt)
	if err != nil {
		t.Fatalf("AddGroupToTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_AddGroupToTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.AddGroupToTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing GroupName should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing GroupName")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		GroupName:    "developers",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		GroupName:    "developers",
		Permission:   "gateadmin", // Not a project permission
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for invalid Permission (non-project permission)")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		GroupName:  "developers",
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// AddProjectCreatorToTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddProjectCreatorToTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/add_project_creator_to_template" {
			t.Errorf("expected path /api/permissions/add_project_creator_to_template, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsAddProjectCreatorToTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddProjectCreatorToTemplate(opt)
	if err != nil {
		t.Fatalf("AddProjectCreatorToTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_AddProjectCreatorToTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.AddProjectCreatorToTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(&PermissionsAddProjectCreatorToTemplateOption{
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(&PermissionsAddProjectCreatorToTemplateOption{
		Permission:   "provisioning", // Not a project permission
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for invalid Permission")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(&PermissionsAddProjectCreatorToTemplateOption{
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// AddUser Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/add_user" {
			t.Errorf("expected path /api/permissions/add_user, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "john.doe" {
			t.Errorf("expected login 'john.doe', got %s", login)
		}

		permission := r.URL.Query().Get("permission")
		if permission != "admin" {
			t.Errorf("expected permission 'admin', got %s", permission)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsAddUserOption{
		Login:      "john.doe",
		Permission: "admin",
	}

	resp, err := client.Permissions.AddUser(opt)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_AddUser_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.AddUser(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, err = client.Permissions.AddUser(&PermissionsAddUserOption{
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddUser(&PermissionsAddUserOption{
		Login: "john.doe",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddUser(&PermissionsAddUserOption{
		Login:      "john.doe",
		Permission: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Permission")
	}
}

// -----------------------------------------------------------------------------
// AddUserToTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddUserToTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/add_user_to_template" {
			t.Errorf("expected path /api/permissions/add_user_to_template, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsAddUserToTemplateOption{
		Login:        "john.doe",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddUserToTemplate(opt)
	if err != nil {
		t.Fatalf("AddUserToTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_AddUserToTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.AddUserToTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Login:        "john.doe",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Login:        "john.doe",
		Permission:   "profileadmin", // Not a project permission
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for invalid Permission")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Login:      "john.doe",
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// ApplyTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_ApplyTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/apply_template" {
			t.Errorf("expected path /api/permissions/apply_template, got %s", r.URL.Path)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsApplyTemplateOption{
		ProjectKey:   "my-project",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.ApplyTemplate(opt)
	if err != nil {
		t.Fatalf("ApplyTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_ApplyTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.ApplyTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing ProjectID and ProjectKey should fail validation.
	_, err = client.Permissions.ApplyTemplate(&PermissionsApplyTemplateOption{
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing ProjectID and ProjectKey")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.ApplyTemplate(&PermissionsApplyTemplateOption{
		ProjectKey: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// BulkApplyTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_BulkApplyTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/bulk_apply_template" {
			t.Errorf("expected path /api/permissions/bulk_apply_template, got %s", r.URL.Path)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsBulkApplyTemplateOption{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.BulkApplyTemplate(opt)
	if err != nil {
		t.Fatalf("BulkApplyTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_BulkApplyTemplate_WithProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projects := r.URL.Query()["projects"]
		if len(projects) != 2 {
			t.Errorf("expected 2 projects, got %d", len(projects))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsBulkApplyTemplateOption{
		TemplateName: "my-template",
		Projects:     []string{"project1", "project2"},
	}

	resp, err := client.Permissions.BulkApplyTemplate(opt)
	if err != nil {
		t.Fatalf("BulkApplyTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_BulkApplyTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.BulkApplyTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.BulkApplyTemplate(&PermissionsBulkApplyTemplateOption{})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}

	// Invalid qualifier should fail validation.
	_, err = client.Permissions.BulkApplyTemplate(&PermissionsBulkApplyTemplateOption{
		TemplateName: "my-template",
		Qualifiers:   "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid Qualifiers")
	}
}

// -----------------------------------------------------------------------------
// CreateTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_CreateTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/create_template" {
			t.Errorf("expected path /api/permissions/create_template, got %s", r.URL.Path)
		}

		name := r.URL.Query().Get("name")
		if name != "my-template" {
			t.Errorf("expected name 'my-template', got %s", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"permissionTemplate": {
				"name": "my-template",
				"description": "Template for my projects",
				"projectKeyPattern": "my-.*"
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsCreateTemplateOption{
		Name:              "my-template",
		Description:       "Template for my projects",
		ProjectKeyPattern: "my-.*",
	}

	result, resp, err := client.Permissions.CreateTemplate(opt)
	if err != nil {
		t.Fatalf("CreateTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.PermissionTemplate.Name != "my-template" {
		t.Errorf("expected template name 'my-template', got %s", result.PermissionTemplate.Name)
	}
}

func TestPermissions_CreateTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Permissions.CreateTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Name should fail validation.
	_, _, err = client.Permissions.CreateTemplate(&PermissionsCreateTemplateOption{})
	if err == nil {
		t.Error("expected error for missing Name")
	}
}

// -----------------------------------------------------------------------------
// DeleteTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_DeleteTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/delete_template" {
			t.Errorf("expected path /api/permissions/delete_template, got %s", r.URL.Path)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsDeleteTemplateOption{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.DeleteTemplate(opt)
	if err != nil {
		t.Fatalf("DeleteTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_DeleteTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.DeleteTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.DeleteTemplate(&PermissionsDeleteTemplateOption{})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// Groups Tests
// -----------------------------------------------------------------------------

func TestPermissions_Groups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/groups" {
			t.Errorf("expected path /api/permissions/groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 2
			},
			"groups": [
				{
					"name": "developers",
					"description": "Developers group",
					"permissions": ["user", "codeviewer"]
				},
				{
					"name": "admins",
					"description": "Admins group",
					"permissions": ["admin", "user", "codeviewer"]
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Permissions.Groups(nil)
	if err != nil {
		t.Fatalf("Groups failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result.Groups))
	}

	if result.Groups[0].Name != "developers" {
		t.Errorf("expected group name 'developers', got %s", result.Groups[0].Name)
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Paging.Total)
	}
}

func TestPermissions_Groups_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		permission := r.URL.Query().Get("permission")
		if permission != "admin" {
			t.Errorf("expected permission 'admin', got %s", permission)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 25, "total": 0}, "groups": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsGroupsOption{
		ProjectKey: "my-project",
		Permission: "admin",
	}

	_, resp, err := client.Permissions.Groups(opt)
	if err != nil {
		t.Fatalf("Groups failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestPermissions_Groups_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Invalid permission should fail validation.
	_, _, err := client.Permissions.Groups(&PermissionsGroupsOption{
		Permission: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Permission")
	}

	// Query too short should fail validation.
	_, _, err = client.Permissions.Groups(&PermissionsGroupsOption{
		Query: "ab",
	})
	if err == nil {
		t.Error("expected error for Query too short")
	}
}

// -----------------------------------------------------------------------------
// RemoveGroup Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/remove_group" {
			t.Errorf("expected path /api/permissions/remove_group, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsRemoveGroupOption{
		GroupName:  "developers",
		Permission: "admin",
	}

	resp, err := client.Permissions.RemoveGroup(opt)
	if err != nil {
		t.Fatalf("RemoveGroup failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_RemoveGroup_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveGroup(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing GroupName should fail validation.
	_, err = client.Permissions.RemoveGroup(&PermissionsRemoveGroupOption{
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing GroupName")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveGroup(&PermissionsRemoveGroupOption{
		GroupName: "developers",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}
}

// -----------------------------------------------------------------------------
// RemoveGroupFromTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveGroupFromTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/remove_group_from_template" {
			t.Errorf("expected path /api/permissions/remove_group_from_template, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsRemoveGroupFromTemplateOption{
		GroupName:    "developers",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveGroupFromTemplate(opt)
	if err != nil {
		t.Fatalf("RemoveGroupFromTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_RemoveGroupFromTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveGroupFromTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing GroupName should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(&PermissionsRemoveGroupFromTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing GroupName")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(&PermissionsRemoveGroupFromTemplateOption{
		GroupName:    "developers",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(&PermissionsRemoveGroupFromTemplateOption{
		GroupName:  "developers",
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// RemoveProjectCreatorFromTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveProjectCreatorFromTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/remove_project_creator_from_template" {
			t.Errorf("expected path /api/permissions/remove_project_creator_from_template, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsRemoveProjectCreatorFromTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(opt)
	if err != nil {
		t.Fatalf("RemoveProjectCreatorFromTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_RemoveProjectCreatorFromTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveProjectCreatorFromTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveProjectCreatorFromTemplate(&PermissionsRemoveProjectCreatorFromTemplateOption{
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveProjectCreatorFromTemplate(&PermissionsRemoveProjectCreatorFromTemplateOption{
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// RemoveUser Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/remove_user" {
			t.Errorf("expected path /api/permissions/remove_user, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsRemoveUserOption{
		Login:      "john.doe",
		Permission: "admin",
	}

	resp, err := client.Permissions.RemoveUser(opt)
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_RemoveUser_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveUser(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, err = client.Permissions.RemoveUser(&PermissionsRemoveUserOption{
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveUser(&PermissionsRemoveUserOption{
		Login: "john.doe",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}
}

// -----------------------------------------------------------------------------
// RemoveUserFromTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveUserFromTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/remove_user_from_template" {
			t.Errorf("expected path /api/permissions/remove_user_from_template, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsRemoveUserFromTemplateOption{
		Login:        "john.doe",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveUserFromTemplate(opt)
	if err != nil {
		t.Fatalf("RemoveUserFromTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_RemoveUserFromTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveUserFromTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(&PermissionsRemoveUserFromTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(&PermissionsRemoveUserFromTemplateOption{
		Login:        "john.doe",
		TemplateName: "my-template",
	})
	if err == nil {
		t.Error("expected error for missing Permission")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(&PermissionsRemoveUserFromTemplateOption{
		Login:      "john.doe",
		Permission: "admin",
	})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}
}

// -----------------------------------------------------------------------------
// SearchTemplates Tests
// -----------------------------------------------------------------------------

func TestPermissions_SearchTemplates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/search_templates" {
			t.Errorf("expected path /api/permissions/search_templates, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"permissionTemplates": [
				{
					"id": "template-1",
					"name": "my-template",
					"description": "My template",
					"projectKeyPattern": "my-.*",
					"createdAt": "2024-01-01T00:00:00+0000",
					"updatedAt": "2024-01-02T00:00:00+0000",
					"permissions": [
						{"key": "admin", "usersCount": 1, "groupsCount": 2, "withProjectCreator": true}
					]
				}
			],
			"defaultTemplates": [
				{"qualifier": "TRK", "templateId": "template-1"}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Permissions.SearchTemplates(nil)
	if err != nil {
		t.Fatalf("SearchTemplates failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.PermissionTemplates) != 1 {
		t.Errorf("expected 1 template, got %d", len(result.PermissionTemplates))
	}

	if result.PermissionTemplates[0].Name != "my-template" {
		t.Errorf("expected template name 'my-template', got %s", result.PermissionTemplates[0].Name)
	}

	if len(result.DefaultTemplates) != 1 {
		t.Errorf("expected 1 default template, got %d", len(result.DefaultTemplates))
	}
}

// -----------------------------------------------------------------------------
// SetDefaultTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_SetDefaultTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/set_default_template" {
			t.Errorf("expected path /api/permissions/set_default_template, got %s", r.URL.Path)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsSetDefaultTemplateOption{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.SetDefaultTemplate(opt)
	if err != nil {
		t.Fatalf("SetDefaultTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPermissions_SetDefaultTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Permissions.SetDefaultTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.SetDefaultTemplate(&PermissionsSetDefaultTemplateOption{})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}

	// Invalid qualifier should fail validation.
	_, err = client.Permissions.SetDefaultTemplate(&PermissionsSetDefaultTemplateOption{
		TemplateName: "my-template",
		Qualifier:    "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid Qualifier")
	}
}

// -----------------------------------------------------------------------------
// TemplateGroups Tests
// -----------------------------------------------------------------------------

func TestPermissions_TemplateGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/template_groups" {
			t.Errorf("expected path /api/permissions/template_groups, got %s", r.URL.Path)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 1
			},
			"groups": [
				{
					"name": "developers",
					"description": "Developers group",
					"permissions": ["user", "codeviewer"]
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsTemplateGroupsOption{
		TemplateName: "my-template",
	}

	result, resp, err := client.Permissions.TemplateGroups(opt)
	if err != nil {
		t.Fatalf("TemplateGroups failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(result.Groups))
	}

	if result.Groups[0].Name != "developers" {
		t.Errorf("expected group name 'developers', got %s", result.Groups[0].Name)
	}
}

func TestPermissions_TemplateGroups_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Permissions.TemplateGroups(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, _, err = client.Permissions.TemplateGroups(&PermissionsTemplateGroupsOption{})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}

	// Invalid permission should fail validation.
	_, _, err = client.Permissions.TemplateGroups(&PermissionsTemplateGroupsOption{
		TemplateName: "my-template",
		Permission:   "gateadmin", // Not a project permission
	})
	if err == nil {
		t.Error("expected error for invalid Permission (non-project permission)")
	}

	// Query too short should fail validation.
	_, _, err = client.Permissions.TemplateGroups(&PermissionsTemplateGroupsOption{
		TemplateName: "my-template",
		Query:        "ab",
	})
	if err == nil {
		t.Error("expected error for Query too short")
	}
}

// -----------------------------------------------------------------------------
// TemplateUsers Tests
// -----------------------------------------------------------------------------

func TestPermissions_TemplateUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/template_users" {
			t.Errorf("expected path /api/permissions/template_users, got %s", r.URL.Path)
		}

		templateName := r.URL.Query().Get("templateName")
		if templateName != "my-template" {
			t.Errorf("expected templateName 'my-template', got %s", templateName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 1
			},
			"users": [
				{
					"login": "john.doe",
					"name": "John Doe",
					"email": "john.doe@example.com",
					"permissions": ["admin", "user"]
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsTemplateUsersOption{
		TemplateName: "my-template",
	}

	result, resp, err := client.Permissions.TemplateUsers(opt)
	if err != nil {
		t.Fatalf("TemplateUsers failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Users) != 1 {
		t.Errorf("expected 1 user, got %d", len(result.Users))
	}

	if result.Users[0].Login != "john.doe" {
		t.Errorf("expected user login 'john.doe', got %s", result.Users[0].Login)
	}
}

func TestPermissions_TemplateUsers_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Permissions.TemplateUsers(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing TemplateID and TemplateName should fail validation.
	_, _, err = client.Permissions.TemplateUsers(&PermissionsTemplateUsersOption{})
	if err == nil {
		t.Error("expected error for missing TemplateID and TemplateName")
	}

	// Invalid permission should fail validation.
	_, _, err = client.Permissions.TemplateUsers(&PermissionsTemplateUsersOption{
		TemplateName: "my-template",
		Permission:   "provisioning", // Not a project permission
	})
	if err == nil {
		t.Error("expected error for invalid Permission (non-project permission)")
	}

	// Query too short should fail validation.
	_, _, err = client.Permissions.TemplateUsers(&PermissionsTemplateUsersOption{
		TemplateName: "my-template",
		Query:        "ab",
	})
	if err == nil {
		t.Error("expected error for Query too short")
	}
}

func TestPermissions_UpdateTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/update_template" {
			t.Errorf("expected path /api/permissions/update_template, got %s", r.URL.Path)
		}

		id := r.URL.Query().Get("id")
		if id != "template-1" {
			t.Errorf("expected id 'template-1', got %s", id)
		}

		name := r.URL.Query().Get("name")
		if name != "new-template-name" {
			t.Errorf("expected name 'new-template-name', got %s", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"permissionTemplate": {
				"id": "template-1",
				"name": "new-template-name",
				"description": "Updated description",
				"projectKeyPattern": "new-.*",
				"createdAt": "2024-01-01T00:00:00+0000",
				"updatedAt": "2024-01-03T00:00:00+0000"
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsUpdateTemplateOption{
		ID:                "template-1",
		Name:              "new-template-name",
		Description:       "Updated description",
		ProjectKeyPattern: "new-.*",
	}

	result, resp, err := client.Permissions.UpdateTemplate(opt)
	if err != nil {
		t.Fatalf("UpdateTemplate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.PermissionTemplate.Name != "new-template-name" {
		t.Errorf("expected template name 'new-template-name', got %s", result.PermissionTemplate.Name)
	}
}

func TestPermissions_UpdateTemplate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Permissions.UpdateTemplate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing ID should fail validation.
	_, _, err = client.Permissions.UpdateTemplate(&PermissionsUpdateTemplateOption{
		Name: "new-name",
	})
	if err == nil {
		t.Error("expected error for missing ID")
	}
}

// -----------------------------------------------------------------------------
// Users Tests
// -----------------------------------------------------------------------------

func TestPermissions_Users(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/permissions/users" {
			t.Errorf("expected path /api/permissions/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 2
			},
			"users": [
				{
					"login": "john.doe",
					"name": "John Doe",
					"email": "john.doe@example.com",
					"permissions": ["admin", "user"],
					"managed": false
				},
				{
					"login": "jane.doe",
					"name": "Jane Doe",
					"email": "jane.doe@example.com",
					"permissions": ["user", "codeviewer"],
					"managed": true
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Permissions.Users(nil)
	if err != nil {
		t.Fatalf("Users failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(result.Users))
	}

	if result.Users[0].Login != "john.doe" {
		t.Errorf("expected user login 'john.doe', got %s", result.Users[0].Login)
	}

	if result.Users[1].Managed != true {
		t.Error("expected jane.doe to be managed")
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Paging.Total)
	}
}

func TestPermissions_Users_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		permission := r.URL.Query().Get("permission")
		if permission != "admin" {
			t.Errorf("expected permission 'admin', got %s", permission)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 25, "total": 0}, "users": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PermissionsUsersOption{
		ProjectKey: "my-project",
		Permission: "admin",
	}

	_, resp, err := client.Permissions.Users(opt)
	if err != nil {
		t.Fatalf("Users failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestPermissions_Users_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Invalid permission should fail validation.
	_, _, err := client.Permissions.Users(&PermissionsUsersOption{
		Permission: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Permission")
	}

	// Query too short should fail validation.
	_, _, err = client.Permissions.Users(&PermissionsUsersOption{
		Query: "ab",
	})
	if err == nil {
		t.Error("expected error for Query too short")
	}
}

// -----------------------------------------------------------------------------
// Helper Function Tests
// -----------------------------------------------------------------------------

func TestPermissions_isValidPermission(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		expected   bool
	}{
		// Global permissions
		{"admin global", "admin", true},
		{"gateadmin global", "gateadmin", true},
		{"profileadmin global", "profileadmin", true},
		{"provisioning global", "provisioning", true},
		{"scan global", "scan", true},
		{"applicationcreator global", "applicationcreator", true},
		{"portfoliocreator global", "portfoliocreator", true},
		// Project permissions
		{"admin project", "admin", true}, // Also global
		{"codeviewer project", "codeviewer", true},
		{"issueadmin project", "issueadmin", true},
		{"securityhotspotadmin project", "securityhotspotadmin", true},
		{"scan project", "scan", true}, // Also global
		{"user project", "user", true},
		// Invalid permissions
		{"empty", "", false},
		{"invalid", "invalid", false},
		{"Admin (case)", "Admin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPermission(tt.permission)
			if result != tt.expected {
				t.Errorf("isValidPermission(%q) = %v, want %v", tt.permission, result, tt.expected)
			}
		})
	}
}
