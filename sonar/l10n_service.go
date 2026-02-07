package sonar

import "net/http"

// L10NService handles communication with the localization related methods
// of the SonarQube API.
// This service manages localization.
type L10NService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// L10NIndex represents the response from getting localization messages.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type L10NIndex struct {
	// Locale is the locale used.
	Locale string `json:"locale,omitempty"`
	// Messages is a map of message keys to their localized values.
	Messages map[string]string `json:"messages,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// L10NIndexOption contains parameters for the Index method.
type L10NIndexOption struct {
	// Locale is the BCP47 language tag, used to override the browser Accept-Language header.
	Locale string `url:"locale,omitempty"`
	// Timestamp is the date of the last cache update.
	Timestamp string `url:"ts,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateIndexOpt validates the options for the Index method.
func (s *L10NService) ValidateIndexOpt(opt *L10NIndexOption) error {
	// Options are optional; nothing to validate.
	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Index gets all localization messages for a given locale.
//
// API endpoint: GET /api/l10n/index.
// Warning: This API is internal and may change without notice.
func (s *L10NService) Index(opt *L10NIndexOption) (*L10NIndex, *http.Response, error) {
	err := s.ValidateIndexOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "l10n/index", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(L10NIndex)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
