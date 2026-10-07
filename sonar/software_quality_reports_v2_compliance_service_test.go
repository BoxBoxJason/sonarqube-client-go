package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSoftwareQualityReportsService_GetComplianceReport(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/software-quality-reports/compliance-reports", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.SoftwareQualityReports.GetComplianceReport(context.Background(), &SoftwareQualityReportsGetComplianceReportOptions{ComponentKey: "test", Standard: "misra", Version: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestSoftwareQualityReportsService_GetComplianceReport_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.SoftwareQualityReports.GetComplianceReport(context.Background(), nil)
	require.Error(t, err)
}
