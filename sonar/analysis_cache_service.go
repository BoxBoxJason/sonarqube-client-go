package sonargo

import (
	"fmt"
	"net/http"
)

// AnalysisCacheService handles communication with the analysis cache related methods
// of the SonarQube API.
// This service provides access to scanner's cached data for branches.
type AnalysisCacheService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AnalysisCacheClearOption contains parameters for the Clear method.
type AnalysisCacheClearOption struct {
	// Branch filters which project's branch's cached data will be cleared.
	// The 'Project' parameter must be set when using this.
	Branch string `url:"branch,omitempty"`
	// Project filters which project's cached data will be cleared.
	Project string `url:"project"`
}

// AnalysisCacheGetOption contains parameters for the Get method.
type AnalysisCacheGetOption struct {
	// Branch key. If not provided, main branch will be used.
	Branch string `url:"branch,omitempty"`
	// Project key.
	// This field is required.
	Project string `url:"project,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateClearOpt validates the options for the Clear method.
// Currently, there are no required fields.
func (s *AnalysisCacheService) ValidateClearOpt(opt *AnalysisCacheClearOption) error {
	// When filtering by branch, Project must be set as documented.
	if opt != nil && opt.Branch != "" && opt.Project == "" {
		return NewValidationError("Project", "Project must be set when Branch is specified", ErrMissingRequired)
	}

	return nil
}

// ValidateGetOpt validates the options for the Get method.
// Currently, there are no required fields.
func (s *AnalysisCacheService) ValidateGetOpt(opt *AnalysisCacheGetOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Check required fields
	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Clear clears all or part of the scanner's cached data.
// Requires global administration permission.
//
// API endpoint: POST /api/analysis_cache/clear.
// WARNING: This is an internal API and may change without notice.
func (s *AnalysisCacheService) Clear(opt *AnalysisCacheClearOption) (*http.Response, error) {
	err := s.ValidateClearOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "analysis_cache/clear", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Get returns the scanner's cached data for a branch.
// Requires scan permission on the project.
// Data is returned gzipped if the corresponding 'Accept-Encoding' header is set in the request.
// The response body contains the raw binary data; the caller is responsible for reading and closing it.
//
// API endpoint: GET /api/analysis_cache/get.
func (s *AnalysisCacheService) Get(opt *AnalysisCacheGetOption) (*http.Response, error) {
	err := s.ValidateGetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "analysis_cache/get", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
