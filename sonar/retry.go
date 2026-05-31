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

const retryBackoffBase = 2

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
	// MaxDelay caps the computed backoff delay.
	MaxDelay time.Duration
	// RetryableStatusCodes lists HTTP status codes that should trigger a retry.
	RetryableStatusCodes []int
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

// RoundTrip executes the request, retrying on configured status codes or network
// errors with exponential backoff and full jitter. Retries stop immediately when
// the request context is cancelled or its deadline is exceeded.
//
// When retries occur, the final response carries an X-Retry-Attempts header with
// the total number of attempts made.
//
//nolint:wrapcheck // pass-through transport: errors from base RoundTripper are intentionally not wrapped
func (r *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.opts.MaxAttempts <= 1 {
		return r.base.RoundTrip(req)
	}

	hasBody := req.Body != nil && req.Body != http.NoBody
	if hasBody && req.GetBody == nil {
		// Non-replayable body: cannot retry safely.
		return r.base.RoundTrip(req)
	}

	return r.retryLoop(req, hasBody)
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
	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		if d, ok := parseRetryAfterHeader(resp.Header.Get("Retry-After")); ok {
			return d
		}
	}

	return r.computeDelay(attempt)
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

		if isLast {
			return true, nil, err
		}

		return false, nil, nil
	}

	if !r.isRetryable(resp.StatusCode) || isLast {
		return true, resp, nil
	}

	drainAndClose(resp)

	return false, nil, nil
}

// isRetryable reports whether statusCode is in the configured retry list.
func (r *retryRoundTripper) isRetryable(statusCode int) bool {
	return slices.Contains(r.opts.RetryableStatusCodes, statusCode)
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
			return time.Duration(secs) * time.Second, true
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
