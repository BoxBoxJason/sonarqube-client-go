package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestL10N_Index(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/l10n/index" {
			t.Errorf("expected path /api/l10n/index, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"locale": "en",
			"messages": {
				"quality_gates.operator.LT": "is less than",
				"quality_gates.operator.GT": "is greater than",
				"projects.no_projects.title": "There are no projects yet",
				"projects.create_project": "Create Project"
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.L10N.Index(nil)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Locale != "en" {
		t.Errorf("expected locale 'en', got %s", result.Locale)
	}

	if len(result.Messages) != 4 {
		t.Errorf("expected 4 messages, got %d", len(result.Messages))
	}

	if result.Messages["quality_gates.operator.LT"] != "is less than" {
		t.Errorf("expected message 'is less than', got %s", result.Messages["quality_gates.operator.LT"])
	}

	if result.Messages["projects.create_project"] != "Create Project" {
		t.Errorf("expected message 'Create Project', got %s", result.Messages["projects.create_project"])
	}
}

func TestL10N_Index_WithLocale(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		if locale != "fr" {
			t.Errorf("expected locale 'fr', got %s", locale)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"locale": "fr",
			"messages": {
				"quality_gates.operator.LT": "est inférieur à",
				"quality_gates.operator.GT": "est supérieur à"
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &L10NIndexOption{
		Locale: "fr",
	}

	result, resp, err := client.L10N.Index(opt)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Locale != "fr" {
		t.Errorf("expected locale 'fr', got %s", result.Locale)
	}

	if result.Messages["quality_gates.operator.LT"] != "est inférieur à" {
		t.Errorf("expected French translation, got %s", result.Messages["quality_gates.operator.LT"])
	}
}

func TestL10N_Index_WithTimestamp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := r.URL.Query().Get("ts")
		if ts != "2024-01-01T00:00:00+0000" {
			t.Errorf("expected ts '2024-01-01T00:00:00+0000', got %s", ts)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"locale": "en", "messages": {}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &L10NIndexOption{
		Ts: "2024-01-01T00:00:00+0000",
	}

	_, _, err = client.L10N.Index(opt)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}
}

func TestL10N_Index_EmptyMessages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"locale": "en", "messages": {}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, _, err := client.L10N.Index(nil)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	if result.Messages == nil {
		t.Error("expected non-nil Messages map")
	}

	if len(result.Messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(result.Messages))
	}
}

func TestL10N_ValidateIndexOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid.
	err := client.L10N.ValidateIndexOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid.
	err = client.L10N.ValidateIndexOpt(&L10NIndexOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}

	// Option with Locale should be valid.
	err = client.L10N.ValidateIndexOpt(&L10NIndexOption{
		Locale: "en",
	})
	if err != nil {
		t.Errorf("expected nil error for option with Locale, got %v", err)
	}

	// Option with Ts should be valid.
	err = client.L10N.ValidateIndexOpt(&L10NIndexOption{
		Ts: "2024-01-01T00:00:00+0000",
	})
	if err != nil {
		t.Errorf("expected nil error for option with Ts, got %v", err)
	}

	// Option with both Locale and Ts should be valid.
	err = client.L10N.ValidateIndexOpt(&L10NIndexOption{
		Locale: "fr",
		Ts:     "2024-01-01T00:00:00+0000",
	})
	if err != nil {
		t.Errorf("expected nil error for option with both Locale and Ts, got %v", err)
	}
}
