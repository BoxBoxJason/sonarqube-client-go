package sonargo

import "net/http"

// ProjectTagsService handles communication with the project tags related methods
// of the SonarQube API.
// This service manages project tags.
type ProjectTagsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectTagsSearch represents the response from searching project tags.
type ProjectTagsSearch struct {
	// Tags is the list of tags found.
	Tags []string `json:"tags,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectTagsSearchOption contains parameters for the Search method.
//
//nolint:govet // Field alignment is less important than logical grouping.
type ProjectTagsSearchOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`

	Query string `url:"q,omitempty"`
}

// ProjectTagsSetOption contains parameters for the Set method.
type ProjectTagsSetOption struct {
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
	// Tags is the list of tags to set on the project.
	Tags []string `url:"tags,comma"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSearchOpt validates the options for the Search method.
func (s *ProjectTagsService) ValidateSearchOpt(opt *ProjectTagsSearchOption) error {
	if opt == nil {
		// Options are optional; nothing to validate.
		return nil
	}
	// Validate pagination arguments (embedded via PaginationArgs).
	err := opt.Validate()
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetOpt validates the options for the Set method.
func (s *ProjectTagsService) ValidateSetOpt(opt *ProjectTagsSetOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Search searches for project tags.
//
// API endpoint: GET /api/project_tags/search.
// Since: 6.4.
func (s *ProjectTagsService) Search(opt *ProjectTagsSearchOption) (*ProjectTagsSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_tags/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectTagsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Set sets tags on a project.
// Requires the 'Administer' rights on the specified project.
//
// API endpoint: POST /api/project_tags/set.
// Since: 6.4.
func (s *ProjectTagsService) Set(opt *ProjectTagsSetOption) (*http.Response, error) {
	err := s.ValidateSetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_tags/set", opt)
	if err != nil {
		return nil, err
	}

	// Handle empty tags array: SonarQube API requires the tags parameter to be present even when empty
	// The query encoder omits empty slices, so we need to add it manually
	if len(opt.Tags) == 0 {
		if req.URL.RawQuery == "" {
			req.URL.RawQuery = "tags="
		} else {
			req.URL.RawQuery += "&tags="
		}
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
