package sonar

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// CountBinding Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CountBinding(t *testing.T) {
	response := &AlmSettingsCountBinding{
		Key:      "my-alm-setting",
		Projects: 5,
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/count_binding", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCountBindingOptions{
		AlmSetting: "my-alm-setting",
	}

	result, resp, err := client.AlmSettings.CountBinding(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "my-alm-setting", result.Key)
	assert.Equal(t, int64(5), result.Projects)
}

func TestAlmSettings_CountBinding_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmSettings.CountBinding(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AlmSetting
	_, _, err = client.AlmSettings.CountBinding(context.Background(), &AlmSettingsCountBindingOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// CreateAzure Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateAzure(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/create_azure", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCreateAzureOptions{
		Key:                 "my-azure-setting",
		PersonalAccessToken: "my-pat",
		URL:                 "https://dev.azure.com/myorg",
	}

	resp, err := client.AlmSettings.CreateAzure(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_CreateAzure_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.CreateAzure(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.CreateAzure(context.Background(), &AlmSettingsCreateAzureOptions{
		PersonalAccessToken: "pat",
		URL:                 "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test Key too long
	_, err = client.AlmSettings.CreateAzure(context.Background(), &AlmSettingsCreateAzureOptions{
		Key:                 strings.Repeat("a", MaxAlmKeyLength+1),
		PersonalAccessToken: "pat",
		URL:                 "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test missing PersonalAccessToken
	_, err = client.AlmSettings.CreateAzure(context.Background(), &AlmSettingsCreateAzureOptions{
		Key: "my-key",
		URL: "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test PersonalAccessToken too long
	_, err = client.AlmSettings.CreateAzure(context.Background(), &AlmSettingsCreateAzureOptions{
		Key:                 "my-key",
		PersonalAccessToken: strings.Repeat("a", MaxPersonalAccessTokenLength+1),
		URL:                 "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.CreateAzure(context.Background(), &AlmSettingsCreateAzureOptions{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
	})
	assert.Error(t, err)

	// Test URL too long
	_, err = client.AlmSettings.CreateAzure(context.Background(), &AlmSettingsCreateAzureOptions{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
		URL:                 strings.Repeat("a", MaxAlmURLLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// CreateBitbucket Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateBitbucket(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/create_bitbucket", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCreateBitbucketOptions{
		Key:                 "my-bitbucket-setting",
		PersonalAccessToken: "my-pat",
		URL:                 "https://bitbucket.example.com",
	}

	resp, err := client.AlmSettings.CreateBitbucket(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_CreateBitbucket_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.CreateBitbucket(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.CreateBitbucket(context.Background(), &AlmSettingsCreateBitbucketOptions{
		PersonalAccessToken: "pat",
		URL:                 "https://bitbucket.example.com",
	})
	assert.Error(t, err)

	// Test Key too long
	_, err = client.AlmSettings.CreateBitbucket(context.Background(), &AlmSettingsCreateBitbucketOptions{
		Key:                 strings.Repeat("a", MaxAlmKeyLength+1),
		PersonalAccessToken: "pat",
		URL:                 "https://bitbucket.example.com",
	})
	assert.Error(t, err)

	// Test missing PersonalAccessToken
	_, err = client.AlmSettings.CreateBitbucket(context.Background(), &AlmSettingsCreateBitbucketOptions{
		Key: "my-key",
		URL: "https://bitbucket.example.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.CreateBitbucket(context.Background(), &AlmSettingsCreateBitbucketOptions{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// CreateBitbucketCloud Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateBitbucketCloud(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/create_bitbucketcloud", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCreateBitbucketCloudOptions{
		ClientID:     "my-client-id",
		ClientSecret: "my-client-secret",
		Key:          "my-bitbucket-cloud-setting",
		Workspace:    "my-workspace",
	}

	resp, err := client.AlmSettings.CreateBitbucketCloud(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_CreateBitbucketCloud_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.CreateBitbucketCloud(context.Background(), nil)
	assert.Error(t, err)

	// Test missing ClientID
	_, err = client.AlmSettings.CreateBitbucketCloud(context.Background(), &AlmSettingsCreateBitbucketCloudOptions{
		ClientSecret: "secret",
		Key:          "my-key",
		Workspace:    "workspace",
	})
	assert.Error(t, err)

	// Test ClientID too long
	_, err = client.AlmSettings.CreateBitbucketCloud(context.Background(), &AlmSettingsCreateBitbucketCloudOptions{
		ClientID:     strings.Repeat("a", MaxBitbucketCloudClientIDLength+1),
		ClientSecret: "secret",
		Key:          "my-key",
		Workspace:    "workspace",
	})
	assert.Error(t, err)

	// Test missing ClientSecret
	_, err = client.AlmSettings.CreateBitbucketCloud(context.Background(), &AlmSettingsCreateBitbucketCloudOptions{
		ClientID:  "client-id",
		Key:       "my-key",
		Workspace: "workspace",
	})
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.CreateBitbucketCloud(context.Background(), &AlmSettingsCreateBitbucketCloudOptions{
		ClientID:     "client-id",
		ClientSecret: "secret",
		Workspace:    "workspace",
	})
	assert.Error(t, err)

	// Test missing Workspace
	_, err = client.AlmSettings.CreateBitbucketCloud(context.Background(), &AlmSettingsCreateBitbucketCloudOptions{
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// CreateGithub Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateGithub(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/create_github", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     "my-client-id",
		ClientSecret: "my-client-secret",
		Key:          "my-github-setting",
		PrivateKey:   "my-private-key",
		URL:          "https://api.github.com",
	}

	resp, err := client.AlmSettings.CreateGithub(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_CreateGithub_WithOptionalWebhookSecret(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/create_github", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCreateGithubOptions{
		AppID:         "12345",
		ClientID:      "my-client-id",
		ClientSecret:  "my-client-secret",
		Key:           "my-github-setting",
		PrivateKey:    "my-private-key",
		URL:           "https://api.github.com",
		WebhookSecret: "my-webhook-secret",
	}

	resp, err := client.AlmSettings.CreateGithub(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_CreateGithub_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.CreateGithub(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AppID
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test AppID too long
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        strings.Repeat("a", MaxGitHubAppIDLength+1),
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing ClientID
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test ClientID too long
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     strings.Repeat("a", MaxGitHubClientIDLength+1),
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing ClientSecret
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:      "12345",
		ClientID:   "client-id",
		Key:        "my-key",
		PrivateKey: "private-key",
		URL:        "https://api.github.com",
	})
	assert.Error(t, err)

	// Test ClientSecret too long
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: strings.Repeat("a", MaxGitHubClientSecretLength+1),
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing PrivateKey
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test PrivateKey too long
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   strings.Repeat("a", MaxGitHubPrivateKeyLength+1),
		URL:          "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
	})
	assert.Error(t, err)

	// Test WebhookSecret too long
	_, err = client.AlmSettings.CreateGithub(context.Background(), &AlmSettingsCreateGithubOptions{
		AppID:         "12345",
		ClientID:      "client-id",
		ClientSecret:  "secret",
		Key:           "my-key",
		PrivateKey:    "private-key",
		URL:           "https://api.github.com",
		WebhookSecret: strings.Repeat("a", MaxGitHubWebhookSecretLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// CreateGitlab Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateGitlab(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/create_gitlab", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsCreateGitlabOptions{
		Key:                 "my-gitlab-setting",
		PersonalAccessToken: "my-pat",
		URL:                 "https://gitlab.example.com",
	}

	resp, err := client.AlmSettings.CreateGitlab(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_CreateGitlab_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.CreateGitlab(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.CreateGitlab(context.Background(), &AlmSettingsCreateGitlabOptions{
		PersonalAccessToken: "pat",
		URL:                 "https://gitlab.example.com",
	})
	assert.Error(t, err)

	// Test missing PersonalAccessToken
	_, err = client.AlmSettings.CreateGitlab(context.Background(), &AlmSettingsCreateGitlabOptions{
		Key: "my-key",
		URL: "https://gitlab.example.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.CreateGitlab(context.Background(), &AlmSettingsCreateGitlabOptions{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Delete Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/delete", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsDeleteOptions{
		Key: "my-alm-setting",
	}

	resp, err := client.AlmSettings.Delete(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.Delete(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.Delete(context.Background(), &AlmSettingsDeleteOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// GetBinding Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_GetBinding(t *testing.T) {
	response := &AlmSettingsGetBinding{
		Alm:                   "github",
		Key:                   "my-github-setting",
		Monorepo:              false,
		Repository:            "my-org/my-repo",
		SummaryCommentEnabled: true,
		URL:                   "https://api.github.com",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/get_binding", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsGetBindingOptions{
		Project: "my-project",
	}

	result, resp, err := client.AlmSettings.GetBinding(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "github", result.Alm)
	assert.Equal(t, "my-github-setting", result.Key)
	assert.True(t, result.SummaryCommentEnabled)
}

func TestAlmSettings_GetBinding_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmSettings.GetBinding(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Project
	_, _, err = client.AlmSettings.GetBinding(context.Background(), &AlmSettingsGetBindingOptions{})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// List Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_List(t *testing.T) {
	response := &AlmSettingsList{
		AlmSettings: []AlmSetting{
			{Alm: "github", Key: "github-setting", URL: "https://api.github.com"},
			{Alm: "azure", Key: "azure-setting", URL: "https://dev.azure.com"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.AlmSettings.List(context.Background(), &AlmSettingsListOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.AlmSettings, 2)
	assert.Equal(t, "github", result.AlmSettings[0].Alm)
}

func TestAlmSettings_List_WithProject(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"almSettings": []}`))
	})
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsListOptions{
		Project: "my-project",
	}

	_, _, err := client.AlmSettings.List(context.Background(), opt)
	require.NoError(t, err)
}

func TestAlmSettings_List_NilOption(t *testing.T) {
	response := &AlmSettingsList{AlmSettings: []AlmSetting{}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, _, err := client.AlmSettings.List(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

// -----------------------------------------------------------------------------
// ListDefinitions Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_ListDefinitions(t *testing.T) {
	response := &AlmSettingsListDefinitions{
		Azure:          []AzureDefinition{{Key: "azure-setting", URL: "https://dev.azure.com"}},
		Bitbucket:      []BitbucketDefinition{},
		BitbucketCloud: []BitbucketCloudDefinition{},
		Github:         []GithubDefinition{{AppID: "12345", ClientID: "client-id", Key: "github-setting", URL: "https://api.github.com"}},
		Gitlab:         []GitlabDefinition{},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/list_definitions", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.AlmSettings.ListDefinitions(context.Background(), )
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Azure, 1)
	assert.Len(t, result.Github, 1)
	assert.Equal(t, "12345", result.Github[0].AppID)
}

// -----------------------------------------------------------------------------
// UpdateAzure Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateAzure(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_azure", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateAzureOptions{
		Key: "my-azure-setting",
		URL: "https://dev.azure.com",
	}

	resp, err := client.AlmSettings.UpdateAzure(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateAzure_WithOptionalFields(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_azure", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateAzureOptions{
		Key:                 "my-azure-setting",
		NewKey:              "new-azure-setting",
		PersonalAccessToken: "new-pat",
		URL:                 "https://dev.azure.com",
	}

	resp, err := client.AlmSettings.UpdateAzure(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateAzure_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.UpdateAzure(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.UpdateAzure(context.Background(), &AlmSettingsUpdateAzureOptions{
		URL: "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test Key too long
	_, err = client.AlmSettings.UpdateAzure(context.Background(), &AlmSettingsUpdateAzureOptions{
		Key: strings.Repeat("a", MaxAlmKeyLength+1),
		URL: "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test NewKey too long
	_, err = client.AlmSettings.UpdateAzure(context.Background(), &AlmSettingsUpdateAzureOptions{
		Key:    "my-key",
		NewKey: strings.Repeat("a", MaxAlmKeyLength+1),
		URL:    "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test PersonalAccessToken too long
	_, err = client.AlmSettings.UpdateAzure(context.Background(), &AlmSettingsUpdateAzureOptions{
		Key:                 "my-key",
		PersonalAccessToken: strings.Repeat("a", MaxPersonalAccessTokenLength+1),
		URL:                 "https://dev.azure.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.UpdateAzure(context.Background(), &AlmSettingsUpdateAzureOptions{
		Key: "my-key",
	})
	assert.Error(t, err)

	// Test URL too long
	_, err = client.AlmSettings.UpdateAzure(context.Background(), &AlmSettingsUpdateAzureOptions{
		Key: "my-key",
		URL: strings.Repeat("a", MaxAlmURLLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// UpdateBitbucket Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateBitbucket(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_bitbucket", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateBitbucketOptions{
		Key: "my-bitbucket-setting",
		URL: "https://bitbucket.example.com",
	}

	resp, err := client.AlmSettings.UpdateBitbucket(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateBitbucket_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.UpdateBitbucket(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.UpdateBitbucket(context.Background(), &AlmSettingsUpdateBitbucketOptions{
		URL: "https://bitbucket.example.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.UpdateBitbucket(context.Background(), &AlmSettingsUpdateBitbucketOptions{
		Key: "my-key",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// UpdateBitbucketCloud Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateBitbucketCloud(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_bitbucketcloud", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateBitbucketCloudOptions{
		ClientID:  "my-client-id",
		Key:       "my-bitbucket-cloud-setting",
		Workspace: "my-workspace",
	}

	resp, err := client.AlmSettings.UpdateBitbucketCloud(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateBitbucketCloud_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.UpdateBitbucketCloud(context.Background(), nil)
	assert.Error(t, err)

	// Test missing ClientID
	_, err = client.AlmSettings.UpdateBitbucketCloud(context.Background(), &AlmSettingsUpdateBitbucketCloudOptions{
		Key:       "my-key",
		Workspace: "workspace",
	})
	assert.Error(t, err)

	// Test ClientID too long
	_, err = client.AlmSettings.UpdateBitbucketCloud(context.Background(), &AlmSettingsUpdateBitbucketCloudOptions{
		ClientID:  strings.Repeat("a", MaxBitbucketCloudClientIDUpdateLength+1),
		Key:       "my-key",
		Workspace: "workspace",
	})
	assert.Error(t, err)

	// Test ClientSecret too long
	_, err = client.AlmSettings.UpdateBitbucketCloud(context.Background(), &AlmSettingsUpdateBitbucketCloudOptions{
		ClientID:     "client-id",
		ClientSecret: strings.Repeat("a", MaxBitbucketCloudClientSecretUpdateLength+1),
		Key:          "my-key",
		Workspace:    "workspace",
	})
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.UpdateBitbucketCloud(context.Background(), &AlmSettingsUpdateBitbucketCloudOptions{
		ClientID:  "client-id",
		Workspace: "workspace",
	})
	assert.Error(t, err)

	// Test missing Workspace
	_, err = client.AlmSettings.UpdateBitbucketCloud(context.Background(), &AlmSettingsUpdateBitbucketCloudOptions{
		ClientID: "client-id",
		Key:      "my-key",
	})
	assert.Error(t, err)

	// Test Workspace too long
	_, err = client.AlmSettings.UpdateBitbucketCloud(context.Background(), &AlmSettingsUpdateBitbucketCloudOptions{
		ClientID:  "client-id",
		Key:       "my-key",
		Workspace: strings.Repeat("a", MaxBitbucketCloudWorkspaceUpdateLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// UpdateGithub Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateGithub(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_github", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateGithubOptions{
		AppID:    "12345",
		ClientID: "my-client-id",
		Key:      "my-github-setting",
		URL:      "https://api.github.com",
	}

	resp, err := client.AlmSettings.UpdateGithub(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateGithub_WithOptionalFields(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_github", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateGithubOptions{
		AppID:         "12345",
		ClientID:      "my-client-id",
		ClientSecret:  "new-client-secret",
		Key:           "my-github-setting",
		NewKey:        "new-github-setting",
		PrivateKey:    "new-private-key",
		URL:           "https://api.github.com",
		WebhookSecret: "new-webhook-secret",
	}

	resp, err := client.AlmSettings.UpdateGithub(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateGithub_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.UpdateGithub(context.Background(), nil)
	assert.Error(t, err)

	// Test missing AppID
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		ClientID: "client-id",
		Key:      "my-key",
		URL:      "https://api.github.com",
	})
	assert.Error(t, err)

	// Test AppID too long
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		AppID:    strings.Repeat("a", MaxGitHubAppIDLength+1),
		ClientID: "client-id",
		Key:      "my-key",
		URL:      "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing ClientID
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		AppID: "12345",
		Key:   "my-key",
		URL:   "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		AppID:    "12345",
		ClientID: "client-id",
		URL:      "https://api.github.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		AppID:    "12345",
		ClientID: "client-id",
		Key:      "my-key",
	})
	assert.Error(t, err)

	// Test PrivateKey too long
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		AppID:      "12345",
		ClientID:   "client-id",
		Key:        "my-key",
		PrivateKey: strings.Repeat("a", MaxGitHubPrivateKeyLength+1),
		URL:        "https://api.github.com",
	})
	assert.Error(t, err)

	// Test WebhookSecret too long
	_, err = client.AlmSettings.UpdateGithub(context.Background(), &AlmSettingsUpdateGithubOptions{
		AppID:         "12345",
		ClientID:      "client-id",
		Key:           "my-key",
		URL:           "https://api.github.com",
		WebhookSecret: strings.Repeat("a", MaxGitHubWebhookSecretLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// UpdateGitlab Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateGitlab(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/alm_settings/update_gitlab", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsUpdateGitlabOptions{
		Key: "my-gitlab-setting",
		URL: "https://gitlab.example.com",
	}

	resp, err := client.AlmSettings.UpdateGitlab(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAlmSettings_UpdateGitlab_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.AlmSettings.UpdateGitlab(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.AlmSettings.UpdateGitlab(context.Background(), &AlmSettingsUpdateGitlabOptions{
		URL: "https://gitlab.example.com",
	})
	assert.Error(t, err)

	// Test missing URL
	_, err = client.AlmSettings.UpdateGitlab(context.Background(), &AlmSettingsUpdateGitlabOptions{
		Key: "my-key",
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Validate Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_Validate(t *testing.T) {
	response := &AlmSettingsValidation{Errors: []AlmValidationError{}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/validate", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsValidateOptions{
		Key: "my-alm-setting",
	}

	result, resp, err := client.AlmSettings.Validate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Empty(t, result.Errors)
}

func TestAlmSettings_Validate_WithErrors(t *testing.T) {
	response := &AlmSettingsValidation{
		Errors: []AlmValidationError{
			{Msg: "Invalid token"},
			{Msg: "Connection refused"},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/alm_settings/validate", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &AlmSettingsValidateOptions{
		Key: "my-alm-setting",
	}

	result, _, err := client.AlmSettings.Validate(context.Background(), opt)
	require.NoError(t, err)
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "Invalid token", result.Errors[0].Msg)
}

func TestAlmSettings_Validate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.AlmSettings.Validate(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, _, err = client.AlmSettings.Validate(context.Background(), &AlmSettingsValidateOptions{})
	assert.Error(t, err)

	// Test Key too long
	_, _, err = client.AlmSettings.Validate(context.Background(), &AlmSettingsValidateOptions{
		Key: strings.Repeat("a", MaxAlmKeyLength+1),
	})
	assert.Error(t, err)
}

// -----------------------------------------------------------------------------
// Validate Function Tests (direct method tests)
// -----------------------------------------------------------------------------

func TestAlmSettings_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil options should be valid
	err := client.AlmSettings.ValidateListOpt(nil)
	assert.NoError(t, err)

	// Empty options should be valid
	err = client.AlmSettings.ValidateListOpt(&AlmSettingsListOptions{})
	assert.NoError(t, err)
}
