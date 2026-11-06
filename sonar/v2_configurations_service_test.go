package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemV2_EmailConfigurations(t *testing.T) {
	ctx := context.Background()
	body := map[string]any{"id": "email-configuration", "host": "h", "authMethod": "BASIC", "username": "u", "oauthScope": "s", "oauthTenant": nil}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/system/email-configurations", http.StatusOK,
		map[string]any{"emailConfigurations": []any{body}, "page": map[string]any{"pageIndex": 1, "pageSize": 1000, "total": 1}}))
	client := newTestClient(t, server.URL)

	list, _, err := client.V2.System.SearchEmailConfigurations(ctx)
	require.NoError(t, err)
	require.Len(t, list.EmailConfigurations, 1)
	assert.Equal(t, "s", *list.EmailConfigurations[0].OauthScope)
	assert.Nil(t, list.EmailConfigurations[0].OauthTenant)

	create := &SystemEmailConfigurationCreateOptions{
		Host: "h", Port: "587", SecurityProtocol: "STARTTLS", FromAddress: "a@b.c", FromName: "n",
		SubjectPrefix: "[SQ]", AuthMethod: "BASIC", Username: "u", BasicPassword: "p", OauthScope: "s",
	}
	server = newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/system/email-configurations", http.StatusOK, create, body))
	client = newTestClient(t, server.URL)

	created, _, err := client.V2.System.CreateEmailConfiguration(ctx, create)
	require.NoError(t, err)
	assert.Equal(t, "email-configuration", created.Id)

	_, _, err = client.V2.System.CreateEmailConfiguration(ctx, &SystemEmailConfigurationCreateOptions{})
	require.Error(t, err)

	bad := *create
	bad.SecurityProtocol = "nope"
	_, _, err = client.V2.System.CreateEmailConfiguration(ctx, &bad)
	require.Error(t, err)

	server = newTestServer(t, mockHandler(t, http.MethodGet, "/v2/system/email-configurations/email-configuration", http.StatusOK, body))
	client = newTestClient(t, server.URL)

	got, _, err := client.V2.System.GetEmailConfiguration(ctx, "email-configuration")
	require.NoError(t, err)
	assert.Equal(t, "h", got.Host)

	_, _, err = client.V2.System.GetEmailConfiguration(ctx, "")
	require.Error(t, err)

	upd := &SystemEmailConfigurationUpdateOptions{OauthScope: "x"}
	server = newTestServer(t, mockPatchHandler(t, "/v2/system/email-configurations/email-configuration", http.StatusOK, upd, body))
	client = newTestClient(t, server.URL)

	_, _, err = client.V2.System.UpdateEmailConfiguration(ctx, "email-configuration", upd)
	require.NoError(t, err)

	server = newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/system/email-configurations/email-configuration", http.StatusNoContent))
	client = newTestClient(t, server.URL)

	resp, err := client.V2.System.DeleteEmailConfiguration(ctx, "email-configuration")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestCleanCodePolicyV2_Mode(t *testing.T) {
	ctx := context.Background()

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/clean-code-policy/mode", http.StatusOK, `{"mode":"MQR","modified":false}`))
	client := newTestClient(t, server.URL)

	mode, _, err := client.V2.CleanCodePolicy.GetMode(ctx)
	require.NoError(t, err)
	assert.Equal(t, "MQR", mode.Mode)

	opt := &CleanCodePolicyUpdateModeOptions{Mode: InstanceModeStandardExperience}
	server = newTestServer(t, mockPatchHandler(t, "/v2/clean-code-policy/mode", http.StatusOK, opt, map[string]any{"mode": "STANDARD_EXPERIENCE", "modified": true}))
	client = newTestClient(t, server.URL)

	mode, _, err = client.V2.CleanCodePolicy.UpdateMode(ctx, opt)
	require.NoError(t, err)
	assert.True(t, mode.Modified)

	_, _, err = client.V2.CleanCodePolicy.UpdateMode(ctx, &CleanCodePolicyUpdateModeOptions{Mode: "bad"})
	require.Error(t, err)
}

func TestDopTranslationV2_Configurations(t *testing.T) {
	ctx := context.Background()
	gh := map[string]any{"id": "github-configuration", "enabled": false, "allowedOrganizations": []string{"o"}, "provisioningType": "JIT", "errorMessage": nil}
	gl := map[string]any{"id": "gitlab-configuration", "allowedGroups": []string{"g"}, "provisioningType": "JIT", "isProvisioningTokenSet": false}

	t.Run("github", func(t *testing.T) {
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/github-configurations", http.StatusOK,
			map[string]any{"githubConfigurations": []any{gh}, "page": map[string]any{"total": 1}}))
		client := newTestClient(t, server.URL)

		list, _, err := client.V2.DopTranslation.SearchGithubConfigurations(ctx)
		require.NoError(t, err)
		require.Len(t, list.GithubConfigurations, 1)
		assert.Equal(t, []string{"o"}, list.GithubConfigurations[0].AllowedOrganizations)

		create := &DopTranslationGithubConfigurationCreateOptions{
			ClientId: "c", ClientSecret: "s", ApplicationId: "1", PrivateKey: "k", ApiUrl: "a", WebUrl: "w",
			AllowedOrganizations: []string{"o"}, ProvisioningType: "JIT",
		}
		server = newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/dop-translation/github-configurations", http.StatusOK, create, gh))
		client = newTestClient(t, server.URL)

		_, _, err = client.V2.DopTranslation.CreateGithubConfiguration(ctx, create)
		require.NoError(t, err)

		_, _, err = client.V2.DopTranslation.CreateGithubConfiguration(ctx, &DopTranslationGithubConfigurationCreateOptions{})
		require.Error(t, err)

		upd := &DopTranslationGithubConfigurationUpdateOptions{AllowedOrganizations: &UpdateFieldListStringV2{Value: []string{"a", "b"}, Defined: true}}
		server = newTestServer(t, mockPatchHandler(t, "/v2/dop-translation/github-configurations/github-configuration", http.StatusOK,
			map[string]any{"allowedOrganizations": []string{"a", "b"}}, gh))
		client = newTestClient(t, server.URL)

		_, _, err = client.V2.DopTranslation.UpdateGithubConfiguration(ctx, "github-configuration", upd)
		require.NoError(t, err)

		server = newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/dop-translation/github-configurations/github-configuration", http.StatusNoContent))
		client = newTestClient(t, server.URL)

		_, err = client.V2.DopTranslation.DeleteGithubConfiguration(ctx, "github-configuration")
		require.NoError(t, err)

		_, _, err = client.V2.DopTranslation.GetGithubConfiguration(ctx, "")
		require.Error(t, err)
	})

	t.Run("gitlab", func(t *testing.T) {
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/gitlab-configurations/gitlab-configuration", http.StatusOK, gl))
		client := newTestClient(t, server.URL)

		got, _, err := client.V2.DopTranslation.GetGitlabConfiguration(ctx, "gitlab-configuration")
		require.NoError(t, err)
		assert.Equal(t, []string{"g"}, got.AllowedGroups)

		create := &DopTranslationGitlabConfigurationCreateOptions{ApplicationId: "1", Url: "u", Secret: "s", AllowedGroups: []string{"g"}, ProvisioningType: "AUTO_PROVISIONING"}
		server = newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/dop-translation/gitlab-configurations", http.StatusOK, create, gl))
		client = newTestClient(t, server.URL)

		_, _, err = client.V2.DopTranslation.CreateGitlabConfiguration(ctx, create)
		require.NoError(t, err)

		bad := *create
		bad.ProvisioningType = "x"
		_, _, err = client.V2.DopTranslation.CreateGitlabConfiguration(ctx, &bad)
		require.Error(t, err)

		server = newTestServer(t, mockHandler(t, http.MethodPost, "/v2/dop-translation/gitlab-synchronization-runs", http.StatusOK, `{"id":"run"}`))
		client = newTestClient(t, server.URL)

		run, _, err := client.V2.DopTranslation.StartGitlabSynchronizationRun(ctx)
		require.NoError(t, err)
		assert.Equal(t, "run", run.Id)
	})
}

func TestDopTranslationV2_PermissionMappingsAndBindings(t *testing.T) {
	ctx := context.Background()
	mapping := map[string]any{"id": "custom", "role": "custom", "baseRole": false,
		"permissions": map[string]any{"user": true, "codeViewer": true, "issueAdmin": false, "securityHotspotAdmin": false, "admin": false, "scan": true}}

	for _, platform := range []string{"github", "gitlab"} {
		base := "/v2/dop-translation/" + platform + "-permission-mappings"
		server := newTestServer(t, mockHandler(t, http.MethodGet, base, http.StatusOK, map[string]any{"permissionMappings": []any{mapping}}))
		client := newTestClient(t, server.URL)

		var (
			list *DopTranslationPermissionMappingsSearch
			err  error
		)

		if platform == "github" {
			list, _, err = client.V2.DopTranslation.SearchGithubPermissionMappings(ctx)
		} else {
			list, _, err = client.V2.DopTranslation.SearchGitlabPermissionMappings(ctx)
		}

		require.NoError(t, err)
		require.Len(t, list.PermissionMappings, 1)
		assert.True(t, list.PermissionMappings[0].Permissions.Scan)

		create := &DopTranslationPermissionMappingCreateOptions{Role: "custom", Permissions: DopTranslationPermissions{User: true, Scan: true}}
		server = newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, base, http.StatusOK, create, mapping))
		client = newTestClient(t, server.URL)

		if platform == "github" {
			_, _, err = client.V2.DopTranslation.CreateGithubPermissionMapping(ctx, create)
		} else {
			_, _, err = client.V2.DopTranslation.CreateGitlabPermissionMapping(ctx, create)
		}

		require.NoError(t, err)

		scan := true
		upd := &DopTranslationPermissionMappingUpdateOptions{Permissions: &DopTranslationPermissionsUpdate{Scan: &scan}}
		server = newTestServer(t, mockPatchHandler(t, base+"/custom", http.StatusOK, map[string]any{"permissions": map[string]any{"scan": true}}, mapping))
		client = newTestClient(t, server.URL)

		if platform == "github" {
			_, _, err = client.V2.DopTranslation.UpdateGithubPermissionMapping(ctx, "custom", upd)
		} else {
			_, _, err = client.V2.DopTranslation.UpdateGitlabPermissionMapping(ctx, "custom", upd)
		}

		require.NoError(t, err)

		server = newTestServer(t, mockEmptyHandler(t, http.MethodDelete, base+"/custom", http.StatusNoContent))
		client = newTestClient(t, server.URL)

		if platform == "github" {
			_, err = client.V2.DopTranslation.DeleteGithubPermissionMapping(ctx, "custom")
		} else {
			_, err = client.V2.DopTranslation.DeleteGitlabPermissionMapping(ctx, "custom")
		}

		require.NoError(t, err)
	}

	binding := map[string]any{"id": "b", "devOpsPlatformSettingId": "s", "projectId": "p", "projectKey": "k", "repository": "o/r", "slug": nil}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/dop-translation/project-bindings", http.StatusOK,
		map[string]string{"repository": "o/r", "pageSize": "1"}, map[string]any{"projectBindings": []any{binding}, "page": map[string]any{"total": 1}}))
	client := newTestClient(t, server.URL)

	list, _, err := client.V2.DopTranslation.SearchProjectBindings(ctx, &DopTranslationProjectBindingsSearchOptions{Repository: "o/r", PageSize: 1})
	require.NoError(t, err)
	require.Len(t, list.ProjectBindings, 1)
	assert.Equal(t, "o/r", *list.ProjectBindings[0].Repository)
	assert.Nil(t, list.ProjectBindings[0].Slug)

	server = newTestServer(t, mockHandler(t, http.MethodGet, "/v2/dop-translation/project-bindings/b", http.StatusOK, binding))
	client = newTestClient(t, server.URL)

	one, _, err := client.V2.DopTranslation.GetProjectBinding(ctx, "b")
	require.NoError(t, err)
	assert.Equal(t, "k", one.ProjectKey)
}

func TestBoundProjectAndUserUpdateNewFields(t *testing.T) {
	ctx := context.Background()
	off := false
	opt := &DopTranslationBoundProjectOptions{
		DevOpsPlatformSettingId: "d", ProjectKey: "k", ProjectName: "n", RepositoryIdentifier: "r", SummaryCommentEnabled: &off,
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/dop-translation/bound-projects", http.StatusOK, opt, map[string]any{"id": "x"}))
	client := newTestClient(t, server.URL)

	_, _, err := client.V2.DopTranslation.CreateBoundProject(ctx, opt)
	require.NoError(t, err)

	local := true
	upd := &UsersUpdateOptionsV2{Local: &local, Name: "n"}
	server = newTestServer(t, mockPatchHandler(t, "/v2/users-management/users/u1", http.StatusOK, map[string]any{"local": true, "name": "n"}, map[string]any{"id": "u1"}))
	client = newTestClient(t, server.URL)

	_, _, err = client.V2.UsersManagement.Update(ctx, "u1", upd)
	require.NoError(t, err)
}
