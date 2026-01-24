package sonargo

import "net/http"

// RulesService handles communication with the Rules related methods of the SonarQube API.
type RulesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// RulesAppResponse contains metadata for rendering the 'Coding Rules' page.
type RulesAppResponse struct {
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

// RulesCreateResponse represents the response from creating a custom rule.
type RulesCreateResponse struct {
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

// RulesRepositoriesResponse contains the list of available rule repositories.
type RulesRepositoriesResponse struct {
	Repositories []RuleRepository `json:"repositories,omitempty"`
}

// RulesSearchResponse represents the response from searching for rules.
// The Actives field is a map because rule keys are dynamic.
type RulesSearchResponse struct {
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

// RulesShowResponse represents the response from showing a specific rule.
type RulesShowResponse struct {
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

// RulesTagsResponse contains the list of available rule tags.
type RulesTagsResponse struct {
	Tags []string `json:"tags,omitempty"`
}

// RulesUpdateResponse represents the response from updating a rule.
type RulesUpdateResponse struct {
	Rule Rule `json:"rule,omitzero"`
}

// RulesCreateOption contains options for creating a custom rule.
type RulesCreateOption struct {
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

// RulesDeleteOption contains options for deleting a custom rule.
type RulesDeleteOption struct {
	Key string `url:"key,omitempty"`
}

// RulesListOption contains options for listing rules.
type RulesListOption struct {
	Asc            string `url:"asc,omitempty"`
	AvailableSince string `url:"available_since,omitempty"`
	P              string `url:"p,omitempty"`
	Ps             string `url:"ps,omitempty"`
	Qprofile       string `url:"qprofile,omitempty"`
	S              string `url:"s,omitempty"`
}

// RulesRepositoriesOption contains options for listing rule repositories.
type RulesRepositoriesOption struct {
	Language string `url:"language,omitempty"`
	Q        string `url:"q,omitempty"`
}

// RulesSearchOption contains options for searching rules.
type RulesSearchOption struct {
	Activation                   string `url:"activation,omitempty"`
	ActiveImpactSeverities       string `url:"active_impactSeverities,omitempty"`
	ActiveSeverities             string `url:"active_severities,omitempty"`
	Asc                          string `url:"asc,omitempty"`
	AvailableSince               string `url:"available_since,omitempty"`
	CleanCodeAttributeCategories string `url:"cleanCodeAttributeCategories,omitempty"`
	CompareToProfile             string `url:"compareToProfile,omitempty"`
	ComplianceStandards          string `url:"complianceStandards,omitempty"`
	Cwe                          string `url:"cwe,omitempty"`
	F                            string `url:"f,omitempty"`
	Facets                       string `url:"facets,omitempty"`
	ImpactSeverities             string `url:"impactSeverities,omitempty"`
	ImpactSoftwareQualities      string `url:"impactSoftwareQualities,omitempty"`
	IncludeExternal              string `url:"include_external,omitempty"`
	Inheritance                  string `url:"inheritance,omitempty"`
	IsTemplate                   string `url:"is_template,omitempty"`
	Languages                    string `url:"languages,omitempty"`
	OwaspMobileTop102024         string `url:"owaspMobileTop10-2024,omitempty"`
	OwaspTop10                   string `url:"owaspTop10,omitempty"`
	OwaspTop102021               string `url:"owaspTop10-2021,omitempty"`
	P                            string `url:"p,omitempty"`
	PrioritizedRule              string `url:"prioritizedRule,omitempty"`
	Ps                           string `url:"ps,omitempty"`
	Q                            string `url:"q,omitempty"`
	Qprofile                     string `url:"qprofile,omitempty"`
	Repositories                 string `url:"repositories,omitempty"`
	RuleKey                      string `url:"rule_key,omitempty"`
	S                            string `url:"s,omitempty"`
	SansTop25                    string `url:"sansTop25,omitempty"`
	Severities                   string `url:"severities,omitempty"`
	SonarsourceSecurity          string `url:"sonarsourceSecurity,omitempty"`
	Statuses                     string `url:"statuses,omitempty"`
	Tags                         string `url:"tags,omitempty"`
	TemplateKey                  string `url:"template_key,omitempty"`
	Types                        string `url:"types,omitempty"`
}

// RulesShowOption contains options for showing a specific rule.
type RulesShowOption struct {
	Actives string `url:"actives,omitempty"`
	Key     string `url:"key,omitempty"`
}

// RulesTagsOption contains options for listing rule tags.
type RulesTagsOption struct {
	Ps string `url:"ps,omitempty"`
	Q  string `url:"q,omitempty"`
}

// RulesUpdateOption contains options for updating a rule.
type RulesUpdateOption struct {
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

// App retrieves data required for rendering the 'Coding Rules' page.
func (s *RulesService) App() (v *RulesAppResponse, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "rules/app", nil)
	if err != nil {
		return
	}

	v = new(RulesAppResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Create creates a custom rule.
// Requires the 'Administer Quality Profiles' permission.
func (s *RulesService) Create(opt *RulesCreateOption) (v *RulesCreateResponse, resp *http.Response, err error) {
	err = s.ValidateCreateOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "rules/create", opt)
	if err != nil {
		return
	}

	v = new(RulesCreateResponse)

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

	req, err := s.client.NewRequest("POST", "rules/delete", opt)
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

	req, err := s.client.NewRequest("GET", "rules/list", opt)
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
func (s *RulesService) Repositories(opt *RulesRepositoriesOption) (v *RulesRepositoriesResponse, resp *http.Response, err error) {
	err = s.ValidateRepositoriesOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "rules/repositories", opt)
	if err != nil {
		return
	}

	v = new(RulesRepositoriesResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Search searches for a collection of relevant rules matching a specified query.
func (s *RulesService) Search(opt *RulesSearchOption) (v *RulesSearchResponse, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "rules/search", opt)
	if err != nil {
		return
	}

	v = new(RulesSearchResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Show retrieves detailed information about a specific rule.
func (s *RulesService) Show(opt *RulesShowOption) (v *RulesShowResponse, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "rules/show", opt)
	if err != nil {
		return
	}

	v = new(RulesShowResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Tags lists all available rule tags.
func (s *RulesService) Tags(opt *RulesTagsOption) (v *RulesTagsResponse, resp *http.Response, err error) {
	err = s.ValidateTagsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "rules/tags", opt)
	if err != nil {
		return
	}

	v = new(RulesTagsResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Update updates an existing rule.
// Requires the 'Administer Quality Profiles' permission.
func (s *RulesService) Update(opt *RulesUpdateOption) (v *RulesUpdateResponse, resp *http.Response, err error) {
	err = s.ValidateUpdateOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "rules/update", opt)
	if err != nil {
		return
	}

	v = new(RulesUpdateResponse)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ValidateCreateOpt validates the options for creating a custom rule.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateCreateOpt(_ *RulesCreateOption) error {
	return nil
}

// ValidateDeleteOpt validates the options for deleting a custom rule.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateDeleteOpt(_ *RulesDeleteOption) error {
	return nil
}

// ValidateListOpt validates the options for listing rules.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateListOpt(_ *RulesListOption) error {
	return nil
}

// ValidateRepositoriesOpt validates the options for listing rule repositories.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateRepositoriesOpt(_ *RulesRepositoriesOption) error {
	return nil
}

// ValidateSearchOpt validates the options for searching rules.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateSearchOpt(_ *RulesSearchOption) error {
	return nil
}

// ValidateShowOpt validates the options for showing a specific rule.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateShowOpt(_ *RulesShowOption) error {
	return nil
}

// ValidateTagsOpt validates the options for listing rule tags.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateTagsOpt(_ *RulesTagsOption) error {
	return nil
}

// ValidateUpdateOpt validates the options for updating a rule.
// Currently a no-op placeholder as there are no specific validations implemented.
func (s *RulesService) ValidateUpdateOpt(_ *RulesUpdateOption) error {
	return nil
}
