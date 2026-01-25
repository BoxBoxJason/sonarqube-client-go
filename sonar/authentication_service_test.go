package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthentication_Login(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/authentication/login" {
			t.Errorf("expected path /api/authentication/login, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "admin" {
			t.Errorf("expected login 'admin', got %s", login)
		}

		password := r.URL.Query().Get("password")
		if password != "secret" {
			t.Errorf("expected password 'secret', got %s", password)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AuthenticationLoginOption{
		Login:    "admin",
		Password: "secret",
	}

	resp, err := client.Authentication.Login(opt)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAuthentication_Login_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Authentication.Login(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, err = client.Authentication.Login(&AuthenticationLoginOption{
		Password: "secret",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Missing Password should fail validation.
	_, err = client.Authentication.Login(&AuthenticationLoginOption{
		Login: "admin",
	})
	if err == nil {
		t.Error("expected error for missing Password")
	}
}

func TestAuthentication_Logout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/authentication/logout" {
			t.Errorf("expected path /api/authentication/logout, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Authentication.Logout()
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAuthentication_Validate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/authentication/validate" {
			t.Errorf("expected path /api/authentication/validate, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"valid": true}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Authentication.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.Valid {
		t.Error("expected valid to be true")
	}
}

func TestAuthentication_Validate_Invalid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"valid": false}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Authentication.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Valid {
		t.Error("expected valid to be false")
	}
}

func TestAuthentication_ValidateLoginOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.Authentication.ValidateLoginOpt(&AuthenticationLoginOption{
		Login:    "admin",
		Password: "secret",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.Authentication.ValidateLoginOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail.
	err = client.Authentication.ValidateLoginOpt(&AuthenticationLoginOption{
		Password: "secret",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Missing Password should fail.
	err = client.Authentication.ValidateLoginOpt(&AuthenticationLoginOption{
		Login: "admin",
	})
	if err == nil {
		t.Error("expected error for missing Password")
	}
}
