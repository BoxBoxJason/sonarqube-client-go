package sonar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
//nolint:wrapcheck // error context is clear from call site
func Do(httpClient *http.Client, req *http.Request, dest any) (*http.Response, error) {
	isText := false

	if _, ok := dest.(*string); ok {
		req.Header.Set("Accept", "text/plain")

		isText = true
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if dest == nil {
		return resp, nil
	}

	if writer, ok := dest.(io.Writer); ok {
		_, err = io.Copy(writer, resp.Body)

		return resp, err
	}

	if isText {
		data, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return resp, readErr
		}

		if strPtr, ok := dest.(*string); ok {
			*strPtr = string(data)
		}

		return resp, nil
	}

	err = json.NewDecoder(resp.Body).Decode(dest)

	return resp, err
}

// ResponseError represents an error response from the SonarQube API.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type ResponseError struct {
	Body     []byte
	Response *http.Response
	Message  string
}

// Error returns the error message.
func (e *ResponseError) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	urlStr := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)

	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, urlStr, e.Response.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors.
//
//nolint:exhaustruct // Body and Message are set conditionally
func CheckResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent, http.StatusNotModified:
		return nil
	}

	errorResponse := &ResponseError{
		Response: resp,
	}

	data, err := io.ReadAll(resp.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw any

		unmarshalErr := json.Unmarshal(data, &raw)
		if unmarshalErr != nil {
			errorResponse.Message = string(data)
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw any) string {
	switch rawTyped := raw.(type) {
	case string:
		return rawTyped

	case []any:
		var errs []string

		for _, v := range rawTyped {
			errs = append(errs, parseError(v))
		}

		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]any:
		var errs []string

		for k, v := range rawTyped {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}

		sort.Strings(errs)

		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

// assignIfNotNil assigns the value pointed to by src to the value pointed to by
// dest if src is not nil.
func assignIfNotNil[T any](dest *T, src *T) {
	if src != nil {
		*dest = *src
	}
}

func assignPtrIfNotNil[T any](dest **T, src *T) {
	if src != nil {
		*dest = src
	}
}
