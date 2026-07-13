package sonar

import (
	"context"
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

	result, resp, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(context.Background(), &DopTranslationBoundProjectOptions{
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
			_, _, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(context.Background(), tt.opt)
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

	result, resp, err := client.V2.DopTranslation.CreateBoundProject(context.Background(), &DopTranslationBoundProjectOptions{
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

	_, _, err := client.V2.DopTranslation.CreateBoundProject(context.Background(), nil)
	assert.Error(t, err)
}

// =============================================================================
// GetDopSettings
// =============================================================================

func TestDopTranslationV2_GetDopSettings(t *testing.T) {
	response := DopTranslationDopSettings{
		DopSettings: []DopSetting{
			{Id: "1", Type: "github", Key: "gh-setting", URL: "https://github.com", AppID: "app-1"},
			{Id: "2", Type: "gitlab", Key: "gl-setting", URL: "https://gitlab.com"},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 2},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/dop-settings", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.GetDopSettings(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.DopSettings, 2)
	assert.Equal(t, "github", result.DopSettings[0].Type)
	assert.Equal(t, int32(2), result.Page.Total)
}

// =============================================================================
// GetJfrogEvidence
// =============================================================================

func TestDopTranslationV2_GetJfrogEvidence(t *testing.T) {
	response := map[string]any{
		"_type":         "https://in-toto.io/Statement/v1",
		"predicateType": "https://sonarsource.com/quality-gate/v1",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/jfrog-evidence/task-1", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.DopTranslation.GetJfrogEvidence(context.Background(), "task-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "https://in-toto.io/Statement/v1", (*result)["_type"])
}

func TestDopTranslationV2_GetJfrogEvidence_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.DopTranslation.GetJfrogEvidence(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
