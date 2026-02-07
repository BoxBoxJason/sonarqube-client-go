package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDump_Export(t *testing.T) {
	response := ProjectDumpExport{
		ProjectID:   "proj-123",
		ProjectKey:  "my-project",
		ProjectName: "My Project",
		TaskID:      "task-456",
	}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/project_dump/export", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	opt := &ProjectDumpExportOption{
		Key: "my-project",
	}

	result, resp, err := client.ProjectDump.Export(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "proj-123", result.ProjectID)
	assert.Equal(t, "my-project", result.ProjectKey)
	assert.Equal(t, "task-456", result.TaskID)
}

func TestProjectDump_Export_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectDump.Export(nil)
	assert.Error(t, err)

	// Missing Key should fail validation.
	_, _, err = client.ProjectDump.Export(&ProjectDumpExportOption{})
	assert.Error(t, err)
}

func TestProjectDump_Status(t *testing.T) {
	response := ProjectDumpStatus{
		CanBeExported: true,
		CanBeImported: false,
		DumpToImport:  "",
		ExportedDump:  "/path/to/dump.zip",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_dump/status", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	opt := &ProjectDumpStatusOption{
		Key: "my-project",
	}

	result, resp, err := client.ProjectDump.Status(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.True(t, result.CanBeExported)
	assert.False(t, result.CanBeImported)
	assert.Equal(t, "/path/to/dump.zip", result.ExportedDump)
}

func TestProjectDump_Status_WithID(t *testing.T) {
	response := ProjectDumpStatus{
		CanBeExported: true,
		CanBeImported: true,
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_dump/status", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	opt := &ProjectDumpStatusOption{
		ID: "proj-123",
	}

	result, resp, err := client.ProjectDump.Status(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestProjectDump_Status_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectDump.Status(nil)
	assert.Error(t, err)

	// Missing both ID and Key should fail validation.
	_, _, err = client.ProjectDump.Status(&ProjectDumpStatusOption{})
	assert.Error(t, err)
}

func TestProjectDump_ValidateExportOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectDump.ValidateExportOpt(&ProjectDumpExportOption{
		Key: "my-project",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.ProjectDump.ValidateExportOpt(nil)
	assert.Error(t, err)

	// Missing Key should fail.
	err = client.ProjectDump.ValidateExportOpt(&ProjectDumpExportOption{})
	assert.Error(t, err)
}

func TestProjectDump_ValidateStatusOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option with Key should pass.
	err := client.ProjectDump.ValidateStatusOpt(&ProjectDumpStatusOption{
		Key: "my-project",
	})
	assert.NoError(t, err)

	// Valid option with ID should pass.
	err = client.ProjectDump.ValidateStatusOpt(&ProjectDumpStatusOption{
		ID: "proj-123",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.ProjectDump.ValidateStatusOpt(nil)
	assert.Error(t, err)

	// Missing both ID and Key should fail.
	err = client.ProjectDump.ValidateStatusOpt(&ProjectDumpStatusOption{})
	assert.Error(t, err)
}
