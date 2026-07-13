package sonar

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// SoftwareQualityReportsAccessibilityStandardWCAG represents the "wcag" accessibility standard.
	SoftwareQualityReportsAccessibilityStandardWCAG = "wcag"

	// SoftwareQualityReportsAccessibilityVersion20 represents accessibility standard version "2.0".
	SoftwareQualityReportsAccessibilityVersion20 = "2.0"
	// SoftwareQualityReportsAccessibilityVersion21 represents accessibility standard version "2.1".
	SoftwareQualityReportsAccessibilityVersion21 = "2.1"
	// SoftwareQualityReportsAccessibilityVersion22 represents accessibility standard version "2.2".
	SoftwareQualityReportsAccessibilityVersion22 = "2.2"
)

// SoftwareQualityReportsService handles communication with the software
// quality reports V2 API endpoints.
// This service is only available in Enterprise Edition. The underlying
// endpoint is marked internal by SonarQube (x-sonar-internal) and its
// request/response contract may change without notice between SonarQube
// versions.
type SoftwareQualityReportsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedSoftwareQualityReportsAccessibilityStandards is the set of supported accessibility standards.
	allowedSoftwareQualityReportsAccessibilityStandards = map[string]struct{}{
		SoftwareQualityReportsAccessibilityStandardWCAG: {},
	}

	// allowedSoftwareQualityReportsAccessibilityVersions is the set of supported accessibility standard versions.
	allowedSoftwareQualityReportsAccessibilityVersions = map[string]struct{}{
		SoftwareQualityReportsAccessibilityVersion20: {},
		SoftwareQualityReportsAccessibilityVersion21: {},
		SoftwareQualityReportsAccessibilityVersion22: {},
	}
)

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SoftwareQualityReportsCategory represents a single accessibility report category.
type SoftwareQualityReportsCategory struct {
	// Key is the category key.
	Key string `json:"key,omitempty"`
	// ActiveRules is the number of active rules for this category.
	ActiveRules int32 `json:"activeRules,omitempty"`
	// Issues is the number of issues found for this category.
	Issues int32 `json:"issues,omitempty"`
}

// SoftwareQualityReportsAccessibilityReport represents the accessibility
// report for a project branch.
type SoftwareQualityReportsAccessibilityReport struct {
	// Categories is the list of accessibility report categories.
	Categories []SoftwareQualityReportsCategory `json:"categories,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SoftwareQualityReportsGetAccessibilityReportOptions contains parameters for
// the GetAccessibilityReport method.
type SoftwareQualityReportsGetAccessibilityReportOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// Standard is the accessibility standard. This field is required.
	// Allowed values: wcag.
	Standard string `json:"standard"`
	// Version is the accessibility standard version. This field is required.
	// Allowed values: 2.0, 2.1, 2.2.
	Version string `json:"version"`
	// BranchKey is the branch key. If not provided, the main branch is used.
	BranchKey string `json:"branchKey,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetAccessibilityReportOpt validates the options for the
// GetAccessibilityReport method.
func (s *SoftwareQualityReportsService) ValidateGetAccessibilityReportOpt(opt *SoftwareQualityReportsGetAccessibilityReportOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Standard, "Standard")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Standard, allowedSoftwareQualityReportsAccessibilityStandards, "Standard")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Version, "Version")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Version, allowedSoftwareQualityReportsAccessibilityVersions, "Version")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetAccessibilityReport returns the accessibility report for a project branch.
//
// API endpoint: GET /api/v2/software-quality-reports/accessibility-reports.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *SoftwareQualityReportsService) GetAccessibilityReport(ctx context.Context, opt *SoftwareQualityReportsGetAccessibilityReportOptions) (*SoftwareQualityReportsAccessibilityReport, *http.Response, error) {
	err := s.ValidateGetAccessibilityReportOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "software-quality-reports/accessibility-reports", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(SoftwareQualityReportsAccessibilityReport)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
