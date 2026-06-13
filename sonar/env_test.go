package sonar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromEnv_Token(t *testing.T) {
	t.Setenv(EnvURL, "https://sonar.example.com/api/")
	t.Setenv(EnvToken, "tok-123")

	client, err := NewClientFromEnv()
	require.NoError(t, err)

	assert.Equal(t, privateToken, client.authType)
	assert.Equal(t, "tok-123", client.token)
	assert.Equal(t, "sonar.example.com", client.baseURL.Host)
}

func TestNewClientFromEnv_BasicAuth(t *testing.T) {
	t.Setenv(EnvURL, "https://sonar.example.com/api/")
	t.Setenv(EnvUsername, "alice")
	t.Setenv(EnvPassword, "s3cret")

	client, err := NewClientFromEnv()
	require.NoError(t, err)

	assert.Equal(t, basicAuth, client.authType)
	assert.Equal(t, "alice", client.username)
	assert.Equal(t, "s3cret", client.password)
}

func TestNewClientFromEnv_TokenTakesPrecedence(t *testing.T) {
	t.Setenv(EnvToken, "tok-123")
	t.Setenv(EnvUsername, "alice")
	t.Setenv(EnvPassword, "s3cret")

	client, err := NewClientFromEnv()
	require.NoError(t, err)

	assert.Equal(t, privateToken, client.authType, "token should take precedence over basic auth")
	assert.Empty(t, client.username)
}

func TestNewClientFromEnv_UserAgentOverride(t *testing.T) {
	t.Setenv(EnvUserAgent, "my-app/1.0")

	client, err := NewClientFromEnv()
	require.NoError(t, err)

	assert.Equal(t, "my-app/1.0", client.userAgent)
}

func TestNewClientFromEnv_OptionsOverrideEnv(t *testing.T) {
	t.Setenv(EnvURL, "https://from-env.example.com/api/")

	client, err := NewClientFromEnv(WithBaseURL("https://from-option.example.com/api/"))
	require.NoError(t, err)

	assert.Equal(t, "from-option.example.com", client.baseURL.Host, "functional options should override the environment")
}

func TestNewClientFromEnv_Timeout(t *testing.T) {
	t.Setenv(EnvTimeout, "90s")

	client, err := NewClientFromEnv()
	require.NoError(t, err)

	assert.Equal(t, 90*time.Second, client.timeout)
}

func TestNewClientFromEnv_TimeoutInvalid(t *testing.T) {
	t.Setenv(EnvTimeout, "not-a-duration")

	_, err := NewClientFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), EnvTimeout)
}

func TestNewClientFromEnv_NoEnvUsesDefaults(t *testing.T) {
	t.Setenv(EnvURL, "")
	t.Setenv(EnvToken, "")
	t.Setenv(EnvUsername, "")
	t.Setenv(EnvPassword, "")

	client, err := NewClientFromEnv()
	require.NoError(t, err)
	require.NotNil(t, client.baseURL)
	assert.Equal(t, "localhost:9000", client.baseURL.Host)
}
