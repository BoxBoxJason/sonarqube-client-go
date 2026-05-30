package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectBadges_Measure(t *testing.T) {
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/project_badges/measure", http.StatusOK, "image/svg+xml", []byte(`<svg>badge content</svg>`)))
	client := newTestClient(t, server.URL)

	opt := &ProjectBadgesMeasureOptions{
		Project: "my-project",
		Metric:  "coverage",
	}

	result, resp, err := client.ProjectBadges.Measure(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "<svg>badge content</svg>", *result)
}

func TestProjectBadges_Measure_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectBadges.Measure(context.Background(), nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, _, err = client.ProjectBadges.Measure(context.Background(), &ProjectBadgesMeasureOptions{
		Metric: "coverage",
	})
	assert.Error(t, err)

	// Missing Metric should fail validation.
	_, _, err = client.ProjectBadges.Measure(context.Background(), &ProjectBadgesMeasureOptions{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Invalid Metric should fail validation.
	_, _, err = client.ProjectBadges.Measure(context.Background(), &ProjectBadgesMeasureOptions{
		Project: "my-project",
		Metric:  "invalid_metric",
	})
	assert.Error(t, err)
}

func TestProjectBadges_QualityGate(t *testing.T) {
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/project_badges/quality_gate", http.StatusOK, "image/svg+xml", []byte(`<svg>quality gate badge</svg>`)))
	client := newTestClient(t, server.URL)

	opt := &ProjectBadgesQualityGateOptions{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBadges.QualityGate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestProjectBadges_QualityGate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectBadges.QualityGate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, _, err = client.ProjectBadges.QualityGate(context.Background(), &ProjectBadgesQualityGateOptions{})
	assert.Error(t, err)
}

func TestProjectBadges_RenewToken(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_badges/renew_token", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &ProjectBadgesRenewTokenOptions{
		Project: "my-project",
	}

	resp, err := client.ProjectBadges.RenewToken(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBadges_RenewToken_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBadges.RenewToken(context.Background(), nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBadges.RenewToken(context.Background(), &ProjectBadgesRenewTokenOptions{})
	assert.Error(t, err)
}

func TestProjectBadges_Token(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_badges/token", http.StatusOK, &ProjectBadgesToken{
		Token: "abc123def456",
	}))
	client := newTestClient(t, server.URL)

	opt := &ProjectBadgesTokenOptions{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBadges.Token(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "abc123def456", result.Token)
}

func TestProjectBadges_Token_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectBadges.Token(context.Background(), nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, _, err = client.ProjectBadges.Token(context.Background(), &ProjectBadgesTokenOptions{})
	assert.Error(t, err)
}

func TestProjectBadges_ValidateMeasureOpt_AllMetrics(t *testing.T) {
	client := newLocalhostClient(t)

	validMetrics := []string{
		"coverage",
		"duplicated_lines_density",
		"ncloc",
		"alert_status",
		"security_hotspots",
		"bugs",
		"code_smells",
		"vulnerabilities",
		"sqale_rating",
		"reliability_rating",
		"security_rating",
		"sqale_index",
		"software_quality_reliability_issues",
		"software_quality_maintainability_issues",
		"software_quality_security_issues",
		"software_quality_maintainability_rating",
		"software_quality_reliability_rating",
		"software_quality_security_rating",
		"software_quality_maintainability_remediation_effort",
	}

	for _, metric := range validMetrics {
		err := client.ProjectBadges.ValidateMeasureOpt(&ProjectBadgesMeasureOptions{
			Project: "my-project",
			Metric:  metric,
		})
		assert.NoError(t, err, "expected nil error for metric '%s'", metric)
	}
}
