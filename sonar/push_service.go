package sonar

import (
	"fmt"
	"net/http"
)

// PushService handles communication with the server-side events related methods
// of the SonarQube API.
// This service provides endpoints for listening to server-side events.
type PushService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SonarlintEvents represents the response from the SonarLint events endpoint.
// The structure of events is dynamic and depends on the event type.
// Currently, it notifies listeners about changes to activation of a rule.
type SonarlintEvents struct{}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// PushSonarlintEventsOption contains parameters for the SonarlintEvents method.
type PushSonarlintEventsOption struct {
	// Languages is a list of languages for which events will be delivered.
	// This field is required.
	Languages []string `url:"languages,comma"`
	// ProjectKeys is a list of project keys for which events will be delivered.
	// This field is required.
	ProjectKeys []string `url:"projectKeys,comma"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSonarlintEventsOpt validates the options for the SonarlintEvents method.
func (s *PushService) ValidateSonarlintEventsOpt(opt *PushSonarlintEventsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if len(opt.Languages) == 0 {
		return NewValidationError("Languages", "is required", ErrMissingRequired)
	}

	if len(opt.ProjectKeys) == 0 {
		return NewValidationError("ProjectKeys", "is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// SonarlintEvents provides an endpoint for listening to server-side events.
// Currently, it notifies listeners about changes to the activation of a rule.
// The response body is left open for streaming events. The caller is responsible
// for reading from resp.Body (e.g., using an event stream parser) and closing
// the response body when finished.
//
// API endpoint: GET /api/push/sonarlint_events.
// WARNING: This is an internal API and may change without notice.
func (s *PushService) SonarlintEvents(opt *PushSonarlintEventsOption) (*http.Response, error) {
	err := s.ValidateSonarlintEventsOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "push/sonarlint_events", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
