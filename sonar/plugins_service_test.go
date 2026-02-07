package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginsService_Available(t *testing.T) {
	response := `{
		"plugins": [
			{
				"key": "sonar-java",
				"name": "SonarJava",
				"category": "Languages",
				"description": "Java static code analyzer",
				"license": "LGPL-3.0",
				"organizationName": "SonarSource",
				"organizationUrl": "https://www.sonarsource.com",
				"editionBundled": false,
				"release": {
					"version": "7.16",
					"date": "2022-01-15",
					"changeLogUrl": "https://example.com/changelog"
				},
				"update": {
					"status": "COMPATIBLE",
					"requires": []
				}
			}
		],
		"updateCenterRefresh": "2022-01-20T10:00:00+0000"
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/plugins/available", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Plugins.Available()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Plugins, 1)
	assert.Equal(t, "sonar-java", result.Plugins[0].Key)
	assert.Equal(t, "7.16", result.Plugins[0].Release.Version)
}

func TestPluginsService_CancelAll(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/plugins/cancel_all", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.Plugins.CancelAll()
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPluginsService_Download(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "sonar-java", r.URL.Query().Get("plugin"))
			w.Header().Set("Content-Type", "application/java-archive")
			_, _ = w.Write([]byte("plugin-jar-binary-data"))
		})
		client := newTestClient(t, server.URL)

		result, resp, err := client.Plugins.Download(&PluginsDownloadOption{
			Plugin: "sonar-java",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Plugins.Download(nil)
		assert.Error(t, err)
	})

	t.Run("missing plugin", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Plugins.Download(&PluginsDownloadOption{})
		assert.Error(t, err)
	})
}

func TestPluginsService_Install(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "sonar-java", r.Form.Get("key"))
			w.WriteHeader(http.StatusNoContent)
		})
		client := newTestClient(t, server.URL)

		resp, err := client.Plugins.Install(&PluginsInstallOption{
			Key: "sonar-java",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.Plugins.Install(nil)
		assert.Error(t, err)
	})

	t.Run("missing key", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.Plugins.Install(&PluginsInstallOption{})
		assert.Error(t, err)
	})
}

func TestPluginsService_Installed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		response := `{
			"plugins": [
				{
					"key": "sonar-java",
					"name": "SonarJava",
					"description": "Java static code analyzer",
					"version": "7.16",
					"license": "LGPL-3.0",
					"organizationName": "SonarSource",
					"organizationUrl": "https://www.sonarsource.com",
					"editionBundled": true,
					"filename": "sonar-java-plugin-7.16.jar",
					"hash": "abc123",
					"sonarLintSupported": true
				}
			]
		}`
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/plugins/installed", http.StatusOK, response))
		client := newTestClient(t, server.URL)

		result, resp, err := client.Plugins.Installed(nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
		assert.Len(t, result.Plugins, 1)
		assert.Equal(t, "sonar-java", result.Plugins[0].Key)
		assert.True(t, result.Plugins[0].EditionBundled)
	})

	t.Run("with fields", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "category", r.URL.Query().Get("f"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"plugins": []}`))
		})
		client := newTestClient(t, server.URL)

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Fields: []string{"category"},
		})
		require.NoError(t, err)
	})

	t.Run("with type", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "BUNDLED", r.URL.Query().Get("type"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"plugins": []}`))
		})
		client := newTestClient(t, server.URL)

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Type: "BUNDLED",
		})
		require.NoError(t, err)
	})

	t.Run("invalid field", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Fields: []string{"invalid"},
		})
		assert.Error(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Type: "INVALID",
		})
		assert.Error(t, err)
	})
}

func TestPluginsService_Pending(t *testing.T) {
	response := `{
		"installing": [
			{
				"key": "sonar-java",
				"name": "SonarJava",
				"version": "7.16"
			}
		],
		"removing": [],
		"updating": []
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/plugins/pending", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Plugins.Pending()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Installing, 1)
	assert.Equal(t, "sonar-java", result.Installing[0].Key)
}

func TestPluginsService_Uninstall(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "sonar-java", r.Form.Get("key"))
			w.WriteHeader(http.StatusNoContent)
		})
		client := newTestClient(t, server.URL)

		resp, err := client.Plugins.Uninstall(&PluginsUninstallOption{
			Key: "sonar-java",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.Plugins.Uninstall(nil)
		assert.Error(t, err)
	})

	t.Run("missing key", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.Plugins.Uninstall(&PluginsUninstallOption{})
		assert.Error(t, err)
	})
}

func TestPluginsService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "sonar-java", r.Form.Get("key"))
			w.WriteHeader(http.StatusNoContent)
		})
		client := newTestClient(t, server.URL)

		resp, err := client.Plugins.Update(&PluginsUpdateOption{
			Key: "sonar-java",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.Plugins.Update(nil)
		assert.Error(t, err)
	})

	t.Run("missing key", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.Plugins.Update(&PluginsUpdateOption{})
		assert.Error(t, err)
	})
}

func TestPluginsService_Updates(t *testing.T) {
	response := `{
		"plugins": [
			{
				"key": "sonar-java",
				"name": "SonarJava",
				"category": "Languages",
				"description": "Java static code analyzer",
				"license": "LGPL-3.0",
				"organizationName": "SonarSource",
				"editionBundled": false,
				"updates": [
					{
						"release": {
							"version": "7.17",
							"date": "2022-02-01"
						},
						"status": "COMPATIBLE",
						"requires": []
					}
				]
			}
		]
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/plugins/updates", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Plugins.Updates()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Plugins, 1)
	assert.Equal(t, "sonar-java", result.Plugins[0].Key)
	assert.Len(t, result.Plugins[0].Updates, 1)
	assert.Equal(t, "7.17", result.Plugins[0].Updates[0].Release.Version)
}

func TestPluginsService_ValidateDownloadOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *PluginsDownloadOption
		wantErr bool
	}{
		{"valid", &PluginsDownloadOption{Plugin: "sonar-java"}, false},
		{"nil option", nil, true},
		{"empty plugin", &PluginsDownloadOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Plugins.ValidateDownloadOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPluginsService_ValidateInstallOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *PluginsInstallOption
		wantErr bool
	}{
		{"valid", &PluginsInstallOption{Key: "sonar-java"}, false},
		{"nil option", nil, true},
		{"empty key", &PluginsInstallOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Plugins.ValidateInstallOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPluginsService_ValidateInstalledOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *PluginsInstalledOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &PluginsInstalledOption{}, false},
		{"valid fields", &PluginsInstalledOption{Fields: []string{"category"}}, false},
		{"invalid fields", &PluginsInstalledOption{Fields: []string{"invalid"}}, true},
		{"valid type BUNDLED", &PluginsInstalledOption{Type: "BUNDLED"}, false},
		{"valid type EXTERNAL", &PluginsInstalledOption{Type: "EXTERNAL"}, false},
		{"invalid type", &PluginsInstalledOption{Type: "INVALID"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Plugins.ValidateInstalledOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPluginsService_ValidateUninstallOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *PluginsUninstallOption
		wantErr bool
	}{
		{"valid", &PluginsUninstallOption{Key: "sonar-java"}, false},
		{"nil option", nil, true},
		{"empty key", &PluginsUninstallOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Plugins.ValidateUninstallOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPluginsService_ValidateUpdateOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *PluginsUpdateOption
		wantErr bool
	}{
		{"valid", &PluginsUpdateOption{Key: "sonar-java"}, false},
		{"nil option", nil, true},
		{"empty key", &PluginsUpdateOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Plugins.ValidateUpdateOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
