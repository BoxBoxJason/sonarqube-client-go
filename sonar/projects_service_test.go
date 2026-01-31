package sonargo

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// ProjectsService Test Suite
// -----------------------------------------------------------------------------

// TestProjectsService_BulkDelete tests the BulkDelete method.
func TestProjectsService_BulkDelete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/projects/bulk_delete", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsBulkDeleteOption{
		Projects: []string{"project1", "project2"},
	}

	resp, err := client.Projects.BulkDelete(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestProjectsService_BulkDelete_ValidationError tests validation for BulkDelete.
func TestProjectsService_BulkDelete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test no filter provided
	opt := &ProjectsBulkDeleteOption{}
	_, err := client.Projects.BulkDelete(opt)
	assert.Error(t, err)

	// Test invalid qualifier
	opt = &ProjectsBulkDeleteOption{
		Query:      "test",
		Qualifiers: []string{"INVALID"},
	}
	_, err = client.Projects.BulkDelete(opt)
	assert.Error(t, err)
}

// TestProjectsService_Create tests the Create method.
func TestProjectsService_Create(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/projects/create", http.StatusOK, &ProjectsCreate{
		Project: Project{
			Key:        "my-project",
			Name:       "My Project",
			Qualifier:  "TRK",
			Visibility: "private",
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsCreateOption{
		Name:       "My Project",
		Project:    "my-project",
		Visibility: "private",
	}

	result, resp, err := client.Projects.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my-project", result.Project.Key)
	assert.Equal(t, "private", result.Project.Visibility)
}

// TestProjectsService_Create_ValidationError tests validation for Create.
func TestProjectsService_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Name
	opt := &ProjectsCreateOption{
		Project: "my-project",
	}
	_, _, err := client.Projects.Create(opt)
	assert.Error(t, err)

	// Test missing Project
	opt = &ProjectsCreateOption{
		Name: "My Project",
	}
	_, _, err = client.Projects.Create(opt)
	assert.Error(t, err)

	// Test Name too long
	opt = &ProjectsCreateOption{
		Name:    strings.Repeat("a", MaxProjectNameLength+1),
		Project: "my-project",
	}
	_, _, err = client.Projects.Create(opt)
	assert.Error(t, err)

	// Test Project key too long
	opt = &ProjectsCreateOption{
		Name:    "My Project",
		Project: strings.Repeat("a", MaxProjectKeyLength+1),
	}
	_, _, err = client.Projects.Create(opt)
	assert.Error(t, err)

	// Test invalid visibility
	opt = &ProjectsCreateOption{
		Name:       "My Project",
		Project:    "my-project",
		Visibility: "invalid",
	}
	_, _, err = client.Projects.Create(opt)
	assert.Error(t, err)
}

// TestProjectsService_Delete tests the Delete method.
func TestProjectsService_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/projects/delete", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsDeleteOption{
		Project: "my-project",
	}

	resp, err := client.Projects.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestProjectsService_Delete_ValidationError tests validation for Delete.
func TestProjectsService_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Project
	opt := &ProjectsDeleteOption{}
	_, err := client.Projects.Delete(opt)
	assert.Error(t, err)
}

// TestProjectsService_Search tests the Search method.
func TestProjectsService_Search(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/projects/search", http.StatusOK, &ProjectsSearch{
		Paging: Paging{PageIndex: 1, PageSize: 100, Total: 2},
		Components: []ProjectComponent{
			{
				Key:              "project1",
				Name:             "Project One",
				Qualifier:        "TRK",
				Visibility:       "public",
				LastAnalysisDate: "2024-01-15T10:30:00+0000",
			},
			{
				Key:        "project2",
				Name:       "Project Two",
				Qualifier:  "TRK",
				Visibility: "private",
				Managed:    true,
			},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsSearchOption{
		Query:      "project",
		Qualifiers: []string{"TRK"},
	}

	result, resp, err := client.Projects.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Components, 2)
	assert.Equal(t, "project1", result.Components[0].Key)
}

// TestProjectsService_Search_ValidationError tests validation for Search.
func TestProjectsService_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test invalid qualifier
	opt := &ProjectsSearchOption{
		Qualifiers: []string{"INVALID"},
	}
	_, _, err := client.Projects.Search(opt)
	assert.Error(t, err)
}

// TestProjectsService_SearchMyProjects tests the SearchMyProjects method.
func TestProjectsService_SearchMyProjects(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/projects/search_my_projects", http.StatusOK, &ProjectsSearchMyProjects{
		Paging: Paging{PageIndex: 1, PageSize: 100, Total: 1},
		Projects: []MyProject{
			{
				Key:         "my-project",
				Name:        "My Project",
				Description: "A test project",
				QualityGate: "OK",
			},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsSearchMyProjectsOption{}

	result, resp, err := client.Projects.SearchMyProjects(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Projects, 1)
}

// TestProjectsService_SearchMyScannableProjects tests the SearchMyScannableProjects method.
func TestProjectsService_SearchMyScannableProjects(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/projects/search_my_scannable_projects", http.StatusOK, &ProjectsSearchMyScannableProjects{
		Projects: []ScannableProject{
			{Key: "project1", Name: "Project One"},
			{Key: "project2", Name: "Project Two"},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	result, resp, err := client.Projects.SearchMyScannableProjects(&ProjectsSearchMyScannableProjectsOption{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Projects, 2)
}

// TestProjectsService_UpdateDefaultVisibility tests the UpdateDefaultVisibility method.
func TestProjectsService_UpdateDefaultVisibility(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/projects/update_default_visibility", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsUpdateDefaultVisibilityOption{
		ProjectVisibility: "private",
	}

	resp, err := client.Projects.UpdateDefaultVisibility(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestProjectsService_UpdateDefaultVisibility_ValidationError tests validation for UpdateDefaultVisibility.
func TestProjectsService_UpdateDefaultVisibility_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing ProjectVisibility
	opt := &ProjectsUpdateDefaultVisibilityOption{}
	_, err := client.Projects.UpdateDefaultVisibility(opt)
	assert.Error(t, err)

	// Test invalid visibility
	opt = &ProjectsUpdateDefaultVisibilityOption{
		ProjectVisibility: "invalid",
	}
	_, err = client.Projects.UpdateDefaultVisibility(opt)
	assert.Error(t, err)
}

// TestProjectsService_UpdateKey tests the UpdateKey method.
func TestProjectsService_UpdateKey(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/projects/update_key", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsUpdateKeyOption{
		From: "old-project-key",
		To:   "new-project-key",
	}

	resp, err := client.Projects.UpdateKey(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestProjectsService_UpdateKey_ValidationError tests validation for UpdateKey.
func TestProjectsService_UpdateKey_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing From
	opt := &ProjectsUpdateKeyOption{
		To: "new-key",
	}
	_, err := client.Projects.UpdateKey(opt)
	assert.Error(t, err)

	// Test missing To
	opt = &ProjectsUpdateKeyOption{
		From: "old-key",
	}
	_, err = client.Projects.UpdateKey(opt)
	assert.Error(t, err)
}

// TestProjectsService_UpdateVisibility tests the UpdateVisibility method.
func TestProjectsService_UpdateVisibility(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/projects/update_visibility", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectsUpdateVisibilityOption{
		Project:    "my-project",
		Visibility: "public",
	}

	resp, err := client.Projects.UpdateVisibility(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestProjectsService_UpdateVisibility_ValidationError tests validation for UpdateVisibility.
func TestProjectsService_UpdateVisibility_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Project
	opt := &ProjectsUpdateVisibilityOption{
		Visibility: "public",
	}
	_, err := client.Projects.UpdateVisibility(opt)
	assert.Error(t, err)

	// Test missing Visibility
	opt = &ProjectsUpdateVisibilityOption{
		Project: "my-project",
	}
	_, err = client.Projects.UpdateVisibility(opt)
	assert.Error(t, err)

	// Test invalid visibility
	opt = &ProjectsUpdateVisibilityOption{
		Project:    "my-project",
		Visibility: "invalid",
	}
	_, err = client.Projects.UpdateVisibility(opt)
	assert.Error(t, err)
}
