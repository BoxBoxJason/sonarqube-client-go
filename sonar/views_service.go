package sonar

import (
	"context"
	"net/http"
)

// ViewsService handles communication with the portfolio related methods of
// the SonarQube API. This service is only available in Enterprise Edition.
type ViewsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// View represents a SonarQube portfolio/view.
type View struct {
	// Description is the portfolio description.
	Description string `json:"description,omitempty"`
	// Key is the portfolio key.
	Key string `json:"key,omitempty"`
	// Name is the portfolio name.
	Name string `json:"name,omitempty"`
	// Qualifier is the qualifier of the component (VW for portfolio, SVW for sub-portfolio).
	Qualifier string `json:"qualifier,omitempty"`
	// Visibility is the visibility of the portfolio (public or private).
	Visibility string `json:"visibility,omitempty"`
}

// ViewDetails is an extended portfolio with sub-portfolios and selection modes.
type ViewDetails struct {
	// Description is the portfolio description.
	Description string `json:"description,omitempty"`
	// Key is the portfolio key.
	Key string `json:"key,omitempty"`
	// Name is the portfolio name.
	Name string `json:"name,omitempty"`
	// Qualifier is the qualifier (VW or SVW).
	Qualifier string `json:"qualifier,omitempty"`
	// Visibility is the visibility (public or private).
	Visibility string `json:"visibility,omitempty"`
	// SelectionMode describes how projects are selected for the portfolio.
	SelectionMode string `json:"selectionMode,omitempty"`
	// SubViews is the list of sub-portfolios.
	SubViews []View `json:"subViews,omitempty"`
}

// ViewProject represents a project entry in a portfolio.
type ViewProject struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Selected indicates whether the project is selected in the portfolio.
	Selected bool `json:"selected,omitempty"`
}

// ViewProjectStatus represents a project with its quality gate status in a portfolio.
type ViewProjectStatus struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Status is the quality gate status of the project.
	Status string `json:"status,omitempty"`
	// BranchKey is the branch key used for the analysis.
	BranchKey string `json:"branchKey,omitempty"`
}

// ViewApplication represents an application entry in a portfolio.
type ViewApplication struct {
	// Key is the application key.
	Key string `json:"key,omitempty"`
	// Name is the application name.
	Name string `json:"name,omitempty"`
}

// ViewDestination represents a possible destination for moving a portfolio.
type ViewDestination struct {
	// Key is the destination portfolio key.
	Key string `json:"key,omitempty"`
	// Name is the destination portfolio name.
	Name string `json:"name,omitempty"`
}

// ViewsList represents the response from the list endpoint.
type ViewsList struct {
	// Views is the list of root portfolios.
	Views []View `json:"views,omitempty"`
}

// ViewsSearch represents the response from the search endpoint.
type ViewsSearch struct {
	// Components is the list of portfolios matching the search.
	Components []View `json:"components,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// ViewsShow represents the response from the show endpoint.
type ViewsShow struct {
	// Portfolio contains the portfolio details.
	Portfolio ViewDetails `json:"portfolio,omitzero"`
}

// ViewsProjects represents the response from the projects endpoint.
type ViewsProjects struct {
	// Projects is the list of projects in the portfolio.
	Projects []ViewProject `json:"projects,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// ViewsProjectsStatus represents the response from the projects_status endpoint.
type ViewsProjectsStatus struct {
	// Projects is the list of projects with their quality gate status.
	Projects []ViewProjectStatus `json:"projects,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// ViewsApplications represents the response from the applications endpoint.
type ViewsApplications struct {
	// Applications is the list of applications in the portfolio.
	Applications []ViewApplication `json:"applications,omitempty"`
}

// ViewsSubViews represents the response from the portfolios endpoint.
type ViewsSubViews struct {
	// SubViews is the list of sub-portfolios.
	SubViews []View `json:"subViews,omitempty"`
}

// ViewsMoveDestinations represents the response from the move_options endpoint.
type ViewsMoveDestinations struct {
	// Views is the list of possible destination portfolios.
	Views []ViewDestination `json:"views,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ViewsCreateOptions contains parameters for the Create method.
type ViewsCreateOptions struct {
	// Description is an optional portfolio description.
	Description string `url:"description,omitempty"`
	// Key is the portfolio key. If not provided, it will be generated from the name.
	Key string `url:"key,omitempty"`
	// Name is the portfolio name. This field is required.
	Name string `url:"name"`
	// Parent is the key of the parent portfolio, when creating a sub-portfolio.
	Parent string `url:"parent,omitempty"`
	// Visibility is the portfolio visibility (public or private).
	Visibility string `url:"visibility,omitempty"`
}

// ViewsDeleteOptions contains parameters for the Delete method.
type ViewsDeleteOptions struct {
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
}

// ViewsSearchOptions contains parameters for the Search method.
//
//nolint:govet // Field alignment is less important than logical grouping
type ViewsSearchOptions struct {
	PaginationArgs

	// Query limits search to portfolios whose name or key contains this value.
	Query string `url:"q,omitempty"`
	// Qualifiers filters by component qualifier (VW for portfolio, SVW for sub-portfolio).
	Qualifiers string `url:"qualifiers,omitempty"`
	// OnlyFavorites if true restricts results to favorite portfolios only.
	OnlyFavorites bool `url:"onlyFavorites,omitempty"`
}

// ViewsShowOptions contains parameters for the Show method.
type ViewsShowOptions struct {
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
}

// ViewsUpdateOptions contains parameters for the Update method.
type ViewsUpdateOptions struct {
	// Description is the new description for the portfolio.
	Description string `url:"description,omitempty"`
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
	// Name is the new name for the portfolio. This field is required.
	Name string `url:"name"`
}

// ViewsAddProjectOptions contains parameters for the AddProject method.
type ViewsAddProjectOptions struct {
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ViewsRemoveProjectOptions contains parameters for the RemoveProject method.
type ViewsRemoveProjectOptions struct {
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ViewsAddProjectBranchOptions contains parameters for the AddProjectBranch method.
type ViewsAddProjectBranchOptions struct {
	// Branch is the branch key. This field is required.
	Branch string `url:"branch"`
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ViewsRemoveProjectBranchOptions contains parameters for the RemoveProjectBranch method.
type ViewsRemoveProjectBranchOptions struct {
	// Branch is the branch key. This field is required.
	Branch string `url:"branch"`
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ViewsAddPortfolioOptions contains parameters for the AddPortfolio method
// (adding a reference portfolio as a sub-portfolio).
type ViewsAddPortfolioOptions struct {
	// Portfolio is the parent portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
	// Reference is the key of the portfolio to add as a sub-portfolio. This field is required.
	Reference string `url:"reference"`
}

// ViewsRemovePortfolioOptions contains parameters for the RemovePortfolio method.
type ViewsRemovePortfolioOptions struct {
	// Portfolio is the parent portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
	// Reference is the key of the sub-portfolio to remove. This field is required.
	Reference string `url:"reference"`
}

// ViewsMoveOptions contains parameters for the Move method.
type ViewsMoveOptions struct {
	// Destination is the key of the destination (new parent) portfolio. This field is required.
	Destination string `url:"destination"`
	// Key is the portfolio key to move. This field is required.
	Key string `url:"key"`
}

// ViewsMoveOptionsOptions contains parameters for the MoveOptions method.
type ViewsMoveOptionsOptions struct {
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
}

// ViewsProjectsOptions contains parameters for the Projects method.
//
//nolint:govet // Field alignment is less important than logical grouping
type ViewsProjectsOptions struct {
	PaginationArgs

	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
	// Query limits results to projects whose key contains this value.
	Query string `url:"query,omitempty"`
	// Selected filters on selected, deselected or all projects.
	Selected string `url:"selected,omitempty"`
}

// ViewsApplicationsOptions contains parameters for the Applications method.
type ViewsApplicationsOptions struct {
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsSubViewsOptions contains parameters for the SubPortfolios method.
type ViewsSubViewsOptions struct {
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsProjectsStatusOptions contains parameters for the ProjectsStatus method.
//
//nolint:govet // Field alignment is less important than logical grouping
type ViewsProjectsStatusOptions struct {
	PaginationArgs

	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
	// Status filters projects by quality gate status.
	Status string `url:"status,omitempty"`
}

// ViewsRefreshOptions contains parameters for the Refresh method.
type ViewsRefreshOptions struct {
	// Key is the root portfolio key. If not specified, all portfolios are refreshed.
	Key string `url:"key,omitempty"`
}

// ViewsAddApplicationOptions contains parameters for the AddApplication method.
type ViewsAddApplicationOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsRemoveApplicationOptions contains parameters for the RemoveApplication method.
type ViewsRemoveApplicationOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsAddApplicationBranchOptions contains parameters for the AddApplicationBranch method.
type ViewsAddApplicationBranchOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the branch key. This field is required.
	Branch string `url:"branch"`
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
}

// ViewsRemoveApplicationBranchOptions contains parameters for the RemoveApplicationBranch method.
type ViewsRemoveApplicationBranchOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the branch key. This field is required.
	Branch string `url:"branch"`
	// Key is the portfolio key. This field is required.
	Key string `url:"key"`
}

// ViewsSetManualModeOptions contains parameters for the SetManualMode method.
type ViewsSetManualModeOptions struct {
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsSetNoneModeOptions contains parameters for the SetNoneMode method.
type ViewsSetNoneModeOptions struct {
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsSetRegexpModeOptions contains parameters for the SetRegexpMode method.
type ViewsSetRegexpModeOptions struct {
	// Branch selects a branch in all matched projects instead of using main branches.
	Branch string `url:"branch,omitempty"`
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
	// Regexp is a valid Java regular expression for project key matching. This field is required.
	Regexp string `url:"regexp"`
}

// ViewsSetRemainingProjectsModeOptions contains parameters for the SetRemainingProjectsMode method.
type ViewsSetRemainingProjectsModeOptions struct {
	// Branch selects a branch in all matched projects instead of using main branches.
	Branch string `url:"branch,omitempty"`
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
}

// ViewsSetTagsModeOptions contains parameters for the SetTagsMode method.
type ViewsSetTagsModeOptions struct {
	// Branch selects a branch in all matched projects instead of using main branches.
	Branch string `url:"branch,omitempty"`
	// Portfolio is the portfolio key. This field is required.
	Portfolio string `url:"portfolio"`
	// Tags is the list of project tags to match. This field is required.
	Tags []string `url:"tags,comma"`
}

// -----------------------------------------------------------------------------
// Allowed Value Sets
// -----------------------------------------------------------------------------

//nolint:gochecknoglobals // constant set of allowed values
var allowedViewVisibilities = map[string]struct{}{
	"private": {},
	"public":  {},
}

//nolint:gochecknoglobals // constant set of allowed values
var allowedViewProjectSelections = map[string]struct{}{
	SelectionFilterAll:        {},
	SelectionFilterSelected:   {},
	SelectionFilterDeselected: {},
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCreateOpt validates the options for the Create method.
func (s *ViewsService) ValidateCreateOpt(opt *ViewsCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Visibility, allowedViewVisibilities, "Visibility")
}

// ValidateDeleteOpt validates the options for the Delete method.
func (s *ViewsService) ValidateDeleteOpt(opt *ViewsDeleteOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateShowOpt validates the options for the Show method.
func (s *ViewsService) ValidateShowOpt(opt *ViewsShowOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateUpdateOpt validates the options for the Update method.
func (s *ViewsService) ValidateUpdateOpt(opt *ViewsUpdateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Name, "Name")
}

// ValidateAddProjectOpt validates the options for the AddProject method.
func (s *ViewsService) ValidateAddProjectOpt(opt *ViewsAddProjectOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateRemoveProjectOpt validates the options for the RemoveProject method.
func (s *ViewsService) ValidateRemoveProjectOpt(opt *ViewsRemoveProjectOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateAddProjectBranchOpt validates the options for the AddProjectBranch method.
func (s *ViewsService) ValidateAddProjectBranchOpt(opt *ViewsAddProjectBranchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Branch, "Branch")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateRemoveProjectBranchOpt validates the options for the RemoveProjectBranch method.
func (s *ViewsService) ValidateRemoveProjectBranchOpt(opt *ViewsRemoveProjectBranchOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Branch, "Branch")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateAddPortfolioOpt validates the options for the AddPortfolio method.
func (s *ViewsService) ValidateAddPortfolioOpt(opt *ViewsAddPortfolioOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Reference, "Reference")
}

// ValidateRemovePortfolioOpt validates the options for the RemovePortfolio method.
func (s *ViewsService) ValidateRemovePortfolioOpt(opt *ViewsRemovePortfolioOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Reference, "Reference")
}

// ValidateMoveOpt validates the options for the Move method.
func (s *ViewsService) ValidateMoveOpt(opt *ViewsMoveOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Destination, "Destination")
}

// ValidateMoveOptionsOpt validates the options for the MoveOptions method.
func (s *ViewsService) ValidateMoveOptionsOpt(opt *ViewsMoveOptionsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateProjectsOpt validates the options for the Projects method.
func (s *ViewsService) ValidateProjectsOpt(opt *ViewsProjectsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = opt.Validate()
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Selected, allowedViewProjectSelections, "Selected")
}

// ValidateApplicationsOpt validates the options for the Applications method.
func (s *ViewsService) ValidateApplicationsOpt(opt *ViewsApplicationsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Portfolio, "Portfolio")
}

// ValidateSubPortfoliosOpt validates the options for the SubPortfolios method.
func (s *ViewsService) ValidateSubPortfoliosOpt(opt *ViewsSubViewsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Portfolio, "Portfolio")
}

// ValidateProjectsStatusOpt validates the options for the ProjectsStatus method.
func (s *ViewsService) ValidateProjectsStatusOpt(opt *ViewsProjectsStatusOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	return opt.Validate()
}

// ValidateAddApplicationOpt validates the options for the AddApplication method.
func (s *ViewsService) ValidateAddApplicationOpt(opt *ViewsAddApplicationOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Application, "Application")
}

// ValidateRemoveApplicationOpt validates the options for the RemoveApplication method.
func (s *ViewsService) ValidateRemoveApplicationOpt(opt *ViewsRemoveApplicationOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Application, "Application")
}

// ValidateAddApplicationBranchOpt validates the options for the AddApplicationBranch method.
func (s *ViewsService) ValidateAddApplicationBranchOpt(opt *ViewsAddApplicationBranchOptions) error {
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

	return ValidateRequired(opt.Key, "Key")
}

// ValidateRemoveApplicationBranchOpt validates the options for the RemoveApplicationBranch method.
func (s *ViewsService) ValidateRemoveApplicationBranchOpt(opt *ViewsRemoveApplicationBranchOptions) error {
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

	return ValidateRequired(opt.Key, "Key")
}

// ValidateSetManualModeOpt validates the options for the SetManualMode method.
func (s *ViewsService) ValidateSetManualModeOpt(opt *ViewsSetManualModeOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Portfolio, "Portfolio")
}

// ValidateSetNoneModeOpt validates the options for the SetNoneMode method.
func (s *ViewsService) ValidateSetNoneModeOpt(opt *ViewsSetNoneModeOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Portfolio, "Portfolio")
}

// ValidateSetRegexpModeOpt validates the options for the SetRegexpMode method.
func (s *ViewsService) ValidateSetRegexpModeOpt(opt *ViewsSetRegexpModeOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Regexp, "Regexp")
}

// ValidateSetRemainingProjectsModeOpt validates the options for the SetRemainingProjectsMode method.
func (s *ViewsService) ValidateSetRemainingProjectsModeOpt(opt *ViewsSetRemainingProjectsModeOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Portfolio, "Portfolio")
}

// ValidateSetTagsModeOpt validates the options for the SetTagsMode method.
func (s *ViewsService) ValidateSetTagsModeOpt(opt *ViewsSetTagsModeOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Portfolio, "Portfolio")
	if err != nil {
		return err
	}

	if len(opt.Tags) == 0 {
		return NewValidationError("Tags", "is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// AddApplication adds an application to a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/add_application.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) AddApplication(ctx context.Context, opt *ViewsAddApplicationOptions) (*http.Response, error) {
	err := s.ValidateAddApplicationOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/add_application", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddApplicationBranch adds a specific application branch to a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/add_application_branch.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) AddApplicationBranch(ctx context.Context, opt *ViewsAddApplicationBranchOptions) (*http.Response, error) {
	err := s.ValidateAddApplicationBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/add_application_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddPortfolio adds an existing portfolio as a sub-portfolio (reference) to a parent portfolio.
// Requires 'Administer' permission on the parent portfolio.
//
// API endpoint: POST /api/views/add_portfolio.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) AddPortfolio(ctx context.Context, opt *ViewsAddPortfolioOptions) (*http.Response, error) {
	err := s.ValidateAddPortfolioOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/add_portfolio", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddProject adds a project to a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/add_project.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) AddProject(ctx context.Context, opt *ViewsAddProjectOptions) (*http.Response, error) {
	err := s.ValidateAddProjectOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/add_project", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddProjectBranch adds a specific project branch to a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/add_project_branch.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) AddProjectBranch(ctx context.Context, opt *ViewsAddProjectBranchOptions) (*http.Response, error) {
	err := s.ValidateAddProjectBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/add_project_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Applications returns the applications in a portfolio.
// Requires 'Browse' permission on the portfolio.
//
// API endpoint: GET /api/views/applications.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) Applications(ctx context.Context, opt *ViewsApplicationsOptions) (*ViewsApplications, *http.Response, error) {
	err := s.ValidateApplicationsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/applications", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsApplications)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Create creates a new root portfolio or sub-portfolio.
// Requires 'Administer System' permission (root) or 'Administer' permission on the parent.
//
// API endpoint: POST /api/views/create.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Create(ctx context.Context, opt *ViewsCreateOptions) (*http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/create", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Delete deletes a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/delete.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Delete(ctx context.Context, opt *ViewsDeleteOptions) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/delete", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// List returns the list of root portfolios.
// Requires authentication.
//
// API endpoint: GET /api/views/list.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) List(ctx context.Context) (*ViewsList, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/list", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Move moves a portfolio to a new parent portfolio.
// Requires 'Administer' permission on both portfolios.
//
// API endpoint: POST /api/views/move.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Move(ctx context.Context, opt *ViewsMoveOptions) (*http.Response, error) {
	err := s.ValidateMoveOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/move", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// MoveOptions returns the list of possible destination portfolios for a move operation.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: GET /api/views/move_options.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) MoveOptions(ctx context.Context, opt *ViewsMoveOptionsOptions) (*ViewsMoveDestinations, *http.Response, error) {
	err := s.ValidateMoveOptionsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/move_options", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsMoveDestinations)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SubPortfolios returns the sub-portfolios of a portfolio.
// Requires 'Browse' permission on the portfolio.
//
// API endpoint: GET /api/views/portfolios.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) SubPortfolios(ctx context.Context, opt *ViewsSubViewsOptions) (*ViewsSubViews, *http.Response, error) {
	err := s.ValidateSubPortfoliosOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/portfolios", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsSubViews)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Projects lists the projects in a portfolio.
// Requires 'Browse' permission on the portfolio.
//
// API endpoint: GET /api/views/projects.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Projects(ctx context.Context, opt *ViewsProjectsOptions) (*ViewsProjects, *http.Response, error) {
	err := s.ValidateProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/projects", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsProjects)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ProjectsStatus returns the quality gate status of projects in a portfolio.
// Requires 'Browse' permission on the portfolio.
//
// API endpoint: GET /api/views/projects_status.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) ProjectsStatus(ctx context.Context, opt *ViewsProjectsStatusOptions) (*ViewsProjectsStatus, *http.Response, error) {
	err := s.ValidateProjectsStatusOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/projects_status", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsProjectsStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Refresh triggers a computation of portfolio measures. If no key is provided,
// all portfolios are refreshed.
// Requires 'Administer System' permission.
//
// API endpoint: POST /api/views/refresh.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Refresh(ctx context.Context, opt *ViewsRefreshOptions) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/refresh", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveApplication removes an application from a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/remove_application.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) RemoveApplication(ctx context.Context, opt *ViewsRemoveApplicationOptions) (*http.Response, error) {
	err := s.ValidateRemoveApplicationOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/remove_application", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveApplicationBranch removes a specific application branch from a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/remove_application_branch.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) RemoveApplicationBranch(ctx context.Context, opt *ViewsRemoveApplicationBranchOptions) (*http.Response, error) {
	err := s.ValidateRemoveApplicationBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/remove_application_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemovePortfolio removes a sub-portfolio from a parent portfolio.
// Requires 'Administer' permission on the parent portfolio.
//
// API endpoint: POST /api/views/remove_portfolio.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) RemovePortfolio(ctx context.Context, opt *ViewsRemovePortfolioOptions) (*http.Response, error) {
	err := s.ValidateRemovePortfolioOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/remove_portfolio", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveProject removes a project from a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/remove_project.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) RemoveProject(ctx context.Context, opt *ViewsRemoveProjectOptions) (*http.Response, error) {
	err := s.ValidateRemoveProjectOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/remove_project", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveProjectBranch removes a specific project branch from a portfolio.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/remove_project_branch.
// Since: 7.1.
// Enterprise Edition only.
func (s *ViewsService) RemoveProjectBranch(ctx context.Context, opt *ViewsRemoveProjectBranchOptions) (*http.Response, error) {
	err := s.ValidateRemoveProjectBranchOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/remove_project_branch", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Search searches portfolios by name, key, or description.
// Requires authentication.
//
// API endpoint: GET /api/views/search.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Search(ctx context.Context, opt *ViewsSearchOptions) (*ViewsSearch, *http.Response, error) {
	if opt != nil {
		err := opt.Validate()
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetManualMode sets a portfolio to manual project selection mode.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/set_manual_mode.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) SetManualMode(ctx context.Context, opt *ViewsSetManualModeOptions) (*http.Response, error) {
	err := s.ValidateSetManualModeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/set_manual_mode", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetNoneMode sets a portfolio to no project selection (empty).
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/set_none_mode.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) SetNoneMode(ctx context.Context, opt *ViewsSetNoneModeOptions) (*http.Response, error) {
	err := s.ValidateSetNoneModeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/set_none_mode", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetRegexpMode sets a portfolio to regexp-based project selection mode.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/set_regexp_mode.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) SetRegexpMode(ctx context.Context, opt *ViewsSetRegexpModeOptions) (*http.Response, error) {
	err := s.ValidateSetRegexpModeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/set_regexp_mode", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetRemainingProjectsMode sets a portfolio to include all remaining (unclassified) projects.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/set_remaining_projects_mode.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) SetRemainingProjectsMode(ctx context.Context, opt *ViewsSetRemainingProjectsModeOptions) (*http.Response, error) {
	err := s.ValidateSetRemainingProjectsModeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/set_remaining_projects_mode", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SetTagsMode sets a portfolio to tag-based project selection mode.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/set_tags_mode.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) SetTagsMode(ctx context.Context, opt *ViewsSetTagsModeOptions) (*http.Response, error) {
	err := s.ValidateSetTagsModeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/set_tags_mode", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Show returns the details of a portfolio, including sub-portfolios.
// Requires 'Browse' permission on the portfolio.
//
// API endpoint: GET /api/views/show.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Show(ctx context.Context, opt *ViewsShowOptions) (*ViewsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "views/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ViewsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Update updates a portfolio's name and/or description.
// Requires 'Administer' permission on the portfolio.
//
// API endpoint: POST /api/views/update.
// Since: 6.6.
// Enterprise Edition only.
func (s *ViewsService) Update(ctx context.Context, opt *ViewsUpdateOptions) (*http.Response, error) {
	err := s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "views/update", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
