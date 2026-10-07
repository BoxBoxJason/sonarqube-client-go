package sonar

import (
	"context"
	"fmt"
	"net/http"
)

//nolint:gochecknoglobals,goconst // constant sets of allowed values
var (
	// allowedHistoryGetIssueDensityHistoryEntityType is the set of allowed values for the corresponding option.
	allowedHistoryGetIssueDensityHistoryEntityType = map[string]struct{}{
		"PORTFOLIO":      {},
		"PROJECT_BRANCH": {},
		"APPLICATION":    {},
	}
	// allowedHistoryGetIssueDensityHistorySliceBy is the set of allowed values for the corresponding option.
	allowedHistoryGetIssueDensityHistorySliceBy = map[string]struct{}{
		"RULE_KEY":         {},
		"SEVERITY":         {},
		"TYPE_SEVERITY":    {},
		"SOFTWARE_QUALITY": {},
		"STATUS":           {},
		"TYPE":             {},
	}
	// allowedHistoryGetIssueResolutionHistoryEntityType is the set of allowed values for the corresponding option.
	allowedHistoryGetIssueResolutionHistoryEntityType = map[string]struct{}{
		"PORTFOLIO":      {},
		"PROJECT_BRANCH": {},
		"APPLICATION":    {},
	}
	// allowedHistoryGetIssueResolutionHistoryStatistic is the set of allowed values for the corresponding option.
	allowedHistoryGetIssueResolutionHistoryStatistic = map[string]struct{}{
		"MTTR":            {},
		"RECENT_MTTR":     {},
		"RESOLVED_ISSUES": {},
	}
	// allowedHistoryGetIssueResolutionHistorySliceBy is the set of allowed values for the corresponding option.
	allowedHistoryGetIssueResolutionHistorySliceBy = map[string]struct{}{
		"SEVERITY":         {},
		"TYPE_SEVERITY":    {},
		"SOFTWARE_QUALITY": {},
		"TYPE":             {},
	}
	// allowedHistoryGetProjectIssueResolutionStatistic is the set of allowed values for the corresponding option.
	allowedHistoryGetProjectIssueResolutionStatistic = map[string]struct{}{
		"MTTR":            {},
		"RECENT_MTTR":     {},
		"RESOLVED_ISSUES": {},
	}
	// allowedHistoryGetProjectIssueResolutionEntityType is the set of allowed values for the corresponding option.
	allowedHistoryGetProjectIssueResolutionEntityType = map[string]struct{}{
		"PORTFOLIO":   {},
		"APPLICATION": {},
	}
	// allowedHistoryGetProjectScaResolutionStatistic is the set of allowed values for the corresponding option.
	allowedHistoryGetProjectScaResolutionStatistic = map[string]struct{}{
		"SCA_MTTR": {},
	}
	// allowedHistoryGetProjectScaResolutionEntityType is the set of allowed values for the corresponding option.
	allowedHistoryGetProjectScaResolutionEntityType = map[string]struct{}{
		"PORTFOLIO":   {},
		"APPLICATION": {},
	}
	// allowedHistoryGetScaResolutionHistoryEntityType is the set of allowed values for the corresponding option.
	allowedHistoryGetScaResolutionHistoryEntityType = map[string]struct{}{
		"PORTFOLIO":      {},
		"PROJECT_BRANCH": {},
		"APPLICATION":    {},
	}
	// allowedHistoryGetScaResolutionHistoryStatistic is the set of allowed values for the corresponding option.
	allowedHistoryGetScaResolutionHistoryStatistic = map[string]struct{}{
		"SCA_MTTR": {},
	}
	// allowedHistoryGetScaResolutionHistorySliceBy is the set of allowed values for the corresponding option.
	allowedHistoryGetScaResolutionHistorySliceBy = map[string]struct{}{
		"SEVERITY": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// HistoryGetIssueDensityHistoryOptions contains parameters for the GetIssueDensityHistory method.
type HistoryGetIssueDensityHistoryOptions struct {
	// EntityId the ID of the entity to retrieve issue density history for. This field is required.
	EntityId string `json:"entityId"`
	// EntityType the entity's type. This field is required. Allowed values: PORTFOLIO,
	// PROJECT_BRANCH, APPLICATION.
	EntityType string `json:"entityType"`
	// StartDate the date and time to begin returning issue density history in ISO 8601 format
	// (inclusive). The start date may not be more than 1 year in the past. This field is required.
	StartDate string `json:"startDate"`
	// EndDate the date and time to stop returning issue density history in ISO 8601 format
	// (inclusive). If not specified, the current date and time at UTC offset 0 will be used.
	EndDate string `json:"endDate,omitempty"`
	// Impacts limit results to issues with specific impacts (comma-separated). An impact is a
	// combination of software quality and severity.
	Impacts []string `json:"impacts,omitempty"`
	// IssueTypes limit results to issues with specific types (comma-separated). Allowed values are
	// the SonarQube issue types.
	IssueTypes []string `json:"issueTypes,omitempty"`
	// RuleKeys comma separated list of ruleKeys.
	RuleKeys []string `json:"ruleKeys,omitempty"`
	// Severities limit results to issues with specific impact severities (comma-separated).
	// Allowed values: BLOCKER, HIGH, INFO, LOW, MEDIUM.
	Severities []string `json:"severities,omitempty"`
	// TypeSeverities limit results to issues with specific stored issue severities (comma-
	// separated). Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER.
	TypeSeverities []string `json:"typeSeverities,omitempty"`
	// SliceBy optionally group issue counts by a specific dimension. If not provided, issue counts
	// will contain a single distribution with the count of all issues. Allowed values: RULE_KEY,
	// SEVERITY, TYPE_SEVERITY, SOFTWARE_QUALITY, STATUS, TYPE.
	SliceBy string `json:"sliceBy,omitempty"`
	// Statuses limit results to issues with specific statuses (comma-separated). Allowed values:
	// ACCEPTED, CONFIRMED, FALSE_POSITIVE, FIXED, OPEN.
	Statuses []string `json:"statuses,omitempty"`
}

// HistoryIssueDensityDistribution represents the IssueDensityDistribution object of the SonarQube
// V2 API.
type HistoryIssueDensityDistribution struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Value is the value.
	Value float64 `json:"value,omitempty"`
}

// HistoryIssueDensityHistory represents the IssueDensityHistory object of the SonarQube V2 API.
type HistoryIssueDensityHistory struct {
	// Date and time in ISO 8601 format, including timezone offset.
	Date string `json:"date,omitempty"`
	// Distribution is the distribution.
	Distribution []HistoryIssueDensityDistribution `json:"distribution,omitempty"`
}

// HistoryIssueDensityHistoryResponse represents the IssueDensityHistoryResponse object of the
// SonarQube V2 API.
type HistoryIssueDensityHistoryResponse struct {
	// IssueDensityHistory is the issue density history.
	IssueDensityHistory []HistoryIssueDensityHistory `json:"issueDensityHistory,omitempty"`
}

// HistoryGetIssueResolutionHistoryOptions contains parameters for the GetIssueResolutionHistory method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryGetIssueResolutionHistoryOptions struct {
	// EntityId the ID of the entity to retrieve issue resolution history for. This field is
	// required.
	EntityId string `json:"entityId"`
	// EntityType the entity's type. This field is required. Allowed values: PORTFOLIO,
	// PROJECT_BRANCH, APPLICATION.
	EntityType string `json:"entityType"`
	// Statistic the issue resolution statistic to return. This field is required. Allowed values:
	// MTTR, RECENT_MTTR, RESOLVED_ISSUES.
	Statistic string `json:"statistic"`
	// StartDate the inclusive start of the history range in ISO 8601 format. This field is
	// required.
	StartDate string `json:"startDate"`
	// EndDate the inclusive end of the history range. Defaults to the current date.
	EndDate string `json:"endDate,omitempty"`
	// Impacts limit results to issues with specific impacts (comma-separated). An impact is a
	// combination of software quality and severity.
	Impacts []string `json:"impacts,omitempty"`
	// IssueTypes limit results to issues with specific types (comma-separated). Allowed values are
	// the SonarQube issue types.
	IssueTypes []string `json:"issueTypes,omitempty"`
	// Severities limit results to issues with specific impact severities (comma-separated).
	// Allowed values: BLOCKER, HIGH, INFO, LOW, MEDIUM.
	Severities []string `json:"severities,omitempty"`
	// TypeSeverities limit results to issues with specific stored issue severities (comma-
	// separated). Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER.
	TypeSeverities []string `json:"typeSeverities,omitempty"`
	// SliceBy optionally group issue resolution metrics by a dimension. Allowed values: SEVERITY,
	// TYPE_SEVERITY, SOFTWARE_QUALITY, TYPE.
	SliceBy string `json:"sliceBy,omitempty"`
}

// HistoryIssueResolutionDistribution represents the IssueResolutionDistribution object of the
// SonarQube V2 API.
type HistoryIssueResolutionDistribution struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Value is the value.
	Value int32 `json:"value,omitempty"`
}

// HistoryIssueResolutionHistory represents the IssueResolutionHistory object of the SonarQube V2
// API.
type HistoryIssueResolutionHistory struct {
	// Date and time in ISO 8601 format, including timezone offset.
	Date string `json:"date,omitempty"`
	// Distribution is the distribution.
	Distribution []HistoryIssueResolutionDistribution `json:"distribution,omitempty"`
}

// HistoryIssueResolutionHistoryResponse represents the IssueResolutionHistoryResponse object of the
// SonarQube V2 API.
type HistoryIssueResolutionHistoryResponse struct {
	// Statistic is the statistic. Allowed values: MTTR, RECENT_MTTR, RESOLVED_ISSUES.
	Statistic string `json:"statistic,omitempty"`
	// IssueResolutionHistory is the issue resolution history.
	IssueResolutionHistory []HistoryIssueResolutionHistory `json:"issueResolutionHistory,omitempty"`
}

// HistoryGetProjectIssueResolutionOptions contains parameters for the GetProjectIssueResolution method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryGetProjectIssueResolutionOptions struct {
	// Statistic the project issue resolution statistic to return. This field is required. Allowed
	// values: MTTR, RECENT_MTTR, RESOLVED_ISSUES.
	Statistic string `json:"statistic"`
	// PortfolioId legacy portfolio ID. Either portfolioId or both entityType and entityId must be
	// provided, but not both forms.
	PortfolioId string `json:"portfolioId,omitempty"`
	// EntityType type of the selection whose projects are returned as a flattened list. Must be
	// provided together with entityId and cannot be combined with portfolioId. Allowed values:
	// PORTFOLIO, APPLICATION.
	EntityType string `json:"entityType,omitempty"`
	// EntityId portfolio ID or application branch ID, according to entityType. Must be provided
	// together with entityType and cannot be combined with portfolioId.
	EntityId string `json:"entityId,omitempty"`
	// Severities limit results to issues with specific impact severities (comma-separated).
	// Allowed values: BLOCKER, HIGH, INFO, LOW, MEDIUM.
	Severities []string `json:"severities,omitempty"`
	// TypeSeverities limit results to issues with specific stored issue severities (comma-
	// separated). Allowed values: INFO, MINOR, MAJOR, CRITICAL, BLOCKER.
	TypeSeverities []string `json:"typeSeverities,omitempty"`
	// IssueTypes limit results to issues with specific types (comma-separated). Allowed values are
	// the SonarQube issue types.
	IssueTypes []string `json:"issueTypes,omitempty"`
	// Impacts limit results to issues with specific impacts (comma-separated). An impact is a
	// combination of software quality and severity.
	Impacts []string `json:"impacts,omitempty"`
	// NameContains search query for project names, case insensitive.
	NameContains string `json:"nameContains,omitempty"`
	// TrendSince date used as the trend comparison point.
	TrendSince string `json:"trendSince,omitempty"`
	// PageIndex index of the page to fetch.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize number of projects to return. A value of 0 returns pagination only.
	// Must be between 0 and 5000. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// HistoryProjectIssueResolution represents the ProjectIssueResolution object of the SonarQube V2
// API.
type HistoryProjectIssueResolution struct {
	// BranchId is the branch id.
	BranchId string `json:"branchId,omitempty"`
	// ProjectName is the project name.
	ProjectName string `json:"projectName,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// BranchName is the branch name.
	BranchName string `json:"branchName,omitempty"`
	// Value is the value.
	Value int64 `json:"value,omitempty"`
	// TrendPercentage percentage change in value relative to trendSince.
	TrendPercentage float64 `json:"trendPercentage,omitempty"`
}

// HistoryProjectIssueResolutionResponse represents the ProjectIssueResolutionResponse object of the
// SonarQube V2 API.
type HistoryProjectIssueResolutionResponse struct {
	// Statistic is the statistic. Allowed values: MTTR, RECENT_MTTR, RESOLVED_ISSUES.
	Statistic string `json:"statistic,omitempty"`
	// ProjectIssueResolution is the project issue resolution.
	ProjectIssueResolution []HistoryProjectIssueResolution `json:"projectIssueResolution,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// HistoryGetProjectScaResolutionOptions contains parameters for the GetProjectScaResolution method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryGetProjectScaResolutionOptions struct {
	// Statistic the SCA resolution statistic to return. This field is required. Allowed values:
	// SCA_MTTR.
	Statistic string `json:"statistic"`
	// PortfolioId legacy portfolio ID. Either portfolioId or both entityType and entityId must be
	// provided, but not both forms.
	PortfolioId string `json:"portfolioId,omitempty"`
	// EntityType type of the selection whose projects are returned as a flattened list. Must be
	// provided together with entityId and cannot be combined with portfolioId. Allowed values:
	// PORTFOLIO, APPLICATION.
	EntityType string `json:"entityType,omitempty"`
	// EntityId portfolio ID or application branch ID, according to entityType. Must be provided
	// together with entityType and cannot be combined with portfolioId.
	EntityId string `json:"entityId,omitempty"`
	// Severities limit results to issues with specific impact severities (comma-separated).
	// Allowed values: BLOCKER, HIGH, INFO, LOW, MEDIUM.
	Severities []string `json:"severities,omitempty"`
	// NameContains search query for project names, case insensitive.
	NameContains string `json:"nameContains,omitempty"`
	// TrendSince the date from which to compare project SCA resolution values. Past values are
	// unbounded. The value must be before the current date in UTC; values on or after the current
	// date are rejected. If omitted, trend is not returned in the response.
	TrendSince string `json:"trendSince,omitempty"`
	// PageIndex index of the page to fetch.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize size of the page to fetch.
	// Must be between 0 and 5000. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// HistoryProjectScaResolution represents the ProjectScaResolution object of the SonarQube V2 API.
type HistoryProjectScaResolution struct {
	// BranchId is the branch id.
	BranchId string `json:"branchId,omitempty"`
	// ProjectName is the project name.
	ProjectName string `json:"projectName,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// BranchName is the branch name.
	BranchName string `json:"branchName,omitempty"`
	// Value is the value.
	Value int64 `json:"value,omitempty"`
	// Trend percentage change in value relative to the value at trendSince.
	Trend float64 `json:"trend,omitempty"`
}

// HistoryProjectScaResolutionResponse represents the ProjectScaResolutionResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryProjectScaResolutionResponse struct {
	// HiddenProjectCount the number of projects excluded from the response because the user does
	// not have permission to view those projects.
	HiddenProjectCount int32 `json:"hiddenProjectCount,omitempty"`
	// Statistic is the statistic. Allowed values: SCA_MTTR.
	Statistic string `json:"statistic,omitempty"`
	// ProjectScaResolution is the project sca resolution.
	ProjectScaResolution []HistoryProjectScaResolution `json:"projectScaResolution,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// HistoryGetScaResolutionHistoryOptions contains parameters for the GetScaResolutionHistory method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type HistoryGetScaResolutionHistoryOptions struct {
	// EntityId the ID of the entity to retrieve SCA resolution history for. This field is
	// required.
	EntityId string `json:"entityId"`
	// EntityType the entity's type. This field is required. Allowed values: PORTFOLIO,
	// PROJECT_BRANCH, APPLICATION.
	EntityType string `json:"entityType"`
	// Statistic the SCA resolution statistic to return. This field is required. Allowed values:
	// SCA_MTTR.
	Statistic string `json:"statistic"`
	// StartDate inclusive start of the date range. Past values are unbounded. This field is
	// required.
	StartDate string `json:"startDate"`
	// EndDate inclusive end of the date range. Defaults to the current date and time in UTC.
	// Future values are clamped to the current time.
	EndDate string `json:"endDate,omitempty"`
	// Severities limit results to issues with specific impact severities (comma-separated).
	// Allowed values: BLOCKER, HIGH, INFO, LOW, MEDIUM.
	Severities []string `json:"severities,omitempty"`
	// SliceBy optionally group SCA resolution metrics by severity. Allowed values: SEVERITY.
	SliceBy string `json:"sliceBy,omitempty"`
}

// HistoryScaResolutionDistribution represents the ScaResolutionDistribution object of the SonarQube
// V2 API.
type HistoryScaResolutionDistribution struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Value is the value.
	Value int32 `json:"value,omitempty"`
}

// HistoryScaResolutionHistory represents the ScaResolutionHistory object of the SonarQube V2 API.
type HistoryScaResolutionHistory struct {
	// Date and time in ISO 8601 format, including timezone offset.
	Date string `json:"date,omitempty"`
	// Distribution is the distribution.
	Distribution []HistoryScaResolutionDistribution `json:"distribution,omitempty"`
}

// HistoryScaResolutionHistoryResponse represents the ScaResolutionHistoryResponse object of the
// SonarQube V2 API.
type HistoryScaResolutionHistoryResponse struct {
	// Statistic is the statistic. Allowed values: SCA_MTTR.
	Statistic string `json:"statistic,omitempty"`
	// ScaResolutionHistory is the sca resolution history.
	ScaResolutionHistory []HistoryScaResolutionHistory `json:"scaResolutionHistory,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetIssueDensityHistoryOpt validates the options for the GetIssueDensityHistory method.
func (s *HistoryService) ValidateGetIssueDensityHistoryOpt(opt *HistoryGetIssueDensityHistoryOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.EntityId, "EntityId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.EntityType, "EntityType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EntityType, allowedHistoryGetIssueDensityHistoryEntityType, "EntityType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.StartDate, "StartDate")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.SliceBy, allowedHistoryGetIssueDensityHistorySliceBy, "SliceBy")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetIssueResolutionHistoryOpt validates the options for the GetIssueResolutionHistory method.
func (s *HistoryService) ValidateGetIssueResolutionHistoryOpt(opt *HistoryGetIssueResolutionHistoryOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.EntityId, "EntityId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.EntityType, "EntityType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EntityType, allowedHistoryGetIssueResolutionHistoryEntityType, "EntityType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Statistic, "Statistic")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Statistic, allowedHistoryGetIssueResolutionHistoryStatistic, "Statistic")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.StartDate, "StartDate")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.SliceBy, allowedHistoryGetIssueResolutionHistorySliceBy, "SliceBy")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetProjectIssueResolutionOpt validates the options for the GetProjectIssueResolution method.
func (s *HistoryService) ValidateGetProjectIssueResolutionOpt(opt *HistoryGetProjectIssueResolutionOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Statistic, "Statistic")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Statistic, allowedHistoryGetProjectIssueResolutionStatistic, "Statistic")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EntityType, allowedHistoryGetProjectIssueResolutionEntityType, "EntityType")
	if err != nil {
		return err
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	if opt.PageSize < 0 || opt.PageSize > 5000 {
		return NewValidationError("PageSize", "must be between 0 and 5000", ErrOutOfRange)
	}

	return nil
}

// ValidateGetProjectScaResolutionOpt validates the options for the GetProjectScaResolution method.
func (s *HistoryService) ValidateGetProjectScaResolutionOpt(opt *HistoryGetProjectScaResolutionOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Statistic, "Statistic")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Statistic, allowedHistoryGetProjectScaResolutionStatistic, "Statistic")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EntityType, allowedHistoryGetProjectScaResolutionEntityType, "EntityType")
	if err != nil {
		return err
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	if opt.PageSize < 0 || opt.PageSize > 5000 {
		return NewValidationError("PageSize", "must be between 0 and 5000", ErrOutOfRange)
	}

	return nil
}

// ValidateGetScaResolutionHistoryOpt validates the options for the GetScaResolutionHistory method.
func (s *HistoryService) ValidateGetScaResolutionHistoryOpt(opt *HistoryGetScaResolutionHistoryOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.EntityId, "EntityId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.EntityType, "EntityType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.EntityType, allowedHistoryGetScaResolutionHistoryEntityType, "EntityType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Statistic, "Statistic")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Statistic, allowedHistoryGetScaResolutionHistoryStatistic, "Statistic")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.StartDate, "StartDate")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.SliceBy, allowedHistoryGetScaResolutionHistorySliceBy, "SliceBy")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetIssueDensityHistory provides issue density history data for an entity. Returns the history of
// issue density (issues per 1000 NCLOC) for an entity. The maximum resolution for history data is 1
// day. All history is recorded in UTC timezone. Issue density values carry-forward for days with no
// changes. Density is calculated as (issue count / NCLOC) * 1000.
//
// API endpoint: GET /api/v2/history/issue-density-history.
// Enterprise Edition only.
func (s *HistoryService) GetIssueDensityHistory(ctx context.Context, opt *HistoryGetIssueDensityHistoryOptions) (*HistoryIssueDensityHistoryResponse, *http.Response, error) {
	err := s.ValidateGetIssueDensityHistoryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/issue-density-history", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryIssueDensityHistoryResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetIssueResolutionHistory provides issue resolution history data for an entity. Returns issue
// resolution metrics for an entity. All history is recorded at UTC midnight. MTTR and RECENT_MTTR
// use a trailing 30-day window ending on each returned day. SQS deliberately permits dates older
// than one year under the approved history date-range policy.
//
// API endpoint: GET /api/v2/history/issue-resolution-history.
// Enterprise Edition only.
func (s *HistoryService) GetIssueResolutionHistory(ctx context.Context, opt *HistoryGetIssueResolutionHistoryOptions) (*HistoryIssueResolutionHistoryResponse, *http.Response, error) {
	err := s.ValidateGetIssueResolutionHistoryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/issue-resolution-history", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryIssueResolutionHistoryResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectIssueResolution gets project issue resolution. Get issue resolution values for the
// flattened list of projects in a portfolio or application. Identify the selection with either the
// legacy portfolioId parameter, or with both entityType and entityId. The two selector forms are
// mutually exclusive.
//
// API endpoint: GET /api/v2/history/project-issue-resolution.
// Enterprise Edition only.
func (s *HistoryService) GetProjectIssueResolution(ctx context.Context, opt *HistoryGetProjectIssueResolutionOptions) (*HistoryProjectIssueResolutionResponse, *http.Response, error) {
	err := s.ValidateGetProjectIssueResolutionOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/project-issue-resolution", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryProjectIssueResolutionResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectScaResolution gets project SCA resolution. Get SCA resolution values for the flattened
// list of projects in a portfolio or application. Identify the selection with either the legacy
// portfolioId parameter, or with both entityType and entityId. The two selector forms are mutually
// exclusive.
//
// API endpoint: GET /api/v2/history/project-sca-resolution.
// Enterprise Edition only.
func (s *HistoryService) GetProjectScaResolution(ctx context.Context, opt *HistoryGetProjectScaResolutionOptions) (*HistoryProjectScaResolutionResponse, *http.Response, error) {
	err := s.ValidateGetProjectScaResolutionOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/project-sca-resolution", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryProjectScaResolutionResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetScaResolutionHistory provides SCA resolution history data for an entity. Returns the history
// of SCA resolution metrics for an entity. All history is recorded in UTC. SCA MTTR values are
// calculated using a trailing 30-day window ending on each returned day.
//
// API endpoint: GET /api/v2/history/sca-resolution-history.
// Enterprise Edition only.
func (s *HistoryService) GetScaResolutionHistory(ctx context.Context, opt *HistoryGetScaResolutionHistoryOptions) (*HistoryScaResolutionHistoryResponse, *http.Response, error) {
	err := s.ValidateGetScaResolutionHistoryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "history/sca-resolution-history", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(HistoryScaResolutionHistoryResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
