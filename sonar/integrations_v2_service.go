package sonar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// IntegrationsService handles communication with the generic third-party
// integrations (e.g. Slack) related methods of the SonarQube V2 API.
type IntegrationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Constants
// -----------------------------------------------------------------------------

const (
	// IntegrationTypeSlack designates the Slack integration type.
	IntegrationTypeSlack = "SLACK"
)

//nolint:gochecknoglobals // constant set of allowed values
var allowedIntegrationsIntegrationTypes = map[string]struct{}{
	IntegrationTypeSlack: {},
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// IntegrationsUserBindingData contains the data required to complete a user
// binding, such as the OAuth authorization code exchanged during the Slack
// connect flow.
type IntegrationsUserBindingData struct {
	// Code is the OAuth authorization code returned by the external
	// integration provider. This field is required.
	Code string `json:"code"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// IntegrationsUserBinding represents a binding between a SonarQube user and an
// external chat identity (e.g. a Slack user).
//
//nolint:tagliatelle // JSON tags match SonarQube API field names (snake_case for this endpoint)
type IntegrationsUserBinding struct {
	// Id is the unique identifier of the user binding.
	Id string `json:"id,omitempty"`
	// UserId is the identifier of the SonarQube user.
	UserId string `json:"user_id,omitempty"`
	// SlackUserId is the identifier of the bound Slack user.
	SlackUserId string `json:"slack_user_id,omitempty"`
	// SlackWorkspaceId is the identifier of the Slack workspace the user belongs to.
	SlackWorkspaceId string `json:"slack_workspace_id,omitempty"`
	// SlackWorkspaceName is the name of the Slack workspace the user belongs to.
	SlackWorkspaceName string `json:"slack_workspace_name,omitempty"`
	// CreatedAt is the Unix timestamp (in milliseconds) when the binding was created.
	CreatedAt int64 `json:"created_at,omitempty"`
}

// IntegrationsIntegrationConfiguration represents a workspace-level
// third-party integration configuration (e.g. a Slack app install).
type IntegrationsIntegrationConfiguration struct {
	// Id is the unique identifier of the integration configuration.
	Id string `json:"id,omitempty"`
	// IntegrationType is the type of integration. Allowed values: SLACK.
	IntegrationType string `json:"integrationType,omitempty"`
	// ClientId is the client ID of the integration application.
	ClientId string `json:"clientId,omitempty"`
	// AppId is the identifier of the integration application.
	AppId string `json:"appId,omitempty"`
}

// IntegrationsIntegrationConfigurationSearch represents the response from
// listing integration configurations.
type IntegrationsIntegrationConfigurationSearch struct {
	// IntegrationConfigurations is the list of integration configurations.
	IntegrationConfigurations []IntegrationsIntegrationConfiguration `json:"integrationConfigurations,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types (Query Parameters)
// -----------------------------------------------------------------------------

// IntegrationsListIntegrationConfigurationsOptions contains query parameters
// for the ListIntegrationConfigurations method.
type IntegrationsListIntegrationConfigurationsOptions struct {
	// IntegrationType filters configurations by integration type.
	// Allowed values: SLACK. This field is required.
	IntegrationType string `json:"integrationType"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// IntegrationsUserBindingCreateOptions contains parameters for creating a
// user binding between a SonarQube user and an external chat identity.
type IntegrationsUserBindingCreateOptions struct {
	// UserId is the identifier of the SonarQube user. This field is required.
	UserId string `json:"userId"`
	// BindingData contains the data required to complete the binding.
	// This field is required.
	BindingData IntegrationsUserBindingData `json:"bindingData"`
}

// IntegrationsIntegrationConfigurationCreateOptions contains parameters for
// creating an integration configuration.
type IntegrationsIntegrationConfigurationCreateOptions struct {
	// IntegrationType is the type of integration to configure.
	// Allowed values: SLACK. This field is required.
	IntegrationType string `json:"integrationType"`
	// ClientId is the client ID of the integration application. This field is required.
	ClientId string `json:"clientId"`
	// ClientSecret is the client secret of the integration application. This field is required.
	ClientSecret string `json:"clientSecret"`
	// SigningSecret is the signing secret used to validate incoming webhook
	// requests from the integration provider. This field is required.
	SigningSecret string `json:"signingSecret"`
}

// IntegrationsIntegrationConfigurationPatchOptions contains parameters for
// updating an integration configuration. All fields are optional (PATCH merge
// semantics).
type IntegrationsIntegrationConfigurationPatchOptions struct {
	// ClientId is the new client ID of the integration application.
	ClientId string `json:"clientId,omitempty"`
	// ClientSecret is the new client secret of the integration application.
	ClientSecret string `json:"clientSecret,omitempty"`
	// SigningSecret is the new signing secret used to validate incoming
	// webhook requests from the integration provider.
	SigningSecret string `json:"signingSecret,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateCreateUserBindingOpt validates the IntegrationsUserBindingCreateOptions.
func (s *IntegrationsService) ValidateCreateUserBindingOpt(opt *IntegrationsUserBindingCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.UserId, "UserId")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.BindingData.Code, "BindingData.Code")
}

// validateSlackSlashCommandForm validates the form-encoded body of an
// incoming Slack slash command request.
func validateSlackSlashCommandForm(form url.Values) error {
	if len(form) == 0 {
		return NewValidationError("form", "must not be empty", ErrMissingRequired)
	}

	return nil
}

// validateSlackEventPayload validates the JSON body of an incoming Slack
// Events API request.
func validateSlackEventPayload(payload json.RawMessage) error {
	if len(payload) == 0 {
		return NewValidationError("payload", "must not be empty", ErrMissingRequired)
	}

	return nil
}

// ValidateListIntegrationConfigurationsOpt validates the
// IntegrationsListIntegrationConfigurationsOptions.
func (s *IntegrationsService) ValidateListIntegrationConfigurationsOpt(opt *IntegrationsListIntegrationConfigurationsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.IntegrationType, "IntegrationType")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.IntegrationType, allowedIntegrationsIntegrationTypes, "IntegrationType")
}

// ValidateCreateIntegrationConfigurationOpt validates the
// IntegrationsIntegrationConfigurationCreateOptions.
func (s *IntegrationsService) ValidateCreateIntegrationConfigurationOpt(opt *IntegrationsIntegrationConfigurationCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.IntegrationType, "IntegrationType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.IntegrationType, allowedIntegrationsIntegrationTypes, "IntegrationType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ClientId, "ClientId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ClientSecret, "ClientSecret")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.SigningSecret, "SigningSecret")
}

// ValidateUpdateIntegrationConfigurationRequest validates the parameters for
// updating an integration configuration.
func (s *IntegrationsService) ValidateUpdateIntegrationConfigurationRequest(configurationID string, opt *IntegrationsIntegrationConfigurationPatchOptions) error {
	err := ValidateRequired(configurationID, "Id")
	if err != nil {
		return err
	}

	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CreateUserBinding creates a new user binding between a SonarQube user and
// their external chat account (e.g. Slack). This endpoint is used during the
// OAuth flow when users connect their accounts via a slash command (e.g.
// "/sonarqube-server connect").
func (s *IntegrationsService) CreateUserBinding(ctx context.Context, opt *IntegrationsUserBindingCreateOptions) (*IntegrationsUserBinding, *http.Response, error) {
	err := s.ValidateCreateUserBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "integrations/user-bindings", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IntegrationsUserBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetUserBinding retrieves a user binding by its ID. Users can only access
// their own bindings unless they have administrative privileges.
func (s *IntegrationsService) GetUserBinding(ctx context.Context, bindingID string) (*IntegrationsUserBinding, *http.Response, error) {
	err := ValidateRequired(bindingID, "Id")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "integrations/user-bindings/"+bindingID, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IntegrationsUserBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// HandleSlackSlashCommand forwards an incoming Slack slash command request
// (e.g. "/sonarqube-server connect") to the SonarQube server for processing.
// The request body is form-encoded (application/x-www-form-urlencoded) and
// typically contains the standard Slack slash-command fields (command, text,
// user_id, team_id, response_url, etc.). Slack signature validation is
// performed server-side before this data is processed.
func (s *IntegrationsService) HandleSlackSlashCommand(ctx context.Context, form url.Values) (*http.Response, error) {
	err := validateSlackSlashCommandForm(form)
	if err != nil {
		return nil, err
	}

	//nolint:exhaustruct // Headers/RawQuery intentionally left at zero value
	req, err := s.client.NewSonarQubeAPIRequest(ctx, SonarAPIRequestParameters{
		Method: http.MethodPost,
		Path:   v2BasePath + "integrations/slack/slash-commands",
		Body:   form,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// HandleSlackEvent forwards an incoming Slack Events API request (e.g. URL
// verification or workspace lifecycle events) to the SonarQube server for
// processing. The payload is the raw JSON body of the Slack event. Slack
// signature validation is performed server-side before this data is
// processed. The response is the raw text/plain body returned by the server
// (e.g. the echoed "challenge" value during URL verification).
func (s *IntegrationsService) HandleSlackEvent(ctx context.Context, payload json.RawMessage) (*string, *http.Response, error) {
	err := validateSlackEventPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "integrations/slack/events", nil, payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result string

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListIntegrationConfigurations lists the integration configurations for the
// given integration type. Requires global administrator permission.
func (s *IntegrationsService) ListIntegrationConfigurations(ctx context.Context, opt *IntegrationsListIntegrationConfigurationsOptions) (*IntegrationsIntegrationConfigurationSearch, *http.Response, error) {
	err := s.ValidateListIntegrationConfigurationsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "integrations/integration-configurations", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IntegrationsIntegrationConfigurationSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateIntegrationConfiguration creates a configuration for a third-party
// integration (e.g. installing a Slack app for a workspace). Requires global
// administrator permission.
func (s *IntegrationsService) CreateIntegrationConfiguration(ctx context.Context, opt *IntegrationsIntegrationConfigurationCreateOptions) (*IntegrationsIntegrationConfiguration, *http.Response, error) {
	err := s.ValidateCreateIntegrationConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "integrations/integration-configurations", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IntegrationsIntegrationConfiguration)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteIntegrationConfiguration permanently deletes an integration
// configuration and all its related data (workspaces, user bindings, and
// subscriptions). Requires global administrator permission.
func (s *IntegrationsService) DeleteIntegrationConfiguration(ctx context.Context, configurationID string) (*http.Response, error) {
	err := ValidateRequired(configurationID, "Id")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "integrations/integration-configurations/"+configurationID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateIntegrationConfiguration updates an integration configuration.
// Requires global administrator permission.
func (s *IntegrationsService) UpdateIntegrationConfiguration(ctx context.Context, configurationID string, opt *IntegrationsIntegrationConfigurationPatchOptions) (*IntegrationsIntegrationConfiguration, *http.Response, error) {
	err := s.ValidateUpdateIntegrationConfigurationRequest(configurationID, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "integrations/integration-configurations/"+configurationID, nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IntegrationsIntegrationConfiguration)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
