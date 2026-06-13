package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixSuggestionsV2_GenerateFixSuggestion(t *testing.T) {
	response := FixSuggestionResponse{
		Id:          "550e8400-e29b-41d4-a716-446655440000",
		IssueId:     "ISSUE-1",
		Explanation: "Replace the unsafe call.",
		Changes: []FixSuggestionChange{
			{StartLine: 10, EndLine: 12, NewCode: "safeCall()"},
		},
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/fix-suggestions/ai-suggestions", http.StatusOK,
		&FixSuggestionPostRequest{IssueId: "ISSUE-1"}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.GenerateFixSuggestion(context.Background(), &FixSuggestionPostRequest{IssueId: "ISSUE-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ISSUE-1", result.IssueId)
	require.Len(t, result.Changes, 1)
	assert.Equal(t, int32(10), result.Changes[0].StartLine)
}

func TestFixSuggestionsV2_GenerateFixSuggestion_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *FixSuggestionPostRequest
	}{
		{"nil opt", nil},
		{"missing issue id", &FixSuggestionPostRequest{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.FixSuggestions.GenerateFixSuggestion(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestFixSuggestionsV2_GetIssueSuggestionAvailability(t *testing.T) {
	response := FixSuggestionIssueResponse{
		IssueId:      "ISSUE-1",
		AiSuggestion: "AVAILABLE",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/issues/ISSUE-1", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.GetIssueSuggestionAvailability(context.Background(), "ISSUE-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "AVAILABLE", result.AiSuggestion)
}

func TestFixSuggestionsV2_GetIssueSuggestionAvailability_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.FixSuggestions.GetIssueSuggestionAvailability(context.Background(), "")
	assert.Error(t, err)
}

func TestFixSuggestionsV2_GetFeatureEnablement(t *testing.T) {
	response := FeatureEnablementResponse{
		Enablement:         "ENABLED_FOR_SOME_PROJECTS",
		EnabledProjectKeys: []string{"project-a", "project-b"},
		Provider:           &ProviderResponse{Key: "provider-1", ModelKey: "model-1", Endpoint: "https://example.com"},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/feature-enablements", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.GetFeatureEnablement(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ENABLED_FOR_SOME_PROJECTS", result.Enablement)
	assert.Equal(t, []string{"project-a", "project-b"}, result.EnabledProjectKeys)
	require.NotNil(t, result.Provider)
	assert.Equal(t, "provider-1", result.Provider.Key)
}

func TestFixSuggestionsV2_UpdateFeatureEnablement(t *testing.T) {
	server := newTestServer(t, mockPatchHandler(t, "/v2/fix-suggestions/feature-enablements", http.StatusNoContent,
		&FeatureEnablementRequest{Enablement: true}, nil))
	client := newTestClient(t, server.url())

	resp, err := client.V2.FixSuggestions.UpdateFeatureEnablement(context.Background(), &FeatureEnablementRequest{Enablement: true})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestFixSuggestionsV2_UpdateFeatureEnablement_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.FixSuggestions.UpdateFeatureEnablement(context.Background(), nil)
	assert.Error(t, err)
}

func TestFixSuggestionsV2_RecordAwarenessBannerInteraction(t *testing.T) {
	response := AwarenessBannerClickedResponse{Id: "interaction-1"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/fix-suggestions/feature-enablements/awareness-banner-interactions", http.StatusOK,
		&AwarenessBannerClickedRequest{BannerType: "ENABLE"}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.RecordAwarenessBannerInteraction(context.Background(), &AwarenessBannerClickedRequest{BannerType: "ENABLE"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "interaction-1", result.Id)
}

func TestFixSuggestionsV2_RecordAwarenessBannerInteraction_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *AwarenessBannerClickedRequest
	}{
		{"nil opt", nil},
		{"missing banner type", &AwarenessBannerClickedRequest{}},
		{"invalid banner type", &AwarenessBannerClickedRequest{BannerType: "OTHER"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.FixSuggestions.RecordAwarenessBannerInteraction(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestFixSuggestionsV2_ListSupportedLLMProviders(t *testing.T) {
	response := []LLMProviderResponse{
		{
			Key:        "provider-1",
			Name:       "Provider One",
			SelfHosted: false,
			Models:     []LLMModel{{Key: "model-1", Name: "Model One", Recommended: true}},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/supported-llm-providers", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.ListSupportedLLMProviders(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result, 1)
	assert.Equal(t, "provider-1", result[0].Key)
	require.Len(t, result[0].Models, 1)
	assert.True(t, result[0].Models[0].Recommended)
}

func TestFixSuggestionsV2_GetServiceInfo(t *testing.T) {
	response := FixSuggestionsServiceInfo{
		Status:           "SUCCESS",
		IsEnabled:        true,
		SubscriptionType: "PAID",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/service-info", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.GetServiceInfo(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "SUCCESS", result.Status)
	assert.True(t, result.IsEnabled)
	assert.Equal(t, "PAID", result.SubscriptionType)
}

func TestFixSuggestionsV2_GetSubscriptionType(t *testing.T) {
	response := SubscriptionTypeResponse{SubscriptionType: "EARLY_ACCESS"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/service-info/subscription-type", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.FixSuggestions.GetSubscriptionType(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "EARLY_ACCESS", result.SubscriptionType)
}
