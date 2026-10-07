package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// OnboardingService handles communication with the onboarding related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type OnboardingService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedOnboardingGetProjectsScanStatus is the set of allowed values for the corresponding option.
	allowedOnboardingGetProjectsScanStatus = map[string]struct{}{
		"SCANNED":     {},
		"NOT_SCANNED": {},
	}
	// allowedOnboardingGetProjectsAnalysisMode is the set of allowed values for the corresponding option.
	allowedOnboardingGetProjectsAnalysisMode = map[string]struct{}{
		"CI":        {},
		"AUTOMATIC": {},
		"NONE":      {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// OnboardingGetOverviewOptions contains parameters for the GetOverview method.
type OnboardingGetOverviewOptions struct {
	// OrganizationKey is the organization key.
	OrganizationKey string `json:"organizationKey,omitempty"`
}

// OnboardingDevopsPlatformsStep represents the DevopsPlatformsStep object of the SonarQube V2 API.
type OnboardingDevopsPlatformsStep struct {
	// Configured is the configured.
	Configured int64 `json:"configured,omitempty"`
}

// OnboardingRepositoriesStep represents the RepositoriesStep object of the SonarQube V2 API.
type OnboardingRepositoriesStep struct {
	// Imported is the imported.
	Imported int64 `json:"imported,omitempty"`
	// Discovered is the discovered.
	Discovered int64 `json:"discovered,omitempty"`
	// Percent is the percent.
	Percent float64 `json:"percent,omitempty"`
}

// OnboardingProjectsStep represents the ProjectsStep object of the SonarQube V2 API.
type OnboardingProjectsStep struct {
	// Analyzed is the analyzed.
	Analyzed int64 `json:"analyzed,omitempty"`
	// NotScanned is the not scanned.
	NotScanned int64 `json:"notScanned,omitempty"`
	// NotImported is the not imported.
	NotImported int64 `json:"notImported,omitempty"`
	// Total is the total.
	Total int64 `json:"total,omitempty"`
	// Percent is the percent.
	Percent float64 `json:"percent,omitempty"`
}

// OnboardingOnboardingSteps represents the OnboardingSteps object of the SonarQube V2 API.
type OnboardingOnboardingSteps struct {
	// DevopsPlatforms is the devops platforms.
	DevopsPlatforms OnboardingDevopsPlatformsStep `json:"devopsPlatforms,omitzero"`
	// Repositories is the repositories.
	Repositories OnboardingRepositoriesStep `json:"repositories,omitzero"`
	// Projects is the projects.
	Projects OnboardingProjectsStep `json:"projects,omitzero"`
}

// OnboardingOnboardingOverviewResponse represents the OnboardingOverviewResponse object of the
// SonarQube V2 API.
type OnboardingOnboardingOverviewResponse struct {
	// ProgressPct is the progress pct.
	ProgressPct float64 `json:"progressPct,omitempty"`
	// Steps is the steps.
	Steps OnboardingOnboardingSteps `json:"steps,omitzero"`
}

// OnboardingGetProjectsOptions contains parameters for the GetProjects method.
type OnboardingGetProjectsOptions struct {
	// OrganizationKey is the organization key.
	OrganizationKey string `json:"organizationKey,omitempty"`
	// Q is the q.
	Q string `json:"q,omitempty"`
	// ScanStatus is the scan status. Allowed values: SCANNED, NOT_SCANNED.
	ScanStatus string `json:"scanStatus,omitempty"`
	// AnalysisMode is the analysis mode. Allowed values: CI, AUTOMATIC, NONE.
	AnalysisMode string `json:"analysisMode,omitempty"`
	// PageIndex is the page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the page size.
	// Must be between 1 and 500. Default: 100.
	PageSize int32 `json:"pageSize,omitempty"`
}

// OnboardingOnboardingProjectRow represents the OnboardingProjectRow object of the SonarQube V2
// API.
type OnboardingOnboardingProjectRow struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Path is the path.
	Path string `json:"path,omitempty"`
	// Alm is the alm.
	Alm string `json:"alm,omitempty"`
	// ScanStatus is the scan status. Allowed values: SCANNED, NOT_SCANNED.
	ScanStatus string `json:"scanStatus,omitempty"`
	// AnalysisMode is the analysis mode. Allowed values: CI, AUTOMATIC, NONE.
	AnalysisMode string `json:"analysisMode,omitempty"`
	// LastScan is the last scan.
	LastScan int64 `json:"lastScan,omitempty"`
}

// OnboardingOnboardingProjectsResponse represents the OnboardingProjectsResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type OnboardingOnboardingProjectsResponse struct {
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
	// Projects is the projects.
	Projects []OnboardingOnboardingProjectRow `json:"projects,omitempty"`
}

// OnboardingGetStatisticsOptions contains parameters for the GetStatistics method.
type OnboardingGetStatisticsOptions struct {
	// OrganizationKey is the organization key.
	OrganizationKey string `json:"organizationKey,omitempty"`
	// StartDate is the start date.
	StartDate string `json:"startDate,omitempty"`
	// EndDate is the end date.
	EndDate string `json:"endDate,omitempty"`
}

// OnboardingTimelinePoint represents the TimelinePoint object of the SonarQube V2 API.
type OnboardingTimelinePoint struct {
	// Date is the date.
	Date string `json:"date,omitempty"`
	// RepositoriesImported is the repositories imported.
	RepositoriesImported int64 `json:"repositoriesImported,omitempty"`
	// ProjectsScanned is the projects scanned.
	ProjectsScanned int64 `json:"projectsScanned,omitempty"`
}

// OnboardingPlatformShare represents the PlatformShare object of the SonarQube V2 API.
type OnboardingPlatformShare struct {
	// Platform is the platform.
	Platform string `json:"platform,omitempty"`
	// Count is the count.
	Count int64 `json:"count,omitempty"`
	// Percentage is the percentage.
	Percentage float64 `json:"percentage,omitempty"`
}

// OnboardingDevopsPlatforms represents the DevopsPlatforms object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type OnboardingDevopsPlatforms struct {
	// Total is the total.
	Total int64 `json:"total,omitempty"`
	// Shares is the shares.
	Shares []OnboardingPlatformShare `json:"shares,omitempty"`
}

// OnboardingOnboardingStatisticsResponse represents the OnboardingStatisticsResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type OnboardingOnboardingStatisticsResponse struct {
	// DiscoveredTotal is the discovered total.
	DiscoveredTotal int64 `json:"discoveredTotal,omitempty"`
	// Timeline is the timeline.
	Timeline []OnboardingTimelinePoint `json:"timeline,omitempty"`
	// DevopsPlatforms is the devops platforms.
	DevopsPlatforms OnboardingDevopsPlatforms `json:"devopsPlatforms,omitzero"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetProjectsOpt validates the options for the GetProjects method.
func (s *OnboardingService) ValidateGetProjectsOpt(opt *OnboardingGetProjectsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.ScanStatus, allowedOnboardingGetProjectsScanStatus, "ScanStatus")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.AnalysisMode, allowedOnboardingGetProjectsAnalysisMode, "AnalysisMode")
	if err != nil {
		return err
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	if opt.PageSize < 0 || opt.PageSize > 500 {
		return NewValidationError("PageSize", "must be between 0 and 500", ErrOutOfRange)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetOverview calls GET /api/v2/onboarding/overview.
//
// API endpoint: GET /api/v2/onboarding/overview.
// Enterprise Edition only.
func (s *OnboardingService) GetOverview(ctx context.Context, opt *OnboardingGetOverviewOptions) (*OnboardingOnboardingOverviewResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "onboarding/overview", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(OnboardingOnboardingOverviewResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjects calls GET /api/v2/onboarding/projects.
//
// API endpoint: GET /api/v2/onboarding/projects.
// Enterprise Edition only.
func (s *OnboardingService) GetProjects(ctx context.Context, opt *OnboardingGetProjectsOptions) (*OnboardingOnboardingProjectsResponse, *http.Response, error) {
	err := s.ValidateGetProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "onboarding/projects", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(OnboardingOnboardingProjectsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetStatistics calls GET /api/v2/onboarding/statistics.
//
// API endpoint: GET /api/v2/onboarding/statistics.
// Enterprise Edition only.
func (s *OnboardingService) GetStatistics(ctx context.Context, opt *OnboardingGetStatisticsOptions) (*OnboardingOnboardingStatisticsResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "onboarding/statistics", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(OnboardingOnboardingStatisticsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
