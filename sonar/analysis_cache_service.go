package sonargo

import "net/http"

// AnalysisCacheService handles communication with the analysis cache related methods
// of the SonarQube API.
// This service provides access to scanner's cached data for branches.
type AnalysisCacheService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AnalysisCacheData represents the cached data returned by the Get method.
// The actual content is binary/gzipped data that should be handled accordingly.
type AnalysisCacheData struct{}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AnalysisCacheClearOption contains parameters for the Clear method.
type AnalysisCacheClearOption struct {
	// Branch filters which project's branch's cached data will be cleared.
	// The 'Project' parameter must be set when using this.
	Branch string `url:"branch,omitempty"`
	// Project filters which project's cached data will be cleared.
	Project string `url:"project,omitempty"`
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
	// No required fields
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
//
// API endpoint: GET /api/analysis_cache/get.
func (s *AnalysisCacheService) Get(opt *AnalysisCacheGetOption) (*AnalysisCacheData, *http.Response, error) {
	err := s.ValidateGetOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "analysis_cache/get", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(AnalysisCacheData)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
