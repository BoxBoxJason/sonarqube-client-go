package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// FixSuggestionsEnablementEnabledForAllProjects represents the "ENABLED_FOR_ALL_PROJECTS" feature enablement state (enabled for all projects).
	FixSuggestionsEnablementEnabledForAllProjects = "ENABLED_FOR_ALL_PROJECTS"
	// FixSuggestionsEnablementDisabled represents the "DISABLED" feature enablement state.
	FixSuggestionsEnablementDisabled = "DISABLED"
	// FixSuggestionsEnablementEnabledForSomeProjects represents the "ENABLED_FOR_SOME_PROJECTS" feature enablement state.
	FixSuggestionsEnablementEnabledForSomeProjects = "ENABLED_FOR_SOME_PROJECTS"

	// FixSuggestionsBannerTypeEnable represents the "ENABLE" awareness banner interaction type.
	FixSuggestionsBannerTypeEnable = "ENABLE"
	// FixSuggestionsBannerTypeLearnMore represents the "LEARN_MORE" awareness banner interaction type.
	FixSuggestionsBannerTypeLearnMore = "LEARN_MORE"
)

// FixSuggestionsService handles communication with the AI fix suggestions V2 API endpoints.
// This service is only available in Enterprise Edition with AI features enabled.
type FixSuggestionsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedFixSuggestionsEnablements is the set of supported feature enablement states.
	allowedFixSuggestionsEnablements = map[string]struct{}{
		FixSuggestionsEnablementEnabledForAllProjects:  {},
		FixSuggestionsEnablementDisabled:               {},
		FixSuggestionsEnablementEnabledForSomeProjects: {},
	}

	// allowedFixSuggestionsBannerTypes is the set of supported awareness banner interaction types.
	allowedFixSuggestionsBannerTypes = map[string]struct{}{
		FixSuggestionsBannerTypeEnable:    {},
		FixSuggestionsBannerTypeLearnMore: {},
	}
)

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// FixSuggestionChange represents a code change proposed by an AI fix suggestion.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type FixSuggestionChange struct {
	// StartLine is the starting line of the code to replace.
	StartLine int `json:"startLine,omitempty"`
	// EndLine is the ending line of the code to replace.
	EndLine int `json:"endLine,omitempty"`
	// NewCode is the replacement code proposed by the AI.
	NewCode string `json:"newCode,omitempty"`
}

// FixSuggestion represents an AI-generated fix suggestion for an issue.
type FixSuggestion struct {
	// Id is the unique identifier of the fix suggestion.
	Id string `json:"id,omitempty"`
	// IssueId is the identifier of the issue this suggestion addresses.
	IssueId string `json:"issueId,omitempty"`
	// Explanation describes the reasoning behind the proposed fix.
	Explanation string `json:"explanation,omitempty"`
	// Changes is the list of code changes to apply.
	Changes []FixSuggestionChange `json:"changes,omitempty"`
}

// FixSuggestionProvider represents an LLM provider configuration for fix suggestions.
type FixSuggestionProvider struct {
	// Key is the provider key.
	Key string `json:"key,omitempty"`
	// ModelKey is the model key used.
	ModelKey string `json:"modelKey,omitempty"`
	// Endpoint is the provider endpoint URL.
	Endpoint string `json:"endpoint,omitempty"`
}

// FixSuggestionsFeatureEnablement represents the AI CodeFix feature enablement state.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type FixSuggestionsFeatureEnablement struct {
	// Enablement is the enablement state.
	Enablement string `json:"enablement,omitempty"`
	// EnabledProjectKeys lists project keys where the feature is enabled.
	EnabledProjectKeys []string `json:"enabledProjectKeys,omitempty"`
	// Provider contains the LLM provider configuration.
	Provider FixSuggestionProvider `json:"provider,omitzero"`
}

// FixSuggestionsAwarenessBanner represents the response from a banner interaction.
type FixSuggestionsAwarenessBanner struct {
	// Id is the interaction identifier.
	Id string `json:"id,omitempty"`
}

// FixSuggestionIssueAvailability represents fix suggestion availability for an issue.
type FixSuggestionIssueAvailability struct {
	// IssueId is the issue identifier.
	IssueId string `json:"issueId,omitempty"`
	// AiSuggestion indicates whether an AI suggestion is available.
	AiSuggestion string `json:"aiSuggestion,omitempty"`
}

// FixSuggestionsServiceInfo represents the AI CodeFix service status.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type FixSuggestionsServiceInfo struct {
	// Status is the current service status.
	Status string `json:"status,omitempty"`
	// IsEnabled indicates whether the service is enabled.
	IsEnabled bool `json:"isEnabled,omitempty"`
	// SubscriptionType is the subscription type.
	SubscriptionType string `json:"subscriptionType,omitempty"`
}

// FixSuggestionsSubscriptionType represents the subscription type response.
type FixSuggestionsSubscriptionType struct {
	// SubscriptionType is the subscription type string.
	SubscriptionType string `json:"subscriptionType,omitempty"`
}

// FixSuggestionsLlmModel represents a supported LLM model.
type FixSuggestionsLlmModel struct {
	// Key is the model key.
	Key string `json:"key,omitempty"`
	// Name is the display name.
	Name string `json:"name,omitempty"`
	// Recommended indicates whether this is the recommended model.
	Recommended bool `json:"recommended,omitempty"`
}

// FixSuggestionsSupportedRules represents the list of rules for which a fix
// suggestion can be generated.
type FixSuggestionsSupportedRules struct {
	// Rules is the list of supported rule keys.
	Rules []string `json:"rules,omitempty"`
}

// FixSuggestionsLlmProvider represents a supported LLM provider.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type FixSuggestionsLlmProvider struct {
	// Key is the provider key.
	Key string `json:"key,omitempty"`
	// Name is the display name.
	Name string `json:"name,omitempty"`
	// SelfHosted indicates whether this is a self-hosted provider.
	SelfHosted bool `json:"selfHosted,omitempty"`
	// Models is the list of models available from this provider.
	Models []FixSuggestionsLlmModel `json:"models,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// FixSuggestionsCreateOptions contains parameters for the CreateSuggestion method.
type FixSuggestionsCreateOptions struct {
	// IssueId is the issue identifier to generate a fix for. This field is required.
	IssueId string `json:"issueId"`
}

// FixSuggestionsSetEnablementOptions contains parameters for the SetEnablement method.
type FixSuggestionsSetEnablementOptions struct {
	// Enablement is the desired feature enablement state. This field is required.
	// Allowed values: ENABLED_FOR_ALL_PROJECTS, DISABLED, ENABLED_FOR_SOME_PROJECTS.
	Enablement string `json:"enablement"`
}

// FixSuggestionsAwarenessBannerOptions contains parameters for the AwarenessBannerInteraction method.
type FixSuggestionsAwarenessBannerOptions struct {
	// BannerType is the type of banner interaction that occurred. This field is required.
	// Allowed values: ENABLE, LEARN_MORE.
	BannerType string `json:"bannerType"`
}

// FixSuggestionsIssueOptions contains parameters for the GetIssueAvailability method.
type FixSuggestionsIssueOptions struct {
	// IssueId is the issue identifier. This field is required.
	IssueId string `json:"issueId"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCreateSuggestionOpt validates the options for the CreateSuggestion method.
func (s *FixSuggestionsService) ValidateCreateSuggestionOpt(opt *FixSuggestionsCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.IssueId, "IssueId")
}

// ValidateSetEnablementOpt validates the options for the SetEnablement method.
func (s *FixSuggestionsService) ValidateSetEnablementOpt(opt *FixSuggestionsSetEnablementOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Enablement, "Enablement")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Enablement, allowedFixSuggestionsEnablements, "Enablement")
}

// ValidateAwarenessBannerOpt validates the options for the AwarenessBannerInteraction method.
func (s *FixSuggestionsService) ValidateAwarenessBannerOpt(opt *FixSuggestionsAwarenessBannerOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.BannerType, "BannerType")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.BannerType, allowedFixSuggestionsBannerTypes, "BannerType")
}

// ValidateGetIssueAvailabilityOpt validates the options for the GetIssueAvailability method.
func (s *FixSuggestionsService) ValidateGetIssueAvailabilityOpt(opt *FixSuggestionsIssueOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.IssueId, "IssueId")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CreateSuggestion generates an AI fix suggestion for the given issue.
// Requires 'Browse' permission on the project.
//
// API endpoint: POST /api/v2/fix-suggestions/ai-suggestions.
// Enterprise Edition only.
func (s *FixSuggestionsService) CreateSuggestion(ctx context.Context, opt *FixSuggestionsCreateOptions) (*FixSuggestion, *http.Response, error) {
	err := s.ValidateCreateSuggestionOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "fix-suggestions/ai-suggestions", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestion)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetEnablement returns the AI CodeFix feature enablement configuration.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/fix-suggestions/feature-enablements.
// Enterprise Edition only.
func (s *FixSuggestionsService) GetEnablement(ctx context.Context) (*FixSuggestionsFeatureEnablement, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/feature-enablements", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionsFeatureEnablement)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetEnablement enables or disables the AI CodeFix feature.
// Requires 'Administer System' permission.
//
// API endpoint: PATCH /api/v2/fix-suggestions/feature-enablements.
// Enterprise Edition only.
func (s *FixSuggestionsService) SetEnablement(ctx context.Context, opt *FixSuggestionsSetEnablementOptions) (*http.Response, error) {
	err := s.ValidateSetEnablementOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "fix-suggestions/feature-enablements", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return s.client.Do(req, nil)
}

// AwarenessBannerInteraction records a user interaction with the AI awareness banner.
// Requires authentication.
//
// API endpoint: POST /api/v2/fix-suggestions/feature-enablements/awareness-banner-interactions.
// Enterprise Edition only.
func (s *FixSuggestionsService) AwarenessBannerInteraction(ctx context.Context, opt *FixSuggestionsAwarenessBannerOptions) (*FixSuggestionsAwarenessBanner, *http.Response, error) {
	err := s.ValidateAwarenessBannerOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "fix-suggestions/feature-enablements/awareness-banner-interactions", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionsAwarenessBanner)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetIssueAvailability checks whether an AI fix suggestion is available for the given issue.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/fix-suggestions/issues/{issueId}.
// Enterprise Edition only.
func (s *FixSuggestionsService) GetIssueAvailability(ctx context.Context, opt *FixSuggestionsIssueOptions) (*FixSuggestionIssueAvailability, *http.Response, error) {
	err := s.ValidateGetIssueAvailabilityOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/issues/"+url.PathEscape(opt.IssueId), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionIssueAvailability)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetServiceInfo returns the AI CodeFix service status and subscription information.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/fix-suggestions/service-info.
// Enterprise Edition only.
func (s *FixSuggestionsService) GetServiceInfo(ctx context.Context) (*FixSuggestionsServiceInfo, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/service-info", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionsServiceInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSubscriptionType returns the AI CodeFix subscription type.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/fix-suggestions/service-info/subscription-type.
// Enterprise Edition only.
//
// Deprecated: This endpoint has been removed from the SonarQube API as of version 2026.3 and will return an error if called.
func (s *FixSuggestionsService) GetSubscriptionType(ctx context.Context) (*FixSuggestionsSubscriptionType, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/service-info/subscription-type", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionsSubscriptionType)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSupportedLlmProviders returns the list of supported LLM providers.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/fix-suggestions/supported-llm-providers.
// The endpoint returns a bare JSON array (confirmed against the SonarQube
// server implementation), not an object wrapping a "providers" field.
// Enterprise Edition only.
func (s *FixSuggestionsService) GetSupportedLlmProviders(ctx context.Context) ([]FixSuggestionsLlmProvider, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/supported-llm-providers", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []FixSuggestionsLlmProvider

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSupportedRules returns the list of rules for which a fix suggestion can be generated.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/fix-suggestions/supported-rules.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *FixSuggestionsService) GetSupportedRules(ctx context.Context) (*FixSuggestionsSupportedRules, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/supported-rules", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionsSupportedRules)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
