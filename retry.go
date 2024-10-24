package workflow

import (
	"time"
)

func Retry(action Action, shouldRetry func(out any, err error) bool, maxRetries int, maxDelay time.Duration) Action {
	return func(in any) (any, error) {
		// calculate starting delay based on the max possible delay
		// max delay divided by max retries to the power of 2
		delay := maxDelay / (1 << maxRetries)

		// first loop is the initial try and does not count as a retry
		for retry := 0; retry <= maxRetries; retry++ {
			out, err := action(in)
			if err != nil && retry >= maxRetries {
				// already retried the maximum number of times, return error
				return out, err
			} else if err == nil {
				// action was successful, return immediately
				return out, nil
			} else if shouldRetry != nil && !shouldRetry(out, err) {
				// determined no retry should be allowed
				return out, err
			}

			// TODO: add random jitter to delay
			// TODO: should jitter be configurable?

			// delay before next retry
			time.Sleep(delay)

			// delay increases exponentially
			delay *= 2
		}

		return nil, nil
	}
}
