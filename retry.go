package workflow

import (
	"math/rand"
	"time"
)

type RetryOptions struct {
	// Maximum number of retries before error is returned.
	MaxRetries int

	// Initial delay between the first failed attempt and the first retry. Will be adjusted according to the backoff strategy for future retries.
	InitialDelay time.Duration

	// Maximum possible delay between retries.
	MaxDelay time.Duration

	// Defines the maximum range of randomness to add to the delay between retries, i.e. jitter of 200 ms with a delay of 500 ms means the range for the actual delay is 300 ms to 700 ms. This helps prevent the thundering herd issue that can happen when a large number of concurrent transactions retry at the exact same time.
	Jitter time.Duration

	// Optional function to determine if a retry should happen. This function is always called, even if no error has occurred.
	ShouldRetry func(out any, err error) bool

	// BackoffStrategy // TODO: should I bother with this now?
}

func Retry(action Action, opts *RetryOptions) Action {
	if opts == nil {
		opts = &RetryOptions{}
	}

	// set default values
	if opts.MaxRetries <= 0 {
		// default to no retries, an error will return immediately
		opts.MaxRetries = 0
	}
	if opts.InitialDelay <= 0 {
		// default to no delay between retries
		opts.InitialDelay = 0
	}
	if opts.MaxDelay <= 0 {
		// default to no delay between retries
		opts.MaxDelay = 0
	}

	return func(in any) (any, error) {
		// calculate starting delay based on the max possible delay
		// max delay divided by max retries to the power of 2
		delay := opts.InitialDelay

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

			// delay before next retry
			time.Sleep(randDuration(delay-opts.Jitter, delay+opts.Jitter))

			// delay increases exponentially
			delay *= 2

			// respect max delay if set
			if opts.MaxDelay > 0 && delay > opts.MaxDelay {
				delay = opts.MaxDelay
			}
		}

		return nil, nil
	}
}

// Returns a random time in the closed range [min, max].
func randDuration(min time.Duration, max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max-min+1)) + int64(min))
}
