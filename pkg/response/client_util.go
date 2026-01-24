package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/google/go-querystring/query"
)

// SetBaseURLUtil sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func SetBaseURLUtil(urlStr string) (*url.URL, error) {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Update the base URL of the client.
	return baseURL, nil
}

// NewRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func NewRequest(method, path string, baseURL *url.URL, username, password string, opt any) (*http.Request, error) {
	// Set the encoded opaque data
	parsedURL := *baseURL

	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape path: %w", err)
	}

	parsedURL.RawPath = parsedURL.Path + path
	parsedURL.Path += unescaped

	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to encode query parameters: %w", err)
		}

		parsedURL.RawQuery = q.Encode()
	}

	//nolint:exhaustruct
	req := &http.Request{
		Method:     method,
		URL:        &parsedURL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       parsedURL.Host,
	}

	if method == http.MethodPost || method == http.MethodPut {
		// SonarQube use RawQuery even when method is POST
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		bodyReader := bytes.NewReader(bodyBytes)

		parsedURL.RawQuery = ""
		req.Body = io.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(username, password)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func Do(c *http.Client, req *http.Request, dest any) error {
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	err = CheckResponse(resp)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if dest != nil {
		if w, ok := dest.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(dest)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// Error holds the error response from the server.
type Error struct {
	Response *http.Response
	Message  string
	Body     []byte
}

func (e *Error) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)

	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
func CheckResponse(resp *http.Response) error {
	if (resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) || resp.StatusCode == http.StatusNotModified {
		return nil
	}

	errorResponse := &Error{Response: resp, Body: nil, Message: ""}

	data, err := io.ReadAll(resp.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw any

		err := json.Unmarshal(data, &raw)
		if err != nil {
			errorResponse.Message = string(data)
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw any) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []any:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}

		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]any:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}

		sort.Strings(errs)

		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}
