package sonar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

// jsonStructToQueryValues marshals a struct using its json tags and converts the
// resulting flat JSON object into url.Values suitable for use as URL query
// parameters. Fields tagged with json:"-" are excluded. Fields with
// omitempty that hold zero values are omitted. Arrays of primitives, numbers
// and booleans are supported. Nested objects are rejected with an error.
func jsonStructToQueryValues(v any) (url.Values, error) {
	if v == nil {
		return url.Values{}, nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to encode query values: %w", err)
	}

	var decodedMap map[string]any

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	err = dec.Decode(&decodedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode query values: %w", err)
	}

	return mapToQueryValues(decodedMap)
}

// mapToQueryValues converts a flat JSON-decoded map into url.Values.
// Returns an error if any value is a nested object.
func mapToQueryValues(decodedMap map[string]any) (url.Values, error) {
	vals := url.Values{}

	for key, val := range decodedMap {
		switch typedVal := val.(type) {
		case json.Number:
			vals.Set(key, typedVal.String())
		case string:
			vals.Set(key, typedVal)
		case bool:
			vals.Set(key, strconv.FormatBool(typedVal))
		case []any:
			for _, item := range typedVal {
				vals.Add(key, fmt.Sprint(item))
			}
		case nil:
			// Skip null values
		case map[string]any:
			return nil, fmt.Errorf("nested objects are not supported as query parameters (field %q)", key)
		}
	}

	return vals, nil
}
