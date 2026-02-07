package sonar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_DefaultBaseURL(t *testing.T) {
	client, err := NewClient(nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, defaultBaseURL, client.BaseURL().String())
}

func TestNewClient_WithToken(t *testing.T) {
	client, err := NewClient(nil, WithToken("token123"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, privateToken, client.authType)
	assert.Equal(t, "token123", client.token)
}

func TestNewClient_WithBasicAuth(t *testing.T) {
	client, err := NewClient(nil, WithBasicAuth("user", "pass"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, basicAuth, client.authType)
	assert.Equal(t, "user", client.username)
	assert.Equal(t, "pass", client.password)
}

func TestNewClient_WithBaseURL(t *testing.T) {
	client, err := NewClient(nil, WithBaseURL("http://example.com/api/"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "http://example.com/api/", client.BaseURL().String())
}

func TestNewClient_WithBaseURL_NoTrailingSlash(t *testing.T) {
	client, err := NewClient(nil, WithBaseURL("http://example.com/api"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "http://example.com/api/", client.BaseURL().String())
}

func TestNewClient_WithCreateOptions(t *testing.T) {
	url := "http://example.com/api/"
	token := "my-token"

	client, err := NewClient(&ClientCreateOption{
		URL:   &url,
		Token: &token,
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "http://example.com/api/", client.BaseURL().String())
	assert.Equal(t, "my-token", client.token)
	assert.Equal(t, privateToken, client.authType)
}

func TestNewClient_ServicesInitialized(t *testing.T) {
	client, err := NewClient(nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify all services are initialized
	assert.NotNil(t, client.AlmIntegrations)
	assert.NotNil(t, client.AlmSettings)
	assert.NotNil(t, client.AnalysisCache)
	assert.NotNil(t, client.Authentication)
	assert.NotNil(t, client.Issues)
	assert.NotNil(t, client.Projects)
	assert.NotNil(t, client.Qualitygates)
	assert.NotNil(t, client.Qualityprofiles)
	assert.NotNil(t, client.Rules)
	assert.NotNil(t, client.Users)
}
