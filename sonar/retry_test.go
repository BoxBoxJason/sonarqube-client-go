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

	// Use an idempotent method so the non-replayable body is the reason the
	// request is not retried, not the method.
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "http://example.com", io.NopCloser(bytes.NewReader([]byte("body"))))
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
	// PUT is idempotent and has a replayable body, so it is retried by default.
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "http://example.com", body)
	// http.NewRequest with bytes.Reader sets GetBody automatically.
	require.NotNil(t, req.GetBody)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, transport.calls)
}

func TestRetryRoundTripper_NonIdempotentNotRetriedByDefault(t *testing.T) {
	t.Parallel()

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
	require.NotNil(t, req.GetBody)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	// POST is not retried by default even with a replayable body and a retryable status.
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, 1, transport.calls)
}

func TestRetryRoundTripper_NonIdempotentRetriedWhenOptedIn(t *testing.T) {
	t.Parallel()

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
			RetryNonIdempotent:   true,
		},
	}

	body := bytes.NewReader([]byte(`{"key":"value"}`))
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://example.com", body)
	require.NotNil(t, req.GetBody)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, transport.calls)
}

func TestIsIdempotentMethod(t *testing.T) {
	t.Parallel()

	idempotent := []string{
		http.MethodGet, http.MethodHead, http.MethodPut,
		http.MethodDelete, http.MethodOptions, http.MethodTrace,
	}
	for _, method := range idempotent {
		assert.True(t, isIdempotentMethod(method), "%s should be idempotent", method)
	}

	for _, method := range []string{http.MethodPost, http.MethodPatch} {
		assert.False(t, isIdempotentMethod(method), "%s should not be idempotent", method)
	}
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

// =============================================
// Retry-After tests (#201)
// =============================================

func TestRetryRoundTripper_RetryAfterSeconds(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponseWithHeader(http.StatusTooManyRequests, "Retry-After", "0"),
			makeResponse(http.StatusOK),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          3,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{429},
		},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, transport.calls)
}

func TestRetryDelay_UsesRetryAfterHTTPDateOnRateLimitResponse(t *testing.T) {
	// Test retryDelay directly to avoid waiting for the actual delay.
	// Use 30s in the future — well beyond http.TimeFormat's 1-second precision so
	// the formatted date is always still in the future when time.Until is called.
	rt := &retryRoundTripper{
		base: http.DefaultTransport,
		opts: RetryOptions{InitialDelay: time.Second, MaxDelay: time.Minute},
	}

	future := time.Now().Add(30 * time.Second).UTC()
	resp := makeResponseWithHeader(http.StatusTooManyRequests, "Retry-After", future.Format(http.TimeFormat))

	delay := rt.retryDelay(resp, 0)
	assert.Greater(t, delay, 25*time.Second)
	assert.LessOrEqual(t, delay, 31*time.Second)
}

func TestRetryRoundTripper_RetryAfterAbsent_FallsBackToBackoff(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponse(http.StatusTooManyRequests),
			makeResponse(http.StatusOK),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          3,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{429},
		},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, transport.calls)
}

func TestRetryRoundTripper_RetryDisabled_429ReturnedAsIs(t *testing.T) {
	transport := &countingTransport{
		responses: []*http.Response{
			makeResponseWithHeader(http.StatusTooManyRequests, "Retry-After", "10"),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{MaxAttempts: 1, RetryableStatusCodes: []int{429}},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", http.NoBody)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, 1, transport.calls)
}

func TestRetryDelay_UsesRetryAfterOnRateLimitResponse(t *testing.T) {
	rt := &retryRoundTripper{
		base: http.DefaultTransport,
		opts: RetryOptions{InitialDelay: time.Second, MaxDelay: time.Minute},
	}

	resp := makeResponseWithHeader(http.StatusTooManyRequests, "Retry-After", "5")
	delay := rt.retryDelay(resp, 0)

	assert.Equal(t, 5*time.Second, delay)
}

func TestRetryDelay_ZeroRetryAfterOverridesBackoff(t *testing.T) {
	rt := &retryRoundTripper{
		base: http.DefaultTransport,
		opts: RetryOptions{InitialDelay: 10 * time.Second, MaxDelay: time.Minute},
	}

	// Retry-After: 0 must result in immediate retry, not fall back to the backoff.
	resp := makeResponseWithHeader(http.StatusTooManyRequests, "Retry-After", "0")
	delay := rt.retryDelay(resp, 0)

	assert.Equal(t, time.Duration(0), delay)
}

func TestRetryDelay_IgnoresRetryAfterOnNon429(t *testing.T) {
	rt := &retryRoundTripper{
		base: http.DefaultTransport,
		opts: RetryOptions{InitialDelay: 10 * time.Millisecond, MaxDelay: time.Second},
	}

	resp := makeResponseWithHeader(http.StatusServiceUnavailable, "Retry-After", "5")
	delay := rt.retryDelay(resp, 0)

	// Should use exponential backoff, not the Retry-After value (5s).
	assert.Less(t, delay, 5*time.Second)
}

func TestParseRetryAfterHeader_Seconds(t *testing.T) {
	tests := []struct {
		header      string
		expected    time.Duration
		expectValid bool
	}{
		{"10", 10 * time.Second, true},
		{"0", 0, true},   // zero is a valid instruction to retry immediately
		{"-1", 0, false}, // negative integers are invalid
		{"3600", 3600 * time.Second, true},
	}

	for _, tc := range tests {
		t.Run(tc.header, func(t *testing.T) {
			d, ok := parseRetryAfterHeader(tc.header)
			assert.Equal(t, tc.expected, d)
			assert.Equal(t, tc.expectValid, ok)
		})
	}
}

func TestParseRetryAfterHeader_HTTPDate(t *testing.T) {
	// Use 30s in the future to survive http.TimeFormat's 1-second precision.
	future := time.Now().Add(30 * time.Second).UTC()
	header := future.Format(http.TimeFormat)

	d, ok := parseRetryAfterHeader(header)
	assert.True(t, ok)
	assert.Greater(t, d, 25*time.Second)
	assert.LessOrEqual(t, d, 31*time.Second)
}

func TestParseRetryAfterHeader_PastDate(t *testing.T) {
	// A past HTTP-date is valid: the server's intended wait has already elapsed,
	// so the result is zero delay (retry immediately) rather than falling back to backoff.
	past := time.Now().Add(-time.Minute).UTC()
	header := past.Format(http.TimeFormat)

	d, ok := parseRetryAfterHeader(header)
	assert.True(t, ok)
	assert.Equal(t, time.Duration(0), d)
}

func TestParseRetryAfterHeader_Empty(t *testing.T) {
	_, ok := parseRetryAfterHeader("")
	assert.False(t, ok)
}

func TestParseRetryAfterHeader_Invalid(t *testing.T) {
	_, ok := parseRetryAfterHeader("not-a-date-or-number")
	assert.False(t, ok)
}

// makeResponseWithHeader creates a test response with a single header set.
func makeResponseWithHeader(statusCode int, key, value string) *http.Response {
	resp := makeResponse(statusCode)
	resp.Header.Set(key, value)

	return resp
}

func TestWithRetry_429ContextCancellationDuringWait(t *testing.T) {
	// The Retry-After header requests a 60s wait; the context times out in 100ms.
	// This verifies that sleepContext honours context cancellation and does not block.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	transport := &countingTransport{
		responses: []*http.Response{
			makeResponseWithHeader(http.StatusTooManyRequests, "Retry-After", "60"),
		},
	}
	rt := &retryRoundTripper{
		base: transport,
		opts: RetryOptions{
			MaxAttempts:          3,
			InitialDelay:         time.Millisecond,
			MaxDelay:             5 * time.Millisecond,
			RetryableStatusCodes: []int{429},
		},
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)

	start := time.Now()
	_, err := rt.RoundTrip(req)
	elapsed := time.Since(start)

	// Must finish well under 60s (the Retry-After value).
	assert.Less(t, elapsed, time.Second)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
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
