package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifications_Add(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/notifications/add", r.URL.Path)
		assert.Equal(t, "NewIssues", r.URL.Query().Get("type"))
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	opt := &NotificationsAddOption{
		Type: "NewIssues",
	}

	resp, err := client.Notifications.Add(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestNotifications_Add_WithAllOptions(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "NewIssues", r.URL.Query().Get("type"))
		assert.Equal(t, "email", r.URL.Query().Get("channel"))
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		assert.Equal(t, "admin", r.URL.Query().Get("login"))
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	opt := &NotificationsAddOption{
		Type:    "NewIssues",
		Channel: "email",
		Project: "my-project",
		Login:   "admin",
	}

	resp, err := client.Notifications.Add(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestNotifications_Add_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Notifications.Add(nil)
	assert.Error(t, err)

	// Missing Type should fail validation.
	_, err = client.Notifications.Add(&NotificationsAddOption{})
	assert.Error(t, err)
}

func TestNotifications_List(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/notifications/list", http.StatusOK, &NotificationsList{
		Channels:        []string{"email"},
		GlobalTypes:     []string{"NewIssues", "NewAlerts"},
		PerProjectTypes: []string{"NewIssues"},
		Notifications: []Notification{
			{
				Type:        "NewIssues",
				Channel:     "email",
				Project:     "my-project",
				ProjectName: "My Project",
			},
		},
	}))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Notifications.List(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Channels, 1)
	assert.Equal(t, "email", result.Channels[0])
	assert.Len(t, result.GlobalTypes, 2)
	assert.Len(t, result.Notifications, 1)
	assert.Equal(t, "NewIssues", result.Notifications[0].Type)
}

func TestNotifications_List_WithLogin(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "admin", r.URL.Query().Get("login"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"channels": [], "globalTypes": [], "perProjectTypes": [], "notifications": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &NotificationsListOption{
		Login: "admin",
	}

	_, resp, err := client.Notifications.List(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNotifications_Remove(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/notifications/remove", r.URL.Path)
		assert.Equal(t, "NewIssues", r.URL.Query().Get("type"))
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	opt := &NotificationsRemoveOption{
		Type: "NewIssues",
	}

	resp, err := client.Notifications.Remove(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestNotifications_Remove_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Notifications.Remove(nil)
	assert.Error(t, err)

	// Missing Type should fail validation.
	_, err = client.Notifications.Remove(&NotificationsRemoveOption{})
	assert.Error(t, err)
}

func TestNotifications_ValidateAddOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.Notifications.ValidateAddOpt(&NotificationsAddOption{
		Type: "NewIssues",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.Notifications.ValidateAddOpt(nil)
	assert.Error(t, err)

	// Missing Type should fail.
	err = client.Notifications.ValidateAddOpt(&NotificationsAddOption{})
	assert.Error(t, err)
}

func TestNotifications_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should be valid.
	err := client.Notifications.ValidateListOpt(nil)
	assert.NoError(t, err)

	// Empty option should be valid.
	err = client.Notifications.ValidateListOpt(&NotificationsListOption{})
	assert.NoError(t, err)
}

func TestNotifications_ValidateRemoveOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.Notifications.ValidateRemoveOpt(&NotificationsRemoveOption{
		Type: "NewIssues",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.Notifications.ValidateRemoveOpt(nil)
	assert.Error(t, err)

	// Missing Type should fail.
	err = client.Notifications.ValidateRemoveOpt(&NotificationsRemoveOption{})
	assert.Error(t, err)
}
