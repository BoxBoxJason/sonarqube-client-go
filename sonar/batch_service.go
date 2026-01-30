package sonargo

import "net/http"

// BatchService provides access to the /batch API endpoints for scanner JAR files and referentials.
type BatchService struct {
	client *Client
}

// BatchFileOption contains parameters for the File endpoint.
type BatchFileOption struct {
	Name string `url:"name,omitempty"` // Description:"File name",ExampleValue:"batch-library-2.3.jar"
}

// BatchProjectOption contains parameters for the Project endpoint.
type BatchProjectOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	Key         string `url:"key,omitempty"`         // Description:"Project key",ExampleValue:"my_project"
	Profile     string `url:"profile,omitempty"`     // Description:"Profile name",ExampleValue:"SonarQube Way"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
}

// BatchProject contains the response from the Project endpoint.
type BatchProject struct {
	FileDataByModuleAndPath map[string]map[string]BatchFileData `json:"fileDataByModuleAndPath,omitempty"`
	LastAnalysisDate        int64                               `json:"lastAnalysisDate,omitempty"`
	Timestamp               int64                               `json:"timestamp,omitempty"`
}

// BatchFileData contains file hash and revision information.
type BatchFileData struct {
	Hash     string `json:"hash,omitempty"`
	Revision string `json:"revision,omitempty"`
}

// ValidateFileOpt validates the options for the File endpoint.
func (s *BatchService) ValidateFileOpt(opt *BatchFileOption) error {
	if opt == nil {
		return nil
	}

	return nil
}

// ValidateProjectOpt validates the options for the Project endpoint.
func (s *BatchService) ValidateProjectOpt(opt *BatchProjectOption) error {
	if opt == nil {
		return nil
	}

	return nil
}

// File downloads a JAR file listed in the index (see batch/index).
// This endpoint returns binary data for the requested JAR file.
func (s *BatchService) File(opt *BatchFileOption) (v *string, resp *http.Response, err error) {
	err = s.ValidateFileOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "batch/file", opt)
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

// Index lists the JAR files to be downloaded by scanners.
// Returns a list of JAR file names and their hashes.
func (s *BatchService) Index() (v *string, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "batch/index", nil)
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

// Project returns project repository information including file hashes
// for incremental analysis.
func (s *BatchService) Project(opt *BatchProjectOption) (v *BatchProject, resp *http.Response, err error) {
	err = s.ValidateProjectOpt(opt)
	if err != nil {
		return
	}

	req, err := s.client.NewRequest("GET", "batch/project", opt)
	if err != nil {
		return
	}

	v = new(BatchProject)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return
}
