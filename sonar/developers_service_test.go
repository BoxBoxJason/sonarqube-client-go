package sonargo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevelopers_SearchEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/developers/search_events" {
			t.Errorf("expected path /api/developers/search_events, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"events": [
				{
					"category": "NEW_ISSUES",
					"link": "https://sonar.example.com/project?id=my-project",
					"message": "10 new issues",
					"project": "my-project"
				}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DevelopersSearchEventsOption{
		From:     []string{"2017-10-19T13:00:00+0200"},
		Projects: []string{"my-project"},
	}

	result, resp, err := client.Developers.SearchEvents(opt)
	if err != nil {
		t.Fatalf("SearchEvents failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(result.Events))
	}

	event := result.Events[0]
	if event.Category != "NEW_ISSUES" {
		t.Errorf("expected category 'NEW_ISSUES', got %s", event.Category)
	}

	if event.Project != "my-project" {
		t.Errorf("expected project 'my-project', got %s", event.Project)
	}
}

func TestDevelopers_SearchEvents_ValidationError_NilOption(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, _, err = client.Developers.SearchEvents(nil)
	if err == nil {
		t.Fatal("expected error for nil option")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestDevelopers_SearchEvents_ValidationError_MissingFrom(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DevelopersSearchEventsOption{
		Projects: []string{"my-project"},
	}

	_, _, err = client.Developers.SearchEvents(opt)
	if err == nil {
		t.Fatal("expected error for missing From")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if validationErr.Field != "From" {
		t.Errorf("expected field 'From', got '%s'", validationErr.Field)
	}
}

func TestDevelopers_SearchEvents_ValidationError_MissingProjects(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DevelopersSearchEventsOption{
		From: []string{"2017-10-19T13:00:00+0200"},
	}

	_, _, err = client.Developers.SearchEvents(opt)
	if err == nil {
		t.Fatal("expected error for missing Projects")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if validationErr.Field != "Projects" {
		t.Errorf("expected field 'Projects', got '%s'", validationErr.Field)
	}
}

func TestDevelopers_ValidateSearchEventsOpt(t *testing.T) {
	tests := []struct {
		name      string
		opt       *DevelopersSearchEventsOption
		wantErr   bool
		wantField string
	}{
		{
			name:      "nil option",
			opt:       nil,
			wantErr:   true,
			wantField: "opt",
		},
		{
			name: "missing from",
			opt: &DevelopersSearchEventsOption{
				Projects: []string{"my-project"},
			},
			wantErr:   true,
			wantField: "From",
		},
		{
			name: "missing projects",
			opt: &DevelopersSearchEventsOption{
				From: []string{"2017-10-19T13:00:00+0200"},
			},
			wantErr:   true,
			wantField: "Projects",
		},
		{
			name: "valid option",
			opt: &DevelopersSearchEventsOption{
				From:     []string{"2017-10-19T13:00:00+0200"},
				Projects: []string{"my-project"},
			},
			wantErr: false,
		},
		{
			name: "valid with multiple values",
			opt: &DevelopersSearchEventsOption{
				From:     []string{"2017-10-19T13:00:00+0200", "2017-10-20T13:00:00+0200"},
				Projects: []string{"my-project", "other-project"},
			},
			wantErr: false,
		},
	}

	client, _ := NewClient("http://localhost/api/", "user", "pass")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Developers.ValidateSearchEventsOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchEventsOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantField != "" {
				var validationErr *ValidationError
				if errors.As(err, &validationErr) {
					if validationErr.Field != tt.wantField {
						t.Errorf("expected field '%s', got '%s'", tt.wantField, validationErr.Field)
					}
				}
			}
		})
	}
}
