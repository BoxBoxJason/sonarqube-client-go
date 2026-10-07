package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DashboardsService handles communication with the custom and built-in dashboards related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type DashboardsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals,goconst // constant sets of allowed values
var (
	// allowedDashboardsListResourceType is the set of allowed values for the corresponding option.
	allowedDashboardsListResourceType = map[string]struct{}{
		"project":   {},
		"portfolio": {},
	}
	// allowedDashboardsListBuiltInResourceType is the set of allowed values for the corresponding option.
	allowedDashboardsListBuiltInResourceType = map[string]struct{}{
		"project":   {},
		"portfolio": {},
	}
	// allowedDashboardsCreateResourceType is the set of allowed values for the corresponding option.
	allowedDashboardsCreateResourceType = map[string]struct{}{
		"project":   {},
		"portfolio": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// DashboardsListOptions contains parameters for the List method.
type DashboardsListOptions struct {
	// ResourceId is the resource id. This field is required.
	ResourceId string `json:"resourceId"`
	// ResourceType is the resource type. This field is required. Allowed values: project,
	// portfolio.
	ResourceType string `json:"resourceType"`
	// Q is the q.
	Q string `json:"q,omitempty"`
	// PageIndex is the page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the page size.
	// Must be between 0 and 5000. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// DashboardsDashboardItem represents the DashboardItem object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DashboardsDashboardItem struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Description is the description.
	Description string `json:"description,omitempty"`
	// CreatedAt is the created at.
	CreatedAt int64 `json:"createdAt,omitempty"`
	// UpdatedAt is the updated at.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// CreatedById is the created by id.
	CreatedById string `json:"createdById,omitempty"`
	// ResourceType is the resource type. Allowed values: project, portfolio.
	ResourceType string `json:"resourceType,omitempty"`
	// ResourceId is the resource id.
	ResourceId string `json:"resourceId,omitempty"`
	// UpdatedById is the updated by id.
	UpdatedById string `json:"updatedById,omitempty"`
}

// DashboardsDashboardsResponse represents the DashboardsResponse object of the SonarQube V2 API.
type DashboardsDashboardsResponse struct {
	// Dashboards is the dashboards.
	Dashboards []DashboardsDashboardItem `json:"dashboards,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// DashboardsListBuiltInOptions contains parameters for the ListBuiltIn method.
type DashboardsListBuiltInOptions struct {
	// Q is the q.
	Q string `json:"q,omitempty"`
	// ResourceType is the resource type. Allowed values: project, portfolio.
	ResourceType string `json:"resourceType,omitempty"`
	// PageIndex is the page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize is the page size.
	// Must be between 0 and 5000. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// DashboardsBuiltInDashboardItem represents the BuiltInDashboardItem object of the SonarQube V2
// API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DashboardsBuiltInDashboardItem struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Description is the description.
	Description string `json:"description,omitempty"`
	// UpdatedAt is the updated at.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// ResourceType is the resource type. Allowed values: project, portfolio.
	ResourceType string `json:"resourceType,omitempty"`
}

// DashboardsBuiltInDashboardsResponse represents the BuiltInDashboardsResponse object of the
// SonarQube V2 API.
type DashboardsBuiltInDashboardsResponse struct {
	// Dashboards is the dashboards.
	Dashboards []DashboardsBuiltInDashboardItem `json:"dashboards,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// DashboardsBuiltInDashboardResponse represents the BuiltInDashboardResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DashboardsBuiltInDashboardResponse struct {
	// Key is the key.
	Key string `json:"key,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Layout is the layout.
	Layout string `json:"layout,omitempty"`
	// Description is the description.
	Description string `json:"description,omitempty"`
	// UpdatedAt is the updated at.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// ResourceType is the resource type. Allowed values: project, portfolio.
	ResourceType string `json:"resourceType,omitempty"`
}

// DashboardsDashboardResponse represents the DashboardResponse object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DashboardsDashboardResponse struct {
	// Id is the id.
	Id string `json:"id,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Layout is the layout.
	Layout string `json:"layout,omitempty"`
	// Description is the description.
	Description string `json:"description,omitempty"`
	// CreatedAt is the created at.
	CreatedAt int64 `json:"createdAt,omitempty"`
	// UpdatedAt is the updated at.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// CreatedById is the created by id.
	CreatedById string `json:"createdById,omitempty"`
	// ResourceType is the resource type. Allowed values: project, portfolio.
	ResourceType string `json:"resourceType,omitempty"`
	// ResourceId is the resource id.
	ResourceId string `json:"resourceId,omitempty"`
	// UpdatedById is the updated by id.
	UpdatedById string `json:"updatedById,omitempty"`
}

// DashboardsUpdateOptions contains the request body for the Update method.
type DashboardsUpdateOptions struct {
	// Name is the name.
	Name *string `json:"name,omitempty"`
	// Layout is the layout.
	Layout *string `json:"layout,omitempty"`
	// Description is the description.
	Description *string `json:"description,omitempty"`
}

// DashboardsCreateOptions contains the request body for the Create method.
type DashboardsCreateOptions struct {
	// Name is the name. This field is required.
	Name string `json:"name"`
	// Layout is the layout. This field is required.
	Layout string `json:"layout"`
	// ResourceType is the resource type. This field is required. Allowed values: project,
	// portfolio.
	ResourceType string `json:"resourceType"`
	// ResourceId is the resource id. This field is required.
	ResourceId string `json:"resourceId"`
	// Description is the description.
	Description string `json:"description,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListOpt validates the options for the List method.
func (s *DashboardsService) ValidateListOpt(opt *DashboardsListOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ResourceId, "ResourceId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ResourceType, "ResourceType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.ResourceType, allowedDashboardsListResourceType, "ResourceType")
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

// ValidateListBuiltInOpt validates the options for the ListBuiltIn method.
func (s *DashboardsService) ValidateListBuiltInOpt(opt *DashboardsListBuiltInOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.ResourceType, allowedDashboardsListBuiltInResourceType, "ResourceType")
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

// ValidateCreateOpt validates the options for the Create method.
func (s *DashboardsService) ValidateCreateOpt(opt *DashboardsCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Layout, "Layout")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ResourceType, "ResourceType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.ResourceType, allowedDashboardsCreateResourceType, "ResourceType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ResourceId, "ResourceId")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Delete deletes a custom dashboard.
//
// API endpoint: DELETE /api/v2/dashboards/{id}.
// Enterprise Edition only.
func (s *DashboardsService) Delete(ctx context.Context, dashboardID string) (*http.Response, error) {
	err := ValidateRequired(dashboardID, "dashboardID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "dashboards/"+url.PathEscape(dashboardID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// List lists dashboards for a resource.
//
// API endpoint: GET /api/v2/dashboards.
// Enterprise Edition only.
func (s *DashboardsService) List(ctx context.Context, opt *DashboardsListOptions) (*DashboardsDashboardsResponse, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "dashboards", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DashboardsDashboardsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListBuiltIn lists built-in dashboards.
//
// API endpoint: GET /api/v2/dashboards/built-ins.
// Enterprise Edition only.
func (s *DashboardsService) ListBuiltIn(ctx context.Context, opt *DashboardsListBuiltInOptions) (*DashboardsBuiltInDashboardsResponse, *http.Response, error) {
	err := s.ValidateListBuiltInOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "dashboards/built-ins", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DashboardsBuiltInDashboardsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetBuiltIn gets a built-in dashboard.
//
// API endpoint: GET /api/v2/dashboards/built-ins/{key}.
// Enterprise Edition only.
func (s *DashboardsService) GetBuiltIn(ctx context.Context, key string) (*DashboardsBuiltInDashboardResponse, *http.Response, error) {
	err := ValidateRequired(key, "key")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "dashboards/built-ins/"+url.PathEscape(key), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DashboardsBuiltInDashboardResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Get gets a custom dashboard.
//
// API endpoint: GET /api/v2/dashboards/{id}.
// Enterprise Edition only.
func (s *DashboardsService) Get(ctx context.Context, dashboardID string) (*DashboardsDashboardResponse, *http.Response, error) {
	err := ValidateRequired(dashboardID, "dashboardID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "dashboards/"+url.PathEscape(dashboardID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DashboardsDashboardResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Update updates a custom dashboard.
//
// API endpoint: PATCH /api/v2/dashboards/{id}.
// Enterprise Edition only.
func (s *DashboardsService) Update(ctx context.Context, dashboardID string, opt *DashboardsUpdateOptions) (*DashboardsDashboardResponse, *http.Response, error) {
	err := ValidateRequired(dashboardID, "dashboardID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "dashboards/"+url.PathEscape(dashboardID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DashboardsDashboardResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Create creates a custom dashboard.
//
// API endpoint: POST /api/v2/dashboards.
// Enterprise Edition only.
func (s *DashboardsService) Create(ctx context.Context, opt *DashboardsCreateOptions) (*DashboardsDashboardResponse, *http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "dashboards", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(DashboardsDashboardResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
