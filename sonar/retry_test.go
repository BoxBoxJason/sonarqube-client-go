package sonar

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// countingTransport records every call and returns pre-configured responses in order.
type countingTransport struct {
	responses []*http.Response
	errors    []error
	calls     int
}

func (c *countingTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	i := c.calls
	c.calls++

	if i < len(c.errors) && c.errors[i] != nil {
		return nil, c.errors[i]
	}

	if i < len(c.responses) {
		return c.responses[i], nil
	}

	return makeResponse(http.StatusOK), nil
}

func makeResponse(statusCode int) *http.Response {
	//nolint:exhaustruct // only status and body matter for tests
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}
}

func TestRetryRoundTripper_SingleAttemptWhenDisabled(t *testing.T) {
	transport := &countingTransport{responses: []*http.Response{makeResponse(http.StatusServiceUnavailable)}}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{MaxAttempts: 0, RetryableStatusCodes: []int{503}},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, 1, transport.calls)
}

func TestRetryRoundTripper_RetriesOnRetryableStatus(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponse(http.StatusServiceUnavailable),
			makeResponse(http.StatusServiceUnavailable),
			makeResponse(http.StatusOK),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          4,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{503},
		},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, transport.calls)
}

func TestRetryRoundTripper_NoRetryOnNonRetryableStatus(t *testing.T) {
	transport := &countingTransport{responses: []*http.Response{makeResponse(http.StatusBadRequest)}}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{MaxAttempts: 4, RetryableStatusCodes: []int{503}},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 1, transport.calls)
}

func TestRetryRoundTripper_ExhaustsAllAttempts(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponse(http.StatusBadGateway),
			makeResponse(http.StatusBadGateway),
			makeResponse(http.StatusBadGateway),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          3,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{502},
		},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, 3, transport.calls)
}

func TestRetryRoundTripper_ContextCancelledStopsRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	transport := &countingTransport{}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          4,
			InitialDelay:         time.Second,
			MaxDelay:             time.Second,
			RetryableStatusCodes: []int{503},
		},
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	// The base transport will also receive the cancelled context and may return an error.
	_, _ = rt.RoundTrip(req)
	// We just verify we don't spin endlessly; the exact error depends on the transport.
	assert.LessOrEqual(t, transport.calls, 1)
}

func TestRetryRoundTripper_NonReplayableBodyNotRetried(t *testing.T) {
	transport := &countingTransport{responses: []*http.Response{makeResponse(http.StatusServiceUnavailable)}}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{MaxAttempts: 4, RetryableStatusCodes: []int{503}},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://example.com", io.NopCloser(bytes.NewReader([]byte("body"))))
	// GetBody is nil because we used io.NopCloser directly, not bytes.NewReader via http.NewRequest.
	require.Nil(t, req.GetBody)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, 1, transport.calls)
}

func TestRetryRoundTripper_ReplayableBodyRetried(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponse(http.StatusServiceUnavailable),
			makeResponse(http.StatusOK),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          3,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{503},
		},
	}

	body := bytes.NewReader([]byte(`{"key":"value"}`))
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://example.com", body)
	// http.NewRequest with bytes.Reader sets GetBody automatically.
	require.NotNil(t, req.GetBody)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, transport.calls)
}

func TestWithRetry_ClientIntegration(t *testing.T) {
	callCount := 0
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		if callCount < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	client, err := NewClient(nil,
		WithBaseURL(ts.url()),
		WithRetry(RetryOptions{
			MaxAttempts:          4,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{503},
		}),
	)
	require.NoError(t, err)

	req, err := client.NewSonarQubeV1APIRequest(context.Background(), http.MethodGet, "api/system/ping", nil)
	require.NoError(t, err)

	resp, err := client.Do(req, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, callCount)
}

func TestComputeDelay_FullJitter(t *testing.T) {
	rt := &retryRoundTripper{
		base: http.DefaultTransport,
		opts: RetryOptions{
			InitialDelay: 100 * time.Millisecond,
			MaxDelay:     10 * time.Second,
		},
	}

	for attempt := range 5 {
		delay := rt.computeDelay(attempt)
		assert.GreaterOrEqual(t, delay, time.Duration(0))

		cap := time.Duration(float64(rt.opts.InitialDelay) * float64(int(1)<<attempt))
		if cap > rt.opts.MaxDelay {
			cap = rt.opts.MaxDelay
		}

		assert.LessOrEqual(t, delay, cap)
	}
}

func TestSleepContext_CancelledContextReturnsFalse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := sleepContext(ctx, time.Second)
	assert.False(t, result)
}

func TestSleepContext_ZeroDelayReturnsTrue(t *testing.T) {
	result := sleepContext(context.Background(), 0)
	assert.True(t, result)
}

func TestSleepContext_TimerReleasedOnContextCancel(t *testing.T) {
	// Cancelling the context while sleeping must unblock immediately and not leak
	// the timer (time.NewTimer + Stop rather than time.After).
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	result := sleepContext(ctx, time.Hour)
	assert.False(t, result)
	assert.Less(t, time.Since(start), time.Second)
}

func TestEvaluate_ContextErrorPreferredOverTransportError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	rt := &retryRoundTripper{
		base: http.DefaultTransport,
		opts: RetryOptions{MaxAttempts: 3, RetryableStatusCodes: []int{503}},
	}

	transportErr := fmt.Errorf("connection refused")
	done, result, err := rt.evaluate(ctx, nil, transportErr, false)

	assert.True(t, done)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestRetryRoundTripper_AttemptsHeaderSetAfterRetry(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponse(http.StatusServiceUnavailable),
			makeResponse(http.StatusServiceUnavailable),
			makeResponse(http.StatusOK),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          4,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{503},
		},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// 3 total attempts → header value "3".
	assert.Equal(t, "3", resp.Header.Get("X-Retry-Attempts"))
}

func TestRetryRoundTripper_AttemptsHeaderAbsentWithNoRetry(t *testing.T) {
	transport := &countingTransport{responses: []*http.Response{makeResponse(http.StatusOK)}}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{MaxAttempts: 4, RetryableStatusCodes: []int{503}},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Empty(t, resp.Header.Get("X-Retry-Attempts"))
}

func TestWithRetry_SliceIsolatedFromCaller(t *testing.T) {
	codes := []int{503}
	client, err := NewClient(nil, WithRetry(RetryOptions{
		MaxAttempts:          2,
		InitialDelay:         time.Millisecond,
		MaxDelay:             5 * time.Millisecond,
		RetryableStatusCodes: codes,
	}))
	require.NoError(t, err)

	// Mutate the original slice after client creation.
	codes[0] = 200

	// The client's retry transport should still use 503, not 200.
	transport := client.httpClient.Transport.(*retryRoundTripper)
	assert.Equal(t, []int{503}, transport.opts.RetryableStatusCodes)
}
