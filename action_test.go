package workflow

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func Test_Unit_Action_NoOp(t *testing.T) {
	// act
	action := NoOp()
	out, err := action(1)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, out)
}

func Test_Unit_Action_Wrap(t *testing.T) {
	// arrange
	action1 := Do(func(in int) (int, error) {
		return in + 1, nil
	})
	action2 := Do(func(in int) (int, error) {
		return in + 2, nil
	})

	actionErr := errors.New("test error")
	actionWithErr := Do(func(in int) (int, error) {
		return 5, actionErr
	})

	testCases := []struct {
		name     string
		action   Action
		next     Action
		in       any
		expected any
		err      error
	}{
		{
			name:     "action and next",
			action:   action1,
			next:     action2,
			in:       1,
			expected: 4,
			err:      nil,
		},
		{
			name:     "no actions",
			action:   nil,
			next:     nil,
			in:       1,
			expected: nil,
			err:      nil,
		},
		{
			name:     "action only",
			action:   action1,
			next:     nil,
			in:       1,
			expected: 2,
			err:      nil,
		},
		{
			name:     "next only",
			action:   nil,
			next:     action2,
			in:       1,
			expected: 3,
			err:      nil,
		},
		{
			name:     "action error",
			action:   actionWithErr,
			next:     action2,
			in:       1,
			expected: 5,
			err:      actionErr,
		},
		{
			name:     "next error",
			action:   action1,
			next:     actionWithErr,
			in:       1,
			expected: 5,
			err:      actionErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			action := wrap(tc.action, tc.next)
			out, err := action(tc.in)

			// assert
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_Unit_Action_Sequential(t *testing.T) {
	// arrange
	action1 := Do(func(in int) (int, error) {
		return in + 1, nil
	})
	action2 := Do(func(in int) (int, error) {
		return in + 2, nil
	})
	action3 := Do(func(in int) (int, error) {
		return in + 3, nil
	})

	actionErr := errors.New("test error")
	actionWithErr := Do(func(in int) (int, error) {
		return 5, actionErr
	})

	testCases := []struct {
		name     string
		actions  []Action
		in       any
		expected any
		err      error
	}{
		{
			name:     "three actions",
			actions:  []Action{action1, action2, action3},
			in:       1,
			expected: 7,
			err:      nil,
		},
		{
			name:     "single actions",
			actions:  []Action{action1},
			in:       1,
			expected: 2,
			err:      nil,
		},
		{
			name:     "no actions",
			actions:  nil,
			in:       1,
			expected: nil,
			err:      nil,
		},
		{
			name:     "action error",
			actions:  []Action{action1, actionWithErr, action3},
			in:       1,
			expected: 5,
			err:      actionErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			action := Sequential(tc.actions...)
			out, err := action(tc.in)

			// assert
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_Unit_Action_Do_Success(t *testing.T) {
	// act
	action := Do(func(in int) (int, error) {
		return in + 1, nil
	})
	out, err := action(1)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 2, out)
}

func Test_Unit_Action_Do_Error(t *testing.T) {
	// act
	actionErr := errors.New("test error")
	action := Do(func(in int) (int, error) {
		return in + 1, actionErr
	})
	out, err := action(1)

	// assert
	assert.Equal(t, actionErr, err)
	assert.Equal(t, 2, out)
}

func Test_Unit_Action_If(t *testing.T) {
	// arrange
	action1 := Do(func(in int) (int, error) {
		return in + 1, nil
	})
	action2 := Do(func(in int) (int, error) {
		return in + 2, nil
	})

	actionErr := errors.New("test error")
	actionWithErr := Do(func(in int) (int, error) {
		return 5, actionErr
	})

	testCases := []struct {
		name      string
		condition func(in int) (bool, error)
		ifTrue    Action
		ifFalse   Action
		in        any
		expected  any
		err       error
	}{
		{
			name:      "nil condition",
			condition: nil,
			ifTrue:    action1,
			ifFalse:   action2,
			in:        1,
			expected:  nil,
			err:       nil,
		},
		{
			name: "nil true",
			condition: func(in int) (bool, error) {
				return true, nil
			},
			ifTrue:   nil,
			ifFalse:  action2,
			in:       1,
			expected: nil,
			err:      nil,
		},
		{
			name: "nil false",
			condition: func(in int) (bool, error) {
				return true, nil
			},
			ifTrue:   action1,
			ifFalse:  nil,
			in:       1,
			expected: nil,
			err:      nil,
		},
		{
			name: "true",
			condition: func(in int) (bool, error) {
				return true, nil
			},
			ifTrue:   action1,
			ifFalse:  action2,
			in:       1,
			expected: 2,
			err:      nil,
		},
		{
			name: "false",
			condition: func(in int) (bool, error) {
				return false, nil
			},
			ifTrue:   action1,
			ifFalse:  action2,
			in:       1,
			expected: 3,
			err:      nil,
		},
		{
			name: "true error",
			condition: func(in int) (bool, error) {
				return true, nil
			},
			ifTrue:   actionWithErr,
			ifFalse:  action2,
			in:       1,
			expected: 5,
			err:      actionErr,
		},
		{
			name: "false error",
			condition: func(in int) (bool, error) {
				return false, nil
			},
			ifTrue:   action1,
			ifFalse:  actionWithErr,
			in:       1,
			expected: 5,
			err:      actionErr,
		},
		{
			name: "condition error",
			condition: func(in int) (bool, error) {
				return false, actionErr
			},
			ifTrue:   action1,
			ifFalse:  action2,
			in:       1,
			expected: nil,
			err:      actionErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			action := If(tc.condition, tc.ifTrue, tc.ifFalse)
			out, err := action(tc.in)

			// assert
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_Unit_Action_Parallel(t *testing.T) {
	// arrange
	action1 := Do(func(in int) (int, error) {
		return in + 1, nil
	})
	action2 := Do(func(in int) (int, error) {
		return in + 2, nil
	})
	action3 := Do(func(in int) (int, error) {
		return in + 3, nil
	})

	actionErr := errors.New("test error")
	actionWithErr := Do(func(in int) (int, error) {
		return 5, actionErr
	})

	reduce := func(in []Result) (int, error) {
		total := 0
		for _, v := range in {
			if v.Err != nil {
				return -1, v.Err
			}
			total += v.Out.(int)
		}
		return total, nil
	}

	testCases := []struct {
		name     string
		reduce   func(in []Result) (int, error)
		actions  []Action
		in       any
		expected any
		err      error
	}{
		{
			name:     "success",
			reduce:   reduce,
			actions:  []Action{action1, action2, action3},
			in:       1,
			expected: 9,
			err:      nil,
		},
		{
			name:     "no actions",
			reduce:   reduce,
			actions:  nil,
			in:       1,
			expected: 0,
			err:      nil,
		},
		{
			name:     "action error",
			reduce:   reduce,
			actions:  []Action{action1, actionWithErr, action3},
			in:       1,
			expected: -1,
			err:      actionErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			action := Parallel(tc.reduce, tc.actions...)
			out, err := action(1)

			// assert
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_Unit_Action_AllActions(t *testing.T) {
	// arrange
	add1 := func(in int) (int, error) {
		return in + 1, nil
	}
	add2 := func(in int) (int, error) {
		return in + 2, nil
	}
	add3 := func(in int) (int, error) {
		return in + 3, nil
	}
	isOdd := func(in int) (bool, error) {
		return in%2 == 1, nil
	}
	sum := func(in []Result) (int, error) {
		total := 0
		for _, v := range in {
			total += v.Out.(int)
		}
		return total, nil
	}

	// act
	action := Sequential(
		Do(add1), // 1 + 1 == 2
		Parallel(sum, // in == 2, result == 3 + 4 + 7 == 14
			Do(add1), // 2 + 1 == 3
			Do(add2), // 2 + 2 == 4
			Sequential( // in == 2
				Do(add1), // 2 + 1 == 3
				Do(add2), // 3 + 2 == 5
				If(isOdd, // in == 5 (true)
					Do(add2), // 5 + 2 == 7
					Do(add3), // skipped
				),
			)),
		If(isOdd, // in == 14 (false)
			NoOp(),   // skipped
			Do(add2), // 14 + 2 == 16
		),
	)
	out, err := action(1)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 16, out)
}
