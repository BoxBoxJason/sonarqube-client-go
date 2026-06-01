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
	// NewInPullRequest indicates whether this dependency is new in the current PR.
	NewInPullRequest bool `json:"newInPullRequest"`
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
	// NewInPullRequest indicates whether this is new in the current PR.
	NewInPullRequest bool `json:"newInPullRequest"`
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
	// NewInPullRequest filters for new-in-PR risks. Optional.
	NewInPullRequest *bool `json:"newInPullRequest,omitempty"`
	// Direct filters for direct dependencies only. Optional.
	Direct *bool `json:"direct,omitempty"`
	// ProductionScope filters by production scope. Optional.
	ProductionScope *bool `json:"productionScope,omitempty"`
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
	// NewInPullRequest filters for new-in-PR releases. Optional.
	NewInPullRequest *bool `json:"newInPullRequest,omitempty"`
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

	// ScaSbomReportFormatJSON requests the report serialized as JSON.
	ScaSbomReportFormatJSON = "json"
	// ScaSbomReportFormatXML requests the report serialized as XML.
	ScaSbomReportFormatXML = "xml"
)

// ScaSbomReportOptions contains parameters for the GetSbomReport method.
type ScaSbomReportOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey filters by branch. Optional.
	BranchKey string `json:"branchKey,omitempty"`
	// Type is the SBOM report type, one of ScaSbomReportTypeCycloneDX or ScaSbomReportTypeSpdx23.
	// This field is required.
	Type string `json:"type"`
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
	case ScaSbomReportTypeSpdx23:
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

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
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
