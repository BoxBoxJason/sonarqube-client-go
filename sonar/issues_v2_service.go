package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// IssuesV2Service handles communication with the Issues sandbox settings V2
// API endpoints (instance-level and project-level issue sandboxing
// configuration). Named IssuesV2Service (rather than IssuesService) to avoid
// colliding with the existing V1 IssuesService.
// This service is only available in Enterprise Edition. The underlying
// endpoints are marked internal by SonarQube (x-sonar-internal) and their
// request/response contract may change without notice between SonarQube
// versions.
type IssuesV2Service struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// IssuesSandboxSoftwareQualityMapping represents a software quality to impact
// severities mapping used by sandbox settings when the instance runs in MQR
// (Multi-Quality Rule) mode.
type IssuesSandboxSoftwareQualityMapping struct {
	// SoftwareQuality is the software quality. This field is required.
	// Allowed values: MAINTAINABILITY, RELIABILITY, SECURITY.
	SoftwareQuality string `json:"softwareQuality"`
	// ImpactSeverities is the list of impact severities. This field is required.
	// Allowed values: INFO, LOW, MEDIUM, HIGH, BLOCKER.
	ImpactSeverities []string `json:"impactSeverities"`
}

// IssuesSandboxRuleTypeMapping represents a rule type to severities mapping
// used by sandbox settings when the instance runs in Standard Experience mode.
type IssuesSandboxRuleTypeMapping struct {
	// Type is the rule type. This field is required.
	// Allowed values: CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT.
	Type string `json:"type"`
	// Severities is the list of severities. This field is required.
	// Allowed values: BLOCKER, CRITICAL, MAJOR, MINOR, INFO.
	Severities []string `json:"severities"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// IssuesSandboxSettings represents the instance-level issue sandbox settings.
type IssuesSandboxSettings struct {
	SoftwareQualities []IssuesSandboxSoftwareQualityMapping `json:"softwareQualities,omitempty"`
	Types             []IssuesSandboxRuleTypeMapping        `json:"types,omitempty"`
	Enabled           bool                                  `json:"enabled,omitempty"`
	DefaultValue      bool                                  `json:"defaultValue,omitempty"`
	AllowOverride     bool                                  `json:"allowOverride,omitempty"`
}

// IssuesSandboxProjectSettings represents the project-level issue sandbox settings.
type IssuesSandboxProjectSettings struct {
	SoftwareQualities []IssuesSandboxSoftwareQualityMapping `json:"softwareQualities,omitempty"`
	Types             []IssuesSandboxRuleTypeMapping        `json:"types,omitempty"`
	Enabled           bool                                  `json:"enabled,omitempty"`
	Overridden        bool                                  `json:"overridden,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// IssuesV2UpdateSandboxSettingsOptions contains parameters for updating the
// instance-level sandbox settings. All fields are optional (PATCH merge
// semantics): a nil field is left unchanged server-side.
type IssuesV2UpdateSandboxSettingsOptions struct {
	// Enabled enables or disables sandbox globally.
	Enabled *bool `json:"enabled,omitempty"`
	// DefaultValue sets the default value applied to new projects.
	DefaultValue *bool `json:"defaultValue,omitempty"`
	// AllowOverride allows or disallows projects to override instance settings.
	// Disabling it deletes all project-level software quality overrides.
	AllowOverride *bool `json:"allowOverride,omitempty"`
	// SoftwareQualities is the list of software quality mappings (MQR mode).
	SoftwareQualities []IssuesSandboxSoftwareQualityMapping `json:"softwareQualities,omitempty"`
	// Types is the list of rule type mappings (Standard Experience mode).
	Types []IssuesSandboxRuleTypeMapping `json:"types,omitempty"`
}

// IssuesV2UpdateProjectSandboxSettingsOptions contains parameters for updating
// the sandbox settings of a specific project. All fields are optional (PATCH
// merge semantics): a nil field is left unchanged server-side.
type IssuesV2UpdateProjectSandboxSettingsOptions struct {
	Enabled           *bool                                 `json:"enabled,omitempty"`
	Overridden        *bool                                 `json:"overridden,omitempty"`
	SoftwareQualities []IssuesSandboxSoftwareQualityMapping `json:"softwareQualities,omitempty"`
	Types             []IssuesSandboxRuleTypeMapping        `json:"types,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// validateSandboxSoftwareQualities validates a list of software quality mappings.
func validateSandboxSoftwareQualities(mappings []IssuesSandboxSoftwareQualityMapping) error {
	for idx, mapping := range mappings {
		err := ValidateRequired(mapping.SoftwareQuality, fmt.Sprintf("SoftwareQualities[%d].SoftwareQuality", idx))
		if err != nil {
			return err
		}

		err = IsValueAuthorized(mapping.SoftwareQuality, allowedImpactSoftwareQualities, fmt.Sprintf("SoftwareQualities[%d].SoftwareQuality", idx))
		if err != nil {
			return err
		}

		if len(mapping.ImpactSeverities) == 0 {
			return NewValidationError(fmt.Sprintf("SoftwareQualities[%d].ImpactSeverities", idx), "is required", ErrMissingRequired)
		}

		err = AreValuesAuthorized(mapping.ImpactSeverities, allowedRuleImpactSeverities, fmt.Sprintf("SoftwareQualities[%d].ImpactSeverities", idx))
		if err != nil {
			return err
		}
	}

	return nil
}

// validateSandboxRuleTypes validates a list of rule type mappings.
func validateSandboxRuleTypes(mappings []IssuesSandboxRuleTypeMapping) error {
	for idx, mapping := range mappings {
		err := ValidateRequired(mapping.Type, fmt.Sprintf("Types[%d].Type", idx))
		if err != nil {
			return err
		}

		err = IsValueAuthorized(mapping.Type, allowedRulesTypes, fmt.Sprintf("Types[%d].Type", idx))
		if err != nil {
			return err
		}

		if len(mapping.Severities) == 0 {
			return NewValidationError(fmt.Sprintf("Types[%d].Severities", idx), "is required", ErrMissingRequired)
		}

		err = AreValuesAuthorized(mapping.Severities, allowedRuleSeverities, fmt.Sprintf("Types[%d].Severities", idx))
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateSandboxSettingsOpt validates the options for the
// UpdateSandboxSettings method.
func (s *IssuesV2Service) ValidateUpdateSandboxSettingsOpt(opt *IssuesV2UpdateSandboxSettingsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := validateSandboxSoftwareQualities(opt.SoftwareQualities)
	if err != nil {
		return err
	}

	return validateSandboxRuleTypes(opt.Types)
}

// ValidateUpdateProjectSandboxSettingsOpt validates the parameters for the
// UpdateProjectSandboxSettings method.
func (s *IssuesV2Service) ValidateUpdateProjectSandboxSettingsOpt(projectKey string, opt *IssuesV2UpdateProjectSandboxSettingsOptions) error {
	err := ValidateRequired(projectKey, "ProjectKey")
	if err != nil {
		return err
	}

	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err = validateSandboxSoftwareQualities(opt.SoftwareQualities)
	if err != nil {
		return err
	}

	return validateSandboxRuleTypes(opt.Types)
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetSandboxSettings fetches the current instance-level sandbox settings.
// Returns configuration including enabled status, default value, allowOverride
// setting, and software qualities or rule types based on current mode.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/issues/sandbox-settings.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *IssuesV2Service) GetSandboxSettings(ctx context.Context) (*IssuesSandboxSettings, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "issues/sandbox-settings", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IssuesSandboxSettings)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateSandboxSettings updates the instance-level sandbox settings. Sandbox
// cannot be enabled without providing software qualities or rule types unless
// they already exist. When disabling allowOverride, all project-level
// software quality overrides are deleted.
// Requires 'Administer System' permission.
//
// API endpoint: PATCH /api/v2/issues/sandbox-settings.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *IssuesV2Service) UpdateSandboxSettings(ctx context.Context, opt *IssuesV2UpdateSandboxSettingsOptions) (*IssuesSandboxSettings, *http.Response, error) {
	err := s.ValidateUpdateSandboxSettingsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "issues/sandbox-settings", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IssuesSandboxSettings)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectSandboxSettings fetches the sandbox settings for a specific
// project. Returns effective configuration including enabled status, software
// quality overrides, and whether settings are inherited from instance or
// overridden at project level.
// Requires 'Administer Project' permission.
//
// API endpoint: GET /api/v2/issues/sandbox-settings/{projectKey}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *IssuesV2Service) GetProjectSandboxSettings(ctx context.Context, projectKey string) (*IssuesSandboxProjectSettings, *http.Response, error) {
	err := ValidateRequired(projectKey, "ProjectKey")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "issues/sandbox-settings/"+url.PathEscape(projectKey), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IssuesSandboxProjectSettings)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateProjectSandboxSettings updates the sandbox settings for a specific
// project. The enabled setting can be changed independently of software
// quality overrides. When allowOverride is disabled at instance level,
// software quality changes are not allowed but enabled changes are still
// permitted. Setting Overridden to false removes project-level software
// quality overrides to inherit from instance settings.
// Requires 'Administer Project' permission and the instance sandbox must be
// enabled.
//
// API endpoint: PATCH /api/v2/issues/sandbox-settings/{projectKey}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *IssuesV2Service) UpdateProjectSandboxSettings(ctx context.Context, projectKey string, opt *IssuesV2UpdateProjectSandboxSettingsOptions) (*IssuesSandboxProjectSettings, *http.Response, error) {
	err := s.ValidateUpdateProjectSandboxSettingsOpt(projectKey, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "issues/sandbox-settings/"+url.PathEscape(projectKey), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(IssuesSandboxProjectSettings)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
