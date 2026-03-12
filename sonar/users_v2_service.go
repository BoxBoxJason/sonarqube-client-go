package sonar

import (
	"fmt"
	"net/http"
)

const (
	// MaxLoginLengthV2 is the maximum length for a user login in V2 API.
	MaxLoginLengthV2 = 100
	// MaxExternalLoginLengthV2 is the maximum length for an external login in V2 API.
	MaxExternalLoginLengthV2 = 255
	// MaxExternalIdLengthV2 is the maximum length for an external ID in V2 API.
	MaxExternalIdLengthV2 = 255
)

// UsersManagementServiceV2 handles communication with the Users Management
// related methods of the SonarQube V2 API.
type UsersManagementServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// UserV2 represents a user returned by V2 API endpoints.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UserV2 struct {
	// Active indicates whether the user is active.
	Active bool `json:"active,omitempty"`
	// Avatar is the user's avatar URL.
	Avatar string `json:"avatar,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// ExternalId is the user's external ID in the authentication system.
	ExternalId string `json:"externalId,omitempty"`
	// ExternalLogin is the user's login in the external authentication system.
	ExternalLogin string `json:"externalLogin,omitempty"`
	// ExternalProvider is the user's external identity provider.
	ExternalProvider string `json:"externalProvider,omitempty"`
	// Id is the user's unique identifier.
	Id string `json:"id,omitempty"`
	// Local indicates whether the user is authenticated locally.
	Local bool `json:"local,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Managed indicates whether the user is managed externally.
	Managed bool `json:"managed,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// ScmAccounts is the list of SCM accounts associated with the user.
	ScmAccounts []string `json:"scmAccounts,omitempty"`
	// SonarLintLastConnectionDate is the user's last SonarLint connection date.
	SonarLintLastConnectionDate string `json:"sonarLintLastConnectionDate,omitempty"`
	// SonarQubeLastConnectionDate is the user's last SonarQube connection date.
	SonarQubeLastConnectionDate string `json:"sonarQubeLastConnectionDate,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// UsersSearchV2 represents the response from searching users via the V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UsersSearchV2 struct {
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
	// Users is the list of users.
	Users []UserV2 `json:"users,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types (Query Parameters)
// -----------------------------------------------------------------------------

// UsersSearchOptionV2 contains query parameters for the Search method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UsersSearchOptionV2 struct {
	PaginationParamsV2

	// Active filters active/inactive users. Default is true.
	Active *bool `json:"active,omitempty"`
	// ExternalIdentity filters by external identity (case-sensitive exact match).
	ExternalIdentity string `json:"externalIdentity,omitempty"`
	// GroupId filters users belonging to a group. Only available for system administrators.
	// Internal: this parameter is marked as internal in the SonarQube API.
	GroupId string `json:"groupId,omitempty"`
	// Managed filters managed or non-managed users.
	Managed *bool `json:"managed,omitempty"`
	// Query filters on login, name and email (partial match, case insensitive).
	Query string `json:"q,omitempty"`
	// SonarLintLastConnectionDateFrom filters users who connected via SonarLint at or after this date (ISO 8601).
	SonarLintLastConnectionDateFrom string `json:"sonarLintLastConnectionDateFrom,omitempty"`
	// SonarLintLastConnectionDateTo filters users who connected via SonarLint at or before this date (ISO 8601).
	SonarLintLastConnectionDateTo string `json:"sonarLintLastConnectionDateTo,omitempty"`
	// SonarQubeLastConnectionDateFrom filters users who connected at or after this date (ISO 8601).
	SonarQubeLastConnectionDateFrom string `json:"sonarQubeLastConnectionDateFrom,omitempty"`
	// SonarQubeLastConnectionDateTo filters users who connected at or before this date (ISO 8601).
	SonarQubeLastConnectionDateTo string `json:"sonarQubeLastConnectionDateTo,omitempty"`
}

// UsersDeactivateOptionsV2 contains parameters for the Deactivate method.
type UsersDeactivateOptionsV2 struct {
	// Id is the user's unique identifier.
	// This field is required.
	Id string `json:"-"`
	// Anonymize specifies whether to anonymize the user in addition to deactivating.
	Anonymize bool `json:"anonymize,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types (JSON Body)
// -----------------------------------------------------------------------------

// UsersCreateOptionsV2 contains parameters for creating a user via the V2 API.
type UsersCreateOptionsV2 struct {
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// Local specifies if the user should be authenticated locally.
	// When false, password should not be set. Default is true.
	Local *bool `json:"local,omitempty"`
	// Login is the user's login.
	// This field is required. Must be between 2 and 100 characters.
	Login string `json:"login"`
	// Name is the user's display name.
	// This field is required. Maximum 200 characters.
	Name string `json:"name"`
	// Password is the user's password. Only required when creating a local user.
	Password string `json:"password,omitempty"`
	// ScmAccounts is the list of SCM accounts.
	ScmAccounts []string `json:"scmAccounts,omitempty"`
}

// UsersUpdateOptionsV2 contains parameters for updating a user via the V2 API.
// All fields are optional (PATCH merge semantics).
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UsersUpdateOptionsV2 struct {
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// ExternalId is the new external ID in the authentication system.
	ExternalId string `json:"externalId,omitempty"`
	// ExternalLogin is the new external login.
	ExternalLogin string `json:"externalLogin,omitempty"`
	// ExternalProvider is the new identity provider.
	ExternalProvider string `json:"externalProvider,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// ScmAccounts is the list of SCM accounts.
	ScmAccounts *UpdateFieldListStringV2 `json:"scmAccounts,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateSearchOpt validates the UsersSearchOptionV2.
func (s *UsersManagementServiceV2) ValidateSearchOpt(opt *UsersSearchOptionV2) error {
	if opt == nil {
		return nil
	}

	return opt.Validate()
}

// ValidateCreateRequest validates the UsersCreateOptionsV2.
func (s *UsersManagementServiceV2) ValidateCreateRequest(opt *UsersCreateOptionsV2) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.Login, MinLoginLength, "Login")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Login, MaxLoginLengthV2, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Email, MaxEmailLength, "Email")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeactivateOpt validates the UsersDeactivateOptionsV2.
func (s *UsersManagementServiceV2) ValidateDeactivateOpt(opt *UsersDeactivateOptionsV2) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.Id, "Id")
}

// ValidateUpdateRequest validates the UsersUpdateOptionsV2.
func (s *UsersManagementServiceV2) ValidateUpdateRequest(userID string, opt *UsersUpdateOptionsV2) error {
	err := ValidateRequired(userID, "Id")
	if err != nil {
		return err
	}

	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err = ValidateMinLength(opt.Login, MinLoginLength, "Login")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Login, MaxLoginLengthV2, "Login")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Email, MaxEmailLength, "Email")
	if err != nil {
		return err
	}

	return validateUsersV2ExternalFields(opt)
}

// validateUsersV2ExternalFields validates external identity fields of a user update request.
func validateUsersV2ExternalFields(opt *UsersUpdateOptionsV2) error {
	err := ValidateMinLength(opt.ExternalLogin, 1, "ExternalLogin")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ExternalLogin, MaxExternalLoginLengthV2, "ExternalLogin")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.ExternalId, 1, "ExternalId")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ExternalId, MaxExternalIdLengthV2, "ExternalId")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Search returns a list of users matching the search criteria.
// By default, only active users are returned.
func (s *UsersManagementServiceV2) Search(opt *UsersSearchOptionV2) (*UsersSearchV2, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "users-management/users", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(UsersSearchV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Create creates a new user. If a deactivated user account exists with the
// given login, it will be reactivated.
// Requires Administer System permission.
func (s *UsersManagementServiceV2) Create(opt *UsersCreateOptionsV2) (*UserV2, *http.Response, error) {
	err := s.ValidateCreateRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPost, "users-management/users", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(UserV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Fetch retrieves a single user by ID.
func (s *UsersManagementServiceV2) Fetch(userID string) (*UserV2, *http.Response, error) {
	err := ValidateRequired(userID, "Id")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "users-management/users/"+userID, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(UserV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Deactivate deactivates a user.
// Requires Administer System permission.
func (s *UsersManagementServiceV2) Deactivate(opt *UsersDeactivateOptionsV2) (*http.Response, error) {
	err := s.ValidateDeactivateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodDelete, "users-management/users/"+opt.Id, opt, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update updates a user's attributes.
func (s *UsersManagementServiceV2) Update(userID string, opt *UsersUpdateOptionsV2) (*UserV2, *http.Response, error) {
	err := s.ValidateUpdateRequest(userID, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPatch, "users-management/users/"+userID, nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(UserV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
