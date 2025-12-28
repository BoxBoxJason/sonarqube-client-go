package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlmSettings_CountBinding(t *testing.T) {
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
	opt := &AlmSettingsCountBindingOption{}
	_, resp, err := client.AlmSettings.CountBinding(opt)
	if err != nil {
		t.Fatalf("CountBinding failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateAzure(t *testing.T) {
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
	opt := &AlmSettingsCreateAzureOption{}
	resp, err := client.AlmSettings.CreateAzure(opt)
	if err != nil {
		t.Fatalf("CreateAzure failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateBitbucket(t *testing.T) {
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
	opt := &AlmSettingsCreateBitbucketOption{}
	resp, err := client.AlmSettings.CreateBitbucket(opt)
	if err != nil {
		t.Fatalf("CreateBitbucket failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateBitbucketcloud(t *testing.T) {
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
	opt := &AlmSettingsCreateBitbucketcloudOption{}
	resp, err := client.AlmSettings.CreateBitbucketcloud(opt)
	if err != nil {
		t.Fatalf("CreateBitbucketcloud failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateGithub(t *testing.T) {
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
	opt := &AlmSettingsCreateGithubOption{}
	resp, err := client.AlmSettings.CreateGithub(opt)
	if err != nil {
		t.Fatalf("CreateGithub failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_CreateGitlab(t *testing.T) {
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
	opt := &AlmSettingsCreateGitlabOption{}
	resp, err := client.AlmSettings.CreateGitlab(opt)
	if err != nil {
		t.Fatalf("CreateGitlab failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_Delete(t *testing.T) {
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
	opt := &AlmSettingsDeleteOption{}
	resp, err := client.AlmSettings.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_GetBinding(t *testing.T) {
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
	opt := &AlmSettingsGetBindingOption{}
	_, resp, err := client.AlmSettings.GetBinding(opt)
	if err != nil {
		t.Fatalf("GetBinding failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_List(t *testing.T) {
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
	opt := &AlmSettingsListOption{}
	_, resp, err := client.AlmSettings.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_ListDefinitions(t *testing.T) {
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
	_, resp, err := client.AlmSettings.ListDefinitions()
	if err != nil {
		t.Fatalf("ListDefinitions failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateAzure(t *testing.T) {
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
	opt := &AlmSettingsUpdateAzureOption{}
	resp, err := client.AlmSettings.UpdateAzure(opt)
	if err != nil {
		t.Fatalf("UpdateAzure failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateBitbucket(t *testing.T) {
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
	opt := &AlmSettingsUpdateBitbucketOption{}
	resp, err := client.AlmSettings.UpdateBitbucket(opt)
	if err != nil {
		t.Fatalf("UpdateBitbucket failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateBitbucketcloud(t *testing.T) {
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
	opt := &AlmSettingsUpdateBitbucketcloudOption{}
	resp, err := client.AlmSettings.UpdateBitbucketcloud(opt)
	if err != nil {
		t.Fatalf("UpdateBitbucketcloud failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateGithub(t *testing.T) {
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
	opt := &AlmSettingsUpdateGithubOption{}
	resp, err := client.AlmSettings.UpdateGithub(opt)
	if err != nil {
		t.Fatalf("UpdateGithub failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_UpdateGitlab(t *testing.T) {
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
	opt := &AlmSettingsUpdateGitlabOption{}
	resp, err := client.AlmSettings.UpdateGitlab(opt)
	if err != nil {
		t.Fatalf("UpdateGitlab failed: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestAlmSettings_Validate(t *testing.T) {
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
	opt := &AlmSettingsValidateOption{}
	_, resp, err := client.AlmSettings.Validate(opt)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
