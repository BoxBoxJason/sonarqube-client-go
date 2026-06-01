package sonar

import (
	"context"
	"net/http"
)

// ApplicationsService handles communication with the application related methods
// of the SonarQube API. This service is only available in Enterprise Edition.
type ApplicationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// Application represents a SonarQube application.
type Application struct {
	// Key is the application key.
	Key string `json:"key,omitempty"`
	// Name is the application name.
	Name string `json:"name,omitempty"`
	// Description is the application description.
	Description string `json:"description,omitempty"`
	// Visibility is the application visibility (public or private).
	Visibility string `json:"visibility,omitempty"`
}

// ApplicationDetails represents detailed information about an application.
type ApplicationDetails struct {
	// Key is the application key.
	Key string `json:"key,omitempty"`
	// Name is the application name.
	Name string `json:"name,omitempty"`
	// Description is the application description.
	Description string `json:"description,omitempty"`
	// Visibility is the application visibility.
	Visibility string `json:"visibility,omitempty"`
	// Tags is the list of tags assigned to the application.
	Tags []string `json:"tags,omitempty"`
}

// ApplicationBranch represents a branch of an application.
type ApplicationBranch struct {
	// Name is the branch name.
	Name string `json:"name,omitempty"`
	// IsMain indicates whether this is the main branch.
	IsMain bool `json:"isMain,omitempty"`
}

// ApplicationProject represents a project within an application.
type ApplicationProject struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Selected indicates whether the project is selected in the application.
	Selected bool `json:"selected,omitempty"`
}

// ApplicationsShow represents the response from the show endpoint.
type ApplicationsShow struct {
	// Application contains the application details.
	Application ApplicationDetails `json:"application,omitzero"`
}

// ApplicationsSearchProjects represents the response from the search_projects endpoint.
type ApplicationsSearchProjects struct {
	// Projects is the list of projects in the application.
	Projects []ApplicationProject `json:"projects,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ApplicationsCreateOptions contains parameters for the Create method.
type ApplicationsCreateOptions struct {
	// Description is the application description. Optional.
	Description string `url:"description,omitempty"`
	// Key is the application key. Optional, will be generated from name if not provided.
	Key string `url:"key,omitempty"`
	// Name is the application name. This field is required.
	Name string `url:"name"`
	// Visibility is the application visibility (public or private). Optional.
	Visibility string `url:"visibility,omitempty"`
}

// ApplicationsDeleteOptions contains parameters for the Delete method.
type ApplicationsDeleteOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
}

// ApplicationsShowOptions contains parameters for the Show method.
type ApplicationsShowOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the branch name. Optional.
	Branch string `url:"branch,omitempty"`
}

// ApplicationsUpdateOptions contains parameters for the Update method.
type ApplicationsUpdateOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Description is the new description for the application. Optional.
	Description string `url:"description,omitempty"`
	// Name is the new name for the application. This field is required.
	Name string `url:"name"`
}

// ApplicationsAddProjectOptions contains parameters for the AddProject method.
type ApplicationsAddProjectOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ApplicationsRemoveProjectOptions contains parameters for the RemoveProject method.
type ApplicationsRemoveProjectOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ApplicationsCreateBranchOptions contains parameters for the CreateBranch method.
type ApplicationsCreateBranchOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the branch name. This field is required.
	Branch string `url:"branch"`
	// Project is the list of project keys (can be specified multiple times). This field is required.
	Project []string `url:"project"`
	// ProjectBranch is the list of project branch names (can be specified multiple times). This field is required.
	ProjectBranch []string `url:"projectBranch"`
}

// ApplicationsDeleteBranchOptions contains parameters for the DeleteBranch method.
type ApplicationsDeleteBranchOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the branch name. This field is required.
	Branch string `url:"branch"`
}

// ApplicationsUpdateBranchOptions contains parameters for the UpdateBranch method.
type ApplicationsUpdateBranchOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the current branch name. This field is required.
	Branch string `url:"branch"`
	// Name is the new branch name. This field is required.
	Name string `url:"name"`
	// Project is the list of project keys. This field is required.
	Project []string `url:"project"`
	// ProjectBranch is the list of project branch names. This field is required.
	ProjectBranch []string `url:"projectBranch"`
}

// ApplicationsSetTagsOptions contains parameters for the SetTags method.
type ApplicationsSetTagsOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Tags is a comma-separated list of tags. This field is required.
	Tags string `url:"tags"`
}

// ApplicationsSearchProjectsOptions contains parameters for the SearchProjects method.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type ApplicationsSearchProjectsOptions struct {
	PaginationArgs

	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Q limits search to project names containing this string. Optional.
	Q string `url:"q,omitempty"`
	// Selected filters by selection state (selected, deselected, all). Optional.
	Selected string `url:"selected,omitempty"`
}

// ApplicationsRefreshOptions contains parameters for the Refresh method.
type ApplicationsRefreshOptions struct {
	// Key is the application key. If not specified, all applications are refreshed. Optional.
	Key string `url:"key,omitempty"`
}

// -----------------------------------------------------------------------------
// Allowed Value Sets
// -----------------------------------------------------------------------------

//nolint:gochecknoglobals // constant set of allowed values
var allowedApplicationVisibilities = map[string]struct{}{
	"private": {},
	"public":  {},
}

//nolint:gochecknoglobals // constant set of allowed values
var allowedApplicationProjectSelections = map[string]struct{}{
	SelectionFilterAll:        {},
	SelectionFilterSelected:   {},
	SelectionFilterDeselected: {},
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCreateOpt validates the options for the Create method.
func (s *ApplicationsService) ValidateCreateOpt(opt *ApplicationsCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Visibility, allowedApplicationVisibilities, "Visibility")
}

// ValidateDeleteOpt validates the options for the Delete method.
func (s *ApplicationsService) ValidateDeleteOpt(opt *ApplicationsDeleteOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Application, "Application")
}

// ValidateShowOpt validates the options for the Show method.
func (s *ApplicationsService) ValidateShowOpt(opt *ApplicationsShowOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Application, "Application")
}

// ValidateUpdateOpt validates the options for the Update method.
func (s *ApplicationsService) ValidateUpdateOpt(opt *ApplicationsUpdateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Name, "Name")
}

// ValidateAddProjectOpt validates the options for the AddProject method.
func (s *ApplicationsService) ValidateAddProjectOpt(opt *ApplicationsAddProjectOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateRemoveProjectOpt validates the options for the RemoveProject method.
func (s *ApplicationsService) ValidateRemoveProjectOpt(opt *ApplicationsRemoveProjectOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateCreateBranchOpt validates the options for the CreateBranch method.
func (s *ApplicationsService) ValidateCreateBranchOpt(opt *ApplicationsCreateBranchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Branch, "Branch")
}

// ValidateDeleteBranchOpt validates the options for the DeleteBranch method.
func (s *ApplicationsService) ValidateDeleteBranchOpt(opt *ApplicationsDeleteBranchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Branch, "Branch")
}

// ValidateUpdateBranchOpt validates the options for the UpdateBranch method.
func (s *ApplicationsService) ValidateUpdateBranchOpt(opt *ApplicationsUpdateBranchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Branch, "Branch")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Name, "Name")
}

// ValidateSetTagsOpt validates the options for the SetTags method.
func (s *ApplicationsService) ValidateSetTagsOpt(opt *ApplicationsSetTagsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Tags, "Tags")
}

// ValidateSearchProjectsOpt validates the options for the SearchProjects method.
func (s *ApplicationsService) ValidateSearchProjectsOpt(opt *ApplicationsSearchProjectsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return err
	}

	err = opt.Validate()
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Selected, allowedApplicationProjectSelections, "Selected")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Create creates a new application.
// Requires 'Administer System' or 'Create Applications' permission.
//
// API endpoint: POST /api/applications/create.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) Create(ctx context.Context, opt *ApplicationsCreateOptions) (*http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/create", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Delete deletes an application definition.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/delete.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) Delete(ctx context.Context, opt *ApplicationsDeleteOptions) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/delete", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Show returns an application and its associated projects.
// Requires 'Browse' permission on the application and on its child projects.
//
// API endpoint: GET /api/applications/show.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) Show(ctx context.Context, opt *ApplicationsShowOptions) (*ApplicationsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "applications/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ApplicationsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Update updates an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/update.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) Update(ctx context.Context, opt *ApplicationsUpdateOptions) (*http.Response, error) {
	err := s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/update", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddProject adds a project to an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/add_project.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) AddProject(ctx context.Context, opt *ApplicationsAddProjectOptions) (*http.Response, error) {
	err := s.ValidateAddProjectOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/add_project", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveProject removes a project from an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/remove_project.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) RemoveProject(ctx context.Context, opt *ApplicationsRemoveProjectOptions) (*http.Response, error) {
	err := s.ValidateRemoveProjectOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/remove_project", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// CreateBranch creates a new branch on an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/create_branch.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) CreateBranch(ctx context.Context, opt *ApplicationsCreateBranchOptions) (*http.Response, error) {
	err := s.ValidateCreateBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/create_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteBranch deletes a branch on an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/delete_branch.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) DeleteBranch(ctx context.Context, opt *ApplicationsDeleteBranchOptions) (*http.Response, error) {
	err := s.ValidateDeleteBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/delete_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateBranch updates a branch on an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: POST /api/applications/update_branch.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) UpdateBranch(ctx context.Context, opt *ApplicationsUpdateBranchOptions) (*http.Response, error) {
	err := s.ValidateUpdateBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/update_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetTags sets tags on an application.
// Requires 'Administer' permission on the application.
//
// API endpoint: POST /api/applications/set_tags.
// Since: 8.3.
// Enterprise Edition only.
func (s *ApplicationsService) SetTags(ctx context.Context, opt *ApplicationsSetTagsOptions) (*http.Response, error) {
	err := s.ValidateSetTagsOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/set_tags", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SearchProjects lists projects manually selected in an application.
// Requires 'Administrator' permission on the application.
//
// API endpoint: GET /api/applications/search_projects.
// Since: 7.3.
// Enterprise Edition only.
func (s *ApplicationsService) SearchProjects(ctx context.Context, opt *ApplicationsSearchProjectsOptions) (*ApplicationsSearchProjects, *http.Response, error) {
	err := s.ValidateSearchProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "applications/search_projects", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ApplicationsSearchProjects)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Refresh triggers a recomputation of an application's measures.
// Requires 'Administer System' or 'Administer' rights on the application.
//
// API endpoint: POST /api/applications/refresh.
// Since: 8.6.
// Enterprise Edition only.
func (s *ApplicationsService) Refresh(ctx context.Context, opt *ApplicationsRefreshOptions) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "applications/refresh", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
