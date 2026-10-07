package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// RemediationAgentService handles communication with the remediation agent related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type RemediationAgentService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedRemediationAgentListJobsStatus is the set of allowed values for the corresponding option.
	allowedRemediationAgentListJobsStatus = map[string]struct{}{
		"PENDING":     {},
		"IN_PROGRESS": {},
		"COMPLETED":   {},
		"FAILED":      {},
	}
	// allowedRemediationAgentCreateBacklogJobAgentTask is the set of allowed values for the corresponding option.
	allowedRemediationAgentCreateBacklogJobAgentTask = map[string]struct{}{
		"SONAR_REPORT": {},
		"HUNTER_FIX":   {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// RemediationAgentJob represents the Job object of the SonarQube V2 API.
type RemediationAgentJob struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// ProjectId is the project id.
	ProjectId string `json:"projectId,omitempty"`
	// BranchId stable identifier of the branch the job ran on.
	BranchId string `json:"branchId,omitempty"`
	// Status is the status. Allowed values: QUEUED, IN_PROGRESS, COMPLETED, FAILED, SKIPPED.
	Status string `json:"status,omitempty"`
	// SubStatus present only when a finer-grained phase is known; status remains the source of
	// truth. Allowed values: AWAITING_FINDINGS, CLONING, EXPLORING, ANALYZING, CONSOLIDATING,
	// REPORTING, CLONE_FAILED, AGENT_FAILED, LLM_NOT_ACCESSIBLE, ANALYSIS_TIMEOUT,
	// RUNTIME_UNAVAILABLE, RESULT_PUBLICATION_FAILED, INSUFFICIENT_CREDITS, NCLOC_UNAVAILABLE,
	// UNSUPPORTED_SCM_PROVIDER, SCM_CREDENTIALS_REJECTED, SCM_BRANCH_NOT_FOUND,
	// PRIOR_FINDINGS_UNUSABLE.
	SubStatus string `json:"subStatus,omitempty"`
	// AnalysisType is the analysis type. Allowed values: FULL, INCREMENTAL.
	AnalysisType string `json:"analysisType,omitempty"`
	// TriggerType is the trigger type. Allowed values: ON_DEMAND, SCHEDULED.
	TriggerType string `json:"triggerType,omitempty"`
	// StartTime epoch millis when the job started; null while PENDING.
	StartTime int64 `json:"startTime,omitempty"`
	// DurationSeconds is the duration seconds.
	DurationSeconds int64 `json:"durationSeconds,omitempty"`
	// ConsumedCredits credits submitted to Billing for this run. Zero when the terminal run was
	// not chargeable; null while non-terminal, for legacy runs, or when billing publication
	// failed.
	ConsumedCredits int32 `json:"consumedCredits,omitempty"`
	// RequestTime is the request time.
	RequestTime int64 `json:"requestTime,omitempty"`
	// UpdatedAt is the updated at.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
}

// RemediationAgentListJobsOptions contains parameters for the ListJobs method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type RemediationAgentListJobsOptions struct {
	// ProjectKey restrict to jobs of this project. This field is required.
	ProjectKey string `json:"projectKey"`
	// Status filter by status. Omit to return all statuses. Allowed values: PENDING, IN_PROGRESS,
	// COMPLETED, FAILED.
	Status string `json:"status,omitempty"`
	// SourcePullRequestKey filter by the key of the original pull request being remediated.
	SourcePullRequestKey string `json:"sourcePullRequestKey,omitempty"`
	// PageIndex 1-based page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize page size. Set to 0 to retrieve only the total count without any items.
	// Must be between 0 and 500. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
	// Sort by a resource field. Prefix '+' for ascending, '-' for descending. Default:
	// '-createdAt' (newest first). Allowed fields: 'createdAt', 'startedAt', 'finishedAt',
	// 'status'.
	// Default: -createdAt.
	Sort string `json:"sort,omitempty"`
}

// RemediationAgentJobSearchResponse represents the JobSearchResponse object of the SonarQube V2
// API.
type RemediationAgentJobSearchResponse struct {
	// Jobs is the jobs.
	Jobs []RemediationAgentJob `json:"jobs,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// RemediationAgentScheduledAgentConfigResource represents the ScheduledAgentConfigResource object
// of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type RemediationAgentScheduledAgentConfigResource struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Enabled is the enabled.
	Enabled bool `json:"enabled,omitempty"`
	// Frequency is the frequency.
	Frequency string `json:"frequency,omitempty"`
	// Hour is the hour.
	Hour int32 `json:"hour,omitempty"`
	// DayOfWeek is the day of week.
	DayOfWeek string `json:"dayOfWeek,omitempty"`
	// Timezone is the timezone.
	Timezone string `json:"timezone,omitempty"`
	// NextRunUtc is the next run utc.
	NextRunUtc int64 `json:"nextRunUtc,omitempty"`
	// MaxIssuesPerRun is the max issues per run.
	MaxIssuesPerRun int32 `json:"maxIssuesPerRun,omitempty"`
	// MaxOpenPRs is the max open p rs.
	MaxOpenPRs int32 `json:"maxOpenPRs,omitempty"`
	// ProjectSelectionMode is the project selection mode.
	ProjectSelectionMode string `json:"projectSelectionMode,omitempty"`
	// SelectedProjectKeys is the selected project keys.
	SelectedProjectKeys []string `json:"selectedProjectKeys,omitempty"`
}

// RemediationAgentScheduledAgentConfigSearchRestResponse represents the
// ScheduledAgentConfigSearchRestResponse object of the SonarQube V2 API.
type RemediationAgentScheduledAgentConfigSearchRestResponse struct {
	// ScheduledAgentConfigs is the scheduled agent configs.
	ScheduledAgentConfigs []RemediationAgentScheduledAgentConfigResource `json:"scheduledAgentConfigs,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// RemediationAgentScheduledAgentDefaultResource represents the ScheduledAgentDefaultResource object
// of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type RemediationAgentScheduledAgentDefaultResource struct {
	// Enabled is the enabled.
	Enabled bool `json:"enabled,omitempty"`
	// Frequency is the frequency.
	Frequency string `json:"frequency,omitempty"`
	// Hour is the hour.
	Hour int32 `json:"hour,omitempty"`
	// DayOfWeek is the day of week.
	DayOfWeek string `json:"dayOfWeek,omitempty"`
	// Timezone is the timezone.
	Timezone string `json:"timezone,omitempty"`
	// MaxIssuesPerRun is the max issues per run.
	MaxIssuesPerRun int32 `json:"maxIssuesPerRun,omitempty"`
	// MaxOpenPRs is the max open p rs.
	MaxOpenPRs int32 `json:"maxOpenPRs,omitempty"`
}

// RemediationAgentScheduledAgentConfigEffectiveRestResponse represents the
// ScheduledAgentConfigEffectiveRestResponse object of the SonarQube V2 API.
type RemediationAgentScheduledAgentConfigEffectiveRestResponse struct {
	// ProjectConfig is the project config.
	ProjectConfig RemediationAgentScheduledAgentConfigResource `json:"projectConfig,omitzero"`
	// OrganizationDefault is the organization default.
	OrganizationDefault RemediationAgentScheduledAgentDefaultResource `json:"organizationDefault,omitzero"`
}

// RemediationAgentCreateBacklogJobOptions contains the request body for the CreateBacklogJob
// method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type RemediationAgentCreateBacklogJobOptions struct {
	// ProjectKey key of the SonarQube project to remediate. This field is required.
	ProjectKey string `json:"projectKey"`
	// RepositoryUrl clone URL of the repository to remediate. When omitted, it is resolved from
	// the project's configuration.
	RepositoryUrl string `json:"repositoryUrl,omitempty"`
	// Branch to remediate. Defaults to the repository's default branch. Required (must be non-
	// blank) when agentTask is HUNTER_FIX, since the finding's branch is not otherwise known.
	Branch string `json:"branch,omitempty"`
	// IssueKeys sonarQube issue keys selected by the caller for remediation. Every key must belong
	// to projectKey, or the request is rejected. Forwarded to the orchestrator, which enriches
	// them into full issue metadata at dispatch time. This field is required.
	IssueKeys []string `json:"issueKeys"`
	// AgentTask which agent subcommand to run for this job. SONAR_REPORT (the default, and the
	// value assumed when the field is absent) is the existing multi-issue Sonar-issue flow.
	// HUNTER_FIX runs the hunter-fix subcommand against a single Hunter Agent finding and requires
	// exactly one entry in issueKeys and a non-blank branch. Allowed values: SONAR_REPORT,
	// HUNTER_FIX.
	AgentTask string `json:"agentTask,omitempty"`
}

// RemediationAgentCreatePullRequestJobOptions contains the request body for the
// CreatePullRequestJob method.
type RemediationAgentCreatePullRequestJobOptions struct {
	// ProjectKey key of the SonarQube project to remediate. This field is required.
	ProjectKey string `json:"projectKey"`
	// PullRequestKey key of the pull request to scope the remediation to. This field is required.
	PullRequestKey string `json:"pullRequestKey"`
	// RepositoryUrl clone URL of the repository to remediate. When omitted, it is resolved from
	// the project's configuration.
	RepositoryUrl string `json:"repositoryUrl,omitempty"`
	// IssueKeys sonarQube issue keys (scoped to the pull request) selected by the caller for
	// remediation. Every key must belong to projectKey, or the request is rejected. Forwarded to
	// the orchestrator, which enriches them into full issue metadata at dispatch time. This field
	// is required.
	IssueKeys []string `json:"issueKeys"`
}

// RemediationAgentCreateScheduledAgentConfigOptions contains the request body for the
// CreateScheduledAgentConfig method.
type RemediationAgentCreateScheduledAgentConfigOptions struct {
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// Enabled is the enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Frequency is the frequency.
	Frequency string `json:"frequency,omitempty"`
	// Hour is the hour.
	Hour *int32 `json:"hour,omitempty"`
	// DayOfWeek is the day of week.
	DayOfWeek string `json:"dayOfWeek,omitempty"`
	// Timezone is the timezone.
	Timezone string `json:"timezone,omitempty"`
	// MaxIssuesPerRun is the max issues per run.
	MaxIssuesPerRun *int32 `json:"maxIssuesPerRun,omitempty"`
	// MaxOpenPRs is the max open p rs.
	MaxOpenPRs *int32 `json:"maxOpenPRs,omitempty"`
	// ProjectSelectionMode is the project selection mode.
	ProjectSelectionMode string `json:"projectSelectionMode,omitempty"`
	// SelectedProjectKeys is the selected project keys.
	SelectedProjectKeys []string `json:"selectedProjectKeys,omitempty"`
}

// RemediationAgentUpdateScheduledAgentConfigOptions contains the request body for the
// UpdateScheduledAgentConfig method.
type RemediationAgentUpdateScheduledAgentConfigOptions struct {
	// Enabled is the enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Frequency is the frequency.
	Frequency string `json:"frequency,omitempty"`
	// Hour is the hour.
	Hour *int32 `json:"hour,omitempty"`
	// DayOfWeek is the day of week.
	DayOfWeek string `json:"dayOfWeek,omitempty"`
	// Timezone is the timezone.
	Timezone string `json:"timezone,omitempty"`
	// MaxIssuesPerRun is the max issues per run.
	MaxIssuesPerRun *int32 `json:"maxIssuesPerRun,omitempty"`
	// MaxOpenPRs is the max open p rs.
	MaxOpenPRs *int32 `json:"maxOpenPRs,omitempty"`
	// ProjectSelectionMode is the project selection mode.
	ProjectSelectionMode string `json:"projectSelectionMode,omitempty"`
	// SelectedProjectKeys is the selected project keys.
	SelectedProjectKeys []string `json:"selectedProjectKeys,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListJobsOpt validates the options for the ListJobs method.
func (s *RemediationAgentService) ValidateListJobsOpt(opt *RemediationAgentListJobsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Status, allowedRemediationAgentListJobsStatus, "Status")
	if err != nil {
		return err
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	if opt.PageSize < 0 || opt.PageSize > 500 {
		return NewValidationError("PageSize", "must be between 0 and 500", ErrOutOfRange)
	}

	return nil
}

// ValidateCreateBacklogJobOpt validates the options for the CreateBacklogJob method.
func (s *RemediationAgentService) ValidateCreateBacklogJobOpt(opt *RemediationAgentCreateBacklogJobOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	if len(opt.IssueKeys) == 0 {
		return NewValidationError("IssueKeys", "is required", ErrMissingRequired)
	}

	err = IsValueAuthorized(opt.AgentTask, allowedRemediationAgentCreateBacklogJobAgentTask, "AgentTask")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreatePullRequestJobOpt validates the options for the CreatePullRequestJob method.
func (s *RemediationAgentService) ValidateCreatePullRequestJobOpt(opt *RemediationAgentCreatePullRequestJobOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.PullRequestKey, "PullRequestKey")
	if err != nil {
		return err
	}

	if len(opt.IssueKeys) == 0 {
		return NewValidationError("IssueKeys", "is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// DeleteScheduledAgentConfig calls DELETE /api/v2/remediation-agent/scheduled-agent-configs/{id}.
//
// API endpoint: DELETE /api/v2/remediation-agent/scheduled-agent-configs/{id}.
// Enterprise Edition only.
func (s *RemediationAgentService) DeleteScheduledAgentConfig(ctx context.Context, scheduledAgentConfigID string) (*http.Response, error) {
	err := ValidateRequired(scheduledAgentConfigID, "scheduledAgentConfigID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "remediation-agent/scheduled-agent-configs/"+url.PathEscape(scheduledAgentConfigID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetBacklogJob gets a remediation backlog job. Retrieve a remediation backlog job, including its
// current status.
//
// API endpoint: GET /api/v2/remediation-agent/backlog-jobs/{jobId}.
// Enterprise Edition only.
func (s *RemediationAgentService) GetBacklogJob(ctx context.Context, jobID string) (*RemediationAgentJob, *http.Response, error) {
	err := ValidateRequired(jobID, "jobID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "remediation-agent/backlog-jobs/"+url.PathEscape(jobID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentJob)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListJobs lists remediation jobs. Returns a paginated list of remediation jobs for a project.
// Default sort is '-createdAt' (newest first).
//
// API endpoint: GET /api/v2/remediation-agent/jobs.
// Enterprise Edition only.
func (s *RemediationAgentService) ListJobs(ctx context.Context, opt *RemediationAgentListJobsOptions) (*RemediationAgentJobSearchResponse, *http.Response, error) {
	err := s.ValidateListJobsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "remediation-agent/jobs", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentJobSearchResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPullRequestJob gets a remediation pull-request job. Retrieve a remediation pull-request job,
// including its current status.
//
// API endpoint: GET /api/v2/remediation-agent/pull-request-jobs/{jobId}.
// Enterprise Edition only.
func (s *RemediationAgentService) GetPullRequestJob(ctx context.Context, jobID string) (*RemediationAgentJob, *http.Response, error) {
	err := ValidateRequired(jobID, "jobID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "remediation-agent/pull-request-jobs/"+url.PathEscape(jobID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentJob)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchScheduledAgentConfigs calls GET /api/v2/remediation-agent/scheduled-agent-configs.
//
// API endpoint: GET /api/v2/remediation-agent/scheduled-agent-configs.
// Enterprise Edition only.
func (s *RemediationAgentService) SearchScheduledAgentConfigs(ctx context.Context) (*RemediationAgentScheduledAgentConfigSearchRestResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "remediation-agent/scheduled-agent-configs", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentScheduledAgentConfigSearchRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetScheduledAgentConfig calls GET /api/v2/remediation-agent/scheduled-agent-configs/{id}.
//
// API endpoint: GET /api/v2/remediation-agent/scheduled-agent-configs/{id}.
// Enterprise Edition only.
func (s *RemediationAgentService) GetScheduledAgentConfig(ctx context.Context, scheduledAgentConfigID string) (*RemediationAgentScheduledAgentConfigResource, *http.Response, error) {
	err := ValidateRequired(scheduledAgentConfigID, "scheduledAgentConfigID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "remediation-agent/scheduled-agent-configs/"+url.PathEscape(scheduledAgentConfigID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentScheduledAgentConfigResource)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetEffectiveScheduledAgentConfig calls GET /api/v2/remediation-agent/scheduled-agent-
// configs/{projectKey}/effective.
//
// API endpoint: GET /api/v2/remediation-agent/scheduled-agent-configs/{projectKey}/effective.
// Enterprise Edition only.
func (s *RemediationAgentService) GetEffectiveScheduledAgentConfig(ctx context.Context, projectKey string) (*RemediationAgentScheduledAgentConfigEffectiveRestResponse, *http.Response, error) {
	err := ValidateRequired(projectKey, "projectKey")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "remediation-agent/scheduled-agent-configs/"+url.PathEscape(projectKey)+"/effective", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentScheduledAgentConfigEffectiveRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateBacklogJob creates a remediation backlog job. Triggers a new full-repository remediation
// job for a project.
//
// API endpoint: POST /api/v2/remediation-agent/backlog-jobs.
// Enterprise Edition only.
func (s *RemediationAgentService) CreateBacklogJob(ctx context.Context, opt *RemediationAgentCreateBacklogJobOptions) (*RemediationAgentJob, *http.Response, error) {
	err := s.ValidateCreateBacklogJobOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "remediation-agent/backlog-jobs", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentJob)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreatePullRequestJob creates a remediation pull-request job. Triggers a new pull-request-scoped
// remediation job for a project.
//
// API endpoint: POST /api/v2/remediation-agent/pull-request-jobs.
// Enterprise Edition only.
func (s *RemediationAgentService) CreatePullRequestJob(ctx context.Context, opt *RemediationAgentCreatePullRequestJobOptions) (*RemediationAgentJob, *http.Response, error) {
	err := s.ValidateCreatePullRequestJobOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "remediation-agent/pull-request-jobs", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentJob)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateScheduledAgentConfig calls POST /api/v2/remediation-agent/scheduled-agent-configs.
//
// API endpoint: POST /api/v2/remediation-agent/scheduled-agent-configs.
// Enterprise Edition only.
func (s *RemediationAgentService) CreateScheduledAgentConfig(ctx context.Context, opt *RemediationAgentCreateScheduledAgentConfigOptions) (*RemediationAgentScheduledAgentConfigResource, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "remediation-agent/scheduled-agent-configs", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentScheduledAgentConfigResource)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateScheduledAgentConfig calls PUT /api/v2/remediation-agent/scheduled-agent-configs/{id}.
//
// API endpoint: PUT /api/v2/remediation-agent/scheduled-agent-configs/{id}.
// Enterprise Edition only.
func (s *RemediationAgentService) UpdateScheduledAgentConfig(ctx context.Context, scheduledAgentConfigID string, opt *RemediationAgentUpdateScheduledAgentConfigOptions) (*RemediationAgentScheduledAgentConfigResource, *http.Response, error) {
	err := ValidateRequired(scheduledAgentConfigID, "scheduledAgentConfigID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPut, "remediation-agent/scheduled-agent-configs/"+url.PathEscape(scheduledAgentConfigID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(RemediationAgentScheduledAgentConfigResource)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
