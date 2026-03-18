package sonar

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// CreateRule
// =============================================================================

func TestCleanCodePolicyV2_CreateRule(t *testing.T) {
	request := &CleanCodePolicyCreateRuleOptions{
		Key:                 "custom:my_rule",
		TemplateKey:         "java:S100",
		Name:                "My Custom Rule",
		MarkdownDescription: "This is a custom rule",
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh},
		},
		Status:             RuleStatusReady,
		CleanCodeAttribute: CleanCodeAttributeClear,
		Type:               RuleTypeCodeSmell,
	}
	response := RuleV2{
		Id:                 "rule-1",
		Key:                "custom:my_rule",
		RepositoryKey:      "custom",
		Name:               "My Custom Rule",
		Status:             RuleStatusReady,
		CleanCodeAttribute: CleanCodeAttributeClear,
		Type:               RuleTypeCodeSmell,
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh},
		},
		MarkdownDescription: "This is a custom rule",
		Template:            false,
		TemplateId:          "template-1",
		LanguageKey:         "java",
		LanguageName:        "Java",
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/clean-code-policy/rules", http.StatusOK, request, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.CleanCodePolicy.CreateRule(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "custom:my_rule", result.Key)
	assert.Equal(t, "My Custom Rule", result.Name)
	assert.Equal(t, CleanCodeAttributeClear, result.CleanCodeAttribute)
	assert.Equal(t, RuleTypeCodeSmell, result.Type)
	assert.Len(t, result.Impacts, 1)
	assert.Equal(t, SoftwareQualityMaintainability, result.Impacts[0].SoftwareQuality)
	assert.Equal(t, RuleImpactSeverityHigh, result.Impacts[0].Severity)
	assert.Equal(t, "java", result.LanguageKey)
}

func TestCleanCodePolicyV2_CreateRule_MinimalRequest(t *testing.T) {
	request := &CleanCodePolicyCreateRuleOptions{
		Key:                 "custom:min_rule",
		TemplateKey:         "java:S100",
		Name:                "Minimal Rule",
		MarkdownDescription: "Minimal description",
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualitySecurity, Severity: RuleImpactSeverityLow},
		},
	}
	response := RuleV2{
		Id:   "rule-2",
		Key:  "custom:min_rule",
		Name: "Minimal Rule",
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualitySecurity, Severity: RuleImpactSeverityLow},
		},
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/clean-code-policy/rules", http.StatusOK, request, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.CleanCodePolicy.CreateRule(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "custom:min_rule", result.Key)
	assert.Equal(t, "Minimal Rule", result.Name)
}

func TestCleanCodePolicyV2_CreateRule_MultipleImpacts(t *testing.T) {
	request := &CleanCodePolicyCreateRuleOptions{
		Key:                 "custom:multi_impact",
		TemplateKey:         "java:S100",
		Name:                "Multi Impact Rule",
		MarkdownDescription: "Rule with multiple impacts",
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityMedium},
			{SoftwareQuality: SoftwareQualityReliability, Severity: RuleImpactSeverityHigh},
			{SoftwareQuality: SoftwareQualitySecurity, Severity: RuleImpactSeverityBlocker},
		},
	}
	response := RuleV2{
		Id:  "rule-3",
		Key: "custom:multi_impact",
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityMedium},
			{SoftwareQuality: SoftwareQualityReliability, Severity: RuleImpactSeverityHigh},
			{SoftwareQuality: SoftwareQualitySecurity, Severity: RuleImpactSeverityBlocker},
		},
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/clean-code-policy/rules", http.StatusOK, request, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.CleanCodePolicy.CreateRule(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Impacts, 3)
}

func TestCleanCodePolicyV2_CreateRule_WithParameters(t *testing.T) {
	request := &CleanCodePolicyCreateRuleOptions{
		Key:                 "custom:params_rule",
		TemplateKey:         "java:S100",
		Name:                "Parameterized Rule",
		MarkdownDescription: "Rule with parameters",
		Impacts: []RuleImpact{
			{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityLow},
		},
		Parameters: []RuleParameterV2{
			{Key: "maxLines", DefaultValue: "100"},
			{Key: "threshold", DefaultValue: "0.8"},
		},
	}
	response := RuleV2{
		Id:  "rule-4",
		Key: "custom:params_rule",
		Parameters: []RuleParameterV2{
			{Key: "maxLines", DefaultValue: "100", Type: "INTEGER"},
			{Key: "threshold", DefaultValue: "0.8", Type: "FLOAT"},
		},
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/clean-code-policy/rules", http.StatusOK, request, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.CleanCodePolicy.CreateRule(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Parameters, 2)
	assert.Equal(t, "maxLines", result.Parameters[0].Key)
	assert.Equal(t, "INTEGER", result.Parameters[0].Type)
}

// =============================================================================
// CreateRule Validation
// =============================================================================

func TestCleanCodePolicyV2_CreateRule_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *CleanCodePolicyCreateRuleOptions
	}{
		{"nil opt", nil},
		{"missing key", &CleanCodePolicyCreateRuleOptions{
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"missing template key", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"missing name", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"missing markdown description", &CleanCodePolicyCreateRuleOptions{
			Key:         "custom:test",
			TemplateKey: "java:S100",
			Name:        "Test",
			Impacts:     []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"missing impacts", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
		}},
		{"empty impacts", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{},
		}},
		{"key too long", &CleanCodePolicyCreateRuleOptions{
			Key:                 strings.Repeat("a", MaxRuleKeyLengthV2+1),
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"template key too long", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         strings.Repeat("a", MaxTemplateKeyLengthV2+1),
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"name too long", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                strings.Repeat("a", MaxRuleNameLengthV2+1),
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
		}},
		{"invalid status", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
			Status:              "INVALID",
		}},
		{"invalid clean code attribute", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
			CleanCodeAttribute:  "INVALID",
		}},
		{"invalid type", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
			Type:                "INVALID",
		}},
		{"invalid impact software quality", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "INVALID", Severity: RuleImpactSeverityHigh}},
		}},
		{"invalid impact severity", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: "INVALID"}},
		}},
		{"missing impact software quality", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{Severity: RuleImpactSeverityHigh}},
		}},
		{"missing impact severity", &CleanCodePolicyCreateRuleOptions{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.CleanCodePolicy.CreateRule(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// CreateRule - All Valid Enum Values
// =============================================================================

func TestCleanCodePolicyV2_CreateRule_AllCleanCodeAttributes(t *testing.T) {
	client := newLocalhostClient(t)

	validAttributes := []string{
		CleanCodeAttributeConventional, CleanCodeAttributeFormatted, CleanCodeAttributeIdentifiable, CleanCodeAttributeClear, CleanCodeAttributeComplete,
		CleanCodeAttributeEfficient, CleanCodeAttributeLogical, CleanCodeAttributeDistinct, CleanCodeAttributeFocused, CleanCodeAttributeModular,
		CleanCodeAttributeTested, CleanCodeAttributeLawful, CleanCodeAttributeRespectful, CleanCodeAttributeTrustworthy,
	}

	for _, attr := range validAttributes {
		t.Run(attr, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleOptions{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
				CleanCodeAttribute:  attr,
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllImpactSeverities(t *testing.T) {
	client := newLocalhostClient(t)

	validSeverities := []string{RuleSeverityInfo, RuleImpactSeverityLow, RuleImpactSeverityMedium, RuleImpactSeverityHigh, RuleSeverityBlocker}

	for _, sev := range validSeverities {
		t.Run(sev, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleOptions{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: sev}},
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllRuleStatuses(t *testing.T) {
	client := newLocalhostClient(t)

	validStatuses := []string{RuleStatusBeta, RuleStatusDeprecated, RuleStatusReady, RuleStatusRemoved}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleOptions{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
				Status:              status,
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllRuleTypes(t *testing.T) {
	client := newLocalhostClient(t)

	validTypes := []string{RuleTypeCodeSmell, RuleTypeBug, RuleTypeVulnerability, RuleTypeSecurityHotspot}

	for _, ruleType := range validTypes {
		t.Run(ruleType, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleOptions{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: SoftwareQualityMaintainability, Severity: RuleImpactSeverityHigh}},
				Type:                ruleType,
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllSoftwareQualities(t *testing.T) {
	client := newLocalhostClient(t)

	validQualities := []string{SoftwareQualityMaintainability, SoftwareQualityReliability, SoftwareQualitySecurity}

	for _, quality := range validQualities {
		t.Run(quality, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleOptions{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: quality, Severity: RuleImpactSeverityHigh}},
			})
			assert.NoError(t, err)
		})
	}
}
