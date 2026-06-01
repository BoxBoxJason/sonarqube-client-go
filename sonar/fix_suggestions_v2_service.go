package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// FixSuggestionsService handles communication with the Fix Suggestions related
// methods of the SonarQube V2 API.
type FixSuggestionsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// FixSuggestionChange represents a code change proposed by a fix suggestion.
//
//nolint:govet // Field alignment less important than matching API field order.
type FixSuggestionChange struct {
	// StartLine is the first line affected by the change.
	StartLine int32 `json:"startLine,omitempty"`
	// EndLine is the last line affected by the change.
	EndLine int32 `json:"endLine,omitempty"`
	// NewCode is the suggested replacement code.
	NewCode string `json:"newCode,omitempty"`
}

// FixSuggestionResponse represents a generated fix suggestion.
type FixSuggestionResponse struct {
	// Id is the unique identifier of the fix suggestion.
	Id string `json:"id,omitempty"`
	// IssueId is the issue key for which the suggestion was generated.
	IssueId string `json:"issueId,omitempty"`
	// Explanation describes the suggested fix.
	Explanation string `json:"explanation,omitempty"`
	// Changes contains the proposed code changes.
	Changes []FixSuggestionChange `json:"changes,omitempty"`
}

// FixSuggestionIssueResponse represents suggestion availability for an issue.
type FixSuggestionIssueResponse struct {
	// IssueId is the issue key.
	IssueId string `json:"issueId,omitempty"`
	// AiSuggestion is the availability status for generating a suggestion.
	AiSuggestion string `json:"aiSuggestion,omitempty"`
}

// ProviderResponse represents the configured provider for fix suggestions.
type ProviderResponse struct {
	// Key is the provider key.
	Key string `json:"key,omitempty"`
	// ModelKey is the configured model key.
	ModelKey string `json:"modelKey,omitempty"`
	// Endpoint is the configured provider endpoint.
	Endpoint string `json:"endpoint,omitempty"`
}

// FeatureEnablementResponse represents the fix suggestions feature state.
//
//nolint:govet // Field alignment less important than matching API field order.
type FeatureEnablementResponse struct {
	// Enablement is the instance/project-level enablement state.
	Enablement string `json:"enablement,omitempty"`
	// EnabledProjectKeys contains project keys for which the feature is enabled.
	EnabledProjectKeys []string `json:"enabledProjectKeys,omitempty"`
	// Provider contains the configured provider, if any.
	Provider *ProviderResponse `json:"provider,omitempty"`
}

// AwarenessBannerClickedResponse represents a recorded awareness banner interaction.
type AwarenessBannerClickedResponse struct {
	// Id is the recorded interaction identifier.
	Id string `json:"id,omitempty"`
}

// LLMModel represents a supported model for a fix suggestions provider.
type LLMModel struct {
	// Key is the model key.
	Key string `json:"key,omitempty"`
	// Name is the display name of the model.
	Name string `json:"name,omitempty"`
	// Recommended indicates whether the model is recommended.
	Recommended bool `json:"recommended,omitempty"`
}

// LLMProviderResponse represents a supported fix suggestions provider.
//
//nolint:govet // Field alignment less important than matching API field order.
type LLMProviderResponse struct {
	// Key is the provider key.
	Key string `json:"key,omitempty"`
	// Name is the display name of the provider.
	Name string `json:"name,omitempty"`
	// SelfHosted indicates whether the provider is self-hosted.
	SelfHosted bool `json:"selfHosted,omitempty"`
	// Models contains supported provider models.
	Models []LLMModel `json:"models,omitempty"`
}

// FixSuggestionsServiceInfo represents status and subscription information.
//
//nolint:govet // Field alignment less important than matching API field order.
type FixSuggestionsServiceInfo struct {
	// Status is the service connectivity status.
	Status string `json:"status,omitempty"`
	// IsEnabled indicates whether fix suggestions are enabled.
	IsEnabled bool `json:"isEnabled,omitempty"`
	// SubscriptionType is the current subscription type.
	SubscriptionType string `json:"subscriptionType,omitempty"`
}

// SubscriptionTypeResponse represents the current fix suggestions subscription type.
type SubscriptionTypeResponse struct {
	// SubscriptionType is the current subscription type.
	SubscriptionType string `json:"subscriptionType,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// FixSuggestionPostRequest contains parameters for generating a fix suggestion.
type FixSuggestionPostRequest struct {
	// IssueId is the issue key for which to generate a fix suggestion.
	// This field is required.
	IssueId string `json:"issueId"`
}

// FeatureEnablementRequest contains parameters for updating feature enablement.
type FeatureEnablementRequest struct {
	// Enablement indicates whether the feature should be enabled.
	Enablement bool `json:"enablement"`
}

// AwarenessBannerClickedRequest contains a recorded awareness banner interaction.
type AwarenessBannerClickedRequest struct {
	// BannerType is the awareness banner interaction type.
	// Allowed values: ENABLE, LEARN_MORE.
	BannerType string `json:"bannerType"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateGenerateFixSuggestionRequest validates the FixSuggestionPostRequest.
func (s *FixSuggestionsService) ValidateGenerateFixSuggestionRequest(opt *FixSuggestionPostRequest) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.IssueId, "IssueId")
}

// ValidateUpdateFeatureEnablementRequest validates the FeatureEnablementRequest.
func (s *FixSuggestionsService) ValidateUpdateFeatureEnablementRequest(opt *FeatureEnablementRequest) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return nil
}

// ValidateAwarenessBannerClickedRequest validates the AwarenessBannerClickedRequest.
func (s *FixSuggestionsService) ValidateAwarenessBannerClickedRequest(opt *AwarenessBannerClickedRequest) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.BannerType, "BannerType")
	if err != nil {
		return err
	}

	switch opt.BannerType {
	case "ENABLE", "LEARN_MORE":
		return nil
	default:
		return NewValidationError("BannerType", "must be one of: ENABLE, LEARN_MORE", ErrInvalidValue)
	}
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GenerateFixSuggestion suggests a fix for the given issue.
func (s *FixSuggestionsService) GenerateFixSuggestion(ctx context.Context, opt *FixSuggestionPostRequest) (*FixSuggestionResponse, *http.Response, error) {
	err := s.ValidateGenerateFixSuggestionRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "fix-suggestions/ai-suggestions", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetIssueSuggestionAvailability fetches fix suggestion availability for an issue.
func (s *FixSuggestionsService) GetIssueSuggestionAvailability(ctx context.Context, issueID string) (*FixSuggestionIssueResponse, *http.Response, error) {
	err := ValidateRequired(issueID, "IssueId")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/issues/"+url.PathEscape(issueID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FixSuggestionIssueResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetFeatureEnablement fetches fix suggestions feature enablement configuration.
func (s *FixSuggestionsService) GetFeatureEnablement(ctx context.Context) (*FeatureEnablementResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/feature-enablements", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(FeatureEnablementResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateFeatureEnablement enables or disables fix suggestions.
func (s *FixSuggestionsService) UpdateFeatureEnablement(ctx context.Context, opt *FeatureEnablementRequest) (*http.Response, error) {
	err := s.ValidateUpdateFeatureEnablementRequest(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "fix-suggestions/feature-enablements", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RecordAwarenessBannerInteraction records a fix suggestions awareness banner interaction.
func (s *FixSuggestionsService) RecordAwarenessBannerInteraction(ctx context.Context, opt *AwarenessBannerClickedRequest) (*AwarenessBannerClickedResponse, *http.Response, error) {
	err := s.ValidateAwarenessBannerClickedRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "fix-suggestions/feature-enablements/awareness-banner-interactions", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AwarenessBannerClickedResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListSupportedLLMProviders lists the supported fix suggestions LLM providers.
func (s *FixSuggestionsService) ListSupportedLLMProviders(ctx context.Context) ([]LLMProviderResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/supported-llm-providers", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []LLMProviderResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetServiceInfo requests status and subscription information.
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

// GetSubscriptionType requests fix suggestions subscription information.
func (s *FixSuggestionsService) GetSubscriptionType(ctx context.Context) (*SubscriptionTypeResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "fix-suggestions/service-info/subscription-type", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(SubscriptionTypeResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
