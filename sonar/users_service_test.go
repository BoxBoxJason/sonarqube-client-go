package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsers_Anonymize(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/anonymize" {
			t.Errorf("expected path /api/users/anonymize, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "deactivated-user" {
			t.Errorf("expected login 'deactivated-user', got %s", login)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersAnonymizeOption{
		Login: "deactivated-user",
	}

	resp, err := client.Users.Anonymize(opt)
	if err != nil {
		t.Fatalf("Anonymize failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUsers_Anonymize_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Users.Anonymize(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, err = client.Users.Anonymize(&UsersAnonymizeOption{})
	if err == nil {
		t.Error("expected error for missing Login")
	}
}

func TestUsers_ChangePassword(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/change_password" {
			t.Errorf("expected path /api/users/change_password, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "myuser" {
			t.Errorf("expected login 'myuser', got %s", login)
		}

		password := r.URL.Query().Get("password")
		if password != "MyNewPassword123!" {
			t.Errorf("expected password 'MyNewPassword123!', got %s", password)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersChangePasswordOption{
		Login:    "myuser",
		Password: "MyNewPassword123!",
	}

	resp, err := client.Users.ChangePassword(opt)
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUsers_ChangePassword_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestUsers_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/create" {
			t.Errorf("expected path /api/users/create, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "newuser" {
			t.Errorf("expected login 'newuser', got %s", login)
		}

		name := r.URL.Query().Get("name")
		if name != "New User" {
			t.Errorf("expected name 'New User', got %s", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"user": {
				"active": true,
				"email": "newuser@example.com",
				"local": true,
				"login": "newuser",
				"name": "New User",
				"scmAccounts": ["scm1", "scm2"]
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersCreateOption{
		Login:       "newuser",
		Name:        "New User",
		Email:       "newuser@example.com",
		Password:    "SecurePassword123!",
		Local:       true,
		ScmAccounts: []string{"scm1", "scm2"},
	}

	result, resp, err := client.Users.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.User.Login != "newuser" {
		t.Errorf("expected login 'newuser', got %s", result.User.Login)
	}

	if result.User.Name != "New User" {
		t.Errorf("expected name 'New User', got %s", result.User.Name)
	}

	if !result.User.Active {
		t.Error("expected user to be active")
	}

	if len(result.User.ScmAccounts) != 2 {
		t.Errorf("expected 2 SCM accounts, got %d", len(result.User.ScmAccounts))
	}
}

func TestUsers_Create_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestUsers_Current(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/current" {
			t.Errorf("expected path /api/users/current, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"avatar": "abc123",
			"dismissedNotices": {
				"educationPrinciples": true,
				"sonarlintAd": false
			},
			"email": "admin@example.com",
			"externalIdentity": "admin",
			"externalProvider": "sonarqube",
			"groups": ["sonar-users", "sonar-administrators"],
			"homepage": {
				"component": "my-project",
				"type": "PROJECT"
			},
			"id": "uuid-123",
			"isLoggedIn": true,
			"local": true,
			"login": "admin",
			"name": "Administrator",
			"permissions": {
				"global": ["admin", "scan"]
			},
			"scmAccounts": ["admin@scm.com"],
			"usingSonarLintConnectedMode": true
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Users.Current()
	if err != nil {
		t.Fatalf("Current failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Login != "admin" {
		t.Errorf("expected login 'admin', got %s", result.Login)
	}

	if !result.IsLoggedIn {
		t.Error("expected user to be logged in")
	}

	if result.Homepage.Type != "PROJECT" {
		t.Errorf("expected homepage type 'PROJECT', got %s", result.Homepage.Type)
	}

	if !result.DismissedNotices.EducationPrinciples {
		t.Error("expected educationPrinciples notice to be dismissed")
	}

	if len(result.Permissions.Global) != 2 {
		t.Errorf("expected 2 global permissions, got %d", len(result.Permissions.Global))
	}
}

func TestUsers_Deactivate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/deactivate" {
			t.Errorf("expected path /api/users/deactivate, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "myuser" {
			t.Errorf("expected login 'myuser', got %s", login)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"user": {
				"active": false,
				"groups": [],
				"local": true,
				"login": "myuser",
				"name": "My User",
				"scmAccounts": []
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersDeactivateOption{
		Login:     "myuser",
		Anonymize: false,
	}

	result, resp, err := client.Users.Deactivate(opt)
	if err != nil {
		t.Fatalf("Deactivate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.User.Active {
		t.Error("expected user to be inactive")
	}

	if result.User.Login != "myuser" {
		t.Errorf("expected login 'myuser', got %s", result.User.Login)
	}
}

func TestUsers_Deactivate_WithAnonymize(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anonymize := r.URL.Query().Get("anonymize")
		if anonymize != "true" {
			t.Errorf("expected anonymize 'true', got %s", anonymize)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"user": {
				"active": false,
				"groups": [],
				"local": true,
				"login": "anonymized",
				"name": "Anonymous",
				"scmAccounts": []
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersDeactivateOption{
		Login:     "myuser",
		Anonymize: true,
	}

	result, resp, err := client.Users.Deactivate(opt)
	if err != nil {
		t.Fatalf("Deactivate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestUsers_Deactivate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Users.Deactivate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, _, err = client.Users.Deactivate(&UsersDeactivateOption{})
	if err == nil {
		t.Error("expected error for missing Login")
	}
}

func TestUsers_DismissNotice(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/dismiss_notice" {
			t.Errorf("expected path /api/users/dismiss_notice, got %s", r.URL.Path)
		}

		notice := r.URL.Query().Get("notice")
		if notice != "educationPrinciples" {
			t.Errorf("expected notice 'educationPrinciples', got %s", notice)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersDismissNoticeOption{
		Notice: "educationPrinciples",
	}

	resp, err := client.Users.DismissNotice(opt)
	if err != nil {
		t.Fatalf("DismissNotice failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUsers_DismissNotice_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestUsers_Groups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/groups" {
			t.Errorf("expected path /api/users/groups, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "myuser" {
			t.Errorf("expected login 'myuser', got %s", login)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"groups": [
				{
					"default": true,
					"description": "Default group",
					"id": 1,
					"name": "sonar-users",
					"selected": true
				},
				{
					"default": false,
					"description": "Administrators",
					"id": 2,
					"name": "sonar-administrators",
					"selected": false
				}
			],
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 2
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersGroupsOption{
		Login: "myuser",
	}

	result, resp, err := client.Users.Groups(opt)
	if err != nil {
		t.Fatalf("Groups failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result.Groups))
	}

	if result.Groups[0].Name != "sonar-users" {
		t.Errorf("expected first group 'sonar-users', got %s", result.Groups[0].Name)
	}

	if !result.Groups[0].Default {
		t.Error("expected first group to be default")
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Paging.Total)
	}
}

func TestUsers_Groups_WithPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get("p")
		if p != "2" {
			t.Errorf("expected p '2', got %s", p)
		}

		ps := r.URL.Query().Get("ps")
		if ps != "10" {
			t.Errorf("expected ps '10', got %s", ps)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"groups": [],
			"paging": {
				"pageIndex": 2,
				"pageSize": 10,
				"total": 2
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersGroupsOption{
		Login: "myuser",
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 10,
		},
	}

	result, resp, err := client.Users.Groups(opt)
	if err != nil {
		t.Fatalf("Groups failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Paging.PageIndex != 2 {
		t.Errorf("expected pageIndex 2, got %d", result.Paging.PageIndex)
	}
}

func TestUsers_Groups_WithFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		selected := r.URL.Query().Get("selected")
		if selected != "selected" {
			t.Errorf("expected selected 'selected', got %s", selected)
		}

		q := r.URL.Query().Get("q")
		if q != "admin" {
			t.Errorf("expected q 'admin', got %s", q)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"groups": [],
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 0
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersGroupsOption{
		Login:    "myuser",
		Q:        "admin",
		Selected: "selected",
	}

	_, resp, err := client.Users.Groups(opt)
	if err != nil {
		t.Fatalf("Groups failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestUsers_Groups_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestUsers_IdentityProviders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/identity_providers" {
			t.Errorf("expected path /api/users/identity_providers, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"identityProviders": [
				{
					"backgroundColor": "#444444",
					"helpMessage": "Use your LDAP credentials",
					"iconPath": "/images/ldap.png",
					"key": "LDAP",
					"name": "LDAP"
				},
				{
					"backgroundColor": "#000000",
					"helpMessage": "Use your GitHub account",
					"iconPath": "/images/github.png",
					"key": "github",
					"name": "GitHub"
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Users.IdentityProviders()
	if err != nil {
		t.Fatalf("IdentityProviders failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.IdentityProviders) != 2 {
		t.Errorf("expected 2 identity providers, got %d", len(result.IdentityProviders))
	}

	if result.IdentityProviders[0].Key != "LDAP" {
		t.Errorf("expected first provider key 'LDAP', got %s", result.IdentityProviders[0].Key)
	}

	if result.IdentityProviders[1].Name != "GitHub" {
		t.Errorf("expected second provider name 'GitHub', got %s", result.IdentityProviders[1].Name)
	}
}

func TestUsers_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/search" {
			t.Errorf("expected path /api/users/search, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 50,
				"total": 2
			},
			"users": [
				{
					"active": true,
					"avatar": "abc123",
					"email": "admin@example.com",
					"externalIdentity": "admin",
					"externalProvider": "sonarqube",
					"groups": ["sonar-users", "sonar-administrators"],
					"lastConnectionDate": "2024-01-01T00:00:00+0000",
					"local": true,
					"login": "admin",
					"managed": false,
					"name": "Administrator",
					"scmAccounts": ["admin@scm.com"],
					"sonarLintLastConnectionDate": "2024-01-01T00:00:00+0000",
					"tokensCount": 3
				},
				{
					"active": true,
					"login": "user1",
					"name": "User One",
					"local": false,
					"managed": true
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Users.Search(nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(result.Users))
	}

	if result.Users[0].Login != "admin" {
		t.Errorf("expected first user login 'admin', got %s", result.Users[0].Login)
	}

	if result.Users[0].TokensCount != 3 {
		t.Errorf("expected tokens count 3, got %d", result.Users[0].TokensCount)
	}

	if !result.Users[1].Managed {
		t.Error("expected second user to be managed")
	}
}

func TestUsers_Search_WithFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deactivated := r.URL.Query().Get("deactivated")
		if deactivated != "true" {
			t.Errorf("expected deactivated 'true', got %s", deactivated)
		}

		q := r.URL.Query().Get("q")
		if q != "test" {
			t.Errorf("expected q 'test', got %s", q)
		}

		p := r.URL.Query().Get("p")
		if p != "1" {
			t.Errorf("expected p '1', got %s", p)
		}

		ps := r.URL.Query().Get("ps")
		if ps != "25" {
			t.Errorf("expected ps '25', got %s", ps)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 25,
				"total": 0
			},
			"users": []
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersSearchOption{
		Deactivated: true,
		Q:           "test",
		PaginationArgs: PaginationArgs{
			Page:     1,
			PageSize: 25,
		},
	}

	_, resp, err := client.Users.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestUsers_Search_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Invalid page size should fail validation.
	opt := &UsersSearchOption{
		PaginationArgs: PaginationArgs{PageSize: 1000},
	}
	_, _, err := client.Users.Search(opt)
	if err == nil {
		t.Error("expected error for invalid page size")
	}
}

func TestUsers_SetHomepage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/set_homepage" {
			t.Errorf("expected path /api/users/set_homepage, got %s", r.URL.Path)
		}

		typeParam := r.URL.Query().Get("type")
		if typeParam != "PROJECT" {
			t.Errorf("expected type 'PROJECT', got %s", typeParam)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersSetHomepageOption{
		Type:      "PROJECT",
		Component: "my-project",
	}

	resp, err := client.Users.SetHomepage(opt)
	if err != nil {
		t.Fatalf("SetHomepage failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUsers_SetHomepage_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestUsers_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/update" {
			t.Errorf("expected path /api/users/update, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "myuser" {
			t.Errorf("expected login 'myuser', got %s", login)
		}

		name := r.URL.Query().Get("name")
		if name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"user": {
				"active": true,
				"email": "updated@example.com",
				"local": true,
				"login": "myuser",
				"name": "Updated Name",
				"scmAccounts": ["scm1"]
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersUpdateOption{
		Login: "myuser",
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	result, resp, err := client.Users.Update(opt)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.User.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %s", result.User.Name)
	}
}

func TestUsers_Update_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Users.Update(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Login should fail validation.
	_, _, err = client.Users.Update(&UsersUpdateOption{})
	if err == nil {
		t.Error("expected error for missing Login")
	}
}

func TestUsers_UpdateIdentityProvider(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/update_identity_provider" {
			t.Errorf("expected path /api/users/update_identity_provider, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "myuser" {
			t.Errorf("expected login 'myuser', got %s", login)
		}

		newExternalProvider := r.URL.Query().Get("newExternalProvider")
		if newExternalProvider != "github" {
			t.Errorf("expected newExternalProvider 'github', got %s", newExternalProvider)
		}

		newExternalIdentity := r.URL.Query().Get("newExternalIdentity")
		if newExternalIdentity != "github-user" {
			t.Errorf("expected newExternalIdentity 'github-user', got %s", newExternalIdentity)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersUpdateIdentityProviderOption{
		Login:               "myuser",
		NewExternalProvider: "github",
		NewExternalIdentity: "github-user",
	}

	resp, err := client.Users.UpdateIdentityProvider(opt)
	if err != nil {
		t.Fatalf("UpdateIdentityProvider failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUsers_UpdateIdentityProvider_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestUsers_UpdateLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/users/update_login" {
			t.Errorf("expected path /api/users/update_login, got %s", r.URL.Path)
		}

		login := r.URL.Query().Get("login")
		if login != "oldlogin" {
			t.Errorf("expected login 'oldlogin', got %s", login)
		}

		newLogin := r.URL.Query().Get("newLogin")
		if newLogin != "newlogin" {
			t.Errorf("expected newLogin 'newlogin', got %s", newLogin)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &UsersUpdateLoginOption{
		Login:    "oldlogin",
		NewLogin: "newlogin",
	}

	resp, err := client.Users.UpdateLogin(opt)
	if err != nil {
		t.Fatalf("UpdateLogin failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUsers_UpdateLogin_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}
