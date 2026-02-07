package sonar

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNavigationService_Component(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/navigation/component", http.StatusOK, response))
		client := newTestClient(t, server.URL)

		opt := &NavigationComponentOption{
			Component: "my-project",
		}

		result, resp, err := client.Navigation.Component(opt)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "my-project", result.Key)
		assert.True(t, result.IsFavorite)
		assert.Len(t, result.Breadcrumbs, 1)
	})

	t.Run("with branch", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/navigation/component", r.URL.Path)
			assert.Equal(t, "feature/my-branch", r.URL.Query().Get("branch"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(NavigationComponent{Key: "my-project"})
		})
		client := newTestClient(t, server.URL)

		opt := &NavigationComponentOption{
			Component: "my-project",
			Branch:    "feature/my-branch",
		}

		_, _, err := client.Navigation.Component(opt)
		require.NoError(t, err)
	})

	t.Run("nil option", func(t *testing.T) {
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/navigation/component", http.StatusOK, NavigationComponent{}))
		client := newTestClient(t, server.URL)

		_, _, err := client.Navigation.Component(nil)
		require.NoError(t, err)
	})
}

func TestNavigationService_Global(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/navigation/global", http.StatusOK, response))
		client := newTestClient(t, server.URL)

		result, resp, err := client.Navigation.Global()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "10.5.0", result.Version)
		assert.True(t, result.CanAdmin)
		assert.Len(t, result.Qualifiers, 3)
	})
}

func TestNavigationService_Marketplace(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		response := NavigationMarketplace{
			Ncloc:    1000000,
			ServerID: "ABC123-XYZ",
		}
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/navigation/marketplace", http.StatusOK, response))
		client := newTestClient(t, server.URL)

		result, resp, err := client.Navigation.Marketplace()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, int64(1000000), result.Ncloc)
		assert.Equal(t, "ABC123-XYZ", result.ServerID)
	})
}

func TestNavigationService_Settings(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		response := NavigationSettings{
			ShowUpdateCenter: true,
			Extensions: []NavigationSettingsExtension{
				{Name: "Plugin Settings", URL: "/admin/settings?category=plugins"},
			},
		}
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/navigation/settings", http.StatusOK, response))
		client := newTestClient(t, server.URL)

		result, resp, err := client.Navigation.Settings()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, result.ShowUpdateCenter)
		assert.Len(t, result.Extensions, 1)
	})
}

func TestNavigationService_ValidateComponentOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option (should be valid)
	err := client.Navigation.ValidateComponentOpt(nil)
	assert.NoError(t, err)

	// Test empty option (should be valid)
	err = client.Navigation.ValidateComponentOpt(&NavigationComponentOption{})
	assert.NoError(t, err)

	// Test with component
	err = client.Navigation.ValidateComponentOpt(&NavigationComponentOption{Component: "my-project"})
	assert.NoError(t, err)
}
