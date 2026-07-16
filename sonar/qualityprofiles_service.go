package sonar

import (
	"context"
	"net/http"
)

const (
	// MaxQualityProfileNameLength is the maximum allowed length for quality profile names.
	MaxQualityProfileNameLength = 100
)

// QualityprofilesService handles communication with the Quality Profiles related methods of the SonarQube API.
// Quality profiles define sets of rules that are applied to projects during analysis.
type QualityprofilesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// QualityprofilesChangelog represents the response from getting a quality profile's changelog.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesChangelog struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Events is the list of changelog events.
	Events []ChangelogEvent `json:"events,omitempty"`
	// Page is the current page number (legacy duplicate of Paging.PageIndex).
	Page int64 `json:"p,omitempty"`
	// PageSize is the page size (legacy duplicate of Paging.PageSize).
	PageSize int64 `json:"ps,omitempty"`
	// Total is the total number of events (legacy duplicate of Paging.Total).
	Total int64 `json:"total,omitempty"`
}

// ChangelogEvent represents a single event in the quality profile changelog.
type ChangelogEvent struct {
	// Action is the type of action performed (e.g., ACTIVATED, DEACTIVATED, UPDATED).
	Action string `json:"action,omitempty"`
	// AuthorLogin is the login of the user who made the change.
	AuthorLogin string `json:"authorLogin,omitempty"`
	// AuthorName is the name of the user who made the change.
	AuthorName string `json:"authorName,omitempty"`
	// CleanCodeAttributeCategory is the clean code attribute category for the rule.
	CleanCodeAttributeCategory string `json:"cleanCodeAttributeCategory,omitempty"`
	// Date is the timestamp of the event.
	Date string `json:"date,omitempty"`
	// RuleKey is the key of the rule that was affected.
	RuleKey string `json:"ruleKey,omitempty"`
	// RuleName is the name of the rule that was affected.
	RuleName string `json:"ruleName,omitempty"`
	// Params contains the change details.
	Params ChangelogEventParams `json:"params,omitzero"`
	// Impacts contains the impact changes.
	Impacts []ChangelogImpact `json:"impacts,omitempty"`
}

// ChangelogEventParams represents the parameters of a changelog event.
type ChangelogEventParams struct {
	// NewCleanCodeAttribute is the new clean code attribute after the change.
	NewCleanCodeAttribute string `json:"newCleanCodeAttribute,omitempty"`
	// NewCleanCodeAttributeCategory is the new clean code attribute category after the change.
	NewCleanCodeAttributeCategory string `json:"newCleanCodeAttributeCategory,omitempty"`
	// OldCleanCodeAttribute is the previous clean code attribute before the change.
	OldCleanCodeAttribute string `json:"oldCleanCodeAttribute,omitempty"`
	// OldCleanCodeAttributeCategory is the previous clean code attribute category before the change.
	OldCleanCodeAttributeCategory string `json:"oldCleanCodeAttributeCategory,omitempty"`
	// PrioritizedRule indicates if the rule was marked as prioritized.
	PrioritizedRule string `json:"prioritizedRule,omitempty"`
	// SonarQubeVersion is the version of SonarQube when the change occurred.
	SonarQubeVersion string `json:"sonarQubeVersion,omitempty"`
	// ImpactChanges contains the impact severity changes.
	ImpactChanges []ImpactChange `json:"impactChanges,omitempty"`
}

// ImpactChange represents a change in impact severity.
type ImpactChange struct {
	// SoftwareQuality is the software quality being impacted.
	SoftwareQuality string `json:"softwareQuality,omitempty"`
	// OldSeverity is the previous severity level.
	OldSeverity string `json:"oldSeverity,omitempty"`
	// NewSeverity is the new severity level.
	NewSeverity string `json:"newSeverity,omitempty"`
}

// ChangelogImpact represents an impact entry in the changelog.
type ChangelogImpact struct {
	// SoftwareQuality is the software quality being impacted.
	SoftwareQuality string `json:"softwareQuality,omitempty"`
	// Severity is the impact severity.
	Severity string `json:"severity,omitempty"`
}

// QualityprofilesCompare represents the response from comparing two quality profiles.
type QualityprofilesCompare struct {
	// Left contains information about the left profile.
	Left QualityprofilesCompareProfile `json:"left,omitzero"`
	// Right contains information about the right profile.
	Right QualityprofilesCompareProfile `json:"right,omitzero"`
	// InLeft contains rules only in the left profile.
	InLeft []QualityprofilesCompareRule `json:"inLeft,omitempty"`
	// InRight contains rules only in the right profile.
	InRight []QualityprofilesCompareRule `json:"inRight,omitempty"`
	// Modified contains rules that differ between profiles.
	Modified []QualityprofilesCompareModifiedRule `json:"modified,omitempty"`
	// Same contains rules that are the same in both profiles.
	Same []QualityprofilesCompareRule `json:"same,omitempty"`
}

// QualityprofilesCompareProfile represents a profile in a comparison response.
type QualityprofilesCompareProfile struct {
	// Key is the unique key of the profile.
	Key string `json:"key,omitempty"`
	// Name is the name of the profile.
	Name string `json:"name,omitempty"`
}

// QualityprofilesCompareRule represents a rule in a comparison response.
type QualityprofilesCompareRule struct {
	// Key is the rule key.
	Key string `json:"key,omitempty"`
	// Name is the rule name.
	Name string `json:"name,omitempty"`
	// PluginKey is the plugin key.
	PluginKey string `json:"pluginKey,omitempty"`
	// PluginName is the plugin name.
	PluginName string `json:"pluginName,omitempty"`
	// LanguageKey is the language key.
	LanguageKey string `json:"languageKey,omitempty"`
	// LanguageName is the language name.
	LanguageName string `json:"languageName,omitempty"`
}

// QualityprofilesCompareModifiedRule represents a modified rule in a comparison response.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesCompareModifiedRule struct {
	// Key is the rule key.
	Key string `json:"key,omitempty"`
	// Name is the rule name.
	Name string `json:"name,omitempty"`
	// PluginKey is the plugin key.
	PluginKey string `json:"pluginKey,omitempty"`
	// PluginName is the plugin name.
	PluginName string `json:"pluginName,omitempty"`
	// LanguageKey is the language key.
	LanguageKey string `json:"languageKey,omitempty"`
	// LanguageName is the language name.
	LanguageName string `json:"languageName,omitempty"`
	// Left contains the left profile's settings for this rule.
	Left QualityprofilesRuleSetting `json:"left,omitzero"`
	// Right contains the right profile's settings for this rule.
	Right QualityprofilesRuleSetting `json:"right,omitzero"`
}

// QualityprofilesRuleSetting represents the settings for a rule in a profile.
type QualityprofilesRuleSetting struct {
	// Params contains the rule parameters.
	Params map[string]string `json:"params,omitempty"`
}

// QualityprofilesCopy represents the response from copying a quality profile.
type QualityprofilesCopy struct {
	// Key is the unique key of the new profile.
	Key string `json:"key,omitempty"`
	// Name is the name of the new profile.
	Name string `json:"name,omitempty"`
	// Language is the language of the new profile.
	Language string `json:"language,omitempty"`
	// LanguageName is the display name of the language.
	LanguageName string `json:"languageName,omitempty"`
	// ParentKey is the key of the parent profile.
	ParentKey string `json:"parentKey,omitempty"`
	// IsDefault indicates if this is the default profile.
	IsDefault bool `json:"isDefault,omitempty"`
	// IsInherited indicates if this profile inherits from another.
	IsInherited bool `json:"isInherited,omitempty"`
}

// QualityprofilesCreate represents the response from creating a quality profile.
type QualityprofilesCreate struct {
	// Profile contains the created profile details.
	Profile QualityprofilesCreatedProfile `json:"profile,omitzero"`
}

// QualityprofilesCreatedProfile represents a newly created quality profile.
type QualityprofilesCreatedProfile struct {
	// Key is the unique key of the profile.
	Key string `json:"key,omitempty"`
	// Name is the name of the profile.
	Name string `json:"name,omitempty"`
	// Language is the language of the profile.
	Language string `json:"language,omitempty"`
	// LanguageName is the display name of the language.
	LanguageName string `json:"languageName,omitempty"`
	// IsDefault indicates if this is the default profile.
	IsDefault bool `json:"isDefault,omitempty"`
	// IsInherited indicates if this profile inherits from another.
	IsInherited bool `json:"isInherited,omitempty"`
}

// QualityprofilesExporters represents the response from listing exporters.
//
// Deprecated: No more custom profile exporters since SonarQube 25.4.
type QualityprofilesExporters struct {
	// Exporters is the list of available exporters.
	Exporters []QualityprofilesExporter `json:"exporters,omitempty"`
}

// QualityprofilesExporter represents a quality profile exporter.
type QualityprofilesExporter struct {
	// Key is the unique key of the exporter.
	Key string `json:"key,omitempty"`
	// Name is the name of the exporter.
	Name string `json:"name,omitempty"`
	// Languages is the list of supported languages.
	Languages []string `json:"languages,omitempty"`
}

// QualityprofilesImporters represents the response from listing importers.
//
// Deprecated: Since SonarQube 25.4.
type QualityprofilesImporters struct {
	// Importers is the list of available importers.
	Importers []QualityprofilesImporter `json:"importers,omitempty"`
}

// QualityprofilesImporter represents a quality profile importer.
type QualityprofilesImporter struct {
	// Key is the unique key of the importer.
	Key string `json:"key,omitempty"`
	// Name is the name of the importer.
	Name string `json:"name,omitempty"`
	// Languages is the list of supported languages.
	Languages []string `json:"languages,omitempty"`
}

// QualityprofilesInheritance represents the response from getting inheritance info.
type QualityprofilesInheritance struct {
	// Profile contains the current profile information.
	Profile QualityprofilesInheritanceProfile `json:"profile,omitzero"`
	// Ancestors contains the ancestor profiles.
	Ancestors []QualityprofilesInheritanceProfile `json:"ancestors,omitempty"`
	// Children contains the child profiles.
	Children []QualityprofilesInheritanceProfile `json:"children,omitempty"`
}

// QualityprofilesInheritanceProfile represents a profile in an inheritance hierarchy.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type QualityprofilesInheritanceProfile struct {
	// Key is the unique key of the profile.
	Key string `json:"key,omitempty"`
	// Name is the name of the profile.
	Name string `json:"name,omitempty"`
	// ActiveRuleCount is the number of active rules.
	ActiveRuleCount int64 `json:"activeRuleCount,omitempty"`
	// InactiveRuleCount is the number of inactive rules.
	InactiveRuleCount int64 `json:"inactiveRuleCount,omitempty"`
	// OverridingRuleCount is the number of overriding rules.
	OverridingRuleCount int64 `json:"overridingRuleCount,omitempty"`
	// IsBuiltIn indicates if this is a built-in profile.
	IsBuiltIn bool `json:"isBuiltIn,omitempty"`
	// Parent is the key of the parent profile, present on the current profile entry.
	Parent string `json:"parent,omitempty"`
}

// QualityprofilesProjects represents the response from listing associated projects.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesProjects struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Results contains the list of projects.
	Results []QualityprofilesProfileProject `json:"results,omitempty"`
}

// QualityprofilesProfileProject represents a project associated with a quality profile.
type QualityprofilesProfileProject struct {
	// Key is the unique key of the project.
	Key string `json:"key,omitempty"`
	// Name is the name of the project.
	Name string `json:"name,omitempty"`
	// Selected indicates if the project is explicitly bound to the profile.
	Selected bool `json:"selected,omitempty"`
}

// QualityprofilesSearch represents the response from searching quality profiles.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesSearch struct {
	// Actions contains global actions available.
	Actions QualityprofilesActions `json:"actions,omitzero"`
	// Profiles is the list of quality profiles.
	Profiles []QualityProfile `json:"profiles,omitempty"`
}

// QualityprofilesActions represents global actions for quality profiles.
type QualityprofilesActions struct {
	// Create indicates if the current user can create quality profiles.
	Create bool `json:"create,omitempty"`
}

// QualityProfile represents a quality profile with its properties.
//
//nolint:govet // Field alignment less important for logical grouping
type QualityProfile struct {
	// Actions contains the actions available for this profile.
	Actions QualityProfileActions `json:"actions,omitzero"`
	// Key is the unique key of the profile.
	Key string `json:"key,omitempty"`
	// Name is the name of the profile.
	Name string `json:"name,omitempty"`
	// Language is the language code of the profile.
	Language string `json:"language,omitempty"`
	// LanguageName is the display name of the language.
	LanguageName string `json:"languageName,omitempty"`
	// ParentKey is the key of the parent profile (for inherited profiles).
	ParentKey string `json:"parentKey,omitempty"`
	// ParentName is the name of the parent profile.
	ParentName string `json:"parentName,omitempty"`
	// LastUsed is the timestamp when the profile was last used.
	LastUsed string `json:"lastUsed,omitempty"`
	// RuleUpdatedAt is the timestamp when rules were last updated.
	RuleUpdatedAt string `json:"rulesUpdatedAt,omitempty"`
	// UserUpdatedAt is the timestamp when the profile was last updated by a user.
	UserUpdatedAt string `json:"userUpdatedAt,omitempty"`
	// ActiveDeprecatedRuleCount is the count of active deprecated rules.
	ActiveDeprecatedRuleCount int64 `json:"activeDeprecatedRuleCount,omitempty"`
	// ActiveRuleCount is the count of active rules.
	ActiveRuleCount int64 `json:"activeRuleCount,omitempty"`
	// ProjectCount is the number of projects using this profile.
	ProjectCount int64 `json:"projectCount,omitempty"`
	// IsBuiltIn indicates if this is a built-in profile.
	IsBuiltIn bool `json:"isBuiltIn,omitempty"`
	// IsDefault indicates if this is the default profile for the language.
	IsDefault bool `json:"isDefault,omitempty"`
	// IsInherited indicates if this profile inherits from another.
	IsInherited bool `json:"isInherited,omitempty"`
}

// QualityProfileActions represents actions available for a specific quality profile.
type QualityProfileActions struct {
	// AssociateProjects indicates if projects can be associated with this profile.
	AssociateProjects bool `json:"associateProjects,omitempty"`
	// Copy indicates if the profile can be copied.
	Copy bool `json:"copy,omitempty"`
	// Delete indicates if the profile can be deleted.
	Delete bool `json:"delete,omitempty"`
	// Edit indicates if the profile can be edited.
	Edit bool `json:"edit,omitempty"`
	// SetAsDefault indicates if the profile can be set as default.
	SetAsDefault bool `json:"setAsDefault,omitempty"`
}

// QualityprofilesSearchGroups represents the response from searching groups.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesSearchGroups struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Groups is the list of groups.
	Groups []QualityprofilesProfileGroup `json:"groups,omitempty"`
}

// QualityprofilesProfileGroup represents a group that can edit a quality profile.
type QualityprofilesProfileGroup struct {
	// Name is the name of the group.
	Name string `json:"name,omitempty"`
	// Description is the description of the group.
	Description string `json:"description,omitempty"`
	// Selected indicates if the group is allowed to edit the profile.
	Selected bool `json:"selected,omitempty"`
}

// QualityprofilesSearchUsers represents the response from searching users.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesSearchUsers struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Users is the list of users.
	Users []QualityprofilesProfileUser `json:"users,omitempty"`
}

// QualityprofilesProfileUser represents a user that can edit a quality profile.
type QualityprofilesProfileUser struct {
	// Login is the login name of the user.
	Login string `json:"login,omitempty"`
	// Name is the display name of the user.
	Name string `json:"name,omitempty"`
	// Avatar is the avatar URL of the user.
	Avatar string `json:"avatar,omitempty"`
	// Selected indicates if the user is allowed to edit the profile.
	Selected bool `json:"selected,omitempty"`
}

// QualityprofilesShow represents the response from showing a quality profile.
type QualityprofilesShow struct {
	// Profile contains the profile details.
	Profile QualityprofilesShownProfile `json:"profile,omitzero"`
}

// QualityprofilesShownProfile represents a quality profile in a show response.
type QualityprofilesShownProfile struct {
	// Key is the unique key of the profile.
	Key string `json:"key,omitempty"`
	// Name is the name of the profile.
	Name string `json:"name,omitempty"`
	// Language is the language code of the profile.
	Language string `json:"language,omitempty"`
	// LanguageName is the display name of the language.
	LanguageName string `json:"languageName,omitempty"`
	// LastUsed is the timestamp when the profile was last used.
	LastUsed string `json:"lastUsed,omitempty"`
	// RulesUpdatedAt is the timestamp when rules were last updated.
	RulesUpdatedAt string `json:"rulesUpdatedAt,omitempty"`
	// ActiveDeprecatedRuleCount is the count of active deprecated rules.
	ActiveDeprecatedRuleCount int64 `json:"activeDeprecatedRuleCount,omitempty"`
	// ActiveRuleCount is the count of active rules.
	ActiveRuleCount int64 `json:"activeRuleCount,omitempty"`
	// ProjectCount is the number of projects using this profile.
	ProjectCount int64 `json:"projectCount,omitempty"`
	// IsBuiltIn indicates if this is a built-in profile.
	IsBuiltIn bool `json:"isBuiltIn,omitempty"`
	// IsDefault indicates if this is the default profile for the language.
	IsDefault bool `json:"isDefault,omitempty"`
	// IsInherited indicates if this profile inherits from another.
	IsInherited bool `json:"isInherited,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// QualityprofilesActivateRuleOptions contains options for activating a rule.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesActivateRuleOptions struct {
	// Key is the quality profile key (required).
	Key string `url:"key,omitempty"`
	// Rule is the rule key (required).
	Rule string `url:"rule,omitempty"`
	// Impacts is the override of impact severities (e.g., "MAINTAINABILITY=HIGH;SECURITY=MEDIUM").
	// Cannot be used at the same time as Severity.
	Impacts map[string]string `url:"impacts,omitempty"`
	// Params is the parameters as semi-colon separated key=value pairs.
	// Ignored if Reset is true.
	Params map[string]string `url:"params,omitempty"`
	// PrioritizedRule marks the activated rule as prioritized.
	PrioritizedRule bool `url:"prioritizedRule,omitempty"`
	// Reset resets severity and parameters to parent profile or rule defaults.
	Reset bool `url:"reset,omitempty"`
	// Severity is the severity level.
	// Cannot be used at the same time as Impacts.
	// Possible values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER
	Severity string `url:"severity,omitempty"`
}

// QualityprofilesActivateRulesOptions contains options for bulk activating rules.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesActivateRulesOptions struct {
	// TargetKey is the quality profile key on which rules are activated (required).
	TargetKey string `url:"targetKey,omitempty"`
	// Activation filters rules that are activated or deactivated on the selected quality profile.
	Activation bool `url:"activation,omitempty"`
	// ActiveImpactSeverities filters by activation software quality severities.
	// Allowed values: 'INFO', 'LOW', 'MEDIUM', 'HIGH', 'BLOCKER'
	ActiveImpactSeverities []string `url:"active_impactSeverities,omitempty,comma"`
	// ActiveSeverities filters by activation severities.
	// Allowed values: 'INFO', 'MINOR', 'MAJOR', 'CRITICAL', 'BLOCKER'
	ActiveSeverities []string `url:"active_severities,omitempty,comma"`
	// Asc indicates ascending sort.
	Asc bool `url:"asc,omitempty"`
	// AvailableSince filters rules added since date (format: yyyy-MM-dd).
	AvailableSince string `url:"available_since,omitempty"`
	// CleanCodeAttributeCategories filters by clean code attribute categories.
	// Allowed values: 'ADAPTABLE', 'CONSISTENT', 'INTENTIONAL', 'RESPONSIBLE'
	CleanCodeAttributeCategories []string `url:"cleanCodeAttributeCategories,omitempty,comma"`
	// CompareToProfile is a quality profile key to filter rules that are activated.
	// WARNING: This parameter is internal and may change without notice.
	CompareToProfile string `url:"compareToProfile,omitempty"`
	// Cwe filters by CWE identifiers.
	Cwe []string `url:"cwe,omitempty,comma"`
	// ImpactSeverities filters by software quality severities.
	// Allowed values: 'INFO', 'LOW', 'MEDIUM', 'HIGH', 'BLOCKER'
	ImpactSeverities []string `url:"impactSeverities,omitempty,comma"`
	// ImpactSoftwareQualities filters by software qualities.
	// Allowed values: 'MAINTAINABILITY', 'RELIABILITY', 'SECURITY',
	ImpactSoftwareQualities []string `url:"impactSoftwareQualities,omitempty,comma"`
	// Inheritance filters by inheritance for a rule within a quality profile.
	// Allowed values: 'NONE', 'INHERITED', 'OVERIDDES'
	Inheritance []string `url:"inheritance,omitempty,comma"`
	// IsTemplate filters template rules.
	IsTemplate bool `url:"is_template,omitempty"`
	// Languages filters by languages.
	Languages []string `url:"languages,omitempty,comma"`
	// OwaspMobileTop102024 filters by OWASP Mobile Top 10 2024 categories.
	// Allowed values: 'm1', 'm2', 'm3', 'm4', 'm5', 'm6', 'm7', 'm8', 'm9', 'm10'
	OwaspMobileTop102024 []string `url:"owaspMobileTop10-2024,omitempty,comma"`
	// OwaspTop10 filters by OWASP Top 10 2017 categories.
	// Allowed values: 'a1', 'a2', 'a3', 'a4', 'a5', 'a6', 'a7', 'a8', 'a9', 'a10'
	OwaspTop10 []string `url:"owaspTop10,omitempty,comma"`
	// OwaspTop102021 filters by OWASP Top 10 2021 categories.
	// Allowed values: 'a1', 'a2', 'a3', 'a4', 'a5', 'a6', 'a7', 'a8', 'a9', 'a10'
	OwaspTop102021 []string `url:"owaspTop10-2021,omitempty,comma"`
	// PrioritizedRule marks activated rules as prioritized.
	PrioritizedRule bool `url:"prioritizedRule,omitempty"`
	// Query is the UTF-8 search query.
	Query string `url:"q,omitempty"`
	// Qprofile is the quality profile key to filter on.
	Qprofile string `url:"qprofile,omitempty"`
	// Repositories filters by repositories.
	Repositories []string `url:"repositories,omitempty,comma"`
	// RuleKey filters by rule key.
	RuleKey string `url:"rule_key,omitempty"`
	// Sort is the sort field.
	// Allowed values: 'updatedAt', 'key', 'name', 'createdAt'
	Sort string `url:"s,omitempty"`
	// SansTop25 filters by SANS Top 25 categories.
	//
	// Deprecated: Since SonarQube 10.0.
	// Allowed values: 'insecure-interaction', 'risky-resource', 'porous-defenses'
	SansTop25 []string `url:"sansTop25,omitempty,comma"`
	// Severities filters by default severities.
	// Allowed values: 'INFO', 'MINOR', 'MAJOR', 'CRITICAL', 'BLOCKER'
	Severities []string `url:"severities,omitempty,comma"`
	// SonarsourceSecurity filters by SonarSource security categories.
	// Allowed values: 'buffer-overflow', 'sql-injection', 'rce', 'object-injection',
	// 'command-injection', 'path-traversal-injection', 'ldap-injection', 'xpath-injection',
	// 'log-injection', 'xxe', 'xss', 'dos', 'ssrf', 'csrf', 'http-response-splitting',
	// 'open-redirect', 'weak-cryptography', 'auth', 'insecure-conf', 'file-manipulation',
	// 'encrypt-data', 'traceability', 'permission', 'others'
	SonarsourceSecurity []string `url:"sonarsourceSecurity,omitempty,comma"`
	// Statuses filters by status codes.
	// Allowed values: 'READY', 'DEPRECATED', 'REMOVED', 'BETA'
	Statuses []string `url:"statuses,omitempty,comma"`
	// Tags filters by tags.
	Tags []string `url:"tags,omitempty,comma"`
	// TargetSeverity is the severity to set on the activated rules.
	// Allowed values: 'INFO', 'MINOR', 'MAJOR', 'CRITICAL', 'BLOCKER'
	TargetSeverity string `url:"targetSeverity,omitempty"`
	// TemplateKey is the template rule key to filter on.
	TemplateKey string `url:"template_key,omitempty"`
	// Types filters by types.
	// Allowed values: 'CODE_SMELL', 'BUG', 'VULNERABILITY', 'SECURITY_HOTSPOT'
	Types []string `url:"types,omitempty,comma"`
}

// QualityprofilesAddGroupOptions contains options for allowing a group to edit a profile.
type QualityprofilesAddGroupOptions struct {
	// Group is the group name (required).
	Group string `url:"group,omitempty"`
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesAddProjectOptions contains options for associating a project.
type QualityprofilesAddProjectOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesAddUserOptions contains options for allowing a user to edit a profile.
type QualityprofilesAddUserOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// Login is the user login (required).
	Login string `url:"login,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesBackupOptions contains options for backing up a profile.
type QualityprofilesBackupOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesChangeParentOptions contains options for changing a profile's parent.
type QualityprofilesChangeParentOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
	// ParentQualityProfile is the new parent profile name.
	// If not provided, breaks the inheritance link with the current parent.
	ParentQualityProfile string `url:"parentQualityProfile,omitempty"`
}

// QualityprofilesChangelogOptions contains options for getting a profile's changelog.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesChangelogOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
	// FilterMode filters events by mode
	// Allowed values: 'MQR', 'STANDARD'
	FilterMode string `url:"filterMode,omitempty"`
	// Since is the start date for the changelog (inclusive).
	Since string `url:"since,omitempty"`
	// To is the end date for the changelog (exclusive).
	To string `url:"to,omitempty"`
}

// QualityprofilesCompareOptions contains options for comparing two profiles.
// WARNING: This endpoint is internal and may change without notice.
type QualityprofilesCompareOptions struct {
	// LeftKey is the left profile key (required).
	LeftKey string `url:"leftKey,omitempty"`
	// RightKey is the right profile key (required).
	RightKey string `url:"rightKey,omitempty"`
}

// QualityprofilesCopyOptions contains options for copying a profile.
type QualityprofilesCopyOptions struct {
	// FromKey is the source quality profile key (required).
	FromKey string `url:"fromKey,omitempty"`
	// ToName is the name for the new quality profile (required).
	// Maximum length: 100 characters
	ToName string `url:"toName,omitempty"`
}

// QualityprofilesCreateOptions contains options for creating a profile.
type QualityprofilesCreateOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// Name is the quality profile name (required).
	// Maximum length: 100 characters
	Name string `url:"name,omitempty"`
}

// QualityprofilesDeactivateRuleOptions contains options for deactivating a rule.
type QualityprofilesDeactivateRuleOptions struct {
	// Key is the quality profile key (required).
	Key string `url:"key,omitempty"`
	// Rule is the rule key (required).
	Rule string `url:"rule,omitempty"`
}

// QualityprofilesDeactivateRulesOptions contains options for bulk deactivating rules.
// Uses the same filter parameters as QualityprofilesActivateRulesOption.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesDeactivateRulesOptions struct {
	// TargetKey is the quality profile key on which rules are deactivated (required).
	TargetKey string `url:"targetKey,omitempty"`
	// Activation filters rules that are activated or deactivated on the selected quality profile.
	Activation bool `url:"activation,omitempty"`
	// ActiveImpactSeverities filters by activation software quality severities.
	ActiveImpactSeverities string `url:"active_impactSeverities,omitempty"`
	// ActiveSeverities filters by activation severities.
	ActiveSeverities string `url:"active_severities,omitempty"`
	// Asc indicates ascending sort.
	Asc bool `url:"asc,omitempty"`
	// AvailableSince filters rules added since date (format: yyyy-MM-dd).
	AvailableSince string `url:"available_since,omitempty"`
	// CleanCodeAttributeCategories filters by clean code attribute categories.
	// Allowed values: 'ADAPTABLE', 'CONSISTENT', 'INTENTIONAL', 'RESPONSIBLE'
	CleanCodeAttributeCategories []string `url:"cleanCodeAttributeCategories,omitempty,comma"`
	// CompareToProfile is a quality profile key to filter rules that are activated.
	// WARNING: This parameter is internal and may change without notice.
	CompareToProfile string `url:"compareToProfile,omitempty"`
	// Cwe filters by CWE identifiers.
	Cwe []string `url:"cwe,omitempty,comma"`
	// ImpactSeverities filters by software quality severities.
	// Allowed values: 'INFO', 'LOW', 'MEDIUM', 'HIGH', 'BLOCKER'
	ImpactSeverities []string `url:"impactSeverities,omitempty,comma"`
	// ImpactSoftwareQualities filters by software qualities.
	// Allowed values: 'MAINTAINABILITY', 'RELIABILITY', 'SECURITY',
	ImpactSoftwareQualities []string `url:"impactSoftwareQualities,omitempty,comma"`
	// Inheritance filters by inheritance for a rule within a quality profile.
	// Allowed values: 'NONE', 'INHERITED', 'OVERIDDES'
	Inheritance []string `url:"inheritance,omitempty,comma"`
	// IsTemplate filters template rules.
	IsTemplate bool `url:"is_template,omitempty"`
	// Languages filters by languages.
	Languages []string `url:"languages,omitempty,comma"`
	// OwaspMobileTop102024 filters by OWASP Mobile Top 10 2024 categories.
	// Allowed values: 'm1', 'm2', 'm3', 'm4', 'm5', 'm6', 'm7', 'm8', 'm9', 'm10'
	OwaspMobileTop102024 []string `url:"owaspMobileTop10-2024,omitempty,comma"`
	// OwaspTop10 filters by OWASP Top 10 2017 categories.
	// Allowed values: 'a1', 'a2', 'a3', 'a4', 'a5', 'a6', 'a7', 'a8', 'a9', 'a10'
	OwaspTop10 []string `url:"owaspTop10,omitempty,comma"`
	// OwaspTop102021 filters by OWASP Top 10 2021 categories.
	// Allowed values: 'a1', 'a2', 'a3', 'a4', 'a5', 'a6', 'a7', 'a8', 'a9', 'a10'
	OwaspTop102021 []string `url:"owaspTop10-2021,omitempty,comma"`
	// Query is the UTF-8 search query.
	Query string `url:"q,omitempty"`
	// Qprofile is the quality profile key to filter on.
	Qprofile string `url:"qprofile,omitempty"`
	// Repositories filters by repositories.
	Repositories []string `url:"repositories,omitempty,comma"`
	// RuleKey filters by rule key.
	RuleKey string `url:"rule_key,omitempty"`
	// Sort is the sort field.
	// Allowed values: 'updatedAt', 'key', 'name', 'createdAt'
	Sort string `url:"s,omitempty"`
	// SansTop25 filters by SANS Top 25 categories.
	//
	// Deprecated: Since SonarQube 10.0.
	// Allowed values: 'insecure-interaction', 'risky-resource', 'porous-defenses'
	SansTop25 []string `url:"sansTop25,omitempty,comma"`
	// Severities filters by default severities.
	// Allowed values: 'INFO', 'MINOR', 'MAJOR', 'CRITICAL', 'BLOCKER'
	Severities []string `url:"severities,omitempty,comma"`
	// SonarsourceSecurity filters by SonarSource security categories.
	// Allowed values: 'buffer-overflow', 'sql-injection', 'rce', 'object-injection',
	// 'command-injection', 'path-traversal-injection', 'ldap-injection', 'xpath-injection',
	// 'log-injection', 'xxe', 'xss', 'dos', 'ssrf', 'csrf', 'http-response-splitting',
	// 'open-redirect', 'weak-cryptography', 'auth', 'insecure-conf', 'file-manipulation',
	// 'encrypt-data', 'traceability', 'permission', 'others'
	SonarsourceSecurity []string `url:"sonarsourceSecurity,omitempty,comma"`
	// Statuses filters by status codes.
	// Allowed values: 'READY', 'DEPRECATED', 'REMOVED', 'BETA'
	Statuses []string `url:"statuses,omitempty,comma"`
	// Tags filters by tags.
	Tags []string `url:"tags,omitempty,comma"`
	// TemplateKey is the template rule key to filter on.
	TemplateKey string `url:"template_key,omitempty"`
	// Types filters by types.
	// Allowed values: 'CODE_SMELL', 'BUG', 'VULNERABILITY', 'SECURITY_HOTSPOT'
	Types []string `url:"types,omitempty,comma"`
}

// QualityprofilesDeleteOptions contains options for deleting a profile.
type QualityprofilesDeleteOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesExportOptions contains options for exporting a profile.
//
// Deprecated: Since SonarQube 25.4. Use Backup instead.
type QualityprofilesExportOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name.
	// If empty, the default profile for the language is exported.
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesInheritanceOptions contains options for getting inheritance info.
type QualityprofilesInheritanceOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesProjectsOptions contains options for listing associated projects.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesProjectsOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// Key is the quality profile key (required).
	Key string `url:"key,omitempty"`
	// Query limits search to projects containing this string.
	Query string `url:"q,omitempty"`
	// Selected filters by selection status.
	// Allowed values: all, deselected, selected
	Selected string `url:"selected,omitempty"`
}

// QualityprofilesRemoveGroupOptions contains options for removing group permissions.
type QualityprofilesRemoveGroupOptions struct {
	// Group is the group name (required).
	Group string `url:"group,omitempty"`
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesRemoveProjectOptions contains options for removing a project association.
type QualityprofilesRemoveProjectOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesRemoveUserOptions contains options for removing user permissions.
type QualityprofilesRemoveUserOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// Login is the user login (required).
	Login string `url:"login,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesRenameOptions contains options for renaming a profile.
type QualityprofilesRenameOptions struct {
	// Key is the quality profile key (required).
	Key string `url:"key,omitempty"`
	// Name is the new quality profile name (required).
	// Maximum length: 100 characters
	Name string `url:"name,omitempty"`
}

// QualityprofilesRestoreOptions contains options for restoring a profile from backup.
type QualityprofilesRestoreOptions struct {
	// Backup is the profile backup file content in XML format (required).
	Backup string `url:"backup,omitempty"`
}

// QualityprofilesSearchOptions contains options for searching profiles.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesSearchOptions struct {
	// Defaults returns only default profiles if true.
	Defaults bool `url:"defaults,omitempty"`
	// Language is the quality profile language.
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// Project is the project key.
	Project string `url:"project,omitempty"`
	// QualityProfile is the quality profile name.
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesSearchGroupsOptions contains options for searching groups.
// WARNING: This endpoint is internal and may change without notice.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesSearchGroupsOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
	// Query limits search to group names containing this string.
	Query string `url:"q,omitempty"`
	// Selected filters by selection status.
	// Allowed values: all, deselected, selected
	Selected string `url:"selected,omitempty"`
}

// QualityprofilesSearchUsersOptions contains options for searching users.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualityprofilesSearchUsersOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// Language is the quality profile language (required).
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
	// Query limits search to names or logins containing this string.
	Query string `url:"q,omitempty"`
	// Selected filters by selection status.
	// Allowed values: all, deselected, selected
	Selected string `url:"selected,omitempty"`
}

// QualityprofilesSetDefaultOptions contains options for setting the default profile.
type QualityprofilesSetDefaultOptions struct {
	// Language is the quality profile language (required).
	// Allowed values: 'kubernetes', 'css', 'scala', 'jsp', 'py', 'js', 'docker', 'rust',
	// 'java', 'web', 'flex', 'xml', 'json', 'ipynb', 'text', 'vbnet', 'cloudformation',
	// 'yaml', 'go', 'kotlin', 'secrets', 'ruby', 'cs', 'php', 'terraform', 'azureresourcemanager', 'ts'
	Language string `url:"language,omitempty"`
	// QualityProfile is the quality profile name (required).
	QualityProfile string `url:"qualityProfile,omitempty"`
}

// QualityprofilesShowOptions contains options for showing a profile.
type QualityprofilesShowOptions struct {
	// Key is the quality profile key (required).
	Key string `url:"key,omitempty"`
	// CompareToSonarWay adds the number of missing rules from related Sonar way profile.
	CompareToSonarWay bool `url:"compareToSonarWay,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ActivateRule activates a rule on a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) ActivateRule(ctx context.Context, opt *QualityprofilesActivateRuleOptions) (resp *http.Response, err error) {
	err = s.ValidateActivateRuleOpt(opt)
	if err != nil {
		return
	}

	// Convert map fields to URL-encodable format
	urlOpt := s.convertActivateRuleOptForURL(opt)

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/activate_rule", urlOpt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// ActivateRules bulk-activates rules on one quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) ActivateRules(ctx context.Context, opt *QualityprofilesActivateRulesOptions) (resp *http.Response, err error) {
	err = s.ValidateActivateRulesOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/activate_rules", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// AddGroup allows a group to edit a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) AddGroup(ctx context.Context, opt *QualityprofilesAddGroupOptions) (resp *http.Response, err error) {
	err = s.ValidateAddGroupOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/add_group", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// AddProject associates a project with a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Administer right on the specified project
func (s *QualityprofilesService) AddProject(ctx context.Context, opt *QualityprofilesAddProjectOptions) (resp *http.Response, err error) {
	err = s.ValidateAddProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/add_project", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// AddUser allows a user to edit a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) AddUser(ctx context.Context, opt *QualityprofilesAddUserOptions) (resp *http.Response, err error) {
	err = s.ValidateAddUserOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/add_user", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Backup backs up a quality profile in XML form.
// The exported profile can be restored through Restore.
func (s *QualityprofilesService) Backup(ctx context.Context, opt *QualityprofilesBackupOptions) (v *string, resp *http.Response, err error) {
	err = s.ValidateBackupOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/backup", opt)
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

// ChangeParent changes a quality profile's parent.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) ChangeParent(ctx context.Context, opt *QualityprofilesChangeParentOptions) (resp *http.Response, err error) {
	err = s.ValidateChangeParentOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/change_parent", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Changelog gets the history of changes on a quality profile.
// Events are ordered by date in descending order (most recent first).
func (s *QualityprofilesService) Changelog(ctx context.Context, opt *QualityprofilesChangelogOptions) (v *QualityprofilesChangelog, resp *http.Response, err error) {
	err = s.ValidateChangelogOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/changelog", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesChangelog)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Compare compares two quality profiles.
func (s *QualityprofilesService) Compare(ctx context.Context, opt *QualityprofilesCompareOptions) (v *QualityprofilesCompare, resp *http.Response, err error) {
	err = s.ValidateCompareOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/compare", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesCompare)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Copy copies a quality profile.
// Requires the 'Administer Quality Profiles' permission.
func (s *QualityprofilesService) Copy(ctx context.Context, opt *QualityprofilesCopyOptions) (v *QualityprofilesCopy, resp *http.Response, err error) {
	err = s.ValidateCopyOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/copy", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesCopy)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Create creates a quality profile.
// Requires the 'Administer Quality Profiles' permission.
func (s *QualityprofilesService) Create(ctx context.Context, opt *QualityprofilesCreateOptions) (v *QualityprofilesCreate, resp *http.Response, err error) {
	err = s.ValidateCreateOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/create", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesCreate)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// DeactivateRule deactivates a rule on a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) DeactivateRule(ctx context.Context, opt *QualityprofilesDeactivateRuleOptions) (resp *http.Response, err error) {
	err = s.ValidateDeactivateRuleOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/deactivate_rule", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// DeactivateRules bulk-deactivates rules on quality profiles.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) DeactivateRules(ctx context.Context, opt *QualityprofilesDeactivateRulesOptions) (resp *http.Response, err error) {
	err = s.ValidateDeactivateRulesOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/deactivate_rules", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Delete deletes a quality profile and all its descendants.
// The default quality profile cannot be deleted.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) Delete(ctx context.Context, opt *QualityprofilesDeleteOptions) (resp *http.Response, err error) {
	err = s.ValidateDeleteOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/delete", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Export exports a quality profile.
//
// Deprecated: Since SonarQube 25.4. Use Backup instead.
func (s *QualityprofilesService) Export(ctx context.Context, opt *QualityprofilesExportOptions) (v *string, resp *http.Response, err error) {
	err = s.ValidateExportOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/export", opt)
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

// Exporters lists quality profile exporters.
//
// Deprecated: No more custom profile exporters since SonarQube 25.4.
func (s *QualityprofilesService) Exporters(ctx context.Context) (v *QualityprofilesExporters, resp *http.Response, err error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/exporters", nil)
	if err != nil {
		return
	}

	v = new(QualityprofilesExporters)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Importers lists supported quality profile importers.
//
// Deprecated: Since SonarQube 25.4.
func (s *QualityprofilesService) Importers(ctx context.Context) (v *QualityprofilesImporters, resp *http.Response, err error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/importers", nil)
	if err != nil {
		return
	}

	v = new(QualityprofilesImporters)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Inheritance shows a quality profile's ancestors and children.
func (s *QualityprofilesService) Inheritance(ctx context.Context, opt *QualityprofilesInheritanceOptions) (v *QualityprofilesInheritance, resp *http.Response, err error) {
	err = s.ValidateInheritanceOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/inheritance", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesInheritance)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Projects lists projects with their association status regarding a quality profile.
// Only projects explicitly bound to the profile are returned.
func (s *QualityprofilesService) Projects(ctx context.Context, opt *QualityprofilesProjectsOptions) (v *QualityprofilesProjects, resp *http.Response, err error) {
	err = s.ValidateProjectsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/projects", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesProjects)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// RemoveGroup removes the ability from a group to edit a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) RemoveGroup(ctx context.Context, opt *QualityprofilesRemoveGroupOptions) (resp *http.Response, err error) {
	err = s.ValidateRemoveGroupOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/remove_group", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// RemoveProject removes a project's association with a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
//   - Administer right on the specified project
func (s *QualityprofilesService) RemoveProject(ctx context.Context, opt *QualityprofilesRemoveProjectOptions) (resp *http.Response, err error) {
	err = s.ValidateRemoveProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/remove_project", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// RemoveUser removes the ability from a user to edit a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) RemoveUser(ctx context.Context, opt *QualityprofilesRemoveUserOptions) (resp *http.Response, err error) {
	err = s.ValidateRemoveUserOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/remove_user", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Rename renames a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) Rename(ctx context.Context, opt *QualityprofilesRenameOptions) (resp *http.Response, err error) {
	err = s.ValidateRenameOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/rename", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Restore restores a quality profile using an XML file.
// The restored profile name is taken from the backup file.
// If a profile with the same name and language exists, it will be overwritten.
// Requires the 'Administer Quality Profiles' permission.
func (s *QualityprofilesService) Restore(ctx context.Context, opt *QualityprofilesRestoreOptions) (resp *http.Response, err error) {
	err = s.ValidateRestoreOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/restore", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Search searches for quality profiles.
func (s *QualityprofilesService) Search(ctx context.Context, opt *QualityprofilesSearchOptions) (v *QualityprofilesSearch, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/search", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesSearch)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchGroups lists the groups that are allowed to edit a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) SearchGroups(ctx context.Context, opt *QualityprofilesSearchGroupsOptions) (v *QualityprofilesSearchGroups, resp *http.Response, err error) {
	err = s.ValidateSearchGroupsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/search_groups", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesSearchGroups)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchUsers lists the users that are allowed to edit a quality profile.
// Requires one of the following permissions:
//   - 'Administer Quality Profiles'
//   - Edit right on the specified quality profile
func (s *QualityprofilesService) SearchUsers(ctx context.Context, opt *QualityprofilesSearchUsersOptions) (v *QualityprofilesSearchUsers, resp *http.Response, err error) {
	err = s.ValidateSearchUsersOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/search_users", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesSearchUsers)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SetDefault selects the default profile for a given language.
// Requires the 'Administer Quality Profiles' permission.
func (s *QualityprofilesService) SetDefault(ctx context.Context, opt *QualityprofilesSetDefaultOptions) (resp *http.Response, err error) {
	err = s.ValidateSetDefaultOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "qualityprofiles/set_default", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Show shows a quality profile.
func (s *QualityprofilesService) Show(ctx context.Context, opt *QualityprofilesShowOptions) (v *QualityprofilesShow, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "qualityprofiles/show", opt)
	if err != nil {
		return
	}

	v = new(QualityprofilesShow)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateActivateRuleOpt validates the options for activating a rule.
func (s *QualityprofilesService) ValidateActivateRuleOpt(opt *QualityprofilesActivateRuleOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesActivateRuleOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Rule, "Rule")
	if err != nil {
		return err
	}

	// Validate that Impacts and Severity are not both set
	if len(opt.Impacts) > 0 && opt.Severity != "" {
		return NewValidationError("QualityprofilesActivateRuleOption", "cannot set both Impacts and Severity", ErrInvalidValue)
	}

	// Validate Severity if provided
	err = IsValueAuthorized(opt.Severity, allowedRuleSeverities, "Severity")
	if err != nil {
		return err
	}

	// Validate Impacts map values
	err = ValidateMapKeys(opt.Impacts, allowedImpactSoftwareQualities, "Impacts")
	if err != nil {
		return err
	}

	err = ValidateMapValues(opt.Impacts, allowedRuleImpactSeverities, "Impacts")
	if err != nil {
		return err
	}

	return nil
}

// ValidateActivateRulesOpt validates the options for bulk activating rules.
//
//nolint:cyclop,funlen // Validation functions are naturally complex due to multiple checks
func (s *QualityprofilesService) ValidateActivateRulesOpt(opt *QualityprofilesActivateRulesOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesActivateRulesOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.TargetKey, "TargetKey")
	if err != nil {
		return err
	}

	// Validate severity values
	err = AreValuesAuthorized(opt.ActiveImpactSeverities, allowedRuleImpactSeverities, "ActiveImpactSeverities")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.ActiveSeverities, allowedRuleSeverities, "ActiveSeverities")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.Severities, allowedRuleSeverities, "Severities")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.ImpactSeverities, allowedRuleImpactSeverities, "ImpactSeverities")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.TargetSeverity, allowedRuleSeverities, "TargetSeverity")
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

	// Validate languages
	err = ValidateLanguages(opt.Languages)
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
		allowed := map[string]struct{}{"key": {}, "name": {}, FieldCreatedAt: {}, "updatedAt": {}}

		err = IsValueAuthorized(opt.Sort, allowed, "Sort")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateAddGroupOpt validates the options for adding a group to a quality profile.
func (s *QualityprofilesService) ValidateAddGroupOpt(opt *QualityprofilesAddGroupOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesAddGroupOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Group, "Group")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAddProjectOpt validates the options for associating a project.
func (s *QualityprofilesService) ValidateAddProjectOpt(opt *QualityprofilesAddProjectOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesAddProjectOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAddUserOpt validates the options for adding a user to a quality profile.
func (s *QualityprofilesService) ValidateAddUserOpt(opt *QualityprofilesAddUserOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesAddUserOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateBackupOpt validates the options for backing up a profile.
func (s *QualityprofilesService) ValidateBackupOpt(opt *QualityprofilesBackupOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesBackupOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateChangeParentOpt validates the options for changing a profile's parent.
func (s *QualityprofilesService) ValidateChangeParentOpt(opt *QualityprofilesChangeParentOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesChangeParentOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	// ParentQualityProfile is optional - if not provided, breaks the inheritance link

	return nil
}

// ValidateChangelogOpt validates the options for getting a profile's changelog.
func (s *QualityprofilesService) ValidateChangelogOpt(opt *QualityprofilesChangelogOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesChangelogOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	// Validate FilterMode if provided
	if opt.FilterMode != "" {
		allowed := map[string]struct{}{"MQR": {}, "STANDARD": {}}

		err = IsValueAuthorized(opt.FilterMode, allowed, "FilterMode")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateCompareOpt validates the options for comparing two profiles.
func (s *QualityprofilesService) ValidateCompareOpt(opt *QualityprofilesCompareOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesCompareOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.LeftKey, "LeftKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RightKey, "RightKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCopyOpt validates the options for copying a profile.
func (s *QualityprofilesService) ValidateCopyOpt(opt *QualityprofilesCopyOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesCopyOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.FromKey, "FromKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ToName, "ToName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ToName, MaxQualityProfileNameLength, "ToName")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateOpt validates the options for creating a profile.
func (s *QualityprofilesService) ValidateCreateOpt(opt *QualityprofilesCreateOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesCreateOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxQualityProfileNameLength, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeactivateRuleOpt validates the options for deactivating a rule.
func (s *QualityprofilesService) ValidateDeactivateRuleOpt(opt *QualityprofilesDeactivateRuleOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesDeactivateRuleOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Rule, "Rule")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeactivateRulesOpt validates the options for bulk deactivating rules.
func (s *QualityprofilesService) ValidateDeactivateRulesOpt(opt *QualityprofilesDeactivateRulesOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesDeactivateRulesOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.TargetKey, "TargetKey")
	if err != nil {
		return err
	}

	// Note: DeactivateRulesOptions uses string fields for filters instead of slices,
	// which limits granular validation. The API will validate the format.

	return nil
}

// ValidateDeleteOpt validates the options for deleting a profile.
func (s *QualityprofilesService) ValidateDeleteOpt(opt *QualityprofilesDeleteOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesDeleteOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateExportOpt validates the options for exporting a profile.
func (s *QualityprofilesService) ValidateExportOpt(opt *QualityprofilesExportOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesExportOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	// QualityProfile is optional - if empty, default profile is exported

	return nil
}

// ValidateInheritanceOpt validates the options for getting inheritance info.
func (s *QualityprofilesService) ValidateInheritanceOpt(opt *QualityprofilesInheritanceOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesInheritanceOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateProjectsOpt validates the options for listing associated projects.
func (s *QualityprofilesService) ValidateProjectsOpt(opt *QualityprofilesProjectsOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesProjectsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	// Validate Selected if provided
	err = IsValueAuthorized(opt.Selected, allowedSelectedFilters, "Selected")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRemoveGroupOpt validates the options for removing group permissions.
func (s *QualityprofilesService) ValidateRemoveGroupOpt(opt *QualityprofilesRemoveGroupOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesRemoveGroupOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Group, "Group")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRemoveProjectOpt validates the options for removing a project association.
func (s *QualityprofilesService) ValidateRemoveProjectOpt(opt *QualityprofilesRemoveProjectOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesRemoveProjectOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRemoveUserOpt validates the options for removing user permissions.
func (s *QualityprofilesService) ValidateRemoveUserOpt(opt *QualityprofilesRemoveUserOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesRemoveUserOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRenameOpt validates the options for renaming a profile.
func (s *QualityprofilesService) ValidateRenameOpt(opt *QualityprofilesRenameOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesRenameOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxQualityProfileNameLength, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRestoreOpt validates the options for restoring a profile.
func (s *QualityprofilesService) ValidateRestoreOpt(opt *QualityprofilesRestoreOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesRestoreOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Backup, "Backup")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for searching profiles.
func (s *QualityprofilesService) ValidateSearchOpt(opt *QualityprofilesSearchOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesSearchOption", "cannot be nil", ErrMissingRequired)
	}

	// Validate Language if provided
	if opt.Language != "" {
		err := ValidateLanguage(opt.Language)
		if err != nil {
			return err
		}
	}

	// All parameters are optional for Search

	return nil
}

// ValidateSearchGroupsOpt validates the options for searching groups.
func (s *QualityprofilesService) ValidateSearchGroupsOpt(opt *QualityprofilesSearchGroupsOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesSearchGroupsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	// Validate Selected if provided
	err = IsValueAuthorized(opt.Selected, allowedSelectedFilters, "Selected")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchUsersOpt validates the options for searching users.
func (s *QualityprofilesService) ValidateSearchUsersOpt(opt *QualityprofilesSearchUsersOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesSearchUsersOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	// Validate Selected if provided
	err = IsValueAuthorized(opt.Selected, allowedSelectedFilters, "Selected")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetDefaultOpt validates the options for setting the default profile.
func (s *QualityprofilesService) ValidateSetDefaultOpt(opt *QualityprofilesSetDefaultOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesSetDefaultOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Language, "Language")
	if err != nil {
		return err
	}

	err = ValidateLanguage(opt.Language)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.QualityProfile, "QualityProfile")
	if err != nil {
		return err
	}

	return nil
}

// ValidateShowOpt validates the options for showing a profile.
func (s *QualityprofilesService) ValidateShowOpt(opt *QualityprofilesShowOptions) error {
	if opt == nil {
		return NewValidationError("QualityprofilesShowOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return nil
}

// ChangelogAll fetches all pages from Changelog and returns a flat slice of events.
func (s *QualityprofilesService) ChangelogAll(ctx context.Context, opt *QualityprofilesChangelogOptions) ([]ChangelogEvent, *http.Response, error) {
	err := s.ValidateChangelogOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	o := *opt

	return allPages(ctx, &o.Page, &o.PageSize, func(ctx context.Context) ([]ChangelogEvent, int64, *http.Response, error) {
		r, resp, err := s.Changelog(ctx, &o)
		if err != nil {
			return nil, 0, resp, err
		}

		return r.Events, r.Paging.Total, resp, nil
	})
}

// ProjectsAll fetches all pages from Projects and returns a flat slice of projects.
func (s *QualityprofilesService) ProjectsAll(ctx context.Context, opt *QualityprofilesProjectsOptions) ([]QualityprofilesProfileProject, *http.Response, error) {
	err := s.ValidateProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	o := *opt

	return allPages(ctx, &o.Page, &o.PageSize, func(ctx context.Context) ([]QualityprofilesProfileProject, int64, *http.Response, error) {
		r, resp, err := s.Projects(ctx, &o)
		if err != nil {
			return nil, 0, resp, err
		}

		return r.Results, r.Paging.Total, resp, nil
	})
}

// SearchGroupsAll fetches all pages from SearchGroups and returns a flat slice of groups.
func (s *QualityprofilesService) SearchGroupsAll(ctx context.Context, opt *QualityprofilesSearchGroupsOptions) ([]QualityprofilesProfileGroup, *http.Response, error) {
	err := s.ValidateSearchGroupsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	o := *opt

	return allPages(ctx, &o.Page, &o.PageSize, func(ctx context.Context) ([]QualityprofilesProfileGroup, int64, *http.Response, error) {
		r, resp, err := s.SearchGroups(ctx, &o)
		if err != nil {
			return nil, 0, resp, err
		}

		return r.Groups, r.Paging.Total, resp, nil
	})
}

// SearchUsersAll fetches all pages from SearchUsers and returns a flat slice of users.
func (s *QualityprofilesService) SearchUsersAll(ctx context.Context, opt *QualityprofilesSearchUsersOptions) ([]QualityprofilesProfileUser, *http.Response, error) {
	err := s.ValidateSearchUsersOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	o := *opt

	return allPages(ctx, &o.Page, &o.PageSize, func(ctx context.Context) ([]QualityprofilesProfileUser, int64, *http.Response, error) {
		r, resp, err := s.SearchUsers(ctx, &o)
		if err != nil {
			return nil, 0, resp, err
		}

		return r.Users, r.Paging.Total, resp, nil
	})
}

// -----------------------------------------------------------------------------
// Conversion Functions
// -----------------------------------------------------------------------------

// convertActivateRuleOptForURL converts QualityprofilesActivateRuleOptions to a URL-encodable format.
func (s *QualityprofilesService) convertActivateRuleOptForURL(opt *QualityprofilesActivateRuleOptions) *qualityprofilesActivateRuleURLOptions {
	//nolint:exhaustruct // Only populate fields that have values
	urlOpt := &qualityprofilesActivateRuleURLOptions{
		Key:             opt.Key,
		Rule:            opt.Rule,
		PrioritizedRule: opt.PrioritizedRule,
		Reset:           opt.Reset,
		Severity:        opt.Severity,
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

// qualityprofilesActivateRuleURLOptions is the URL-encodable version of QualityprofilesActivateRuleOption.
//
//nolint:govet // Field alignment is less important than logical grouping
type qualityprofilesActivateRuleURLOptions struct {
	Key             string `url:"key,omitempty"`
	Rule            string `url:"rule,omitempty"`
	Impacts         string `url:"impacts,omitempty"`
	Params          string `url:"params,omitempty"`
	PrioritizedRule bool   `url:"prioritizedRule,omitempty"`
	Reset           bool   `url:"reset,omitempty"`
	Severity        string `url:"severity,omitempty"`
}
