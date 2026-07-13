package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// GetSandboxSettings
// =============================================================================

func TestIssuesV2_GetSandboxSettings(t *testing.T) {
	response := IssuesSandboxSettings{
		Enabled:       true,
		DefaultValue:  false,
		AllowOverride: true,
		SoftwareQualities: []IssuesSandboxSoftwareQualityMapping{
			{SoftwareQuality: SoftwareQualityMaintainability, ImpactSeverities: []string{RuleImpactSeverityHigh}},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/issues/sandbox-settings", http.StatusOK, response))
	client := newTestClient(t, server.URL)
	svc := &IssuesV2Service{client: client}

	result, resp, err := svc.GetSandboxSettings(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Enabled)
	assert.True(t, result.AllowOverride)
	require.Len(t, result.SoftwareQualities, 1)
	assert.Equal(t, SoftwareQualityMaintainability, result.SoftwareQualities[0].SoftwareQuality)
}

// =============================================================================
// UpdateSandboxSettings
// =============================================================================

func TestIssuesV2_UpdateSandboxSettings(t *testing.T) {
	enabled := true
	response := IssuesSandboxSettings{
		Enabled: true,
		Types: []IssuesSandboxRuleTypeMapping{
			{Type: RuleTypeBug, Severities: []string{RuleSeverityBlocker}},
		},
	}
	server := newTestServer(t, mockPatchHandler(t, "/v2/issues/sandbox-settings", http.StatusOK,
		map[string]any{
			"enabled": true,
			"types": []map[string]any{
				{"type": RuleTypeBug, "severities": []string{RuleSeverityBlocker}},
			},
		}, response))
	client := newTestClient(t, server.URL)
	svc := &IssuesV2Service{client: client}

	result, resp, err := svc.UpdateSandboxSettings(context.Background(), &IssuesV2UpdateSandboxSettingsOptions{
		Enabled: &enabled,
		Types: []IssuesSandboxRuleTypeMapping{
			{Type: RuleTypeBug, Severities: []string{RuleSeverityBlocker}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Enabled)
	require.Len(t, result.Types, 1)
	assert.Equal(t, RuleTypeBug, result.Types[0].Type)
}

func TestIssuesV2_UpdateSandboxSettings_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IssuesV2Service{client: client}

	tests := []struct {
		opt  *IssuesV2UpdateSandboxSettingsOptions
		name string
	}{
		{nil, "nil opt"},
		{
			&IssuesV2UpdateSandboxSettingsOptions{
				SoftwareQualities: []IssuesSandboxSoftwareQualityMapping{{SoftwareQuality: "BOGUS", ImpactSeverities: []string{RuleImpactSeverityHigh}}},
			},
			"invalid software quality",
		},
		{
			&IssuesV2UpdateSandboxSettingsOptions{
				SoftwareQualities: []IssuesSandboxSoftwareQualityMapping{{SoftwareQuality: SoftwareQualityMaintainability}},
			},
			"missing impact severities",
		},
		{
			&IssuesV2UpdateSandboxSettingsOptions{
				Types: []IssuesSandboxRuleTypeMapping{{Type: "NOT_A_TYPE", Severities: []string{RuleSeverityBlocker}}},
			},
			"invalid rule type",
		},
		{
			&IssuesV2UpdateSandboxSettingsOptions{
				Types: []IssuesSandboxRuleTypeMapping{{Type: RuleTypeBug}},
			},
			"missing severities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, resp, err := svc.UpdateSandboxSettings(context.Background(), tt.opt)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Nil(t, resp)
		})
	}
}

// =============================================================================
// GetProjectSandboxSettings
// =============================================================================

func TestIssuesV2_GetProjectSandboxSettings(t *testing.T) {
	response := IssuesSandboxProjectSettings{
		Enabled:    true,
		Overridden: true,
		SoftwareQualities: []IssuesSandboxSoftwareQualityMapping{
			{SoftwareQuality: SoftwareQualitySecurity, ImpactSeverities: []string{RuleImpactSeverityBlocker}},
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/issues/sandbox-settings/my-project", http.StatusOK, response))
	client := newTestClient(t, server.URL)
	svc := &IssuesV2Service{client: client}

	result, resp, err := svc.GetProjectSandboxSettings(context.Background(), "my-project")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Enabled)
	assert.True(t, result.Overridden)
	require.Len(t, result.SoftwareQualities, 1)
	assert.Equal(t, SoftwareQualitySecurity, result.SoftwareQualities[0].SoftwareQuality)
}

func TestIssuesV2_GetProjectSandboxSettings_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IssuesV2Service{client: client}

	result, resp, err := svc.GetProjectSandboxSettings(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

// =============================================================================
// UpdateProjectSandboxSettings
// =============================================================================

func TestIssuesV2_UpdateProjectSandboxSettings(t *testing.T) {
	overridden := false
	response := IssuesSandboxProjectSettings{Enabled: false, Overridden: false}
	server := newTestServer(t, mockPatchHandler(t, "/v2/issues/sandbox-settings/my-project", http.StatusOK,
		map[string]any{"overridden": false}, response))
	client := newTestClient(t, server.URL)
	svc := &IssuesV2Service{client: client}

	result, resp, err := svc.UpdateProjectSandboxSettings(context.Background(), "my-project", &IssuesV2UpdateProjectSandboxSettingsOptions{
		Overridden: &overridden,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, result.Enabled)
	assert.False(t, result.Overridden)
}

func TestIssuesV2_UpdateProjectSandboxSettings_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IssuesV2Service{client: client}

	result, resp, err := svc.UpdateProjectSandboxSettings(context.Background(), "", &IssuesV2UpdateProjectSandboxSettingsOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = svc.UpdateProjectSandboxSettings(context.Background(), "my-project", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
