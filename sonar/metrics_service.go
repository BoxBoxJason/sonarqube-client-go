// Get information on automatic metrics
package sonargo

import "net/http"

type MetricsService struct {
	client *Client
}

type MetricsSearchObject struct {
	Metrics []MetricsSearchObject_sub1 `json:"metrics,omitempty"`
	Paging  MetricsSearchObject_sub2   `json:"paging,omitempty"`
}

type MetricsSearchObject_sub1 struct {
	Custom      bool   `json:"custom,omitempty"`
	Description string `json:"description,omitempty"`
	Direction   int64  `json:"direction,omitempty"`
	Domain      string `json:"domain,omitempty"`
	Hidden      bool   `json:"hidden,omitempty"`
	ID          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Qualitative bool   `json:"qualitative,omitempty"`
	Type        string `json:"type,omitempty"`
}

type MetricsSearchObject_sub2 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type MetricsTypesObject struct {
	Types []string `json:"types,omitempty"`
}

type MetricsSearchOption struct {
	P  string `url:"p,omitempty"`  // Description:"1-based page number",ExampleValue:"42"
	Ps string `url:"ps,omitempty"` // Description:"Page size. Must be greater than 0 and less or equal than 500",ExampleValue:"20"
}

// Search Search for metrics
func (s *MetricsService) Search(opt *MetricsSearchOption) (v *MetricsSearchObject, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "metrics/search", opt)
	if err != nil {
		return
	}
	v = new(MetricsSearchObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// Types List all available metric types.
func (s *MetricsService) Types() (v *MetricsTypesObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "metrics/types", nil)
	if err != nil {
		return
	}
	v = new(MetricsTypesObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
