package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DetectionAgentService handles communication with the detection agent related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type DetectionAgentService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// DetectionAgentInstanceConfiguration represents the InstanceConfiguration object of the SonarQube
// V2 API.
type DetectionAgentInstanceConfiguration struct {
	// Enabled is the enabled.
	Enabled bool `json:"enabled,omitempty"`
	// MaxConcurrentJobs is the max concurrent jobs.
	MaxConcurrentJobs int32 `json:"maxConcurrentJobs,omitempty"`
}

// DetectionAgentListJobsOptions contains parameters for the ListJobs method.
type DetectionAgentListJobsOptions struct {
	// ProjectId restrict to jobs of this project. This field is required.
	ProjectId string `json:"projectId"`
	// BranchId restrict to jobs of this branch.
	BranchId string `json:"branchId,omitempty"`
	// Status job statuses to include. Repeat the parameter to provide multiple statuses. Allowed
	// values: QUEUED, IN_PROGRESS, COMPLETED, FAILED, SKIPPED.
	Status []string `json:"status,omitempty"`
	// PageIndex 1-based page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize number of items per page.
	// Must be between 1 and 500. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// DetectionAgentJob represents the Job object of the SonarQube V2 API.
type DetectionAgentJob struct {
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

// DetectionAgentJobsResponse represents the JobsResponse object of the SonarQube V2 API.
type DetectionAgentJobsResponse struct {
	// Jobs is the jobs.
	Jobs []DetectionAgentJob `json:"jobs,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// DetectionAgentGetCreditCostOptions contains parameters for the GetCreditCost method.
type DetectionAgentGetCreditCostOptions struct {
	// BranchId stable identifier of the branch to estimate. This field is required.
	BranchId string `json:"branchId"`
}

// DetectionAgentCreditCostBranch represents the CreditCostBranch object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentCreditCostBranch struct {
	// Id stable identifier of the branch.
	Id string `json:"id,omitempty"`
	// Ncloc latest persisted non-comment lines of code for the branch.
	Ncloc int64 `json:"ncloc,omitempty"`
	// Size is the size. Allowed values: XS, S, M, L, XL, XXL.
	Size string `json:"size,omitempty"`
}

// DetectionAgentCreditCostBreakdown represents the CreditCostBreakdown object of the SonarQube V2
// API.
type DetectionAgentCreditCostBreakdown struct {
	// Total credits required by the branch.
	Total int64 `json:"total,omitempty"`
	// Base branch cost assigned to available base credits.
	Base int64 `json:"base,omitempty"`
	// Overage credits needed beyond the available base quota, including when the request will be
	// refused.
	Overage int64 `json:"overage,omitempty"`
}

// DetectionAgentAvailableCreditsBreakdown represents the AvailableCreditsBreakdown object of the
// SonarQube V2 API.
type DetectionAgentAvailableCreditsBreakdown struct {
	// Total base and active overage credits available after usage and reservations.
	Total int64 `json:"total,omitempty"`
	// Base credits available after usage and reservations.
	Base int64 `json:"base,omitempty"`
	// Overage active overage credits available after usage and reservations. Null when overage is
	// disabled.
	Overage int64 `json:"overage,omitempty"`
}

// DetectionAgentCreditCostCredits represents the CreditCostCredits object of the SonarQube V2 API.
type DetectionAgentCreditCostCredits struct {
	// BranchCost is the branch cost.
	BranchCost DetectionAgentCreditCostBreakdown `json:"branchCost,omitzero"`
	// Available is the available.
	Available DetectionAgentAvailableCreditsBreakdown `json:"available,omitzero"`
}

// DetectionAgentCreditCostResponse represents the CreditCostResponse object of the SonarQube V2
// API.
type DetectionAgentCreditCostResponse struct {
	// Branch is the branch.
	Branch DetectionAgentCreditCostBranch `json:"branch,omitzero"`
	// Credits is the credits.
	Credits DetectionAgentCreditCostCredits `json:"credits,omitzero"`
	// CreditCost estimated credit cost for the branch size tier.
	CreditCost int32 `json:"creditCost,omitempty"`
	// AvailableCredits available credits, when consumption tracking is configured.
	AvailableCredits int64 `json:"availableCredits,omitempty"`
}

// DetectionAgentScheduling represents the scheduled-run configuration of a project as returned
// by the detection-agent project configuration endpoints.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentScheduling struct {
	// BranchId is the identifier of the long-lived branch scanned on scheduled runs.
	BranchId string `json:"branchId,omitempty"`
	// Enabled indicates whether scheduled runs are effectively enabled. It reads as false while
	// the detection agent is disabled instance-wide.
	Enabled bool `json:"enabled,omitempty"`
	// NextRunTime is the server-computed date and time of the next scheduled run.
	NextRunTime string `json:"nextRunTime,omitempty"`
	// Recurrence is how often the run recurs.
	Recurrence *DetectionAgentScheduleRecurrence `json:"recurrence,omitempty"`
	// Time is the time of day of the scheduled run.
	Time *DetectionAgentScheduleTime `json:"time,omitempty"`
}

// DetectionAgentProjectConfiguration represents the detection-agent configuration of a project.
//
// Note: the SonarQube V2 API spec references a colliding "ProjectConfiguration" schema (the
// architecture scanner configuration); this type reflects what the endpoint actually returns.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentProjectConfiguration struct {
	// ProjectId is the project identifier.
	ProjectId string `json:"projectId,omitempty"`
	// Scheduling is the scheduled-run configuration of the project.
	Scheduling DetectionAgentScheduling `json:"scheduling,omitzero"`
}

// DetectionAgentScheduleTime represents a SonarQube V2 API object.
//
// Time of day of the scheduled run, in the given IANA timezone. Runs fire on the hour.
type DetectionAgentScheduleTime struct {
	// Hour of day (0-23) in the given timezone.
	Hour *int32 `json:"hour,omitempty"`
	// Timezone IANA timezone id the schedule fires in (DST-aware).
	Timezone string `json:"timezone,omitempty"`
}

// DetectionAgentScheduleRecurrence represents a SonarQube V2 API object.
//
// How often the run recurs. WEEKLY requires exactly one dayOfWeek; DAILY and MONTHLY must omit it.
// MONTHLY always fires on the last day of the month.
type DetectionAgentScheduleRecurrence struct {
	// Type the recurrence cadence. Allowed values: DAILY, WEEKLY, MONTHLY.
	Type string `json:"type,omitempty"`
	// DayOfWeek day of week for WEEKLY recurrence; omitted otherwise. Allowed values: MONDAY,
	// TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, SUNDAY.
	DayOfWeek string `json:"dayOfWeek,omitempty"`
}

// DetectionAgentSchedulingRequest represents a SonarQube V2 API object.
//
// Scheduled-run configuration to apply to the project. The structured config (branch, time,
// recurrence) is optional but all-or-nothing: supply a complete config to overwrite the stored
// schedule, or omit it to keep the stored schedule. Omitting it while enabling a project that has
// no stored schedule provisions a default one (daily at 02:00 Etc/UTC on the main branch) instead
// of failing, so enabling never requires the client to pick a schedule. This is a request-only
// model: the server-derived nextRunTime is never accepted here.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentSchedulingRequest struct {
	// Enabled whether scheduled detection-agent runs are desired for the project. This is the
	// project's own setting, not the effective state: unlike the identically named response field
	// it is never combined with the instance-wide detection-agent setting. Setting it to true
	// while the agent is disabled instance-wide is refused with 403; setting it to false is always
	// accepted. Do not populate it by echoing a SchedulingResponse read back — while the agent is
	// disabled instance-wide that response reports false as a mask, and echoing it turns the
	// project's own scheduling off for real.
	Enabled *bool `json:"enabled,omitempty"`
	// BranchId stable identifier of the long-lived branch to scan on scheduled runs.
	BranchId string `json:"branchId,omitempty"`
	// Time is the time of day of the scheduled run.
	Time *DetectionAgentScheduleTime `json:"time,omitempty"`
	// Recurrence is how often the run recurs.
	Recurrence *DetectionAgentScheduleRecurrence `json:"recurrence,omitempty"`
}

// DetectionAgentUpsertProjectConfigurationOptions contains the request body for the
// UpsertProjectConfiguration method.
type DetectionAgentUpsertProjectConfigurationOptions struct {
	// Scheduling is the scheduling.
	Scheduling *DetectionAgentSchedulingRequest `json:"scheduling,omitempty"`
}

// DetectionAgentCreateJobOptions contains the request body for the CreateJob method.
type DetectionAgentCreateJobOptions struct {
	// ProjectId stable identifier of the SonarQube project to analyse. This field is required.
	ProjectId string `json:"projectId"`
	// BranchId stable identifier of the branch to analyse. This field is required.
	BranchId string `json:"branchId"`
}

// DetectionAgentAgenticFlowLocation represents the AgenticFlowLocation object of the SonarQube V2
// API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentAgenticFlowLocation struct {
	// FilePath is the file path.
	FilePath string `json:"filePath,omitempty"`
	// Line is the line.
	Line int32 `json:"line,omitempty"`
	// Message is the message.
	Message string `json:"message,omitempty"`
}

// DetectionAgentAgenticFlow represents the AgenticFlow object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentAgenticFlow struct {
	// Locations is the locations.
	Locations []DetectionAgentAgenticFlowLocation `json:"locations,omitempty"`
	// Description is the description.
	Description string `json:"description,omitempty"`
	// Type is the type.
	Type string `json:"type,omitempty"`
}

// DetectionAgentAgenticFinding represents the AgenticFinding object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentAgenticFinding struct {
	// FilePath is the file path.
	FilePath string `json:"filePath,omitempty"`
	// Line is the line.
	Line int32 `json:"line,omitempty"`
	// Message is the message.
	Message string `json:"message,omitempty"`
	// RuleKey is the rule key.
	RuleKey string `json:"ruleKey,omitempty"`
	// RuleTitle is the rule title.
	RuleTitle string `json:"ruleTitle,omitempty"`
	// Severity is the severity.
	Severity string `json:"severity,omitempty"`
	// WhyIsThisAnIssue is the why is this an issue.
	WhyIsThisAnIssue string `json:"whyIsThisAnIssue,omitempty"`
	// HowToFix is the how to fix.
	HowToFix string `json:"howToFix,omitempty"`
	// AgentContext is the agent context.
	AgentContext string `json:"agentContext,omitempty"`
	// Flows is the flows.
	Flows []DetectionAgentAgenticFlow `json:"flows,omitempty"`
	// Status is the status.
	Status string `json:"status,omitempty"`
}

// DetectionAgentSubmitJobResultOptions contains the request body for the SubmitJobResult method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DetectionAgentSubmitJobResultOptions struct {
	// JobId is the job id.
	JobId string `json:"jobId,omitempty"`
	// ProjectId is the project id.
	ProjectId string `json:"projectId,omitempty"`
	// BranchName is the branch name.
	BranchName string `json:"branchName,omitempty"`
	// AnalyzedCommitSha is the analyzed commit sha.
	AnalyzedCommitSha string `json:"analyzedCommitSha,omitempty"`
	// AnalyzedAt is the analyzed at.
	AnalyzedAt *int64 `json:"analyzedAt,omitempty"`
	// Findings is the findings.
	Findings []DetectionAgentAgenticFinding `json:"findings,omitempty"`
	// DurationSeconds is the duration seconds.
	DurationSeconds *int64 `json:"durationSeconds,omitempty"`
	// CostUsd is the cost usd.
	CostUsd *float64 `json:"costUsd,omitempty"`
	// PlaybookId is the playbook id.
	PlaybookId string `json:"playbookId,omitempty"`
	// PlaybookKey is the playbook key.
	PlaybookKey string `json:"playbookKey,omitempty"`
	// PlaybookVersion is the playbook version.
	PlaybookVersion string `json:"playbookVersion,omitempty"`
	// RevalidatesPriors is the revalidates priors.
	RevalidatesPriors *bool `json:"revalidatesPriors,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListJobsOpt validates the options for the ListJobs method.
func (s *DetectionAgentService) ValidateListJobsOpt(opt *DetectionAgentListJobsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectId, "ProjectId")
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

// ValidateGetCreditCostOpt validates the options for the GetCreditCost method.
func (s *DetectionAgentService) ValidateGetCreditCostOpt(opt *DetectionAgentGetCreditCostOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.BranchId, "BranchId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateJobOpt validates the options for the CreateJob method.
func (s *DetectionAgentService) ValidateCreateJobOpt(opt *DetectionAgentCreateJobOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.BranchId, "BranchId")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetInstanceConfiguration gets instance configuration. Retrieve the instance-wide detection-agent
// configuration.
//
// API endpoint: GET /api/v2/detection-agent/instance-config.
// Enterprise Edition only.
func (s *DetectionAgentService) GetInstanceConfiguration(ctx context.Context) (*DetectionAgentInstanceConfiguration, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "detection-agent/instance-config", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DetectionAgentInstanceConfiguration)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListJobs lists detection jobs. Returns a paginated list of detection jobs for a project.
//
// API endpoint: GET /api/v2/detection-agent/jobs.
// Enterprise Edition only.
func (s *DetectionAgentService) ListJobs(ctx context.Context, opt *DetectionAgentListJobsOptions) (*DetectionAgentJobsResponse, *http.Response, error) {
	err := s.ValidateListJobsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "detection-agent/jobs", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DetectionAgentJobsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetCreditCost estimates detection job credit cost. Returns the LOC-based credit cost for the
// latest persisted NCLOC of a branch.
//
// API endpoint: GET /api/v2/detection-agent/jobs/credit-cost.
// Enterprise Edition only.
func (s *DetectionAgentService) GetCreditCost(ctx context.Context, opt *DetectionAgentGetCreditCostOptions) (*DetectionAgentCreditCostResponse, *http.Response, error) {
	err := s.ValidateGetCreditCostOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "detection-agent/jobs/credit-cost", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DetectionAgentCreditCostResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectConfiguration gets project configuration. Retrieve the detection-agent configuration
// for a project. Scheduling reflects the instance-wide detection-agent setting, so it reads as
// disabled while the detection agent is disabled instance-wide. Entitlement and credit availability
// are not reflected here: a scheduled run also requires the instance to be entitled and to hold
// enough credits at fire time, so scheduling can read as enabled on an instance where a run would
// be refused for one of those reasons.
//
// API endpoint: GET /api/v2/detection-agent/project-configs/{projectId}.
// Enterprise Edition only.
func (s *DetectionAgentService) GetProjectConfiguration(ctx context.Context, projectID string) (*DetectionAgentProjectConfiguration, *http.Response, error) {
	err := ValidateRequired(projectID, "projectID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "detection-agent/project-configs/"+url.PathEscape(projectID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DetectionAgentProjectConfiguration)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpsertProjectConfiguration upserts project configuration. Create or update the detection-agent
// configuration (including scheduling) for a project. Enabling scheduling with no schedule
// available — neither supplied in the request nor already stored — provisions a default one rather
// than failing: a daily run at 02:00 Etc/UTC on the project's main branch, persisted and returned
// in the response. Each missing piece is defaulted independently; a partial structured config is
// still refused with 400. While the detection agent is disabled instance-wide, a request that sets
// scheduling.enabled to true is refused with 403: nothing can dispatch, so scheduling cannot be
// switched on. Disabling is always accepted, and a request that carries no scheduling intent writes
// nothing and behaves as a read. The instance-wide setting never touches a project's stored
// schedule, which takes effect again once that setting is turned back on.
//
// API endpoint: PATCH /api/v2/detection-agent/project-configs/{projectId}.
// Enterprise Edition only.
func (s *DetectionAgentService) UpsertProjectConfiguration(ctx context.Context, projectID string, opt *DetectionAgentUpsertProjectConfigurationOptions) (*DetectionAgentProjectConfiguration, *http.Response, error) {
	err := ValidateRequired(projectID, "projectID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Unlike most V2 PATCH endpoints, this one only accepts application/json (it answers 415 to
	// application/merge-patch+json), so the default PATCH content type is overridden.
	//nolint:exhaustruct // only the fields relevant to this request are set
	req, err := s.client.NewSonarQubeAPIRequest(ctx, SonarAPIRequestParameters{
		Method:  http.MethodPatch,
		Path:    v2BasePath + "detection-agent/project-configs/" + url.PathEscape(projectID),
		Headers: map[string]string{headerContentType: "application/json"},
		Body:    opt,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DetectionAgentProjectConfiguration)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateJob creates a detection job. Triggers a new agentic detection job for a project.
//
// API endpoint: POST /api/v2/detection-agent/jobs.
// Enterprise Edition only.
func (s *DetectionAgentService) CreateJob(ctx context.Context, opt *DetectionAgentCreateJobOptions) (*DetectionAgentJob, *http.Response, error) {
	err := s.ValidateCreateJobOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "detection-agent/jobs", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DetectionAgentJob)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SubmitJobResult calls POST /api/v2/detection-agent/jobs/results.
//
// API endpoint: POST /api/v2/detection-agent/jobs/results.
// Enterprise Edition only.
func (s *DetectionAgentService) SubmitJobResult(ctx context.Context, jobID string, opt *DetectionAgentSubmitJobResultOptions) (*http.Response, error) {
	err := ValidateRequired(jobID, "jobID")
	if err != nil {
		return nil, err
	}

	if opt == nil {
		return nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "detection-agent/jobs/results", map[string]string{"jobId": jobID}, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
