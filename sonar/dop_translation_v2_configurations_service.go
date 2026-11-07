package sonar

import (
	"context"
	"net/http"
)

const (
	// DopProvisioningTypeJIT is the just-in-time provisioning type.
	DopProvisioningTypeJIT = "JIT"
	// DopProvisioningTypeAuto is the automatic provisioning type.
	DopProvisioningTypeAuto = "AUTO_PROVISIONING"
)

//nolint:gochecknoglobals // allowed-value sets
var allowedDopProvisioningTypes = map[string]struct{}{DopProvisioningTypeJIT: {}, DopProvisioningTypeAuto: {}}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// DopTranslationGithubConfiguration represents a GitHub authentication/provisioning configuration.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationGithubConfiguration struct {
	// Id is the configuration identifier.
	Id string `json:"id,omitempty"`
	// Enabled indicates whether the configuration is enabled.
	Enabled bool `json:"enabled,omitempty"`
	// ApplicationId is the GitHub App id.
	ApplicationId string `json:"applicationId,omitempty"`
	// SynchronizeGroups indicates whether groups are synchronized.
	SynchronizeGroups bool `json:"synchronizeGroups,omitempty"`
	// ApiUrl is the GitHub API URL.
	ApiUrl string `json:"apiUrl,omitempty"`
	// WebUrl is the GitHub web URL.
	WebUrl string `json:"webUrl,omitempty"`
	// AllowedOrganizations is the list of allowed GitHub organizations.
	AllowedOrganizations []string `json:"allowedOrganizations,omitempty"`
	// ProvisioningType is the provisioning type. Allowed values: JIT, AUTO_PROVISIONING.
	ProvisioningType string `json:"provisioningType,omitempty"`
	// AllowUsersToSignUp indicates whether users may sign up.
	AllowUsersToSignUp bool `json:"allowUsersToSignUp,omitempty"`
	// ProjectVisibility indicates whether the project visibility is synchronized.
	ProjectVisibility bool `json:"projectVisibility,omitempty"`
	// UserConsentRequiredAfterUpgrade indicates whether user consent is required after an upgrade.
	UserConsentRequiredAfterUpgrade bool `json:"userConsentRequiredAfterUpgrade,omitempty"`
	// ErrorMessage is the last error message, if any.
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// DopTranslationGithubConfigurationsSearch represents the response from searching GitHub configurations.
type DopTranslationGithubConfigurationsSearch struct {
	// GithubConfigurations is the list of GitHub configurations.
	GithubConfigurations []DopTranslationGithubConfiguration `json:"githubConfigurations,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// DopTranslationGitlabConfiguration represents a GitLab authentication/provisioning configuration.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationGitlabConfiguration struct {
	// Id is the configuration identifier.
	Id string `json:"id,omitempty"`
	// Enabled indicates whether the configuration is enabled.
	Enabled bool `json:"enabled,omitempty"`
	// ApplicationId is the GitLab application id.
	ApplicationId string `json:"applicationId,omitempty"`
	// Url is the GitLab URL.
	Url string `json:"url,omitempty"`
	// SynchronizeGroups indicates whether groups are synchronized.
	SynchronizeGroups bool `json:"synchronizeGroups,omitempty"`
	// AllowedGroups is the list of allowed GitLab groups.
	AllowedGroups []string `json:"allowedGroups,omitempty"`
	// AllowAllGroups indicates whether all groups are allowed.
	AllowAllGroups bool `json:"allowAllGroups,omitempty"`
	// AllowUsersToSignUp indicates whether users may sign up.
	AllowUsersToSignUp bool `json:"allowUsersToSignUp,omitempty"`
	// ProvisioningType is the provisioning type. Allowed values: JIT, AUTO_PROVISIONING.
	ProvisioningType string `json:"provisioningType,omitempty"`
	// IsProvisioningTokenSet indicates whether a provisioning token is set.
	IsProvisioningTokenSet bool `json:"isProvisioningTokenSet,omitempty"`
	// ErrorMessage is the last error message, if any.
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// DopTranslationGitlabConfigurationsSearch represents the response from searching GitLab configurations.
type DopTranslationGitlabConfigurationsSearch struct {
	// GitlabConfigurations is the list of GitLab configurations.
	GitlabConfigurations []DopTranslationGitlabConfiguration `json:"gitlabConfigurations,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// DopTranslationPermissions represents the SonarQube permissions granted to a DevOps Platform role.
type DopTranslationPermissions struct {
	// User is the 'user' permission.
	User bool `json:"user"`
	// CodeViewer is the 'codeviewer' permission.
	CodeViewer bool `json:"codeViewer"`
	// IssueAdmin is the 'issueadmin' permission.
	IssueAdmin bool `json:"issueAdmin"`
	// SecurityHotspotAdmin is the 'securityhotspotadmin' permission.
	SecurityHotspotAdmin bool `json:"securityHotspotAdmin"`
	// Admin is the 'admin' permission.
	Admin bool `json:"admin"`
	// Scan is the 'scan' permission.
	Scan bool `json:"scan"`
}

// DopTranslationPermissionMapping represents the mapping between a DevOps Platform role and SonarQube permissions.
type DopTranslationPermissionMapping struct {
	// Id is the mapping identifier.
	Id string `json:"id,omitempty"`
	// Role is the DevOps Platform role.
	Role string `json:"role,omitempty"`
	// Permissions are the SonarQube permissions granted to the role.
	Permissions DopTranslationPermissions `json:"permissions,omitzero"`
	// BaseRole indicates whether the role is a built-in role of the DevOps Platform.
	BaseRole bool `json:"baseRole,omitempty"`
}

// DopTranslationPermissionMappingsSearch represents the response from fetching permission mappings.
type DopTranslationPermissionMappingsSearch struct {
	// PermissionMappings is the list of permission mappings.
	PermissionMappings []DopTranslationPermissionMapping `json:"permissionMappings,omitempty"`
}

// DopTranslationProjectBinding represents the binding of a project to a DevOps Platform repository.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationProjectBinding struct {
	// Id is the binding identifier.
	Id string `json:"id,omitempty"`
	// DevOpsPlatformSettingId is the identifier of the DevOps Platform setting.
	DevOpsPlatformSettingId string `json:"devOpsPlatformSettingId,omitempty"`
	// ProjectId is the project UUID.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Repository is the repository identifier.
	Repository *string `json:"repository,omitempty"`
	// Slug is the repository slug.
	Slug *string `json:"slug,omitempty"`
}

// DopTranslationProjectBindingsSearch represents the response from searching project bindings.
type DopTranslationProjectBindingsSearch struct {
	// ProjectBindings is the list of project bindings.
	ProjectBindings []DopTranslationProjectBinding `json:"projectBindings,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// DopTranslationGitlabSynchronizationRun represents a started GitLab synchronization run.
type DopTranslationGitlabSynchronizationRun struct {
	// Id is the synchronization run identifier.
	Id string `json:"id,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// DopTranslationGithubConfigurationCreateOptions contains parameters for the CreateGithubConfiguration method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationGithubConfigurationCreateOptions struct {
	// Enabled indicates whether the configuration is enabled. This field is required.
	Enabled bool `json:"enabled"`
	// ClientId is the GitHub App client id. This field is required.
	ClientId string `json:"clientId"`
	// ClientSecret is the GitHub App client secret. This field is required.
	ClientSecret string `json:"clientSecret"`
	// ApplicationId is the GitHub App id. This field is required.
	ApplicationId string `json:"applicationId"`
	// PrivateKey is the GitHub App private key. This field is required.
	PrivateKey string `json:"privateKey"`
	// SynchronizeGroups indicates whether groups are synchronized. This field is required.
	SynchronizeGroups bool `json:"synchronizeGroups"`
	// ApiUrl is the GitHub API URL. This field is required.
	ApiUrl string `json:"apiUrl"`
	// WebUrl is the GitHub web URL. This field is required.
	WebUrl string `json:"webUrl"`
	// AllowedOrganizations is the list of allowed GitHub organizations. This field is required.
	AllowedOrganizations []string `json:"allowedOrganizations"`
	// ProvisioningType is the provisioning type. Allowed values: JIT, AUTO_PROVISIONING. This field is required.
	ProvisioningType string `json:"provisioningType"`
	// AllowUsersToSignUp indicates whether users may sign up.
	AllowUsersToSignUp *bool `json:"allowUsersToSignUp,omitempty"`
	// ProjectVisibility indicates whether the project visibility is synchronized.
	ProjectVisibility *bool `json:"projectVisibility,omitempty"`
	// UserConsentRequiredAfterUpgrade indicates whether user consent is required after an upgrade.
	UserConsentRequiredAfterUpgrade *bool `json:"userConsentRequiredAfterUpgrade,omitempty"`
}

// DopTranslationGithubConfigurationUpdateOptions contains parameters for the UpdateGithubConfiguration method.
// All fields are optional (PATCH merge semantics).
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationGithubConfigurationUpdateOptions struct {
	// Enabled indicates whether the configuration is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// ClientId is the GitHub App client id.
	ClientId string `json:"clientId,omitempty"`
	// ClientSecret is the GitHub App client secret.
	ClientSecret string `json:"clientSecret,omitempty"`
	// ApplicationId is the GitHub App id.
	ApplicationId string `json:"applicationId,omitempty"`
	// PrivateKey is the GitHub App private key. Required when changing ApiUrl or WebUrl.
	PrivateKey string `json:"privateKey,omitempty"`
	// SynchronizeGroups indicates whether groups are synchronized.
	SynchronizeGroups *bool `json:"synchronizeGroups,omitempty"`
	// ApiUrl is the GitHub API URL.
	ApiUrl string `json:"apiUrl,omitempty"`
	// WebUrl is the GitHub web URL.
	WebUrl string `json:"webUrl,omitempty"`
	// AllowedOrganizations is the list of allowed GitHub organizations.
	AllowedOrganizations *UpdateFieldListStringV2 `json:"allowedOrganizations,omitempty"`
	// ProvisioningType is the provisioning type. Allowed values: JIT, AUTO_PROVISIONING.
	ProvisioningType string `json:"provisioningType,omitempty"`
	// AllowUsersToSignUp indicates whether users may sign up.
	AllowUsersToSignUp *bool `json:"allowUsersToSignUp,omitempty"`
	// ProjectVisibility indicates whether the project visibility is synchronized.
	ProjectVisibility *bool `json:"projectVisibility,omitempty"`
	// UserConsentRequiredAfterUpgrade indicates whether user consent is required after an upgrade.
	UserConsentRequiredAfterUpgrade *bool `json:"userConsentRequiredAfterUpgrade,omitempty"`
}

// DopTranslationGitlabConfigurationCreateOptions contains parameters for the CreateGitlabConfiguration method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationGitlabConfigurationCreateOptions struct {
	// Enabled indicates whether the configuration is enabled. This field is required.
	Enabled bool `json:"enabled"`
	// ApplicationId is the GitLab application id. This field is required.
	ApplicationId string `json:"applicationId"`
	// Url is the GitLab URL. This field is required.
	Url string `json:"url"`
	// Secret is the GitLab application secret. This field is required.
	Secret string `json:"secret"`
	// SynchronizeGroups indicates whether groups are synchronized. This field is required.
	SynchronizeGroups bool `json:"synchronizeGroups"`
	// AllowedGroups is the list of allowed GitLab groups. This field is required.
	AllowedGroups []string `json:"allowedGroups"`
	// AllowAllGroups indicates whether all groups are allowed. Requires auto-provisioning.
	AllowAllGroups *bool `json:"allowAllGroups,omitempty"`
	// ProvisioningType is the provisioning type. Allowed values: JIT, AUTO_PROVISIONING. This field is required.
	ProvisioningType string `json:"provisioningType"`
	// ProvisioningToken is the provisioning token, for auto-provisioning.
	ProvisioningToken string `json:"provisioningToken,omitempty"`
	// AllowUsersToSignUp indicates whether users may sign up.
	AllowUsersToSignUp *bool `json:"allowUsersToSignUp,omitempty"`
}

// DopTranslationGitlabConfigurationUpdateOptions contains parameters for the UpdateGitlabConfiguration method.
// All fields are optional (PATCH merge semantics).
type DopTranslationGitlabConfigurationUpdateOptions struct {
	// Enabled indicates whether the configuration is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// ApplicationId is the GitLab application id.
	ApplicationId string `json:"applicationId,omitempty"`
	// Url is the GitLab URL.
	Url string `json:"url,omitempty"`
	// Secret is the GitLab application secret.
	Secret string `json:"secret,omitempty"`
	// SynchronizeGroups indicates whether groups are synchronized.
	SynchronizeGroups *bool `json:"synchronizeGroups,omitempty"`
	// AllowedGroups is the list of allowed GitLab groups.
	AllowedGroups *UpdateFieldListStringV2 `json:"allowedGroups,omitempty"`
	// AllowAllGroups indicates whether all groups are allowed. Requires auto-provisioning.
	AllowAllGroups *bool `json:"allowAllGroups,omitempty"`
	// ProvisioningType is the provisioning type. Allowed values: JIT, AUTO_PROVISIONING.
	ProvisioningType string `json:"provisioningType,omitempty"`
	// AllowUsersToSignUp indicates whether users may sign up.
	AllowUsersToSignUp *bool `json:"allowUsersToSignUp,omitempty"`
	// ProvisioningToken is the provisioning token, for auto-provisioning.
	ProvisioningToken string `json:"provisioningToken,omitempty"`
}

// DopTranslationPermissionMappingCreateOptions contains parameters for creating a permission mapping for a custom role.
type DopTranslationPermissionMappingCreateOptions struct {
	// Role is the custom role name. This field is required.
	Role string `json:"role"`
	// Permissions are the SonarQube permissions granted to the role. This field is required.
	Permissions DopTranslationPermissions `json:"permissions"`
}

// DopTranslationPermissionsUpdate holds the permissions to change in a mapping. Nil fields are left unchanged.
type DopTranslationPermissionsUpdate struct {
	// User is the 'user' permission.
	User *bool `json:"user,omitempty"`
	// CodeViewer is the 'codeviewer' permission.
	CodeViewer *bool `json:"codeViewer,omitempty"`
	// IssueAdmin is the 'issueadmin' permission.
	IssueAdmin *bool `json:"issueAdmin,omitempty"`
	// SecurityHotspotAdmin is the 'securityhotspotadmin' permission.
	SecurityHotspotAdmin *bool `json:"securityHotspotAdmin,omitempty"`
	// Admin is the 'admin' permission.
	Admin *bool `json:"admin,omitempty"`
	// Scan is the 'scan' permission.
	Scan *bool `json:"scan,omitempty"`
}

// DopTranslationPermissionMappingUpdateOptions contains parameters for updating a permission mapping.
type DopTranslationPermissionMappingUpdateOptions struct {
	// Permissions are the permissions to change.
	Permissions *DopTranslationPermissionsUpdate `json:"permissions,omitempty"`
}

// DopTranslationProjectBindingsSearchOptions contains parameters for the SearchProjectBindings method.
type DopTranslationProjectBindingsSearchOptions struct {
	// Repository filters on the repository name (exact, case insensitive match).
	Repository string `json:"repository,omitempty"`
	// DopSettingId filters on the DevOps Platform setting id.
	DopSettingId string `json:"dopSettingId,omitempty"`
	// RepositoryUrl filters on the repository URL (traditional URL or git format).
	RepositoryUrl string `json:"repositoryUrl,omitempty"`
	// PageIndex is the 1-based page index. Default is 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the number of results per page. Default is 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

func validateDopProvisioningType(value string, required bool) error {
	if value == "" && !required {
		return nil
	}

	err := ValidateRequired(value, "ProvisioningType")
	if err != nil {
		return err
	}

	return IsValueAuthorized(value, allowedDopProvisioningTypes, "ProvisioningType")
}

// ValidateCreateGithubConfigurationOpt validates the options for the CreateGithubConfiguration method.
func (s *DopTranslationService) ValidateCreateGithubConfigurationOpt(opt *DopTranslationGithubConfigurationCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	for _, field := range []struct{ value, name string }{
		{opt.ClientId, "ClientId"},
		{opt.ClientSecret, "ClientSecret"},
		{opt.ApplicationId, "ApplicationId"},
		{opt.PrivateKey, "PrivateKey"},
		{opt.ApiUrl, "ApiUrl"},
		{opt.WebUrl, "WebUrl"},
	} {
		err := ValidateRequired(field.value, field.name)
		if err != nil {
			return err
		}
	}

	return validateDopProvisioningType(opt.ProvisioningType, true)
}

// ValidateUpdateGithubConfigurationOpt validates the options for the UpdateGithubConfiguration method.
func (s *DopTranslationService) ValidateUpdateGithubConfigurationOpt(opt *DopTranslationGithubConfigurationUpdateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return validateDopProvisioningType(opt.ProvisioningType, false)
}

// ValidateCreateGitlabConfigurationOpt validates the options for the CreateGitlabConfiguration method.
func (s *DopTranslationService) ValidateCreateGitlabConfigurationOpt(opt *DopTranslationGitlabConfigurationCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	for _, field := range []struct{ value, name string }{
		{opt.ApplicationId, "ApplicationId"}, {opt.Url, "Url"}, {opt.Secret, "Secret"},
	} {
		err := ValidateRequired(field.value, field.name)
		if err != nil {
			return err
		}
	}

	return validateDopProvisioningType(opt.ProvisioningType, true)
}

// ValidateUpdateGitlabConfigurationOpt validates the options for the UpdateGitlabConfiguration method.
func (s *DopTranslationService) ValidateUpdateGitlabConfigurationOpt(opt *DopTranslationGitlabConfigurationUpdateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return validateDopProvisioningType(opt.ProvisioningType, false)
}

// ValidatePermissionMappingCreateOpt validates the options for creating a permission mapping.
func (s *DopTranslationService) ValidatePermissionMappingCreateOpt(opt *DopTranslationPermissionMappingCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Role, "Role")
}

// ValidateProjectBindingsSearchOpt validates the options for the SearchProjectBindings method.
func (s *DopTranslationService) ValidateProjectBindingsSearchOpt(opt *DopTranslationProjectBindingsSearchOptions) error {
	if opt == nil {
		return nil
	}

	return (&PaginationParamsV2{PageIndex: opt.PageIndex, PageSize: opt.PageSize}).Validate()
}

// -----------------------------------------------------------------------------
// GitHub configurations
// -----------------------------------------------------------------------------

// SearchGithubConfigurations lists the GitHub configurations.
//
// API endpoint: GET /api/v2/dop-translation/github-configurations.
func (s *DopTranslationService) SearchGithubConfigurations(ctx context.Context) (*DopTranslationGithubConfigurationsSearch, *http.Response, error) {
	return callV2[DopTranslationGithubConfigurationsSearch](ctx, s.client, http.MethodGet, "dop-translation/github-configurations", nil, nil)
}

// CreateGithubConfiguration creates a GitHub configuration.
//
// API endpoint: POST /api/v2/dop-translation/github-configurations.
func (s *DopTranslationService) CreateGithubConfiguration(ctx context.Context, opt *DopTranslationGithubConfigurationCreateOptions) (*DopTranslationGithubConfiguration, *http.Response, error) {
	err := s.ValidateCreateGithubConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationGithubConfiguration](ctx, s.client, http.MethodPost, "dop-translation/github-configurations", nil, opt)
}

// GetGithubConfiguration fetches a GitHub configuration.
//
// API endpoint: GET /api/v2/dop-translation/github-configurations/{id}.
func (s *DopTranslationService) GetGithubConfiguration(ctx context.Context, id string) (*DopTranslationGithubConfiguration, *http.Response, error) {
	path, err := v2PathWithID("dop-translation/github-configurations", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationGithubConfiguration](ctx, s.client, http.MethodGet, path, nil, nil)
}

// UpdateGithubConfiguration updates a GitHub configuration.
//
// API endpoint: PATCH /api/v2/dop-translation/github-configurations/{id}.
func (s *DopTranslationService) UpdateGithubConfiguration(ctx context.Context, id string, opt *DopTranslationGithubConfigurationUpdateOptions) (*DopTranslationGithubConfiguration, *http.Response, error) {
	path, err := v2PathWithID("dop-translation/github-configurations", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateUpdateGithubConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationGithubConfiguration](ctx, s.client, http.MethodPatch, path, nil, opt)
}

// DeleteGithubConfiguration deletes a GitHub configuration.
//
// API endpoint: DELETE /api/v2/dop-translation/github-configurations/{id}.
func (s *DopTranslationService) DeleteGithubConfiguration(ctx context.Context, id string) (*http.Response, error) {
	path, err := v2PathWithID("dop-translation/github-configurations", id, "Id")
	if err != nil {
		return nil, err
	}

	return deleteV2(ctx, s.client, path)
}

// -----------------------------------------------------------------------------
// GitLab configurations
// -----------------------------------------------------------------------------

// SearchGitlabConfigurations lists the GitLab configurations.
//
// API endpoint: GET /api/v2/dop-translation/gitlab-configurations.
func (s *DopTranslationService) SearchGitlabConfigurations(ctx context.Context) (*DopTranslationGitlabConfigurationsSearch, *http.Response, error) {
	return callV2[DopTranslationGitlabConfigurationsSearch](ctx, s.client, http.MethodGet, "dop-translation/gitlab-configurations", nil, nil)
}

// CreateGitlabConfiguration creates a GitLab configuration.
//
// API endpoint: POST /api/v2/dop-translation/gitlab-configurations.
func (s *DopTranslationService) CreateGitlabConfiguration(ctx context.Context, opt *DopTranslationGitlabConfigurationCreateOptions) (*DopTranslationGitlabConfiguration, *http.Response, error) {
	err := s.ValidateCreateGitlabConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationGitlabConfiguration](ctx, s.client, http.MethodPost, "dop-translation/gitlab-configurations", nil, opt)
}

// GetGitlabConfiguration fetches a GitLab configuration.
//
// API endpoint: GET /api/v2/dop-translation/gitlab-configurations/{id}.
func (s *DopTranslationService) GetGitlabConfiguration(ctx context.Context, id string) (*DopTranslationGitlabConfiguration, *http.Response, error) {
	path, err := v2PathWithID("dop-translation/gitlab-configurations", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationGitlabConfiguration](ctx, s.client, http.MethodGet, path, nil, nil)
}

// UpdateGitlabConfiguration updates a GitLab configuration.
//
// API endpoint: PATCH /api/v2/dop-translation/gitlab-configurations/{id}.
func (s *DopTranslationService) UpdateGitlabConfiguration(ctx context.Context, id string, opt *DopTranslationGitlabConfigurationUpdateOptions) (*DopTranslationGitlabConfiguration, *http.Response, error) {
	path, err := v2PathWithID("dop-translation/gitlab-configurations", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateUpdateGitlabConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationGitlabConfiguration](ctx, s.client, http.MethodPatch, path, nil, opt)
}

// DeleteGitlabConfiguration deletes a GitLab configuration.
//
// API endpoint: DELETE /api/v2/dop-translation/gitlab-configurations/{id}.
func (s *DopTranslationService) DeleteGitlabConfiguration(ctx context.Context, id string) (*http.Response, error) {
	path, err := v2PathWithID("dop-translation/gitlab-configurations", id, "Id")
	if err != nil {
		return nil, err
	}

	return deleteV2(ctx, s.client, path)
}

// StartGitlabSynchronizationRun starts a GitLab synchronization run.
//
// API endpoint: POST /api/v2/dop-translation/gitlab-synchronization-runs.
func (s *DopTranslationService) StartGitlabSynchronizationRun(ctx context.Context) (*DopTranslationGitlabSynchronizationRun, *http.Response, error) {
	return callV2[DopTranslationGitlabSynchronizationRun](ctx, s.client, http.MethodPost, "dop-translation/gitlab-synchronization-runs", nil, nil)
}

// -----------------------------------------------------------------------------
// Permission mappings (GitHub and GitLab)
// -----------------------------------------------------------------------------

func (s *DopTranslationService) searchPermissionMappings(ctx context.Context, platform string) (*DopTranslationPermissionMappingsSearch, *http.Response, error) {
	return callV2[DopTranslationPermissionMappingsSearch](ctx, s.client, http.MethodGet, "dop-translation/"+platform+"-permission-mappings", nil, nil)
}

func (s *DopTranslationService) createPermissionMapping(ctx context.Context, platform string, opt *DopTranslationPermissionMappingCreateOptions) (*DopTranslationPermissionMapping, *http.Response, error) {
	err := s.ValidatePermissionMappingCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationPermissionMapping](ctx, s.client, http.MethodPost, "dop-translation/"+platform+"-permission-mappings", nil, opt)
}

func (s *DopTranslationService) updatePermissionMapping(ctx context.Context, platform, role string, opt *DopTranslationPermissionMappingUpdateOptions) (*DopTranslationPermissionMapping, *http.Response, error) {
	path, err := v2PathWithID("dop-translation/"+platform+"-permission-mappings", role, "Role")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return callV2[DopTranslationPermissionMapping](ctx, s.client, http.MethodPatch, path, nil, opt)
}

func (s *DopTranslationService) deletePermissionMapping(ctx context.Context, platform, role string) (*http.Response, error) {
	path, err := v2PathWithID("dop-translation/"+platform+"-permission-mappings", role, "Role")
	if err != nil {
		return nil, err
	}

	return deleteV2(ctx, s.client, path)
}

// SearchGithubPermissionMappings fetches the GitHub permission mappings.
//
// API endpoint: GET /api/v2/dop-translation/github-permission-mappings.
func (s *DopTranslationService) SearchGithubPermissionMappings(ctx context.Context) (*DopTranslationPermissionMappingsSearch, *http.Response, error) {
	return s.searchPermissionMappings(ctx, "github")
}

// CreateGithubPermissionMapping creates a permission mapping for a custom GitHub role.
//
// API endpoint: POST /api/v2/dop-translation/github-permission-mappings.
func (s *DopTranslationService) CreateGithubPermissionMapping(ctx context.Context, opt *DopTranslationPermissionMappingCreateOptions) (*DopTranslationPermissionMapping, *http.Response, error) {
	return s.createPermissionMapping(ctx, "github", opt)
}

// UpdateGithubPermissionMapping updates the permission mapping of a GitHub role.
//
// API endpoint: PATCH /api/v2/dop-translation/github-permission-mappings/{role}.
func (s *DopTranslationService) UpdateGithubPermissionMapping(ctx context.Context, role string, opt *DopTranslationPermissionMappingUpdateOptions) (*DopTranslationPermissionMapping, *http.Response, error) {
	return s.updatePermissionMapping(ctx, "github", role, opt)
}

// DeleteGithubPermissionMapping deletes the permission mapping of a custom GitHub role.
//
// API endpoint: DELETE /api/v2/dop-translation/github-permission-mappings/{role}.
func (s *DopTranslationService) DeleteGithubPermissionMapping(ctx context.Context, role string) (*http.Response, error) {
	return s.deletePermissionMapping(ctx, "github", role)
}

// SearchGitlabPermissionMappings fetches the GitLab permission mappings.
//
// API endpoint: GET /api/v2/dop-translation/gitlab-permission-mappings.
func (s *DopTranslationService) SearchGitlabPermissionMappings(ctx context.Context) (*DopTranslationPermissionMappingsSearch, *http.Response, error) {
	return s.searchPermissionMappings(ctx, "gitlab")
}

// CreateGitlabPermissionMapping creates a permission mapping for a custom GitLab role.
//
// API endpoint: POST /api/v2/dop-translation/gitlab-permission-mappings.
func (s *DopTranslationService) CreateGitlabPermissionMapping(ctx context.Context, opt *DopTranslationPermissionMappingCreateOptions) (*DopTranslationPermissionMapping, *http.Response, error) {
	return s.createPermissionMapping(ctx, "gitlab", opt)
}

// UpdateGitlabPermissionMapping updates the permission mapping of a GitLab role.
//
// API endpoint: PATCH /api/v2/dop-translation/gitlab-permission-mappings/{role}.
func (s *DopTranslationService) UpdateGitlabPermissionMapping(ctx context.Context, role string, opt *DopTranslationPermissionMappingUpdateOptions) (*DopTranslationPermissionMapping, *http.Response, error) {
	return s.updatePermissionMapping(ctx, "gitlab", role, opt)
}

// DeleteGitlabPermissionMapping deletes the permission mapping of a custom GitLab role.
//
// API endpoint: DELETE /api/v2/dop-translation/gitlab-permission-mappings/{role}.
func (s *DopTranslationService) DeleteGitlabPermissionMapping(ctx context.Context, role string) (*http.Response, error) {
	return s.deletePermissionMapping(ctx, "gitlab", role)
}

// -----------------------------------------------------------------------------
// Project bindings
// -----------------------------------------------------------------------------

// SearchProjectBindings searches across project bindings.
//
// API endpoint: GET /api/v2/dop-translation/project-bindings.
func (s *DopTranslationService) SearchProjectBindings(ctx context.Context, opt *DopTranslationProjectBindingsSearchOptions) (*DopTranslationProjectBindingsSearch, *http.Response, error) {
	err := s.ValidateProjectBindingsSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationProjectBindingsSearch](ctx, s.client, http.MethodGet, "dop-translation/project-bindings", opt, nil)
}

// GetProjectBinding fetches a single project binding.
//
// API endpoint: GET /api/v2/dop-translation/project-bindings/{id}.
func (s *DopTranslationService) GetProjectBinding(ctx context.Context, id string) (*DopTranslationProjectBinding, *http.Response, error) {
	path, err := v2PathWithID("dop-translation/project-bindings", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	return callV2[DopTranslationProjectBinding](ctx, s.client, http.MethodGet, path, nil, nil)
}
