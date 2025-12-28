package response

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		username    string
		password    string
		expectError bool
	}{
		{
			name:        "valid endpoint",
			endpoint:    "http://localhost:9000/api",
			username:    "admin",
			password:    "admin",
			expectError: false,
		},
		{
			name:        "empty endpoint uses default",
			endpoint:    "",
			username:    "admin",
			password:    "admin",
			expectError: false,
		},
		{
			name:        "endpoint with trailing slash",
			endpoint:    "http://localhost:9000/api/",
			username:    "admin",
			password:    "admin",
			expectError: false,
		},
		{
			name:        "invalid endpoint",
			endpoint:    "://invalid",
			username:    "admin",
			password:    "admin",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.endpoint, tt.username, tt.password)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if client == nil {
				t.Error("Expected client, got nil")
				return
			}
			if client.Webservices == nil {
				t.Error("Expected Webservices to be initialized")
			}
		})
	}
}

func TestSetBaseURLUtil(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectPath  string
	}{
		{
			name:        "without trailing slash",
			input:       "http://localhost:9000/api",
			expectError: false,
			expectPath:  "/api/",
		},
		{
			name:        "with trailing slash",
			input:       "http://localhost:9000/api/",
			expectError: false,
			expectPath:  "/api/",
		},
		{
			name:        "root path",
			input:       "http://localhost:9000",
			expectError: false,
			expectPath:  "/",
		},
		{
			name:        "invalid URL",
			input:       "://invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := SetBaseURLUtil(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if url.Path != tt.expectPath {
				t.Errorf("Expected path %q, got %q", tt.expectPath, url.Path)
			}
		})
	}
}

func TestNewExampleFetcher(t *testing.T) {
	fetcher := NewExampleFetcher("http://localhost:9000/api", "admin", "password")
	if fetcher == nil {
		t.Fatal("Expected fetcher, got nil")
	}
	if fetcher.endpoint != "http://localhost:9000/api" {
		t.Errorf("Expected endpoint 'http://localhost:9000/api', got '%s'", fetcher.endpoint)
	}
	if fetcher.username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", fetcher.username)
	}
	if fetcher.password != "password" {
		t.Errorf("Expected password 'password', got '%s'", fetcher.password)
	}
}

func TestWebservicesService(t *testing.T) {
	// Create a test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Setup mock response
	mux.HandleFunc("/webservices/response_example", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"format": "json", "example": "{\"key\": \"value\"}"}`))
	})

	// Create client pointing to test server
	client, err := NewClient(server.URL, "admin", "admin")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test ResponseExample
	opt := &WebservicesResponseExampleOption{
		Action:     "search",
		Controller: "api/issues",
	}
	resp, err := client.Webservices.ResponseExample(opt)
	if err != nil {
		t.Fatalf("ResponseExample failed: %v", err)
	}
	if resp.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", resp.Format)
	}
	if resp.Name != "search" {
		t.Errorf("Expected name 'search', got '%s'", resp.Name)
	}
}

func TestCheckResponse(t *testing.T) {
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
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       http.NoBody,
				Request:    &http.Request{URL: &url.URL{Path: "/test"}},
			}
			err := CheckResponse(resp)
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
