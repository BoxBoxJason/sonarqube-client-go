package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectAnalysesService_CreateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("analysis") != "AU-TpxcA-iU5OvuD2FL0" {
				t.Errorf("expected analysis 'AU-TpxcA-iU5OvuD2FL0', got '%s'", r.Form.Get("analysis"))
			}
			if r.Form.Get("name") != "1.0" {
				t.Errorf("expected name '1.0', got '%s'", r.Form.Get("name"))
			}
			if r.Form.Get("category") != "VERSION" {
				t.Errorf("expected category 'VERSION', got '%s'", r.Form.Get("category"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"event": {
					"analysis": "AU-TpxcA-iU5OvuD2FL0",
					"key": "AU-TpxcA-iU5OvuD2FL1",
					"category": "VERSION",
					"name": "1.0"
				}
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
			Category: "VERSION",
			Name:     "1.0",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.CreateEvent(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing analysis", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Name: "1.0",
		})
		if err == nil {
			t.Error("expected error for missing analysis")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
		})
		if err == nil {
			t.Error("expected error for missing name")
		}
	})

	t.Run("invalid category", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
			Category: "INVALID",
			Name:     "1.0",
		})
		if err == nil {
			t.Error("expected error for invalid category")
		}
	})
}

func TestProjectAnalysesService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("analysis") != "AU-TpxcA-iU5OvuD2FL0" {
				t.Errorf("expected analysis 'AU-TpxcA-iU5OvuD2FL0', got '%s'", r.Form.Get("analysis"))
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		resp, err := client.ProjectAnalyses.Delete(&ProjectAnalysesDeleteOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.ProjectAnalyses.Delete(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing analysis", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.ProjectAnalyses.Delete(&ProjectAnalysesDeleteOption{})
		if err == nil {
			t.Error("expected error for missing analysis")
		}
	})
}

func TestProjectAnalysesService_DeleteEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("event") != "AU-TpxcA-iU5OvuD2FL1" {
				t.Errorf("expected event 'AU-TpxcA-iU5OvuD2FL1', got '%s'", r.Form.Get("event"))
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		resp, err := client.ProjectAnalyses.DeleteEvent(&ProjectAnalysesDeleteEventOption{
			Event: "AU-TpxcA-iU5OvuD2FL1",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.ProjectAnalyses.DeleteEvent(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing event", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, err := client.ProjectAnalyses.DeleteEvent(&ProjectAnalysesDeleteEventOption{})
		if err == nil {
			t.Error("expected error for missing event")
		}
	})
}

func TestProjectAnalysesService_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Query().Get("project") != "my-project" {
				t.Errorf("expected project 'my-project', got '%s'", r.URL.Query().Get("project"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"paging": {
					"pageIndex": 1,
					"pageSize": 100,
					"total": 1
				},
				"analyses": [
					{
						"key": "AU-TpxcA-iU5OvuD2FL0",
						"date": "2022-01-15T10:00:00+0000",
						"projectVersion": "1.0",
						"revision": "abc123",
						"events": [
							{
								"key": "AU-TpxcA-iU5OvuD2FL1",
								"category": "VERSION",
								"name": "1.0"
							}
						]
					}
				]
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if len(result.Analyses) != 1 {
			t.Errorf("expected 1 analysis, got %d", len(result.Analyses))
		}
		if result.Analyses[0].Key != "AU-TpxcA-iU5OvuD2FL0" {
			t.Errorf("expected key 'AU-TpxcA-iU5OvuD2FL0', got '%s'", result.Analyses[0].Key)
		}
		if len(result.Analyses[0].Events) != 1 {
			t.Errorf("expected 1 event, got %d", len(result.Analyses[0].Events))
		}
	})

	t.Run("with filters", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("project") != "my-project" {
				t.Errorf("expected project 'my-project', got '%s'", r.URL.Query().Get("project"))
			}
			if r.URL.Query().Get("branch") != "main" {
				t.Errorf("expected branch 'main', got '%s'", r.URL.Query().Get("branch"))
			}
			if r.URL.Query().Get("category") != "VERSION" {
				t.Errorf("expected category 'VERSION', got '%s'", r.URL.Query().Get("category"))
			}
			if r.URL.Query().Get("from") != "2022-01-01" {
				t.Errorf("expected from '2022-01-01', got '%s'", r.URL.Query().Get("from"))
			}
			if r.URL.Query().Get("to") != "2022-12-31" {
				t.Errorf("expected to '2022-12-31', got '%s'", r.URL.Query().Get("to"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}, "analyses": []}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project:  "my-project",
			Branch:   "main",
			Category: "VERSION",
			From:     "2022-01-01",
			To:       "2022-12-31",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with datetime format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}, "analyses": []}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
			From:    "2022-01-01T00:00:00Z",
			To:      "2022-12-31T23:59:59Z",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing project", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{})
		if err == nil {
			t.Error("expected error for missing project")
		}
	})

	t.Run("invalid category", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project:  "my-project",
			Category: "INVALID",
		})
		if err == nil {
			t.Error("expected error for invalid category")
		}
	})

	t.Run("invalid from date", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
			From:    "invalid-date",
		})
		if err == nil {
			t.Error("expected error for invalid from date")
		}
	})

	t.Run("invalid to date", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
			To:      "invalid-date",
		})
		if err == nil {
			t.Error("expected error for invalid to date")
		}
	})
}

func TestProjectAnalysesService_SearchAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{
					"paging": {"pageIndex": 1, "pageSize": 1, "total": 2},
					"analyses": [{"key": "analysis1"}]
				}`))
			} else {
				_, _ = w.Write([]byte(`{
					"paging": {"pageIndex": 2, "pageSize": 1, "total": 2},
					"analyses": [{"key": "analysis2"}]
				}`))
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		opt := &ProjectAnalysesSearchOption{
			Project: "my-project",
		}
		opt.PageSize = 1

		result, _, err := client.ProjectAnalyses.SearchAll(opt)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 analyses, got %d", len(result))
		}
		if callCount != 2 {
			t.Errorf("expected 2 calls, got %d", callCount)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.SearchAll(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})
}

func TestProjectAnalysesService_UpdateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.Form.Get("event") != "AU-TpxcA-iU5OvuD2FL1" {
				t.Errorf("expected event 'AU-TpxcA-iU5OvuD2FL1', got '%s'", r.Form.Get("event"))
			}
			if r.Form.Get("name") != "2.0" {
				t.Errorf("expected name '2.0', got '%s'", r.Form.Get("name"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"event": {
					"key": "AU-TpxcA-iU5OvuD2FL1",
					"category": "VERSION",
					"name": "2.0"
				}
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.ProjectAnalyses.UpdateEvent(&ProjectAnalysesUpdateEventOption{
			Event: "AU-TpxcA-iU5OvuD2FL1",
			Name:  "2.0",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("nil option", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.UpdateEvent(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing event", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.UpdateEvent(&ProjectAnalysesUpdateEventOption{
			Name: "2.0",
		})
		if err == nil {
			t.Error("expected error for missing event")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.ProjectAnalyses.UpdateEvent(&ProjectAnalysesUpdateEventOption{
			Event: "AU-TpxcA-iU5OvuD2FL1",
		})
		if err == nil {
			t.Error("expected error for missing name")
		}
	})
}

func TestProjectAnalysesService_ValidateCreateEventOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *ProjectAnalysesCreateEventOption
		wantErr bool
	}{
		{"valid minimal", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0"}, false},
		{"valid with VERSION", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0", Category: "VERSION"}, false},
		{"valid with OTHER", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0", Category: "OTHER"}, false},
		{"nil option", nil, true},
		{"missing analysis", &ProjectAnalysesCreateEventOption{Name: "1.0"}, true},
		{"missing name", &ProjectAnalysesCreateEventOption{Analysis: "a1"}, true},
		{"invalid category", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0", Category: "INVALID"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateCreateEventOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreateEventOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateDeleteOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *ProjectAnalysesDeleteOption
		wantErr bool
	}{
		{"valid", &ProjectAnalysesDeleteOption{Analysis: "a1"}, false},
		{"nil option", nil, true},
		{"empty analysis", &ProjectAnalysesDeleteOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateDeleteOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDeleteOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateDeleteEventOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *ProjectAnalysesDeleteEventOption
		wantErr bool
	}{
		{"valid", &ProjectAnalysesDeleteEventOption{Event: "e1"}, false},
		{"nil option", nil, true},
		{"empty event", &ProjectAnalysesDeleteEventOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateDeleteEventOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDeleteEventOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateSearchOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *ProjectAnalysesSearchOption
		wantErr bool
	}{
		{"valid minimal", &ProjectAnalysesSearchOption{Project: "p1"}, false},
		{"valid with VERSION", &ProjectAnalysesSearchOption{Project: "p1", Category: "VERSION"}, false},
		{"valid with QUALITY_GATE", &ProjectAnalysesSearchOption{Project: "p1", Category: "QUALITY_GATE"}, false},
		{"valid with date", &ProjectAnalysesSearchOption{Project: "p1", From: "2022-01-01"}, false},
		{"valid with datetime", &ProjectAnalysesSearchOption{Project: "p1", From: "2022-01-01T00:00:00Z"}, false},
		{"nil option", nil, true},
		{"missing project", &ProjectAnalysesSearchOption{}, true},
		{"invalid category", &ProjectAnalysesSearchOption{Project: "p1", Category: "INVALID"}, true},
		{"invalid from", &ProjectAnalysesSearchOption{Project: "p1", From: "bad"}, true},
		{"invalid to", &ProjectAnalysesSearchOption{Project: "p1", To: "bad"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateSearchOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateUpdateEventOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *ProjectAnalysesUpdateEventOption
		wantErr bool
	}{
		{"valid", &ProjectAnalysesUpdateEventOption{Event: "e1", Name: "2.0"}, false},
		{"nil option", nil, true},
		{"missing event", &ProjectAnalysesUpdateEventOption{Name: "2.0"}, true},
		{"missing name", &ProjectAnalysesUpdateEventOption{Event: "e1"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateUpdateEventOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateEventOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"2022-01-15", true},
		{"2022-12-31", true},
		{"22-01-15", false},
		{"2022/01/15", false},
		{"2022-1-15", false},
		{"2022-01-1", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isValidDate(tt.input)
			if got != tt.want {
				t.Errorf("isValidDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidDateTime(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"2022-01-15T10:30:00Z", true},
		{"2022-12-31T23:59:59+0000", true},
		{"2022-01-15", false},
		{"invalid", false},
		{"", false},
		{"2022-01-15 10:30:00", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isValidDateTime(tt.input)
			if got != tt.want {
				t.Errorf("isValidDateTime(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
