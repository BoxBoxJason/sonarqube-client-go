package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectLinks_Create(t *testing.T) {
	response := &ProjectLinksCreate{
		Link: ProjectLink{
			ID:   "1",
			Name: "Homepage",
			Type: "homepage",
			URL:  "https://example.com",
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/project_links/create", http.StatusOK, response))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectLinksCreateOption{
		Name:       "Homepage",
		ProjectKey: "my-project",
		URL:        "https://example.com",
	}

	result, resp, err := client.ProjectLinks.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "1", result.Link.ID)
	assert.Equal(t, "Homepage", result.Link.Name)
}

func TestProjectLinks_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectLinks.Create(nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, _, err = client.ProjectLinks.Create(&ProjectLinksCreateOption{
		ProjectKey: "my-project",
		URL:        "https://example.com",
	})
	assert.Error(t, err)

	// Missing URL should fail validation.
	_, _, err = client.ProjectLinks.Create(&ProjectLinksCreateOption{
		Name:       "Homepage",
		ProjectKey: "my-project",
	})
	assert.Error(t, err)

	// Missing ProjectID and ProjectKey should fail validation.
	_, _, err = client.ProjectLinks.Create(&ProjectLinksCreateOption{
		Name: "Homepage",
		URL:  "https://example.com",
	})
	assert.Error(t, err)
}

func TestProjectLinks_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_links/delete", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectLinksDeleteOption{
		ID: "1",
	}

	resp, err := client.ProjectLinks.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectLinks_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectLinks.Delete(nil)
	assert.Error(t, err)

	// Missing ID should fail validation.
	_, err = client.ProjectLinks.Delete(&ProjectLinksDeleteOption{})
	assert.Error(t, err)
}

func TestProjectLinks_Search(t *testing.T) {
	response := &ProjectLinksSearch{
		Links: []ProjectLink{
			{ID: "1", Name: "Homepage", Type: "homepage", URL: "https://example.com"},
			{ID: "2", Name: "CI", Type: "ci", URL: "https://ci.example.com"},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_links/search", http.StatusOK, response))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectLinksSearchOption{
		ProjectKey: "my-project",
	}

	result, resp, err := client.ProjectLinks.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Links, 2)
	assert.Equal(t, "1", result.Links[0].ID)
}

func TestProjectLinks_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectLinks.Search(nil)
	assert.Error(t, err)

	// Missing ProjectID and ProjectKey should fail validation.
	_, _, err = client.ProjectLinks.Search(&ProjectLinksSearchOption{})
	assert.Error(t, err)
}

func TestProjectLinks_ValidateCreateOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option with ProjectKey should pass.
	err := client.ProjectLinks.ValidateCreateOpt(&ProjectLinksCreateOption{
		Name:       "Homepage",
		ProjectKey: "my-project",
		URL:        "https://example.com",
	})
	assert.NoError(t, err)

	// Valid option with ProjectID should pass.
	err = client.ProjectLinks.ValidateCreateOpt(&ProjectLinksCreateOption{
		Name:      "Homepage",
		ProjectID: "project-id",
		URL:       "https://example.com",
	})
	assert.NoError(t, err)

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxLinkNameLength+1; i++ {
		longName += "a"
	}
	err = client.ProjectLinks.ValidateCreateOpt(&ProjectLinksCreateOption{
		Name:       longName,
		ProjectKey: "my-project",
		URL:        "https://example.com",
	})
	assert.Error(t, err)
}
