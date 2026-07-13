package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixSuggestionsService_CreateSuggestion(t *testing.T) {
	response := FixSuggestion{
		Id:          "suggestion-1",
		IssueId:     "issue-123",
		Explanation: "Fix the null pointer dereference",
		Changes:     []FixSuggestionChange{{StartLine: 10, EndLine: 12, NewCode: "if obj != nil {"}},
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/fix-suggestions/ai-suggestions", http.StatusOK,
		map[string]any{"issueId": "issue-123"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.CreateSuggestion(context.Background(), &FixSuggestionsCreateOptions{
		IssueId: "issue-123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "suggestion-1", result.Id)
}

func TestFixSuggestionsService_CreateSuggestion_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.FixSuggestions.CreateSuggestion(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.FixSuggestions.CreateSuggestion(context.Background(), &FixSuggestionsCreateOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestFixSuggestionsService_GetEnablement(t *testing.T) {
	response := FixSuggestionsFeatureEnablement{
		Enablement:         "ENABLED_FOR_ALL_PROJECTS",
		EnabledProjectKeys: []string{"proj-1"},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/feature-enablements", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.GetEnablement(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ENABLED_FOR_ALL_PROJECTS", result.Enablement)
}

func TestFixSuggestionsService_SetEnablement(t *testing.T) {
	server := newTestServer(t, mockPatchHandler(t, "/v2/fix-suggestions/feature-enablements", http.StatusNoContent,
		map[string]any{"enablement": FixSuggestionsEnablementEnabledForAllProjects}, nil))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.FixSuggestions.SetEnablement(context.Background(), &FixSuggestionsSetEnablementOptions{
		Enablement: FixSuggestionsEnablementEnabledForAllProjects,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestFixSuggestionsService_SetEnablement_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.FixSuggestions.SetEnablement(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.FixSuggestions.SetEnablement(context.Background(), &FixSuggestionsSetEnablementOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.FixSuggestions.SetEnablement(context.Background(), &FixSuggestionsSetEnablementOptions{
		Enablement: "NOT_A_VALID_STATE",
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestFixSuggestionsService_AwarenessBannerInteraction(t *testing.T) {
	response := FixSuggestionsAwarenessBanner{Id: "interaction-1"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/fix-suggestions/feature-enablements/awareness-banner-interactions",
		http.StatusOK, map[string]any{"bannerType": FixSuggestionsBannerTypeEnable}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.AwarenessBannerInteraction(context.Background(), &FixSuggestionsAwarenessBannerOptions{
		BannerType: FixSuggestionsBannerTypeEnable,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "interaction-1", result.Id)
}

func TestFixSuggestionsService_AwarenessBannerInteraction_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.FixSuggestions.AwarenessBannerInteraction(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.FixSuggestions.AwarenessBannerInteraction(context.Background(), &FixSuggestionsAwarenessBannerOptions{
		BannerType: "INFO",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestFixSuggestionsService_GetIssueAvailability(t *testing.T) {
	response := FixSuggestionIssueAvailability{IssueId: "issue-123", AiSuggestion: "AVAILABLE"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/issues/issue-123", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.GetIssueAvailability(context.Background(), &FixSuggestionsIssueOptions{
		IssueId: "issue-123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "AVAILABLE", result.AiSuggestion)
}

func TestFixSuggestionsService_GetIssueAvailability_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.FixSuggestions.GetIssueAvailability(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestFixSuggestionsService_GetServiceInfo(t *testing.T) {
	response := FixSuggestionsServiceInfo{Status: "UP", IsEnabled: true, SubscriptionType: "PAID"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/service-info", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.GetServiceInfo(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "UP", result.Status)
	assert.True(t, result.IsEnabled)
}

func TestFixSuggestionsService_GetSubscriptionType(t *testing.T) {
	response := FixSuggestionsSubscriptionType{SubscriptionType: "PAID"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/service-info/subscription-type", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.GetSubscriptionType(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "PAID", result.SubscriptionType)
}

func TestFixSuggestionsService_GetSupportedRules(t *testing.T) {
	response := FixSuggestionsSupportedRules{Rules: []string{"java:S1234", "python:S5678"}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/supported-rules", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.GetSupportedRules(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []string{"java:S1234", "python:S5678"}, result.Rules)
}

func TestFixSuggestionsService_GetSupportedLlmProviders(t *testing.T) {
	// The real endpoint returns a bare JSON array (verified live against a
	// SonarQube 2025.2 Enterprise instance's OpenAPI schema and server
	// bytecode), not an object wrapping a "providers" field.
	response := []FixSuggestionsLlmProvider{
		{Key: "openai", Name: "OpenAI", SelfHosted: false},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/fix-suggestions/supported-llm-providers", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.FixSuggestions.GetSupportedLlmProviders(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
}
