package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testServer wraps httptest.Server for testing.
type testServer struct {
	*httptest.Server
	t *testing.T
}

// newTestServer creates a new test server with automatic cleanup.
func newTestServer(t *testing.T, handler http.HandlerFunc) *testServer {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	return &testServer{Server: ts, t: t}
}

// url returns the server URL with /api/ suffix.
func (s *testServer) url() string {
	return s.Server.URL + "/"
}

// mockHandler creates a handler that validates method and path, then returns a JSON response.
func mockHandler(t *testing.T, method, path string, statusCode int, response any) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method, "unexpected HTTP method")
		assert.Equal(t, path, r.URL.Path, "unexpected URL path")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if response != nil {
			switch v := response.(type) {
			case string:
				_, _ = w.Write([]byte(v))
			case []byte:
				_, _ = w.Write(v)
			default:
				_ = json.NewEncoder(w).Encode(response)
			}
		}
	}
}

// mockEmptyHandler creates a handler that validates method and path, then returns no body.
func mockEmptyHandler(t *testing.T, method, path string, statusCode int) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method, "unexpected HTTP method")
		assert.Equal(t, path, r.URL.Path, "unexpected URL path")

		w.WriteHeader(statusCode)
	}
}

// mockHandlerWithParams creates a handler that validates method, path, and query parameters, then returns a JSON response.
func mockHandlerWithParams(t *testing.T, method, path string, statusCode int, expectedParams map[string]string, response any) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method, "unexpected HTTP method")
		assert.Equal(t, path, r.URL.Path, "unexpected URL path")

		// Verify expected query parameters
		query := r.URL.Query()
		for key, expectedValue := range expectedParams {
			assert.Equal(t, expectedValue, query.Get(key), "unexpected value for query param %q", key)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if response != nil {
			switch v := response.(type) {
			case string:
				_, _ = w.Write([]byte(v))
			case []byte:
				_, _ = w.Write(v)
			default:
				_ = json.NewEncoder(w).Encode(response)
			}
		}
	}
}

// mockEmptyHandlerWithParams creates a handler that validates method, path, and query parameters, then returns no body.
func mockEmptyHandlerWithParams(t *testing.T, method, path string, statusCode int, expectedParams map[string]string) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method, "unexpected HTTP method")
		assert.Equal(t, path, r.URL.Path, "unexpected URL path")

		// Verify expected query parameters
		query := r.URL.Query()
		for key, expectedValue := range expectedParams {
			assert.Equal(t, expectedValue, query.Get(key), "unexpected value for query param %q", key)
		}

		w.WriteHeader(statusCode)
	}
}

// mockBinaryHandler creates a handler that returns binary content.
func mockBinaryHandler(t *testing.T, method, path string, statusCode int, contentType string, data []byte) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method, "unexpected HTTP method")
		assert.Equal(t, path, r.URL.Path, "unexpected URL path")

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(statusCode)
		_, _ = w.Write(data)
	}
}

// newTestClient creates a new client for the given test server URL.
func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()

	client, err := NewClient(nil, WithBaseURL(serverURL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client
}

// newLocalhostClient creates a client for validation-only tests (no real server needed).
func newLocalhostClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewClient(nil, WithBaseURL("http://localhost/api/"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client
}
