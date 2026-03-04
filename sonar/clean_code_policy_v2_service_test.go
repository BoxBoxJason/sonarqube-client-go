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
	request := &CleanCodePolicyCreateRuleRequestV2{
		Key:                 "custom:my_rule",
		TemplateKey:         "java:S100",
		Name:                "My Custom Rule",
		MarkdownDescription: "This is a custom rule",
		Impacts: []RuleImpact{
			{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"},
		},
		Status:             "READY",
		CleanCodeAttribute: "CLEAR",
		Type:               "CODE_SMELL",
	}
	response := RuleV2{
		Id:                 "rule-1",
		Key:                "custom:my_rule",
		RepositoryKey:      "custom",
		Name:               "My Custom Rule",
		Status:             "READY",
		CleanCodeAttribute: "CLEAR",
		Type:               "CODE_SMELL",
		Impacts: []RuleImpact{
			{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"},
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
	assert.Equal(t, "CLEAR", result.CleanCodeAttribute)
	assert.Equal(t, "CODE_SMELL", result.Type)
	assert.Len(t, result.Impacts, 1)
	assert.Equal(t, "MAINTAINABILITY", result.Impacts[0].SoftwareQuality)
	assert.Equal(t, "HIGH", result.Impacts[0].Severity)
	assert.Equal(t, "java", result.LanguageKey)
}

func TestCleanCodePolicyV2_CreateRule_MinimalRequest(t *testing.T) {
	request := &CleanCodePolicyCreateRuleRequestV2{
		Key:                 "custom:min_rule",
		TemplateKey:         "java:S100",
		Name:                "Minimal Rule",
		MarkdownDescription: "Minimal description",
		Impacts: []RuleImpact{
			{SoftwareQuality: "SECURITY", Severity: "LOW"},
		},
	}
	response := RuleV2{
		Id:   "rule-2",
		Key:  "custom:min_rule",
		Name: "Minimal Rule",
		Impacts: []RuleImpact{
			{SoftwareQuality: "SECURITY", Severity: "LOW"},
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
	request := &CleanCodePolicyCreateRuleRequestV2{
		Key:                 "custom:multi_impact",
		TemplateKey:         "java:S100",
		Name:                "Multi Impact Rule",
		MarkdownDescription: "Rule with multiple impacts",
		Impacts: []RuleImpact{
			{SoftwareQuality: "MAINTAINABILITY", Severity: "MEDIUM"},
			{SoftwareQuality: "RELIABILITY", Severity: "HIGH"},
			{SoftwareQuality: "SECURITY", Severity: "BLOCKER"},
		},
	}
	response := RuleV2{
		Id:  "rule-3",
		Key: "custom:multi_impact",
		Impacts: []RuleImpact{
			{SoftwareQuality: "MAINTAINABILITY", Severity: "MEDIUM"},
			{SoftwareQuality: "RELIABILITY", Severity: "HIGH"},
			{SoftwareQuality: "SECURITY", Severity: "BLOCKER"},
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
	request := &CleanCodePolicyCreateRuleRequestV2{
		Key:                 "custom:params_rule",
		TemplateKey:         "java:S100",
		Name:                "Parameterized Rule",
		MarkdownDescription: "Rule with parameters",
		Impacts: []RuleImpact{
			{SoftwareQuality: "MAINTAINABILITY", Severity: "LOW"},
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
		opt  *CleanCodePolicyCreateRuleRequestV2
	}{
		{"nil opt", nil},
		{"missing key", &CleanCodePolicyCreateRuleRequestV2{
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"missing template key", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"missing name", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"missing markdown description", &CleanCodePolicyCreateRuleRequestV2{
			Key:         "custom:test",
			TemplateKey: "java:S100",
			Name:        "Test",
			Impacts:     []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"missing impacts", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
		}},
		{"empty impacts", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{},
		}},
		{"key too long", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 strings.Repeat("a", MaxRuleKeyLengthV2+1),
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"template key too long", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         strings.Repeat("a", MaxTemplateKeyLengthV2+1),
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"name too long", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                strings.Repeat("a", MaxRuleNameLengthV2+1),
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
		}},
		{"invalid status", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
			Status:              "INVALID",
		}},
		{"invalid clean code attribute", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
			CleanCodeAttribute:  "INVALID",
		}},
		{"invalid type", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
			Type:                "INVALID",
		}},
		{"invalid impact software quality", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "INVALID", Severity: "HIGH"}},
		}},
		{"invalid impact severity", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "INVALID"}},
		}},
		{"missing impact software quality", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{Severity: "HIGH"}},
		}},
		{"missing impact severity", &CleanCodePolicyCreateRuleRequestV2{
			Key:                 "custom:test",
			TemplateKey:         "java:S100",
			Name:                "Test",
			MarkdownDescription: "desc",
			Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY"}},
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
		"CONVENTIONAL", "FORMATTED", "IDENTIFIABLE", "CLEAR", "COMPLETE",
		"EFFICIENT", "LOGICAL", "DISTINCT", "FOCUSED", "MODULAR",
		"TESTED", "LAWFUL", "RESPECTFUL", "TRUSTWORTHY",
	}

	for _, attr := range validAttributes {
		t.Run(attr, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleRequestV2{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				CleanCodeAttribute:  attr,
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllImpactSeverities(t *testing.T) {
	client := newLocalhostClient(t)

	validSeverities := []string{"INFO", "LOW", "MEDIUM", "HIGH", "BLOCKER"}

	for _, sev := range validSeverities {
		t.Run(sev, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleRequestV2{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: sev}},
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllRuleStatuses(t *testing.T) {
	client := newLocalhostClient(t)

	validStatuses := []string{"BETA", "DEPRECATED", "READY", "REMOVED"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleRequestV2{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				Status:              status,
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllRuleTypes(t *testing.T) {
	client := newLocalhostClient(t)

	validTypes := []string{"CODE_SMELL", "BUG", "VULNERABILITY", "SECURITY_HOTSPOT"}

	for _, ruleType := range validTypes {
		t.Run(ruleType, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleRequestV2{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				Type:                ruleType,
			})
			assert.NoError(t, err)
		})
	}
}

func TestCleanCodePolicyV2_CreateRule_AllSoftwareQualities(t *testing.T) {
	client := newLocalhostClient(t)

	validQualities := []string{"MAINTAINABILITY", "RELIABILITY", "SECURITY"}

	for _, quality := range validQualities {
		t.Run(quality, func(t *testing.T) {
			err := client.V2.CleanCodePolicy.ValidateCreateRuleRequest(&CleanCodePolicyCreateRuleRequestV2{
				Key:                 "custom:test",
				TemplateKey:         "java:S100",
				Name:                "Test",
				MarkdownDescription: "desc",
				Impacts:             []RuleImpact{{SoftwareQuality: quality, Severity: "HIGH"}},
			})
			assert.NoError(t, err)
		})
	}
}
