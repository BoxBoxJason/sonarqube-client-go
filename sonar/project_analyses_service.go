package sonargo

import (
	"net/http"
	"strings"
)

// Date format constants.
const (
	dateLen        = 10 // Length of YYYY-MM-DD
	dateTimeMinLen = 20 // Minimum length of YYYY-MM-DDTHH:mm:ssZ
	dateParts      = 3  // Number of parts in a date (year, month, day)
	dateTimeParts  = 2  // Number of parts when splitting by T
	yearLen        = 4  // Length of year part
	monthDayLen    = 2  // Length of month and day parts
)

// ProjectAnalysesService handles communication with the project analyses related
// methods of the SonarQube API.
// This service provides management of project analyses, including events.
//
// Since: 6.3.
type ProjectAnalysesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectAnalysesEvent represents a project analysis event.
type ProjectAnalysesEvent struct {
	// Analysis is the analysis key.
	Analysis string `json:"analysis,omitempty"`
	// Category is the event category.
	Category string `json:"category,omitempty"`
	// Description is the event description.
	Description string `json:"description,omitempty"`
	// Key is the event key.
	Key string `json:"key,omitempty"`
	// Name is the event name.
	Name string `json:"name,omitempty"`
	// QualityGate contains quality gate event details.
	QualityGate ProjectAnalysesQualityGate `json:"qualityGate,omitzero"`
}

// ProjectAnalysesQualityGate represents quality gate details in an event.
//
//nolint:govet // fieldalignment - structure kept for readability
type ProjectAnalysesQualityGate struct {
	// Failing lists failing conditions.
	Failing []ProjectAnalysesCondition `json:"failing,omitempty"`
	// Status is the quality gate status.
	Status string `json:"status,omitempty"`
	// StillFailing indicates if still failing.
	StillFailing bool `json:"stillFailing,omitempty"`
}

// ProjectAnalysesCondition represents a quality gate condition.
type ProjectAnalysesCondition struct {
	// Branch is the branch name.
	Branch string `json:"branch,omitempty"`
	// ErrorThreshold is the error threshold.
	ErrorThreshold string `json:"errorThreshold,omitempty"`
	// Metric is the condition metric.
	Metric string `json:"metric,omitempty"`
	// PullRequest is the pull request key.
	PullRequest string `json:"pullRequest,omitempty"`
}

// ProjectAnalysesSearch represents the response from searching analyses.
type ProjectAnalysesSearch struct {
	// Analyses is the list of analyses.
	Analyses []ProjectAnalysis `json:"analyses,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// ProjectAnalysis represents a project analysis.
//
//nolint:govet,tagliatelle // fieldalignment - structure kept for readability, API-defined JSON fields
type ProjectAnalysis struct {
	// BuildString is the build string.
	BuildString string `json:"buildString,omitempty"`
	// Date is the analysis date.
	Date string `json:"date,omitempty"`
	// DetectedCI is the detected CI.
	DetectedCI string `json:"detectedCI,omitempty"`
	// Events is the list of events.
	Events []ProjectAnalysesEvent `json:"events,omitempty"`
	// Key is the analysis key.
	Key string `json:"key,omitempty"`
	// ManualNewCodePeriodBaseline indicates if it's a manual baseline.
	ManualNewCodePeriodBaseline bool `json:"manualNewCodePeriodBaseline,omitempty"`
	// ProjectVersion is the project version.
	ProjectVersion string `json:"projectVersion,omitempty"`
	// Revision is the revision.
	Revision string `json:"revision,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectAnalysesCreateEventOption represents options for creating an event.
type ProjectAnalysesCreateEventOption struct {
	// Analysis is the analysis key (required).
	Analysis string `url:"analysis,omitempty"`
	// Category is the event category.
	// Possible values: VERSION, OTHER.
	// Default: OTHER.
	Category string `url:"category,omitempty"`
	// Name is the event name (required).
	Name string `url:"name,omitempty"`
}

// ProjectAnalysesDeleteOption represents options for deleting an analysis.
type ProjectAnalysesDeleteOption struct {
	// Analysis is the analysis key (required).
	Analysis string `url:"analysis,omitempty"`
}

// ProjectAnalysesDeleteEventOption represents options for deleting an event.
type ProjectAnalysesDeleteEventOption struct {
	// Event is the event key (required).
	Event string `url:"event,omitempty"`
}

// ProjectAnalysesSearchOption represents options for searching analyses.
//
//nolint:govet // fieldalignment - structure kept for readability
type ProjectAnalysesSearchOption struct {
	// PaginationArgs embeds pagination parameters.
	PaginationArgs

	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Category filters events by category.
	// Possible values: VERSION, OTHER, QUALITY_PROFILE, QUALITY_GATE, DEFINITION_CHANGE, SQ_UPGRADE.
	Category string `url:"category,omitempty"`
	// From is the filter by date (inclusive).
	// Format: date or datetime (YYYY-MM-DD or YYYY-MM-DDTHH:mm:ssZ).
	From string `url:"from,omitempty"`
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
	// PullRequest is the pull request key.
	PullRequest string `url:"pullRequest,omitempty"`
	// To is the filter by date (inclusive).
	// Format: date or datetime (YYYY-MM-DD or YYYY-MM-DDTHH:mm:ssZ).
	To string `url:"to,omitempty"`
}

// ProjectAnalysesUpdateEventOption represents options for updating an event.
type ProjectAnalysesUpdateEventOption struct {
	// Event is the event key (required).
	Event string `url:"event,omitempty"`
	// Name is the new event name (required).
	Name string `url:"name,omitempty"`
}

// -----------------------------------------------------------------------------
// Allowed Values
// -----------------------------------------------------------------------------

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	allowedEventCategories = map[string]struct{}{
		"VERSION": {},
		"OTHER":   {},
	}

	allowedSearchCategories = map[string]struct{}{
		"VERSION":           {},
		"OTHER":             {},
		"QUALITY_PROFILE":   {},
		"QUALITY_GATE":      {},
		"DEFINITION_CHANGE": {},
		"SQ_UPGRADE":        {},
	}
)

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCreateEventOpt validates the options for the CreateEvent method.
func (s *ProjectAnalysesService) ValidateCreateEventOpt(opt *ProjectAnalysesCreateEventOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Analysis, "Analysis")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	if opt.Category != "" {
		err = IsValueAuthorized(opt.Category, allowedEventCategories, "Category")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteOpt validates the options for the Delete method.
func (s *ProjectAnalysesService) ValidateDeleteOpt(opt *ProjectAnalysesDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Analysis, "Analysis")
}

// ValidateDeleteEventOpt validates the options for the DeleteEvent method.
func (s *ProjectAnalysesService) ValidateDeleteEventOpt(opt *ProjectAnalysesDeleteEventOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Event, "Event")
}

// validateDateFormat validates that a date string is in a valid format.
func validateDateFormat(value, fieldName string) error {
	if value != "" && !isValidDate(value) && !isValidDateTime(value) {
		return NewValidationError(fieldName, "must be a valid date (YYYY-MM-DD) or datetime (YYYY-MM-DDTHH:mm:ssZ)", ErrInvalidFormat)
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
func (s *ProjectAnalysesService) ValidateSearchOpt(opt *ProjectAnalysesSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	if opt.Category != "" {
		err = IsValueAuthorized(opt.Category, allowedSearchCategories, "Category")
		if err != nil {
			return err
		}
	}

	err = validateDateFormat(opt.From, "From")
	if err != nil {
		return err
	}

	return validateDateFormat(opt.To, "To")
}

// ValidateUpdateEventOpt validates the options for the UpdateEvent method.
func (s *ProjectAnalysesService) ValidateUpdateEventOpt(opt *ProjectAnalysesUpdateEventOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Event, "Event")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Name, "Name")
}

// isValidDate checks if a string is a valid date in YYYY-MM-DD format.
func isValidDate(dateStr string) bool {
	if len(dateStr) != dateLen {
		return false
	}

	// Simple validation: YYYY-MM-DD format
	parts := strings.Split(dateStr, "-")
	if len(parts) != dateParts {
		return false
	}

	if len(parts[0]) != yearLen || len(parts[1]) != monthDayLen || len(parts[2]) != monthDayLen {
		return false
	}

	return true
}

// isValidDateTime checks if a string is a valid datetime in YYYY-MM-DDTHH:mm:ssZ format.
func isValidDateTime(dateTimeStr string) bool {
	// Minimum length: YYYY-MM-DDTHH:mm:ssZ = 20
	if len(dateTimeStr) < dateTimeMinLen {
		return false
	}

	// Check for T separator
	if !strings.Contains(dateTimeStr, "T") {
		return false
	}

	// Check date part
	parts := strings.Split(dateTimeStr, "T")
	if len(parts) != dateTimeParts {
		return false
	}

	return isValidDate(parts[0])
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CreateEvent creates an event on a project analysis.
// Only events of category 'VERSION' and 'OTHER' can be created.
// Requires 'Administer' permission on the project.
//
// API endpoint: POST /api/project_analyses/create_event.
// Since: 6.3.
func (s *ProjectAnalysesService) CreateEvent(opt *ProjectAnalysesCreateEventOption) (*ProjectAnalysesEvent, *http.Response, error) {
	err := s.ValidateCreateEventOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_analyses/create_event", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectAnalysesEvent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes a project analysis.
// Requires 'Administer' permission on the project.
//
// API endpoint: POST /api/project_analyses/delete.
// Since: 6.3.
func (s *ProjectAnalysesService) Delete(opt *ProjectAnalysesDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_analyses/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteEvent deletes a project analysis event.
// Only events of category 'VERSION' and 'OTHER' can be deleted.
// Requires 'Administer' permission on the project.
//
// API endpoint: POST /api/project_analyses/delete_event.
// Since: 6.3.
func (s *ProjectAnalysesService) DeleteEvent(opt *ProjectAnalysesDeleteEventOption) (*http.Response, error) {
	err := s.ValidateDeleteEventOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_analyses/delete_event", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Search searches project analyses and attached events.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/project_analyses/search.
// Since: 6.3.
func (s *ProjectAnalysesService) Search(opt *ProjectAnalysesSearchOption) (*ProjectAnalysesSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_analyses/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectAnalysesSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchAll is a convenience method to iterate over all analyses.
// It handles pagination automatically.
func (s *ProjectAnalysesService) SearchAll(opt *ProjectAnalysesSearchOption) ([]ProjectAnalysis, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	var allAnalyses []ProjectAnalysis

	// Copy opt to avoid modifying the original
	searchOpt := *opt

	searchOpt.Page = 1
	if searchOpt.PageSize == 0 {
		searchOpt.PageSize = 100
	}

	for {
		result, resp, err := s.Search(&searchOpt)
		if err != nil {
			return nil, resp, err
		}

		allAnalyses = append(allAnalyses, result.Analyses...)

		if int64(len(allAnalyses)) >= result.Paging.Total {
			return allAnalyses, resp, nil
		}

		searchOpt.Page++
	}
}

// UpdateEvent updates an event name.
// Only events of category 'VERSION' and 'OTHER' can be updated.
// Requires 'Administer' permission on the project.
//
// API endpoint: POST /api/project_analyses/update_event.
// Since: 6.3.
func (s *ProjectAnalysesService) UpdateEvent(opt *ProjectAnalysesUpdateEventOption) (*ProjectAnalysesEvent, *http.Response, error) {
	err := s.ValidateUpdateEventOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_analyses/update_event", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectAnalysesEvent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
