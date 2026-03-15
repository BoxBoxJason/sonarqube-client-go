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

func TestIssues_AddComment(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/add_comment", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"issue": {"key": "AU-Tpxb--iU5OvuD2FLy", "message": "Test issue"},
			"components": [],
			"rules": [],
			"users": []
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesAddCommentOptions{
		Issue: "AU-Tpxb--iU5OvuD2FLy",
		Text:  "This is a comment",
	}
	result, resp, err := client.Issues.AddComment(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "AU-Tpxb--iU5OvuD2FLy", result.Issue.Key)
}

func TestIssues_AddComment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesAddCommentOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing Issue",
			opt:  &IssuesAddCommentOptions{Text: "test"},
		},
		{
			name: "missing Text",
			opt:  &IssuesAddCommentOptions{Issue: "key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Issues.AddComment(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Assign Tests
// =============================================================================

func TestIssues_Assign(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/assign", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"issue": {"key": "test-key", "assignee": "admin"},
			"components": [],
			"rules": [],
			"users": []
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesAssignOptions{
		Issue:    "test-key",
		Assignee: "admin",
	}
	result, resp, err := client.Issues.Assign(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "admin", result.Issue.Assignee)
}

// =============================================================================
// Authors Tests
// =============================================================================

func TestIssues_Authors(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/authors", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"authors": ["john@example.com", "jane@example.com"]}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesAuthorsOptions{
		Project:  "my-project",
		PageSize: 50,
	}
	result, resp, err := client.Issues.Authors(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Authors, 2)
}

func TestIssues_Authors_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Issues.Authors(&IssuesAuthorsOptions{PageSize: 150})
	assert.Error(t, err)
}

// =============================================================================
// BulkChange Tests
// =============================================================================

func TestIssues_BulkChange(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/bulk_change", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total": 10, "success": 8, "ignored": 1, "failures": 1}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesBulkChangeOptions{
		Issues:      []string{"issue1", "issue2"},
		SetSeverity: RuleSeverityMajor,
		AddTags:     []string{"tag1", "tag2"},
	}
	result, resp, err := client.Issues.BulkChange(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(8), result.Success)
}

func TestIssues_BulkChange_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesBulkChangeOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "empty issues",
			opt:  &IssuesBulkChangeOptions{},
		},
		{
			name: "invalid severity",
			opt: &IssuesBulkChangeOptions{
				Issues:      []string{"issue1"},
				SetSeverity: "INVALID",
			},
		},
		{
			name: "invalid type",
			opt: &IssuesBulkChangeOptions{
				Issues:  []string{"issue1"},
				SetType: "INVALID",
			},
		},
		{
			name: "invalid transition",
			opt: &IssuesBulkChangeOptions{
				Issues:       []string{"issue1"},
				DoTransition: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Issues.BulkChange(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Changelog Tests
// =============================================================================

func TestIssues_Changelog(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/changelog", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"changelog": [
				{
					"user": "admin",
					"userName": "Admin User",
					"creationDate": "2023-01-01T00:00:00+0000",
					"diffs": [
						{"key": "severity", "oldValue": "MAJOR", "newValue": "CRITICAL"}
					]
				}
			]
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesChangelogOptions{Issue: "test-key"}
	result, resp, err := client.Issues.Changelog(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Changelog, 1)
}

// =============================================================================
// ComponentTags Tests
// =============================================================================

func TestIssues_ComponentTags(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/component_tags", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tags": [{"key": "security", "value": 10}]}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesComponentTagsOptions{ComponentUuid: "uuid-123"}
	result, resp, err := client.Issues.ComponentTags(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Tags, 1)
}

// =============================================================================
// DeleteComment Tests
// =============================================================================

func TestIssues_DeleteComment(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/delete_comment", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issue": {"key": "test-key"}, "components": [], "rules": [], "users": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesDeleteCommentOptions{Comment: "comment-key"}
	result, resp, err := client.Issues.DeleteComment(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "test-key", result.Issue.Key)
}

// =============================================================================
// DoTransition Tests
// =============================================================================

func TestIssues_DoTransition(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/do_transition", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issue": {"key": "test-key"}, "components": [], "rules": [], "users": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesDoTransitionOptions{
		Issue:      "test-key",
		Transition: IssueTransitionConfirm,
	}
	result, resp, err := client.Issues.DoTransition(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "test-key", result.Issue.Key)
}

func TestIssues_DoTransition_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Issues.DoTransition(&IssuesDoTransitionOptions{
		Issue:      "test-key",
		Transition: "invalid",
	})
	assert.Error(t, err)
}

// =============================================================================
// EditComment Tests
// =============================================================================

func TestIssues_EditComment(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/edit_comment", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issue": {"key": "test-key"}, "components": [], "rules": [], "users": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesEditCommentOptions{
		Comment: "comment-key",
		Text:    "Updated comment",
	}
	result, resp, err := client.Issues.EditComment(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "test-key", result.Issue.Key)
}

// =============================================================================
// List Tests
// =============================================================================

func TestIssues_List(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/list", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"issues": [{"key": "issue-1"}],
			"components": [],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesListOptions{
		Project: "my-project",
	}
	result, resp, err := client.Issues.List(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Issues, 1)
}

func TestIssues_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesListOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project and component",
			opt:  &IssuesListOptions{},
		},
		{
			name: "invalid type",
			opt: &IssuesListOptions{
				Project: "my-project",
				Types:   []string{"INVALID"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Issues.List(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Search Tests
// =============================================================================

func TestIssues_Search(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/search", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"issues": [{"key": "issue-1", "message": "Test issue"}],
			"components": [],
			"rules": [],
			"users": [],
			"facets": [],
			"paging": {"pageIndex": 1, "pageSize": 100, "total": 1}
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesSearchOptions{
		Projects:         []string{"my-project"},
		Severities:       []string{RuleSeverityBlocker, RuleSeverityCritical},
		Types:            []string{RuleTypeBug},
		ImpactSeverities: []string{RuleImpactSeverityHigh},
	}
	result, resp, err := client.Issues.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Issues, 1)
}

func TestIssues_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesSearchOptions
	}{
		{
			name: "invalid severity",
			opt: &IssuesSearchOptions{
				Severities: []string{"INVALID"},
			},
		},
		{
			name: "invalid type",
			opt: &IssuesSearchOptions{
				Types: []string{"INVALID"},
			},
		},
		{
			name: "invalid impact severity",
			opt: &IssuesSearchOptions{
				ImpactSeverities: []string{"INVALID"},
			},
		},
		{
			name: "invalid impact software quality",
			opt: &IssuesSearchOptions{
				ImpactSoftwareQualities: []string{"INVALID"},
			},
		},
		{
			name: "invalid clean code attribute category",
			opt: &IssuesSearchOptions{
				CleanCodeAttributeCategories: []string{"INVALID"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Issues.Search(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// SetSeverity Tests
// =============================================================================

func TestIssues_SetSeverity(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/set_severity", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issue": {"key": "test-key", "severity": "BLOCKER"}, "components": [], "rules": [], "users": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesSetSeverityOptions{
		Issue:    "test-key",
		Severity: RuleSeverityBlocker,
	}
	result, resp, err := client.Issues.SetSeverity(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, RuleSeverityBlocker, result.Issue.Severity)
}

func TestIssues_SetSeverity_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.Issues.SetSeverity(&IssuesSetSeverityOptions{
		Issue:    "test-key",
		Severity: "INVALID",
	})
	assert.Error(t, err)
}

// =============================================================================
// SetTags Tests
// =============================================================================

func TestIssues_SetTags(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/set_tags", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issue": {"key": "test-key", "tags": ["security", "cwe"]}, "components": [], "rules": [], "users": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesSetTagsOptions{
		Issue: "test-key",
		Tags:  []string{"security", "cwe"},
	}
	result, resp, err := client.Issues.SetTags(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Issue.Tags, 2)
}

// =============================================================================
// SetType Tests
// =============================================================================

func TestIssues_SetType(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/set_type", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issue": {"key": "test-key", "type": "BUG"}, "components": [], "rules": [], "users": []}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesSetTypeOptions{
		Issue: "test-key",
		Type:  RuleTypeBug,
	}
	result, resp, err := client.Issues.SetType(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, RuleTypeBug, result.Issue.Type)
}

func TestIssues_SetType_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesSetTypeOptions
	}{
		{
			name: "missing issue",
			opt:  &IssuesSetTypeOptions{Type: RuleTypeBug},
		},
		{
			name: "missing type",
			opt:  &IssuesSetTypeOptions{Issue: "test-key"},
		},
		{
			name: "invalid type",
			opt: &IssuesSetTypeOptions{
				Issue: "test-key",
				Type:  "INVALID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Issues.SetType(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Tags Tests
// =============================================================================

func TestIssues_Tags(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/tags", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tags": ["security", "cwe", "java"]}`))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesTagsOptions{
		Project:  "my-project",
		PageSize: 100,
	}
	result, resp, err := client.Issues.Tags(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Tags, 3)
}

// =============================================================================
// Pull Tests
// =============================================================================

func TestIssues_Pull(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/pull", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesPullOptions{
		ProjectKey: "my-project",
		BranchName: "main",
	}
	result, resp, err := client.Issues.Pull(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestIssues_Pull_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesPullOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project key",
			opt:  &IssuesPullOptions{},
		},
		{
			name: "invalid language",
			opt: &IssuesPullOptions{
				ProjectKey: "my-project",
				Languages:  []string{"invalid-language"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Issues.Pull(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// PullTaint Tests
// =============================================================================

func TestIssues_PullTaint(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/issues/pull_taint", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesPullTaintOptions{
		ProjectKey: "my-project",
		BranchName: "main",
	}
	result, resp, err := client.Issues.PullTaint(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

// =============================================================================
// Reindex Tests
// =============================================================================

func TestIssues_Reindex(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/reindex", r.URL.Path)

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesReindexOptions{Project: "my-project"}
	resp, err := client.Issues.Reindex(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// =============================================================================
// AnticipatedTransitions Tests
// =============================================================================

func TestIssues_AnticipatedTransitions(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/issues/anticipated_transitions", r.URL.Path)

		w.WriteHeader(http.StatusAccepted)
	})

	client := newTestClient(t, server.URL)

	opt := &IssuesAnticipatedTransitionsOptions{ProjectKey: "my-project"}
	resp, err := client.Issues.AnticipatedTransitions(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestIssues_AnticipatedTransitions_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *IssuesAnticipatedTransitionsOptions
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project key",
			opt:  &IssuesAnticipatedTransitionsOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Issues.AnticipatedTransitions(tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// Validation Functions Tests
// =============================================================================

func TestIssues_ValidateAssignOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesAssignOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Issue",
			opt:     &IssuesAssignOptions{},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesAssignOptions{Issue: "issue-key"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateAssignOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateChangelogOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesChangelogOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Issue",
			opt:     &IssuesChangelogOptions{},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesChangelogOptions{Issue: "issue-key"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateChangelogOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateComponentTagsOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesComponentTagsOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing ComponentUuid",
			opt:     &IssuesComponentTagsOptions{},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesComponentTagsOptions{ComponentUuid: "uuid"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateComponentTagsOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateDeleteCommentOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesDeleteCommentOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Comment",
			opt:     &IssuesDeleteCommentOptions{},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesDeleteCommentOptions{Comment: "comment-key"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateDeleteCommentOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateEditCommentOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesEditCommentOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Comment",
			opt:     &IssuesEditCommentOptions{Text: "new text"},
			wantErr: true,
		},
		{
			name:    "missing Text",
			opt:     &IssuesEditCommentOptions{Comment: "comment-key"},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesEditCommentOptions{Comment: "comment-key", Text: "new text"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateEditCommentOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateReindexOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesReindexOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Project",
			opt:     &IssuesReindexOptions{},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesReindexOptions{Project: "project-key"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateReindexOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateSetTagsOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesSetTagsOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Issue",
			opt:     &IssuesSetTagsOptions{},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesSetTagsOptions{Issue: "issue-key"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateSetTagsOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidatePullTaintOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesPullTaintOptions
		wantErr bool
	}{
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing ProjectKey",
			opt:     &IssuesPullTaintOptions{},
			wantErr: true,
		},
		{
			name:    "invalid language",
			opt:     &IssuesPullTaintOptions{ProjectKey: "project", Languages: []string{"invalid"}},
			wantErr: true,
		},
		{
			name:    "valid option",
			opt:     &IssuesPullTaintOptions{ProjectKey: "project", Languages: []string{"java"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidatePullTaintOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIssues_ValidateTagsOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *IssuesTagsOptions
		wantErr bool
	}{
		{
			name:    "nil option (valid)",
			opt:     nil,
			wantErr: false,
		},
		{
			name:    "PageSize 0 (valid)",
			opt:     &IssuesTagsOptions{PageSize: 0},
			wantErr: false,
		},
		{
			name:    "PageSize 501 (invalid)",
			opt:     &IssuesTagsOptions{PageSize: 501},
			wantErr: true,
		},
		{
			name:    "PageSize 100 (valid)",
			opt:     &IssuesTagsOptions{PageSize: 100},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Issues.ValidateTagsOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
