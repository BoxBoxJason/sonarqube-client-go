package sonar

import (
	"net/http"
)

// SourcesService handles communication with the Sources related methods of the SonarQube API.
// Get source code.
//
// Since: 4.4.
type SourcesService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SourcesIndex represents the response from getting source file lines.
// The key is the line number as a string, and the value is the source code.
type SourcesIndex struct {
	Sources map[string]string `json:"sources,omitempty"`
}

// SourcesIssueSnippets represents the response from getting issue snippets.
type SourcesIssueSnippets map[string]SourcesIssueSnippet

// SourcesIssueSnippet represents a single issue snippet.
type SourcesIssueSnippet struct {
	// Component is the component information.
	Component SourcesComponent `json:"component,omitzero"`
	// Sources is the source lines for the issue snippet.
	Sources []SourcesLine `json:"sources,omitempty"`
}

// SourcesComponent represents source component information.
type SourcesComponent struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// LongName is the long name of the component.
	LongName string `json:"longName,omitempty"`
	// Name is the short name of the component.
	Name string `json:"name,omitempty"`
	// Path is the path to the component.
	Path string `json:"path,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// UUID is the component UUID.
	UUID string `json:"uuid,omitempty"`
}

// SourcesLines represents the response from getting source file lines.
type SourcesLines struct {
	// Sources is the list of source lines.
	Sources []SourcesLine `json:"sources,omitempty"`
}

// SourcesLine represents a single source line.
//
//nolint:govet // Field alignment is less important than logical grouping
type SourcesLine struct {
	// Line is the line number.
	Line int64 `json:"line,omitempty"`
	// Code is the HTML-formatted source code.
	Code string `json:"code,omitempty"`
	// SCMAuthor is the author who last modified this line.
	SCMAuthor string `json:"scmAuthor,omitempty"`
	// SCMDate is the date when the line was last modified.
	SCMDate string `json:"scmDate,omitempty"`
	// SCMRevision is the SCM revision for this line.
	SCMRevision string `json:"scmRevision,omitempty"`
	// LineHits is the number of times this line was executed by tests.
	LineHits int64 `json:"lineHits,omitempty"`
	// Conditions is the number of conditions on this line.
	Conditions int64 `json:"conditions,omitempty"`
	// CoveredConditions is the number of covered conditions on this line.
	CoveredConditions int64 `json:"coveredConditions,omitempty"`
	// Duplicated indicates if the line is duplicated.
	Duplicated bool `json:"duplicated,omitempty"`
	// IsNew indicates if this is a new line.
	IsNew bool `json:"isNew,omitempty"`
}

// SourcesScm represents the response from getting SCM data.
type SourcesScm struct {
	// Scm is the SCM information by line number.
	// Each entry contains: [line number, author, commit date, revision]
	Scm [][]any `json:"scm,omitempty"`
}

// SourcesShow represents the response from getting raw source file content.
type SourcesShow struct {
	// Sources is the list of source lines (index and code pairs).
	Sources [][]any `json:"sources,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SourcesIndexOption represents options for getting source file lines.
type SourcesIndexOption struct {
	// Resource is the file key (API parameter: "resource", required).
	Resource string `url:"resource,omitempty"`
	// From is the starting line number (optional, default: 1).
	From int64 `url:"from,omitempty"`
	// To is the ending line number (optional, default: end of file).
	To int64 `url:"to,omitempty"`
}

// SourcesIssueSnippetsOption represents options for getting issue snippets.
type SourcesIssueSnippetsOption struct {
	// IssueKey is the issue key (required).
	IssueKey string `url:"issueKey,omitempty"`
}

// SourcesLinesOption represents options for getting source file lines.
type SourcesLinesOption struct {
	// Key is the file key (required).
	Key string `url:"key,omitempty"`
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// PullRequest is the pull request identifier (optional).
	PullRequest string `url:"pullRequest,omitempty"`
	// From is the starting line number (optional, default: 1).
	From int64 `url:"from,omitempty"`
	// To is the ending line number (optional, default: end of file).
	To int64 `url:"to,omitempty"`
}

// SourcesRawOption represents options for getting raw source file content.
type SourcesRawOption struct {
	// Key is the file key (required).
	Key string `url:"key,omitempty"`
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// PullRequest is the pull request identifier (optional).
	PullRequest string `url:"pullRequest,omitempty"`
}

// SourcesScmOption represents options for getting SCM data.
type SourcesScmOption struct {
	// Key is the file key (required).
	Key string `url:"key,omitempty"`
	// CommitsByLine indicates whether to group commits by line (optional, default: false).
	CommitsByLine bool `url:"commits_by_line,omitempty"`
	// From is the starting line number (optional, default: 1).
	From int64 `url:"from,omitempty"`
	// To is the ending line number (optional, default: end of file).
	To int64 `url:"to,omitempty"`
}

// SourcesShowOption represents options for showing source file content.
type SourcesShowOption struct {
	// Key is the file key (required).
	Key string `url:"key,omitempty"`
	// From is the starting line number (optional, default: 1).
	From int64 `url:"from,omitempty"`
	// To is the ending line number (optional, default: end of file).
	To int64 `url:"to,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Methods
// -----------------------------------------------------------------------------

// ValidateIndexOpt validates the options for Index.
func (s *SourcesService) ValidateIndexOpt(opt *SourcesIndexOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Resource, "Resource")
}

// ValidateIssueSnippetsOpt validates the options for IssueSnippets.
func (s *SourcesService) ValidateIssueSnippetsOpt(opt *SourcesIssueSnippetsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.IssueKey, "IssueKey")
}

// ValidateLinesOpt validates the options for Lines.
func (s *SourcesService) ValidateLinesOpt(opt *SourcesLinesOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateRawOpt validates the options for Raw.
func (s *SourcesService) ValidateRawOpt(opt *SourcesRawOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateScmOpt validates the options for Scm.
func (s *SourcesService) ValidateScmOpt(opt *SourcesScmOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateShowOpt validates the options for Show.
func (s *SourcesService) ValidateShowOpt(opt *SourcesShowOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Index gets source code.
// Requires 'See Source Code' permission on file's project.
// Each element is a key/value pair where the key is the line number
// and the value is the content of the line.
//
// Since: 5.0.
//
// Deprecated: This web service is deprecated since 5.1. Use api/sources/lines instead.
func (s *SourcesService) Index(opt *SourcesIndexOption) (*SourcesIndex, *http.Response, error) {
	err := s.ValidateIndexOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "sources/index", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SourcesIndex)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IssueSnippets gets source code snippets for issues.
// Requires 'Browse' permission on the file's project.
// Returns source code snippets relevant to the given issue.
//
// Since: 7.8.
func (s *SourcesService) IssueSnippets(opt *SourcesIssueSnippetsOption) (*SourcesIssueSnippets, *http.Response, error) {
	err := s.ValidateIssueSnippetsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "sources/issue_snippets", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SourcesIssueSnippets)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Lines shows source code with line-oriented info.
// Requires 'See Source Code' permission on file's project.
// Returns source code with additional info like SCM data, coverage, and duplications.
//
// Since: 5.0.
func (s *SourcesService) Lines(opt *SourcesLinesOption) (*SourcesLines, *http.Response, error) {
	err := s.ValidateLinesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "sources/lines", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SourcesLines)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Raw gets the raw source code.
// Returns the raw source code as plain text.
// Requires 'See Source Code' permission on file's project.
//
// Since: 5.0.
func (s *SourcesService) Raw(opt *SourcesRawOption) (string, *http.Response, error) {
	err := s.ValidateRawOpt(opt)
	if err != nil {
		return "", nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "sources/raw", opt)
	if err != nil {
		return "", nil, err
	}

	var result string

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return "", resp, err
	}

	return result, resp, nil
}

// Scm gets SCM information of source files.
// Requires 'See Source Code' permission on file's project.
// Returns source code modification information.
//
// Since: 4.4.
func (s *SourcesService) Scm(opt *SourcesScmOption) (*SourcesScm, *http.Response, error) {
	err := s.ValidateScmOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "sources/scm", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SourcesScm)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Show gets source code as line number / text pairs.
// Requires 'See Source Code' permission on file's project.
//
// Since: 4.4.
func (s *SourcesService) Show(opt *SourcesShowOption) (*SourcesShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "sources/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SourcesShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
