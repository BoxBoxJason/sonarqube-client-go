package sonar

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"slices"
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
		c.retryOptions = &opts

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
func (r *retryRoundTripper) retryLoop(req *http.Request, hasBody bool) (*http.Response, error) {
	for attempt := range r.opts.MaxAttempts {
		resp, err := r.doAttempt(req, hasBody)
		isLast := attempt >= r.opts.MaxAttempts-1

		if done, result, resultErr := r.evaluate(req.Context(), resp, err, isLast); done {
			return result, resultErr
		}

		if !sleepContext(req.Context(), r.computeDelay(attempt)) {
			return nil, req.Context().Err() //nolint:wrapcheck // context error is the direct cause
		}
	}

	return nil, fmt.Errorf("retry: exhausted %d attempts", r.opts.MaxAttempts)
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
func (r *retryRoundTripper) evaluate(ctx context.Context, resp *http.Response, err error, isLast bool) (bool, *http.Response, error) {
	if err != nil {
		if ctx.Err() != nil || isLast {
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
// or becomes done while waiting.
func sleepContext(ctx context.Context, delay time.Duration) bool {
	if ctx.Err() != nil {
		return false
	}

	if delay <= 0 {
		return true
	}

	select {
	case <-ctx.Done():
		return false
	case <-time.After(delay):
		return true
	}
}

// drainAndClose discards resp's body and closes it so the connection can be reused.
func drainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}
