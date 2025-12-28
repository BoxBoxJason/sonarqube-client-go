// Access the analysis cache
package sonargo

import "net/http"

type AnalysisCacheService struct {
	client *Client
}

type AnalysisCacheClearOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Filter which project's branch's cached data will be cleared with the provided key. 'project' parameter must be set.",ExampleValue:"6468"
	Project string `url:"project,omitempty"` // Description:"Filter which project's cached data will be cleared with the provided key.",ExampleValue:"org.sonarsource.sonarqube:sonarqube-private"
}

// Clear Clear all or part of the scanner's cached data. Requires global administration permission.
func (s *AnalysisCacheService) Clear(opt *AnalysisCacheClearOption) (resp *http.Response, err error) {
	err = s.ValidateClearOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "analysis_cache/clear", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AnalysisCacheGetOption struct {
	Branch  string `url:"branch,omitempty"`  // Description:"Branch key. If not provided, main branch will be used.",ExampleValue:"feature/my_branch"
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:"my_project"
}

// Get Get the scanner's cached data for a branch. Requires scan permission on the project. Data is returned gzipped if the corresponding 'Accept-Encoding' header is set in the request.
func (s *AnalysisCacheService) Get(opt *AnalysisCacheGetOption) (resp *http.Response, err error) {
	err = s.ValidateGetOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "analysis_cache/get", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}
