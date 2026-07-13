package sonar

import (
	"bytes"
	"context"
	"net/http"
)

// SecurityReportsService handles communication with the security reports related
// methods of the SonarQube API. This service is internal and enterprise-only.
type SecurityReportsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SecurityReportsShow represents the response from the show endpoint.
type SecurityReportsShow struct {
	// Categories is the list of security report categories.
	Categories []SecurityReportCategory `json:"categories,omitempty"`
}

// SecurityReportCategory represents a category in a security report.
type SecurityReportCategory struct {
	// Category is the category name.
	Category string `json:"category,omitempty"`
	// Version is the standard revision this category belongs to (e.g. "2017", "2021"
	// for owaspTop10, "3.2", "4.0" for pciDss). Empty for standards without revisions.
	Version string `json:"version,omitempty"`
	// Distribution is the CWE breakdown for this category. Only populated when
	// includeDistribution is requested (and non-empty for CWE-oriented standards).
	Distribution []SecurityReportCWEDistribution `json:"distribution,omitempty"`
	// Vulnerabilities is the number of vulnerabilities in this category.
	Vulnerabilities int `json:"vulnerabilities,omitempty"`
	// SecurityReviewRating is the security review rating (1-5, A-E) for this category.
	SecurityReviewRating int `json:"securityReviewRating,omitempty"`
	// ActiveRules is the number of active rules for this category.
	ActiveRules int `json:"activeRules,omitempty"`
	// TotalRules is the total number of rules (active or not) for this category.
	TotalRules int `json:"totalRules,omitempty"`
	// ToReviewSecurityHotspots is the number of security hotspots to review.
	ToReviewSecurityHotspots int `json:"toReviewSecurityHotspots,omitempty"`
	// ReviewedSecurityHotspots is the number of reviewed security hotspots.
	ReviewedSecurityHotspots int `json:"reviewedSecurityHotspots,omitempty"`
	// HasMoreRules indicates whether more rules exist beyond TotalRules.
	HasMoreRules bool `json:"hasMoreRules,omitempty"`
}

// SecurityReportCWEDistribution represents the per-CWE breakdown within a security
// report category, returned when includeDistribution is set.
type SecurityReportCWEDistribution struct {
	// CWE is the CWE identifier (e.g. "89").
	CWE string `json:"cwe,omitempty"`
	// Vulnerabilities is the number of vulnerabilities for this CWE.
	Vulnerabilities int `json:"vulnerabilities,omitempty"`
	// VulnerabilityRating is the vulnerability rating (1-5, A-E) for this CWE.
	VulnerabilityRating int `json:"vulnerabilityRating,omitempty"`
	// SecurityReviewRating is the security review rating (1-5, A-E) for this CWE.
	SecurityReviewRating int `json:"securityReviewRating,omitempty"`
	// ActiveRules is the number of active rules for this CWE.
	ActiveRules int `json:"activeRules,omitempty"`
	// TotalRules is the total number of rules (active or not) for this CWE.
	TotalRules int `json:"totalRules,omitempty"`
	// ToReviewSecurityHotspots is the number of security hotspots to review.
	ToReviewSecurityHotspots int `json:"toReviewSecurityHotspots,omitempty"`
	// ReviewedSecurityHotspots is the number of reviewed security hotspots.
	ReviewedSecurityHotspots int `json:"reviewedSecurityHotspots,omitempty"`
	// HasMoreRules indicates whether more rules exist beyond TotalRules.
	HasMoreRules bool `json:"hasMoreRules,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SecurityReportsDownloadOptions contains parameters for the Download method.
type SecurityReportsDownloadOptions struct {
	// Branch filters the results by branch. Optional.
	Branch string `url:"branch,omitempty"`
	// Project is the project key. This field is required.
	Project string `url:"project"`

	// Standards is the list of standards to include in the report.
	// If omitted, all standards are included. Optional.
	Standards []string `url:"standards,omitempty,comma"`
}

// SecurityReportsShowOptions contains parameters for the Show method.
type SecurityReportsShowOptions struct {
	// Branch filters the results by branch. Optional.
	Branch string `url:"branch,omitempty"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
	// Standard is the security standard to report on. This field is required.
	Standard string `url:"standard"`
	// IncludeDistribution includes distribution information when set. Optional.
	IncludeDistribution bool `url:"includeDistribution,omitempty"`
}

// -----------------------------------------------------------------------------
// Allowed Value Sets
// -----------------------------------------------------------------------------

//nolint:gochecknoglobals // constant set of allowed values
var allowedSecurityStandards = map[string]struct{}{
	"owaspTop10":          {},
	"sonarsourceSecurity": {},
	"cweTop25":            {},
	"pciDss":              {},
	"owaspAsvs":           {},
	"stig":                {},
	"casa":                {},
	"owaspMasvs-v2":       {},
	"owaspLlmTop10":       {},
	"owaspMobileTop10":    {},
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateDownloadOpt validates the options for the Download method.
func (s *SecurityReportsService) ValidateDownloadOpt(opt *SecurityReportsDownloadOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	if len(opt.Standards) > 0 {
		return AreValuesAuthorized(opt.Standards, allowedSecurityStandards, "Standards")
	}

	return nil
}

// ValidateShowOpt validates the options for the Show method.
func (s *SecurityReportsService) ValidateShowOpt(opt *SecurityReportsShowOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Standard, "Standard")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Standard, allowedSecurityStandards, "Standard")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Download downloads a security report as a PDF document.
// Requires Browse permission on the project.
//
// API endpoint: GET /api/security_reports/download.
// Since: 8.8.
// Internal endpoint.
func (s *SecurityReportsService) Download(ctx context.Context, opt *SecurityReportsDownloadOptions) ([]byte, *http.Response, error) {
	err := s.ValidateDownloadOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "security_reports/download", opt)
	if err != nil {
		return nil, nil, err
	}

	var buf bytes.Buffer

	resp, err := s.client.Do(req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, nil
}

// Show returns the security report for a project.
// Requires Browse permission on the project.
//
// API endpoint: GET /api/security_reports/show.
// Since: 7.3.
// Internal endpoint.
func (s *SecurityReportsService) Show(ctx context.Context, opt *SecurityReportsShowOptions) (*SecurityReportsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "security_reports/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SecurityReportsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
