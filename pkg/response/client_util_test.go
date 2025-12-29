package response

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSetBaseURLUtilFunction(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPath    string
		expectError bool
	}{
		{
			name:        "URL without trailing slash",
			input:       "http://localhost:9000/api",
			wantPath:    "/api/",
			expectError: false,
		},
		{
			name:        "URL with trailing slash",
			input:       "http://localhost:9000/api/",
			wantPath:    "/api/",
			expectError: false,
		},
		{
			name:        "simple URL",
			input:       "http://example.com",
			wantPath:    "/",
			expectError: false,
		},
		{
			name:        "invalid URL",
			input:       "://invalid",
			wantPath:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SetBaseURLUtil(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Path != tt.wantPath {
				t.Errorf("expected path %q, got %q", tt.wantPath, result.Path)
			}
		})
	}
}

func TestNewRequestFunction(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:9000/api/")

	tests := []struct {
		name     string
		method   string
		path     string
		opt      interface{}
		wantPath string
	}{
		{
			name:     "GET request without options",
			method:   "GET",
			path:     "projects/search",
			opt:      nil,
			wantPath: "/api/projects/search",
		},
		{
			name:   "GET request with options",
			method: "GET",
			path:   "projects/search",
			opt: struct {
				Query string `url:"q"`
			}{Query: "test"},
			wantPath: "/api/projects/search",
		},
		{
			name:   "POST request",
			method: "POST",
			path:   "projects/create",
			opt: struct {
				Name string `url:"name"`
			}{Name: "my-project"},
			wantPath: "/api/projects/create",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest(tt.method, tt.path, baseURL, "admin", "admin", tt.opt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if req.Method != tt.method {
				t.Errorf("expected method %q, got %q", tt.method, req.Method)
			}

			if req.URL.Path != tt.wantPath {
				t.Errorf("expected path %q, got %q", tt.wantPath, req.URL.Path)
			}

			// Check basic auth is set
			username, password, ok := req.BasicAuth()
			if !ok {
				t.Error("expected basic auth to be set")
			}
			if username != "admin" || password != "admin" {
				t.Errorf("expected admin:admin, got %s:%s", username, password)
			}

			// Check Accept header
			if req.Header.Get("Accept") != "application/json" {
				t.Errorf("expected Accept header 'application/json', got %q", req.Header.Get("Accept"))
			}

			// For POST/PUT, check Content-Type
			if tt.method == "POST" || tt.method == "PUT" {
				if req.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type 'application/json', got %q", req.Header.Get("Content-Type"))
				}
			}
		})
	}
}

func TestDoFunction(t *testing.T) {
	t.Run("successful JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "test", "value": 42}`))
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)
		var result map[string]interface{}
		err := Do(server.Client(), req, &result)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["name"] != "test" {
			t.Errorf("expected name 'test', got %v", result["name"])
		}
	})

	t.Run("error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "bad request"}`))
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)
		var result map[string]interface{}
		err := Do(server.Client(), req, &result)

		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("io.Writer response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("raw data"))
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)
		var buf bytes.Buffer
		err := Do(server.Client(), req, &buf)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if buf.String() != "raw data" {
			t.Errorf("expected 'raw data', got %q", buf.String())
		}
	})

	t.Run("nil response value", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL, nil)
		err := Do(server.Client(), req, nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCheckResponseFunction(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        string
		expectError bool
	}{
		{"200 OK", 200, "", false},
		{"201 Created", 201, "", false},
		{"202 Accepted", 202, "", false},
		{"204 No Content", 204, "", false},
		{"304 Not Modified", 304, "", false},
		{"400 Bad Request", 400, `{"error": "bad request"}`, true},
		{"401 Unauthorized", 401, `{"error": "unauthorized"}`, true},
		{"403 Forbidden", 403, `{"error": "forbidden"}`, true},
		{"404 Not Found", 404, `{"error": "not found"}`, true},
		{"500 Internal Server Error", 500, `{"error": "internal error"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.ReadCloser
			if tt.body != "" {
				body = io.NopCloser(bytes.NewBufferString(tt.body))
			} else {
				body = http.NoBody
			}

			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       body,
				Request:    &http.Request{URL: &url.URL{Path: "/test", Scheme: "http", Host: "localhost"}},
			}

			err := CheckResponse(resp)

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestResponseErrorError(t *testing.T) {
	errResp := &Error{
		Response: &http.Response{
			StatusCode: 404,
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Scheme: "http", Host: "localhost", Path: "/api/test"},
			},
		},
		Message: "not found",
	}

	errStr := errResp.Error()
	if errStr == "" {
		t.Error("expected non-empty error string")
	}

	// Check that it contains expected parts
	if !bytes.Contains([]byte(errStr), []byte("404")) {
		t.Error("error string should contain status code")
	}
	if !bytes.Contains([]byte(errStr), []byte("not found")) {
		t.Error("error string should contain message")
	}
}

func TestParseErrorFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		contains string
	}{
		{
			name:     "string error",
			input:    "simple error",
			contains: "simple error",
		},
		{
			name:     "array error",
			input:    []interface{}{"error1", "error2"},
			contains: "error1",
		},
		{
			name:     "map error",
			input:    map[string]interface{}{"key": "value"},
			contains: "key",
		},
		{
			name:     "unknown type",
			input:    123,
			contains: "failed to parse unexpected error type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseError(tt.input)
			if !bytes.Contains([]byte(result), []byte(tt.contains)) {
				t.Errorf("expected result to contain %q, got %q", tt.contains, result)
			}
		})
	}
}
