package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLlmConnectivityService_DeleteProvider(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/llm-connectivity/llm-providers/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.LlmConnectivity.DeleteProvider(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestLlmConnectivityService_DeleteProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.LlmConnectivity.DeleteProvider(context.Background(), "")
	require.Error(t, err)
}

func TestLlmConnectivityService_GetProviderDefinitions(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/llm-connectivity/llm-provider-definitions", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.GetProviderDefinitions(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_GetProviderMappings(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/llm-connectivity/llm-provider-mappings", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.GetProviderMappings(context.Background(), &LlmConnectivityGetProviderMappingsOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_GetProviderMappings_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.GetProviderMappings(context.Background(), nil)
	require.Error(t, err)
}

func TestLlmConnectivityService_ListProviders(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/llm-connectivity/llm-providers", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.ListProviders(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_GetProvider(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/llm-connectivity/llm-providers/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.GetProvider(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_GetProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.GetProvider(context.Background(), "")
	require.Error(t, err)
}

func TestLlmConnectivityService_GetSupportedModels(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/llm-connectivity/supported-models", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.GetSupportedModels(context.Background(), &LlmConnectivityGetSupportedModelsOptions{AiCapability: "AI_CODEFIX"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_GetSupportedModels_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.GetSupportedModels(context.Background(), nil)
	require.Error(t, err)
}

func TestLlmConnectivityService_UpdateProvider(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/llm-connectivity/llm-providers/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.UpdateProvider(context.Background(), "test-id", &LlmConnectivityUpdateProviderOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_UpdateProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.UpdateProvider(context.Background(), "", nil)
	require.Error(t, err)
}

func TestLlmConnectivityService_UpsertProviderMapping(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/llm-connectivity/llm-provider-mappings", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.UpsertProviderMapping(context.Background(), &LlmConnectivityUpsertProviderMappingOptions{AiCapability: "AI_CODEFIX", LlmProviderId: "test", ModelIdentifier: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_UpsertProviderMapping_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.UpsertProviderMapping(context.Background(), nil)
	require.Error(t, err)
}

func TestLlmConnectivityService_ValidateProvider(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/llm-connectivity/llm-provider-validations", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.ValidateProvider(context.Background(), &LlmConnectivityValidateProviderOptions{LlmProviderId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_ValidateProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.ValidateProvider(context.Background(), nil)
	require.Error(t, err)
}

func TestLlmConnectivityService_CreateProvider(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/llm-connectivity/llm-providers", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.LlmConnectivity.CreateProvider(context.Background(), &LlmConnectivityCreateProviderOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestLlmConnectivityService_CreateProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.LlmConnectivity.CreateProvider(context.Background(), nil)
	require.Error(t, err)
}
