package sonar

import (
	"bytes"
	"context"
	"net/http"
)

// AuditLogsService handles communication with the audit logs related methods of
// the SonarQube API. This service is only available in Enterprise Edition.
type AuditLogsService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AuditLogsDownloadOptions contains parameters for the Download method.
type AuditLogsDownloadOptions struct {
	// From is the start datetime in ISO 8601 format. This field is required.
	From string `url:"from"`
	// To is the end datetime in ISO 8601 format. This field is required.
	To string `url:"to"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateDownloadOpt validates the options for the Download method.
func (s *AuditLogsService) ValidateDownloadOpt(opt *AuditLogsDownloadOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.From, "From")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.To, "To")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Download downloads the audit logs for the given time range as raw JSON bytes.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/audit_logs/download.
// Since: 9.1.
// Enterprise Edition only.
func (s *AuditLogsService) Download(ctx context.Context, opt *AuditLogsDownloadOptions) ([]byte, *http.Response, error) {
	err := s.ValidateDownloadOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "audit_logs/download", opt)
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
