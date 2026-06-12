package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

func TestApplicationsService_Create(t *testing.T) {
	response := ApplicationsCreate{
		Application: Application{
			Key:  "my-application",
			Name: "My Application",
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/applications/create", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Applications.Create(context.Background(), &ApplicationsCreateOptions{
		Name: "My Application",
		Key:  "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my-application", result.Application.Key)
}

func TestApplicationsService_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.Applications.Create(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.Applications.Create(context.Background(), &ApplicationsCreateOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.Applications.Create(context.Background(), &ApplicationsCreateOptions{
		Name:       "My App",
		Visibility: "invalid",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

func TestApplicationsService_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/delete", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.Delete(context.Background(), &ApplicationsDeleteOptions{
		Application: "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.Delete(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.Delete(context.Background(), &ApplicationsDeleteOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Show
// -----------------------------------------------------------------------------

func TestApplicationsService_Show(t *testing.T) {
	response := ApplicationsShow{
		Application: ApplicationDetails{
			Key:  "my-application",
			Name: "My Application",
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/applications/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Applications.Show(context.Background(), &ApplicationsShowOptions{
		Application: "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my-application", result.Application.Key)
}

func TestApplicationsService_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.Applications.Show(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

func TestApplicationsService_Update(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/update", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.Update(context.Background(), &ApplicationsUpdateOptions{
		Application: "my-application",
		Name:        "Updated Application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.Update(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.Update(context.Background(), &ApplicationsUpdateOptions{Name: "name"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.Update(context.Background(), &ApplicationsUpdateOptions{Application: "app"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// AddProject / RemoveProject
// -----------------------------------------------------------------------------

func TestApplicationsService_AddProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/add_project", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.AddProject(context.Background(), &ApplicationsAddProjectOptions{
		Application: "my-application",
		Project:     "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_AddProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.AddProject(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.AddProject(context.Background(), &ApplicationsAddProjectOptions{Project: "proj"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestApplicationsService_RemoveProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/remove_project", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.RemoveProject(context.Background(), &ApplicationsRemoveProjectOptions{
		Application: "my-application",
		Project:     "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_RemoveProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.RemoveProject(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.RemoveProject(context.Background(), &ApplicationsRemoveProjectOptions{Project: "proj"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.RemoveProject(context.Background(), &ApplicationsRemoveProjectOptions{Application: "app"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// CreateBranch / DeleteBranch / UpdateBranch
// -----------------------------------------------------------------------------

func TestApplicationsService_CreateBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/create_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.CreateBranch(context.Background(), &ApplicationsCreateBranchOptions{
		Application:   "my-application",
		Branch:        "feature-branch",
		Project:       []string{"proj1"},
		ProjectBranch: []string{"main"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_CreateBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.CreateBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.CreateBranch(context.Background(), &ApplicationsCreateBranchOptions{Branch: "b"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.CreateBranch(context.Background(), &ApplicationsCreateBranchOptions{
		Application: "app",
		Branch:      "b",
	})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.CreateBranch(context.Background(), &ApplicationsCreateBranchOptions{
		Application: "app",
		Branch:      "b",
		Project:     []string{"proj1"},
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestApplicationsService_DeleteBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/delete_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.DeleteBranch(context.Background(), &ApplicationsDeleteBranchOptions{
		Application: "my-application",
		Branch:      "feature-branch",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_DeleteBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.DeleteBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.DeleteBranch(context.Background(), &ApplicationsDeleteBranchOptions{Branch: "b"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.DeleteBranch(context.Background(), &ApplicationsDeleteBranchOptions{Application: "app"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestApplicationsService_UpdateBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/update_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.UpdateBranch(context.Background(), &ApplicationsUpdateBranchOptions{
		Application:   "my-application",
		Branch:        "feature-branch",
		Name:          "new-branch-name",
		Project:       []string{"proj1"},
		ProjectBranch: []string{"main"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_UpdateBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.UpdateBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.UpdateBranch(context.Background(), &ApplicationsUpdateBranchOptions{
		Application: "app",
		Branch:      "branch",
	})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.UpdateBranch(context.Background(), &ApplicationsUpdateBranchOptions{
		Application: "app",
		Branch:      "branch",
		Name:        "new-name",
	})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.UpdateBranch(context.Background(), &ApplicationsUpdateBranchOptions{
		Application: "app",
		Branch:      "branch",
		Name:        "new-name",
		Project:     []string{"proj1"},
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// SetTags
// -----------------------------------------------------------------------------

func TestApplicationsService_SetTags(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/set_tags", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.SetTags(context.Background(), &ApplicationsSetTagsOptions{
		Application: "my-application",
		Tags:        "tag1,tag2",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestApplicationsService_SetTags_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Applications.SetTags(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Applications.SetTags(context.Background(), &ApplicationsSetTagsOptions{Application: "app"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// SearchProjects
// -----------------------------------------------------------------------------

func TestApplicationsService_SearchProjects(t *testing.T) {
	response := ApplicationsSearchProjects{
		Projects: []ApplicationProject{
			{Key: "proj1", Name: "Project 1", Selected: true},
		},
		Paging: Paging{Total: 1, PageIndex: 1, PageSize: 10},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/applications/search_projects", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Applications.SearchProjects(context.Background(), &ApplicationsSearchProjectsOptions{
		Application: "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Projects, 1)
}

func TestApplicationsService_SearchProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.Applications.SearchProjects(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.Applications.SearchProjects(context.Background(), &ApplicationsSearchProjectsOptions{
		Application: "app",
		Selected:    "invalid",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// SearchAllProjects
// -----------------------------------------------------------------------------

func TestApplicationsService_SearchAllProjects(t *testing.T) {
	response := ApplicationsSearchProjects{
		Projects: []ApplicationProject{
			{Key: "proj1", Name: "Project 1", Selected: true},
			{Key: "proj2", Name: "Project 2", Selected: false},
		},
		Paging: Paging{Total: 2, PageIndex: 1, PageSize: 100},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/applications/search_projects", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	results, resp, err := client.Applications.SearchAllProjects(context.Background(), &ApplicationsSearchProjectsOptions{
		Application: "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, results, 2)
}

// -----------------------------------------------------------------------------
// Refresh
// -----------------------------------------------------------------------------

func TestApplicationsService_Refresh(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/applications/refresh", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Applications.Refresh(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
