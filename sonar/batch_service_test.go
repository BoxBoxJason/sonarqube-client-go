package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchService_File(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := mockBinaryHandler(t, http.MethodGet, "/batch/file", http.StatusOK, "application/java-archive", []byte("jar-binary-content"))
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		result, resp, err := client.Batch.File(&BatchFileOption{
			Name: "batch-library-2.3.jar",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, result)
	})

	t.Run("nil option", func(t *testing.T) {
		handler := mockBinaryHandler(t, http.MethodGet, "/batch/file", http.StatusOK, "application/java-archive", []byte("jar-binary-content"))
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Batch.File(nil)

		require.NoError(t, err)
	})

	t.Run("empty option", func(t *testing.T) {
		handler := mockBinaryHandler(t, http.MethodGet, "/batch/file", http.StatusOK, "application/java-archive", []byte("jar-binary-content"))
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Batch.File(&BatchFileOption{})

		require.NoError(t, err)
	})
}

func TestBatchService_Index(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := mockBinaryHandler(t, http.MethodGet, "/batch/index", http.StatusOK, "text/plain", []byte("batch-library-2.3.jar|abc123def456\nscanner-engine-9.0.jar|789xyz"))
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		result, resp, err := client.Batch.Index()

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, result)
	})
}

func TestBatchService_Project(t *testing.T) {
	projectJSON := `{
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
	}`

	t.Run("success", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/batch/project", http.StatusOK, projectJSON)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		result, resp, err := client.Batch.Project(&BatchProjectOption{
			Key: "my-project",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
		assert.Equal(t, int64(1640000000000), result.LastAnalysisDate)
	})

	t.Run("with branch", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/batch/project", http.StatusOK, `{"fileDataByModuleAndPath": {}, "lastAnalysisDate": 0, "timestamp": 0}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Batch.Project(&BatchProjectOption{
			Key:    "my-project",
			Branch: "feature/my-branch",
		})

		require.NoError(t, err)
	})

	t.Run("with pull request", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/batch/project", http.StatusOK, `{"fileDataByModuleAndPath": {}, "lastAnalysisDate": 0, "timestamp": 0}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Batch.Project(&BatchProjectOption{
			Key:         "my-project",
			PullRequest: "5461",
		})

		require.NoError(t, err)
	})

	t.Run("nil option", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/batch/project", http.StatusOK, `{"fileDataByModuleAndPath": {}, "lastAnalysisDate": 0, "timestamp": 0}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Batch.Project(nil)

		require.NoError(t, err)
	})
}

func TestBatchService_ValidateFileOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBatchService_ValidateProjectOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *BatchProjectOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &BatchProjectOption{}, true},
		{"with key", &BatchProjectOption{Key: "my-project"}, false},
		{"with branch", &BatchProjectOption{Key: "my-project", Branch: "main"}, false},
		{"with pull request", &BatchProjectOption{Key: "my-project", PullRequest: "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Batch.ValidateProjectOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
