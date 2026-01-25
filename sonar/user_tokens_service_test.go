package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserTokens_Generate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/user_tokens/generate" {
			t.Errorf("expected path /api/user_tokens/generate, got %s", r.URL.Path)
		}

		name := r.URL.Query().Get("name")
		if name != "my-token" {
			t.Errorf("expected name 'my-token', got %s", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"login": "admin",
			"name": "my-token",
			"token": "secret-token-value",
			"createdAt": "2024-01-01T00:00:00+0000",
			"type": "USER_TOKEN"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserTokensGenerateOption{
		Name: "my-token",
	}

	result, resp, err := client.UserTokens.Generate(opt)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Name != "my-token" {
		t.Errorf("expected name 'my-token', got %s", result.Name)
	}

	if result.Token != "secret-token-value" {
		t.Errorf("expected token 'secret-token-value', got %s", result.Token)
	}

	if result.Type != "USER_TOKEN" {
		t.Errorf("expected type 'USER_TOKEN', got %s", result.Type)
	}
}

func TestUserTokens_Generate_WithType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenType := r.URL.Query().Get("type")
		if tokenType != "PROJECT_ANALYSIS_TOKEN" {
			t.Errorf("expected type 'PROJECT_ANALYSIS_TOKEN', got %s", tokenType)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"login": "admin",
			"name": "project-token",
			"token": "project-token-value",
			"createdAt": "2024-01-01T00:00:00+0000",
			"type": "PROJECT_ANALYSIS_TOKEN"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserTokensGenerateOption{
		Name:       "project-token",
		Type:       "PROJECT_ANALYSIS_TOKEN",
		ProjectKey: "my-project",
	}

	result, resp, err := client.UserTokens.Generate(opt)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Type != "PROJECT_ANALYSIS_TOKEN" {
		t.Errorf("expected type 'PROJECT_ANALYSIS_TOKEN', got %s", result.Type)
	}
}

func TestUserTokens_Generate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.UserTokens.Generate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Name should fail validation.
	_, _, err = client.UserTokens.Generate(&UserTokensGenerateOption{})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Invalid Type should fail validation.
	_, _, err = client.UserTokens.Generate(&UserTokensGenerateOption{
		Name: "my-token",
		Type: "INVALID_TYPE",
	})
	if err == nil {
		t.Error("expected error for invalid Type")
	}

	// PROJECT_ANALYSIS_TOKEN without ProjectKey should fail validation.
	_, _, err = client.UserTokens.Generate(&UserTokensGenerateOption{
		Name: "my-token",
		Type: "PROJECT_ANALYSIS_TOKEN",
	})
	if err == nil {
		t.Error("expected error for PROJECT_ANALYSIS_TOKEN without ProjectKey")
	}
}

func TestUserTokens_Revoke(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/user_tokens/revoke" {
			t.Errorf("expected path /api/user_tokens/revoke, got %s", r.URL.Path)
		}

		name := r.URL.Query().Get("name")
		if name != "my-token" {
			t.Errorf("expected name 'my-token', got %s", name)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserTokensRevokeOption{
		Name: "my-token",
	}

	resp, err := client.UserTokens.Revoke(opt)
	if err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUserTokens_Revoke_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.UserTokens.Revoke(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Name should fail validation.
	_, err = client.UserTokens.Revoke(&UserTokensRevokeOption{})
	if err == nil {
		t.Error("expected error for missing Name")
	}
}

func TestUserTokens_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/user_tokens/search" {
			t.Errorf("expected path /api/user_tokens/search, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
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
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.UserTokens.Search(nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Login != "admin" {
		t.Errorf("expected login 'admin', got %s", result.Login)
	}

	if len(result.UserTokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(result.UserTokens))
	}

	if result.UserTokens[0].Name != "token1" {
		t.Errorf("expected first token name 'token1', got %s", result.UserTokens[0].Name)
	}

	if result.UserTokens[1].IsExpired != true {
		t.Error("expected second token to be expired")
	}
}

func TestUserTokens_Search_WithLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		if login != "testuser" {
			t.Errorf("expected login 'testuser', got %s", login)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"login": "testuser", "userTokens": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserTokensSearchOption{
		Login: "testuser",
	}

	result, _, err := client.UserTokens.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Login != "testuser" {
		t.Errorf("expected login 'testuser', got %s", result.Login)
	}
}

func TestUserTokens_ValidateGenerateOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
		Name: "my-token",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// All valid token types should pass.
	validTypes := []string{"USER_TOKEN", "GLOBAL_ANALYSIS_TOKEN"}
	for _, tokenType := range validTypes {
		err := client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
			Name: "my-token",
			Type: tokenType,
		})
		if err != nil {
			t.Errorf("expected nil error for type '%s', got %v", tokenType, err)
		}
	}

	// PROJECT_ANALYSIS_TOKEN with ProjectKey should pass.
	err = client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
		Name:       "project-token",
		Type:       "PROJECT_ANALYSIS_TOKEN",
		ProjectKey: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error for PROJECT_ANALYSIS_TOKEN with ProjectKey, got %v", err)
	}

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxTokenNameLength+1; i++ {
		longName += "a"
	}
	err = client.UserTokens.ValidateGenerateOpt(&UserTokensGenerateOption{
		Name: longName,
	})
	if err == nil {
		t.Error("expected error for name exceeding max length")
	}
}
