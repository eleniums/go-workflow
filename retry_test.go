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
		return func(a any) (any, error) {
			if numErrs <= 0 {
				return 3, nil
			}
			numErrs--
			return 5, actionErr
		}
	}

	testCases := []struct {
		name        string
		action      Action
		shouldRetry func(out any, err error) bool
		maxRetries  int
		maxDelay    time.Duration
		in          any
		expected    any
		err         error
	}{
		{
			name:        "action success",
			action:      action,
			shouldRetry: nil,
			maxRetries:  0,
			maxDelay:    0,
			in:          1,
			expected:    2,
			err:         nil,
		},
		{
			name:        "action error with no retries",
			action:      actionWithNumErrs(1),
			shouldRetry: nil,
			maxRetries:  0,
			maxDelay:    time.Millisecond * 10,
			in:          1,
			expected:    5,
			err:         actionErr,
		},
		{
			name:        "action error with no delay",
			action:      actionWithNumErrs(4),
			shouldRetry: nil,
			maxRetries:  3,
			maxDelay:    0,
			in:          1,
			expected:    5,
			err:         actionErr,
		},
		{
			name:        "action error retries and then succeeds",
			action:      actionWithNumErrs(3),
			shouldRetry: nil,
			maxRetries:  3,
			maxDelay:    time.Millisecond * 10,
			in:          1,
			expected:    3,
			err:         nil,
		},
		{
			name:   "action error should not retry",
			action: actionWithNumErrs(3),
			shouldRetry: func(out any, err error) bool {
				return false
			},
			maxRetries: 3,
			maxDelay:   time.Millisecond * 10,
			in:         1,
			expected:   5,
			err:        actionErr,
		},
		{
			name:   "action error should retry",
			action: actionWithNumErrs(3),
			shouldRetry: func(out any, err error) bool {
				return true
			},
			maxRetries: 3,
			maxDelay:   time.Millisecond * 10,
			in:         1,
			expected:   3,
			err:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			action := Retry(tc.action, tc.shouldRetry, tc.maxRetries, tc.maxDelay)
			out, err := action(tc.in)

			// assert w
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}
