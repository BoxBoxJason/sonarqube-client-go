package sonar

import (
	"encoding/json"
	"fmt"
)

const (
	// defaultBaseURL is the default base URL for the SonarQube API.
	defaultBaseURL = "http://localhost:9000/api/"
	// defaultUserAgent is the default User-Agent header value.
	defaultUserAgent = "sonarqube-client-go"

	// MaxPageSize is the maximum allowed page size for pagination.
	MaxPageSize = 500
	// MinPageSize is the minimum allowed page size for pagination.
	MinPageSize = 1

	// MaxLinkNameLength is the maximum length for a project link name.
	MaxLinkNameLength = 128
	// MaxLinkURLLength is the maximum length for a project link URL.
	MaxLinkURLLength = 2048
	// MaxTokenNameLength is the maximum length for a user token name.
	MaxTokenNameLength = 100
	// MaxBranchNameLength is the maximum length for a branch name.
	MaxBranchNameLength = 255

	// LanguageAzureResourceManager is the language key for Azure Resource Manager.
	LanguageAzureResourceManager = "azureresourcemanager"
	// LanguageCloudFormation is the language key for CloudFormation.
	LanguageCloudFormation = "cloudformation"
	// LanguageCS is the language key for C#.
	LanguageCS = "cs"
	// LanguageCSS is the language key for CSS.
	LanguageCSS = "css"
	// LanguageDocker is the language key for Docker.
	LanguageDocker = "docker"
	// LanguageFlex is the language key for Flex.
	LanguageFlex = "flex"
	// LanguageGo is the language key for Go.
	LanguageGo = "go"
	// LanguageIPYNB is the language key for Jupyter Notebooks.
	LanguageIPYNB = "ipynb"
	// LanguageJava is the language key for Java.
	LanguageJava = "java"
	// LanguageJS is the language key for JavaScript.
	LanguageJS = "js"
	// LanguageJSON is the language key for JSON.
	LanguageJSON = "json"
	// LanguageJSP is the language key for JSP.
	LanguageJSP = "jsp"
	// LanguageKotlin is the language key for Kotlin.
	LanguageKotlin = "kotlin"
	// LanguageKubernetes is the language key for Kubernetes.
	LanguageKubernetes = "kubernetes"
	// LanguagePHP is the language key for PHP.
	LanguagePHP = "php"
	// LanguagePython is the language key for Python.
	LanguagePython = "py"
	// LanguageRuby is the language key for Ruby On Rails.
	LanguageRuby = "ruby"
	// LanguageRust is the language key for Rust.
	LanguageRust = "rust"
	// LanguageScala is the language key for Scala.
	LanguageScala = "scala"
	// LanguageSecrets is the language key for Secrets.
	LanguageSecrets = "secrets"
	// LanguageTerraform is the language key for Terraform.
	LanguageTerraform = "terraform"
	// LanguageText is the language key for Text files.
	LanguageText = "text"
	// LanguageTypeScript is the language key for TypeScript.
	LanguageTypeScript = "ts"
	// LanguageVBNet is the language key for Visual Basic .NET.
	LanguageVBNet = "vbnet"
	// LanguageWeb is the language key for Web files.
	LanguageWeb = "web"
	// LanguageXML is the language key for XML.
	LanguageXML = "xml"
	// LanguageYAML is the language key for YAML.
	LanguageYAML = "yaml"

	// CleanCodeAttributeCategoryAdaptable is the Clean Code attribute category for adaptability.
	CleanCodeAttributeCategoryAdaptable = "ADAPTABLE"
	// CleanCodeAttributeCategoryConsistent is the Clean Code attribute category for consistency.
	CleanCodeAttributeCategoryConsistent = "CONSISTENT"
	// CleanCodeAttributeCategoryIntentional is the Clean Code attribute category for intentionality.
	CleanCodeAttributeCategoryIntentional = "INTENTIONAL"
	// CleanCodeAttributeCategoryResponsible is the Clean Code attribute category for responsibility.
	CleanCodeAttributeCategoryResponsible = "RESPONSIBLE"

	// CleanCodeAttributeConventional is the Clean Code attribute for conventional code.
	CleanCodeAttributeConventional = "CONVENTIONAL"
	// CleanCodeAttributeFormatted is the Clean Code attribute for formatted code.
	CleanCodeAttributeFormatted = "FORMATTED"
	// CleanCodeAttributeIdentifiable is the Clean Code attribute for identifiable code.
	CleanCodeAttributeIdentifiable = "IDENTIFIABLE"
	// CleanCodeAttributeClear is the Clean Code attribute for clear code.
	CleanCodeAttributeClear = "CLEAR"
	// CleanCodeAttributeComplete is the Clean Code attribute for complete code.
	CleanCodeAttributeComplete = "COMPLETE"
	// CleanCodeAttributeEfficient is the Clean Code attribute for efficient code.
	CleanCodeAttributeEfficient = "EFFICIENT"
	// CleanCodeAttributeLogical is the Clean Code attribute for logical code.
	CleanCodeAttributeLogical = "LOGICAL"
	// CleanCodeAttributeDistinct is the Clean Code attribute for distinct code.
	CleanCodeAttributeDistinct = "DISTINCT"
	// CleanCodeAttributeFocused is the Clean Code attribute for focused code.
	CleanCodeAttributeFocused = "FOCUSED"
	// CleanCodeAttributeModular is the Clean Code attribute for modular code.
	CleanCodeAttributeModular = "MODULAR"
	// CleanCodeAttributeTested is the Clean Code attribute for tested code.
	CleanCodeAttributeTested = "TESTED"
	// CleanCodeAttributeLawful is the Clean Code attribute for lawful code.
	CleanCodeAttributeLawful = "LAWFUL"
	// CleanCodeAttributeRespectful is the Clean Code attribute for respectful code.
	CleanCodeAttributeRespectful = "RESPECTFUL"
	// CleanCodeAttributeTrustworthy is the Clean Code attribute for trustworthy code.
	CleanCodeAttributeTrustworthy = "TRUSTWORTHY"

	// OwaspCategoryA1 is the OWASP category for Injection.
	OwaspCategoryA1 = "a1"
	// OwaspCategoryA2 is the OWASP category for Broken Authentication.
	OwaspCategoryA2 = "a2"
	// OwaspCategoryA3 is the OWASP category for Sensitive Data Exposure.
	OwaspCategoryA3 = "a3"
	// OwaspCategoryA4 is the OWASP category for XML External Entities (XXE).
	OwaspCategoryA4 = "a4"
	// OwaspCategoryA5 is the OWASP category for Broken Access Control.
	OwaspCategoryA5 = "a5"
	// OwaspCategoryA6 is the OWASP category for Security Misconfiguration.
	OwaspCategoryA6 = "a6"
	// OwaspCategoryA7 is the OWASP category for Cross-Site Scripting (XSS).
	OwaspCategoryA7 = "a7"
	// OwaspCategoryA8 is the OWASP category for Insecure Deserialization.
	OwaspCategoryA8 = "a8"
	// OwaspCategoryA9 is the OWASP category for Using Components with Known Vulnerabilities.
	OwaspCategoryA9 = "a9"
	// OwaspCategoryA10 is the OWASP category for Insufficient Logging & Monitoring.
	OwaspCategoryA10 = "a10"

	// OwaspMobileCategoryM1 is the OWASP Mobile category for Improper Platform Usage.
	OwaspMobileCategoryM1 = "m1"
	// OwaspMobileCategoryM2 is the OWASP Mobile category for Insecure Data Storage.
	OwaspMobileCategoryM2 = "m2"
	// OwaspMobileCategoryM3 is the OWASP Mobile category for Insecure Communication.
	OwaspMobileCategoryM3 = "m3"
	// OwaspMobileCategoryM4 is the OWASP Mobile category for Insecure Authentication.
	OwaspMobileCategoryM4 = "m4"
	// OwaspMobileCategoryM5 is the OWASP Mobile category for Insufficient Cryptography.
	OwaspMobileCategoryM5 = "m5"
	// OwaspMobileCategoryM6 is the OWASP Mobile category for Insecure Authorization.
	OwaspMobileCategoryM6 = "m6"
	// OwaspMobileCategoryM7 is the OWASP Mobile category for Client Code Quality.
	OwaspMobileCategoryM7 = "m7"
	// OwaspMobileCategoryM8 is the OWASP Mobile category for Code Tampering.
	OwaspMobileCategoryM8 = "m8"
	// OwaspMobileCategoryM9 is the OWASP Mobile category for Reverse Engineering.
	OwaspMobileCategoryM9 = "m9"
	// OwaspMobileCategoryM10 is the OWASP Mobile category for Extraneous Functionality.
	OwaspMobileCategoryM10 = "m10"

	// SansTop25CategoryInsecureInteraction is the SANS Top 25 category for Insecure Interaction.
	SansTop25CategoryInsecureInteraction = "insecure-interaction"
	// SansTop25CategoryRiskyResource is the SANS Top 25 category for Risky Resource.
	SansTop25CategoryRiskyResource = "risky-resource"
	// SansTop25CategoryPorousDefenses is the SANS Top 25 category for Porous Defenses.
	SansTop25CategoryPorousDefenses = "porous-defenses"

	// SoftwareQualityMaintainability is the software quality characteristic for maintainability.
	SoftwareQualityMaintainability = "MAINTAINABILITY"
	// SoftwareQualityReliability is the software quality characteristic for reliability.
	SoftwareQualityReliability = "RELIABILITY"
	// SoftwareQualitySecurity is the software quality characteristic for security.
	SoftwareQualitySecurity = "SECURITY"

	// InheritanceTypeNone represents no inheritance.
	InheritanceTypeNone = "NONE"
	// InheritanceTypeInherited represents inherited.
	InheritanceTypeInherited = "INHERITED"
	// InheritanceTypeOverrides represents overrides.
	InheritanceTypeOverrides = "OVERRIDES"

	// SelectionFilterAll represents the "all" selection filter.
	SelectionFilterAll = "all"
	// SelectionFilterSelected represents the "selected" selection filter.
	SelectionFilterSelected = "selected"
	// SelectionFilterDeselected represents the "deselected" selection filter.
	SelectionFilterDeselected = "deselected"
)

type authType int

const (
	basicAuth authType = iota
	oAuthToken
	privateToken
)

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedLanguages is the set of supported programming languages.
	allowedLanguages = map[string]struct{}{
		LanguageAzureResourceManager: {},
		LanguageCloudFormation:       {},
		LanguageCS:                   {},
		LanguageCSS:                  {},
		LanguageDocker:               {},
		LanguageFlex:                 {},
		LanguageGo:                   {},
		LanguageIPYNB:                {},
		LanguageJava:                 {},
		LanguageJS:                   {},
		LanguageJSON:                 {},
		LanguageJSP:                  {},
		LanguageKotlin:               {},
		LanguageKubernetes:           {},
		LanguagePHP:                  {},
		LanguagePython:               {},
		LanguageRuby:                 {},
		LanguageRust:                 {},
		LanguageScala:                {},
		LanguageSecrets:              {},
		LanguageTerraform:            {},
		LanguageText:                 {},
		LanguageTypeScript:           {},
		LanguageVBNet:                {},
		LanguageWeb:                  {},
		LanguageXML:                  {},
		LanguageYAML:                 {},
	}

	// allowedRuleSeverities is the set of supported severity levels.
	allowedRuleSeverities = map[string]struct{}{
		RuleSeverityBlocker:  {},
		RuleSeverityCritical: {},
		RuleSeverityMajor:    {},
		RuleSeverityMinor:    {},
		RuleSeverityInfo:     {},
	}

	// allowedRuleImpactSeverities is the set of supported impact severity levels.
	allowedRuleImpactSeverities = map[string]struct{}{
		RuleImpactSeverityBlocker: {},
		RuleImpactSeverityHigh:    {},
		RuleImpactSeverityMedium:  {},
		RuleImpactSeverityLow:     {},
		RuleImpactSeverityInfo:    {},
	}

	// allowedCleanCodeAttributesCategories is the set of supported Clean Code attribute categories.
	allowedCleanCodeAttributesCategories = map[string]struct{}{
		CleanCodeAttributeCategoryAdaptable:   {},
		CleanCodeAttributeCategoryConsistent:  {},
		CleanCodeAttributeCategoryIntentional: {},
		CleanCodeAttributeCategoryResponsible: {},
	}

	// allowedCleanCodeAttributes is the set of supported Clean Code attributes.
	allowedCleanCodeAttributes = map[string]struct{}{
		CleanCodeAttributeConventional: {},
		CleanCodeAttributeFormatted:    {},
		CleanCodeAttributeIdentifiable: {},
		CleanCodeAttributeClear:        {},
		CleanCodeAttributeComplete:     {},
		CleanCodeAttributeEfficient:    {},
		CleanCodeAttributeLogical:      {},
		CleanCodeAttributeDistinct:     {},
		CleanCodeAttributeFocused:      {},
		CleanCodeAttributeModular:      {},
		CleanCodeAttributeTested:       {},
		CleanCodeAttributeLawful:       {},
		CleanCodeAttributeRespectful:   {},
		CleanCodeAttributeTrustworthy:  {},
	}

	// allowedImpactSoftwareQualities is the set of supported impact software qualities.
	allowedImpactSoftwareQualities = map[string]struct{}{
		SoftwareQualityMaintainability: {},
		SoftwareQualityReliability:     {},
		SoftwareQualitySecurity:        {},
	}

	// allowedInheritanceTypes is the set of supported inheritance types.
	allowedInheritanceTypes = map[string]struct{}{
		InheritanceTypeNone:      {},
		InheritanceTypeInherited: {},
		InheritanceTypeOverrides: {},
	}

	// allowedOwaspCategories is the set of supported OWASP categories.
	allowedOwaspCategories = map[string]struct{}{
		OwaspCategoryA1:  {},
		OwaspCategoryA2:  {},
		OwaspCategoryA3:  {},
		OwaspCategoryA4:  {},
		OwaspCategoryA5:  {},
		OwaspCategoryA6:  {},
		OwaspCategoryA7:  {},
		OwaspCategoryA8:  {},
		OwaspCategoryA9:  {},
		OwaspCategoryA10: {},
	}

	// allowedOwaspMobileCategories is the set of supported OWASP Mobile categories.
	allowedOwaspMobileCategories = map[string]struct{}{
		OwaspMobileCategoryM1:  {},
		OwaspMobileCategoryM2:  {},
		OwaspMobileCategoryM3:  {},
		OwaspMobileCategoryM4:  {},
		OwaspMobileCategoryM5:  {},
		OwaspMobileCategoryM6:  {},
		OwaspMobileCategoryM7:  {},
		OwaspMobileCategoryM8:  {},
		OwaspMobileCategoryM9:  {},
		OwaspMobileCategoryM10: {},
	}

	// allowedRulesStatuses is the set of supported statuses.
	allowedRulesStatuses = map[string]struct{}{
		RuleStatusReady:      {},
		RuleStatusDeprecated: {},
		RuleStatusRemoved:    {},
		RuleStatusBeta:       {},
	}

	// allowedRulesExistingStatuses is the set of supported existing statuses.
	allowedRulesExistingStatuses = map[string]struct{}{
		RuleStatusReady:      {},
		RuleStatusDeprecated: {},
		RuleStatusBeta:       {},
	}

	// allowedRulesTypes is the set of supported rule types.
	allowedRulesTypes = map[string]struct{}{
		RuleTypeCodeSmell:       {},
		RuleTypeBug:             {},
		RuleTypeVulnerability:   {},
		RuleTypeSecurityHotspot: {},
	}

	// allowedSansTop25Categories is the set of supported SANS Top 25 categories.
	allowedSansTop25Categories = map[string]struct{}{
		SansTop25CategoryInsecureInteraction: {},
		SansTop25CategoryRiskyResource:       {},
		SansTop25CategoryPorousDefenses:      {},
	}

	// allowedSelectedFilters is the set of supported selected filters.
	allowedSelectedFilters = map[string]struct{}{
		SelectionFilterAll:        {},
		SelectionFilterSelected:   {},
		SelectionFilterDeselected: {},
	}

	// allowedIssueTypes is the set of supported issue types.
	allowedIssueTypes = map[string]struct{}{
		RuleTypeCodeSmell:       {},
		RuleTypeBug:             {},
		RuleTypeVulnerability:   {},
		RuleTypeSecurityHotspot: {},
	}

	// allowedIssueTransitions is the set of supported issue transitions.
	allowedIssueTransitions = map[string]struct{}{
		IssueTransitionConfirm:           {},
		IssueTransitionUnconfirm:         {},
		IssueTransitionReopen:            {},
		IssueTransitionResolve:           {},
		IssueTransitionFalsePositive:     {},
		IssueTransitionWontFix:           {},
		IssueTransitionAccept:            {},
		IssueTransitionClose:             {},
		IssueTransitionResolveAsReviewed: {},
		IssueTransitionResetAsReviewed:   {},
	}

	// allowedIssueStatuses is the set of supported issue statuses.
	allowedIssueStatuses = map[string]struct{}{
		IssueStatusOpen:          {},
		IssueStatusConfirmed:     {},
		IssueStatusFalsePositive: {},
		IssueStatusAccepted:      {},
		IssueStatusFixed:         {},
		IssueStatusInSandbox:     {},
	}

	// allowedIssueResolutions is the set of supported issue resolutions.
	allowedIssueResolutions = map[string]struct{}{
		IssueResolutionFixed:         {},
		IssueResolutionRemoved:       {},
		IssueResolutionFalsePositive: {},
		IssueResolutionWontFix:       {},
	}

	// allowedIssueScopes is the set of supported issue scopes.
	allowedIssueScopes = map[string]struct{}{
		IssueScopeMain: {},
		IssueScopeTest: {},
	}
)

// Paging is used in many APIs.
type Paging struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

// PaginationArgs contains common pagination parameters for API requests.
type PaginationArgs struct {
	// Page is the response page number. Must be strictly greater than 0.
	Page int64 `url:"p,omitempty"`
	// PageSize is the response page size. Must be greater than 0 and less than or equal to 500.
	PageSize int64 `url:"ps,omitempty"`
}

// Validate validates the pagination arguments.
func (p *PaginationArgs) Validate() error {
	return ValidatePagination(p.Page, p.PageSize)
}

// =============================================
// V2 API COMMON TYPES
// =============================================

const (
	// v2BasePath is the base path segment for V2 API endpoints, appended to the
	// client base URL (e.g. "http://localhost:9000/api/" + "v2/" →
	// "http://localhost:9000/api/v2/").
	v2BasePath = "v2/"
)

// PageResponseV2 represents the pagination information returned by V2 API endpoints.
type PageResponseV2 struct {
	// PageIndex is the 1-based page index.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the number of items per page.
	PageSize int32 `json:"pageSize,omitempty"`
	// Total is the total number of items.
	Total int32 `json:"total,omitempty"`
}

// PaginationParamsV2 contains common pagination query parameters for V2 API requests.
type PaginationParamsV2 struct {
	// PageIndex is the 1-based page index. Default is 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the number of results per page. A value of 0 will only return
	// pagination information. Default is 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// Validate validates the V2 pagination parameters.
func (p *PaginationParamsV2) Validate() error {
	if p.PageIndex != 0 && p.PageIndex < 1 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	if p.PageSize != 0 && (p.PageSize < 0 || p.PageSize > MaxPageSize) {
		return NewValidationError("PageSize", fmt.Sprintf("must be between 0 and %d", MaxPageSize), ErrOutOfRange)
	}

	return nil
}

// UpdateFieldListStringV2 represents a field that can be explicitly set or cleared
// in a V2 PATCH request. When used as a pointer field with omitempty, a nil pointer
// means the field is not included in the request (no change). A non-nil pointer
// serializes the Value slice directly as a JSON array (or null if Value is nil).
type UpdateFieldListStringV2 struct {
	// Value is the list of string values.
	Value []string
	// Defined indicates whether this field has been explicitly set.
	Defined bool
}

// MarshalJSON implements the json.Marshaler interface. It serializes the
// UpdateFieldListStringV2 as the bare Value slice so that merge-patch+json
// requests send `"field":["a","b"]` instead of `"field":{"value":...,"defined":...}`.
func (u UpdateFieldListStringV2) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(u.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UpdateFieldListStringV2 value: %w", err)
	}

	return data, nil
}
