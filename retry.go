package workflow

import (
	"time"
)

type RetryOptions struct {
	MaxRetries  int
	MaxDelay    time.Duration
	ShouldRetry func(out any, err error) bool
}

func Retry(action Action, opts *RetryOptions) Action {
	if opts == nil {
		opts = &RetryOptions{
			MaxRetries:  3,
			MaxDelay:    time.Second * 30,
			ShouldRetry: nil,
		}
	}

	return func(in any) (any, error) {
		// calculate starting delay based on the max possible delay
		// max delay divided by max retries to the power of 2
		delay := opts.MaxDelay / (1 << opts.MaxRetries)

		// first loop is the initial try and does not count as a retry
		for retry := 0; retry <= opts.MaxRetries; retry++ {
			out, err := action(in)
			if err != nil && retry >= opts.MaxRetries {
				// already retried the maximum number of times, return error
				return out, err
			} else if err == nil {
				// action was successful, return immediately
				return out, nil
			} else if opts.ShouldRetry != nil && !opts.ShouldRetry(out, err) {
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
