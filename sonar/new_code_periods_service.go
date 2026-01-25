package sonargo

import (
	"net/http"
	"strconv"
)

// NewCodePeriodsService handles communication with the new code periods related methods
// of the SonarQube API.
// This service manages new code definitions.
type NewCodePeriodsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedNewCodePeriodTypes is the set of supported new code period types.
	allowedNewCodePeriodTypes = map[string]struct{}{
		"SPECIFIC_ANALYSIS": {},
		"PREVIOUS_VERSION":  {},
		"NUMBER_OF_DAYS":    {},
		"REFERENCE_BRANCH":  {},
	}
)

const (
	// maxNumberOfDays is the maximum allowed number of days for NUMBER_OF_DAYS type.
	maxNumberOfDays = 90
)

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// NewCodePeriod represents a new code period definition.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type NewCodePeriod struct {
	// BranchKey is the branch key.
	BranchKey string `json:"branchKey,omitempty"`
	// EffectiveValue is the effective value of the new code period.
	EffectiveValue string `json:"effectiveValue,omitempty"`
	// Inherited indicates whether the value is inherited from a parent.
	Inherited bool `json:"inherited,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Type is the type of the new code period.
	Type string `json:"type,omitempty"`
	// Value is the value of the new code period.
	Value string `json:"value,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// NewCodePeriodsList represents the response from listing new code periods.
type NewCodePeriodsList struct {
	// NewCodePeriods is the list of new code periods.
	NewCodePeriods []NewCodePeriod `json:"newCodePeriods,omitempty"`
}

// NewCodePeriodsShow represents the response from showing a new code period.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type NewCodePeriodsShow struct {
	// BranchKey is the branch key.
	BranchKey string `json:"branchKey,omitempty"`
	// Inherited indicates whether the value is inherited from a parent.
	Inherited bool `json:"inherited,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Type is the type of the new code period.
	Type string `json:"type,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// NewCodePeriodsListOption contains parameters for the List method.
type NewCodePeriodsListOption struct {
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// NewCodePeriodsSetOption contains parameters for the Set method.
type NewCodePeriodsSetOption struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Project is the project key.
	Project string `url:"project,omitempty"`
	// Type is the new code period type.
	// This field is required.
	// Allowed values: SPECIFIC_ANALYSIS, PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH.
	// - SPECIFIC_ANALYSIS: can be set at branch level only
	// - PREVIOUS_VERSION: can be set at any level (global, project, branch)
	// - NUMBER_OF_DAYS: can be set at any level (global, project, branch)
	// - REFERENCE_BRANCH: can only be set for projects and branches
	Type string `url:"type"`
	// Value is the new code period value.
	// For SPECIFIC_ANALYSIS: the uuid of an analysis
	// For PREVIOUS_VERSION: no value
	// For NUMBER_OF_DAYS: a number between 1 and 90
	// For REFERENCE_BRANCH: a string (branch name)
	Value string `url:"value,omitempty"`
}

// NewCodePeriodsShowOption contains parameters for the Show method.
type NewCodePeriodsShowOption struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Project is the project key.
	Project string `url:"project,omitempty"`
}

// NewCodePeriodsUnsetOption contains parameters for the Unset method.
type NewCodePeriodsUnsetOption struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Project is the project key.
	Project string `url:"project,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListOpt validates the options for the List method.
func (s *NewCodePeriodsService) ValidateListOpt(opt *NewCodePeriodsListOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetOpt validates the options for the Set method.
func (s *NewCodePeriodsService) ValidateSetOpt(opt *NewCodePeriodsSetOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Type, allowedNewCodePeriodTypes, "Type")
	if err != nil {
		return err
	}

	switch opt.Type {
	case "NUMBER_OF_DAYS":
		return s.validateNumberOfDays(opt)
	case "SPECIFIC_ANALYSIS":
		return s.validateSpecificAnalysis(opt)
	case "REFERENCE_BRANCH":
		return s.validateReferenceBranch(opt)
	case "PREVIOUS_VERSION":
		return s.validatePreviousVersion(opt)
	default:
		return NewValidationError("Type", "unsupported type", ErrInvalidValue)
	}
}

// ValidateShowOpt validates the options for the Show method.
func (s *NewCodePeriodsService) ValidateShowOpt(opt *NewCodePeriodsShowOption) error {
	// Options are optional; nothing to validate.
	return nil
}

// ValidateUnsetOpt validates the options for the Unset method.
func (s *NewCodePeriodsService) ValidateUnsetOpt(opt *NewCodePeriodsUnsetOption) error {
	// Options are optional; nothing to validate.
	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// List lists the new code definition for all branches in a project.
// Requires the permission to browse the project.
//
// API endpoint: GET /api/new_code_periods/list.
// Since: 8.0.
func (s *NewCodePeriodsService) List(opt *NewCodePeriodsListOption) (*NewCodePeriodsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "new_code_periods/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(NewCodePeriodsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Set updates the new code definition on different levels:
//   - Not providing a project key and a branch key will update the default value at global level
//   - Project key must be provided to update the value for a project
//   - Both project and branch keys must be provided to update the value for a branch
//
// Requires one of the following permissions:
//   - 'Administer System' to change the global setting
//   - 'Administer' rights on the specified project to change the project setting
//
// API endpoint: POST /api/new_code_periods/set.
// Since: 8.0.
func (s *NewCodePeriodsService) Set(opt *NewCodePeriodsSetOption) (*http.Response, error) {
	err := s.ValidateSetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "new_code_periods/set", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Show shows the new code definition.
// If the component requested doesn't exist or if no new code definition is set for it,
// a value is inherited from the project or from the global setting.
//
// Requires one of the following permissions if a component is specified:
//   - 'Administer' rights on the specified component
//   - 'Execute analysis' rights on the specified component
//
// API endpoint: GET /api/new_code_periods/show.
// Since: 8.0.
func (s *NewCodePeriodsService) Show(opt *NewCodePeriodsShowOption) (*NewCodePeriodsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "new_code_periods/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(NewCodePeriodsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Unset unsets the new code definition for a branch, project or global.
// It requires the inherited New Code Definition to be compatible with the Clean as You Code methodology.
//
// Requires one of the following permissions:
//   - 'Administer System' to change the global setting
//   - 'Administer' rights for a specified component
//
// API endpoint: POST /api/new_code_periods/unset.
// Since: 8.0.
func (s *NewCodePeriodsService) Unset(opt *NewCodePeriodsUnsetOption) (*http.Response, error) {
	err := s.ValidateUnsetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "new_code_periods/unset", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// validateNumberOfDays validates the NUMBER_OF_DAYS type.
func (s *NewCodePeriodsService) validateNumberOfDays(opt *NewCodePeriodsSetOption) error {
	// Convert Value to int64 and validate range
	intValue, parseErr := strconv.ParseInt(opt.Value, 10, 64)
	if parseErr != nil {
		return NewValidationError("Value", "must be a valid number", ErrInvalidValue)
	}
	// Value must be a number between 1 and 90
	return ValidateRange(intValue, 1, maxNumberOfDays, "Value")
}

// validateSpecificAnalysis validates the SPECIFIC_ANALYSIS type.
func (s *NewCodePeriodsService) validateSpecificAnalysis(opt *NewCodePeriodsSetOption) error {
	// Branch is required
	return ValidateRequired(opt.Branch, "Branch")
}

// validateReferenceBranch validates the REFERENCE_BRANCH type.
func (s *NewCodePeriodsService) validateReferenceBranch(opt *NewCodePeriodsSetOption) error {
	// Project is required
	return ValidateRequired(opt.Project, "Project")
}

// validatePreviousVersion validates the PREVIOUS_VERSION type.
func (s *NewCodePeriodsService) validatePreviousVersion(opt *NewCodePeriodsSetOption) error {
	// No special requirements
	return nil
}
