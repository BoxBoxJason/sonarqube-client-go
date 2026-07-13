package sonar

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// GetApplicationConfiguration
// =============================================================================

func TestAtlassianV2_GetApplicationConfiguration(t *testing.T) {
	response := AtlassianAuthenticationDetails{ClientId: "client-1"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/atlassian/application-configuration", http.StatusOK, response))
	client := newTestClient(t, server.url())
	svc := &AtlassianService{client: client}

	result, resp, err := svc.GetApplicationConfiguration(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "client-1", result.ClientId)
}

// =============================================================================
// CreateOrUpdateApplicationConfiguration
// =============================================================================

func TestAtlassianV2_CreateOrUpdateApplicationConfiguration(t *testing.T) {
	response := AtlassianAuthenticationDetails{ClientId: "client-1"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/atlassian/application-configuration", http.StatusOK,
		&AtlassianAuthenticationConfigureOptions{
			ClientId: "client-1",
			Secret:   "super-secret",
		}, response))
	client := newTestClient(t, server.url())
	svc := &AtlassianService{client: client}

	result, resp, err := svc.CreateOrUpdateApplicationConfiguration(context.Background(), &AtlassianAuthenticationConfigureOptions{
		ClientId: "client-1",
		Secret:   "super-secret",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "client-1", result.ClientId)
}

func TestAtlassianV2_CreateOrUpdateApplicationConfiguration_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &AtlassianService{client: client}

	tests := []struct {
		opt  *AtlassianAuthenticationConfigureOptions
		name string
	}{
		{nil, "nil opt"},
		{&AtlassianAuthenticationConfigureOptions{Secret: "s"}, "missing ClientId"},
		{&AtlassianAuthenticationConfigureOptions{ClientId: "c"}, "missing Secret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.CreateOrUpdateApplicationConfiguration(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// GetAuthURL
// =============================================================================

func TestAtlassianV2_GetAuthURL(t *testing.T) {
	authURL := "https://auth.atlassian.com/authorize?client_id=abc"

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "unexpected HTTP method")
		assert.Equal(t, "/v2/atlassian/auth-url", r.URL.Path, "unexpected URL path")
		assert.Equal(t, "my-org", r.URL.Query().Get("sonarOrganizationKey"))
		assert.Equal(t, "11111111-1111-1111-1111-111111111111", r.URL.Query().Get("sonarOrganizationUuid"))
		// The API spec declares this endpoint's 200 response as
		// "application/json" with a string schema, so the client must not
		// request text/plain (see atlassianAuthURLResponse doc comment).
		assert.Equal(t, "application/json", r.Header.Get("Accept"), "GetAuthURL must not request text/plain")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(authURL))
	})
	client := newTestClient(t, server.URL)
	svc := &AtlassianService{client: client}

	result, resp, err := svc.GetAuthURL(context.Background(), &AtlassianAuthURLOptions{
		SonarOrganizationKey:  "my-org",
		SonarOrganizationUuid: "11111111-1111-1111-1111-111111111111",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, authURL, *result)
}

func TestAtlassianV2_GetAuthURL_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &AtlassianService{client: client}

	tests := []struct {
		opt  *AtlassianAuthURLOptions
		name string
	}{
		{nil, "nil opt"},
		{&AtlassianAuthURLOptions{SonarOrganizationUuid: "11111111-1111-1111-1111-111111111111"}, "missing SonarOrganizationKey"},
		{&AtlassianAuthURLOptions{SonarOrganizationKey: "my-org"}, "missing SonarOrganizationUuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.GetAuthURL(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}
