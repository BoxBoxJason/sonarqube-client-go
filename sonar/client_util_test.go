package sonar

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
// Use (*Client).NewSonarQubeV1APIRequest or (*Client).NewSonarQubeV2APIRequest for production code.
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

func TestDo_DecodeError_WrapsRequestContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not valid json`))
	}))
	t.Cleanup(ts.Close)

	baseURL, _ := url.Parse(ts.URL + "/")
	req, err := newRequest(http.MethodGet, "projects/search", baseURL, "u", "p", nil)
	require.NoError(t, err)

	var v map[string]any
	_, err = Do(http.DefaultClient, req, &v)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode", "decode error should be wrapped with context")
	assert.Contains(t, err.Error(), http.MethodGet, "wrapped error should include the method")
	assert.Contains(t, err.Error(), "projects/search", "wrapped error should include the endpoint path")
}

// drainTrackingBody records whether it was fully read (to EOF) and closed.
type drainTrackingBody struct {
	r       io.Reader
	readEOF bool
	closed  bool
}

func (b *drainTrackingBody) Read(p []byte) (int, error) {
	n, err := b.r.Read(p)
	if err == io.EOF {
		b.readEOF = true
	}

	return n, err //nolint:wrapcheck // test stub mirrors the underlying reader
}

func (b *drainTrackingBody) Close() error {
	b.closed = true

	return nil
}

// fixedRoundTripper returns a fixed response regardless of the request.
type fixedRoundTripper struct{ resp *http.Response }

func (t fixedRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return t.resp, nil
}

func TestDo_DrainsBodyOnSuccess(t *testing.T) {
	// A JSON value followed by trailing bytes the decoder will not consume, so the
	// drain (not the decode) is what reaches EOF.
	payload := append([]byte(`{"foo":"bar"}`), bytes.Repeat([]byte(" "), 4096)...)
	body := &drainTrackingBody{r: bytes.NewReader(payload)} //nolint:exhaustruct

	resp := &http.Response{ //nolint:exhaustruct // only fields needed for the test
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       body,
	}
	client := &http.Client{Transport: fixedRoundTripper{resp: resp}} //nolint:exhaustruct

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)

	var v map[string]any
	_, err := Do(client, req, &v)
	require.NoError(t, err)
	assert.Equal(t, "bar", v["foo"])
	assert.True(t, body.readEOF, "body should be drained to EOF before close")
	assert.True(t, body.closed, "body should be closed")
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

func TestCheckResponse_PopulatesStatusCode(t *testing.T) {
	for _, code := range []int{http.StatusNotFound, http.StatusUnauthorized, http.StatusForbidden, http.StatusConflict, http.StatusTooManyRequests, http.StatusInternalServerError} {
		t.Run(fmt.Sprintf("%d", code), func(t *testing.T) {
			resp := &http.Response{
				StatusCode: code,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(nil)),
				Request:    httptest.NewRequest(http.MethodGet, "/api/test", nil),
			}

			err := CheckResponse(resp)
			require.Error(t, err)

			var re *ResponseError
			require.ErrorAs(t, err, &re)
			assert.Equal(t, code, re.StatusCode)
		})
	}
}

func TestCheckResponse_PreservesJSONErrorBody(t *testing.T) {
	body := []byte(`"api error"`)
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    httptest.NewRequest(http.MethodGet, "/api/test", nil),
	}

	err := CheckResponse(resp)
	require.Error(t, err)

	var re *ResponseError
	require.ErrorAs(t, err, &re)
	assert.Equal(t, body, re.Body)
	assert.Equal(t, "api error", re.Message)
}

func TestCheckResponse_PreservesTextErrorBody(t *testing.T) {
	body := []byte("plain error")
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    httptest.NewRequest(http.MethodGet, "/api/test", nil),
	}

	err := CheckResponse(resp)
	require.Error(t, err)

	var re *ResponseError
	require.ErrorAs(t, err, &re)
	assert.Equal(t, body, re.Body)
	assert.Equal(t, string(body), re.Message)
}

func TestCheckResponse_TruncatesLargeErrorBody(t *testing.T) {
	body := bytes.Repeat([]byte("x"), maxErrorResponseBodyBytes+1)
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    httptest.NewRequest(http.MethodGet, "/api/test", nil),
	}

	err := CheckResponse(resp)
	require.Error(t, err)

	var re *ResponseError
	require.ErrorAs(t, err, &re)
	assert.Len(t, re.Body, maxErrorResponseBodyBytes)
	assert.Equal(t, body[:maxErrorResponseBodyBytes], re.Body)
	assert.Contains(t, re.Message, "(response body truncated after 1048576 bytes)")
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(makeAPIError(http.StatusNotFound)))
	assert.False(t, IsNotFound(makeAPIError(http.StatusUnauthorized)))
	assert.False(t, IsNotFound(nil))
	assert.False(t, IsNotFound(fmt.Errorf("plain error")))
}

func TestIsUnauthorized(t *testing.T) {
	assert.True(t, IsUnauthorized(makeAPIError(http.StatusUnauthorized)))
	assert.False(t, IsUnauthorized(makeAPIError(http.StatusNotFound)))
}

func TestIsForbidden(t *testing.T) {
	assert.True(t, IsForbidden(makeAPIError(http.StatusForbidden)))
	assert.False(t, IsForbidden(makeAPIError(http.StatusUnauthorized)))
}

func TestIsConflict(t *testing.T) {
	assert.True(t, IsConflict(makeAPIError(http.StatusConflict)))
	assert.False(t, IsConflict(makeAPIError(http.StatusNotFound)))
}

func TestIsRateLimited(t *testing.T) {
	assert.True(t, IsRateLimited(makeAPIError(http.StatusTooManyRequests)))
	assert.False(t, IsRateLimited(makeAPIError(http.StatusForbidden)))
}

func TestIsServerError(t *testing.T) {
	assert.True(t, IsServerError(makeAPIError(http.StatusInternalServerError)))
	assert.True(t, IsServerError(makeAPIError(http.StatusBadGateway)))
	assert.True(t, IsServerError(makeAPIError(599)))
	assert.False(t, IsServerError(makeAPIError(http.StatusNotFound)))
	assert.False(t, IsServerError(makeAPIError(http.StatusTooManyRequests)))
}

func TestSentinels_WorkThroughWrappingChain(t *testing.T) {
	inner := makeAPIError(http.StatusNotFound)
	wrapped := fmt.Errorf("operation failed: %w", inner)

	assert.True(t, IsNotFound(wrapped))
	assert.False(t, IsUnauthorized(wrapped))
}

// makeAPIError creates a *ResponseError with the given status code for testing.
func makeAPIError(statusCode int) *ResponseError {
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	//nolint:exhaustruct // only fields needed for test populated
	return &ResponseError{
		Response: &http.Response{
			StatusCode: statusCode,
			Request:    req,
		},
		StatusCode: statusCode,
		Message:    fmt.Sprintf("status %d", statusCode),
	}
}

func TestResponseError_Error_NilResponse(t *testing.T) {
	// Must not panic when Response or Request is nil.
	re := &ResponseError{StatusCode: http.StatusNotFound, Message: "not found"}
	assert.Equal(t, "404 not found", re.Error())
}

func TestResponseError_Error_NilRequest(t *testing.T) {
	//nolint:exhaustruct // only fields needed for test populated
	re := &ResponseError{
		Response:   &http.Response{StatusCode: http.StatusInternalServerError},
		StatusCode: http.StatusInternalServerError,
		Message:    "internal error",
	}
	assert.Equal(t, "500 internal error", re.Error())
}

func TestResponseError_Error_UsesStatusCodeField(t *testing.T) {
	// StatusCode field is canonical; Response.StatusCode is not used in Error().
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	//nolint:exhaustruct // only fields needed for test populated
	re := &ResponseError{
		Response:   &http.Response{StatusCode: 999, Request: req},
		StatusCode: http.StatusNotFound,
		Message:    "not found",
	}

	assert.Contains(t, re.Error(), "404")
	assert.NotContains(t, re.Error(), "999")
}

func TestResponseError_Error_PathUnescape(t *testing.T) {
	// Encoded path segments must be unescaped correctly in the error string.
	req := httptest.NewRequest(http.MethodGet, "/api/my%20resource", nil)

	//nolint:exhaustruct // only fields needed for test populated
	re := &ResponseError{
		Response:   &http.Response{StatusCode: http.StatusNotFound, Request: req},
		StatusCode: http.StatusNotFound,
		Message:    "not found",
	}

	assert.Contains(t, re.Error(), "/api/my resource")
}

func TestJsonStructToQueryValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    url.Values
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  url.Values{},
		},
		{
			name: "flat struct with primitives",
			input: struct {
				Name    string `json:"name"`
				Count   int    `json:"count"`
				Enabled bool   `json:"enabled"`
			}{
				Name:    "test",
				Count:   42,
				Enabled: true,
			},
			want: url.Values{
				FieldName: []string{"test"},
				"count":   []string{"42"},
				"enabled": []string{"true"},
			},
		},
		{
			name: "omitempty zero value omitted",
			input: struct {
				Name  string `json:"name,omitempty"`
				Count int    `json:"count,omitempty"`
			}{
				Name: "only-name",
			},
			want: url.Values{
				FieldName: []string{"only-name"},
			},
		},
		{
			name: "string slice",
			input: struct {
				Tags []string `json:"tags"`
			}{
				Tags: []string{"a", "b"},
			},
			want: url.Values{
				"tags": []string{"a", "b"},
			},
		},
		{
			name: "nested object rejected",
			input: struct {
				Nested map[string]string `json:"nested"`
			}{
				Nested: map[string]string{"key": "val"},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := jsonStructToQueryValues(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
