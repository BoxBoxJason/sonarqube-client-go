package sonar

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"time"
)

const (
	retryBackoffBase = 2
	// maxRetryAfterSeconds is the largest delay-seconds value a time.Duration
	// can hold. A larger one is clamped to it rather than overflowing into a
	// negative delay, which would retry immediately.
	maxRetryAfterSeconds = int64(math.MaxInt64 / time.Second)
)

// RetryOptions configures opt-in retry with exponential backoff and full jitter.
// The zero value disables retrying (default behaviour).
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type RetryOptions struct {
	// MaxAttempts is the total number of attempts including the first.
	// Values of 0 or 1 disable retries.
	MaxAttempts int
	// InitialDelay is the base delay used in the exponential backoff calculation.
	InitialDelay time.Duration
	// MaxDelay caps the computed backoff delay. It does not apply to a delay
	// requested by a Retry-After header: see MaxRetryAfter.
	MaxDelay time.Duration
	// MaxRetryAfter is the longest Retry-After wait that is honoured on a 429
	// response. A response asking for more is returned to the caller at once
	// instead of being retried, as retrying sooner would ignore the server and
	// waiting would hold the request for that long. Zero sets no limit.
	MaxRetryAfter time.Duration
	// RetryableStatusCodes lists HTTP status codes that should trigger a retry.
	RetryableStatusCodes []int
	// RetryableError reports whether a transport error should trigger a retry.
	// Returning false fails the request with that error at once, for example
	// to avoid repeating a timeout that took long to expire. When nil, every
	// transport error is retried. It is not called once the request context
	// is done, as that always stops the retries.
	RetryableError func(err error) bool
	// RetryNonIdempotent enables automatic retries for non-idempotent methods
	// (POST, PATCH). It is disabled by default: only idempotent methods (GET,
	// HEAD, PUT, DELETE, OPTIONS, TRACE) are retried, because resending a
	// non-idempotent request after a network error or 5xx can duplicate a
	// server-side side effect (for example creating a resource twice). Enable
	// this only for endpoints you know are safe to repeat.
	RetryNonIdempotent bool
}

// WithRetry is a ClientOptionFunc that enables opt-in retry with exponential
// backoff and full jitter around the HTTP transport. Retry is disabled by default.
func WithRetry(opts RetryOptions) ClientOptionFunc {
	return func(c *Client) error {
		copied := opts
		copied.RetryableStatusCodes = slices.Clone(opts.RetryableStatusCodes)
		c.retryOptions = &copied

		return nil
	}
}

// retryRoundTripper wraps a base http.RoundTripper with configurable retry logic.
type retryRoundTripper struct {
	base http.RoundTripper
	opts RetryOptions
}

// RoundTrip executes the request, retrying on configured status codes or
// transport errors with exponential backoff and full jitter. Retries stop
// immediately when the request context is cancelled or its deadline is exceeded.
//
// Only idempotent methods are retried by default; non-idempotent methods (POST,
// PATCH) are passed straight through unless RetryNonIdempotent is set, so a
// resend cannot accidentally duplicate a server-side side effect.
//
// When retries occur, the final response carries an X-Retry-Attempts header with
// the total number of attempts made.
//
//nolint:wrapcheck // pass-through transport: errors from base RoundTripper are intentionally not wrapped
func (r *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.opts.MaxAttempts <= 1 {
		return r.base.RoundTrip(req)
	}

	// Non-idempotent methods are not retried unless explicitly opted in, to avoid
	// duplicating server-side side effects on a resend.
	if !r.opts.RetryNonIdempotent && !isIdempotentMethod(req.Method) {
		return r.base.RoundTrip(req)
	}

	hasBody := req.Body != nil && req.Body != http.NoBody
	if hasBody && req.GetBody == nil {
		// Non-replayable body: cannot retry safely.
		return r.base.RoundTrip(req)
	}

	return r.retryLoop(req, hasBody)
}

// isIdempotentMethod reports whether an HTTP method is safe to retry
// automatically. Per RFC 7231, GET, HEAD, PUT, DELETE, OPTIONS and TRACE are
// idempotent; POST and PATCH are not.
func isIdempotentMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut,
		http.MethodDelete, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// retryLoop runs up to MaxAttempts, sleeping between retries.
// On the final response, X-Retry-Attempts is set to the total attempt count
// when more than one attempt was made.
func (r *retryRoundTripper) retryLoop(req *http.Request, hasBody bool) (*http.Response, error) {
	for attempt := range r.opts.MaxAttempts {
		resp, err := r.doAttempt(req, hasBody)
		isLast := attempt >= r.opts.MaxAttempts-1

		if done, result, resultErr := r.evaluate(req.Context(), resp, err, isLast); done {
			if result != nil && attempt > 0 {
				result.Header.Set("X-Retry-Attempts", strconv.Itoa(attempt+1))
			}

			return result, resultErr
		}

		// resp headers are still accessible after evaluate drains the body.
		if !sleepContext(req.Context(), r.retryDelay(resp, attempt)) {
			return nil, req.Context().Err() //nolint:wrapcheck // context error is the direct cause
		}
	}

	return nil, fmt.Errorf("retry: exhausted %d attempts", r.opts.MaxAttempts)
}

// retryDelay returns the duration to wait before the next attempt.
// For 429 responses with a valid Retry-After header, that value takes precedence
// over the computed exponential backoff. A Retry-After of zero is honoured as an
// explicit instruction to retry immediately rather than falling back to backoff.
func (r *retryRoundTripper) retryDelay(resp *http.Response, attempt int) time.Duration {
	if delay, ok := requestedRetryDelay(resp); ok {
		return delay
	}

	return r.computeDelay(attempt)
}

// requestedRetryDelay returns the wait a 429 response asks for in its
// Retry-After header, if it carries a valid one.
func requestedRetryDelay(resp *http.Response) (time.Duration, bool) {
	if resp == nil || resp.StatusCode != http.StatusTooManyRequests {
		return 0, false
	}

	return parseRetryAfterHeader(resp.Header.Get("Retry-After"))
}

// doAttempt clones the request (replaying the body when present) and executes it.
func (r *retryRoundTripper) doAttempt(req *http.Request, hasBody bool) (*http.Response, error) {
	clonedReq := req.Clone(req.Context())

	if hasBody {
		var err error

		clonedReq.Body, err = req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("retry: failed to replay request body: %w", err)
		}
	}

	return r.base.RoundTrip(clonedReq) //nolint:wrapcheck
}

// evaluate decides whether the loop should stop after an attempt.
// Returns (true, resp, err) to stop, (false, nil, nil) to sleep and continue.
//
// When the context is cancelled, the context error is returned rather than the
// transport error so callers can reliably detect cancellation/timeouts.
func (r *retryRoundTripper) evaluate(ctx context.Context, resp *http.Response, err error, isLast bool) (bool, *http.Response, error) {
	if err != nil {
		ctxErr := ctx.Err()
		if ctxErr != nil {
			return true, nil, ctxErr //nolint:wrapcheck // context error is the direct cause
		}

		if isLast || !r.isRetryableError(err) {
			return true, nil, err
		}

		return false, nil, nil
	}

	if !r.isRetryable(resp.StatusCode) || isLast || r.exceedsMaxRetryAfter(resp) {
		return true, resp, nil
	}

	drainAndClose(resp)

	return false, nil, nil
}

// isRetryable reports whether statusCode is in the configured retry list.
func (r *retryRoundTripper) isRetryable(statusCode int) bool {
	return slices.Contains(r.opts.RetryableStatusCodes, statusCode)
}

// isRetryableError reports whether a transport error should trigger a retry.
func (r *retryRoundTripper) isRetryableError(err error) bool {
	return r.opts.RetryableError == nil || r.opts.RetryableError(err)
}

// exceedsMaxRetryAfter reports whether resp asks for a longer wait than
// MaxRetryAfter allows.
func (r *retryRoundTripper) exceedsMaxRetryAfter(resp *http.Response) bool {
	if r.opts.MaxRetryAfter <= 0 {
		return false
	}

	delay, ok := requestedRetryDelay(resp)

	return ok && delay > r.opts.MaxRetryAfter
}

// computeDelay returns the backoff duration for the given attempt index using full
// jitter: a random value in [0, min(InitialDelay * 2^attempt, MaxDelay)).
func (r *retryRoundTripper) computeDelay(attempt int) time.Duration {
	maxDelay := float64(r.opts.MaxDelay)
	base := float64(r.opts.InitialDelay) * math.Pow(retryBackoffBase, float64(attempt))

	if base > maxDelay {
		base = maxDelay
	}

	return time.Duration(rand.Float64() * base) //nolint:gosec // math/rand/v2 sufficient for jitter
}

// sleepContext waits for delay, returning false immediately if ctx is already done
// or becomes done while waiting. It uses time.NewTimer rather than time.After so
// the timer is stopped and its resources released as soon as the context fires.
func sleepContext(ctx context.Context, delay time.Duration) bool {
	if ctx.Err() != nil {
		return false
	}

	if delay <= 0 {
		return true
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// parseRetryAfterHeader parses a Retry-After header value. It supports both the
// delay-seconds form ("30") and the HTTP-date form ("Wed, 21 Oct 2025 07:28:00 GMT").
//
// Returns (duration, true) when the header is present and parseable. A value of
// zero is valid and means "retry immediately" (e.g. Retry-After: 0 or a past
// HTTP-date). Returns (0, false) when the header is absent, malformed, or carries
// a negative integer.
func parseRetryAfterHeader(header string) (time.Duration, bool) {
	if header == "" {
		return 0, false
	}

	// delay-seconds form: non-negative integer.
	secs, secsErr := strconv.Atoi(header)
	if secsErr == nil {
		if secs >= 0 {
			return time.Duration(min(int64(secs), maxRetryAfterSeconds)) * time.Second, true
		}

		return 0, false
	}

	// HTTP-date form. A past-dated value is treated as an instruction to retry
	// immediately (zero delay) rather than falling back to backoff.
	t, dateErr := http.ParseTime(header)
	if dateErr == nil {
		return max(time.Until(t), 0), true
	}

	return 0, false
}

// drainAndClose discards resp's body and closes it so the connection can be reused.
func drainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}
