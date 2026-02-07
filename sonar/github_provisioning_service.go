package sonar

import "net/http"

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
	Application ProvisioningApplicationStatus `json:"application,omitzero"`
	// Installations contains provisioning status for each GitHub organization/installation.
	Installations []ProvisioningInstallation `json:"installations,omitempty"`
}

// ProvisioningApplicationStatus represents the provisioning status at the application level.
type ProvisioningApplicationStatus struct {
	// AutoProvisioning contains the auto-provisioning status.
	AutoProvisioning ProvisioningStatus `json:"autoProvisioning,omitzero"`
	// Jit contains the Just-In-Time provisioning status.
	Jit JitStatus `json:"jit,omitzero"`
}

// ProvisioningInstallation represents the provisioning status for a GitHub organization/installation.
type ProvisioningInstallation struct {
	// AutoProvisioning contains the auto-provisioning status for this installation.
	AutoProvisioning ProvisioningStatus `json:"autoProvisioning,omitzero"`
	// Jit contains the Just-In-Time provisioning status for this installation.
	Jit JitStatus `json:"jit,omitzero"`
	// Organization is the name of the GitHub organization.
	Organization string `json:"organization,omitempty"`
}

// ProvisioningStatus represents the status of auto-provisioning.
type ProvisioningStatus struct {
	// ErrorMessage contains an error message if provisioning failed.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Status is the provisioning status.
	Status string `json:"status,omitempty"`
}

// JitStatus represents the status of Just-In-Time provisioning.
type JitStatus struct {
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
func (s *GithubProvisioningService) Check() (*GithubProvisioningCheck, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "github_provisioning/check", nil)
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
