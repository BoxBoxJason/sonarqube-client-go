package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectBadges_Measure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_badges/measure" {
			t.Errorf("expected path /api/project_badges/measure, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		metric := r.URL.Query().Get("metric")
		if metric != "coverage" {
			t.Errorf("expected metric 'coverage', got %s", metric)
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<svg>badge content</svg>`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBadgesMeasureOption{
		Project: "my-project",
		Metric:  "coverage",
	}

	result, resp, err := client.ProjectBadges.Measure(opt)
	if err != nil {
		t.Fatalf("Measure failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if *result != "<svg>badge content</svg>" {
		t.Errorf("expected SVG content, got %s", *result)
	}
}

func TestProjectBadges_Measure_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectBadges.Measure(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, _, err = client.ProjectBadges.Measure(&ProjectBadgesMeasureOption{
		Metric: "coverage",
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}

	// Missing Metric should fail validation.
	_, _, err = client.ProjectBadges.Measure(&ProjectBadgesMeasureOption{
		Project: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing Metric")
	}

	// Invalid Metric should fail validation.
	_, _, err = client.ProjectBadges.Measure(&ProjectBadgesMeasureOption{
		Project: "my-project",
		Metric:  "invalid_metric",
	})
	if err == nil {
		t.Error("expected error for invalid Metric")
	}
}

func TestProjectBadges_QualityGate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_badges/quality_gate" {
			t.Errorf("expected path /api/project_badges/quality_gate, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<svg>quality gate badge</svg>`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBadgesQualityGateOption{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBadges.QualityGate(opt)
	if err != nil {
		t.Fatalf("QualityGate failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestProjectBadges_QualityGate_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectBadges.QualityGate(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, _, err = client.ProjectBadges.QualityGate(&ProjectBadgesQualityGateOption{})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBadges_RenewToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_badges/renew_token" {
			t.Errorf("expected path /api/project_badges/renew_token, got %s", r.URL.Path)
		}

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

	opt := &ProjectBadgesRenewTokenOption{
		Project: "my-project",
	}

	resp, err := client.ProjectBadges.RenewToken(opt)
	if err != nil {
		t.Fatalf("RenewToken failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestProjectBadges_RenewToken_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.ProjectBadges.RenewToken(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, err = client.ProjectBadges.RenewToken(&ProjectBadgesRenewTokenOption{})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBadges_Token(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/project_badges/token" {
			t.Errorf("expected path /api/project_badges/token, got %s", r.URL.Path)
		}

		project := r.URL.Query().Get("project")
		if project != "my-project" {
			t.Errorf("expected project 'my-project', got %s", project)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token": "abc123def456"}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &ProjectBadgesTokenOption{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBadges.Token(opt)
	if err != nil {
		t.Fatalf("Token failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Token != "abc123def456" {
		t.Errorf("expected token 'abc123def456', got %s", result.Token)
	}
}

func TestProjectBadges_Token_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.ProjectBadges.Token(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Project should fail validation.
	_, _, err = client.ProjectBadges.Token(&ProjectBadgesTokenOption{})
	if err == nil {
		t.Error("expected error for missing Project")
	}
}

func TestProjectBadges_ValidateMeasureOpt_AllMetrics(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	validMetrics := []string{
		"coverage",
		"duplicated_lines_density",
		"ncloc",
		"alert_status",
		"security_hotspots",
		"bugs",
		"code_smells",
		"vulnerabilities",
		"sqale_rating",
		"reliability_rating",
		"security_rating",
		"sqale_index",
		"software_quality_reliability_issues",
		"software_quality_maintainability_issues",
		"software_quality_security_issues",
		"software_quality_maintainability_rating",
		"software_quality_reliability_rating",
		"software_quality_security_rating",
		"software_quality_maintainability_remediation_effort",
	}

	for _, metric := range validMetrics {
		err := client.ProjectBadges.ValidateMeasureOpt(&ProjectBadgesMeasureOption{
			Project: "my-project",
			Metric:  metric,
		})
		if err != nil {
			t.Errorf("expected nil error for metric '%s', got %v", metric, err)
		}
	}
}
