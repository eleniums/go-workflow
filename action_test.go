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
		expected any
	}{
		{
			name:     "action and next",
			action:   action1,
			next:     action2,
			expected: 4,
		},
		{
			name:     "no actions",
			action:   nil,
			next:     nil,
			expected: nil,
		},
		{
			name:     "action only",
			action:   action1,
			next:     nil,
			expected: 2,
		},
		{
			name:     "next only",
			action:   nil,
			next:     action2,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			combo := Wrap(tc.action, tc.next)

			// assert
			assert.Equal(t, tc.expected, combo(1))
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
		expected any
	}{
		{
			name:     "three actions",
			actions:  []Action{action1, action2, action3},
			expected: 7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			combo := Combine(tc.actions...)

			// assert
			assert.Equal(t, tc.expected, combo(1))
		})
	}
}

// func Test_Unit_Action_Combine(t *testing.T) {
// 	// arrange
// 	action1 := Do(func(in int) int {
// 		fmt.Println("action1")
// 		return in + 1
// 	})
// 	action2 := Do(func(in int) int {
// 		fmt.Println("action2")
// 		return in + 2
// 	})
// 	action3 := Do(func(in int) int {
// 		fmt.Println("action3")
// 		return in + 3
// 	})

// 	// act
// 	combo := Combine(action1, action2, action3)

// 	// assert
// 	assert.Equal(t, 7, combo(1))
// }

// func Test_Unit_Create(t *testing.T) {
// 	// arrange
// 	workflow := Create().Do(func(in int) int {
// 		return in * 5
// 	})

// 	// act
// 	assert.Equal(t, 5, workflow.Start(1))
// 	assert.Equal(t, 10, workflow.Start(2))
// 	assert.Equal(t, 15, workflow.Start(3))
// }

// func Test_Unit_Next(t *testing.T) {
// 	// arrange
// 	workflow := Create().Do(func(in int) int {
// 		return in * 5
// 	}).Next(Do(func(in int) int {
// 		return in + 2
// 	}))

// 	// act
// 	assert.Equal(t, 7, workflow.Start(1))
// 	assert.Equal(t, 12, workflow.Start(2))
// 	assert.Equal(t, 17, workflow.Start(3))
// }

// func Test_Unit_MultipleNext(t *testing.T) {
// 	// arrange
// 	workflow := Create().Do(func(in int) int {
// 		return in * 5
// 	}).Next(Do(func(in int) int {
// 		return in + 2
// 	})).Next(Do(func(in int) int {
// 		return in - 1
// 	}))

// 	// act
// 	assert.Equal(t, 6, workflow.Start(1))
// 	assert.Equal(t, 11, workflow.Start(2))
// 	assert.Equal(t, 16, workflow.Start(3))
// }

// func Test_Unit_If(t *testing.T) {
// 	// arrange
// 	workflow := Create().Do(func(in int) int {
// 		return in * 5
// 	}).Next(Do(func(in int) int {
// 		return in + 2
// 	})).If(func(in any) bool {
// 		v := in.(int)
// 		return v%2 == 1
// 	}, Create().Do(func(in int) int {
// 		return in + 4
// 	})), Create().Do(func(in int) int {
// 		return in - 4
// 	})

// 	// act
// 	assert.Equal(t, 11, workflow.Start(1))
// 	assert.Equal(t, 8, workflow.Start(2))
// 	assert.Equal(t, 21, workflow.Start(3))
// }

// func Test_Unit_Parallel(t *testing.T) {
// 	// arrange
// 	workflow := Create().Do(func(in int) int {
// 		return in * 5
// 	}).Parallel(func(results []any) any {
// 		total := 0
// 		for _, v := range results {
// 			total += v.(int)
// 		}
// 		return total
// 	}, Create().Do(func(in int) int {
// 		return in + 2
// 	})), Create().Do(func(in int) int {
// 		return in - 1
// 	}))).Next(Do(func(total int) int {
// 		fmt.Println(total)
// 		return total
// 	}))

// 	// act
// 	assert.Equal(t, 11, workflow.Start(1))
// 	assert.Equal(t, 21, workflow.Start(2))
// 	assert.Equal(t, 31, workflow.Start(3))
// }
