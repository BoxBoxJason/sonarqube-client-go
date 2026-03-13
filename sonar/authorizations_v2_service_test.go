package sonar

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// SearchGroups
// =============================================================================

func TestAuthorizationsV2_SearchGroups(t *testing.T) {
	response := AuthorizationsGroupsSearch{
		Groups: []Group{
			{Id: "g1", Name: "admins", Description: "Administrator group", Managed: false},
			{Id: "g2", Name: "members", Default: true},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 2},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/authorizations/groups", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Authorizations.SearchGroups(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Groups, 2)
	assert.Equal(t, "admins", result.Groups[0].Name)
}

func TestAuthorizationsV2_SearchGroups_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Authorizations.SearchGroups(&AuthorizationsSearchGroupsOptions{
		PaginationParamsV2: PaginationParamsV2{PageSize: 600},
	})
	assert.Error(t, err)
}

// =============================================================================
// CreateGroup
// =============================================================================

func TestAuthorizationsV2_CreateGroup(t *testing.T) {
	response := Group{
		Id:          "g-new",
		Name:        "new-group",
		Description: "A new group",
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/authorizations/groups", http.StatusOK,
		&AuthorizationsCreateGroupOptions{
			Name:        "new-group",
			Description: "A new group",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Authorizations.CreateGroup(&AuthorizationsCreateGroupOptions{
		Name:        "new-group",
		Description: "A new group",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "new-group", result.Name)
}

func TestAuthorizationsV2_CreateGroup_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *AuthorizationsCreateGroupOptions
	}{
		{"nil opt", nil},
		{"missing name", &AuthorizationsCreateGroupOptions{}},
		{"name too long", &AuthorizationsCreateGroupOptions{Name: strings.Repeat("a", MaxGroupNameLength+1)}},
		{"description too long", &AuthorizationsCreateGroupOptions{Name: "ok", Description: strings.Repeat("a", MaxGroupDescriptionLength+1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.Authorizations.CreateGroup(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// FetchGroup
// =============================================================================

func TestAuthorizationsV2_FetchGroup(t *testing.T) {
	response := Group{
		Id:   "g1",
		Name: "admins",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/authorizations/groups/g1", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Authorizations.FetchGroup("g1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "admins", result.Name)
}

func TestAuthorizationsV2_FetchGroup_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Authorizations.FetchGroup("")
	assert.Error(t, err)
}

// =============================================================================
// DeleteGroup
// =============================================================================

func TestAuthorizationsV2_DeleteGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/authorizations/groups/g1", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Authorizations.DeleteGroup("g1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAuthorizationsV2_DeleteGroup_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Authorizations.DeleteGroup("")
	assert.Error(t, err)
}

// =============================================================================
// UpdateGroup
// =============================================================================

func TestAuthorizationsV2_UpdateGroup(t *testing.T) {
	response := Group{
		Id:          "g1",
		Name:        "renamed-group",
		Description: "Updated description",
	}
	server := newTestServer(t, mockPatchHandler(t, "/v2/authorizations/groups/g1", http.StatusOK,
		&AuthorizationsUpdateGroupOptions{
			Name:        "renamed-group",
			Description: "Updated description",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Authorizations.UpdateGroup("g1", &AuthorizationsUpdateGroupOptions{
		Name:        "renamed-group",
		Description: "Updated description",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "renamed-group", result.Name)
}

func TestAuthorizationsV2_UpdateGroup_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		id   string
		opt  *AuthorizationsUpdateGroupOptions
	}{
		{"missing id", "", &AuthorizationsUpdateGroupOptions{Name: "ok"}},
		{"nil opt", "g1", nil},
		{"name too long", "g1", &AuthorizationsUpdateGroupOptions{Name: strings.Repeat("a", MaxGroupNameLength+1)}},
		{"description too long", "g1", &AuthorizationsUpdateGroupOptions{Description: strings.Repeat("a", MaxGroupDescriptionLength+1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.Authorizations.UpdateGroup(tt.id, tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// SearchGroupMemberships
// =============================================================================

func TestAuthorizationsV2_SearchGroupMemberships(t *testing.T) {
	response := AuthorizationsGroupMembershipsSearch{
		GroupMemberships: []GroupMembership{
			{Id: "m1", GroupId: "g1", UserId: "u1"},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/authorizations/group-memberships", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Authorizations.SearchGroupMemberships(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.GroupMemberships, 1)
	assert.Equal(t, "g1", result.GroupMemberships[0].GroupId)
}

func TestAuthorizationsV2_SearchGroupMemberships_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Authorizations.SearchGroupMemberships(&AuthorizationsSearchGroupMembershipsOptions{
		PaginationParamsV2: PaginationParamsV2{PageSize: 600},
	})
	assert.Error(t, err)
}

// =============================================================================
// CreateGroupMembership
// =============================================================================

func TestAuthorizationsV2_CreateGroupMembership(t *testing.T) {
	response := GroupMembership{
		Id:      "m-new",
		GroupId: "g1",
		UserId:  "u1",
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/authorizations/group-memberships", http.StatusOK,
		&AuthorizationsCreateGroupMembershipOptions{
			GroupId: "g1",
			UserId:  "u1",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Authorizations.CreateGroupMembership(&AuthorizationsCreateGroupMembershipOptions{
		GroupId: "g1",
		UserId:  "u1",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "g1", result.GroupId)
	assert.Equal(t, "u1", result.UserId)
}

func TestAuthorizationsV2_CreateGroupMembership_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *AuthorizationsCreateGroupMembershipOptions
	}{
		{"nil request", nil},
		{"missing group ID", &AuthorizationsCreateGroupMembershipOptions{UserId: "u1"}},
		{"missing user ID", &AuthorizationsCreateGroupMembershipOptions{GroupId: "g1"}},
		{"empty both", &AuthorizationsCreateGroupMembershipOptions{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.Authorizations.CreateGroupMembership(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// DeleteGroupMembership
// =============================================================================

func TestAuthorizationsV2_DeleteGroupMembership(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/authorizations/group-memberships/m1", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Authorizations.DeleteGroupMembership("m1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAuthorizationsV2_DeleteGroupMembership_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Authorizations.DeleteGroupMembership("")
	assert.Error(t, err)
}
