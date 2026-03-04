package sonar

import (
	"fmt"
	"net/http"
)

// DopTranslationServiceV2 handles communication with the DevOps Platform
// Translation related methods of the SonarQube V2 API.
type DopTranslationServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// DopSetting represents a DevOps Platform integration setting.
type DopSetting struct {
	// AppId is the application ID of the DevOps Platform setting.
	AppId string `json:"appId,omitempty"`
	// Id is the unique identifier of the DevOps Platform setting.
	Id string `json:"id,omitempty"`
	// Key is the key of the DevOps Platform setting.
	Key string `json:"key,omitempty"`
	// Type is the type of DevOps Platform (github, gitlab, azure, bitbucketcloud, bitbucket_server).
	Type string `json:"type,omitempty"`
	// Url is the URL of the DevOps Platform instance.
	Url string `json:"url,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// DopTranslationBoundProject represents the response from creating or updating a bound project.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationBoundProject struct {
	// BindingId is the identifier of the binding between the project and the DevOps platform.
	BindingId string `json:"bindingId,omitempty"`
	// NewProjectCreated is true if a new project was created, false if an existing project was bound.
	NewProjectCreated bool `json:"newProjectCreated,omitempty"`
	// ProjectId is the identifier of the created project.
	ProjectId string `json:"projectId,omitempty"`
}

// DopTranslationDopSettings represents the response from listing all DevOps Platform settings.
type DopTranslationDopSettings struct {
	// DopSettings is the list of DevOps Platform settings.
	DopSettings []DopSetting `json:"dopSettings,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// DopTranslationBoundProjectOptions contains parameters for creating or updating a bound project.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DopTranslationBoundProjectOptions struct {
	// DevOpsPlatformSettingId is the identifier of the DevOps platform configuration.
	// This field is required.
	DevOpsPlatformSettingId string `json:"devOpsPlatformSettingId"`
	// Monorepo indicates whether the project is part of a monorepo.
	// This field is required.
	Monorepo bool `json:"monorepo"`
	// NewCodeDefinitionType is the project new code definition type.
	// Allowed values: PREVIOUS_VERSION, NUMBER_OF_DAYS, REFERENCE_BRANCH.
	NewCodeDefinitionType string `json:"newCodeDefinitionType,omitempty"`
	// NewCodeDefinitionValue is the project new code definition value.
	NewCodeDefinitionValue string `json:"newCodeDefinitionValue,omitempty"`
	// ProjectIdentifier is the identifier of the DevOps platform project.
	// Only needed for Azure and BitBucket Server platforms.
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
	// ProjectKey is the key of the project to create.
	// This field is required.
	ProjectKey string `json:"projectKey"`
	// ProjectName is the name of the project to create.
	// This field is required.
	ProjectName string `json:"projectName"`
	// RepositoryIdentifier is the identifier of the DevOps platform repository to import.
	// This field is required.
	RepositoryIdentifier string `json:"repositoryIdentifier"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateCreateBoundProjectRequest validates the DopTranslationBoundProjectOptions.
func (s *DopTranslationServiceV2) ValidateCreateBoundProjectRequest(opt *DopTranslationBoundProjectOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.DevOpsPlatformSettingId, "DevOpsPlatformSettingId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectName, "ProjectName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RepositoryIdentifier, "RepositoryIdentifier")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CreateOrUpdateBoundProject creates a SonarQube project bound to a DevOps platform
// repository, or updates the binding if the project already exists.
// This is an idempotent operation.
// Requires the 'Create Projects' permission and a configured Personal Access Token.
func (s *DopTranslationServiceV2) CreateOrUpdateBoundProject(opt *DopTranslationBoundProjectOptions) (*DopTranslationBoundProject, *http.Response, error) {
	err := s.ValidateCreateBoundProjectRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPut, "dop-translation/bound-projects", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DopTranslationBoundProject)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateBoundProject creates a SonarQube project with the information from the
// provided DevOps platform project.
// Requires the 'Create Projects' permission and a configured Personal Access Token.
func (s *DopTranslationServiceV2) CreateBoundProject(opt *DopTranslationBoundProjectOptions) (*DopTranslationBoundProject, *http.Response, error) {
	err := s.ValidateCreateBoundProjectRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPost, "dop-translation/bound-projects", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DopTranslationBoundProject)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// FetchAllDopSettings lists all DevOps Platform Integration settings.
// Requires the 'Create Projects' permission.
func (s *DopTranslationServiceV2) FetchAllDopSettings() (*DopTranslationDopSettings, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "dop-translation/dop-settings", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DopTranslationDopSettings)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
