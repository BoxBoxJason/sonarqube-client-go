package sonargo

import (
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

	opt := &AlmIntegrationsCheckPatOption{
		AlmSetting: "my-azure-setting",
	}

	_, resp, err := client.AlmIntegrations.CheckPat(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAlmIntegrations_CheckPat_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.CheckPat(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.CheckPat(&AlmIntegrationsCheckPatOption{})
	assert.Error(t, err)

	// Test AlmSetting too long
	_, _, err = client.AlmIntegrations.CheckPat(&AlmIntegrationsCheckPatOption{
		AlmSetting: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// GetGithubClientId Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_GetGithubClientId(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/get_github_client_id", http.StatusOK, `{"clientId":"my-client-id"}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsGetGithubClientIdOption{
		AlmSetting: "my-github-setting",
	}

	result, resp, err := client.AlmIntegrations.GetGithubClientId(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Equal(t, "my-client-id", result.ClientID)
}

func TestAlmIntegrations_GetGithubClientId_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.GetGithubClientId(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.GetGithubClientId(&AlmIntegrationsGetGithubClientIdOption{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ImportAzureProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportAzureProject(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_azure_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportAzureProjectOption{
		ProjectName:    "my-azure-project",
		RepositoryName: "my-azure-repo",
	}

	resp, err := client.AlmIntegrations.ImportAzureProject(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportAzureProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportAzureProject(nil)
	assert.Error(t, err)

	// Test missing ProjectName
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		RepositoryName: "repo",
	})
	assert.Error(t, err)

	// Test missing RepositoryName
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName: "project",
	})
	assert.Error(t, err)

	// Test invalid NewCodeDefinitionType
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "INVALID_TYPE",
	})
	assert.Error(t, err)

	// Test NUMBER_OF_DAYS without value
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "NUMBER_OF_DAYS",
	})
	assert.Error(t, err)

	// Test PREVIOUS_VERSION with value
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName:            "project",
		RepositoryName:         "repo",
		NewCodeDefinitionType:  "PREVIOUS_VERSION",
		NewCodeDefinitionValue: 30,
	})
	assert.Error(t, err)
}

func TestAlmIntegrations_ImportAzureProject_WithNewCodeDefinition(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_azure_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	// Test with PREVIOUS_VERSION
	opt := &AlmIntegrationsImportAzureProjectOption{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "PREVIOUS_VERSION",
	}

	resp, err := client.AlmIntegrations.ImportAzureProject(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Test with NUMBER_OF_DAYS
	opt = &AlmIntegrationsImportAzureProjectOption{
		ProjectName:            "project",
		RepositoryName:         "repo",
		NewCodeDefinitionType:  "NUMBER_OF_DAYS",
		NewCodeDefinitionValue: 30,
	}

	resp, err = client.AlmIntegrations.ImportAzureProject(opt)
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

	opt := &AlmIntegrationsImportBitbucketCloudRepoOption{
		RepositorySlug: "my-repo",
	}

	resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportBitbucketCloudRepo_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportBitbucketCloudRepo(nil)
	assert.Error(t, err)

	// Test missing RepositorySlug
	_, err = client.AlmIntegrations.ImportBitbucketCloudRepo(&AlmIntegrationsImportBitbucketCloudRepoOption{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ImportBitbucketServerProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportBitbucketServerProject(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/import_bitbucketserver_project", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsImportBitbucketServerProjectOption{
		ProjectKey:     "PRJ",
		RepositorySlug: "my-repo",
	}

	resp, err := client.AlmIntegrations.ImportBitbucketServerProject(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportBitbucketServerProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportBitbucketServerProject(nil)
	assert.Error(t, err)

	// Test missing ProjectKey
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(&AlmIntegrationsImportBitbucketServerProjectOption{
		RepositorySlug: "repo",
	})
	assert.Error(t, err)

	// Test missing RepositorySlug
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(&AlmIntegrationsImportBitbucketServerProjectOption{
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

	opt := &AlmIntegrationsImportGithubProjectOption{
		RepositoryKey: "octocat/hello-world",
	}

	resp, err := client.AlmIntegrations.ImportGithubProject(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportGithubProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportGithubProject(nil)
	assert.Error(t, err)

	// Test missing RepositoryKey
	_, err = client.AlmIntegrations.ImportGithubProject(&AlmIntegrationsImportGithubProjectOption{})
	assert.Error(t, err)

	// Test RepositoryKey too long
	_, err = client.AlmIntegrations.ImportGithubProject(&AlmIntegrationsImportGithubProjectOption{
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

	opt := &AlmIntegrationsImportGitlabProjectOption{
		GitlabProjectId: "12345",
	}

	resp, err := client.AlmIntegrations.ImportGitlabProject(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_ImportGitlabProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.ImportGitlabProject(nil)
	assert.Error(t, err)

	// Test missing GitlabProjectId
	_, err = client.AlmIntegrations.ImportGitlabProject(&AlmIntegrationsImportGitlabProjectOption{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ListAzureProjects Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListAzureProjects(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/list_azure_projects", http.StatusOK, `{"projects":[{"name":"Project1","description":"Description1"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsListAzureProjectsOption{
		AlmSetting: "my-azure-setting",
	}

	result, resp, err := client.AlmIntegrations.ListAzureProjects(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Projects, 1)
}

func TestAlmIntegrations_ListAzureProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListAzureProjects(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListAzureProjects(&AlmIntegrationsListAzureProjectsOption{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// ListBitbucketServerProjects Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListBitbucketServerProjects(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/alm_integrations/list_bitbucketserver_projects", http.StatusOK, `{"projects":[{"key":"PRJ","name":"Project1"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AlmIntegrationsListBitbucketServerProjectsOption{
		AlmSetting: "my-bitbucket-setting",
		PageSize:   25,
	}

	result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Projects, 1)
}

func TestAlmIntegrations_ListBitbucketServerProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListBitbucketServerProjects(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(&AlmIntegrationsListBitbucketServerProjectsOption{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(&AlmIntegrationsListBitbucketServerProjectsOption{
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

	opt := &AlmIntegrationsListGithubOrganizationsOption{
		AlmSetting: "my-github-setting",
	}

	result, resp, err := client.AlmIntegrations.ListGithubOrganizations(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Organizations, 1)
}

func TestAlmIntegrations_ListGithubOrganizations_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListGithubOrganizations(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(&AlmIntegrationsListGithubOrganizationsOption{})
	assert.Error(t, err)

	// Test Token too long
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(&AlmIntegrationsListGithubOrganizationsOption{
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

	opt := &AlmIntegrationsListGithubRepositoriesOption{
		AlmSetting:   "my-github-setting",
		Organization: "octocat",
	}

	result, resp, err := client.AlmIntegrations.ListGithubRepositories(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_ListGithubRepositories_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.ListGithubRepositories(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListGithubRepositories(&AlmIntegrationsListGithubRepositoriesOption{
		Organization: "octocat",
	})
	assert.Error(t, err)

	// Test missing Organization
	_, _, err = client.AlmIntegrations.ListGithubRepositories(&AlmIntegrationsListGithubRepositoriesOption{
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

	opt := &AlmIntegrationsSearchAzureReposOption{
		AlmSetting: "my-azure-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchAzureRepos(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_SearchAzureRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchAzureRepos(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchAzureRepos(&AlmIntegrationsSearchAzureReposOption{})
	assert.Error(t, err)

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchAzureRepos(&AlmIntegrationsSearchAzureReposOption{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)

	// Test SearchQuery too long
	_, _, err = client.AlmIntegrations.SearchAzureRepos(&AlmIntegrationsSearchAzureReposOption{
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

	opt := &AlmIntegrationsSearchBitbucketCloudReposOption{
		AlmSetting: "my-bitbucket-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
	assert.True(t, result.IsLastPage)
}

func TestAlmIntegrations_SearchBitbucketCloudRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchBitbucketCloudRepos(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(&AlmIntegrationsSearchBitbucketCloudReposOption{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(&AlmIntegrationsSearchBitbucketCloudReposOption{
		AlmSetting: "setting",
		PaginationArgs: PaginationArgs{
			PageSize: MaxPageSizeAlmIntegrations + 1,
		},
	})
	assert.Error(t, err)

	// Test RepositoryName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(&AlmIntegrationsSearchBitbucketCloudReposOption{
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

	opt := &AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting: "my-bitbucket-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_SearchBitbucketServerRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchBitbucketServerRepos(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting: "setting",
		PageSize:   MaxPageSizeAlmIntegrations + 1,
	})
	assert.Error(t, err)

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	assert.Error(t, err)

	// Test RepositoryName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{
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

	opt := &AlmIntegrationsSearchGitlabReposOption{
		AlmSetting: "my-gitlab-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchGitlabRepos(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestAlmIntegrations_SearchGitlabRepos_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmIntegrations.SearchGitlabRepos(nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(&AlmIntegrationsSearchGitlabReposOption{})
	assert.Error(t, err)

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(&AlmIntegrationsSearchGitlabReposOption{
		AlmSetting: "setting",
		PaginationArgs: PaginationArgs{
			PageSize: MaxPageSizeAlmIntegrations + 1,
		},
	})
	assert.Error(t, err)

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(&AlmIntegrationsSearchGitlabReposOption{
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

	opt := &AlmIntegrationsSetPatOption{
		AlmSetting: "my-setting",
		Pat:        "my-personal-access-token",
	}

	resp, err := client.AlmIntegrations.SetPat(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_SetPat_WithUsername(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/alm_integrations/set_pat", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	// Test with username (for Bitbucket Cloud)
	opt := &AlmIntegrationsSetPatOption{
		AlmSetting: "my-bitbucket-cloud-setting",
		Pat:        "my-app-password",
		Username:   "my-username",
	}

	resp, err := client.AlmIntegrations.SetPat(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmIntegrations_SetPat_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmIntegrations.SetPat(nil)
	assert.Error(t, err)

	// Test missing Pat
	_, err = client.AlmIntegrations.SetPat(&AlmIntegrationsSetPatOption{
		AlmSetting: "setting",
	})
	assert.Error(t, err)

	// Test Pat too long
	_, err = client.AlmIntegrations.SetPat(&AlmIntegrationsSetPatOption{
		Pat: strings.Repeat("a", MaxPatLength+1),
	})
	assert.Error(t, err)

	// Test Username too long
	_, err = client.AlmIntegrations.SetPat(&AlmIntegrationsSetPatOption{
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
			definitionType:  "PREVIOUS_VERSION",
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "PREVIOUS_VERSION with value is invalid",
			definitionType:  "PREVIOUS_VERSION",
			definitionValue: 30,
			wantErr:         true,
		},
		{
			name:            "REFERENCE_BRANCH without value",
			definitionType:  "REFERENCE_BRANCH",
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "REFERENCE_BRANCH with value is invalid",
			definitionType:  "REFERENCE_BRANCH",
			definitionValue: 30,
			wantErr:         true,
		},
		{
			name:            "NUMBER_OF_DAYS with value",
			definitionType:  "NUMBER_OF_DAYS",
			definitionValue: 30,
			wantErr:         false,
		},
		{
			name:            "NUMBER_OF_DAYS without value is invalid",
			definitionType:  "NUMBER_OF_DAYS",
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
