package sonargo

import "net/http"

const (
	// MaxCePageSize is the maximum page size for CE activity queries.
	MaxCePageSize = 1000
	// MaxProjectKeyLength is the maximum length for a project key.
	MaxProjectKeyLength = 400
	// MaxProjectNameLength is the maximum length for a project name (truncated if longer).
	MaxProjectNameLength = 500
)

// CeService handles communication with the Compute Engine related methods
// of the SonarQube API.
// This service provides information on Compute Engine tasks, including
// task activity, status, and management operations.
type CeService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedTaskStatuses is the set of supported task statuses.
	allowedTaskStatuses = map[string]struct{}{
		"SUCCESS":     {},
		"FAILED":      {},
		"CANCELED":    {},
		"PENDING":     {},
		"IN_PROGRESS": {},
	}

	// allowedTaskTypes is the set of supported task types.
	allowedTaskTypes = map[string]struct{}{
		"REPORT":         {},
		"ISSUE_SYNC":     {},
		"AUDIT_PURGE":    {},
		"PROJECT_EXPORT": {},
	}

	// allowedAdditionalFields is the set of supported additional fields for task queries.
	allowedAdditionalFields = map[string]struct{}{
		"stacktrace":     {},
		"scannerContext": {},
		"warnings":       {},
	}
)

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// CeTask represents a Compute Engine task.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CeTask struct {
	// AnalysisID is the unique identifier of the analysis.
	AnalysisID string `json:"analysisId,omitempty"`
	// Branch is the branch name (if applicable).
	Branch string `json:"branch,omitempty"`
	// BranchType is the type of branch (e.g., "BRANCH", "PULL_REQUEST").
	BranchType string `json:"branchType,omitempty"`
	// ComponentID is the component identifier.
	ComponentID string `json:"componentId,omitempty"`
	// ComponentKey is the component key.
	ComponentKey string `json:"componentKey,omitempty"`
	// ComponentName is the component name.
	ComponentName string `json:"componentName,omitempty"`
	// ComponentQualifier is the component qualifier (e.g., "TRK" for project).
	ComponentQualifier string `json:"componentQualifier,omitempty"`
	// ErrorMessage is the error message if the task failed.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// ErrorStacktrace is the error stacktrace if available and requested.
	ErrorStacktrace string `json:"errorStacktrace,omitempty"`
	// ErrorType is the type of error if the task failed.
	ErrorType string `json:"errorType,omitempty"`
	// ExecutedAt is the timestamp when the task finished execution.
	ExecutedAt string `json:"executedAt,omitempty"`
	// ExecutionTimeMs is the execution time in milliseconds.
	ExecutionTimeMs int64 `json:"executionTimeMs,omitempty"`
	// FinishedAt is the timestamp when the task finished.
	FinishedAt string `json:"finishedAt,omitempty"`
	// HasErrorStacktrace indicates if an error stacktrace is available.
	HasErrorStacktrace bool `json:"hasErrorStacktrace,omitempty"`
	// HasScannerContext indicates if scanner context is available.
	HasScannerContext bool `json:"hasScannerContext,omitempty"`
	// ID is the unique identifier of the task.
	ID string `json:"id,omitempty"`
	// InfoMessages contains informational messages about the task.
	InfoMessages []string `json:"infoMessages,omitempty"`
	// Organization is the organization key (deprecated).
	Organization string `json:"organization,omitempty"`
	// PullRequest is the pull request identifier (if applicable).
	PullRequest string `json:"pullRequest,omitempty"`
	// ScannerContext is the scanner context if requested.
	ScannerContext string `json:"scannerContext,omitempty"`
	// StartedAt is the timestamp when the task started.
	StartedAt string `json:"startedAt,omitempty"`
	// Status is the task status (SUCCESS, FAILED, CANCELED, PENDING, IN_PROGRESS).
	Status string `json:"status,omitempty"`
	// SubmittedAt is the timestamp when the task was submitted.
	SubmittedAt string `json:"submittedAt,omitempty"`
	// SubmitterLogin is the login of the user who submitted the task.
	SubmitterLogin string `json:"submitterLogin,omitempty"`
	// Type is the task type (e.g., "REPORT").
	Type string `json:"type,omitempty"`
	// WarningCount is the number of warnings.
	WarningCount int64 `json:"warningCount,omitempty"`
	// Warnings contains warning messages.
	Warnings []string `json:"warnings,omitempty"`
}

// CeQueuedTask represents a task in the queue (pending or in-progress).
type CeQueuedTask struct {
	// ComponentID is the component identifier.
	ComponentID string `json:"componentId,omitempty"`
	// ComponentKey is the component key.
	ComponentKey string `json:"componentKey,omitempty"`
	// ComponentName is the component name.
	ComponentName string `json:"componentName,omitempty"`
	// ComponentQualifier is the component qualifier.
	ComponentQualifier string `json:"componentQualifier,omitempty"`
	// ID is the unique identifier of the task.
	ID string `json:"id,omitempty"`
	// Organization is the organization key (deprecated).
	Organization string `json:"organization,omitempty"`
	// Status is the task status.
	Status string `json:"status,omitempty"`
	// SubmittedAt is the timestamp when the task was submitted.
	SubmittedAt string `json:"submittedAt,omitempty"`
	// Type is the task type.
	Type string `json:"type,omitempty"`
}

// CePaging represents pagination information for CE queries.
type CePaging struct {
	// PageIndex is the current page index (1-based).
	PageIndex int64 `json:"pageIndex,omitempty"`
	// PageSize is the number of items per page.
	PageSize int64 `json:"pageSize,omitempty"`
	// Total is the total number of items.
	Total int64 `json:"total,omitempty"`
}

// AnalysisWarning represents a warning from analysis.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type AnalysisWarning struct {
	// Dismissable indicates if the warning can be dismissed.
	Dismissable bool `json:"dismissable,omitempty"`
	// Key is the unique key of the warning.
	Key string `json:"key,omitempty"`
	// Message is the warning message.
	Message string `json:"message,omitempty"`
}

// AnalysisComponent represents a component with its analysis warnings.
type AnalysisComponent struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Warnings is the list of analysis warnings.
	Warnings []AnalysisWarning `json:"warnings,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// CeActivity represents the response from searching CE tasks.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CeActivity struct {
	// Paging contains pagination information.
	Paging CePaging `json:"paging,omitzero"`
	// Tasks is the list of tasks.
	Tasks []CeTask `json:"tasks,omitempty"`
}

// CeActivityStatus represents CE activity metrics.
type CeActivityStatus struct {
	// Failing is the number of failing tasks.
	Failing int64 `json:"failing,omitempty"`
	// InProgress is the number of in-progress tasks.
	InProgress int64 `json:"inProgress,omitempty"`
	// Pending is the number of pending tasks.
	Pending int64 `json:"pending,omitempty"`
	// PendingTime is the total pending time in milliseconds (only included when there are pending tasks).
	PendingTime int64 `json:"pendingTime,omitempty"`
}

// CeAnalysisStatus represents the analysis status of a component.
type CeAnalysisStatus struct {
	// Component contains the component information with warnings.
	Component AnalysisComponent `json:"component,omitzero"`
}

// CeComponent represents the pending, in-progress, and last executed tasks for a component.
type CeComponent struct {
	// Current is the last executed task.
	Current CeTask `json:"current,omitzero"`
	// Queue contains pending and in-progress tasks.
	Queue []CeQueuedTask `json:"queue,omitempty"`
}

// CeIndexationStatus represents the indexation status.
type CeIndexationStatus struct {
	// CompletedCount is the number of projects with completed indexing.
	CompletedCount int64 `json:"completedCount,omitempty"`
	// HasFailures indicates if there are indexation failures.
	HasFailures bool `json:"hasFailures,omitempty"`
	// IsCompleted indicates if indexation is complete.
	IsCompleted bool `json:"isCompleted,omitempty"`
	// Total is the total number of projects to index.
	Total int64 `json:"total,omitempty"`
}

// CeInfo represents Compute Engine information.
type CeInfo struct {
	// WorkersPauseStatus is the pause status of CE workers.
	WorkersPauseStatus string `json:"workersPauseStatus,omitempty"`
}

// CeSubmit represents the response from submitting a scanner report.
type CeSubmit struct {
	// ProjectID is the project identifier.
	ProjectID string `json:"projectId,omitempty"`
	// TaskID is the created task identifier.
	TaskID string `json:"taskId,omitempty"`
}

// CeTaskDetails represents the response from getting task details.
type CeTaskDetails struct {
	// Task is the task details.
	Task CeTask `json:"task,omitzero"`
}

// CeTaskTypes represents the available task types.
type CeTaskTypes struct {
	// TaskTypes is the list of available task types.
	TaskTypes []string `json:"taskTypes,omitempty"`
}

// CeWorkerCount represents the CE worker count.
type CeWorkerCount struct {
	// CanSetWorkerCount indicates if the worker count can be modified.
	CanSetWorkerCount bool `json:"canSetWorkerCount,omitempty"`
	// Value is the number of workers.
	Value int64 `json:"value,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// CeActivityOption contains parameters for the Activity method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CeActivityOption struct {
	CePaginationArgs

	// Component is the key of the component (project) to filter on.
	Component string `url:"component,omitempty"`
	// MaxExecutedAt is the maximum date of end of task processing (inclusive).
	// Format: ISO 8601 datetime (e.g., 2017-10-19T13:00:00+0200).
	MaxExecutedAt string `url:"maxExecutedAt,omitempty"`
	// MinSubmittedAt is the minimum date of task submission (inclusive).
	// Format: ISO 8601 datetime (e.g., 2017-10-19T13:00:00+0200).
	MinSubmittedAt string `url:"minSubmittedAt,omitempty"`
	// OnlyCurrents filters on the last tasks (only the most recent finished task by project).
	OnlyCurrents bool `url:"onlyCurrents,omitempty"`
	// Q limits search to component names containing the string, component keys matching exactly,
	// or task IDs matching exactly.
	Q string `url:"q,omitempty"`
	// Statuses filters by task statuses.
	// Allowed values: SUCCESS, FAILED, CANCELED, PENDING, IN_PROGRESS.
	Statuses []string `url:"status,omitempty,comma"`
	// Type filters by task type.
	// Allowed values: REPORT, ISSUE_SYNC, AUDIT_PURGE, PROJECT_EXPORT.
	Type string `url:"type,omitempty"`
}

// CePaginationArgs contains pagination parameters specific to CE (supports up to 1000).
type CePaginationArgs struct {
	// Page is the response page number. Must be strictly greater than 0.
	Page int64 `url:"p,omitempty"`
	// PageSize is the response page size. Must be greater than 0 and less than or equal to 1000.
	PageSize int64 `url:"ps,omitempty"`
}

// Validate validates the CE pagination arguments.
func (p *CePaginationArgs) Validate() error {
	if p.Page != 0 && p.Page < MinPageSize {
		return NewValidationError("Page", "must be greater than 0", ErrOutOfRange)
	}

	if p.PageSize != 0 && (p.PageSize < MinPageSize || p.PageSize > MaxCePageSize) {
		return NewValidationError("PageSize", "must be between 1 and 1000", ErrOutOfRange)
	}

	return nil
}

// CeActivityStatusOption contains parameters for the ActivityStatus method.
type CeActivityStatusOption struct {
	// Component is the key of the component (project) to filter on.
	Component string `url:"component,omitempty"`
}

// CeAnalysisStatusOption contains parameters for the AnalysisStatus method.
type CeAnalysisStatusOption struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Component is the component key.
	// This field is required.
	Component string `url:"component"`
	// PullRequest is the pull request ID.
	PullRequest string `url:"pullRequest,omitempty"`
}

// CeCancelOption contains parameters for the Cancel method.
type CeCancelOption struct {
	// ID is the ID of the task to cancel.
	// This field is required.
	ID string `url:"id"`
}

// CeComponentOption contains parameters for the Component method.
type CeComponentOption struct {
	// Component is the component key.
	// This field is required.
	Component string `url:"component"`
}

// CeDismissAnalysisWarningOption contains parameters for the DismissAnalysisWarning method.
type CeDismissAnalysisWarningOption struct {
	// Component is the key of the project.
	// This field is required.
	Component string `url:"component"`
	// Warning is the key of the warning to dismiss.
	// This field is required.
	Warning string `url:"warning"`
}

// CeSubmitOption contains parameters for the Submit method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CeSubmitOption struct {
	// Characteristics contains optional characteristics of the analysis.
	// Can contain multiple key=value pairs.
	Characteristics []string `url:"characteristic,omitempty,comma"`
	// ProjectKey is the key of the project.
	// This field is required. Maximum length is 400 characters.
	ProjectKey string `url:"projectKey"`
	// ProjectName is the optional name of the project.
	// Used only if the project does not exist yet. If longer than 500, it is abbreviated.
	ProjectName string `url:"projectName,omitempty"`
	// Report is the report file content.
	// This field is required. Format is not an API, it changes among SonarQube versions.
	Report string `url:"report"`
}

// CeTaskOption contains parameters for the Task method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CeTaskOption struct {
	// AdditionalFields is a list of optional fields to be returned in response.
	// Allowed values: stacktrace, scannerContext, warnings.
	AdditionalFields []string `url:"additionalFields,omitempty,comma"`
	// ID is the ID of the task.
	// This field is required.
	ID string `url:"id"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateActivityOpt validates the options for the Activity method.
func (s *CeService) ValidateActivityOpt(opt *CeActivityOption) error {
	if opt == nil {
		// Activity with no options is valid
		return nil
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	if len(opt.Statuses) > 0 {
		err = AreValuesAuthorized(opt.Statuses, allowedTaskStatuses, "Statuses")
		if err != nil {
			return err
		}
	}

	if opt.Type != "" {
		err = IsValueAuthorized(opt.Type, allowedTaskTypes, "Type")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateActivityStatusOpt validates the options for the ActivityStatus method.
func (s *CeService) ValidateActivityStatusOpt(opt *CeActivityStatusOption) error {
	// Options are optional; nothing to validate.
	return nil
}

// ValidateAnalysisStatusOpt validates the options for the AnalysisStatus method.
func (s *CeService) ValidateAnalysisStatusOpt(opt *CeAnalysisStatusOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCancelOpt validates the options for the Cancel method.
func (s *CeService) ValidateCancelOpt(opt *CeCancelOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ID, "ID")
	if err != nil {
		return err
	}

	return nil
}

// ValidateComponentOpt validates the options for the Component method.
func (s *CeService) ValidateComponentOpt(opt *CeComponentOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDismissAnalysisWarningOpt validates the options for the DismissAnalysisWarning method.
func (s *CeService) ValidateDismissAnalysisWarningOpt(opt *CeDismissAnalysisWarningOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Warning, "Warning")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSubmitOpt validates the options for the Submit method.
func (s *CeService) ValidateSubmitOpt(opt *CeSubmitOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ProjectKey, MaxProjectKeyLength, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Report, "Report")
	if err != nil {
		return err
	}

	return nil
}

// ValidateTaskOpt validates the options for the Task method.
func (s *CeService) ValidateTaskOpt(opt *CeTaskOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ID, "ID")
	if err != nil {
		return err
	}

	if len(opt.AdditionalFields) > 0 {
		err = AreValuesAuthorized(opt.AdditionalFields, allowedAdditionalFields, "AdditionalFields")
		if err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Activity searches for CE tasks.
// Requires the system administration permission, or project administration permission
// if component is set.
//
// API endpoint: GET /api/ce/activity.
// Since: 5.2.
func (s *CeService) Activity(opt *CeActivityOption) (*CeActivity, *http.Response, error) {
	err := s.ValidateActivityOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "ce/activity", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeActivity)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ActivityStatus returns CE activity related metrics.
// Requires 'Administer System' permission or 'Administer' rights on the specified project.
//
// API endpoint: GET /api/ce/activity_status.
// Since: 5.5.
func (s *CeService) ActivityStatus(opt *CeActivityStatusOption) (*CeActivityStatus, *http.Response, error) {
	err := s.ValidateActivityStatusOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "ce/activity_status", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeActivityStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// AnalysisStatus gets the analysis status of a given component: a project, branch or pull request.
// Requires 'Browse' permission on the specified component.
//
// API endpoint: GET /api/ce/analysis_status.
// Since: 7.4.
func (s *CeService) AnalysisStatus(opt *CeAnalysisStatusOption) (*CeAnalysisStatus, *http.Response, error) {
	err := s.ValidateAnalysisStatusOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "ce/analysis_status", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeAnalysisStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Cancel cancels a pending task.
// In-progress tasks cannot be canceled.
// Requires 'Administer System' or 'Administer' rights on the project related to the task.
//
// API endpoint: POST /api/ce/cancel.
// Since: 5.2.
func (s *CeService) Cancel(opt *CeCancelOption) (*http.Response, error) {
	err := s.ValidateCancelOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "ce/cancel", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CancelAll cancels all pending tasks.
// Requires system administration permission. In-progress tasks are not canceled.
//
// API endpoint: POST /api/ce/cancel_all.
// Since: 5.2.
func (s *CeService) CancelAll() (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "ce/cancel_all", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Component gets the pending tasks, in-progress tasks and the last executed task
// of a given component (usually a project).
// Requires 'Browse' permission on the specified component.
//
// API endpoint: GET /api/ce/component.
// Since: 5.2.
func (s *CeService) Component(opt *CeComponentOption) (*CeComponent, *http.Response, error) {
	err := s.ValidateComponentOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "ce/component", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeComponent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DismissAnalysisWarning permanently dismisses a specific analysis warning.
// Requires authentication and 'Browse' permission on the specified project.
//
// API endpoint: POST /api/ce/dismiss_analysis_warning.
// Since: 8.5.
func (s *CeService) DismissAnalysisWarning(opt *CeDismissAnalysisWarningOption) (*http.Response, error) {
	err := s.ValidateDismissAnalysisWarningOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "ce/dismiss_analysis_warning", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// IndexationStatus returns the count of projects with completed issue indexing.
//
// API endpoint: GET /api/ce/indexation_status.
// Since: 8.4.
func (s *CeService) IndexationStatus() (*CeIndexationStatus, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "ce/indexation_status", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeIndexationStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Info gets information about Compute Engine.
// Requires the system administration permission or system passcode.
//
// API endpoint: GET /api/ce/info.
// Since: 7.2.
func (s *CeService) Info() (*CeInfo, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "ce/info", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Pause requests pause of Compute Engine workers.
// Requires the system administration permission or system passcode.
//
// API endpoint: POST /api/ce/pause.
// Since: 7.2.
func (s *CeService) Pause() (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "ce/pause", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Resume resumes pause of Compute Engine workers.
// Requires the system administration permission or system passcode.
//
// API endpoint: POST /api/ce/resume.
// Since: 7.2.
func (s *CeService) Resume() (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "ce/resume", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Submit submits a scanner report to the queue.
// Report is processed asynchronously. Requires analysis permission.
// If the project does not exist, then the provisioning permission is also required.
//
// API endpoint: POST /api/ce/submit.
// Since: 5.2.
func (s *CeService) Submit(opt *CeSubmitOption) (*CeSubmit, *http.Response, error) {
	err := s.ValidateSubmitOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "ce/submit", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeSubmit)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Task gives Compute Engine task details such as type, status, duration
// and associated component.
// Requires 'Administer' at global or project level, or 'Execute Analysis' at global or project level.
//
// API endpoint: GET /api/ce/task.
// Since: 5.2.
func (s *CeService) Task(opt *CeTaskOption) (*CeTaskDetails, *http.Response, error) {
	err := s.ValidateTaskOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "ce/task", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeTaskDetails)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// TaskTypes lists available task types.
//
// API endpoint: GET /api/ce/task_types.
// Since: 5.5.
func (s *CeService) TaskTypes() (*CeTaskTypes, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "ce/task_types", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeTaskTypes)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// WorkerCount returns the number of Compute Engine workers.
// Requires the system administration permission.
//
// API endpoint: GET /api/ce/worker_count.
// Since: 6.5.
func (s *CeService) WorkerCount() (*CeWorkerCount, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "ce/worker_count", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(CeWorkerCount)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
