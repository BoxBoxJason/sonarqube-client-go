package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFavorites_Add(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/favorites/add" {
			t.Errorf("expected path /api/favorites/add, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &FavoritesAddOption{
		Component: "my-project",
	}

	resp, err := client.Favorites.Add(opt)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestFavorites_Add_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Favorites.Add(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Component should fail validation.
	_, err = client.Favorites.Add(&FavoritesAddOption{})
	if err == nil {
		t.Error("expected error for missing Component")
	}
}

func TestFavorites_Remove(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/favorites/remove" {
			t.Errorf("expected path /api/favorites/remove, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &FavoritesRemoveOption{
		Component: "my-project",
	}

	resp, err := client.Favorites.Remove(opt)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestFavorites_Remove_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Favorites.Remove(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Component should fail validation.
	_, err = client.Favorites.Remove(&FavoritesRemoveOption{})
	if err == nil {
		t.Error("expected error for missing Component")
	}
}

func TestFavorites_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/favorites/search" {
			t.Errorf("expected path /api/favorites/search, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"favorites": [
				{"key": "project-1", "name": "Project One", "qualifier": "TRK"},
				{"key": "project-2", "name": "Project Two", "qualifier": "TRK"}
			],
			"paging": {
				"pageIndex": 1,
				"pageSize": 100,
				"total": 2
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Favorites.Search(nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Favorites) != 2 {
		t.Errorf("expected 2 favorites, got %d", len(result.Favorites))
	}

	if result.Favorites[0].Key != "project-1" {
		t.Errorf("expected first favorite key 'project-1', got %s", result.Favorites[0].Key)
	}

	if result.Favorites[0].Qualifier != "TRK" {
		t.Errorf("expected first favorite qualifier 'TRK', got %s", result.Favorites[0].Qualifier)
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected paging total 2, got %d", result.Paging.Total)
	}
}

func TestFavorites_Search_WithPagination(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"favorites": [], "paging": {"pageIndex": 2, "pageSize": 50, "total": 0}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &FavoritesSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	}

	_, resp, err := client.Favorites.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestFavorites_ValidateAddOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.Favorites.ValidateAddOpt(&FavoritesAddOption{
		Component: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.Favorites.ValidateAddOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Component should fail.
	err = client.Favorites.ValidateAddOpt(&FavoritesAddOption{})
	if err == nil {
		t.Error("expected error for missing Component")
	}
}

func TestFavorites_ValidateRemoveOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.Favorites.ValidateRemoveOpt(&FavoritesRemoveOption{
		Component: "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.Favorites.ValidateRemoveOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Component should fail.
	err = client.Favorites.ValidateRemoveOpt(&FavoritesRemoveOption{})
	if err == nil {
		t.Error("expected error for missing Component")
	}
}

func TestFavorites_ValidateSearchOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should be valid.
	err := client.Favorites.ValidateSearchOpt(nil)
	if err != nil {
		t.Errorf("expected nil error for nil option, got %v", err)
	}

	// Empty option should be valid.
	err = client.Favorites.ValidateSearchOpt(&FavoritesSearchOption{})
	if err != nil {
		t.Errorf("expected nil error for empty option, got %v", err)
	}

	// Valid pagination should be valid.
	err = client.Favorites.ValidateSearchOpt(&FavoritesSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     1,
			PageSize: 100,
		},
	})
	if err != nil {
		t.Errorf("expected nil error for valid pagination, got %v", err)
	}
}
