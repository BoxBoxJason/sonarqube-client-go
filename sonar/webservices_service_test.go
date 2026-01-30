package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebservicesService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"webServices": [
					{
						"path": "api/issues",
						"description": "Issues management",
						"since": "3.6",
						"actions": [
							{
								"key": "search",
								"description": "Search for issues",
								"since": "3.6",
								"post": false,
								"internal": false,
								"hasResponseExample": true,
								"params": [
									{
										"key": "severities",
										"description": "Comma-separated list of severities",
										"required": false,
										"possibleValues": ["INFO", "MINOR", "MAJOR", "CRITICAL", "BLOCKER"]
									}
								]
							}
						]
					}
				]
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Webservices.List(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if len(result.Webservices) != 1 {
			t.Errorf("expected 1 webservice, got %d", len(result.Webservices))
		}
		if result.Webservices[0].Path != "api/issues" {
			t.Errorf("expected path 'api/issues', got '%s'", result.Webservices[0].Path)
		}
		if len(result.Webservices[0].Actions) != 1 {
			t.Errorf("expected 1 action, got %d", len(result.Webservices[0].Actions))
		}
		if result.Webservices[0].Actions[0].Key != "search" {
			t.Errorf("expected key 'search', got '%s'", result.Webservices[0].Actions[0].Key)
		}
		if len(result.Webservices[0].Actions[0].Params) != 1 {
			t.Errorf("expected 1 param, got %d", len(result.Webservices[0].Actions[0].Params))
		}
	})

	t.Run("with include_internals", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("include_internals") != "true" {
				t.Errorf("expected include_internals 'true', got '%s'", r.URL.Query().Get("include_internals"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"webServices": []}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Webservices.List(&WebservicesListOption{
			IncludeInternals: true,
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty option", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"webServices": []}`))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		_, _, err := client.Webservices.List(&WebservicesListOption{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestWebservicesService_ResponseExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Query().Get("action") != "search" {
				t.Errorf("expected action 'search', got '%s'", r.URL.Query().Get("action"))
			}
			if r.URL.Query().Get("controller") != "api/issues" {
				t.Errorf("expected controller 'api/issues', got '%s'", r.URL.Query().Get("controller"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"issues": [
					{
						"key": "AVfN9MxQTN6qjVMfZpW-",
						"rule": "squid:S2259"
					}
				]
			}`))
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client, _ := NewClient(server.URL+"/api/", "user", "pass")

		result, resp, err := client.Webservices.ResponseExample(&WebservicesResponseExampleOption{
			Action:     "search",
			Controller: "api/issues",
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
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Webservices.ResponseExample(nil)
		if err == nil {
			t.Error("expected error for nil option")
		}
	})

	t.Run("missing action", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Webservices.ResponseExample(&WebservicesResponseExampleOption{
			Controller: "api/issues",
		})
		if err == nil {
			t.Error("expected error for missing action")
		}
	})

	t.Run("missing controller", func(t *testing.T) {
		client, _ := NewClient("http://localhost/api/", "user", "pass")

		_, _, err := client.Webservices.ResponseExample(&WebservicesResponseExampleOption{
			Action: "search",
		})
		if err == nil {
			t.Error("expected error for missing controller")
		}
	})
}

func TestWebservicesService_ValidateListOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *WebservicesListOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &WebservicesListOption{}, false},
		{"with include internals", &WebservicesListOption{IncludeInternals: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Webservices.ValidateListOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateListOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebservicesService_ValidateResponseExampleOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	tests := []struct {
		name    string
		opt     *WebservicesResponseExampleOption
		wantErr bool
	}{
		{"valid", &WebservicesResponseExampleOption{Action: "search", Controller: "api/issues"}, false},
		{"nil option", nil, true},
		{"missing action", &WebservicesResponseExampleOption{Controller: "api/issues"}, true},
		{"missing controller", &WebservicesResponseExampleOption{Action: "search"}, true},
		{"empty both", &WebservicesResponseExampleOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Webservices.ValidateResponseExampleOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResponseExampleOpt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
