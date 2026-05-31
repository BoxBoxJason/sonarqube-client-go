package sonar

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// CheckPat Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_CheckPat(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/check_pat", http.StatusOK, "{}")
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsCheckPatOptions{
		AlmSetting: "my-azure-setting",
	}

	_, resp, err := client.AlmIntegrations.CheckPat(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAlmIntegrations_CheckPat_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.CheckPat(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.CheckPat(context.Background(), &AlmIntegrationsCheckPatOptions{})
	assert.Error(t, err)

	// Test AlmSetting too long
	_, _, err = client.AlmIntegrations.CheckPat(context.Background(), &AlmIntegrationsCheckPatOptions{
		AlmSetting: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// GetGithubClientID Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_GetGithubClientID(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/get_github_client_id", http.StatusOK, `{"clientId":"my-client-id"}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsGetGithubClientIDOptions{
		AlmSetting: "my-github-setting",
	}

	result, resp, err := client.AlmIntegrations.GetGithubClientID(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Equal(t, "my-client-id", result.ClientID)
}

func TestAlmIntegrations_GetGithubClientID_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.GetGithubClientID(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.GetGithubClientID(context.Background(), &AlmIntegrationsGetGithubClientIDOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ImportAzureProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportAzureProject(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_azure_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportAzureProjectOptions{
		ProjectName:    "my-azure-project",
		RepositoryName: "my-azure-repo",
	}

	resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportAzureProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportAzureProject(context.Background(), nil)
	assert.Error(t, err)

	// Test missing ProjectName
	_, err = client.AlmIntegrations.ImportAzureProject(context.Background(), &AlmIntegrationsImportAzureProjectOptions{
		RepositoryName: "repo",
	})
	assert.Error(t, err)

	// Test missing RepositoryName
	_, err = client.AlmIntegrations.ImportAzureProject(context.Background(), &AlmIntegrationsImportAzureProjectOptions{
		ProjectName: "project",
	})
	assert.Error(t, err)

	// Test invalid NewCodeDefinitionType
	_, err = client.AlmIntegrations.ImportAzureProject(context.Background(), &AlmIntegrationsImportAzureProjectOptions{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "INVALID_TYPE",
	})
	assert.Error(t, err)

	// Test NUMBER_OF_DAYS without value
	_, err = client.AlmIntegrations.ImportAzureProject(context.Background(), &AlmIntegrationsImportAzureProjectOptions{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: NewCodePeriodTypeNumberOfDays,
	})
	assert.Error(t, err)

	// Test PREVIOUS_VERSION with value
	_, err = client.AlmIntegrations.ImportAzureProject(context.Background(), &AlmIntegrationsImportAzureProjectOptions{
		ProjectName:            "project",
		RepositoryName:         "repo",
		NewCodeDefinitionType:  NewCodePeriodTypePreviousVersion,
		NewCodeDefinitionValue: 30,
	})
	assert.Error(t, err)
}

func TestAlmIntegrations_ImportAzureProject_WithNewCodeDefinition(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_azure_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	// Test with PREVIOUS_VERSION
	opt := &AlmIntegrationsImportAzureProjectOptions{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: NewCodePeriodTypePreviousVersion,
	}

	resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Test with NUMBER_OF_DAYS
	opt = &AlmIntegrationsImportAzureProjectOptions{
		ProjectName:            "project",
		RepositoryName:         "repo",
		NewCodeDefinitionType:  NewCodePeriodTypeNumberOfDays,
		NewCodeDefinitionValue: 30,
	}

	resp, err = client.AlmIntegrations.ImportAzureProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// ImportBitbucketCloudRepo Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportBitbucketCloudRepo(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_bitbucketcloud_repo", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportBitbucketCloudRepoOptions{
		RepositorySlug: "my-repo",
	}

	resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportBitbucketCloudRepo_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportBitbucketCloudRepo(context.Background(), nil)
	assert.Error(t, err)

	// Test missing RepositorySlug
	_, err = client.AlmIntegrations.ImportBitbucketCloudRepo(context.Background(), &AlmIntegrationsImportBitbucketCloudRepoOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ImportBitbucketServerProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportBitbucketServerProject(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_bitbucketserver_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportBitbucketServerProjectOptions{
		ProjectKey:     "PRJ",
		RepositorySlug: "my-repo",
	}

	resp, err := client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportBitbucketServerProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), nil)
	assert.Error(t, err)

	// Test missing ProjectKey
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), &AlmIntegrationsImportBitbucketServerProjectOptions{
		RepositorySlug: "repo",
	})
	assert.Error(t, err)

	// Test missing RepositorySlug
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), &AlmIntegrationsImportBitbucketServerProjectOptions{
		ProjectKey: "PRJ",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ImportGithubProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportGithubProject(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_github_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportGithubProjectOptions{
		RepositoryKey: "octocat/hello-world",
	}

	resp, err := client.AlmIntegrations.ImportGithubProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportGithubProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportGithubProject(context.Background(), nil)
	assert.Error(t, err)

	// Test missing RepositoryKey
	_, err = client.AlmIntegrations.ImportGithubProject(context.Background(), &AlmIntegrationsImportGithubProjectOptions{})
	assert.Error(t, err)

	// Test RepositoryKey too long
	_, err = client.AlmIntegrations.ImportGithubProject(context.Background(), &AlmIntegrationsImportGithubProjectOptions{
		RepositoryKey: strings.Repeat("a", MaxGitHubRepoKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ImportGitlabProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportGitlabProject(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_gitlab_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportGitlabProjectOptions{
		GitlabProjectId: "12345",
	}

	resp, err := client.AlmIntegrations.ImportGitlabProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportGitlabProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportGitlabProject(context.Background(), nil)
	assert.Error(t, err)

	// Test missing GitlabProjectId
	_, err = client.AlmIntegrations.ImportGitlabProject(context.Background(), &AlmIntegrationsImportGitlabProjectOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ListAzureProjects Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListAzureProjects(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/list_azure_projects", http.StatusOK, `{"projects":[{"name":"Project1","description":"Description1"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsListAzureProjectsOptions{
		AlmSetting: "my-azure-setting",
	}

	result, resp, err := client.AlmIntegrations.ListAzureProjects(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Projects, 1)
}

func TestAlmIntegrations_ListAzureProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListAzureProjects(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListAzureProjects(context.Background(), &AlmIntegrationsListAzureProjectsOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ListBitbucketServerProjects Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListBitbucketServerProjects(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/list_bitbucketserver_projects", http.StatusOK, `{"projects":[{"key":"PRJ","name":"Project1"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsListBitbucketServerProjectsOptions{
		AlmSetting: "my-bitbucket-setting",
		PageSize:   25,
	}

	result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Projects, 1)
}

func TestAlmIntegrations_ListBitbucketServerProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListBitbucketServerProjects(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(context.Background(), &AlmIntegrationsListBitbucketServerProjectsOptions{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(context.Background(), &AlmIntegrationsListBitbucketServerProjectsOptions{
		AlmSetting: "setting",
		PageSize:   MaxPageSizeAlmIntegrations + 1,
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ListGithubOrganizations Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListGithubOrganizations(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/list_github_organizations", http.StatusOK, `{"organizations":[{"key":"octocat","name":"Octocat"}],"paging":{"pageIndex":1,"pageSize":100,"total":1}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsListGithubOrganizationsOptions{
		AlmSetting: "my-github-setting",
	}

	result, resp, err := client.AlmIntegrations.ListGithubOrganizations(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Organizations, 1)
}

func TestAlmIntegrations_ListGithubOrganizations_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListGithubOrganizations(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(context.Background(), &AlmIntegrationsListGithubOrganizationsOptions{})
	assert.Error(t, err)

	// Test Token too long
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(context.Background(), &AlmIntegrationsListGithubOrganizationsOptions{
		AlmSetting: "setting",
		Token:      strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ListGithubRepositories Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListGithubRepositories(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/list_github_repositories", http.StatusOK, `{"repositories":[{"id":1,"key":"octocat/hello-world","name":"hello-world","url":"https://github.com/octocat/hello-world"}],"paging":{"pageIndex":1,"pageSize":100,"total":1}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsListGithubRepositoriesOptions{
		AlmSetting:   "my-github-setting",
		Organization: "octocat",
	}

	result, resp, err := client.AlmIntegrations.ListGithubRepositories(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_ListGithubRepositories_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListGithubRepositories(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListGithubRepositories(context.Background(), &AlmIntegrationsListGithubRepositoriesOptions{
		Organization: "octocat",
	})
	assert.Error(t, err)

	// Test missing Organization
	_, _, err = client.AlmIntegrations.ListGithubRepositories(context.Background(), &AlmIntegrationsListGithubRepositoriesOptions{
		AlmSetting: "setting",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SearchAzureRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchAzureRepos(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/search_azure_repos", http.StatusOK, `{"repositories":[{"name":"repo1","projectName":"Project1"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsSearchAzureReposOptions{
		AlmSetting: "my-azure-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchAzureRepos(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_SearchAzureRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchAzureRepos(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchAzureRepos(context.Background(), &AlmIntegrationsSearchAzureReposOptions{})
	assert.Error(t, err)

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchAzureRepos(context.Background(), &AlmIntegrationsSearchAzureReposOptions{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)

	// Test SearchQuery too long
	_, _, err = client.AlmIntegrations.SearchAzureRepos(context.Background(), &AlmIntegrationsSearchAzureReposOptions{
		AlmSetting:  "setting",
		SearchQuery: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SearchBitbucketCloudRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchBitbucketCloudRepos(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/search_bitbucketcloud_repos", http.StatusOK, `{"isLastPage":true,"repositories":[{"name":"repo1","slug":"repo1","uuid":"uuid1"}],"paging":{"pageIndex":1,"pageSize":20}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsSearchBitbucketCloudReposOptions{
		AlmSetting: "my-bitbucket-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
	assert.True(t, result.IsLastPage)
}

func TestAlmIntegrations_SearchBitbucketCloudRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), &AlmIntegrationsSearchBitbucketCloudReposOptions{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), &AlmIntegrationsSearchBitbucketCloudReposOptions{
		AlmSetting: "setting",
		PaginationArgs: PaginationArgs{
			PageSize: MaxPageSizeAlmIntegrations + 1,
		},
	})
	assert.Error(t, err)

	// Test RepositoryName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), &AlmIntegrationsSearchBitbucketCloudReposOptions{
		AlmSetting:     "setting",
		RepositoryName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SearchBitbucketServerRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchBitbucketServerRepos(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/search_bitbucketserver_repos", http.StatusOK, `{"isLastPage":false,"repositories":[{"name":"repo1","slug":"repo1","projectKey":"PRJ"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsSearchBitbucketServerReposOptions{
		AlmSetting: "my-bitbucket-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_SearchBitbucketServerRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), &AlmIntegrationsSearchBitbucketServerReposOptions{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), &AlmIntegrationsSearchBitbucketServerReposOptions{
		AlmSetting: "setting",
		PageSize:   MaxPageSizeAlmIntegrations + 1,
	})
	assert.Error(t, err)

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), &AlmIntegrationsSearchBitbucketServerReposOptions{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)

	// Test RepositoryName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), &AlmIntegrationsSearchBitbucketServerReposOptions{
		AlmSetting:     "setting",
		RepositoryName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SearchGitlabRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchGitlabRepos(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/search_gitlab_repos", http.StatusOK, `{"repositories":[{"id":1,"name":"project1","pathName":"group/project1","slug":"project1","url":"https://gitlab.com/group/project1"}],"paging":{"pageIndex":1,"pageSize":20,"total":1}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsSearchGitlabReposOptions{
		AlmSetting: "my-gitlab-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchGitlabRepos(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_SearchGitlabRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchGitlabRepos(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(context.Background(), &AlmIntegrationsSearchGitlabReposOptions{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(context.Background(), &AlmIntegrationsSearchGitlabReposOptions{
		AlmSetting: "setting",
		PaginationArgs: PaginationArgs{
			PageSize: MaxPageSizeAlmIntegrations + 1,
		},
	})
	assert.Error(t, err)

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(context.Background(), &AlmIntegrationsSearchGitlabReposOptions{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// SetPat Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SetPat(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/set_pat", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsSetPatOptions{
		AlmSetting: "my-setting",
		Pat:        "my-personal-access-token",
	}

	resp, err := client.AlmIntegrations.SetPat(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_SetPat_WithUsername(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/set_pat", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	// Test with username (for Bitbucket Cloud)
	opt := &AlmIntegrationsSetPatOptions{
		AlmSetting: "my-bitbucket-cloud-setting",
		Pat:        "my-app-password",
		Username:   "my-username",
	}

	resp, err := client.AlmIntegrations.SetPat(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_SetPat_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.SetPat(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Pat
	_, err = client.AlmIntegrations.SetPat(context.Background(), &AlmIntegrationsSetPatOptions{
		AlmSetting: "setting",
	})
	assert.Error(t, err)

	// Test Pat too long
	_, err = client.AlmIntegrations.SetPat(context.Background(), &AlmIntegrationsSetPatOptions{
		Pat: strings.Repeat("a", MaxPatLength+1),
	})
	assert.Error(t, err)

	// Test Username too long
	_, err = client.AlmIntegrations.SetPat(context.Background(), &AlmIntegrationsSetPatOptions{
		Pat:      "token",
		Username: strings.Repeat("a", MaxUsernameLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Helper Function Tests
// -----------------------------------------------------------------------------

func TestValidateNewCodeDefinition(t *testing.T) {
	tests := []struct {
		name            string
		definitionType  string
		definitionValue int64
		wantErr         bool
	}{
		{
			name:            "empty type is valid",
			definitionType:  "",
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "PREVIOUS_VERSION without value",
			definitionType:  NewCodePeriodTypePreviousVersion,
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "PREVIOUS_VERSION with value is invalid",
			definitionType:  NewCodePeriodTypePreviousVersion,
			definitionValue: 30,
			wantErr:         true,
		},
		{
			name:            "REFERENCE_BRANCH without value",
			definitionType:  NewCodePeriodTypeReferenceBranch,
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "REFERENCE_BRANCH with value is invalid",
			definitionType:  NewCodePeriodTypeReferenceBranch,
			definitionValue: 30,
			wantErr:         true,
		},
		{
			name:            "NUMBER_OF_DAYS with value",
			definitionType:  NewCodePeriodTypeNumberOfDays,
			definitionValue: 30,
			wantErr:         false,
		},
		{
			name:            "NUMBER_OF_DAYS without value is invalid",
			definitionType:  NewCodePeriodTypeNumberOfDays,
			definitionValue: 0,
			wantErr:         true,
		},
		{
			name:            "invalid type",
			definitionType:  "INVALID_TYPE",
			definitionValue: 0,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNewCodeDefinition(tt.definitionType, tt.definitionValue)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
