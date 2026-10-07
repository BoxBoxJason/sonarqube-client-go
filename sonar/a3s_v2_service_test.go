package sonar

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestA3sService_GetAnalysisStats(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/analysis-stats/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.GetAnalysisStats(context.Background(), "test-id", &A3sGetAnalysisStatsOptions{Period: "LAST_DAY"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetAnalysisStats_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.GetAnalysisStats(context.Background(), "", nil)
	require.Error(t, err)
}

func TestA3sService_GetPublicCollectionConfig(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/collection-config", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.GetPublicCollectionConfig(context.Background(), &A3sGetPublicCollectionConfigOptions{OrganizationKey: "test", ProjectKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetPublicCollectionConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.GetPublicCollectionConfig(context.Background(), nil)
	require.Error(t, err)
}

func TestA3sService_GetPublicOrgEntitlement(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/org-entitlement/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.GetPublicOrgEntitlement(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetPublicOrgEntitlement_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.GetPublicOrgEntitlement(context.Background(), "")
	require.Error(t, err)
}

func TestA3sService_GetPrivateOrgEntitlement(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/private/a3s-analysis/org-entitlement/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.GetPrivateOrgEntitlement(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetPrivateOrgEntitlement_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.GetPrivateOrgEntitlement(context.Background(), "")
	require.Error(t, err)
}

func TestA3sService_GetPrivateCollectionConfig(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/private/a3s-frontend-context/collection-config", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.GetPrivateCollectionConfig(context.Background(), &A3sGetPrivateCollectionConfigOptions{OrganizationKey: "test", ProjectKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetPrivateCollectionConfig_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.GetPrivateCollectionConfig(context.Background(), nil)
	require.Error(t, err)
}

func TestA3sService_SearchPrivateContexts(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/private/a3s-frontend-context/contexts", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.SearchPrivateContexts(context.Background(), &A3sSearchPrivateContextsOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetPrivateContext(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/a3s/private/a3s-frontend-context/contexts/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.GetPrivateContext(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_GetPrivateContext_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.GetPrivateContext(context.Background(), "")
	require.Error(t, err)
}

func TestA3sService_CreatePublicAnalysis(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/a3s/analyses", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.CreatePublicAnalysis(context.Background(), "test-value", &A3sCreatePublicAnalysisOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_CreatePublicAnalysis_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.CreatePublicAnalysis(context.Background(), "", nil)
	require.Error(t, err)
}

func TestA3sService_CreatePublicContext(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/a3s/contexts", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.CreatePublicContext(context.Background(), &A3sCreatePublicContextOptions{AnalysisId: "test", Kind: "test", Metadata: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_CreatePublicContext_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.CreatePublicContext(context.Background(), nil)
	require.Error(t, err)
}

func TestA3sService_CreatePrivateContext(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/a3s/private/a3s-frontend-context/contexts", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.CreatePrivateContext(context.Background(), &A3sCreatePrivateContextOptions{AnalysisId: "test", Kind: "test", Metadata: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_CreatePrivateContext_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.CreatePrivateContext(context.Background(), nil)
	require.Error(t, err)
}

func TestA3sService_CreatePrivateAnalysis(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/a3s/private/analyses", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.A3s.CreatePrivateAnalysis(context.Background(), &A3sCreatePrivateAnalysisOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestA3sService_CreatePrivateAnalysis_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.A3s.CreatePrivateAnalysis(context.Background(), nil)
	require.Error(t, err)
}

func TestA3sService_UploadContext(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPut, "/v2/a3s/context-uploads/test-id", http.StatusOK))
	client := newTestClient(t, server.url())

	resp, err := client.V2.A3s.UploadContext(context.Background(), "test-id", strings.NewReader("payload"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
