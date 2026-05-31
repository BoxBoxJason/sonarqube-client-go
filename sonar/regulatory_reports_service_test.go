package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegulatoryReportsService_Download(t *testing.T) {
	data := []byte("PK\x03\x04zip-content")
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/regulatory_reports/download", http.StatusOK, "application/zip", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.RegulatoryReports.Download(context.Background(), &RegulatoryReportsDownloadOptions{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestRegulatoryReportsService_Download_WithBranch(t *testing.T) {
	data := []byte("PK\x03\x04zip-content")
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/regulatory_reports/download", http.StatusOK, "application/zip", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.RegulatoryReports.Download(context.Background(), &RegulatoryReportsDownloadOptions{
		Project: "my-project",
		Branch:  "main",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestRegulatoryReportsService_Download_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.RegulatoryReports.Download(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, result)

	result, resp, err = client.RegulatoryReports.Download(context.Background(), &RegulatoryReportsDownloadOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Nil(t, result)
}
