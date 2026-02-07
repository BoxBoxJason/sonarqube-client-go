package sonar

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalysisCache_Clear(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/analysis_cache/clear", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	resp, err := client.AnalysisCache.Clear(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAnalysisCache_Clear_WithOptions(t *testing.T) {
	expectedParams := map[string]string{
		"project": "my-project",
		"branch":  "feature",
	}
	handler := mockEmptyHandlerWithParams(t, http.MethodPost, "/analysis_cache/clear", http.StatusNoContent, expectedParams)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AnalysisCacheClearOption{
		Project: "my-project",
		Branch:  "feature",
	}

	resp, err := client.AnalysisCache.Clear(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAnalysisCache_Get(t *testing.T) {
	mockData := []byte("mock cached data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/analysis_cache/get", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"), "expected project query param")

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	t.Cleanup(ts.Close)

	client := newTestClient(t, ts.URL+"/")

	opt := &AnalysisCacheGetOption{
		Project: "my-project",
	}

	resp, err := client.AnalysisCache.Get(opt)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(body, mockData))
}

func TestAnalysisCache_Get_WithOptions(t *testing.T) {
	mockData := []byte("mock cached data with options")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/analysis_cache/get", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"), "expected project query param")
		assert.Equal(t, "main", r.URL.Query().Get("branch"), "expected branch query param")

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	t.Cleanup(ts.Close)

	client := newTestClient(t, ts.URL+"/")

	opt := &AnalysisCacheGetOption{
		Project: "my-project",
		Branch:  "main",
	}

	resp, err := client.AnalysisCache.Get(opt)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(body, mockData))
}

func TestAnalysisCache_ValidateClearOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *AnalysisCacheClearOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &AnalysisCacheClearOption{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.AnalysisCache.ValidateClearOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAnalysisCache_ValidateGetOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *AnalysisCacheGetOption
		wantErr bool
	}{
		{"nil option", nil, true},
		{"empty option", &AnalysisCacheGetOption{}, true},
		{"with Project", &AnalysisCacheGetOption{Project: "my-project"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.AnalysisCache.ValidateGetOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
