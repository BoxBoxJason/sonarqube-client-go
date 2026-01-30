package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebservicesService_List(t *testing.T) {
	listJSON := `{
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
	}`

	t.Run("success", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/webservices/list", http.StatusOK, listJSON)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		result, resp, err := client.Webservices.List(nil)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
		require.Len(t, result.Webservices, 1)
		assert.Equal(t, "api/issues", result.Webservices[0].Path)
		require.Len(t, result.Webservices[0].Actions, 1)
		assert.Equal(t, "search", result.Webservices[0].Actions[0].Key)
		assert.Len(t, result.Webservices[0].Actions[0].Params, 1)
	})

	t.Run("with include_internals", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/webservices/list", http.StatusOK, `{"webServices": []}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Webservices.List(&WebservicesListOption{
			IncludeInternals: true,
		})

		require.NoError(t, err)
	})

	t.Run("empty option", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/webservices/list", http.StatusOK, `{"webServices": []}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		_, _, err := client.Webservices.List(&WebservicesListOption{})

		require.NoError(t, err)
	})
}

func TestWebservicesService_ResponseExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := mockHandler(t, http.MethodGet, "/webservices/response_example", http.StatusOK, `{
			"issues": [
				{
					"key": "AVfN9MxQTN6qjVMfZpW-",
					"rule": "squid:S2259"
				}
			]
		}`)
		server := newTestServer(t, handler)
		client := newTestClient(t, server.URL)

		result, resp, err := client.Webservices.ResponseExample(&WebservicesResponseExampleOption{
			Action:     "search",
			Controller: "api/issues",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, result)
	})

	t.Run("nil option fails validation", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Webservices.ResponseExample(nil)

		assert.Error(t, err)
	})

	t.Run("missing action fails validation", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Webservices.ResponseExample(&WebservicesResponseExampleOption{
			Controller: "api/issues",
		})

		assert.Error(t, err)
	})

	t.Run("missing controller fails validation", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.Webservices.ResponseExample(&WebservicesResponseExampleOption{
			Action: "search",
		})

		assert.Error(t, err)
	})
}

func TestWebservicesService_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWebservicesService_ValidateResponseExampleOpt(t *testing.T) {
	client := newLocalhostClient(t)

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
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
