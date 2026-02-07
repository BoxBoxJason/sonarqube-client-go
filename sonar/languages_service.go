package sonar

import "net/http"

// LanguagesService handles communication with the languages related methods
// of the SonarQube API.
// This service provides information about programming languages supported by SonarQube.
type LanguagesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// LanguagesList represents the response from the languages list endpoint.
type LanguagesList struct {
	// Languages is the list of supported programming languages.
	Languages []Language `json:"languages,omitempty"`
}

// Language represents a programming language supported by SonarQube.
type Language struct {
	// Key is the unique identifier for the language.
	Key string `json:"key,omitempty"`
	// Name is the display name of the language.
	Name string `json:"name,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// LanguagesListOption contains parameters for the List method.
type LanguagesListOption struct {
	// Query is a pattern to match language keys/names against.
	Query string `url:"q,omitempty"`
	// PageSize is the size of the list to return. Use 0 for all languages.
	// Default is 0.
	PageSize int64 `url:"ps,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListOpt validates the options for the List method.
// Currently, there are no validation rules for this method,
// but this function is provided for consistency and future extensibility.
func (s *LanguagesService) ValidateListOpt(opt *LanguagesListOption) error {
	// No required fields, no validation needed
	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// List returns the list of programming languages supported in this SonarQube instance.
//
// API endpoint: GET /api/languages/list.
// Since: 5.1.
func (s *LanguagesService) List(opt *LanguagesListOption) (*LanguagesList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "languages/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(LanguagesList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
