package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGithubProvisioning_Check(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/github_provisioning/check" {
			t.Errorf("expected path /api/github_provisioning/check, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"application": {
				"autoProvisioning": {
					"status": "success"
				},
				"jit": {
					"status": "success"
				}
			},
			"installations": [
				{
					"autoProvisioning": {
						"status": "success"
					},
					"jit": {
						"status": "success"
					},
					"organization": "my-org"
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.GithubProvisioning.Check()
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Application.AutoProvisioning.Status != "success" {
		t.Errorf("expected application autoProvisioning status 'success', got %s", result.Application.AutoProvisioning.Status)
	}

	if result.Application.Jit.Status != "success" {
		t.Errorf("expected application jit status 'success', got %s", result.Application.Jit.Status)
	}

	if len(result.Installations) != 1 {
		t.Errorf("expected 1 installation, got %d", len(result.Installations))
	}

	installation := result.Installations[0]
	if installation.Organization != "my-org" {
		t.Errorf("expected organization 'my-org', got %s", installation.Organization)
	}
}

func TestGithubProvisioning_Check_WithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"application": {
				"autoProvisioning": {
					"status": "failed",
					"errorMessage": "Invalid token"
				},
				"jit": {
					"status": "success"
				}
			},
			"installations": []
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.GithubProvisioning.Check()
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Application.AutoProvisioning.Status != "failed" {
		t.Errorf("expected status 'failed', got %s", result.Application.AutoProvisioning.Status)
	}

	if result.Application.AutoProvisioning.ErrorMessage != "Invalid token" {
		t.Errorf("expected error message 'Invalid token', got %s", result.Application.AutoProvisioning.ErrorMessage)
	}
}

func TestGithubProvisioning_Check_EmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.GithubProvisioning.Check()
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
