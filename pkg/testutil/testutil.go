// Package testutil provides common test utilities and mock implementations
// for unit testing across the sonarqube-client-go project.
package testutil

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestServer wraps httptest.Server with additional helper methods
type TestServer struct {
	*httptest.Server
	Mux *http.ServeMux
}

// NewTestServer creates a new test server with a mux for routing
func NewTestServer() *TestServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &TestServer{
		Server: server,
		Mux:    mux,
	}
}

// Close shuts down the test server
func (s *TestServer) Close() {
	s.Server.Close()
}

// URL returns the base URL of the test server
func (s *TestServer) URL() string {
	return s.Server.URL
}

// HandleFunc registers a handler function for a given pattern
func (s *TestServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.Mux.HandleFunc(pattern, handler)
}

// JSONResponse is a helper to respond with JSON content
func JSONResponse(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

// TextResponse is a helper to respond with text content
func TextResponse(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

// TempDir creates a temporary directory for testing and returns a cleanup function
func TempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "sonarqube-client-go-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir, func() {
		os.RemoveAll(dir)
	}
}

// WriteFile writes content to a file in the specified directory
func WriteFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
	return path
}

// ReadFile reads content from a file
func ReadFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(content)
}

// MockReadCloser creates a ReadCloser from a string for testing
func MockReadCloser(content string) io.ReadCloser {
	return io.NopCloser(bytes.NewBufferString(content))
}

// AssertEqual is a simple assertion helper
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNoError asserts that an error is nil
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: unexpected error: %v", msg, err)
	}
}

// AssertError asserts that an error is not nil
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error, got nil", msg)
	}
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, s, substr, msg string) {
	t.Helper()
	if !bytes.Contains([]byte(s), []byte(substr)) {
		t.Errorf("%s: expected %q to contain %q", msg, s, substr)
	}
}
