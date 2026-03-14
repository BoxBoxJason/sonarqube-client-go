package sonar

import (
	"fmt"
	"net/http"
)

// AuthorizationsService handles communication with the Authorizations related
// methods of the SonarQube V2 API.
type AuthorizationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// Group represents a group returned by V2 API endpoints.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type Group struct {
	// Default indicates whether this is a default group.
	Default bool `json:"default,omitempty"`
	// Description is the group description.
	Description string `json:"description,omitempty"`
	// Id is the group's unique identifier.
	Id string `json:"id,omitempty"`
	// Managed indicates whether the group is managed externally.
	Managed bool `json:"managed,omitempty"`
	// Name is the group name.
	Name string `json:"name,omitempty"`
}

// GroupMembership represents a group membership returned by V2 API endpoints.
type GroupMembership struct {
	// GroupId is the group's unique identifier.
	GroupId string `json:"groupId,omitempty"`
	// Id is the membership's unique identifier.
	Id string `json:"id,omitempty"`
	// UserId is the user's unique identifier.
	UserId string `json:"userId,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AuthorizationsGroupsSearch represents the response from searching groups.
type AuthorizationsGroupsSearch struct {
	// Groups is the list of groups.
	Groups []Group `json:"groups,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// AuthorizationsGroupMembershipsSearch represents the response from searching group memberships.
type AuthorizationsGroupMembershipsSearch struct {
	// GroupMemberships is the list of group memberships.
	GroupMemberships []GroupMembership `json:"groupMemberships,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// -----------------------------------------------------------------------------
// Option Types (Query Parameters)
// -----------------------------------------------------------------------------

// AuthorizationsSearchGroupsOptions contains query parameters for the SearchGroups method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type AuthorizationsSearchGroupsOptions struct {
	PaginationParamsV2

	// Managed filters managed or non-managed groups.
	Managed *bool `json:"managed,omitempty"`
	// Query filters on group name (partial match, case insensitive).
	Query string `json:"q,omitempty"`
	// UserId filters groups containing the user. Only available for system administrators.
	// Internal: this parameter is marked as internal in the SonarQube API.
	UserId string `json:"userId,omitempty"`
}

// AuthorizationsSearchGroupMembershipsOptions contains query parameters for the SearchGroupMemberships method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type AuthorizationsSearchGroupMembershipsOptions struct {
	PaginationParamsV2

	// GroupId filters memberships by group ID.
	GroupId string `json:"groupId,omitempty"`
	// UserId filters memberships by user ID.
	UserId string `json:"userId,omitempty"`
}

// -----------------------------------------------------------------------------
// Request Types (JSON Body)
// -----------------------------------------------------------------------------

// AuthorizationsCreateGroupOptions contains parameters for creating a group.
type AuthorizationsCreateGroupOptions struct {
	// Description is the group description. Maximum 200 characters.
	Description string `json:"description,omitempty"`
	// Name is the group name. Must be unique. The value 'anyone' is reserved.
	// This field is required. Must be between 1 and 255 characters.
	Name string `json:"name"`
}

// AuthorizationsUpdateGroupOptions contains parameters for updating a group.
// All fields are optional (PATCH merge semantics).
type AuthorizationsUpdateGroupOptions struct {
	// Description is the group description. Maximum 200 characters.
	// Use nil to leave unchanged, or a pointer to an empty string to clear it.
	Description *string `json:"description,omitempty"`
	// Name is the group name. Must be between 1 and 255 characters.
	Name string `json:"name,omitempty"`
}

// AuthorizationsCreateGroupMembershipOptions contains parameters for creating a group membership.
type AuthorizationsCreateGroupMembershipOptions struct {
	// GroupId is the ID of the group.
	GroupId string `json:"groupId,omitempty"`
	// UserId is the ID of the user.
	UserId string `json:"userId,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateSearchGroupsOpt validates the AuthorizationsSearchGroupsOptions.
func (s *AuthorizationsService) ValidateSearchGroupsOpt(opt *AuthorizationsSearchGroupsOptions) error {
	if opt == nil {
		return nil
	}

	return opt.Validate()
}

// ValidateSearchGroupMembershipsOpt validates the AuthorizationsSearchGroupMembershipsOptions.
func (s *AuthorizationsService) ValidateSearchGroupMembershipsOpt(opt *AuthorizationsSearchGroupMembershipsOptions) error {
	if opt == nil {
		return nil
	}

	return opt.Validate()
}

// ValidateCreateGroupRequest validates the AuthorizationsCreateGroupOptions.
func (s *AuthorizationsService) ValidateCreateGroupRequest(opt *AuthorizationsCreateGroupOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxGroupNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Description, MaxGroupDescriptionLength, "Description")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateGroupRequest validates the AuthorizationsUpdateGroupOptions.
func (s *AuthorizationsService) ValidateUpdateGroupRequest(groupID string, opt *AuthorizationsUpdateGroupOptions) error {
	err := ValidateRequired(groupID, "Id")
	if err != nil {
		return err
	}

	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err = ValidateMinLength(opt.Name, 1, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxGroupNameLength, "Name")
	if err != nil {
		return err
	}

	if opt.Description != nil {
		err = ValidateMaxLength(*opt.Description, MaxGroupDescriptionLength, "Description")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateCreateGroupMembershipRequest validates the AuthorizationsCreateGroupMembershipOptions.
func (s *AuthorizationsService) ValidateCreateGroupMembershipRequest(opt *AuthorizationsCreateGroupMembershipOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GroupId, "GroupId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.UserId, "UserId")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// SearchGroups returns a list of groups matching the search criteria.
// The results are sorted alphabetically by group name.
func (s *AuthorizationsService) SearchGroups(opt *AuthorizationsSearchGroupsOptions) (*AuthorizationsGroupsSearch, *http.Response, error) {
	err := s.ValidateSearchGroupsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "authorizations/groups", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AuthorizationsGroupsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateGroup creates a new group.
func (s *AuthorizationsService) CreateGroup(opt *AuthorizationsCreateGroupOptions) (*Group, *http.Response, error) {
	err := s.ValidateCreateGroupRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPost, "authorizations/groups", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(Group)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// FetchGroup retrieves a single group by ID.
func (s *AuthorizationsService) FetchGroup(groupID string) (*Group, *http.Response, error) {
	err := ValidateRequired(groupID, "Id")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "authorizations/groups/"+groupID, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(Group)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteGroup deletes a group by ID.
func (s *AuthorizationsService) DeleteGroup(groupID string) (*http.Response, error) {
	err := ValidateRequired(groupID, "Id")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodDelete, "authorizations/groups/"+groupID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateGroup updates a group's name or description.
func (s *AuthorizationsService) UpdateGroup(groupID string, opt *AuthorizationsUpdateGroupOptions) (*Group, *http.Response, error) {
	err := s.ValidateUpdateGroupRequest(groupID, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPatch, "authorizations/groups/"+groupID, nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(Group)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchGroupMemberships returns a list of group memberships matching the search criteria.
func (s *AuthorizationsService) SearchGroupMemberships(opt *AuthorizationsSearchGroupMembershipsOptions) (*AuthorizationsGroupMembershipsSearch, *http.Response, error) {
	err := s.ValidateSearchGroupMembershipsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "authorizations/group-memberships", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AuthorizationsGroupMembershipsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateGroupMembership adds a user to a group.
func (s *AuthorizationsService) CreateGroupMembership(opt *AuthorizationsCreateGroupMembershipOptions) (*GroupMembership, *http.Response, error) {
	err := s.ValidateCreateGroupMembershipRequest(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPost, "authorizations/group-memberships", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(GroupMembership)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteGroupMembership removes a user from a group.
func (s *AuthorizationsService) DeleteGroupMembership(membershipID string) (*http.Response, error) {
	err := ValidateRequired(membershipID, "Id")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodDelete, "authorizations/group-memberships/"+membershipID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
