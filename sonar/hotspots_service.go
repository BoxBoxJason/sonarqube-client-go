package sonar

import "net/http"

const (
	// MaxHotspotCommentLength is the maximum length for a hotspot comment.
	MaxHotspotCommentLength = 1000
	// MaxHotspotListPageSize is the maximum page size for the List endpoint.
	MaxHotspotListPageSize = 500
)

// HotspotsService handles communication with the Security Hotspots related methods
// of the SonarQube API.
// Security Hotspots are security-sensitive pieces of code that require a review
// to assess whether they are vulnerabilities.
type HotspotsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// =============================================================================
// Allowed Values
// =============================================================================

//nolint:gochecknoglobals // constant set of allowed values
var (
	// allowedHotspotStatuses is the set of allowed hotspot statuses.
	allowedHotspotStatuses = map[string]struct{}{
		"TO_REVIEW": {},
		"REVIEWED":  {},
	}

	// allowedHotspotResolutions is the set of allowed hotspot resolutions.
	allowedHotspotResolutions = map[string]struct{}{
		"FIXED":        {},
		"SAFE":         {},
		"ACKNOWLEDGED": {},
	}

	// allowedOwaspAsvsLevels is the set of allowed OWASP ASVS levels.
	allowedOwaspAsvsLevels = map[string]struct{}{
		"1": {},
		"2": {},
		"3": {},
	}
)

// =============================================================================
// Shared Types
// =============================================================================

// HotspotComponent represents a component in a hotspots response.
type HotspotComponent struct {
	// Key is the unique identifier of the component.
	Key string `json:"key,omitempty"`
	// LongName is the long name of the component.
	LongName string `json:"longName,omitempty"`
	// Name is the short name of the component.
	Name string `json:"name,omitempty"`
	// Path is the path to the component file.
	Path string `json:"path,omitempty"`
	// Qualifier is the component qualifier (FIL, DIR, etc.).
	Qualifier string `json:"qualifier,omitempty"`
}

// HotspotPaging represents pagination information in a hotspots response.
type HotspotPaging struct {
	// PageIndex is the current page index (1-based).
	PageIndex int64 `json:"pageIndex,omitempty"`
	// PageSize is the number of items per page.
	PageSize int64 `json:"pageSize,omitempty"`
	// Total is the total number of items.
	Total int64 `json:"total,omitempty"`
}

// HotspotSummary represents a hotspot summary in list/search results.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HotspotSummary struct {
	// Assignee is the login of the user assigned to this hotspot.
	Assignee string `json:"assignee,omitempty"`
	// Author is the SCM author of the code where the hotspot was found.
	Author string `json:"author,omitempty"`
	// Component is the key of the component where the hotspot was found.
	Component string `json:"component,omitempty"`
	// CreationDate is the timestamp when the hotspot was created.
	CreationDate string `json:"creationDate,omitempty"`
	// Flows is the list of code flows related to this hotspot.
	Flows []any `json:"flows,omitempty"`
	// Key is the unique identifier of the hotspot.
	Key string `json:"key,omitempty"`
	// Line is the line number where the hotspot was found.
	Line int64 `json:"line,omitempty"`
	// Message is the main message describing the hotspot.
	Message string `json:"message,omitempty"`
	// MessageFormattings is the list of message formatting rules.
	MessageFormattings []any `json:"messageFormattings,omitempty"`
	// Project is the key of the project containing the hotspot.
	Project string `json:"project,omitempty"`
	// RuleKey is the key of the rule that raised this hotspot.
	RuleKey string `json:"ruleKey,omitempty"`
	// SecurityCategory is the security category of the hotspot.
	SecurityCategory string `json:"securityCategory,omitempty"`
	// Status is the current status of the hotspot.
	Status string `json:"status,omitempty"`
	// UpdateDate is the timestamp when the hotspot was last updated.
	UpdateDate string `json:"updateDate,omitempty"`
	// VulnerabilityProbability is the probability that this hotspot is a vulnerability.
	VulnerabilityProbability string `json:"vulnerabilityProbability,omitempty"`
}

// HotspotComment represents a comment on a hotspot.
type HotspotComment struct {
	// CreatedAt is the timestamp when the comment was created.
	CreatedAt string `json:"createdAt,omitempty"`
	// HTMLText is the HTML-formatted text of the comment.
	HTMLText string `json:"htmlText,omitempty"`
	// Key is the unique identifier of the comment.
	Key string `json:"key,omitempty"`
	// Login is the login of the user who created the comment.
	Login string `json:"login,omitempty"`
	// Markdown is the Markdown-formatted text of the comment.
	Markdown string `json:"markdown,omitempty"`
	// Updatable indicates if the comment can be updated by the current user.
	Updatable bool `json:"updatable,omitempty"`
}

// HotspotUser represents a user in a hotspots response.
type HotspotUser struct {
	// Login is the login of the user.
	Login string `json:"login,omitempty"`
	// Name is the display name of the user.
	Name string `json:"name,omitempty"`
	// Active indicates if the user is active.
	Active bool `json:"active,omitempty"`
}

// HotspotChangelogEntry represents a changelog entry for a hotspot.
//
//nolint:govet // Field order maintained for JSON serialization consistency
type HotspotChangelogEntry struct {
	// Diffs is the list of changes made.
	Diffs []HotspotDiff `json:"diffs,omitempty"`
	// Avatar is the URL to the user's avatar.
	Avatar string `json:"avatar,omitempty"`
	// CreationDate is the timestamp when the change was made.
	CreationDate string `json:"creationDate,omitempty"`
	// User is the login of the user who made the change.
	User string `json:"user,omitempty"`
	// UserName is the display name of the user who made the change.
	UserName string `json:"userName,omitempty"`
	// IsUserActive indicates if the user is active.
	IsUserActive bool `json:"isUserActive,omitempty"`
}

// HotspotDiff represents a single change in a changelog entry.
type HotspotDiff struct {
	// Key is the identifier of the changed field.
	Key string `json:"key,omitempty"`
	// NewValue is the new value after the change.
	NewValue string `json:"newValue,omitempty"`
	// OldValue is the old value before the change.
	OldValue string `json:"oldValue,omitempty"`
}

// HotspotMessageFormatting represents a formatting rule for messages.
type HotspotMessageFormatting struct {
	// Type is the type of formatting.
	Type string `json:"type,omitempty"`
	// End is the ending position of the formatting.
	End int64 `json:"end,omitempty"`
	// Start is the starting position of the formatting.
	Start int64 `json:"start,omitempty"`
}

// HotspotProject represents a project in a hotspot show response.
type HotspotProject struct {
	// Key is the unique identifier of the project.
	Key string `json:"key,omitempty"`
	// LongName is the long name of the project.
	LongName string `json:"longName,omitempty"`
	// Name is the short name of the project.
	Name string `json:"name,omitempty"`
	// Qualifier is the project qualifier.
	Qualifier string `json:"qualifier,omitempty"`
}

// HotspotRule represents a rule in a hotspot show response.
type HotspotRule struct {
	// Key is the rule key.
	Key string `json:"key,omitempty"`
	// Name is the rule name.
	Name string `json:"name,omitempty"`
	// SecurityCategory is the security category of the rule.
	SecurityCategory string `json:"securityCategory,omitempty"`
	// VulnerabilityProbability is the probability of vulnerability.
	VulnerabilityProbability string `json:"vulnerabilityProbability,omitempty"`
}

// =============================================================================
// Response Types
// =============================================================================

// HotspotsEditComment represents the response from editing a hotspot comment.
type HotspotsEditComment struct {
	// CreatedAt is the timestamp when the comment was created.
	CreatedAt string `json:"createdAt,omitempty"`
	// HTMLText is the HTML-formatted text of the comment.
	HTMLText string `json:"htmlText,omitempty"`
	// Key is the unique identifier of the comment.
	Key string `json:"key,omitempty"`
	// Login is the login of the user who created the comment.
	Login string `json:"login,omitempty"`
	// Markdown is the Markdown-formatted text of the comment.
	Markdown string `json:"markdown,omitempty"`
	// Updatable indicates if the comment can be updated by the current user.
	Updatable bool `json:"updatable,omitempty"`
}

// HotspotsList represents the response from listing hotspots.
type HotspotsList struct {
	// Components is the list of components referenced by the hotspots.
	Components []HotspotComponent `json:"components,omitempty"`
	// Hotspots is the list of hotspots.
	Hotspots []HotspotSummary `json:"hotspots,omitempty"`
	// Paging contains pagination information.
	Paging HotspotPaging `json:"paging"`
}

// HotspotsSearch represents the response from searching hotspots.
type HotspotsSearch struct {
	// Components is the list of components referenced by the hotspots.
	Components []HotspotComponent `json:"components,omitempty"`
	// Hotspots is the list of hotspots.
	Hotspots []HotspotSummary `json:"hotspots,omitempty"`
	// Paging contains pagination information.
	Paging HotspotPaging `json:"paging"`
}

// HotspotsShow represents the response from showing hotspot details.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HotspotsShow struct {
	// Assignee is the login of the user assigned to this hotspot.
	Assignee string `json:"assignee,omitempty"`
	// Author is the SCM author of the code where the hotspot was found.
	Author string `json:"author,omitempty"`
	// CanChangeStatus indicates if the current user can change the hotspot status.
	CanChangeStatus bool `json:"canChangeStatus,omitempty"`
	// Changelog is the list of changes made to this hotspot.
	Changelog []HotspotChangelogEntry `json:"changelog,omitempty"`
	// CodeVariants is the list of code variants affected by this hotspot.
	CodeVariants []string `json:"codeVariants,omitempty"`
	// Comment is the list of comments on this hotspot.
	Comment []HotspotComment `json:"comment,omitempty"`
	// Component is the component where the hotspot was found.
	Component HotspotComponent `json:"component"`
	// CreationDate is the timestamp when the hotspot was created.
	CreationDate string `json:"creationDate,omitempty"`
	// Hash is the hash of the line content for tracking purposes.
	Hash string `json:"hash,omitempty"`
	// Key is the unique identifier of the hotspot.
	Key string `json:"key,omitempty"`
	// Line is the line number where the hotspot was found.
	Line int64 `json:"line,omitempty"`
	// Message is the main message describing the hotspot.
	Message string `json:"message,omitempty"`
	// MessageFormattings is the list of message formatting rules.
	MessageFormattings []HotspotMessageFormatting `json:"messageFormattings,omitempty"`
	// Project is the project containing the hotspot.
	Project HotspotProject `json:"project"`
	// Rule is the rule that raised this hotspot.
	Rule HotspotRule `json:"rule"`
	// Status is the current status of the hotspot.
	Status string `json:"status,omitempty"`
	// UpdateDate is the timestamp when the hotspot was last updated.
	UpdateDate string `json:"updateDate,omitempty"`
	// Users is the list of users referenced by this hotspot.
	Users []HotspotUser `json:"users,omitempty"`
}

// =============================================================================
// Option Types
// =============================================================================

// HotspotsAddCommentOption contains parameters for the AddComment method.
type HotspotsAddCommentOption struct {
	// Comment is the comment text.
	// This field is required. Maximum length: 1000 characters.
	Comment string `url:"comment"`
	// Hotspot is the key of the Security Hotspot.
	// This field is required.
	Hotspot string `url:"hotspot"`
}

// HotspotsAssignOption contains parameters for the Assign method.
type HotspotsAssignOption struct {
	// Assignee is the login of the assignee with 'Browse' project permission.
	// This field is optional (since 8.9).
	Assignee string `url:"assignee,omitempty"`
	// Comment is a comment provided with the assign action.
	// This field is optional.
	Comment string `url:"comment,omitempty"`
	// Hotspot is the key of the Security Hotspot.
	// This field is required.
	Hotspot string `url:"hotspot"`
}

// HotspotsChangeStatusOption contains parameters for the ChangeStatus method.
type HotspotsChangeStatusOption struct {
	// Comment is optional comment text.
	// This field is optional.
	Comment string `url:"comment,omitempty"`
	// Hotspot is the key of the Security Hotspot.
	// This field is required.
	Hotspot string `url:"hotspot"`
	// Resolution is the resolution when new status is REVIEWED.
	// Allowed values: FIXED, SAFE, ACKNOWLEDGED.
	// This field is optional.
	Resolution string `url:"resolution,omitempty"`
	// Status is the new status of the Security Hotspot.
	// Allowed values: TO_REVIEW, REVIEWED.
	// This field is required.
	Status string `url:"status"`
}

// HotspotsDeleteCommentOption contains parameters for the DeleteComment method.
type HotspotsDeleteCommentOption struct {
	// Comment is the key of the comment to delete.
	// This field is required.
	Comment string `url:"comment"`
}

// HotspotsEditCommentOption contains parameters for the EditComment method.
type HotspotsEditCommentOption struct {
	// Comment is the key of the comment to edit.
	// This field is required.
	Comment string `url:"comment"`
	// Text is the new comment text.
	// This field is required. Maximum length: 1000 characters.
	Text string `url:"text"`
}

// HotspotsListOption contains parameters for the List method.
//
//nolint:govet // Field order maintained for API parameter consistency
type HotspotsListOption struct {
	// PaginationArgs contains the pagination parameters.
	PaginationArgs `url:",inline"`

	// Branch is the branch key. Not available in the community edition.
	// This field is optional.
	Branch string `url:"branch,omitempty"`
	// Project is the key of the project.
	// This field is required.
	Project string `url:"project"`
	// PullRequest is the pull request id. Not available in the community edition.
	// This field is optional.
	PullRequest string `url:"pullRequest,omitempty"`
	// Resolution filters by resolution when status is REVIEWED.
	// Allowed values: FIXED, SAFE, ACKNOWLEDGED.
	// This field is optional.
	Resolution string `url:"resolution,omitempty"`
	// Status filters by hotspot status.
	// Allowed values: TO_REVIEW, REVIEWED.
	// This field is optional.
	Status string `url:"status,omitempty"`
	// InNewCodePeriod filters to only Security Hotspots created in the new code period.
	// This field is optional. Default: false.
	InNewCodePeriod bool `url:"inNewCodePeriod,omitempty"`
}

// HotspotsPullOption contains parameters for the Pull method.
//
//nolint:govet // Field order maintained for API parameter consistency
type HotspotsPullOption struct {
	// Languages is a comma-separated list of languages.
	// If not present, all hotspots regardless of their language are returned.
	// This field is optional.
	Languages []string `url:"languages,omitempty,comma"`
	// BranchName is the branch name for which hotspots are fetched.
	// This field is required.
	BranchName string `url:"branchName"`
	// ProjectKey is the project key for which hotspots are fetched.
	// This field is required.
	ProjectKey string `url:"projectKey"`
	// ChangedSince is a timestamp. If present, only hotspots modified after this
	// timestamp are returned (both open and closed).
	// This field is optional.
	ChangedSince int64 `url:"changedSince,omitempty"`
}

// HotspotsSearchOption contains parameters for the Search method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HotspotsSearchOption struct {
	// PaginationArgs contains the pagination parameters.
	PaginationArgs `url:",inline"`

	// Casa is a comma-separated list of CASA categories.
	// This field is optional. Since: 10.7.
	Casa []string `url:"casa,omitempty,comma"`
	// ComplianceStandards is a list of compliance standards.
	// This field is optional. Since: 2025.6.
	ComplianceStandards []string `url:"complianceStandards,omitempty,comma"`
	// Cwe is a comma-separated list of CWE numbers.
	// This field is optional. Since: 8.8.
	Cwe []string `url:"cwe,omitempty,comma"`
	// Files is a comma-separated list of files to filter hotspots.
	// This field is optional. Since: 9.0.
	Files []string `url:"files,omitempty,comma"`
	// Hotspots is a comma-separated list of Security Hotspot keys.
	// This parameter is required unless Project is provided.
	// This field is optional.
	Hotspots []string `url:"hotspots,omitempty,comma"`
	// OwaspAsvs40 is a comma-separated list of OWASP ASVS v4.0 categories or rules.
	// This field is optional. Since: 9.7.
	OwaspAsvs40 []string `url:"owaspAsvs-4.0,omitempty,comma"`
	// OwaspTop10 is a comma-separated list of OWASP 2017 Top 10 lowercase categories.
	// This field is optional. Since: 8.6.
	OwaspTop10 []string `url:"owaspTop10,omitempty,comma"`
	// OwaspTop102021 is a comma-separated list of OWASP 2021 Top 10 lowercase categories.
	// This field is optional. Since: 9.4.
	OwaspTop102021 []string `url:"owaspTop10-2021,omitempty,comma"`
	// PciDss32 is a comma-separated list of PCI DSS v3.2 categories.
	// This field is optional. Since: 9.6.
	PciDss32 []string `url:"pciDss-3.2,omitempty,comma"`
	// PciDss40 is a comma-separated list of PCI DSS v4.0 categories.
	// This field is optional. Since: 9.6.
	PciDss40 []string `url:"pciDss-4.0,omitempty,comma"`
	// SansTop25 is a comma-separated list of SANS Top 25 categories.
	// This field is optional. Since: 8.6. Deprecated since: 10.0.
	SansTop25 []string `url:"sansTop25,omitempty,comma"`
	// SonarsourceSecurity is a comma-separated list of SonarSource security categories.
	// Use 'others' to select issues not associated with any category.
	// This field is optional. Since: 8.6.
	SonarsourceSecurity []string `url:"sonarsourceSecurity,omitempty,comma"`
	// StigASDV5R3 is a comma-separated list of STIG V5R3 lowercase categories.
	// This field is optional. Since: 10.7.
	StigASDV5R3 []string `url:"stig-ASD_V5R3,omitempty,comma"`
	// Branch is the branch key. Not available in the community edition.
	// This field is optional.
	Branch string `url:"branch,omitempty"`
	// OwaspAsvsLevel filters hotspots with lower or equal OWASP ASVS level.
	// Should be used in combination with OwaspAsvs40.
	// Allowed values: 1, 2, 3.
	// This field is optional. Since: 9.7.
	OwaspAsvsLevel string `url:"owaspAsvsLevel,omitempty"`
	// Project is the key of the project or application.
	// This parameter is required unless Hotspots is provided.
	// This field is optional.
	Project string `url:"project,omitempty"`
	// PullRequest is the pull request id. Not available in the community edition.
	// This field is optional.
	PullRequest string `url:"pullRequest,omitempty"`
	// Resolution filters by resolution when status is REVIEWED.
	// Allowed values: FIXED, SAFE, ACKNOWLEDGED.
	// This field is optional.
	Resolution string `url:"resolution,omitempty"`
	// Status filters by hotspot status.
	// Allowed values: TO_REVIEW, REVIEWED.
	// This field is optional.
	Status string `url:"status,omitempty"`
	// InNewCodePeriod filters to only Security Hotspots created in the new code period.
	// This field is optional. Default: false. Since: 9.5.
	InNewCodePeriod bool `url:"inNewCodePeriod,omitempty"`
	// OnlyMine filters to only hotspots assigned to the current user.
	// This field is optional.
	OnlyMine bool `url:"onlyMine,omitempty"`
}

// HotspotsShowOption contains parameters for the Show method.
type HotspotsShowOption struct {
	// Hotspot is the key of the Security Hotspot.
	// This field is required.
	Hotspot string `url:"hotspot"`
}

// =============================================================================
// Validation Functions
// =============================================================================

// ValidateAddCommentOpt validates the options for the AddComment method.
func (s *HotspotsService) ValidateAddCommentOpt(opt *HotspotsAddCommentOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Comment, "Comment")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Comment, MaxHotspotCommentLength, "Comment")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Hotspot, "Hotspot")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAssignOpt validates the options for the Assign method.
func (s *HotspotsService) ValidateAssignOpt(opt *HotspotsAssignOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Hotspot, "Hotspot")
	if err != nil {
		return err
	}

	return nil
}

// ValidateChangeStatusOpt validates the options for the ChangeStatus method.
func (s *HotspotsService) ValidateChangeStatusOpt(opt *HotspotsChangeStatusOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Hotspot, "Hotspot")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Status, "Status")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Status, allowedHotspotStatuses, "Status")
	if err != nil {
		return err
	}

	if opt.Resolution != "" {
		err = IsValueAuthorized(opt.Resolution, allowedHotspotResolutions, "Resolution")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteCommentOpt validates the options for the DeleteComment method.
func (s *HotspotsService) ValidateDeleteCommentOpt(opt *HotspotsDeleteCommentOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Comment, "Comment")
	if err != nil {
		return err
	}

	return nil
}

// ValidateEditCommentOpt validates the options for the EditComment method.
func (s *HotspotsService) ValidateEditCommentOpt(opt *HotspotsEditCommentOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Comment, "Comment")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Text, "Text")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Text, MaxHotspotCommentLength, "Text")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for the List method.
func (s *HotspotsService) ValidateListOpt(opt *HotspotsListOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	if opt.Status != "" {
		err = IsValueAuthorized(opt.Status, allowedHotspotStatuses, "Status")
		if err != nil {
			return err
		}
	}

	if opt.Resolution != "" {
		err = IsValueAuthorized(opt.Resolution, allowedHotspotResolutions, "Resolution")
		if err != nil {
			return err
		}
	}

	// Validate page size for List endpoint (max 500)
	if opt.PageSize > MaxHotspotListPageSize {
		return NewValidationError("PageSize", "must be less than or equal to 500", ErrInvalidValue)
	}

	err = opt.Validate()
	if err != nil {
		return err
	}

	return nil
}

// ValidatePullOpt validates the options for the Pull method.
func (s *HotspotsService) ValidatePullOpt(opt *HotspotsPullOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.BranchName, "BranchName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
//
//nolint:cyclop // Validation functions are naturally complex due to multiple checks
func (s *HotspotsService) ValidateSearchOpt(opt *HotspotsSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Either project or hotspots must be provided
	if opt.Project == "" && len(opt.Hotspots) == 0 {
		return NewValidationError("Project", "either project or hotspots is required", ErrMissingRequired)
	}

	if opt.Status != "" {
		err := IsValueAuthorized(opt.Status, allowedHotspotStatuses, "Status")
		if err != nil {
			return err
		}
	}

	if opt.Resolution != "" {
		err := IsValueAuthorized(opt.Resolution, allowedHotspotResolutions, "Resolution")
		if err != nil {
			return err
		}
	}

	if opt.OwaspAsvsLevel != "" {
		err := IsValueAuthorized(opt.OwaspAsvsLevel, allowedOwaspAsvsLevels, "OwaspAsvsLevel")
		if err != nil {
			return err
		}
	}

	if len(opt.OwaspTop10) > 0 {
		err := AreValuesAuthorized(opt.OwaspTop10, allowedOwaspCategories, "OwaspTop10")
		if err != nil {
			return err
		}
	}

	if len(opt.OwaspTop102021) > 0 {
		err := AreValuesAuthorized(opt.OwaspTop102021, allowedOwaspCategories, "OwaspTop102021")
		if err != nil {
			return err
		}
	}

	if len(opt.SansTop25) > 0 {
		err := AreValuesAuthorized(opt.SansTop25, allowedSansTop25Categories, "SansTop25")
		if err != nil {
			return err
		}
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	return nil
}

// ValidateShowOpt validates the options for the Show method.
func (s *HotspotsService) ValidateShowOpt(opt *HotspotsShowOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Hotspot, "Hotspot")
	if err != nil {
		return err
	}

	return nil
}

// =============================================================================
// Service Methods
// =============================================================================

// AddComment adds a comment to a Security Hotspot.
// Requires authentication and the 'Browse' permission on the project.
//
// API endpoint: POST /api/hotspots/add_comment.
// Since: 8.1.
// Internal: true.
func (s *HotspotsService) AddComment(opt *HotspotsAddCommentOption) (*http.Response, error) {
	err := s.ValidateAddCommentOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "hotspots/add_comment", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Assign assigns a hotspot to an active user.
// Requires authentication and 'Browse' permission on the project.
//
// API endpoint: POST /api/hotspots/assign.
// Since: 8.2.
// Internal: true.
func (s *HotspotsService) Assign(opt *HotspotsAssignOption) (*http.Response, error) {
	err := s.ValidateAssignOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "hotspots/assign", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ChangeStatus changes the status of a Security Hotspot.
// Requires the 'Administer Security Hotspot' permission.
//
// API endpoint: POST /api/hotspots/change_status.
// Since: 8.1.
func (s *HotspotsService) ChangeStatus(opt *HotspotsChangeStatusOption) (*http.Response, error) {
	err := s.ValidateChangeStatusOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "hotspots/change_status", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteComment deletes a comment from a Security Hotspot.
// Requires authentication and the 'Browse' permission on the project.
//
// API endpoint: POST /api/hotspots/delete_comment.
// Since: 8.2.
// Internal: true.
func (s *HotspotsService) DeleteComment(opt *HotspotsDeleteCommentOption) (*http.Response, error) {
	err := s.ValidateDeleteCommentOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "hotspots/delete_comment", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// EditComment edits a comment on a Security Hotspot.
// Requires authentication and the 'Browse' permission on the project.
//
// API endpoint: POST /api/hotspots/edit_comment.
// Since: 8.2.
// Internal: true.
func (s *HotspotsService) EditComment(opt *HotspotsEditCommentOption) (*HotspotsEditComment, *http.Response, error) {
	err := s.ValidateEditCommentOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "hotspots/edit_comment", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(HotspotsEditComment)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// List lists Security Hotspots. This endpoint is used in degraded mode when
// issue indexing is running.
// Requires the 'Browse' permission on the specified project.
//
// Note: Total number of Security Hotspots will always equal the page size,
// as counting all issues is not supported.
//
// API endpoint: GET /api/hotspots/list.
// Since: 10.2.
// Internal: true.
func (s *HotspotsService) List(opt *HotspotsListOption) (*HotspotsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "hotspots/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(HotspotsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Pull fetches and returns all (unless filtered) hotspots for a given branch.
// The hotspots returned are not paginated, so the response size can be big.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/hotspots/pull.
// Since: 10.1.
// Internal: true.
func (s *HotspotsService) Pull(opt *HotspotsPullOption) ([]byte, *http.Response, error) {
	err := s.ValidatePullOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "hotspots/pull", opt)
	if err != nil {
		return nil, nil, err
	}

	var result []byte

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Search searches for Security Hotspots.
// Requires the 'Browse' permission on the specified project(s).
// For applications, also requires 'Browse' permission on child projects.
// Returns 503 Service Unavailable when issue indexing is in progress.
//
// API endpoint: GET /api/hotspots/search.
// Since: 8.1.
func (s *HotspotsService) Search(opt *HotspotsSearchOption) (*HotspotsSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "hotspots/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(HotspotsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Show provides the details of a Security Hotspot.
//
// API endpoint: GET /api/hotspots/show.
// Since: 8.1.
func (s *HotspotsService) Show(opt *HotspotsShowOption) (*HotspotsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "hotspots/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(HotspotsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
