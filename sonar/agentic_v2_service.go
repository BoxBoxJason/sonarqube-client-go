package sonar

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// AgenticService handles communication with the Agentic orchestration jobs (Hunter and Remediation agents) related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type AgenticService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedAgenticSearchJobsStatus is the set of allowed values for the corresponding option.
	allowedAgenticSearchJobsStatus = map[string]struct{}{
		"PENDING":     {},
		"IN_PROGRESS": {},
		"COMPLETED":   {},
		"FAILED":      {},
	}
	// allowedAgenticSearchJobsType is the set of allowed values for the corresponding option.
	allowedAgenticSearchJobsType = map[string]struct{}{
		"HUNTER":      {},
		"REMEDIATION": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// AgenticSearchJobsOptions contains parameters for the SearchJobs method.
type AgenticSearchJobsOptions struct {
	// Id comma-separated list of job ids to filter on.
	Id string `json:"id,omitempty"`
	// Status comma-separated list of statuses to filter on. Allowed values: PENDING, IN_PROGRESS,
	// COMPLETED, FAILED.
	Status string `json:"status,omitempty"`
	// Type comma-separated list of agent types to filter on. Allowed values: HUNTER, REMEDIATION.
	Type string `json:"type,omitempty"`
	// CreatedAfter only return jobs created on or after this date (or date-time), e.g. 2026-01-01
	// or 2026-01-01T12:00:00+0000.
	CreatedAfter string `json:"createdAfter,omitempty"`
	// CreatedBefore only return jobs created on or before this date (or date-time), e.g.
	// 2026-01-01 or 2026-01-01T12:00:00+0000.
	CreatedBefore string `json:"createdBefore,omitempty"`
	// PageSize number of results per page. A value of 0 will only return the pagination
	// information.
	// Must be between 0 and 500. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
	// PageIndex 1-based page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
}

// AgenticAgenticJobRestResponse represents the AgenticJobRestResponse object of the SonarQube V2
// API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type AgenticAgenticJobRestResponse struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// ProjectId is the project id.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectKey not set if the project has been deleted.
	ProjectKey string `json:"projectKey,omitempty"`
	// ProjectName not set if the project has been deleted.
	ProjectName string `json:"projectName,omitempty"`
	// Type is the type. Allowed values: HUNTER, REMEDIATION.
	Type string `json:"type,omitempty"`
	// AnalysisType is the analysis type. Allowed values: FULL, INCREMENTAL.
	AnalysisType string `json:"analysisType,omitempty"`
	// Status is the status. Allowed values: PENDING, IN_PROGRESS, COMPLETED, FAILED.
	Status string `json:"status,omitempty"`
	// Branch is the branch.
	Branch string `json:"branch,omitempty"`
	// RepositoryUrl is the repository url.
	RepositoryUrl string `json:"repositoryUrl,omitempty"`
	// Revision is the revision.
	Revision string `json:"revision,omitempty"`
	// WorkflowType is the workflow type.
	WorkflowType string `json:"workflowType,omitempty"`
	// FindingsCount is the findings count.
	FindingsCount int32 `json:"findingsCount,omitempty"`
	// CreatedAt is the created at.
	CreatedAt int64 `json:"createdAt,omitempty"`
	// UpdatedAt is the updated at.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// StartedAt is the started at.
	StartedAt int64 `json:"startedAt,omitempty"`
	// FinishedAt is the finished at.
	FinishedAt int64 `json:"finishedAt,omitempty"`
	// FailureReason deprecated. Use errorKey instead.
	FailureReason string `json:"failureReason,omitempty"`
	// ErrorKey stable customer-safe failure identifier. Only set when status is FAILED.
	ErrorKey string `json:"errorKey,omitempty"`
}

// AgenticAgenticJobsSearchRestResponse represents the AgenticJobsSearchRestResponse object of the
// SonarQube V2 API.
type AgenticAgenticJobsSearchRestResponse struct {
	// Jobs is the jobs.
	Jobs []AgenticAgenticJobRestResponse `json:"jobs,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// AgenticDownloadJobLogsOptions contains parameters for the DownloadJobLogs method.
type AgenticDownloadJobLogsOptions struct {
	// JobId is the job id. This field is required.
	JobId string `json:"jobId"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSearchJobsOpt validates the options for the SearchJobs method.
func (s *AgenticService) ValidateSearchJobsOpt(opt *AgenticSearchJobsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.Status, allowedAgenticSearchJobsStatus, "Status")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Type, allowedAgenticSearchJobsType, "Type")
	if err != nil {
		return err
	}

	if opt.PageSize < 0 || opt.PageSize > 500 {
		return NewValidationError("PageSize", "must be between 0 and 500", ErrOutOfRange)
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	return nil
}

// ValidateDownloadJobLogsOpt validates the options for the DownloadJobLogs method.
func (s *AgenticService) ValidateDownloadJobLogsOpt(opt *AgenticDownloadJobLogsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.JobId, "JobId")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// SearchJobs searches agentic jobs. List agentic orchestration jobs (Hunter and Remediation
// agents), most recent first. Requires system administration permission.
//
// API endpoint: GET /api/v2/agentic/jobs.
// Enterprise Edition only.
func (s *AgenticService) SearchJobs(ctx context.Context, opt *AgenticSearchJobsOptions) (*AgenticAgenticJobsSearchRestResponse, *http.Response, error) {
	err := s.ValidateSearchJobsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "agentic/jobs", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AgenticAgenticJobsSearchRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DownloadJobLogs downloads agent job logs. Download the stdout/stderr logs archive for a
// Hunter/Remediation agent job, as a zip file. Requires system administration permission.
//
// API endpoint: GET /api/v2/agentic/jobs/logs.
// Enterprise Edition only.
func (s *AgenticService) DownloadJobLogs(ctx context.Context, opt *AgenticDownloadJobLogsOptions) ([]byte, *http.Response, error) {
	err := s.ValidateDownloadJobLogsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "agentic/jobs/logs", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/zip")

	var buf bytes.Buffer

	resp, err := s.client.Do(req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, nil
}
