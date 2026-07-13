package sonar

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// JiraResourceTypeSonarIssue represents the "SONAR_ISSUE" Jira work item resource type.
	JiraResourceTypeSonarIssue = "SONAR_ISSUE"
	// JiraResourceTypeDependencyRisk represents the "DEPENDENCY_RISK" Jira work item resource type.
	JiraResourceTypeDependencyRisk = "DEPENDENCY_RISK"
)

// JiraService handles communication with the Jira integration related
// methods of the SonarQube V2 API.
// This service is only available in Enterprise Edition with the Jira
// integration feature enabled.
type JiraService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant set of allowed values
var (
	// allowedJiraResourceTypes is the set of supported Jira work item resource types.
	allowedJiraResourceTypes = map[string]struct{}{
		JiraResourceTypeSonarIssue:     {},
		JiraResourceTypeDependencyRisk: {},
	}
)

// -----------------------------------------------------------------------------
// Shared / Response Types
// -----------------------------------------------------------------------------
//
// The SonarQube V2 API spec does not publish concrete schemas for the Jira
// integration request/response bodies (they are declared as an empty schema
// `{}`). Following the existing convention used elsewhere in this codebase
// for unspecified schemas (see e.g. SupportInfo in support_service.go), these
// are represented as loosely-typed maps.

// JiraWorkItem represents a Jira work item linked to a SonarQube issue or
// dependency risk. The exact shape of this object is not published by the
// SonarQube API.
type JiraWorkItem map[string]any

// JiraProjectBinding represents the binding between a SonarQube project and a
// Jira project. The exact shape of this object is not published by the
// SonarQube API.
type JiraProjectBinding map[string]any

// JiraOrganizationBinding represents the binding between a SonarQube instance
// (or organization) and a Jira/Atlassian organization. The exact shape of
// this object is not published by the SonarQube API.
type JiraOrganizationBinding map[string]any

// JiraWorkType represents a Jira work (issue) type available for a Jira
// project, optionally including field metadata. The exact shape of this
// object is not published by the SonarQube API.
type JiraWorkType map[string]any

// JiraWorkTypeSelection represents the payload used to save the Jira work
// types selected for a project. The exact shape of this object is not
// published by the SonarQube API.
type JiraWorkTypeSelection map[string]any

// JiraProject represents a Jira project available for binding. The exact
// shape of this object is not published by the SonarQube API.
type JiraProject map[string]any

// JiraLinkedIssuesCount represents the count of Jira issues linked to a
// SonarQube project. The exact shape of this object is not published by the
// SonarQube API.
type JiraLinkedIssuesCount map[string]any

// -----------------------------------------------------------------------------
// Option Types (Query Parameters)
// -----------------------------------------------------------------------------

// JiraWorkItemsOptions contains parameters identifying the Jira work items
// linked to a SonarQube resource.
type JiraWorkItemsOptions struct {
	// SonarProjectId is the identifier of the Sonar project. This field is required.
	SonarProjectId string `json:"sonarProjectId"`
	// ResourceId is the identifier of the resource (Sonar issue or dependency risk). This field is required.
	ResourceId string `json:"resourceId"`
	// ResourceType is the type of the resource. This field is required.
	// Allowed values: SONAR_ISSUE, DEPENDENCY_RISK.
	ResourceType string `json:"resourceType"`
}

// JiraProjectBindingOptions contains parameters identifying the Jira project
// binding of a SonarQube project.
type JiraProjectBindingOptions struct {
	// SonarProjectId is the identifier of the Sonar project. This field is required.
	SonarProjectId string `json:"sonarProjectId"`
}

// JiraOrganizationBindingOptions contains parameters identifying the Jira
// instance binding.
type JiraOrganizationBindingOptions struct {
	// SonarOrganizationUuid is the UUID of the Sonar organization/instance. This field is required.
	SonarOrganizationUuid string `json:"sonarOrganizationUuid"`
}

// JiraWorkTypesOptions contains parameters for listing the Jira work types
// available for a Jira project.
type JiraWorkTypesOptions struct {
	// JiraProjectKey is the key of the Jira project. This field is required.
	JiraProjectKey string `json:"jiraProjectKey"`
	// SonarOrganizationUuid is the UUID of the Sonar organization/instance. This field is required.
	SonarOrganizationUuid string `json:"sonarOrganizationUuid"`
	// SonarProjectId is the identifier of the Sonar project, used to determine which work
	// types are already selected. Optional.
	SonarProjectId string `json:"sonarProjectId,omitempty"`
	// IncludeFields indicates whether field metadata should be included in the response. Optional.
	IncludeFields bool `json:"includeFields,omitempty"`
}

// JiraUserActionsOptions contains parameters for listing the Jira related
// user actions available for a project.
type JiraUserActionsOptions struct {
	// SonarProjectId is the identifier of the Sonar project. This field is required.
	SonarProjectId string `json:"sonarProjectId"`
}

// JiraProjectsOptions contains parameters for listing the Jira projects
// available for binding.
type JiraProjectsOptions struct {
	// SonarOrganizationUuid is the UUID of the Sonar organization/instance. This field is required.
	SonarOrganizationUuid string `json:"sonarOrganizationUuid"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateWorkItemsOpt validates the JiraWorkItemsOptions.
func (s *JiraService) ValidateWorkItemsOpt(opt *JiraWorkItemsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.SonarProjectId, "SonarProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ResourceId, "ResourceId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ResourceType, "ResourceType")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.ResourceType, allowedJiraResourceTypes, "ResourceType")
}

// ValidateProjectBindingOpt validates the JiraProjectBindingOptions.
func (s *JiraService) ValidateProjectBindingOpt(opt *JiraProjectBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.SonarProjectId, "SonarProjectId")
}

// ValidateOrganizationBindingOpt validates the JiraOrganizationBindingOptions.
func (s *JiraService) ValidateOrganizationBindingOpt(opt *JiraOrganizationBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.SonarOrganizationUuid, "SonarOrganizationUuid")
}

// ValidateWorkTypesOpt validates the JiraWorkTypesOptions.
func (s *JiraService) ValidateWorkTypesOpt(opt *JiraWorkTypesOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.JiraProjectKey, "JiraProjectKey")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.SonarOrganizationUuid, "SonarOrganizationUuid")
}

// ValidateUserActionsOpt validates the JiraUserActionsOptions.
func (s *JiraService) ValidateUserActionsOpt(opt *JiraUserActionsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.SonarProjectId, "SonarProjectId")
}

// ValidateProjectsOpt validates the JiraProjectsOptions.
func (s *JiraService) ValidateProjectsOpt(opt *JiraProjectsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.SonarOrganizationUuid, "SonarOrganizationUuid")
}

// -----------------------------------------------------------------------------
// Service Methods - Work Items
// -----------------------------------------------------------------------------

// GetWorkItems fetches the Jira work items for a specific Sonar project and resource.
// Accepts only authenticated requests.
func (s *JiraService) GetWorkItems(ctx context.Context, opt *JiraWorkItemsOptions) (*JiraWorkItem, *http.Response, error) {
	err := s.ValidateWorkItemsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/work-items", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraWorkItem)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateWorkItem creates a Jira work item for a specific resource.
// Requires the issue administration permission.
func (s *JiraService) CreateWorkItem(ctx context.Context, body JiraWorkItem) (*JiraWorkItem, *http.Response, error) {
	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "jira/work-items", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraWorkItem)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteWorkItems deletes the Jira work items associated with a specific resource.
func (s *JiraService) DeleteWorkItems(ctx context.Context, opt *JiraWorkItemsOptions) (*http.Response, error) {
	err := s.ValidateWorkItemsOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "jira/work-items", opt, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------
// Service Methods - Project Bindings
// -----------------------------------------------------------------------------

// GetProjectBinding fetches the Jira project binding for a specific Sonar project.
// Accepts only authenticated requests.
func (s *JiraService) GetProjectBinding(ctx context.Context, opt *JiraProjectBindingOptions) (*JiraProjectBinding, *http.Response, error) {
	err := s.ValidateProjectBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/project-bindings", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraProjectBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateProjectBinding creates or updates a Jira project binding for a specific Sonar project.
// Accepts only authenticated requests.
func (s *JiraService) CreateProjectBinding(ctx context.Context, body JiraProjectBinding) (*JiraProjectBinding, *http.Response, error) {
	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "jira/project-bindings", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraProjectBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateProjectBinding updates an existing Jira project binding for a specific Sonar project.
// Accepts only authenticated requests.
func (s *JiraService) UpdateProjectBinding(ctx context.Context, body JiraProjectBinding) (*JiraProjectBinding, *http.Response, error) {
	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "jira/project-bindings", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraProjectBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteProjectBinding deletes the Jira project binding for a specific Sonar project.
// Accepts only authenticated requests.
func (s *JiraService) DeleteProjectBinding(ctx context.Context, opt *JiraProjectBindingOptions) (*http.Response, error) {
	err := s.ValidateProjectBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "jira/project-bindings", opt, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------
// Service Methods - Organization Bindings
// -----------------------------------------------------------------------------

// GetOrganizationBinding fetches the Jira instance binding.
// Requires the global administrator permission.
func (s *JiraService) GetOrganizationBinding(ctx context.Context, opt *JiraOrganizationBindingOptions) (*JiraOrganizationBinding, *http.Response, error) {
	err := s.ValidateOrganizationBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/organization-bindings", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraOrganizationBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateOrganizationBinding receives the 3LO state and authorization code parameters and
// attempts to create a Jira instance binding. If successful, returns the binding. Otherwise,
// returns a list of available resources.
// Requires the global administrator permission.
func (s *JiraService) CreateOrganizationBinding(ctx context.Context, body JiraOrganizationBinding) (*JiraOrganizationBinding, *http.Response, error) {
	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "jira/organization-bindings", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraOrganizationBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// BindOrganizationBinding binds the specified pending instance binding with the specified Jira cloud ID.
func (s *JiraService) BindOrganizationBinding(ctx context.Context, body JiraOrganizationBinding) (*JiraOrganizationBinding, *http.Response, error) {
	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "jira/organization-bindings", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraOrganizationBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteOrganizationBinding deletes the Jira instance binding from the database.
// Requires the global administrator permission.
func (s *JiraService) DeleteOrganizationBinding(ctx context.Context, opt *JiraOrganizationBindingOptions) (*http.Response, error) {
	err := s.ValidateOrganizationBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "jira/organization-bindings", opt, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateOrganizationBinding updates properties of an existing Jira organization binding, such as
// token sharing configuration.
// Requires the global administrator permission.
func (s *JiraService) UpdateOrganizationBinding(ctx context.Context, body JiraOrganizationBinding) (*JiraOrganizationBinding, *http.Response, error) {
	if body == nil {
		return nil, nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "jira/organization-binding-edit", nil, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraOrganizationBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// -----------------------------------------------------------------------------
// Service Methods - Work Types
// -----------------------------------------------------------------------------

// GetWorkTypes returns a list of all the available Jira work types for a specific Jira project.
// Also checks which work types are selected for a given Sonar project. Conditionally, also
// includes field metadata.
func (s *JiraService) GetWorkTypes(ctx context.Context, opt *JiraWorkTypesOptions) ([]JiraWorkType, *http.Response, error) {
	err := s.ValidateWorkTypesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/work-types", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []JiraWorkType

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateWorkTypes receives a Jira project key and a list of work types and stores the
// selected work types for the project.
func (s *JiraService) UpdateWorkTypes(ctx context.Context, body JiraWorkTypeSelection) (*http.Response, error) {
	if body == nil {
		return nil, NewValidationError("body", "must not be nil", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "jira/work-types", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------
// Service Methods - User Actions, Projects & Linked Issues
// -----------------------------------------------------------------------------

// GetUserActions returns a list of available user actions for the authenticated user in the
// context of a specific project. If the user is not authenticated, returns an empty list.
func (s *JiraService) GetUserActions(ctx context.Context, opt *JiraUserActionsOptions) ([]string, *http.Response, error) {
	err := s.ValidateUserActionsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/user-actions", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []string

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjects returns a list of all the available Jira projects.
func (s *JiraService) GetProjects(ctx context.Context, opt *JiraProjectsOptions) ([]JiraProject, *http.Response, error) {
	err := s.ValidateProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/projects", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []JiraProject

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetLinkedIssuesCount counts the number of Jira issues linked to a specific Sonar project.
// Accepts only authenticated requests.
func (s *JiraService) GetLinkedIssuesCount(ctx context.Context, sonarProjectID string) (*JiraLinkedIssuesCount, *http.Response, error) {
	err := ValidateRequired(sonarProjectID, "SonarProjectId")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "jira/linked-issues-count/"+sonarProjectID, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(JiraLinkedIssuesCount)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
