package sonargo

import (
	"net/http"
)

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedProjectQualifiers is the set of allowed values for project qualifiers.
	allowedProjectQualifiers = map[string]struct{}{
		"TRK": {},
		"VW":  {},
		"APP": {},
	}

	// allowedProjectVisibility is the set of allowed values for project visibility.
	allowedProjectVisibility = map[string]struct{}{
		"private": {},
		"public":  {},
	}
)

// ProjectsService handles communication with the Projects related methods of the SonarQube API.
// Manage project existence.
//
// Since: 2.10.
type ProjectsService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectsCreate represents the response from creating a project.
type ProjectsCreate struct {
	// Project is the created project.
	Project Project `json:"project,omitzero"`
}

// Project represents a project.
type Project struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Qualifier is the project qualifier (TRK, VW, APP).
	Qualifier string `json:"qualifier,omitempty"`
	// Visibility is the project visibility (public or private).
	Visibility string `json:"visibility,omitempty"`
}

// ProjectsSearch represents the response from searching projects.
type ProjectsSearch struct {
	// Components is the list of projects.
	Components []ProjectComponent `json:"components,omitempty"`
	// Paging is the pagination info.
	Paging Paging `json:"paging,omitzero"`
}

// ProjectComponent represents a project component.
type ProjectComponent struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Qualifier is the project qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// Visibility is the project visibility.
	Visibility string `json:"visibility,omitempty"`
	// LastAnalysisDate is the date of the last analysis.
	LastAnalysisDate string `json:"lastAnalysisDate,omitempty"`
	// Revision is the last analysis revision.
	Revision string `json:"revision,omitempty"`
	// Managed indicates if the project is managed.
	Managed bool `json:"managed,omitempty"`
}

// ProjectsSearchMyProjects represents the response from searching my projects.
type ProjectsSearchMyProjects struct {
	// Projects is the list of projects.
	Projects []MyProject `json:"projects,omitempty"`
	// Paging is the pagination info.
	Paging Paging `json:"paging,omitzero"`
}

// MyProject represents a user's project.
type MyProject struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Description is the project description.
	Description string `json:"description,omitempty"`
	// LastAnalysisDate is the date of the last analysis.
	LastAnalysisDate string `json:"lastAnalysisDate,omitempty"`
	// QualityGate is the project's quality gate status.
	QualityGate string `json:"qualityGate,omitempty"`
	// Links is the list of project links.
	Links []MyProjectLink `json:"links,omitempty"`
}

// MyProjectLink represents a project link in my projects search.
type MyProjectLink struct {
	// Name is the link name.
	Name string `json:"name,omitempty"`
	// Type is the link type.
	Type string `json:"type,omitempty"`
	// Href is the link URL.
	Href string `json:"href,omitempty"`
}

// ProjectsSearchMyScannableProjects represents the response from searching my scannable projects.
type ProjectsSearchMyScannableProjects struct {
	// Projects is the list of scannable projects.
	Projects []ScannableProject `json:"projects,omitempty"`
}

// ScannableProject represents a scannable project.
type ScannableProject struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectsBulkDeleteOption represents options for bulk deleting projects.
type ProjectsBulkDeleteOption struct {
	// AnalyzedBefore filters projects analyzed before the given date (optional).
	// Format: YYYY-MM-DD.
	AnalyzedBefore string `url:"analyzedBefore,omitempty"`
	// Projects is the list of project keys to delete (optional).
	Projects []string `url:"projects,omitempty,comma"`
	// Query is used to filter projects by name (optional).
	Query string `url:"q,omitempty"`
	// Qualifiers filters by project qualifiers (optional).
	// Possible values: TRK, VW, APP.
	Qualifiers []string `url:"qualifiers,omitempty,comma"`
	// OnProvisionedOnly filters only provisioned projects (optional).
	OnProvisionedOnly bool `url:"onProvisionedOnly,omitempty"`
}

// ProjectsCreateOption represents options for creating a project.
type ProjectsCreateOption struct {
	// Name is the project name (required).
	// Maximum length: 500 characters.
	Name string `url:"name,omitempty"`
	// Project is the project key (required).
	// Maximum length: 400 characters.
	Project string `url:"project,omitempty"`
	// MainBranch is the name of the main branch (optional).
	// If not provided, the default main branch name will be used.
	MainBranch string `url:"mainBranch,omitempty"`
	// NewCodeDefinitionType is the type of new code definition (optional).
	// Possible values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH.
	NewCodeDefinitionType string `url:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the value of new code definition (optional).
	// Required if newCodeDefinitionType is NUMBER_OF_DAYS or REFERENCE_BRANCH.
	NewCodeDefinitionValue string `url:"newCodeDefinitionValue,omitempty"`
	// Visibility is the project visibility (optional).
	// If no visibility is specified, the default visibility will be used.
	// Possible values: private, public.
	Visibility string `url:"visibility,omitempty"`
}

// ProjectsDeleteOption represents options for deleting a project.
type ProjectsDeleteOption struct {
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
}

// ProjectsSearchOption represents options for searching projects.
//
//nolint:govet // Field alignment is less important than logical grouping
type ProjectsSearchOption struct {
	PaginationArgs

	// AnalyzedBefore filters projects analyzed before the given date (optional).
	// Format: YYYY-MM-DD.
	AnalyzedBefore string `url:"analyzedBefore,omitempty"`
	// OnProvisionedOnly filters only provisioned projects (optional).
	OnProvisionedOnly bool `url:"onProvisionedOnly,omitempty"`
	// Projects is the list of project keys to return (optional).
	Projects []string `url:"projects,omitempty,comma"`
	// Query is used to filter projects by name or key (optional).
	Query string `url:"q,omitempty"`
	// Qualifiers filters by project qualifiers (optional).
	// Possible values: TRK, VW, APP.
	Qualifiers []string `url:"qualifiers,omitempty,comma"`
}

// ProjectsSearchMyProjectsOption represents options for searching my projects.
type ProjectsSearchMyProjectsOption struct {
	PaginationArgs
}

// ProjectsUpdateDefaultVisibilityOption represents options for updating default visibility.
type ProjectsUpdateDefaultVisibilityOption struct {
	// ProjectVisibility is the new default visibility (required).
	// Possible values: private, public.
	ProjectVisibility string `url:"projectVisibility,omitempty"`
}

// ProjectsUpdateKeyOption represents options for updating a project key.
type ProjectsUpdateKeyOption struct {
	// From is the current project key (required).
	From string `url:"from,omitempty"`
	// To is the new project key (required).
	To string `url:"to,omitempty"`
}

// ProjectsUpdateVisibilityOption represents options for updating project visibility.
type ProjectsUpdateVisibilityOption struct {
	// Project is the project key (required).
	Project string `url:"project,omitempty"`
	// Visibility is the new visibility (required).
	// Possible values: private, public.
	Visibility string `url:"visibility,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Methods
// -----------------------------------------------------------------------------

// ValidateBulkDeleteOpt validates the options for BulkDelete.
func (s *ProjectsService) ValidateBulkDeleteOpt(opt *ProjectsBulkDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// At least one filter is required
	if opt.AnalyzedBefore == "" && len(opt.Projects) == 0 && opt.Query == "" {
		return NewValidationError("Projects", "at least one of analyzedBefore, projects or q is required", ErrMissingRequired)
	}

	if len(opt.Qualifiers) > 0 {
		err := AreValuesAuthorized(opt.Qualifiers, allowedProjectQualifiers, "Qualifiers")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateCreateOpt validates the options for Create.
func (s *ProjectsService) ValidateCreateOpt(opt *ProjectsCreateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxProjectNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Project, MaxProjectKeyLength, "Project")
	if err != nil {
		return err
	}

	if opt.Visibility != "" {
		err = IsValueAuthorized(opt.Visibility, allowedProjectVisibility, "Visibility")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteOpt validates the options for Delete.
func (s *ProjectsService) ValidateDeleteOpt(opt *ProjectsDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateSearchOpt validates the options for Search.
func (s *ProjectsService) ValidateSearchOpt(opt *ProjectsSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if len(opt.Qualifiers) > 0 {
		err := AreValuesAuthorized(opt.Qualifiers, allowedProjectQualifiers, "Qualifiers")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchMyProjectsOpt validates the options for SearchMyProjects.
func (s *ProjectsService) ValidateSearchMyProjectsOpt(opt *ProjectsSearchMyProjectsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return nil
}

// ValidateUpdateDefaultVisibilityOpt validates the options for UpdateDefaultVisibility.
func (s *ProjectsService) ValidateUpdateDefaultVisibilityOpt(opt *ProjectsUpdateDefaultVisibilityOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectVisibility, "ProjectVisibility")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.ProjectVisibility, allowedProjectVisibility, "ProjectVisibility")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateKeyOpt validates the options for UpdateKey.
func (s *ProjectsService) ValidateUpdateKeyOpt(opt *ProjectsUpdateKeyOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.From, "From")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.To, "To")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateVisibilityOpt validates the options for UpdateVisibility.
func (s *ProjectsService) ValidateUpdateVisibilityOpt(opt *ProjectsUpdateVisibilityOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Visibility, "Visibility")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Visibility, allowedProjectVisibility, "Visibility")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// BulkDelete deletes one or several projects.
// At least one parameter is required among 'analyzedBefore', 'projects', and 'q'.
// Requires 'Administer System' permission.
//
// Since: 5.2.
func (s *ProjectsService) BulkDelete(opt *ProjectsBulkDeleteOption) (*http.Response, error) {
	err := s.ValidateBulkDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "projects/bulk_delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Create provisions a new project.
// Requires 'Create Projects' permission.
//
// Since: 4.0.
func (s *ProjectsService) Create(opt *ProjectsCreateOption) (*ProjectsCreate, *http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "projects/create", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectsCreate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes a project.
// Requires 'Administer System' permission or 'Administer' permission on the project.
//
// Since: 5.2.
func (s *ProjectsService) Delete(opt *ProjectsDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "projects/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Search searches for projects.
// Requires 'Browse' permission on the returned projects.
//
// Since: 6.3.
func (s *ProjectsService) Search(opt *ProjectsSearchOption) (*ProjectsSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "projects/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchMyProjects returns the list of projects the current user can scan.
// Only the favorite and recently analyzed projects will be returned.
//
// Since: 6.4.
func (s *ProjectsService) SearchMyProjects(opt *ProjectsSearchMyProjectsOption) (*ProjectsSearchMyProjects, *http.Response, error) {
	err := s.ValidateSearchMyProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "projects/search_my_projects", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectsSearchMyProjects)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchMyScannableProjects returns the list of projects the current user can analyze.
// Requires authentication.
//
// Since: 9.5.
func (s *ProjectsService) SearchMyScannableProjects() (*ProjectsSearchMyScannableProjects, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "projects/search_my_scannable_projects", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectsSearchMyScannableProjects)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateDefaultVisibility updates the default visibility of new projects.
// Requires 'Administer System' permission.
//
// Since: 6.4.
func (s *ProjectsService) UpdateDefaultVisibility(opt *ProjectsUpdateDefaultVisibilityOption) (*http.Response, error) {
	err := s.ValidateUpdateDefaultVisibilityOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "projects/update_default_visibility", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateKey updates a project key.
// Requires 'Administer' permission on the project.
//
// Since: 6.1.
func (s *ProjectsService) UpdateKey(opt *ProjectsUpdateKeyOption) (*http.Response, error) {
	err := s.ValidateUpdateKeyOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "projects/update_key", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateVisibility updates a project visibility.
// Requires 'Project administer' permission on the project.
//
// Since: 6.4.
func (s *ProjectsService) UpdateVisibility(opt *ProjectsUpdateVisibilityOption) (*http.Response, error) {
	err := s.ValidateUpdateVisibilityOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "projects/update_visibility", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
