package cli

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeResponse is a test response struct for PatternResponseBody.
type fakeResponse struct {
	Name string
}

// fakeService is a mock service for testing InvokeMethod.
type fakeService struct{}

// ResponseBodyMethod simulates a (*Response, *http.Response, error) method.
func (f *fakeService) ResponseBodyMethod(opt *struct{}) (*fakeResponse, *http.Response, error) {
	return &fakeResponse{Name: "ok"}, nil, nil
}

// NoBodyMethod simulates a (*http.Response, error) method.
func (f *fakeService) NoBodyMethod() (*http.Response, error) {
	return nil, nil
}

// RawBytesMethod simulates a ([]byte, *http.Response, error) method.
func (f *fakeService) RawBytesMethod(opt *struct{}) ([]byte, *http.Response, error) {
	return []byte("hello"), nil, nil
}

// RawStringMethod simulates a (*string, *http.Response, error) method.
func (f *fakeService) RawStringMethod(opt *struct{}) (*string, *http.Response, error) {
	s := "world"

	return &s, nil, nil
}

// SliceMethod simulates a ([]SomeStruct, *http.Response, error) method.
func (f *fakeService) SliceMethod(opt *struct{}) ([]fakeResponse, *http.Response, error) {
	return []fakeResponse{{Name: "a"}, {Name: "b"}}, nil, nil
}

// ErrorMethod simulates a method that returns an error.
func (f *fakeService) ErrorMethod(opt *struct{}) (*fakeResponse, *http.Response, error) {
	return nil, nil, errors.New("something failed")
}

// TestClassifyMethod tests method pattern classification.
func TestClassifyMethod(t *testing.T) {
	svcType := reflect.TypeOf(&fakeService{})

	tests := []struct {
		name   string
		method string
		want   MethodReturnPattern
	}{
		{name: "response body", method: "ResponseBodyMethod", want: PatternResponseBody},
		{name: "no body", method: "NoBodyMethod", want: PatternNoBody},
		{name: "raw bytes", method: "RawBytesMethod", want: PatternRawBytes},
		{name: "raw string", method: "RawStringMethod", want: PatternRawString},
		{name: "slice", method: "SliceMethod", want: PatternSlice},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			method, found := svcType.MethodByName(tc.method)
			require.True(t, found, "method %s not found", tc.method)

			got := ClassifyMethod(method)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestInvokeMethod_ResponseBody tests invocation of a response body method.
func TestInvokeMethod_ResponseBody(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.New(reflect.TypeOf(struct{}{}))

	result, _, err := InvokeMethod(svc, "ResponseBodyMethod", opt, PatternResponseBody, true)
	require.NoError(t, err)

	resp, ok := result.(*fakeResponse)
	require.True(t, ok)
	assert.Equal(t, "ok", resp.Name)
}

// TestInvokeMethod_NoBody tests invocation of a no-body method.
func TestInvokeMethod_NoBody(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.Value{}

	result, _, err := InvokeMethod(svc, "NoBodyMethod", opt, PatternNoBody, false)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// TestInvokeMethod_RawBytes tests invocation of a raw bytes method.
func TestInvokeMethod_RawBytes(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.New(reflect.TypeOf(struct{}{}))

	result, _, err := InvokeMethod(svc, "RawBytesMethod", opt, PatternRawBytes, true)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), result)
}

// TestInvokeMethod_RawString tests invocation of a raw string method.
func TestInvokeMethod_RawString(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.New(reflect.TypeOf(struct{}{}))

	result, _, err := InvokeMethod(svc, "RawStringMethod", opt, PatternRawString, true)
	require.NoError(t, err)

	s, ok := result.(*string)
	require.True(t, ok)
	assert.Equal(t, "world", *s)
}

// TestInvokeMethod_Slice tests invocation of a slice method.
func TestInvokeMethod_Slice(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.New(reflect.TypeOf(struct{}{}))

	result, _, err := InvokeMethod(svc, "SliceMethod", opt, PatternSlice, true)
	require.NoError(t, err)

	items, ok := result.([]fakeResponse)
	require.True(t, ok)
	assert.Len(t, items, 2)
	assert.Equal(t, "a", items[0].Name)
}

// TestInvokeMethod_Error tests that errors are propagated correctly.
func TestInvokeMethod_Error(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.New(reflect.TypeOf(struct{}{}))

	_, _, err := InvokeMethod(svc, "ErrorMethod", opt, PatternResponseBody, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "something failed")
}

// TestInvokeMethod_NotFound tests error for missing methods.
func TestInvokeMethod_NotFound(t *testing.T) {
	svc := reflect.ValueOf(&fakeService{})
	opt := reflect.Value{}

	_, _, err := InvokeMethod(svc, "DoesNotExist", opt, PatternNoBody, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestExtractError_Nil tests that nil error values return nil.
func TestExtractError_Nil(t *testing.T) {
	err := extractError(reflect.ValueOf((*error)(nil)))
	assert.NoError(t, err)
}

// TestExtractHTTPResponse_Nil tests that nil response values return nil.
func TestExtractHTTPResponse_Nil(t *testing.T) {
	resp := extractHTTPResponse(reflect.ValueOf((*http.Response)(nil)))
	assert.Nil(t, resp)
}

// TestExtractHTTPResponse_Valid tests extracting a valid response.
func TestExtractHTTPResponse_Valid(t *testing.T) {
	expected := &http.Response{StatusCode: http.StatusOK}
	resp := extractHTTPResponse(reflect.ValueOf(expected))
	assert.Equal(t, expected, resp)
}
