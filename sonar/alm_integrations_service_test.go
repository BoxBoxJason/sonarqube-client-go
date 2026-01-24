package sonargo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// CheckPat Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_CheckPat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "alm_integrations/check_pat") {
			t.Errorf("expected path to contain alm_integrations/check_pat, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsCheckPatOption{
		AlmSetting: "my-azure-setting",
	}

	_, resp, err := client.AlmIntegrations.CheckPat(opt)
	if err != nil {
		t.Fatalf("CheckPat failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_CheckPat_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.CheckPat(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.CheckPat(&AlmIntegrationsCheckPatOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test AlmSetting too long
	_, _, err = client.AlmIntegrations.CheckPat(&AlmIntegrationsCheckPatOption{
		AlmSetting: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for AlmSetting exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// GetGithubClientId Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_GetGithubClientId(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"clientId":"my-client-id"}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsGetGithubClientIdOption{
		AlmSetting: "my-github-setting",
	}

	result, resp, err := client.AlmIntegrations.GetGithubClientId(opt)
	if err != nil {
		t.Fatalf("GetGithubClientId failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.ClientID != "my-client-id" {
		t.Error("expected ClientID 'my-client-id'")
	}
}

func TestAlmIntegrations_GetGithubClientId_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.GetGithubClientId(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.GetGithubClientId(&AlmIntegrationsGetGithubClientIdOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}
}

// -----------------------------------------------------------------------------
// ImportAzureProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportAzureProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsImportAzureProjectOption{
		ProjectName:    "my-azure-project",
		RepositoryName: "my-azure-repo",
	}

	resp, err := client.AlmIntegrations.ImportAzureProject(opt)
	if err != nil {
		t.Fatalf("ImportAzureProject failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportAzureProject_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmIntegrations.ImportAzureProject(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing ProjectName
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		RepositoryName: "repo",
	})
	if err == nil {
		t.Error("expected error for missing ProjectName")
	}

	// Test missing RepositoryName
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName: "project",
	})
	if err == nil {
		t.Error("expected error for missing RepositoryName")
	}

	// Test invalid NewCodeDefinitionType
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "INVALID_TYPE",
	})
	if err == nil {
		t.Error("expected error for invalid NewCodeDefinitionType")
	}

	// Test NUMBER_OF_DAYS without value
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "NUMBER_OF_DAYS",
	})
	if err == nil {
		t.Error("expected error for NUMBER_OF_DAYS without value")
	}

	// Test PREVIOUS_VERSION with value
	_, err = client.AlmIntegrations.ImportAzureProject(&AlmIntegrationsImportAzureProjectOption{
		ProjectName:            "project",
		RepositoryName:         "repo",
		NewCodeDefinitionType:  "PREVIOUS_VERSION",
		NewCodeDefinitionValue: 30,
	})
	if err == nil {
		t.Error("expected error for PREVIOUS_VERSION with value")
	}
}

func TestAlmIntegrations_ImportAzureProject_WithNewCodeDefinition(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test with PREVIOUS_VERSION
	opt := &AlmIntegrationsImportAzureProjectOption{
		ProjectName:           "project",
		RepositoryName:        "repo",
		NewCodeDefinitionType: "PREVIOUS_VERSION",
	}

	resp, err := client.AlmIntegrations.ImportAzureProject(opt)
	if err != nil {
		t.Fatalf("ImportAzureProject with PREVIOUS_VERSION failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}

	// Test with NUMBER_OF_DAYS
	opt = &AlmIntegrationsImportAzureProjectOption{
		ProjectName:            "project",
		RepositoryName:         "repo",
		NewCodeDefinitionType:  "NUMBER_OF_DAYS",
		NewCodeDefinitionValue: 30,
	}

	resp, err = client.AlmIntegrations.ImportAzureProject(opt)
	if err != nil {
		t.Fatalf("ImportAzureProject with NUMBER_OF_DAYS failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// -----------------------------------------------------------------------------
// ImportBitbucketCloudRepo Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportBitbucketCloudRepo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsImportBitbucketCloudRepoOption{
		RepositorySlug: "my-repo",
	}

	resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(opt)
	if err != nil {
		t.Fatalf("ImportBitbucketCloudRepo failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportBitbucketCloudRepo_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmIntegrations.ImportBitbucketCloudRepo(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing RepositorySlug
	_, err = client.AlmIntegrations.ImportBitbucketCloudRepo(&AlmIntegrationsImportBitbucketCloudRepoOption{})
	if err == nil {
		t.Error("expected error for missing RepositorySlug")
	}
}

// -----------------------------------------------------------------------------
// ImportBitbucketServerProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportBitbucketServerProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsImportBitbucketServerProjectOption{
		ProjectKey:     "PRJ",
		RepositorySlug: "my-repo",
	}

	resp, err := client.AlmIntegrations.ImportBitbucketServerProject(opt)
	if err != nil {
		t.Fatalf("ImportBitbucketServerProject failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportBitbucketServerProject_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing ProjectKey
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(&AlmIntegrationsImportBitbucketServerProjectOption{
		RepositorySlug: "repo",
	})
	if err == nil {
		t.Error("expected error for missing ProjectKey")
	}

	// Test missing RepositorySlug
	_, err = client.AlmIntegrations.ImportBitbucketServerProject(&AlmIntegrationsImportBitbucketServerProjectOption{
		ProjectKey: "PRJ",
	})
	if err == nil {
		t.Error("expected error for missing RepositorySlug")
	}
}

// -----------------------------------------------------------------------------
// ImportGithubProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportGithubProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsImportGithubProjectOption{
		RepositoryKey: "octocat/hello-world",
	}

	resp, err := client.AlmIntegrations.ImportGithubProject(opt)
	if err != nil {
		t.Fatalf("ImportGithubProject failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportGithubProject_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmIntegrations.ImportGithubProject(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing RepositoryKey
	_, err = client.AlmIntegrations.ImportGithubProject(&AlmIntegrationsImportGithubProjectOption{})
	if err == nil {
		t.Error("expected error for missing RepositoryKey")
	}

	// Test RepositoryKey too long
	_, err = client.AlmIntegrations.ImportGithubProject(&AlmIntegrationsImportGithubProjectOption{
		RepositoryKey: strings.Repeat("a", MaxGitHubRepoKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for RepositoryKey exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// ImportGitlabProject Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ImportGitlabProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsImportGitlabProjectOption{
		GitlabProjectId: "12345",
	}

	resp, err := client.AlmIntegrations.ImportGitlabProject(opt)
	if err != nil {
		t.Fatalf("ImportGitlabProject failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportGitlabProject_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmIntegrations.ImportGitlabProject(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing GitlabProjectId
	_, err = client.AlmIntegrations.ImportGitlabProject(&AlmIntegrationsImportGitlabProjectOption{})
	if err == nil {
		t.Error("expected error for missing GitlabProjectId")
	}
}

// -----------------------------------------------------------------------------
// ListAzureProjects Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListAzureProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"projects":[{"name":"Project1","description":"Description1"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsListAzureProjectsOption{
		AlmSetting: "my-azure-setting",
	}

	result, resp, err := client.AlmIntegrations.ListAzureProjects(opt)
	if err != nil {
		t.Fatalf("ListAzureProjects failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Projects) != 1 {
		t.Error("expected 1 project")
	}
}

func TestAlmIntegrations_ListAzureProjects_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.ListAzureProjects(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListAzureProjects(&AlmIntegrationsListAzureProjectsOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}
}

// -----------------------------------------------------------------------------
// ListBitbucketServerProjects Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListBitbucketServerProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"projects":[{"key":"PRJ","name":"Project1"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsListBitbucketServerProjectsOption{
		AlmSetting: "my-bitbucket-setting",
		PageSize:   25,
	}

	result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(opt)
	if err != nil {
		t.Fatalf("ListBitbucketServerProjects failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Projects) != 1 {
		t.Error("expected 1 project")
	}
}

func TestAlmIntegrations_ListBitbucketServerProjects_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(&AlmIntegrationsListBitbucketServerProjectsOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.ListBitbucketServerProjects(&AlmIntegrationsListBitbucketServerProjectsOption{
		AlmSetting: "setting",
		PageSize:   MaxPageSizeAlmIntegrations + 1,
	})
	if err == nil {
		t.Error("expected error for PageSize out of range")
	}
}

// -----------------------------------------------------------------------------
// ListGithubOrganizations Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListGithubOrganizations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"organizations":[{"key":"octocat","name":"Octocat"}],"paging":{"pageIndex":1,"pageSize":100,"total":1}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsListGithubOrganizationsOption{
		AlmSetting: "my-github-setting",
	}

	result, resp, err := client.AlmIntegrations.ListGithubOrganizations(opt)
	if err != nil {
		t.Fatalf("ListGithubOrganizations failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Organizations) != 1 {
		t.Error("expected 1 organization")
	}
}

func TestAlmIntegrations_ListGithubOrganizations_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(&AlmIntegrationsListGithubOrganizationsOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test Token too long
	_, _, err = client.AlmIntegrations.ListGithubOrganizations(&AlmIntegrationsListGithubOrganizationsOption{
		AlmSetting: "setting",
		Token:      strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for Token exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// ListGithubRepositories Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_ListGithubRepositories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"repositories":[{"id":1,"key":"octocat/hello-world","name":"hello-world","url":"https://github.com/octocat/hello-world"}],"paging":{"pageIndex":1,"pageSize":100,"total":1}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsListGithubRepositoriesOption{
		AlmSetting:   "my-github-setting",
		Organization: "octocat",
	}

	result, resp, err := client.AlmIntegrations.ListGithubRepositories(opt)
	if err != nil {
		t.Fatalf("ListGithubRepositories failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}
}

func TestAlmIntegrations_ListGithubRepositories_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.ListGithubRepositories(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.ListGithubRepositories(&AlmIntegrationsListGithubRepositoriesOption{
		Organization: "octocat",
	})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test missing Organization
	_, _, err = client.AlmIntegrations.ListGithubRepositories(&AlmIntegrationsListGithubRepositoriesOption{
		AlmSetting: "setting",
	})
	if err == nil {
		t.Error("expected error for missing Organization")
	}
}

// -----------------------------------------------------------------------------
// SearchAzureRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchAzureRepos(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"repositories":[{"name":"repo1","projectName":"Project1"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsSearchAzureReposOption{
		AlmSetting: "my-azure-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchAzureRepos(opt)
	if err != nil {
		t.Fatalf("SearchAzureRepos failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}
}

func TestAlmIntegrations_SearchAzureRepos_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.SearchAzureRepos(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchAzureRepos(&AlmIntegrationsSearchAzureReposOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchAzureRepos(&AlmIntegrationsSearchAzureReposOption{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for ProjectName exceeding max length")
	}

	// Test SearchQuery too long
	_, _, err = client.AlmIntegrations.SearchAzureRepos(&AlmIntegrationsSearchAzureReposOption{
		AlmSetting:  "setting",
		SearchQuery: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for SearchQuery exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// SearchBitbucketCloudRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchBitbucketCloudRepos(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"isLastPage":true,"repositories":[{"name":"repo1","slug":"repo1","uuid":"uuid1"}],"paging":{"pageIndex":1,"pageSize":20}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsSearchBitbucketCloudReposOption{
		AlmSetting: "my-bitbucket-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(opt)
	if err != nil {
		t.Fatalf("SearchBitbucketCloudRepos failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}

	if !result.IsLastPage {
		t.Error("expected IsLastPage to be true")
	}
}

func TestAlmIntegrations_SearchBitbucketCloudRepos_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(&AlmIntegrationsSearchBitbucketCloudReposOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(&AlmIntegrationsSearchBitbucketCloudReposOption{
		AlmSetting: "setting",
		PaginationArgs: PaginationArgs{
			PageSize: MaxPageSizeAlmIntegrations + 1,
		},
	})
	if err == nil {
		t.Error("expected error for PageSize out of range")
	}

	// Test RepositoryName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketCloudRepos(&AlmIntegrationsSearchBitbucketCloudReposOption{
		AlmSetting:     "setting",
		RepositoryName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for RepositoryName exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// SearchBitbucketServerRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchBitbucketServerRepos(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"isLastPage":false,"repositories":[{"name":"repo1","slug":"repo1","projectKey":"PRJ"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting: "my-bitbucket-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(opt)
	if err != nil {
		t.Fatalf("SearchBitbucketServerRepos failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}
}

func TestAlmIntegrations_SearchBitbucketServerRepos_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting: "setting",
		PageSize:   MaxPageSizeAlmIntegrations + 1,
	})
	if err == nil {
		t.Error("expected error for PageSize out of range")
	}

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for ProjectName exceeding max length")
	}

	// Test RepositoryName too long
	_, _, err = client.AlmIntegrations.SearchBitbucketServerRepos(&AlmIntegrationsSearchBitbucketServerReposOption{
		AlmSetting:     "setting",
		RepositoryName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for RepositoryName exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// SearchGitlabRepos Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SearchGitlabRepos(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"repositories":[{"id":1,"name":"project1","pathName":"group/project1","slug":"project1","url":"https://gitlab.com/group/project1"}],"paging":{"pageIndex":1,"pageSize":20,"total":1}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsSearchGitlabReposOption{
		AlmSetting: "my-gitlab-setting",
	}

	result, resp, err := client.AlmIntegrations.SearchGitlabRepos(opt)
	if err != nil {
		t.Fatalf("SearchGitlabRepos failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}
}

func TestAlmIntegrations_SearchGitlabRepos_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing AlmSetting
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(&AlmIntegrationsSearchGitlabReposOption{})
	if err == nil {
		t.Error("expected error for missing AlmSetting")
	}

	// Test PageSize out of range
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(&AlmIntegrationsSearchGitlabReposOption{
		AlmSetting: "setting",
		PaginationArgs: PaginationArgs{
			PageSize: MaxPageSizeAlmIntegrations + 1,
		},
	})
	if err == nil {
		t.Error("expected error for PageSize out of range")
	}

	// Test ProjectName too long
	_, _, err = client.AlmIntegrations.SearchGitlabRepos(&AlmIntegrationsSearchGitlabReposOption{
		AlmSetting:  "setting",
		ProjectName: strings.Repeat("a", MaxAlmSettingKeyLength+1),
	})
	if err == nil {
		t.Error("expected error for ProjectName exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// SetPat Tests
// -----------------------------------------------------------------------------

func TestAlmIntegrations_SetPat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &AlmIntegrationsSetPatOption{
		AlmSetting: "my-setting",
		Pat:        "my-personal-access-token",
	}

	resp, err := client.AlmIntegrations.SetPat(opt)
	if err != nil {
		t.Fatalf("SetPat failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SetPat_WithUsername(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test with username (for Bitbucket Cloud)
	opt := &AlmIntegrationsSetPatOption{
		AlmSetting: "my-bitbucket-cloud-setting",
		Pat:        "my-app-password",
		Username:   "my-username",
	}

	resp, err := client.AlmIntegrations.SetPat(opt)
	if err != nil {
		t.Fatalf("SetPat with username failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SetPat_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.AlmIntegrations.SetPat(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Pat
	_, err = client.AlmIntegrations.SetPat(&AlmIntegrationsSetPatOption{
		AlmSetting: "setting",
	})
	if err == nil {
		t.Error("expected error for missing Pat")
	}

	// Test Pat too long
	_, err = client.AlmIntegrations.SetPat(&AlmIntegrationsSetPatOption{
		Pat: strings.Repeat("a", MaxPatLength+1),
	})
	if err == nil {
		t.Error("expected error for Pat exceeding max length")
	}

	// Test Username too long
	_, err = client.AlmIntegrations.SetPat(&AlmIntegrationsSetPatOption{
		Pat:      "token",
		Username: strings.Repeat("a", MaxUsernameLength+1),
	})
	if err == nil {
		t.Error("expected error for Username exceeding max length")
	}
}

// -----------------------------------------------------------------------------
// Helper Function Tests
// -----------------------------------------------------------------------------

func TestValidateNewCodeDefinition(t *testing.T) {
	tests := []struct {
		name            string
		definitionType  string
		definitionValue int64
		wantErr         bool
	}{
		{
			name:            "empty type is valid",
			definitionType:  "",
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "PREVIOUS_VERSION without value",
			definitionType:  "PREVIOUS_VERSION",
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "PREVIOUS_VERSION with value is invalid",
			definitionType:  "PREVIOUS_VERSION",
			definitionValue: 30,
			wantErr:         true,
		},
		{
			name:            "REFERENCE_BRANCH without value",
			definitionType:  "REFERENCE_BRANCH",
			definitionValue: 0,
			wantErr:         false,
		},
		{
			name:            "REFERENCE_BRANCH with value is invalid",
			definitionType:  "REFERENCE_BRANCH",
			definitionValue: 30,
			wantErr:         true,
		},
		{
			name:            "NUMBER_OF_DAYS with value",
			definitionType:  "NUMBER_OF_DAYS",
			definitionValue: 30,
			wantErr:         false,
		},
		{
			name:            "NUMBER_OF_DAYS without value is invalid",
			definitionType:  "NUMBER_OF_DAYS",
			definitionValue: 0,
			wantErr:         true,
		},
		{
			name:            "invalid type",
			definitionType:  "INVALID_TYPE",
			definitionValue: 0,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNewCodeDefinition(tt.definitionType, tt.definitionValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNewCodeDefinition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
