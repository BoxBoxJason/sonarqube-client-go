package sonargo

import (
	"net/http"
)

// MeasuresService handles communication with the Measures related methods of the SonarQube API.
// Get components or children with specified measures.
//
// Since: 5.4.
type MeasuresService struct {
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedMeasuresMetricSortFilter is the set of allowed values for metric sort filtering.
	allowedMeasuresMetricSortFilter = map[string]struct{}{
		"all":              {},
		"withMeasuresOnly": {},
	}

	// allowedMeasuresStrategy is the set of allowed values for strategy.
	allowedMeasuresStrategy = map[string]struct{}{
		"all":      {},
		"children": {},
		"leaves":   {},
	}
)

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// MeasuresComponent represents the response from getting component measures.
//
//nolint:govet // MeasureComponent embedded struct alignment is complex
type MeasuresComponent struct {
	// Component is the component with its measures.
	Component MeasureComponent `json:"component,omitzero"`
	// Metrics is the list of metrics definitions.
	Metrics []MeasureMetric `json:"metrics,omitempty"`
	// Period is the period definition.
	Period *MeasurePeriod `json:"period,omitempty"`
}

// MeasuresComponentTree represents the response from getting component tree measures.
//
//nolint:govet // Multiple embedded structs make optimal alignment impractical
type MeasuresComponentTree struct {
	// BaseComponent is the base component from which the tree is computed.
	BaseComponent MeasureComponent `json:"baseComponent,omitzero"`
	// Paging is the pagination info.
	Paging Paging `json:"paging,omitzero"`
	// Components is the list of components with their measures.
	Components []MeasureComponent `json:"components,omitempty"`
	// Metrics is the list of metrics definitions.
	Metrics []MeasureMetric `json:"metrics,omitempty"`
	// Period is the period definition.
	Period *MeasurePeriod `json:"period,omitempty"`
}

// MeasuresSearch represents the response from searching measures.
type MeasuresSearch struct {
	// Measures is the list of measures.
	Measures []MeasureSearchResult `json:"measures,omitempty"`
}

// MeasureSearchResult represents a measure search result.
type MeasureSearchResult struct {
	// Component is the component key.
	Component string `json:"component,omitempty"`
	// Metric is the metric key.
	Metric string `json:"metric,omitempty"`
	// Value is the measure value.
	Value string `json:"value,omitempty"`
	// BestValue indicates if this is the best possible value.
	BestValue bool `json:"bestValue,omitempty"`
}

// MeasuresSearchHistory represents the response from getting measures history.
type MeasuresSearchHistory struct {
	// Measures is the list of measures with their history.
	Measures []MeasureHistory `json:"measures,omitempty"`
	// Paging is the pagination info.
	Paging Paging `json:"paging,omitzero"`
}

// MeasureHistory represents a measure with its history.
type MeasureHistory struct {
	// Metric is the metric key.
	Metric string `json:"metric,omitempty"`
	// History is the list of historical values.
	History []MeasureHistoryValue `json:"history,omitempty"`
}

// MeasureHistoryValue represents a historical measure value.
type MeasureHistoryValue struct {
	// Date is the date of the measure.
	Date string `json:"date,omitempty"`
	// Value is the measure value at that date.
	Value string `json:"value,omitempty"`
}

// MeasureComponent represents a component with its measures.
type MeasureComponent struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Description is the component description.
	Description string `json:"description,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// Path is the path within the project.
	Path string `json:"path,omitempty"`
	// Language is the programming language.
	Language string `json:"language,omitempty"`
	// Measures is the list of measures for this component.
	Measures []Measure `json:"measures,omitempty"`
}

// Measure represents a single measure.
type Measure struct {
	// Period is the period measure value (deprecated).
	Period *MeasurePeriodValue `json:"period,omitempty"`
	// Metric is the metric key.
	Metric string `json:"metric,omitempty"`
	// Value is the measure value.
	Value string `json:"value,omitempty"`
	// BestValue indicates if this is the best possible value.
	BestValue bool `json:"bestValue,omitempty"`
}

// MeasurePeriodValue represents a measure value for a period.
type MeasurePeriodValue struct {
	// Value is the measure value for this period.
	Value string `json:"value,omitempty"`
	// BestValue indicates if this is the best possible value.
	BestValue bool `json:"bestValue,omitempty"`
}

// MeasureMetric represents a metric definition.
//
//nolint:govet // Field alignment is less important than logical grouping
type MeasureMetric struct {
	// Key is the unique metric key.
	Key string `json:"key,omitempty"`
	// Name is the display name.
	Name string `json:"name,omitempty"`
	// Description is the metric description.
	Description string `json:"description,omitempty"`
	// Domain is the metric domain.
	Domain string `json:"domain,omitempty"`
	// Type is the metric type.
	Type string `json:"type,omitempty"`
	// HigherValuesAreBetter indicates if higher values are better.
	HigherValuesAreBetter bool `json:"higherValuesAreBetter,omitempty"`
	// Qualitative indicates if the metric is qualitative.
	Qualitative bool `json:"qualitative,omitempty"`
	// Hidden indicates if the metric is hidden.
	Hidden bool `json:"hidden,omitempty"`
	// DecimalScale is the number of decimal places.
	DecimalScale int64 `json:"decimalScale,omitempty"`
	// BestValue is the best possible value.
	BestValue string `json:"bestValue,omitempty"`
	// WorstValue is the worst possible value.
	WorstValue string `json:"worstValue,omitempty"`
}

// MeasurePeriod represents a period definition.
type MeasurePeriod struct {
	// Mode is the period mode.
	Mode string `json:"mode,omitempty"`
	// Date is the period date.
	Date string `json:"date,omitempty"`
	// Parameter is the period parameter (e.g., branch name, days).
	Parameter string `json:"parameter,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// MeasuresComponentOption represents options for getting component measures.
//
//nolint:govet // Field alignment is less important than logical grouping
type MeasuresComponentOption struct {
	// Component is the component key (required).
	Component string `url:"component,omitempty"`
	// MetricKeys is the list of metric keys to fetch (required).
	MetricKeys []string `url:"metricKeys,omitempty,comma"`
	// AdditionalFields is the list of additional fields to return (optional).
	// Possible values: metrics, period, periods.
	AdditionalFields []string `url:"additionalFields,omitempty,comma"`
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// PullRequest is the pull request identifier (optional).
	PullRequest string `url:"pullRequest,omitempty"`
}

// MeasuresComponentTreeOption represents options for getting component tree measures.
//
//nolint:govet // Field alignment is less important than logical grouping
type MeasuresComponentTreeOption struct {
	PaginationArgs

	// Component is the base component key (required).
	Component string `url:"component,omitempty"`
	// MetricKeys is the list of metric keys to fetch (required).
	MetricKeys []string `url:"metricKeys,omitempty,comma"`
	// AdditionalFields is the list of additional fields to return (optional).
	// Possible values: metrics, period, periods.
	AdditionalFields []string `url:"additionalFields,omitempty,comma"`
	// Asc indicates ascending sort when true, descending when false (optional).
	Asc bool `url:"asc,omitempty"`
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// MetricPeriodSort is the period to use for metric sorting (optional).
	MetricPeriodSort int64 `url:"metricPeriodSort,omitempty"`
	// MetricSort is the metric key to use for sorting (optional).
	MetricSort string `url:"metricSort,omitempty"`
	// MetricSortFilter specifies which components to sort (optional).
	// Possible values: all, withMeasuresOnly.
	MetricSortFilter string `url:"metricSortFilter,omitempty"`
	// PullRequest is the pull request identifier (optional).
	PullRequest string `url:"pullRequest,omitempty"`
	// Qualifiers filters by component qualifiers (optional).
	// Possible values: BRC, DIR, FIL, TRK, UTS.
	Qualifiers []string `url:"qualifiers,omitempty,comma"`
	// Query is used to filter components by name (optional).
	Query string `url:"q,omitempty"`
	// Sort is the list of sort fields (optional).
	// Possible values: metric, metricPeriod, name, path, qualifier.
	Sort []string `url:"s,omitempty,comma"`
	// Strategy specifies how to traverse the tree (optional).
	// Possible values: all, children, leaves.
	Strategy string `url:"strategy,omitempty"`
}

// MeasuresSearchOption represents options for searching measures.
type MeasuresSearchOption struct {
	// MetricKeys is the list of metric keys to fetch (required).
	MetricKeys []string `url:"metricKeys,omitempty,comma"`
	// ProjectKeys is the list of project keys (required).
	ProjectKeys []string `url:"projectKeys,omitempty,comma"`
}

// MeasuresSearchHistoryOption represents options for getting measures history.
//
//nolint:govet // Field alignment is less important than logical grouping
type MeasuresSearchHistoryOption struct {
	PaginationArgs

	// Component is the component key (required).
	Component string `url:"component,omitempty"`
	// Metrics is the list of metric keys (required).
	Metrics []string `url:"metrics,omitempty,comma"`
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// From is the start date filter (optional).
	// Format: YYYY-MM-DD or YYYY-MM-DDTHH:mm:ssZ.
	From string `url:"from,omitempty"`
	// PullRequest is the pull request identifier (optional).
	PullRequest string `url:"pullRequest,omitempty"`
	// To is the end date filter (optional).
	// Format: YYYY-MM-DD or YYYY-MM-DDTHH:mm:ssZ.
	To string `url:"to,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Methods
// -----------------------------------------------------------------------------

// ValidateComponentOpt validates the options for Component.
func (s *MeasuresService) ValidateComponentOpt(opt *MeasuresComponentOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	if len(opt.MetricKeys) == 0 {
		return NewValidationError("MetricKeys", "at least one metric key is required", ErrMissingRequired)
	}

	return nil
}

// ValidateComponentTreeOpt validates the options for ComponentTree.
func (s *MeasuresService) ValidateComponentTreeOpt(opt *MeasuresComponentTreeOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	if len(opt.MetricKeys) == 0 {
		return NewValidationError("MetricKeys", "at least one metric key is required", ErrMissingRequired)
	}

	if opt.MetricSortFilter != "" {
		err = IsValueAuthorized(opt.MetricSortFilter, allowedMeasuresMetricSortFilter, "MetricSortFilter")
		if err != nil {
			return err
		}
	}

	if opt.Strategy != "" {
		err = IsValueAuthorized(opt.Strategy, allowedMeasuresStrategy, "Strategy")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchOpt validates the options for Search.
func (s *MeasuresService) ValidateSearchOpt(opt *MeasuresSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if len(opt.MetricKeys) == 0 {
		return NewValidationError("MetricKeys", "at least one metric key is required", ErrMissingRequired)
	}

	if len(opt.ProjectKeys) == 0 {
		return NewValidationError("ProjectKeys", "at least one project key is required", ErrMissingRequired)
	}

	return nil
}

// ValidateSearchHistoryOpt validates the options for SearchHistory.
func (s *MeasuresService) ValidateSearchHistoryOpt(opt *MeasuresSearchHistoryOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	if len(opt.Metrics) == 0 {
		return NewValidationError("Metrics", "at least one metric is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Component returns measures for a component.
// Requires the following permission: 'Browse' on the project of specified component.
//
// Since: 5.4.
func (s *MeasuresService) Component(opt *MeasuresComponentOption) (*MeasuresComponent, *http.Response, error) {
	err := s.ValidateComponentOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "measures/component", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(MeasuresComponent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ComponentTree returns measures for a component tree.
// Navigate through components based on the chosen strategy with specified measures.
// Requires the following permission: 'Browse' on the specified project.
// When limiting search with the q parameter, directories are not returned.
//
// Since: 5.4.
func (s *MeasuresService) ComponentTree(opt *MeasuresComponentTreeOption) (*MeasuresComponentTree, *http.Response, error) {
	err := s.ValidateComponentTreeOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "measures/component_tree", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(MeasuresComponentTree)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Search returns the measures for multiple projects.
// Requires 'Browse' permission on the projects.
// At most 100 projects can be provided.
//
// Since: 6.2.
func (s *MeasuresService) Search(opt *MeasuresSearchOption) (*MeasuresSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "measures/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(MeasuresSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchHistory returns the history of measures for a component.
// Pagination is available for the last 1000 analyses.
// Requires the following permission: 'Browse' on the specified component.
//
// Since: 6.3.
func (s *MeasuresService) SearchHistory(opt *MeasuresSearchHistoryOption) (*MeasuresSearchHistory, *http.Response, error) {
	err := s.ValidateSearchHistoryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "measures/search_history", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(MeasuresSearchHistory)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
