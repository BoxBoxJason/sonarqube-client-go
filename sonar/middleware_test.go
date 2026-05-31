package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordingTransport is an http.RoundTripper that records whether it was called.
//
//nolint:exhaustruct // next is set explicitly; called is intentionally zero-value false
type recordingTransport struct {
	next   http.RoundTripper
	called bool
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.called = true

	return r.next.RoundTrip(req)
}

// TestWithMiddleware_WrapsTransport verifies that middleware is invoked on each request.
func TestWithMiddleware_WrapsTransport(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, mockHandler(t, http.MethodGet, "/authentication/validate", http.StatusOK, nil))

	recorder := &recordingTransport{}

	mw := func(next http.RoundTripper) http.RoundTripper {
		recorder.next = next

		return recorder
	}

	client, err := NewClient(nil, WithBaseURL(ts.url()), WithMiddleware(mw))
	require.NoError(t, err)

	_, _, _ = client.Authentication.Validate(context.Background())

	assert.True(t, recorder.called, "middleware RoundTripper should have been called")
}

// TestWithMiddleware_OutermostFirst verifies that the first middleware provided is
// the outermost wrapper (first to execute on a request).
func TestWithMiddleware_OutermostFirst(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, mockHandler(t, http.MethodGet, "/authentication/validate", http.StatusOK, nil))

	var callOrder []int

	makeOrderedMiddleware := func(id int) Middleware {
		return func(next http.RoundTripper) http.RoundTripper {
			return roundTripFunc(func(req *http.Request) (*http.Response, error) {
				callOrder = append(callOrder, id)

				return next.RoundTrip(req)
			})
		}
	}

	client, err := NewClient(nil, WithBaseURL(ts.url()), WithMiddleware(makeOrderedMiddleware(1), makeOrderedMiddleware(2)))
	require.NoError(t, err)

	_, _, _ = client.Authentication.Validate(context.Background())

	assert.Equal(t, []int{1, 2}, callOrder, "middleware should execute outermost (index 0) first")
}

// TestWithMiddleware_WithHTTPClient verifies that middleware wraps the transport of
// a custom http.Client provided via WithHTTPClient.
func TestWithMiddleware_WithHTTPClient(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, mockHandler(t, http.MethodGet, "/authentication/validate", http.StatusOK, nil))

	innerRecorder := &recordingTransport{next: http.DefaultTransport}
	customClient := &http.Client{Transport: innerRecorder}

	outerRecorder := &recordingTransport{}

	mw := func(next http.RoundTripper) http.RoundTripper {
		outerRecorder.next = next

		return outerRecorder
	}

	client, err := NewClient(nil, WithBaseURL(ts.url()), WithHTTPClient(customClient), WithMiddleware(mw))
	require.NoError(t, err)

	_, _, _ = client.Authentication.Validate(context.Background())

	assert.True(t, outerRecorder.called, "outer middleware should have been called")
	assert.True(t, innerRecorder.called, "inner custom transport should have been called through middleware chain")
}

// TestWithMiddleware_NoMiddleware_DefaultUnchanged verifies that when no middleware is
// provided the client uses http.DefaultClient without wrapping it.
func TestWithMiddleware_NoMiddleware_DefaultUnchanged(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil)
	require.NoError(t, err)

	assert.Same(t, http.DefaultClient, client.httpClient, "httpClient should be http.DefaultClient when no middleware is provided")
}

// roundTripFunc is a convenience type that implements http.RoundTripper via a function.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
