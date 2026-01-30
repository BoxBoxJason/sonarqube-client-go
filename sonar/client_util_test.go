package sonargo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo_JSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"foo":"bar"}`))
	}))
	t.Cleanup(ts.Close)

	baseURL, _ := url.Parse(ts.URL + "/")
	req, err := NewRequest(http.MethodGet, "test", baseURL, "u", "p", nil)
	require.NoError(t, err)

	var v map[string]any
	resp, err := Do(http.DefaultClient, req, &v)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "bar", v["foo"])
}

func TestDo_NoBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(ts.Close)

	baseURL, _ := url.Parse(ts.URL + "/")
	req, err := NewRequest(http.MethodPost, "test", baseURL, "u", "p", nil)
	require.NoError(t, err)

	resp, err := Do(http.DefaultClient, req, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestNewRequest_WithQueryParams(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost/api/")
	opt := struct {
		Query string `url:"q,omitempty"`
		Page  int    `url:"page,omitempty"`
	}{
		Query: "test",
		Page:  1,
	}

	req, err := NewRequest(http.MethodGet, "search", baseURL, "u", "p", opt)
	require.NoError(t, err)
	assert.Contains(t, req.URL.RawQuery, "q=test")
	assert.Contains(t, req.URL.RawQuery, "page=1")
}

func TestNewRequest_BasicAuth(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost/api/")
	req, err := NewRequest(http.MethodGet, "test", baseURL, "user", "pass", nil)
	require.NoError(t, err)

	username, password, ok := req.BasicAuth()
	require.True(t, ok)
	assert.Equal(t, "user", username)
	assert.Equal(t, "pass", password)
}
