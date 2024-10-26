package workflow

import (
	"errors"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func Test_Unit_Action_Retry(t *testing.T) {
	// arrange
	action := Do(func(in int) (int, error) {
		return in + 1, nil
	})
	actionErr := errors.New("test error")
	actionWithNumErrs := func(numErrs int) Action {
		return func(in any) (any, error) {
			if numErrs <= 0 {
				return 3, nil
			}
			numErrs--
			return 5, actionErr
		}
	}

	testCases := []struct {
		name     string
		action   Action
		opts     *RetryOptions
		in       any
		expected any
		err      error
	}{
		{
			name:   "action success",
			action: action,
			opts: &RetryOptions{
				ShouldRetry: nil,
				MaxRetries:  0,
				MaxDelay:    0,
			},
			in:       1,
			expected: 2,
			err:      nil,
		},
		{
			name:   "action error with no retries",
			action: actionWithNumErrs(1),
			opts: &RetryOptions{
				ShouldRetry: nil,
				MaxRetries:  0,
				MaxDelay:    time.Millisecond * 10,
			},
			in:       1,
			expected: 5,
			err:      actionErr,
		},
		{
			name:   "action error with no delay",
			action: actionWithNumErrs(4),
			opts: &RetryOptions{
				ShouldRetry: nil,
				MaxRetries:  3,
				MaxDelay:    0,
			},
			in:       1,
			expected: 5,
			err:      actionErr,
		},
		{
			name:   "action error retries and then succeeds",
			action: actionWithNumErrs(3),
			opts: &RetryOptions{
				ShouldRetry: nil,
				MaxRetries:  3,
				MaxDelay:    time.Millisecond * 10,
			},
			in:       1,
			expected: 3,
			err:      nil,
		},
		{
			name:   "action error should not retry",
			action: actionWithNumErrs(3),
			opts: &RetryOptions{
				ShouldRetry: func(out any, err error) bool {
					return false
				},
				MaxRetries: 3,
				MaxDelay:   time.Millisecond * 10,
			},
			in:       1,
			expected: 5,
			err:      actionErr,
		},
		{
			name:   "action error should retry",
			action: actionWithNumErrs(3),
			opts: &RetryOptions{
				ShouldRetry: func(out any, err error) bool {
					return true
				},
				MaxRetries: 3,
				MaxDelay:   time.Millisecond * 10,
			},
			in:       1,
			expected: 3,
			err:      nil,
		},
		{
			name:     "default options",
			action:   actionWithNumErrs(3),
			opts:     nil,
			in:       1,
			expected: 3,
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			action := Retry(tc.action, tc.opts)
			out, err := action(tc.in)

			// assert w
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}
