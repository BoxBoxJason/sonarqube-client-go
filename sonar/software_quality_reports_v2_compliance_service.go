package sonar

import (
	"context"
	"fmt"
	"net/http"
)

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedSoftwareQualityReportsGetComplianceReportStandard is the set of allowed values for the corresponding option.
	allowedSoftwareQualityReportsGetComplianceReportStandard = map[string]struct{}{
		"misra": {},
		"wcag":  {},
	}
	// allowedSoftwareQualityReportsGetComplianceReportLevel is the set of allowed values for the corresponding option.
	allowedSoftwareQualityReportsGetComplianceReportLevel = map[string]struct{}{
		"A":   {},
		"AA":  {},
		"AAA": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// SoftwareQualityReportsGetComplianceReportOptions contains parameters for the GetComplianceReport method.
type SoftwareQualityReportsGetComplianceReportOptions struct {
	// BranchKey branch key. If not provided, the main branch is used for projects and
	// applications.
	BranchKey string `json:"branchKey,omitempty"`
	// ComponentKey project, application, or portfolio key. This field is required.
	ComponentKey string `json:"componentKey"`
	// Standard compliance standard. This field is required. Allowed values: misra, wcag.
	Standard string `json:"standard"`
	// Version compliance standard version. This field is required.
	Version string `json:"version"`
	// Level optional WCAG conformance level. Valid values are A, AA, and AAA and they are
	// inclusive. If omitted, all levels are included. Not supported for MISRA. Allowed values: A,
	// AA, AAA.
	Level string `json:"level,omitempty"`
}

// SoftwareQualityReportsComplianceCategory represents the ComplianceCategory object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SoftwareQualityReportsComplianceCategory struct {
	// ActiveRules is the active rules.
	ActiveRules int32 `json:"activeRules,omitempty"`
	// Issues is the issues.
	Issues int32 `json:"issues,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Categories is the list of sub-categories. The API spec omits this field, but the
	// server returns it recursively (e.g. WCAG principle -> guideline -> success criterion).
	Categories []SoftwareQualityReportsComplianceCategory `json:"categories,omitempty"`
}

// SoftwareQualityReportsGetComplianceReportResponse represents the GetComplianceReportResponse
// object of the SonarQube V2 API.
type SoftwareQualityReportsGetComplianceReportResponse struct {
	// Categories is the categories.
	Categories []SoftwareQualityReportsComplianceCategory `json:"categories,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetComplianceReportOpt validates the options for the GetComplianceReport method.
func (s *SoftwareQualityReportsService) ValidateGetComplianceReportOpt(opt *SoftwareQualityReportsGetComplianceReportOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ComponentKey, "ComponentKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Standard, "Standard")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Standard, allowedSoftwareQualityReportsGetComplianceReportStandard, "Standard")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Version, "Version")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Level, allowedSoftwareQualityReportsGetComplianceReportLevel, "Level")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetComplianceReport gets compliance report. Get the compliance report for a project, application,
// or portfolio.
//
// API endpoint: GET /api/v2/software-quality-reports/compliance-reports.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *SoftwareQualityReportsService) GetComplianceReport(ctx context.Context, opt *SoftwareQualityReportsGetComplianceReportOptions) (*SoftwareQualityReportsGetComplianceReportResponse, *http.Response, error) {
	err := s.ValidateGetComplianceReportOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "software-quality-reports/compliance-reports", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(SoftwareQualityReportsGetComplianceReportResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
