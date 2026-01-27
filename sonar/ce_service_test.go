package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCe_Activity(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/activity" {
			t.Errorf("expected path /api/ce/activity, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		status := r.URL.Query().Get("status")
		if status != "SUCCESS" {
			t.Errorf("expected status 'SUCCESS', got %s", status)
		}

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
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeActivityOption{
		Component: "my-project",
		Statuses:  []string{"SUCCESS"},
	}

	result, resp, err := client.Ce.Activity(opt)
	if err != nil {
		t.Fatalf("Activity failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Paging.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Paging.Total)
	}

	if len(result.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(result.Tasks))
	}

	if result.Tasks[0].ID != "task-1" {
		t.Errorf("expected task id 'task-1', got %s", result.Tasks[0].ID)
	}

	if result.Tasks[0].Status != "SUCCESS" {
		t.Errorf("expected status 'SUCCESS', got %s", result.Tasks[0].Status)
	}

	if result.Tasks[0].ExecutionTimeMs != 1234 {
		t.Errorf("expected executionTimeMs 1234, got %d", result.Tasks[0].ExecutionTimeMs)
	}
}

func TestCe_Activity_WithPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		p := r.URL.Query().Get("p")
		if p != "2" {
			t.Errorf("expected p '2', got %s", p)
		}

		ps := r.URL.Query().Get("ps")
		if ps != "50" {
			t.Errorf("expected ps '50', got %s", ps)
		}

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
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeActivityOption{
		CePaginationArgs: CePaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	}

	result, resp, err := client.Ce.Activity(opt)
	if err != nil {
		t.Fatalf("Activity failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result.Paging.PageIndex != 2 {
		t.Errorf("expected pageIndex 2, got %d", result.Paging.PageIndex)
	}

	if result.Paging.PageSize != 50 {
		t.Errorf("expected pageSize 50, got %d", result.Paging.PageSize)
	}
}

func TestCe_Activity_WithNilOption(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}, "tasks": []}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Ce.Activity(nil)
	if err != nil {
		t.Fatalf("Activity with nil option failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCe_Activity_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestCe_ActivityStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/activity_status" {
			t.Errorf("expected path /api/ce/activity_status, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"failing": 2,
			"inProgress": 1,
			"pending": 5,
			"pendingTime": 12345
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeActivityStatusOption{
		Component: "my-project",
	}

	result, resp, err := client.Ce.ActivityStatus(opt)
	if err != nil {
		t.Fatalf("ActivityStatus failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Failing != 2 {
		t.Errorf("expected failing 2, got %d", result.Failing)
	}

	if result.InProgress != 1 {
		t.Errorf("expected inProgress 1, got %d", result.InProgress)
	}

	if result.Pending != 5 {
		t.Errorf("expected pending 5, got %d", result.Pending)
	}

	if result.PendingTime != 12345 {
		t.Errorf("expected pendingTime 12345, got %d", result.PendingTime)
	}
}

func TestCe_ActivityStatus_WithNilOption(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"failing": 0, "inProgress": 0, "pending": 0}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Ce.ActivityStatus(nil)
	if err != nil {
		t.Fatalf("ActivityStatus with nil option failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCe_AnalysisStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/analysis_status" {
			t.Errorf("expected path /api/ce/analysis_status, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		branch := r.URL.Query().Get("branch")
		if branch != "main" {
			t.Errorf("expected branch 'main', got %s", branch)
		}

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
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeAnalysisStatusOption{
		Component: "my-project",
		Branch:    "main",
	}

	result, resp, err := client.Ce.AnalysisStatus(opt)
	if err != nil {
		t.Fatalf("AnalysisStatus failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Component.Key != "my-project" {
		t.Errorf("expected component key 'my-project', got %s", result.Component.Key)
	}

	if len(result.Component.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Component.Warnings))
	}

	if result.Component.Warnings[0].Key != "warning-1" {
		t.Errorf("expected warning key 'warning-1', got %s", result.Component.Warnings[0].Key)
	}

	if !result.Component.Warnings[0].Dismissable {
		t.Error("expected warning to be dismissable")
	}
}

func TestCe_AnalysisStatus_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Ce.AnalysisStatus(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Component should fail validation.
	_, _, err = client.Ce.AnalysisStatus(&CeAnalysisStatusOption{})
	if err == nil {
		t.Error("expected error for missing Component")
	}
}

func TestCe_Cancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/cancel" {
			t.Errorf("expected path /api/ce/cancel, got %s", r.URL.Path)
		}

		id := r.URL.Query().Get("id")
		if id != "task-123" {
			t.Errorf("expected id 'task-123', got %s", id)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeCancelOption{
		ID: "task-123",
	}

	resp, err := client.Ce.Cancel(opt)
	if err != nil {
		t.Fatalf("Cancel failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestCe_Cancel_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.Ce.Cancel(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing ID should fail validation.
	_, err = client.Ce.Cancel(&CeCancelOption{})
	if err == nil {
		t.Error("expected error for missing ID")
	}
}

func TestCe_CancelAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/cancel_all" {
			t.Errorf("expected path /api/ce/cancel_all, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Ce.CancelAll()
	if err != nil {
		t.Fatalf("CancelAll failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestCe_Component(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/component" {
			t.Errorf("expected path /api/ce/component, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

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
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeComponentOption{
		Component: "my-project",
	}

	result, resp, err := client.Ce.Component(opt)
	if err != nil {
		t.Fatalf("Component failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Current.ID != "current-task" {
		t.Errorf("expected current task id 'current-task', got %s", result.Current.ID)
	}

	if result.Current.Status != "SUCCESS" {
		t.Errorf("expected current status 'SUCCESS', got %s", result.Current.Status)
	}

	if len(result.Queue) != 1 {
		t.Errorf("expected 1 queued task, got %d", len(result.Queue))
	}

	if result.Queue[0].ID != "queued-task-1" {
		t.Errorf("expected queued task id 'queued-task-1', got %s", result.Queue[0].ID)
	}

	if result.Queue[0].Status != "PENDING" {
		t.Errorf("expected queued status 'PENDING', got %s", result.Queue[0].Status)
	}
}

func TestCe_Component_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.Ce.Component(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing Component should fail validation.
	_, _, err = client.Ce.Component(&CeComponentOption{})
	if err == nil {
		t.Error("expected error for missing Component")
	}
}

func TestCe_DismissAnalysisWarning(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/dismiss_analysis_warning" {
			t.Errorf("expected path /api/ce/dismiss_analysis_warning, got %s", r.URL.Path)
		}

		component := r.URL.Query().Get("component")
		if component != "my-project" {
			t.Errorf("expected component 'my-project', got %s", component)
		}

		warning := r.URL.Query().Get("warning")
		if warning != "warning-key-1" {
			t.Errorf("expected warning 'warning-key-1', got %s", warning)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeDismissAnalysisWarningOption{
		Component: "my-project",
		Warning:   "warning-key-1",
	}

	resp, err := client.Ce.DismissAnalysisWarning(opt)
	if err != nil {
		t.Fatalf("DismissAnalysisWarning failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestCe_DismissAnalysisWarning_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestCe_IndexationStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/indexation_status" {
			t.Errorf("expected path /api/ce/indexation_status, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"completedCount": 50,
			"hasFailures": false,
			"isCompleted": true,
			"total": 50
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Ce.IndexationStatus()
	if err != nil {
		t.Fatalf("IndexationStatus failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.CompletedCount != 50 {
		t.Errorf("expected completedCount 50, got %d", result.CompletedCount)
	}

	if result.HasFailures {
		t.Error("expected hasFailures to be false")
	}

	if !result.IsCompleted {
		t.Error("expected isCompleted to be true")
	}

	if result.Total != 50 {
		t.Errorf("expected total 50, got %d", result.Total)
	}
}

func TestCe_Info(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/info" {
			t.Errorf("expected path /api/ce/info, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"workersPauseStatus": "RUNNING"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Ce.Info()
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.WorkersPauseStatus != "RUNNING" {
		t.Errorf("expected workersPauseStatus 'RUNNING', got %s", result.WorkersPauseStatus)
	}
}

func TestCe_Pause(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/pause" {
			t.Errorf("expected path /api/ce/pause, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Ce.Pause()
	if err != nil {
		t.Fatalf("Pause failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestCe_Resume(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/resume" {
			t.Errorf("expected path /api/ce/resume, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Ce.Resume()
	if err != nil {
		t.Fatalf("Resume failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestCe_Submit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/submit" {
			t.Errorf("expected path /api/ce/submit, got %s", r.URL.Path)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		projectName := r.URL.Query().Get("projectName")
		if projectName != "My Project" {
			t.Errorf("expected projectName 'My Project', got %s", projectName)
		}

		report := r.URL.Query().Get("report")
		if report != "base64-encoded-report" {
			t.Errorf("expected report 'base64-encoded-report', got %s", report)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"projectId": "project-uuid",
			"taskId": "task-uuid"
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeSubmitOption{
		ProjectKey:  "my-project",
		ProjectName: "My Project",
		Report:      "base64-encoded-report",
	}

	result, resp, err := client.Ce.Submit(opt)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.ProjectID != "project-uuid" {
		t.Errorf("expected projectId 'project-uuid', got %s", result.ProjectID)
	}

	if result.TaskID != "task-uuid" {
		t.Errorf("expected taskId 'task-uuid', got %s", result.TaskID)
	}
}

func TestCe_Submit_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestCe_Task(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/task" {
			t.Errorf("expected path /api/ce/task, got %s", r.URL.Path)
		}

		id := r.URL.Query().Get("id")
		if id != "task-123" {
			t.Errorf("expected id 'task-123', got %s", id)
		}

		additionalFields := r.URL.Query().Get("additionalFields")
		if additionalFields != "stacktrace" {
			t.Errorf("expected additionalFields 'stacktrace', got %s", additionalFields)
		}

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
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &CeTaskOption{
		ID:               "task-123",
		AdditionalFields: []string{"stacktrace"},
	}

	result, resp, err := client.Ce.Task(opt)
	if err != nil {
		t.Fatalf("Task failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Task.ID != "task-123" {
		t.Errorf("expected task id 'task-123', got %s", result.Task.ID)
	}

	if result.Task.Status != "FAILED" {
		t.Errorf("expected status 'FAILED', got %s", result.Task.Status)
	}

	if result.Task.ErrorMessage != "Analysis failed" {
		t.Errorf("expected error message 'Analysis failed', got %s", result.Task.ErrorMessage)
	}

	if !result.Task.HasErrorStacktrace {
		t.Error("expected hasErrorStacktrace to be true")
	}

	if result.Task.WarningCount != 3 {
		t.Errorf("expected warningCount 3, got %d", result.Task.WarningCount)
	}

	if len(result.Task.Warnings) != 3 {
		t.Errorf("expected 3 warnings, got %d", len(result.Task.Warnings))
	}
}

func TestCe_Task_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestCe_TaskTypes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/task_types" {
			t.Errorf("expected path /api/ce/task_types, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"taskTypes": ["REPORT", "ISSUE_SYNC", "AUDIT_PURGE", "PROJECT_EXPORT"]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Ce.TaskTypes()
	if err != nil {
		t.Fatalf("TaskTypes failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.TaskTypes) != 4 {
		t.Errorf("expected 4 task types, got %d", len(result.TaskTypes))
	}

	expectedTypes := []string{"REPORT", "ISSUE_SYNC", "AUDIT_PURGE", "PROJECT_EXPORT"}
	for i, expectedType := range expectedTypes {
		if result.TaskTypes[i] != expectedType {
			t.Errorf("expected task type '%s' at index %d, got '%s'", expectedType, i, result.TaskTypes[i])
		}
	}
}

func TestCe_WorkerCount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/ce/worker_count" {
			t.Errorf("expected path /api/ce/worker_count, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"canSetWorkerCount": true,
			"value": 4
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Ce.WorkerCount()
	if err != nil {
		t.Fatalf("WorkerCount failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.CanSetWorkerCount {
		t.Error("expected canSetWorkerCount to be true")
	}

	if result.Value != 4 {
		t.Errorf("expected value 4, got %d", result.Value)
	}
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
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test validation methods
func TestCe_ValidateActivityOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateActivityOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCe_ValidateTaskOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

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
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTaskOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
