package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLanguages_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/languages/list" {
			t.Errorf("expected path /api/languages/list, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"languages": [
				{"key": "java", "name": "Java"},
				{"key": "go", "name": "Go"},
				{"key": "py", "name": "Python"}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &LanguagesListOption{}

	result, resp, err := client.Languages.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Languages) != 3 {
		t.Errorf("expected 3 languages, got %d", len(result.Languages))
	}

	if result.Languages[0].Key != "java" {
		t.Errorf("expected first language key to be 'java', got %s", result.Languages[0].Key)
	}

	if result.Languages[0].Name != "Java" {
		t.Errorf("expected first language name to be 'Java', got %s", result.Languages[0].Name)
	}
}

func TestLanguages_List_WithQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query != "java" {
			t.Errorf("expected query 'java', got %s", query)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"languages": [
				{"key": "java", "name": "Java"}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &LanguagesListOption{
		Query: "java",
	}

	result, resp, err := client.Languages.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Languages) != 1 {
		t.Errorf("expected 1 language, got %d", len(result.Languages))
	}
}

func TestLanguages_List_NilOption(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"languages": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Languages.List(nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestLanguages_ValidateListOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid
	err := client.Languages.ValidateListOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid
	err = client.Languages.ValidateListOpt(&LanguagesListOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}

	// Option with values should be valid
	err = client.Languages.ValidateListOpt(&LanguagesListOption{
		PageSize: 25,
		Query:    "java",
	})
	if err != nil {
		t.Errorf("expected nil error for valid option, got %v", err)
	}
}
