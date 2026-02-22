package cli

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// MethodReturnPattern describes the return signature of a service method.
type MethodReturnPattern int

const (
	// PatternResponseBody is for methods returning (*TypedResponse, *http.Response, error).
	PatternResponseBody MethodReturnPattern = iota
	// PatternNoBody is for methods returning (*http.Response, error).
	PatternNoBody
	// PatternRawBytes is for methods returning ([]byte, *http.Response, error).
	PatternRawBytes
	// PatternRawString is for methods returning (*string, *http.Response, error).
	PatternRawString
	// PatternSlice is for methods returning ([]SomeStruct, *http.Response, error).
	PatternSlice
)

const (
	// expectedTripleReturn is the number of return values for most service methods.
	expectedTripleReturn = 3
	// expectedDoubleReturn is the number of return values for no-body methods.
	expectedDoubleReturn = 2
)

// ClassifyMethod determines the return pattern of a service method from its reflect.Type.
func ClassifyMethod(method reflect.Method) MethodReturnPattern {
	methodType := method.Type
	numOut := methodType.NumOut()

	if numOut == expectedDoubleReturn {
		return PatternNoBody
	}

	if numOut == expectedTripleReturn {
		first := methodType.Out(0)

		// []byte
		if first.Kind() == reflect.Slice && first.Elem().Kind() == reflect.Uint8 {
			return PatternRawBytes
		}

		// []SomeStruct (not []byte)
		if first.Kind() == reflect.Slice {
			return PatternSlice
		}

		// *string
		if first.Kind() == reflect.Ptr && first.Elem().Kind() == reflect.String {
			return PatternRawString
		}

		return PatternResponseBody
	}

	return PatternNoBody
}

// InvokeMethod calls a service method via reflection and returns the result.
// It handles all four return patterns, extracting the response value and error.
func InvokeMethod(service reflect.Value, methodName string, opt reflect.Value, pattern MethodReturnPattern, hasOpt bool) (any, *http.Response, error) {
	method := service.MethodByName(methodName)
	if !method.IsValid() {
		return nil, nil, fmt.Errorf("method %q not found on service", methodName)
	}

	var results []reflect.Value

	if hasOpt {
		results = method.Call([]reflect.Value{opt})
	} else {
		results = method.Call(nil)
	}

	switch pattern {
	case PatternResponseBody:
		return extractTripleReturn(results)
	case PatternNoBody:
		return extractDoubleReturn(results)
	case PatternRawBytes:
		return extractRawBytesReturn(results)
	case PatternRawString:
		return extractRawStringReturn(results)
	case PatternSlice:
		return extractSliceReturn(results)
	default:
		return extractDoubleReturn(results)
	}
}

// InvokeStreamingMethod calls a streaming service method (like Push.SonarlintEvents)
// and pipes the response body to the writer. The response body is left open for streaming.
func InvokeStreamingMethod(service reflect.Value, methodName string, opt reflect.Value, writer io.Writer) error {
	method := service.MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("method %q not found on service", methodName)
	}

	results := method.Call([]reflect.Value{opt})

	// Push.SonarlintEvents returns (*http.Response, error)
	if len(results) != expectedDoubleReturn {
		return fmt.Errorf("unexpected return count %d from streaming method", len(results))
	}

	if !results[1].IsNil() {
		//nolint:forcetypeassert // second return is always error
		return results[1].Interface().(error)
	}

	//nolint:forcetypeassert // first return is always *http.Response
	resp := results[0].Interface().(*http.Response)
	if resp == nil || resp.Body == nil {
		return nil
	}

	defer func() { _ = resp.Body.Close() }()

	_, err := io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("streaming copy failed: %w", err)
	}

	return nil
}

// extractTripleReturn extracts (*TypedResponse, *http.Response, error) from reflect results.
func extractTripleReturn(results []reflect.Value) (any, *http.Response, error) {
	err := extractError(results[2])
	if err != nil {
		return nil, extractHTTPResponse(results[1]), err
	}

	var val any

	if !results[0].IsNil() {
		val = results[0].Interface()
	}

	return val, extractHTTPResponse(results[1]), nil
}

// extractDoubleReturn extracts (*http.Response, error) from reflect results.
func extractDoubleReturn(results []reflect.Value) (any, *http.Response, error) {
	err := extractError(results[1])

	return nil, extractHTTPResponse(results[0]), err
}

// extractRawBytesReturn extracts ([]byte, *http.Response, error) from reflect results.
func extractRawBytesReturn(results []reflect.Value) (any, *http.Response, error) {
	err := extractError(results[2])
	if err != nil {
		return nil, extractHTTPResponse(results[1]), err
	}

	if results[0].IsNil() {
		return nil, extractHTTPResponse(results[1]), nil
	}

	return results[0].Bytes(), extractHTTPResponse(results[1]), nil
}

// extractRawStringReturn extracts (*string, *http.Response, error) from reflect results.
func extractRawStringReturn(results []reflect.Value) (any, *http.Response, error) {
	err := extractError(results[2])
	if err != nil {
		return nil, extractHTTPResponse(results[1]), err
	}

	if results[0].IsNil() {
		return nil, extractHTTPResponse(results[1]), nil
	}

	return results[0].Interface(), extractHTTPResponse(results[1]), nil
}

// extractSliceReturn extracts ([]SomeStruct, *http.Response, error) from reflect results.
func extractSliceReturn(results []reflect.Value) (any, *http.Response, error) {
	err := extractError(results[2])
	if err != nil {
		return nil, extractHTTPResponse(results[1]), err
	}

	if results[0].IsNil() {
		return nil, extractHTTPResponse(results[1]), nil
	}

	return results[0].Interface(), extractHTTPResponse(results[1]), nil
}

// extractError safely extracts an error from a reflect.Value.
func extractError(val reflect.Value) error {
	if val.IsNil() {
		return nil
	}

	err, ok := val.Interface().(error)
	if !ok {
		return nil
	}

	return err
}

// extractHTTPResponse safely extracts an *http.Response from a reflect.Value.
func extractHTTPResponse(val reflect.Value) *http.Response {
	if val.IsNil() {
		return nil
	}

	resp, ok := val.Interface().(*http.Response)
	if !ok {
		return nil
	}

	return resp
}
