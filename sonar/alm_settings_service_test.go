package sonargo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// CountBinding Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CountBinding(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/count_binding") {
			t.Errorf("expected path to contain alm_settings/count_binding, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"key":"my-alm-setting","projects":5}`))
	}))

	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCountBindingOption{
		AlmSetting: "my-alm-setting",
	}

	result, resp, err := client.AlmSettings.CountBinding(opt)
	if err != nil {
		t.Fatalf("CountBinding failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Key != "my-alm-setting" || result.Projects != 5 {
		t.Error("unexpected result from CountBinding")
	}
}

func TestAlmSettings_CountBinding_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmSettings.CountBinding(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmSettings.CountBinding(&AlmSettingsCountBindingOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}
}

// -----------------------------------------------------------------------------
// CreateAzure Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateAzure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/create_azure") {
			t.Errorf("expected path to contain alm_settings/create_azure, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCreateAzureOption{
		Key:                 "my-azure-setting",
		PersonalAccessToken: "my-pat",
		URL:                 "https://dev.azure.com/myorg",
	}

	resp, err := client.AlmSettings.CreateAzure(opt)
	if err != nil {
		t.Fatalf("CreateAzure failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateAzure_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.CreateAzure(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.CreateAzure(&AlmSettingsCreateAzureOption{
		PersonalAccessToken: "pat",
		URL:                 "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test Key too long
	_, err = client.AlmSettings.CreateAzure(&AlmSettingsCreateAzureOption{
		Key:                 strings.Repeat("a", MaxAlmKeyLength+1),
		PersonalAccessToken: "pat",
		URL:                 "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for Key exceeding max length")
	}

	// Test missing PersonalAccessToken
	_, err = client.AlmSettings.CreateAzure(&AlmSettingsCreateAzureOption{
		Key: "my-key",
		URL: "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for missing PersonalAccessToken")
	}

	// Test PersonalAccessToken too long
	_, err = client.AlmSettings.CreateAzure(&AlmSettingsCreateAzureOption{
		Key:                 "my-key",
		PersonalAccessToken: strings.Repeat("a", MaxPersonalAccessTokenLength+1),
		URL:                 "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for PersonalAccessToken exceeding max length")
	}

	// Test missing URL
	_, err = client.AlmSettings.CreateAzure(&AlmSettingsCreateAzureOption{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Test URL too long
	_, err = client.AlmSettings.CreateAzure(&AlmSettingsCreateAzureOption{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
		URL:                 strings.Repeat("a", MaxAlmURLLength+1),
	})
	if err == nil {
		t.Error("expected error for URL exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// CreateBitbucket Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateBitbucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/create_bitbucket") {
			t.Errorf("expected path to contain alm_settings/create_bitbucket, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCreateBitbucketOption{
		Key:                 "my-bitbucket-setting",
		PersonalAccessToken: "my-pat",
		URL:                 "https://bitbucket.example.com",
	}

	resp, err := client.AlmSettings.CreateBitbucket(opt)
	if err != nil {
		t.Fatalf("CreateBitbucket failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateBitbucket_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.CreateBitbucket(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.CreateBitbucket(&AlmSettingsCreateBitbucketOption{
		PersonalAccessToken: "pat",
		URL:                 "https://bitbucket.example.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test Key too long
	_, err = client.AlmSettings.CreateBitbucket(&AlmSettingsCreateBitbucketOption{
		Key:                 strings.Repeat("a", MaxAlmKeyLength+1),
		PersonalAccessToken: "pat",
		URL:                 "https://bitbucket.example.com",
	})
	if err == nil {
		t.Error("expected error for Key exceeding max length")
	}

	// Test missing PersonalAccessToken
	_, err = client.AlmSettings.CreateBitbucket(&AlmSettingsCreateBitbucketOption{
		Key: "my-key",
		URL: "https://bitbucket.example.com",
	})
	if err == nil {
		t.Error("expected error for missing PersonalAccessToken")
	}

	// Test missing URL
	_, err = client.AlmSettings.CreateBitbucket(&AlmSettingsCreateBitbucketOption{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

// -----------------------------------------------------------------------------
// CreateBitbucketCloud Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateBitbucketCloud(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/create_bitbucketcloud") {
			t.Errorf("expected path to contain alm_settings/create_bitbucketcloud, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCreateBitbucketCloudOption{
		ClientID:     "my-client-id",
		ClientSecret: "my-client-secret",
		Key:          "my-bitbucket-cloud-setting",
		Workspace:    "my-workspace",
	}

	resp, err := client.AlmSettings.CreateBitbucketCloud(opt)
	if err != nil {
		t.Fatalf("CreateBitbucketCloud failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateBitbucketCloud_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.CreateBitbucketCloud(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing ClientID
	_, err = client.AlmSettings.CreateBitbucketCloud(&AlmSettingsCreateBitbucketCloudOption{
		ClientSecret: "secret",
		Key:          "my-key",
		Workspace:    "workspace",
	})
	if err == nil {
		t.Error("expected error for missing ClientID")
	}

	// Test ClientID too long
	_, err = client.AlmSettings.CreateBitbucketCloud(&AlmSettingsCreateBitbucketCloudOption{
		ClientID:     strings.Repeat("a", MaxBitbucketCloudClientIDLength+1),
		ClientSecret: "secret",
		Key:          "my-key",
		Workspace:    "workspace",
	})
	if err == nil {
		t.Error("expected error for ClientID exceeding max length")
	}

	// Test missing ClientSecret
	_, err = client.AlmSettings.CreateBitbucketCloud(&AlmSettingsCreateBitbucketCloudOption{
		ClientID:  "client-id",
		Key:       "my-key",
		Workspace: "workspace",
	})
	if err == nil {
		t.Error("expected error for missing ClientSecret")
	}

	// Test missing Key
	_, err = client.AlmSettings.CreateBitbucketCloud(&AlmSettingsCreateBitbucketCloudOption{
		ClientID:     "client-id",
		ClientSecret: "secret",
		Workspace:    "workspace",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing Workspace
	_, err = client.AlmSettings.CreateBitbucketCloud(&AlmSettingsCreateBitbucketCloudOption{
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
	})
	if err == nil {
		t.Error("expected error for missing Workspace")
	}
}

// -----------------------------------------------------------------------------
// CreateGithub Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateGithub(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/create_github") {
			t.Errorf("expected path to contain alm_settings/create_github, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     "my-client-id",
		ClientSecret: "my-client-secret",
		Key:          "my-github-setting",
		PrivateKey:   "my-private-key",
		URL:          "https://api.github.com",
	}

	resp, err := client.AlmSettings.CreateGithub(opt)
	if err != nil {
		t.Fatalf("CreateGithub failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateGithub_WithOptionalWebhookSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCreateGithubOption{
		AppID:         "12345",
		ClientID:      "my-client-id",
		ClientSecret:  "my-client-secret",
		Key:           "my-github-setting",
		PrivateKey:    "my-private-key",
		URL:           "https://api.github.com",
		WebhookSecret: "my-webhook-secret",
	}

	resp, err := client.AlmSettings.CreateGithub(opt)
	if err != nil {
		t.Fatalf("CreateGithub failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateGithub_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.CreateGithub(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AppID
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing AppID")
	}

	// Test AppID too long
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        strings.Repeat("a", MaxGitHubAppIDLength+1),
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for AppID exceeding max length")
	}

	// Test missing ClientID
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing ClientID")
	}

	// Test ClientID too long
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     strings.Repeat("a", MaxGitHubClientIDLength+1),
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for ClientID exceeding max length")
	}

	// Test missing ClientSecret
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:      "12345",
		ClientID:   "client-id",
		Key:        "my-key",
		PrivateKey: "private-key",
		URL:        "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing ClientSecret")
	}

	// Test ClientSecret too long
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: strings.Repeat("a", MaxGitHubClientSecretLength+1),
		Key:          "my-key",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for ClientSecret exceeding max length")
	}

	// Test missing Key
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		PrivateKey:   "private-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing PrivateKey
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing PrivateKey")
	}

	// Test PrivateKey too long
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   strings.Repeat("a", MaxGitHubPrivateKeyLength+1),
		URL:          "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for PrivateKey exceeding max length")
	}

	// Test missing URL
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:        "12345",
		ClientID:     "client-id",
		ClientSecret: "secret",
		Key:          "my-key",
		PrivateKey:   "private-key",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Test WebhookSecret too long
	_, err = client.AlmSettings.CreateGithub(&AlmSettingsCreateGithubOption{
		AppID:         "12345",
		ClientID:      "client-id",
		ClientSecret:  "secret",
		Key:           "my-key",
		PrivateKey:    "private-key",
		URL:           "https://api.github.com",
		WebhookSecret: strings.Repeat("a", MaxGitHubWebhookSecretLength+1),
	})
	if err == nil {
		t.Error("expected error for WebhookSecret exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// CreateGitlab Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_CreateGitlab(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/create_gitlab") {
			t.Errorf("expected path to contain alm_settings/create_gitlab, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsCreateGitlabOption{
		Key:                 "my-gitlab-setting",
		PersonalAccessToken: "my-pat",
		URL:                 "https://gitlab.example.com",
	}

	resp, err := client.AlmSettings.CreateGitlab(opt)
	if err != nil {
		t.Fatalf("CreateGitlab failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateGitlab_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.CreateGitlab(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.CreateGitlab(&AlmSettingsCreateGitlabOption{
		PersonalAccessToken: "pat",
		URL:                 "https://gitlab.example.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing PersonalAccessToken
	_, err = client.AlmSettings.CreateGitlab(&AlmSettingsCreateGitlabOption{
		Key: "my-key",
		URL: "https://gitlab.example.com",
	})
	if err == nil {
		t.Error("expected error for missing PersonalAccessToken")
	}

	// Test missing URL
	_, err = client.AlmSettings.CreateGitlab(&AlmSettingsCreateGitlabOption{
		Key:                 "my-key",
		PersonalAccessToken: "pat",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

// -----------------------------------------------------------------------------
// Delete Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/delete") {
			t.Errorf("expected path to contain alm_settings/delete, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsDeleteOption{
		Key: "my-alm-setting",
	}

	resp, err := client.AlmSettings.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_Delete_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.Delete(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.Delete(&AlmSettingsDeleteOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}
}

// -----------------------------------------------------------------------------
// GetBinding Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_GetBinding(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/get_binding") {
			t.Errorf("expected path to contain alm_settings/get_binding, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"alm": "github",
			"key": "my-github-setting",
			"monorepo": false,
			"repository": "my-org/my-repo",
			"summaryCommentEnabled": true,
			"url": "https://api.github.com"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsGetBindingOption{
		Project: "my-project",
	}

	result, resp, err := client.AlmSettings.GetBinding(opt)
	if err != nil {
		t.Fatalf("GetBinding failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Alm != "github" || result.Key != "my-github-setting" {
		t.Error("unexpected result from GetBinding")
	}

	if !result.SummaryCommentEnabled {
		t.Error("expected SummaryCommentEnabled to be true")
	}
}

func TestAlmSettings_GetBinding_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmSettings.GetBinding(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Project
	_, _, err = client.AlmSettings.GetBinding(&AlmSettingsGetBindingOption{})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

// -----------------------------------------------------------------------------
// List Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/list") {
			t.Errorf("expected path to contain alm_settings/list, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"almSettings": [
				{"alm": "github", "key": "github-setting", "url": "https://api.github.com"},
				{"alm": "azure", "key": "azure-setting", "url": "https://dev.azure.com"}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.AlmSettings.List(&AlmSettingsListOption{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.AlmSettings) != 2 {
		t.Error("expected 2 alm settings")
	}

	if result.AlmSettings[0].Alm != "github" {
		t.Errorf("expected first setting to be github, got %s", result.AlmSettings[0].Alm)
	}
}

func TestAlmSettings_List_WithProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("project") != "my-project" {
			t.Errorf("expected project parameter, got %s", r.URL.Query().Get("project"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"almSettings": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsListOption{
		Project: "my-project",
	}

	_, _, err = client.AlmSettings.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestAlmSettings_List_NilOption(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"almSettings": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, _, err := client.AlmSettings.List(nil)
	if err != nil {
		t.Fatalf("List should accept nil option: %v", err)
	}

	if result == nil {
		t.Error("expected result, got nil")
	}
}

// -----------------------------------------------------------------------------
// ListDefinitions Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_ListDefinitions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/list_definitions") {
			t.Errorf("expected path to contain alm_settings/list_definitions, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"azure": [{"key": "azure-setting", "url": "https://dev.azure.com"}],
			"bitbucket": [],
			"bitbucketcloud": [],
			"github": [{"appId": "12345", "clientId": "client-id", "key": "github-setting", "url": "https://api.github.com"}],
			"gitlab": []
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.AlmSettings.ListDefinitions()
	if err != nil {
		t.Fatalf("ListDefinitions failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if len(result.Azure) != 1 {
		t.Errorf("expected 1 azure setting, got %d", len(result.Azure))
	}

	if len(result.Github) != 1 {
		t.Errorf("expected 1 github setting, got %d", len(result.Github))
	}

	if result.Github[0].AppID != "12345" {
		t.Errorf("expected AppID 12345, got %s", result.Github[0].AppID)
	}
}

// -----------------------------------------------------------------------------
// UpdateAzure Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateAzure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/update_azure") {
			t.Errorf("expected path to contain alm_settings/update_azure, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateAzureOption{
		Key: "my-azure-setting",
		URL: "https://dev.azure.com",
	}

	resp, err := client.AlmSettings.UpdateAzure(opt)
	if err != nil {
		t.Fatalf("UpdateAzure failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateAzure_WithOptionalFields(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateAzureOption{
		Key:                 "my-azure-setting",
		NewKey:              "new-azure-setting",
		PersonalAccessToken: "new-pat",
		URL:                 "https://dev.azure.com",
	}

	resp, err := client.AlmSettings.UpdateAzure(opt)
	if err != nil {
		t.Fatalf("UpdateAzure failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateAzure_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.UpdateAzure(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.UpdateAzure(&AlmSettingsUpdateAzureOption{
		URL: "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test Key too long
	_, err = client.AlmSettings.UpdateAzure(&AlmSettingsUpdateAzureOption{
		Key: strings.Repeat("a", MaxAlmKeyLength+1),
		URL: "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for Key exceeding max length")
	}

	// Test NewKey too long
	_, err = client.AlmSettings.UpdateAzure(&AlmSettingsUpdateAzureOption{
		Key:    "my-key",
		NewKey: strings.Repeat("a", MaxAlmKeyLength+1),
		URL:    "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for NewKey exceeding max length")
	}

	// Test PersonalAccessToken too long
	_, err = client.AlmSettings.UpdateAzure(&AlmSettingsUpdateAzureOption{
		Key:                 "my-key",
		PersonalAccessToken: strings.Repeat("a", MaxPersonalAccessTokenLength+1),
		URL:                 "https://dev.azure.com",
	})
	if err == nil {
		t.Error("expected error for PersonalAccessToken exceeding max length")
	}

	// Test missing URL
	_, err = client.AlmSettings.UpdateAzure(&AlmSettingsUpdateAzureOption{
		Key: "my-key",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Test URL too long
	_, err = client.AlmSettings.UpdateAzure(&AlmSettingsUpdateAzureOption{
		Key: "my-key",
		URL: strings.Repeat("a", MaxAlmURLLength+1),
	})
	if err == nil {
		t.Error("expected error for URL exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// UpdateBitbucket Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateBitbucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/update_bitbucket") {
			t.Errorf("expected path to contain alm_settings/update_bitbucket, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateBitbucketOption{
		Key: "my-bitbucket-setting",
		URL: "https://bitbucket.example.com",
	}

	resp, err := client.AlmSettings.UpdateBitbucket(opt)
	if err != nil {
		t.Fatalf("UpdateBitbucket failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateBitbucket_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.UpdateBitbucket(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.UpdateBitbucket(&AlmSettingsUpdateBitbucketOption{
		URL: "https://bitbucket.example.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing URL
	_, err = client.AlmSettings.UpdateBitbucket(&AlmSettingsUpdateBitbucketOption{
		Key: "my-key",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

// -----------------------------------------------------------------------------
// UpdateBitbucketCloud Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateBitbucketCloud(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/update_bitbucketcloud") {
			t.Errorf("expected path to contain alm_settings/update_bitbucketcloud, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateBitbucketCloudOption{
		ClientID:  "my-client-id",
		Key:       "my-bitbucket-cloud-setting",
		Workspace: "my-workspace",
	}

	resp, err := client.AlmSettings.UpdateBitbucketCloud(opt)
	if err != nil {
		t.Fatalf("UpdateBitbucketCloud failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateBitbucketCloud_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.UpdateBitbucketCloud(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing ClientID
	_, err = client.AlmSettings.UpdateBitbucketCloud(&AlmSettingsUpdateBitbucketCloudOption{
		Key:       "my-key",
		Workspace: "workspace",
	})
	if err == nil {
		t.Error("expected error for missing ClientID")
	}

	// Test ClientID too long
	_, err = client.AlmSettings.UpdateBitbucketCloud(&AlmSettingsUpdateBitbucketCloudOption{
		ClientID:  strings.Repeat("a", MaxBitbucketCloudClientIDUpdateLength+1),
		Key:       "my-key",
		Workspace: "workspace",
	})
	if err == nil {
		t.Error("expected error for ClientID exceeding max length")
	}

	// Test ClientSecret too long
	_, err = client.AlmSettings.UpdateBitbucketCloud(&AlmSettingsUpdateBitbucketCloudOption{
		ClientID:     "client-id",
		ClientSecret: strings.Repeat("a", MaxBitbucketCloudClientSecretUpdateLength+1),
		Key:          "my-key",
		Workspace:    "workspace",
	})
	if err == nil {
		t.Error("expected error for ClientSecret exceeding max length")
	}

	// Test missing Key
	_, err = client.AlmSettings.UpdateBitbucketCloud(&AlmSettingsUpdateBitbucketCloudOption{
		ClientID:  "client-id",
		Workspace: "workspace",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing Workspace
	_, err = client.AlmSettings.UpdateBitbucketCloud(&AlmSettingsUpdateBitbucketCloudOption{
		ClientID: "client-id",
		Key:      "my-key",
	})
	if err == nil {
		t.Error("expected error for missing Workspace")
	}

	// Test Workspace too long
	_, err = client.AlmSettings.UpdateBitbucketCloud(&AlmSettingsUpdateBitbucketCloudOption{
		ClientID:  "client-id",
		Key:       "my-key",
		Workspace: strings.Repeat("a", MaxBitbucketCloudWorkspaceUpdateLength+1),
	})
	if err == nil {
		t.Error("expected error for Workspace exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// UpdateGithub Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateGithub(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/update_github") {
			t.Errorf("expected path to contain alm_settings/update_github, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateGithubOption{
		AppID:    "12345",
		ClientID: "my-client-id",
		Key:      "my-github-setting",
		URL:      "https://api.github.com",
	}

	resp, err := client.AlmSettings.UpdateGithub(opt)
	if err != nil {
		t.Fatalf("UpdateGithub failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateGithub_WithOptionalFields(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateGithubOption{
		AppID:         "12345",
		ClientID:      "my-client-id",
		ClientSecret:  "new-client-secret",
		Key:           "my-github-setting",
		NewKey:        "new-github-setting",
		PrivateKey:    "new-private-key",
		URL:           "https://api.github.com",
		WebhookSecret: "new-webhook-secret",
	}

	resp, err := client.AlmSettings.UpdateGithub(opt)
	if err != nil {
		t.Fatalf("UpdateGithub failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateGithub_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.UpdateGithub(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AppID
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		ClientID: "client-id",
		Key:      "my-key",
		URL:      "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing AppID")
	}

	// Test AppID too long
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		AppID:    strings.Repeat("a", MaxGitHubAppIDLength+1),
		ClientID: "client-id",
		Key:      "my-key",
		URL:      "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for AppID exceeding max length")
	}

	// Test missing ClientID
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		AppID: "12345",
		Key:   "my-key",
		URL:   "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing ClientID")
	}

	// Test missing Key
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		AppID:    "12345",
		ClientID: "client-id",
		URL:      "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing URL
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		AppID:    "12345",
		ClientID: "client-id",
		Key:      "my-key",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}

	// Test PrivateKey too long
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		AppID:      "12345",
		ClientID:   "client-id",
		Key:        "my-key",
		PrivateKey: strings.Repeat("a", MaxGitHubPrivateKeyLength+1),
		URL:        "https://api.github.com",
	})
	if err == nil {
		t.Error("expected error for PrivateKey exceeding max length")
	}

	// Test WebhookSecret too long
	_, err = client.AlmSettings.UpdateGithub(&AlmSettingsUpdateGithubOption{
		AppID:         "12345",
		ClientID:      "client-id",
		Key:           "my-key",
		URL:           "https://api.github.com",
		WebhookSecret: strings.Repeat("a", MaxGitHubWebhookSecretLength+1),
	})
	if err == nil {
		t.Error("expected error for WebhookSecret exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// UpdateGitlab Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_UpdateGitlab(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/update_gitlab") {
			t.Errorf("expected path to contain alm_settings/update_gitlab, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsUpdateGitlabOption{
		Key: "my-gitlab-setting",
		URL: "https://gitlab.example.com",
	}

	resp, err := client.AlmSettings.UpdateGitlab(opt)
	if err != nil {
		t.Fatalf("UpdateGitlab failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateGitlab_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmSettings.UpdateGitlab(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.AlmSettings.UpdateGitlab(&AlmSettingsUpdateGitlabOption{
		URL: "https://gitlab.example.com",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing URL
	_, err = client.AlmSettings.UpdateGitlab(&AlmSettingsUpdateGitlabOption{
		Key: "my-key",
	})
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

// -----------------------------------------------------------------------------
// Validate Tests
// -----------------------------------------------------------------------------

func TestAlmSettings_Validate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_settings/validate") {
			t.Errorf("expected path to contain alm_settings/validate, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"errors": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsValidateOption{
		Key: "my-alm-setting",
	}

	result, resp, err := client.AlmSettings.Validate(opt)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Errors == nil {
		t.Error("expected empty errors array, got nil")
	}
}

func TestAlmSettings_Validate_WithErrors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"errors": [{"msg": "Invalid token"}, {"msg": "Connection refused"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmSettingsValidateOption{
		Key: "my-alm-setting",
	}

	result, _, err := client.AlmSettings.Validate(opt)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}

	if result.Errors[0].Msg != "Invalid token" {
		t.Errorf("expected first error to be 'Invalid token', got %s", result.Errors[0].Msg)
	}
}

func TestAlmSettings_Validate_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmSettings.Validate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, _, err = client.AlmSettings.Validate(&AlmSettingsValidateOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test Key too long
	_, _, err = client.AlmSettings.Validate(&AlmSettingsValidateOption{
		Key: strings.Repeat("a", MaxAlmKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for Key exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// Validate Function Tests (direct method tests)
// -----------------------------------------------------------------------------

func TestAlmSettings_ValidateListOpt(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Nil options should be valid
	err = client.AlmSettings.ValidateListOpt(nil)
	if err != nil {
		t.Errorf("expected nil option to be valid, got error: %v", err)
	}

	// Empty options should be valid
	err = client.AlmSettings.ValidateListOpt(&AlmSettingsListOption{})
	if err != nil {
		t.Errorf("expected empty option to be valid, got error: %v", err)
	}
}
