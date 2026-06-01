package sonar

import (
	"bytes"
	"context"
	"net/http"
)

// SamlService handles communication with the SAML related methods of
// the SonarQube API. This service is internal.
type SamlService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SamlValidationOptions contains parameters for the Validation method.
type SamlValidationOptions struct {
	// SAMLResponse is the SAML response to validate. This field is required.
	SAMLResponse string `url:"SAMLResponse"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateValidationOpt validates the options for the Validation method.
func (s *SamlService) ValidateValidationOpt(opt *SamlValidationOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.SAMLResponse, "SAMLResponse")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Validation validates a SAML response and returns the resulting HTML page.
// Requires authentication.
//
// API endpoint: POST saml/validation.
// Since: 9.7.
// Internal endpoint.
func (s *SamlService) Validation(ctx context.Context, opt *SamlValidationOptions) ([]byte, *http.Response, error) {
	err := s.ValidateValidationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "saml/validation", opt)
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

// ValidationInit initiates the SAML validation flow.
// Requires authentication.
//
// API endpoint: GET saml/validation_init.
// Since: 9.7.
// Internal endpoint.
func (s *SamlService) ValidationInit(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "saml/validation_init", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
