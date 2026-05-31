package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// AddGroup Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_group", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddGroupOptions{
		GroupName:  "developers",
		Permission: "admin",
	}

	resp, err := client.Permissions.AddGroup(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddGroup_WithProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_group", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddGroupOptions{
		GroupName:  "developers",
		Permission: "user",
		ProjectKey: "my-project",
	}

	resp, err := client.Permissions.AddGroup(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddGroup(context.Background(), nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.AddGroup(context.Background(), &PermissionsAddGroupOptions{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddGroup(context.Background(), &PermissionsAddGroupOptions{
		GroupName: "developers",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddGroup(context.Background(), &PermissionsAddGroupOptions{
		GroupName:  "developers",
		Permission: "invalid",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// AddGroupToTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddGroupToTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_group_to_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddGroupToTemplateOptions{
		GroupName:    "developers",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddGroupToTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddGroupToTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddGroupToTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(context.Background(), &PermissionsAddGroupToTemplateOptions{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(context.Background(), &PermissionsAddGroupToTemplateOptions{
		GroupName:    "developers",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(context.Background(), &PermissionsAddGroupToTemplateOptions{
		GroupName:    "developers",
		Permission:   "gateadmin", // Not a project permission
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(context.Background(), &PermissionsAddGroupToTemplateOptions{
		GroupName:  "developers",
		Permission: "admin",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// AddProjectCreatorToTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddProjectCreatorToTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_project_creator_to_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddProjectCreatorToTemplateOptions{
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddProjectCreatorToTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddProjectCreatorToTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddProjectCreatorToTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(context.Background(), &PermissionsAddProjectCreatorToTemplateOptions{
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(context.Background(), &PermissionsAddProjectCreatorToTemplateOptions{
		Permission:   "provisioning", // Not a project permission
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(context.Background(), &PermissionsAddProjectCreatorToTemplateOptions{
		Permission: "admin",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// AddUser Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_user", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddUserOptions{
		Login:      "john.doe",
		Permission: "admin",
	}

	resp, err := client.Permissions.AddUser(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddUser(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.AddUser(context.Background(), &PermissionsAddUserOptions{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddUser(context.Background(), &PermissionsAddUserOptions{
		Login: "john.doe",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddUser(context.Background(), &PermissionsAddUserOptions{
		Login:      "john.doe",
		Permission: "invalid",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// AddUserToTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_AddUserToTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_user_to_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddUserToTemplateOptions{
		Login:        "john.doe",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddUserToTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddUserToTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddUserToTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.AddUserToTemplate(context.Background(), &PermissionsAddUserToTemplateOptions{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddUserToTemplate(context.Background(), &PermissionsAddUserToTemplateOptions{
		Login:        "john.doe",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddUserToTemplate(context.Background(), &PermissionsAddUserToTemplateOptions{
		Login:        "john.doe",
		Permission:   "profileadmin", // Not a project permission
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddUserToTemplate(context.Background(), &PermissionsAddUserToTemplateOptions{
		Login:      "john.doe",
		Permission: "admin",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ApplyTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_ApplyTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/apply_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsApplyTemplateOptions{
		ProjectKey:   "my-project",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.ApplyTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_ApplyTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.ApplyTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing ProjectID and ProjectKey should fail validation.
	_, err = client.Permissions.ApplyTemplate(context.Background(), &PermissionsApplyTemplateOptions{
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.ApplyTemplate(context.Background(), &PermissionsApplyTemplateOptions{
		ProjectKey: "my-project",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// BulkApplyTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_BulkApplyTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/bulk_apply_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsBulkApplyTemplateOptions{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.BulkApplyTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_BulkApplyTemplate_WithProjects(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/bulk_apply_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsBulkApplyTemplateOptions{
		TemplateName: "my-template",
		Projects:     []string{"project1", "project2"},
	}

	resp, err := client.Permissions.BulkApplyTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_BulkApplyTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.BulkApplyTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.BulkApplyTemplate(context.Background(), &PermissionsBulkApplyTemplateOptions{})
	assert.Error(t, err)

	// Invalid qualifier should fail validation.
	_, err = client.Permissions.BulkApplyTemplate(context.Background(), &PermissionsBulkApplyTemplateOptions{
		TemplateName: "my-template",
		Qualifiers:   "INVALID",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// CreateTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_CreateTemplate(t *testing.T) {
	response := PermissionsCreateTemplate{
		PermissionTemplate: PermissionsTemplateBasic{
			Name:              "my-template",
			Description:       "Template for my projects",
			ProjectKeyPattern: "my-.*",
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/permissions/create_template", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsCreateTemplateOptions{
		Name:              "my-template",
		Description:       "Template for my projects",
		ProjectKeyPattern: "my-.*",
	}

	result, resp, err := client.Permissions.CreateTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "my-template", result.PermissionTemplate.Name)
}

func TestPermissions_CreateTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.CreateTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, _, err = client.Permissions.CreateTemplate(context.Background(), &PermissionsCreateTemplateOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// DeleteTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_DeleteTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/delete_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsDeleteTemplateOptions{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.DeleteTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_DeleteTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.DeleteTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.DeleteTemplate(context.Background(), &PermissionsDeleteTemplateOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Groups Tests
// -----------------------------------------------------------------------------

func TestPermissions_Groups(t *testing.T) {
	response := PermissionsGroups{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     2,
		},
		Groups: []PermissionGroup{
			{
				Name:        "developers",
				Description: "Developers group",
				Permissions: []string{"user", "codeviewer"},
			},
			{
				Name:        "admins",
				Description: "Admins group",
				Permissions: []string{"admin", "user", "codeviewer"},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Permissions.Groups(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 2)
	assert.Equal(t, "developers", result.Groups[0].Name)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestPermissions_Groups_WithOptions(t *testing.T) {
	response := PermissionsGroups{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
		Groups: []PermissionGroup{},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsGroupsOptions{
		ProjectKey: "my-project",
		Permission: "admin",
	}

	_, resp, err := client.Permissions.Groups(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPermissions_Groups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid permission should fail validation.
	_, _, err := client.Permissions.Groups(context.Background(), &PermissionsGroupsOptions{
		Permission: "invalid",
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.Groups(context.Background(), &PermissionsGroupsOptions{
		Query: "ab",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// RemoveGroup Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/remove_group", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsRemoveGroupOptions{
		GroupName:  "developers",
		Permission: "admin",
	}

	resp, err := client.Permissions.RemoveGroup(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveGroup(context.Background(), nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.RemoveGroup(context.Background(), &PermissionsRemoveGroupOptions{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveGroup(context.Background(), &PermissionsRemoveGroupOptions{
		GroupName: "developers",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// RemoveGroupFromTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveGroupFromTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/remove_group_from_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsRemoveGroupFromTemplateOptions{
		GroupName:    "developers",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveGroupFromTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveGroupFromTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveGroupFromTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(context.Background(), &PermissionsRemoveGroupFromTemplateOptions{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(context.Background(), &PermissionsRemoveGroupFromTemplateOptions{
		GroupName:    "developers",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(context.Background(), &PermissionsRemoveGroupFromTemplateOptions{
		GroupName:  "developers",
		Permission: "admin",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// RemoveProjectCreatorFromTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveProjectCreatorFromTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/remove_project_creator_from_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsRemoveProjectCreatorFromTemplateOptions{
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveProjectCreatorFromTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveProjectCreatorFromTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveProjectCreatorFromTemplate(context.Background(), &PermissionsRemoveProjectCreatorFromTemplateOptions{
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveProjectCreatorFromTemplate(context.Background(), &PermissionsRemoveProjectCreatorFromTemplateOptions{
		Permission: "admin",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// RemoveUser Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/remove_user", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsRemoveUserOptions{
		Login:      "john.doe",
		Permission: "admin",
	}

	resp, err := client.Permissions.RemoveUser(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveUser(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.RemoveUser(context.Background(), &PermissionsRemoveUserOptions{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveUser(context.Background(), &PermissionsRemoveUserOptions{
		Login: "john.doe",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// RemoveUserFromTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_RemoveUserFromTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/remove_user_from_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsRemoveUserFromTemplateOptions{
		Login:        "john.doe",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveUserFromTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveUserFromTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveUserFromTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(context.Background(), &PermissionsRemoveUserFromTemplateOptions{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(context.Background(), &PermissionsRemoveUserFromTemplateOptions{
		Login:        "john.doe",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(context.Background(), &PermissionsRemoveUserFromTemplateOptions{
		Login:      "john.doe",
		Permission: "admin",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SearchTemplates Tests
// -----------------------------------------------------------------------------

func TestPermissions_SearchTemplates(t *testing.T) {
	response := PermissionsSearchTemplates{
		PermissionTemplates: []PermissionTemplate{
			{
				ID:                "template-1",
				Name:              "my-template",
				Description:       "My template",
				ProjectKeyPattern: "my-.*",
				CreatedAt:         "2024-01-01T00:00:00+0000",
				UpdatedAt:         "2024-01-02T00:00:00+0000",
				Permissions: []PermissionsTemplatePermission{
					{Key: "admin", UsersCount: 1, GroupsCount: 2, WithProjectCreator: true},
				},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/search_templates", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Permissions.SearchTemplates(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.PermissionTemplates, 1)
	assert.Equal(t, "my-template", result.PermissionTemplates[0].Name)
}

func TestPermissions_ValidateSearchTemplatesOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *PermissionsSearchTemplatesOptions
		wantErr bool
	}{
		{name: "nil option", opt: nil, wantErr: false},
		{name: "empty option", opt: &PermissionsSearchTemplatesOptions{}, wantErr: false},
		{name: "query only", opt: &PermissionsSearchTemplatesOptions{Query: "template"}, wantErr: false},
		{name: "valid pagination", opt: &PermissionsSearchTemplatesOptions{PaginationArgs: PaginationArgs{Page: 1, PageSize: 100}}, wantErr: false},
		{name: "invalid page", opt: &PermissionsSearchTemplatesOptions{PaginationArgs: PaginationArgs{Page: -1, PageSize: 100}}, wantErr: true},
		{name: "invalid page size", opt: &PermissionsSearchTemplatesOptions{PaginationArgs: PaginationArgs{Page: 1, PageSize: 501}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Permissions.ValidateSearchTemplatesOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// SetDefaultTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_SetDefaultTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/set_default_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsSetDefaultTemplateOptions{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.SetDefaultTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_SetDefaultTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.SetDefaultTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.SetDefaultTemplate(context.Background(), &PermissionsSetDefaultTemplateOptions{})
	assert.Error(t, err)

	// Invalid qualifier should fail validation.
	_, err = client.Permissions.SetDefaultTemplate(context.Background(), &PermissionsSetDefaultTemplateOptions{
		TemplateName: "my-template",
		Qualifier:    "INVALID",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// TemplateGroups Tests
// -----------------------------------------------------------------------------

func TestPermissions_TemplateGroups(t *testing.T) {
	response := PermissionsTemplateGroups{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     1,
		},
		Groups: []PermissionsTemplateGroup{
			{
				Name:        "developers",
				Description: "Developers group",
				Permissions: []string{"user", "codeviewer"},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/template_groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsTemplateGroupsOptions{
		TemplateName: "my-template",
	}

	result, resp, err := client.Permissions.TemplateGroups(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 1)
	assert.Equal(t, "developers", result.Groups[0].Name)
}

func TestPermissions_TemplateGroups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.TemplateGroups(context.Background(), nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, _, err = client.Permissions.TemplateGroups(context.Background(), &PermissionsTemplateGroupsOptions{})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, _, err = client.Permissions.TemplateGroups(context.Background(), &PermissionsTemplateGroupsOptions{
		TemplateName: "my-template",
		Permission:   "gateadmin", // Not a project permission
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.TemplateGroups(context.Background(), &PermissionsTemplateGroupsOptions{
		TemplateName: "my-template",
		Query:        "ab",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// TemplateUsers Tests
// -----------------------------------------------------------------------------

func TestPermissions_TemplateUsers(t *testing.T) {
	response := PermissionsTemplateUsers{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     1,
		},
		Users: []PermissionsTemplateUser{
			{
				Login:       "john.doe",
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Permissions: []string{"admin", "user"},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/template_users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsTemplateUsersOptions{
		TemplateName: "my-template",
	}

	result, resp, err := client.Permissions.TemplateUsers(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
}

func TestPermissions_TemplateUsers_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.TemplateUsers(context.Background(), nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, _, err = client.Permissions.TemplateUsers(context.Background(), &PermissionsTemplateUsersOptions{})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, _, err = client.Permissions.TemplateUsers(context.Background(), &PermissionsTemplateUsersOptions{
		TemplateName: "my-template",
		Permission:   "provisioning", // Not a project permission
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.TemplateUsers(context.Background(), &PermissionsTemplateUsersOptions{
		TemplateName: "my-template",
		Query:        "ab",
	})
	assert.Error(t, err)
}

func TestPermissions_UpdateTemplate(t *testing.T) {
	response := PermissionsUpdateTemplate{
		PermissionTemplate: PermissionsTemplateUpdated{
			ID:                "template-1",
			Name:              "new-template-name",
			Description:       "Updated description",
			ProjectKeyPattern: "new-.*",
			CreatedAt:         "2024-01-01T00:00:00+0000",
			UpdatedAt:         "2024-01-03T00:00:00+0000",
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/permissions/update_template", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsUpdateTemplateOptions{
		ID:                "template-1",
		Name:              "new-template-name",
		Description:       "Updated description",
		ProjectKeyPattern: "new-.*",
	}

	result, resp, err := client.Permissions.UpdateTemplate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "new-template-name", result.PermissionTemplate.Name)
}

func TestPermissions_UpdateTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.UpdateTemplate(context.Background(), nil)
	assert.Error(t, err)

	// Missing ID should fail validation.
	_, _, err = client.Permissions.UpdateTemplate(context.Background(), &PermissionsUpdateTemplateOptions{
		Name: "new-name",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Users Tests
// -----------------------------------------------------------------------------

func TestPermissions_Users(t *testing.T) {
	response := PermissionsUsers{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     2,
		},
		Users: []PermissionUser{
			{
				Login:       "john.doe",
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Permissions: []string{"admin", "user"},
				Managed:     false,
			},
			{
				Login:       "jane.doe",
				Name:        "Jane Doe",
				Email:       "jane.doe@example.com",
				Permissions: []string{"user", "codeviewer"},
				Managed:     true,
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Permissions.Users(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, "john.doe", result.Users[0].Login)
	assert.True(t, result.Users[1].Managed)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestPermissions_Users_WithOptions(t *testing.T) {
	response := PermissionsUsers{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
		Users: []PermissionUser{},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsUsersOptions{
		ProjectKey: "my-project",
		Permission: "admin",
	}

	_, resp, err := client.Permissions.Users(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPermissions_Users_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid permission should fail validation.
	_, _, err := client.Permissions.Users(context.Background(), &PermissionsUsersOptions{
		Permission: "invalid",
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.Users(context.Background(), &PermissionsUsersOptions{
		Query: "ab",
	})
	assert.Error(t, err)
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
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionsService_GroupsAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"groups":[{"name":"g1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"groups":[{"name":"g2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &PermissionsGroupsOptions{}
		result, _, err := client.Permissions.GroupsAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

}

func TestPermissionsService_TemplateGroupsAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"groups":[{"name":"g1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"groups":[{"name":"g2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &PermissionsTemplateGroupsOptions{TemplateID: "tmpl1"}
		result, _, err := client.Permissions.TemplateGroupsAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Permissions.TemplateGroupsAll(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestPermissionsService_TemplateUsersAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"users":[{"login":"u1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"users":[{"login":"u2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &PermissionsTemplateUsersOptions{TemplateID: "tmpl1"}
		result, _, err := client.Permissions.TemplateUsersAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Permissions.TemplateUsersAll(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestPermissionsService_UsersAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"users":[{"login":"u1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"users":[{"login":"u2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &PermissionsUsersOptions{}
		result, _, err := client.Permissions.UsersAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})
}
