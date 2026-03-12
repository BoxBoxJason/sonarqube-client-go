package sonar

import (
	"fmt"
	"net/http"
)

// SystemServiceV2 handles communication with the System related methods of the
// SonarQube V2 API.
type SystemServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SystemDbMigrationsStatusV2 represents the response from getting database migration status.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SystemDbMigrationsStatusV2 struct {
	// CompletedSteps is the number of migration steps completed.
	CompletedSteps int32 `json:"completedSteps,omitempty"`
	// ExpectedFinishTimestamp is the estimated finish time for the migration.
	ExpectedFinishTimestamp string `json:"expectedFinishTimestamp,omitempty"`
	// Message is a descriptive message about the migration status.
	Message string `json:"message,omitempty"`
	// StartedAt is the datetime when the migration started.
	StartedAt string `json:"startedAt,omitempty"`
	// Status is the current status of the migration.
	Status string `json:"status,omitempty"`
	// TotalSteps is the total number of migration steps.
	TotalSteps int32 `json:"totalSteps,omitempty"`
}

// SystemHealthV2 represents the response from getting system health.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SystemHealthV2 struct {
	// Causes lists the reasons for the current health status.
	Causes []string `json:"causes,omitempty"`
	// Status is the health status (GREEN, YELLOW, RED).
	Status string `json:"status,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SystemPasscodeOptionV2 contains the optional passcode header for system endpoints.
type SystemPasscodeOptionV2 struct {
	// Passcode is the value for the X-Sonar-Passcode header.
	// Can be provided instead of system admin credentials.
	Passcode string
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetMigrationsStatus returns the detailed status of ongoing database migrations.
// If no migration is ongoing or needed, it still returns appropriate information.
func (s *SystemServiceV2) GetMigrationsStatus() (*SystemDbMigrationsStatusV2, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "system/migrations-status", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(SystemDbMigrationsStatusV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CheckLiveness provides liveness of SonarQube, meant to be used as a liveness
// probe on Kubernetes. Returns a 204 status when alive.
// Requires 'Administer System' permission or authentication with passcode.
func (s *SystemServiceV2) CheckLiveness(opt *SystemPasscodeOptionV2) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "system/liveness", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if opt != nil && opt.Passcode != "" {
		req.Header.Set("X-Sonar-Passcode", opt.Passcode)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetHealth returns the health status of the SonarQube instance.
// Requires 'Administer System' permission or authentication with passcode.
func (s *SystemServiceV2) GetHealth(opt *SystemPasscodeOptionV2) (*SystemHealthV2, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "system/health", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if opt != nil && opt.Passcode != "" {
		req.Header.Set("X-Sonar-Passcode", opt.Passcode)
	}

	result := new(SystemHealthV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
