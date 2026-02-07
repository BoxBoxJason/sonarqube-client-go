package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavorites_Add(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/favorites/add", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &FavoritesAddOption{
		Component: "my-project",
	}

	resp, err := client.Favorites.Add(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestFavorites_Add_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Favorites.Add(nil)
	assert.Error(t, err)

	// Missing Component should fail validation.
	_, err = client.Favorites.Add(&FavoritesAddOption{})
	assert.Error(t, err)
}

func TestFavorites_Remove(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/favorites/remove", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &FavoritesRemoveOption{
		Component: "my-project",
	}

	resp, err := client.Favorites.Remove(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestFavorites_Remove_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Favorites.Remove(nil)
	assert.Error(t, err)

	// Missing Component should fail validation.
	_, err = client.Favorites.Remove(&FavoritesRemoveOption{})
	assert.Error(t, err)
}

func TestFavorites_Search(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/favorites/search", http.StatusOK, &FavoritesSearch{
		Favorites: []Favorite{
			{Key: "project-1", Name: "Project One", Qualifier: "TRK"},
			{Key: "project-2", Name: "Project Two", Qualifier: "TRK"},
		},
		Paging: Paging{
			PageIndex: 1,
			PageSize:  100,
			Total:     2,
		},
	}))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Favorites.Search(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Favorites, 2)
	assert.Equal(t, "project-1", result.Favorites[0].Key)
	assert.Equal(t, "TRK", result.Favorites[0].Qualifier)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestFavorites_Search_WithPagination(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/favorites/search", http.StatusOK, &FavoritesSearch{
		Favorites: []Favorite{},
		Paging: Paging{
			PageIndex: 2,
			PageSize:  50,
			Total:     0,
		},
	}))
	client := newTestClient(t, server.URL)

	opt := &FavoritesSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	}

	_, resp, err := client.Favorites.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFavorites_ValidateAddOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.Favorites.ValidateAddOpt(&FavoritesAddOption{
		Component: "my-project",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.Favorites.ValidateAddOpt(nil)
	assert.Error(t, err)

	// Missing Component should fail.
	err = client.Favorites.ValidateAddOpt(&FavoritesAddOption{})
	assert.Error(t, err)
}

func TestFavorites_ValidateRemoveOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.Favorites.ValidateRemoveOpt(&FavoritesRemoveOption{
		Component: "my-project",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.Favorites.ValidateRemoveOpt(nil)
	assert.Error(t, err)

	// Missing Component should fail.
	err = client.Favorites.ValidateRemoveOpt(&FavoritesRemoveOption{})
	assert.Error(t, err)
}

func TestFavorites_ValidateSearchOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should be valid.
	err := client.Favorites.ValidateSearchOpt(nil)
	assert.NoError(t, err)

	// Empty option should be valid.
	err = client.Favorites.ValidateSearchOpt(&FavoritesSearchOption{})
	assert.NoError(t, err)

	// Valid pagination should be valid.
	err = client.Favorites.ValidateSearchOpt(&FavoritesSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     1,
			PageSize: 100,
		},
	})
	assert.NoError(t, err)
}
