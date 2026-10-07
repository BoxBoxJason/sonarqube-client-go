package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// LlmConnectivityService handles communication with the LLM provider connectivity related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type LlmConnectivityService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals,goconst // constant sets of allowed values
var (
	// allowedLlmConnectivityGetProviderMappingsAiCapability is the set of allowed values for the corresponding option.
	allowedLlmConnectivityGetProviderMappingsAiCapability = map[string]struct{}{
		"AI_CODEFIX":        {},
		"HUNTER_AGENT":      {},
		"REMEDIATION_AGENT": {},
	}
	// allowedLlmConnectivityGetSupportedModelsAiCapability is the set of allowed values for the corresponding option.
	allowedLlmConnectivityGetSupportedModelsAiCapability = map[string]struct{}{
		"AI_CODEFIX":        {},
		"HUNTER_AGENT":      {},
		"REMEDIATION_AGENT": {},
	}
	// allowedLlmConnectivityUpsertProviderMappingAiCapability is the set of allowed values for the corresponding option.
	allowedLlmConnectivityUpsertProviderMappingAiCapability = map[string]struct{}{
		"AI_CODEFIX":        {},
		"HUNTER_AGENT":      {},
		"REMEDIATION_AGENT": {},
	}
	// allowedLlmConnectivityCreateProviderProvider is the set of allowed values for the corresponding option.
	allowedLlmConnectivityCreateProviderProvider = map[string]struct{}{
		"AZURE_OPENAI":       {},
		"AWS_BEDROCK":        {},
		"AWS_BEDROCK_MANTLE": {},
		"VERTEX_AI":          {},
		"CUSTOM_PROXY":       {},
		"OPENAI":             {},
		"ANTHROPIC":          {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// LlmConnectivityLlmProviderFieldResponse represents the LlmProviderFieldResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type LlmConnectivityLlmProviderFieldResponse struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Label is the label.
	Label string `json:"label,omitempty"`
	// Required is the required.
	Required bool `json:"required,omitempty"`
	// Type is the type. Allowed values: STRING, HTTP_HEADERS.
	Type string `json:"type,omitempty"`
	// Secret is the secret.
	Secret bool `json:"secret,omitempty"`
}

// LlmConnectivityLlmProviderDefinitionResponse represents the LlmProviderDefinitionResponse object
// of the SonarQube V2 API.
type LlmConnectivityLlmProviderDefinitionResponse struct {
	// Provider is the provider. Allowed values: AZURE_OPENAI, AWS_BEDROCK, AWS_BEDROCK_MANTLE,
	// VERTEX_AI, CUSTOM_PROXY, OPENAI, ANTHROPIC.
	Provider string `json:"provider,omitempty"`
	// Label is the label.
	Label string `json:"label,omitempty"`
	// Fields is the fields.
	Fields []LlmConnectivityLlmProviderFieldResponse `json:"fields,omitempty"`
}

// LlmConnectivityLlmProviderDefinitionsResponse represents the LlmProviderDefinitionsResponse
// object of the SonarQube V2 API.
type LlmConnectivityLlmProviderDefinitionsResponse struct {
	// ProviderDefinitions is the provider definitions.
	ProviderDefinitions []LlmConnectivityLlmProviderDefinitionResponse `json:"providerDefinitions,omitempty"`
}

// LlmConnectivityGetProviderMappingsOptions contains parameters for the GetProviderMappings method.
type LlmConnectivityGetProviderMappingsOptions struct {
	// AiCapability is the ai capability. Allowed values: AI_CODEFIX, HUNTER_AGENT,
	// REMEDIATION_AGENT.
	AiCapability string `json:"aiCapability,omitempty"`
}

// LlmConnectivityProviderMappingResponse represents the ProviderMappingResponse object of the
// SonarQube V2 API.
type LlmConnectivityProviderMappingResponse struct {
	// AiCapability is the ai capability. Allowed values: AI_CODEFIX, HUNTER_AGENT,
	// REMEDIATION_AGENT.
	AiCapability string `json:"aiCapability,omitempty"`
	// LlmProviderId is the llm provider id.
	LlmProviderId string `json:"llmProviderId,omitempty"`
	// ModelIdentifier is the model identifier.
	ModelIdentifier string `json:"modelIdentifier,omitempty"`
	// SecondaryModelIdentifier is the secondary model identifier.
	SecondaryModelIdentifier string `json:"secondaryModelIdentifier,omitempty"`
}

// LlmConnectivityProviderMappingsResponse represents the ProviderMappingsResponse object of the
// SonarQube V2 API.
type LlmConnectivityProviderMappingsResponse struct {
	// ProviderMappings is the provider mappings.
	ProviderMappings []LlmConnectivityProviderMappingResponse `json:"providerMappings,omitempty"`
}

// LlmConnectivityLlmProviderResponse represents the LlmProviderResponse object of the SonarQube V2
// API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type LlmConnectivityLlmProviderResponse struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// Provider is the provider. Allowed values: AZURE_OPENAI, AWS_BEDROCK, AWS_BEDROCK_MANTLE,
	// VERTEX_AI, CUSTOM_PROXY, OPENAI, ANTHROPIC.
	Provider string `json:"provider,omitempty"`
	// Label is the label.
	Label string `json:"label,omitempty"`
	// Configuration is the configuration.
	Configuration map[string]any `json:"configuration,omitempty"`
}

// LlmConnectivityLlmProvidersResponse represents the LlmProvidersResponse object of the SonarQube
// V2 API.
type LlmConnectivityLlmProvidersResponse struct {
	// Providers is the providers.
	Providers []LlmConnectivityLlmProviderResponse `json:"providers,omitempty"`
}

// LlmConnectivityGetSupportedModelsOptions contains parameters for the GetSupportedModels method.
type LlmConnectivityGetSupportedModelsOptions struct {
	// AiCapability is the ai capability. This field is required. Allowed values: AI_CODEFIX,
	// HUNTER_AGENT, REMEDIATION_AGENT.
	AiCapability string `json:"aiCapability"`
}

// LlmConnectivitySupportedModelResponse represents the SupportedModelResponse object of the
// SonarQube V2 API.
type LlmConnectivitySupportedModelResponse struct {
	// ModelKey is the model key.
	ModelKey string `json:"modelKey,omitempty"`
	// ModelDisplayName is the model display name.
	ModelDisplayName string `json:"modelDisplayName,omitempty"`
	// ModelIdentifiers is the model identifiers.
	ModelIdentifiers []string `json:"modelIdentifiers,omitempty"`
	// SupportedSecondaryModels is the supported secondary models.
	SupportedSecondaryModels []LlmConnectivitySupportedModelResponse `json:"supportedSecondaryModels,omitempty"`
}

// LlmConnectivitySupportedModelsResponse represents the SupportedModelsResponse object of the
// SonarQube V2 API.
type LlmConnectivitySupportedModelsResponse struct {
	// SupportedModels is the supported models.
	SupportedModels []LlmConnectivitySupportedModelResponse `json:"supportedModels,omitempty"`
}

// LlmConnectivityUpdateProviderOptions contains the request body for the UpdateProvider method.
type LlmConnectivityUpdateProviderOptions struct {
	// Label is the label.
	Label *string `json:"label,omitempty"`
	// Configuration is the configuration.
	Configuration map[string]any `json:"configuration,omitempty"`
}

// LlmConnectivityUpsertProviderMappingOptions contains the request body for the
// UpsertProviderMapping method.
type LlmConnectivityUpsertProviderMappingOptions struct {
	// AiCapability is the ai capability. This field is required. Allowed values: AI_CODEFIX,
	// HUNTER_AGENT, REMEDIATION_AGENT.
	AiCapability string `json:"aiCapability"`
	// LlmProviderId is the llm provider id. This field is required.
	LlmProviderId string `json:"llmProviderId"`
	// ModelIdentifier is the model identifier. This field is required.
	ModelIdentifier string `json:"modelIdentifier"`
	// SecondaryModelIdentifier is the secondary model identifier.
	SecondaryModelIdentifier string `json:"secondaryModelIdentifier,omitempty"`
}

// LlmConnectivityValidateProviderOptions contains the request body for the ValidateProvider method.
type LlmConnectivityValidateProviderOptions struct {
	// LlmProviderId is the llm provider id. This field is required.
	LlmProviderId string `json:"llmProviderId"`
}

// LlmConnectivityProviderValidationResponse represents the ProviderValidationResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type LlmConnectivityProviderValidationResponse struct {
	// LlmProviderId is the llm provider id.
	LlmProviderId string `json:"llmProviderId,omitempty"`
	// IsValid is the is valid.
	IsValid bool `json:"isValid,omitempty"`
	// Error is the error.
	Error string `json:"error,omitempty"`
}

// LlmConnectivityCreateProviderOptions contains the request body for the CreateProvider method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type LlmConnectivityCreateProviderOptions struct {
	// Provider is the provider. Allowed values: AZURE_OPENAI, AWS_BEDROCK, AWS_BEDROCK_MANTLE,
	// VERTEX_AI, CUSTOM_PROXY, OPENAI, ANTHROPIC.
	Provider string `json:"provider,omitempty"`
	// Label is the label.
	Label string `json:"label,omitempty"`
	// Configuration is the configuration.
	Configuration map[string]any `json:"configuration,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetProviderMappingsOpt validates the options for the GetProviderMappings method.
func (s *LlmConnectivityService) ValidateGetProviderMappingsOpt(opt *LlmConnectivityGetProviderMappingsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.AiCapability, allowedLlmConnectivityGetProviderMappingsAiCapability, "AiCapability")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetSupportedModelsOpt validates the options for the GetSupportedModels method.
func (s *LlmConnectivityService) ValidateGetSupportedModelsOpt(opt *LlmConnectivityGetSupportedModelsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AiCapability, "AiCapability")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.AiCapability, allowedLlmConnectivityGetSupportedModelsAiCapability, "AiCapability")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpsertProviderMappingOpt validates the options for the UpsertProviderMapping method.
func (s *LlmConnectivityService) ValidateUpsertProviderMappingOpt(opt *LlmConnectivityUpsertProviderMappingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AiCapability, "AiCapability")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.AiCapability, allowedLlmConnectivityUpsertProviderMappingAiCapability, "AiCapability")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.LlmProviderId, "LlmProviderId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ModelIdentifier, "ModelIdentifier")
	if err != nil {
		return err
	}

	return nil
}

// ValidateValidateProviderOpt validates the options for the ValidateProvider method.
func (s *LlmConnectivityService) ValidateValidateProviderOpt(opt *LlmConnectivityValidateProviderOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.LlmProviderId, "LlmProviderId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateProviderOpt validates the options for the CreateProvider method.
func (s *LlmConnectivityService) ValidateCreateProviderOpt(opt *LlmConnectivityCreateProviderOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.Provider, allowedLlmConnectivityCreateProviderProvider, "Provider")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// DeleteProvider deletes an LLM provider. Blocked when the provider is selected by an AI
// capability. Requires system administration.
//
// API endpoint: DELETE /api/v2/llm-connectivity/llm-providers/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) DeleteProvider(ctx context.Context, providerID string) (*http.Response, error) {
	err := ValidateRequired(providerID, "providerID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "llm-connectivity/llm-providers/"+url.PathEscape(providerID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetProviderDefinitions lists available LLM provider types. Returns the provider types and the
// configuration fields each requires. Requires system administration.
//
// API endpoint: GET /api/v2/llm-connectivity/llm-provider-definitions.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) GetProviderDefinitions(ctx context.Context) (*LlmConnectivityLlmProviderDefinitionsResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "llm-connectivity/llm-provider-definitions", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityLlmProviderDefinitionsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProviderMappings lists per-capability provider mappings. Optionally filtered to a single AI
// capability. Requires system administration.
//
// API endpoint: GET /api/v2/llm-connectivity/llm-provider-mappings.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) GetProviderMappings(ctx context.Context, opt *LlmConnectivityGetProviderMappingsOptions) (*LlmConnectivityProviderMappingsResponse, *http.Response, error) {
	err := s.ValidateGetProviderMappingsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "llm-connectivity/llm-provider-mappings", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityProviderMappingsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListProviders lists configured LLM providers. Requires system administration.
//
// API endpoint: GET /api/v2/llm-connectivity/llm-providers.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) ListProviders(ctx context.Context) (*LlmConnectivityLlmProvidersResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "llm-connectivity/llm-providers", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityLlmProvidersResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProvider gets one LLM provider. Requires system administration.
//
// API endpoint: GET /api/v2/llm-connectivity/llm-providers/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) GetProvider(ctx context.Context, providerID string) (*LlmConnectivityLlmProviderResponse, *http.Response, error) {
	err := ValidateRequired(providerID, "providerID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "llm-connectivity/llm-providers/"+url.PathEscape(providerID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityLlmProviderResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSupportedModels lists supported models for an AI capability. Proxies the agent-orchestrator's
// supported-models endpoint for the given AI capability. Requires system administration.
//
// API endpoint: GET /api/v2/llm-connectivity/supported-models.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) GetSupportedModels(ctx context.Context, opt *LlmConnectivityGetSupportedModelsOptions) (*LlmConnectivitySupportedModelsResponse, *http.Response, error) {
	err := s.ValidateGetSupportedModelsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "llm-connectivity/supported-models", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivitySupportedModelsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateProvider updates an LLM provider. The connection is re-tested before changes are persisted.
// Requires system administration.
//
// API endpoint: PATCH /api/v2/llm-connectivity/llm-providers/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) UpdateProvider(ctx context.Context, providerID string, opt *LlmConnectivityUpdateProviderOptions) (*LlmConnectivityLlmProviderResponse, *http.Response, error) {
	err := ValidateRequired(providerID, "providerID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "llm-connectivity/llm-providers/"+url.PathEscape(providerID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityLlmProviderResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpsertProviderMapping creates or replace the provider mapping for an AI capability. Creates or
// replaces the provider mapping for the given AI capability. Rejected when the provider or
// configured models cannot be validated, or when Hunter does not have both models. Requires system
// administration.
//
// API endpoint: POST /api/v2/llm-connectivity/llm-provider-mappings.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) UpsertProviderMapping(ctx context.Context, opt *LlmConnectivityUpsertProviderMappingOptions) (*LlmConnectivityProviderMappingResponse, *http.Response, error) {
	err := s.ValidateUpsertProviderMappingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "llm-connectivity/llm-provider-mappings", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityProviderMappingResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ValidateProvider runs a connection test for an existing provider. Runs a fresh live connection
// test. Nothing is persisted. Requires system administration.
//
// API endpoint: POST /api/v2/llm-connectivity/llm-provider-validations.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) ValidateProvider(ctx context.Context, opt *LlmConnectivityValidateProviderOptions) (*LlmConnectivityProviderValidationResponse, *http.Response, error) {
	err := s.ValidateValidateProviderOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "llm-connectivity/llm-provider-validations", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityProviderValidationResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateProvider creates an LLM provider. The connection is tested before the provider is
// persisted. Requires system administration.
//
// API endpoint: POST /api/v2/llm-connectivity/llm-providers.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *LlmConnectivityService) CreateProvider(ctx context.Context, opt *LlmConnectivityCreateProviderOptions) (*LlmConnectivityLlmProviderResponse, *http.Response, error) {
	err := s.ValidateCreateProviderOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "llm-connectivity/llm-providers", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LlmConnectivityLlmProviderResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
