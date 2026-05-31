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

func TestGovernanceReportsService_Download(t *testing.T) {
	data := []byte(`%PDF-1.4 report content`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/governance_reports/download", http.StatusOK, "application/pdf", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.GovernanceReports.Download(context.Background(), &GovernanceReportsDownloadOptions{
		ComponentKey: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestGovernanceReportsService_Download_NoOptions(t *testing.T) {
	data := []byte(`%PDF-1.4 report content`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/governance_reports/download", http.StatusOK, "application/pdf", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.GovernanceReports.Download(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

// -----------------------------------------------------------------------------
// Status
// -----------------------------------------------------------------------------

func TestGovernanceReportsService_Status(t *testing.T) {
	response := GovernanceReportsStatus{
		HasFile:      true,
		Subscribed:   false,
		ComponentKey: "my-portfolio",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/governance_reports/status", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.GovernanceReports.Status(context.Background(), &GovernanceReportsStatusOptions{
		ComponentKey: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
	assert.True(t, result.HasFile)
	assert.Equal(t, "my-portfolio", result.ComponentKey)
}

// -----------------------------------------------------------------------------
// Subscribe / Unsubscribe
// -----------------------------------------------------------------------------

func TestGovernanceReportsService_Subscribe(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/governance_reports/subscribe", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.GovernanceReports.Subscribe(context.Background(), &GovernanceReportsSubscribeOptions{
		ComponentKey: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGovernanceReportsService_Unsubscribe(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/governance_reports/unsubscribe", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.GovernanceReports.Unsubscribe(context.Background(), &GovernanceReportsUnsubscribeOptions{
		ComponentKey: "my-portfolio",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// UpdateFrequency
// -----------------------------------------------------------------------------

func TestGovernanceReportsService_UpdateFrequency(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/governance_reports/update_frequency", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.GovernanceReports.UpdateFrequency(context.Background(), &GovernanceReportsUpdateFrequencyOptions{
		ComponentKey: "my-portfolio",
		Frequency:    "WEEKLY",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// UpdateRecipients
// -----------------------------------------------------------------------------

func TestGovernanceReportsService_UpdateRecipients(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/governance_reports/update_recipients", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.GovernanceReports.UpdateRecipients(context.Background(), &GovernanceReportsUpdateRecipientsOptions{
		ComponentKey: "my-portfolio",
		Recipients:   "user@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGovernanceReportsService_UpdateRecipients_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.GovernanceReports.UpdateRecipients(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.GovernanceReports.UpdateRecipients(context.Background(), &GovernanceReportsUpdateRecipientsOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}
