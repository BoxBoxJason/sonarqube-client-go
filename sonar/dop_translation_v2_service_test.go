package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// CreateOrUpdateBoundProject
// =============================================================================

func TestDopTranslationV2_CreateOrUpdateBoundProject(t *testing.T) {
	response := DopTranslationBoundProject{
		ProjectId:         "proj-1",
		BindingId:         "binding-1",
		NewProjectCreated: false,
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPut, "/v2/dop-translation/bound-projects", http.StatusOK,
		&DopTranslationBoundProjectOptions{
			DevOpsPlatformSettingId: "dop-1",
			Monorepo:                false,
			ProjectKey:              "my-project",
			ProjectName:             "My Project",
			RepositoryIdentifier:    "my-org/my-repo",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(&DopTranslationBoundProjectOptions{
		DevOpsPlatformSettingId: "dop-1",
		Monorepo:                false,
		ProjectKey:              "my-project",
		ProjectName:             "My Project",
		RepositoryIdentifier:    "my-org/my-repo",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "proj-1", result.ProjectId)
	assert.Equal(t, "binding-1", result.BindingId)
	assert.False(t, result.NewProjectCreated)
}

func TestDopTranslationV2_CreateOrUpdateBoundProject_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *DopTranslationBoundProjectOptions
	}{
		{"nil opt", nil},
		{"missing DevOpsPlatformSettingId", &DopTranslationBoundProjectOptions{ProjectKey: "k", ProjectName: "n", RepositoryIdentifier: "r"}},
		{"missing ProjectKey", &DopTranslationBoundProjectOptions{DevOpsPlatformSettingId: "d", ProjectName: "n", RepositoryIdentifier: "r"}},
		{"missing ProjectName", &DopTranslationBoundProjectOptions{DevOpsPlatformSettingId: "d", ProjectKey: "k", RepositoryIdentifier: "r"}},
		{"missing RepositoryIdentifier", &DopTranslationBoundProjectOptions{DevOpsPlatformSettingId: "d", ProjectKey: "k", ProjectName: "n"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// CreateBoundProject
// =============================================================================

func TestDopTranslationV2_CreateBoundProject(t *testing.T) {
	response := DopTranslationBoundProject{
		ProjectId:         "proj-2",
		BindingId:         "binding-2",
		NewProjectCreated: true,
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/dop-translation/bound-projects", http.StatusCreated,
		&DopTranslationBoundProjectOptions{
			DevOpsPlatformSettingId: "dop-1",
			Monorepo:                true,
			ProjectKey:              "mono-project",
			ProjectName:             "Mono Project",
			RepositoryIdentifier:    "my-org/mono",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.CreateBoundProject(&DopTranslationBoundProjectOptions{
		DevOpsPlatformSettingId: "dop-1",
		Monorepo:                true,
		ProjectKey:              "mono-project",
		ProjectName:             "Mono Project",
		RepositoryIdentifier:    "my-org/mono",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "proj-2", result.ProjectId)
	assert.True(t, result.NewProjectCreated)
}

func TestDopTranslationV2_CreateBoundProject_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.DopTranslation.CreateBoundProject(nil)
	assert.Error(t, err)
}

// =============================================================================
// FetchAllDopSettings
// =============================================================================

func TestDopTranslationV2_FetchAllDopSettings(t *testing.T) {
	response := DopTranslationDopSettings{
		DopSettings: []DopSetting{
			{Id: "1", Type: "github", Key: "gh-setting", Url: "https://github.com", AppId: "app-1"},
			{Id: "2", Type: "gitlab", Key: "gl-setting", Url: "https://gitlab.com"},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 2},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/dop-settings", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.FetchAllDopSettings()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.DopSettings, 2)
	assert.Equal(t, "github", result.DopSettings[0].Type)
	assert.Equal(t, int32(2), result.Page.Total)
}
