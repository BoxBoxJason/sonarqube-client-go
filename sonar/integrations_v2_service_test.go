package sonar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// CreateUserBinding
// =============================================================================

func TestIntegrationsV2_CreateUserBinding(t *testing.T) {
	response := IntegrationsUserBinding{
		Id:                 "binding-1",
		UserId:             "user-1",
		SlackUserId:        "U123",
		SlackWorkspaceId:   "W123",
		SlackWorkspaceName: "my-workspace",
		CreatedAt:          1700000000000,
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/integrations/user-bindings", http.StatusCreated,
		&IntegrationsUserBindingCreateOptions{
			UserId:      "user-1",
			BindingData: IntegrationsUserBindingData{Code: "auth-code"},
		}, response))
	client := newTestClient(t, server.url())
	svc := &IntegrationsService{client: client}

	result, resp, err := svc.CreateUserBinding(context.Background(), &IntegrationsUserBindingCreateOptions{
		UserId:      "user-1",
		BindingData: IntegrationsUserBindingData{Code: "auth-code"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "binding-1", result.Id)
	assert.Equal(t, "U123", result.SlackUserId)
	assert.Equal(t, int64(1700000000000), result.CreatedAt)
}

func TestIntegrationsV2_CreateUserBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	tests := []struct {
		opt  *IntegrationsUserBindingCreateOptions
		name string
	}{
		{nil, "nil opt"},
		{&IntegrationsUserBindingCreateOptions{BindingData: IntegrationsUserBindingData{Code: "c"}}, "missing UserId"},
		{&IntegrationsUserBindingCreateOptions{UserId: "u"}, "missing BindingData.Code"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.CreateUserBinding(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// GetUserBinding
// =============================================================================

func TestIntegrationsV2_GetUserBinding(t *testing.T) {
	response := IntegrationsUserBinding{
		Id:     "binding-1",
		UserId: "user-1",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/integrations/user-bindings/binding-1", http.StatusOK, response))
	client := newTestClient(t, server.url())
	svc := &IntegrationsService{client: client}

	result, resp, err := svc.GetUserBinding(context.Background(), "binding-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "binding-1", result.Id)
}

func TestIntegrationsV2_GetUserBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	_, _, err := svc.GetUserBinding(context.Background(), "")
	assert.Error(t, err)
}

// =============================================================================
// HandleSlackSlashCommand
// =============================================================================

func TestIntegrationsV2_HandleSlackSlashCommand(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "unexpected HTTP method")
		assert.Equal(t, "/v2/integrations/slack/slash-commands", r.URL.Path, "unexpected URL path")
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"), "unexpected Content-Type")

		require.NoError(t, r.ParseForm())
		assert.Equal(t, "/sonarqube-server", r.FormValue("command"))
		assert.Equal(t, "connect", r.FormValue("text"))
		assert.Equal(t, "U123", r.FormValue("user_id"))

		w.WriteHeader(http.StatusOK)
	})
	client := newTestClient(t, server.URL)
	svc := &IntegrationsService{client: client}

	form := url.Values{
		"command": {"/sonarqube-server"},
		"text":    {"connect"},
		"user_id": {"U123"},
	}

	resp, err := svc.HandleSlackSlashCommand(context.Background(), form)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIntegrationsV2_HandleSlackSlashCommand_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	tests := []struct {
		form url.Values
		name string
	}{
		{nil, "nil form"},
		{url.Values{}, "empty form"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.HandleSlackSlashCommand(context.Background(), tt.form)
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	}
}

// =============================================================================
// HandleSlackEvent
// =============================================================================

func TestIntegrationsV2_HandleSlackEvent(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "unexpected HTTP method")
		assert.Equal(t, "/v2/integrations/slack/events", r.URL.Path, "unexpected URL path")

		var actual map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&actual))
		assert.Equal(t, "url_verification", actual["type"])
		assert.Equal(t, "some-challenge", actual["challenge"])

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("some-challenge"))
	})
	client := newTestClient(t, server.URL)
	svc := &IntegrationsService{client: client}

	payload := json.RawMessage(`{"type":"url_verification","challenge":"some-challenge"}`)

	result, resp, err := svc.HandleSlackEvent(context.Background(), payload)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "some-challenge", *result)
}

func TestIntegrationsV2_HandleSlackEvent_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	tests := []struct {
		name    string
		payload json.RawMessage
	}{
		{"nil payload", nil},
		{"empty payload", json.RawMessage{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.HandleSlackEvent(context.Background(), tt.payload)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// ListIntegrationConfigurations
// =============================================================================

func TestIntegrationsV2_ListIntegrationConfigurations(t *testing.T) {
	response := IntegrationsIntegrationConfigurationSearch{
		IntegrationConfigurations: []IntegrationsIntegrationConfiguration{
			{Id: "cfg-1", IntegrationType: IntegrationTypeSlack, ClientId: "client-1", AppId: "app-1"},
		},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/integrations/integration-configurations", http.StatusOK,
		map[string]string{"integrationType": "SLACK"},
		response))
	client := newTestClient(t, server.url())
	svc := &IntegrationsService{client: client}

	result, resp, err := svc.ListIntegrationConfigurations(context.Background(), &IntegrationsListIntegrationConfigurationsOptions{
		IntegrationType: IntegrationTypeSlack,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.IntegrationConfigurations, 1)
	assert.Equal(t, "cfg-1", result.IntegrationConfigurations[0].Id)
}

func TestIntegrationsV2_ListIntegrationConfigurations_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	tests := []struct {
		opt  *IntegrationsListIntegrationConfigurationsOptions
		name string
	}{
		{nil, "nil opt"},
		{&IntegrationsListIntegrationConfigurationsOptions{}, "missing IntegrationType"},
		{&IntegrationsListIntegrationConfigurationsOptions{IntegrationType: "TEAMS"}, "invalid IntegrationType"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.ListIntegrationConfigurations(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// CreateIntegrationConfiguration
// =============================================================================

func TestIntegrationsV2_CreateIntegrationConfiguration(t *testing.T) {
	response := IntegrationsIntegrationConfiguration{
		Id:              "cfg-1",
		IntegrationType: IntegrationTypeSlack,
		ClientId:        "client-1",
		AppId:           "app-1",
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/integrations/integration-configurations", http.StatusCreated,
		&IntegrationsIntegrationConfigurationCreateOptions{
			IntegrationType: IntegrationTypeSlack,
			ClientId:        "client-1",
			ClientSecret:    "client-secret",
			SigningSecret:   "signing-secret",
		}, response))
	client := newTestClient(t, server.url())
	svc := &IntegrationsService{client: client}

	result, resp, err := svc.CreateIntegrationConfiguration(context.Background(), &IntegrationsIntegrationConfigurationCreateOptions{
		IntegrationType: IntegrationTypeSlack,
		ClientId:        "client-1",
		ClientSecret:    "client-secret",
		SigningSecret:   "signing-secret",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "cfg-1", result.Id)
	assert.Equal(t, "app-1", result.AppId)
}

func TestIntegrationsV2_CreateIntegrationConfiguration_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	base := IntegrationsIntegrationConfigurationCreateOptions{
		IntegrationType: IntegrationTypeSlack,
		ClientId:        "client-1",
		ClientSecret:    "client-secret",
		SigningSecret:   "signing-secret",
	}

	tests := []struct {
		opt  *IntegrationsIntegrationConfigurationCreateOptions
		name string
	}{
		{nil, "nil opt"},
		{func() *IntegrationsIntegrationConfigurationCreateOptions {
			o := base
			o.IntegrationType = ""
			return &o
		}(), "missing IntegrationType"},
		{func() *IntegrationsIntegrationConfigurationCreateOptions {
			o := base
			o.IntegrationType = "TEAMS"
			return &o
		}(), "invalid IntegrationType"},
		{func() *IntegrationsIntegrationConfigurationCreateOptions { o := base; o.ClientId = ""; return &o }(), "missing ClientId"},
		{func() *IntegrationsIntegrationConfigurationCreateOptions { o := base; o.ClientSecret = ""; return &o }(), "missing ClientSecret"},
		{func() *IntegrationsIntegrationConfigurationCreateOptions { o := base; o.SigningSecret = ""; return &o }(), "missing SigningSecret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.CreateIntegrationConfiguration(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// DeleteIntegrationConfiguration
// =============================================================================

func TestIntegrationsV2_DeleteIntegrationConfiguration(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/integrations/integration-configurations/cfg-1", http.StatusNoContent))
	client := newTestClient(t, server.url())
	svc := &IntegrationsService{client: client}

	resp, err := svc.DeleteIntegrationConfiguration(context.Background(), "cfg-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestIntegrationsV2_DeleteIntegrationConfiguration_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	resp, err := svc.DeleteIntegrationConfiguration(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// =============================================================================
// UpdateIntegrationConfiguration
// =============================================================================

func TestIntegrationsV2_UpdateIntegrationConfiguration(t *testing.T) {
	response := IntegrationsIntegrationConfiguration{
		Id:              "cfg-1",
		IntegrationType: IntegrationTypeSlack,
		ClientId:        "new-client-id",
		AppId:           "app-1",
	}
	server := newTestServer(t, mockPatchHandler(t, "/v2/integrations/integration-configurations/cfg-1", http.StatusOK,
		&IntegrationsIntegrationConfigurationPatchOptions{
			ClientId: "new-client-id",
		}, response))
	client := newTestClient(t, server.url())
	svc := &IntegrationsService{client: client}

	result, resp, err := svc.UpdateIntegrationConfiguration(context.Background(), "cfg-1", &IntegrationsIntegrationConfigurationPatchOptions{
		ClientId: "new-client-id",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "new-client-id", result.ClientId)
}

func TestIntegrationsV2_UpdateIntegrationConfiguration_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &IntegrationsService{client: client}

	tests := []struct {
		opt  *IntegrationsIntegrationConfigurationPatchOptions
		name string
		id   string
	}{
		{&IntegrationsIntegrationConfigurationPatchOptions{ClientId: "c"}, "missing id", ""},
		{nil, "nil opt", "cfg-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.UpdateIntegrationConfiguration(context.Background(), tt.id, tt.opt)
			assert.Error(t, err)
		})
	}
}
