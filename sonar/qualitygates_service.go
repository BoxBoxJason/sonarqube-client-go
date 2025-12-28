// Manage quality gates, including conditions and project association.
package sonargo

import "net/http"

type QualitygatesService struct {
	client *Client
}

type QualitygatesCreateObject struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type QualitygatesCreateConditionObject struct {
	Error   string `json:"error,omitempty"`
	ID      string `json:"id,omitempty"`
	Metric  string `json:"metric,omitempty"`
	Op      string `json:"op,omitempty"`
	Warning string `json:"warning,omitempty"`
}

type QualitygatesGetByProjectObject struct {
	QualityGate QualitygatesGetByProjectObject_sub1 `json:"qualityGate,omitempty"`
}

type QualitygatesGetByProjectObject_sub1 struct {
	Default bool   `json:"default,omitempty"`
	Name    string `json:"name,omitempty"`
}

type QualitygatesListObject struct {
	Actions      QualitygatesListObject_sub1   `json:"actions,omitempty"`
	Qualitygates []QualitygatesListObject_sub3 `json:"qualitygates,omitempty"`
}

type QualitygatesListObject_sub3 struct {
	Actions               QualitygatesListObject_sub2 `json:"actions,omitempty"`
	CaycStatus            string                      `json:"caycStatus,omitempty"`
	HasMQRConditions      bool                        `json:"hasMQRConditions,omitempty"`
	HasStandardConditions bool                        `json:"hasStandardConditions,omitempty"`
	IsAiCodeSupported     bool                        `json:"isAiCodeSupported,omitempty"`
	IsBuiltIn             bool                        `json:"isBuiltIn,omitempty"`
	IsDefault             bool                        `json:"isDefault,omitempty"`
	Name                  string                      `json:"name,omitempty"`
}

type QualitygatesListObject_sub2 struct {
	AssociateProjects     bool `json:"associateProjects,omitempty"`
	Copy                  bool `json:"copy,omitempty"`
	Delegate              bool `json:"delegate,omitempty"`
	Delete                bool `json:"delete,omitempty"`
	ManageAiCodeAssurance bool `json:"manageAiCodeAssurance,omitempty"`
	ManageConditions      bool `json:"manageConditions,omitempty"`
	Rename                bool `json:"rename,omitempty"`
	SetAsDefault          bool `json:"setAsDefault,omitempty"`
}

type QualitygatesListObject_sub1 struct {
	Create bool `json:"create,omitempty"`
}

type QualitygatesProjectStatusObject struct {
	ProjectStatus QualitygatesProjectStatusObject_sub3 `json:"projectStatus,omitempty"`
}

type QualitygatesProjectStatusObject_sub1 struct {
	ActualValue    string `json:"actualValue,omitempty"`
	Comparator     string `json:"comparator,omitempty"`
	ErrorThreshold string `json:"errorThreshold,omitempty"`
	MetricKey      string `json:"metricKey,omitempty"`
	Status         string `json:"status,omitempty"`
}

type QualitygatesProjectStatusObject_sub3 struct {
	CaycStatus        string                                 `json:"caycStatus,omitempty"`
	Conditions        []QualitygatesProjectStatusObject_sub1 `json:"conditions,omitempty"`
	IgnoredConditions bool                                   `json:"ignoredConditions,omitempty"`
	Period            QualitygatesProjectStatusObject_sub2   `json:"period,omitempty"`
	Status            string                                 `json:"status,omitempty"`
}

type QualitygatesProjectStatusObject_sub2 struct {
	Date      string `json:"date,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type QualitygatesSearchObject struct {
	Paging  QualitygatesSearchObject_sub1   `json:"paging,omitempty"`
	Results []QualitygatesSearchObject_sub2 `json:"results,omitempty"`
}

type QualitygatesSearchObject_sub2 struct {
	AiCodeAssurance string `json:"aiCodeAssurance,omitempty"`
	Key             string `json:"key,omitempty"`
	Name            string `json:"name,omitempty"`
	Selected        bool   `json:"selected,omitempty"`
}

type QualitygatesSearchObject_sub1 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type QualitygatesSearchGroupsObject struct {
	Groups []QualitygatesSearchGroupsObject_sub1 `json:"groups,omitempty"`
	Paging QualitygatesSearchGroupsObject_sub2   `json:"paging,omitempty"`
}

type QualitygatesSearchGroupsObject_sub1 struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Selected    bool   `json:"selected,omitempty"`
}

type QualitygatesSearchGroupsObject_sub2 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type QualitygatesSearchUsersObject struct {
	Paging QualitygatesSearchUsersObject_sub1   `json:"paging,omitempty"`
	Users  []QualitygatesSearchUsersObject_sub2 `json:"users,omitempty"`
}

type QualitygatesSearchUsersObject_sub2 struct {
	Avatar   string `json:"avatar,omitempty"`
	Login    string `json:"login,omitempty"`
	Name     string `json:"name,omitempty"`
	Selected bool   `json:"selected,omitempty"`
}

type QualitygatesSearchUsersObject_sub1 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type QualitygatesShowObject struct {
	Actions           QualitygatesShowObject_sub1   `json:"actions,omitempty"`
	CaycStatus        string                        `json:"caycStatus,omitempty"`
	Conditions        []QualitygatesShowObject_sub2 `json:"conditions,omitempty"`
	IsAiCodeSupported bool                          `json:"isAiCodeSupported,omitempty"`
	IsBuiltIn         bool                          `json:"isBuiltIn,omitempty"`
	IsDefault         bool                          `json:"isDefault,omitempty"`
	Name              string                        `json:"name,omitempty"`
}

type QualitygatesShowObject_sub1 struct {
	AssociateProjects     bool `json:"associateProjects,omitempty"`
	Copy                  bool `json:"copy,omitempty"`
	Delegate              bool `json:"delegate,omitempty"`
	Delete                bool `json:"delete,omitempty"`
	ManageAiCodeAssurance bool `json:"manageAiCodeAssurance,omitempty"`
	ManageConditions      bool `json:"manageConditions,omitempty"`
	Rename                bool `json:"rename,omitempty"`
	SetAsDefault          bool `json:"setAsDefault,omitempty"`
}

type QualitygatesShowObject_sub2 struct {
	Error  string `json:"error,omitempty"`
	ID     string `json:"id,omitempty"`
	Metric string `json:"metric,omitempty"`
	Op     string `json:"op,omitempty"`
}

type QualitygatesAddGroupOption struct {
	GateName  string `url:"gateName,omitempty"`  // Description:"Quality Gate name",ExampleValue:"SonarSource Way"
	GroupName string `url:"groupName,omitempty"` // Description:"Group name or 'anyone' (case insensitive)",ExampleValue:"sonar-administrators"
}

// AddGroup Allow a group of users to edit a Quality Gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>Edit right on the specified quality gate</li></ul>
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

type QualitygatesAddUserOption struct {
	GateName string `url:"gateName,omitempty"` // Description:"Quality Gate name",ExampleValue:"SonarSource Way"
	Login    string `url:"login,omitempty"`    // Description:"User login",ExampleValue:"john.doe"
}

// AddUser Allow a user to edit a Quality Gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>Edit right on the specified quality gate</li></ul>
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

type QualitygatesCopyOption struct {
	Name       string `url:"name,omitempty"`       // Description:"The name of the quality gate to create",ExampleValue:"My New Quality Gate"
	SourceName string `url:"sourceName,omitempty"` // Description:"The name of the quality gate to copy",ExampleValue:"My Quality Gate"
}

// Copy Copy a Quality Gate.<br>'sourceName' must be provided. Requires the 'Administer Quality Gates' permission.
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

type QualitygatesCreateOption struct {
	Name string `url:"name,omitempty"` // Description:"The name of the quality gate to create",ExampleValue:"My Quality Gate"
}

// Create Create a Quality Gate.<br>Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) Create(opt *QualitygatesCreateOption) (v *QualitygatesCreateObject, resp *http.Response, err error) {
	err = s.ValidateCreateOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "qualitygates/create", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesCreateObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesCreateConditionOption struct {
	Error    string `url:"error,omitempty"`    // Description:"Condition error threshold",ExampleValue:"10"
	GateName string `url:"gateName,omitempty"` // Description:"Name of the quality gate",ExampleValue:"SonarSource way"
	Metric   string `url:"metric,omitempty"`   // Description:"Condition metric.<br/> Only metric of the following types are allowed:<ul><li>INT</li><li>MILLISEC</li><li>RATING</li><li>WORK_DUR</li><li>FLOAT</li><li>PERCENT</li><li>LEVEL</li></ul>Following metrics are forbidden:<ul><li>alert_status</li><li>security_hotspots</li><li>new_security_hotspots</li></ul>",ExampleValue:"blocker_violations, vulnerabilities, new_code_smells"
	Op       string `url:"op,omitempty"`       // Description:"Condition operator:<br/><ul><li>LT = is lower than</li><li>GT = is greater than</li></ul>",ExampleValue:"GT"
}

// CreateCondition Add a new condition to a quality gate.<br>Parameter 'gateName' must be provided. Requires the 'Administer Quality Gates' permission.
func (s *QualitygatesService) CreateCondition(opt *QualitygatesCreateConditionOption) (v *QualitygatesCreateConditionObject, resp *http.Response, err error) {
	err = s.ValidateCreateConditionOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "qualitygates/create_condition", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesCreateConditionObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesDeleteConditionOption struct {
	Id string `url:"id,omitempty"` // Description:"Condition UUID",ExampleValue:"2"
}

// DeleteCondition Delete a condition from a quality gate.<br>Requires the 'Administer Quality Gates' permission.
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

type QualitygatesDeselectOption struct {
	ProjectKey string `url:"projectKey,omitempty"` // Description:"Project key",ExampleValue:"my_project"
}

// Deselect Remove the association of a project from a quality gate.<br>Requires one of the following permissions:<ul><li>'Administer Quality Gates'</li><li>'Administer' rights on the project</li></ul>
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

type QualitygatesDestroyOption struct {
	Name string `url:"name,omitempty"` // Description:"Name of the quality gate to delete",ExampleValue:"SonarSource Way"
}

// Destroy Delete a Quality Gate.<br>Parameter 'name' must be specified. Requires the 'Administer Quality Gates' permission.
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

type QualitygatesGetByProjectOption struct {
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:"my_project"
}

// GetByProject Get the quality gate of a project.<br />Requires one of the following permissions:<ul><li>'Administer System'</li><li>'Administer' rights on the specified project</li><li>'Browse' on the specified project</li></ul>
func (s *QualitygatesService) GetByProject(opt *QualitygatesGetByProjectOption) (v *QualitygatesGetByProjectObject, resp *http.Response, err error) {
	err = s.ValidateGetByProjectOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "qualitygates/get_by_project", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesGetByProjectObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// List Get a list of quality gates
func (s *QualitygatesService) List() (v *QualitygatesListObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "qualitygates/list", nil)
	if err != nil {
		return
	}
	v = new(QualitygatesListObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesProjectStatusOption struct {
	AnalysisId  string `url:"analysisId,omitempty"`  // Description:"Analysis id",ExampleValue:"AU-TpxcA-iU5OvuD2FL1"
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	ProjectId   string `url:"projectId,omitempty"`   // Description:"Project UUID. Doesn't work with branches or pull requests",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
	ProjectKey  string `url:"projectKey,omitempty"`  // Description:"Project key",ExampleValue:"my_project"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
}

// ProjectStatus Get the quality gate status of a project or a Compute Engine task.<br />Either 'analysisId', 'projectId' or 'projectKey' must be provided <br />The different statuses returned are: OK, WARN, ERROR, NONE. The NONE status is returned when there is no quality gate associated with the analysis.<br />Returns an HTTP code 404 if the analysis associated with the task is not found or does not exist.<br />Requires one of the following permissions:<ul><li>'Administer System'</li><li>'Administer' rights on the specified project</li><li>'Browse' on the specified project</li><li>'Execute Analysis' on the specified project</li></ul>
func (s *QualitygatesService) ProjectStatus(opt *QualitygatesProjectStatusOption) (v *QualitygatesProjectStatusObject, resp *http.Response, err error) {
	err = s.ValidateProjectStatusOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "qualitygates/project_status", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesProjectStatusObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesRemoveGroupOption struct {
	GateName  string `url:"gateName,omitempty"`  // Description:"Quality Gate name",ExampleValue:"SonarSource Way"
	GroupName string `url:"groupName,omitempty"` // Description:"Group name or 'anyone' (case insensitive)",ExampleValue:"sonar-administrators"
}

// RemoveGroup Remove the ability from a group to edit a Quality Gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>Edit right on the specified quality gate</li></ul>
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

type QualitygatesRemoveUserOption struct {
	GateName string `url:"gateName,omitempty"` // Description:"Quality Gate name",ExampleValue:"SonarSource Way"
	Login    string `url:"login,omitempty"`    // Description:"User login",ExampleValue:"john.doe"
}

// RemoveUser Remove the ability from an user to edit a Quality Gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>Edit right on the specified quality gate</li></ul>
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

type QualitygatesRenameOption struct {
	CurrentName string `url:"currentName,omitempty"` // Description:"Current name of the quality gate",ExampleValue:"My Quality Gate"
	Name        string `url:"name,omitempty"`        // Description:"New name of the quality gate",ExampleValue:"My New Quality Gate"
}

// Rename Rename a Quality Gate.<br>'currentName' must be specified. Requires the 'Administer Quality Gates' permission.
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

type QualitygatesSearchOption struct {
	GateName string `url:"gateName,omitempty"` // Description:"Quality Gate name",ExampleValue:"SonarSource Way"
	Page     string `url:"page,omitempty"`     // Description:"Page number",ExampleValue:"2"
	PageSize string `url:"pageSize,omitempty"` // Description:"Page size",ExampleValue:"10"
	Query    string `url:"query,omitempty"`    // Description:"To search for projects containing this string. If this parameter is set, "selected" is set to "all".",ExampleValue:"abc"
	Selected string `url:"selected,omitempty"` // Description:"Depending on the value, show only selected items (selected=selected), deselected items (selected=deselected), or all items with their selection status (selected=all).",ExampleValue:""
}

// Search Search for projects associated (or not) to a quality gate.<br/>Only authorized projects for the current user will be returned.
func (s *QualitygatesService) Search(opt *QualitygatesSearchOption) (v *QualitygatesSearchObject, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "qualitygates/search", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesSearchObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesSearchGroupsOption struct {
	GateName string `url:"gateName,omitempty"` // Description:"Quality Gate name",ExampleValue:"SonarSource Way"
	P        string `url:"p,omitempty"`        // Description:"1-based page number",ExampleValue:"42"
	Ps       string `url:"ps,omitempty"`       // Description:"Page size. Must be greater than 0.",ExampleValue:"20"
	Q        string `url:"q,omitempty"`        // Description:"Limit search to group names that contain the supplied string.",ExampleValue:"sonar"
	Selected string `url:"selected,omitempty"` // Description:"Depending on the value, show only selected items (selected=selected), deselected items (selected=deselected), or all items with their selection status (selected=all).",ExampleValue:""
}

// SearchGroups List the groups that are allowed to edit a Quality Gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>Edit right on the specified quality gate</li></ul>
func (s *QualitygatesService) SearchGroups(opt *QualitygatesSearchGroupsOption) (v *QualitygatesSearchGroupsObject, resp *http.Response, err error) {
	err = s.ValidateSearchGroupsOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "qualitygates/search_groups", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesSearchGroupsObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesSearchUsersOption struct {
	GateName string `url:"gateName,omitempty"` // Description:"Quality Gate name",ExampleValue:"Recommended quality gate"
	P        string `url:"p,omitempty"`        // Description:"1-based page number",ExampleValue:"42"
	Ps       string `url:"ps,omitempty"`       // Description:"Page size. Must be greater than 0.",ExampleValue:"20"
	Q        string `url:"q,omitempty"`        // Description:"Limit search to names or logins that contain the supplied string.",ExampleValue:"freddy"
	Selected string `url:"selected,omitempty"` // Description:"Depending on the value, show only selected items (selected=selected), deselected items (selected=deselected), or all items with their selection status (selected=all).",ExampleValue:""
}

// SearchUsers List the users that are allowed to edit a Quality Gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>Edit right on the specified quality gate</li></ul>
func (s *QualitygatesService) SearchUsers(opt *QualitygatesSearchUsersOption) (v *QualitygatesSearchUsersObject, resp *http.Response, err error) {
	err = s.ValidateSearchUsersOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "qualitygates/search_users", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesSearchUsersObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesSelectOption struct {
	GateName   string `url:"gateName,omitempty"`   // Description:"Name of the quality gate",ExampleValue:"SonarSource way"
	ProjectKey string `url:"projectKey,omitempty"` // Description:"Project key",ExampleValue:"my_project"
}

// Select Associate a project to a quality gate.<br>Requires one of the following permissions:<ul>  <li>'Administer Quality Gates'</li>  <li>'Administer' right on the specified project</li></ul>
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

type QualitygatesSetAsDefaultOption struct {
	Name string `url:"name,omitempty"` // Description:"Name of the quality gate to set as default",ExampleValue:"SonarSource Way"
}

// SetAsDefault Set a quality gate as the default quality gate.<br>Parameter 'name' must be specified. Requires the 'Administer Quality Gates' permission.
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

type QualitygatesShowOption struct {
	Name string `url:"name,omitempty"` // Description:"Name of the quality gate. Either id or name must be set",ExampleValue:"My Quality Gate"
}

// Show Display the details of a quality gate
func (s *QualitygatesService) Show(opt *QualitygatesShowOption) (v *QualitygatesShowObject, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "qualitygates/show", opt)
	if err != nil {
		return
	}
	v = new(QualitygatesShowObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type QualitygatesUpdateConditionOption struct {
	Error  string `url:"error,omitempty"`  // Description:"Condition error threshold",ExampleValue:"10"
	Id     string `url:"id,omitempty"`     // Description:"Condition ID",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
	Metric string `url:"metric,omitempty"` // Description:"Condition metric.<br/> Only metric of the following types are allowed:<ul><li>INT</li><li>MILLISEC</li><li>RATING</li><li>WORK_DUR</li><li>FLOAT</li><li>PERCENT</li><li>LEVEL</li></ul>Following metrics are forbidden:<ul><li>alert_status</li><li>security_hotspots</li><li>new_security_hotspots</li></ul>",ExampleValue:"blocker_violations, vulnerabilities, new_code_smells"
	Op     string `url:"op,omitempty"`     // Description:"Condition operator:<br/><ul><li>LT = is lower than</li><li>GT = is greater than</li></ul>",ExampleValue:"GT"
}

// UpdateCondition Update a condition attached to a quality gate.<br>Requires the 'Administer Quality Gates' permission.
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
