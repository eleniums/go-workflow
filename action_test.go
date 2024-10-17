package workflow

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func Test_Unit_Action_NoOp(t *testing.T) {
	// arrange
	action := NoOp()

	// act
	assert.Nil(t, action(1))
}

func Test_Unit_Action_Wrap(t *testing.T) {
	// arrange
	action1 := Do(func(in int) int {
		return in + 1
	})
	action2 := Do(func(in int) int {
		return in + 2
	})

	testCases := []struct {
		name     string
		action   Action
		next     Action
		in       any
		expected any
	}{
		{
			name:     "action and next",
			action:   action1,
			next:     action2,
			in:       1,
			expected: 4,
		},
		{
			name:     "no actions",
			action:   nil,
			next:     nil,
			in:       1,
			expected: nil,
		},
		{
			name:     "action only",
			action:   action1,
			next:     nil,
			in:       1,
			expected: 2,
		},
		{
			name:     "next only",
			action:   nil,
			next:     action2,
			in:       1,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			combo := Wrap(tc.action, tc.next)

			// assert
			assert.Equal(t, tc.expected, combo(tc.in))
		})
	}
}

func Test_Unit_Action_Combine(t *testing.T) {
	// arrange
	action1 := Do(func(in int) int {
		return in + 1
	})
	action2 := Do(func(in int) int {
		return in + 2
	})
	action3 := Do(func(in int) int {
		return in + 3
	})

	testCases := []struct {
		name     string
		actions  []Action
		in       any
		expected any
	}{
		{
			name:     "three actions",
			actions:  []Action{action1, action2, action3},
			in:       1,
			expected: 7,
		},
		{
			name:     "single actions",
			actions:  []Action{action1},
			in:       1,
			expected: 2,
		},
		{
			name:     "no actions",
			actions:  nil,
			in:       1,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			combo := Combine(tc.actions...)

			// assert
			assert.Equal(t, tc.expected, combo(tc.in))
		})
	}
}

func Test_Unit_Action_Do(t *testing.T) {
	// arrange
	action := Do(func(in int) int {
		return in + 1
	})

	// act
	assert.Equal(t, 2, action(1))
}

func Test_Unit_Action_If(t *testing.T) {
	// arrange
	action1 := Do(func(in int) int {
		return in + 1
	})
	action2 := Do(func(in int) int {
		return in + 2
	})

	testCases := []struct {
		name      string
		condition func(in int) bool
		ifTrue    Action
		ifFalse   Action
		in        any
		expected  any
	}{
		{
			name:      "nil condition",
			condition: nil,
			ifTrue:    action1,
			ifFalse:   action2,
			in:        1,
			expected:  nil,
		},
		{
			name: "nil true",
			condition: func(in int) bool {
				return true
			},
			ifTrue:   nil,
			ifFalse:  action2,
			in:       1,
			expected: nil,
		},
		{
			name: "nil false",
			condition: func(in int) bool {
				return true
			},
			ifTrue:   action1,
			ifFalse:  nil,
			in:       1,
			expected: nil,
		},
		{
			name: "true",
			condition: func(in int) bool {
				return true
			},
			ifTrue:   action1,
			ifFalse:  action2,
			in:       1,
			expected: 2,
		},
		{
			name: "false",
			condition: func(in int) bool {
				return false
			},
			ifTrue:   action1,
			ifFalse:  action2,
			in:       1,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			combo := If(tc.condition, tc.ifTrue, tc.ifFalse)

			// assert
			assert.Equal(t, tc.expected, combo(tc.in))
		})
	}
}

func Test_Unit_Action_Parallel(t *testing.T) {
	// arrange
	action1 := Do(func(in int) int {
		return in + 1
	})
	action2 := Do(func(in int) int {
		return in + 2
	})
	action3 := Do(func(in int) int {
		return in + 3
	})
	result := func(in []any) int {
		total := 0
		for _, v := range in {
			total += v.(int)
		}
		return total
	}

	action := Parallel(result, action1, action2, action3)

	// act
	assert.Equal(t, 7, action(1))
}
