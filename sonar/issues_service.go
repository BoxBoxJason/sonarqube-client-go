package sonargo

import (
	"net/http"
)

// IssuesService handles communication with the Issues related methods of the SonarQube API.
// Issues represent code problems detected by SonarQube during analysis.
type IssuesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// =============================================================================
// Shared Types
// =============================================================================

// Issue represents an issue detected by SonarQube.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type Issue struct {
	// Key is the unique identifier of the issue.
	Key string `json:"key,omitempty"`
	// Component is the key of the component where the issue was found.
	Component string `json:"component,omitempty"`
	// Project is the key of the project containing the issue.
	Project string `json:"project,omitempty"`
	// Rule is the key of the rule that raised this issue.
	Rule string `json:"rule,omitempty"`
	// Message is the main message describing the issue.
	Message string `json:"message,omitempty"`
	// Line is the line number where the issue was found.
	Line int64 `json:"line,omitempty"`
	// Hash is the hash of the line content for tracking purposes.
	Hash string `json:"hash,omitempty"`
	// IssueStatus is the current status of the issue (OPEN, CONFIRMED, RESOLVED, etc.).
	IssueStatus string `json:"issueStatus,omitempty"`
	// Author is the SCM author of the code where the issue was found.
	Author string `json:"author,omitempty"`
	// Assignee is the login of the user assigned to this issue.
	Assignee string `json:"assignee,omitempty"`
	// Effort is the estimated time to fix the issue.
	Effort string `json:"effort,omitempty"`
	// Debt is the estimated time to fix the issue (legacy field, same as Effort).
	Debt string `json:"debt,omitempty"`
	// CreationDate is the timestamp when the issue was created.
	CreationDate string `json:"creationDate,omitempty"`
	// UpdateDate is the timestamp when the issue was last updated.
	UpdateDate string `json:"updateDate,omitempty"`
	// CleanCodeAttribute is the clean code attribute of the issue.
	CleanCodeAttribute string `json:"cleanCodeAttribute,omitempty"`
	// CleanCodeAttributeCategory is the category of the clean code attribute.
	CleanCodeAttributeCategory string `json:"cleanCodeAttributeCategory,omitempty"`
	// RuleDescriptionContextKey is the context key for the rule description.
	RuleDescriptionContextKey string `json:"ruleDescriptionContextKey,omitempty"`
	// Severity is the severity level of the issue (BLOCKER, CRITICAL, MAJOR, MINOR, INFO).
	Severity string `json:"severity,omitempty"`
	// Status is the status of the issue (legacy field).
	Status string `json:"status,omitempty"`
	// Type is the type of the issue (BUG, VULNERABILITY, CODE_SMELL).
	Type string `json:"type,omitempty"`
	// LinkedTicketStatus is the status of the linked ticket (if any).
	LinkedTicketStatus string `json:"linkedTicketStatus,omitempty"`
	// TextRange is the text range where the issue was found.
	TextRange TextRange `json:"textRange,omitzero"`
	// Actions are the available actions for this issue.
	Actions []string `json:"actions,omitempty"`
	// Transitions are the available workflow transitions for this issue.
	Transitions []string `json:"transitions,omitempty"`
	// Tags is the list of tags associated with this issue.
	Tags []string `json:"tags,omitempty"`
	// InternalTags is the list of internal tags.
	InternalTags []string `json:"internalTags,omitempty"`
	// CodeVariants is the list of code variants affected by this issue.
	CodeVariants []string `json:"codeVariants,omitempty"`
	// Comments is the list of comments on this issue.
	Comments []IssueComment `json:"comments,omitempty"`
	// Impacts is the list of impacts on software quality.
	Impacts []IssueImpact `json:"impacts,omitempty"`
	// Flows is the list of code flows related to this issue.
	Flows []IssueFlow `json:"flows,omitempty"`
	// MessageFormattings is the list of message formatting rules.
	MessageFormattings []MessageFormatting `json:"messageFormattings,omitempty"`
	// QuickFixAvailable indicates if a quick fix is available for this issue.
	QuickFixAvailable bool `json:"quickFixAvailable,omitempty"`
	// PrioritizedRule indicates if the rule is prioritized.
	PrioritizedRule bool `json:"prioritizedRule,omitempty"`
}

// IssueComment represents a comment on an issue.
type IssueComment struct {
	// Key is the unique identifier of the comment.
	Key string `json:"key,omitempty"`
	// Login is the login of the user who created the comment.
	Login string `json:"login,omitempty"`
	// HTMLText is the HTML-formatted text of the comment.
	HTMLText string `json:"htmlText,omitempty"`
	// Markdown is the Markdown-formatted text of the comment.
	Markdown string `json:"markdown,omitempty"`
	// CreatedAt is the timestamp when the comment was created.
	CreatedAt string `json:"createdAt,omitempty"`
	// Updatable indicates if the comment can be updated by the current user.
	Updatable bool `json:"updatable,omitempty"`
}

// IssueImpact represents the impact of an issue on software quality.
type IssueImpact struct {
	// Severity is the severity of the impact (BLOCKER, HIGH, MEDIUM, LOW, INFO).
	Severity string `json:"severity,omitempty"`
	// SoftwareQuality is the software quality characteristic affected (MAINTAINABILITY, RELIABILITY, SECURITY).
	SoftwareQuality string `json:"softwareQuality,omitempty"`
}

// IssueFlow represents a code flow related to an issue.
type IssueFlow struct {
	// Locations is the list of locations in the flow.
	Locations []FlowLocation `json:"locations,omitempty"`
}

// FlowLocation represents a location in a code flow.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type FlowLocation struct {
	// Msg is the message for this location.
	Msg string `json:"msg,omitempty"`
	// TextRange is the text range of this location.
	TextRange TextRange `json:"textRange,omitzero"`
	// MsgFormattings is the list of message formatting rules.
	MsgFormattings []MessageFormatting `json:"msgFormattings,omitempty"`
}

// TextRange represents a range of text in a source file.
type TextRange struct {
	// StartLine is the starting line number (1-based).
	StartLine int64 `json:"startLine,omitempty"`
	// EndLine is the ending line number (1-based).
	EndLine int64 `json:"endLine,omitempty"`
	// StartOffset is the starting offset within the line (0-based).
	StartOffset int64 `json:"startOffset,omitempty"`
	// EndOffset is the ending offset within the line (0-based).
	EndOffset int64 `json:"endOffset,omitempty"`
}

// MessageFormatting represents a formatting rule for messages.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type MessageFormatting struct {
	// Start is the starting position of the formatting.
	Start int64 `json:"start,omitempty"`
	// End is the ending position of the formatting.
	End int64 `json:"end,omitempty"`
	// Type is the type of formatting (CODE, etc.).
	Type string `json:"type,omitempty"`
}

// IssueComponent represents a component in issue responses.
type IssueComponent struct {
	// Key is the unique identifier of the component.
	Key string `json:"key,omitempty"`
	// Name is the short name of the component.
	Name string `json:"name,omitempty"`
	// LongName is the long name of the component.
	LongName string `json:"longName,omitempty"`
	// Path is the path of the component.
	Path string `json:"path,omitempty"`
	// Qualifier is the component qualifier (TRK, DIR, FIL, etc.).
	Qualifier string `json:"qualifier,omitempty"`
	// ID is the internal identifier of the component.
	ID int64 `json:"id,omitempty"`
	// ProjectID is the internal identifier of the project.
	ProjectID int64 `json:"projectId,omitempty"`
	// SubProjectID is the internal identifier of the sub-project.
	SubProjectID int64 `json:"subProjectId,omitempty"`
	// Enabled indicates if the component is enabled.
	Enabled bool `json:"enabled,omitempty"`
}

// IssueRule represents a rule in issue responses.
type IssueRule struct {
	// Key is the unique identifier of the rule.
	Key string `json:"key,omitempty"`
	// Name is the display name of the rule.
	Name string `json:"name,omitempty"`
	// Lang is the language key of the rule.
	Lang string `json:"lang,omitempty"`
	// LangName is the display name of the language.
	LangName string `json:"langName,omitempty"`
	// Status is the status of the rule.
	Status string `json:"status,omitempty"`
}

// IssueUser represents a user in issue responses.
type IssueUser struct {
	// Login is the login of the user.
	Login string `json:"login,omitempty"`
	// Name is the display name of the user.
	Name string `json:"name,omitempty"`
	// Email is the email of the user.
	Email string `json:"email,omitempty"`
	// Avatar is the avatar hash of the user.
	Avatar string `json:"avatar,omitempty"`
	// Active indicates if the user is active.
	Active bool `json:"active,omitempty"`
}

// ChangelogEntry represents an entry in the issue changelog.
type ChangelogEntry struct {
	// User is the login of the user who made the change.
	User string `json:"user,omitempty"`
	// UserName is the display name of the user.
	UserName string `json:"userName,omitempty"`
	// ExternalUser is the external user identifier.
	ExternalUser string `json:"externalUser,omitempty"`
	// Avatar is the avatar hash of the user.
	Avatar string `json:"avatar,omitempty"`
	// CreationDate is the timestamp of the change.
	CreationDate string `json:"creationDate,omitempty"`
	// WebhookSource is the source of the webhook that triggered the change.
	WebhookSource string `json:"webhookSource,omitempty"`
	// Diffs is the list of field changes.
	Diffs []ChangelogDiff `json:"diffs,omitempty"`
	// IsUserActive indicates if the user is active.
	IsUserActive bool `json:"isUserActive,omitempty"`
}

// ChangelogDiff represents a single field change in the changelog.
type ChangelogDiff struct {
	// Key is the name of the field that changed.
	Key string `json:"key,omitempty"`
	// OldValue is the previous value of the field.
	OldValue string `json:"oldValue,omitempty"`
	// NewValue is the new value of the field.
	NewValue string `json:"newValue,omitempty"`
}

// IssueFacet represents a facet in search results.
type IssueFacet struct {
	// Property is the name of the facet property.
	Property string `json:"property,omitempty"`
	// Values is the list of facet values with their counts.
	Values []IssueFacetValue `json:"values,omitempty"`
}

// IssueFacetValue represents a single facet value with its count.
type IssueFacetValue struct {
	// Val is the facet value.
	Val string `json:"val,omitempty"`
	// Count is the number of items matching this facet value.
	Count int64 `json:"count,omitempty"`
}

// ComponentTag represents a tag with its count.
type ComponentTag struct {
	// Key is the tag name.
	Key string `json:"key,omitempty"`
	// Value is the count of issues with this tag.
	Value int64 `json:"value,omitempty"`
}

// =============================================================================
// Response Types
// =============================================================================

// IssuesAddComment represents the response from adding a comment to an issue.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesAddComment struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesAssign represents the response from assigning an issue.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesAssign struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesAuthors represents the response from the authors endpoint.
type IssuesAuthors struct {
	// Authors is the list of SCM authors.
	Authors []string `json:"authors,omitempty"`
}

// IssuesBulkChange represents the response from a bulk change operation.
type IssuesBulkChange struct {
	// Total is the total number of issues processed.
	Total int64 `json:"total,omitempty"`
	// Success is the number of successfully changed issues.
	Success int64 `json:"success,omitempty"`
	// Ignored is the number of ignored issues.
	Ignored int64 `json:"ignored,omitempty"`
	// Failures is the number of failed changes.
	Failures int64 `json:"failures,omitempty"`
}

// IssuesChangelog represents the response from the changelog endpoint.
type IssuesChangelog struct {
	// Changelog is the list of changelog entries.
	Changelog []ChangelogEntry `json:"changelog,omitempty"`
}

// IssuesComponentTags represents the response from the component tags endpoint.
type IssuesComponentTags struct {
	// Tags is the list of tags with their counts.
	Tags []ComponentTag `json:"tags,omitempty"`
}

// IssuesDeleteComment represents the response from deleting a comment.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesDeleteComment struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesDoTransition represents the response from performing a transition.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesDoTransition struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesEditComment represents the response from editing a comment.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesEditComment struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesList represents the response from the list endpoint.
type IssuesList struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issues     []Issue          `json:"issues,omitempty"`
	Paging     Paging           `json:"paging,omitzero"`
}

// IssuesSearch represents the response from the search endpoint.
type IssuesSearch struct {
	Components []IssueComponent `json:"components,omitempty"`
	Facets     []IssueFacet     `json:"facets,omitempty"`
	Issues     []Issue          `json:"issues,omitempty"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
	Paging     Paging           `json:"paging,omitzero"`
}

// IssuesSetSeverity represents the response from setting severity.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesSetSeverity struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesSetTags represents the response from setting tags.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesSetTags struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesSetType represents the response from setting type.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesSetType struct {
	Components []IssueComponent `json:"components,omitempty"`
	Issue      Issue            `json:"issue,omitzero"`
	Rules      []IssueRule      `json:"rules,omitempty"`
	Users      []IssueUser      `json:"users,omitempty"`
}

// IssuesTags represents the response from the tags endpoint.
type IssuesTags struct {
	// Tags is the list of available tags.
	Tags []string `json:"tags,omitempty"`
}

// =============================================================================
// Option Types
// =============================================================================

// IssuesAddCommentOption contains options for adding a comment to an issue.
type IssuesAddCommentOption struct {
	// Issue is the key of the issue to comment on (required).
	Issue string `url:"issue,omitempty"`
	// Text is the comment text (required).
	Text string `url:"text,omitempty"`
}

// IssuesAnticipatedTransitionsOption contains options for anticipated transitions.
type IssuesAnticipatedTransitionsOption struct {
	// ProjectKey is the key of the project (required).
	ProjectKey string `url:"projectKey,omitempty"`
}

// IssuesAssignOption contains options for assigning an issue.
type IssuesAssignOption struct {
	// Issue is the key of the issue to assign (required).
	Issue string `url:"issue,omitempty"`
	// Assignee is the login of the assignee. When not set, it will unassign the issue.
	// Use '_me' to assign to the current user.
	Assignee string `url:"assignee,omitempty"`
}

// IssuesAuthorsOption contains options for listing authors.
type IssuesAuthorsOption struct {
	// Project is the project key to limit the search.
	Project string `url:"project,omitempty"`
	// Query limits the search to authors that contain the supplied string.
	Query string `url:"q,omitempty"`
	// PageSize is the maximum number of results to return (must be between 1 and 100).
	PageSize int64 `url:"ps,omitempty"`
}

// IssuesBulkChangeOption contains options for bulk changing issues.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesBulkChangeOption struct {
	// Issues is the list of issue keys to change (required).
	Issues []string `url:"issues,omitempty,comma"`
	// AddTags is the list of tags to add.
	AddTags []string `url:"add_tags,omitempty,comma"`
	// RemoveTags is the list of tags to remove.
	RemoveTags []string `url:"remove_tags,omitempty,comma"`
	// Assign is the login of the user to assign, or empty to unassign.
	Assign string `url:"assign,omitempty"`
	// Comment is the comment to add (only added if type or severity changes).
	Comment string `url:"comment,omitempty"`
	// DoTransition is the transition to perform.
	// Allowed values: confirm, unconfirm, reopen, resolve, falsepositive, wontfix, resolveasreviewed, resetastoreview, accept
	DoTransition string `url:"do_transition,omitempty"`
	// SetSeverity is the new severity to set.
	// Allowed values: BLOCKER, CRITICAL, MAJOR, MINOR, INFO
	SetSeverity string `url:"set_severity,omitempty"`
	// SetType is the new type to set.
	// Allowed values: BUG, VULNERABILITY, CODE_SMELL, SECURITY_HOTSPOT
	SetType string `url:"set_type,omitempty"`
	// SendNotifications indicates whether to send notifications.
	SendNotifications bool `url:"sendNotifications,omitempty"`
}

// IssuesChangelogOption contains options for retrieving issue changelog.
type IssuesChangelogOption struct {
	// Issue is the key of the issue (required).
	Issue string `url:"issue,omitempty"`
}

// IssuesComponentTagsOption contains options for listing component tags.
type IssuesComponentTagsOption struct {
	// ComponentUuid is the UUID of the component (required).
	ComponentUuid string `url:"componentUuid,omitempty"`
	// CreatedAfter filters issues created after the given date.
	CreatedAfter string `url:"createdAfter,omitempty"`
	// PageSize is the maximum number of tags to return.
	PageSize int64 `url:"ps,omitempty"`
}

// IssuesDeleteCommentOption contains options for deleting a comment.
type IssuesDeleteCommentOption struct {
	// Comment is the key of the comment to delete (required).
	Comment string `url:"comment,omitempty"`
}

// IssuesDoTransitionOption contains options for performing a transition.
type IssuesDoTransitionOption struct {
	// Issue is the key of the issue (required).
	Issue string `url:"issue,omitempty"`
	// Transition is the transition to perform (required).
	// Allowed values: confirm, unconfirm, reopen, resolve, falsepositive, wontfix, resolveasreviewed, resetastoreview, accept, close
	Transition string `url:"transition,omitempty"`
}

// IssuesEditCommentOption contains options for editing a comment.
type IssuesEditCommentOption struct {
	// Comment is the key of the comment to edit (required).
	Comment string `url:"comment,omitempty"`
	// Text is the new comment text (required).
	Text string `url:"text,omitempty"`
}

// IssuesListOption contains options for listing issues.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesListOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// Project is the project key (one of Project or Component is required).
	Project string `url:"project,omitempty"`
	// Component is the component key.
	Component string `url:"component,omitempty"`
	// Branch is the branch key (not available in Community Edition).
	Branch string `url:"branch,omitempty"`
	// PullRequest is the pull request ID (not available in Community Edition).
	PullRequest string `url:"pullRequest,omitempty"`
	// Types is the list of issue types to filter by.
	// Allowed values: BUG, VULNERABILITY, CODE_SMELL, SECURITY_HOTSPOT
	Types []string `url:"types,omitempty,comma"`
	// Resolved filters by resolved status. If nil, returns all issues.
	Resolved *bool `url:"resolved,omitempty"`
	// InNewCodePeriod filters issues created in the new code period.
	InNewCodePeriod bool `url:"inNewCodePeriod,omitempty"`
}

// IssuesPullOption contains options for pulling issues.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesPullOption struct {
	// ProjectKey is the project key (required).
	ProjectKey string `url:"projectKey,omitempty"`
	// BranchName is the branch name to fetch issues for.
	BranchName string `url:"branchName,omitempty"`
	// Languages is the list of languages to filter by.
	Languages []string `url:"languages,omitempty,comma"`
	// RuleRepositories is the list of rule repositories to filter by.
	RuleRepositories []string `url:"ruleRepositories,omitempty,comma"`
	// ChangedSince is the timestamp to filter issues modified after.
	ChangedSince string `url:"changedSince,omitempty"`
	// ResolvedOnly returns only resolved issues if true.
	ResolvedOnly bool `url:"resolvedOnly,omitempty"`
}

// IssuesPullTaintOption contains options for pulling taint vulnerabilities.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesPullTaintOption struct {
	// ProjectKey is the project key (required).
	ProjectKey string `url:"projectKey,omitempty"`
	// BranchName is the branch name to fetch taint vulnerabilities for.
	BranchName string `url:"branchName,omitempty"`
	// Languages is the list of languages to filter by.
	Languages []string `url:"languages,omitempty,comma"`
	// ChangedSince is the timestamp to filter issues modified after.
	ChangedSince string `url:"changedSince,omitempty"`
}

// IssuesReindexOption contains options for reindexing issues.
type IssuesReindexOption struct {
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
}

// IssuesSearchOption contains options for searching issues.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type IssuesSearchOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// AdditionalFields is the list of optional fields to return.
	// Allowed values: _all, comments, languages, rules, ruleDescriptionContextKey, transitions, actions, users
	AdditionalFields []string `url:"additionalFields,omitempty,comma"`
	// Assignees is the list of assignee logins. Use '__me__' for current user.
	Assignees []string `url:"assignees,omitempty,comma"`
	// Author is the SCM account to filter by.
	Author string `url:"author,omitempty"`
	// Branch is the branch key (not available in Community Edition).
	Branch string `url:"branch,omitempty"`
	// Casa is the list of CASA categories.
	Casa []string `url:"casa,omitempty,comma"`
	// CleanCodeAttributeCategories is the list of clean code attribute categories.
	// Allowed values: ADAPTABLE, CONSISTENT, INTENTIONAL, RESPONSIBLE
	CleanCodeAttributeCategories []string `url:"cleanCodeAttributeCategories,omitempty,comma"`
	// CodeVariants is the list of code variants.
	CodeVariants []string `url:"codeVariants,omitempty,comma"`
	// ComplianceStandards is the list of compliance standards.
	ComplianceStandards []string `url:"complianceStandards,omitempty,comma"`
	// Components is the list of component keys.
	Components []string `url:"components,omitempty,comma"`
	// CreatedAfter filters issues created after the given date.
	CreatedAfter string `url:"createdAfter,omitempty"`
	// CreatedAt filters issues created during a specific analysis.
	CreatedAt string `url:"createdAt,omitempty"`
	// CreatedBefore filters issues created before the given date.
	CreatedBefore string `url:"createdBefore,omitempty"`
	// CreatedInLast filters issues created in the last time span.
	CreatedInLast string `url:"createdInLast,omitempty"`
	// Cwe is the list of CWE identifiers.
	Cwe []string `url:"cwe,omitempty,comma"`
	// Directories is the list of directories to filter by.
	Directories []string `url:"directories,omitempty,comma"`
	// Facets is the list of facets to compute.
	// Allowed values: projects, files, assigned_to_me, severities, statuses, resolutions, rules, assignees, author, directories, scopes, languages, tags, types, pciDss-3.2, pciDss-4.0, owaspAsvs-4.0, owaspMobileTop10-2024, owaspTop10, owaspTop10-2021, stig-ASD_V5R3, casa, sansTop25, cwe, createdAt, sonarsourceSecurity, codeVariants, cleanCodeAttributeCategories, impactSoftwareQualities, impactSeverities, issueStatuses, prioritizedRule, complianceStandards
	Facets []string `url:"facets,omitempty,comma"`
	// Files is the list of files to filter by.
	Files []string `url:"files,omitempty,comma"`
	// FixedInPullRequest filters issues fixed in the specified pull request.
	FixedInPullRequest string `url:"fixedInPullRequest,omitempty"`
	// ImpactSeverities is the list of impact severities to filter by.
	// Allowed values: BLOCKER, HIGH, MEDIUM, LOW, INFO
	ImpactSeverities []string `url:"impactSeverities,omitempty,comma"`
	// ImpactSoftwareQualities is the list of software qualities to filter by.
	// Allowed values: MAINTAINABILITY, RELIABILITY, SECURITY
	ImpactSoftwareQualities []string `url:"impactSoftwareQualities,omitempty,comma"`
	// InNewCodePeriod filters issues created in the new code period.
	InNewCodePeriod bool `url:"inNewCodePeriod,omitempty"`
	// IssueStatuses is the list of issue statuses to filter by.
	// Allowed values: OPEN, CONFIRMED, FALSE_POSITIVE, ACCEPTED, FIXED, IN_SANDBOX
	IssueStatuses []string `url:"issueStatuses,omitempty,comma"`
	// Issues is the list of issue keys to retrieve.
	Issues []string `url:"issues,omitempty,comma"`
	// Languages is the list of languages to filter by.
	Languages []string `url:"languages,omitempty,comma"`
	// OnComponentOnly returns only issues at the component level.
	OnComponentOnly bool `url:"onComponentOnly,omitempty"`
	// OwaspAsvs40 is the list of OWASP ASVS v4.0 categories.
	OwaspAsvs40 []string `url:"owaspAsvs-4.0,omitempty,comma"`
	// OwaspAsvsLevel is the OWASP ASVS level.
	// Allowed values: 1, 2, 3
	OwaspAsvsLevel int64 `url:"owaspAsvsLevel,omitempty"`
	// OwaspMobileTop102024 is the list of OWASP Mobile Top 10 2024 categories.
	// Allowed values: m1, m2, m3, m4, m5, m6, m7, m8, m9, m10
	OwaspMobileTop102024 []string `url:"owaspMobileTop10-2024,omitempty,comma"`
	// OwaspTop10 is the list of OWASP Top 10 2017 categories.
	// Allowed values: a1, a2, a3, a4, a5, a6, a7, a8, a9, a10
	OwaspTop10 []string `url:"owaspTop10,omitempty,comma"`
	// OwaspTop102021 is the list of OWASP Top 10 2021 categories.
	// Allowed values: a1, a2, a3, a4, a5, a6, a7, a8, a9, a10
	OwaspTop102021 []string `url:"owaspTop10-2021,omitempty,comma"`
	// PciDss32 is the list of PCI DSS v3.2 categories.
	PciDss32 []string `url:"pciDss-3.2,omitempty,comma"`
	// PciDss40 is the list of PCI DSS v4.0 categories.
	PciDss40 []string `url:"pciDss-4.0,omitempty,comma"`
	// PrioritizedRule filters by prioritized rule status.
	PrioritizedRule bool `url:"prioritizedRule,omitempty"`
	// Projects is the list of project keys.
	Projects []string `url:"projects,omitempty,comma"`
	// PullRequest is the pull request ID (not available in Community Edition).
	PullRequest string `url:"pullRequest,omitempty"`
	// Resolutions is the list of resolutions to filter by.
	// Allowed values: FALSE-POSITIVE, WONTFIX, FIXED, REMOVED
	Resolutions []string `url:"resolutions,omitempty,comma"`
	// Resolved filters by resolved status. If nil, returns all issues.
	Resolved *bool `url:"resolved,omitempty"`
	// Rules is the list of rule keys to filter by.
	Rules []string `url:"rules,omitempty,comma"`
	// Sort is the sort field.
	// Allowed values: CREATION_DATE, CLOSE_DATE, SEVERITY, STATUS, FILE_LINE, HOTSPOTS, UPDATE_DATE
	Sort string `url:"s,omitempty"`
	// SansTop25 is the list of SANS Top 25 categories.
	// Allowed values: insecure-interaction, risky-resource, porous-defenses
	//
	// Deprecated: This filter is deprecated since 10.0
	SansTop25 []string `url:"sansTop25,omitempty,comma"`
	// Scopes is the list of scopes to filter by.
	// Allowed values: MAIN, TEST
	Scopes []string `url:"scopes,omitempty,comma"`
	// Severities is the list of severities to filter by.
	// Allowed values: BLOCKER, CRITICAL, MAJOR, MINOR, INFO
	Severities []string `url:"severities,omitempty,comma"`
	// SonarsourceSecurity is the list of SonarSource security categories.
	// Allowed values: buffer-overflow, sql-injection, rce, object-injection, command-injection, path-traversal-injection, ldap-injection, xpath-injection, log-injection, xxe, xss, dos, ssrf, csrf, http-response-splitting, open-redirect, weak-cryptography, auth, insecure-conf, file-manipulation, encrypt-data, traceability, permission, others
	SonarsourceSecurity []string `url:"sonarsourceSecurity,omitempty,comma"`
	// Statuses is the list of statuses to filter by.
	Statuses []string `url:"statuses,omitempty,comma"`
	// StigAsdV5R3 is the list of STIG V5R3 categories.
	StigAsdV5R3 []string `url:"stig-ASD_V5R3,omitempty,comma"`
	// Tags is the list of tags to filter by.
	Tags []string `url:"tags,omitempty,comma"`
	// TimeZone is the timezone for date resolution.
	TimeZone string `url:"timeZone,omitempty"`
	// Types is the list of issue types to filter by.
	// Allowed values: BUG, VULNERABILITY, CODE_SMELL
	Types []string `url:"types,omitempty,comma"`
	// Assigned filters by assigned status.
	Assigned *bool `url:"assigned,omitempty"`
	// Asc sorts in ascending order.
	Asc bool `url:"asc,omitempty"`
}

// IssuesSetSeverityOption contains options for setting severity.
type IssuesSetSeverityOption struct {
	// Issue is the key of the issue (required).
	Issue string `url:"issue,omitempty"`
	// Severity is the new severity level.
	// Allowed values: BLOCKER, CRITICAL, MAJOR, MINOR, INFO
	Severity string `url:"severity,omitempty"`
	// Impact is the override of impact severity (cannot be used with Severity).
	Impact string `url:"impact,omitempty"`
}

// IssuesSetTagsOption contains options for setting tags.
type IssuesSetTagsOption struct {
	// Issue is the key of the issue (required).
	Issue string `url:"issue,omitempty"`
	// Tags is the list of tags to set. Empty list removes all tags.
	Tags []string `url:"tags,omitempty,comma"`
}

// IssuesSetTypeOption contains options for setting type.
type IssuesSetTypeOption struct {
	// Issue is the key of the issue (required).
	Issue string `url:"issue,omitempty"`
	// Type is the new issue type (required).
	// Allowed values: BUG, VULNERABILITY, CODE_SMELL
	Type string `url:"type,omitempty"`
}

// IssuesTagsOption contains options for listing tags.
type IssuesTagsOption struct {
	// Project is the project key.
	Project string `url:"project,omitempty"`
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Query limits the search to tags containing the supplied string.
	Query string `url:"q,omitempty"`
	// PageSize is the maximum number of tags to return.
	PageSize int64 `url:"ps,omitempty"`
	// All includes tags from all branches if true.
	All bool `url:"all,omitempty"`
}

// =============================================================================
// Service Methods
// =============================================================================

// AddComment adds a comment to an issue.
// Requires authentication and 'Browse' permission on the project of the specified issue.
func (s *IssuesService) AddComment(opt *IssuesAddCommentOption) (v *IssuesAddComment, resp *http.Response, err error) {
	err = s.ValidateAddCommentOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/add_comment", opt)
	if err != nil {
		return
	}

	v = new(IssuesAddComment)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// AnticipatedTransitions receives a list of anticipated transitions for not yet discovered issues.
// Requires 'Administer Issues' permission on the specified project.
// Only 'falsepositive', 'wontfix' and 'accept' transitions are supported.
// Upon successful execution, the HTTP status code returned is 202 (Accepted).
func (s *IssuesService) AnticipatedTransitions(opt *IssuesAnticipatedTransitionsOption) (resp *http.Response, err error) {
	err = s.ValidateAnticipatedTransitionsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/anticipated_transitions", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Assign assigns or unassigns an issue.
// Requires authentication and 'Browse' permission on the project.
func (s *IssuesService) Assign(opt *IssuesAssignOption) (v *IssuesAssign, resp *http.Response, err error) {
	err = s.ValidateAssignOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/assign", opt)
	if err != nil {
		return
	}

	v = new(IssuesAssign)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Authors searches SCM accounts which match a given query.
// Requires authentication. Returns 503 when issue indexing is in progress.
func (s *IssuesService) Authors(opt *IssuesAuthorsOption) (v *IssuesAuthors, resp *http.Response, err error) {
	err = s.ValidateAuthorsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/authors", opt)
	if err != nil {
		return
	}

	v = new(IssuesAuthors)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// BulkChange performs bulk changes on issues. Up to 500 issues can be updated.
// Requires authentication.
func (s *IssuesService) BulkChange(opt *IssuesBulkChangeOption) (v *IssuesBulkChange, resp *http.Response, err error) {
	err = s.ValidateBulkChangeOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/bulk_change", opt)
	if err != nil {
		return
	}

	v = new(IssuesBulkChange)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Changelog displays the changelog of an issue.
// Requires 'Browse' permission on the project of the specified issue.
func (s *IssuesService) Changelog(opt *IssuesChangelogOption) (v *IssuesChangelog, resp *http.Response, err error) {
	err = s.ValidateChangelogOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/changelog", opt)
	if err != nil {
		return
	}

	v = new(IssuesChangelog)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ComponentTags lists tags for issues under a given component.
// Returns 503 when issue indexing is in progress.
func (s *IssuesService) ComponentTags(opt *IssuesComponentTagsOption) (v *IssuesComponentTags, resp *http.Response, err error) {
	err = s.ValidateComponentTagsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/component_tags", opt)
	if err != nil {
		return
	}

	v = new(IssuesComponentTags)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// DeleteComment deletes a comment.
// Requires authentication and 'Browse' permission on the project of the specified issue.
func (s *IssuesService) DeleteComment(opt *IssuesDeleteCommentOption) (v *IssuesDeleteComment, resp *http.Response, err error) {
	err = s.ValidateDeleteCommentOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/delete_comment", opt)
	if err != nil {
		return
	}

	v = new(IssuesDeleteComment)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// DoTransition performs a workflow transition on an issue.
// Requires authentication and 'Browse' permission on the project.
// Transitions 'accept', 'wontfix', and 'falsepositive' require 'Administer Issues' permission.
// Security hotspot transitions require 'Administer Security Hotspot' permission.
func (s *IssuesService) DoTransition(opt *IssuesDoTransitionOption) (v *IssuesDoTransition, resp *http.Response, err error) {
	err = s.ValidateDoTransitionOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/do_transition", opt)
	if err != nil {
		return
	}

	v = new(IssuesDoTransition)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// EditComment edits a comment.
// Requires authentication and 'Browse' permission on the project of the specified issue.
func (s *IssuesService) EditComment(opt *IssuesEditCommentOption) (v *IssuesEditComment, resp *http.Response, err error) {
	err = s.ValidateEditCommentOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/edit_comment", opt)
	if err != nil {
		return
	}

	v = new(IssuesEditComment)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// List lists issues in degraded mode when issue indexing is running.
// Either 'project' or 'component' parameter is required.
// Requires 'Browse' permission on the specified project.
func (s *IssuesService) List(opt *IssuesListOption) (v *IssuesList, resp *http.Response, err error) {
	err = s.ValidateListOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/list", opt)
	if err != nil {
		return
	}

	v = new(IssuesList)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Pull fetches all issues for a given branch.
// The issues returned are not paginated, so the response size can be big.
// Requires 'Browse' permission on the project.
func (s *IssuesService) Pull(opt *IssuesPullOption) (v []byte, resp *http.Response, err error) {
	err = s.ValidatePullOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/pull", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, &v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// PullTaint fetches all taint vulnerabilities for a given branch.
// The vulnerabilities returned are not paginated, so the response size can be big.
// Requires 'Browse' permission on the project.
func (s *IssuesService) PullTaint(opt *IssuesPullTaintOption) (v []byte, resp *http.Response, err error) {
	err = s.ValidatePullTaintOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/pull_taint", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, &v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Reindex triggers reindexing of issues for a project.
// Requires 'Administer System' permission.
func (s *IssuesService) Reindex(opt *IssuesReindexOption) (resp *http.Response, err error) {
	err = s.ValidateReindexOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/reindex", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Search searches for issues.
// Requires 'Browse' permission on the specified project(s).
// For applications, it also requires 'Browse' permission on child projects.
// Returns 503 when issue indexing is in progress.
func (s *IssuesService) Search(opt *IssuesSearchOption) (v *IssuesSearch, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/search", opt)
	if err != nil {
		return
	}

	v = new(IssuesSearch)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SetSeverity changes the severity of an issue.
// Requires authentication, 'Browse' and 'Administer Issues' permissions on the project.
func (s *IssuesService) SetSeverity(opt *IssuesSetSeverityOption) (v *IssuesSetSeverity, resp *http.Response, err error) {
	err = s.ValidateSetSeverityOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/set_severity", opt)
	if err != nil {
		return
	}

	v = new(IssuesSetSeverity)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SetTags sets tags on an issue.
// Requires authentication and 'Browse' permission on the project.
func (s *IssuesService) SetTags(opt *IssuesSetTagsOption) (v *IssuesSetTags, resp *http.Response, err error) {
	err = s.ValidateSetTagsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/set_tags", opt)
	if err != nil {
		return
	}

	v = new(IssuesSetTags)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SetType changes the type of an issue.
// Requires authentication, 'Browse' and 'Administer Issues' permissions on the project.
func (s *IssuesService) SetType(opt *IssuesSetTypeOption) (v *IssuesSetType, resp *http.Response, err error) {
	err = s.ValidateSetTypeOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodPost, "issues/set_type", opt)
	if err != nil {
		return
	}

	v = new(IssuesSetType)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Tags lists tags matching a given query.
func (s *IssuesService) Tags(opt *IssuesTagsOption) (v *IssuesTags, resp *http.Response, err error) {
	err = s.ValidateTagsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest(http.MethodGet, "issues/tags", opt)
	if err != nil {
		return
	}

	v = new(IssuesTags)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// =============================================================================
// Validation Methods
// =============================================================================

// ValidateAddCommentOpt validates the options for adding a comment.
func (s *IssuesService) ValidateAddCommentOpt(opt *IssuesAddCommentOption) error {
	if opt == nil {
		return NewValidationError("IssuesAddCommentOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Text, "Text")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAnticipatedTransitionsOpt validates the options for anticipated transitions.
func (s *IssuesService) ValidateAnticipatedTransitionsOpt(opt *IssuesAnticipatedTransitionsOption) error {
	if opt == nil {
		return NewValidationError("IssuesAnticipatedTransitionsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAssignOpt validates the options for assigning an issue.
func (s *IssuesService) ValidateAssignOpt(opt *IssuesAssignOption) error {
	if opt == nil {
		return NewValidationError("IssuesAssignOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAuthorsOpt validates the options for listing authors.
func (s *IssuesService) ValidateAuthorsOpt(opt *IssuesAuthorsOption) error {
	if opt == nil {
		return nil
	}

	// PageSize must be between 1 and 100 for the authors endpoint
	if opt.PageSize != 0 {
		const maxAuthorsPageSize = 100

		err := ValidateRange(opt.PageSize, MinPageSize, maxAuthorsPageSize, "PageSize")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateBulkChangeOpt validates the options for bulk change.
func (s *IssuesService) ValidateBulkChangeOpt(opt *IssuesBulkChangeOption) error {
	if opt == nil {
		return NewValidationError("IssuesBulkChangeOption", "cannot be nil", ErrMissingRequired)
	}

	if len(opt.Issues) == 0 {
		return NewValidationError("Issues", "is required", ErrMissingRequired)
	}

	// Validate severity if set
	if opt.SetSeverity != "" {
		err := IsValueAuthorized(opt.SetSeverity, allowedSeverities, "SetSeverity")
		if err != nil {
			return err
		}
	}

	// Validate type if set
	if opt.SetType != "" {
		err := IsValueAuthorized(opt.SetType, allowedIssueTypes, "SetType")
		if err != nil {
			return err
		}
	}

	// Validate transition if set
	if opt.DoTransition != "" {
		err := IsValueAuthorized(opt.DoTransition, allowedIssueTransitions, "DoTransition")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateChangelogOpt validates the options for changelog.
func (s *IssuesService) ValidateChangelogOpt(opt *IssuesChangelogOption) error {
	if opt == nil {
		return NewValidationError("IssuesChangelogOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	return nil
}

// ValidateComponentTagsOpt validates the options for component tags.
func (s *IssuesService) ValidateComponentTagsOpt(opt *IssuesComponentTagsOption) error {
	if opt == nil {
		return NewValidationError("IssuesComponentTagsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ComponentUuid, "ComponentUuid")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeleteCommentOpt validates the options for deleting a comment.
func (s *IssuesService) ValidateDeleteCommentOpt(opt *IssuesDeleteCommentOption) error {
	if opt == nil {
		return NewValidationError("IssuesDeleteCommentOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Comment, "Comment")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDoTransitionOpt validates the options for doing a transition.
func (s *IssuesService) ValidateDoTransitionOpt(opt *IssuesDoTransitionOption) error {
	if opt == nil {
		return NewValidationError("IssuesDoTransitionOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Transition, "Transition")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Transition, allowedIssueTransitions, "Transition")
	if err != nil {
		return err
	}

	return nil
}

// ValidateEditCommentOpt validates the options for editing a comment.
func (s *IssuesService) ValidateEditCommentOpt(opt *IssuesEditCommentOption) error {
	if opt == nil {
		return NewValidationError("IssuesEditCommentOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Comment, "Comment")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Text, "Text")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for listing issues.
func (s *IssuesService) ValidateListOpt(opt *IssuesListOption) error {
	if opt == nil {
		return NewValidationError("IssuesListOption", "cannot be nil", ErrMissingRequired)
	}

	// Either project or component is required
	if opt.Project == "" && opt.Component == "" {
		return NewValidationError("Project", "either Project or Component is required", ErrMissingRequired)
	}

	// Validate pagination
	err := opt.Validate()
	if err != nil {
		return err
	}

	// Validate types if set
	if len(opt.Types) > 0 {
		err := AreValuesAuthorized(opt.Types, allowedIssueTypes, "Types")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidatePullOpt validates the options for pulling issues.
func (s *IssuesService) ValidatePullOpt(opt *IssuesPullOption) error {
	if opt == nil {
		return NewValidationError("IssuesPullOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	// Validate languages if set
	if len(opt.Languages) > 0 {
		err := AreValuesAuthorized(opt.Languages, allowedLanguages, "Languages")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidatePullTaintOpt validates the options for pulling taint vulnerabilities.
func (s *IssuesService) ValidatePullTaintOpt(opt *IssuesPullTaintOption) error {
	if opt == nil {
		return NewValidationError("IssuesPullTaintOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	// Validate languages if set
	if len(opt.Languages) > 0 {
		err := AreValuesAuthorized(opt.Languages, allowedLanguages, "Languages")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateReindexOpt validates the options for reindexing.
func (s *IssuesService) ValidateReindexOpt(opt *IssuesReindexOption) error {
	if opt == nil {
		return NewValidationError("IssuesReindexOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for searching issues.
//
//nolint:cyclop,funlen,gocognit,gocyclo // Validation functions are naturally complex due to multiple checks
func (s *IssuesService) ValidateSearchOpt(opt *IssuesSearchOption) error {
	if opt == nil {
		return nil
	}

	// Validate pagination
	err := opt.Validate()
	if err != nil {
		return err
	}

	// Validate impact severities
	if len(opt.ImpactSeverities) > 0 {
		err := AreValuesAuthorized(opt.ImpactSeverities, allowedImpactSeverities, "ImpactSeverities")
		if err != nil {
			return err
		}
	}

	// Validate impact software qualities
	if len(opt.ImpactSoftwareQualities) > 0 {
		err := AreValuesAuthorized(opt.ImpactSoftwareQualities, allowedImpactSoftwareQualities, "ImpactSoftwareQualities")
		if err != nil {
			return err
		}
	}

	// Validate clean code attribute categories
	if len(opt.CleanCodeAttributeCategories) > 0 {
		err := AreValuesAuthorized(opt.CleanCodeAttributeCategories, allowedCleanCodeAttributesCategories, "CleanCodeAttributeCategories")
		if err != nil {
			return err
		}
	}

	// Validate severities
	if len(opt.Severities) > 0 {
		err := AreValuesAuthorized(opt.Severities, allowedSeverities, "Severities")
		if err != nil {
			return err
		}
	}

	// Validate types
	if len(opt.Types) > 0 {
		err := AreValuesAuthorized(opt.Types, allowedIssueTypes, "Types")
		if err != nil {
			return err
		}
	}

	// Validate statuses
	if len(opt.Statuses) > 0 {
		err := AreValuesAuthorized(opt.Statuses, allowedIssueStatuses, "Statuses")
		if err != nil {
			return err
		}
	}

	// Validate issue statuses
	if len(opt.IssueStatuses) > 0 {
		err := AreValuesAuthorized(opt.IssueStatuses, allowedIssueStatuses, "IssueStatuses")
		if err != nil {
			return err
		}
	}

	// Validate resolutions
	if len(opt.Resolutions) > 0 {
		err := AreValuesAuthorized(opt.Resolutions, allowedIssueResolutions, "Resolutions")
		if err != nil {
			return err
		}
	}

	// Validate scopes
	if len(opt.Scopes) > 0 {
		err := AreValuesAuthorized(opt.Scopes, allowedIssueScopes, "Scopes")
		if err != nil {
			return err
		}
	}

	// Validate languages
	if len(opt.Languages) > 0 {
		err := AreValuesAuthorized(opt.Languages, allowedLanguages, "Languages")
		if err != nil {
			return err
		}
	}

	// Validate OWASP Top 10 categories
	if len(opt.OwaspTop10) > 0 {
		err := AreValuesAuthorized(opt.OwaspTop10, allowedOwaspCategories, "OwaspTop10")
		if err != nil {
			return err
		}
	}

	// Validate OWASP Top 10 2021 categories
	if len(opt.OwaspTop102021) > 0 {
		err := AreValuesAuthorized(opt.OwaspTop102021, allowedOwaspCategories, "OwaspTop102021")
		if err != nil {
			return err
		}
	}

	// Validate OWASP Mobile Top 10 2024 categories
	if len(opt.OwaspMobileTop102024) > 0 {
		err := AreValuesAuthorized(opt.OwaspMobileTop102024, allowedOwaspMobileCategories, "OwaspMobileTop102024")
		if err != nil {
			return err
		}
	}

	// Validate SANS Top 25 categories
	if len(opt.SansTop25) > 0 {
		err := AreValuesAuthorized(opt.SansTop25, allowedSansTop25Categories, "SansTop25")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSetSeverityOpt validates the options for setting severity.
func (s *IssuesService) ValidateSetSeverityOpt(opt *IssuesSetSeverityOption) error {
	if opt == nil {
		return NewValidationError("IssuesSetSeverityOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	// Validate severity if set
	if opt.Severity != "" {
		err := IsValueAuthorized(opt.Severity, allowedSeverities, "Severity")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSetTagsOpt validates the options for setting tags.
func (s *IssuesService) ValidateSetTagsOpt(opt *IssuesSetTagsOption) error {
	if opt == nil {
		return NewValidationError("IssuesSetTagsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetTypeOpt validates the options for setting type.
func (s *IssuesService) ValidateSetTypeOpt(opt *IssuesSetTypeOption) error {
	if opt == nil {
		return NewValidationError("IssuesSetTypeOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Issue, "Issue")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Type, allowedIssueTypes, "Type")
	if err != nil {
		return err
	}

	return nil
}

// ValidateTagsOpt validates the options for listing tags.
func (s *IssuesService) ValidateTagsOpt(opt *IssuesTagsOption) error {
	if opt == nil {
		return nil
	}

	// Validate pagination
	if opt.PageSize != 0 {
		err := ValidateRange(opt.PageSize, MinPageSize, MaxPageSize, "PageSize")
		if err != nil {
			return err
		}
	}

	return nil
}
