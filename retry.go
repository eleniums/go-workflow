package workflow

import (
	"math/rand"
	"time"
)

// Options for configuring a retry action.
type RetryOptions struct {
	// Maximum number of retries before error is returned. Can set to 0 for no retries.
	MaxRetries int

	// Initial delay between the first failed attempt and the first retry. Will be adjusted according to the backoff strategy for future retries. Can set to 0 for no delay.
	InitialDelay time.Duration

	// Maximum possible delay between retries. Can set to 0 for no maximum delay.
	MaxDelay time.Duration

	// Defines the maximum range of randomness to add to the delay between retries, i.e. jitter of 200 ms with a delay of 500 ms means the range for the actual delay is 300 ms to 700 ms. This helps prevent the thundering herd issue that can happen when a large number of concurrent transactions retry at the exact same time. Can set to 0 for no jitter.
	Jitter time.Duration

	// Optional function to determine if a retry should happen. This function is always called, even if no error has occurred.
	ShouldRetry func(out any, err error) bool

	// Optional function to determine backoff strategy. Can set to nil for no backoff.
	BackoffStrategy func(delay time.Duration) time.Duration
}

// Retry an action if it returns an error.
func Retry(action Action, opts *RetryOptions) Action {
	if opts == nil {
		// set some defaults if no options provided
		opts = &RetryOptions{
			MaxRetries:      3,
			InitialDelay:    time.Millisecond * 200,
			MaxDelay:        time.Second * 30,
			Jitter:          time.Millisecond * 50,
			BackoffStrategy: BackoffStrategyExponential(),
		}
	}

	return func(in any) (any, error) {
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

			// increase delay as required by the backoff strategy
			if opts.BackoffStrategy != nil {
				delay = opts.BackoffStrategy(delay)
			}

			// respect max delay if set
			if opts.MaxDelay > 0 && delay > opts.MaxDelay {
				delay = opts.MaxDelay
			}
		}

		return nil, nil
	}
}

// Backoff strategy that does nothing. The delay is consistent between retries.
func BackoffStrategyNone() func(delay time.Duration) time.Duration {
	return func(delay time.Duration) time.Duration {
		return delay
	}
}

// Backoff strategy that increases delay by a predefined amount after each retry.
func BackoffStrategyLinear(increment time.Duration) func(delay time.Duration) time.Duration {
	return func(delay time.Duration) time.Duration {
		delay += increment
		return delay
	}
}

// Backoff strategy that increases delay exponentially after each retry.
func BackoffStrategyExponential() func(delay time.Duration) time.Duration {
	return func(delay time.Duration) time.Duration {
		delay *= 2
		return delay
	}
}

// Returns a random time in the closed range [min, max].
func randDuration(min time.Duration, max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max-min+1)) + int64(min))
}
