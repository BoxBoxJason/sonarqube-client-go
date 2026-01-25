package sonargo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPush_SonarlintEvents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/push/sonarlint_events" {
			t.Errorf("expected path /api/push/sonarlint_events, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PushSonarlintEventsOption{
		Languages:   []string{"java", "go"},
		ProjectKeys: []string{"my-project"},
	}

	resp, err := client.Push.SonarlintEvents(opt)
	if err != nil {
		t.Fatalf("SonarlintEvents failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestPush_SonarlintEvents_ValidationError_NilOption(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Push.SonarlintEvents(nil)
	if err == nil {
		t.Fatal("expected error for nil option")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestPush_SonarlintEvents_ValidationError_MissingLanguages(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PushSonarlintEventsOption{
		ProjectKeys: []string{"my-project"},
	}

	_, err = client.Push.SonarlintEvents(opt)
	if err == nil {
		t.Fatal("expected error for missing Languages")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if validationErr.Field != "Languages" {
		t.Errorf("expected field 'Languages', got '%s'", validationErr.Field)
	}
}

func TestPush_SonarlintEvents_ValidationError_MissingProjectKeys(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &PushSonarlintEventsOption{
		Languages: []string{"java"},
	}

	_, err = client.Push.SonarlintEvents(opt)
	if err == nil {
		t.Fatal("expected error for missing ProjectKeys")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if validationErr.Field != "ProjectKeys" {
		t.Errorf("expected field 'ProjectKeys', got '%s'", validationErr.Field)
	}
}

func TestPush_ValidateSonarlintEventsOpt(t *testing.T) {
	tests := []struct {
		name      string
		opt       *PushSonarlintEventsOption
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
			name: "missing languages",
			opt: &PushSonarlintEventsOption{
				ProjectKeys: []string{"my-project"},
			},
			wantErr:   true,
			wantField: "Languages",
		},
		{
			name: "missing project keys",
			opt: &PushSonarlintEventsOption{
				Languages: []string{"java"},
			},
			wantErr:   true,
			wantField: "ProjectKeys",
		},
		{
			name: "valid option",
			opt: &PushSonarlintEventsOption{
				Languages:   []string{"java", "go"},
				ProjectKeys: []string{"my-project"},
			},
			wantErr: false,
		},
	}

	client, _ := NewClient("http://localhost/api/", "user", "pass")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Push.ValidateSonarlintEventsOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSonarlintEventsOpt() error = %v, wantErr %v", err, tt.wantErr)
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
