package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// AddApplication / RemoveApplication
// -----------------------------------------------------------------------------

func TestViewsService_AddApplication(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/add_application", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.AddApplication(context.Background(), &ViewsAddApplicationOptions{
		Portfolio:   "my-portfolio",
		Application: "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_AddApplication_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.AddApplication(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddApplication(context.Background(), &ViewsAddApplicationOptions{Application: "my-app"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddApplication(context.Background(), &ViewsAddApplicationOptions{Portfolio: "my-portfolio"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_RemoveApplication(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/remove_application", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.RemoveApplication(context.Background(), &ViewsRemoveApplicationOptions{
		Portfolio:   "my-portfolio",
		Application: "my-application",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_RemoveApplication_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.RemoveApplication(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.RemoveApplication(context.Background(), &ViewsRemoveApplicationOptions{Application: "my-app"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.RemoveApplication(context.Background(), &ViewsRemoveApplicationOptions{Portfolio: "my-portfolio"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// AddApplicationBranch / RemoveApplicationBranch
// -----------------------------------------------------------------------------

func TestViewsService_AddApplicationBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/add_application_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.AddApplicationBranch(context.Background(), &ViewsAddApplicationBranchOptions{
		Application: "my-app",
		Branch:      "main",
		Key:         "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_AddApplicationBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.AddApplicationBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddApplicationBranch(context.Background(), &ViewsAddApplicationBranchOptions{Branch: "main", Key: "pf"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddApplicationBranch(context.Background(), &ViewsAddApplicationBranchOptions{Application: "app", Key: "pf"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddApplicationBranch(context.Background(), &ViewsAddApplicationBranchOptions{Application: "app", Branch: "main"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_RemoveApplicationBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/remove_application_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.RemoveApplicationBranch(context.Background(), &ViewsRemoveApplicationBranchOptions{
		Application: "my-app",
		Branch:      "main",
		Key:         "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_RemoveApplicationBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.RemoveApplicationBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// AddPortfolio / RemovePortfolio
// -----------------------------------------------------------------------------

func TestViewsService_AddPortfolio(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/add_portfolio", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.AddPortfolio(context.Background(), &ViewsAddPortfolioOptions{
		Portfolio: "parent-portfolio",
		Reference: "ref-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_AddPortfolio_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.AddPortfolio(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddPortfolio(context.Background(), &ViewsAddPortfolioOptions{Reference: "ref"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddPortfolio(context.Background(), &ViewsAddPortfolioOptions{Portfolio: "parent"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_RemovePortfolio(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/remove_portfolio", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.RemovePortfolio(context.Background(), &ViewsRemovePortfolioOptions{
		Portfolio: "parent-portfolio",
		Reference: "sub-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_RemovePortfolio_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.RemovePortfolio(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.RemovePortfolio(context.Background(), &ViewsRemovePortfolioOptions{Reference: "sub"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.RemovePortfolio(context.Background(), &ViewsRemovePortfolioOptions{Portfolio: "parent"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// AddProject / RemoveProject
// -----------------------------------------------------------------------------

func TestViewsService_AddProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/add_project", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.AddProject(context.Background(), &ViewsAddProjectOptions{
		Key:     "my-portfolio",
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_AddProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.AddProject(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddProject(context.Background(), &ViewsAddProjectOptions{Project: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddProject(context.Background(), &ViewsAddProjectOptions{Key: "my-portfolio"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_RemoveProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/remove_project", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.RemoveProject(context.Background(), &ViewsRemoveProjectOptions{
		Key:     "my-portfolio",
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_RemoveProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.RemoveProject(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.RemoveProject(context.Background(), &ViewsRemoveProjectOptions{Project: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.RemoveProject(context.Background(), &ViewsRemoveProjectOptions{Key: "my-portfolio"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// AddProjectBranch / RemoveProjectBranch
// -----------------------------------------------------------------------------

func TestViewsService_AddProjectBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/add_project_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.AddProjectBranch(context.Background(), &ViewsAddProjectBranchOptions{
		Branch:  "main",
		Key:     "my-portfolio",
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_AddProjectBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.AddProjectBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.AddProjectBranch(context.Background(), &ViewsAddProjectBranchOptions{Key: "pf", Project: "proj"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_RemoveProjectBranch(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/remove_project_branch", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.RemoveProjectBranch(context.Background(), &ViewsRemoveProjectBranchOptions{
		Branch:  "main",
		Key:     "my-portfolio",
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_RemoveProjectBranch_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.RemoveProjectBranch(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Applications / SubPortfolios
// -----------------------------------------------------------------------------

func TestViewsService_Applications(t *testing.T) {
	response := ViewsApplications{
		Applications: []ViewApplication{
			{Key: "app-1", Name: "Application 1"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/applications", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.Applications(context.Background(), &ViewsApplicationsOptions{
		Portfolio: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Applications, 1)
	assert.Equal(t, "app-1", result.Applications[0].Key)
}

func TestViewsService_Applications_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.Applications(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.Views.Applications(context.Background(), &ViewsApplicationsOptions{})
	assert.Error(t, err)
}

func TestViewsService_SubPortfolios(t *testing.T) {
	response := ViewsSubViews{
		SubViews: []View{
			{Key: "sub-1", Name: "Sub Portfolio 1", Qualifier: "SVW"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/portfolios", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.SubPortfolios(context.Background(), &ViewsSubViewsOptions{
		Portfolio: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.SubViews, 1)
}

func TestViewsService_SubPortfolios_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.SubPortfolios(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.Views.SubPortfolios(context.Background(), &ViewsSubViewsOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

func TestViewsService_Create(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/create", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Create(context.Background(), &ViewsCreateOptions{
		Name: "My Portfolio",
		Key:  "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Create_WithParent(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/create", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Create(context.Background(), &ViewsCreateOptions{
		Name:   "Sub Portfolio",
		Parent: "parent-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Create_WithVisibility(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/create", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Create(context.Background(), &ViewsCreateOptions{
		Name:        "My Portfolio",
		Description: "A test portfolio",
		Visibility:  "private",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.Create(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Create(context.Background(), &ViewsCreateOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Create(context.Background(), &ViewsCreateOptions{Name: "My Portfolio", Visibility: "invalid"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

func TestViewsService_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/delete", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Delete(context.Background(), &ViewsDeleteOptions{
		Key: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.Delete(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Delete(context.Background(), &ViewsDeleteOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// List
// -----------------------------------------------------------------------------

func TestViewsService_List(t *testing.T) {
	response := ViewsList{
		Views: []View{
			{Key: "pf-1", Name: "Portfolio 1", Qualifier: "VW"},
			{Key: "pf-2", Name: "Portfolio 2", Qualifier: "VW"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.List(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Views, 2)
	assert.Equal(t, "pf-1", result.Views[0].Key)
	assert.Equal(t, "Portfolio 1", result.Views[0].Name)
}

func TestViewsService_List_Empty(t *testing.T) {
	response := ViewsList{}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.List(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Empty(t, result.Views)
}

// -----------------------------------------------------------------------------
// Move
// -----------------------------------------------------------------------------

func TestViewsService_Move(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/move", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Move(context.Background(), &ViewsMoveOptions{
		Key:         "my-portfolio",
		Destination: "destination-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Move_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.Move(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Move(context.Background(), &ViewsMoveOptions{Destination: "dest"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Move(context.Background(), &ViewsMoveOptions{Key: "my-portfolio"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// MoveOptions
// -----------------------------------------------------------------------------

func TestViewsService_MoveOptions(t *testing.T) {
	response := ViewsMoveDestinations{
		Views: []ViewDestination{
			{Key: "dest-1", Name: "Destination 1"},
			{Key: "dest-2", Name: "Destination 2"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/move_options", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.MoveOptions(context.Background(), &ViewsMoveOptionsOptions{
		Key: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Views, 2)
	assert.Equal(t, "dest-1", result.Views[0].Key)
}

func TestViewsService_MoveOptions_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.MoveOptions(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.Views.MoveOptions(context.Background(), &ViewsMoveOptionsOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Projects / ProjectsStatus
// -----------------------------------------------------------------------------

func TestViewsService_Projects(t *testing.T) {
	response := ViewsProjects{
		Projects: []ViewProject{
			{Key: "proj-1", Name: "Project 1", Selected: true},
			{Key: "proj-2", Name: "Project 2", Selected: false},
		},
		Paging: Paging{PageIndex: 1, PageSize: 100, Total: 2},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/projects", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.Projects(context.Background(), &ViewsProjectsOptions{
		Key:      "my-portfolio",
		Selected: "all",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Projects, 2)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestViewsService_Projects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.Projects(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.Views.Projects(context.Background(), &ViewsProjectsOptions{})
	assert.Error(t, err)

	_, _, err = client.Views.Projects(context.Background(), &ViewsProjectsOptions{Key: "pf", Selected: "invalid"})
	assert.Error(t, err)
}

func TestViewsService_ProjectsStatus(t *testing.T) {
	response := ViewsProjectsStatus{
		Projects: []ViewProjectStatus{
			{Key: "proj-1", Name: "Project 1", Status: "OK"},
		},
		Paging: Paging{PageIndex: 1, PageSize: 100, Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/projects_status", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.ProjectsStatus(context.Background(), &ViewsProjectsStatusOptions{
		Portfolio: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Projects, 1)
	assert.Equal(t, "OK", result.Projects[0].Status)
}

func TestViewsService_ProjectsStatus_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.ProjectsStatus(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.Views.ProjectsStatus(context.Background(), &ViewsProjectsStatusOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Refresh
// -----------------------------------------------------------------------------

func TestViewsService_Refresh(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/refresh", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Refresh(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Refresh_WithKey(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/refresh", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Refresh(context.Background(), &ViewsRefreshOptions{Key: "my-portfolio"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Search
// -----------------------------------------------------------------------------

func TestViewsService_Search(t *testing.T) {
	response := ViewsSearch{
		Components: []View{
			{Key: "pf-1", Name: "Portfolio 1"},
		},
		Paging: Paging{PageIndex: 1, PageSize: 100, Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.Search(context.Background(), &ViewsSearchOptions{Query: "Portfolio"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Components, 1)
	assert.Equal(t, int64(1), result.Paging.Total)
}

func TestViewsService_Search_Nil(t *testing.T) {
	response := ViewsSearch{}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.Search(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestViewsService_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.Search(context.Background(), &ViewsSearchOptions{PaginationArgs: PaginationArgs{Page: -1}})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SetManualMode / SetNoneMode / SetRegexpMode / SetRemainingProjectsMode / SetTagsMode
// -----------------------------------------------------------------------------

func TestViewsService_SetManualMode(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/set_manual_mode", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.SetManualMode(context.Background(), &ViewsSetManualModeOptions{Portfolio: "my-portfolio"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_SetManualMode_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.SetManualMode(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.SetManualMode(context.Background(), &ViewsSetManualModeOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_SetNoneMode(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/set_none_mode", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.SetNoneMode(context.Background(), &ViewsSetNoneModeOptions{Portfolio: "my-portfolio"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_SetNoneMode_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.SetNoneMode(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_SetRegexpMode(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/set_regexp_mode", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.SetRegexpMode(context.Background(), &ViewsSetRegexpModeOptions{
		Portfolio: "my-portfolio",
		Regexp:    ".*my-project.*",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_SetRegexpMode_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.SetRegexpMode(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.SetRegexpMode(context.Background(), &ViewsSetRegexpModeOptions{Portfolio: "pf"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.SetRegexpMode(context.Background(), &ViewsSetRegexpModeOptions{Regexp: ".*"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_SetRemainingProjectsMode(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/set_remaining_projects_mode", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.SetRemainingProjectsMode(context.Background(), &ViewsSetRemainingProjectsModeOptions{
		Portfolio: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_SetRemainingProjectsMode_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.SetRemainingProjectsMode(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewsService_SetTagsMode(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/set_tags_mode", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.SetTagsMode(context.Background(), &ViewsSetTagsModeOptions{
		Portfolio: "my-portfolio",
		Tags:      []string{"java", "security"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_SetTagsMode_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.SetTagsMode(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.SetTagsMode(context.Background(), &ViewsSetTagsModeOptions{Portfolio: "pf"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Show
// -----------------------------------------------------------------------------

func TestViewsService_Show(t *testing.T) {
	response := ViewsShow{
		Portfolio: ViewDetails{
			Key:           "pf-1",
			Name:          "Portfolio 1",
			Qualifier:     "VW",
			SelectionMode: "MANUAL",
			SubViews: []View{
				{Key: "sub-pf-1", Name: "Sub Portfolio 1", Qualifier: "SVW"},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/views/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Views.Show(context.Background(), &ViewsShowOptions{Key: "pf-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "pf-1", result.Portfolio.Key)
	assert.Equal(t, "Portfolio 1", result.Portfolio.Name)
	assert.Len(t, result.Portfolio.SubViews, 1)
}

func TestViewsService_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Views.Show(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = client.Views.Show(context.Background(), &ViewsShowOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

func TestViewsService_Update(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/views/update", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Views.Update(context.Background(), &ViewsUpdateOptions{
		Key:  "my-portfolio",
		Name: "My Renamed Portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestViewsService_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.Views.Update(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Update(context.Background(), &ViewsUpdateOptions{Name: "New Name"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.Views.Update(context.Background(), &ViewsUpdateOptions{Key: "my-portfolio"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Validation method tests
// -----------------------------------------------------------------------------

func TestViewsService_ValidateCreateOpt(t *testing.T) {
	client := newLocalhostClient(t)

	err := client.Views.ValidateCreateOpt(&ViewsCreateOptions{Name: "My Portfolio"})
	assert.NoError(t, err)

	err = client.Views.ValidateCreateOpt(nil)
	assert.Error(t, err)

	err = client.Views.ValidateCreateOpt(&ViewsCreateOptions{})
	assert.Error(t, err)

	err = client.Views.ValidateCreateOpt(&ViewsCreateOptions{Name: "X", Visibility: "bad"})
	assert.Error(t, err)
}

func TestViewsService_ValidateProjectsOpt(t *testing.T) {
	client := newLocalhostClient(t)

	err := client.Views.ValidateProjectsOpt(&ViewsProjectsOptions{Key: "pf", Selected: "selected"})
	assert.NoError(t, err)

	err = client.Views.ValidateProjectsOpt(nil)
	assert.Error(t, err)

	err = client.Views.ValidateProjectsOpt(&ViewsProjectsOptions{Key: "pf", Selected: "bad-value"})
	assert.Error(t, err)
}
