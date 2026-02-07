package sonar

import "net/http"

// FavoritesService handles communication with the favorites related methods
// of the SonarQube API.
// This service manages user favorites.
type FavoritesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// FavoritesSearch represents the response from searching favorites.
type FavoritesSearch struct {
	// Favorites is the list of favorited components.
	Favorites []Favorite `json:"favorites,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// Favorite represents a favorited component.
type Favorite struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Qualifier is the component qualifier (e.g., TRK for project).
	Qualifier string `json:"qualifier,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// FavoritesAddOption contains parameters for the Add method.
type FavoritesAddOption struct {
	// Component is the component key.
	// Only components with qualifiers TRK, VW, SVW, APP are supported.
	// This field is required.
	Component string `url:"component"`
}

// FavoritesRemoveOption contains parameters for the Remove method.
type FavoritesRemoveOption struct {
	// Component is the component key.
	// This field is required.
	Component string `url:"component"`
}

// FavoritesSearchOption contains parameters for the Search method.
type FavoritesSearchOption struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs `url:",inline"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateAddOpt validates the options for the Add method.
func (s *FavoritesService) ValidateAddOpt(opt *FavoritesAddOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	return nil
}

// ValidateRemoveOpt validates the options for the Remove method.
func (s *FavoritesService) ValidateRemoveOpt(opt *FavoritesRemoveOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
func (s *FavoritesService) ValidateSearchOpt(opt *FavoritesSearchOption) error {
	if opt == nil {
		return nil
	}

	return ValidatePagination(opt.Page, opt.PageSize)
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Add adds a component as favorite for the authenticated user.
// Only 100 components by qualifier can be added as favorite.
// Requires authentication and the 'Browse' permission on the component.
//
// API endpoint: POST /api/favorites/add.
// Since: 6.3.
func (s *FavoritesService) Add(opt *FavoritesAddOption) (*http.Response, error) {
	err := s.ValidateAddOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "favorites/add", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Remove removes a component as favorite for the authenticated user.
// Requires authentication.
//
// API endpoint: POST /api/favorites/remove.
// Since: 6.3.
func (s *FavoritesService) Remove(opt *FavoritesRemoveOption) (*http.Response, error) {
	err := s.ValidateRemoveOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "favorites/remove", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Search searches for the authenticated user favorites.
// Requires authentication.
//
// API endpoint: GET /api/favorites/search.
// Since: 6.3.
func (s *FavoritesService) Search(opt *FavoritesSearchOption) (*FavoritesSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "favorites/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(FavoritesSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
