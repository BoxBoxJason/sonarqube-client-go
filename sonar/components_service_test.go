package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// App Tests
// =============================================================================

func TestComponents_App(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/components/app" {
			t.Errorf("expected path /api/components/app, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my_project:src/file.go" {
			t.Errorf("expected component 'my_project:src/file.go', got %s", component)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "feature/my_branch" {
			t.Errorf("expected branch 'feature/my_branch', got %s", branch)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsApp{
			Key:               "my_project:src/file.go",
			Name:              "file.go",
			LongName:          "src/file.go",
			Q:                 "FIL",
			Project:           "my_project",
			ProjectName:       "My Project",
			Fav:               true,
			CanMarkAsFavorite: true,
			Measures: ComponentMeasures{
				Lines:       "100",
				Issues:      "5",
				SqaleRating: "A",
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Components.App(&ComponentsAppOption{
		Component: "my_project:src/file.go",
		Branch:    "feature/my_branch",
	})
	if err != nil {
		t.Fatalf("App failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Key != "my_project:src/file.go" {
		t.Errorf("expected key 'my_project:src/file.go', got %s", result.Key)
	}

	if !result.Fav {
		t.Error("expected Fav to be true")
	}
}

func TestComponents_App_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Search Tests
// =============================================================================

func TestComponents_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/components/search" {
			t.Errorf("expected path /api/components/search, got %s", r.URL.Path)
		}

		qualifiers := r.URL.Query().Get("qualifiers")
		if qualifiers != "TRK" {
			t.Errorf("expected qualifiers 'TRK', got %s", qualifiers)
		}

		query := r.URL.Query().Get("q")
		if query != "sonar" {
			t.Errorf("expected query 'sonar', got %s", query)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsSearch{
			Components: []ComponentSearchItem{
				{Key: "my_project", Name: "My Project", Qualifier: "TRK"},
				{Key: "another_project", Name: "Another Project", Qualifier: "TRK"},
			},
			Paging: Paging{
				PageIndex: 1,
				PageSize:  100,
				Total:     2,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Components.Search(&ComponentsSearchOption{
		Qualifiers: []string{"TRK"},
		Query:      "sonar",
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(result.Components))
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Paging.Total)
	}
}

func TestComponents_Search_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// SearchProjects Tests
// =============================================================================

func TestComponents_SearchProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/components/search_projects" {
			t.Errorf("expected path /api/components/search_projects, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsSearchProjects{
			Components: []ComponentProject{
				{
					Key:        "my_project",
					Name:       "My Project",
					Qualifier:  "TRK",
					Visibility: "public",
					IsFavorite: true,
				},
			},
			Facets: []ComponentFacet{
				{
					Property: "alert_status",
					Values: []ComponentFacetValue{
						{Val: "OK", Count: 10},
						{Val: "ERROR", Count: 2},
					},
				},
			},
			Paging: Paging{
				PageIndex: 1,
				PageSize:  100,
				Total:     1,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Components.SearchProjects(&ComponentsSearchProjectsOption{
		Facets: []string{"alert_status"},
		Filter: "coverage >= 80",
	})
	if err != nil {
		t.Fatalf("SearchProjects failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Components) != 1 {
		t.Errorf("expected 1 component, got %d", len(result.Components))
	}

	if len(result.Facets) != 1 {
		t.Errorf("expected 1 facet, got %d", len(result.Facets))
	}

	if result.Facets[0].Property != "alert_status" {
		t.Errorf("expected facet property 'alert_status', got %s", result.Facets[0].Property)
	}
}

func TestComponents_SearchProjects_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Show Tests
// =============================================================================

func TestComponents_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/components/show" {
			t.Errorf("expected path /api/components/show, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my_project:src/file.go" {
			t.Errorf("expected component 'my_project:src/file.go', got %s", component)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsShow{
			Component: ComponentDetails{
				Key:       "my_project:src/file.go",
				Name:      "file.go",
				Path:      "src/file.go",
				Qualifier: "FIL",
				Language:  "go",
			},
			Ancestors: []ComponentAncestor{
				{
					Key:       "my_project:src",
					Name:      "src",
					Qualifier: "DIR",
				},
				{
					Key:       "my_project",
					Name:      "My Project",
					Qualifier: "TRK",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Components.Show(&ComponentsShowOption{
		Component: "my_project:src/file.go",
	})
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Component.Key != "my_project:src/file.go" {
		t.Errorf("expected component key 'my_project:src/file.go', got %s", result.Component.Key)
	}

	if len(result.Ancestors) != 2 {
		t.Errorf("expected 2 ancestors, got %d", len(result.Ancestors))
	}

	if result.Ancestors[1].Qualifier != "TRK" {
		t.Errorf("expected root ancestor qualifier 'TRK', got %s", result.Ancestors[1].Qualifier)
	}
}

func TestComponents_Show_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Suggestions Tests
// =============================================================================

func TestComponents_Suggestions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/components/suggestions" {
			t.Errorf("expected path /api/components/suggestions, got %s", r.URL.Path)
		}

		search := r.URL.Query().Get("s")
		if search != "sonar" {
			t.Errorf("expected search 'sonar', got %s", search)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsSuggestions{
			Results: []ComponentSuggestionGroup{
				{
					Q:    "TRK",
					More: 5,
					Items: []ComponentSuggestionItem{
						{
							Key:        "sonarqube",
							Name:       "SonarQube",
							Match:      "<mark>Sonar</mark>Qube",
							IsFavorite: true,
						},
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Components.Suggestions(&ComponentsSuggestionsOption{
		Search: "sonar",
	})
	if err != nil {
		t.Fatalf("Suggestions failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if len(result.Results) != 1 {
		t.Errorf("expected 1 result group, got %d", len(result.Results))
	}

	if result.Results[0].Q != "TRK" {
		t.Errorf("expected qualifier 'TRK', got %s", result.Results[0].Q)
	}

	if len(result.Results[0].Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Results[0].Items))
	}
}

func TestComponents_Suggestions_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Tree Tests
// =============================================================================

func TestComponents_Tree(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/components/tree" {
			t.Errorf("expected path /api/components/tree, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my_project" {
			t.Errorf("expected component 'my_project', got %s", component)
		}

		strategy := r.URL.Query().Get("strategy")
		if strategy != "children" {
			t.Errorf("expected strategy 'children', got %s", strategy)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsTree{
			BaseComponent: ComponentTreeBase{
				Key:        "my_project",
				Qualifier:  "TRK",
				Visibility: "public",
			},
			Components: []ComponentTreeItem{
				{
					Key:       "my_project:src",
					Name:      "src",
					Path:      "src",
					Qualifier: "DIR",
				},
				{
					Key:       "my_project:main.go",
					Name:      "main.go",
					Path:      "main.go",
					Qualifier: "FIL",
					Language:  "go",
				},
			},
			Paging: Paging{
				PageIndex: 1,
				PageSize:  100,
				Total:     2,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Components.Tree(&ComponentsTreeOption{
		Component: "my_project",
		Strategy:  "children",
	})
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.BaseComponent.Key != "my_project" {
		t.Errorf("expected base component key 'my_project', got %s", result.BaseComponent.Key)
	}

	if len(result.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(result.Components))
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Paging.Total)
	}
}

func TestComponents_Tree_WithQualifiers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		qualifiers := r.URL.Query().Get("qualifiers")
		if qualifiers != "FIL,UTS" {
			t.Errorf("expected qualifiers 'FIL,UTS', got %s", qualifiers)
		}

		sortFields := r.URL.Query().Get("s")
		if sortFields != "name,path" {
			t.Errorf("expected sort 'name,path', got %s", sortFields)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &ComponentsTree{
			BaseComponent: ComponentTreeBase{Key: "my_project"},
			Components:    []ComponentTreeItem{},
			Paging:        Paging{PageIndex: 1, PageSize: 100, Total: 0},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, resp, err := client.Components.Tree(&ComponentsTreeOption{
		Component:  "my_project",
		Qualifiers: []string{"FIL", "UTS"},
		Sort:       []string{"name", "path"},
	})
	if err != nil {
		t.Fatalf("Tree failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestComponents_Tree_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// =============================================================================
// Pagination Tests
// =============================================================================

func TestComponents_Search_WithPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("p")
		if page != "2" {
			t.Errorf("expected page '2', got %s", page)
		}

		pageSize := r.URL.Query().Get("ps")
		if pageSize != "50" {
			t.Errorf("expected pageSize '50', got %s", pageSize)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&ComponentsSearch{
			Components: []ComponentSearchItem{},
			Paging:     Paging{PageIndex: 2, PageSize: 50, Total: 100},
		})
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, _, err := client.Components.Search(&ComponentsSearchOption{
		Qualifiers: []string{"TRK"},
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Paging.PageIndex != 2 {
		t.Errorf("expected page index 2, got %d", result.Paging.PageIndex)
	}
}

// =============================================================================
// Allowed Values Tests
// =============================================================================

func TestComponents_Tree_AllowedQualifiers(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	validQualifiers := []string{"UTS", "FIL", "DIR", "TRK"}
	for _, q := range validQualifiers {
		opt := &ComponentsTreeOption{
			Component:  "my_project",
			Qualifiers: []string{q},
		}
		if err := client.Components.ValidateTreeOpt(opt); err != nil {
			t.Errorf("expected qualifier %s to be valid, got error: %v", q, err)
		}
	}
}

func TestComponents_Tree_AllowedStrategies(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	validStrategies := []string{"all", "children", "leaves"}
	for _, s := range validStrategies {
		opt := &ComponentsTreeOption{
			Component: "my_project",
			Strategy:  s,
		}
		if err := client.Components.ValidateTreeOpt(opt); err != nil {
			t.Errorf("expected strategy %s to be valid, got error: %v", s, err)
		}
	}
}

func TestComponents_Tree_AllowedSortFields(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	validSortFields := []string{"name", "path", "qualifier"}
	for _, s := range validSortFields {
		opt := &ComponentsTreeOption{
			Component: "my_project",
			Sort:      []string{s},
		}
		if err := client.Components.ValidateTreeOpt(opt); err != nil {
			t.Errorf("expected sort field %s to be valid, got error: %v", s, err)
		}
	}
}

func TestComponents_Suggestions_AllowedMore(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	validMoreValues := []string{"VW", "SVW", "APP", "TRK"}
	for _, m := range validMoreValues {
		opt := &ComponentsSuggestionsOption{
			More: m,
		}
		if err := client.Components.ValidateSuggestionsOpt(opt); err != nil {
			t.Errorf("expected more value %s to be valid, got error: %v", m, err)
		}
	}
}

func TestComponents_SearchProjects_AllowedFields(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	validFields := []string{"analysisDate", "leakPeriodDate", "_all"}
	for _, f := range validFields {
		opt := &ComponentsSearchProjectsOption{
			Fields: []string{f},
		}
		if err := client.Components.ValidateSearchProjectsOpt(opt); err != nil {
			t.Errorf("expected field %s to be valid, got error: %v", f, err)
		}
	}
}
