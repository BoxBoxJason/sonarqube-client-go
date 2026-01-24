package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMonitoring_Metrics(t *testing.T) {
	// Prometheus-format metrics with escaped newlines for JSON encoding
	metricsContent := `"# HELP sonarqube_health Health check status\\n# TYPE sonarqube_health gauge\\nsonarqube_health 1"`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/monitoring/metrics" {
			t.Errorf("expected path /api/monitoring/metrics, got %s", r.URL.Path)
		}

		// Return as JSON-encoded string since the client decodes the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(metricsContent))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Monitoring.Metrics()
	if err != nil {
		t.Fatalf("Metrics failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if *result == "" {
		t.Error("expected non-empty metrics")
	}
}
