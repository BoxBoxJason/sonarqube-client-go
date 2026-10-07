package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// DopTranslationCheckPermissionsOptions contains parameters for the CheckPermissions method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationCheckPermissionsOptions struct {
	// Project key of the project whose bound DevOps Platform should be checked.
	Project string `json:"project,omitempty"`
	// Configuration key of a single DevOps Platform configuration to check, instead of all of
	// them. Cannot be combined with 'project'.
	Configuration string `json:"configuration,omitempty"`
	// Refresh bypass the cache and run a live check, then store its result. Requires
	// 'configuration' and the 'Administer System' permission.
	// Default: False.
	Refresh *bool `json:"refresh,omitempty"`
}

// DopTranslationPermissionDeficitResource represents a SonarQube V2 API object.
//
// A permission the Remediation Agent requires that is not granted at the level it needs.
type DopTranslationPermissionDeficitResource struct {
	// Permission name of the DevOps Platform permission, e.g. 'contents'.
	Permission string `json:"permission,omitempty"`
	// Required access level the Remediation Agent requires, e.g. 'write'.
	Required string `json:"required,omitempty"`
	// Granted access level actually granted, or null when the permission is not granted at all.
	Granted string `json:"granted,omitempty"`
}

// DopTranslationAffectedInstallationResource represents a SonarQube V2 API object.
//
// A GitHub App installation that has not approved every permission the Remediation Agent requires.
// Only permissions the app itself requests appear here; a permission the app does not request is
// reported once in 'appMissingPermissions' instead, because no installation owner can approve it.
// Suspension alone does not put an installation in this list.
type DopTranslationAffectedInstallationResource struct {
	// InstallationId identifier of the GitHub App installation.
	InstallationId string `json:"installationId,omitempty"`
	// InstallationOwner login of whoever the app is installed on. May be an organization or a
	// personal account.
	InstallationOwner string `json:"installationOwner,omitempty"`
	// SettingsUrl URL of the installation settings page, where the missing permissions can be
	// approved.
	SettingsUrl string `json:"settingsUrl,omitempty"`
	// MissingPermissions permissions this installation still has to approve.
	MissingPermissions []DopTranslationPermissionDeficitResource `json:"missingPermissions,omitempty"`
}

// DopTranslationPermissionCheckResource represents the PermissionCheckResource object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationPermissionCheckResource struct {
	// Key of the DevOps Platform configuration.
	Key string `json:"key,omitempty"`
	// Type devOps Platform type of the configuration (github, gitlab, azure).
	Type string `json:"type,omitempty"`
	// Status whether the configuration grants the write permissions the Remediation Agent needs to
	// clone and open pull requests. Allowed values: SUFFICIENT, INSUFFICIENT, UNKNOWN,
	// CHECK_FAILED, UNSUPPORTED_TOKEN_TYPE.
	Status string `json:"status,omitempty"`
	// CheckedAt epoch milliseconds when this status was last computed. May reflect a cached
	// result.
	CheckedAt int64 `json:"checkedAt,omitempty"`
	// AppMissingPermissions gitHub only. Permissions that the GitHub App does not request. Add
	// these permissions to the app before an installation can approve them. If the system cannot
	// read the app configuration, 'status' is CHECK_FAILED. In this case, an empty list means that
	// no app permission result is available.
	AppMissingPermissions []DopTranslationPermissionDeficitResource `json:"appMissingPermissions,omitempty"`
	// InstallationCheckStatus gitHub only. Whether the per-installation scan produced a
	// trustworthy answer. The counts and the list are authoritative only when this is COMPLETE. On
	// a project request, NOT_INSTALLED means the app is not installed on that project's repository
	// — a definite answer, unlike FAILED. NOT_RUN means that the system did not check an
	// installation because the project has no bound repository. In this case, 'status' covers only
	// the app configuration. Allowed values: COMPLETE, FAILED, NOT_INSTALLED, NOT_RUN.
	InstallationCheckStatus string `json:"installationCheckStatus,omitempty"`
	// TotalInstallationCount gitHub administrator requests only. How many installations the app
	// has at all. Tells an app nobody has installed apart from an app whose installations have all
	// approved, which both report zero affected installations. Omitted unless
	// 'installationCheckStatus' is COMPLETE.
	TotalInstallationCount int32 `json:"totalInstallationCount,omitempty"`
	// AffectedInstallationCount gitHub administrator requests only. Exact number of installations
	// with missing permissions. Omitted unless 'installationCheckStatus' is COMPLETE.
	AffectedInstallationCount int32 `json:"affectedInstallationCount,omitempty"`
	// AffectedInstallations gitHub administrator requests only. Every installation with missing
	// permissions, in a stable order. Omitted unless 'installationCheckStatus' is COMPLETE.
	AffectedInstallations []DopTranslationAffectedInstallationResource `json:"affectedInstallations,omitempty"`
}

// DopTranslationPermissionChecksRestResponse represents the PermissionChecksRestResponse object of
// the SonarQube V2 API.
type DopTranslationPermissionChecksRestResponse struct {
	// PermissionChecks is the permission checks.
	PermissionChecks []DopTranslationPermissionCheckResource `json:"permissionChecks,omitempty"`
}

// DopTranslationGenerateScmAccessTokenOptions contains parameters for the GenerateScmAccessToken method.
type DopTranslationGenerateScmAccessTokenOptions struct {
	// Project key. This field is required.
	Project string `json:"project"`
}

// DopTranslationScmAccessTokenRestResponse represents the ScmAccessTokenRestResponse object of the
// SonarQube V2 API.
type DopTranslationScmAccessTokenRestResponse struct {
	// Alm is the alm.
	Alm string `json:"alm,omitempty"`
	// Username is the username.
	Username string `json:"username,omitempty"`
	// Secret is the secret.
	Secret string `json:"secret,omitempty"`
	// ExpiresAt is the expires at.
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGenerateScmAccessTokenOpt validates the options for the GenerateScmAccessToken method.
func (s *DopTranslationService) ValidateGenerateScmAccessTokenOpt(opt *DopTranslationGenerateScmAccessTokenOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CheckPermissions checks DevOps Platform permissions for the Remediation Agent. Validates whether
// the configured DevOps Platforms (GitHub, GitLab, Azure DevOps) grant the write permissions the
// SonarQube Remediation Agent needs to clone a repository, push a branch and open a pull/merge
// request. Results are cached for a short time. Internal endpoint used by the Remediation Agent UI.
// Without 'project' it checks configurations instance-wide and requires the 'Administer System'
// permission or trusted privileged-service authentication; 'configuration' narrows that to a single
// configuration. With 'project' it checks the platform bound to that project and requires 'Browse'
// permission on the project; a GitHub project is checked against the installation covering its own
// repository, so two projects sharing one configuration can get different results. 'project' and
// 'configuration' cannot be combined. For GitHub, the response also reports permissions the app
// itself is not configured to request, and — for administrator requests — every installation that
// has not approved the permissions it does request, whether an organization or a personal account
// owns it. The counts and list are valid only when 'installationCheckStatus' is COMPLETE. If the
// scan cannot finish, 'installationCheckStatus' is FAILED and 'status' is CHECK_FAILED. The
// response does not include the counts or list.
//
// API endpoint: GET /api/v2/dop-translation/permission-checks.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *DopTranslationService) CheckPermissions(ctx context.Context, opt *DopTranslationCheckPermissionsOptions) (*DopTranslationPermissionChecksRestResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "dop-translation/permission-checks", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DopTranslationPermissionChecksRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GenerateScmAccessToken provides a scoped SCM access token. Provide a scoped git credential for
// the given project's bound DevOps Platform (GitHub or GitLab). GitHub credentials are short-lived
// and minted per call; GitLab credentials are cached in memory per SonarQube Server node and
// refreshed before expiry. Internal endpoint used by the agentic-workflows remediation orchestrator
// (SONAR-31165). Requires the 'Administer System' permission or trusted privileged-service
// authentication.
//
// API endpoint: POST /api/v2/dop-translation/scm-access-tokens.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *DopTranslationService) GenerateScmAccessToken(ctx context.Context, opt *DopTranslationGenerateScmAccessTokenOptions) (*DopTranslationScmAccessTokenRestResponse, *http.Response, error) {
	err := s.ValidateGenerateScmAccessTokenOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "dop-translation/scm-access-tokens", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DopTranslationScmAccessTokenRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
