package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlmIntegrations_CheckPat(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsCheckPatOption{}
	resp, err := client.AlmIntegrations.CheckPat(opt)
	if err != nil {
		t.Fatalf("CheckPat failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_GetGithubClientId(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsGetGithubClientIdOption{}
	_, resp, err := client.AlmIntegrations.GetGithubClientId(opt)
	if err != nil {
		t.Fatalf("GetGithubClientId failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportAzureProject(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsImportAzureProjectOption{}
	resp, err := client.AlmIntegrations.ImportAzureProject(opt)
	if err != nil {
		t.Fatalf("ImportAzureProject failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportBitbucketcloudRepo(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsImportBitbucketcloudRepoOption{}
	resp, err := client.AlmIntegrations.ImportBitbucketcloudRepo(opt)
	if err != nil {
		t.Fatalf("ImportBitbucketcloudRepo failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportBitbucketserverProject(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsImportBitbucketserverProjectOption{}
	resp, err := client.AlmIntegrations.ImportBitbucketserverProject(opt)
	if err != nil {
		t.Fatalf("ImportBitbucketserverProject failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportGithubProject(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsImportGithubProjectOption{}
	resp, err := client.AlmIntegrations.ImportGithubProject(opt)
	if err != nil {
		t.Fatalf("ImportGithubProject failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ImportGitlabProject(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsImportGitlabProjectOption{}
	resp, err := client.AlmIntegrations.ImportGitlabProject(opt)
	if err != nil {
		t.Fatalf("ImportGitlabProject failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ListAzureProjects(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsListAzureProjectsOption{}
	_, resp, err := client.AlmIntegrations.ListAzureProjects(opt)
	if err != nil {
		t.Fatalf("ListAzureProjects failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ListBitbucketserverProjects(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsListBitbucketserverProjectsOption{}
	_, resp, err := client.AlmIntegrations.ListBitbucketserverProjects(opt)
	if err != nil {
		t.Fatalf("ListBitbucketserverProjects failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ListGithubOrganizations(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsListGithubOrganizationsOption{}
	_, resp, err := client.AlmIntegrations.ListGithubOrganizations(opt)
	if err != nil {
		t.Fatalf("ListGithubOrganizations failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_ListGithubRepositories(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsListGithubRepositoriesOption{}
	_, resp, err := client.AlmIntegrations.ListGithubRepositories(opt)
	if err != nil {
		t.Fatalf("ListGithubRepositories failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SearchAzureRepos(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsSearchAzureReposOption{}
	_, resp, err := client.AlmIntegrations.SearchAzureRepos(opt)
	if err != nil {
		t.Fatalf("SearchAzureRepos failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SearchBitbucketcloudRepos(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsSearchBitbucketcloudReposOption{}
	_, resp, err := client.AlmIntegrations.SearchBitbucketcloudRepos(opt)
	if err != nil {
		t.Fatalf("SearchBitbucketcloudRepos failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SearchBitbucketserverRepos(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsSearchBitbucketserverReposOption{}
	_, resp, err := client.AlmIntegrations.SearchBitbucketserverRepos(opt)
	if err != nil {
		t.Fatalf("SearchBitbucketserverRepos failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SearchGitlabRepos(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("null"))
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsSearchGitlabReposOption{}
	_, resp, err := client.AlmIntegrations.SearchGitlabRepos(opt)
	if err != nil {
		t.Fatalf("SearchGitlabRepos failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmIntegrations_SetPat(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		// Return mock response
		w.WriteHeader(204)
	}))
	defer ts.Close()
	// Create client pointing to mock server
	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Call service method
	opt := &AlmIntegrationsSetPatOption{}
	resp, err := client.AlmIntegrations.SetPat(opt)
	if err != nil {
		t.Fatalf("SetPat failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}
