package sonargo

import (
	"fmt"
	"net/http"
)

const (
	// MaxRuleKeyLength is the maximum allowed length for rule keys and names.
	MaxRuleKeyLength = 200
	// MinSearchQueryLength is the minimum required length for search queries.
	MinSearchQueryLength = 2
)

// RulesService handles communication with the Rules related methods of the SonarQube API.
type RulesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// RulesApp contains metadata for rendering the 'Coding Rules' page.
type RulesApp struct {
	// Languages is a map listing languages keys and their associated display names.
	Languages map[string]string `json:"languages,omitempty"`
	// Statuses is a map of statuses keys and their associated display names.
	Statuses map[string]string `json:"statuses,omitempty"`
	// Characteristics is the list of available rule characteristics.
	Characteristics []RuleCharacteristic `json:"characteristics,omitempty"`
	// Repositories is the list of available rule repositories.
	Repositories []RuleRepository `json:"repositories,omitempty"`
	// CanWrite indicates if the current user has permission to modify rules.
	CanWrite bool `json:"canWrite,omitempty"`
}

// RuleCharacteristic represents a characteristic that can be associated with rules.
type RuleCharacteristic struct {
	// Key is the unique identifier of the characteristic.
	Key string `json:"key,omitempty"`
	// Name is the display name of the characteristic.
	Name string `json:"name,omitempty"`
	// Parent is the key of the parent characteristic, if any.
	Parent string `json:"parent,omitempty"`
}

// RuleRepository represents a rules repository.
type RuleRepository struct {
	// Key is the unique identifier of the repository.
	Key string `json:"key,omitempty"`
	// Language is the programming language of the repository.
	Language string `json:"language,omitempty"`
	// Name is the display name of the repository.
	Name string `json:"name,omitempty"`
}

// RulesCreate represents the response from creating a custom rule.
type RulesCreate struct {
	Rule Rule `json:"rule,omitzero"`
}

// Rule represents a SonarQube rule.
type Rule struct {
	// Key is the unique identifier of the rule.
	Key string `json:"key,omitempty"`
	// Severity indicates the severity level of the rule.
	Severity string `json:"severity,omitempty"`
	// CreatedAt is the timestamp when the rule was created.
	CreatedAt string `json:"createdAt,omitempty"`
	// UpdatedAt is the timestamp when the rule was last updated.
	UpdatedAt string `json:"updatedAt,omitempty"`
	// HTMLDesc is the HTML-formatted description of the rule.
	HTMLDesc string `json:"htmlDesc,omitempty"`
	// MdDesc is the Markdown-formatted description of the rule.
	MdDesc string `json:"mdDesc,omitempty"`
	// HTMLNote is the HTML-formatted note for the rule.
	HTMLNote string `json:"htmlNote,omitempty"`
	// MdNote is the Markdown-formatted note for the rule.
	MdNote string `json:"mdNote,omitempty"`
	// NoteLogin is the login of the user who created the note.
	NoteLogin string `json:"noteLogin,omitempty"`
	// InternalKey is the internal key used by the rule engine.
	InternalKey string `json:"internalKey,omitempty"`
	// Type is the rule type (CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT).
	Type string `json:"type,omitempty"`
	// TemplateKey is the key of the template rule, if this is a custom rule.
	TemplateKey string `json:"templateKey,omitempty"`
	// CleanCodeAttributeCategory is the category of the clean code attribute.
	CleanCodeAttributeCategory string `json:"cleanCodeAttributeCategory,omitempty"`
	// Lang is the language key of the rule.
	Lang string `json:"lang,omitempty"`
	// Scope is the scope of the rule (MAIN, TEST, ALL).
	Scope string `json:"scope,omitempty"`
	// Name is the display name of the rule.
	Name string `json:"name,omitempty"`
	// Status is the status of the rule (READY, DEPRECATED, BETA, REMOVED).
	Status string `json:"status,omitempty"`
	// Repo is the repository key of the rule.
	Repo string `json:"repo,omitempty"`
	// LangName is the display name of the language.
	LangName string `json:"langName,omitempty"`
	// CleanCodeAttribute is the clean code attribute of the rule.
	CleanCodeAttribute string `json:"cleanCodeAttribute,omitempty"`
	// Params is the list of parameters that can be configured for the rule.
	Params []RuleParam `json:"params,omitempty"`
	// SysTags is the list of system-defined tags.
	SysTags []string `json:"sysTags,omitempty"`
	// Tags is the list of user-defined tags.
	Tags []any `json:"tags,omitempty"`
	// Impacts is the list of impacts on software quality.
	Impacts []RuleImpact `json:"impacts,omitempty"`
	// IsTemplate indicates if this is a template rule that can be used to create custom rules.
	IsTemplate bool `json:"isTemplate,omitempty"`
	// IsExternal indicates if this is an external rule.
	IsExternal bool `json:"isExternal,omitempty"`
}

// RuleImpact represents the impact of a rule on software quality.
type RuleImpact struct {
	// Severity is the severity of the impact (HIGH, MEDIUM, LOW).
	Severity string `json:"severity,omitempty"`
	// SoftwareQuality is the software quality characteristic affected (MAINTAINABILITY, RELIABILITY, SECURITY).
	SoftwareQuality string `json:"softwareQuality,omitempty"`
}

// RuleParam represents a parameter that can be configured for a rule.
type RuleParam struct {
	// DefaultValue is the default value of the parameter.
	DefaultValue string `json:"defaultValue,omitempty"`
	// HTMLDesc is the HTML-formatted description of the parameter.
	HTMLDesc string `json:"htmlDesc,omitempty"`
	// Desc is the plain text description of the parameter.
	Desc string `json:"desc,omitempty"`
	// Key is the unique identifier of the parameter.
	Key string `json:"key,omitempty"`
	// Type is the data type of the parameter (STRING, TEXT, BOOLEAN, INTEGER, FLOAT).
	Type string `json:"type,omitempty"`
}

// RulesRepositories contains the list of available rule repositories.
type RulesRepositories struct {
	Repositories []RuleRepository `json:"repositories,omitempty"`
}

// RulesSearch represents the response from searching for rules.
// The Actives field is a map because rule keys are dynamic.
type RulesSearch struct {
	Actives map[string][]RuleActivation `json:"actives,omitempty"`
	Facets  []SearchFacet               `json:"facets,omitempty"`
	Rules   []RuleDetails               `json:"rules,omitempty"`
	Paging  Paging                      `json:"paging,omitzero"`
}

// RuleActivation represents how a rule is activated in a quality profile.
type RuleActivation struct {
	// Inherit indicates how the rule is inherited (NONE, INHERITED, OVERRIDES).
	Inherit string `json:"inherit,omitempty"`
	// QProfile is the key of the quality profile where the rule is activated.
	QProfile string `json:"qProfile,omitempty"`
	// Severity is the severity level of the activated rule.
	Severity string `json:"severity,omitempty"`
	// Params is the list of parameter values for the activated rule.
	Params []ParamKV `json:"params,omitempty"`
}

// ParamKV represents a key-value pair for rule parameters.
type ParamKV struct {
	// Key is the parameter name.
	Key string `json:"key,omitempty"`
	// Value is the parameter value.
	Value string `json:"value,omitempty"`
}

// SearchFacet represents a facet in search results.
type SearchFacet struct {
	// Name is the facet name (e.g., languages, repositories, tags).
	Name string `json:"name,omitempty"`
	// Values is the list of facet values with their counts.
	Values []FacetItem `json:"values,omitempty"`
}

// FacetItem represents a single facet value with its count.
type FacetItem struct {
	// Val is the facet value.
	Val string `json:"val,omitempty"`
	// Count is the number of items matching this facet value.
	Count int64 `json:"count,omitempty"`
}

// RuleDetails contains comprehensive information about a rule.
type RuleDetails struct {
	Name                       string               `json:"name,omitempty"`
	Key                        string               `json:"key,omitempty"`
	CreatedAt                  string               `json:"createdAt,omitempty"`
	UpdatedAt                  string               `json:"updatedAt,omitempty"`
	RemFnType                  string               `json:"remFnType,omitempty"`
	HTMLDesc                   string               `json:"htmlDesc,omitempty"`
	HTMLNote                   string               `json:"htmlNote,omitempty"`
	MdNote                     string               `json:"mdNote,omitempty"`
	NoteLogin                  string               `json:"noteLogin,omitempty"`
	CleanCodeAttribute         string               `json:"cleanCodeAttribute,omitempty"`
	InternalKey                string               `json:"internalKey,omitempty"`
	RemFnGapMultiplier         string               `json:"remFnGapMultiplier,omitempty"`
	RemFnBaseEffort            string               `json:"remFnBaseEffort,omitempty"`
	DefaultRemFnBaseEffort     string               `json:"defaultRemFnBaseEffort,omitempty"`
	Lang                       string               `json:"lang,omitempty"`
	LangName                   string               `json:"langName,omitempty"`
	CleanCodeAttributeCategory string               `json:"cleanCodeAttributeCategory,omitempty"`
	GapDescription             string               `json:"gapDescription,omitempty"`
	Repo                       string               `json:"repo,omitempty"`
	Scope                      string               `json:"scope,omitempty"`
	Severity                   string               `json:"severity,omitempty"`
	Status                     string               `json:"status,omitempty"`
	DefaultRemFnType           string               `json:"defaultRemFnType,omitempty"`
	DefaultRemFnGapMultiplier  string               `json:"defaultRemFnGapMultiplier,omitempty"`
	TemplateKey                string               `json:"templateKey,omitempty"`
	Type                       string               `json:"type,omitempty"`
	Impacts                    []RuleImpact         `json:"impacts,omitempty"`
	Tags                       []any                `json:"tags,omitempty"`
	SysTags                    []string             `json:"sysTags,omitempty"`
	Params                     []RuleParam          `json:"params,omitempty"`
	DescriptionSections        []DescriptionSection `json:"descriptionSections,omitempty"`
	IsTemplate                 bool                 `json:"isTemplate,omitempty"`
	IsExternal                 bool                 `json:"isExternal,omitempty"`
	RemFnOverloaded            bool                 `json:"remFnOverloaded,omitempty"`
	Template                   bool                 `json:"template,omitempty"`
}

// DescriptionSection represents a section of a rule's description.
type DescriptionSection struct {
	// Content is the HTML content of the section.
	Content string `json:"content,omitempty"`
	// Context provides additional context for the section.
	Context DescriptionContext `json:"context,omitzero"`
	// Key is the unique identifier of the section.
	Key string `json:"key,omitempty"`
}

// DescriptionContext provides context for a description section.
type DescriptionContext struct {
	// DisplayName is the human-readable name of the context.
	DisplayName string `json:"displayName,omitempty"`
	// Key is the unique identifier of the context.
	Key string `json:"key,omitempty"`
}

// RulesShow represents the response from showing a specific rule.
type RulesShow struct {
	Actives []RuleActivationDetailed `json:"actives,omitempty"`
	Rule    RuleDetails              `json:"rule,omitzero"`
}

// RuleActivationDetailed contains detailed information about a rule activation.
type RuleActivationDetailed struct {
	// Inherit indicates how the rule is inherited (NONE, INHERITED, OVERRIDES).
	Inherit string `json:"inherit,omitempty"`
	// QProfile is the key of the quality profile where the rule is activated.
	QProfile string `json:"qProfile,omitempty"`
	// Severity is the severity level of the activated rule.
	Severity string `json:"severity,omitempty"`
	// Params is the list of parameter values for the activated rule.
	Params []ParamKV `json:"params,omitempty"`
	// PrioritizedRule indicates if the rule is prioritized in this profile.
	PrioritizedRule bool `json:"prioritizedRule,omitempty"`
}

// RulesTags contains the list of available rule tags.
type RulesTags struct {
	Tags []string `json:"tags,omitempty"`
}

// RulesUpdate represents the response from updating a rule.
type RulesUpdate struct {
	Rule Rule `json:"rule,omitzero"`
}

// RulesCreateOption contains options for creating a custom rule.
type RulesCreateOption struct {
	// CleanCodeAttribute represents the Clean Code Attribute associated with the rule.
	// Allowed values: CONVENTIONAL, FORMATTED, IDENTIFIABLE, CLEAR, COMPLETE, EFFICIENT,
	// LOGICAL, DISTINCT, FOCUSED, MODULAR, TESTED, LAWFUL, RESPECTFUL, TRUSTWORTHY
	CleanCodeAttribute string `url:"cleanCodeAttribute,omitempty"`
	// CustomKey is the unique identifier for the custom rule (required).
	// Maximum length: 200 characters
	CustomKey string `url:"customKey,omitempty"`
	// Impacts is a map of software quality to severity (e.g., MAINTAINABILITY: HIGH, SECURITY: LOW).
	// Allowed keys: MAINTAINABILITY, RELIABILITY, SECURITY
	// Allowed values: INFO, LOW, MEDIUM, HIGH, BLOCKER
	Impacts map[string]string `url:"impacts,omitempty"`
	// MarkdownDescription is the Markdown-formatted description of the rule (required).
	MarkdownDescription string `url:"markdownDescription,omitempty"`
	// Name is the display name of the rule (required).
	// Maximum length: 200 characters
	Name string `url:"name,omitempty"`
	// Params is a map of parameter names to values (e.g., key1: v1, key2: v2).
	Params map[string]string `url:"params,omitempty"`
	// PreventReactivation prevents reactivation in profiles where the rule was deactivated.
	//
	// Deprecated: Since SonarQube 10.4
	PreventReactivation string `url:"preventReactivation,omitempty"`
	// Severity is the rule severity.
	// Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER
	Severity string `url:"severity,omitempty"`
	// Status is the status of the rule.
	// Allowed values: READY, DEPRECATED, BETA
	Status string `url:"status,omitempty"`
	// TemplateKey is the key of the template rule from which to create the custom rule (required).
	TemplateKey string `url:"templateKey,omitempty"`
	// Type is the rule type.
	// Allowed values: CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT
	Type string `url:"type,omitempty"`
}

// RulesDeleteOption contains options for deleting a custom rule.
type RulesDeleteOption struct {
	// Key is the unique identifier of the rule to be deleted (required).
	Key string `url:"key,omitempty"`
}

// RulesListOption contains options for listing rules.
// WARNING: Internal endpoint, may change without notice.
//
//nolint:govet // Field alignment is less important than logical grouping and readability
type RulesListOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// AvailableSince filters rules added since the specified date (format: yyyy-MM-dd).
	// If not set, all rules are returned.
	AvailableSince string `url:"available_since,omitempty"`
	// Qprofile is the key of the quality profile to filter by activated rules.
	Qprofile string `url:"qprofile,omitempty"`
	// Sort is the sort field.
	// Allowed values: createdAt
	Sort string `url:"s,omitempty"`
	// Asc indicates whether to sort results in ascending order.
	Asc bool `url:"asc,omitempty"`
}

// RulesRepositoriesOption contains options for listing rule repositories.
type RulesRepositoriesOption struct {
	// Language filters repositories by programming language.
	// If provided, only repositories for the given language will be returned.
	Language string `url:"language,omitempty"`
	// Query is a pattern to match repository keys/names against.
	Query string `url:"q,omitempty"`
}

// RulesSearchOption contains options for searching rules.
//
//nolint:govet // Field alignment is less important than logical grouping and readability
type RulesSearchOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// Activation filters rules by their activation status on the selected quality profile.
	// If Qprofile is not set, this parameter is ignored.
	Activation bool `url:"activation,omitempty"`
	// ActiveImpactSeverities filters by impact severity of rules in quality profiles.
	// Allowed values: INFO, LOW, MEDIUM, HIGH, BLOCKER
	ActiveImpactSeverities []string `url:"active_impactSeverities,omitempty,comma"`
	// ActiveSeverities filters by severity of rules in quality profiles.
	// Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER
	ActiveSeverities []string `url:"active_severities,omitempty,comma"`
	// Asc indicates whether to sort results in ascending order.
	Asc bool `url:"asc,omitempty"`
	// AvailableSince filters rules added since the specified date (format: yyyy-MM-dd).
	AvailableSince string `url:"available_since,omitempty"`
	// CleanCodeAttributeCategories filters by clean code attribute categories.
	// Allowed values: ADAPTABLE, CONSISTENT, INTENTIONAL, RESPONSIBLE
	CleanCodeAttributeCategories []string `url:"cleanCodeAttributeCategories,omitempty,comma"`
	// CompareToProfile is the key of the quality profile to compare against (internal parameter).
	CompareToProfile string `url:"compareToProfile,omitempty"`
	// ComplianceStandards filters by compliance standards.
	ComplianceStandards []string `url:"complianceStandards,omitempty,comma"`
	// Cwe filters by CWE identifiers. Use 'unknown' to select rules not associated with any CWE.
	Cwe []string `url:"cwe,omitempty,comma"`
	// Fields specifies which fields to return in the response.
	// Allowed values: actives, cleanCodeAttributeCategory, cleanCodeAttribute, createdAt, debtRemFn,
	// defaultDebtRemFn, defaultRemFn, deprecatedKeys, descriptionSections, educationPrinciples,
	// gapDescription, htmlDesc, htmlNote, internalKey, isExternal, isTemplate, lang, langName,
	// mdDesc, mdNote, name, noteLogin, params, remFn, remFnOverloaded, repo, scope, severity,
	// status, sysTags, tags, templateKey, updatedAt
	Fields []string `url:"f,omitempty,comma"`
	// Facets specifies which facets to compute and return.
	// Allowed values: languages, repositories, tags, severities, active_severities, statuses, types,
	// true, cwe, owaspTop10, owaspTop10-2021, owaspMobileTop10-2024, sansTop25, sonarsourceSecurity,
	// cleanCodeAttributeCategories, impactSeverities, impactSoftwareQualities, active_impactSeverities,
	// complianceStandards
	Facets []string `url:"facets,omitempty,comma"`
	// ImpactSeverities filters by impact severity of rules.
	// Allowed values: INFO, LOW, MEDIUM, HIGH, BLOCKER
	ImpactSeverities []string `url:"impactSeverities,omitempty,comma"`
	// ImpactSoftwareQualities filters by impact software quality of rules.
	// Allowed values: MAINTAINABILITY, RELIABILITY, SECURITY
	ImpactSoftwareQualities []string `url:"impactSoftwareQualities,omitempty,comma"`
	// IncludeExternal determines whether to include external rules in the results.
	IncludeExternal bool `url:"include_external,omitempty"`
	// Inheritance filters by inheritance status within a quality profile.
	// Used only if Activation parameter is set.
	// Allowed values: NONE, INHERITED, OVERRIDES
	Inheritance []string `url:"inheritance,omitempty,comma"`
	// IsTemplate filters rules based on whether they are templates.
	IsTemplate bool `url:"is_template,omitempty"`
	// Languages filters by programming languages.
	Languages []string `url:"languages,omitempty,comma"`
	// OwaspMobileTop102024 filters by OWASP Mobile Top 10 - 2024 categories.
	// Allowed values: m1, m2, m3, m4, m5, m6, m7, m8, m9, m10
	OwaspMobileTop102024 []string `url:"owaspMobileTop10-2024,omitempty,comma"`
	// OwaspTop10 filters by OWASP Top 10 2017 categories.
	// Allowed values: a1, a2, a3, a4, a5, a6, a7, a8, a9, a10
	OwaspTop10 []string `url:"owaspTop10,omitempty,comma"`
	// OwaspTop102021 filters by OWASP Top 10 2021 categories.
	// Allowed values: a1, a2, a3, a4, a5, a6, a7, a8, a9, a10
	OwaspTop102021 []string `url:"owaspTop10-2021,omitempty,comma"`
	// PrioritizedRule filters rules based on whether they are prioritized in the selected quality profile.
	// If Qprofile is not set, this parameter is ignored.
	PrioritizedRule bool `url:"prioritizedRule,omitempty"`
	// Query is a free text search query to filter rules (must be at least 2 characters).
	// Searches in rule name, description, note, tags, and key.
	Query string `url:"q,omitempty"`
	// Qprofile is the key of the quality profile to filter by activation status.
	// Only rules of the same language as this profile are returned.
	Qprofile string `url:"qprofile,omitempty"`
	// Repositories filters by rule repositories.
	Repositories []string `url:"repositories,omitempty,comma"`
	// RuleKey is the unique identifier of a specific rule to search for.
	RuleKey string `url:"rule_key,omitempty"`
	// Sort specifies the sort field.
	// Allowed values: key, name, createdAt, updatedAt
	Sort string `url:"s,omitempty"`
	// SansTop25 filters by SANS Top 25 categories.
	// Allowed values: insecure-interaction, risky-resource, porous-defenses
	//
	// Deprecated: Since SonarQube 10.0
	SansTop25 []string `url:"sansTop25,omitempty,comma"`
	// Severities filters by rule severities.
	// Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER
	Severities []string `url:"severities,omitempty,comma"`
	// SonarsourceSecurity filters by SonarSource security categories.
	// Use 'others' to select rules not associated with any category.
	// Allowed values: buffer-overflow, sql-injection, rce, object-injection, command-injection,
	// path-traversal-injection, ldap-injection, xpath-injection, log-injection, xxe, xss, dos,
	// ssrf, csrf, http-response-splitting, open-redirect, weak-cryptography, auth, insecure-conf,
	// file-manipulation, encrypt-data, traceability, permission, others
	SonarsourceSecurity []string `url:"sonarsourceSecurity,omitempty,comma"`
	// Statuses filters by rule statuses.
	// Allowed values: READY, DEPRECATED, REMOVED, BETA
	Statuses []string `url:"statuses,omitempty,comma"`
	// Tags filters by rule tags (OR filter - rules having at least one of the tags will be returned).
	Tags []string `url:"tags,omitempty,comma"`
	// TemplateKey filters custom rules based on the specified template key.
	TemplateKey string `url:"template_key,omitempty"`
	// Types filters by rule types (OR filter - rules matching at least one type will be returned).
	// Allowed values: CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT
	Types []string `url:"types,omitempty,comma"`
}

// RulesShowOption contains options for showing a specific rule.
type RulesShowOption struct {
	// Key is the unique identifier of the rule to be retrieved (required).
	Key string `url:"key,omitempty"`
	// Actives determines whether to include the list of quality profiles where the rule is active.
	Actives bool `url:"actives,omitempty"`
}

// RulesTagsOption contains options for listing rule tags.
type RulesTagsOption struct {
	// Query limits the search to tags containing the supplied string.
	Query string `url:"q,omitempty"`
	// PageSize is the response page size (must be greater than 0 and less than or equal to 500).
	PageSize int64 `url:"ps,omitempty"`
}

// RulesUpdateOption contains options for updating a rule.
type RulesUpdateOption struct {
	// Impacts is a map of software quality to severity (e.g., MAINTAINABILITY: HIGH, SECURITY: LOW).
	// Allowed keys: MAINTAINABILITY, RELIABILITY, SECURITY
	// Allowed values: INFO, LOW, MEDIUM, HIGH, BLOCKER
	Impacts map[string]string `url:"impacts,omitempty"`
	// Key is the unique identifier of the rule to be updated (required).
	// Maximum length: 200 characters
	Key string `url:"key,omitempty"`
	// MarkdownDescription is the Markdown-formatted description of the rule.
	// Mandatory for custom and manual rules.
	MarkdownDescription string `url:"markdownDescription,omitempty"`
	// MarkdownNote is the optional note in Markdown format.
	// Use empty value to remove current note. Note is not changed if parameter is not set.
	MarkdownNote string `url:"markdown_note,omitempty"`
	// Name is the name of the rule (mandatory for custom rules).
	// Maximum length: 200 characters
	Name string `url:"name,omitempty"`
	// Params is a map of parameter names to values.
	// Only applicable when updating a custom rule.
	Params map[string]string `url:"params,omitempty"`
	// RemediationFnBaseEffort is the base effort of the remediation function (e.g., '1d').
	RemediationFnBaseEffort string `url:"remediation_fn_base_effort,omitempty"`
	// RemediationFnType is the type of the remediation function.
	// Allowed values: LINEAR, CONSTANT, LINEAR_OFFSET
	RemediationFnType string `url:"remediation_fn_type,omitempty"`
	// RemediationFyGapMultiplier is the gap multiplier of the remediation function (e.g., '2min').
	RemediationFyGapMultiplier string `url:"remediation_fy_gap_multiplier,omitempty"`
	// Severity is the rule severity (only when updating a custom rule).
	// Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER
	Severity string `url:"severity,omitempty"`
	// Status is the status of the rule.
	// Allowed values: READY, DEPRECATED, REMOVED, BETA
	Status string `url:"status,omitempty"`
	// Tags is a list of tags to associate with the rule.
	// Use empty slice to remove current tags. Tags are not changed if parameter is not set.
	Tags []string `url:"tags,omitempty,comma"`
}

// App retrieves data required for rendering the 'Coding Rules' page.
// WARNING: This is an internal endpoint, may change without notice.
func (s *RulesService) App() (v *RulesApp, resp *http.Response, err error) {
	req, err := s.client.NewRequest(http.MethodGet, "rules/app", nil)
	if err != nil {
		return
	}

	v = new(RulesApp)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Create creates a custom rule.
// Requires the 'Administer Quality Profiles' permission.
func (s *RulesService) Create(opt *RulesCreateOption) (v *RulesCreate, resp *http.Response, err error) {
	err = s.ValidateCreateOpt(opt)
	if err != nil {
		return
	}

	// Convert to URL-encodable format
	urlOpt := s.convertCreateOptForURL(opt)

	req, err := s.client.NewRequest(http.MethodPost, "rules/create", urlOpt)
	if err != nil {
		return
	}

	v = new(RulesCreate)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Delete deletes a custom rule.
// Requires the 'Administer Quality Profiles' permission.
func (s *RulesService) Delete(opt *RulesDeleteOption) (resp *http.Response, err error) {
	err = s.ValidateDeleteOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "rules/delete", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// List lists rules, excluding external rules and rules with status REMOVED.
func (s *RulesService) List(opt *RulesListOption) (v *string, resp *http.Response, err error) {
	err = s.ValidateListOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "rules/list", opt)
	if err != nil {
		return
	}

	v = new(string)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Repositories lists available rule repositories.
func (s *RulesService) Repositories(opt *RulesRepositoriesOption) (v *RulesRepositories, resp *http.Response, err error) {
	err = s.ValidateRepositoriesOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "rules/repositories", opt)
	if err != nil {
		return
	}

	v = new(RulesRepositories)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Search searches for a collection of relevant rules matching a specified query.
func (s *RulesService) Search(opt *RulesSearchOption) (v *RulesSearch, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}

	// Convert to URL-encodable format
	urlOpt := s.convertSearchOptForURL(opt)

	req, err := s.client.NewRequest(http.MethodGet, "rules/search", urlOpt)
	if err != nil {
		return
	}

	v = new(RulesSearch)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Show retrieves detailed information about a specific rule.
func (s *RulesService) Show(opt *RulesShowOption) (v *RulesShow, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "rules/show", opt)
	if err != nil {
		return
	}

	v = new(RulesShow)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Tags lists all available rule tags.
func (s *RulesService) Tags(opt *RulesTagsOption) (v *RulesTags, resp *http.Response, err error) {
	err = s.ValidateTagsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "rules/tags", opt)
	if err != nil {
		return
	}

	v = new(RulesTags)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Update updates an existing rule.
// Requires the 'Administer Quality Profiles' permission.
func (s *RulesService) Update(opt *RulesUpdateOption) (v *RulesUpdate, resp *http.Response, err error) {
	err = s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	// Convert to URL-encodable format
	urlOpt := s.convertUpdateOptForURL(opt)

	req, err := s.client.NewRequest(http.MethodPost, "rules/update", urlOpt)
	if err != nil {
		return nil, nil, err
	}

	// Special handling: If Tags is explicitly set to an empty array, we need to append "tags="
	// to the query string because the URL encoder omits empty values with omitempty.
	// This allows clearing all tags from a rule.
	if opt.Tags != nil && len(opt.Tags) == 0 {
		if req.URL.RawQuery == "" {
			req.URL.RawQuery = "tags="
		} else {
			req.URL.RawQuery += "&tags="
		}
	}

	v = new(RulesUpdate)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// ValidateCreateOpt validates the options for creating a custom rule.
//
//nolint:cyclop,funlen // Validation functions are naturally complex due to multiple checks
func (s *RulesService) ValidateCreateOpt(opt *RulesCreateOption) error {
	if opt == nil {
		return NewValidationError("RulesCreateOption", "cannot be nil", ErrMissingRequired)
	}

	// Required fields
	err := ValidateRequired(opt.CustomKey, "CustomKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.MarkdownDescription, "MarkdownDescription")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.TemplateKey, "TemplateKey")
	if err != nil {
		return err
	}

	// Length validations
	err = ValidateMaxLength(opt.CustomKey, MaxRuleKeyLength, "CustomKey")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxRuleKeyLength, "Name")
	if err != nil {
		return err
	}

	// Value validations
	if opt.CleanCodeAttribute != "" {
		err := IsValueAuthorized(opt.CleanCodeAttribute, allowedCleanCodeAttributes, "CleanCodeAttribute")
		if err != nil {
			return err
		}
	}

	if opt.Severity != "" {
		err := IsValueAuthorized(opt.Severity, allowedSeverities, "Severity")
		if err != nil {
			return err
		}
	}

	if opt.Status != "" {
		err := IsValueAuthorized(opt.Status, allowedRulesExistingStatuses, "Status")
		if err != nil {
			return err
		}
	}

	if opt.Type != "" {
		err := IsValueAuthorized(opt.Type, allowedRulesTypes, "Type")
		if err != nil {
			return err
		}
	}

	// Validate Impacts map
	if len(opt.Impacts) > 0 {
		err := ValidateMapKeys(opt.Impacts, allowedImpactSoftwareQualities, "Impacts")
		if err != nil {
			return err
		}

		err = ValidateMapValues(opt.Impacts, allowedImpactSeverities, "Impacts")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteOpt validates the options for deleting a custom rule.
func (s *RulesService) ValidateDeleteOpt(opt *RulesDeleteOption) error {
	if opt == nil {
		return NewValidationError("RulesDeleteOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for listing rules.
func (s *RulesService) ValidateListOpt(opt *RulesListOption) error {
	if opt == nil {
		return nil
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	if opt.Sort != "" {
		allowed := map[string]struct{}{"createdAt": {}}

		err := IsValueAuthorized(opt.Sort, allowed, "Sort")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateRepositoriesOpt validates the options for listing rule repositories.
func (s *RulesService) ValidateRepositoriesOpt(opt *RulesRepositoriesOption) error {
	// No specific validations needed for this endpoint
	return nil
}

// ValidateSearchOpt validates the options for searching rules.
//
//nolint:cyclop,funlen // Validation functions are naturally complex due to multiple checks
func (s *RulesService) ValidateSearchOpt(opt *RulesSearchOption) error {
	if opt == nil {
		return nil
	}

	// Validate pagination
	err := opt.Validate()
	if err != nil {
		return err
	}

	// Validate Query minimum length
	err = ValidateMinLength(opt.Query, MinSearchQueryLength, "Query")
	if err != nil {
		return err
	}

	// Validate severity values
	err = AreValuesAuthorized(opt.ActiveImpactSeverities, allowedImpactSeverities, "ActiveImpactSeverities")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.ActiveSeverities, allowedSeverities, "ActiveSeverities")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.Severities, allowedSeverities, "Severities")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.ImpactSeverities, allowedImpactSeverities, "ImpactSeverities")
	if err != nil {
		return err
	}

	// Validate clean code attribute categories
	err = AreValuesAuthorized(opt.CleanCodeAttributeCategories, allowedCleanCodeAttributesCategories, "CleanCodeAttributeCategories")
	if err != nil {
		return err
	}

	// Validate impact software qualities
	err = AreValuesAuthorized(opt.ImpactSoftwareQualities, allowedImpactSoftwareQualities, "ImpactSoftwareQualities")
	if err != nil {
		return err
	}

	// Validate inheritance
	err = AreValuesAuthorized(opt.Inheritance, allowedInheritanceTypes, "Inheritance")
	if err != nil {
		return err
	}

	// Validate OWASP categories
	err = AreValuesAuthorized(opt.OwaspTop10, allowedOwaspCategories, "OwaspTop10")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.OwaspTop102021, allowedOwaspCategories, "OwaspTop102021")
	if err != nil {
		return err
	}

	// Validate OWASP Mobile Top 10
	err = AreValuesAuthorized(opt.OwaspMobileTop102024, allowedOwaspMobileCategories, "OwaspMobileTop102024")
	if err != nil {
		return err
	}

	// Validate SANS Top 25
	err = AreValuesAuthorized(opt.SansTop25, allowedSansTop25Categories, "SansTop25")
	if err != nil {
		return err
	}

	// Validate statuses
	err = AreValuesAuthorized(opt.Statuses, allowedRulesStatuses, "Statuses")
	if err != nil {
		return err
	}

	// Validate types
	err = AreValuesAuthorized(opt.Types, allowedRulesTypes, "Types")
	if err != nil {
		return err
	}

	// Validate sort field
	if opt.Sort != "" {
		allowed := map[string]struct{}{"key": {}, "name": {}, "createdAt": {}, "updatedAt": {}}

		err = IsValueAuthorized(opt.Sort, allowed, "Sort")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateShowOpt validates the options for showing a specific rule.
func (s *RulesService) ValidateShowOpt(opt *RulesShowOption) error {
	if opt == nil {
		return NewValidationError("RulesShowOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return nil
}

// ValidateTagsOpt validates the options for listing rule tags.
func (s *RulesService) ValidateTagsOpt(opt *RulesTagsOption) error {
	if opt == nil {
		return nil
	}

	// Validate page size
	if opt.PageSize != 0 && (opt.PageSize < MinPageSize || opt.PageSize > MaxPageSize) {
		return NewValidationError("PageSize", fmt.Sprintf("must be between %d and %d", MinPageSize, MaxPageSize), ErrOutOfRange)
	}

	return nil
}

// ValidateUpdateOpt validates the options for updating a rule.
//
//nolint:cyclop // Validation functions are naturally complex due to multiple checks
func (s *RulesService) ValidateUpdateOpt(opt *RulesUpdateOption) error {
	if opt == nil {
		return NewValidationError("RulesUpdateOption", "cannot be nil", ErrMissingRequired)
	}

	// Required field
	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	// Length validations
	err = ValidateMaxLength(opt.Key, MaxRuleKeyLength, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxRuleKeyLength, "Name")
	if err != nil {
		return err
	}

	// Value validations
	err = IsValueAuthorized(opt.Severity, allowedSeverities, "Severity")
	if err != nil {
		return err
	}

	if opt.Status != "" {
		err = IsValueAuthorized(opt.Status, allowedRulesStatuses, "Status")
		if err != nil {
			return err
		}
	}

	if opt.RemediationFnType != "" {
		allowed := map[string]struct{}{"LINEAR": {}, "CONSTANT": {}, "LINEAR_OFFSET": {}}

		err = IsValueAuthorized(opt.RemediationFnType, allowed, "RemediationFnType")
		if err != nil {
			return err
		}
	}

	// Validate Impacts map
	if len(opt.Impacts) > 0 {
		err = ValidateMapKeys(opt.Impacts, allowedImpactSoftwareQualities, "Impacts")
		if err != nil {
			return err
		}

		err = ValidateMapValues(opt.Impacts, allowedImpactSeverities, "Impacts")
		if err != nil {
			return err
		}
	}

	return nil
}

// convertCreateOptForURL converts RulesCreateOption to a URL-encodable format.
func (s *RulesService) convertCreateOptForURL(opt *RulesCreateOption) *rulesCreateURLOption {
	//nolint:exhaustruct // Only populate fields that have values
	urlOpt := &rulesCreateURLOption{
		CleanCodeAttribute:  opt.CleanCodeAttribute,
		CustomKey:           opt.CustomKey,
		MarkdownDescription: opt.MarkdownDescription,
		Name:                opt.Name,
		PreventReactivation: opt.PreventReactivation,
		Severity:            opt.Severity,
		Status:              opt.Status,
		TemplateKey:         opt.TemplateKey,
		Type:                opt.Type,
	}

	// Convert maps to semicolon-separated strings
	if len(opt.Impacts) > 0 {
		urlOpt.Impacts = MapToSeparatedString(opt.Impacts, ";", "=")
	}

	if len(opt.Params) > 0 {
		urlOpt.Params = MapToSeparatedString(opt.Params, ";", "=")
	}

	return urlOpt
}

// rulesCreateURLOption is the URL-encodable version of RulesCreateOption.
type rulesCreateURLOption struct {
	CleanCodeAttribute  string `url:"cleanCodeAttribute,omitempty"`
	CustomKey           string `url:"customKey,omitempty"`
	Impacts             string `url:"impacts,omitempty"`
	MarkdownDescription string `url:"markdownDescription,omitempty"`
	Name                string `url:"name,omitempty"`
	Params              string `url:"params,omitempty"`
	PreventReactivation string `url:"preventReactivation,omitempty"`
	Severity            string `url:"severity,omitempty"`
	Status              string `url:"status,omitempty"`
	TemplateKey         string `url:"templateKey,omitempty"`
	Type                string `url:"type,omitempty"`
}

// convertUpdateOptForURL converts RulesUpdateOption to a URL-encodable format.
func (s *RulesService) convertUpdateOptForURL(opt *RulesUpdateOption) *rulesUpdateURLOption {
	//nolint:exhaustruct // Only populate fields that have values
	urlOpt := &rulesUpdateURLOption{
		Key:                        opt.Key,
		MarkdownDescription:        opt.MarkdownDescription,
		MarkdownNote:               opt.MarkdownNote,
		Name:                       opt.Name,
		RemediationFnBaseEffort:    opt.RemediationFnBaseEffort,
		RemediationFnType:          opt.RemediationFnType,
		RemediationFyGapMultiplier: opt.RemediationFyGapMultiplier,
		Severity:                   opt.Severity,
		Status:                     opt.Status,
	}

	// Convert maps to semicolon-separated strings
	if len(opt.Impacts) > 0 {
		urlOpt.Impacts = MapToSeparatedString(opt.Impacts, ";", "=")
	}

	if len(opt.Params) > 0 {
		urlOpt.Params = MapToSeparatedString(opt.Params, ";", "=")
	}

	// Convert tags slice to comma-separated string
	if len(opt.Tags) > 0 {
		urlOpt.Tags = ListToSeparatedString(opt.Tags, ",")
	}

	return urlOpt
}

// rulesUpdateURLOption is the URL-encodable version of RulesUpdateOption.
type rulesUpdateURLOption struct {
	Impacts                    string `url:"impacts,omitempty"`
	Key                        string `url:"key,omitempty"`
	MarkdownDescription        string `url:"markdownDescription,omitempty"`
	MarkdownNote               string `url:"markdown_note,omitempty"`
	Name                       string `url:"name,omitempty"`
	Params                     string `url:"params,omitempty"`
	RemediationFnBaseEffort    string `url:"remediation_fn_base_effort,omitempty"`
	RemediationFnType          string `url:"remediation_fn_type,omitempty"`
	RemediationFyGapMultiplier string `url:"remediation_fy_gap_multiplier,omitempty"`
	Severity                   string `url:"severity,omitempty"`
	Status                     string `url:"status,omitempty"`
	Tags                       string `url:"tags,omitempty"`
}

// convertSearchOptForURL converts RulesSearchOption to a URL-encodable format.
//
//nolint:cyclop,funlen // Conversion functions need to handle many optional fields
func (s *RulesService) convertSearchOptForURL(opt *RulesSearchOption) *rulesSearchURLOption {
	if opt == nil {
		return nil
	}

	//nolint:exhaustruct // Only populate fields that have values
	urlOpt := &rulesSearchURLOption{
		Page:             opt.Page,
		PageSize:         opt.PageSize,
		Activation:       opt.Activation,
		Asc:              opt.Asc,
		AvailableSince:   opt.AvailableSince,
		CompareToProfile: opt.CompareToProfile,
		IncludeExternal:  opt.IncludeExternal,
		IsTemplate:       opt.IsTemplate,
		PrioritizedRule:  opt.PrioritizedRule,
		Query:            opt.Query,
		Qprofile:         opt.Qprofile,
		RuleKey:          opt.RuleKey,
		Sort:             opt.Sort,
		TemplateKey:      opt.TemplateKey,
	}

	// Convert all slices to comma-separated strings
	if len(opt.ActiveImpactSeverities) > 0 {
		urlOpt.ActiveImpactSeverities = ListToSeparatedString(opt.ActiveImpactSeverities, ",")
	}

	if len(opt.ActiveSeverities) > 0 {
		urlOpt.ActiveSeverities = ListToSeparatedString(opt.ActiveSeverities, ",")
	}

	if len(opt.CleanCodeAttributeCategories) > 0 {
		urlOpt.CleanCodeAttributeCategories = ListToSeparatedString(opt.CleanCodeAttributeCategories, ",")
	}

	if len(opt.ComplianceStandards) > 0 {
		urlOpt.ComplianceStandards = ListToSeparatedString(opt.ComplianceStandards, ",")
	}

	if len(opt.Cwe) > 0 {
		urlOpt.Cwe = ListToSeparatedString(opt.Cwe, ",")
	}

	if len(opt.Fields) > 0 {
		urlOpt.Fields = ListToSeparatedString(opt.Fields, ",")
	}

	if len(opt.Facets) > 0 {
		urlOpt.Facets = ListToSeparatedString(opt.Facets, ",")
	}

	if len(opt.ImpactSeverities) > 0 {
		urlOpt.ImpactSeverities = ListToSeparatedString(opt.ImpactSeverities, ",")
	}

	if len(opt.ImpactSoftwareQualities) > 0 {
		urlOpt.ImpactSoftwareQualities = ListToSeparatedString(opt.ImpactSoftwareQualities, ",")
	}

	if len(opt.Inheritance) > 0 {
		urlOpt.Inheritance = ListToSeparatedString(opt.Inheritance, ",")
	}

	if len(opt.Languages) > 0 {
		urlOpt.Languages = ListToSeparatedString(opt.Languages, ",")
	}

	if len(opt.OwaspMobileTop102024) > 0 {
		urlOpt.OwaspMobileTop102024 = ListToSeparatedString(opt.OwaspMobileTop102024, ",")
	}

	if len(opt.OwaspTop10) > 0 {
		urlOpt.OwaspTop10 = ListToSeparatedString(opt.OwaspTop10, ",")
	}

	if len(opt.OwaspTop102021) > 0 {
		urlOpt.OwaspTop102021 = ListToSeparatedString(opt.OwaspTop102021, ",")
	}

	if len(opt.Repositories) > 0 {
		urlOpt.Repositories = ListToSeparatedString(opt.Repositories, ",")
	}

	if len(opt.SansTop25) > 0 {
		urlOpt.SansTop25 = ListToSeparatedString(opt.SansTop25, ",")
	}

	if len(opt.Severities) > 0 {
		urlOpt.Severities = ListToSeparatedString(opt.Severities, ",")
	}

	if len(opt.SonarsourceSecurity) > 0 {
		urlOpt.SonarsourceSecurity = ListToSeparatedString(opt.SonarsourceSecurity, ",")
	}

	if len(opt.Statuses) > 0 {
		urlOpt.Statuses = ListToSeparatedString(opt.Statuses, ",")
	}

	if len(opt.Tags) > 0 {
		urlOpt.Tags = ListToSeparatedString(opt.Tags, ",")
	}

	if len(opt.Types) > 0 {
		urlOpt.Types = ListToSeparatedString(opt.Types, ",")
	}

	return urlOpt
}

// rulesSearchURLOption is the URL-encodable version of RulesSearchOption.
//
//nolint:govet // Field alignment less important than maintaining consistent field order
type rulesSearchURLOption struct {
	Page                         int64  `url:"p,omitempty"`
	PageSize                     int64  `url:"ps,omitempty"`
	Activation                   bool   `url:"activation,omitempty"`
	ActiveImpactSeverities       string `url:"active_impactSeverities,omitempty"`
	ActiveSeverities             string `url:"active_severities,omitempty"`
	Asc                          bool   `url:"asc,omitempty"`
	AvailableSince               string `url:"available_since,omitempty"`
	CleanCodeAttributeCategories string `url:"cleanCodeAttributeCategories,omitempty"`
	CompareToProfile             string `url:"compareToProfile,omitempty"`
	ComplianceStandards          string `url:"complianceStandards,omitempty"`
	Cwe                          string `url:"cwe,omitempty"`
	Fields                       string `url:"f,omitempty"`
	Facets                       string `url:"facets,omitempty"`
	ImpactSeverities             string `url:"impactSeverities,omitempty"`
	ImpactSoftwareQualities      string `url:"impactSoftwareQualities,omitempty"`
	IncludeExternal              bool   `url:"include_external,omitempty"`
	Inheritance                  string `url:"inheritance,omitempty"`
	IsTemplate                   bool   `url:"is_template,omitempty"`
	Languages                    string `url:"languages,omitempty"`
	OwaspMobileTop102024         string `url:"owaspMobileTop10-2024,omitempty"`
	OwaspTop10                   string `url:"owaspTop10,omitempty"`
	OwaspTop102021               string `url:"owaspTop10-2021,omitempty"`
	PrioritizedRule              bool   `url:"prioritizedRule,omitempty"`
	Query                        string `url:"q,omitempty"`
	Qprofile                     string `url:"qprofile,omitempty"`
	Repositories                 string `url:"repositories,omitempty"`
	RuleKey                      string `url:"rule_key,omitempty"`
	Sort                         string `url:"s,omitempty"`
	SansTop25                    string `url:"sansTop25,omitempty"`
	Severities                   string `url:"severities,omitempty"`
	SonarsourceSecurity          string `url:"sonarsourceSecurity,omitempty"`
	Statuses                     string `url:"statuses,omitempty"`
	Tags                         string `url:"tags,omitempty"`
	TemplateKey                  string `url:"template_key,omitempty"`
	Types                        string `url:"types,omitempty"`
}
