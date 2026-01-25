package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifications_Add(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/notifications/add" {
			t.Errorf("expected path /api/notifications/add, got %s", r.URL.Path)
		}

		notifType := r.URL.Query().Get("type")
		if notifType != "NewIssues" {
			t.Errorf("expected type 'NewIssues', got %s", notifType)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NotificationsAddOption{
		Type: "NewIssues",
	}

	resp, err := client.Notifications.Add(opt)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestNotifications_Add_WithAllOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notifType := r.URL.Query().Get("type")
		if notifType != "NewIssues" {
			t.Errorf("expected type 'NewIssues', got %s", notifType)
		}

		channel := r.URL.Query().Get("channel")
		if channel != "email" {
			t.Errorf("expected channel 'email', got %s", channel)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		login := r.URL.Query().Get("login")
		if login != "admin" {
			t.Errorf("expected login 'admin', got %s", login)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NotificationsAddOption{
		Type:    "NewIssues",
		Channel: "email",
		Project: "my-project",
		Login:   "admin",
	}

	resp, err := client.Notifications.Add(opt)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestNotifications_Add_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Notifications.Add(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Type should fail validation.
	_, err = client.Notifications.Add(&NotificationsAddOption{})
	if err == nil {
		t.Error("expected error for missing Type")
	}
}

func TestNotifications_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/notifications/list" {
			t.Errorf("expected path /api/notifications/list, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"channels": ["email"],
			"globalTypes": ["NewIssues", "NewAlerts"],
			"perProjectTypes": ["NewIssues"],
			"notifications": [
				{
					"type": "NewIssues",
					"channel": "email",
					"project": "my-project",
					"projectName": "My Project"
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Notifications.List(nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Channels) != 1 {
		t.Errorf("expected 1 channel, got %d", len(result.Channels))
	}

	if result.Channels[0] != "email" {
		t.Errorf("expected channel 'email', got %s", result.Channels[0])
	}

	if len(result.GlobalTypes) != 2 {
		t.Errorf("expected 2 global types, got %d", len(result.GlobalTypes))
	}

	if len(result.Notifications) != 1 {
		t.Errorf("expected 1 notification, got %d", len(result.Notifications))
	}

	if result.Notifications[0].Type != "NewIssues" {
		t.Errorf("expected notification type 'NewIssues', got %s", result.Notifications[0].Type)
	}
}

func TestNotifications_List_WithLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := r.URL.Query().Get("login")
		if login != "admin" {
			t.Errorf("expected login 'admin', got %s", login)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"channels": [], "globalTypes": [], "perProjectTypes": [], "notifications": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NotificationsListOption{
		Login: "admin",
	}

	_, resp, err := client.Notifications.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestNotifications_Remove(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/notifications/remove" {
			t.Errorf("expected path /api/notifications/remove, got %s", r.URL.Path)
		}

		notifType := r.URL.Query().Get("type")
		if notifType != "NewIssues" {
			t.Errorf("expected type 'NewIssues', got %s", notifType)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NotificationsRemoveOption{
		Type: "NewIssues",
	}

	resp, err := client.Notifications.Remove(opt)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestNotifications_Remove_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Notifications.Remove(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Type should fail validation.
	_, err = client.Notifications.Remove(&NotificationsRemoveOption{})
	if err == nil {
		t.Error("expected error for missing Type")
	}
}

func TestNotifications_ValidateAddOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.Notifications.ValidateAddOpt(&NotificationsAddOption{
		Type: "NewIssues",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.Notifications.ValidateAddOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Type should fail.
	err = client.Notifications.ValidateAddOpt(&NotificationsAddOption{})
	if err == nil {
		t.Error("expected error for missing Type")
	}
}

func TestNotifications_ValidateListOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid.
	err := client.Notifications.ValidateListOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid.
	err = client.Notifications.ValidateListOpt(&NotificationsListOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}
}

func TestNotifications_ValidateRemoveOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.Notifications.ValidateRemoveOpt(&NotificationsRemoveOption{
		Type: "NewIssues",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.Notifications.ValidateRemoveOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Type should fail.
	err = client.Notifications.ValidateRemoveOpt(&NotificationsRemoveOption{})
	if err == nil {
		t.Error("expected error for missing Type")
	}
}
