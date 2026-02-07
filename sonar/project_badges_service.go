package sonar

import "net/http"

// ProjectBadgesService handles communication with the project badges related methods
// of the SonarQube API.
// This service generates badges based on quality gates or measures.
type ProjectBadgesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedBadgeMetrics is the set of supported metrics for badges.
	allowedBadgeMetrics = map[string]struct{}{
		"coverage":                            {},
		"duplicated_lines_density":            {},
		"ncloc":                               {},
		"alert_status":                        {},
		"security_hotspots":                   {},
		"bugs":                                {},
		"code_smells":                         {},
		"vulnerabilities":                     {},
		"sqale_rating":                        {},
		"reliability_rating":                  {},
		"security_rating":                     {},
		"sqale_index":                         {},
		"software_quality_reliability_issues": {},
		"software_quality_maintainability_issues":             {},
		"software_quality_security_issues":                    {},
		"software_quality_maintainability_rating":             {},
		"software_quality_reliability_rating":                 {},
		"software_quality_security_rating":                    {},
		"software_quality_maintainability_remediation_effort": {},
	}
)

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectBadgesToken represents the response from retrieving a badge token.
type ProjectBadgesToken struct {
	// Token is the badge access token.
	Token string `json:"token,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectBadgesMeasureOption contains parameters for the Measure method.
type ProjectBadgesMeasureOption struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Metric is the metric key.
	// This field is required.
	// Allowed values: coverage, duplicated_lines_density, ncloc, alert_status, security_hotspots,
	// bugs, code_smells, vulnerabilities, sqale_rating, reliability_rating, security_rating,
	// sqale_index, software_quality_reliability_issues, software_quality_maintainability_issues,
	// software_quality_security_issues, software_quality_maintainability_rating,
	// software_quality_reliability_rating, software_quality_security_rating,
	// software_quality_maintainability_remediation_effort.
	Metric string `url:"metric"`
	// Project is the project or application key.
	// This field is required.
	Project string `url:"project"`
	// Token is the project badge token.
	Token string `url:"token,omitempty"`
}

// ProjectBadgesQualityGateOption contains parameters for the QualityGate method.
type ProjectBadgesQualityGateOption struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Project is the project or application key.
	// This field is required.
	Project string `url:"project"`
	// Token is the project badge token.
	Token string `url:"token,omitempty"`
}

// ProjectBadgesRenewTokenOption contains parameters for the RenewToken method.
type ProjectBadgesRenewTokenOption struct {
	// Project is the project or application key.
	// This field is required.
	Project string `url:"project"`
}

// ProjectBadgesTokenOption contains parameters for the Token method.
type ProjectBadgesTokenOption struct {
	// Project is the project or application key.
	// This field is required.
	Project string `url:"project"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateMeasureOpt validates the options for the Measure method.
func (s *ProjectBadgesService) ValidateMeasureOpt(opt *ProjectBadgesMeasureOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Metric, "Metric")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Metric, allowedBadgeMetrics, "Metric")
	if err != nil {
		return err
	}

	return nil
}

// ValidateQualityGateOpt validates the options for the QualityGate method.
func (s *ProjectBadgesService) ValidateQualityGateOpt(opt *ProjectBadgesQualityGateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRenewTokenOpt validates the options for the RenewToken method.
func (s *ProjectBadgesService) ValidateRenewTokenOpt(opt *ProjectBadgesRenewTokenOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateTokenOpt validates the options for the Token method.
func (s *ProjectBadgesService) ValidateTokenOpt(opt *ProjectBadgesTokenOption) error {
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

// Measure generates a badge for project's measure as an SVG.
// Requires 'Browse' permission on the specified project.
//
// API endpoint: GET /api/project_badges/measure.
// Since: 7.1.
func (s *ProjectBadgesService) Measure(opt *ProjectBadgesMeasureOption) (*string, *http.Response, error) {
	err := s.ValidateMeasureOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_badges/measure", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(string)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// QualityGate generates a badge for project's quality gate as an SVG.
// Requires 'Browse' permission on the specified project.
//
// API endpoint: GET /api/project_badges/quality_gate.
// Since: 7.1.
func (s *ProjectBadgesService) QualityGate(opt *ProjectBadgesQualityGateOption) (*string, *http.Response, error) {
	err := s.ValidateQualityGateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_badges/quality_gate", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(string)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RenewToken creates a new token replacing any existing token for project or application badge access.
// This token can be used to authenticate with api/project_badges/quality_gate and api/project_badges/measure endpoints.
// Requires 'Administer' permission on the specified project or application.
//
// API endpoint: POST /api/project_badges/renew_token.
// Since: 9.2.
func (s *ProjectBadgesService) RenewToken(opt *ProjectBadgesRenewTokenOption) (*http.Response, error) {
	err := s.ValidateRenewTokenOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_badges/renew_token", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Token retrieves a token to use for project or application badge access.
// This token can be used to authenticate with api/project_badges/quality_gate and api/project_badges/measure endpoints.
// Requires 'Browse' permission on the specified project or application.
//
// API endpoint: GET /api/project_badges/token.
// Since: 9.2.
func (s *ProjectBadgesService) Token(opt *ProjectBadgesTokenOption) (*ProjectBadgesToken, *http.Response, error) {
	err := s.ValidateTokenOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_badges/token", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectBadgesToken)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
