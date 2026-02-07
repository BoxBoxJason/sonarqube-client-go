package sonar

import "net/http"

// ProjectBranchesService handles communication with the project branches related methods
// of the SonarQube API.
// This service manages project branches.
type ProjectBranchesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// BranchStatus represents the status of a branch.
type BranchStatus struct {
	// QualityGateStatus is the quality gate status of the branch.
	QualityGateStatus string `json:"qualityGateStatus,omitempty"`
}

// Branch represents a project branch.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type Branch struct {
	// AnalysisDate is the date of the last analysis.
	AnalysisDate string `json:"analysisDate,omitempty"`
	// BranchID is the unique identifier of the branch.
	BranchID string `json:"branchId,omitempty"`
	// ExcludedFromPurge indicates whether the branch is excluded from automatic purge.
	ExcludedFromPurge bool `json:"excludedFromPurge,omitempty"`
	// IsMain indicates whether this is the main branch.
	IsMain bool `json:"isMain,omitempty"`
	// Name is the name of the branch.
	Name string `json:"name,omitempty"`
	// Status is the status of the branch.
	Status BranchStatus `json:"status,omitzero"`
	// Type is the type of the branch.
	Type string `json:"type,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectBranchesList represents the response from listing branches.
type ProjectBranchesList struct {
	// Branches is the list of branches.
	Branches []Branch `json:"branches,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectBranchesDeleteOption contains parameters for the Delete method.
type ProjectBranchesDeleteOption struct {
	// Branch is the branch key.
	// This field is required.
	Branch string `url:"branch"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// ProjectBranchesListOption contains parameters for the List method.
type ProjectBranchesListOption struct {
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// ProjectBranchesRenameOption contains parameters for the Rename method.
type ProjectBranchesRenameOption struct {
	// Name is the new name of the main branch.
	// This field is required. Maximum length is 255 characters.
	Name string `url:"name"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// ProjectBranchesSetAutomaticDeletionProtectionOption contains parameters for the SetAutomaticDeletionProtection method.
type ProjectBranchesSetAutomaticDeletionProtectionOption struct {
	// Branch is the branch key.
	// This field is required.
	Branch string `url:"branch"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
	// Value sets whether the branch should be protected from automatic deletion.
	// This field is required.
	Value bool `url:"value"`
}

// ProjectBranchesSetMainOption contains parameters for the SetMain method.
type ProjectBranchesSetMainOption struct {
	// Branch is the branch key.
	// This field is required.
	Branch string `url:"branch"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateDeleteOpt validates the options for the Delete method.
func (s *ProjectBranchesService) ValidateDeleteOpt(opt *ProjectBranchesDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Branch, "Branch")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for the List method.
func (s *ProjectBranchesService) ValidateListOpt(opt *ProjectBranchesListOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRenameOpt validates the options for the Rename method.
func (s *ProjectBranchesService) ValidateRenameOpt(opt *ProjectBranchesRenameOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxBranchNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetAutomaticDeletionProtectionOpt validates the options for the SetAutomaticDeletionProtection method.
func (s *ProjectBranchesService) ValidateSetAutomaticDeletionProtectionOpt(opt *ProjectBranchesSetAutomaticDeletionProtectionOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Branch, "Branch")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	// Value is a bool, so it's always valid (true or false)

	return nil
}

// ValidateSetMainOpt validates the options for the SetMain method.
func (s *ProjectBranchesService) ValidateSetMainOpt(opt *ProjectBranchesSetMainOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Branch, "Branch")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Delete deletes a non-main branch of a project or application.
// Requires 'Administer' rights on the specified project or application.
//
// API endpoint: POST /api/project_branches/delete.
// Since: 6.6.
func (s *ProjectBranchesService) Delete(opt *ProjectBranchesDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_branches/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// List lists the branches of a project or application.
// Requires 'Browse' or 'Execute analysis' rights on the specified project or application.
//
// API endpoint: GET /api/project_branches/list.
// Since: 6.6.
func (s *ProjectBranchesService) List(opt *ProjectBranchesListOption) (*ProjectBranchesList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_branches/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectBranchesList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Rename renames the main branch of a project or application.
// Requires 'Administer' permission on the specified project or application.
//
// API endpoint: POST /api/project_branches/rename.
// Since: 6.6.
func (s *ProjectBranchesService) Rename(opt *ProjectBranchesRenameOption) (*http.Response, error) {
	err := s.ValidateRenameOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_branches/rename", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetAutomaticDeletionProtection protects a specific branch from automatic deletion.
// Protection can't be disabled for the main branch.
// Requires 'Administer' permission on the specified project or application.
//
// API endpoint: POST /api/project_branches/set_automatic_deletion_protection.
// Since: 8.1.
func (s *ProjectBranchesService) SetAutomaticDeletionProtection(opt *ProjectBranchesSetAutomaticDeletionProtectionOption) (*http.Response, error) {
	err := s.ValidateSetAutomaticDeletionProtectionOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_branches/set_automatic_deletion_protection", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetMain allows setting a new main branch.
// Caution: only applicable on projects.
// Requires 'Administer' rights on the specified project or application.
//
// API endpoint: POST /api/project_branches/set_main.
// Since: 10.2.
func (s *ProjectBranchesService) SetMain(opt *ProjectBranchesSetMainOption) (*http.Response, error) {
	err := s.ValidateSetMainOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_branches/set_main", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
