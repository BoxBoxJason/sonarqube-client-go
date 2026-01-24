package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// -----------------------------------------------------------------------------
// Service Method Tests
// -----------------------------------------------------------------------------

func TestSystem_ChangeLogLevel(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		// Return 204 No Content (success)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	opt := &SystemChangeLogLevelOption{Level: "INFO"}
	resp, err := client.System.ChangeLogLevel(opt)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestSystem_ChangeLogLevel_ValidationErrors(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name string
		opt  *SystemChangeLogLevelOption
	}{
		{"nil options", nil},
		{"missing level", &SystemChangeLogLevelOption{}},
		{"invalid level", &SystemChangeLogLevelOption{Level: "INVALID"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.System.ChangeLogLevel(tc.opt)
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}

func TestSystem_DbMigrationStatus(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"state":"NO_MIGRATION","message":"Database is up to date."}`))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.DbMigrationStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result.State != "NO_MIGRATION" {
		t.Errorf("expected state NO_MIGRATION, got %s", result.State)
	}
}

func TestSystem_Health(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"health":"GREEN","nodes":[{"name":"node1","host":"localhost","port":9000,"health":"GREEN","type":"APPLICATION"}]}`))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.Health()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result.Health != "GREEN" {
		t.Errorf("expected health GREEN, got %s", result.Health)
	}

	if len(result.Nodes) == 0 {
		t.Error("expected nodes in response")
	}
}

func TestSystem_Health_WithCauses(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"health":"YELLOW","causes":[{"message":"Cluster has only one search node"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, _, err := client.System.Health()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Health != "YELLOW" {
		t.Errorf("expected health YELLOW, got %s", result.Health)
	}

	if len(result.Causes) != 1 {
		t.Errorf("expected 1 cause, got %d", len(result.Causes))
	}
}

func TestSystem_Info(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"Health": "GREEN",
			"System": {
				"Version": "10.0.0",
				"Edition": "Community",
				"Docker": false,
				"Processors": 8
			},
			"Database": {
				"Database": "PostgreSQL",
				"Database Version": "14.0"
			}
		}`))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.Info()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result.Health != "GREEN" {
		t.Errorf("expected health GREEN, got %s", result.Health)
	}

	if result.System.Version != "10.0.0" {
		t.Errorf("expected version 10.0.0, got %s", result.System.Version)
	}
}

func TestSystem_Liveness(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return 204 No Content (liveness success)
		w.WriteHeader(http.StatusNoContent)

		_, _ = w.Write([]byte("{}"))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	_, resp, err := client.System.Liveness()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestSystem_Logs(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return plain text response
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("2024-01-01 12:00:00 INFO Sample log message"))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method with specific log name
	opt := &SystemLogsOption{Name: "app"}
	result, resp, err := client.System.Logs(opt)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result == nil || *result == "" {
		t.Error("expected log content, got empty")
	}
}

func TestSystem_Logs_NilOption(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("default log content"))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call with nil options (should use defaults)
	result, _, err := client.System.Logs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil || *result == "" {
		t.Error("expected log content, got empty")
	}
}

func TestSystem_Logs_ValidationErrors(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Invalid log name
	opt := &SystemLogsOption{Name: "invalid_log"}
	_, _, err = client.System.Logs(opt)

	if err == nil {
		t.Error("expected validation error for invalid log name, got nil")
	}
}

func TestSystem_MigrateDb(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"state":"MIGRATION_RUNNING","startedAt":"2024-01-01T12:00:00+0000"}`))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.MigrateDb()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result.State != "MIGRATION_RUNNING" {
		t.Errorf("expected state MIGRATION_RUNNING, got %s", result.State)
	}
}

func TestSystem_Ping(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return plain text response
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("pong"))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result == nil || *result != "pong" {
		t.Errorf("expected 'pong', got '%v'", result)
	}
}

func TestSystem_Restart(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		// Return 204 No Content (success)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	resp, err := client.System.Restart()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestSystem_Status(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"ABC123","status":"UP","version":"10.0.0"}`))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.Status()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if result.Status != "UP" {
		t.Errorf("expected status UP, got %s", result.Status)
	}

	if result.Version != "10.0.0" {
		t.Errorf("expected version 10.0.0, got %s", result.Version)
	}
}

func TestSystem_Upgrades(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"installedVersionActive": true,
			"latestLTA": "10.1.0",
			"updateCenterRefresh": "2024-01-01T12:00:00+0000",
			"upgrades": [
				{
					"version": "10.1.0",
					"description": "New version",
					"releaseDate": "2024-01-01",
					"changeLogUrl": "https://example.com/changelog",
					"plugins": {
						"incompatible": [
							{"key": "old-plugin", "name": "Old Plugin"}
						],
						"requireUpdate": [
							{"key": "update-plugin", "name": "Update Plugin", "version": "2.0.0"}
						]
					}
				}
			]
		}`))
	}))
	defer ts.Close()

	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Call service method
	result, resp, err := client.System.Upgrades()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if !result.InstalledVersionActive {
		t.Error("expected InstalledVersionActive to be true")
	}

	if result.LatestLTA != "10.1.0" {
		t.Errorf("expected latestLTA 10.1.0, got %s", result.LatestLTA)
	}

	if len(result.Upgrades) != 1 {
		t.Errorf("expected 1 upgrade, got %d", len(result.Upgrades))
	}

	if len(result.Upgrades[0].Plugins.Incompatible) != 1 {
		t.Errorf("expected 1 incompatible plugin, got %d", len(result.Upgrades[0].Plugins.Incompatible))
	}

	if len(result.Upgrades[0].Plugins.RequireUpdate) != 1 {
		t.Errorf("expected 1 requireUpdate plugin, got %d", len(result.Upgrades[0].Plugins.RequireUpdate))
	}
}

func TestSystem_Upgrades_NoUpgrades(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"installedVersionActive": true,
			"latestLTA": "10.0.0",
			"upgrades": []
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, _, err := client.System.Upgrades()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Upgrades) != 0 {
		t.Errorf("expected 0 upgrades, got %d", len(result.Upgrades))
	}
}

// -----------------------------------------------------------------------------
// Validation Function Tests
// -----------------------------------------------------------------------------

func TestValidateChangeLogLevelOpt(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		opt     *SystemChangeLogLevelOption
		wantErr bool
	}{
		{
			name:    "nil options",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing level",
			opt:     &SystemChangeLogLevelOption{},
			wantErr: true,
		},
		{
			name:    "invalid level",
			opt:     &SystemChangeLogLevelOption{Level: "VERBOSE"},
			wantErr: true,
		},
		{
			name:    "valid TRACE",
			opt:     &SystemChangeLogLevelOption{Level: "TRACE"},
			wantErr: false,
		},
		{
			name:    "valid DEBUG",
			opt:     &SystemChangeLogLevelOption{Level: "DEBUG"},
			wantErr: false,
		},
		{
			name:    "valid INFO",
			opt:     &SystemChangeLogLevelOption{Level: "INFO"},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.System.ValidateChangeLogLevelOpt(tc.opt)
			if tc.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateLogsOpt(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		opt     *SystemLogsOption
		wantErr bool
	}{
		{
			name:    "nil options (allowed)",
			opt:     nil,
			wantErr: false,
		},
		{
			name:    "empty name (allowed, uses default)",
			opt:     &SystemLogsOption{},
			wantErr: false,
		},
		{
			name:    "invalid name",
			opt:     &SystemLogsOption{Name: "invalid_log_type"},
			wantErr: true,
		},
		{
			name:    "valid access",
			opt:     &SystemLogsOption{Name: "access"},
			wantErr: false,
		},
		{
			name:    "valid app",
			opt:     &SystemLogsOption{Name: "app"},
			wantErr: false,
		},
		{
			name:    "valid ce",
			opt:     &SystemLogsOption{Name: "ce"},
			wantErr: false,
		},
		{
			name:    "valid deprecation",
			opt:     &SystemLogsOption{Name: "deprecation"},
			wantErr: false,
		},
		{
			name:    "valid es",
			opt:     &SystemLogsOption{Name: "es"},
			wantErr: false,
		},
		{
			name:    "valid web",
			opt:     &SystemLogsOption{Name: "web"},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.System.ValidateLogsOpt(tc.opt)
			if tc.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
