package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// AddComment Tests
// =============================================================================

func TestHotspots_AddComment(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/add_comment", r.URL.Path)
		assert.Equal(t, "hotspot123", r.URL.Query().Get("hotspot"))
		assert.Equal(t, "This is a comment", r.URL.Query().Get("comment"))

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	resp, err := client.Hotspots.AddComment(&HotspotsAddCommentOption{
		Hotspot: "hotspot123",
		Comment: "This is a comment",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspots_AddComment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsAddCommentOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsAddCommentOption{Comment: "This is a comment"},
		},
		{
			name: "missing comment",
			opt:  &HotspotsAddCommentOption{Hotspot: "hotspot123"},
		},
		{
			name: "comment too long",
			opt: &HotspotsAddCommentOption{
				Hotspot: "hotspot123",
				Comment: string(make([]byte, MaxHotspotCommentLength+1)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.AddComment(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Assign Tests
// =============================================================================

func TestHotspots_Assign(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/assign", r.URL.Path)
		assert.Equal(t, "hotspot123", r.URL.Query().Get("hotspot"))
		assert.Equal(t, "john.doe", r.URL.Query().Get("assignee"))

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	resp, err := client.Hotspots.Assign(&HotspotsAssignOption{
		Hotspot:  "hotspot123",
		Assignee: "john.doe",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspots_Assign_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsAssignOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsAssignOption{Assignee: "john.doe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.Assign(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// ChangeStatus Tests
// =============================================================================

func TestHotspots_ChangeStatus(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/change_status", r.URL.Path)
		assert.Equal(t, "hotspot123", r.URL.Query().Get("hotspot"))
		assert.Equal(t, "REVIEWED", r.URL.Query().Get("status"))
		assert.Equal(t, "SAFE", r.URL.Query().Get("resolution"))

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	resp, err := client.Hotspots.ChangeStatus(&HotspotsChangeStatusOption{
		Hotspot:    "hotspot123",
		Status:     "REVIEWED",
		Resolution: "SAFE",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspots_ChangeStatus_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsChangeStatusOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsChangeStatusOption{Status: "REVIEWED"},
		},
		{
			name: "missing status",
			opt:  &HotspotsChangeStatusOption{Hotspot: "hotspot123"},
		},
		{
			name: "invalid status",
			opt:  &HotspotsChangeStatusOption{Hotspot: "hotspot123", Status: "INVALID"},
		},
		{
			name: "invalid resolution",
			opt:  &HotspotsChangeStatusOption{Hotspot: "hotspot123", Status: "REVIEWED", Resolution: "INVALID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.ChangeStatus(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// DeleteComment Tests
// =============================================================================

func TestHotspots_DeleteComment(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/hotspots/delete_comment", r.URL.Path)
		assert.Equal(t, "comment123", r.URL.Query().Get("comment"))

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	resp, err := client.Hotspots.DeleteComment(&HotspotsDeleteCommentOption{
		Comment: "comment123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHotspots_DeleteComment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsDeleteCommentOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing comment",
			opt:  &HotspotsDeleteCommentOption{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Hotspots.DeleteComment(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// EditComment Tests
// =============================================================================

func TestHotspots_EditComment(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/hotspots/edit_comment", http.StatusOK, &HotspotsEditComment{
		CreatedAt: "2024-01-01T12:00:00+0000",
		HTMLText:  "<p>Updated comment</p>",
		Key:       "comment123",
		Login:     "john.doe",
		Markdown:  "Updated comment",
		Updatable: true,
	}))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Hotspots.EditComment(&HotspotsEditCommentOption{
		Comment: "comment123",
		Text:    "Updated comment",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "comment123", result.Key)
	assert.Equal(t, "Updated comment", result.Markdown)
}

func TestHotspots_EditComment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsEditCommentOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing comment",
			opt:  &HotspotsEditCommentOption{Text: "Updated comment"},
		},
		{
			name: "missing text",
			opt:  &HotspotsEditCommentOption{Comment: "comment123"},
		},
		{
			name: "text too long",
			opt: &HotspotsEditCommentOption{
				Comment: "comment123",
				Text:    string(make([]byte, MaxHotspotCommentLength+1)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.EditComment(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// List Tests
// =============================================================================

func TestHotspots_List(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/list", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [
				{"key": "project:src/main.go", "name": "main.go", "qualifier": "FIL"}
			],
			"hotspots": [
				{
					"key": "hotspot123",
					"component": "project:src/main.go",
					"status": "TO_REVIEW",
					"vulnerabilityProbability": "HIGH"
				}
			],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Hotspots.List(&HotspotsListOption{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Hotspots, 1)
	assert.Equal(t, "hotspot123", result.Hotspots[0].Key)
}

func TestHotspots_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsListOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project",
			opt:  &HotspotsListOption{},
		},
		{
			name: "invalid status",
			opt:  &HotspotsListOption{Project: "my-project", Status: "INVALID"},
		},
		{
			name: "invalid resolution",
			opt:  &HotspotsListOption{Project: "my-project", Resolution: "INVALID"},
		},
		{
			name: "page size too large",
			opt:  &HotspotsListOption{Project: "my-project", PaginationArgs: PaginationArgs{PageSize: MaxHotspotListPageSize + 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.List(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Pull Tests
// =============================================================================

func TestHotspots_Pull(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/pull", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("projectKey"))
		assert.Equal(t, "main", r.URL.Query().Get("branchName"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Hotspots.Pull(&HotspotsPullOption{
		ProjectKey: "my-project",
		BranchName: "main",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestHotspots_Pull_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsPullOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project key",
			opt:  &HotspotsPullOption{BranchName: "main"},
		},
		{
			name: "missing branch name",
			opt:  &HotspotsPullOption{ProjectKey: "my-project"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.Pull(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Search Tests
// =============================================================================

func TestHotspots_Search_WithProject(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/search", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"components": [
				{"key": "project:src/main.go", "name": "main.go", "qualifier": "FIL"}
			],
			"hotspots": [
				{
					"key": "hotspot123",
					"component": "project:src/main.go",
					"status": "TO_REVIEW",
					"vulnerabilityProbability": "HIGH"
				}
			],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Hotspots.Search(&HotspotsSearchOption{
		Project: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Hotspots, 1)
	assert.Equal(t, "hotspot123", result.Hotspots[0].Key)
}

func TestHotspots_Search_WithHotspots(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "hotspot1,hotspot2", r.URL.Query().Get("hotspots"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"components": [], "hotspots": [], "paging": {"pageIndex": 1, "pageSize": 100, "total": 0}}`))
	})

	client := newTestClient(t, server.URL)

	_, resp, err := client.Hotspots.Search(&HotspotsSearchOption{
		Hotspots: []string{"hotspot1", "hotspot2"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHotspots_Search_WithFilters(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		assert.Equal(t, "REVIEWED", r.URL.Query().Get("status"))
		assert.Equal(t, "SAFE", r.URL.Query().Get("resolution"))
		assert.Equal(t, "true", r.URL.Query().Get("inNewCodePeriod"))
		assert.Equal(t, "true", r.URL.Query().Get("onlyMine"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"components": [], "hotspots": [], "paging": {"pageIndex": 1, "pageSize": 100, "total": 0}}`))
	})

	client := newTestClient(t, server.URL)

	_, resp, err := client.Hotspots.Search(&HotspotsSearchOption{
		Project:         "my-project",
		Status:          "REVIEWED",
		Resolution:      "SAFE",
		InNewCodePeriod: true,
		OnlyMine:        true,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHotspots_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsSearchOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project and hotspots",
			opt:  &HotspotsSearchOption{},
		},
		{
			name: "invalid status",
			opt:  &HotspotsSearchOption{Project: "my-project", Status: "INVALID"},
		},
		{
			name: "invalid resolution",
			opt:  &HotspotsSearchOption{Project: "my-project", Resolution: "INVALID"},
		},
		{
			name: "invalid owasp asvs level",
			opt:  &HotspotsSearchOption{Project: "my-project", OwaspAsvsLevel: "5"},
		},
		{
			name: "invalid owasp top 10",
			opt:  &HotspotsSearchOption{Project: "my-project", OwaspTop10: []string{"a1", "invalid"}},
		},
		{
			name: "invalid sans top 25",
			opt:  &HotspotsSearchOption{Project: "my-project", SansTop25: []string{"insecure-interaction", "invalid"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.Search(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Show Tests
// =============================================================================

func TestHotspots_Show(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/hotspots/show", r.URL.Path)
		assert.Equal(t, "hotspot123", r.URL.Query().Get("hotspot"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"key": "hotspot123",
			"status": "TO_REVIEW",
			"canChangeStatus": true,
			"message": "Security issue found",
			"component": {
				"key": "project:src/main.go",
				"name": "main.go"
			},
			"project": {
				"key": "my-project",
				"name": "My Project"
			},
			"rule": {
				"key": "java:S2092",
				"name": "Cookies should be secure"
			},
			"users": [
				{"login": "john.doe", "name": "John Doe", "active": true}
			],
			"changelog": [
				{
					"user": "john.doe",
					"creationDate": "2024-01-01T12:00:00+0000",
					"diffs": [
						{"key": "status", "oldValue": "TO_REVIEW", "newValue": "REVIEWED"}
					]
				}
			],
			"comment": [
				{
					"key": "comment123",
					"login": "john.doe",
					"markdown": "This is safe"
				}
			]
		}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Hotspots.Show(&HotspotsShowOption{
		Hotspot: "hotspot123",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "hotspot123", result.Key)
	assert.Equal(t, "TO_REVIEW", result.Status)
	assert.True(t, result.CanChangeStatus)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
	assert.Len(t, result.Changelog, 1)
	assert.Equal(t, "john.doe", result.Changelog[0].User)
	assert.Len(t, result.Changelog[0].Diffs, 1)
	assert.Equal(t, "status", result.Changelog[0].Diffs[0].Key)
}

func TestHotspots_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *HotspotsShowOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing hotspot",
			opt:  &HotspotsShowOption{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Hotspots.Show(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Validation Function Tests
// =============================================================================

func TestHotspots_ValidateAddCommentOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *HotspotsAddCommentOption
		wantErr bool
	}{
		{
			name: "valid option",
			opt: &HotspotsAddCommentOption{
				Comment: "Valid comment",
				Hotspot: "hotspot123",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name: "missing comment",
			opt: &HotspotsAddCommentOption{
				Hotspot: "hotspot123",
			},
			wantErr: true,
		},
		{
			name: "missing hotspot",
			opt: &HotspotsAddCommentOption{
				Comment: "Valid comment",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateAddCommentOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHotspots_ValidateChangeStatusOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *HotspotsChangeStatusOption
		wantErr bool
	}{
		{
			name: "valid TO_REVIEW status",
			opt: &HotspotsChangeStatusOption{
				Hotspot: "hotspot123",
				Status:  "TO_REVIEW",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with SAFE resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "SAFE",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with FIXED resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "FIXED",
			},
			wantErr: false,
		},
		{
			name: "valid REVIEWED with ACKNOWLEDGED resolution",
			opt: &HotspotsChangeStatusOption{
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
		},
		{
			name: "invalid status",
			opt: &HotspotsChangeStatusOption{
				Hotspot: "hotspot123",
				Status:  "INVALID_STATUS",
			},
			wantErr: true,
		},
		{
			name: "invalid resolution",
			opt: &HotspotsChangeStatusOption{
				Hotspot:    "hotspot123",
				Status:     "REVIEWED",
				Resolution: "INVALID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateChangeStatusOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHotspots_ValidateSearchOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *HotspotsSearchOption
		wantErr bool
	}{
		{
			name: "valid with project",
			opt: &HotspotsSearchOption{
				Project: "my-project",
			},
			wantErr: false,
		},
		{
			name: "valid with hotspots",
			opt: &HotspotsSearchOption{
				Hotspots: []string{"hotspot1", "hotspot2"},
			},
			wantErr: false,
		},
		{
			name: "valid with all OWASP filters",
			opt: &HotspotsSearchOption{
				Project:        "my-project",
				OwaspTop10:     []string{"a1", "a2"},
				OwaspTop102021: []string{"a3", "a4"},
				OwaspAsvsLevel: "2",
			},
			wantErr: false,
		},
		{
			name: "valid with SANS filter",
			opt: &HotspotsSearchOption{
				Project:   "my-project",
				SansTop25: []string{"insecure-interaction", "porous-defenses"},
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing project and hotspots",
			opt:     &HotspotsSearchOption{},
			wantErr: true,
		},
		{
			name: "invalid OwaspAsvsLevel",
			opt: &HotspotsSearchOption{
				Project:        "my-project",
				OwaspAsvsLevel: "4",
			},
			wantErr: true,
		},
		{
			name: "invalid OwaspTop10 value",
			opt: &HotspotsSearchOption{
				Project:    "my-project",
				OwaspTop10: []string{"a11"},
			},
			wantErr: true,
		},
		{
			name: "invalid SansTop25 value",
			opt: &HotspotsSearchOption{
				Project:   "my-project",
				SansTop25: []string{"invalid-category"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateSearchOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHotspots_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *HotspotsListOption
		wantErr bool
	}{
		{
			name: "valid basic option",
			opt: &HotspotsListOption{
				Project: "my-project",
			},
			wantErr: false,
		},
		{
			name: "valid with all optional params",
			opt: &HotspotsListOption{
				Project:         "my-project",
				Branch:          "main",
				InNewCodePeriod: true,
				Status:          "TO_REVIEW",
				Resolution:      "SAFE",
				PaginationArgs:  PaginationArgs{PageSize: 100},
			},
			wantErr: false,
		},
		{
			name: "valid max page size",
			opt: &HotspotsListOption{
				Project:        "my-project",
				PaginationArgs: PaginationArgs{PageSize: MaxHotspotListPageSize},
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing project",
			opt:     &HotspotsListOption{},
			wantErr: true,
		},
		{
			name: "page size exceeds max",
			opt: &HotspotsListOption{
				Project:        "my-project",
				PaginationArgs: PaginationArgs{PageSize: MaxHotspotListPageSize + 1},
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			opt: &HotspotsListOption{
				Project: "my-project",
				Status:  "CLOSED",
			},
			wantErr: true,
		},
		{
			name: "invalid resolution",
			opt: &HotspotsListOption{
				Project:    "my-project",
				Resolution: "WONTFIX",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Hotspots.ValidateListOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
