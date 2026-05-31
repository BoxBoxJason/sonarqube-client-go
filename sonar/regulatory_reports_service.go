package sonar

import (
	"bytes"
	"context"
	"net/http"
)

// RegulatoryReportsService handles communication with the regulatory reports related methods of
// the SonarQube API. This service is only available in Enterprise Edition.
type RegulatoryReportsService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// RegulatoryReportsDownloadOptions contains parameters for the Download method.
type RegulatoryReportsDownloadOptions struct {
	// Branch is the branch key. Optional.
	Branch string `url:"branch,omitempty"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateDownloadOpt validates the options for the Download method.
func (s *RegulatoryReportsService) ValidateDownloadOpt(opt *RegulatoryReportsDownloadOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Project, "Project")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Download downloads the regulatory report for a project as a zip file.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/regulatory_reports/download.
// Since: 9.5.
// Enterprise Edition only.
func (s *RegulatoryReportsService) Download(ctx context.Context, opt *RegulatoryReportsDownloadOptions) ([]byte, *http.Response, error) {
	err := s.ValidateDownloadOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "regulatory_reports/download", opt)
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
