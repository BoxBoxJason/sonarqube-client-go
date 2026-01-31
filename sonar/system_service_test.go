package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Service Method Tests
// -----------------------------------------------------------------------------

func TestSystem_ChangeLogLevel(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/system/change_log_level", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &SystemChangeLogLevelOption{Level: "INFO"}
	resp, err := client.System.ChangeLogLevel(opt)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestSystem_ChangeLogLevel_ValidationErrors(t *testing.T) {
	client := newLocalhostClient(t)

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
			assert.Error(t, err, "expected validation error")
		})
	}
}

func TestSystem_DbMigrationStatus(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/db_migration_status", http.StatusOK,
		`{"state":"NO_MIGRATION","message":"Database is up to date."}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.DbMigrationStatus()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "NO_MIGRATION", result.State)
}

func TestSystem_Health(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/health", http.StatusOK,
		`{"health":"GREEN","nodes":[{"name":"node1","host":"localhost","port":9000,"health":"GREEN","type":"APPLICATION"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.Health()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GREEN", result.Health)
	assert.NotEmpty(t, result.Nodes)
}

func TestSystem_Health_WithCauses(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/health", http.StatusOK,
		`{"health":"YELLOW","causes":[{"message":"Cluster has only one search node"}]}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, _, err := client.System.Health()

	require.NoError(t, err)
	assert.Equal(t, "YELLOW", result.Health)
	assert.Len(t, result.Causes, 1)
}

func TestSystem_Info(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/info", http.StatusOK, `{
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
	}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.Info()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GREEN", result.Health)
	assert.Equal(t, "10.0.0", result.System.Version)
}

func TestSystem_Liveness(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/liveness", http.StatusNoContent, "{}")
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	_, resp, err := client.System.Liveness()

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestSystem_Logs(t *testing.T) {
	handler := mockBinaryHandler(t, http.MethodGet, "/system/logs", http.StatusOK, "text/plain",
		[]byte("2024-01-01 12:00:00 INFO Sample log message"),
	)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &SystemLogsOption{Name: "app"}
	result, resp, err := client.System.Logs(opt)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)
}

func TestSystem_Logs_NilOption(t *testing.T) {
	handler := mockBinaryHandler(t, http.MethodGet, "/system/logs", http.StatusOK, "text/plain",
		[]byte("default log content"))
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, _, err := client.System.Logs(nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)
}

func TestSystem_Logs_ValidationErrors(t *testing.T) {
	client := newLocalhostClient(t)

	opt := &SystemLogsOption{Name: "invalid_log"}
	_, _, err := client.System.Logs(opt)

	assert.Error(t, err, "expected validation error for invalid log name")
}

func TestSystem_MigrateDb(t *testing.T) {
	handler := mockHandler(t, http.MethodPost, "/system/migrate_db", http.StatusOK,
		`{"state":"MIGRATION_RUNNING","startedAt":"2024-01-01T12:00:00+0000"}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.MigrateDb()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "MIGRATION_RUNNING", result.State)
}

func TestSystem_Ping(t *testing.T) {
	handler := mockBinaryHandler(t, http.MethodGet, "/system/ping", http.StatusOK, "text/plain", []byte("pong"))
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.Ping()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "pong", *result)
}

func TestSystem_Restart(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/system/restart", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	resp, err := client.System.Restart()

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestSystem_Status(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/status", http.StatusOK,
		`{"id":"ABC123","status":"UP","version":"10.0.0"}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.Status()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "UP", result.Status)
	assert.Equal(t, "10.0.0", result.Version)
}

func TestSystem_Upgrades(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/upgrades", http.StatusOK, `{
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
	}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.System.Upgrades()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.InstalledVersionActive)
	assert.Equal(t, "10.1.0", result.LatestLTA)
	assert.Len(t, result.Upgrades, 1)
	assert.Len(t, result.Upgrades[0].Plugins.Incompatible, 1)
	assert.Len(t, result.Upgrades[0].Plugins.RequireUpdate, 1)
}

func TestSystem_Upgrades_NoUpgrades(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/system/upgrades", http.StatusOK, `{
		"installedVersionActive": true,
		"latestLTA": "10.0.0",
		"upgrades": []
	}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, _, err := client.System.Upgrades()

	require.NoError(t, err)
	assert.Empty(t, result.Upgrades)
}

// -----------------------------------------------------------------------------
// Validation Function Tests
// -----------------------------------------------------------------------------

func TestValidateChangeLogLevelOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLogsOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
