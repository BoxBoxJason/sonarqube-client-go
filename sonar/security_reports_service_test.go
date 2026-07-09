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
				Category:                 "a1",
				Version:                  "2017",
				Vulnerabilities:          5,
				SecurityReviewRating:     2,
				ActiveRules:              10,
				TotalRules:               12,
				ToReviewSecurityHotspots: 3,
				ReviewedSecurityHotspots: 2,
				HasMoreRules:             true,
				Distribution: []SecurityReportCWEDistribution{
					{
						CWE:                      "89",
						Vulnerabilities:          1,
						VulnerabilityRating:      3,
						SecurityReviewRating:     2,
						ActiveRules:              4,
						TotalRules:               5,
						ToReviewSecurityHotspots: 1,
						ReviewedSecurityHotspots: 0,
						HasMoreRules:             false,
					},
				},
			},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/security_reports/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.SecurityReports.Show(context.Background(), &SecurityReportsShowOptions{
		Project:             "my-project",
		Standard:            "owaspTop10",
		IncludeDistribution: true,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.Len(t, result.Categories, 1)
	assert.Equal(t, "a1", result.Categories[0].Category)
	assert.Equal(t, "2017", result.Categories[0].Version)
	assert.Equal(t, 5, result.Categories[0].Vulnerabilities)
	assert.Equal(t, 2, result.Categories[0].SecurityReviewRating)
	assert.Equal(t, 10, result.Categories[0].ActiveRules)
	assert.Equal(t, 12, result.Categories[0].TotalRules)
	assert.Equal(t, 3, result.Categories[0].ToReviewSecurityHotspots)
	assert.Equal(t, 2, result.Categories[0].ReviewedSecurityHotspots)
	assert.True(t, result.Categories[0].HasMoreRules)
	require.Len(t, result.Categories[0].Distribution, 1)
	assert.Equal(t, "89", result.Categories[0].Distribution[0].CWE)
	assert.Equal(t, 1, result.Categories[0].Distribution[0].Vulnerabilities)
	assert.Equal(t, 3, result.Categories[0].Distribution[0].VulnerabilityRating)
	assert.Equal(t, 2, result.Categories[0].Distribution[0].SecurityReviewRating)
	assert.Equal(t, 4, result.Categories[0].Distribution[0].ActiveRules)
	assert.Equal(t, 5, result.Categories[0].Distribution[0].TotalRules)
	assert.Equal(t, 1, result.Categories[0].Distribution[0].ToReviewSecurityHotspots)
	assert.Equal(t, 0, result.Categories[0].Distribution[0].ReviewedSecurityHotspots)
	assert.False(t, result.Categories[0].Distribution[0].HasMoreRules)
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
