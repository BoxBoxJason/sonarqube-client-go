package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SecurityAlertsService handles communication with the security alerts related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type SecurityAlertsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedSecurityAlertsSearchSort is the set of allowed values for the corresponding option.
	allowedSecurityAlertsSearchSort = map[string]struct{}{
		"FIRST_DETECTED_AT": {},
		"LAST_DETECTED_AT":  {},
	}
	// allowedSecurityAlertsSearchDirection is the set of allowed values for the corresponding option.
	allowedSecurityAlertsSearchDirection = map[string]struct{}{
		"ASC":  {},
		"DESC": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// SecurityAlertsSearchOptions contains parameters for the Search method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SecurityAlertsSearchOptions struct {
	// AlertTypes filter by alert type, repeatable. Allowed values: DEPENDENCY_RISK.
	AlertTypes []string `json:"alertTypes,omitempty"`
	// Statuses filter by aggregate alert status, repeatable. Allowed values: OPEN, RESOLVED.
	Statuses []string `json:"statuses,omitempty"`
	// Sort field to sort by. Allowed values: FIRST_DETECTED_AT, LAST_DETECTED_AT.
	// Default: LAST_DETECTED_AT.
	Sort string `json:"sort,omitempty"`
	// Direction sort direction. Allowed values: ASC, DESC.
	// Default: DESC.
	Direction string `json:"direction,omitempty"`
	// PageSize number of results per page. A value of 0 will only return the pagination
	// information.
	// Must be between 0 and 500. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
	// PageIndex 1-based page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
}

// SecurityAlertsCweInfo represents a SonarQube V2 API object.
//
// CWE entry with resolved name.
type SecurityAlertsCweInfo struct {
	// Code CWE identifier (e.g. CWE-89).
	Code string `json:"code,omitempty"`
	// Name human-readable CWE name, or null if not found in the catalog.
	Name string `json:"name,omitempty"`
}

// SecurityAlertsScaIssue represents a SonarQube V2 API object.
//
// SCA issue linked to the security alert.
type SecurityAlertsScaIssue struct {
	// ScaIssueUuid SCA issue UUID.
	ScaIssueUuid string `json:"scaIssueUuid,omitempty"`
	// VulnerabilityId CVE ID (e.g. CVE-2024-12345).
	VulnerabilityId string `json:"vulnerabilityId,omitempty"`
	// PackageUrl package URL (PURL).
	PackageUrl string `json:"packageUrl,omitempty"`
	// IssueType issue type (VULNERABILITY, MALWARE, PROHIBITED_LICENSE).
	IssueType string `json:"issueType,omitempty"`
	// RiskId representative SCA risk ID (highest severity among affected branches, tie-broken by
	// uuid; null if no risk exists on any affected branch; may reference a non-OPEN e.g. FIXED
	// risk).
	RiskId string `json:"riskId,omitempty"`
	// Cwes CWE codes with resolved names.
	Cwes []SecurityAlertsCweInfo `json:"cwes,omitempty"`
}

// SecurityAlertsScaRiskSummaryRestResponse represents the ScaRiskSummaryRestResponse object of the
// SonarQube V2 API.
type SecurityAlertsScaRiskSummaryRestResponse struct {
	// RiskId is the risk id.
	RiskId string `json:"riskId,omitempty"`
	// PackageName is the package name.
	PackageName string `json:"packageName,omitempty"`
	// PackageVersion is the package version.
	PackageVersion string `json:"packageVersion,omitempty"`
	// PackageEcosystem is the package ecosystem.
	PackageEcosystem string `json:"packageEcosystem,omitempty"`
}

// SecurityAlertsAffectedBranchRestResponse represents the AffectedBranchRestResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SecurityAlertsAffectedBranchRestResponse struct {
	// ProjectId is the project id.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// ProjectName is the project name.
	ProjectName string `json:"projectName,omitempty"`
	// BranchId is the branch id.
	BranchId string `json:"branchId,omitempty"`
	// BranchKey not set if the branch has been deleted.
	BranchKey string `json:"branchKey,omitempty"`
	// IsMain whether this is the project's primary branch.
	IsMain bool `json:"isMain,omitempty"`
	// Status is the status. Allowed values: OPEN, RESOLVED.
	Status string `json:"status,omitempty"`
	// FirstDetectedAt ISO-8601 UTC timestamp.
	FirstDetectedAt string `json:"firstDetectedAt,omitempty"`
	// LastDetectedAt ISO-8601 UTC timestamp.
	LastDetectedAt string `json:"lastDetectedAt,omitempty"`
	// ScaRisks SCA risk details per affected package/version.
	ScaRisks []SecurityAlertsScaRiskSummaryRestResponse `json:"scaRisks,omitempty"`
	// Severity SCA severity on this branch. Allowed values: BLOCKER, HIGH, MEDIUM, LOW, INFO.
	Severity string `json:"severity,omitempty"`
	// ReachabilityAnalyzed whether reachability has been analyzed.
	ReachabilityAnalyzed bool `json:"reachabilityAnalyzed,omitempty"`
}

// SecurityAlertsSecurityAlertRestResponse represents the SecurityAlertRestResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SecurityAlertsSecurityAlertRestResponse struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// AlertType is the alert type. Allowed values: DEPENDENCY_RISK.
	AlertType string `json:"alertType,omitempty"`
	// Status aggregate status across all affected branches: OPEN if any affected branch is OPEN.
	// Allowed values: OPEN, RESOLVED.
	Status string `json:"status,omitempty"`
	// Severity highest severity among the affected branches. Allowed values: BLOCKER, HIGH,
	// MEDIUM, LOW, INFO.
	Severity string `json:"severity,omitempty"`
	// AffectedAuthorizedBranchesTotal total number of affected branches that the current user is
	// authorized to view.
	AffectedAuthorizedBranchesTotal int32 `json:"affectedAuthorizedBranchesTotal,omitempty"`
	// FirstDetectedAt ISO-8601 UTC timestamp.
	FirstDetectedAt string `json:"firstDetectedAt,omitempty"`
	// LastDetectedAt ISO-8601 UTC timestamp.
	LastDetectedAt string `json:"lastDetectedAt,omitempty"`
	// ScaIssue linked SCA issue details (if any).
	ScaIssue SecurityAlertsScaIssue `json:"scaIssue,omitzero"`
	// AffectedBranches is the affected branches.
	AffectedBranches []SecurityAlertsAffectedBranchRestResponse `json:"affectedBranches,omitempty"`
}

// SecurityAlertsSecurityAlertSearchRestResponse represents the SecurityAlertSearchRestResponse
// object of the SonarQube V2 API.
type SecurityAlertsSecurityAlertSearchRestResponse struct {
	// SecurityAlerts is the security alerts.
	SecurityAlerts []SecurityAlertsSecurityAlertRestResponse `json:"securityAlerts,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSearchOpt validates the options for the Search method.
func (s *SecurityAlertsService) ValidateSearchOpt(opt *SecurityAlertsSearchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.Sort, allowedSecurityAlertsSearchSort, "Sort")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Direction, allowedSecurityAlertsSearchDirection, "Direction")
	if err != nil {
		return err
	}

	if opt.PageSize < 0 || opt.PageSize > 500 {
		return NewValidationError("PageSize", "must be between 0 and 500", ErrOutOfRange)
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Search searches security alerts. List security alerts, one row per alert, aggregating its
// affected branches. Requires 'Browse' permission on the affected projects; results are limited to
// alerts on projects visible to the current user.
//
// API endpoint: GET /api/v2/security-alerts/alerts.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *SecurityAlertsService) Search(ctx context.Context, opt *SecurityAlertsSearchOptions) (*SecurityAlertsSecurityAlertSearchRestResponse, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "security-alerts/alerts", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(SecurityAlertsSecurityAlertSearchRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get gets a security alert by ID. Returns a single security alert by its ID. Returns 404 if the
// alert does not exist or the current user does not have 'Browse' permission on any affected
// project.
//
// API endpoint: GET /api/v2/security-alerts/alerts/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *SecurityAlertsService) Get(ctx context.Context, securityAlertID string) (*SecurityAlertsSecurityAlertRestResponse, *http.Response, error) {
	err := ValidateRequired(securityAlertID, "securityAlertID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "security-alerts/alerts/"+url.PathEscape(securityAlertID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(SecurityAlertsSecurityAlertRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
