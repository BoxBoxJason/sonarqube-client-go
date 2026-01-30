package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNavigationService_Component(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/navigation/component" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("expected method GET, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := NavigationComponent{
				Key:          "my-project",
				Name:         "My Project",
				AnalysisDate: "2023-05-01T12:00:00+0000",
				IsFavorite:   true,
				Breadcrumbs: []NavigationBreadcrumb{
					{Key: "my-project", Name: "My Project", Qualifier: "TRK"},
				},
				QualityGate: NavigationQualityGate{
					Key:       "1",
					Name:      "Sonar way",
					IsDefault: true,
				},
				QualityProfiles: []NavigationQualityProfile{
					{Key: "AX-xyz", Language: "java", Name: "Sonar way"},
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client, err := NewClient(server.URL+"/api/", "user", "pass")
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		opt := &NavigationComponentOption{
			Component: "my-project",
		}

		result, resp, err := client.Navigation.Component(opt)
		if err != nil {
			t.Fatalf("Component failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result.Key != "my-project" {
			t.Errorf("expected key 'my-project', got %s", result.Key)
		}
		if !result.IsFavorite {
			t.Error("expected IsFavorite to be true")
		}
		if len(result.Breadcrumbs) != 1 {
			t.Errorf("expected 1 breadcrumb, got %d", len(result.Breadcrumbs))
		}
	})

	t.Run("with branch", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("branch") != "feature/my-branch" {
				t.Errorf("expected branch 'feature/my-branch', got %s", r.URL.Query().Get("branch"))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(NavigationComponent{Key: "my-project"})
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		opt := &NavigationComponentOption{
			Component: "my-project",
			Branch:    "feature/my-branch",
		}

		_, _, err := client.Navigation.Component(opt)
		if err != nil {
			t.Fatalf("Component failed: %v", err)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(NavigationComponent{})
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Navigation.Component(nil)
		if err != nil {
			t.Fatalf("Component with nil option failed: %v", err)
		}
	})
}

func TestNavigationService_Global(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/navigation/global" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("expected method GET, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := NavigationGlobal{
				Version:            "10.5.0",
				VersionEOL:         "2025-12-31",
				Edition:            "enterprise",
				CanAdmin:           true,
				Standalone:         true,
				ProductionDatabase: true,
				Qualifiers:         []string{"TRK", "VW", "APP"},
				GlobalPages: []NavigationExtension{
					{Key: "page1", Name: "Page 1"},
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client, err := NewClient(server.URL+"/api/", "user", "pass")
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		result, resp, err := client.Navigation.Global()
		if err != nil {
			t.Fatalf("Global failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result.Version != "10.5.0" {
			t.Errorf("expected version '10.5.0', got %s", result.Version)
		}
		if !result.CanAdmin {
			t.Error("expected CanAdmin to be true")
		}
		if len(result.Qualifiers) != 3 {
			t.Errorf("expected 3 qualifiers, got %d", len(result.Qualifiers))
		}
	})
}

func TestNavigationService_Marketplace(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/navigation/marketplace" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("expected method GET, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := NavigationMarketplace{
				Ncloc:    1000000,
				ServerID: "ABC123-XYZ",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client, err := NewClient(server.URL+"/api/", "user", "pass")
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		result, resp, err := client.Navigation.Marketplace()
		if err != nil {
			t.Fatalf("Marketplace failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result.Ncloc != 1000000 {
			t.Errorf("expected ncloc 1000000, got %d", result.Ncloc)
		}
		if result.ServerID != "ABC123-XYZ" {
			t.Errorf("expected serverID 'ABC123-XYZ', got %s", result.ServerID)
		}
	})
}

func TestNavigationService_Settings(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/navigation/settings" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodGet {
				t.Errorf("expected method GET, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := NavigationSettings{
				ShowUpdateCenter: true,
				Extensions: []NavigationSettingsExtension{
					{Name: "Plugin Settings", URL: "/admin/settings?category=plugins"},
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client, err := NewClient(server.URL+"/api/", "user", "pass")
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		result, resp, err := client.Navigation.Settings()
		if err != nil {
			t.Fatalf("Settings failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if !result.ShowUpdateCenter {
			t.Error("expected ShowUpdateCenter to be true")
		}
		if len(result.Extensions) != 1 {
			t.Errorf("expected 1 extension, got %d", len(result.Extensions))
		}
	})
}

func TestNavigationService_ValidateComponentOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test nil option (should be valid)
	err := client.Navigation.ValidateComponentOpt(nil)
	if err != nil {
		t.Errorf("ValidateComponentOpt with nil should not return error, got: %v", err)
	}

	// Test empty option (should be valid)
	err = client.Navigation.ValidateComponentOpt(&NavigationComponentOption{})
	if err != nil {
		t.Errorf("ValidateComponentOpt with empty option should not return error, got: %v", err)
	}

	// Test with component
	err = client.Navigation.ValidateComponentOpt(&NavigationComponentOption{Component: "my-project"})
	if err != nil {
		t.Errorf("ValidateComponentOpt with component should not return error, got: %v", err)
	}
}
