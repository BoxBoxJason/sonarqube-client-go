package helpers

import (
	"context"
	"fmt"
	"time"
)

// WaitCondition is a function that returns true when the condition is met.
type WaitCondition func() (bool, error)

// WaitForCondition waits for a condition to be met with polling.
func WaitForCondition(condition WaitCondition, timeout, interval time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Try immediately
	done, err := condition()
	if err != nil {
		return err
	}

	if done {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for condition after %v", timeout)
		case <-ticker.C:
			done, err := condition()
			if err != nil {
				return err
			}

			if done {
				return nil
			}
		}
	}
}

// WaitForConditionWithDefault waits for a condition with default timeout and interval.
func WaitForConditionWithDefault(condition WaitCondition) error {
	return WaitForCondition(condition, DefaultTimeout, DefaultPollInterval)
}

// Retry executes a function with retries.
func Retry(fn func() error, maxRetries int, delay time.Duration) error {
	var attempts int

	for {
		err := fn()
		if err == nil {
			return nil
		}

		attempts++

		if attempts > maxRetries {
			return fmt.Errorf("failed after %d attempts: %w", attempts, err)
		}

		time.Sleep(delay)
	}
}
