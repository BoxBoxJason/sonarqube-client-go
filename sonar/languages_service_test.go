package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLanguages_List(t *testing.T) {
	response := `{
		"languages": [
			{"key": "java", "name": "Java"},
			{"key": "go", "name": "Go"},
			{"key": "py", "name": "Python"}
		]
	}`
	handler := mockHandler(t, http.MethodGet, "/languages/list", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Languages.List(&LanguagesListOption{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Languages, 3)
	assert.Equal(t, "java", result.Languages[0].Key)
	assert.Equal(t, "Java", result.Languages[0].Name)
}

func TestLanguages_List_WithQuery(t *testing.T) {
	response := `{"languages": [{"key": "java", "name": "Java"}]}`
	handler := mockHandler(t, http.MethodGet, "/languages/list", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Languages.List(&LanguagesListOption{Query: "java"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Languages, 1)
}

func TestLanguages_List_NilOption(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/languages/list", http.StatusOK, `{"languages": []}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Languages.List(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestLanguages_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *LanguagesListOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &LanguagesListOption{}, false},
		{"with query", &LanguagesListOption{Query: "java", PageSize: 25}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Languages.ValidateListOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
