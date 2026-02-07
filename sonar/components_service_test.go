package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// App Tests
// =============================================================================

func TestComponents_App(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/components/app", r.URL.Path)
		assert.Equal(t, "my_project:src/file.go", r.URL.Query().Get("component"))
		assert.Equal(t, "feature/my_branch", r.URL.Query().Get("branch"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"key": "my_project:src/file.go",
			"name": "file.go",
			"longName": "src/file.go",
			"q": "FIL",
			"project": "my_project",
			"projectName": "My Project",
			"fav": true,
			"canMarkAsFavorite": true,
			"measures": {
				"lines": "100",
				"issues": "5",
				"sqaleRating": "A"
			}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Components.App(&ComponentsAppOption{
		Component: "my_project:src/file.go",
		Branch:    "feature/my_branch",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_project:src/file.go", result.Key)
	assert.True(t, result.Fav)
}

func TestComponents_App_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *ComponentsAppOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing component",
			opt:  &ComponentsAppOption{Branch: "main"},
		},
		{
			name: "branch and pullRequest both set",
			opt: &ComponentsAppOption{
				Component:   "my_project",
				Branch:      "main",
				PullRequest: "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Components.App(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Search Tests
// =============================================================================

func TestComponents_Search(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/components/search", r.URL.Path)
		assert.Equal(t, "TRK", r.URL.Query().Get("qualifiers"))
		assert.Equal(t, "sonar", r.URL.Query().Get("q"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [
				{"key": "my_project", "name": "My Project", "qualifier": "TRK"},
				{"key": "another_project", "name": "Another Project", "qualifier": "TRK"}
			],
			"paging": {
				"pageIndex": 1,
				"pageSize": 100,
				"total": 2
			}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Components.Search(&ComponentsSearchOption{
		Qualifiers: []string{"TRK"},
		Query:      "sonar",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Components, 2)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestComponents_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *ComponentsSearchOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing qualifiers",
			opt:  &ComponentsSearchOption{Query: "test"},
		},
		{
			name: "empty qualifiers slice",
			opt:  &ComponentsSearchOption{Qualifiers: []string{}},
		},
		{
			name: "invalid qualifier",
			opt:  &ComponentsSearchOption{Qualifiers: []string{"INVALID"}},
		},
		{
			name: "query too short",
			opt: &ComponentsSearchOption{
				Qualifiers: []string{"TRK"},
				Query:      "a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Components.Search(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// SearchProjects Tests
// =============================================================================

func TestComponents_SearchProjects(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/components/search_projects", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [
				{
					"key": "my_project",
					"name": "My Project",
					"qualifier": "TRK",
					"visibility": "public",
					"isFavorite": true
				}
			],
			"facets": [
				{
					"property": "alert_status",
					"values": [
						{"val": "OK", "count": 10},
						{"val": "ERROR", "count": 2}
					]
				}
			],
			"paging": {
				"pageIndex": 1,
				"pageSize": 100,
				"total": 1
			}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Components.SearchProjects(&ComponentsSearchProjectsOption{
		Facets: []string{"alert_status"},
		Filter: "coverage >= 80",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Components, 1)
	assert.Len(t, result.Facets, 1)
	assert.Equal(t, "alert_status", result.Facets[0].Property)
}

func TestComponents_SearchProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *ComponentsSearchProjectsOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "invalid field",
			opt:  &ComponentsSearchProjectsOption{Fields: []string{"invalid_field"}},
		},
		{
			name: "invalid facet",
			opt:  &ComponentsSearchProjectsOption{Facets: []string{"invalid_facet"}},
		},
		{
			name: "invalid sort field",
			opt:  &ComponentsSearchProjectsOption{Sort: "invalid_sort"},
		},
		{
			name: "filter too short",
			opt:  &ComponentsSearchProjectsOption{Filter: "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Components.SearchProjects(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Show Tests
// =============================================================================

func TestComponents_Show(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/components/show", r.URL.Path)
		assert.Equal(t, "my_project:src/file.go", r.URL.Query().Get("component"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"component": {
				"key": "my_project:src/file.go",
				"name": "file.go",
				"path": "src/file.go",
				"qualifier": "FIL",
				"language": "go"
			},
			"ancestors": [
				{
					"key": "my_project:src",
					"name": "src",
					"qualifier": "DIR"
				},
				{
					"key": "my_project",
					"name": "My Project",
					"qualifier": "TRK"
				}
			]
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Components.Show(&ComponentsShowOption{
		Component: "my_project:src/file.go",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_project:src/file.go", result.Component.Key)
	assert.Len(t, result.Ancestors, 2)
	assert.Equal(t, "TRK", result.Ancestors[1].Qualifier)
}

func TestComponents_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *ComponentsShowOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing component",
			opt:  &ComponentsShowOption{Branch: "main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Components.Show(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Suggestions Tests
// =============================================================================

func TestComponents_Suggestions(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/components/suggestions", r.URL.Path)
		assert.Equal(t, "sonar", r.URL.Query().Get("s"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"q": "TRK",
					"more": 5,
					"items": [
						{
							"key": "sonarqube",
							"name": "SonarQube",
							"match": "<mark>Sonar</mark>Qube",
							"isFavorite": true
						}
					]
				}
			]
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Components.Suggestions(&ComponentsSuggestionsOption{
		Search: "sonar",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Results, 1)
	assert.Equal(t, "TRK", result.Results[0].Q)
	assert.Len(t, result.Results[0].Items, 1)
}

func TestComponents_Suggestions_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *ComponentsSuggestionsOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "invalid more value",
			opt:  &ComponentsSuggestionsOption{More: "INVALID"},
		},
		{
			name: "search too short",
			opt:  &ComponentsSuggestionsOption{Search: "a"},
		},
		{
			name: "too many recently browsed items",
			opt: &ComponentsSuggestionsOption{
				RecentlyBrowsed: make([]string, MaxRecentlyBrowsedItems+1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Components.Suggestions(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Tree Tests
// =============================================================================

func TestComponents_Tree(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/components/tree", r.URL.Path)
		assert.Equal(t, "my_project", r.URL.Query().Get("component"))
		assert.Equal(t, "children", r.URL.Query().Get("strategy"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"baseComponent": {
				"key": "my_project",
				"qualifier": "TRK",
				"visibility": "public"
			},
			"components": [
				{
					"key": "my_project:src",
					"name": "src",
					"path": "src",
					"qualifier": "DIR"
				},
				{
					"key": "my_project:main.go",
					"name": "main.go",
					"path": "main.go",
					"qualifier": "FIL",
					"language": "go"
				}
			],
			"paging": {
				"pageIndex": 1,
				"pageSize": 100,
				"total": 2
			}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Components.Tree(&ComponentsTreeOption{
		Component: "my_project",
		Strategy:  "children",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_project", result.BaseComponent.Key)
	assert.Len(t, result.Components, 2)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestComponents_Tree_WithQualifiers(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "FIL,UTS", r.URL.Query().Get("qualifiers"))
		assert.Equal(t, "name,path", r.URL.Query().Get("s"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"baseComponent": {"key": "my_project"},
			"components": [],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}
		}`))
	})

	client := newTestClient(t, server.URL)

	_, resp, err := client.Components.Tree(&ComponentsTreeOption{
		Component:  "my_project",
		Qualifiers: []string{"FIL", "UTS"},
		Sort:       []string{"name", "path"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestComponents_Tree_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *ComponentsTreeOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing component",
			opt:  &ComponentsTreeOption{Strategy: "all"},
		},
		{
			name: "query too short",
			opt: &ComponentsTreeOption{
				Component: "my_project",
				Query:     "ab",
			},
		},
		{
			name: "invalid qualifier",
			opt: &ComponentsTreeOption{
				Component:  "my_project",
				Qualifiers: []string{"INVALID"},
			},
		},
		{
			name: "invalid sort field",
			opt: &ComponentsTreeOption{
				Component: "my_project",
				Sort:      []string{"invalid"},
			},
		},
		{
			name: "invalid strategy",
			opt: &ComponentsTreeOption{
				Component: "my_project",
				Strategy:  "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Components.Tree(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Pagination Tests
// =============================================================================

func TestComponents_Search_WithPagination(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("p"))
		assert.Equal(t, "50", r.URL.Query().Get("ps"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [],
			"paging": {"pageIndex": 2, "pageSize": 50, "total": 100}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, _, err := client.Components.Search(&ComponentsSearchOption{
		Qualifiers: []string{"TRK"},
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(2), result.Paging.PageIndex)
}

// =============================================================================
// Allowed Values Tests
// =============================================================================

func TestComponents_Tree_AllowedQualifiers(t *testing.T) {
	client := newLocalhostClient(t)

	validQualifiers := []string{"UTS", "FIL", "DIR", "TRK"}
	for _, q := range validQualifiers {
		opt := &ComponentsTreeOption{
			Component:  "my_project",
			Qualifiers: []string{q},
		}
		err := client.Components.ValidateTreeOpt(opt)
		assert.NoError(t, err, "expected qualifier %s to be valid", q)
	}
}

func TestComponents_Tree_AllowedStrategies(t *testing.T) {
	client := newLocalhostClient(t)

	validStrategies := []string{"all", "children", "leaves"}
	for _, s := range validStrategies {
		opt := &ComponentsTreeOption{
			Component: "my_project",
			Strategy:  s,
		}
		err := client.Components.ValidateTreeOpt(opt)
		assert.NoError(t, err, "expected strategy %s to be valid", s)
	}
}

func TestComponents_Tree_AllowedSortFields(t *testing.T) {
	client := newLocalhostClient(t)

	validSortFields := []string{"name", "path", "qualifier"}
	for _, s := range validSortFields {
		opt := &ComponentsTreeOption{
			Component: "my_project",
			Sort:      []string{s},
		}
		err := client.Components.ValidateTreeOpt(opt)
		assert.NoError(t, err, "expected sort field %s to be valid", s)
	}
}

func TestComponents_Suggestions_AllowedMore(t *testing.T) {
	client := newLocalhostClient(t)

	validMoreValues := []string{"VW", "SVW", "APP", "TRK"}
	for _, m := range validMoreValues {
		opt := &ComponentsSuggestionsOption{
			More: m,
		}
		err := client.Components.ValidateSuggestionsOpt(opt)
		assert.NoError(t, err, "expected more value %s to be valid", m)
	}
}

func TestComponents_SearchProjects_AllowedFields(t *testing.T) {
	client := newLocalhostClient(t)

	validFields := []string{"analysisDate", "leakPeriodDate", "_all"}
	for _, f := range validFields {
		opt := &ComponentsSearchProjectsOption{
			Fields: []string{f},
		}
		err := client.Components.ValidateSearchProjectsOpt(opt)
		assert.NoError(t, err, "expected field %s to be valid", f)
	}
}
