package sonar

import (
	"context"
	"net/http"
)

// ScimManagementService handles communication with the SCIM management related
// methods of the SonarQube API. This service is internal.
type ScimManagementService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ScimManagementStatus represents the SCIM provisioning status.
type ScimManagementStatus struct {
	// Enabled indicates whether SCIM provisioning is enabled.
	Enabled bool `json:"enabled,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Disable disables SCIM provisioning.
// Requires Global Admin permission.
//
// API endpoint: POST /api/scim_management/disable.
// Since: 10.0.
// Internal endpoint.
func (s *ScimManagementService) Disable(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "scim_management/disable", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Enable enables SCIM provisioning.
// Requires Global Admin permission.
//
// API endpoint: POST /api/scim_management/enable.
// Since: 10.0.
// Internal endpoint.
func (s *ScimManagementService) Enable(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "scim_management/enable", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Status returns the SCIM provisioning status.
// Requires Global Admin permission.
//
// API endpoint: GET /api/scim_management/status.
// Since: 10.0.
// Internal endpoint.
func (s *ScimManagementService) Status(ctx context.Context) (*ScimManagementStatus, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "scim_management/status", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ScimManagementStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
