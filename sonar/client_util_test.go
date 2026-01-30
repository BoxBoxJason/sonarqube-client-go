package sonargo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newRequest is an internal helper that creates an API request with custom semantics.
// For POST/PUT requests, it marshals opt into a JSON body. For GET requests, it encodes
// opt as URL query parameters. This is kept for internal testing purposes only.
// Use (*Client).NewRequest for production code.
//
//nolint:unused // Used in test files only
func newRequest(method, path string, baseURL *url.URL, username, password string, opt any) (*http.Request, error) {
	baseURLCopy := *baseURL

	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape path: %w", err)
	}

	baseURLCopy.RawPath = baseURLCopy.Path + path
	baseURLCopy.Path += unescaped

	if opt != nil {
		queryValues, queryErr := query.Values(opt)
		if queryErr != nil {
			return nil, fmt.Errorf("failed to encode query values: %w", queryErr)
		}

		baseURLCopy.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequestWithContext(context.Background(), method, baseURLCopy.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if method == http.MethodPost || method == http.MethodPut {
		// SonarQube uses RawQuery even when method is POST
		bodyBytes, marshalErr := json.Marshal(opt)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", marshalErr)
		}

		bodyReader := bytes.NewReader(bodyBytes)
		baseURLCopy.RawQuery = ""
		req.Body = io.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(username, password)

	return req, nil
}

func TestDo_JSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"foo":"bar"}`))
	}))
	t.Cleanup(ts.Close)

	baseURL, _ := url.Parse(ts.URL + "/")
	req, err := newRequest(http.MethodGet, "test", baseURL, "u", "p", nil)
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
	req, err := newRequest(http.MethodPost, "test", baseURL, "u", "p", nil)
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

	req, err := newRequest(http.MethodGet, "search", baseURL, "u", "p", opt)
	require.NoError(t, err)
	assert.Contains(t, req.URL.RawQuery, "q=test")
	assert.Contains(t, req.URL.RawQuery, "page=1")
}

func TestNewRequest_BasicAuth(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost/api/")
	req, err := newRequest(http.MethodGet, "test", baseURL, "user", "pass", nil)
	require.NoError(t, err)

	username, password, ok := req.BasicAuth()
	require.True(t, ok)
	assert.Equal(t, "user", username)
	assert.Equal(t, "pass", password)
}
