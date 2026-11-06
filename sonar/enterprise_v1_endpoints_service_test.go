package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnterpriseV1Endpoints(t *testing.T) {
	ctx := context.Background()

	t.Run("AlmSettings.ValidateBinding", func(t *testing.T) {
		server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/alm_settings/validate_binding", http.StatusOK,
			map[string]string{"project": "p"}, `{"scope":"GLOBAL","errors":[{"msg":"bad key"}]}`))
		client := newTestClient(t, server.URL)

		result, _, err := client.AlmSettings.ValidateBinding(ctx, &AlmSettingsValidateBindingOptions{Project: "p"})
		require.NoError(t, err)
		assert.Equal(t, "GLOBAL", result.Scope)
		require.Len(t, result.Errors, 1)
		assert.Equal(t, "bad key", result.Errors[0].Msg)

		_, _, err = client.AlmSettings.ValidateBinding(ctx, &AlmSettingsValidateBindingOptions{})
		require.Error(t, err)
	})

	t.Run("Ce.SetWorkerCount", func(t *testing.T) {
		server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/ce/set_worker_count", http.StatusNoContent))
		client := newTestClient(t, server.URL)

		resp, err := client.Ce.SetWorkerCount(ctx, &CeSetWorkerCountOptions{Count: 2})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, err = client.Ce.SetWorkerCount(ctx, &CeSetWorkerCountOptions{Count: 0})
		require.Error(t, err)
	})

	t.Run("GithubProvisioning.Status and Sync", func(t *testing.T) {
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/github_provisioning/status", http.StatusOK, `{"enabled":true}`))
		client := newTestClient(t, server.URL)

		result, _, err := client.GithubProvisioning.Status(ctx)
		require.NoError(t, err)
		assert.True(t, result.Enabled)

		server = newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/github_provisioning/sync", http.StatusNoContent))
		client = newTestClient(t, server.URL)

		_, err = client.GithubProvisioning.Sync(ctx)
		require.NoError(t, err)
	})

	t.Run("Issues.GitlabSastExport", func(t *testing.T) {
		server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/issues/gitlab_sast_export", http.StatusOK,
			map[string]string{"projectKey": "p"},
			`{"version":"15.0.0","scan":{"analyzer":{"id":"a","name":"n","vendor":{"name":"Sonar"},"version":"1"},"start_time":"s","end_time":"e","status":"success","type":"sast"},"vulnerabilities":[{"id":"1"}]}`))
		client := newTestClient(t, server.URL)

		result, _, err := client.Issues.GitlabSastExport(ctx, &IssuesGitlabSastExportOptions{ProjectKey: "p"})
		require.NoError(t, err)
		assert.Equal(t, "15.0.0", result.Version)
		assert.Equal(t, "Sonar", result.Scan.Analyzer.Vendor.Name)
		assert.Equal(t, "success", result.Scan.Status)
		assert.Len(t, result.Vulnerabilities, 1)

		_, _, err = client.Issues.GitlabSastExport(ctx, &IssuesGitlabSastExportOptions{ProjectKey: "p", Branch: "b", PullRequest: "1"})
		require.Error(t, err)
	})

	t.Run("ProjectBranches.GetAiCodeAssurance", func(t *testing.T) {
		server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_branches/get_ai_code_assurance", http.StatusOK, `{"aiCodeAssurance":"NONE"}`))
		client := newTestClient(t, server.URL)

		result, _, err := client.ProjectBranches.GetAiCodeAssurance(ctx, &ProjectBranchesGetAiCodeAssuranceOptions{Project: "p"})
		require.NoError(t, err)
		assert.Equal(t, "NONE", result.AiCodeAssurance)
	})

	t.Run("ProjectDump.Import", func(t *testing.T) {
		server := newTestServer(t, mockHandler(t, http.MethodPost, "/project_dump/import", http.StatusOK, `{"taskId":"t","projectId":"i","projectKey":"k","projectName":"n"}`))
		client := newTestClient(t, server.URL)

		result, _, err := client.ProjectDump.Import(ctx, &ProjectDumpImportOptions{Key: "k"})
		require.NoError(t, err)
		assert.Equal(t, "t", result.TaskID)

		_, _, err = client.ProjectDump.Import(ctx, &ProjectDumpImportOptions{})
		require.Error(t, err)
	})

	t.Run("Projects AI code and findings", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			switch r.URL.Path {
			case "/projects/get_contains_ai_code":
				_, _ = w.Write([]byte(`{"containsAiCode":true}`))
			case "/projects/get_detected_ai_code":
				_, _ = w.Write([]byte(`{"detectedAiCode":true}`))
			case "/projects/export_findings":
				_, _ = w.Write([]byte(`{"export_findings":[{"key":"k"}]}`))
			case "/projects/set_contains_ai_code":
				require.NoError(t, r.ParseForm())
				assert.Equal(t, "true", r.Form.Get("contains_ai_code"))
				w.WriteHeader(http.StatusNoContent)
			default:
				t.Errorf("unexpected path %s", r.URL.Path)
			}
		})
		client := newTestClient(t, server.URL)

		contains, _, err := client.Projects.GetContainsAiCode(ctx, &ProjectsProjectOptions{Project: "p"})
		require.NoError(t, err)
		assert.True(t, contains.ContainsAiCode)

		detected, _, err := client.Projects.GetDetectedAiCode(ctx, &ProjectsProjectOptions{Project: "p"})
		require.NoError(t, err)
		assert.True(t, detected.DetectedAiCode)

		findings, _, err := client.Projects.ExportFindings(ctx, &ProjectsExportFindingsOptions{Project: "p"})
		require.NoError(t, err)
		assert.Len(t, findings.ExportFindings, 1)

		_, err = client.Projects.SetContainsAiCode(ctx, &ProjectsSetContainsAiCodeOptions{Project: "p", ContainsAiCode: true})
		require.NoError(t, err)

		_, _, err = client.Projects.GetContainsAiCode(ctx, &ProjectsProjectOptions{})
		require.Error(t, err)
	})

	t.Run("Qualitygates", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/qualitygates/application_status":
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"status":"NONE","projects":[],"metrics":[]}`))
			case "/qualitygates/set_ai_code_assurance":
				require.NoError(t, r.ParseForm())
				assert.Equal(t, "true", r.Form.Get("aiCodeAssurance"))
				w.WriteHeader(http.StatusNoContent)
			default:
				t.Errorf("unexpected path %s", r.URL.Path)
			}
		})
		client := newTestClient(t, server.URL)

		status, _, err := client.Qualitygates.ApplicationStatus(ctx, &QualitygatesApplicationStatusOptions{Application: "a"})
		require.NoError(t, err)
		assert.Equal(t, "NONE", status.Status)

		_, err = client.Qualitygates.SetAiCodeAssurance(ctx, &QualitygatesSetAiCodeAssuranceOptions{GateName: "g", AiCodeAssurance: true})
		require.NoError(t, err)

		_, _, err = client.Qualitygates.ApplicationStatus(ctx, &QualitygatesApplicationStatusOptions{})
		require.Error(t, err)
	})

	t.Run("ProjectBadges.AiCodeAssurance", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/project_badges/ai_code_assurance", r.URL.Path)
			w.Header().Set("Content-Type", "image/svg+xml")
			_, _ = w.Write([]byte(`<svg></svg>`))
		})
		client := newTestClient(t, server.URL)

		result, _, err := client.ProjectBadges.AiCodeAssurance(ctx, &ProjectBadgesAiCodeAssuranceOptions{Project: "p"})
		require.NoError(t, err)
		assert.Equal(t, "<svg></svg>", *result)
	})
}
