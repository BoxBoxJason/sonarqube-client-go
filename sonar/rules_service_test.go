package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRules_App(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"canWrite":true,"languages":{"java":"Java"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Rules.App()
	if err != nil {
		t.Fatalf("App failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || !result.CanWrite {
		t.Error("expected CanWrite to be true")
	}
}

func TestRules_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rule":{"key":"java:MyRule","name":"My Rule"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesCreateOption{
		CustomKey:           "MyRule",
		Name:                "My Rule",
		MarkdownDescription: "Test description",
		TemplateKey:         "java:TemplateRule",
	}
	result, resp, err := client.Rules.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Rule.Key != "java:MyRule" {
		t.Error("unexpected rule key")
	}
}

func TestRules_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesDeleteOption{Key: "java:MyRule"}
	resp, err := client.Rules.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestRules_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"paging":{"total":0},"rules":[],"actives":{}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesSearchOption{Languages: []string{"java"}}
	result, resp, err := client.Rules.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

// TestRulesSearchResponse_DynamicActives verifies that the Actives field
// correctly handles dynamic rule keys as a map[string][]RuleActivation.
func TestRulesSearchResponse_DynamicActives(t *testing.T) {
	jsonData := `{
		"actives": {
			"squid:S1067": [
				{"qProfile": "profile1", "severity": "MAJOR"}
			],
			"squid:ClassCyclomaticComplexity": [
				{"qProfile": "profile2", "severity": "CRITICAL"}
			],
			"custom:MyRule": [
				{"qProfile": "profile3", "severity": "INFO"}
			]
		},
		"paging": {"pageIndex": 1, "pageSize": 10, "total": 3},
		"rules": []
	}`

	var response RulesSearchResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(response.Actives) != 3 {
		t.Errorf("expected 3 active rule keys, got %d", len(response.Actives))
	}

	// Verify dynamic keys exist
	if _, exists := response.Actives["squid:S1067"]; !exists {
		t.Error("expected 'squid:S1067' in actives")
	}

	if _, exists := response.Actives["squid:ClassCyclomaticComplexity"]; !exists {
		t.Error("expected 'squid:ClassCyclomaticComplexity' in actives")
	}

	if _, exists := response.Actives["custom:MyRule"]; !exists {
		t.Error("expected 'custom:MyRule' in actives")
	}

	// Verify activation details
	if len(response.Actives["squid:S1067"]) != 1 {
		t.Error("expected 1 activation for squid:S1067")
	}

	if response.Actives["squid:S1067"][0].QProfile != "profile1" {
		t.Errorf("expected qProfile 'profile1', got '%s'",
			response.Actives["squid:S1067"][0].QProfile)
	}
}

func TestRules_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rule":{"key":"java:S1067","name":"Test Rule"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesShowOption{Key: "java:S1067"}
	result, resp, err := client.Rules.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Rule.Key != "java:S1067" {
		t.Error("unexpected rule key")
	}
}

func TestRules_Tags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"tags":["security","bug"]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesTagsOption{PageSize: 100}
	result, resp, err := client.Rules.Tags(opt)
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Tags) != 2 {
		t.Error("expected 2 tags")
	}
}

func TestRules_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rule":{"key":"java:MyRule","name":"Updated Rule"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesUpdateOption{Key: "java:MyRule", Name: "Updated Rule"}
	result, resp, err := client.Rules.Update(opt)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Rule.Name != "Updated Rule" {
		t.Error("unexpected rule name")
	}
}

func TestRules_Repositories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"repositories":[{"key":"java","language":"java","name":"SonarAnalyzer"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesRepositoriesOption{}
	result, resp, err := client.Rules.Repositories(opt)
	if err != nil {
		t.Fatalf("Repositories failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}
}

func TestRules_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rules":[]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesListOption{PageSize: "50"}
	_, resp, err := client.Rules.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

// Validation Tests

func TestValidateCreateOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *RulesCreateOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name: "missing CustomKey",
			opt: &RulesCreateOption{
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
			},
			wantErr: true,
			errMsg:  "CustomKey",
		},
		{
			name: "missing Name",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
			},
			wantErr: true,
			errMsg:  "Name",
		},
		{
			name: "missing MarkdownDescription",
			opt: &RulesCreateOption{
				CustomKey:   "MyRule",
				Name:        "Test",
				TemplateKey: "java:Template",
			},
			wantErr: true,
			errMsg:  "MarkdownDescription",
		},
		{
			name: "missing TemplateKey",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test",
				MarkdownDescription: "Test",
			},
			wantErr: true,
			errMsg:  "TemplateKey",
		},
		{
			name: "CustomKey too long",
			opt: &RulesCreateOption{
				CustomKey:           strings.Repeat("a", 201),
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
			},
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		{
			name: "invalid Severity",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
				Severity:            "INVALID",
			},
			wantErr: true,
			errMsg:  "Severity",
		},
		{
			name: "invalid Status",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
				Status:              "INVALID",
			},
			wantErr: true,
			errMsg:  "Status",
		},
		{
			name: "invalid Type",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
				Type:                "INVALID",
			},
			wantErr: true,
			errMsg:  "Type",
		},
		{
			name: "invalid Impacts key",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
				Impacts: map[string]string{
					"INVALID": "HIGH",
				},
			},
			wantErr: true,
			errMsg:  "Impacts",
		},
		{
			name: "invalid Impacts value",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test",
				MarkdownDescription: "Test",
				TemplateKey:         "java:Template",
				Impacts: map[string]string{
					"MAINTAINABILITY": "INVALID",
				},
			},
			wantErr: true,
			errMsg:  "Impacts",
		},
		{
			name: "valid option",
			opt: &RulesCreateOption{
				CustomKey:           "MyRule",
				Name:                "Test Rule",
				MarkdownDescription: "Test description",
				TemplateKey:         "java:Template",
				Severity:            "MAJOR",
				Status:              "READY",
				Type:                "BUG",
				Impacts: map[string]string{
					"MAINTAINABILITY": "HIGH",
					"SECURITY":        "LOW",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateCreateOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreateOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateCreateOpt() error = %v, expected to contain %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateSearchOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *RulesSearchOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: false,
		},
		{
			name: "invalid Page",
			opt: &RulesSearchOption{
				PaginationArgs: PaginationArgs{Page: -1},
			},
			wantErr: true,
			errMsg:  "Page",
		},
		{
			name: "invalid PageSize",
			opt: &RulesSearchOption{
				PaginationArgs: PaginationArgs{PageSize: 600},
			},
			wantErr: true,
			errMsg:  "PageSize",
		},
		{
			name: "Q too short",
			opt: &RulesSearchOption{
				Q: "a",
			},
			wantErr: true,
			errMsg:  "Q",
		},
		{
			name: "invalid ActiveSeverities",
			opt: &RulesSearchOption{
				ActiveSeverities: []string{"INVALID"},
			},
			wantErr: true,
			errMsg:  "ActiveSeverities",
		},
		{
			name: "invalid CleanCodeAttributeCategories",
			opt: &RulesSearchOption{
				CleanCodeAttributeCategories: []string{"INVALID"},
			},
			wantErr: true,
			errMsg:  "CleanCodeAttributeCategories",
		},
		{
			name: "invalid ImpactSoftwareQualities",
			opt: &RulesSearchOption{
				ImpactSoftwareQualities: []string{"INVALID"},
			},
			wantErr: true,
			errMsg:  "ImpactSoftwareQualities",
		},
		{
			name: "invalid Inheritance",
			opt: &RulesSearchOption{
				Inheritance: []string{"INVALID"},
			},
			wantErr: true,
			errMsg:  "Inheritance",
		},
		{
			name: "invalid OwaspTop10",
			opt: &RulesSearchOption{
				OwaspTop10: []string{"a11"},
			},
			wantErr: true,
			errMsg:  "OwaspTop10",
		},
		{
			name: "invalid Statuses",
			opt: &RulesSearchOption{
				Statuses: []string{"INVALID"},
			},
			wantErr: true,
			errMsg:  "Statuses",
		},
		{
			name: "invalid Types",
			opt: &RulesSearchOption{
				Types: []string{"INVALID"},
			},
			wantErr: true,
			errMsg:  "Types",
		},
		{
			name: "invalid Sort",
			opt: &RulesSearchOption{
				Sort: "invalid",
			},
			wantErr: true,
			errMsg:  "Sort",
		},
		{
			name: "valid option",
			opt: &RulesSearchOption{
				PaginationArgs: PaginationArgs{Page: 1, PageSize: 50},
				Q:              "test",
				Languages:      []string{"java", "go"},
				Severities:     []string{"MAJOR", "CRITICAL"},
				Statuses:       []string{"READY", "DEPRECATED"},
				Types:          []string{"BUG", "CODE_SMELL"},
				Sort:           "name",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateSearchOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateSearchOpt() error = %v, expected to contain %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateUpdateOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *RulesUpdateOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name:    "missing Key",
			opt:     &RulesUpdateOption{},
			wantErr: true,
			errMsg:  "Key",
		},
		{
			name: "Key too long",
			opt: &RulesUpdateOption{
				Key: strings.Repeat("a", 201),
			},
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		{
			name: "invalid Severity",
			opt: &RulesUpdateOption{
				Key:      "java:MyRule",
				Severity: "INVALID",
			},
			wantErr: true,
			errMsg:  "Severity",
		},
		{
			name: "invalid Status",
			opt: &RulesUpdateOption{
				Key:    "java:MyRule",
				Status: "INVALID",
			},
			wantErr: true,
			errMsg:  "Status",
		},
		{
			name: "invalid RemediationFnType",
			opt: &RulesUpdateOption{
				Key:               "java:MyRule",
				RemediationFnType: "INVALID",
			},
			wantErr: true,
			errMsg:  "RemediationFnType",
		},
		{
			name: "invalid Impacts key",
			opt: &RulesUpdateOption{
				Key: "java:MyRule",
				Impacts: map[string]string{
					"INVALID": "HIGH",
				},
			},
			wantErr: true,
			errMsg:  "Impacts",
		},
		{
			name: "valid option",
			opt: &RulesUpdateOption{
				Key:      "java:MyRule",
				Name:     "Updated Rule",
				Severity: "CRITICAL",
				Status:   "READY",
				Impacts: map[string]string{
					"SECURITY": "HIGH",
				},
				Tags: []string{"security", "bug"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateUpdateOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateUpdateOpt() error = %v, expected to contain %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateShowOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *RulesShowOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name:    "missing Key",
			opt:     &RulesShowOption{},
			wantErr: true,
			errMsg:  "Key",
		},
		{
			name: "valid option",
			opt: &RulesShowOption{
				Key:     "java:S1067",
				Actives: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateShowOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShowOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateShowOpt() error = %v, expected to contain %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateDeleteOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *RulesDeleteOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name:    "missing Key",
			opt:     &RulesDeleteOption{},
			wantErr: true,
			errMsg:  "Key",
		},
		{
			name: "valid option",
			opt: &RulesDeleteOption{
				Key: "java:MyRule",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateDeleteOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDeleteOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateDeleteOpt() error = %v, expected to contain %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateTagsOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *RulesTagsOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: false,
		},
		{
			name: "invalid PageSize - too large",
			opt: &RulesTagsOption{
				PageSize: 600,
			},
			wantErr: true,
			errMsg:  "PageSize",
		},
		{
			name: "invalid PageSize - negative",
			opt: &RulesTagsOption{
				PageSize: -1,
			},
			wantErr: true,
			errMsg:  "PageSize",
		},
		{
			name: "valid option",
			opt: &RulesTagsOption{
				PageSize: 100,
				Q:        "security",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateTagsOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTagsOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateTagsOpt() error = %v, expected to contain %q", err, tt.errMsg)
			}
		})
	}
}

// URL Conversion Tests

func TestConvertCreateOptForURL(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	opt := &RulesCreateOption{
		CustomKey:           "MyRule",
		Name:                "Test Rule",
		MarkdownDescription: "Description",
		TemplateKey:         "java:Template",
		Impacts: map[string]string{
			"MAINTAINABILITY": "HIGH",
			"SECURITY":        "LOW",
		},
		Params: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
	}

	urlOpt := client.Rules.convertCreateOptForURL(opt)

	if urlOpt.CustomKey != opt.CustomKey {
		t.Errorf("CustomKey mismatch: got %q, want %q", urlOpt.CustomKey, opt.CustomKey)
	}

	// Check that impacts are formatted correctly (semicolon-separated)
	if urlOpt.Impacts == "" {
		t.Error("Impacts should not be empty")
	}
	if !strings.Contains(urlOpt.Impacts, "=") || !strings.Contains(urlOpt.Impacts, ";") {
		t.Errorf("Impacts not properly formatted: %q", urlOpt.Impacts)
	}

	// Check that params are formatted correctly
	if urlOpt.Params == "" {
		t.Error("Params should not be empty")
	}
	if !strings.Contains(urlOpt.Params, "=") {
		t.Errorf("Params not properly formatted: %q", urlOpt.Params)
	}
}

func TestConvertSearchOptForURL(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	opt := &RulesSearchOption{
		PaginationArgs: PaginationArgs{Page: 1, PageSize: 50},
		Languages:      []string{"java", "go"},
		Severities:     []string{"MAJOR", "CRITICAL"},
		Tags:           []string{"security", "bug"},
	}

	urlOpt := client.Rules.convertSearchOptForURL(opt)

	if urlOpt.Page != 1 {
		t.Errorf("Page mismatch: got %d, want 1", urlOpt.Page)
	}

	if urlOpt.PageSize != 50 {
		t.Errorf("PageSize mismatch: got %d, want 50", urlOpt.PageSize)
	}

	// Check that slices are formatted as comma-separated
	if urlOpt.Languages == "" {
		t.Error("Languages should not be empty")
	}
	if !strings.Contains(urlOpt.Languages, ",") {
		t.Errorf("Languages not properly formatted: %q", urlOpt.Languages)
	}

	if !strings.Contains(urlOpt.Languages, "java") || !strings.Contains(urlOpt.Languages, "go") {
		t.Errorf("Languages missing values: %q", urlOpt.Languages)
	}
}

func TestConvertUpdateOptForURL(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	opt := &RulesUpdateOption{
		Key:  "java:MyRule",
		Name: "Updated",
		Impacts: map[string]string{
			"SECURITY": "HIGH",
		},
		Params: map[string]string{
			"key1": "val1",
		},
		Tags: []string{"security", "performance"},
	}

	urlOpt := client.Rules.convertUpdateOptForURL(opt)

	if urlOpt.Key != opt.Key {
		t.Errorf("Key mismatch: got %q, want %q", urlOpt.Key, opt.Key)
	}

	// Check that impacts are formatted correctly
	if urlOpt.Impacts == "" {
		t.Error("Impacts should not be empty")
	}
	if !strings.Contains(urlOpt.Impacts, "=") {
		t.Errorf("Impacts not properly formatted: %q", urlOpt.Impacts)
	}

	// Check that tags are formatted as comma-separated
	if urlOpt.Tags == "" {
		t.Error("Tags should not be empty")
	}
	if !strings.Contains(urlOpt.Tags, ",") {
		t.Errorf("Tags not properly formatted: %q", urlOpt.Tags)
	}
}

// Edge Case Tests

func TestPaginationValidation(t *testing.T) {
	tests := []struct {
		name     string
		page     int64
		pageSize int64
		wantErr  bool
	}{
		{"valid pagination", 1, 50, false},
		{"zero values", 0, 0, false},
		{"invalid page", -1, 50, true},
		{"invalid page size - too small", 1, -1, true},
		{"invalid page size - too large", 1, 600, true},
		{"valid max page size", 1, 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePagination(tt.page, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePagination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmptySlicesAndMaps(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test with empty slices and maps
	opt := &RulesSearchOption{
		Languages:  []string{},
		Severities: []string{},
		Tags:       []string{},
	}

	urlOpt := client.Rules.convertSearchOptForURL(opt)

	// Empty slices should result in empty strings
	if urlOpt.Languages != "" {
		t.Error("Empty Languages slice should result in empty string")
	}
	if urlOpt.Severities != "" {
		t.Error("Empty Severities slice should result in empty string")
	}
	if urlOpt.Tags != "" {
		t.Error("Empty Tags slice should result in empty string")
	}
}
