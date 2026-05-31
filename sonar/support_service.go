package sonar

import (
	"bytes"
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

// Info returns raw support information as bytes.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/support/info.
// Since: 3.1.
// Enterprise Edition only.
func (s *SupportService) Info(ctx context.Context) ([]byte, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "support/info", nil)
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
