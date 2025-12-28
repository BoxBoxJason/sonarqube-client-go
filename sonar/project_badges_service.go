// Generate badges based on quality gates or measures
package sonargo

import "net/http"

type ProjectBadgesService struct {
	client *Client
}

type ProjectBadgesTokenObject struct {
	Token string `json:"token,omitempty"`
}

type ProjectBadgesMeasureOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Branch key",ExampleValue:"feature/my_branch"
	Metric  string `url:"metric,omitempty"`  // Description:"Metric key",ExampleValue:""
	Project string `url:"project,omitempty"` // Description:"Project or application key",ExampleValue:"my_project"
	Token   string `url:"token,omitempty"`   // Description:"Project badge token",ExampleValue:"8bb493196edb5896ccb64582499895f187a2ae8f"
}

// Measure Generate badge for project's measure as an SVG.<br/>Requires 'Browse' permission on the specified project.
func (s *ProjectBadgesService) Measure(opt *ProjectBadgesMeasureOption) (v *string, resp *http.Response, err error) {
	err = s.ValidateMeasureOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "project_badges/measure", opt)
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

type ProjectBadgesQualityGateOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Branch key",ExampleValue:"feature/my_branch"
	Project string `url:"project,omitempty"` // Description:"Project or application key",ExampleValue:"my_project"
	Token   string `url:"token,omitempty"`   // Description:"Project badge token",ExampleValue:"8bb493196edb5896ccb64582499895f187a2ae8f"
}

// QualityGate Generate badge for project's quality gate as an SVG.<br/>Requires 'Browse' permission on the specified project.
func (s *ProjectBadgesService) QualityGate(opt *ProjectBadgesQualityGateOption) (v *string, resp *http.Response, err error) {
	err = s.ValidateQualityGateOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "project_badges/quality_gate", opt)
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

type ProjectBadgesRenewTokenOption struct {
	Project string `url:"project,omitempty"` // Description:"Project or application key",ExampleValue:"my_project"
}

// RenewToken Creates new token replacing any existing token for project or application badge access for private projects and applications.<br/>This token can be used to authenticate with api/project_badges/quality_gate and api/project_badges/measure endpoints.<br/>Requires 'Administer' permission on the specified project or application.
func (s *ProjectBadgesService) RenewToken(opt *ProjectBadgesRenewTokenOption) (resp *http.Response, err error) {
	err = s.ValidateRenewTokenOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "project_badges/renew_token", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type ProjectBadgesTokenOption struct {
	Project string `url:"project,omitempty"` // Description:"Project or application key",ExampleValue:"my_project"
}

// Token Retrieve a token to use for project or application badge access for private projects or applications.<br/>This token can be used to authenticate with api/project_badges/quality_gate and api/project_badges/measure endpoints.<br/>Requires 'Browse' permission on the specified project or application.
func (s *ProjectBadgesService) Token(opt *ProjectBadgesTokenOption) (v *ProjectBadgesTokenObject, resp *http.Response, err error) {
	err = s.ValidateTokenOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "project_badges/token", opt)
	if err != nil {
		return
	}
	v = new(ProjectBadgesTokenObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
