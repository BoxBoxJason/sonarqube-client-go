package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDopTranslationService_CheckPermissions(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/permission-checks", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.CheckPermissions(context.Background(), &DopTranslationCheckPermissionsOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDopTranslationService_GenerateScmAccessToken(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/dop-translation/scm-access-tokens", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.GenerateScmAccessToken(context.Background(), &DopTranslationGenerateScmAccessTokenOptions{Project: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDopTranslationService_GenerateScmAccessToken_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DopTranslation.GenerateScmAccessToken(context.Background(), nil)
	require.Error(t, err)
}
