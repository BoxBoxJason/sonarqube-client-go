package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPluginsService_Available(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
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
		}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Plugins.Available()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if len(result.Plugins) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(result.Plugins))
	}
	if result.Plugins[0].Key != "sonar-java" {
		t.Errorf("expected key 'sonar-java', got '%s'", result.Plugins[0].Key)
	}
	if result.Plugins[0].Release.Version != "7.16" {
		t.Errorf("expected version '7.16', got '%s'", result.Plugins[0].Release.Version)
	}
}

func TestPluginsService_CancelAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	resp, err := client.Plugins.CancelAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPluginsService_Download(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Query().Get("plugin") != "sonar-java" {
				t.Errorf("expected plugin 'sonar-java', got '%s'", r.URL.Query().Get("plugin"))
			}
			w.Header().Set("Content-Type", "application/java-archive")
			_, _ = w.Write([]byte("plugin-jar-binary-data"))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Plugins.Download(&PluginsDownloadOption{
			Plugin: "sonar-java",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Plugins.Download(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing plugin", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Plugins.Download(&PluginsDownloadOption{})
		if err == nil {
			t.Error("expected error for missing plugin")
		}
	})
}

func TestPluginsService_Install(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("key") != "sonar-java" {
				t.Errorf("expected key 'sonar-java', got '%s'", r.Form.Get("key"))
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		resp, err := client.Plugins.Install(&PluginsInstallOption{
			Key: "sonar-java",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.Plugins.Install(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing key", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.Plugins.Install(&PluginsInstallOption{})
		if err == nil {
			t.Error("expected error for missing key")
		}
	})
}

func TestPluginsService_Installed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
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
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Plugins.Installed(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if len(result.Plugins) != 1 {
			t.Errorf("expected 1 plugin, got %d", len(result.Plugins))
		}
		if result.Plugins[0].Key != "sonar-java" {
			t.Errorf("expected key 'sonar-java', got '%s'", result.Plugins[0].Key)
		}
		if !result.Plugins[0].EditionBundled {
			t.Error("expected EditionBundled to be true")
		}
	})

	t.Run("with fields", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("f") != "category" {
				t.Errorf("expected f 'category', got '%s'", r.URL.Query().Get("f"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"plugins": []}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Fields: []string{"category"},
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("type") != "BUNDLED" {
				t.Errorf("expected type 'BUNDLED', got '%s'", r.URL.Query().Get("type"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"plugins": []}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Type: "BUNDLED",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Fields: []string{"invalid"},
		})
		if err == nil {
			t.Error("expected error for invalid field")
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Plugins.Installed(&PluginsInstalledOption{
			Type: "INVALID",
		})
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestPluginsService_Pending(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"installing": [
				{
					"key": "sonar-java",
					"name": "SonarJava",
					"version": "7.16"
				}
			],
			"removing": [],
			"updating": []
		}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Plugins.Pending()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if len(result.Installing) != 1 {
		t.Errorf("expected 1 installing, got %d", len(result.Installing))
	}
	if result.Installing[0].Key != "sonar-java" {
		t.Errorf("expected key 'sonar-java', got '%s'", result.Installing[0].Key)
	}
}

func TestPluginsService_Uninstall(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("key") != "sonar-java" {
				t.Errorf("expected key 'sonar-java', got '%s'", r.Form.Get("key"))
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		resp, err := client.Plugins.Uninstall(&PluginsUninstallOption{
			Key: "sonar-java",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.Plugins.Uninstall(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing key", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.Plugins.Uninstall(&PluginsUninstallOption{})
		if err == nil {
			t.Error("expected error for missing key")
		}
	})
}

func TestPluginsService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("key") != "sonar-java" {
				t.Errorf("expected key 'sonar-java', got '%s'", r.Form.Get("key"))
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		resp, err := client.Plugins.Update(&PluginsUpdateOption{
			Key: "sonar-java",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.Plugins.Update(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing key", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.Plugins.Update(&PluginsUpdateOption{})
		if err == nil {
			t.Error("expected error for missing key")
		}
	})
}

func TestPluginsService_Updates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
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
		}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Plugins.Updates()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if len(result.Plugins) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(result.Plugins))
	}
	if result.Plugins[0].Key != "sonar-java" {
		t.Errorf("expected key 'sonar-java', got '%s'", result.Plugins[0].Key)
	}
	if len(result.Plugins[0].Updates) != 1 {
		t.Errorf("expected 1 update, got %d", len(result.Plugins[0].Updates))
	}
	if result.Plugins[0].Updates[0].Release.Version != "7.17" {
		t.Errorf("expected version '7.17', got '%s'", result.Plugins[0].Updates[0].Release.Version)
	}
}

func TestPluginsService_ValidateDownloadOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDownloadOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginsService_ValidateInstallOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInstallOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginsService_ValidateInstalledOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInstalledOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginsService_ValidateUninstallOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUninstallOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginsService_ValidateUpdateOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
