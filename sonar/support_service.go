package sonar

import (
	"context"
	"net/http"
)

// SupportService handles communication with the support related methods of
// the SonarQube API. This service is only available in Enterprise Edition.
type SupportService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SupportInfo contains system support information returned by the Info endpoint.
//
//nolint:tagliatelle // SonarQube API uses PascalCase JSON keys for these fields
type SupportInfo struct {
	// Database contains database-related support information.
	Database map[string]any `json:"Database,omitempty"`
	// SonarQube contains SonarQube application-level support information.
	SonarQube map[string]any `json:"SonarQube,omitempty"`
	// Statistics contains usage statistics information.
	Statistics map[string]any `json:"Statistics,omitempty"`
	// System contains system-level support information.
	System map[string]any `json:"System,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Info returns support information about the system and the currently
// installed license.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/support/info.
// Since: 3.1.
// Enterprise Edition only.
// WARNING: This is an internal API and may change without notice.
// WARNING: Live-verified against SonarQube 2025.2 Enterprise: this endpoint
// requires an actual installed commercial license to succeed, independent of
// Enterprise Edition + 'Administer System' permission. Without one it fails
// with HTTP 400 and body {"errors":[{"msg":"License not found"}]}, so the
// success-path JSON shape of SupportInfo below could not be confirmed against
// a live server and remains a best-effort guess (see .github/reviews/pr-265.md).
func (s *SupportService) Info(ctx context.Context) (*SupportInfo, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "support/info", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SupportInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
