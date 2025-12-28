// Get duplication information for a project.
package sonargo

import "net/http"

type DuplicationsService struct {
	client *Client
}

type DuplicationsShowObject struct {
	Duplications []DuplicationsShowObject_sub2 `json:"duplications,omitempty"`
	Files        DuplicationsShowObject_sub4   `json:"files,omitempty"`
}

type DuplicationsShowObject_sub2 struct {
	Blocks []DuplicationsShowObject_sub1 `json:"blocks,omitempty"`
}

type DuplicationsShowObject_sub3 struct {
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
}

type DuplicationsShowObject_sub4 struct {
	One   DuplicationsShowObject_sub3 `json:"1,omitempty"`
	Two   DuplicationsShowObject_sub3 `json:"2,omitempty"`
	Three DuplicationsShowObject_sub3 `json:"3,omitempty"`
}

type DuplicationsShowObject_sub1 struct {
	Ref  string `json:"_ref,omitempty"`
	From int64  `json:"from,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type DuplicationsShowOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	Key         string `url:"key,omitempty"`         // Description:"File key",ExampleValue:"my_project:/src/foo/Bar.php"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
}

// Show Get duplications. Require Browse permission on file's project
func (s *DuplicationsService) Show(opt *DuplicationsShowOption) (v *DuplicationsShowObject, resp *http.Response, err error) {
	err = s.ValidateShowOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "duplications/show", opt)
	if err != nil {
		return
	}
	v = new(DuplicationsShowObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
