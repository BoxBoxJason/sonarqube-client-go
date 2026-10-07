package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// ScaListReachabilityPathsOptions contains parameters for the ListReachabilityPaths method.
type ScaListReachabilityPathsOptions struct {
	// LanguageKey is the language key. This field is required.
	LanguageKey string `json:"languageKey"`
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// PageSize number of results per page. A value of 0 will only return the pagination
	// information.
	// Must be between 0 and 10000. Default: 10000.
	PageSize int32 `json:"pageSize,omitempty"`
	// PageIndex 1-based page index.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
}

// ScaReachabilityPathsRestResponse represents the ReachabilityPathsRestResponse object of the
// SonarQube V2 API.
type ScaReachabilityPathsRestResponse struct {
	// Paths is the paths.
	Paths []string `json:"paths,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListReachabilityPathsOpt validates the options for the ListReachabilityPaths method.
func (s *ScaService) ValidateListReachabilityPathsOpt(opt *ScaListReachabilityPathsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.LanguageKey, "LanguageKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	if opt.PageSize < 0 || opt.PageSize > 10000 {
		return NewValidationError("PageSize", "must be between 0 and 10000", ErrOutOfRange)
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ListReachabilityPaths gets reachability paths. Returns a page of distinct reachability paths.
//
// API endpoint: GET /api/v2/sca/reachability/list-paths.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ScaService) ListReachabilityPaths(ctx context.Context, opt *ScaListReachabilityPathsOptions) (*ScaReachabilityPathsRestResponse, *http.Response, error) {
	err := s.ValidateListReachabilityPathsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "sca/reachability/list-paths", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ScaReachabilityPathsRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
