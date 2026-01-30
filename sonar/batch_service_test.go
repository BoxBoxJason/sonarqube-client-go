package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBatchService_File(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Query().Get("name") != "batch-library-2.3.jar" {
				t.Errorf("expected name 'batch-library-2.3.jar', got '%s'", r.URL.Query().Get("name"))
			}
			w.Header().Set("Content-Type", "application/java-archive")
			_, _ = w.Write([]byte("jar-binary-content"))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Batch.File(&BatchFileOption{
			Name: "batch-library-2.3.jar",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})

	t.Run("nil option", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/java-archive")
			_, _ = w.Write([]byte("jar-binary-content"))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Batch.File(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty option", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/java-archive")
			_, _ = w.Write([]byte("jar-binary-content"))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Batch.File(&BatchFileOption{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestBatchService_Index(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "text/plain")
			_, err := w.Write([]byte("batch-library-2.3.jar|abc123def456\nscanner-engine-9.0.jar|789xyz"))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Batch.Index()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
	})
}

func TestBatchService_Project(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Query().Get("key") != "my-project" {
				t.Errorf("expected key 'my-project', got '%s'", r.URL.Query().Get("key"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"fileDataByModuleAndPath": {
					"my-project": {
						"src/main/java/App.java": {
							"hash": "abc123",
							"revision": "1"
						}
					}
				},
				"lastAnalysisDate": 1640000000000,
				"timestamp": 1640000001000
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Batch.Project(&BatchProjectOption{
			Key: "my-project",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if result.LastAnalysisDate != 1640000000000 {
			t.Errorf("expected lastAnalysisDate 1640000000000, got %d", result.LastAnalysisDate)
		}
	})

	t.Run("with branch", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("branch") != "feature/my-branch" {
				t.Errorf("expected branch 'feature/my-branch', got '%s'", r.URL.Query().Get("branch"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"fileDataByModuleAndPath": {}, "lastAnalysisDate": 0, "timestamp": 0}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Batch.Project(&BatchProjectOption{
			Key:    "my-project",
			Branch: "feature/my-branch",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with pull request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("pullRequest") != "5461" {
				t.Errorf("expected pullRequest '5461', got '%s'", r.URL.Query().Get("pullRequest"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"fileDataByModuleAndPath": {}, "lastAnalysisDate": 0, "timestamp": 0}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Batch.Project(&BatchProjectOption{
			Key:         "my-project",
			PullRequest: "5461",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"fileDataByModuleAndPath": {}, "lastAnalysisDate": 0, "timestamp": 0}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Batch.Project(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestBatchService_ValidateFileOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *BatchFileOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &BatchFileOption{}, false},
		{"with name", &BatchFileOption{Name: "test.jar"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Batch.ValidateFileOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBatchService_ValidateProjectOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *BatchProjectOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &BatchProjectOption{}, false},
		{"with key", &BatchProjectOption{Key: "my-project"}, false},
		{"with branch", &BatchProjectOption{Key: "my-project", Branch: "main"}, false},
		{"with pull request", &BatchProjectOption{Key: "my-project", PullRequest: "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Batch.ValidateProjectOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
