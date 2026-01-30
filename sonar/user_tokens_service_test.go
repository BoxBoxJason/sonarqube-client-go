package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserTokens_Generate(t *testing.T) {
	response := `{
		"login": "admin",
		"name": "my-token",
		"token": "secret-token-value",
		"createdAt": "2024-01-01T00:00:00+0000",
		"type": "USER_TOKEN"
	}`

	handler := mockHandler(t, http.MethodPost, "/user_tokens/generate", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &UserTokensGenerateOption{
		Name: "my-token",
	}

	result, resp, err := client.UserTokens.Generate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "my-token", result.Name)
	assert.Equal(t, "secret-token-value", result.Token)
	assert.Equal(t, "USER_TOKEN", result.Type)
}

func TestUserTokens_Generate_WithType(t *testing.T) {
	response := `{
		"login": "admin",
		"name": "project-token",
		"token": "project-token-value",
		"createdAt": "2024-01-01T00:00:00+0000",
		"type": "PROJECT_ANALYSIS_TOKEN"
	}`

	handler := mockHandler(t, http.MethodPost, "/user_tokens/generate", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &UserTokensGenerateOption{
		Name:       "project-token",
		Type:       "PROJECT_ANALYSIS_TOKEN",
		ProjectKey: "my-project",
	}

	result, resp, err := client.UserTokens.Generate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "PROJECT_ANALYSIS_TOKEN", result.Type)
}

func TestUserTokens_Generate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.UserTokens.Generate(nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, _, err = client.UserTokens.Generate(&UserTokensGenerateOption{})
	assert.Error(t, err)

	// Invalid Type should fail validation.
	_, _, err = client.UserTokens.Generate(&UserTokensGenerateOption{
		Name: "my-token",
		Type: "INVALID_TYPE",
	})
	assert.Error(t, err)

	// PROJECT_ANALYSIS_TOKEN without ProjectKey should fail validation.
	_, _, err = client.UserTokens.Generate(&UserTokensGenerateOption{
		Name: "my-token",
		Type: "PROJECT_ANALYSIS_TOKEN",
	})
	assert.Error(t, err)
}

func TestUserTokens_Revoke(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/user_tokens/revoke", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &UserTokensRevokeOption{
		Name: "my-token",
	}

	resp, err := client.UserTokens.Revoke(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserTokens_Revoke_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.UserTokens.Revoke(nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, err = client.UserTokens.Revoke(&UserTokensRevokeOption{})
	assert.Error(t, err)
}

func TestUserTokens_Search(t *testing.T) {
	response := `{
		"login": "admin",
		"userTokens": [
			{
				"name": "token1",
				"type": "USER_TOKEN",
				"createdAt": "2024-01-01T00:00:00+0000",
				"isExpired": false
			},
			{
				"name": "token2",
				"type": "GLOBAL_ANALYSIS_TOKEN",
				"createdAt": "2024-01-02T00:00:00+0000",
				"isExpired": true
			}
		]
	}`

	handler := mockHandler(t, http.MethodGet, "/user_tokens/search", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.UserTokens.Search(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "admin", result.Login)
	assert.Len(t, result.UserTokens, 2)
	assert.Equal(t, "token1", result.UserTokens[0].Name)
	assert.True(t, result.UserTokens[1].IsExpired)
}

func TestUserTokens_Search_WithLogin(t *testing.T) {
	response := `{"login": "testuser", "userTokens": []}`

	handler := mockHandler(t, http.MethodGet, "/user_tokens/search", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &UserTokensSearchOption{
		Login: "testuser",
	}

	result, _, err := client.UserTokens.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Login)
}

func TestUserTokens_ValidateGenerateOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
		Name: "my-token",
	})
	assert.NoError(t, err)

	// All valid token types should pass.
	validTypes := []string{"USER_TOKEN", "GLOBAL_ANALYSIS_TOKEN"}
	for _, tokenType := range validTypes {
		err := client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
			Name: "my-token",
			Type: tokenType,
		})
		assert.NoError(t, err, "expected nil error for type '%s'", tokenType)
	}

	// PROJECT_ANALYSIS_TOKEN with ProjectKey should pass.
	err = client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
		Name:       "project-token",
		Type:       "PROJECT_ANALYSIS_TOKEN",
		ProjectKey: "my-project",
	})
	assert.NoError(t, err)

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxTokenNameLength+1; i++ {
		longName += "a"
	}
	err = client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
		Name: longName,
	})
	assert.Error(t, err)
}
