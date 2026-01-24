package sonargo

import "net/http"

// DevelopersService handles communication with the developers related methods
// of the SonarQube API.
// This service provides data needed by SonarLint.
type DevelopersService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// DevelopersSearchEvents represents the response from the search events endpoint.
type DevelopersSearchEvents struct {
	// Events is the list of developer events.
	Events []DeveloperEvent `json:"events,omitempty"`
}

// DeveloperEvent represents an event in the developer activity feed.
type DeveloperEvent struct {
	// Category is the type of event.
	Category string `json:"category,omitempty"`
	// Link is the URL to more information about the event.
	Link string `json:"link,omitempty"`
	// Message is the event message.
	Message string `json:"message,omitempty"`
	// Project is the project key associated with the event.
	Project string `json:"project,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// DevelopersSearchEventsOption contains parameters for the SearchEvents method.
type DevelopersSearchEventsOption struct {
	// From is a comma-separated list of datetimes.
	// Filter events created after the given date (exclusive).
	// This field is required.
	From []string `url:"from,comma"`
	// Projects is a comma-separated list of project keys to search notifications for.
	// This field is required.
	Projects []string `url:"projects,comma"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSearchEventsOpt validates the options for the SearchEvents method.
func (s *DevelopersService) ValidateSearchEventsOpt(opt *DevelopersSearchEventsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if len(opt.From) == 0 {
		return NewValidationError("From", "is required", ErrMissingRequired)
	}

	if len(opt.Projects) == 0 {
		return NewValidationError("Projects", "is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// SearchEvents searches for developer events.
// Requires authentication.
// When issue indexing is in progress, returns 503 service unavailable HTTP code.
//
// API endpoint: GET /api/developers/search_events.
// WARNING: This is an internal API and may change without notice.
func (s *DevelopersService) SearchEvents(opt *DevelopersSearchEventsOption) (*DevelopersSearchEvents, *http.Response, error) {
	err := s.ValidateSearchEventsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "developers/search_events", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(DevelopersSearchEvents)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
