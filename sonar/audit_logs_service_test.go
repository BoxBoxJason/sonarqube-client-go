package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditLogsService_Download(t *testing.T) {
	data := []byte(`{"events":[{"category":"USER_AUTHENTICATION","action":"LOGIN"}]}`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/audit_logs/download", http.StatusOK, "application/json", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.AuditLogs.Download(context.Background(), &AuditLogsDownloadOptions{
		From: "2024-01-01T00:00:00+00:00",
		To:   "2024-12-31T23:59:59+00:00",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestAuditLogsService_Download_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.AuditLogs.Download(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, result)

	result, resp, err = client.AuditLogs.Download(context.Background(), &AuditLogsDownloadOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, result)

	result, resp, err = client.AuditLogs.Download(context.Background(), &AuditLogsDownloadOptions{From: "2024-01-01T00:00:00+00:00"})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, result)

	result, resp, err = client.AuditLogs.Download(context.Background(), &AuditLogsDownloadOptions{To: "2024-12-31T23:59:59+00:00"})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, result)
}
