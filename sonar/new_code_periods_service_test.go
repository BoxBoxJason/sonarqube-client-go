package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCodePeriods_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/new_code_periods/list" {
			t.Errorf("expected path /api/new_code_periods/list, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"newCodePeriods": [
				{
					"projectKey": "my-project",
					"branchKey": "main",
					"type": "PREVIOUS_VERSION",
					"inherited": false
				},
				{
					"projectKey": "my-project",
					"branchKey": "feature-1",
					"type": "NUMBER_OF_DAYS",
					"value": "30",
					"inherited": true
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NewCodePeriodsListOption{
		Project: "my-project",
	}

	result, resp, err := client.NewCodePeriods.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.NewCodePeriods) != 2 {
		t.Errorf("expected 2 new code periods, got %d", len(result.NewCodePeriods))
	}

	if result.NewCodePeriods[0].Type != "PREVIOUS_VERSION" {
		t.Errorf("expected type 'PREVIOUS_VERSION', got %s", result.NewCodePeriods[0].Type)
	}

	if result.NewCodePeriods[1].Value != "30" {
		t.Errorf("expected value '30', got %s", result.NewCodePeriods[1].Value)
	}
}

func TestNewCodePeriods_List_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.NewCodePeriods.List(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, _, err = client.NewCodePeriods.List(&NewCodePeriodsListOption{})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestNewCodePeriods_Set(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/new_code_periods/set" {
			t.Errorf("expected path /api/new_code_periods/set, got %s", r.URL.Path)
		}

		periodType := r.URL.Query().Get("type")
		if periodType != "NUMBER_OF_DAYS" {
			t.Errorf("expected type 'NUMBER_OF_DAYS', got %s", periodType)
		}

		value := r.URL.Query().Get("value")
		if value != "30" {
			t.Errorf("expected value '30', got %s", value)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NewCodePeriodsSetOption{
		Type:  "NUMBER_OF_DAYS",
		Value: "30",
	}

	resp, err := client.NewCodePeriods.Set(opt)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestNewCodePeriods_Set_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.NewCodePeriods.Set(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Type should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOption{})
	if err == nil {
		t.Error("expected error for missing Type")
	}

	// Invalid Type should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOption{
		Type: "INVALID_TYPE",
	})
	if err == nil {
		t.Error("expected error for invalid Type")
	}

	// SPECIFIC_ANALYSIS without Branch should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOption{
		Type: "SPECIFIC_ANALYSIS",
	})
	if err == nil {
		t.Error("expected error for SPECIFIC_ANALYSIS without Branch")
	}

	// REFERENCE_BRANCH without Project should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOption{
		Type: "REFERENCE_BRANCH",
	})
	if err == nil {
		t.Error("expected error for REFERENCE_BRANCH without Project")
	}
}

func TestNewCodePeriods_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/new_code_periods/show" {
			t.Errorf("expected path /api/new_code_periods/show, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"type": "NUMBER_OF_DAYS",
			"inherited": false
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.NewCodePeriods.Show(nil)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Type != "NUMBER_OF_DAYS" {
		t.Errorf("expected type 'NUMBER_OF_DAYS', got %s", result.Type)
	}
}

func TestNewCodePeriods_Show_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "main" {
			t.Errorf("expected branch 'main', got %s", branch)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"projectKey": "my-project",
			"branchKey": "main",
			"type": "REFERENCE_BRANCH",
			"inherited": true
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NewCodePeriodsShowOption{
		Project: "my-project",
		Branch:  "main",
	}

	result, _, err := client.NewCodePeriods.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if result.ProjectKey != "my-project" {
		t.Errorf("expected projectKey 'my-project', got %s", result.ProjectKey)
	}
}

func TestNewCodePeriods_Unset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/new_code_periods/unset" {
			t.Errorf("expected path /api/new_code_periods/unset, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.NewCodePeriods.Unset(nil)
	if err != nil {
		t.Fatalf("Unset failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestNewCodePeriods_Unset_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &NewCodePeriodsUnsetOption{
		Project: "my-project",
	}

	resp, err := client.NewCodePeriods.Unset(opt)
	if err != nil {
		t.Fatalf("Unset failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestNewCodePeriods_ValidateSetOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// All valid types without special requirements should pass.
	validTypes := []string{"PREVIOUS_VERSION"}
	for _, periodType := range validTypes {
		err := client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOption{
			Type: periodType,
		})
		if err != nil {
			t.Errorf("expected nil error for type '%s', got %v", periodType, err)
		}
	}

	// NUMBER_OF_DAYS with valid Value should pass.
	err := client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOption{
		Type:  "NUMBER_OF_DAYS",
		Value: "30",
	})
	if err != nil {
		t.Errorf("expected nil error for NUMBER_OF_DAYS with Value, got %v", err)
	}

	// SPECIFIC_ANALYSIS with Branch should pass.
	err = client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOption{
		Type:   "SPECIFIC_ANALYSIS",
		Branch: "main",
	})
	if err != nil {
		t.Errorf("expected nil error for SPECIFIC_ANALYSIS with Branch, got %v", err)
	}

	// REFERENCE_BRANCH with Project should pass.
	err = client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOption{
		Type:    "REFERENCE_BRANCH",
		Project: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error for REFERENCE_BRANCH with Project, got %v", err)
	}
}
