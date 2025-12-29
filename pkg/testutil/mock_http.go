package testutil

import (
	"net/http"
	"net/http/httptest"
)

// MockHTTPClient provides a mock HTTP client for testing.
type MockHTTPClient struct {
	Server    *httptest.Server
	Mux       *http.ServeMux
	Responses map[string]MockResponse
}

// MockResponse defines a mock HTTP response.
type MockResponse struct {
	Headers    map[string]string
	Body       string
	StatusCode int
}

// NewMockHTTPClient creates a new mock HTTP client.
func NewMockHTTPClient() *MockHTTPClient {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return &MockHTTPClient{
		Server:    server,
		Mux:       mux,
		Responses: make(map[string]MockResponse),
	}
}

// SetResponse sets a mock response for a specific path.
func (m *MockHTTPClient) SetResponse(path string, resp MockResponse) {
	m.Responses[path] = resp
	m.Mux.HandleFunc(path, func(respWriter http.ResponseWriter, _ *http.Request) {
		for k, v := range resp.Headers {
			respWriter.Header().Set(k, v)
		}

		if resp.Headers["Content-Type"] == "" {
			respWriter.Header().Set("Content-Type", "application/json")
		}

		respWriter.WriteHeader(resp.StatusCode)

		_, _ = respWriter.Write([]byte(resp.Body))
	})
}

// URL returns the base URL of the mock server.
func (m *MockHTTPClient) URL() string {
	return m.Server.URL
}

// Close shuts down the mock server.
func (m *MockHTTPClient) Close() {
	m.Server.Close()
}

// Client returns an HTTP client configured to use the mock server.
func (m *MockHTTPClient) Client() *http.Client {
	return m.Server.Client()
}
