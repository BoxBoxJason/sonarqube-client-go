package sonar

import (
	"net/http"
)

const (
	// MaxGroupNameLength is the maximum allowed length for group names.
	MaxGroupNameLength = 255
	// MaxGroupDescriptionLength is the maximum allowed length for group descriptions.
	MaxGroupDescriptionLength = 200
)

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedUserGroupSearchFields is the set of fields that can be returned in search response.
	allowedUserGroupSearchFields = map[string]struct{}{
		"name":         {},
		"description":  {},
		"membersCount": {},
		"managed":      {},
	}

	// allowedUserGroupsUsersSelected is the set of allowed values for user selection filter.
	allowedUserGroupsUsersSelected = map[string]struct{}{
		"all":        {},
		"deselected": {},
		"selected":   {},
	}
)

// UserGroupsService handles communication with the User Groups related methods of the SonarQube API.
// Manage user groups.
//
// Deprecated: Since 10.4. Use v2 API endpoints instead.
//
// Since: 5.2.
type UserGroupsService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// UserGroupDetail represents a user group with its properties from user_groups API.
type UserGroupDetail struct {
	// Description is the description of the group.
	Description string `json:"description,omitempty"`
	// ID is the unique identifier of the group.
	ID string `json:"id,omitempty"`
	// Name is the name of the group.
	Name string `json:"name,omitempty"`
	// Organization is the organization the group belongs to.
	Organization string `json:"organization,omitempty"`
	// MembersCount is the number of members in the group.
	MembersCount int64 `json:"membersCount,omitempty"`
	// Default indicates if this is a default group.
	Default bool `json:"default,omitempty"`
	// Managed indicates if the group is managed externally.
	Managed bool `json:"managed,omitempty"`
}

// UserGroupsCreate represents the response from creating a group.
type UserGroupsCreate struct {
	// Group contains the created group details.
	Group UserGroupDetail `json:"group,omitzero"`
}

// UserGroupsSearch represents the response from searching groups.
type UserGroupsSearch struct {
	// Groups is the list of groups.
	Groups []UserGroupDetail `json:"groups,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// UserGroupUser represents a user within a group.
type UserGroupUser struct {
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// Managed indicates if the user is managed externally.
	Managed bool `json:"managed,omitempty"`
	// Selected indicates if the user is selected/member of the group.
	Selected bool `json:"selected,omitempty"`
}

// UserGroupsUsers represents the response from listing users in a group.
type UserGroupsUsers struct {
	// Users is the list of users.
	Users []UserGroupUser `json:"users,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// UserGroupsAddUserOption represents options for adding a user to a group.
type UserGroupsAddUserOption struct {
	// Login is the user login (optional - for internal accounts).
	Login string `url:"login,omitempty"`
	// Name is the group name (required).
	Name string `url:"name,omitempty"`
}

// UserGroupsCreateOption represents options for creating a group.
type UserGroupsCreateOption struct {
	// Description is the description for the new group.
	// Maximum length: 200 characters.
	Description string `url:"description,omitempty"`
	// Name is the name for the new group (required).
	// Must be unique and cannot be 'anyone'.
	// Maximum length: 255 characters.
	Name string `url:"name,omitempty"`
}

// UserGroupsDeleteOption represents options for deleting a group.
type UserGroupsDeleteOption struct {
	// Name is the group name (required).
	Name string `url:"name,omitempty"`
}

// UserGroupsRemoveUserOption represents options for removing a user from a group.
type UserGroupsRemoveUserOption struct {
	// Login is the user login (optional - for internal accounts).
	Login string `url:"login,omitempty"`
	// Name is the group name (required).
	Name string `url:"name,omitempty"`
}

// UserGroupsSearchOption represents options for searching groups.
//
//nolint:govet // Embedded PaginationArgs makes optimal alignment impractical
type UserGroupsSearchOption struct {
	PaginationArgs

	// Managed filters by managed status.
	// Only available for managed instances.
	Managed *bool `url:"managed,omitempty"`
	// Fields is the list of fields to return in response.
	// Possible values: name, description, membersCount, managed.
	Fields []string `url:"f,omitempty,comma"`
	// Query limits search to names that contain the supplied string.
	Query string `url:"q,omitempty"`
}

// UserGroupsUpdateOption represents options for updating a group.
type UserGroupsUpdateOption struct {
	// CurrentName is the current name of the group to update (required).
	CurrentName string `url:"currentName,omitempty"`
	// Description is the new optional description for the group.
	// Maximum length: 200 characters.
	Description string `url:"description,omitempty"`
	// Name is the new optional name for the group.
	// Cannot be 'anyone'.
	// Maximum length: 255 characters.
	Name string `url:"name,omitempty"`
}

// UserGroupsUsersOption represents options for listing users in a group.
//
//nolint:govet // Embedded PaginationArgs makes optimal alignment impractical
type UserGroupsUsersOption struct {
	PaginationArgs

	// Name is the group name (required).
	Name string `url:"name,omitempty"`
	// Query limits search to names or logins that contain the supplied string.
	Query string `url:"q,omitempty"`
	// Selected filters by selection status.
	// Possible values: all, deselected, selected.
	Selected string `url:"selected,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Methods
// -----------------------------------------------------------------------------

// ValidateAddUserOpt validates the options for AddUser.
func (s *UserGroupsService) ValidateAddUserOpt(opt *UserGroupsAddUserOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Name, "Name")
}

// ValidateCreateOpt validates the options for Create.
func (s *UserGroupsService) ValidateCreateOpt(opt *UserGroupsCreateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxGroupNameLength, "Name")
	if err != nil {
		return err
	}

	if opt.Description != "" {
		err = ValidateMaxLength(opt.Description, MaxGroupDescriptionLength, "Description")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteOpt validates the options for Delete.
func (s *UserGroupsService) ValidateDeleteOpt(opt *UserGroupsDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Name, "Name")
}

// ValidateRemoveUserOpt validates the options for RemoveUser.
func (s *UserGroupsService) ValidateRemoveUserOpt(opt *UserGroupsRemoveUserOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Name, "Name")
}

// ValidateSearchOpt validates the options for Search.
func (s *UserGroupsService) ValidateSearchOpt(opt *UserGroupsSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	return AreValuesAuthorized(opt.Fields, allowedUserGroupSearchFields, "Fields")
}

// ValidateUpdateOpt validates the options for Update.
func (s *UserGroupsService) ValidateUpdateOpt(opt *UserGroupsUpdateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.CurrentName, "CurrentName")
	if err != nil {
		return err
	}

	if opt.Name != "" {
		err = ValidateMaxLength(opt.Name, MaxGroupNameLength, "Name")
		if err != nil {
			return err
		}
	}

	if opt.Description != "" {
		err = ValidateMaxLength(opt.Description, MaxGroupDescriptionLength, "Description")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUsersOpt validates the options for Users.
func (s *UserGroupsService) ValidateUsersOpt(opt *UserGroupsUsersOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Selected, allowedUserGroupsUsersSelected, "Selected")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// AddUser adds a user to a group.
// 'name' must be provided.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use POST /api/v2/authorizations/group-memberships instead.
//
// Since: 5.2.
func (s *UserGroupsService) AddUser(opt *UserGroupsAddUserOption) (*http.Response, error) {
	err := s.ValidateAddUserOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_groups/add_user", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Create creates a group.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use POST /api/v2/authorizations/groups instead.
//
// Since: 5.2.
func (s *UserGroupsService) Create(opt *UserGroupsCreateOption) (*UserGroupsCreate, *http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_groups/create", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserGroupsCreate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes a group. The default groups cannot be deleted.
// 'name' must be provided.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use DELETE /api/v2/authorizations/groups instead.
//
// Since: 5.2.
func (s *UserGroupsService) Delete(opt *UserGroupsDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_groups/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveUser removes a user from a group.
// 'name' must be provided.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use DELETE /api/v2/authorizations/group-memberships instead.
//
// Since: 5.2.
func (s *UserGroupsService) RemoveUser(opt *UserGroupsRemoveUserOption) (*http.Response, error) {
	err := s.ValidateRemoveUserOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_groups/remove_user", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Search searches for user groups.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use GET /api/v2/authorizations/groups instead.
//
// Since: 5.2.
func (s *UserGroupsService) Search(opt *UserGroupsSearchOption) (*UserGroupsSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "user_groups/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserGroupsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Update updates a group.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use PATCH /api/v2/authorizations/groups instead.
//
// Since: 5.2.
func (s *UserGroupsService) Update(opt *UserGroupsUpdateOption) (*http.Response, error) {
	err := s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_groups/update", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Users searches for users with membership information with respect to a group.
// Requires the following permission: 'Administer System'.
//
// Deprecated: Since 10.4. Use GET /api/v2/authorizations/group-memberships instead.
//
// Since: 5.2.
func (s *UserGroupsService) Users(opt *UserGroupsUsersOption) (*UserGroupsUsers, *http.Response, error) {
	err := s.ValidateUsersOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "user_groups/users", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserGroupsUsers)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
