package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectTagsService_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/project_tags/search", http.StatusOK, `{
			"tags": ["security", "performance", "bug"]
		}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		result, resp, err := client.ProjectTags.Search(nil)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
		assert.Len(t, result.Tags, 3)
		assert.Equal(t, "security", result.Tags[0])
	})

	t.Run("with pagination", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/project_tags/search", http.StatusOK, `{"tags": ["security"]}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		opt := &ProjectTagsSearchOption{
			PaginationArgs: PaginationArgs{
				Page:     2,
				PageSize: 10,
			},
			Query: "sec",
		}

		result, resp, err := client.ProjectTags.Search(opt)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, result.Tags, 1)
	})
}

func TestProjectTagsService_Set(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := mockEmptyHandler(t, http.MethodPost, "/project_tags/set", http.StatusNoContent)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		opt := &ProjectTagsSetOption{
			Project: "my-project",
			Tags:    []string{"security", "performance"},
		}

		resp, err := client.ProjectTags.Set(opt)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nil option fails validation", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.ProjectTags.Set(nil)

		assert.Error(t, err)
	})

	t.Run("missing project fails validation", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.ProjectTags.Set(&ProjectTagsSetOption{
			Tags: []string{"tag1"},
		})

		assert.Error(t, err)
	})
}

func TestProjectTagsService_ValidateSearchOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectTagsSearchOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &ProjectTagsSearchOption{}, false},
		{"with query", &ProjectTagsSearchOption{Query: "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectTags.ValidateSearchOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectTagsService_ValidateSetOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectTagsSetOption
		wantErr bool
	}{
		{"valid", &ProjectTagsSetOption{Project: "my-project", Tags: []string{"tag1", "tag2"}}, false},
		{"nil option", nil, true},
		{"missing project", &ProjectTagsSetOption{Tags: []string{"tag1"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectTags.ValidateSetOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
