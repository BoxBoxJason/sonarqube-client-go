package sonar

import (
	"context"
	"net/http"
)

// GithubProvisioningService handles communication with the GitHub provisioning related methods
// of the SonarQube API.
// This service manages GitHub provisioning operations.
type GithubProvisioningService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// GithubProvisioningCheck represents the response from checking GitHub provisioning configuration.
type GithubProvisioningCheck struct {
	// Application contains provisioning status for the application level.
	Application GithubProvisioningApplicationStatus `json:"application,omitzero"`
	// Installations contains provisioning status for each GitHub organization/installation.
	Installations []GithubProvisioningInstallation `json:"installations,omitempty"`
}

// GithubProvisioningApplicationStatus represents the provisioning status at the application level.
type GithubProvisioningApplicationStatus struct {
	// AutoProvisioning contains the auto-provisioning status.
	AutoProvisioning GithubProvisioningStatus `json:"autoProvisioning,omitzero"`
	// Jit contains the Just-In-Time provisioning status.
	Jit GithubProvisioningJitStatus `json:"jit,omitzero"`
}

// GithubProvisioningInstallation represents the provisioning status for a GitHub organization/installation.
type GithubProvisioningInstallation struct {
	// AutoProvisioning contains the auto-provisioning status for this installation.
	AutoProvisioning GithubProvisioningStatus `json:"autoProvisioning,omitzero"`
	// Jit contains the Just-In-Time provisioning status for this installation.
	Jit GithubProvisioningJitStatus `json:"jit,omitzero"`
	// Organization is the name of the GitHub organization.
	Organization string `json:"organization,omitempty"`
}

// GithubProvisioningStatus represents the status of auto-provisioning.
type GithubProvisioningStatus struct {
	// ErrorMessage contains an error message if provisioning failed.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Status is the provisioning status.
	Status string `json:"status,omitempty"`
}

// GithubProvisioningJitStatus represents the status of Just-In-Time provisioning.
type GithubProvisioningJitStatus struct {
	// Status is the JIT provisioning status.
	Status string `json:"status,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Check validates the GitHub provisioning configuration.
//
// API endpoint: POST /api/github_provisioning/check.
// WARNING: This is an internal API and may change without notice.
func (s *GithubProvisioningService) Check(ctx context.Context) (*GithubProvisioningCheck, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "github_provisioning/check", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(GithubProvisioningCheck)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
