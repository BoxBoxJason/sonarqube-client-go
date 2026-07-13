package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// AtlassianService handles communication with the Atlassian (Jira/Confluence
// Connect) related methods of the SonarQube V2 API.
type AtlassianService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AtlassianAuthenticationDetails represents the Atlassian OAuth application
// configuration. Only the client ID is ever returned by the API; the secret
// is write-only and never echoed back for security purposes.
type AtlassianAuthenticationDetails struct {
	// ClientId is the Atlassian 3LO App Client ID.
	ClientId string `json:"clientId,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// AtlassianAuthenticationConfigureOptions contains parameters for creating or
// updating the Atlassian OAuth application configuration.
type AtlassianAuthenticationConfigureOptions struct {
	// ClientId is the Atlassian 3LO App Client ID. This field is required.
	ClientId string `json:"clientId"`
	// Secret is the Atlassian 3LO App Secret. This field is required.
	Secret string `json:"secret"`
}

// AtlassianAuthURLOptions contains query parameters for generating the Jira
// OAuth authentication URL.
type AtlassianAuthURLOptions struct {
	// SonarOrganizationKey is the key of the SonarQube organization initiating
	// the Jira OAuth flow. This field is required.
	SonarOrganizationKey string `json:"sonarOrganizationKey"`
	// SonarOrganizationUuid is the UUID of the SonarQube organization initiating
	// the Jira OAuth flow. This field is required.
	SonarOrganizationUuid string `json:"sonarOrganizationUuid"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateConfigureApplicationOpt validates the AtlassianAuthenticationConfigureOptions.
func (s *AtlassianService) ValidateConfigureApplicationOpt(opt *AtlassianAuthenticationConfigureOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ClientId, "ClientId")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Secret, "Secret")
}

// ValidateAuthURLOpt validates the AtlassianAuthURLOptions.
func (s *AtlassianService) ValidateAuthURLOpt(opt *AtlassianAuthURLOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.SonarOrganizationKey, "SonarOrganizationKey")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.SonarOrganizationUuid, "SonarOrganizationUuid")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetApplicationConfiguration fetches the Atlassian Authentication details if
// they exist. Only the client ID is returned, for security purposes.
// Requires authenticated user.
func (s *AtlassianService) GetApplicationConfiguration(ctx context.Context) (*AtlassianAuthenticationDetails, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "atlassian/application-configuration", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AtlassianAuthenticationDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateOrUpdateApplicationConfiguration creates or updates the Atlassian
// Authentication details (client ID and secret) used for the Jira/Confluence
// OAuth integration.
func (s *AtlassianService) CreateOrUpdateApplicationConfiguration(ctx context.Context, opt *AtlassianAuthenticationConfigureOptions) (*AtlassianAuthenticationDetails, *http.Response, error) {
	err := s.ValidateConfigureApplicationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "atlassian/application-configuration", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AtlassianAuthenticationDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// atlassianAuthURLResponse is a defined string type (rather than a bare
// string) used solely to decode the GetAuthURL response body. client.Do
// treats a destination of type *string as an opaque text/plain payload and
// forces an "Accept: text/plain" request header for it. The API spec
// declares this endpoint's 200 response as "application/json" with a string
// schema, not "text/plain", and V2 endpoints strictly enforce their declared
// content type (see architecture_v2_service.go's architectureFileGraphResponse
// for the same reasoning verified against a live SonarQube instance). Using a
// distinct named type keeps client.Do on its default JSON-decode path, which
// matches the endpoint's contract.
type atlassianAuthURLResponse string

// GetAuthURL generates the authentication URL to start the Jira OAuth flow
// for the given SonarQube organization. Requires authenticated user.
func (s *AtlassianService) GetAuthURL(ctx context.Context, opt *AtlassianAuthURLOptions) (*string, *http.Response, error) {
	err := s.ValidateAuthURLOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "atlassian/auth-url", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result atlassianAuthURLResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	strResult := string(result)

	return &strResult, resp, nil
}
