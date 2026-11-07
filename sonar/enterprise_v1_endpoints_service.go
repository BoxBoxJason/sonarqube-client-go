package sonar

import (
	"context"
	"net/http"
)

// This file holds V1 endpoints that complement their respective services: AI code assurance,
// findings exports, application quality gate status, provisioning status and a few
// administration endpoints.

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AlmSettingsValidateBinding represents the response from validating a project binding.
type AlmSettingsValidateBinding struct {
	// Scope is the scope of the failing validation (e.g. GLOBAL or PROJECT).
	Scope string `json:"scope,omitempty"`
	// Errors is the list of validation errors. Empty when the binding is valid.
	Errors []AlmSettingsValidateBindingError `json:"errors,omitempty"`
}

// AlmSettingsValidateBindingError represents a validation error of a project binding.
type AlmSettingsValidateBindingError struct {
	// Msg is the error message.
	Msg string `json:"msg,omitempty"`
}

// GithubProvisioningActivation represents the response from the GitHub provisioning status endpoint.
type GithubProvisioningActivation struct {
	// Enabled indicates whether GitHub provisioning is enabled.
	Enabled bool `json:"enabled,omitempty"`
}

// IssuesGitlabSastExport represents a vulnerability report in the GitLab SAST JSON format.
type IssuesGitlabSastExport struct {
	// Version is the GitLab SAST report schema version.
	Version string `json:"version,omitempty"`
	// Scan describes the scan that produced the report.
	Scan IssuesGitlabSastScan `json:"scan,omitzero"`
	// Vulnerabilities is the list of vulnerabilities, in the GitLab SAST format.
	Vulnerabilities []map[string]any `json:"vulnerabilities,omitempty"`
}

// IssuesGitlabSastScan describes the scan of a GitLab SAST report.
type IssuesGitlabSastScan struct {
	// Analyzer is the analyzer that produced the report.
	Analyzer IssuesGitlabSastTool `json:"analyzer,omitzero"`
	// Scanner is the scanner that produced the report.
	Scanner IssuesGitlabSastTool `json:"scanner,omitzero"`
	// StartTime is the scan start time.
	StartTime string `json:"start_time,omitempty"` //nolint:tagliatelle // GitLab SAST format
	// EndTime is the scan end time.
	EndTime string `json:"end_time,omitempty"` //nolint:tagliatelle // GitLab SAST format
	// Status is the scan status.
	Status string `json:"status,omitempty"`
	// Type is the scan type (sast).
	Type string `json:"type,omitempty"`
	// Messages is the list of scan messages.
	Messages []map[string]any `json:"messages,omitempty"`
}

// IssuesGitlabSastTool describes an analyzer or scanner of a GitLab SAST report.
type IssuesGitlabSastTool struct {
	// Id is the tool identifier.
	Id string `json:"id,omitempty"`
	// Name is the tool name.
	Name string `json:"name,omitempty"`
	// Vendor is the tool vendor.
	Vendor IssuesGitlabSastVendor `json:"vendor,omitzero"`
	// Version is the tool version.
	Version string `json:"version,omitempty"`
}

// IssuesGitlabSastVendor describes the vendor of a GitLab SAST tool.
type IssuesGitlabSastVendor struct {
	// Name is the vendor name.
	Name string `json:"name,omitempty"`
}

// ProjectBranchesAiCodeAssurance represents the AI code assurance status of a project branch.
type ProjectBranchesAiCodeAssurance struct {
	// AiCodeAssurance is the AI code assurance status (e.g. NONE, AI_CODE_ASSURANCE_OFF).
	AiCodeAssurance string `json:"aiCodeAssurance,omitempty"`
}

// ProjectsContainsAiCode represents whether a project is flagged as containing AI code.
type ProjectsContainsAiCode struct {
	// ContainsAiCode indicates whether the project contains AI code.
	ContainsAiCode bool `json:"containsAiCode,omitempty"`
}

// ProjectsDetectedAiCode represents whether AI code was detected in a project.
type ProjectsDetectedAiCode struct {
	// DetectedAiCode indicates whether AI code was detected.
	DetectedAiCode bool `json:"detectedAiCode,omitempty"`
}

// ProjectsExportFindings represents all the findings (issues and hotspots) of a project branch.
type ProjectsExportFindings struct {
	// ExportFindings is the list of findings.
	ExportFindings []map[string]any `json:"export_findings,omitempty"` //nolint:tagliatelle // API-defined JSON field
}

// QualitygatesApplicationStatus represents the quality gate status of an application.
type QualitygatesApplicationStatus struct {
	// Status is the overall status (OK, WARN, ERROR, NONE).
	Status string `json:"status,omitempty"`
	// Projects is the status of each child project.
	Projects []map[string]any `json:"projects,omitempty"`
	// Metrics is the list of metrics involved in the status.
	Metrics []map[string]any `json:"metrics,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AlmSettingsValidateBindingOptions contains parameters for the ValidateBinding method.
type AlmSettingsValidateBindingOptions struct {
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// CeSetWorkerCountOptions contains parameters for the SetWorkerCount method.
type CeSetWorkerCountOptions struct {
	// Count is the number of workers to use in the Compute Engine. This field is required.
	Count int64 `url:"count"`
}

// IssuesGitlabSastExportOptions contains parameters for the GitlabSastExport method.
type IssuesGitlabSastExportOptions struct {
	// Branch is the branch key. Cannot be combined with PullRequest.
	Branch string `url:"branch,omitempty"`
	// ProjectKey is the project key. This field is required.
	ProjectKey string `url:"projectKey"`
	// PullRequest is the pull request id. Cannot be combined with Branch.
	PullRequest string `url:"pullRequest,omitempty"`
}

// ProjectBadgesAiCodeAssuranceOptions contains parameters for the AiCodeAssurance method.
type ProjectBadgesAiCodeAssuranceOptions struct {
	// Project is the project or application key. This field is required.
	Project string `url:"project"`
	// Token is the project badge token. Required for private projects.
	Token string `url:"token,omitempty"`
}

// ProjectBranchesGetAiCodeAssuranceOptions contains parameters for the GetAiCodeAssurance method.
type ProjectBranchesGetAiCodeAssuranceOptions struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ProjectDumpImportOptions contains parameters for the Import method.
type ProjectDumpImportOptions struct {
	// Key is the key of the (existing, empty) project to import the dump into. This field is required.
	Key string `url:"key"`
}

// ProjectsProjectOptions contains a project key parameter, shared by several Projects methods.
type ProjectsProjectOptions struct {
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// ProjectsExportFindingsOptions contains parameters for the ExportFindings method.
type ProjectsExportFindingsOptions struct {
	// Branch is the branch key.
	Branch string `url:"branch,omitempty"`
	// Project is the project key. This field is required.
	Project string `url:"project"`
	// PullRequest is the pull request key. When omitted, the branch data is returned.
	PullRequest string `url:"pullRequest,omitempty"`
}

// ProjectsSetContainsAiCodeOptions contains parameters for the SetContainsAiCode method.
type ProjectsSetContainsAiCodeOptions struct {
	// Project is the project key. This field is required.
	Project string `url:"project"`
	// ContainsAiCode sets whether the project contains AI code. This field is required.
	ContainsAiCode bool `url:"contains_ai_code"`
}

// QualitygatesApplicationStatusOptions contains parameters for the ApplicationStatus method.
type QualitygatesApplicationStatusOptions struct {
	// Application is the application key. This field is required.
	Application string `url:"application"`
	// Branch is the application branch name.
	Branch string `url:"branch,omitempty"`
}

// QualitygatesSetAiCodeAssuranceOptions contains parameters for the SetAiCodeAssurance method.
type QualitygatesSetAiCodeAssuranceOptions struct {
	// GateName is the quality gate name. This field is required.
	GateName string `url:"gateName"`
	// AiCodeAssurance sets whether the quality gate is AI Code Assured. This field is required.
	AiCodeAssurance bool `url:"aiCodeAssurance"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateValidateBindingOpt validates the options for the ValidateBinding method.
func (s *AlmSettingsService) ValidateValidateBindingOpt(opt *AlmSettingsValidateBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Project, "Project")
}

// ValidateSetWorkerCountOpt validates the options for the SetWorkerCount method.
func (s *CeService) ValidateSetWorkerCountOpt(opt *CeSetWorkerCountOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if opt.Count < 1 {
		return NewValidationError("Count", "must be greater than 0", ErrOutOfRange)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ValidateBinding validates a project binding by checking connectivity and permissions.
// Requires the 'Browse' permission on the project.
//
// API endpoint: GET /api/alm_settings/validate_binding.
// Internal endpoint.
func (s *AlmSettingsService) ValidateBinding(ctx context.Context, opt *AlmSettingsValidateBindingOptions) (*AlmSettingsValidateBinding, *http.Response, error) {
	err := s.ValidateValidateBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV1[AlmSettingsValidateBinding](ctx, s.client, http.MethodGet, "alm_settings/validate_binding", opt)
}

// SetWorkerCount sets the number of workers in the Compute Engine.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/ce/set_worker_count.
// Internal endpoint.
func (s *CeService) SetWorkerCount(ctx context.Context, opt *CeSetWorkerCountOptions) (*http.Response, error) {
	err := s.ValidateSetWorkerCountOpt(opt)
	if err != nil {
		return nil, err
	}

	return postV1NoContent(ctx, s.client, "ce/set_worker_count", opt)
}

// Status returns the GitHub provisioning activation status.
// Requires authentication.
//
// API endpoint: GET /api/github_provisioning/status.
// Internal endpoint.
func (s *GithubProvisioningService) Status(ctx context.Context) (*GithubProvisioningActivation, *http.Response, error) {
	return callV1[GithubProvisioningActivation](ctx, s.client, http.MethodGet, "github_provisioning/status", nil)
}

// Sync schedules a GitHub provisioning task.
// Requires the GitHub provisioning feature to be enabled and the 'Administer System' permission.
//
// API endpoint: POST /api/github_provisioning/sync.
// Internal endpoint.
func (s *GithubProvisioningService) Sync(ctx context.Context) (*http.Response, error) {
	return postV1NoContent(ctx, s.client, "github_provisioning/sync", nil)
}

// GitlabSastExport returns the vulnerabilities of a project in the GitLab SAST JSON format.
// Requires the 'Browse' or 'Scan' permission on the project.
//
// API endpoint: GET /api/issues/gitlab_sast_export.
func (s *IssuesService) GitlabSastExport(ctx context.Context, opt *IssuesGitlabSastExportOptions) (*IssuesGitlabSastExport, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return nil, nil, err
	}

	if opt.Branch != "" && opt.PullRequest != "" {
		return nil, nil, NewValidationError("Branch", "cannot be combined with PullRequest", ErrInvalidValue)
	}

	return callV1[IssuesGitlabSastExport](ctx, s.client, http.MethodGet, "issues/gitlab_sast_export", opt)
}

// AiCodeAssurance generates a badge for the AI code assurance of a project as an SVG.
// Requires the 'Browse' permission on the project.
//
// API endpoint: GET /api/project_badges/ai_code_assurance.
func (s *ProjectBadgesService) AiCodeAssurance(ctx context.Context, opt *ProjectBadgesAiCodeAssuranceOptions) (*string, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return nil, nil, err
	}

	return callV1[string](ctx, s.client, http.MethodGet, "project_badges/ai_code_assurance", opt)
}

// GetAiCodeAssurance returns whether a project is AI Code Assured.
//
// API endpoint: GET /api/project_branches/get_ai_code_assurance.
// Internal endpoint.
func (s *ProjectBranchesService) GetAiCodeAssurance(ctx context.Context, opt *ProjectBranchesGetAiCodeAssuranceOptions) (*ProjectBranchesAiCodeAssurance, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return nil, nil, err
	}

	return callV1[ProjectBranchesAiCodeAssurance](ctx, s.client, http.MethodGet, "project_branches/get_ai_code_assurance", opt)
}

// Import triggers the import of a project dump. Requires the 'Administer' permission.
//
// API endpoint: POST /api/project_dump/import.
func (s *ProjectDumpService) Import(ctx context.Context, opt *ProjectDumpImportOptions) (*ProjectDumpExport, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return nil, nil, err
	}

	return callV1[ProjectDumpExport](ctx, s.client, http.MethodPost, "project_dump/import", opt)
}

// ExportFindings exports all the findings (issues and hotspots) of a project branch.
// Requires the 'Administer System' permission. Can be slow and resource intensive.
//
// API endpoint: GET /api/projects/export_findings.
func (s *ProjectsService) ExportFindings(ctx context.Context, opt *ProjectsExportFindingsOptions) (*ProjectsExportFindings, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return nil, nil, err
	}

	return callV1[ProjectsExportFindings](ctx, s.client, http.MethodGet, "projects/export_findings", opt)
}

// GetContainsAiCode returns whether a project is flagged as containing AI code.
//
// API endpoint: GET /api/projects/get_contains_ai_code.
func (s *ProjectsService) GetContainsAiCode(ctx context.Context, opt *ProjectsProjectOptions) (*ProjectsContainsAiCode, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return nil, nil, err
	}

	return callV1[ProjectsContainsAiCode](ctx, s.client, http.MethodGet, "projects/get_contains_ai_code", opt)
}

// GetDetectedAiCode returns whether AI code was detected in a project.
//
// API endpoint: GET /api/projects/get_detected_ai_code.
// Internal endpoint.
func (s *ProjectsService) GetDetectedAiCode(ctx context.Context, opt *ProjectsProjectOptions) (*ProjectsDetectedAiCode, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return nil, nil, err
	}

	return callV1[ProjectsDetectedAiCode](ctx, s.client, http.MethodGet, "projects/get_detected_ai_code", opt)
}

// SetContainsAiCode sets whether a project contains AI code.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/projects/set_contains_ai_code.
func (s *ProjectsService) SetContainsAiCode(ctx context.Context, opt *ProjectsSetContainsAiCodeOptions) (*http.Response, error) {
	if opt == nil {
		return nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return nil, err
	}

	return postV1NoContent(ctx, s.client, "projects/set_contains_ai_code", opt)
}

// ApplicationStatus returns the quality gate status of an application.
// Requires the 'Browse' permission on the application and its child projects.
//
// API endpoint: GET /api/qualitygates/application_status.
// Internal endpoint.
func (s *QualitygatesService) ApplicationStatus(ctx context.Context, opt *QualitygatesApplicationStatusOptions) (*QualitygatesApplicationStatus, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Application, "Application")
	if err != nil {
		return nil, nil, err
	}

	return callV1[QualitygatesApplicationStatus](ctx, s.client, http.MethodGet, "qualitygates/application_status", opt)
}

// SetAiCodeAssurance qualifies or disqualifies a custom quality gate as AI Code Assured.
// Requires the 'Administer Quality Gates' permission.
//
// API endpoint: POST /api/qualitygates/set_ai_code_assurance.
// Internal endpoint.
func (s *QualitygatesService) SetAiCodeAssurance(ctx context.Context, opt *QualitygatesSetAiCodeAssuranceOptions) (*http.Response, error) {
	if opt == nil {
		return nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GateName, "GateName")
	if err != nil {
		return nil, err
	}

	return postV1NoContent(ctx, s.client, "qualitygates/set_ai_code_assurance", opt)
}
