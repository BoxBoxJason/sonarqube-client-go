package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// GetWorkItems
// =============================================================================

func TestJiraV2_GetWorkItems(t *testing.T) {
	response := JiraWorkItem{"key": "PROJ-123", "status": "In Progress"}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/jira/work-items", http.StatusOK,
		map[string]string{
			"sonarProjectId": "proj-1",
			"resourceId":     "issue-1",
			"resourceType":   JiraResourceTypeSonarIssue,
		}, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetWorkItems(context.Background(), &JiraWorkItemsOptions{
		SonarProjectId: "proj-1",
		ResourceId:     "issue-1",
		ResourceType:   JiraResourceTypeSonarIssue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "PROJ-123", (*result)["key"])
}

func TestJiraV2_GetWorkItems_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	tests := []struct {
		opt  *JiraWorkItemsOptions
		name string
	}{
		{nil, "nil opt"},
		{&JiraWorkItemsOptions{ResourceId: "r1", ResourceType: JiraResourceTypeSonarIssue}, "missing SonarProjectId"},
		{&JiraWorkItemsOptions{SonarProjectId: "p1", ResourceType: JiraResourceTypeSonarIssue}, "missing ResourceId"},
		{&JiraWorkItemsOptions{SonarProjectId: "p1", ResourceId: "r1"}, "missing ResourceType"},
		{&JiraWorkItemsOptions{SonarProjectId: "p1", ResourceId: "r1", ResourceType: "BOGUS"}, "invalid ResourceType"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.GetWorkItems(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// CreateWorkItem
// =============================================================================

func TestJiraV2_CreateWorkItem(t *testing.T) {
	body := JiraWorkItem{
		"sonarProjectId": "proj-1",
		"resourceId":     "issue-1",
		"resourceType":   JiraResourceTypeSonarIssue,
	}
	response := JiraWorkItem{"key": "PROJ-123"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/jira/work-items", http.StatusCreated, body, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.CreateWorkItem(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "PROJ-123", (*result)["key"])
}

func TestJiraV2_CreateWorkItem_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.CreateWorkItem(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// DeleteWorkItems
// =============================================================================

func TestJiraV2_DeleteWorkItems(t *testing.T) {
	server := newTestServer(t, mockEmptyHandlerWithParams(t, http.MethodDelete, "/v2/jira/work-items", http.StatusNoContent,
		map[string]string{
			"sonarProjectId": "proj-1",
			"resourceId":     "issue-1",
			"resourceType":   JiraResourceTypeDependencyRisk,
		}))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	resp, err := svc.DeleteWorkItems(context.Background(), &JiraWorkItemsOptions{
		SonarProjectId: "proj-1",
		ResourceId:     "issue-1",
		ResourceType:   JiraResourceTypeDependencyRisk,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestJiraV2_DeleteWorkItems_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, err := svc.DeleteWorkItems(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// GetProjectBinding
// =============================================================================

func TestJiraV2_GetProjectBinding(t *testing.T) {
	response := JiraProjectBinding{"jiraProjectKey": "PROJ"}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/jira/project-bindings", http.StatusOK,
		map[string]string{"sonarProjectId": "proj-1"}, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetProjectBinding(context.Background(), &JiraProjectBindingOptions{SonarProjectId: "proj-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "PROJ", (*result)["jiraProjectKey"])
}

func TestJiraV2_GetProjectBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.GetProjectBinding(context.Background(), &JiraProjectBindingOptions{})
	assert.Error(t, err)

	_, _, err = svc.GetProjectBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// CreateProjectBinding
// =============================================================================

func TestJiraV2_CreateProjectBinding(t *testing.T) {
	body := JiraProjectBinding{"sonarProjectId": "proj-1", "jiraProjectKey": "PROJ"}
	response := JiraProjectBinding{"sonarProjectId": "proj-1", "jiraProjectKey": "PROJ"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/jira/project-bindings", http.StatusCreated, body, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.CreateProjectBinding(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "PROJ", (*result)["jiraProjectKey"])
}

func TestJiraV2_CreateProjectBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.CreateProjectBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// UpdateProjectBinding
// =============================================================================

func TestJiraV2_UpdateProjectBinding(t *testing.T) {
	body := JiraProjectBinding{"sonarProjectId": "proj-1", "jiraProjectKey": "PROJ2"}
	response := JiraProjectBinding{"sonarProjectId": "proj-1", "jiraProjectKey": "PROJ2"}
	server := newTestServer(t, mockPatchHandler(t, "/v2/jira/project-bindings", http.StatusOK, body, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.UpdateProjectBinding(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "PROJ2", (*result)["jiraProjectKey"])
}

func TestJiraV2_UpdateProjectBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.UpdateProjectBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// DeleteProjectBinding
// =============================================================================

func TestJiraV2_DeleteProjectBinding(t *testing.T) {
	server := newTestServer(t, mockEmptyHandlerWithParams(t, http.MethodDelete, "/v2/jira/project-bindings", http.StatusNoContent,
		map[string]string{"sonarProjectId": "proj-1"}))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	resp, err := svc.DeleteProjectBinding(context.Background(), &JiraProjectBindingOptions{SonarProjectId: "proj-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestJiraV2_DeleteProjectBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, err := svc.DeleteProjectBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// GetOrganizationBinding
// =============================================================================

func TestJiraV2_GetOrganizationBinding(t *testing.T) {
	response := JiraOrganizationBinding{"cloudId": "cloud-1"}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/jira/organization-bindings", http.StatusOK,
		map[string]string{"sonarOrganizationUuid": "org-uuid-1"}, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetOrganizationBinding(context.Background(), &JiraOrganizationBindingOptions{SonarOrganizationUuid: "org-uuid-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "cloud-1", (*result)["cloudId"])
}

func TestJiraV2_GetOrganizationBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.GetOrganizationBinding(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = svc.GetOrganizationBinding(context.Background(), &JiraOrganizationBindingOptions{})
	assert.Error(t, err)
}

// =============================================================================
// CreateOrganizationBinding
// =============================================================================

func TestJiraV2_CreateOrganizationBinding(t *testing.T) {
	body := JiraOrganizationBinding{"state": "abc", "code": "xyz"}
	response := JiraOrganizationBinding{"cloudId": "cloud-1"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/jira/organization-bindings", http.StatusCreated, body, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.CreateOrganizationBinding(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "cloud-1", (*result)["cloudId"])
}

func TestJiraV2_CreateOrganizationBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.CreateOrganizationBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// BindOrganizationBinding
// =============================================================================

func TestJiraV2_BindOrganizationBinding(t *testing.T) {
	body := JiraOrganizationBinding{"cloudId": "cloud-1"}
	response := JiraOrganizationBinding{"cloudId": "cloud-1", "bound": true}
	server := newTestServer(t, mockPatchHandler(t, "/v2/jira/organization-bindings", http.StatusCreated, body, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.BindOrganizationBinding(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "cloud-1", (*result)["cloudId"])
}

func TestJiraV2_BindOrganizationBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.BindOrganizationBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// DeleteOrganizationBinding
// =============================================================================

func TestJiraV2_DeleteOrganizationBinding(t *testing.T) {
	server := newTestServer(t, mockEmptyHandlerWithParams(t, http.MethodDelete, "/v2/jira/organization-bindings", http.StatusNoContent,
		map[string]string{"sonarOrganizationUuid": "org-uuid-1"}))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	resp, err := svc.DeleteOrganizationBinding(context.Background(), &JiraOrganizationBindingOptions{SonarOrganizationUuid: "org-uuid-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestJiraV2_DeleteOrganizationBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, err := svc.DeleteOrganizationBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// UpdateOrganizationBinding
// =============================================================================

func TestJiraV2_UpdateOrganizationBinding(t *testing.T) {
	body := JiraOrganizationBinding{"tokenSharingEnabled": true}
	response := JiraOrganizationBinding{"tokenSharingEnabled": true}
	server := newTestServer(t, mockPatchHandler(t, "/v2/jira/organization-binding-edit", http.StatusCreated, body, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.UpdateOrganizationBinding(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, true, (*result)["tokenSharingEnabled"])
}

func TestJiraV2_UpdateOrganizationBinding_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.UpdateOrganizationBinding(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// GetWorkTypes
// =============================================================================

func TestJiraV2_GetWorkTypes(t *testing.T) {
	response := []JiraWorkType{
		{"id": "10001", "name": "Bug"},
		{"id": "10002", "name": "Task"},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/jira/work-types", http.StatusOK,
		map[string]string{
			"jiraProjectKey":        "PROJ",
			"sonarOrganizationUuid": "org-uuid-1",
			"sonarProjectId":        "proj-1",
			"includeFields":         "true",
		}, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetWorkTypes(context.Background(), &JiraWorkTypesOptions{
		JiraProjectKey:        "PROJ",
		SonarOrganizationUuid: "org-uuid-1",
		SonarProjectId:        "proj-1",
		IncludeFields:         true,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 2)
	assert.Equal(t, "Bug", result[0]["name"])
}

func TestJiraV2_GetWorkTypes_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	tests := []struct {
		opt  *JiraWorkTypesOptions
		name string
	}{
		{nil, "nil opt"},
		{&JiraWorkTypesOptions{SonarOrganizationUuid: "org-1"}, "missing JiraProjectKey"},
		{&JiraWorkTypesOptions{JiraProjectKey: "PROJ"}, "missing SonarOrganizationUuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.GetWorkTypes(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// UpdateWorkTypes
// =============================================================================

func TestJiraV2_UpdateWorkTypes(t *testing.T) {
	body := JiraWorkTypeSelection{"jiraProjectKey": "PROJ", "workTypeIds": []any{"10001", "10002"}}
	server := newTestServer(t, mockPatchHandler(t, "/v2/jira/work-types", http.StatusOK, body, nil))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	resp, err := svc.UpdateWorkTypes(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestJiraV2_UpdateWorkTypes_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, err := svc.UpdateWorkTypes(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// GetUserActions
// =============================================================================

func TestJiraV2_GetUserActions(t *testing.T) {
	response := []string{"CREATE_WORK_ITEM", "VIEW_WORK_ITEM"}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/jira/user-actions", http.StatusOK,
		map[string]string{"sonarProjectId": "proj-1"}, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetUserActions(context.Background(), &JiraUserActionsOptions{SonarProjectId: "proj-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []string{"CREATE_WORK_ITEM", "VIEW_WORK_ITEM"}, result)
}

func TestJiraV2_GetUserActions_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.GetUserActions(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = svc.GetUserActions(context.Background(), &JiraUserActionsOptions{})
	assert.Error(t, err)
}

// =============================================================================
// GetProjects
// =============================================================================

func TestJiraV2_GetProjects(t *testing.T) {
	response := []JiraProject{
		{"key": "PROJ", "name": "Project"},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/jira/projects", http.StatusOK,
		map[string]string{"sonarOrganizationUuid": "org-uuid-1"}, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetProjects(context.Background(), &JiraProjectsOptions{SonarOrganizationUuid: "org-uuid-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
	assert.Equal(t, "PROJ", result[0]["key"])
}

func TestJiraV2_GetProjects_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.GetProjects(context.Background(), nil)
	assert.Error(t, err)

	_, _, err = svc.GetProjects(context.Background(), &JiraProjectsOptions{})
	assert.Error(t, err)
}

// =============================================================================
// GetLinkedIssuesCount
// =============================================================================

func TestJiraV2_GetLinkedIssuesCount(t *testing.T) {
	response := JiraLinkedIssuesCount{"count": float64(3)}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/jira/linked-issues-count/proj-1", http.StatusOK, response))
	client := newTestClient(t, server.url())
	svc := &JiraService{client: client}

	result, resp, err := svc.GetLinkedIssuesCount(context.Background(), "proj-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.InDelta(t, float64(3), (*result)["count"], 0)
}

func TestJiraV2_GetLinkedIssuesCount_Validation(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &JiraService{client: client}

	_, _, err := svc.GetLinkedIssuesCount(context.Background(), "")
	assert.Error(t, err)
}
