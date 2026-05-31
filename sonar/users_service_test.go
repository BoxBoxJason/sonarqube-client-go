package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers_Anonymize(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/anonymize", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersAnonymizeOptions{
		Login: "deactivated-user",
	}

	resp, err := client.Users.Anonymize(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_Anonymize_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Users.Anonymize(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Users.Anonymize(context.Background(), &UsersAnonymizeOptions{})
	assert.Error(t, err)
}

func TestUsers_ChangePassword(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/change_password", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersChangePasswordOptions{
		Login:    "myuser",
		Password: "MyNewPassword123!",
	}

	resp, err := client.Users.ChangePassword(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_ChangePassword_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersChangePasswordOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersChangePasswordOptions{Password: "MyNewPassword123!"},
		},
		{
			name: "missing password",
			opt:  &UsersChangePasswordOptions{Login: "myuser"},
		},
		{
			name: "password too short",
			opt:  &UsersChangePasswordOptions{Login: "myuser", Password: "short"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.ChangePassword(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_Create(t *testing.T) {
	response := UsersCreate{
		User: User{
			Active:      true,
			Email:       "newuser@example.com",
			Local:       true,
			Login:       "newuser",
			Name:        "New User",
			ScmAccounts: []string{"scm1", "scm2"},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/users/create", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersCreateOptions{
		Login:       "newuser",
		Name:        "New User",
		Email:       "newuser@example.com",
		Password:    "SecurePassword123!",
		Local:       true,
		ScmAccounts: []string{"scm1", "scm2"},
	}

	result, resp, err := client.Users.Create(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "newuser", result.User.Login)
	assert.Equal(t, "New User", result.User.Name)
	assert.True(t, result.User.Active)
	assert.Len(t, result.User.ScmAccounts, 2)
}

func TestUsers_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersCreateOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersCreateOptions{Name: "Test User"},
		},
		{
			name: "missing name",
			opt:  &UsersCreateOptions{Login: "testuser"},
		},
		{
			name: "login too short",
			opt:  &UsersCreateOptions{Login: "x", Name: "Test User"},
		},
		{
			name: "local user without password",
			opt:  &UsersCreateOptions{Login: "testuser", Name: "Test User", Local: true},
		},
		{
			name: "password too short",
			opt:  &UsersCreateOptions{Login: "testuser", Name: "Test User", Local: true, Password: "short"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Users.Create(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_Current(t *testing.T) {
	response := UsersCurrentProfile{
		Avatar: "abc123",
		DismissedNotices: UsersDismissedNotices{
			EducationPrinciples: true,
			SonarlintAd:         false,
		},
		Email:            "admin@example.com",
		ExternalIdentity: "admin",
		ExternalProvider: "sonarqube",
		Groups:           []string{"sonar-users", "sonar-administrators"},
		Homepage: UsersHomepage{
			Component: "my-project",
			Type:      "PROJECT",
		},
		ID:                          "uuid-123",
		IsLoggedIn:                  true,
		Local:                       true,
		Login:                       "admin",
		Name:                        "Administrator",
		Permissions:                 UsersPermissions{Global: []string{"admin", "scan"}},
		ScmAccounts:                 []string{"admin@scm.com"},
		UsingSonarLintConnectedMode: true,
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/current", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Users.Current(context.Background(), )
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "admin", result.Login)
	assert.True(t, result.IsLoggedIn)
	assert.Equal(t, "PROJECT", result.Homepage.Type)
	assert.True(t, result.DismissedNotices.EducationPrinciples)
	assert.Len(t, result.Permissions.Global, 2)
}

func TestUsers_Deactivate(t *testing.T) {
	response := UsersDeactivate{
		User: UsersDeactivateResult{
			Active: false,
			Groups: []any{},
			Local:  true,
			Login:  "myuser",
			Name:   "My User",
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/users/deactivate", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersDeactivateOptions{
		Login:     "myuser",
		Anonymize: false,
	}

	result, resp, err := client.Users.Deactivate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.False(t, result.User.Active)
	assert.Equal(t, "myuser", result.User.Login)
}

func TestUsers_Deactivate_WithAnonymize(t *testing.T) {
	response := UsersDeactivate{
		User: UsersDeactivateResult{
			Active: false,
			Groups: []any{},
			Local:  true,
			Login:  "anonymized",
			Name:   "Anonymous",
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/users/deactivate", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersDeactivateOptions{
		Login:     "myuser",
		Anonymize: true,
	}

	result, resp, err := client.Users.Deactivate(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestUsers_Deactivate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Users.Deactivate(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, _, err = client.Users.Deactivate(context.Background(), &UsersDeactivateOptions{})
	assert.Error(t, err)
}

func TestUsers_DismissNotice(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/dismiss_notice", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersDismissNoticeOptions{
		Notice: "educationPrinciples",
	}

	resp, err := client.Users.DismissNotice(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_DismissNotice_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersDismissNoticeOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing notice",
			opt:  &UsersDismissNoticeOptions{},
		},
		{
			name: "invalid notice",
			opt:  &UsersDismissNoticeOptions{Notice: "invalidNotice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.DismissNotice(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_Groups(t *testing.T) {
	response := UsersGroups{
		Groups: []UsersGroup{
			{
				Default:     true,
				Description: "Default group",
				ID:          "1",
				Name:        "sonar-users",
				Selected:    true,
			},
			{
				Default:     false,
				Description: "Administrators",
				ID:          "2",
				Name:        "sonar-administrators",
				Selected:    false,
			},
		},
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     2,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersGroupsOptions{
		Login: "myuser",
	}

	result, resp, err := client.Users.Groups(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 2)
	assert.Equal(t, "sonar-users", result.Groups[0].Name)
	assert.True(t, result.Groups[0].Default)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestUsers_Groups_WithPagination(t *testing.T) {
	response := UsersGroups{
		Groups: []UsersGroup{},
		Paging: Paging{
			PageIndex: 2,
			PageSize:  10,
			Total:     2,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersGroupsOptions{
		Login: "myuser",
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 10,
		},
	}

	result, resp, err := client.Users.Groups(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(2), result.Paging.PageIndex)
}

func TestUsers_Groups_WithFilter(t *testing.T) {
	response := UsersGroups{
		Groups: []UsersGroup{},
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersGroupsOptions{
		Login:    "myuser",
		Query:    "admin",
		Selected: SelectionFilterSelected,
	}

	_, resp, err := client.Users.Groups(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUsers_Groups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersGroupsOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersGroupsOptions{},
		},
		{
			name: "invalid selected",
			opt:  &UsersGroupsOptions{Login: "myuser", Selected: "invalid"},
		},
		{
			name: "invalid page size",
			opt:  &UsersGroupsOptions{Login: "myuser", PaginationArgs: PaginationArgs{PageSize: 1000}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Users.Groups(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_IdentityProviders(t *testing.T) {
	response := UsersIdentityProviders{
		IdentityProviders: []UsersIdentityProvider{
			{
				BackgroundColor: "#444444",
				HelpMessage:     "Use your LDAP credentials",
				IconPath:        "/images/ldap.png",
				Key:             "LDAP",
				Name:            "LDAP",
			},
			{
				BackgroundColor: "#000000",
				HelpMessage:     "Use your GitHub account",
				IconPath:        "/images/github.png",
				Key:             "github",
				Name:            "GitHub",
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/identity_providers", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Users.IdentityProviders(context.Background(), )
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.IdentityProviders, 2)
	assert.Equal(t, "LDAP", result.IdentityProviders[0].Key)
	assert.Equal(t, "GitHub", result.IdentityProviders[1].Name)
}

func TestUsers_Search(t *testing.T) {
	response := UsersSearch{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  50,
			Total:     2,
		},
		Users: []UsersSearchResult{
			{
				Active:                      true,
				Avatar:                      "abc123",
				Email:                       "admin@example.com",
				ExternalIdentity:            "admin",
				ExternalProvider:            "sonarqube",
				Groups:                      []string{"sonar-users", "sonar-administrators"},
				LastConnectionDate:          "2024-01-01T00:00:00+0000",
				Local:                       true,
				Login:                       "admin",
				Managed:                     false,
				Name:                        "Administrator",
				ScmAccounts:                 []string{"admin@scm.com"},
				SonarLintLastConnectionDate: "2024-01-01T00:00:00+0000",
				TokensCount:                 3,
			},
			{
				Active:  true,
				Login:   "user1",
				Name:    "User One",
				Local:   false,
				Managed: true,
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Users.Search(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, "admin", result.Users[0].Login)
	assert.Equal(t, int64(3), result.Users[0].TokensCount)
	assert.True(t, result.Users[1].Managed)
}

func TestUsers_Search_WithFilters(t *testing.T) {
	response := UsersSearch{
		Paging: Paging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
		Users: []UsersSearchResult{},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersSearchOptions{
		Deactivated: true,
		Query:       "test",
		PaginationArgs: PaginationArgs{
			Page:     1,
			PageSize: 25,
		},
	}

	_, resp, err := client.Users.Search(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUsers_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid page size should fail validation.
	opt := &UsersSearchOptions{
		PaginationArgs: PaginationArgs{PageSize: 1000},
	}
	_, _, err := client.Users.Search(context.Background(), opt)
	assert.Error(t, err)
}

func TestUsers_SetHomepage(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/set_homepage", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersSetHomepageOptions{
		Type:      "PROJECT",
		Component: "my-project",
	}

	resp, err := client.Users.SetHomepage(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_SetHomepage_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersSetHomepageOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing type",
			opt:  &UsersSetHomepageOptions{},
		},
		{
			name: "invalid type",
			opt:  &UsersSetHomepageOptions{Type: "INVALID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.SetHomepage(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_Update(t *testing.T) {
	response := UsersUpdate{
		User: User{
			Active:      true,
			Email:       "updated@example.com",
			Local:       true,
			Login:       "myuser",
			Name:        "Updated Name",
			ScmAccounts: []string{"scm1"},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/users/update", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersUpdateOptions{
		Login: "myuser",
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	result, resp, err := client.Users.Update(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Name", result.User.Name)
}

func TestUsers_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Users.Update(context.Background(), nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, _, err = client.Users.Update(context.Background(), &UsersUpdateOptions{})
	assert.Error(t, err)
}

func TestUsers_UpdateIdentityProvider(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/update_identity_provider", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersUpdateIdentityProviderOptions{
		Login:               "myuser",
		NewExternalProvider: "github",
		NewExternalIdentity: "github-user",
	}

	resp, err := client.Users.UpdateIdentityProvider(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_UpdateIdentityProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersUpdateIdentityProviderOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersUpdateIdentityProviderOptions{NewExternalProvider: "github"},
		},
		{
			name: "missing newExternalProvider",
			opt:  &UsersUpdateIdentityProviderOptions{Login: "myuser"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.UpdateIdentityProvider(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_UpdateLogin(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/update_login", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersUpdateLoginOptions{
		Login:    "oldlogin",
		NewLogin: "newlogin",
	}

	resp, err := client.Users.UpdateLogin(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_UpdateLogin_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersUpdateLoginOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersUpdateLoginOptions{NewLogin: "newlogin"},
		},
		{
			name: "missing newLogin",
			opt:  &UsersUpdateLoginOptions{Login: "oldlogin"},
		},
		{
			name: "newLogin too short",
			opt:  &UsersUpdateLoginOptions{Login: "oldlogin", NewLogin: "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.UpdateLogin(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsersService_SearchAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"users":[{"login":"user1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"users":[{"login":"user2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &UsersSearchOptions{}
		result, _, err := client.Users.SearchAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})
}

func TestUsersService_GroupsAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"groups":[{"name":"g1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"groups":[{"name":"g2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &UsersGroupsOptions{Login: "user1"}
		result, _, err := client.Users.GroupsAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Users.GroupsAll(context.Background(), nil)
		assert.Error(t, err)
	})
}
