package sonar

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Search
// =============================================================================

func TestUsersV2_Search(t *testing.T) {
	response := UsersSearchV2{
		Users: []UserV2{
			{Id: "user-1", Login: "jdoe", Name: "John Doe", Email: "jdoe@example.com", Active: true},
			{Id: "user-2", Login: "asmith", Name: "Alice Smith", Active: true, Local: true},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 2},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/users-management/users", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.UsersManagement.Search(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, "jdoe", result.Users[0].Login)
}

func TestUsersV2_Search_WithOptions(t *testing.T) {
	active := true
	response := UsersSearchV2{
		Users: []UserV2{{Id: "user-1", Login: "jdoe", Name: "John Doe", Active: true}},
		Page:  PageResponseV2{PageIndex: 1, PageSize: 10, Total: 1},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/users-management/users", http.StatusOK,
		map[string]string{"active": "true", "q": "jdoe", "pageSize": "10"},
		response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.UsersManagement.Search(&UsersSearchOptionV2{
		PaginationParamsV2: PaginationParamsV2{PageSize: 10},
		Active:             &active,
		Query:              "jdoe",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Users, 1)
}

func TestUsersV2_Search_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	// Invalid page size
	_, _, err := client.V2.UsersManagement.Search(&UsersSearchOptionV2{
		PaginationParamsV2: PaginationParamsV2{PageSize: 600},
	})
	assert.Error(t, err)
}

// =============================================================================
// Create
// =============================================================================

func TestUsersV2_Create(t *testing.T) {
	response := UserV2{
		Id:    "user-new",
		Login: "newuser",
		Name:  "New User",
		Email: "new@example.com",
		Local: true,
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/users-management/users", http.StatusOK,
		&UsersCreateOptionsV2{
			Login:    "newuser",
			Name:     "New User",
			Email:    "new@example.com",
			Password: "secret123",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.UsersManagement.Create(&UsersCreateOptionsV2{
		Login:    "newuser",
		Name:     "New User",
		Email:    "new@example.com",
		Password: "secret123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "newuser", result.Login)
}

func TestUsersV2_Create_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersCreateOptionsV2
	}{
		{"nil opt", nil},
		{"missing login", &UsersCreateOptionsV2{Name: "Test"}},
		{"missing name", &UsersCreateOptionsV2{Login: "test"}},
		{"login too short", &UsersCreateOptionsV2{Login: "a", Name: "Test"}},
		{"login too long", &UsersCreateOptionsV2{Login: strings.Repeat("a", MaxLoginLengthV2+1), Name: "Test"}},
		{"name too long", &UsersCreateOptionsV2{Login: "test", Name: strings.Repeat("a", MaxNameLength+1)}},
		{"email too long", &UsersCreateOptionsV2{Login: "test", Name: "Test", Email: strings.Repeat("a", MaxEmailLength+1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.UsersManagement.Create(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Fetch
// =============================================================================

func TestUsersV2_Fetch(t *testing.T) {
	response := UserV2{
		Id:    "user-1",
		Login: "jdoe",
		Name:  "John Doe",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/users-management/users/user-1", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.UsersManagement.Fetch("user-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "jdoe", result.Login)
}

func TestUsersV2_Fetch_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.UsersManagement.Fetch("")
	assert.Error(t, err)
}

// =============================================================================
// Deactivate
// =============================================================================

func TestUsersV2_Deactivate(t *testing.T) {
	server := newTestServer(t, mockEmptyHandlerWithParams(t, http.MethodDelete, "/v2/users-management/users/user-1", http.StatusNoContent,
		map[string]string{"anonymize": "true"}))
	client := newTestClient(t, server.url())

	resp, err := client.V2.UsersManagement.Deactivate(&UsersDeactivateOptionsV2{
		Id:        "user-1",
		Anonymize: true,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestUsersV2_Deactivate_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *UsersDeactivateOptionsV2
	}{
		{"nil opt", nil},
		{"missing id", &UsersDeactivateOptionsV2{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.V2.UsersManagement.Deactivate(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Update
// =============================================================================

func TestUsersV2_Update(t *testing.T) {
	response := UserV2{
		Id:    "user-1",
		Login: "jdoe",
		Name:  "John Updated",
		Email: "updated@example.com",
	}
	server := newTestServer(t, mockPatchHandler(t, "/v2/users-management/users/user-1", http.StatusOK,
		&UsersUpdateOptionsV2{
			Name:  "John Updated",
			Email: "updated@example.com",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.UsersManagement.Update("user-1", &UsersUpdateOptionsV2{
		Name:  "John Updated",
		Email: "updated@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "John Updated", result.Name)
}

func TestUsersV2_Update_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		id   string
		opt  *UsersUpdateOptionsV2
	}{
		{"missing id", "", &UsersUpdateOptionsV2{}},
		{"nil opt", "user-1", nil},
		{"login too long", "user-1", &UsersUpdateOptionsV2{Login: strings.Repeat("a", MaxLoginLengthV2+1)}},
		{"name too long", "user-1", &UsersUpdateOptionsV2{Name: strings.Repeat("a", MaxNameLength+1)}},
		{"email too long", "user-1", &UsersUpdateOptionsV2{Email: strings.Repeat("a", MaxEmailLength+1)}},
		{"external login too long", "user-1", &UsersUpdateOptionsV2{ExternalLogin: strings.Repeat("a", MaxExternalLoginLengthV2+1)}},
		{"external id too long", "user-1", &UsersUpdateOptionsV2{ExternalId: strings.Repeat("a", MaxExternalIdLengthV2+1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.UsersManagement.Update(tt.id, tt.opt)
			assert.Error(t, err)
		})
	}
}
