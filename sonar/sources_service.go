// Get details on source files. See also api/tests.
package sonargo

import "net/http"

type SourcesService struct {
	client *Client
}

type SourcesIndexObject []struct {
	Four8 string `json:"48,omitempty"`
	Four9 string `json:"49,omitempty"`
	Five0 string `json:"50,omitempty"`
	Five1 string `json:"51,omitempty"`
	Five2 string `json:"52,omitempty"`
	Five3 string `json:"53,omitempty"`
	Five4 string `json:"54,omitempty"`
	Five5 string `json:"55,omitempty"`
	Five6 string `json:"56,omitempty"`
	Five7 string `json:"57,omitempty"`
	Five8 string `json:"58,omitempty"`
	Five9 string `json:"59,omitempty"`
	Six0  string `json:"60,omitempty"`
}

type SourcesIssueSnippetsObject struct {
	Test_project_src_main_java_org_example_Main_java SourcesIssueSnippetsObject_sub4 `json:"test-project:src/main/java/org/example/Main.java,omitempty"`
}

type SourcesIssueSnippetsObject_sub3 struct {
	Code        string `json:"code,omitempty"`
	Duplicated  bool   `json:"duplicated,omitempty"`
	IsNew       bool   `json:"isNew,omitempty"`
	Line        int64  `json:"line,omitempty"`
	LineHits    int64  `json:"lineHits,omitempty"`
	ScmAuthor   string `json:"scmAuthor,omitempty"`
	ScmDate     string `json:"scmDate,omitempty"`
	ScmRevision string `json:"scmRevision,omitempty"`
	UtLineHits  int64  `json:"utLineHits,omitempty"`
}

type SourcesIssueSnippetsObject_sub4 struct {
	Component SourcesIssueSnippetsObject_sub2   `json:"component,omitempty"`
	Sources   []SourcesIssueSnippetsObject_sub3 `json:"sources,omitempty"`
}

type SourcesIssueSnippetsObject_sub1 struct {
	Coverage string `json:"coverage,omitempty"`
	Issues   string `json:"issues,omitempty"`
	Lines    string `json:"lines,omitempty"`
}

type SourcesIssueSnippetsObject_sub2 struct {
	Key         string                          `json:"key,omitempty"`
	LongName    string                          `json:"longName,omitempty"`
	Measures    SourcesIssueSnippetsObject_sub1 `json:"measures,omitempty"`
	Name        string                          `json:"name,omitempty"`
	Path        string                          `json:"path,omitempty"`
	Project     string                          `json:"project,omitempty"`
	ProjectName string                          `json:"projectName,omitempty"`
	Q           string                          `json:"q,omitempty"`
	UUID        string                          `json:"uuid,omitempty"`
}

type SourcesLinesObject struct {
	Sources []SourcesLinesObject_sub1 `json:"sources,omitempty"`
}

type SourcesLinesObject_sub1 struct {
	Code        string `json:"code,omitempty"`
	IsNew       bool   `json:"isNew,omitempty"`
	Line        int64  `json:"line,omitempty"`
	ScmAuthor   string `json:"scmAuthor,omitempty"`
	ScmDate     string `json:"scmDate,omitempty"`
	ScmRevision string `json:"scmRevision,omitempty"`
}

type SourcesScmObject struct {
	Scm [][]string `json:"scm,omitempty"`
}

type SourcesShowObject struct {
	Sources [][]interface{} `json:"sources,omitempty"`
}

type SourcesIndexOption struct {
	From     string `url:"from,omitempty"`     // Description:"First line",ExampleValue:""
	Resource string `url:"resource,omitempty"` // Description:"File key",ExampleValue:"my_project:/src/foo/Bar.php"
	To       string `url:"to,omitempty"`       // Description:"Last line (excluded). If not specified, all lines are returned until end of file",ExampleValue:""
}

// Index Get source code as line number / text pairs. Require See Source Code permission on file
func (s *SourcesService) Index(opt *SourcesIndexOption) (v *SourcesIndexObject, resp *http.Response, err error) {
	err = s.ValidateIndexOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "sources/index", opt)
	if err != nil {
		return
	}
	v = new(SourcesIndexObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SourcesIssueSnippetsOption struct {
	IssueKey string `url:"issueKey,omitempty"` // Description:"Issue or hotspot key",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
}

// IssueSnippets Get code snippets involved in an issue or hotspot. Requires 'See Source Code permission' permission on the project<br/>
func (s *SourcesService) IssueSnippets(opt *SourcesIssueSnippetsOption) (v *SourcesIssueSnippetsObject, resp *http.Response, err error) {
	err = s.ValidateIssueSnippetsOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "sources/issue_snippets", opt)
	if err != nil {
		return
	}
	v = new(SourcesIssueSnippetsObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SourcesLinesOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	From        string `url:"from,omitempty"`        // Description:"First line to return. Starts from 1",ExampleValue:"10"
	Key         string `url:"key,omitempty"`         // Description:"File key. Mandatory if param 'uuid' is not set. Available since 5.2",ExampleValue:"my_project:/src/foo/Bar.php"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
	To          string `url:"to,omitempty"`          // Description:"Optional last line to return (inclusive). It must be greater than or equal to parameter 'from'. If unset, then all the lines greater than or equal to 'from' are returned.",ExampleValue:"20"
	Uuid        string `url:"uuid,omitempty"`        // Description:"File uuid. Mandatory if param 'key' is not set",ExampleValue:"f333aab4-7e3a-4d70-87e1-f4c491f05e5c"
}

// Lines Show source code with line oriented info. Requires See Source Code permission on file's project<br/>Each element of the result array is an object which contains:<ol><li>Line number</li><li>Content of the line</li><li>Author of the line (from SCM information)</li><li>Revision of the line (from SCM information)</li><li>Last commit date of the line (from SCM information)</li><li>Line hits from coverage</li><li>Number of conditions to cover in tests</li><li>Number of conditions covered by tests</li><li>Whether the line is new</li></ol>
func (s *SourcesService) Lines(opt *SourcesLinesOption) (v *SourcesLinesObject, resp *http.Response, err error) {
	err = s.ValidateLinesOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "sources/lines", opt)
	if err != nil {
		return
	}
	v = new(SourcesLinesObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SourcesRawOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	Key         string `url:"key,omitempty"`         // Description:"File key",ExampleValue:"my_project:src/foo/Bar.php"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
}

// Raw Get source code as raw text. Require 'See Source Code' permission on file
func (s *SourcesService) Raw(opt *SourcesRawOption) (v *string, resp *http.Response, err error) {
	err = s.ValidateRawOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "sources/raw", opt)
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

type SourcesScmOption struct {
	CommitsByLine string `url:"commits_by_line,omitempty"` // Description:"Group lines by SCM commit if value is false, else display commits for each line, even if two consecutive lines relate to the same commit.",ExampleValue:""
	From          string `url:"from,omitempty"`            // Description:"First line to return. Starts at 1",ExampleValue:"10"
	Key           string `url:"key,omitempty"`             // Description:"File key",ExampleValue:"my_project:/src/foo/Bar.php"
	To            string `url:"to,omitempty"`              // Description:"Last line to return (inclusive)",ExampleValue:"20"
}

// Scm Get SCM information of source files. Require See Source Code permission on file's project<br/>Each element of the result array is composed of:<ol><li>Line number</li><li>Author of the commit</li><li>Datetime of the commit (before 5.2 it was only the Date)</li><li>Revision of the commit (added in 5.2)</li></ol>
func (s *SourcesService) Scm(opt *SourcesScmOption) (v *SourcesScmObject, resp *http.Response, err error) {
	err = s.ValidateScmOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "sources/scm", opt)
	if err != nil {
		return
	}
	v = new(SourcesScmObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SourcesShowOption struct {
	From string `url:"from,omitempty"` // Description:"First line to return. Starts at 1",ExampleValue:"10"
	Key  string `url:"key,omitempty"`  // Description:"File key",ExampleValue:"my_project:/src/foo/Bar.php"
	To   string `url:"to,omitempty"`   // Description:"Last line to return (inclusive)",ExampleValue:"20"
}

// Show Get source code. Requires See Source Code permission on file's project<br/>Each element of the result array is composed of:<ol><li>Line number</li><li>Content of the line</li></ol>
func (s *SourcesService) Show(opt *SourcesShowOption) (v *SourcesShowObject, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "sources/show", opt)
	if err != nil {
		return
	}
	v = new(SourcesShowObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
