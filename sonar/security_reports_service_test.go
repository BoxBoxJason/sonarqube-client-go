package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Download
// -----------------------------------------------------------------------------

func TestSecurityReportsService_Download(t *testing.T) {
	data := []byte(`%PDF test content`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/security_reports/download", http.StatusOK, "application/pdf", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.SecurityReports.Download(context.Background(), &SecurityReportsDownloadOptions{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestSecurityReportsService_Download_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.SecurityReports.Download(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.SecurityReports.Download(context.Background(), &SecurityReportsDownloadOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// -----------------------------------------------------------------------------
// Show
// -----------------------------------------------------------------------------

func TestSecurityReportsService_Show(t *testing.T) {
	response := SecurityReportsShow{
		Categories: []SecurityReportCategory{
			{
				Category:                 "A1",
				Vulnerabilities:          5,
				ActiveRules:              10,
				ToReviewSecurityHotspots: 3,
				ReviewedSecurityHotspots: 2,
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/security_reports/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.SecurityReports.Show(context.Background(), &SecurityReportsShowOptions{
		Project:  "my-project",
		Standard: "owaspTop10",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Categories, 1)
	assert.Equal(t, "A1", result.Categories[0].Category)
}

func TestSecurityReportsService_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.SecurityReports.Show(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.SecurityReports.Show(context.Background(), &SecurityReportsShowOptions{Standard: "owaspTop10"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.SecurityReports.Show(context.Background(), &SecurityReportsShowOptions{Project: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.SecurityReports.Show(context.Background(), &SecurityReportsShowOptions{Project: "my-project", Standard: "invalid-standard"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
