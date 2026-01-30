package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCe_Activity(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ce/activity", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("component"))
		assert.Equal(t, "SUCCESS", r.URL.Query().Get("status"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 1,
				"pageSize": 100,
				"total": 2
			},
			"tasks": [
				{
					"id": "task-1",
					"type": "REPORT",
					"status": "SUCCESS",
					"componentKey": "my-project",
					"componentName": "My Project",
					"executionTimeMs": 1234,
					"submittedAt": "2023-01-01T10:00:00+0000",
					"startedAt": "2023-01-01T10:00:01+0000",
					"executedAt": "2023-01-01T10:00:02+0000"
				},
				{
					"id": "task-2",
					"type": "REPORT",
					"status": "SUCCESS",
					"componentKey": "my-project",
					"componentName": "My Project",
					"executionTimeMs": 2345
				}
			]
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeActivityOption{
		Component: "my-project",
		Statuses:  []string{"SUCCESS"},
	}

	result, resp, err := client.Ce.Activity(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, int64(2), result.Paging.Total)
	assert.Len(t, result.Tasks, 2)
	assert.Equal(t, "task-1", result.Tasks[0].ID)
	assert.Equal(t, "SUCCESS", result.Tasks[0].Status)
	assert.Equal(t, int64(1234), result.Tasks[0].ExecutionTimeMs)
}

func TestCe_Activity_WithPagination(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "2", r.URL.Query().Get("p"))
		assert.Equal(t, "50", r.URL.Query().Get("ps"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"paging": {
				"pageIndex": 2,
				"pageSize": 50,
				"total": 100
			},
			"tasks": []
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeActivityOption{
		CePaginationArgs: CePaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	}

	result, resp, err := client.Ce.Activity(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(2), result.Paging.PageIndex)
	assert.Equal(t, int64(50), result.Paging.PageSize)
}

func TestCe_Activity_WithNilOption(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}, "tasks": []}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Ce.Activity(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestCe_Activity_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *CeActivityOption
	}{
		{
			name: "invalid status",
			opt:  &CeActivityOption{Statuses: []string{"INVALID_STATUS"}},
		},
		{
			name: "invalid type",
			opt:  &CeActivityOption{Type: "INVALID_TYPE"},
		},
		{
			name: "page size too large",
			opt:  &CeActivityOption{CePaginationArgs: CePaginationArgs{PageSize: 1001}},
		},
		{
			name: "page size too small",
			opt:  &CeActivityOption{CePaginationArgs: CePaginationArgs{PageSize: 0}},
		},
		{
			name: "page too small",
			opt:  &CeActivityOption{CePaginationArgs: CePaginationArgs{Page: 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Ce.Activity(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestCe_ActivityStatus(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ce/activity_status", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("component"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"failing": 2,
			"inProgress": 1,
			"pending": 5,
			"pendingTime": 12345
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeActivityStatusOption{
		Component: "my-project",
	}

	result, resp, err := client.Ce.ActivityStatus(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, int64(2), result.Failing)
	assert.Equal(t, int64(1), result.InProgress)
	assert.Equal(t, int64(5), result.Pending)
	assert.Equal(t, int64(12345), result.PendingTime)
}

func TestCe_ActivityStatus_WithNilOption(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"failing": 0, "inProgress": 0, "pending": 0}`))
	})

	client := newTestClient(t, server.URL)

	result, resp, err := client.Ce.ActivityStatus(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestCe_AnalysisStatus(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ce/analysis_status", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("component"))
		assert.Equal(t, "main", r.URL.Query().Get("branch"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"component": {
				"key": "my-project",
				"name": "My Project",
				"warnings": [
					{
						"key": "warning-1",
						"message": "This is a warning",
						"dismissable": true
					}
				]
			}
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeAnalysisStatusOption{
		Component: "my-project",
		Branch:    "main",
	}

	result, resp, err := client.Ce.AnalysisStatus(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "my-project", result.Component.Key)
	assert.Len(t, result.Component.Warnings, 1)
	assert.Equal(t, "warning-1", result.Component.Warnings[0].Key)
	assert.True(t, result.Component.Warnings[0].Dismissable)
}

func TestCe_AnalysisStatus_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Ce.AnalysisStatus(nil)
	assert.Error(t, err)

	// Missing Component should fail validation.
	_, _, err = client.Ce.AnalysisStatus(&CeAnalysisStatusOption{})
	assert.Error(t, err)
}

func TestCe_Cancel(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/ce/cancel", r.URL.Path)
		assert.Equal(t, "task-123", r.URL.Query().Get("id"))

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	opt := &CeCancelOption{
		ID: "task-123",
	}

	resp, err := client.Ce.Cancel(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestCe_Cancel_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.Ce.Cancel(nil)
	assert.Error(t, err)

	// Missing ID should fail validation.
	_, err = client.Ce.Cancel(&CeCancelOption{})
	assert.Error(t, err)
}

func TestCe_CancelAll(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/ce/cancel_all", http.StatusNoContent))

	client := newTestClient(t, server.URL)

	resp, err := client.Ce.CancelAll()
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestCe_Component(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ce/component", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("component"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"current": {
				"id": "current-task",
				"type": "REPORT",
				"status": "SUCCESS",
				"componentKey": "my-project",
				"executionTimeMs": 5000
			},
			"queue": [
				{
					"id": "queued-task-1",
					"type": "REPORT",
					"status": "PENDING",
					"componentKey": "my-project"
				}
			]
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeComponentOption{
		Component: "my-project",
	}

	result, resp, err := client.Ce.Component(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "current-task", result.Current.ID)
	assert.Equal(t, "SUCCESS", result.Current.Status)
	assert.Len(t, result.Queue, 1)
	assert.Equal(t, "queued-task-1", result.Queue[0].ID)
	assert.Equal(t, "PENDING", result.Queue[0].Status)
}

func TestCe_Component_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.Ce.Component(nil)
	assert.Error(t, err)

	// Missing Component should fail validation.
	_, _, err = client.Ce.Component(&CeComponentOption{})
	assert.Error(t, err)
}

func TestCe_DismissAnalysisWarning(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/ce/dismiss_analysis_warning", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("component"))
		assert.Equal(t, "warning-key-1", r.URL.Query().Get("warning"))

		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, server.URL)

	opt := &CeDismissAnalysisWarningOption{
		Component: "my-project",
		Warning:   "warning-key-1",
	}

	resp, err := client.Ce.DismissAnalysisWarning(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestCe_DismissAnalysisWarning_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *CeDismissAnalysisWarningOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing component",
			opt:  &CeDismissAnalysisWarningOption{Warning: "warning-key"},
		},
		{
			name: "missing warning",
			opt:  &CeDismissAnalysisWarningOption{Component: "my-project"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Ce.DismissAnalysisWarning(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestCe_IndexationStatus(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/ce/indexation_status", http.StatusOK, `{
		"completedCount": 50,
		"hasFailures": false,
		"isCompleted": true,
		"total": 50
	}`))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Ce.IndexationStatus()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, int64(50), result.CompletedCount)
	assert.False(t, result.HasFailures)
	assert.True(t, result.IsCompleted)
	assert.Equal(t, int64(50), result.Total)
}

func TestCe_Info(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/ce/info", http.StatusOK, `{
		"workersPauseStatus": "RUNNING"
	}`))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Ce.Info()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "RUNNING", result.WorkersPauseStatus)
}

func TestCe_Pause(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/ce/pause", http.StatusNoContent))

	client := newTestClient(t, server.URL)

	resp, err := client.Ce.Pause()
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestCe_Resume(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/ce/resume", http.StatusNoContent))

	client := newTestClient(t, server.URL)

	resp, err := client.Ce.Resume()
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestCe_Submit(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/ce/submit", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("projectKey"))
		assert.Equal(t, "My Project", r.URL.Query().Get("projectName"))
		assert.Equal(t, "base64-encoded-report", r.URL.Query().Get("report"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"projectId": "project-uuid",
			"taskId": "task-uuid"
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeSubmitOption{
		ProjectKey:  "my-project",
		ProjectName: "My Project",
		Report:      "base64-encoded-report",
	}

	result, resp, err := client.Ce.Submit(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "project-uuid", result.ProjectID)
	assert.Equal(t, "task-uuid", result.TaskID)
}

func TestCe_Submit_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *CeSubmitOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing project key",
			opt:  &CeSubmitOption{Report: "report"},
		},
		{
			name: "missing report",
			opt:  &CeSubmitOption{ProjectKey: "my-project"},
		},
		{
			name: "project key too long",
			opt: &CeSubmitOption{
				ProjectKey: string(make([]byte, 401)),
				Report:     "report",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Ce.Submit(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestCe_Task(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ce/task", r.URL.Path)
		assert.Equal(t, "task-123", r.URL.Query().Get("id"))
		assert.Equal(t, "stacktrace", r.URL.Query().Get("additionalFields"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"task": {
				"id": "task-123",
				"type": "REPORT",
				"status": "FAILED",
				"componentKey": "my-project",
				"componentName": "My Project",
				"errorMessage": "Analysis failed",
				"errorStacktrace": "java.lang.Exception: error\n  at ...",
				"hasErrorStacktrace": true,
				"executionTimeMs": 12345,
				"warningCount": 3,
				"warnings": ["Warning 1", "Warning 2", "Warning 3"]
			}
		}`))
	})

	client := newTestClient(t, server.URL)

	opt := &CeTaskOption{
		ID:               "task-123",
		AdditionalFields: []string{"stacktrace"},
	}

	result, resp, err := client.Ce.Task(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "task-123", result.Task.ID)
	assert.Equal(t, "FAILED", result.Task.Status)
	assert.Equal(t, "Analysis failed", result.Task.ErrorMessage)
	assert.True(t, result.Task.HasErrorStacktrace)
	assert.Equal(t, int64(3), result.Task.WarningCount)
	assert.Len(t, result.Task.Warnings, 3)
}

func TestCe_Task_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *CeTaskOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing id",
			opt:  &CeTaskOption{},
		},
		{
			name: "invalid additional field",
			opt:  &CeTaskOption{ID: "task-123", AdditionalFields: []string{"invalid"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Ce.Task(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestCe_TaskTypes(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/ce/task_types", http.StatusOK, `{
		"taskTypes": ["REPORT", "ISSUE_SYNC", "AUDIT_PURGE", "PROJECT_EXPORT"]
	}`))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Ce.TaskTypes()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.TaskTypes, 4)

	expectedTypes := []string{"REPORT", "ISSUE_SYNC", "AUDIT_PURGE", "PROJECT_EXPORT"}
	assert.Equal(t, expectedTypes, result.TaskTypes)
}

func TestCe_WorkerCount(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/ce/worker_count", http.StatusOK, `{
		"canSetWorkerCount": true,
		"value": 4
	}`))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Ce.WorkerCount()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.True(t, result.CanSetWorkerCount)
	assert.Equal(t, int64(4), result.Value)
}

// Test CePaginationArgs validation
func TestCePaginationArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    CePaginationArgs
		wantErr bool
	}{
		{
			name:    "empty is valid",
			args:    CePaginationArgs{},
			wantErr: false,
		},
		{
			name:    "valid page and size",
			args:    CePaginationArgs{Page: 1, PageSize: 100},
			wantErr: false,
		},
		{
			name:    "max page size 1000 is valid",
			args:    CePaginationArgs{PageSize: 1000},
			wantErr: false,
		},
		{
			name:    "page size 1001 is invalid",
			args:    CePaginationArgs{PageSize: 1001},
			wantErr: true,
		},
		{
			name:    "page 0 is invalid when set",
			args:    CePaginationArgs{Page: 0, PageSize: 100},
			wantErr: false, // 0 means not set
		},
		{
			name:    "negative page is invalid",
			args:    CePaginationArgs{Page: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validation methods
func TestCe_ValidateActivityOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *CeActivityOption
		wantErr bool
	}{
		{
			name:    "nil option is valid",
			opt:     nil,
			wantErr: false,
		},
		{
			name:    "empty option is valid",
			opt:     &CeActivityOption{},
			wantErr: false,
		},
		{
			name: "valid statuses",
			opt: &CeActivityOption{
				Statuses: []string{"SUCCESS", "FAILED"},
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			opt: &CeActivityOption{
				Statuses: []string{"SUCCESS", "INVALID"},
			},
			wantErr: true,
		},
		{
			name: "valid type",
			opt: &CeActivityOption{
				Type: "REPORT",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			opt: &CeActivityOption{
				Type: "INVALID_TYPE",
			},
			wantErr: true,
		},
		{
			name: "valid all fields",
			opt: &CeActivityOption{
				Component:    "my-project",
				Statuses:     []string{"SUCCESS", "FAILED", "PENDING"},
				Type:         "REPORT",
				OnlyCurrents: true,
				CePaginationArgs: CePaginationArgs{
					Page:     1,
					PageSize: 100,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Ce.ValidateActivityOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCe_ValidateTaskOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *CeTaskOption
		wantErr bool
	}{
		{
			name:    "nil option is invalid",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing id is invalid",
			opt:     &CeTaskOption{},
			wantErr: true,
		},
		{
			name: "valid with id only",
			opt: &CeTaskOption{
				ID: "task-123",
			},
			wantErr: false,
		},
		{
			name: "valid with all additional fields",
			opt: &CeTaskOption{
				ID:               "task-123",
				AdditionalFields: []string{"stacktrace", "scannerContext", "warnings"},
			},
			wantErr: false,
		},
		{
			name: "invalid additional field",
			opt: &CeTaskOption{
				ID:               "task-123",
				AdditionalFields: []string{"invalid"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Ce.ValidateTaskOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
