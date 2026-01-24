package sonargo

import (
	"net/http"
)

const (
	// MaxQualityGateNameLength is the maximum allowed length for quality gate names.
	MaxQualityGateNameLength = 100
	// MaxConditionErrorLength is the maximum allowed length for condition error thresholds.
	MaxConditionErrorLength = 64
)

// QualitygatesService handles communication with the Quality Gates related methods of the SonarQube API.
// Quality gates define conditions that projects must meet to pass analysis.
type QualitygatesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// QualitygatesCreate represents the response from creating a quality gate.
type QualitygatesCreate struct {
	// ID is the deprecated unique identifier of the quality gate.
	//
	// Deprecated: Since SonarQube 10.0 - use Name instead.
	ID string `json:"id,omitempty"`
	// Name is the name of the created quality gate.
	Name string `json:"name,omitempty"`
}

// QualitygatesCreateCondition represents the response from creating a condition.
type QualitygatesCreateCondition struct {
	// Error is the error threshold value.
	Error string `json:"error,omitempty"`
	// ID is the unique identifier of the condition.
	ID string `json:"id,omitempty"`
	// Metric is the metric key for the condition.
	Metric string `json:"metric,omitempty"`
	// Op is the comparison operator (LT, GT).
	Op string `json:"op,omitempty"`
	// Warning is the deprecated warning threshold value.
	//
	// Deprecated: Warning thresholds are no longer supported.
	Warning string `json:"warning,omitempty"`
}

// QualitygatesGetByProject represents the response from getting a project's quality gate.
type QualitygatesGetByProject struct {
	// QualityGate contains the quality gate details for the project.
	QualityGate ProjectQualityGate `json:"qualityGate,omitzero"`
}

// ProjectQualityGate represents the quality gate associated with a project.
type ProjectQualityGate struct {
	// Name is the name of the quality gate.
	Name string `json:"name,omitempty"`
	// Default indicates if this is the default quality gate.
	Default bool `json:"default,omitempty"`
}

// QualitygatesList represents the response from listing quality gates.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesList struct {
	// Actions contains the global actions available to the current user.
	Actions QualitygatesActions `json:"actions,omitzero"`
	// Qualitygates is the list of quality gates.
	Qualitygates []QualityGate `json:"qualitygates,omitempty"`
}

// QualitygatesActions represents global actions available for quality gates.
type QualitygatesActions struct {
	// Create indicates if the current user can create quality gates.
	Create bool `json:"create,omitempty"`
}

// QualityGate represents a quality gate with its properties and available actions.
//
//nolint:govet,tagliatelle // Field alignment less important; JSON follows SonarQube API naming
type QualityGate struct {
	// Actions contains the actions available for this quality gate.
	Actions QualityGateActions `json:"actions,omitzero"`
	// CaycStatus is the Clean As You Code status.
	// Possible values: compliant, non-compliant, over-compliant.
	CaycStatus string `json:"caycStatus,omitempty"`
	// Name is the name of the quality gate.
	Name string `json:"name,omitempty"`
	// HasMQRConditions indicates if the gate has MQR (Multi-Quality Rule) conditions.
	HasMQRConditions bool `json:"hasMQRConditions,omitempty"`
	// HasStandardConditions indicates if the gate has standard conditions.
	HasStandardConditions bool `json:"hasStandardConditions,omitempty"`
	// IsAiCodeSupported indicates if the gate supports AI-generated code.
	IsAiCodeSupported bool `json:"isAiCodeSupported,omitempty"`
	// IsBuiltIn indicates if this is a built-in quality gate.
	IsBuiltIn bool `json:"isBuiltIn,omitempty"`
	// IsDefault indicates if this is the default quality gate.
	IsDefault bool `json:"isDefault,omitempty"`
}

// QualityGateActions represents actions available for a specific quality gate.
type QualityGateActions struct {
	// AssociateProjects indicates if projects can be associated with this gate.
	AssociateProjects bool `json:"associateProjects,omitempty"`
	// Copy indicates if the gate can be copied.
	Copy bool `json:"copy,omitempty"`
	// Delegate indicates if permissions can be delegated for this gate.
	Delegate bool `json:"delegate,omitempty"`
	// Delete indicates if the gate can be deleted.
	Delete bool `json:"delete,omitempty"`
	// ManageAiCodeAssurance indicates if AI code assurance can be managed.
	ManageAiCodeAssurance bool `json:"manageAiCodeAssurance,omitempty"`
	// ManageConditions indicates if conditions can be managed.
	ManageConditions bool `json:"manageConditions,omitempty"`
	// Rename indicates if the gate can be renamed.
	Rename bool `json:"rename,omitempty"`
	// SetAsDefault indicates if the gate can be set as default.
	SetAsDefault bool `json:"setAsDefault,omitempty"`
}

// QualitygatesProjectStatus represents the quality gate status of a project.
type QualitygatesProjectStatus struct {
	// ProjectStatus contains the detailed project status information.
	ProjectStatus ProjectStatus `json:"projectStatus,omitzero"`
}

// ProjectStatus represents the detailed status of a project's quality gate.
//
//nolint:govet // Field alignment is less important than logical grouping
type ProjectStatus struct {
	// Status is the overall quality gate status (OK, WARN, ERROR, NONE).
	Status string `json:"status,omitempty"`
	// CaycStatus is the Clean As You Code status.
	CaycStatus string `json:"caycStatus,omitempty"`
	// Conditions is the list of condition evaluations.
	Conditions []ConditionStatus `json:"conditions,omitempty"`
	// Period contains information about the analysis period.
	Period AnalysisPeriod `json:"period,omitzero"`
	// IgnoredConditions indicates if some conditions were ignored.
	IgnoredConditions bool `json:"ignoredConditions,omitempty"`
}

// ConditionStatus represents the evaluation result of a single condition.
type ConditionStatus struct {
	// ActualValue is the actual measured value.
	ActualValue string `json:"actualValue,omitempty"`
	// Comparator is the comparison operator used.
	Comparator string `json:"comparator,omitempty"`
	// ErrorThreshold is the threshold that triggers an error.
	ErrorThreshold string `json:"errorThreshold,omitempty"`
	// MetricKey is the key of the evaluated metric.
	MetricKey string `json:"metricKey,omitempty"`
	// Status is the status of this condition (OK, WARN, ERROR).
	Status string `json:"status,omitempty"`
}

// AnalysisPeriod represents information about the analysis period.
type AnalysisPeriod struct {
	// Date is the date of the period.
	Date string `json:"date,omitempty"`
	// Mode is the period mode.
	Mode string `json:"mode,omitempty"`
	// Parameter is an optional parameter for the period.
	Parameter string `json:"parameter,omitempty"`
}

// QualitygatesSearch represents the response from searching for projects.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesSearch struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Results is the list of projects.
	Results []QualityGateProject `json:"results,omitempty"`
}

// QualityGateProject represents a project in quality gate search results.
type QualityGateProject struct {
	// Key is the unique key of the project.
	Key string `json:"key,omitempty"`
	// Name is the name of the project.
	Name string `json:"name,omitempty"`
	// ContainsAiCode indicates if the project contains AI-generated code.
	ContainsAiCode bool `json:"containsAiCode,omitempty"`
	// Selected indicates if the project is associated with the quality gate.
	Selected bool `json:"selected,omitempty"`
}

// QualitygatesSearchGroups represents the response from searching for groups.
type QualitygatesSearchGroups struct {
	// Groups is the list of groups.
	Groups []QualityGateGroup `json:"groups,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// QualityGateGroup represents a group that can edit a quality gate.
type QualityGateGroup struct {
	// Description is the description of the group.
	Description string `json:"description,omitempty"`
	// Name is the name of the group.
	Name string `json:"name,omitempty"`
	// Selected indicates if the group is allowed to edit the quality gate.
	Selected bool `json:"selected,omitempty"`
}

// QualitygatesSearchUsers represents the response from searching for users.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesSearchUsers struct {
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
	// Users is the list of users.
	Users []QualityGateUser `json:"users,omitempty"`
}

// QualityGateUser represents a user that can edit a quality gate.
type QualityGateUser struct {
	// Avatar is the avatar URL of the user.
	Avatar string `json:"avatar,omitempty"`
	// Login is the login name of the user.
	Login string `json:"login,omitempty"`
	// Name is the display name of the user.
	Name string `json:"name,omitempty"`
	// Selected indicates if the user is allowed to edit the quality gate.
	Selected bool `json:"selected,omitempty"`
}

// QualitygatesShow represents the response from showing a quality gate.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesShow struct {
	// Actions contains the actions available for this quality gate.
	Actions QualityGateActions `json:"actions,omitzero"`
	// Name is the name of the quality gate.
	Name string `json:"name,omitempty"`
	// CaycStatus is the Clean As You Code status.
	CaycStatus string `json:"caycStatus,omitempty"`
	// Conditions is the list of conditions defined for this quality gate.
	Conditions []QualityGateCondition `json:"conditions,omitempty"`
	// IsAiCodeSupported indicates if the gate supports AI-generated code.
	IsAiCodeSupported bool `json:"isAiCodeSupported,omitempty"`
	// IsBuiltIn indicates if this is a built-in quality gate.
	IsBuiltIn bool `json:"isBuiltIn,omitempty"`
	// IsDefault indicates if this is the default quality gate.
	IsDefault bool `json:"isDefault,omitempty"`
}

// QualityGateCondition represents a condition in a quality gate.
type QualityGateCondition struct {
	// Error is the error threshold value.
	Error string `json:"error,omitempty"`
	// ID is the unique identifier of the condition.
	ID string `json:"id,omitempty"`
	// Metric is the metric key for the condition.
	Metric string `json:"metric,omitempty"`
	// Op is the comparison operator (LT, GT).
	Op string `json:"op,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// QualitygatesAddGroupOption contains options for allowing a group to edit a quality gate.
type QualitygatesAddGroupOption struct {
	// GateName is the name of the quality gate (required).
	// Maximum length: 100 characters
	GateName string `url:"gateName,omitempty"`
	// GroupName is the group name or 'anyone' (case insensitive) (required).
	GroupName string `url:"groupName,omitempty"`
}

// QualitygatesAddUserOption contains options for allowing a user to edit a quality gate.
type QualitygatesAddUserOption struct {
	// GateName is the name of the quality gate (required).
	// Maximum length: 100 characters
	GateName string `url:"gateName,omitempty"`
	// Login is the user login (required).
	Login string `url:"login,omitempty"`
}

// QualitygatesCopyOption contains options for copying a quality gate.
type QualitygatesCopyOption struct {
	// Name is the name of the new quality gate to create (required).
	Name string `url:"name,omitempty"`
	// SourceName is the name of the quality gate to copy (required).
	// Maximum length: 100 characters
	SourceName string `url:"sourceName,omitempty"`
}

// QualitygatesCreateOption contains options for creating a quality gate.
type QualitygatesCreateOption struct {
	// Name is the name of the quality gate to create (required).
	// Maximum length: 100 characters
	Name string `url:"name,omitempty"`
}

// QualitygatesCreateConditionOption contains options for creating a condition.
type QualitygatesCreateConditionOption struct {
	// Error is the condition error threshold (required).
	// Maximum length: 64 characters
	Error string `url:"error,omitempty"`
	// GateName is the name of the quality gate (required).
	GateName string `url:"gateName,omitempty"`
	// Metric is the condition metric (required).
	// Only metrics of the following types are allowed: INT, MILLISEC, RATING, WORK_DUR, FLOAT, PERCENT, LEVEL.
	// Forbidden metrics: alert_status, security_hotspots, new_security_hotspots.
	Metric string `url:"metric,omitempty"`
	// Op is the condition operator (optional).
	// Allowed values: LT (is lower than), GT (is greater than)
	Op string `url:"op,omitempty"`
}

// QualitygatesDeleteConditionOption contains options for deleting a condition.
type QualitygatesDeleteConditionOption struct {
	// ID is the condition UUID (required).
	ID string `url:"id,omitempty"`
}

// QualitygatesDeselectOption contains options for removing a project association.
type QualitygatesDeselectOption struct {
	// ProjectKey is the project key (required).
	ProjectKey string `url:"projectKey,omitempty"`
}

// QualitygatesDestroyOption contains options for deleting a quality gate.
type QualitygatesDestroyOption struct {
	// Name is the name of the quality gate to delete (required).
	// Maximum length: 100 characters
	Name string `url:"name,omitempty"`
}

// QualitygatesGetByProjectOption contains options for getting a project's quality gate.
type QualitygatesGetByProjectOption struct {
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
}

// QualitygatesProjectStatusOption contains options for getting project status.
type QualitygatesProjectStatusOption struct {
	// AnalysisID is the analysis id (optional).
	// Either AnalysisID, ProjectID, or ProjectKey must be provided.
	AnalysisID string `url:"analysisId,omitempty"`
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// ProjectID is the project UUID (optional).
	// Doesn't work with branches or pull requests.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key (optional).
	ProjectKey string `url:"projectKey,omitempty"`
	// PullRequest is the pull request id (optional).
	PullRequest string `url:"pullRequest,omitempty"`
}

// QualitygatesRemoveGroupOption contains options for removing group permissions.
type QualitygatesRemoveGroupOption struct {
	// GateName is the name of the quality gate (required).
	// Maximum length: 100 characters
	GateName string `url:"gateName,omitempty"`
	// GroupName is the group name or 'anyone' (case insensitive) (required).
	GroupName string `url:"groupName,omitempty"`
}

// QualitygatesRemoveUserOption contains options for removing user permissions.
type QualitygatesRemoveUserOption struct {
	// GateName is the name of the quality gate (required).
	// Maximum length: 100 characters
	GateName string `url:"gateName,omitempty"`
	// Login is the user login (required).
	Login string `url:"login,omitempty"`
}

// QualitygatesRenameOption contains options for renaming a quality gate.
type QualitygatesRenameOption struct {
	// CurrentName is the current name of the quality gate (required).
	// Maximum length: 100 characters
	CurrentName string `url:"currentName,omitempty"`
	// Name is the new name of the quality gate (required).
	// Maximum length: 100 characters
	Name string `url:"name,omitempty"`
}

// QualitygatesSearchOption contains options for searching projects.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesSearchOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// GateName is the name of the quality gate (required).
	// Maximum length: 100 characters
	GateName string `url:"gateName,omitempty"`
	// Query is the search query for projects (optional).
	// If set, "selected" is set to "all".
	Query string `url:"query,omitempty"`
	// Selected filters by selection status (optional, default: selected).
	// Allowed values: all, deselected, selected
	Selected string `url:"selected,omitempty"`
}

// QualitygatesSearchGroupsOption contains options for searching groups.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesSearchGroupsOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// GateName is the name of the quality gate (required).
	GateName string `url:"gateName,omitempty"`
	// Query limits search to group names containing this string (optional).
	Query string `url:"q,omitempty"`
	// Selected filters by selection status (optional, default: selected).
	// Allowed values: all, deselected, selected
	Selected string `url:"selected,omitempty"`
}

// QualitygatesSearchUsersOption contains options for searching users.
//
//nolint:govet // Field alignment is less important than logical grouping
type QualitygatesSearchUsersOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	// GateName is the name of the quality gate (required).
	GateName string `url:"gateName,omitempty"`
	// Query limits search to names or logins containing this string (optional).
	Query string `url:"q,omitempty"`
	// Selected filters by selection status (optional, default: selected).
	// Allowed values: all, deselected, selected
	Selected string `url:"selected,omitempty"`
}

// QualitygatesSelectOption contains options for associating a project.
type QualitygatesSelectOption struct {
	// GateName is the name of the quality gate (required).
	// Maximum length: 100 characters
	GateName string `url:"gateName,omitempty"`
	// ProjectKey is the project key (required).
	ProjectKey string `url:"projectKey,omitempty"`
}

// QualitygatesSetAsDefaultOption contains options for setting the default gate.
type QualitygatesSetAsDefaultOption struct {
	// Name is the name of the quality gate to set as default (required).
	// Maximum length: 100 characters
	Name string `url:"name,omitempty"`
}

// QualitygatesShowOption contains options for showing a quality gate.
type QualitygatesShowOption struct {
	// Name is the name of the quality gate (required).
	Name string `url:"name,omitempty"`
}

// QualitygatesUpdateConditionOption contains options for updating a condition.
type QualitygatesUpdateConditionOption struct {
	// Error is the condition error threshold (required).
	// Maximum length: 64 characters
	Error string `url:"error,omitempty"`
	// ID is the condition ID (required).
	ID string `url:"id,omitempty"`
	// Metric is the condition metric (required).
	// Only metrics of the following types are allowed: INT, MILLISEC, RATING, WORK_DUR, FLOAT, PERCENT, LEVEL.
	// Forbidden metrics: alert_status, security_hotspots, new_security_hotspots.
	Metric string `url:"metric,omitempty"`
	// Op is the condition operator (optional).
	// Allowed values: LT (is lower than), GT (is greater than)
	Op string `url:"op,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// AddGroup allows a group of users to edit a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - Edit right on the specified quality gate
func (s *QualitygatesService) AddGroup(opt *QualitygatesAddGroupOption) (resp *http.Response, err error) {
	err = s.ValidateAddGroupOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/add_group", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// AddUser allows a user to edit a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - Edit right on the specified quality gate
func (s *QualitygatesService) AddUser(opt *QualitygatesAddUserOption) (resp *http.Response, err error) {
	err = s.ValidateAddUserOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/add_user", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Copy copies a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) Copy(opt *QualitygatesCopyOption) (resp *http.Response, err error) {
	err = s.ValidateCopyOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/copy", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Create creates a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) Create(opt *QualitygatesCreateOption) (v *QualitygatesCreate, resp *http.Response, err error) {
	err = s.ValidateCreateOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/create", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesCreate)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// CreateCondition adds a new condition to a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) CreateCondition(opt *QualitygatesCreateConditionOption) (v *QualitygatesCreateCondition, resp *http.Response, err error) {
	err = s.ValidateCreateConditionOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/create_condition", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesCreateCondition)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// DeleteCondition deletes a condition from a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) DeleteCondition(opt *QualitygatesDeleteConditionOption) (resp *http.Response, err error) {
	err = s.ValidateDeleteConditionOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/delete_condition", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Deselect removes the association of a project from a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - 'Administer' rights on the project
func (s *QualitygatesService) Deselect(opt *QualitygatesDeselectOption) (resp *http.Response, err error) {
	err = s.ValidateDeselectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/deselect", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Destroy deletes a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) Destroy(opt *QualitygatesDestroyOption) (resp *http.Response, err error) {
	err = s.ValidateDestroyOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/destroy", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// GetByProject gets the quality gate of a project.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//   - 'Browse' on the specified project
func (s *QualitygatesService) GetByProject(opt *QualitygatesGetByProjectOption) (v *QualitygatesGetByProject, resp *http.Response, err error) {
	err = s.ValidateGetByProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "qualitygates/get_by_project", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesGetByProject)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// List gets a list of quality gates.
func (s *QualitygatesService) List() (v *QualitygatesList, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "qualitygates/list", nil)
	if err != nil {
		return
	}

	v = new(QualitygatesList)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ProjectStatus gets the quality gate status of a project or a Compute Engine task.
// Either AnalysisID, ProjectID, or ProjectKey must be provided.
// Returns status: OK, WARN, ERROR, NONE.
// NONE is returned when there is no quality gate associated with the analysis.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//   - 'Browse' on the specified project
//   - 'Execute Analysis' on the specified project
func (s *QualitygatesService) ProjectStatus(opt *QualitygatesProjectStatusOption) (v *QualitygatesProjectStatus, resp *http.Response, err error) {
	err = s.ValidateProjectStatusOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "qualitygates/project_status", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesProjectStatus)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// RemoveGroup removes the ability from a group to edit a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - Edit right on the specified quality gate
func (s *QualitygatesService) RemoveGroup(opt *QualitygatesRemoveGroupOption) (resp *http.Response, err error) {
	err = s.ValidateRemoveGroupOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/remove_group", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// RemoveUser removes the ability from a user to edit a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - Edit right on the specified quality gate
func (s *QualitygatesService) RemoveUser(opt *QualitygatesRemoveUserOption) (resp *http.Response, err error) {
	err = s.ValidateRemoveUserOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/remove_user", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Rename renames a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) Rename(opt *QualitygatesRenameOption) (resp *http.Response, err error) {
	err = s.ValidateRenameOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/rename", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Search searches for projects associated (or not) to a quality gate.
// Only authorized projects for the current user will be returned.
func (s *QualitygatesService) Search(opt *QualitygatesSearchOption) (v *QualitygatesSearch, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "qualitygates/search", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesSearch)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchGroups lists the groups that are allowed to edit a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - Edit right on the specified quality gate
func (s *QualitygatesService) SearchGroups(opt *QualitygatesSearchGroupsOption) (v *QualitygatesSearchGroups, resp *http.Response, err error) {
	err = s.ValidateSearchGroupsOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "qualitygates/search_groups", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesSearchGroups)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// SearchUsers lists the users that are allowed to edit a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - Edit right on the specified quality gate
func (s *QualitygatesService) SearchUsers(opt *QualitygatesSearchUsersOption) (v *QualitygatesSearchUsers, resp *http.Response, err error) {
	err = s.ValidateSearchUsersOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "qualitygates/search_users", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesSearchUsers)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// Select associates a project to a quality gate.
// Requires one of the following permissions:
//   - 'Administer Quality Gates'
//   - 'Administer' right on the specified project
func (s *QualitygatesService) Select(opt *QualitygatesSelectOption) (resp *http.Response, err error) {
	err = s.ValidateSelectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/select", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// SetAsDefault sets a quality gate as the default quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) SetAsDefault(opt *QualitygatesSetAsDefaultOption) (resp *http.Response, err error) {
	err = s.ValidateSetAsDefaultOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/set_as_default", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// Show displays the details of a quality gate.
func (s *QualitygatesService) Show(opt *QualitygatesShowOption) (v *QualitygatesShow, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "qualitygates/show", opt)
	if err != nil {
		return
	}

	v = new(QualitygatesShow)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}

// UpdateCondition updates a condition attached to a quality gate.
// Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) UpdateCondition(opt *QualitygatesUpdateConditionOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateConditionOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("POST", "qualitygates/update_condition", opt)
	if err != nil {
		return
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}

	return
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateAddGroupOpt validates the options for adding a group to a quality gate.
func (s *QualitygatesService) ValidateAddGroupOpt(opt *QualitygatesAddGroupOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesAddGroupOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.GateName, MaxQualityGateNameLength, "GateName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.GroupName, "GroupName")
	if err != nil {
		return err
	}

	return nil
}

// ValidateAddUserOpt validates the options for adding a user to a quality gate.
func (s *QualitygatesService) ValidateAddUserOpt(opt *QualitygatesAddUserOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesAddUserOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.GateName, MaxQualityGateNameLength, "GateName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCopyOpt validates the options for copying a quality gate.
func (s *QualitygatesService) ValidateCopyOpt(opt *QualitygatesCopyOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesCopyOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.SourceName, "SourceName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.SourceName, MaxQualityGateNameLength, "SourceName")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateOpt validates the options for creating a quality gate.
func (s *QualitygatesService) ValidateCreateOpt(opt *QualitygatesCreateOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesCreateOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxQualityGateNameLength, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateConditionOpt validates the options for creating a condition.
func (s *QualitygatesService) ValidateCreateConditionOpt(opt *QualitygatesCreateConditionOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesCreateConditionOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Error, "Error")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Error, MaxConditionErrorLength, "Error")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Metric, "Metric")
	if err != nil {
		return err
	}

	// Validate operator if provided
	if opt.Op != "" {
		allowed := []string{"LT", "GT"}

		err = ValidateInSlice(opt.Op, allowed, "Op")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteConditionOpt validates the options for deleting a condition.
func (s *QualitygatesService) ValidateDeleteConditionOpt(opt *QualitygatesDeleteConditionOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesDeleteConditionOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ID, "ID")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeselectOpt validates the options for deselecting a project.
func (s *QualitygatesService) ValidateDeselectOpt(opt *QualitygatesDeselectOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesDeselectOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDestroyOpt validates the options for destroying a quality gate.
func (s *QualitygatesService) ValidateDestroyOpt(opt *QualitygatesDestroyOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesDestroyOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxQualityGateNameLength, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetByProjectOpt validates the options for getting a project's quality gate.
func (s *QualitygatesService) ValidateGetByProjectOpt(opt *QualitygatesGetByProjectOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesGetByProjectOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateProjectStatusOpt validates the options for getting project status.
func (s *QualitygatesService) ValidateProjectStatusOpt(opt *QualitygatesProjectStatusOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesProjectStatusOption", "cannot be nil", ErrMissingRequired)
	}

	// At least one of AnalysisID, ProjectID, or ProjectKey must be provided
	if opt.AnalysisID == "" && opt.ProjectID == "" && opt.ProjectKey == "" {
		return NewValidationError(
			"QualitygatesProjectStatusOption",
			"at least one of AnalysisID, ProjectID, or ProjectKey must be provided",
			ErrMissingRequired,
		)
	}

	return nil
}

// ValidateRemoveGroupOpt validates the options for removing a group from a quality gate.
func (s *QualitygatesService) ValidateRemoveGroupOpt(opt *QualitygatesRemoveGroupOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesRemoveGroupOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.GateName, MaxQualityGateNameLength, "GateName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.GroupName, "GroupName")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRemoveUserOpt validates the options for removing a user from a quality gate.
func (s *QualitygatesService) ValidateRemoveUserOpt(opt *QualitygatesRemoveUserOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesRemoveUserOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.GateName, MaxQualityGateNameLength, "GateName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRenameOpt validates the options for renaming a quality gate.
func (s *QualitygatesService) ValidateRenameOpt(opt *QualitygatesRenameOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesRenameOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.CurrentName, "CurrentName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.CurrentName, MaxQualityGateNameLength, "CurrentName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxQualityGateNameLength, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for searching projects.
func (s *QualitygatesService) ValidateSearchOpt(opt *QualitygatesSearchOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesSearchOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.GateName, MaxQualityGateNameLength, "GateName")
	if err != nil {
		return err
	}

	// Validate Selected if provided
	if opt.Selected != "" {
		allowed := []string{"all", "deselected", "selected"}

		err = ValidateInSlice(opt.Selected, allowed, "Selected")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchGroupsOpt validates the options for searching groups.
func (s *QualitygatesService) ValidateSearchGroupsOpt(opt *QualitygatesSearchGroupsOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesSearchGroupsOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	// Validate Selected if provided
	if opt.Selected != "" {
		allowed := []string{"all", "deselected", "selected"}

		err = ValidateInSlice(opt.Selected, allowed, "Selected")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchUsersOpt validates the options for searching users.
func (s *QualitygatesService) ValidateSearchUsersOpt(opt *QualitygatesSearchUsersOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesSearchUsersOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	// Validate Selected if provided
	if opt.Selected != "" {
		allowed := []string{"all", "deselected", "selected"}

		err = ValidateInSlice(opt.Selected, allowed, "Selected")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSelectOpt validates the options for selecting a project.
func (s *QualitygatesService) ValidateSelectOpt(opt *QualitygatesSelectOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesSelectOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.GateName, MaxQualityGateNameLength, "GateName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetAsDefaultOpt validates the options for setting a default quality gate.
func (s *QualitygatesService) ValidateSetAsDefaultOpt(opt *QualitygatesSetAsDefaultOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesSetAsDefaultOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxQualityGateNameLength, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateShowOpt validates the options for showing a quality gate.
func (s *QualitygatesService) ValidateShowOpt(opt *QualitygatesShowOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesShowOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateConditionOpt validates the options for updating a condition.
func (s *QualitygatesService) ValidateUpdateConditionOpt(opt *QualitygatesUpdateConditionOption) error {
	if opt == nil {
		return NewValidationError("QualitygatesUpdateConditionOption", "cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Error, "Error")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Error, MaxConditionErrorLength, "Error")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ID, "ID")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Metric, "Metric")
	if err != nil {
		return err
	}

	// Validate operator if provided
	if opt.Op != "" {
		allowed := []string{"LT", "GT"}

		err = ValidateInSlice(opt.Op, allowed, "Op")
		if err != nil {
			return err
		}
	}

	return nil
}
