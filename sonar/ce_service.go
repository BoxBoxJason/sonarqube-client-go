// Get information on Compute Engine tasks.
package sonargo

import "net/http"

type CeService struct {
	client *Client
}

type CeActivityObject struct {
	Paging CeActivityObject_sub1   `json:"paging,omitempty"`
	Tasks  []CeActivityObject_sub2 `json:"tasks,omitempty"`
}

type CeActivityObject_sub2 struct {
	AnalysisID         string   `json:"analysisId,omitempty"`
	ComponentID        string   `json:"componentId,omitempty"`
	ComponentKey       string   `json:"componentKey,omitempty"`
	ComponentName      string   `json:"componentName,omitempty"`
	ComponentQualifier string   `json:"componentQualifier,omitempty"`
	ExecutedAt         string   `json:"executedAt,omitempty"`
	ExecutionTimeMs    int64    `json:"executionTimeMs,omitempty"`
	HasScannerContext  bool     `json:"hasScannerContext,omitempty"`
	ID                 string   `json:"id,omitempty"`
	InfoMessages       []string `json:"infoMessages,omitempty"`
	StartedAt          string   `json:"startedAt,omitempty"`
	Status             string   `json:"status,omitempty"`
	SubmittedAt        string   `json:"submittedAt,omitempty"`
	SubmitterLogin     string   `json:"submitterLogin,omitempty"`
	Type               string   `json:"type,omitempty"`
	WarningCount       int64    `json:"warningCount,omitempty"`
	Warnings           []string `json:"warnings,omitempty"`
}

type CeActivityObject_sub1 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type CeActivityStatusObject struct {
	Failing     int64 `json:"failing,omitempty"`
	InProgress  int64 `json:"inProgress,omitempty"`
	Pending     int64 `json:"pending,omitempty"`
	PendingTime int64 `json:"pendingTime,omitempty"`
}

type CeAnalysisStatusObject struct {
	Component CeAnalysisStatusObject_sub2 `json:"component,omitempty"`
}

type CeAnalysisStatusObject_sub1 struct {
	Dismissable bool   `json:"dismissable,omitempty"`
	Key         string `json:"key,omitempty"`
	Message     string `json:"message,omitempty"`
}

type CeAnalysisStatusObject_sub2 struct {
	Key      string                        `json:"key,omitempty"`
	Name     string                        `json:"name,omitempty"`
	Warnings []CeAnalysisStatusObject_sub1 `json:"warnings,omitempty"`
}

type CeComponentObject struct {
	Current CeComponentObject_sub1   `json:"current,omitempty"`
	Queue   []CeComponentObject_sub2 `json:"queue,omitempty"`
}

type CeComponentObject_sub1 struct {
	AnalysisID         string `json:"analysisId,omitempty"`
	ComponentID        string `json:"componentId,omitempty"`
	ComponentKey       string `json:"componentKey,omitempty"`
	ComponentName      string `json:"componentName,omitempty"`
	ComponentQualifier string `json:"componentQualifier,omitempty"`
	ErrorMessage       string `json:"errorMessage,omitempty"`
	ErrorType          string `json:"errorType,omitempty"`
	ExecutionTimeMs    int64  `json:"executionTimeMs,omitempty"`
	FinishedAt         string `json:"finishedAt,omitempty"`
	HasErrorStacktrace bool   `json:"hasErrorStacktrace,omitempty"`
	HasScannerContext  bool   `json:"hasScannerContext,omitempty"`
	ID                 string `json:"id,omitempty"`
	Organization       string `json:"organization,omitempty"`
	StartedAt          string `json:"startedAt,omitempty"`
	Status             string `json:"status,omitempty"`
	SubmittedAt        string `json:"submittedAt,omitempty"`
	Type               string `json:"type,omitempty"`
}

type CeComponentObject_sub2 struct {
	ComponentID        string `json:"componentId,omitempty"`
	ComponentKey       string `json:"componentKey,omitempty"`
	ComponentName      string `json:"componentName,omitempty"`
	ComponentQualifier string `json:"componentQualifier,omitempty"`
	ID                 string `json:"id,omitempty"`
	Organization       string `json:"organization,omitempty"`
	Status             string `json:"status,omitempty"`
	SubmittedAt        string `json:"submittedAt,omitempty"`
	Type               string `json:"type,omitempty"`
}

type CeIndexationStatusObject struct {
	CompletedCount int64 `json:"completedCount,omitempty"`
	HasFailures    bool  `json:"hasFailures,omitempty"`
	IsCompleted    bool  `json:"isCompleted,omitempty"`
	Total          int64 `json:"total,omitempty"`
}

type CeInfoObject struct {
	WorkersPauseStatus string `json:"workersPauseStatus,omitempty"`
}

type CeSubmitObject struct {
	ProjectID string `json:"projectId,omitempty"`
	TaskID    string `json:"taskId,omitempty"`
}

type CeTaskObject struct {
	Task CeTaskObject_sub1 `json:"task,omitempty"`
}

type CeTaskObject_sub1 struct {
	AnalysisID         string   `json:"analysisId,omitempty"`
	ComponentID        string   `json:"componentId,omitempty"`
	ComponentKey       string   `json:"componentKey,omitempty"`
	ComponentName      string   `json:"componentName,omitempty"`
	ComponentQualifier string   `json:"componentQualifier,omitempty"`
	ExecutedAt         string   `json:"executedAt,omitempty"`
	ExecutionTimeMs    int64    `json:"executionTimeMs,omitempty"`
	HasScannerContext  bool     `json:"hasScannerContext,omitempty"`
	ID                 string   `json:"id,omitempty"`
	StartedAt          string   `json:"startedAt,omitempty"`
	Status             string   `json:"status,omitempty"`
	SubmittedAt        string   `json:"submittedAt,omitempty"`
	SubmitterLogin     string   `json:"submitterLogin,omitempty"`
	Type               string   `json:"type,omitempty"`
	WarningCount       int64    `json:"warningCount,omitempty"`
	Warnings           []string `json:"warnings,omitempty"`
}

type CeTaskTypesObject struct {
	TaskTypes []string `json:"taskTypes,omitempty"`
}

type CeWorkerCountObject struct {
	CanSetWorkerCount bool  `json:"canSetWorkerCount,omitempty"`
	Value             int64 `json:"value,omitempty"`
}

type CeActivityOption struct {
	Component      string `url:"component,omitempty"`      // Description:"Key of the component (project) to filter on",ExampleValue:"projectKey"
	MaxExecutedAt  string `url:"maxExecutedAt,omitempty"`  // Description:"Maximum date of end of task processing (inclusive)",ExampleValue:"2017-10-19T13:00:00+0200"
	MinSubmittedAt string `url:"minSubmittedAt,omitempty"` // Description:"Minimum date of task submission (inclusive)",ExampleValue:"2017-10-19T13:00:00+0200"
	OnlyCurrents   string `url:"onlyCurrents,omitempty"`   // Description:"Filter on the last tasks (only the most recent finished task by project)",ExampleValue:""
	P              string `url:"p,omitempty"`              // Description:"1-based page number",ExampleValue:"42"
	Ps             string `url:"ps,omitempty"`             // Description:"Page size. Must be greater than 0 and less or equal than 1000",ExampleValue:"20"
	Q              string `url:"q,omitempty"`              // Description:"Limit search to: <ul><li>component names that contain the supplied string</li><li>component keys that are exactly the same as the supplied string</li><li>task ids that are exactly the same as the supplied string</li></ul>",ExampleValue:"Apache"
	Status         string `url:"status,omitempty"`         // Description:"Comma separated list of task statuses",ExampleValue:"IN_PROGRESS,SUCCESS"
	Type           string `url:"type,omitempty"`           // Description:"Task type",ExampleValue:"REPORT"
}

// Activity Search for tasks.<br> Requires the system administration permission, or project administration permission if component is set.
func (s *CeService) Activity(opt *CeActivityOption) (v *CeActivityObject, resp *http.Response, err error) {
	err = s.ValidateActivityOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "ce/activity", opt)
	if err != nil {
		return
	}
	v = new(CeActivityObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type CeActivityStatusOption struct {
	Component string `url:"component,omitempty"` // Description:"Key of the component (project) to filter on",ExampleValue:"my_project"
}

// ActivityStatus Returns CE activity related metrics.<br>Requires 'Administer System' permission or 'Administer' rights on the specified project.
func (s *CeService) ActivityStatus(opt *CeActivityStatusOption) (v *CeActivityStatusObject, resp *http.Response, err error) {
	err = s.ValidateActivityStatusOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "ce/activity_status", opt)
	if err != nil {
		return
	}
	v = new(CeActivityStatusObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type CeAnalysisStatusOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	Component   string `url:"component,omitempty"`   // Description:"",ExampleValue:"my_project"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
}

// AnalysisStatus Get the analysis status of a given component: a project, branch or pull request.<br>Requires the following permission: 'Browse' on the specified component.
func (s *CeService) AnalysisStatus(opt *CeAnalysisStatusOption) (v *CeAnalysisStatusObject, resp *http.Response, err error) {
	err = s.ValidateAnalysisStatusOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "ce/analysis_status", opt)
	if err != nil {
		return
	}
	v = new(CeAnalysisStatusObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type CeCancelOption struct {
	Id string `url:"id,omitempty"` // Description:"Id of the task to cancel.",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
}

// Cancel Cancels a pending task.<br/>In-progress tasks cannot be canceled.<br/>Requires one of the following permissions:<ul><li>'Administer System'</li><li>'Administer' rights on the project related to the task</li></ul>
func (s *CeService) Cancel(opt *CeCancelOption) (resp *http.Response, err error) {
	err = s.ValidateCancelOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "ce/cancel", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

// CancelAll Cancels all pending tasks. Requires system administration permission. In-progress tasks are not canceled.
func (s *CeService) CancelAll() (resp *http.Response, err error) {
	req, err := s.client.NewRequest("POST", "ce/cancel_all", nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type CeComponentOption struct {
	Component string `url:"component,omitempty"` // Description:"",ExampleValue:"my_project"
}

// Component Get the pending tasks, in-progress tasks and the last executed task of a given component (usually a project).<br>Requires the following permission: 'Browse' on the specified component.
func (s *CeService) Component(opt *CeComponentOption) (v *CeComponentObject, resp *http.Response, err error) {
	err = s.ValidateComponentOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "ce/component", opt)
	if err != nil {
		return
	}
	v = new(CeComponentObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type CeDismissAnalysisWarningOption struct {
	Component string `url:"component,omitempty"` // Description:"Key of the project",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
	Warning   string `url:"warning,omitempty"`   // Description:"Key of the warning",ExampleValue:"AU-TpxcA-iU5OvuD2FLz"
}

// DismissAnalysisWarning Permanently dismiss a specific analysis warning. Requires authentication and 'Browse' permission on the specified project.
func (s *CeService) DismissAnalysisWarning(opt *CeDismissAnalysisWarningOption) (resp *http.Response, err error) {
	err = s.ValidateDismissAnalysisWarningOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "ce/dismiss_analysis_warning", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

// IndexationStatus Returns the count of projects with completed issue indexing.
func (s *CeService) IndexationStatus() (v *CeIndexationStatusObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "ce/indexation_status", nil)
	if err != nil {
		return
	}
	v = new(CeIndexationStatusObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// Info Gets information about Compute Engine. Requires the system administration permission or system passcode (see WEB_SYSTEM_PASS_CODE in sonar.properties).
func (s *CeService) Info() (v *CeInfoObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "ce/info", nil)
	if err != nil {
		return
	}
	v = new(CeInfoObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// Pause Requests pause of Compute Engine workers. Requires the system administration permission or system passcode (see WEB_SYSTEM_PASS_CODE in sonar.properties).
func (s *CeService) Pause() (resp *http.Response, err error) {
	req, err := s.client.NewRequest("POST", "ce/pause", nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

// Resume Resumes pause of Compute Engine workers. Requires the system administration permission or system passcode (see WEB_SYSTEM_PASS_CODE in sonar.properties).
func (s *CeService) Resume() (resp *http.Response, err error) {
	req, err := s.client.NewRequest("POST", "ce/resume", nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type CeSubmitOption struct {
	Characteristic string `url:"characteristic,omitempty"` // Description:"Optional characteristic of the analysis. Can be repeated to define multiple characteristics.",ExampleValue:"branchType=long"
	ProjectKey     string `url:"projectKey,omitempty"`     // Description:"Key of project",ExampleValue:"my_project"
	ProjectName    string `url:"projectName,omitempty"`    // Description:"Optional name of the project, used only if the project does not exist yet. If name is longer than 500, it is abbreviated.",ExampleValue:"My Project"
	Report         string `url:"report,omitempty"`         // Description:"Report file. Format is not an API, it changes among SonarQube versions.",ExampleValue:""
}

// Submit Submits a scanner report to the queue. Report is processed asynchronously. Requires analysis permission. If the project does not exist, then the provisioning permission is also required.
func (s *CeService) Submit(opt *CeSubmitOption) (v *CeSubmitObject, resp *http.Response, err error) {
	err = s.ValidateSubmitOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "ce/submit", opt)
	if err != nil {
		return
	}
	v = new(CeSubmitObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type CeTaskOption struct {
	AdditionalFields string `url:"additionalFields,omitempty"` // Description:"Comma-separated list of the optional fields to be returned in response.",ExampleValue:""
	Id               string `url:"id,omitempty"`               // Description:"Id of task",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
}

// Task Give Compute Engine task details such as type, status, duration and associated component.<br/>Requires one of the following permissions: <ul><li>'Administer' at global or project level</li><li>'Execute Analysis' at global or project level</li></ul>Since 6.1, field "logs" is deprecated and its value is always false.
func (s *CeService) Task(opt *CeTaskOption) (v *CeTaskObject, resp *http.Response, err error) {
	err = s.ValidateTaskOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "ce/task", opt)
	if err != nil {
		return
	}
	v = new(CeTaskObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// TaskTypes List available task types
func (s *CeService) TaskTypes() (v *CeTaskTypesObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "ce/task_types", nil)
	if err != nil {
		return
	}
	v = new(CeTaskTypesObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// WorkerCount Return number of Compute Engine workers.<br/>Requires the system administration permission
func (s *CeService) WorkerCount() (v *CeWorkerCountObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "ce/worker_count", nil)
	if err != nil {
		return
	}
	v = new(CeWorkerCountObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
