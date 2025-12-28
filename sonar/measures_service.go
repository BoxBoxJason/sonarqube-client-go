// Get components or children with specified measures.
package sonargo

import "net/http"

type MeasuresService struct {
	client *Client
}

type MeasuresComponentObject struct {
	Component MeasuresComponentObject_sub3   `json:"component,omitempty"`
	Metrics   []MeasuresComponentObject_sub4 `json:"metrics,omitempty"`
	Period    MeasuresComponentObject_sub5   `json:"period,omitempty"`
}

type MeasuresComponentObject_sub1 struct {
	BestValue bool   `json:"bestValue,omitempty"`
	Value     string `json:"value,omitempty"`
}

type MeasuresComponentObject_sub5 struct {
	Date      string `json:"date,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type MeasuresComponentObject_sub4 struct {
	Description           string `json:"description,omitempty"`
	Domain                string `json:"domain,omitempty"`
	Hidden                bool   `json:"hidden,omitempty"`
	HigherValuesAreBetter bool   `json:"higherValuesAreBetter,omitempty"`
	Key                   string `json:"key,omitempty"`
	Name                  string `json:"name,omitempty"`
	Qualitative           bool   `json:"qualitative,omitempty"`
	Type                  string `json:"type,omitempty"`
}

type MeasuresComponentObject_sub3 struct {
	Key       string                         `json:"key,omitempty"`
	Language  string                         `json:"language,omitempty"`
	Measures  []MeasuresComponentObject_sub2 `json:"measures,omitempty"`
	Name      string                         `json:"name,omitempty"`
	Path      string                         `json:"path,omitempty"`
	Qualifier string                         `json:"qualifier,omitempty"`
}

type MeasuresComponentObject_sub2 struct {
	Metric string                       `json:"metric,omitempty"`
	Period MeasuresComponentObject_sub1 `json:"period,omitempty"`
	Value  string                       `json:"value,omitempty"`
}

type MeasuresComponentTreeObject struct {
	BaseComponent MeasuresComponentTreeObject_sub3   `json:"baseComponent,omitempty"`
	Components    []MeasuresComponentTreeObject_sub4 `json:"components,omitempty"`
	Metrics       []MeasuresComponentTreeObject_sub5 `json:"metrics,omitempty"`
	Paging        MeasuresComponentTreeObject_sub6   `json:"paging,omitempty"`
	Period        MeasuresComponentTreeObject_sub7   `json:"period,omitempty"`
}

type MeasuresComponentTreeObject_sub5 struct {
	BestValue             string `json:"bestValue,omitempty"`
	Description           string `json:"description,omitempty"`
	Domain                string `json:"domain,omitempty"`
	Hidden                bool   `json:"hidden,omitempty"`
	HigherValuesAreBetter bool   `json:"higherValuesAreBetter,omitempty"`
	Key                   string `json:"key,omitempty"`
	Name                  string `json:"name,omitempty"`
	Qualitative           bool   `json:"qualitative,omitempty"`
	Type                  string `json:"type,omitempty"`
}

type MeasuresComponentTreeObject_sub7 struct {
	Date      string `json:"date,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type MeasuresComponentTreeObject_sub4 struct {
	Key       string                             `json:"key,omitempty"`
	Language  string                             `json:"language,omitempty"`
	Measures  []MeasuresComponentTreeObject_sub2 `json:"measures,omitempty"`
	Name      string                             `json:"name,omitempty"`
	Path      string                             `json:"path,omitempty"`
	Qualifier string                             `json:"qualifier,omitempty"`
}

type MeasuresComponentTreeObject_sub3 struct {
	Key       string                             `json:"key,omitempty"`
	Measures  []MeasuresComponentTreeObject_sub2 `json:"measures,omitempty"`
	Name      string                             `json:"name,omitempty"`
	Qualifier string                             `json:"qualifier,omitempty"`
}

type MeasuresComponentTreeObject_sub2 struct {
	Metric string                           `json:"metric,omitempty"`
	Period MeasuresComponentTreeObject_sub1 `json:"period,omitempty"`
	Value  string                           `json:"value,omitempty"`
}

type MeasuresComponentTreeObject_sub6 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type MeasuresComponentTreeObject_sub1 struct {
	Value string `json:"value,omitempty"`
}

type MeasuresSearchObject struct {
	Measures []MeasuresSearchObject_sub2 `json:"measures,omitempty"`
}

type MeasuresSearchObject_sub2 struct {
	BestValue bool                      `json:"bestValue,omitempty"`
	Component string                    `json:"component,omitempty"`
	Metric    string                    `json:"metric,omitempty"`
	Period    MeasuresSearchObject_sub1 `json:"period,omitempty"`
	Value     string                    `json:"value,omitempty"`
}

type MeasuresSearchObject_sub1 struct {
	BestValue bool   `json:"bestValue,omitempty"`
	Index     int64  `json:"index,omitempty"`
	Value     string `json:"value,omitempty"`
}

type MeasuresSearchHistoryObject struct {
	Measures []MeasuresSearchHistoryObject_sub2 `json:"measures,omitempty"`
	Paging   MeasuresSearchHistoryObject_sub3   `json:"paging,omitempty"`
}

type MeasuresSearchHistoryObject_sub1 struct {
	Date  string `json:"date,omitempty"`
	Value string `json:"value,omitempty"`
}

type MeasuresSearchHistoryObject_sub2 struct {
	History []MeasuresSearchHistoryObject_sub1 `json:"history,omitempty"`
	Metric  string                             `json:"metric,omitempty"`
}

type MeasuresSearchHistoryObject_sub3 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type MeasuresComponentOption struct {
	AdditionalFields string `url:"additionalFields,omitempty"` // Description:"Comma-separated list of additional fields that can be returned in the response.",ExampleValue:"period,metrics"
	Branch           string `url:"branch,omitempty"`           // Description:"Branch key. Not available in the community edition.",ExampleValue:"feature/my_branch"
	Component        string `url:"component,omitempty"`        // Description:"Component key",ExampleValue:"my_project"
	MetricKeys       string `url:"metricKeys,omitempty"`       // Description:"Comma-separated list of metric keys",ExampleValue:"ncloc,complexity,violations"
	PullRequest      string `url:"pullRequest,omitempty"`      // Description:"Pull request id. Not available in the community edition.",ExampleValue:"5461"
}

// Component Return component with specified measures.<br>Requires one of the following permissions:<ul><li>'Browse' on the project of the specified component</li><li>'Execute Analysis' on the project of the specified component</li></ul>
func (s *MeasuresService) Component(opt *MeasuresComponentOption) (v *MeasuresComponentObject, resp *http.Response, err error) {
	err = s.ValidateComponentOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "measures/component", opt)
	if err != nil {
		return
	}
	v = new(MeasuresComponentObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type MeasuresComponentTreeOption struct {
	AdditionalFields string `url:"additionalFields,omitempty"` // Description:"Comma-separated list of additional fields that can be returned in the response.",ExampleValue:"period,metrics"
	Asc              string `url:"asc,omitempty"`              // Description:"Ascending sort",ExampleValue:""
	Branch           string `url:"branch,omitempty"`           // Description:"Branch key. Not available in the community edition.",ExampleValue:"feature/my_branch"
	Component        string `url:"component,omitempty"`        // Description:"Component key. The search is based on this component.",ExampleValue:"my_project"
	MetricKeys       string `url:"metricKeys,omitempty"`       // Description:"Comma-separated list of metric keys. Types DISTRIB are not allowed. For type DATA only security_issues, maintainability_issues, reliability_issues, new_maintainability_issues, new_security_issues, new_reliability_issues metrics are supported",ExampleValue:"ncloc,complexity,violations"
	MetricPeriodSort string `url:"metricPeriodSort,omitempty"` // Description:"Sort measures by leak period or not ?. The 's' parameter must contain the 'metricPeriod' value.",ExampleValue:""
	MetricSort       string `url:"metricSort,omitempty"`       // Description:"Metric key to sort by. The 's' parameter must contain the 'metric' or 'metricPeriod' value. It must be part of the 'metricKeys' parameter",ExampleValue:"ncloc"
	MetricSortFilter string `url:"metricSortFilter,omitempty"` // Description:"Filter components. Sort must be on a metric. Possible values are: <ul><li>all: return all components</li><li>withMeasuresOnly: filter out components that do not have a measure on the sorted metric</li></ul>",ExampleValue:""
	P                string `url:"p,omitempty"`                // Description:"1-based page number",ExampleValue:"42"
	Ps               string `url:"ps,omitempty"`               // Description:"Page size. Must be greater than 0 and less or equal than 500",ExampleValue:"20"
	PullRequest      string `url:"pullRequest,omitempty"`      // Description:"Pull request id. Not available in the community edition.",ExampleValue:"5461"
	Q                string `url:"q,omitempty"`                // Description:"Limit search to: <ul><li>component names that contain the supplied string</li><li>component keys that are exactly the same as the supplied string</li></ul>",ExampleValue:"FILE_NAM"
	Qualifiers       string `url:"qualifiers,omitempty"`       // Description:"Comma-separated list of component qualifiers. Filter the results with the specified qualifiers. Possible values are:<ul><li>UTS - Test Files</li><li>FIL - Files</li><li>DIR - Directories</li><li>TRK - Projects</li></ul>",ExampleValue:""
	S                string `url:"s,omitempty"`                // Description:"Comma-separated list of sort fields",ExampleValue:"name,path"
	Strategy         string `url:"strategy,omitempty"`         // Description:"Strategy to search for base component descendants:<ul><li>children: return the children components of the base component. Grandchildren components are not returned</li><li>all: return all the descendants components of the base component. Grandchildren are returned.</li><li>leaves: return all the descendant components (files, in general) which don't have other children. They are the leaves of the component tree.</li></ul>",ExampleValue:""
}

// ComponentTree Navigate through components based on the chosen strategy with specified measures.<br>Requires the following permission: 'Browse' on the specified project.<br>For applications, it also requires 'Browse' permission on its child projects. <br>When limiting search with the q parameter, directories are not returned.
func (s *MeasuresService) ComponentTree(opt *MeasuresComponentTreeOption) (v *MeasuresComponentTreeObject, resp *http.Response, err error) {
	err = s.ValidateComponentTreeOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "measures/component_tree", opt)
	if err != nil {
		return
	}
	v = new(MeasuresComponentTreeObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type MeasuresSearchOption struct {
	MetricKeys  string `url:"metricKeys,omitempty"`  // Description:"Comma-separated list of metric keys",ExampleValue:"ncloc,complexity,violations"
	ProjectKeys string `url:"projectKeys,omitempty"` // Description:"Comma-separated list of project, view or sub-view keys",ExampleValue:"my_project,another_project"
}

// Search Search for project measures ordered by project names.<br>At most 100 projects can be provided.<br>Returns the projects with the 'Browse' permission.
func (s *MeasuresService) Search(opt *MeasuresSearchOption) (v *MeasuresSearchObject, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "measures/search", opt)
	if err != nil {
		return
	}
	v = new(MeasuresSearchObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type MeasuresSearchHistoryOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key. Not available in the community edition.",ExampleValue:"feature/my_branch"
	Component   string `url:"component,omitempty"`   // Description:"Component key",ExampleValue:"my_project"
	From        string `url:"from,omitempty"`        // Description:"Filter measures created after the given date (inclusive). <br>Either a date (server timezone) or datetime can be provided",ExampleValue:"2017-10-19 or 2017-10-19T13:00:00+0200"
	Metrics     string `url:"metrics,omitempty"`     // Description:"Comma-separated list of metric keys",ExampleValue:"ncloc,coverage,new_violations"
	P           string `url:"p,omitempty"`           // Description:"1-based page number",ExampleValue:"42"
	Ps          string `url:"ps,omitempty"`          // Description:"Page size. Must be greater than 0 and less or equal than 1000",ExampleValue:"20"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id. Not available in the community edition.",ExampleValue:"5461"
	To          string `url:"to,omitempty"`          // Description:"Filter measures created before the given date (inclusive). <br>Either a date (server timezone) or datetime can be provided",ExampleValue:"2017-10-19 or 2017-10-19T13:00:00+0200"
}

// SearchHistory Search measures history of a component.<br>Measures are ordered chronologically.<br>Pagination applies to the number of measures for each metric.<br>Requires the following permission: 'Browse' on the specified component. <br>For applications, it also requires 'Browse' permission on its child projects.
func (s *MeasuresService) SearchHistory(opt *MeasuresSearchHistoryOption) (v *MeasuresSearchHistoryObject, resp *http.Response, err error) {
	err = s.ValidateSearchHistoryOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "measures/search_history", opt)
	if err != nil {
		return
	}
	v = new(MeasuresSearchHistoryObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
