package sonar

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// HistoryEntityTypePortfolio identifies a portfolio entity.
	HistoryEntityTypePortfolio = "PORTFOLIO"
	// HistoryEntityTypeProjectBranch identifies a project branch entity.
	HistoryEntityTypeProjectBranch = "PROJECT_BRANCH"
)

// HistoryService handles communication with the history related methods of the
// SonarQube V2 API (the "/history" endpoints: measures history, issue count
// history, and the flattened per-project measure and issue count listings for a
// portfolio or application).
type HistoryService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant set of allowed values
var allowedHistoryMeasuresEntityTypes = map[string]struct{}{
	HistoryEntityTypePortfolio:     {},
	HistoryEntityTypeProjectBranch: {},
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// HistoryIssueCountDistribution represents a single bucket of an issue count
// distribution.
type HistoryIssueCountDistribution struct {
	// Key is the distribution bucket key. It is empty when the counts are not
	// grouped by any dimension.
	Key string `json:"key,omitempty"`
	// Value is the number of issues in the bucket.
	Value int32 `json:"value,omitempty"`
}

// HistoryIssueCount represents the issue count distribution recorded at a given
// point in time.
type HistoryIssueCount struct {
	// Date is the date and time in ISO 8601 format, including timezone offset.
	Date string `json:"date,omitempty"`
	// Distribution is the list of issue count buckets for the date.
	Distribution []HistoryIssueCountDistribution `json:"distribution,omitempty"`
}

// HistoryIssueCountHistoryResponse represents the response from GetIssueCountHistory.
type HistoryIssueCountHistoryResponse struct {
	// IssueCountHistory is the ordered list of issue count snapshots.
	IssueCountHistory []HistoryIssueCount `json:"issueCountHistory,omitempty"`
}

// HistoryMeasureEntry represents a single measure value inside a measures
// history item.
type HistoryMeasureEntry struct {
	// Metric is the metric whose value is being recorded (e.g. "coverage").
	Metric string `json:"metric,omitempty"`
	// Type is the type of value the measure contains (e.g. "PERCENT").
	Type string `json:"type,omitempty"`
	// Value is the value of the measure encoded as a string.
	Value string `json:"value,omitempty"`
}

// HistoryMeasureItem represents the measures recorded at a given point in time.
type HistoryMeasureItem struct {
	// Date is the date the measures were taken.
	Date string `json:"date,omitempty"`
	// Measures is the list of measure values for the date.
	Measures []HistoryMeasureEntry `json:"measures,omitempty"`
}

// HistoryMeasuresHistoryResponse represents the response from GetMeasuresHistory.
type HistoryMeasuresHistoryResponse struct {
	// MeasuresHistory is the ordered list of measure snapshots.
	MeasuresHistory []HistoryMeasureItem `json:"measuresHistory,omitempty"`
}

// HistoryProjectIssueCount represents the issue count for a single project in a
// flattened portfolio or application listing.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryProjectIssueCount struct {
	// BranchID is the identifier of the project branch.
	BranchID string `json:"branchId,omitempty"`
	// ProjectName is the display name of the project.
	ProjectName string `json:"projectName,omitempty"`
	// ProjectKey is the key of the project.
	ProjectKey string `json:"projectKey,omitempty"`
	// BranchName is the name of the branch.
	BranchName string `json:"branchName,omitempty"`
	// IssueCount is the number of issues in the project branch.
	IssueCount int64 `json:"issueCount,omitempty"`
	// ReferenceIssueCount is the issue count at the reference date, used for
	// trend computation. It is nil when no reference date was supplied.
	ReferenceIssueCount *int64 `json:"referenceIssueCount,omitempty"`
}

// HistoryProjectIssueCountsResponse represents the response from GetProjectIssueCounts.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryProjectIssueCountsResponse struct {
	// HiddenProjectCount is the number of projects excluded from the response
	// because the user does not have permission to view them.
	HiddenProjectCount int32 `json:"hiddenProjectCount,omitempty"`
	// ProjectIssueCounts is the flattened list of per-project issue counts.
	ProjectIssueCounts []HistoryProjectIssueCount `json:"projectIssueCounts,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// HistoryProjectMeasureMetric represents the current and reference values of a
// single metric for a project.
type HistoryProjectMeasureMetric struct {
	// CurrentValue is the current value of the measure encoded as a string. It
	// is empty when the measure has no current value.
	CurrentValue string `json:"currentValue,omitempty"`
	// Metric is the metric key (e.g. "coverage").
	Metric string `json:"metric,omitempty"`
	// ReferenceValue is the value of the measure at the reference date, used for
	// trend computation. It is empty when there was no value on that date.
	ReferenceValue string `json:"referenceValue,omitempty"`
	// Type identifies the measure's data type (e.g. "PERCENT").
	Type string `json:"type,omitempty"`
}

// HistoryProjectMeasure represents the measure value for a single project in a
// flattened portfolio or application listing.
type HistoryProjectMeasure struct {
	// BranchID is the identifier of the project branch.
	BranchID string `json:"branchId,omitempty"`
	// BranchName is the name of the branch.
	BranchName string `json:"branchName,omitempty"`
	// Measure is the metric value for the project.
	Measure HistoryProjectMeasureMetric `json:"measure,omitzero"`
	// ProjectKey is the key of the project.
	ProjectKey string `json:"projectKey,omitempty"`
	// ProjectName is the display name of the project.
	ProjectName string `json:"projectName,omitempty"`
}

// HistoryProjectMeasuresResponse represents the response from GetProjectMeasures.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryProjectMeasuresResponse struct {
	// HiddenProjectCount is the number of projects excluded from the response
	// because the user does not have permission to view them.
	HiddenProjectCount int32 `json:"hiddenProjectCount,omitempty"`
	// ProjectMeasures is the flattened list of per-project measure values.
	ProjectMeasures []HistoryProjectMeasure `json:"projectMeasures,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// HistoryIssueCountHistoryOptions contains parameters for GetIssueCountHistory.
type HistoryIssueCountHistoryOptions struct {
	// EntityID is the UUID of the entity to retrieve issue count history for.
	// This field is required.
	EntityID string `json:"entityId"`
	// EntityType is the entity's type. This field is required.
	EntityType string `json:"entityType"`
	// StartDate is the inclusive start of the date range in ISO 8601 format.
	// This field is required and may not be more than 1 year in the past.
	StartDate string `json:"startDate"`
	// EndDate is the inclusive end of the date range in ISO 8601 format.
	// Defaults to the current date and time (UTC) when omitted.
	EndDate string `json:"endDate,omitempty"`
	// Impacts limits results to issues with specific impacts.
	Impacts []string `json:"impacts,omitempty"`
	// IssueTypes limits results to issues with specific types.
	IssueTypes []string `json:"issueTypes,omitempty"`
	// RuleKeys limits results to issues raised by specific rules.
	RuleKeys []string `json:"ruleKeys,omitempty"`
	// Severities limits results to issues with specific severities.
	Severities []string `json:"severities,omitempty"`
	// SliceBy optionally groups issue counts by a specific dimension.
	SliceBy string `json:"sliceBy,omitempty"`
	// Statuses limits results to issues with specific statuses.
	Statuses []string `json:"statuses,omitempty"`
}

// HistoryMeasuresHistoryOptions contains parameters for GetMeasuresHistory.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryMeasuresHistoryOptions struct {
	// EntityType is the type of entity to retrieve measures history for.
	// This field is required. Allowed values: PORTFOLIO, PROJECT_BRANCH.
	EntityType string `json:"entityType"`
	// EntityID is the UUID of the entity. This field is required.
	EntityID string `json:"entityId"`
	// MetricKeys is the list of metric keys to filter by. This field is required.
	MetricKeys []string `json:"metricKeys"`
	// StartDate is the inclusive start of the date range in ISO 8601 format.
	// This field is required.
	StartDate string `json:"startDate"`
	// EndDate is the inclusive end of the date range in ISO 8601 format.
	// Defaults to the current date when omitted.
	EndDate string `json:"endDate,omitempty"`
}

// HistoryProjectIssueCountsOptions contains parameters for GetProjectIssueCounts.
//
// The selection is identified with either the legacy PortfolioID parameter, or
// with both EntityType and EntityID. The two selector forms are mutually
// exclusive.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryProjectIssueCountsOptions struct {
	// PortfolioID is the legacy portfolio ID.
	PortfolioID string `json:"portfolioId,omitempty"`
	// EntityType is the type of the selection whose projects are returned as a
	// flattened list. Must be provided together with EntityID.
	EntityType string `json:"entityType,omitempty"`
	// EntityID is the portfolio ID or application branch ID, according to
	// EntityType. Must be provided together with EntityType.
	EntityID string `json:"entityId,omitempty"`
	// RuleKeys limits results to issues raised by specific rules.
	RuleKeys []string `json:"ruleKeys,omitempty"`
	// Severities limits results to issues with specific severities.
	Severities []string `json:"severities,omitempty"`
	// IssueTypes limits results to issues with specific types.
	IssueTypes []string `json:"issueTypes,omitempty"`
	// Statuses limits results to issues with specific statuses.
	Statuses []string `json:"statuses,omitempty"`
	// Impacts limits results to issues with specific impacts.
	Impacts []string `json:"impacts,omitempty"`
	// NameContains is a case-insensitive search query for project names.
	NameContains string `json:"nameContains,omitempty"`
	// ReferenceDate is the date from which reference measure values are returned
	// for computing trends. Reference values are omitted when this is empty.
	ReferenceDate string `json:"referenceDate,omitempty"`
	// PageIndex is the 1-based index of the page to fetch. Default is 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the size of the page to fetch. Default is 50.
	PageSize int32 `json:"pageSize,omitempty"`
	// Sort defines the sort order. Each field is prefixed with '+' for ascending
	// or '-' for descending. Valid fields: issueCount, projectName.
	Sort []string `json:"sort,omitempty"`
	// RequireIssues only returns projects with an issue count greater than 0.
	RequireIssues bool `json:"requireIssues,omitempty"`
}

// HistoryProjectMeasuresOptions contains parameters for GetProjectMeasures.
//
// The selection is identified with either the legacy PortfolioID parameter, or
// with both EntityType and EntityID. The two selector forms are mutually
// exclusive.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryProjectMeasuresOptions struct {
	// MetricKey is the metric key. This field is required.
	MetricKey string `json:"metricKey"`
	// MetricValue filters by exact metric value. Only primitive metric types are
	// supported (not DISTR or DATA).
	MetricValue string `json:"metricValue,omitempty"`
	// NameContains is a case-insensitive search query for project names.
	NameContains string `json:"nameContains,omitempty"`
	// PageIndex is the 1-based index of the page to fetch. Default is 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the size of the page to fetch. Default is 50.
	PageSize int32 `json:"pageSize,omitempty"`
	// PortfolioID is the legacy portfolio ID.
	PortfolioID string `json:"portfolioId,omitempty"`
	// EntityType is the type of the selection whose projects are returned as a
	// flattened list. Must be provided together with EntityID.
	EntityType string `json:"entityType,omitempty"`
	// EntityID is the portfolio ID or application branch ID, according to
	// EntityType. Must be provided together with EntityType.
	EntityID string `json:"entityId,omitempty"`
	// ReferenceDate is the date from which reference measure values are returned
	// for computing trends. Reference values are omitted when this is empty.
	ReferenceDate string `json:"referenceDate,omitempty"`
	// Sort defines the sort order. Each field is prefixed with '+' for ascending
	// or '-' for descending. Valid fields: measure.currentValue, projectName.
	Sort []string `json:"sort,omitempty"`
	// RequireValue only returns projects with a non-null metric value.
	RequireValue bool `json:"requireValue,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateIssueCountHistoryOpt validates the options for GetIssueCountHistory.
func (s *HistoryService) ValidateIssueCountHistoryOpt(opt *HistoryIssueCountHistoryOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.EntityID, "EntityID")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.EntityType, "EntityType")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.StartDate, "StartDate")
}

// ValidateMeasuresHistoryOpt validates the options for GetMeasuresHistory.
func (s *HistoryService) ValidateMeasuresHistoryOpt(opt *HistoryMeasuresHistoryOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.EntityType, "EntityType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EntityType, allowedHistoryMeasuresEntityTypes, "EntityType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.EntityID, "EntityID")
	if err != nil {
		return err
	}

	if len(opt.MetricKeys) == 0 {
		return NewValidationError("MetricKeys", "must not be empty", ErrMissingRequired)
	}

	return ValidateRequired(opt.StartDate, "StartDate")
}

// validateHistorySelector validates the mutually exclusive portfolioId /
// (entityType, entityId) selector shared by the project-scoped endpoints.
func validateHistorySelector(portfolioID, entityType, entityID string) error {
	entityForm := entityType != "" || entityID != ""

	if portfolioID != "" && entityForm {
		return NewValidationError("PortfolioID", "cannot be combined with EntityType/EntityID", ErrInvalidValue)
	}

	if portfolioID == "" && !entityForm {
		return NewValidationError("PortfolioID", "either PortfolioID or both EntityType and EntityID must be provided", ErrMissingRequired)
	}

	if entityForm && (entityType == "" || entityID == "") {
		return NewValidationError("EntityType", "EntityType and EntityID must be provided together", ErrMissingRequired)
	}

	return nil
}

// ValidateProjectIssueCountsOpt validates the options for GetProjectIssueCounts.
func (s *HistoryService) ValidateProjectIssueCountsOpt(opt *HistoryProjectIssueCountsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return validateHistorySelector(opt.PortfolioID, opt.EntityType, opt.EntityID)
}

// ValidateProjectMeasuresOpt validates the options for GetProjectMeasures.
func (s *HistoryService) ValidateProjectMeasuresOpt(opt *HistoryProjectMeasuresOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.MetricKey, "MetricKey")
	if err != nil {
		return err
	}

	return validateHistorySelector(opt.PortfolioID, opt.EntityType, opt.EntityID)
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetIssueCountHistory returns the history of issue counts for an entity. The
// maximum resolution for history data is one day and all history is recorded in
// the UTC timezone. Issue counts carry forward for days with no changes.
//
// API endpoint: GET /api/v2/history/issue-count-history.
func (s *HistoryService) GetIssueCountHistory(ctx context.Context, opt *HistoryIssueCountHistoryOptions) (*HistoryIssueCountHistoryResponse, *http.Response, error) {
	err := s.ValidateIssueCountHistoryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/issue-count-history", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryIssueCountHistoryResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetMeasuresHistory returns the history of measures for the given entity.
//
// API endpoint: GET /api/v2/history/measures-history.
func (s *HistoryService) GetMeasuresHistory(ctx context.Context, opt *HistoryMeasuresHistoryOptions) (*HistoryMeasuresHistoryResponse, *http.Response, error) {
	err := s.ValidateMeasuresHistoryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/measures-history", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryMeasuresHistoryResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectIssueCounts returns issue counts for the flattened list of projects
// in a portfolio or application. The selection is identified with either the
// legacy PortfolioID parameter, or with both EntityType and EntityID.
//
// API endpoint: GET /api/v2/history/project-issue-counts.
func (s *HistoryService) GetProjectIssueCounts(ctx context.Context, opt *HistoryProjectIssueCountsOptions) (*HistoryProjectIssueCountsResponse, *http.Response, error) {
	err := s.ValidateProjectIssueCountsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/project-issue-counts", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryProjectIssueCountsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectMeasures returns current measure values for the flattened list of
// projects in a portfolio or application. The selection is identified with
// either the legacy PortfolioID parameter, or with both EntityType and EntityID.
//
// API endpoint: GET /api/v2/history/project-measures.
func (s *HistoryService) GetProjectMeasures(ctx context.Context, opt *HistoryProjectMeasuresOptions) (*HistoryProjectMeasuresResponse, *http.Response, error) {
	err := s.ValidateProjectMeasuresOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/project-measures", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryProjectMeasuresResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
