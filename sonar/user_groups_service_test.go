package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserGroups_AddUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "user_groups/add_user") {
			t.Errorf("expected path to contain user_groups/add_user, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsAddUserOption{
		Name:  "sonar-administrators",
		Login: "g.hopper",
	}

	resp, err := client.UserGroups.AddUser(opt)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUserGroups_AddUser_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.UserGroups.AddUser(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Name
	_, err = client.UserGroups.AddUser(&UserGroupsAddUserOption{
		Login: "user",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}
}

func TestUserGroups_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := UserGroupsCreate{
			Group: UserGroupDetail{
				ID:           "uuid-group-1",
				Name:         "sonar-users",
				Description:  "Default group",
				MembersCount: 0,
				Default:      false,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsCreateOption{
		Name:        "sonar-users",
		Description: "Default group",
	}

	result, resp, err := client.UserGroups.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Group.Name != "sonar-users" {
		t.Errorf("expected name sonar-users, got %s", result.Group.Name)
	}
}

func TestUserGroups_Create_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.UserGroups.Create(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Name
	_, _, err = client.UserGroups.Create(&UserGroupsCreateOption{
		Description: "test",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Test Name too long
	_, _, err = client.UserGroups.Create(&UserGroupsCreateOption{
		Name: strings.Repeat("a", MaxGroupNameLength+1),
	})
	if err == nil {
		t.Error("expected error for Name exceeding max length")
	}

	// Test Description too long
	_, _, err = client.UserGroups.Create(&UserGroupsCreateOption{
		Name:        "test",
		Description: strings.Repeat("a", MaxGroupDescriptionLength+1),
	})
	if err == nil {
		t.Error("expected error for Description exceeding max length")
	}
}

func TestUserGroups_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsDeleteOption{
		Name: "sonar-users",
	}

	resp, err := client.UserGroups.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUserGroups_Delete_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.UserGroups.Delete(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Name
	_, err = client.UserGroups.Delete(&UserGroupsDeleteOption{})
	if err == nil {
		t.Error("expected error for missing Name")
	}
}

func TestUserGroups_RemoveUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsRemoveUserOption{
		Name:  "sonar-administrators",
		Login: "g.hopper",
	}

	resp, err := client.UserGroups.RemoveUser(opt)
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUserGroups_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsSearchOption{
		Query:  "admin",
		Fields: []string{"name", "description"},
	}

	result, resp, err := client.UserGroups.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(result.Groups))
	}
}

func TestUserGroups_Search_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.UserGroups.Search(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test invalid field
	_, _, err = client.UserGroups.Search(&UserGroupsSearchOption{
		Fields: []string{"invalid_field"},
	})
	if err == nil {
		t.Error("expected error for invalid field")
	}
}

func TestUserGroups_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsUpdateOption{
		CurrentName: "old-group",
		Name:        "new-group",
		Description: "Updated description",
	}

	resp, err := client.UserGroups.Update(opt)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUserGroups_Update_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.UserGroups.Update(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing CurrentName
	_, err = client.UserGroups.Update(&UserGroupsUpdateOption{
		Name: "new-group",
	})
	if err == nil {
		t.Error("expected error for missing CurrentName")
	}
}

func TestUserGroups_Users(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UserGroupsUsersOption{
		Name:     "sonar-administrators",
		Selected: "selected",
	}

	result, resp, err := client.UserGroups.Users(opt)
	if err != nil {
		t.Fatalf("Users failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Users) != 1 {
		t.Errorf("expected 1 user, got %d", len(result.Users))
	}
}

func TestUserGroups_Users_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.UserGroups.Users(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Name
	_, _, err = client.UserGroups.Users(&UserGroupsUsersOption{})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Test invalid Selected value
	_, _, err = client.UserGroups.Users(&UserGroupsUsersOption{
		Name:     "test",
		Selected: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected value")
	}
}
