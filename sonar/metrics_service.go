package sonargo

import "net/http"

// MetricsService handles communication with the metrics related methods
// of the SonarQube API.
// This service provides information about metrics.
type MetricsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// MetricsSearch represents the response from searching metrics.
type MetricsSearch struct {
	// Metrics is the list of metrics found.
	Metrics []Metric `json:"metrics,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// Metric represents a SonarQube metric.
//
//nolint:govet // Field alignment is less important than logical grouping.
type Metric struct {
	// Custom indicates whether this is a custom metric.
	Custom bool `json:"custom,omitempty"`
	// Description is the metric description.
	Description string `json:"description,omitempty"`
	// Direction indicates the metric direction (-1: lower is better, 0: neutral, 1: higher is better).
	Direction int64 `json:"direction,omitempty"`
	// Domain is the metric domain.
	Domain string `json:"domain,omitempty"`
	// Hidden indicates whether the metric is hidden.
	Hidden bool `json:"hidden,omitempty"`
	// ID is the metric identifier.
	ID string `json:"id,omitempty"`
	// Key is the metric key.
	Key string `json:"key,omitempty"`
	// Name is the metric display name.
	Name string `json:"name,omitempty"`
	// Qualitative indicates whether the metric is qualitative.
	Qualitative bool `json:"qualitative,omitempty"`
	// Type is the metric value type.
	Type string `json:"type,omitempty"`
}

// MetricsTypes represents the response from listing metric types.
type MetricsTypes struct {
	// Types is the list of available metric types.
	Types []string `json:"types,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// MetricsSearchOption contains parameters for the Search method.
type MetricsSearchOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSearchOpt validates the options for the Search method.
func (s *MetricsService) ValidateSearchOpt(opt *MetricsSearchOption) error {
	if opt == nil {
		return nil
	}

	return ValidatePagination(opt.Page, opt.PageSize)
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Search searches for metrics.
//
// API endpoint: GET /api/metrics/search.
// Since: 2.6.
func (s *MetricsService) Search(opt *MetricsSearchOption) (*MetricsSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "metrics/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(MetricsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Types lists all available metric types.
//
// API endpoint: GET /api/metrics/types.
// Since: 2.6.
func (s *MetricsService) Types() (*MetricsTypes, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "metrics/types", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(MetricsTypes)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
