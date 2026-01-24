package sonargo

import (
	"net/http"
)

const (
	// MaxAlmSettingKeyLength is the maximum allowed length for DevOps Platform setting keys.
	MaxAlmSettingKeyLength = 200
	// MaxPatLength is the maximum allowed length for Personal Access Tokens.
	MaxPatLength = 2000
	// MaxUsernameLength is the maximum allowed length for usernames.
	MaxUsernameLength = 2000
	// MaxGitHubRepoKeyLength is the maximum allowed length for GitHub repository keys.
	MaxGitHubRepoKeyLength = 256
	// MaxPageSizeAlmIntegrations is the maximum allowed page size for ALM integrations APIs.
	MaxPageSizeAlmIntegrations = 100
	// MinNewCodeDefinitionDays is the minimum number of days for NUMBER_OF_DAYS new code definition.
	MinNewCodeDefinitionDays = 1
	// MaxNewCodeDefinitionDays is the maximum number of days for NUMBER_OF_DAYS new code definition.
	MaxNewCodeDefinitionDays = 90
)

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedNewCodeDefinitionTypes is the set of allowed new code definition types.
	allowedNewCodeDefinitionTypes = map[string]struct{}{
		"PREVIOUS_VERSION": {},
		"NUMBER_OF_DAYS":   {},
		"REFERENCE_BRANCH": {},
	}
)

// AlmIntegrationsService handles communication with the DevOps Platform Integration related methods
// of the SonarQube API.
// This service manages integrations with Azure DevOps, Bitbucket, GitHub, and GitLab.
type AlmIntegrationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AlmIntegrationsCheckPatResponse represents the response from checking a Personal Access Token.
// Note: The API does not return a response body, validation is indicated by HTTP status.
type AlmIntegrationsCheckPatResponse struct{}

// AlmIntegrationsGetGithubClientIdResponse represents the response from getting a GitHub client ID.
type AlmIntegrationsGetGithubClientIdResponse struct {
	// ClientID is the GitHub OAuth client ID for the integration.
	ClientID string `json:"clientId,omitempty"`
}

// AlmIntegrationsListAzureProjectsResponse represents the response from listing Azure projects.
type AlmIntegrationsListAzureProjectsResponse struct {
	// Projects is the list of Azure projects.
	Projects []AzureProject `json:"projects,omitempty"`
}

// AzureProject represents an Azure DevOps project.
type AzureProject struct {
	// Description is the project description.
	Description string `json:"description,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
}

// AlmIntegrationsListBitbucketServerProjectsResponse represents the response from listing Bitbucket Server projects.
type AlmIntegrationsListBitbucketServerProjectsResponse struct {
	// Projects is the list of Bitbucket Server projects.
	Projects []BitbucketServerProject `json:"projects,omitempty"`
}

// BitbucketServerProject represents a Bitbucket Server project.
type BitbucketServerProject struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
}

// AlmIntegrationsListGithubOrganizationsResponse represents the response from listing GitHub organizations.
type AlmIntegrationsListGithubOrganizationsResponse struct {
	// Organizations is the list of GitHub organizations.
	Organizations []GithubOrganization `json:"organizations,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// GithubOrganization represents a GitHub organization.
type GithubOrganization struct {
	// Key is the organization key (login).
	Key string `json:"key,omitempty"`
	// Name is the organization display name.
	Name string `json:"name,omitempty"`
}

// AlmIntegrationsListGithubRepositoriesResponse represents the response from listing GitHub repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsListGithubRepositoriesResponse struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Repositories is the list of GitHub repositories.
	Repositories []GithubRepository `json:"repositories,omitempty"`
}

// GithubRepository represents a GitHub repository.
//
//nolint:govet // Field alignment is less important than logical grouping
type GithubRepository struct {
	// ID is the repository ID.
	ID int64 `json:"id,omitempty"`
	// Key is the repository key (owner/repo).
	Key string `json:"key,omitempty"`
	// Name is the repository name.
	Name string `json:"name,omitempty"`
	// URL is the repository URL.
	URL string `json:"url,omitempty"`
}

// AlmIntegrationsSearchAzureReposResponse represents the response from searching Azure repositories.
type AlmIntegrationsSearchAzureReposResponse struct {
	// Repositories is the list of Azure repositories.
	Repositories []AzureRepository `json:"repositories,omitempty"`
}

// AzureRepository represents an Azure DevOps repository.
type AzureRepository struct {
	// Name is the repository name.
	Name string `json:"name,omitempty"`
	// ProjectName is the parent project name.
	ProjectName string `json:"projectName,omitempty"`
}

// AlmIntegrationsSearchBitbucketCloudReposResponse represents the response from searching Bitbucket Cloud repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsSearchBitbucketCloudReposResponse struct {
	// IsLastPage indicates if this is the last page of results.
	IsLastPage bool `json:"isLastPage,omitempty"`
	// Paging contains pagination information.
	Paging BitbucketCloudPaging `json:"paging,omitzero"`
	// Repositories is the list of Bitbucket Cloud repositories.
	Repositories []BitbucketCloudRepository `json:"repositories,omitempty"`
}

// BitbucketCloudPaging represents pagination info for Bitbucket Cloud APIs.
type BitbucketCloudPaging struct {
	// PageIndex is the current page index.
	PageIndex int64 `json:"pageIndex,omitempty"`
	// PageSize is the page size.
	PageSize int64 `json:"pageSize,omitempty"`
}

// BitbucketCloudRepository represents a Bitbucket Cloud repository.
type BitbucketCloudRepository struct {
	// Name is the repository name.
	Name string `json:"name,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Slug is the repository slug.
	Slug string `json:"slug,omitempty"`
	// SqProjectKey is the SonarQube project key if already imported.
	SqProjectKey string `json:"sqProjectKey,omitempty"`
	// UUID is the repository UUID.
	UUID string `json:"uuid,omitempty"`
	// Workspace is the workspace slug.
	Workspace string `json:"workspace,omitempty"`
}

// AlmIntegrationsSearchBitbucketServerReposResponse represents the response from searching Bitbucket Server repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsSearchBitbucketServerReposResponse struct {
	// IsLastPage indicates if this is the last page of results.
	IsLastPage bool `json:"isLastPage,omitempty"`
	// Repositories is the list of Bitbucket Server repositories.
	Repositories []BitbucketServerRepository `json:"repositories,omitempty"`
}

// BitbucketServerRepository represents a Bitbucket Server repository.
type BitbucketServerRepository struct {
	// Name is the repository name.
	Name string `json:"name,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Slug is the repository slug.
	Slug string `json:"slug,omitempty"`
	// SqProjectKey is the SonarQube project key if already imported.
	SqProjectKey string `json:"sqProjectKey,omitempty"`
	// UUID is the repository UUID.
	UUID string `json:"uuid,omitempty"`
	// Workspace is the workspace slug.
	Workspace string `json:"workspace,omitempty"`
}

// AlmIntegrationsSearchGitlabReposResponse represents the response from searching GitLab repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsSearchGitlabReposResponse struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Repositories is the list of GitLab repositories.
	Repositories []GitlabRepository `json:"repositories,omitempty"`
}

// GitlabRepository represents a GitLab repository (project).
//
//nolint:govet // Field alignment is less important than logical grouping
type GitlabRepository struct {
	// ID is the GitLab project ID.
	ID int64 `json:"id,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// PathName is the full path with namespace.
	PathName string `json:"pathName,omitempty"`
	// PathSlug is the path slug.
	PathSlug string `json:"pathSlug,omitempty"`
	// Slug is the project slug.
	Slug string `json:"slug,omitempty"`
	// URL is the project URL.
	URL string `json:"url,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AlmIntegrationsCheckPatOption contains options for checking a Personal Access Token.
type AlmIntegrationsCheckPatOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
}

// AlmIntegrationsGetGithubClientIdOption contains options for getting a GitHub client ID.
type AlmIntegrationsGetGithubClientIdOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
}

// AlmIntegrationsImportAzureProjectOption contains options for importing an Azure DevOps project.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
type AlmIntegrationsImportAzureProjectOption struct {
	// AlmSetting is the DevOps Platform configuration key.
	// This parameter is optional if you have only one Azure integration.
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// NewCodeDefinitionType is the project new code definition type.
	// Allowed values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH
	NewCodeDefinitionType string `url:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the project new code definition value.
	// Required when NewCodeDefinitionType is NUMBER_OF_DAYS (value between 1 and 90).
	// No value expected for PREVIOUS_VERSION and REFERENCE_BRANCH.
	NewCodeDefinitionValue string `url:"newCodeDefinitionValue,omitempty"`
	// ProjectName is the Azure project name (required).
	// Maximum length: 200 characters
	ProjectName string `url:"projectName,omitempty"`
	// RepositoryName is the Azure repository name (required).
	// Maximum length: 200 characters
	RepositoryName string `url:"repositoryName,omitempty"`
}

// AlmIntegrationsImportBitbucketCloudRepoOption contains options for importing a Bitbucket Cloud repository.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
type AlmIntegrationsImportBitbucketCloudRepoOption struct {
	// AlmSetting is the DevOps Platform configuration key.
	// This parameter is optional if you have only one Bitbucket Cloud integration.
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// NewCodeDefinitionType is the project new code definition type.
	// Allowed values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH
	NewCodeDefinitionType string `url:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the project new code definition value.
	// Required when NewCodeDefinitionType is NUMBER_OF_DAYS (value between 1 and 90).
	// No value expected for PREVIOUS_VERSION and REFERENCE_BRANCH.
	NewCodeDefinitionValue string `url:"newCodeDefinitionValue,omitempty"`
	// RepositorySlug is the Bitbucket Cloud repository slug (required).
	// Maximum length: 200 characters
	RepositorySlug string `url:"repositorySlug,omitempty"`
}

// AlmIntegrationsImportBitbucketServerProjectOption contains options for importing a Bitbucket Server project.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
type AlmIntegrationsImportBitbucketServerProjectOption struct {
	// AlmSetting is the DevOps Platform configuration key.
	// This parameter is optional if you have only one Bitbucket Server integration.
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// NewCodeDefinitionType is the project new code definition type.
	// Allowed values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH
	NewCodeDefinitionType string `url:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the project new code definition value.
	// Required when NewCodeDefinitionType is NUMBER_OF_DAYS (value between 1 and 90).
	// No value expected for PREVIOUS_VERSION and REFERENCE_BRANCH.
	NewCodeDefinitionValue string `url:"newCodeDefinitionValue,omitempty"`
	// ProjectKey is the Bitbucket Server project key (required).
	// Maximum length: 200 characters
	ProjectKey string `url:"projectKey,omitempty"`
	// RepositorySlug is the Bitbucket Server repository slug (required).
	// Maximum length: 200 characters
	RepositorySlug string `url:"repositorySlug,omitempty"`
}

// AlmIntegrationsImportGithubProjectOption contains options for importing a GitHub project.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
type AlmIntegrationsImportGithubProjectOption struct {
	// AlmSetting is the DevOps Platform configuration key.
	// This parameter is optional if you have only one GitHub integration.
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// NewCodeDefinitionType is the project new code definition type.
	// Allowed values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH
	NewCodeDefinitionType string `url:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the project new code definition value.
	// Required when NewCodeDefinitionType is NUMBER_OF_DAYS (value between 1 and 90).
	// No value expected for PREVIOUS_VERSION and REFERENCE_BRANCH.
	NewCodeDefinitionValue string `url:"newCodeDefinitionValue,omitempty"`
	// RepositoryKey is the GitHub repository key (organization/repoSlug) (required).
	// Maximum length: 256 characters
	RepositoryKey string `url:"repositoryKey,omitempty"`
}

// AlmIntegrationsImportGitlabProjectOption contains options for importing a GitLab project.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
type AlmIntegrationsImportGitlabProjectOption struct {
	// AlmSetting is the DevOps Platform configuration key.
	// This parameter is optional if you have only one GitLab integration.
	AlmSetting string `url:"almSetting,omitempty"`
	// GitlabProjectId is the GitLab project ID (required).
	GitlabProjectId string `url:"gitlabProjectId,omitempty"`
	// NewCodeDefinitionType is the project new code definition type.
	// Allowed values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH
	NewCodeDefinitionType string `url:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the project new code definition value.
	// Required when NewCodeDefinitionType is NUMBER_OF_DAYS (value between 1 and 90).
	// No value expected for PREVIOUS_VERSION and REFERENCE_BRANCH.
	NewCodeDefinitionValue string `url:"newCodeDefinitionValue,omitempty"`
}

// AlmIntegrationsListAzureProjectsOption contains options for listing Azure projects.
type AlmIntegrationsListAzureProjectsOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
}

// AlmIntegrationsListBitbucketServerProjectsOption contains options for listing Bitbucket Server projects.
type AlmIntegrationsListBitbucketServerProjectsOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// PageSize is the number of items to return (optional, default: 25, max: 100).
	PageSize int64 `url:"pageSize,omitempty"`
	// Start is the start number for the page (inclusive, optional).
	// If not passed, the first page is assumed.
	Start int64 `url:"start,omitempty"`
}

// AlmIntegrationsListGithubOrganizationsOption contains options for listing GitHub organizations.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsListGithubOrganizationsOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// P is the index of the page to display (optional, default: 1).
	P int64 `url:"p,omitempty"`
	// Ps is the size for the paging to apply (optional, default: 100).
	Ps int64 `url:"ps,omitempty"`
	// Token is the GitHub authorization code (optional).
	// Maximum length: 200 characters
	Token string `url:"token,omitempty"`
}

// AlmIntegrationsListGithubRepositoriesOption contains options for listing GitHub repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsListGithubRepositoriesOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// Organization is the GitHub organization (required).
	// Maximum length: 200 characters
	Organization string `url:"organization,omitempty"`
	// P is the index of the page to display (optional, default: 1).
	P int64 `url:"p,omitempty"`
	// Ps is the size for the paging to apply (optional, default: 100).
	Ps int64 `url:"ps,omitempty"`
	// Q is a filter to limit search to repositories that contain the supplied string (optional).
	Q string `url:"q,omitempty"`
}

// AlmIntegrationsSearchAzureReposOption contains options for searching Azure repositories.
type AlmIntegrationsSearchAzureReposOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// ProjectName is the project name filter (optional).
	// Maximum length: 200 characters
	ProjectName string `url:"projectName,omitempty"`
	// SearchQuery is the search query filter (optional).
	// Maximum length: 200 characters
	SearchQuery string `url:"searchQuery,omitempty"`
}

// AlmIntegrationsSearchBitbucketCloudReposOption contains options for searching Bitbucket Cloud repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsSearchBitbucketCloudReposOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// P is the 1-based page number (optional, default: 1).
	P int64 `url:"p,omitempty"`
	// Ps is the page size (optional, default: 20, max: 100).
	Ps int64 `url:"ps,omitempty"`
	// RepositoryName is the repository name filter (optional).
	// Maximum length: 200 characters
	RepositoryName string `url:"repositoryName,omitempty"`
}

// AlmIntegrationsSearchBitbucketServerReposOption contains options for searching Bitbucket Server repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsSearchBitbucketServerReposOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// PageSize is the number of items to return (optional, default: 25, max: 100).
	PageSize int64 `url:"pageSize,omitempty"`
	// ProjectName is the project name filter (optional).
	// Maximum length: 200 characters
	ProjectName string `url:"projectName,omitempty"`
	// RepositoryName is the repository name filter (optional).
	// Maximum length: 200 characters
	RepositoryName string `url:"repositoryName,omitempty"`
	// Start is the start number for the page (inclusive, optional).
	// If not passed, the first page is assumed.
	Start int64 `url:"start,omitempty"`
}

// AlmIntegrationsSearchGitlabReposOption contains options for searching GitLab repositories.
//
//nolint:govet // Field alignment is less important than logical grouping
type AlmIntegrationsSearchGitlabReposOption struct {
	// AlmSetting is the DevOps Platform setting key (required).
	// Maximum length: 200 characters
	AlmSetting string `url:"almSetting,omitempty"`
	// P is the 1-based page number (optional, default: 1).
	P int64 `url:"p,omitempty"`
	// ProjectName is the project name filter (optional).
	// Maximum length: 200 characters
	ProjectName string `url:"projectName,omitempty"`
	// Ps is the page size (optional, default: 20, max: 100).
	Ps int64 `url:"ps,omitempty"`
}

// AlmIntegrationsSetPatOption contains options for setting a Personal Access Token.
type AlmIntegrationsSetPatOption struct {
	// AlmSetting is the DevOps Platform configuration key.
	// This parameter is optional if you have only one single DevOps Platform integration.
	AlmSetting string `url:"almSetting,omitempty"`
	// Pat is the Personal Access Token (required).
	// Maximum length: 2000 characters
	Pat string `url:"pat,omitempty"`
	// Username is the username (optional, used for Bitbucket Cloud).
	// Maximum length: 2000 characters
	Username string `url:"username,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CheckPat checks the validity of a Personal Access Token for the given DevOps Platform setting.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) CheckPat(opt *AlmIntegrationsCheckPatOption) (v *AlmIntegrationsCheckPatResponse, resp *http.Response, err error) {
	err = s.ValidateCheckPatOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/check_pat", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsCheckPatResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// GetGithubClientId gets the client ID of a GitHub Integration.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) GetGithubClientId(opt *AlmIntegrationsGetGithubClientIdOption) (v *AlmIntegrationsGetGithubClientIdResponse, resp *http.Response, err error) {
	err = s.ValidateGetGithubClientIdOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/get_github_client_id", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsGetGithubClientIdResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ImportAzureProject creates a SonarQube project with the information from the provided Azure DevOps project.
// Autoconfigures pull request decoration mechanism.
// Requires the 'Create Projects' permission.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
func (s *AlmIntegrationsService) ImportAzureProject(opt *AlmIntegrationsImportAzureProjectOption) (resp *http.Response, err error) {
	err = s.ValidateImportAzureProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_integrations/import_azure_project", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// ImportBitbucketCloudRepo creates a SonarQube project with the information from the provided Bitbucket Cloud repository.
// Autoconfigures pull request decoration mechanism.
// Requires the 'Create Projects' permission.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
func (s *AlmIntegrationsService) ImportBitbucketCloudRepo(opt *AlmIntegrationsImportBitbucketCloudRepoOption) (resp *http.Response, err error) {
	err = s.ValidateImportBitbucketCloudRepoOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_integrations/import_bitbucketcloud_repo", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// ImportBitbucketServerProject creates a SonarQube project with the information from the provided Bitbucket Server project.
// Autoconfigures pull request decoration mechanism.
// Requires the 'Create Projects' permission.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
func (s *AlmIntegrationsService) ImportBitbucketServerProject(opt *AlmIntegrationsImportBitbucketServerProjectOption) (resp *http.Response, err error) {
	err = s.ValidateImportBitbucketServerProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_integrations/import_bitbucketserver_project", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// ImportGithubProject creates a SonarQube project with the information from the provided GitHub repository.
// Autoconfigures pull request decoration mechanism.
// If Automatic Provisioning is enabled for GitHub, it will also synchronize permissions from the repository.
// Requires the 'Create Projects' permission.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
func (s *AlmIntegrationsService) ImportGithubProject(opt *AlmIntegrationsImportGithubProjectOption) (resp *http.Response, err error) {
	err = s.ValidateImportGithubProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_integrations/import_github_project", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// ImportGitlabProject imports a GitLab project to SonarQube, creating a new project and configuring MR decoration.
// Requires the 'Create Projects' permission.
//
// Deprecated: Since 10.5 - use /api/v2/dop-translation/bound-projects instead.
func (s *AlmIntegrationsService) ImportGitlabProject(opt *AlmIntegrationsImportGitlabProjectOption) (resp *http.Response, err error) {
	err = s.ValidateImportGitlabProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_integrations/import_gitlab_project", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// ListAzureProjects lists Azure projects.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) ListAzureProjects(opt *AlmIntegrationsListAzureProjectsOption) (v *AlmIntegrationsListAzureProjectsResponse, resp *http.Response, err error) {
	err = s.ValidateListAzureProjectsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/list_azure_projects", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsListAzureProjectsResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ListBitbucketServerProjects lists the Bitbucket Server projects.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) ListBitbucketServerProjects(opt *AlmIntegrationsListBitbucketServerProjectsOption) (v *AlmIntegrationsListBitbucketServerProjectsResponse, resp *http.Response, err error) {
	err = s.ValidateListBitbucketServerProjectsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/list_bitbucketserver_projects", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsListBitbucketServerProjectsResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ListGithubOrganizations lists GitHub organizations.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) ListGithubOrganizations(opt *AlmIntegrationsListGithubOrganizationsOption) (v *AlmIntegrationsListGithubOrganizationsResponse, resp *http.Response, err error) {
	err = s.ValidateListGithubOrganizationsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/list_github_organizations", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsListGithubOrganizationsResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ListGithubRepositories lists the GitHub repositories for an organization.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) ListGithubRepositories(opt *AlmIntegrationsListGithubRepositoriesOption) (v *AlmIntegrationsListGithubRepositoriesResponse, resp *http.Response, err error) {
	err = s.ValidateListGithubRepositoriesOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/list_github_repositories", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsListGithubRepositoriesResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchAzureRepos searches the Azure repositories.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) SearchAzureRepos(opt *AlmIntegrationsSearchAzureReposOption) (v *AlmIntegrationsSearchAzureReposResponse, resp *http.Response, err error) {
	err = s.ValidateSearchAzureReposOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/search_azure_repos", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsSearchAzureReposResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchBitbucketCloudRepos searches the Bitbucket Cloud repositories.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) SearchBitbucketCloudRepos(opt *AlmIntegrationsSearchBitbucketCloudReposOption) (v *AlmIntegrationsSearchBitbucketCloudReposResponse, resp *http.Response, err error) {
	err = s.ValidateSearchBitbucketCloudReposOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/search_bitbucketcloud_repos", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsSearchBitbucketCloudReposResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchBitbucketServerRepos searches the Bitbucket Server repositories with REPO_ADMIN access.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) SearchBitbucketServerRepos(opt *AlmIntegrationsSearchBitbucketServerReposOption) (v *AlmIntegrationsSearchBitbucketServerReposResponse, resp *http.Response, err error) {
	err = s.ValidateSearchBitbucketServerReposOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/search_bitbucketserver_repos", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsSearchBitbucketServerReposResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchGitlabRepos searches the GitLab projects.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) SearchGitlabRepos(opt *AlmIntegrationsSearchGitlabReposOption) (v *AlmIntegrationsSearchGitlabReposResponse, resp *http.Response, err error) {
	err = s.ValidateSearchGitlabReposOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_integrations/search_gitlab_repos", opt)
	if err != nil {
		return
	}

	v = new(AlmIntegrationsSearchGitlabReposResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SetPat sets a Personal Access Token for the given DevOps Platform setting.
// Requires the 'Create Projects' permission.
func (s *AlmIntegrationsService) SetPat(opt *AlmIntegrationsSetPatOption) (resp *http.Response, err error) {
	err = s.ValidateSetPatOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_integrations/set_pat", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCheckPatOpt validates the options for checking a Personal Access Token.
func (s *AlmIntegrationsService) ValidateCheckPatOpt(opt *AlmIntegrationsCheckPatOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsCheckPatOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetGithubClientIdOpt validates the options for getting a GitHub client ID.
func (s *AlmIntegrationsService) ValidateGetGithubClientIdOpt(opt *AlmIntegrationsGetGithubClientIdOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsGetGithubClientIdOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	return nil
}

// ValidateImportAzureProjectOpt validates the options for importing an Azure DevOps project.
func (s *AlmIntegrationsService) ValidateImportAzureProjectOpt(opt *AlmIntegrationsImportAzureProjectOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsImportAzureProjectOption", "cannot be nil", ErrMissingRequired)
	}

	// AlmSetting is optional if only one Azure integration exists
	if opt.AlmSetting != "" {
		err := ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
		if err != nil {
			return err
		}
	}

	// Validate new code definition if provided
	err := validateNewCodeDefinition(opt.NewCodeDefinitionType, opt.NewCodeDefinitionValue)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectName, "ProjectName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ProjectName, MaxAlmSettingKeyLength, "ProjectName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RepositoryName, "RepositoryName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.RepositoryName, MaxAlmSettingKeyLength, "RepositoryName")
	if err != nil {
		return err
	}

	return nil
}

// ValidateImportBitbucketCloudRepoOpt validates the options for importing a Bitbucket Cloud repository.
func (s *AlmIntegrationsService) ValidateImportBitbucketCloudRepoOpt(opt *AlmIntegrationsImportBitbucketCloudRepoOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsImportBitbucketCloudRepoOption", "cannot be nil", ErrMissingRequired)
	}

	// AlmSetting is optional if only one Bitbucket Cloud integration exists
	if opt.AlmSetting != "" {
		err := ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
		if err != nil {
			return err
		}
	}

	// Validate new code definition if provided
	err := validateNewCodeDefinition(opt.NewCodeDefinitionType, opt.NewCodeDefinitionValue)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RepositorySlug, "RepositorySlug")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.RepositorySlug, MaxAlmSettingKeyLength, "RepositorySlug")
	if err != nil {
		return err
	}

	return nil
}

// ValidateImportBitbucketServerProjectOpt validates the options for importing a Bitbucket Server project.
func (s *AlmIntegrationsService) ValidateImportBitbucketServerProjectOpt(opt *AlmIntegrationsImportBitbucketServerProjectOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsImportBitbucketServerProjectOption", "cannot be nil", ErrMissingRequired)
	}

	// AlmSetting is optional if only one Bitbucket Server integration exists
	if opt.AlmSetting != "" {
		err := ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
		if err != nil {
			return err
		}
	}

	// Validate new code definition if provided
	err := validateNewCodeDefinition(opt.NewCodeDefinitionType, opt.NewCodeDefinitionValue)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ProjectKey, MaxAlmSettingKeyLength, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RepositorySlug, "RepositorySlug")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.RepositorySlug, MaxAlmSettingKeyLength, "RepositorySlug")
	if err != nil {
		return err
	}

	return nil
}

// ValidateImportGithubProjectOpt validates the options for importing a GitHub project.
func (s *AlmIntegrationsService) ValidateImportGithubProjectOpt(opt *AlmIntegrationsImportGithubProjectOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsImportGithubProjectOption", "cannot be nil", ErrMissingRequired)
	}

	// AlmSetting is optional if only one GitHub integration exists
	if opt.AlmSetting != "" {
		err := ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
		if err != nil {
			return err
		}
	}

	// Validate new code definition if provided
	err := validateNewCodeDefinition(opt.NewCodeDefinitionType, opt.NewCodeDefinitionValue)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RepositoryKey, "RepositoryKey")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.RepositoryKey, MaxGitHubRepoKeyLength, "RepositoryKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateImportGitlabProjectOpt validates the options for importing a GitLab project.
func (s *AlmIntegrationsService) ValidateImportGitlabProjectOpt(opt *AlmIntegrationsImportGitlabProjectOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsImportGitlabProjectOption", "cannot be nil", ErrMissingRequired)
	}

	// AlmSetting is optional if only one GitLab integration exists

	// Validate new code definition if provided
	err := validateNewCodeDefinition(opt.NewCodeDefinitionType, opt.NewCodeDefinitionValue)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.GitlabProjectId, "GitlabProjectId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListAzureProjectsOpt validates the options for listing Azure projects.
func (s *AlmIntegrationsService) ValidateListAzureProjectsOpt(opt *AlmIntegrationsListAzureProjectsOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsListAzureProjectsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListBitbucketServerProjectsOpt validates the options for listing Bitbucket Server projects.
func (s *AlmIntegrationsService) ValidateListBitbucketServerProjectsOpt(opt *AlmIntegrationsListBitbucketServerProjectsOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsListBitbucketServerProjectsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	if opt.PageSize != 0 {
		err = ValidateRange(opt.PageSize, MinPageSize, MaxPageSizeAlmIntegrations, "PageSize")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateListGithubOrganizationsOpt validates the options for listing GitHub organizations.
func (s *AlmIntegrationsService) ValidateListGithubOrganizationsOpt(opt *AlmIntegrationsListGithubOrganizationsOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsListGithubOrganizationsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	if opt.Token != "" {
		err = ValidateMaxLength(opt.Token, MaxAlmSettingKeyLength, "Token")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateListGithubRepositoriesOpt validates the options for listing GitHub repositories.
func (s *AlmIntegrationsService) ValidateListGithubRepositoriesOpt(opt *AlmIntegrationsListGithubRepositoriesOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsListGithubRepositoriesOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Organization, "Organization")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Organization, MaxAlmSettingKeyLength, "Organization")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchAzureReposOpt validates the options for searching Azure repositories.
func (s *AlmIntegrationsService) ValidateSearchAzureReposOpt(opt *AlmIntegrationsSearchAzureReposOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsSearchAzureReposOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	if opt.ProjectName != "" {
		err = ValidateMaxLength(opt.ProjectName, MaxAlmSettingKeyLength, "ProjectName")
		if err != nil {
			return err
		}
	}

	if opt.SearchQuery != "" {
		err = ValidateMaxLength(opt.SearchQuery, MaxAlmSettingKeyLength, "SearchQuery")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchBitbucketCloudReposOpt validates the options for searching Bitbucket Cloud repositories.
func (s *AlmIntegrationsService) ValidateSearchBitbucketCloudReposOpt(opt *AlmIntegrationsSearchBitbucketCloudReposOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsSearchBitbucketCloudReposOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	if opt.Ps != 0 {
		err = ValidateRange(opt.Ps, MinPageSize, MaxPageSizeAlmIntegrations, "Ps")
		if err != nil {
			return err
		}
	}

	if opt.RepositoryName != "" {
		err = ValidateMaxLength(opt.RepositoryName, MaxAlmSettingKeyLength, "RepositoryName")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchBitbucketServerReposOpt validates the options for searching Bitbucket Server repositories.
func (s *AlmIntegrationsService) ValidateSearchBitbucketServerReposOpt(opt *AlmIntegrationsSearchBitbucketServerReposOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsSearchBitbucketServerReposOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	if opt.PageSize != 0 {
		err = ValidateRange(opt.PageSize, MinPageSize, MaxPageSizeAlmIntegrations, "PageSize")
		if err != nil {
			return err
		}
	}

	if opt.ProjectName != "" {
		err = ValidateMaxLength(opt.ProjectName, MaxAlmSettingKeyLength, "ProjectName")
		if err != nil {
			return err
		}
	}

	if opt.RepositoryName != "" {
		err = ValidateMaxLength(opt.RepositoryName, MaxAlmSettingKeyLength, "RepositoryName")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchGitlabReposOpt validates the options for searching GitLab repositories.
func (s *AlmIntegrationsService) ValidateSearchGitlabReposOpt(opt *AlmIntegrationsSearchGitlabReposOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsSearchGitlabReposOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmSettingKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	if opt.Ps != 0 {
		err = ValidateRange(opt.Ps, MinPageSize, MaxPageSizeAlmIntegrations, "Ps")
		if err != nil {
			return err
		}
	}

	if opt.ProjectName != "" {
		err = ValidateMaxLength(opt.ProjectName, MaxAlmSettingKeyLength, "ProjectName")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSetPatOpt validates the options for setting a Personal Access Token.
func (s *AlmIntegrationsService) ValidateSetPatOpt(opt *AlmIntegrationsSetPatOption) error {
	if opt == nil {
		return NewValidationError("AlmIntegrationsSetPatOption", "cannot be nil", ErrMissingRequired)
	}

	// AlmSetting is optional if only one DevOps Platform integration exists

	err := ValidateRequired(opt.Pat, "Pat")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Pat, MaxPatLength, "Pat")
	if err != nil {
		return err
	}

	if opt.Username != "" {
		err = ValidateMaxLength(opt.Username, MaxUsernameLength, "Username")
		if err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------------

// validateNewCodeDefinition validates the new code definition type and value.
func validateNewCodeDefinition(definitionType, definitionValue string) error {
	if definitionType == "" {
		return nil
	}

	err := IsValueAuthorized(definitionType, allowedNewCodeDefinitionTypes, "NewCodeDefinitionType")
	if err != nil {
		return err
	}

	// NUMBER_OF_DAYS requires a value between 1 and 90
	if definitionType == "NUMBER_OF_DAYS" {
		if definitionValue == "" {
			return NewValidationError(
				"NewCodeDefinitionValue",
				"is required when NewCodeDefinitionType is NUMBER_OF_DAYS",
				ErrMissingRequired,
			)
		}
		// The value is passed as string, so we just validate it's not empty
		// The API will validate the numeric range
	}

	// PREVIOUS_VERSION and REFERENCE_BRANCH should not have a value
	if (definitionType == "PREVIOUS_VERSION" || definitionType == "REFERENCE_BRANCH") && definitionValue != "" {
		return NewValidationError(
			"NewCodeDefinitionValue",
			"should not be provided when NewCodeDefinitionType is "+definitionType,
			ErrInvalidValue,
		)
	}

	return nil
}
