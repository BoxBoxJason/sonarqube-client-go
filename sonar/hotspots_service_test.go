package sonargo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// AddComment Tests
// =============================================================================

func TestHotspotsAddComment_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/add_comment", r.URL.Path)
		assert.Equal(t, "hotspot123", r.FormValue("hotspot"))
		assert.Equal(t, "This is a comment", r.FormValue("comment"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	resp, err := client.Hotspots.AddComment(&sonargo.HotspotsAddCommentOption{
		Hotspot: "hotspot123",
		Comment: "This is a comment",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspotsAddComment_MissingHotspot(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.AddComment(&sonargo.HotspotsAddCommentOption{
		Comment: "This is a comment",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Hotspot")
}

func TestHotspotsAddComment_MissingComment(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.AddComment(&sonargo.HotspotsAddCommentOption{
		Hotspot: "hotspot123",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Comment")
}

func TestHotspotsAddComment_CommentTooLong(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	longComment := make([]byte, 1001)
	for i := range longComment {
		longComment[i] = 'a'
	}

	_, err = client.Hotspots.AddComment(&sonargo.HotspotsAddCommentOption{
		Hotspot: "hotspot123",
		Comment: string(longComment),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Comment")
}

func TestHotspotsAddComment_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.AddComment(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// Assign Tests
// =============================================================================

func TestHotspotsAssign_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/assign", r.URL.Path)
		assert.Equal(t, "hotspot123", r.FormValue("hotspot"))
		assert.Equal(t, "john.doe", r.FormValue("assignee"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	resp, err := client.Hotspots.Assign(&sonargo.HotspotsAssignOption{
		Hotspot:  "hotspot123",
		Assignee: "john.doe",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspotsAssign_MissingHotspot(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.Assign(&sonargo.HotspotsAssignOption{
		Assignee: "john.doe",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Hotspot")
}

func TestHotspotsAssign_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.Assign(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// ChangeStatus Tests
// =============================================================================

func TestHotspotsChangeStatus_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/change_status", r.URL.Path)
		assert.Equal(t, "hotspot123", r.FormValue("hotspot"))
		assert.Equal(t, "REVIEWED", r.FormValue("status"))
		assert.Equal(t, "SAFE", r.FormValue("resolution"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	resp, err := client.Hotspots.ChangeStatus(&sonargo.HotspotsChangeStatusOption{
		Hotspot:    "hotspot123",
		Status:     "REVIEWED",
		Resolution: "SAFE",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspotsChangeStatus_InvalidStatus(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.ChangeStatus(&sonargo.HotspotsChangeStatusOption{
		Hotspot: "hotspot123",
		Status:  "INVALID",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Status")
}

func TestHotspotsChangeStatus_InvalidResolution(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.ChangeStatus(&sonargo.HotspotsChangeStatusOption{
		Hotspot:    "hotspot123",
		Status:     "REVIEWED",
		Resolution: "INVALID",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Resolution")
}

func TestHotspotsChangeStatus_MissingHotspot(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.ChangeStatus(&sonargo.HotspotsChangeStatusOption{
		Status: "REVIEWED",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Hotspot")
}

func TestHotspotsChangeStatus_MissingStatus(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.ChangeStatus(&sonargo.HotspotsChangeStatusOption{
		Hotspot: "hotspot123",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Status")
}

func TestHotspotsChangeStatus_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.ChangeStatus(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// DeleteComment Tests
// =============================================================================

func TestHotspotsDeleteComment_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/delete_comment", r.URL.Path)
		assert.Equal(t, "comment123", r.FormValue("comment"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	resp, err := client.Hotspots.DeleteComment(&sonargo.HotspotsDeleteCommentOption{
		Comment: "comment123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspotsDeleteComment_MissingComment(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.DeleteComment(&sonargo.HotspotsDeleteCommentOption{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Comment")
}

func TestHotspotsDeleteComment_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, err = client.Hotspots.DeleteComment(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// EditComment Tests
// =============================================================================

func TestHotspotsEditComment_Success(t *testing.T) {
	t.Parallel()

	expectedResponse := &sonargo.HotspotsEditComment{
		CreatedAt: "2024-01-01T12:00:00+0000",
		HTMLText:  "<p>Updated comment</p>",
		Key:       "comment123",
		Login:     "john.doe",
		Markdown:  "Updated comment",
		Updatable: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/edit_comment", r.URL.Path)
		assert.Equal(t, "comment123", r.FormValue("comment"))
		assert.Equal(t, "Updated comment", r.FormValue("text"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expectedResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	result, resp, err := client.Hotspots.EditComment(&sonargo.HotspotsEditCommentOption{
		Comment: "comment123",
		Text:    "Updated comment",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedResponse.Key, result.Key)
	assert.Equal(t, expectedResponse.Markdown, result.Markdown)
}

func TestHotspotsEditComment_MissingComment(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.EditComment(&sonargo.HotspotsEditCommentOption{
		Text: "Updated comment",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Comment")
}

func TestHotspotsEditComment_MissingText(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.EditComment(&sonargo.HotspotsEditCommentOption{
		Comment: "comment123",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Text")
}

func TestHotspotsEditComment_TextTooLong(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	longText := make([]byte, 1001)
	for i := range longText {
		longText[i] = 'a'
	}

	_, _, err = client.Hotspots.EditComment(&sonargo.HotspotsEditCommentOption{
		Comment: "comment123",
		Text:    string(longText),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Text")
}

func TestHotspotsEditComment_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.EditComment(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// List Tests
// =============================================================================

func TestHotspotsList_Success(t *testing.T) {
	t.Parallel()

	expectedResponse := &sonargo.HotspotsList{
		Components: []sonargo.HotspotComponent{
			{Key: "project:src/main.go", Name: "main.go", Qualifier: "FIL"},
		},
		Hotspots: []sonargo.HotspotSummary{
			{
				Key:                      "hotspot123",
				Component:                "project:src/main.go",
				Status:                   "TO_REVIEW",
				VulnerabilityProbability: "HIGH",
			},
		},
		Paging: sonargo.HotspotPaging{PageIndex: 1, PageSize: 100, Total: 1},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/list", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expectedResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	result, resp, err := client.Hotspots.List(&sonargo.HotspotsListOption{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Hotspots, 1)
	assert.Equal(t, "hotspot123", result.Hotspots[0].Key)
}

func TestHotspotsList_MissingProject(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.List(&sonargo.HotspotsListOption{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Project")
}

func TestHotspotsList_InvalidStatus(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.List(&sonargo.HotspotsListOption{
		Project: "my-project",
		Status:  "INVALID",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Status")
}

func TestHotspotsList_InvalidResolution(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.List(&sonargo.HotspotsListOption{
		Project:    "my-project",
		Resolution: "INVALID",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Resolution")
}

func TestHotspotsList_PageSizeTooLarge(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.List(&sonargo.HotspotsListOption{
		Project: "my-project",
		PaginationArgs: sonargo.PaginationArgs{
			PageSize: 501,
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "PageSize")
}

func TestHotspotsList_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.List(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// Pull Tests
// =============================================================================

func TestHotspotsPull_Success(t *testing.T) {
	t.Parallel()

	// The Pull endpoint returns raw bytes - use an empty JSON array for compatibility
	// with the client's JSON decoder when used with []byte type
	expectedData := `[]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/pull", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("projectKey"))
		assert.Equal(t, "main", r.URL.Query().Get("branchName"))
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(expectedData))
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	result, resp, err := client.Hotspots.Pull(&sonargo.HotspotsPullOption{
		ProjectKey: "my-project",
		BranchName: "main",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestHotspotsPull_MissingProjectKey(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Pull(&sonargo.HotspotsPullOption{
		BranchName: "main",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ProjectKey")
}

func TestHotspotsPull_MissingBranchName(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Pull(&sonargo.HotspotsPullOption{
		ProjectKey: "my-project",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "BranchName")
}

func TestHotspotsPull_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Pull(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// Search Tests
// =============================================================================

func TestHotspotsSearch_SuccessWithProject(t *testing.T) {
	t.Parallel()

	expectedResponse := &sonargo.HotspotsSearch{
		Components: []sonargo.HotspotComponent{
			{Key: "project:src/main.go", Name: "main.go", Qualifier: "FIL"},
		},
		Hotspots: []sonargo.HotspotSummary{
			{
				Key:                      "hotspot123",
				Component:                "project:src/main.go",
				Status:                   "TO_REVIEW",
				VulnerabilityProbability: "HIGH",
			},
		},
		Paging: sonargo.HotspotPaging{PageIndex: 1, PageSize: 100, Total: 1},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/search", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expectedResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	result, resp, err := client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Hotspots, 1)
	assert.Equal(t, "hotspot123", result.Hotspots[0].Key)
}

func TestHotspotsSearch_SuccessWithHotspots(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/search", r.URL.Path)
		assert.Equal(t, "hotspot1,hotspot2", r.URL.Query().Get("hotspots"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(&sonargo.HotspotsSearch{})
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	_, resp, err := client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Hotspots: []string{"hotspot1", "hotspot2"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHotspotsSearch_MissingProjectAndHotspots(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(&sonargo.HotspotsSearchOption{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "project or hotspots is required")
}

func TestHotspotsSearch_InvalidStatus(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project: "my-project",
		Status:  "INVALID",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Status")
}

func TestHotspotsSearch_InvalidResolution(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project:    "my-project",
		Resolution: "INVALID",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Resolution")
}

func TestHotspotsSearch_InvalidOwaspAsvsLevel(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project:        "my-project",
		OwaspAsvsLevel: "5",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OwaspAsvsLevel")
}

func TestHotspotsSearch_InvalidOwaspTop10(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project:    "my-project",
		OwaspTop10: []string{"a1", "invalid"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OwaspTop10")
}

func TestHotspotsSearch_InvalidSansTop25(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project:   "my-project",
		SansTop25: []string{"insecure-interaction", "invalid"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SansTop25")
}

func TestHotspotsSearch_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Search(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

func TestHotspotsSearch_WithFilters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/search", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		assert.Equal(t, "REVIEWED", r.URL.Query().Get("status"))
		assert.Equal(t, "SAFE", r.URL.Query().Get("resolution"))
		assert.Equal(t, "true", r.URL.Query().Get("inNewCodePeriod"))
		assert.Equal(t, "true", r.URL.Query().Get("onlyMine"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(&sonargo.HotspotsSearch{})
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	_, resp, err := client.Hotspots.Search(&sonargo.HotspotsSearchOption{
		Project:         "my-project",
		Status:          "REVIEWED",
		Resolution:      "SAFE",
		InNewCodePeriod: true,
		OnlyMine:        true,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// =============================================================================
// Show Tests
// =============================================================================

func TestHotspotsShow_Success(t *testing.T) {
	t.Parallel()

	expectedResponse := &sonargo.HotspotsShow{
		Key:             "hotspot123",
		Status:          "TO_REVIEW",
		CanChangeStatus: true,
		Message:         "Security issue found",
		Component: sonargo.HotspotComponent{
			Key:  "project:src/main.go",
			Name: "main.go",
		},
		Project: sonargo.HotspotProject{
			Key:  "my-project",
			Name: "My Project",
		},
		Rule: sonargo.HotspotRule{
			Key:  "java:S2092",
			Name: "Cookies should be secure",
		},
		Users: []sonargo.HotspotUser{
			{Login: "john.doe", Name: "John Doe", Active: true},
		},
		Changelog: []sonargo.HotspotChangelogEntry{
			{
				User:         "john.doe",
				CreationDate: "2024-01-01T12:00:00+0000",
				Diffs: []sonargo.HotspotDiff{
					{Key: "status", OldValue: "TO_REVIEW", NewValue: "REVIEWED"},
				},
			},
		},
		Comment: []sonargo.HotspotComment{
			{
				Key:      "comment123",
				Login:    "john.doe",
				Markdown: "This is safe",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/show", r.URL.Path)
		assert.Equal(t, "hotspot123", r.URL.Query().Get("hotspot"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expectedResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	client, err := sonargo.NewClient(server.URL, "", "")
	require.NoError(t, err)

	result, resp, err := client.Hotspots.Show(&sonargo.HotspotsShowOption{
		Hotspot: "hotspot123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hotspot123", result.Key)
	assert.Equal(t, "TO_REVIEW", result.Status)
	assert.True(t, result.CanChangeStatus)
	require.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
	require.Len(t, result.Changelog, 1)
	assert.Equal(t, "john.doe", result.Changelog[0].User)
	require.Len(t, result.Changelog[0].Diffs, 1)
	assert.Equal(t, "status", result.Changelog[0].Diffs[0].Key)
}

func TestHotspotsShow_MissingHotspot(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Show(&sonargo.HotspotsShowOption{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Hotspot")
}

func TestHotspotsShow_NilOption(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	_, _, err = client.Hotspots.Show(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "option struct")
}

// =============================================================================
// Validation Function Tests
// =============================================================================

func TestHotspotsValidateAddCommentOpt(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	testCases := []struct {
		name    string
		opt     *sonargo.HotspotsAddCommentOption
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid option",
			opt: &sonargo.HotspotsAddCommentOption{
				Comment: "Valid comment",
				Hotspot: "hotspot123",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "option struct",
		},
		{
			name: "missing comment",
			opt: &sonargo.HotspotsAddCommentOption{
				Hotspot: "hotspot123",
			},
			wantErr: true,
			errMsg:  "Comment",
		},
		{
			name: "missing hotspot",
			opt: &sonargo.HotspotsAddCommentOption{
				Comment: "Valid comment",
			},
			wantErr: true,
			errMsg:  "Hotspot",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := client.Hotspots.ValidateAddCommentOpt(tc.opt)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHotspotsValidateChangeStatusOpt(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	testCases := []struct {
		name    string
		opt     *sonargo.HotspotsChangeStatusOption
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid TO_REVIEW status",
			opt: &sonargo.HotspotsChangeStatusOption{
				Hotspot: "hotspot123",
				Status:  "TO_REVIEW",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with SAFE resolution",
			opt: &sonargo.HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "SAFE",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with FIXED resolution",
			opt: &sonargo.HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "FIXED",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with ACKNOWLEDGED resolution",
			opt: &sonargo.HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "ACKNOWLEDGED",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "option struct",
		},
		{
			name: "invalid status",
			opt: &sonargo.HotspotsChangeStatusOption{
				Hotspot: "hotspot123",
				Status:  "INVALID_STATUS",
			},
			wantErr: true,
			errMsg:  "Status",
		},
		{
			name: "invalid resolution",
			opt: &sonargo.HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "INVALID",
			},
			wantErr: true,
			errMsg:  "Resolution",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := client.Hotspots.ValidateChangeStatusOpt(tc.opt)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHotspotsValidateSearchOpt(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	testCases := []struct {
		name    string
		opt     *sonargo.HotspotsSearchOption
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid with project",
			opt: &sonargo.HotspotsSearchOption{
				Project: "my-project",
			},
			wantErr: false,
		},
		{
			name: "valid with hotspots",
			opt: &sonargo.HotspotsSearchOption{
				Hotspots: []string{"hotspot1", "hotspot2"},
			},
			wantErr: false,
		},
		{
			name: "valid with all OWASP filters",
			opt: &sonargo.HotspotsSearchOption{
				Project:        "my-project",
				OwaspTop10:     []string{"a1", "a2"},
				OwaspTop102021: []string{"a3", "a4"},
				OwaspAsvsLevel: "2",
			},
			wantErr: false,
		},
		{
			name: "valid with SANS filter",
			opt: &sonargo.HotspotsSearchOption{
				Project:   "my-project",
				SansTop25: []string{"insecure-interaction", "porous-defenses"},
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "option struct",
		},
		{
			name:    "missing project and hotspots",
			opt:     &sonargo.HotspotsSearchOption{},
			wantErr: true,
			errMsg:  "project or hotspots is required",
		},
		{
			name: "invalid OwaspAsvsLevel",
			opt: &sonargo.HotspotsSearchOption{
				Project:        "my-project",
				OwaspAsvsLevel: "4",
			},
			wantErr: true,
			errMsg:  "OwaspAsvsLevel",
		},
		{
			name: "invalid OwaspTop10 value",
			opt: &sonargo.HotspotsSearchOption{
				Project:    "my-project",
				OwaspTop10: []string{"a11"},
			},
			wantErr: true,
			errMsg:  "OwaspTop10",
		},
		{
			name: "invalid SansTop25 value",
			opt: &sonargo.HotspotsSearchOption{
				Project:   "my-project",
				SansTop25: []string{"invalid-category"},
			},
			wantErr: true,
			errMsg:  "SansTop25",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := client.Hotspots.ValidateSearchOpt(tc.opt)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHotspotsValidateListOpt(t *testing.T) {
	t.Parallel()

	client, err := sonargo.NewClient("http://localhost", "", "")
	require.NoError(t, err)

	testCases := []struct {
		name    string
		opt     *sonargo.HotspotsListOption
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid basic option",
			opt: &sonargo.HotspotsListOption{
				Project: "my-project",
			},
			wantErr: false,
		},
		{
			name: "valid with all optional params",
			opt: &sonargo.HotspotsListOption{
				Project:         "my-project",
				Branch:          "main",
				InNewCodePeriod: true,
				Status:          "TO_REVIEW",
				Resolution:      "SAFE",
				PaginationArgs:  sonargo.PaginationArgs{PageSize: 100},
			},
			wantErr: false,
		},
		{
			name: "valid max page size",
			opt: &sonargo.HotspotsListOption{
				Project:        "my-project",
				PaginationArgs: sonargo.PaginationArgs{PageSize: 500},
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
			errMsg:  "option struct",
		},
		{
			name:    "missing project",
			opt:     &sonargo.HotspotsListOption{},
			wantErr: true,
			errMsg:  "Project",
		},
		{
			name: "page size exceeds max",
			opt: &sonargo.HotspotsListOption{
				Project:        "my-project",
				PaginationArgs: sonargo.PaginationArgs{PageSize: 501},
			},
			wantErr: true,
			errMsg:  "PageSize",
		},
		{
			name: "invalid status",
			opt: &sonargo.HotspotsListOption{
				Project: "my-project",
				Status:  "CLOSED",
			},
			wantErr: true,
			errMsg:  "Status",
		},
		{
			name: "invalid resolution",
			opt: &sonargo.HotspotsListOption{
				Project:    "my-project",
				Resolution: "WONTFIX",
			},
			wantErr: true,
			errMsg:  "Resolution",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := client.Hotspots.ValidateListOpt(tc.opt)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
