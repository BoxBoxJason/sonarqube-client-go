package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardsService_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/dashboards/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Dashboards.Delete(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDashboardsService_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Dashboards.Delete(context.Background(), "")
	require.Error(t, err)
}

func TestDashboardsService_List(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dashboards", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Dashboards.List(context.Background(), &DashboardsListOptions{ResourceId: "test", ResourceType: "project"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDashboardsService_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Dashboards.List(context.Background(), nil)
	require.Error(t, err)
}

func TestDashboardsService_ListBuiltIn(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dashboards/built-ins", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Dashboards.ListBuiltIn(context.Background(), &DashboardsListBuiltInOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDashboardsService_ListBuiltIn_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Dashboards.ListBuiltIn(context.Background(), nil)
	require.Error(t, err)
}

func TestDashboardsService_GetBuiltIn(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dashboards/built-ins/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Dashboards.GetBuiltIn(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDashboardsService_GetBuiltIn_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Dashboards.GetBuiltIn(context.Background(), "")
	require.Error(t, err)
}

func TestDashboardsService_Get(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dashboards/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Dashboards.Get(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDashboardsService_Get_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Dashboards.Get(context.Background(), "")
	require.Error(t, err)
}

func TestDashboardsService_Update(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/dashboards/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Dashboards.Update(context.Background(), "test-id", &DashboardsUpdateOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDashboardsService_Update_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Dashboards.Update(context.Background(), "", nil)
	require.Error(t, err)
}

func TestDashboardsService_Create(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/dashboards", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Dashboards.Create(context.Background(), &DashboardsCreateOptions{Name: "test", Layout: "test", ResourceType: "project", ResourceId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestDashboardsService_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Dashboards.Create(context.Background(), nil)
	require.Error(t, err)
}
