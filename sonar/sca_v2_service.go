package sonar

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// ScaService handles communication with the Software Composition Analysis (SCA) V2 API endpoints.
// This service is only available in Enterprise Edition with SCA enabled.
type ScaService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ScaCliInfo represents metadata about a downloadable SCA CLI binary.
type ScaCliInfo struct {
	// Id is the unique identifier of this CLI release.
	Id string `json:"id,omitempty"`
	// Filename is the filename of the CLI binary.
	Filename string `json:"filename,omitempty"`
	// Sha256 is the SHA-256 hash of the binary.
	Sha256 string `json:"sha256,omitempty"`
	// Os is the operating system target.
	Os string `json:"os,omitempty"`
	// Arch is the CPU architecture target.
	Arch string `json:"arch,omitempty"`
}

// ScaFeatureEnablement represents the SCA enablement state.
type ScaFeatureEnablement struct {
	// Enablement indicates whether SCA is enabled.
	Enablement bool `json:"enablement"`
}

// ScaReleaseSearchResource represents a release in the SCA search results.
type ScaReleaseSearchResource struct {
	// Key is the unique release key.
	Key string `json:"key,omitempty"`
	// PackageUrl is the package URL (purl).
	PackageUrl string `json:"packageUrl,omitempty"`
	// PackageManager is the package manager (e.g., npm, maven).
	PackageManager string `json:"packageManager,omitempty"`
	// PackageName is the package name.
	PackageName string `json:"packageName,omitempty"`
	// Version is the package version.
	Version string `json:"version,omitempty"`
	// LicenseExpression is the SPDX license expression.
	LicenseExpression string `json:"licenseExpression,omitempty"`
	// ScopeSummary is a summary of dependency scopes.
	ScopeSummary string `json:"scopeSummary,omitempty"`
	// DependencyFilePaths lists the dependency files that reference this release.
	DependencyFilePaths []string `json:"dependencyFilePaths,omitempty"`
	// Known indicates whether the package is known to SonarQube SCA.
	Known bool `json:"known"`
	// NewlyIntroduced indicates whether this release was newly introduced vs. the target branch.
	NewlyIntroduced bool `json:"newlyIntroduced"`
	// DirectSummary indicates whether this is a direct dependency.
	DirectSummary bool `json:"directSummary"`
}

// ScaVulnerabilityResource represents a vulnerability associated with a release.
type ScaVulnerabilityResource struct {
	// VulnerabilityId is the CVE or other vulnerability identifier.
	VulnerabilityId string `json:"vulnerabilityId,omitempty"`
	// Description is the vulnerability description.
	Description string `json:"description,omitempty"`
	// EpssPercentile is the EPSS percentile score.
	EpssPercentile string `json:"epssPercentile,omitempty"`
	// EpssProbability is the EPSS probability score.
	EpssProbability string `json:"epssProbability,omitempty"`
	// CweIds lists the associated CWE identifiers.
	CweIds []string `json:"cweIds,omitempty"`
	// KnownExploited indicates whether this vulnerability is known to be exploited.
	KnownExploited bool `json:"knownExploited"`
}

// ScaDependencyRisk represents an issue-release pair (a vulnerability in a dependency).
type ScaDependencyRisk struct {
	// Key is the unique identifier for this issue-release pair.
	Key string `json:"key,omitempty"`
	// Severity is the severity level.
	Severity string `json:"severity,omitempty"`
	// Type is the risk type (e.g., VULNERABILITY, LICENSE).
	Type string `json:"type,omitempty"`
	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`
	// SpdxLicenseId is the SPDX license identifier for license risks.
	SpdxLicenseId string `json:"spdxLicenseId,omitempty"`
	// Release is the release this risk is associated with.
	Release ScaReleaseSearchResource `json:"release,omitzero"`
	// Vulnerability contains vulnerability details for CVE-type risks.
	Vulnerability ScaVulnerabilityResource `json:"vulnerability,omitzero"`
}

// ScaDependencyRisksSearch represents the response from the issues-releases search endpoint.
type ScaDependencyRisksSearch struct {
	// IssuesReleases is the list of issue-release pairs.
	IssuesReleases []ScaDependencyRisk `json:"issuesReleases,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ScaReleasesSearch represents the response from the releases search endpoint.
type ScaReleasesSearch struct {
	// Releases is the list of releases.
	Releases []ScaReleaseSearchResource `json:"releases,omitempty"`
	// PackageManagerCounts summarizes the number of releases per package manager.
	PackageManagerCounts []ScaReleasePackageManagerCount `json:"packageManagerCounts,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ScaReleasePackageManagerCount represents the number of releases for a given package manager.
type ScaReleasePackageManagerCount struct {
	// PackageManager is the package manager (e.g., npm, maven).
	PackageManager string `json:"packageManager,omitempty"`
	// ReleaseCount is the number of releases for this package manager.
	ReleaseCount int32 `json:"releaseCount"`
}

// ScaReleaseDetail represents detailed information about a release.
type ScaReleaseDetail struct {
	// Key is the unique release key.
	Key string `json:"key,omitempty"`
	// PackageUrl is the package URL (purl).
	PackageUrl string `json:"packageUrl,omitempty"`
	// PackageManager is the package manager.
	PackageManager string `json:"packageManager,omitempty"`
	// PackageName is the package name.
	PackageName string `json:"packageName,omitempty"`
	// Version is the package version.
	Version string `json:"version,omitempty"`
	// LicenseExpression is the SPDX license expression.
	LicenseExpression string `json:"licenseExpression,omitempty"`
	// ScopeSummary is a summary of dependency scopes.
	ScopeSummary string `json:"scopeSummary,omitempty"`
	// Known indicates whether the package is known.
	Known bool `json:"known"`
	// NewlyIntroduced indicates whether this release was newly introduced vs. the target branch.
	NewlyIntroduced bool `json:"newlyIntroduced"`
	// DirectSummary indicates whether this is a direct dependency.
	DirectSummary bool `json:"directSummary"`
	// ProductionScopeSummary indicates whether this is in production scope.
	ProductionScopeSummary bool `json:"productionScopeSummary"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ScaCliListOptions contains parameters for the ListClis method.
type ScaCliListOptions struct {
	// Os filters by operating system. Optional.
	Os string `json:"os,omitempty"`
	// Arch filters by CPU architecture. Optional.
	Arch string `json:"arch,omitempty"`
}

// ScaCliGetOptions contains parameters for the GetCli method.
type ScaCliGetOptions struct {
	// Id is the CLI download identifier. This field is required.
	Id string `json:"id"`
}

// ScaSetEnablementOptions contains parameters for the SetEnablement method.
type ScaSetEnablementOptions struct {
	// Enablement indicates whether to enable SCA. This field is required.
	Enablement bool `json:"enablement"`
}

// ScaDependencyRisksSearchOptions contains parameters for the SearchDependencyRisks method.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type ScaDependencyRisksSearchOptions struct {
	PaginationParamsV2

	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey filters by branch. Optional.
	BranchKey string `json:"branchKey,omitempty"`
	// PullRequestKey filters by pull request. Optional.
	PullRequestKey string `json:"pullRequestKey,omitempty"`
	// PackageManagers filters by package managers (comma-separated). Optional.
	PackageManagers string `json:"packageManagers,omitempty"`
	// Types filters by risk type. Optional.
	Types string `json:"types,omitempty"`
	// Severities filters by severity. Optional.
	Severities string `json:"severities,omitempty"`
	// PackageName filters by package name. Optional.
	PackageName string `json:"packageName,omitempty"`
	// VulnerabilityId filters by CVE/vulnerability ID. Optional.
	VulnerabilityId string `json:"vulnerabilityId,omitempty"`
	// Sort specifies the sort order. Optional.
	Sort string `json:"sort,omitempty"`
	// NewlyIntroduced filters for risks newly introduced vs. the target branch. Optional.
	NewlyIntroduced *bool `json:"newlyIntroduced,omitempty"`
	// Direct filters for direct dependencies only. Optional.
	Direct *bool `json:"direct,omitempty"`
	// ProductionScope filters by production scope. Optional.
	ProductionScope *bool `json:"productionScope,omitempty"`
	// Assigned filters for assigned (or unassigned, when false) issues. Optional.
	Assigned *bool `json:"assigned,omitempty"`
	// OnlyShowConfirmedReachable filters for issues confirmed as reachable. Optional.
	OnlyShowConfirmedReachable *bool `json:"onlyShowConfirmedReachable,omitempty"`
	// SecurityStandard filters by security standard. Optional.
	SecurityStandard string `json:"securityStandard,omitempty"`
	// SecurityStandardCategory filters by security standard category. Optional.
	SecurityStandardCategory string `json:"securityStandardCategory,omitempty"`
	// SecurityStandardVersion filters by security standard version. Optional.
	SecurityStandardVersion string `json:"securityStandardVersion,omitempty"`
	// FixedInPullRequestKey filters for issues fixed in the given pull request. Optional.
	FixedInPullRequestKey string `json:"fixedInPullRequestKey,omitempty"`
	// Qualities filters by component qualifier. Optional.
	Qualities []string `json:"qualities,omitempty"`
	// Statuses filters by issue status. Optional.
	Statuses []string `json:"statuses,omitempty"`
	// Assignees filters by assignee logins. Optional.
	Assignees []string `json:"assignees,omitempty"`
	// Projects filters by project keys. Optional.
	Projects []string `json:"projects,omitempty"`
}

// ScaDependencyRiskGetOptions contains parameters for the GetDependencyRisk method.
type ScaDependencyRiskGetOptions struct {
	// Key is the issue-release key. This field is required.
	Key string `json:"key"`
}

// ScaReleasesSearchOptions contains parameters for the SearchReleases method.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type ScaReleasesSearchOptions struct {
	PaginationParamsV2

	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey filters by branch. Optional.
	BranchKey string `json:"branchKey,omitempty"`
	// PullRequestKey filters by pull request. Optional.
	PullRequestKey string `json:"pullRequestKey,omitempty"`
	// PackageManagers filters by package managers. Optional.
	PackageManagers string `json:"packageManagers,omitempty"`
	// Q filters by package name query. Optional.
	Q string `json:"q,omitempty"`
	// NewlyIntroduced filters for releases newly introduced vs. the target branch. Optional.
	NewlyIntroduced *bool `json:"newlyIntroduced,omitempty"`
	// Direct filters for direct dependencies. Optional.
	Direct *bool `json:"direct,omitempty"`
	// ProductionScope filters by production scope. Optional.
	ProductionScope *bool `json:"productionScope,omitempty"`
}

// ScaReleaseGetOptions contains parameters for the GetRelease method.
type ScaReleaseGetOptions struct {
	// Key is the release key. This field is required.
	Key string `json:"key"`
}

// SBOM report type values accepted by the "type" query parameter of GetSbomReport.
const (
	// ScaSbomReportTypeCycloneDX requests a CycloneDX-formatted SBOM report.
	ScaSbomReportTypeCycloneDX = "cyclonedx"
	// ScaSbomReportTypeSpdx23 requests an SPDX 2.3-formatted SBOM report.
	ScaSbomReportTypeSpdx23 = "spdx_23"
	// ScaSbomReportTypeSpdx30 requests an SPDX 3.0-formatted SBOM report. Only the JSON format is
	// supported for this report type.
	ScaSbomReportTypeSpdx30 = "spdx_30"

	// ScaSbomReportFormatJSON requests the report serialized as JSON.
	ScaSbomReportFormatJSON = "json"
	// ScaSbomReportFormatXML requests the report serialized as XML.
	ScaSbomReportFormatXML = "xml"
)

// ScaSbomReportOptions contains parameters for the GetSbomReport method.
type ScaSbomReportOptions struct {
	// Component is the key of the component (project, application, portfolio) to build the
	// report for. This field is required.
	Component string `json:"component"`
	// Branch filters by branch. Optional.
	Branch string `json:"branch,omitempty"`
	// Type is the SBOM report type, one of ScaSbomReportTypeCycloneDX, ScaSbomReportTypeSpdx23 or
	// ScaSbomReportTypeSpdx30. This field is required.
	Type string `json:"type"`
	// OnlyProductionScope filters to production dependencies only when true. If false, all
	// dependencies are included. Optional; defaults to true server-side when absent.
	OnlyProductionScope *bool `json:"onlyProductionScope,omitempty"`
	// Format is the desired report serialization, one of ScaSbomReportFormatJSON or
	// ScaSbomReportFormatXML. It is sent as the request's Accept header rather than a query
	// parameter, since the SonarQube API selects the report's MIME type via content negotiation.
	// This field is required.
	Format string `json:"-"`
}

// scaSbomAcceptHeader builds the Accept header value expected by the SBOM report endpoint for
// the given report type and format.
func scaSbomAcceptHeader(reportType, format string) string {
	switch reportType {
	case ScaSbomReportTypeSpdx23, ScaSbomReportTypeSpdx30:
		return "application/spdx+" + format
	case ScaSbomReportTypeCycloneDX:
		return "application/vnd.cyclonedx+" + format
	default:
		return "application/" + reportType + "+" + format
	}
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCliGetOpt validates the options for the GetCli method.
func (s *ScaService) ValidateCliGetOpt(opt *ScaCliGetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Id, "Id")
}

// ValidateSetEnablementOpt validates the options for the SetEnablement method.
func (s *ScaService) ValidateSetEnablementOpt(opt *ScaSetEnablementOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return nil
}

// ValidateSearchDependencyRisksOpt validates the options for the SearchDependencyRisks method.
func (s *ScaService) ValidateSearchDependencyRisksOpt(opt *ScaDependencyRisksSearchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return opt.Validate()
}

// ValidateGetDependencyRiskOpt validates the options for the GetDependencyRisk method.
func (s *ScaService) ValidateGetDependencyRiskOpt(opt *ScaDependencyRiskGetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateSearchReleasesOpt validates the options for the SearchReleases method.
func (s *ScaService) ValidateSearchReleasesOpt(opt *ScaReleasesSearchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return opt.Validate()
}

// ValidateGetReleaseOpt validates the options for the GetRelease method.
func (s *ScaService) ValidateGetReleaseOpt(opt *ScaReleaseGetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateSbomReportOpt validates the options for the GetSbomReport method.
func (s *ScaService) ValidateSbomReportOpt(opt *ScaSbomReportOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Format, "Format")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ListClis returns available SCA CLI downloads.
// Requires 'Browse' permission.
//
// API endpoint: GET /api/v2/sca/clis.
// Enterprise Edition only.
func (s *ScaService) ListClis(ctx context.Context, opt *ScaCliListOptions) ([]ScaCliInfo, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/clis", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ScaCliInfo

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetCli returns metadata for a specific SCA CLI download.
// Requires 'Browse' permission.
//
// API endpoint: GET /api/v2/sca/clis/{id}.
// Enterprise Edition only.
func (s *ScaService) GetCli(ctx context.Context, opt *ScaCliGetOptions) (*ScaCliInfo, *http.Response, error) {
	err := s.ValidateCliGetOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/clis/"+opt.Id, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaCliInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetEnablement returns the SCA feature enablement status.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/sca/feature-enablements.
// Enterprise Edition only.
func (s *ScaService) GetEnablement(ctx context.Context) (*ScaFeatureEnablement, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/feature-enablements", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaFeatureEnablement)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetEnablement updates the SCA feature enablement settings.
// Requires 'Administer System' permission.
//
// API endpoint: PATCH /api/v2/sca/feature-enablements.
// Enterprise Edition only.
func (s *ScaService) SetEnablement(ctx context.Context, opt *ScaSetEnablementOptions) (*ScaFeatureEnablement, *http.Response, error) {
	err := s.ValidateSetEnablementOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/feature-enablements", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaFeatureEnablement)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchDependencyRisks searches for issue-release pairs (dependency risks).
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/sca/issues-releases.
// Enterprise Edition only.
func (s *ScaService) SearchDependencyRisks(ctx context.Context, opt *ScaDependencyRisksSearchOptions) (*ScaDependencyRisksSearch, *http.Response, error) {
	err := s.ValidateSearchDependencyRisksOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/issues-releases", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaDependencyRisksSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetDependencyRisk returns details for a single issue-release pair.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/sca/issues-releases/{key}.
// Enterprise Edition only.
func (s *ScaService) GetDependencyRisk(ctx context.Context, opt *ScaDependencyRiskGetOptions) (*ScaDependencyRisk, *http.Response, error) {
	err := s.ValidateGetDependencyRiskOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/issues-releases/"+opt.Key, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaDependencyRisk)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchReleases searches for dependency releases in a project.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/sca/releases.
// Enterprise Edition only.
func (s *ScaService) SearchReleases(ctx context.Context, opt *ScaReleasesSearchOptions) (*ScaReleasesSearch, *http.Response, error) {
	err := s.ValidateSearchReleasesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/releases", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaReleasesSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetRelease returns details for a single dependency release.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/sca/releases/{key}.
// Enterprise Edition only.
func (s *ScaService) GetRelease(ctx context.Context, opt *ScaReleaseGetOptions) (*ScaReleaseDetail, *http.Response, error) {
	err := s.ValidateGetReleaseOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/releases/"+opt.Key, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaReleaseDetail)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSbomReport returns a Software Bill of Materials (SBOM) report for a project. The report's
// MIME type is selected via content negotiation: opt.Type picks the report standard
// (ScaSbomReportTypeCycloneDX or ScaSbomReportTypeSpdx23) and opt.Format picks the serialization
// (ScaSbomReportFormatJSON or ScaSbomReportFormatXML); together they set the request's Accept
// header, which the server requires to disambiguate the response format.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/sca/sbom-reports.
// Enterprise Edition only.
func (s *ScaService) GetSbomReport(ctx context.Context, opt *ScaSbomReportOptions) ([]byte, *http.Response, error) {
	err := s.ValidateSbomReportOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/sbom-reports", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", scaSbomAcceptHeader(opt.Type, opt.Format))

	var buf bytes.Buffer

	resp, err := s.client.Do(req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, nil
}

// -----------------------------------------------------------------------------
// License Profiles - Response Types
// -----------------------------------------------------------------------------

// ScaLicenseProfileActions describes the actions available to the current user on a license
// profile.
type ScaLicenseProfileActions struct {
	// Edit is true if the current user can edit the license profile.
	Edit bool `json:"edit"`
	// SetAsDefault is true if the current user can set the license profile as the default.
	SetAsDefault bool `json:"setAsDefault"`
	// AssociateProjects is true if the current user can associate projects with the license
	// profile.
	AssociateProjects bool `json:"associateProjects"`
	// Delete is true if the current user can delete the license profile.
	Delete bool `json:"delete"`
}

// ScaLicenseProfile represents a license profile used to evaluate license issues.
type ScaLicenseProfile struct {
	// Id is the license profile's unique identifier.
	Id string `json:"id,omitempty"`
	// Name is the license profile's display name.
	Name string `json:"name,omitempty"`
	// UpdatedAt is the timestamp of the license profile's last update.
	UpdatedAt string `json:"updatedAt,omitempty"`
	// Actions describes the actions available to the current user on this license profile.
	Actions ScaLicenseProfileActions `json:"actions,omitzero"`
	// Default is true if this is the default license profile.
	Default bool `json:"default"`
}

// ScaLicenseProfileCollectionActions describes the actions available to the current user on the
// collection of license profiles.
type ScaLicenseProfileCollectionActions struct {
	// Create is true if the current user can create a new license profile.
	Create bool `json:"create"`
}

// ScaLicenseProfileIndex represents the response from the license profiles list endpoint.
type ScaLicenseProfileIndex struct {
	// LicenseProfiles is the list of configured license profiles.
	LicenseProfiles []ScaLicenseProfile `json:"licenseProfiles,omitempty"`
	// Actions describes the actions available to the current user on the license profiles
	// collection.
	Actions ScaLicenseProfileCollectionActions `json:"actions,omitzero"`
}

// ScaLicenseProfileCreateRequest contains the parameters for the CreateLicenseProfile method.
type ScaLicenseProfileCreateRequest struct {
	// Name is the name of the license profile to create. This field is required.
	Name string `json:"name"`
	// Organization is the key of the organization the license profile belongs to. This field is
	// required.
	Organization string `json:"organization"`
	// Default, if true, sets the license profile as the default. Optional.
	Default bool `json:"default,omitempty"`
}

// ScaLicenseProfileAssignableProject represents a project that can be assigned to a license
// profile.
type ScaLicenseProfileAssignableProject struct {
	// Id is the project's unique identifier.
	Id string `json:"id,omitempty"`
	// ProjectKey is the project's key.
	ProjectKey string `json:"projectKey,omitempty"`
	// ProjectName is the project's display name.
	ProjectName string `json:"projectName,omitempty"`
	// AssignedToLicenseProfile is true if the project is already assigned to the license profile.
	AssignedToLicenseProfile bool `json:"assignedToLicenseProfile"`
}

// ScaLicenseProfileAssignableProjectsResponse represents the response from the license profile
// assignable projects endpoint.
type ScaLicenseProfileAssignableProjectsResponse struct {
	// AssignableProjects is the list of projects that can be assigned to the license profile.
	AssignableProjects []ScaLicenseProfileAssignableProject `json:"assignableProjects,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ScaLicenseProfileAssignmentRequest contains the parameters for the AssignLicenseProfileProject
// method.
type ScaLicenseProfileAssignmentRequest struct {
	// LicenseProfileUuid is the key of the license profile to assign the project to. Optional.
	//
	// Deprecated: use LicenseProfileId instead.
	LicenseProfileUuid string `json:"licenseProfileUuid,omitempty"`
	// LicenseProfileId is the key of the license profile to assign the project to. Optional.
	LicenseProfileId string `json:"licenseProfileId,omitempty"`
	// ProjectKey is the key of the project to assign to the license profile. This field is
	// required.
	ProjectKey string `json:"projectKey"`
}

// ScaLicenseProfileCategory represents the policy configured for an entire license category
// within a license profile.
type ScaLicenseProfileCategory struct {
	// Id is the category's unique identifier.
	Id string `json:"id,omitempty"`
	// Key is the license category key.
	// Allowed values: ScaLicenseCategoryUnknown, ScaLicenseCategoryCopyleftWeak,
	// ScaLicenseCategoryCopyleftStrong, ScaLicenseCategoryCopyleftNetwork,
	// ScaLicenseCategoryCopyleftMaximal, ScaLicenseCategoryPermissiveStandard,
	// ScaLicenseCategoryPermissiveAmateur.
	Key string `json:"key,omitempty"`
	// Policy is the policy configured for this category.
	// Allowed values: ScaLicensePolicyDeny, ScaLicensePolicyAllow.
	Policy string `json:"policy,omitempty"`
}

// ScaLicensePolicyLicense represents the policy configured for a single license within a license
// profile.
type ScaLicensePolicyLicense struct {
	// Id is the license policy's unique identifier.
	Id string `json:"id,omitempty"`
	// SpdxLicenseId is the SPDX license identifier.
	SpdxLicenseId string `json:"spdxLicenseId,omitempty"`
	// Name is the license's display name.
	Name string `json:"name,omitempty"`
	// Category is the license category key.
	// Allowed values: ScaLicenseCategoryUnknown, ScaLicenseCategoryCopyleftWeak,
	// ScaLicenseCategoryCopyleftStrong, ScaLicenseCategoryCopyleftNetwork,
	// ScaLicenseCategoryCopyleftMaximal, ScaLicenseCategoryPermissiveStandard,
	// ScaLicenseCategoryPermissiveAmateur.
	Category string `json:"category,omitempty"`
	// Policy is the policy configured for this license.
	// Allowed values: ScaLicensePolicyDeny, ScaLicensePolicyAllow.
	Policy string `json:"policy,omitempty"`
}

// ScaLicenseProfileDetails represents the response from the get license profile endpoint,
// containing the profile itself plus its category and license policies.
type ScaLicenseProfileDetails struct {
	// Profile is the license profile.
	Profile ScaLicenseProfile `json:"profile,omitzero"`
	// Categories is the list of category policies configured for the license profile.
	Categories []ScaLicenseProfileCategory `json:"categories,omitempty"`
	// Licenses is the list of individual license policies configured for the license profile.
	Licenses []ScaLicensePolicyLicense `json:"licenses,omitempty"`
}

// ScaLicenseProfileUpdateRequest contains the parameters for the UpdateLicenseProfile method.
type ScaLicenseProfileUpdateRequest struct {
	// Default, if set, changes whether the license profile is the default. Optional.
	Default *bool `json:"default,omitempty"`
	// Name is the new name of the license profile. Optional.
	Name string `json:"name,omitempty"`
}

// ScaLicenseProfileCategoryUpdateRequest contains the parameters for the
// UpdateLicenseProfileCategory method.
type ScaLicenseProfileCategoryUpdateRequest struct {
	// Policy is the policy to apply to the category. This field is required.
	// Allowed values: ScaLicensePolicyDeny, ScaLicensePolicyAllow.
	Policy string `json:"policy"`
}

// ScaLicensePolicyLicenseUpdateRequest contains the parameters for the
// UpdateLicenseProfileLicense method.
type ScaLicensePolicyLicenseUpdateRequest struct {
	// Policy is the policy to apply to the license. This field is required.
	// Allowed values: ScaLicensePolicyDeny, ScaLicensePolicyAllow.
	Policy string `json:"policy"`
}

// -----------------------------------------------------------------------------
// License Profiles - Option Types
// -----------------------------------------------------------------------------

// ScaLicenseProfileListOptions contains parameters for the ListLicenseProfiles method.
type ScaLicenseProfileListOptions struct {
	// ProjectKey, if provided, filters to the license profile that is used to analyze the project
	// (if one exists). Optional.
	ProjectKey string `json:"projectKey,omitempty"`
}

// ScaLicenseProfileAssignableProjectsOptions contains parameters for the
// ListLicenseProfileAssignableProjects method.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type ScaLicenseProfileAssignableProjectsOptions struct {
	PaginationParamsV2

	AssignedToLicenseProfile *bool  `json:"assignedToLicenseProfile,omitempty"`
	LicenseProfileUuid       string `json:"licenseProfileUuid,omitempty"`
	LicenseProfileId         string `json:"licenseProfileId,omitempty"`
	Q                        string `json:"q,omitempty"`
}

// ScaLicenseProfileAssignmentDeleteOptions contains parameters for the
// DeleteLicenseProfileAssignment method.
type ScaLicenseProfileAssignmentDeleteOptions struct {
	// ProjectKey is the key of the project whose license profile assignment should be removed. This
	// field is required.
	ProjectKey string `json:"projectKey"`
}

// ScaLicenseProfileGetOptions contains parameters for the GetLicenseProfile and
// DeleteLicenseProfile methods.
type ScaLicenseProfileGetOptions struct {
	// Key is the license profile key. This field is required.
	Key string `json:"key"`
}

// Allowed values for the "category-key" path segment of UpdateLicenseProfileCategory.
const (
	// ScaLicenseCategoryUnknown represents the "UNKNOWN" license category.
	ScaLicenseCategoryUnknown = "UNKNOWN"
	// ScaLicenseCategoryCopyleftWeak represents the "COPYLEFT_WEAK" license category.
	ScaLicenseCategoryCopyleftWeak = "COPYLEFT_WEAK"
	// ScaLicenseCategoryCopyleftStrong represents the "COPYLEFT_STRONG" license category.
	ScaLicenseCategoryCopyleftStrong = "COPYLEFT_STRONG"
	// ScaLicenseCategoryCopyleftNetwork represents the "COPYLEFT_NETWORK" license category.
	ScaLicenseCategoryCopyleftNetwork = "COPYLEFT_NETWORK"
	// ScaLicenseCategoryCopyleftMaximal represents the "COPYLEFT_MAXIMAL" license category.
	ScaLicenseCategoryCopyleftMaximal = "COPYLEFT_MAXIMAL"
	// ScaLicenseCategoryPermissiveStandard represents the "PERMISSIVE_STANDARD" license category.
	ScaLicenseCategoryPermissiveStandard = "PERMISSIVE_STANDARD"
	// ScaLicenseCategoryPermissiveAmateur represents the "PERMISSIVE_AMATEUR" license category.
	ScaLicenseCategoryPermissiveAmateur = "PERMISSIVE_AMATEUR"
)

//nolint:gochecknoglobals // constant set of allowed values
var allowedScaLicenseCategories = map[string]struct{}{
	ScaLicenseCategoryUnknown:            {},
	ScaLicenseCategoryCopyleftWeak:       {},
	ScaLicenseCategoryCopyleftStrong:     {},
	ScaLicenseCategoryCopyleftNetwork:    {},
	ScaLicenseCategoryCopyleftMaximal:    {},
	ScaLicenseCategoryPermissiveStandard: {},
	ScaLicenseCategoryPermissiveAmateur:  {},
}

// Allowed values for the "policy" field of license category and license policy update requests.
const (
	// ScaLicensePolicyDeny represents the "DENY" license policy.
	ScaLicensePolicyDeny = "DENY"
	// ScaLicensePolicyAllow represents the "ALLOW" license policy.
	ScaLicensePolicyAllow = "ALLOW"
)

//nolint:gochecknoglobals // constant set of allowed values
var allowedScaLicensePolicies = map[string]struct{}{
	ScaLicensePolicyDeny:  {},
	ScaLicensePolicyAllow: {},
}

// ScaLicenseProfileCategoryOptions contains parameters for the UpdateLicenseProfileCategory method.
type ScaLicenseProfileCategoryOptions struct {
	// LicenseProfileKey is the license profile key. This field is required.
	LicenseProfileKey string `json:"licenseProfileKey"`
	// CategoryKey is the license category key. This field is required.
	// Allowed values: ScaLicenseCategoryUnknown, ScaLicenseCategoryCopyleftWeak,
	// ScaLicenseCategoryCopyleftStrong, ScaLicenseCategoryCopyleftNetwork,
	// ScaLicenseCategoryCopyleftMaximal, ScaLicenseCategoryPermissiveStandard,
	// ScaLicenseCategoryPermissiveAmateur.
	CategoryKey string `json:"categoryKey"`
}

// ScaLicenseProfileLicenseOptions contains parameters for the UpdateLicenseProfileLicense method.
type ScaLicenseProfileLicenseOptions struct {
	// LicenseProfileKey is the license profile key. This field is required.
	LicenseProfileKey string `json:"licenseProfileKey"`
	// LicensePolicyId is the license ID. This field is required.
	LicensePolicyId string `json:"licensePolicyId"`
}

// -----------------------------------------------------------------------------
// License Profiles - Validation Functions
// -----------------------------------------------------------------------------

// ValidateLicenseProfileAssignmentDeleteOpt validates the options for the
// DeleteLicenseProfileAssignment method.
func (s *ScaService) ValidateLicenseProfileAssignmentDeleteOpt(opt *ScaLicenseProfileAssignmentDeleteOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectKey, "ProjectKey")
}

// ValidateLicenseProfileGetOpt validates the options for the GetLicenseProfile and
// DeleteLicenseProfile methods.
func (s *ScaService) ValidateLicenseProfileGetOpt(opt *ScaLicenseProfileGetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateLicenseProfileCategoryOpt validates the options for the
// UpdateLicenseProfileCategory method.
func (s *ScaService) ValidateLicenseProfileCategoryOpt(opt *ScaLicenseProfileCategoryOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.LicenseProfileKey, "LicenseProfileKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.CategoryKey, "CategoryKey")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.CategoryKey, allowedScaLicenseCategories, "CategoryKey")
}

// ValidateLicenseProfileLicenseOpt validates the options for the
// UpdateLicenseProfileLicense method.
func (s *ScaService) ValidateLicenseProfileLicenseOpt(opt *ScaLicenseProfileLicenseOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.LicenseProfileKey, "LicenseProfileKey")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.LicensePolicyId, "LicensePolicyId")
}

// ValidateLicenseProfileCreateOpt validates the request body for the CreateLicenseProfile method.
func (s *ScaService) ValidateLicenseProfileCreateOpt(body *ScaLicenseProfileCreateRequest) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(body.Name, "Name")
	if err != nil {
		return err
	}

	return ValidateRequired(body.Organization, "Organization")
}

// ValidateLicenseProfileAssignmentOpt validates the request body for the
// AssignLicenseProfileProject method.
func (s *ScaService) ValidateLicenseProfileAssignmentOpt(body *ScaLicenseProfileAssignmentRequest) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(body.ProjectKey, "ProjectKey")
}

// ValidateLicenseProfileCategoryUpdateOpt validates the options and request body for the
// UpdateLicenseProfileCategory method.
func (s *ScaService) ValidateLicenseProfileCategoryUpdateOpt(
	opt *ScaLicenseProfileCategoryOptions, body *ScaLicenseProfileCategoryUpdateRequest,
) error {
	err := s.ValidateLicenseProfileCategoryOpt(opt)
	if err != nil {
		return err
	}

	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err = ValidateRequired(body.Policy, "Policy")
	if err != nil {
		return err
	}

	return IsValueAuthorized(body.Policy, allowedScaLicensePolicies, "Policy")
}

// ValidateLicenseProfileLicenseUpdateOpt validates the options and request body for the
// UpdateLicenseProfileLicense method.
func (s *ScaService) ValidateLicenseProfileLicenseUpdateOpt(
	opt *ScaLicenseProfileLicenseOptions, body *ScaLicensePolicyLicenseUpdateRequest,
) error {
	err := s.ValidateLicenseProfileLicenseOpt(opt)
	if err != nil {
		return err
	}

	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err = ValidateRequired(body.Policy, "Policy")
	if err != nil {
		return err
	}

	return IsValueAuthorized(body.Policy, allowedScaLicensePolicies, "Policy")
}

// -----------------------------------------------------------------------------
// License Profiles - Service Methods
// -----------------------------------------------------------------------------

// ListLicenseProfiles lists the license profiles that have been configured for use when evaluating
// license issues.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/license-profiles.
// Enterprise Edition only.
func (s *ScaService) ListLicenseProfiles(ctx context.Context, opt *ScaLicenseProfileListOptions) (*ScaLicenseProfileIndex, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/license-profiles", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfileIndex)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateLicenseProfile creates a new license profile for projects to use when evaluating license
// issues.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: POST /api/v2/sca/license-profiles.
// Enterprise Edition only.
func (s *ScaService) CreateLicenseProfile(ctx context.Context, body *ScaLicenseProfileCreateRequest) (*ScaLicenseProfile, *http.Response, error) {
	err := s.ValidateLicenseProfileCreateOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "sca/license-profiles", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfile)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListLicenseProfileAssignableProjects lists the projects that can be assigned to a license
// profile. Assigning a project to a license profile will cause that license profile to be used
// when analyzing the project for license issues.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/license-profiles/assignable-projects.
// Enterprise Edition only.
func (s *ScaService) ListLicenseProfileAssignableProjects(ctx context.Context, opt *ScaLicenseProfileAssignableProjectsOptions) (*ScaLicenseProfileAssignableProjectsResponse, *http.Response, error) {
	if opt != nil {
		err := opt.Validate()
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/license-profiles/assignable-projects", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfileAssignableProjectsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AssignLicenseProfileProject configures which license profile should be used when analyzing a
// project for license issues.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/license-profiles/assigned-projects.
// Enterprise Edition only.
func (s *ScaService) AssignLicenseProfileProject(ctx context.Context, body *ScaLicenseProfileAssignmentRequest) (*ScaLicenseProfile, *http.Response, error) {
	err := s.ValidateLicenseProfileAssignmentOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/license-profiles/assigned-projects", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfile)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteLicenseProfileAssignment removes the project's assignment to the license profile. If there
// is a default license profile, it will be used to analyze license issues for this project.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: DELETE /api/v2/sca/license-profiles/assigned-projects/{project-key}.
// Enterprise Edition only.
func (s *ScaService) DeleteLicenseProfileAssignment(ctx context.Context, opt *ScaLicenseProfileAssignmentDeleteOptions) (*http.Response, error) {
	err := s.ValidateLicenseProfileAssignmentDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "sca/license-profiles/assigned-projects/"+opt.ProjectKey, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetLicenseProfile returns the license policy for a given profile.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/license-profiles/{license-profile-key}.
// Enterprise Edition only.
func (s *ScaService) GetLicenseProfile(ctx context.Context, opt *ScaLicenseProfileGetOptions) (*ScaLicenseProfileDetails, *http.Response, error) {
	err := s.ValidateLicenseProfileGetOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/license-profiles/"+opt.Key, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfileDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteLicenseProfile deletes the license profile.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: DELETE /api/v2/sca/license-profiles/{license-profile-key}.
// Enterprise Edition only.
func (s *ScaService) DeleteLicenseProfile(ctx context.Context, opt *ScaLicenseProfileGetOptions) (*http.Response, error) {
	err := s.ValidateLicenseProfileGetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "sca/license-profiles/"+opt.Key, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateLicenseProfile updates the license profile for projects to use when evaluating license
// issues.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/license-profiles/{license-profile-key}.
// Enterprise Edition only.
func (s *ScaService) UpdateLicenseProfile(ctx context.Context, key string, body *ScaLicenseProfileUpdateRequest) (*ScaLicenseProfile, *http.Response, error) {
	err := ValidateRequired(key, "key")
	if err != nil {
		return nil, nil, err
	}

	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/license-profiles/"+key, nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfile)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateLicenseProfileCategory updates the policy for an entire category in the license profile.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/license-profiles/{license-profile-key}/categories/{category-key}.
// Enterprise Edition only.
func (s *ScaService) UpdateLicenseProfileCategory(ctx context.Context, opt *ScaLicenseProfileCategoryOptions, body *ScaLicenseProfileCategoryUpdateRequest) (*ScaLicenseProfileCategory, *http.Response, error) {
	err := s.ValidateLicenseProfileCategoryUpdateOpt(opt, body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch,
		"sca/license-profiles/"+opt.LicenseProfileKey+"/categories/"+opt.CategoryKey, nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicenseProfileCategory)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateLicenseProfileLicense updates the policy for a single license in the license profile.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/license-profiles/{license-profile-key}/licenses/{license-policy-id}.
// Enterprise Edition only.
func (s *ScaService) UpdateLicenseProfileLicense(ctx context.Context, opt *ScaLicenseProfileLicenseOptions, body *ScaLicensePolicyLicenseUpdateRequest) (*ScaLicensePolicyLicense, *http.Response, error) {
	err := s.ValidateLicenseProfileLicenseUpdateOpt(opt, body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch,
		"sca/license-profiles/"+opt.LicenseProfileKey+"/licenses/"+opt.LicensePolicyId, nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaLicensePolicyLicense)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// -----------------------------------------------------------------------------
// Issues-Releases Mutations - Response Types
// -----------------------------------------------------------------------------

// Analysis status values for ScaAnalysisResource.Status.
const (
	// ScaAnalysisStatusFailed indicates the analysis failed.
	ScaAnalysisStatusFailed = "FAILED"
	// ScaAnalysisStatusOutdated indicates the analysis is outdated.
	ScaAnalysisStatusOutdated = "OUTDATED"
	// ScaAnalysisStatusCompleted indicates the analysis completed successfully.
	ScaAnalysisStatusCompleted = "COMPLETED"
)

// ScaAnalysisError represents an error encountered while parsing dependency files during analysis.
type ScaAnalysisError struct {
	// Id is the error identifier.
	Id string `json:"id,omitempty"`
	// Code is the error code.
	Code string `json:"code,omitempty"`
	// Path is the path of the file that caused the error.
	Path string `json:"path,omitempty"`
	// Message is a human-readable description of the error.
	Message string `json:"message,omitempty"`
}

// ScaAnalysisResource represents the SCA analysis status of a project branch or pull request.
type ScaAnalysisResource struct {
	Status       string             `json:"status,omitempty"`
	FailedReason string             `json:"failedReason,omitempty"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
	Errors       []ScaAnalysisError `json:"errors,omitempty"`
	ParsedFiles  []string           `json:"parsedFiles,omitempty"`
}

// ScaFeatureEnabledResult represents whether SCA is licensed and set up for the resource's
// organization or enterprise.
type ScaFeatureEnabledResult struct {
	// Enabled indicates whether SCA is licensed and set up.
	Enabled bool `json:"enabled"`
}

// ScaUserResource represents a SonarQube user referenced by SCA issues-releases endpoints.
type ScaUserResource struct {
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// Avatar is the URL of the user's avatar.
	Avatar string `json:"avatar,omitempty"`
	// Active indicates whether the user account is active.
	Active bool `json:"active"`
}

// ScaIssueReleaseBranch describes the branch an issue-release pair belongs to.
type ScaIssueReleaseBranch struct {
	// Uuid is the branch UUID.
	Uuid string `json:"uuid,omitempty"`
	// Key is the branch key.
	Key string `json:"key,omitempty"`
	// ProjectKey is the key of the project the branch belongs to.
	ProjectKey string `json:"projectKey,omitempty"`
	// ProjectName is the name of the project the branch belongs to.
	ProjectName string `json:"projectName,omitempty"`
	// PullRequest indicates whether this branch is a pull request.
	PullRequest bool `json:"pullRequest"`
}

// ScaIssueReleaseDetails represents detailed information about an issue-release pair (dependency risk).
type ScaIssueReleaseDetails struct {
	Type                         string                   `json:"type,omitempty"`
	CreatedAt                    string                   `json:"createdAt,omitempty"`
	OriginalSeverity             string                   `json:"originalSeverity,omitempty"`
	ManualSeverity               string                   `json:"manualSeverity,omitempty"`
	SpdxLicenseId                string                   `json:"spdxLicenseId,omitempty"`
	Quality                      string                   `json:"quality,omitempty"`
	Status                       string                   `json:"status,omitempty"`
	Key                          string                   `json:"key,omitempty"`
	Severity                     string                   `json:"severity,omitempty"`
	Branch                       ScaIssueReleaseBranch    `json:"branch,omitzero"`
	Assignee                     ScaUserResource          `json:"assignee,omitzero"`
	Actions                      []string                 `json:"actions,omitempty"`
	Transitions                  []string                 `json:"transitions,omitempty"`
	Release                      ScaReleaseSearchResource `json:"release,omitzero"`
	Vulnerability                ScaVulnerabilityResource `json:"vulnerability,omitzero"`
	CommentCount                 int32                    `json:"commentCount"`
	ShowIncreasedSeverityWarning bool                     `json:"showIncreasedSeverityWarning"`
}

// ScaIssueReleaseChangeDiff represents a single field change recorded in an issue-release
// changelog entry.
type ScaIssueReleaseChangeDiff struct {
	// FieldName is the name of the field that changed.
	FieldName string `json:"fieldName,omitempty"`
	// OldValue is the value before the change.
	OldValue string `json:"oldValue,omitempty"`
	// NewValue is the value after the change.
	NewValue string `json:"newValue,omitempty"`
}

// ScaIssueReleaseChange represents a single changelog entry for an issue-release pair.
type ScaIssueReleaseChange struct {
	// Key is the unique identifier of this changelog entry.
	Key string `json:"key,omitempty"`
	// CreatedAt is the creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`
	// User is the user who made the change.
	User ScaUserResource `json:"user,omitzero"`
	// MarkdownComment is the comment text in Markdown format, if any.
	MarkdownComment string `json:"markdownComment,omitempty"`
	// HtmlComment is the comment text rendered as HTML, if any.
	HtmlComment string `json:"htmlComment,omitempty"`
	// ChangeData lists the field-level changes recorded by this entry.
	ChangeData []ScaIssueReleaseChangeDiff `json:"changeData,omitempty"`
	// Actions lists the actions available to the current user for this entry (e.g. EDIT_COMMENT,
	// DELETE_COMMENT).
	Actions []string `json:"actions,omitempty"`
}

// ScaIssueReleaseChangelog represents the response from the GetChangelog method.
type ScaIssueReleaseChangelog struct {
	// Changelog is the list of changelog entries.
	Changelog []ScaIssueReleaseChange `json:"changelog,omitempty"`
}

// ScaSelfTestHTTPCall represents the result of a single diagnostic HTTP call performed by SelfTest.
type ScaSelfTestHTTPCall struct {
	AttemptedUrl             string     `json:"attemptedUrl,omitempty"`
	AttemptedMethod          string     `json:"attemptedMethod,omitempty"`
	ResponseBody             string     `json:"responseBody,omitempty"`
	ResponseHeaders          [][]string `json:"responseHeaders,omitempty"`
	ResponseCode             int32      `json:"responseCode"`
	ResponseBodyAppearsValid bool       `json:"responseBodyAppearsValid"`
}

// ScaSelfTestResult represents the response from the SelfTest method.
type ScaSelfTestResult struct {
	// CliVersionCheck is the result of the CLI version check.
	CliVersionCheck ScaSelfTestHTTPCall `json:"cliVersionCheck,omitzero"`
	// VulnerabilityDetailsCheck is the result of the vulnerability details check.
	VulnerabilityDetailsCheck ScaSelfTestHTTPCall `json:"vulnerabilityDetailsCheck,omitzero"`
	// FeatureEnabled indicates whether SCA is enabled.
	FeatureEnabled bool `json:"featureEnabled"`
	// SelfTestPassed indicates whether the overall self-test passed.
	SelfTestPassed bool `json:"selfTestPassed"`
}

// ScaReleaseResearchIssue represents a potential issue identified for a researched release.
type ScaReleaseResearchIssue struct {
	Key                          string   `json:"key,omitempty"`
	Severity                     string   `json:"severity,omitempty"`
	Type                         string   `json:"type,omitempty"`
	Quality                      string   `json:"quality,omitempty"`
	Status                       string   `json:"status,omitempty"`
	VulnerabilityId              string   `json:"vulnerabilityId,omitempty"`
	CvssScore                    string   `json:"cvssScore,omitempty"`
	SpdxLicenseId                string   `json:"spdxLicenseId,omitempty"`
	CweIds                       []string `json:"cweIds,omitempty"`
	ShowIncreasedSeverityWarning bool     `json:"showIncreasedSeverityWarning"`
}

// ScaReleaseResearchEntry represents a single researched release.
type ScaReleaseResearchEntry struct {
	Key               string                    `json:"key,omitempty"`
	PackageUrl        string                    `json:"packageUrl,omitempty"`
	PackageManager    string                    `json:"packageManager,omitempty"`
	PackageName       string                    `json:"packageName,omitempty"`
	Version           string                    `json:"version,omitempty"`
	LicenseExpression string                    `json:"licenseExpression,omitempty"`
	Issues            []ScaReleaseResearchIssue `json:"issues,omitempty"`
	Known             bool                      `json:"known"`
	KnownPackage      bool                      `json:"knownPackage"`
	NewlyIntroduced   bool                      `json:"newlyIntroduced"`
}

// ScaReleaseResearch represents the response from the ResearchReleases method.
type ScaReleaseResearch struct {
	// Releases is the list of researched releases.
	Releases []ScaReleaseResearchEntry `json:"releases,omitempty"`
}

// ScaReanalysisResult represents the response from the TriggerReanalysis method.
type ScaReanalysisResult struct {
	// BranchesQueued is the number of branches that were queued for re-analysis.
	BranchesQueued int32 `json:"branchesQueued"`
}

// ScaLicenseInfo represents license information about a package, as returned by GetPackageInfo.
type ScaLicenseInfo struct {
	// Expression is the SPDX license expression.
	Expression string `json:"expression,omitempty"`
	// Allowed indicates whether the license is allowed under the applicable license policy.
	Allowed bool `json:"allowed"`
}

// ScaVersionFix represents a package version that fixes a known vulnerability.
type ScaVersionFix struct {
	// Version is the fixed version.
	Version string `json:"version,omitempty"`
	// FixLevel describes the level of fix (e.g. direct, transitive).
	FixLevel string `json:"fixLevel,omitempty"`
	// DescriptionCode is a code identifying the fix description.
	DescriptionCode string `json:"descriptionCode,omitempty"`
}

// ScaPackageVulnerability represents a vulnerability affecting a package, as returned by
// GetPackageInfo.
type ScaPackageVulnerability struct {
	Id                 string          `json:"id,omitempty"`
	RiskSeverity       string          `json:"riskSeverity,omitempty"`
	PublishedOn        string          `json:"publishedOn,omitempty"`
	UnaffectedVersions string          `json:"unaffectedVersions,omitempty"`
	CweIds             []string        `json:"cweIds,omitempty"`
	FixedVersions      []ScaVersionFix `json:"fixedVersions,omitempty"`
	CvssScore          float64         `json:"cvssScore"`
	Withdrawn          bool            `json:"withdrawn"`
}

// ScaPackageInfoEntry represents analyzed information about a single package, as returned by
// GetPackageInfo.
type ScaPackageInfoEntry struct {
	// PackageUrl is the package URL (purl) that was analyzed.
	PackageUrl string `json:"purl,omitempty"`
	// License contains license information about the package.
	License ScaLicenseInfo `json:"license,omitzero"`
	// Vulnerabilities lists the vulnerabilities affecting this package.
	Vulnerabilities []ScaPackageVulnerability `json:"vulnerabilities,omitempty"`
	// Malicious indicates whether the package is known to be malicious.
	Malicious bool `json:"malicious"`
	// KnownPackage indicates whether the package (regardless of version) is known.
	KnownPackage bool `json:"knownPackage"`
	// KnownRelease indicates whether this specific package version is known.
	KnownRelease bool `json:"knownRelease"`
}

// ScaPackageInfo represents the response from the GetPackageInfo method when analyzing one or
// more packages by their purls.
type ScaPackageInfo struct {
	// Packages is the list of analyzed packages.
	Packages []ScaPackageInfoEntry `json:"packages,omitempty"`
}

// ScaReleaseSearchBranch represents a branch on which a matching release was found, as returned
// by SearchReleasesByPurl.
type ScaReleaseSearchBranch struct {
	Id                  string   `json:"id,omitempty"`
	Key                 string   `json:"key,omitempty"`
	ProjectKey          string   `json:"projectKey,omitempty"`
	ProjectName         string   `json:"projectName,omitempty"`
	PackageUrl          string   `json:"packageUrl,omitempty"`
	DependencyFilePaths []string `json:"dependencyFilePaths,omitempty"`
	PullRequest         bool     `json:"pullRequest"`
	Direct              bool     `json:"direct"`
	ProductionScope     bool     `json:"productionScope"`
}

// ScaReleaseSearchByPurl represents the response from the SearchReleasesByPurl method.
type ScaReleaseSearchByPurl struct {
	// Branches is the list of branches on which a matching release was found.
	Branches []ScaReleaseSearchBranch `json:"branches,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ScaRiskReportStatusChange represents a status change recorded in a risk report item's history.
type ScaRiskReportStatusChange struct {
	// Comment explains why the status changed.
	Comment string `json:"comment,omitempty"`
	// NewStatus is the status the risk transitioned to.
	NewStatus string `json:"newStatus,omitempty"`
	// CreatedAt is the timestamp of the change.
	CreatedAt string `json:"createdAt,omitempty"`
}

// ScaRiskReportItem represents a single SCA dependency risk in a risk report.
type ScaRiskReportItem struct {
	CreatedAt          string                      `json:"createdAt,omitempty"`
	PackageUrl         string                      `json:"packageUrl,omitempty"`
	BranchKey          string                      `json:"branchKey,omitempty"`
	RiskTitle          string                      `json:"riskTitle,omitempty"`
	RiskType           string                      `json:"riskType,omitempty"`
	RiskSeverity       string                      `json:"riskSeverity,omitempty"`
	RiskStatus         string                      `json:"riskStatus,omitempty"`
	Scope              string                      `json:"scope,omitempty"`
	VulnerabilityId    string                      `json:"vulnerabilityId,omitempty"`
	RiskUrl            string                      `json:"riskUrl,omitempty"`
	ProjectName        string                      `json:"projectName,omitempty"`
	ProjectKey         string                      `json:"projectKey,omitempty"`
	PublishedOn        string                      `json:"publishedOn,omitempty"`
	CweIds             []string                    `json:"cweIds,omitempty"`
	DependencyChains   [][]string                  `json:"dependencyChains,omitempty"`
	StatusChanges      []ScaRiskReportStatusChange `json:"statusChanges,omitempty"`
	CvssScore          float64                     `json:"cvssScore"`
	EpssPercentile     float64                     `json:"epssPercentile"`
	EpssScore          float64                     `json:"epssScore"`
	KnownExploited     bool                        `json:"knownExploited"`
	ProductionScope    bool                        `json:"productionScope"`
	ConfirmedReachable bool                        `json:"confirmedReachable"`
}

// -----------------------------------------------------------------------------
// Issues-Releases Mutations - Option Types
// -----------------------------------------------------------------------------

// ScaAnalysisGetOptions contains parameters for the GetAnalysis method.
type ScaAnalysisGetOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey filters by branch. Optional; if not provided, the default branch is used unless
	// PullRequestKey is provided.
	BranchKey string `json:"branchKey,omitempty"`
	// PullRequestKey filters by pull request. Optional.
	PullRequestKey string `json:"pullRequestKey,omitempty"`
}

// ScaAllAssigneesOptions contains parameters for the GetAllAssignees method.
type ScaAllAssigneesOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey filters by branch. Optional; if not provided, the default branch is used unless
	// PullRequestKey is provided.
	BranchKey string `json:"branchKey,omitempty"`
	// PullRequestKey filters by pull request. Optional.
	PullRequestKey string `json:"pullRequestKey,omitempty"`
}

// Transition key values accepted by ScaChangeStatusRequest.TransitionKey.
const (
	// ScaTransitionConfirm confirms the risk.
	ScaTransitionConfirm = "CONFIRM"
	// ScaTransitionReopen reopens the risk.
	ScaTransitionReopen = "REOPEN"
	// ScaTransitionSafe marks the risk as safe.
	ScaTransitionSafe = "SAFE"
	// ScaTransitionFixed marks the risk as fixed.
	ScaTransitionFixed = "FIXED"
	// ScaTransitionAccept accepts the risk.
	ScaTransitionAccept = "ACCEPT"
)

//nolint:gochecknoglobals // constant set of allowed values
var allowedScaTransitions = map[string]struct{}{
	ScaTransitionConfirm: {},
	ScaTransitionReopen:  {},
	ScaTransitionSafe:    {},
	ScaTransitionFixed:   {},
	ScaTransitionAccept:  {},
}

// ScaUpdateAssigneeRequest contains the request body for the UpdateAssignee method.
type ScaUpdateAssigneeRequest struct {
	// IssueReleaseKey is the issue-release key. This field is required.
	IssueReleaseKey string `json:"issueReleaseKey"`
	// AssigneeLogin is the login of the user to assign. Optional; omit to unassign.
	AssigneeLogin string `json:"assigneeLogin,omitempty"`
}

// ScaChangeStatusRequest contains the request body for the ChangeStatus method.
type ScaChangeStatusRequest struct {
	// IssueReleaseKey is the issue-release key. This field is required.
	IssueReleaseKey string `json:"issueReleaseKey"`
	// TransitionKey is the transition to apply. This field is required.
	// Allowed values: ScaTransitionConfirm, ScaTransitionReopen, ScaTransitionSafe,
	// ScaTransitionFixed, ScaTransitionAccept.
	TransitionKey string `json:"transitionKey"`
	// Comment explains why the status is changing. Optional.
	Comment string `json:"comment,omitempty"`
}

// ScaSetSeverityRequest contains the request body for the SetSeverity and
// UpdateIssueReleaseSeverity methods.
type ScaSetSeverityRequest struct {
	// IssueReleaseKey is the issue-release key. Required for SetSeverity; ignored by
	// UpdateIssueReleaseSeverity, which identifies the issue-release pair via its "id" parameter
	// instead.
	IssueReleaseKey string `json:"issueReleaseKey,omitempty"`
	// Quality is the software quality to set the severity for. Optional.
	// Allowed values: SoftwareQualityMaintainability, SoftwareQualityReliability,
	// SoftwareQualitySecurity.
	Quality string `json:"quality,omitempty"`
	// Severity is the manual severity to set. Optional.
	// Allowed values: RuleImpactSeverityInfo, RuleImpactSeverityLow, RuleImpactSeverityMedium,
	// RuleImpactSeverityHigh, RuleImpactSeverityBlocker.
	Severity string `json:"severity,omitempty"`
}

// ScaAddCommentRequest contains the request body for the AddComment method.
type ScaAddCommentRequest struct {
	// IssueReleaseKey is the issue-release key. This field is required.
	IssueReleaseKey string `json:"issueReleaseKey"`
	// Comment is the comment text. Optional.
	Comment string `json:"comment,omitempty"`
}

// ScaChangelogGetOptions contains parameters for the GetChangelog method.
type ScaChangelogGetOptions struct {
	// Key is the issue-release key. This field is required.
	Key string `json:"key"`
}

// ScaChangelogDeleteOptions contains parameters for the DeleteChangelogEntry method.
type ScaChangelogDeleteOptions struct {
	// Key is the issue-release key. This field is required. Sent as a path parameter.
	Key string `json:"-"`
	// IssueReleaseChangeKey is the changelog entry key to delete. This field is required.
	IssueReleaseChangeKey string `json:"issueReleaseChangeKey"`
}

// ScaChangelogUpdateOptions contains parameters for the UpdateChangelogComment method.
type ScaChangelogUpdateOptions struct {
	// Key is the issue-release key. This field is required.
	Key string `json:"key"`
}

// ScaChangelogUpdateRequest contains the request body for the UpdateChangelogComment method.
type ScaChangelogUpdateRequest struct {
	// IssueReleaseChangeKey is the changelog entry key to update. This field is required.
	IssueReleaseChangeKey string `json:"issueReleaseChangeKey"`
	// Comment is the new comment text. This field is required.
	Comment string `json:"comment"`
}

// ScaReleaseResearchOptions contains the request body for the ResearchReleases method.
type ScaReleaseResearchOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// Purls is the list of versioned package URLs to research. This field is required. Maximum
	// 100 per request.
	Purls []string `json:"purls"`
}

// ScaReanalysisOptions contains the request body for the TriggerReanalysis method.
type ScaReanalysisOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey is the branch to re-analyze. Optional; if omitted, all branches matching the
	// organization's rescan branch type setting will be queued.
	BranchKey string `json:"branchKey,omitempty"`
}

// ScaPackageInfoOptions contains the request body for the GetPackageInfo method.
type ScaPackageInfoOptions struct {
	IncludeVulnerabilityDetails *bool    `json:"includeVulnerabilityDetails,omitempty"`
	ProjectKey                  string   `json:"projectKey"`
	Purls                       []string `json:"purls"`
}

// ScaReleaseSearchByPurlOptions contains parameters for the SearchReleasesByPurl method.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type ScaReleaseSearchByPurlOptions struct {
	PaginationParamsV2

	Purl string `json:"purl"`
}

// ScaRiskReportOptions contains parameters for the GetRiskReport method.
type ScaRiskReportOptions struct {
	// Component is the key of the component (project, application, portfolio) to build the
	// report for. This field is required.
	Component string `json:"component"`
	// Branch filters by branch. Optional.
	Branch string `json:"branch,omitempty"`
	// RiskType filters by risk type. Optional; if not provided, all risk types are included.
	RiskType []string `json:"riskType,omitempty"`
}

// -----------------------------------------------------------------------------
// Issues-Releases Mutations - Validation Functions
// -----------------------------------------------------------------------------

// ValidateAnalysisGetOpt validates the options for the GetAnalysis method.
func (s *ScaService) ValidateAnalysisGetOpt(opt *ScaAnalysisGetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectKey, "ProjectKey")
}

// ValidateAllAssigneesOpt validates the options for the GetAllAssignees method.
func (s *ScaService) ValidateAllAssigneesOpt(opt *ScaAllAssigneesOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectKey, "ProjectKey")
}

// ValidateUpdateAssigneeOpt validates the request body for the UpdateAssignee method.
func (s *ScaService) ValidateUpdateAssigneeOpt(body *ScaUpdateAssigneeRequest) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(body.IssueReleaseKey, "IssueReleaseKey")
}

// ValidateChangeStatusOpt validates the request body for the ChangeStatus method.
func (s *ScaService) ValidateChangeStatusOpt(body *ScaChangeStatusRequest) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(body.IssueReleaseKey, "IssueReleaseKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(body.TransitionKey, "TransitionKey")
	if err != nil {
		return err
	}

	return IsValueAuthorized(body.TransitionKey, allowedScaTransitions, "TransitionKey")
}

// validateScaSetSeverityFields validates the optional Quality/Severity enum fields shared by the
// SetSeverity and UpdateIssueReleaseSeverity request bodies.
func validateScaSetSeverityFields(body *ScaSetSeverityRequest) error {
	if body.Quality != "" {
		err := IsValueAuthorized(body.Quality, allowedImpactSoftwareQualities, "Quality")
		if err != nil {
			return err
		}
	}

	if body.Severity != "" {
		return IsValueAuthorized(body.Severity, allowedRuleImpactSeverities, "Severity")
	}

	return nil
}

// ValidateSetSeverityOpt validates the request body for the SetSeverity method.
func (s *ScaService) ValidateSetSeverityOpt(body *ScaSetSeverityRequest) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(body.IssueReleaseKey, "IssueReleaseKey")
	if err != nil {
		return err
	}

	return validateScaSetSeverityFields(body)
}

// ValidateUpdateIssueReleaseSeverityOpt validates the key and request body for the
// UpdateIssueReleaseSeverity method.
func (s *ScaService) ValidateUpdateIssueReleaseSeverityOpt(key string, body *ScaSetSeverityRequest) error {
	err := ValidateRequired(key, "key")
	if err != nil {
		return err
	}

	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	return validateScaSetSeverityFields(body)
}

// ValidateAddCommentOpt validates the request body for the AddComment method.
func (s *ScaService) ValidateAddCommentOpt(body *ScaAddCommentRequest) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(body.IssueReleaseKey, "IssueReleaseKey")
}

// ValidateChangelogGetOpt validates the options for the GetChangelog method.
func (s *ScaService) ValidateChangelogGetOpt(opt *ScaChangelogGetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateChangelogDeleteOpt validates the options for the DeleteChangelogEntry method.
func (s *ScaService) ValidateChangelogDeleteOpt(opt *ScaChangelogDeleteOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.IssueReleaseChangeKey, "IssueReleaseChangeKey")
}

// ValidateChangelogUpdateOpt validates the options and request body for the
// UpdateChangelogComment method.
func (s *ScaService) ValidateChangelogUpdateOpt(opt *ScaChangelogUpdateOptions, body *ScaChangelogUpdateRequest) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err = ValidateRequired(body.IssueReleaseChangeKey, "IssueReleaseChangeKey")
	if err != nil {
		return err
	}

	return ValidateRequired(body.Comment, "Comment")
}

// ValidateReleaseResearchOpt validates the request body for the ResearchReleases method.
func (s *ScaService) ValidateReleaseResearchOpt(body *ScaReleaseResearchOptions) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(body.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	if len(body.Purls) == 0 {
		return NewValidationError("Purls", "must not be empty", ErrMissingRequired)
	}

	return nil
}

// ValidateReanalysisOpt validates the request body for the TriggerReanalysis method.
func (s *ScaService) ValidateReanalysisOpt(body *ScaReanalysisOptions) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(body.ProjectKey, "ProjectKey")
}

// ValidatePackageInfoOpt validates the request body for the GetPackageInfo method.
func (s *ScaService) ValidatePackageInfoOpt(body *ScaPackageInfoOptions) error {
	if body == nil {
		return NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(body.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	if len(body.Purls) == 0 {
		return NewValidationError("Purls", "must not be empty", ErrMissingRequired)
	}

	return nil
}

// ValidateReleaseSearchByPurlOpt validates the options for the SearchReleasesByPurl method.
func (s *ScaService) ValidateReleaseSearchByPurlOpt(opt *ScaReleaseSearchByPurlOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Purl, "Purl")
	if err != nil {
		return err
	}

	return opt.Validate()
}

// ValidateRiskReportOpt validates the options for the GetRiskReport method.
func (s *ScaService) ValidateRiskReportOpt(opt *ScaRiskReportOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Component, "Component")
}

// -----------------------------------------------------------------------------
// Issues-Releases Mutations - Service Methods
// -----------------------------------------------------------------------------

// GetAnalysis fetches the analysis status for a project branch or pull request.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/analyses.
// Enterprise Edition only.
func (s *ScaService) GetAnalysis(ctx context.Context, opt *ScaAnalysisGetOptions) (*ScaAnalysisResource, *http.Response, error) {
	err := s.ValidateAnalysisGetOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/analyses", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaAnalysisResource)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IsScaEnabled returns whether SCA is licensed and set up for the resource's organization or
// enterprise.
//
// Deprecated: prefer GetFeatureEnabled in new code.
//
// API endpoint: GET /api/v2/sca/enabled.
// Enterprise Edition only.
func (s *ScaService) IsScaEnabled(ctx context.Context) (*ScaFeatureEnabledResult, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/enabled", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaFeatureEnabledResult)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetFeatureEnabled returns whether SCA is licensed and set up for the resource's organization or
// enterprise.
//
// API endpoint: GET /api/v2/sca/feature-enabled.
// Enterprise Edition only.
func (s *ScaService) GetFeatureEnabled(ctx context.Context) (*ScaFeatureEnabledResult, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/feature-enabled", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaFeatureEnabledResult)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SelfTest runs an SCA self-test, checking connectivity to backing services.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/self-test.
// Enterprise Edition only.
func (s *ScaService) SelfTest(ctx context.Context) (*ScaSelfTestResult, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/self-test", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaSelfTestResult)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetAllAssignees returns the list of users assigned to at least one issue in the given branch.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/issues-releases/all-assignees.
// Enterprise Edition only.
func (s *ScaService) GetAllAssignees(ctx context.Context, opt *ScaAllAssigneesOptions) ([]ScaUserResource, *http.Response, error) {
	err := s.ValidateAllAssigneesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/issues-releases/all-assignees", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ScaUserResource

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateAssignee updates the assignee of an issue-release pair (dependency risk).
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/issues-releases/update-assignee.
// Enterprise Edition only.
func (s *ScaService) UpdateAssignee(ctx context.Context, body *ScaUpdateAssigneeRequest) (*ScaIssueReleaseDetails, *http.Response, error) {
	err := s.ValidateUpdateAssigneeOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/issues-releases/update-assignee", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaIssueReleaseDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ChangeStatus transitions a single issue-release pair (dependency risk) to a new status, with an
// optional comment explaining why.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/issues-releases/change-status.
// Enterprise Edition only.
func (s *ScaService) ChangeStatus(ctx context.Context, body *ScaChangeStatusRequest) (*ScaIssueReleaseDetails, *http.Response, error) {
	err := s.ValidateChangeStatusOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/issues-releases/change-status", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaIssueReleaseDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetSeverity manually changes the severity of an issue-release pair.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: POST /api/v2/sca/issues-releases/set-severity.
// Enterprise Edition only.
func (s *ScaService) SetSeverity(ctx context.Context, body *ScaSetSeverityRequest) (*ScaIssueReleaseDetails, *http.Response, error) {
	err := s.ValidateSetSeverityOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "sca/issues-releases/set-severity", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaIssueReleaseDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateIssueReleaseSeverity updates an issue-release pair's information, specifically its
// severity, identifying the pair via the given id rather than a request body field.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/issues-releases/{id}.
// Enterprise Edition only.
func (s *ScaService) UpdateIssueReleaseSeverity(ctx context.Context, key string, body *ScaSetSeverityRequest) (*ScaIssueReleaseDetails, *http.Response, error) {
	err := s.ValidateUpdateIssueReleaseSeverityOpt(key, body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/issues-releases/"+key, nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaIssueReleaseDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AddComment adds a comment to an issue-release pair.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: POST /api/v2/sca/issues-releases/comments.
// Enterprise Edition only.
func (s *ScaService) AddComment(ctx context.Context, body *ScaAddCommentRequest) (*http.Response, error) {
	err := s.ValidateAddCommentOpt(body)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "sca/issues-releases/comments", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetChangelog returns the changelog for a single issue-release pair.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/issues-releases/{key}/changelogs.
// Enterprise Edition only.
func (s *ScaService) GetChangelog(ctx context.Context, opt *ScaChangelogGetOptions) (*ScaIssueReleaseChangelog, *http.Response, error) {
	err := s.ValidateChangelogGetOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/issues-releases/"+opt.Key+"/changelogs", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaIssueReleaseChangelog)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteChangelogEntry deletes a comment from an issue-release pair's changelog.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: DELETE /api/v2/sca/issues-releases/{key}/changelog.
// Enterprise Edition only.
func (s *ScaService) DeleteChangelogEntry(ctx context.Context, opt *ScaChangelogDeleteOptions) (*http.Response, error) {
	err := s.ValidateChangelogDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "sca/issues-releases/"+opt.Key+"/changelog", opt, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateChangelogComment updates a comment on an issue-release pair's changelog.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: PATCH /api/v2/sca/issues-releases/{key}/changelog.
// Enterprise Edition only.
func (s *ScaService) UpdateChangelogComment(ctx context.Context, opt *ScaChangelogUpdateOptions, body *ScaChangelogUpdateRequest) (*http.Response, error) {
	err := s.ValidateChangelogUpdateOpt(opt, body)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "sca/issues-releases/"+opt.Key+"/changelog", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ResearchReleases researches release information for a set of versioned package URLs, without
// requiring them to already be present in a project's dependency graph.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: POST /api/v2/sca/release-research/releases.
// Enterprise Edition only.
func (s *ScaService) ResearchReleases(ctx context.Context, body *ScaReleaseResearchOptions) (*ScaReleaseResearch, *http.Response, error) {
	err := s.ValidateReleaseResearchOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "sca/release-research/releases", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaReleaseResearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// TriggerReanalysis triggers an immediate re-analysis of a project or branch.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: POST /api/v2/sca/reanalysis.
// Enterprise Edition only.
func (s *ScaService) TriggerReanalysis(ctx context.Context, body *ScaReanalysisOptions) (*ScaReanalysisResult, *http.Response, error) {
	err := s.ValidateReanalysisOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "sca/reanalysis", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaReanalysisResult)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPackageInfo analyzes one or more packages by their package URLs (purls), returning license,
// vulnerability and malware information.
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: POST /api/v2/sca/package-info.
// Enterprise Edition only.
func (s *ScaService) GetPackageInfo(ctx context.Context, body *ScaPackageInfoOptions) (*ScaPackageInfo, *http.Response, error) {
	err := s.ValidatePackageInfoOpt(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "sca/package-info", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaPackageInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchReleasesByPurl searches, across all projects visible to the current user, for releases
// matching the given package URL (purl).
// Accepts only authenticated requests. This is an internal API and subject to change without notice.
//
// API endpoint: GET /api/v2/sca/releases/search.
// Enterprise Edition only.
func (s *ScaService) SearchReleasesByPurl(ctx context.Context, opt *ScaReleaseSearchByPurlOptions) (*ScaReleaseSearchByPurl, *http.Response, error) {
	err := s.ValidateReleaseSearchByPurlOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/releases/search", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaReleaseSearchByPurl)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetRiskReport returns a report of all current SCA dependency risks for a given component and
// branch.
//
// API endpoint: GET /api/v2/sca/risk-reports.
// Enterprise Edition only.
func (s *ScaService) GetRiskReport(ctx context.Context, opt *ScaRiskReportOptions) ([]ScaRiskReportItem, *http.Response, error) {
	err := s.ValidateRiskReportOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/risk-reports", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ScaRiskReportItem

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
