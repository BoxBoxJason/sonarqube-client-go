package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers_Anonymize(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/anonymize", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersAnonymizeOption{
		Login: "deactivated-user",
	}

	resp, err := client.Users.Anonymize(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_Anonymize_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Users.Anonymize(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, err = client.Users.Anonymize(&UsersAnonymizeOption{})
	assert.Error(t, err)
}

func TestUsers_ChangePassword(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/change_password", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersChangePasswordOption{
		Login:    "myuser",
		Password: "MyNewPassword123!",
	}

	resp, err := client.Users.ChangePassword(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_ChangePassword_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersChangePasswordOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersChangePasswordOption{Password: "MyNewPassword123!"},
		},
		{
			name: "missing password",
			opt:  &UsersChangePasswordOption{Login: "myuser"},
		},
		{
			name: "password too short",
			opt:  &UsersChangePasswordOption{Login: "myuser", Password: "short"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.ChangePassword(tt.opt)
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

	opt := &UsersCreateOption{
		Login:       "newuser",
		Name:        "New User",
		Email:       "newuser@example.com",
		Password:    "SecurePassword123!",
		Local:       true,
		ScmAccounts: []string{"scm1", "scm2"},
	}

	result, resp, err := client.Users.Create(opt)
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
		opt  *UsersCreateOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersCreateOption{Name: "Test User"},
		},
		{
			name: "missing name",
			opt:  &UsersCreateOption{Login: "testuser"},
		},
		{
			name: "login too short",
			opt:  &UsersCreateOption{Login: "x", Name: "Test User"},
		},
		{
			name: "local user without password",
			opt:  &UsersCreateOption{Login: "testuser", Name: "Test User", Local: true},
		},
		{
			name: "password too short",
			opt:  &UsersCreateOption{Login: "testuser", Name: "Test User", Local: true, Password: "short"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Users.Create(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_Current(t *testing.T) {
	response := CurrentUser{
		Avatar: "abc123",
		DismissedNotices: DismissedNotices{
			EducationPrinciples: true,
			SonarlintAd:         false,
		},
		Email:            "admin@example.com",
		ExternalIdentity: "admin",
		ExternalProvider: "sonarqube",
		Groups:           []string{"sonar-users", "sonar-administrators"},
		Homepage: Homepage{
			Component: "my-project",
			Type:      "PROJECT",
		},
		ID:                          "uuid-123",
		IsLoggedIn:                  true,
		Local:                       true,
		Login:                       "admin",
		Name:                        "Administrator",
		Permissions:                 UserPermissions{Global: []string{"admin", "scan"}},
		ScmAccounts:                 []string{"admin@scm.com"},
		UsingSonarLintConnectedMode: true,
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/current", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Users.Current()
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
		User: DeactivatedUser{
			Active: false,
			Groups: []any{},
			Local:  true,
			Login:  "myuser",
			Name:   "My User",
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/users/deactivate", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersDeactivateOption{
		Login:     "myuser",
		Anonymize: false,
	}

	result, resp, err := client.Users.Deactivate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.False(t, result.User.Active)
	assert.Equal(t, "myuser", result.User.Login)
}

func TestUsers_Deactivate_WithAnonymize(t *testing.T) {
	response := UsersDeactivate{
		User: DeactivatedUser{
			Active: false,
			Groups: []any{},
			Local:  true,
			Login:  "anonymized",
			Name:   "Anonymous",
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/users/deactivate", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersDeactivateOption{
		Login:     "myuser",
		Anonymize: true,
	}

	result, resp, err := client.Users.Deactivate(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestUsers_Deactivate_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Users.Deactivate(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, _, err = client.Users.Deactivate(&UsersDeactivateOption{})
	assert.Error(t, err)
}

func TestUsers_DismissNotice(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/dismiss_notice", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersDismissNoticeOption{
		Notice: "educationPrinciples",
	}

	resp, err := client.Users.DismissNotice(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_DismissNotice_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersDismissNoticeOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing notice",
			opt:  &UsersDismissNoticeOption{},
		},
		{
			name: "invalid notice",
			opt:  &UsersDismissNoticeOption{Notice: "invalidNotice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.DismissNotice(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_Groups(t *testing.T) {
	response := UsersGroups{
		Groups: []UserGroup{
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
		Paging: UsersPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     2,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersGroupsOption{
		Login: "myuser",
	}

	result, resp, err := client.Users.Groups(opt)
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
		Groups: []UserGroup{},
		Paging: UsersPaging{
			PageIndex: 2,
			PageSize:  10,
			Total:     2,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersGroupsOption{
		Login: "myuser",
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 10,
		},
	}

	result, resp, err := client.Users.Groups(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(2), result.Paging.PageIndex)
}

func TestUsers_Groups_WithFilter(t *testing.T) {
	response := UsersGroups{
		Groups: []UserGroup{},
		Paging: UsersPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersGroupsOption{
		Login:    "myuser",
		Q:        "admin",
		Selected: "selected",
	}

	_, resp, err := client.Users.Groups(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUsers_Groups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersGroupsOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersGroupsOption{},
		},
		{
			name: "invalid selected",
			opt:  &UsersGroupsOption{Login: "myuser", Selected: "invalid"},
		},
		{
			name: "invalid page size",
			opt:  &UsersGroupsOption{Login: "myuser", PaginationArgs: PaginationArgs{PageSize: 1000}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Users.Groups(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_IdentityProviders(t *testing.T) {
	response := UsersIdentityProviders{
		IdentityProviders: []IdentityProvider{
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

	result, resp, err := client.Users.IdentityProviders()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.IdentityProviders, 2)
	assert.Equal(t, "LDAP", result.IdentityProviders[0].Key)
	assert.Equal(t, "GitHub", result.IdentityProviders[1].Name)
}

func TestUsers_Search(t *testing.T) {
	response := UsersSearch{
		Paging: UsersPaging{
			PageIndex: 1,
			PageSize:  50,
			Total:     2,
		},
		Users: []SearchedUser{
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

	result, resp, err := client.Users.Search(nil)
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
		Paging: UsersPaging{
			PageIndex: 1,
			PageSize:  25,
			Total:     0,
		},
		Users: []SearchedUser{},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/users/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &UsersSearchOption{
		Deactivated: true,
		Q:           "test",
		PaginationArgs: PaginationArgs{
			Page:     1,
			PageSize: 25,
		},
	}

	_, resp, err := client.Users.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUsers_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid page size should fail validation.
	opt := &UsersSearchOption{
		PaginationArgs: PaginationArgs{PageSize: 1000},
	}
	_, _, err := client.Users.Search(opt)
	assert.Error(t, err)
}

func TestUsers_SetHomepage(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/set_homepage", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersSetHomepageOption{
		Type:      "PROJECT",
		Component: "my-project",
	}

	resp, err := client.Users.SetHomepage(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_SetHomepage_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersSetHomepageOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing type",
			opt:  &UsersSetHomepageOption{},
		},
		{
			name: "invalid type",
			opt:  &UsersSetHomepageOption{Type: "INVALID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.SetHomepage(tt.opt)
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

	opt := &UsersUpdateOption{
		Login: "myuser",
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	result, resp, err := client.Users.Update(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Name", result.User.Name)
}

func TestUsers_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Users.Update(nil)
	assert.Error(t, err)

	// Missing Login should fail validation.
	_, _, err = client.Users.Update(&UsersUpdateOption{})
	assert.Error(t, err)
}

func TestUsers_UpdateIdentityProvider(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/update_identity_provider", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersUpdateIdentityProviderOption{
		Login:               "myuser",
		NewExternalProvider: "github",
		NewExternalIdentity: "github-user",
	}

	resp, err := client.Users.UpdateIdentityProvider(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_UpdateIdentityProvider_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersUpdateIdentityProviderOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersUpdateIdentityProviderOption{NewExternalProvider: "github"},
		},
		{
			name: "missing newExternalProvider",
			opt:  &UsersUpdateIdentityProviderOption{Login: "myuser"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.UpdateIdentityProvider(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestUsers_UpdateLogin(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/users/update_login", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &UsersUpdateLoginOption{
		Login:    "oldlogin",
		NewLogin: "newlogin",
	}

	resp, err := client.Users.UpdateLogin(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsers_UpdateLogin_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersUpdateLoginOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing login",
			opt:  &UsersUpdateLoginOption{NewLogin: "newlogin"},
		},
		{
			name: "missing newLogin",
			opt:  &UsersUpdateLoginOption{Login: "oldlogin"},
		},
		{
			name: "newLogin too short",
			opt:  &UsersUpdateLoginOption{Login: "oldlogin", NewLogin: "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Users.UpdateLogin(tt.opt)
			assert.Error(t, err)
		})
	}
}
