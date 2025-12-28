// Manage <a href="https://docs.sonarsource.com/sonarqube-community-build/project-administration/configuring-new-code-calculation#setting-specific-new-code-definition-for-project" target="_blank" rel="noopener noreferrer">new code definition</a>.
package sonargo

import "net/http"

type NewCodePeriodsService struct {
	client *Client
}

type NewCodePeriodsListObject struct {
	NewCodePeriods []NewCodePeriodsListObject_sub1 `json:"newCodePeriods,omitempty"`
}

type NewCodePeriodsListObject_sub1 struct {
	BranchKey      string `json:"branchKey,omitempty"`
	EffectiveValue string `json:"effectiveValue,omitempty"`
	Inherited      bool   `json:"inherited,omitempty"`
	ProjectKey     string `json:"projectKey,omitempty"`
	Type           string `json:"type,omitempty"`
	Value          string `json:"value,omitempty"`
}

type NewCodePeriodsShowObject struct {
	BranchKey  string `json:"branchKey,omitempty"`
	Inherited  bool   `json:"inherited,omitempty"`
	ProjectKey string `json:"projectKey,omitempty"`
	Type       string `json:"type,omitempty"`
}

type NewCodePeriodsListOption struct {
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:""
}

// List Lists the <a href="https://docs.sonarsource.com/sonarqube-community-build/project-administration/configuring-new-code-calculation#setting-specific-new-code-definition-for-project" target="_blank" rel="noopener noreferrer">new code definition</a> for all branches in a project.<br>Requires the permission to browse the project
func (s *NewCodePeriodsService) List(opt *NewCodePeriodsListOption) (v *NewCodePeriodsListObject, resp *http.Response, err error) {
	err = s.ValidateListOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "new_code_periods/list", opt)
	if err != nil {
		return
	}
	v = new(NewCodePeriodsListObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type NewCodePeriodsSetOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Branch key",ExampleValue:""
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:""
	Type    string `url:"type,omitempty"`    // Description:"Type<br/>New code definitions of the following types are allowed:<ul><li>SPECIFIC_ANALYSIS - can be set at branch level only</li><li>PREVIOUS_VERSION - can be set at any level (global, project, branch)</li><li>NUMBER_OF_DAYS - can be set at any level (global, project, branch)</li><li>REFERENCE_BRANCH - can only be set for projects and branches</li></ul>",ExampleValue:""
	Value   string `url:"value,omitempty"`   // Description:"Value<br/>For each type, a different value is expected:<ul><li>the uuid of an analysis, when type is SPECIFIC_ANALYSIS</li><li>no value, when type is PREVIOUS_VERSION</li><li>a number between 1 and 90, when type is NUMBER_OF_DAYS</li><li>a string, when type is REFERENCE_BRANCH</li></ul>",ExampleValue:""
}

// Set Updates the <a href="https://docs.sonarsource.com/sonarqube-community-build/project-administration/configuring-new-code-calculation#setting-specific-new-code-definition-for-project" target="_blank" rel="noopener noreferrer">new code definition</a> on different levels:<br><ul><li>Not providing a project key and a branch key will update the default value at global level. Existing projects or branches having a specific new code definition will not be impacted</li><li>Project key must be provided to update the value for a project</li><li>Both project and branch keys must be provided to update the value for a branch</li></ul>Requires one of the following permissions: <ul><li>'Administer System' to change the global setting</li><li>'Administer' rights on the specified project to change the project setting</li></ul>
func (s *NewCodePeriodsService) Set(opt *NewCodePeriodsSetOption) (resp *http.Response, err error) {
	err = s.ValidateSetOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "new_code_periods/set", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type NewCodePeriodsShowOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Branch key",ExampleValue:""
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:""
}

// Show Shows the <a href="https://docs.sonarsource.com/sonarqube-community-build/project-administration/configuring-new-code-calculation#setting-specific-new-code-definition-for-project" target="_blank" rel="noopener noreferrer">new code definition</a>.<br> If the component requested doesn't exist or if no new code definition is set for it, a value is inherited from the project or from the global setting.Requires one of the following permissions if a component is specified: <ul><li>'Administer' rights on the specified component</li><li>'Execute analysis' rights on the specified component</li></ul>
func (s *NewCodePeriodsService) Show(opt *NewCodePeriodsShowOption) (v *NewCodePeriodsShowObject, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "new_code_periods/show", opt)
	if err != nil {
		return
	}
	v = new(NewCodePeriodsShowObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type NewCodePeriodsUnsetOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Branch key",ExampleValue:""
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:""
}

// Unset Unsets the <a href="https://docs.sonarsource.com/sonarqube-community-build/project-administration/configuring-new-code-calculation#setting-specific-new-code-definition-for-project" target="_blank" rel="noopener noreferrer">new code definition</a> for a branch, project or global. It requires the inherited New Code Definition to be compatible with the Clean as You Code methodology, and one of the following permissions: <ul><li>'Administer System' to change the global setting</li><li>'Administer' rights for a specified component</li></ul>
func (s *NewCodePeriodsService) Unset(opt *NewCodePeriodsUnsetOption) (resp *http.Response, err error) {
	err = s.ValidateUnsetOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "new_code_periods/unset", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}
