package sonar

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGroups_AddUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/user_groups/add_user", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsAddUserOptions{
		Name:  "sonar-administrators",
		Login: "g.hopper",
	}

	resp, err := client.UserGroups.AddUser(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserGroups_AddUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.UserGroups.AddUser(nil)
	assert.Error(t, err)

	// Test missing Name
	_, err = client.UserGroups.AddUser(&UserGroupsAddUserOptions{
		Login: "user",
	})
	assert.Error(t, err)
}

func TestUserGroups_Create(t *testing.T) {
	response := UserGroupsCreate{
		Group: UserGroupDetail{
			ID:           "uuid-group-1",
			Name:         "sonar-users",
			Description:  "Default group",
			MembersCount: 0,
			Default:      false,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/user_groups/create", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsCreateOptions{
		Name:        "sonar-users",
		Description: "Default group",
	}

	result, resp, err := client.UserGroups.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "sonar-users", result.Group.Name)
}

func TestUserGroups_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.UserGroups.Create(nil)
	assert.Error(t, err)

	// Test missing Name
	_, _, err = client.UserGroups.Create(&UserGroupsCreateOptions{
		Description: "test",
	})
	assert.Error(t, err)

	// Test Name too long
	_, _, err = client.UserGroups.Create(&UserGroupsCreateOptions{
		Name: strings.Repeat("a", MaxGroupNameLength+1),
	})
	assert.Error(t, err)

	// Test Description too long
	_, _, err = client.UserGroups.Create(&UserGroupsCreateOptions{
		Name:        "test",
		Description: strings.Repeat("a", MaxGroupDescriptionLength+1),
	})
	assert.Error(t, err)
}

func TestUserGroups_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/user_groups/delete", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsDeleteOptions{
		Name: "sonar-users",
	}

	resp, err := client.UserGroups.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserGroups_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.UserGroups.Delete(nil)
	assert.Error(t, err)

	// Test missing Name
	_, err = client.UserGroups.Delete(&UserGroupsDeleteOptions{})
	assert.Error(t, err)
}

func TestUserGroups_RemoveUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/user_groups/remove_user", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsRemoveUserOptions{
		Name:  "sonar-administrators",
		Login: "g.hopper",
	}

	resp, err := client.UserGroups.RemoveUser(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserGroups_Search(t *testing.T) {
	response := UserGroupsSearch{
		Groups: []UserGroupDetail{
			{
				Name:         "sonar-administrators",
				Description:  "Admins",
				MembersCount: 3,
				Default:      false,
				Managed:      false,
			},
		},
		Paging: Paging{
			PageIndex: 1,
			PageSize:  100,
			Total:     1,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/user_groups/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsSearchOptions{
		Query:  "admin",
		Fields: []string{"name", "description"},
	}

	result, resp, err := client.UserGroups.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Groups, 1)
}

func TestUserGroups_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.UserGroups.Search(nil)
	assert.Error(t, err)

	// Test invalid field
	_, _, err = client.UserGroups.Search(&UserGroupsSearchOptions{
		Fields: []string{"invalid_field"},
	})
	assert.Error(t, err)
}

func TestUserGroups_Update(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/user_groups/update", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsUpdateOptions{
		CurrentName: "old-group",
		Name:        "new-group",
		Description: "Updated description",
	}

	resp, err := client.UserGroups.Update(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUserGroups_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.UserGroups.Update(nil)
	assert.Error(t, err)

	// Test missing CurrentName
	_, err = client.UserGroups.Update(&UserGroupsUpdateOptions{
		Name: "new-group",
	})
	assert.Error(t, err)
}

func TestUserGroups_Users(t *testing.T) {
	response := UserGroupsUsers{
		Users: []UserGroupUser{
			{
				Login:    "john.doe",
				Name:     "John Doe",
				Selected: true,
				Managed:  false,
			},
		},
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     1,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/user_groups/users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UserGroupsUsersOptions{
		Name:     "sonar-administrators",
		Selected: SelectionFilterSelected,
	}

	result, resp, err := client.UserGroups.Users(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Users, 1)
}

func TestUserGroups_Users_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.UserGroups.Users(nil)
	assert.Error(t, err)

	// Test missing Name
	_, _, err = client.UserGroups.Users(&UserGroupsUsersOptions{})
	assert.Error(t, err)

	// Test invalid Selected value
	_, _, err = client.UserGroups.Users(&UserGroupsUsersOptions{
		Name:     "test",
		Selected: "invalid",
	})
	assert.Error(t, err)
}
