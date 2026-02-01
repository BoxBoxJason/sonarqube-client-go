package sonargo

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRules_App(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "GET", "/rules/app", 200, `{"canWrite":true,"languages":{"java":"Java"}}`))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Rules.App()
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	assert.True(t, result.CanWrite)
}

func TestRules_Create(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "POST", "/rules/create", 200, `{"rule":{"key":"java:MyRule","name":"My Rule"}}`))
	client := newTestClient(t, server.URL)

	opt := &RulesCreateOption{
		CustomKey:           "MyRule",
		Name:                "My Rule",
		MarkdownDescription: "Test description",
		TemplateKey:         "java:TemplateRule",
	}
	result, resp, err := client.Rules.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "java:MyRule", result.Rule.Key)
}

func TestRules_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, "POST", "/rules/delete", 204))
	client := newTestClient(t, server.URL)

	opt := &RulesDeleteOption{Key: "java:MyRule"}
	resp, err := client.Rules.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestRules_Search(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "GET", "/rules/search", 200, `{"paging":{"total":0},"rules":[],"actives":{}}`))
	client := newTestClient(t, server.URL)

	opt := &RulesSearchOption{Languages: []string{"java"}}
	result, resp, err := client.Rules.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
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

	var response RulesSearch
	err := json.Unmarshal([]byte(jsonData), &response)
	require.NoError(t, err)

	assert.Len(t, response.Actives, 3)

	// Verify dynamic keys exist
	assert.Contains(t, response.Actives, "squid:S1067")
	assert.Contains(t, response.Actives, "squid:ClassCyclomaticComplexity")
	assert.Contains(t, response.Actives, "custom:MyRule")

	// Verify activation details
	assert.Len(t, response.Actives["squid:S1067"], 1)
	assert.Equal(t, "profile1", response.Actives["squid:S1067"][0].QProfile)
}

func TestRules_Show(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "GET", "/rules/show", 200, `{"rule":{"key":"java:S1067","name":"Test Rule"}}`))
	client := newTestClient(t, server.URL)

	opt := &RulesShowOption{Key: "java:S1067"}
	result, resp, err := client.Rules.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "java:S1067", result.Rule.Key)
}

func TestRules_Tags(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "GET", "/rules/tags", 200, `{"tags":["security","bug"]}`))
	client := newTestClient(t, server.URL)

	opt := &RulesTagsOption{PageSize: 100}
	result, resp, err := client.Rules.Tags(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Tags, 2)
}

func TestRules_Update(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "POST", "/rules/update", 200, `{"rule":{"key":"java:MyRule","name":"Updated Rule"}}`))
	client := newTestClient(t, server.URL)

	opt := &RulesUpdateOption{Key: "java:MyRule", Name: "Updated Rule"}
	result, resp, err := client.Rules.Update(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Rule", result.Rule.Name)
}

func TestRules_UpdateClearTags(t *testing.T) {
	var capturedURL string
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"rule":{"key":"java:MyRule","tags":[]}}`))
	}))
	client := newTestClient(t, server.URL)

	opt := &RulesUpdateOption{Key: "java:MyRule", Tags: []string{}}
	result, resp, err := client.Rules.Update(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	// Verify that "tags=" is present in the query string even though Tags is empty
	assert.Contains(t, capturedURL, "tags=")
	assert.Empty(t, result.Rule.Tags)
}

func TestRules_Repositories(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "GET", "/rules/repositories", 200, `{"repositories":[{"key":"java","language":"java","name":"SonarAnalyzer"}]}`))
	client := newTestClient(t, server.URL)

	opt := &RulesRepositoriesOption{}
	result, resp, err := client.Rules.Repositories(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Repositories, 1)
}

func TestRules_List(t *testing.T) {
	server := newTestServer(t, mockHandler(t, "GET", "/rules/list", 200, `{"rules":[]}`))
	client := newTestClient(t, server.URL)

	opt := &RulesListOption{
		PaginationArgs: PaginationArgs{
			PageSize: 50,
		},
	}
	_, resp, err := client.Rules.List(opt)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// Validation Tests

func TestValidateCreateOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSearchOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			name: "Query too short",
			opt: &RulesSearchOption{
				Query: "a",
			},
			wantErr: true,
			errMsg:  "Query",
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
				Query:          "test",
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
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateShowOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDeleteOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTagsOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
				Query:    "security",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Rules.ValidateTagsOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// URL Conversion Tests

func TestConvertCreateOptForURL(t *testing.T) {
	client := newLocalhostClient(t)

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

	assert.Equal(t, opt.CustomKey, urlOpt.CustomKey)

	// Check that impacts are formatted correctly (semicolon-separated)
	assert.NotEmpty(t, urlOpt.Impacts)
	assert.Contains(t, urlOpt.Impacts, "=")
	assert.Contains(t, urlOpt.Impacts, ";")

	// Check that params are formatted correctly
	assert.NotEmpty(t, urlOpt.Params)
	assert.Contains(t, urlOpt.Params, "=")
}

func TestConvertSearchOptForURL(t *testing.T) {
	client := newLocalhostClient(t)

	opt := &RulesSearchOption{
		PaginationArgs: PaginationArgs{Page: 1, PageSize: 50},
		Languages:      []string{"java", "go"},
		Severities:     []string{"MAJOR", "CRITICAL"},
		Tags:           []string{"security", "bug"},
	}

	urlOpt := client.Rules.convertSearchOptForURL(opt)

	assert.Equal(t, int64(1), urlOpt.Page)
	assert.Equal(t, int64(50), urlOpt.PageSize)

	// Check that slices are formatted as comma-separated
	assert.NotEmpty(t, urlOpt.Languages)
	assert.Contains(t, urlOpt.Languages, ",")
	assert.Contains(t, urlOpt.Languages, "java")
	assert.Contains(t, urlOpt.Languages, "go")
}

func TestConvertUpdateOptForURL(t *testing.T) {
	client := newLocalhostClient(t)

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

	assert.Equal(t, opt.Key, urlOpt.Key)

	// Check that impacts are formatted correctly
	assert.NotEmpty(t, urlOpt.Impacts)
	assert.Contains(t, urlOpt.Impacts, "=")

	// Check that tags are formatted as comma-separated
	assert.NotEmpty(t, urlOpt.Tags)
	assert.Contains(t, urlOpt.Tags, ",")
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
		{"valid max page size", MinPageSize, MaxPageSize, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePagination(tt.page, tt.pageSize)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmptySlicesAndMaps(t *testing.T) {
	client := newLocalhostClient(t)

	// Test with empty slices and maps
	opt := &RulesSearchOption{
		Languages:  []string{},
		Severities: []string{},
		Tags:       []string{},
	}

	urlOpt := client.Rules.convertSearchOptForURL(opt)

	// Empty slices should result in empty strings
	assert.Empty(t, urlOpt.Languages)
	assert.Empty(t, urlOpt.Severities)
	assert.Empty(t, urlOpt.Tags)
}
