package sonar

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

// SamlService handles communication with the SAML related methods of
// the SonarQube API. This service is only available in Enterprise Edition.
type SamlService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SamlValidationOptions contains parameters for the Validation method.
type SamlValidationOptions struct {
	// SAMLResponse is the base64-encoded SAML assertion returned by the
	// identity provider. This field is required.
	//
	// This value is sent in the POST request body (application/x-www-form-urlencoded),
	// not as a URL query parameter: real-world SAML assertions are typically
	// 2-10 KB once base64-encoded, which comfortably exceeds the URL length
	// limits enforced by most HTTP servers and proxies (commonly 2-8 KB).
	SAMLResponse string
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
// No authentication is required: this endpoint is the SAML assertion
// consumer service (ACS) that the identity provider posts to as part of the
// SSO callback, so the caller cannot be authenticated yet.
//
// The SAMLResponse value is sent as a form-urlencoded POST body rather than a
// URL query parameter, since real SAML assertions are too large to fit
// reliably within URL length limits.
//
// API endpoint: POST saml/validation.
// Since: 9.7.
// Internal endpoint.
func (s *SamlService) Validation(ctx context.Context, opt *SamlValidationOptions) ([]byte, *http.Response, error) {
	err := s.ValidateValidationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	formBody := url.Values{
		"SAMLResponse": {opt.SAMLResponse},
	}

	//nolint:exhaustruct // RawQuery intentionally unset: SAMLResponse is sent in the body, not the query string
	req, err := s.client.NewSonarQubeAPIRequest(ctx, SonarAPIRequestParameters{
		Method: http.MethodPost,
		Path:   "saml/validation",
		Body:   formBody,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	})
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

// ValidationInit initiates the SAML validation flow. No authentication is
// required: this endpoint kicks off the SSO redirect to the identity
// provider so that an administrator can validate the SAML configuration
// before it is enforced.
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
