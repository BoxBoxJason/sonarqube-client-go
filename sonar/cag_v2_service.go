package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CagService handles communication with the Context Augmentation (CAG) related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type CagService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedCagGetUsageStatsPeriod is the set of allowed values for the corresponding option.
	allowedCagGetUsageStatsPeriod = map[string]struct{}{
		"LAST_DAY":     {},
		"LAST_7_DAYS":  {},
		"LAST_30_DAYS": {},
	}
	// allowedCagRecordImpactEventEventType is the set of allowed values for the corresponding option.
	allowedCagRecordImpactEventEventType = map[string]struct{}{
		"NAVIGATION":            {},
		"OUTPUT_COMPRESSION":    {},
		"OUTPUT_RESTORE":        {},
		"GUIDELINES":            {},
		"DEPENDENCY_CHECK":      {},
		"INTENDED_ARCHITECTURE": {},
	}
	// allowedCagRecordImpactEventTransport is the set of allowed values for the corresponding option.
	allowedCagRecordImpactEventTransport = map[string]struct{}{
		"MCP": {},
		"CLI": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// CagCagEntitlementConsumption represents the CagEntitlementConsumption object of the SonarQube V2
// API.
type CagCagEntitlementConsumption struct {
	// Consumed is the consumed.
	Consumed int64 `json:"consumed,omitempty"`
	// Limit is the limit.
	Limit int64 `json:"limit,omitempty"`
}

// CagCagEntitlement represents the CagEntitlement object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CagCagEntitlement struct {
	// Allowed whether the organization is currently allowed to use Context Augmentation.
	Allowed bool `json:"allowed,omitempty"`
	// HasEntitlement whether the organization has any valid access path to Context Augmentation.
	// When allowed is false, true means the organization is still entitled but cannot use it right
	// now, typically because it has reached its current limit; false means the organization is not
	// entitled at all, or Context Augmentation is disabled instance-wide (reason disabled), in
	// which case entitlement is not evaluated. See reason for the exact cause.
	HasEntitlement bool `json:"hasEntitlement,omitempty"`
	// FeatureKey feature key to use for consumption reporting when relevant. This is always
	// 'contextAugmentation' for the CAG entitlement response.
	FeatureKey string `json:"featureKey,omitempty"`
	// Consumption is the consumption.
	Consumption CagCagEntitlementConsumption `json:"consumption,omitzero"`
	// Reason is the reason. Allowed values: not_entitled, disabled, blocked.
	Reason string `json:"reason,omitempty"`
}

// CagGetUsageStatsOptions contains parameters for the GetUsageStats method.
type CagGetUsageStatsOptions struct {
	// Period for the usage statistics. This field is required. Allowed values: LAST_DAY,
	// LAST_7_DAYS, LAST_30_DAYS.
	Period string `json:"period"`
}

// CagProjectCagUsageCount represents a SonarQube V2 API object.
//
// CAG tool call count for a project.
type CagProjectCagUsageCount struct {
	// Identifier of the project; an opaque string on SonarQube Server, a UUID string on SonarQube
	// Cloud.
	Id string `json:"id,omitempty"`
	// Name the project name.
	Name string `json:"name,omitempty"`
	// ToolCallCount number of CAG tool calls for the project over the requested period.
	ToolCallCount int64 `json:"toolCallCount,omitempty"`
}

// CagCagUsageStats represents a SonarQube V2 API object.
//
// CAG tool usage statistics for an organization.
type CagCagUsageStats struct {
	// Id the organization ID these statistics are about.
	Id string `json:"id,omitempty"`
	// Projects all projects with their CAG tool call count, ordered by tool call count descending.
	Projects []CagProjectCagUsageCount `json:"projects,omitempty"`
}

// CagGetGuideMetricsOptions contains parameters for the GetGuideMetrics method.
type CagGetGuideMetricsOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// From inclusive start of the calendar date range (YYYY-MM-DD, UTC). This field is required.
	From string `json:"from"`
	// To inclusive end of the calendar date range (YYYY-MM-DD, UTC). Must not precede from. This
	// field is required.
	To string `json:"to"`
}

// CagGuideMetrics represents a SonarQube V2 API object.
//
// Guide impact metrics for the requested period.
type CagGuideMetrics struct {
	// PreventedIssues estimated issues prevented by supplied guidelines, rounded to a whole
	// number. Each successful get_guidelines call that returns at least one rule contributes
	// N_prevented.
	PreventedIssues int64 `json:"preventedIssues,omitempty"`
	// TokensSaved estimated tokens saved by navigation tools and output compression, preserving
	// fractional values. Aggregate estimates before rounding for display.
	TokensSaved float64 `json:"tokensSaved,omitempty"`
	// EngineeringMinutesSaved guide contribution to estimated engineering time saved in minutes,
	// preserving fractional values. Add this to the verifyMetrics value before rounding the
	// combined total for display.
	EngineeringMinutesSaved float64 `json:"engineeringMinutesSaved,omitempty"`
}

// CagGetProjectActivityOptions contains parameters for the GetProjectActivity method.
type CagGetProjectActivityOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// From inclusive start of the calendar date range (YYYY-MM-DD, UTC). This field is required.
	From string `json:"from"`
	// To inclusive end of the calendar date range (YYYY-MM-DD, UTC). Must not precede from. This
	// field is required.
	To string `json:"to"`
}

// CagVortexProjectUsage represents a SonarQube V2 API object.
//
// Projects classified by successful Vortex usage during the requested period.
type CagVortexProjectUsage struct {
	// WithUsageProjectCount projects with successful CAG or Verify usage.
	WithUsageProjectCount int32 `json:"withUsageProjectCount,omitempty"`
	// WithoutUsageProjectCount projects without successful CAG or Verify usage.
	WithoutUsageProjectCount int32 `json:"withoutUsageProjectCount,omitempty"`
	// CagOnlyProjectCount projects with successful CAG usage but no successful Verify usage.
	CagOnlyProjectCount int32 `json:"cagOnlyProjectCount,omitempty"`
	// VerifyOnlyProjectCount projects with successful Verify usage but no successful CAG usage.
	VerifyOnlyProjectCount int32 `json:"verifyOnlyProjectCount,omitempty"`
	// CagAndVerifyProjectCount projects with both successful CAG and Verify usage.
	CagAndVerifyProjectCount int32 `json:"cagAndVerifyProjectCount,omitempty"`
}

// CagVortexProjectAnalysisCohort represents a SonarQube V2 API object.
//
// Project count and Vortex usage for one analysis cohort.
type CagVortexProjectAnalysisCohort struct {
	// ProjectCount projects in this analysis cohort.
	ProjectCount int32 `json:"projectCount,omitempty"`
	// VortexUsage is the vortex usage.
	VortexUsage CagVortexProjectUsage `json:"vortexUsage,omitzero"`
}

// CagVortexProjectActivity represents a SonarQube V2 API object.
//
// Vortex usage grouped by project analysis presence during the requested period.
type CagVortexProjectActivity struct {
	// TotalProjectCount all projects across the four analysis cohorts.
	TotalProjectCount int32 `json:"totalProjectCount,omitempty"`
	// WithMainBranchAnalysisInPeriod projects with at least one processed main-branch analysis in
	// the requested period.
	WithMainBranchAnalysisInPeriod CagVortexProjectAnalysisCohort `json:"withMainBranchAnalysisInPeriod,omitzero"`
	// WithLongLivedBranchAnalysisInPeriod projects with at least one processed long-lived-branch
	// analysis, but no processed main-branch analysis, in the requested period. Pull-request
	// analyses may also be present.
	WithLongLivedBranchAnalysisInPeriod CagVortexProjectAnalysisCohort `json:"withLongLivedBranchAnalysisInPeriod,omitzero"`
	// WithOnlyPullRequestAnalysisInPeriod projects with at least one processed pull-request
	// analysis, but no processed main-branch or long-lived-branch analysis, in the requested
	// period.
	WithOnlyPullRequestAnalysisInPeriod CagVortexProjectAnalysisCohort `json:"withOnlyPullRequestAnalysisInPeriod,omitzero"`
	// WithoutAnalysisInPeriod projects without a processed analysis on any branch in the requested
	// period.
	WithoutAnalysisInPeriod CagVortexProjectAnalysisCohort `json:"withoutAnalysisInPeriod,omitzero"`
}

// CagGetProjectImpactOptions contains parameters for the GetProjectImpact method.
type CagGetProjectImpactOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// From inclusive start of the calendar date range (YYYY-MM-DD, UTC). This field is required.
	From string `json:"from"`
	// To inclusive end of the calendar date range (YYYY-MM-DD, UTC). Must not precede from. This
	// field is required.
	To string `json:"to"`
	// PageIndex page to retrieve, 1-based. Defaults to 1.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize number of projects per page. Defaults to 50. Use 0 to return only the total count.
	// Must be between 0 and 500. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// CagVortexProjectImpact represents a SonarQube V2 API object.
//
// Vortex impact metrics for a single project in the requested period. Successful output compression
// and restoration both establish CAG activity and latest activity, but neither increases the Guide
// call count. Compression contributes token savings; restoration adds no impact and does not
// reverse compression savings in V0. Neither contributes engineering savings or resolved issues. A
// project using only these capabilities can therefore be active with zero tool calls; a
// restoration-only project also has zero calculated savings.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CagVortexProjectImpact struct {
	// ProjectId identifier of the project.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectName display name of the project.
	ProjectName string `json:"projectName,omitempty"`
	// ActivityClassification is the activity classification. Allowed values: NONE, CAG_ONLY,
	// VERIFY_ONLY, BOTH.
	ActivityClassification string `json:"activityClassification,omitempty"`
	// GuideCallCount successful CAG tool calls in the period, counted from stored event metadata
	// independently of payload decoding. Compression and restoration count as activity but not as
	// Guide calls. Shell compression is hook-driven and expected to occur much more frequently
	// than model-initiated tool calls, so including it would inflate this count. Successful events
	// with malformed or unsupported payloads remain counted but contribute no calculated impact.
	GuideCallCount int64 `json:"guideCallCount,omitempty"`
	// VerifyCallCount successful public SQAA analyses in the period, counted from stored event
	// metadata independently of payload decoding. Successful events with malformed or unsupported
	// payloads remain counted but contribute no calculated impact.
	VerifyCallCount int64 `json:"verifyCallCount,omitempty"`
	// ResolvedIssues issues inferred to be resolved in the period.
	ResolvedIssues int64 `json:"resolvedIssues,omitempty"`
	// EngineeringMinutesSaved combined Guide and Verify contribution to estimated engineering time
	// saved in minutes, preserving fractional values. Aggregate estimates before rounding for
	// display.
	EngineeringMinutesSaved float64 `json:"engineeringMinutesSaved,omitempty"`
	// TokensSaved estimated tokens saved by navigation tools and output compression for this
	// project, preserving fractional values. Aggregate estimates before rounding for display.
	// Output restoration does not add to or reverse compression savings in V0.
	TokensSaved float64 `json:"tokensSaved,omitempty"`
	// LatestActivityAt server-recorded time of the latest successful Vortex event in the period.
	// Null when there is no activity.
	LatestActivityAt string `json:"latestActivityAt,omitempty"`
}

// CagVortexProjectPage represents a SonarQube V2 API object.
//
// Paginated Vortex impact per project for the requested period.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CagVortexProjectPage struct {
	// PageIndex current page index (1-based).
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize number of projects per page as resolved for this response.
	PageSize int32 `json:"pageSize,omitempty"`
	// Total number of projects, including those with no Vortex activity.
	Total int32 `json:"total,omitempty"`
	// Projects is the projects.
	Projects []CagVortexProjectImpact `json:"projects,omitempty"`
}

// CagGetVerifyMetricsOptions contains parameters for the GetVerifyMetrics method.
type CagGetVerifyMetricsOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// From inclusive start of the calendar date range (YYYY-MM-DD, UTC). This field is required.
	From string `json:"from"`
	// To inclusive end of the calendar date range (YYYY-MM-DD, UTC). Must not precede from. This
	// field is required.
	To string `json:"to"`
}

// CagVerifyMetrics represents a SonarQube V2 API object.
//
// Verify impact metrics for the requested period.
type CagVerifyMetrics struct {
	// ResolvedIssues issues inferred to be resolved across all SQAA analysis sequences. Detected
	// as issue-count decreases between successive successful analyses of the same (file, rule) in
	// inferred sessions.
	ResolvedIssues int64 `json:"resolvedIssues,omitempty"`
	// EngineeringMinutesSaved verify contribution to estimated engineering time saved in minutes,
	// preserving fractional values. Add this to the guideMetrics value before rounding the
	// combined total for display.
	EngineeringMinutesSaved float64 `json:"engineeringMinutesSaved,omitempty"`
}

// CagPingResponse represents the PingResponse object of the SonarQube V2 API.
type CagPingResponse struct {
	// Status service status.
	Status string `json:"status,omitempty"`
}

// CagRecordImpactEventOptions contains the request body for the RecordImpactEvent method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CagRecordImpactEventOptions struct {
	// OrganizationId organization UUID. This field is required.
	OrganizationId string `json:"organizationId"`
	// ProjectId project identifier. Not format: uuid — Server project ids are opaque strings;
	// Cloud sends a UUID string, accepted unchanged. This field is required.
	ProjectId string `json:"projectId"`
	// InvocationId invocation identifier. Natural deduplication key. This field is required.
	InvocationId string `json:"invocationId"`
	// CagInstanceId shared by all calls from the same CAG instance (MCP server or CLI daemon).
	// This field is required.
	CagInstanceId string `json:"cagInstanceId"`
	// EventType is the event type. This field is required. Allowed values: NAVIGATION,
	// OUTPUT_COMPRESSION, OUTPUT_RESTORE, GUIDELINES, DEPENDENCY_CHECK, INTENDED_ARCHITECTURE.
	EventType string `json:"eventType"`
	// EventVersion payload schema version. Only 1 is supported. This field is required.
	EventVersion int32 `json:"eventVersion"`
	// Success true when the call completed successfully. This field is required.
	Success bool `json:"success"`
	// Transport is the transport. This field is required. Allowed values: MCP, CLI.
	Transport string `json:"transport"`
	// BranchName git branch at invocation time. Optional.
	BranchName string `json:"branchName,omitempty"`
	// Payload event-specific fields. Schema is selected by (eventType, eventVersion). See
	// NavigationPayloadV1, OutputCompressionPayloadV1, OutputRestorePayloadV1,
	// GuidelinesPayloadV1, DependencyCheckPayloadV1, IntendedArchitecturePayloadV1. This field is
	// required.
	Payload map[string]any `json:"payload"`
}

// CagRecordUsageEventOptions contains the request body for the RecordUsageEvent method.
type CagRecordUsageEventOptions struct {
	// OrganizationId UUID v4 of the organization. On SonarQube Server, the '@OrganizationId'
	// annotation resolves this value server-side to the caller's default organization id, so
	// callers on that host may pass any syntactically valid UUID as a placeholder — the value must
	// still parse as a UUID, but its content is otherwise ignored. This field is required.
	OrganizationId string `json:"organizationId"`
	// ProjectId identifier of the project the tool call was about. Deliberately not 'format:
	// uuid': SonarQube Server project ids are opaque strings, not UUIDs, so a UUID-typed field
	// would reject every real Server project id. SonarQube Cloud sends a UUID string, which this
	// schema accepts unchanged since the JSON wire format is identical either way. This field is
	// required.
	ProjectId string `json:"projectId"`
	// InvocationId UUID v4 of the tool call, generated by the client. Natural deduplication key.
	// This field is required.
	InvocationId string `json:"invocationId"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetUsageStatsOpt validates the options for the GetUsageStats method.
func (s *CagService) ValidateGetUsageStatsOpt(opt *CagGetUsageStatsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Period, "Period")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Period, allowedCagGetUsageStatsPeriod, "Period")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetGuideMetricsOpt validates the options for the GetGuideMetrics method.
func (s *CagService) ValidateGetGuideMetricsOpt(opt *CagGetGuideMetricsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.From, "From")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.To, "To")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetProjectActivityOpt validates the options for the GetProjectActivity method.
func (s *CagService) ValidateGetProjectActivityOpt(opt *CagGetProjectActivityOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.From, "From")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.To, "To")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetProjectImpactOpt validates the options for the GetProjectImpact method.
func (s *CagService) ValidateGetProjectImpactOpt(opt *CagGetProjectImpactOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.From, "From")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.To, "To")
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

// ValidateGetVerifyMetricsOpt validates the options for the GetVerifyMetrics method.
func (s *CagService) ValidateGetVerifyMetricsOpt(opt *CagGetVerifyMetricsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.From, "From")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.To, "To")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRecordImpactEventOpt validates the options for the RecordImpactEvent method.
func (s *CagService) ValidateRecordImpactEventOpt(opt *CagRecordImpactEventOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.InvocationId, "InvocationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.CagInstanceId, "CagInstanceId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.EventType, "EventType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EventType, allowedCagRecordImpactEventEventType, "EventType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Transport, "Transport")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Transport, allowedCagRecordImpactEventTransport, "Transport")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRecordUsageEventOpt validates the options for the RecordUsageEvent method.
func (s *CagService) ValidateRecordUsageEventOpt(opt *CagRecordUsageEventOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.InvocationId, "InvocationId")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetEntitlement gets the CAG entitlement for an organization. Returns the Context Augmentation
// entitlement decision for an organization. Response shape (allowed / hasEntitlement / featureKey /
// consumption) matches the existing SonarQube Cloud contract so clients are unaffected by which
// host answers the request. An unresolvable check (license unreadable, billing unreachable) is 503,
// not a not-entitled 200.
//
// API endpoint: GET /api/v2/cag/cag-entitlement/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) GetEntitlement(ctx context.Context, organizationID string) (*CagCagEntitlement, *http.Response, error) {
	err := ValidateRequired(organizationID, "organizationID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/cag-entitlement/"+url.PathEscape(organizationID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagCagEntitlement)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetUsageStats gets CAG tool usage statistics for an organization. Returns, for every project, the
// number of Context Augmentation tool calls recorded over the requested period. Response shape (id
// / projects[].id,name,toolCallCount) matches the existing SonarQube Cloud contract so clients are
// unaffected by which host answers the request. Projects with no usage in the period are still
// returned, with a toolCallCount of zero, and the list is ordered by toolCallCount descending.
//
// API endpoint: GET /api/v2/cag/cag-usage-stats/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) GetUsageStats(ctx context.Context, organizationID string, opt *CagGetUsageStatsOptions) (*CagCagUsageStats, *http.Response, error) {
	err := ValidateRequired(organizationID, "organizationID")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateGetUsageStatsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/cag-usage-stats/"+url.PathEscape(organizationID), opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagCagUsageStats)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetGuideMetrics gets Guide impact metrics for the organization. Returns prevented issues, tokens
// saved, and the Guide contribution to engineering time saved for the requested calendar date
// range. Covers navigation tools, output compression, get_guidelines, check_dependency, and
// get_intended_architecture events. Only successful events contribute to calculated values. Dates
// are inclusive UTC calendar boundaries: from=2026-08-01 to=2026-08-31 selects [2026-08-01T00:00Z,
// 2026-09-01T00:00Z).
//
// API endpoint: GET /api/v2/cag/impact/guide-metrics.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) GetGuideMetrics(ctx context.Context, opt *CagGetGuideMetricsOptions) (*CagGuideMetrics, *http.Response, error) {
	err := s.ValidateGetGuideMetricsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/impact/guide-metrics", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagGuideMetrics)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectActivity gets Vortex usage by project analysis cohort for the organization. Groups
// projects into mutually exclusive main-branch, long-lived-branch, pull-request-only, and no-
// analysis cohorts for the requested period, then reports Vortex usage within each cohort. Analysis
// precedence follows that order. Successful output-compression and restoration events establish CAG
// usage; failed events do not. Date boundaries follow the same UTC-calendar interpretation as
// /cag/impact/guide-metrics. Only periods ending on the current date are supported. To tolerate
// browser and Server timezone differences, a 'to' value within one day of the current UTC date is
// accepted and normalized to the current UTC date.
//
// API endpoint: GET /api/v2/cag/impact/project-activity.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) GetProjectActivity(ctx context.Context, opt *CagGetProjectActivityOptions) (*CagVortexProjectActivity, *http.Response, error) {
	err := s.ValidateGetProjectActivityOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/impact/project-activity", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagVortexProjectActivity)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectImpact gets paginated Vortex impact per project for the organization. Returns every
// project, including projects with no Vortex activity in the period. Each row combines Guide and
// Verify impact. The page is ordered by total Guide plus Verify calls descending, then latest
// successful activity descending with no activity last, then project key ascending. Successful
// compression and restoration establish activity without increasing Guide calls. Sorting and
// filtering by calculated impact values are outside the V0 contract. Date boundaries follow the
// same UTC-calendar interpretation as /cag/impact/guide-metrics.
//
// API endpoint: GET /api/v2/cag/impact/projects.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) GetProjectImpact(ctx context.Context, opt *CagGetProjectImpactOptions) (*CagVortexProjectPage, *http.Response, error) {
	err := s.ValidateGetProjectImpactOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/impact/projects", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagVortexProjectPage)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetVerifyMetrics gets Verify impact metrics for the organization. Returns resolved issues, and
// the Verify contribution to engineering time saved, for the requested calendar date range. Groups
// public SQAA analysis events into inferred sessions and detects issue-count decreases between
// successive successful analyses. Date boundaries follow the same UTC-calendar interpretation as
// /cag/impact/guide-metrics.
//
// API endpoint: GET /api/v2/cag/impact/verify-metrics.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) GetVerifyMetrics(ctx context.Context, opt *CagGetVerifyMetricsOptions) (*CagVerifyMetrics, *http.Response, error) {
	err := s.ValidateGetVerifyMetricsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/impact/verify-metrics", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagVerifyMetrics)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Ping authenticateds liveness probe for the CAG Hub service. Validates the caller's SonarCloud /
// SonarQube Server user token and returns a generic success payload. Zero-dependency authenticated
// liveness probe shared by both hosts. No business logic.
//
// API endpoint: GET /api/v2/cag/ping.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) Ping(ctx context.Context) (*CagPingResponse, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "cag/ping", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(CagPingResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RecordImpactEvent records a CAG Vortex impact event. Records one versioned Vortex impact event.
// Null or missing eventType and eventVersion return 400. The payload must conform to the named
// payload schema for the (eventType, eventVersion) pair: NavigationPayloadV1,
// OutputCompressionPayloadV1, OutputRestorePayloadV1, GuidelinesPayloadV1,
// DependencyCheckPayloadV1, or IntendedArchitecturePayloadV1 (all in components/schemas).
//
// API endpoint: POST /api/v2/cag/cag-impact-events.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) RecordImpactEvent(ctx context.Context, opt *CagRecordImpactEventOptions) (*http.Response, error) {
	err := s.ValidateRecordImpactEventOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "cag/cag-impact-events", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RecordUsageEvent records a CAG tool usage event. Records a single Context Augmentation tool-call
// event. The event is acknowledged with 202 once accepted; a storage failure answers 500. Request
// shape (organizationId / projectId / invocationId) matches the existing SonarQube Cloud contract
// so clients are unaffected by which host answers the request.
//
// API endpoint: POST /api/v2/cag/cag-usage-events.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *CagService) RecordUsageEvent(ctx context.Context, opt *CagRecordUsageEventOptions) (*http.Response, error) {
	err := s.ValidateRecordUsageEventOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "cag/cag-usage-events", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
