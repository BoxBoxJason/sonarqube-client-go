package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// GetAccessibilityReport
// -----------------------------------------------------------------------------

func TestSoftwareQualityReportsService_GetAccessibilityReport(t *testing.T) {
	response := SoftwareQualityReportsAccessibilityReport{
		Categories: []SoftwareQualityReportsCategory{
			{Key: "perceivable", ActiveRules: 5, Issues: 2},
		},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/software-quality-reports/accessibility-reports", http.StatusOK,
		map[string]string{
			"projectKey": "my-project",
			"standard":   SoftwareQualityReportsAccessibilityStandardWCAG,
			"version":    SoftwareQualityReportsAccessibilityVersion21,
			"branchKey":  "my-branch",
		}, response))
	client := newTestClient(t, server.URL)
	svc := &SoftwareQualityReportsService{client: client}

	result, resp, err := svc.GetAccessibilityReport(context.Background(), &SoftwareQualityReportsGetAccessibilityReportOptions{
		ProjectKey: "my-project",
		Standard:   SoftwareQualityReportsAccessibilityStandardWCAG,
		Version:    SoftwareQualityReportsAccessibilityVersion21,
		BranchKey:  "my-branch",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Categories, 1)
	assert.Equal(t, "perceivable", result.Categories[0].Key)
	assert.Equal(t, int32(5), result.Categories[0].ActiveRules)
	assert.Equal(t, int32(2), result.Categories[0].Issues)
}

func TestSoftwareQualityReportsService_GetAccessibilityReport_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &SoftwareQualityReportsService{client: client}

	tests := []struct {
		opt  *SoftwareQualityReportsGetAccessibilityReportOptions
		name string
	}{
		{nil, "nil opt"},
		{&SoftwareQualityReportsGetAccessibilityReportOptions{Standard: SoftwareQualityReportsAccessibilityStandardWCAG, Version: SoftwareQualityReportsAccessibilityVersion20}, "missing ProjectKey"},
		{&SoftwareQualityReportsGetAccessibilityReportOptions{ProjectKey: "p1", Version: SoftwareQualityReportsAccessibilityVersion20}, "missing Standard"},
		{&SoftwareQualityReportsGetAccessibilityReportOptions{ProjectKey: "p1", Standard: "section508", Version: SoftwareQualityReportsAccessibilityVersion20}, "invalid Standard"},
		{&SoftwareQualityReportsGetAccessibilityReportOptions{ProjectKey: "p1", Standard: SoftwareQualityReportsAccessibilityStandardWCAG}, "missing Version"},
		{&SoftwareQualityReportsGetAccessibilityReportOptions{ProjectKey: "p1", Standard: SoftwareQualityReportsAccessibilityStandardWCAG, Version: "3.0"}, "invalid Version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, resp, err := svc.GetAccessibilityReport(context.Background(), tt.opt)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Nil(t, resp)
		})
	}
}
