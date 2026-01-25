package sonargo

import "net/http"

// ProjectLinksService handles communication with the project links related methods
// of the SonarQube API.
// This service manages project links (custom URLs associated with projects).
type ProjectLinksService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// ProjectLink represents a link associated with a project.
type ProjectLink struct {
	// ID is the unique identifier of the link.
	ID string `json:"id,omitempty"`
	// Name is the display name of the link.
	Name string `json:"name,omitempty"`
	// Type is the type of the link (e.g., "homepage", "ci", "issue", "scm").
	Type string `json:"type,omitempty"`
	// URL is the target URL of the link.
	URL string `json:"url,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectLinksCreate represents the response from creating a project link.
type ProjectLinksCreate struct {
	// Link is the created link.
	Link ProjectLink `json:"link,omitzero"`
}

// ProjectLinksSearch represents the response from searching project links.
type ProjectLinksSearch struct {
	// Links is the list of links found.
	Links []ProjectLink `json:"links,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectLinksCreateOption contains parameters for the Create method.
type ProjectLinksCreateOption struct {
	// Name is the display name of the link.
	// This field is required. Maximum length is 128 characters.
	Name string `url:"name"`
	// ProjectID is the project id.
	// Either ProjectID or ProjectKey must be provided.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key.
	// Either ProjectID or ProjectKey must be provided.
	ProjectKey string `url:"projectKey,omitempty"`
	// URL is the target URL of the link.
	// This field is required. Maximum length is 2048 characters.
	URL string `url:"url"`
}

// ProjectLinksDeleteOption contains parameters for the Delete method.
type ProjectLinksDeleteOption struct {
	// ID is the unique identifier of the link to delete.
	// This field is required.
	ID string `url:"id"`
}

// ProjectLinksSearchOption contains parameters for the Search method.
type ProjectLinksSearchOption struct {
	// ProjectID is the project id.
	// Either ProjectID or ProjectKey must be provided.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key.
	// Either ProjectID or ProjectKey must be provided.
	ProjectKey string `url:"projectKey,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCreateOpt validates the options for the Create method.
func (s *ProjectLinksService) ValidateCreateOpt(opt *ProjectLinksCreateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxLinkNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxLinkURLLength, "URL")
	if err != nil {
		return err
	}

	// Either ProjectID or ProjectKey must be provided
	if opt.ProjectID == "" && opt.ProjectKey == "" {
		return NewValidationError("ProjectID/ProjectKey", "either ProjectID or ProjectKey must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateDeleteOpt validates the options for the Delete method.
func (s *ProjectLinksService) ValidateDeleteOpt(opt *ProjectLinksDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ID, "ID")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
func (s *ProjectLinksService) ValidateSearchOpt(opt *ProjectLinksSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Either ProjectID or ProjectKey must be provided
	if opt.ProjectID == "" && opt.ProjectKey == "" {
		return NewValidationError("ProjectID/ProjectKey", "either ProjectID or ProjectKey must be provided", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Create creates a new project link.
// Requires 'Administer' permission on the specified project, or global 'Administer' permission.
//
// API endpoint: POST /api/project_links/create.
// Since: 6.1.
func (s *ProjectLinksService) Create(opt *ProjectLinksCreateOption) (*ProjectLinksCreate, *http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_links/create", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectLinksCreate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes an existing project link.
// Requires 'Administer' permission on the specified project, or global 'Administer' permission.
//
// API endpoint: POST /api/project_links/delete.
// Since: 6.1.
func (s *ProjectLinksService) Delete(opt *ProjectLinksDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_links/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Search lists links of a project.
// The 'projectId' or 'projectKey' must be provided.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//   - 'Browse' on the specified project
//
// API endpoint: GET /api/project_links/search.
// Since: 6.1.
func (s *ProjectLinksService) Search(opt *ProjectLinksSearchOption) (*ProjectLinksSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_links/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectLinksSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
