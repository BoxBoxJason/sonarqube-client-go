package sonar

import (
	"fmt"
	"net/http"
)

const (
	// MaxRuleKeyLengthV2 is the maximum length for a custom rule key in V2 API.
	MaxRuleKeyLengthV2 = 200
	// MaxRuleNameLengthV2 is the maximum length for a custom rule name in V2 API.
	MaxRuleNameLengthV2 = 200
	// MaxTemplateKeyLengthV2 is the maximum length for a rule template key in V2 API.
	MaxTemplateKeyLengthV2 = 200
)

// CleanCodePolicyService handles communication with the Clean Code Policy
// related methods of the SonarQube V2 API.
type CleanCodePolicyService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// RuleParameterV2 represents a parameter for a custom rule.
type RuleParameterV2 struct {
	// Key is the parameter key.
	Key string `json:"key,omitempty"`
	// HtmlDescription is the HTML description of the parameter (read-only).
	HtmlDescription string `json:"htmlDescription,omitempty"`
	// DefaultValue is the default value for the parameter.
	DefaultValue string `json:"defaultValue,omitempty"`
	// Type is the parameter type (read-only).
	Type string `json:"type,omitempty"`
}

// RuleV2 represents a rule returned by V2 API endpoints.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type RuleV2 struct {
	// Id is the rule's unique identifier.
	Id string `json:"id,omitempty"`
	// Key is the rule key.
	Key string `json:"key,omitempty"`
	// RepositoryKey is the repository key.
	RepositoryKey string `json:"repositoryKey,omitempty"`
	// Name is the rule name.
	Name string `json:"name,omitempty"`
	// Severity is the rule severity.
	Severity string `json:"severity,omitempty"`
	// Type is the rule type.
	Type string `json:"type,omitempty"`
	// Impacts is the list of software quality impacts.
	Impacts []RuleImpact `json:"impacts,omitempty"`
	// CleanCodeAttribute is the clean code attribute for the rule.
	CleanCodeAttribute string `json:"cleanCodeAttribute,omitempty"`
	// CleanCodeAttributeCategory is the clean code attribute category.
	CleanCodeAttributeCategory string `json:"cleanCodeAttributeCategory,omitempty"`
	// Status is the rule status.
	Status string `json:"status,omitempty"`
	// External indicates whether the rule is external.
	External bool `json:"external,omitempty"`
	// CreatedAt is the rule creation timestamp.
	CreatedAt string `json:"createdAt,omitempty"`
	// DescriptionSections contains the rule description sections.
	DescriptionSections []RuleDescriptionSection `json:"descriptionSections,omitempty"`
	// MarkdownDescription is the rule description in markdown format.
	MarkdownDescription string `json:"markdownDescription,omitempty"`
	// GapDescription is the gap description for the rule.
	GapDescription string `json:"gapDescription,omitempty"`
	// HtmlNote is the HTML note for the rule.
	HtmlNote string `json:"htmlNote,omitempty"`
	// MarkdownNote is the markdown note for the rule.
	MarkdownNote string `json:"markdownNote,omitempty"`
	// EducationPrinciples is the list of education principles.
	EducationPrinciples []string `json:"educationPrinciples,omitempty"`
	// Template indicates whether this rule is a template.
	Template bool `json:"template,omitempty"`
	// TemplateId is the ID of the template this rule is based on.
	TemplateId string `json:"templateId,omitempty"`
	// Tags is the list of user-defined tags.
	Tags []string `json:"tags,omitempty"`
	// SystemTags is the list of system tags.
	SystemTags []string `json:"systemTags,omitempty"`
	// LanguageKey is the language key.
	LanguageKey string `json:"languageKey,omitempty"`
	// LanguageName is the language display name.
	LanguageName string `json:"languageName,omitempty"`
	// Parameters is the list of rule parameters.
	Parameters []RuleParameterV2 `json:"parameters,omitempty"`
	// RemediationFunctionType is the remediation function type.
	RemediationFunctionType string `json:"remediationFunctionType,omitempty"`
	// RemediationFunctionGapMultiplier is the gap multiplier for remediation.
	RemediationFunctionGapMultiplier string `json:"remediationFunctionGapMultiplier,omitempty"`
	// RemediationFunctionBaseEffort is the base effort for remediation.
	RemediationFunctionBaseEffort string `json:"remediationFunctionBaseEffort,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// CleanCodePolicyCreateRuleOptions contains parameters for creating a custom rule
// via the V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CleanCodePolicyCreateRuleOptions struct {
	// Key is the key of the custom rule to create (must include the repository).
	// This field is required. Maximum 200 characters.
	Key string `json:"key"`
	// TemplateKey is the key of the rule template to base the custom rule on.
	// This field is required. Maximum 200 characters.
	TemplateKey string `json:"templateKey"`
	// Name is the rule name.
	// This field is required. Maximum 200 characters.
	Name string `json:"name"`
	// MarkdownDescription is the rule description in markdown format.
	// This field is required.
	MarkdownDescription string `json:"markdownDescription"`
	// Impacts is the list of software quality impacts.
	// This field is required (at least one impact).
	Impacts []RuleImpact `json:"impacts"`
	// Status is the rule status. Default is "READY".
	// Valid values: BETA, DEPRECATED, READY, REMOVED.
	Status string `json:"status,omitempty"`
	// Parameters is the list of custom rule parameters.
	Parameters []RuleParameterV2 `json:"parameters,omitempty"`
	// CleanCodeAttribute is the clean code attribute for the rule.
	// Valid values: CONVENTIONAL, FORMATTED, IDENTIFIABLE, CLEAR, COMPLETE,
	// EFFICIENT, LOGICAL, DISTINCT, FOCUSED, MODULAR, TESTED, LAWFUL,
	// RESPECTFUL, TRUSTWORTHY.
	CleanCodeAttribute string `json:"cleanCodeAttribute,omitempty"`
	// Severity is the rule severity.
	Severity string `json:"severity,omitempty"`
	// Type is the rule type.
	// Valid values: CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT.
	Type string `json:"type,omitempty"`
}

// -----------------------------------------------------------------------------
// Allowed Values
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// validateRuleImpacts validates the impacts list in a create rule request.
func validateRuleImpacts(impacts []RuleImpact) error {
	if len(impacts) == 0 {
		return NewValidationError("Impacts", "at least one impact is required", ErrMissingRequired)
	}

	for impactIdx, impact := range impacts {
		err := ValidateRequired(impact.SoftwareQuality, fmt.Sprintf("Impacts[%d].SoftwareQuality", impactIdx))
		if err != nil {
			return err
		}

		err = IsValueAuthorized(impact.SoftwareQuality, allowedImpactSoftwareQualities, fmt.Sprintf("Impacts[%d].SoftwareQuality", impactIdx))
		if err != nil {
			return err
		}

		err = ValidateRequired(impact.Severity, fmt.Sprintf("Impacts[%d].Severity", impactIdx))
		if err != nil {
			return err
		}

		err = IsValueAuthorized(impact.Severity, allowedImpactSeverities, fmt.Sprintf("Impacts[%d].Severity", impactIdx))
		if err != nil {
			return err
		}
	}

	return nil
}

// validateRuleRequiredFields validates the required string fields of a create rule request.
func validateRuleRequiredFields(opt *CleanCodePolicyCreateRuleOptions) error {
	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxRuleKeyLengthV2, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.TemplateKey, "TemplateKey")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.TemplateKey, MaxTemplateKeyLengthV2, "TemplateKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxRuleNameLengthV2, "Name")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.MarkdownDescription, "MarkdownDescription")
}

// ValidateCreateRuleRequest validates the CleanCodePolicyCreateRuleOptions.
func (s *CleanCodePolicyService) ValidateCreateRuleRequest(opt *CleanCodePolicyCreateRuleOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := validateRuleRequiredFields(opt)
	if err != nil {
		return err
	}

	err = validateRuleImpacts(opt.Impacts)
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Severity, allowedSeverities, "Severity")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Status, allowedRulesStatuses, "Status")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.CleanCodeAttribute, allowedCleanCodeAttributes, "CleanCodeAttribute")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Type, allowedRulesTypes, "Type")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CreateRule creates a custom rule based on a template.
// Requires the 'Administer Quality Profiles' permission.
func (s *CleanCodePolicyService) CreateRule(opt *CleanCodePolicyCreateRuleOptions) (*RuleV2, *http.Response, error) {
	err := s.ValidateCreateRuleRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPost, "clean-code-policy/rules", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RuleV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
