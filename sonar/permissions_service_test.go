package sonar

import (
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

	opt := &PermissionsAddGroupOption{
		GroupName:  "developers",
		Permission: "admin",
	}

	resp, err := client.Permissions.AddGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddGroup_WithProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/add_group", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsAddGroupOption{
		GroupName:  "developers",
		Permission: "user",
		ProjectKey: "my-project",
	}

	resp, err := client.Permissions.AddGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddGroup(nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.AddGroup(&PermissionsAddGroupOption{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddGroup(&PermissionsAddGroupOption{
		GroupName: "developers",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddGroup(&PermissionsAddGroupOption{
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

	opt := &PermissionsAddGroupToTemplateOption{
		GroupName:    "developers",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddGroupToTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddGroupToTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddGroupToTemplate(nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		GroupName:    "developers",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
		GroupName:    "developers",
		Permission:   "gateadmin", // Not a project permission
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddGroupToTemplate(&PermissionsAddGroupToTemplateOption{
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

	opt := &PermissionsAddProjectCreatorToTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddProjectCreatorToTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddProjectCreatorToTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddProjectCreatorToTemplate(nil)
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(&PermissionsAddProjectCreatorToTemplateOption{
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(&PermissionsAddProjectCreatorToTemplateOption{
		Permission:   "provisioning", // Not a project permission
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddProjectCreatorToTemplate(&PermissionsAddProjectCreatorToTemplateOption{
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

	opt := &PermissionsAddUserOption{
		Login:      "john.doe",
		Permission: "admin",
	}

	resp, err := client.Permissions.AddUser(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddUser(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.AddUser(&PermissionsAddUserOption{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddUser(&PermissionsAddUserOption{
		Login: "john.doe",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddUser(&PermissionsAddUserOption{
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

	opt := &PermissionsAddUserToTemplateOption{
		Login:        "john.doe",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.AddUserToTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_AddUserToTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.AddUserToTemplate(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Login:        "john.doe",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
		Login:        "john.doe",
		Permission:   "profileadmin", // Not a project permission
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.AddUserToTemplate(&PermissionsAddUserToTemplateOption{
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

	opt := &PermissionsApplyTemplateOption{
		ProjectKey:   "my-project",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.ApplyTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_ApplyTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.ApplyTemplate(nil)
	assert.Error(t, err)

	// Missing ProjectID and ProjectKey should fail validation.
	_, err = client.Permissions.ApplyTemplate(&PermissionsApplyTemplateOption{
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.ApplyTemplate(&PermissionsApplyTemplateOption{
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

	opt := &PermissionsBulkApplyTemplateOption{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.BulkApplyTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_BulkApplyTemplate_WithProjects(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/bulk_apply_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsBulkApplyTemplateOption{
		TemplateName: "my-template",
		Projects:     []string{"project1", "project2"},
	}

	resp, err := client.Permissions.BulkApplyTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_BulkApplyTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.BulkApplyTemplate(nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.BulkApplyTemplate(&PermissionsBulkApplyTemplateOption{})
	assert.Error(t, err)

	// Invalid qualifier should fail validation.
	_, err = client.Permissions.BulkApplyTemplate(&PermissionsBulkApplyTemplateOption{
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
		PermissionTemplate: PermissionTemplateBasic{
			Name:              "my-template",
			Description:       "Template for my projects",
			ProjectKeyPattern: "my-.*",
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/permissions/create_template", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsCreateTemplateOption{
		Name:              "my-template",
		Description:       "Template for my projects",
		ProjectKeyPattern: "my-.*",
	}

	result, resp, err := client.Permissions.CreateTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "my-template", result.PermissionTemplate.Name)
}

func TestPermissions_CreateTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.CreateTemplate(nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, _, err = client.Permissions.CreateTemplate(&PermissionsCreateTemplateOption{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// DeleteTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_DeleteTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/delete_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsDeleteTemplateOption{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.DeleteTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_DeleteTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.DeleteTemplate(nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.DeleteTemplate(&PermissionsDeleteTemplateOption{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Groups Tests
// -----------------------------------------------------------------------------

func TestPermissions_Groups(t *testing.T) {
	response := PermissionsGroups{
		Paging: PermissionsPaging{
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

	result, resp, err := client.Permissions.Groups(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 2)
	assert.Equal(t, "developers", result.Groups[0].Name)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestPermissions_Groups_WithOptions(t *testing.T) {
	response := PermissionsGroups{
		Paging: PermissionsPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
		Groups: []PermissionGroup{},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsGroupsOption{
		ProjectKey: "my-project",
		Permission: "admin",
	}

	_, resp, err := client.Permissions.Groups(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPermissions_Groups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid permission should fail validation.
	_, _, err := client.Permissions.Groups(&PermissionsGroupsOption{
		Permission: "invalid",
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.Groups(&PermissionsGroupsOption{
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

	opt := &PermissionsRemoveGroupOption{
		GroupName:  "developers",
		Permission: "admin",
	}

	resp, err := client.Permissions.RemoveGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveGroup(nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.RemoveGroup(&PermissionsRemoveGroupOption{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveGroup(&PermissionsRemoveGroupOption{
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

	opt := &PermissionsRemoveGroupFromTemplateOption{
		GroupName:    "developers",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveGroupFromTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveGroupFromTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveGroupFromTemplate(nil)
	assert.Error(t, err)

	// Missing GroupName should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(&PermissionsRemoveGroupFromTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(&PermissionsRemoveGroupFromTemplateOption{
		GroupName:    "developers",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveGroupFromTemplate(&PermissionsRemoveGroupFromTemplateOption{
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

	opt := &PermissionsRemoveProjectCreatorFromTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveProjectCreatorFromTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveProjectCreatorFromTemplate(nil)
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveProjectCreatorFromTemplate(&PermissionsRemoveProjectCreatorFromTemplateOption{
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveProjectCreatorFromTemplate(&PermissionsRemoveProjectCreatorFromTemplateOption{
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

	opt := &PermissionsRemoveUserOption{
		Login:      "john.doe",
		Permission: "admin",
	}

	resp, err := client.Permissions.RemoveUser(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveUser(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.RemoveUser(&PermissionsRemoveUserOption{
		Permission: "admin",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveUser(&PermissionsRemoveUserOption{
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

	opt := &PermissionsRemoveUserFromTemplateOption{
		Login:        "john.doe",
		Permission:   "admin",
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.RemoveUserFromTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_RemoveUserFromTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.RemoveUserFromTemplate(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(&PermissionsRemoveUserFromTemplateOption{
		Permission:   "admin",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing Permission should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(&PermissionsRemoveUserFromTemplateOption{
		Login:        "john.doe",
		TemplateName: "my-template",
	})
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.RemoveUserFromTemplate(&PermissionsRemoveUserFromTemplateOption{
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
				Permissions: []TemplatePermission{
					{Key: "admin", UsersCount: 1, GroupsCount: 2, WithProjectCreator: true},
				},
			},
		},
		DefaultTemplates: []DefaultTemplate{
			{Qualifier: "TRK", TemplateID: "template-1"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/search_templates", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Permissions.SearchTemplates(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.PermissionTemplates, 1)
	assert.Equal(t, "my-template", result.PermissionTemplates[0].Name)
	assert.Len(t, result.DefaultTemplates, 1)
}

// -----------------------------------------------------------------------------
// SetDefaultTemplate Tests
// -----------------------------------------------------------------------------

func TestPermissions_SetDefaultTemplate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/permissions/set_default_template", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &PermissionsSetDefaultTemplateOption{
		TemplateName: "my-template",
	}

	resp, err := client.Permissions.SetDefaultTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPermissions_SetDefaultTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Permissions.SetDefaultTemplate(nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, err = client.Permissions.SetDefaultTemplate(&PermissionsSetDefaultTemplateOption{})
	assert.Error(t, err)

	// Invalid qualifier should fail validation.
	_, err = client.Permissions.SetDefaultTemplate(&PermissionsSetDefaultTemplateOption{
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
		Paging: PermissionsPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     1,
		},
		Groups: []TemplateGroup{
			{
				Name:        "developers",
				Description: "Developers group",
				Permissions: []string{"user", "codeviewer"},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/template_groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsTemplateGroupsOption{
		TemplateName: "my-template",
	}

	result, resp, err := client.Permissions.TemplateGroups(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 1)
	assert.Equal(t, "developers", result.Groups[0].Name)
}

func TestPermissions_TemplateGroups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.TemplateGroups(nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, _, err = client.Permissions.TemplateGroups(&PermissionsTemplateGroupsOption{})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, _, err = client.Permissions.TemplateGroups(&PermissionsTemplateGroupsOption{
		TemplateName: "my-template",
		Permission:   "gateadmin", // Not a project permission
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.TemplateGroups(&PermissionsTemplateGroupsOption{
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
		Paging: PermissionsPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     1,
		},
		Users: []TemplateUser{
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

	opt := &PermissionsTemplateUsersOption{
		TemplateName: "my-template",
	}

	result, resp, err := client.Permissions.TemplateUsers(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
}

func TestPermissions_TemplateUsers_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.TemplateUsers(nil)
	assert.Error(t, err)

	// Missing TemplateID and TemplateName should fail validation.
	_, _, err = client.Permissions.TemplateUsers(&PermissionsTemplateUsersOption{})
	assert.Error(t, err)

	// Invalid permission should fail validation.
	_, _, err = client.Permissions.TemplateUsers(&PermissionsTemplateUsersOption{
		TemplateName: "my-template",
		Permission:   "provisioning", // Not a project permission
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.TemplateUsers(&PermissionsTemplateUsersOption{
		TemplateName: "my-template",
		Query:        "ab",
	})
	assert.Error(t, err)
}

func TestPermissions_UpdateTemplate(t *testing.T) {
	response := PermissionsUpdateTemplate{
		PermissionTemplate: PermissionTemplateUpdated{
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

	opt := &PermissionsUpdateTemplateOption{
		ID:                "template-1",
		Name:              "new-template-name",
		Description:       "Updated description",
		ProjectKeyPattern: "new-.*",
	}

	result, resp, err := client.Permissions.UpdateTemplate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "new-template-name", result.PermissionTemplate.Name)
}

func TestPermissions_UpdateTemplate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Permissions.UpdateTemplate(nil)
	assert.Error(t, err)

	// Missing ID should fail validation.
	_, _, err = client.Permissions.UpdateTemplate(&PermissionsUpdateTemplateOption{
		Name: "new-name",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Users Tests
// -----------------------------------------------------------------------------

func TestPermissions_Users(t *testing.T) {
	response := PermissionsUsers{
		Paging: PermissionsPaging{
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

	result, resp, err := client.Permissions.Users(nil)
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
		Paging: PermissionsPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
		Users: []PermissionUser{},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/permissions/users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &PermissionsUsersOption{
		ProjectKey: "my-project",
		Permission: "admin",
	}

	_, resp, err := client.Permissions.Users(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPermissions_Users_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid permission should fail validation.
	_, _, err := client.Permissions.Users(&PermissionsUsersOption{
		Permission: "invalid",
	})
	assert.Error(t, err)

	// Query too short should fail validation.
	_, _, err = client.Permissions.Users(&PermissionsUsersOption{
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
