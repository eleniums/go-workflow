package workflow

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func Test_Unit_Wrap(t *testing.T) {
	// arrange
	action1 := Do(func(in int) int {
		fmt.Println("action1")
		return in + 1
	})
	action2 := Do(func(in int) int {
		fmt.Println("action2")
		return in + 1
	})
	action3 := Do(func(in int) int {
		fmt.Println("action3")
		return in + 1
	})

	combo := Combine(action1, action2, action3)
	// combo = Wrap(combo, action3)

	// act
	assert.Equal(t, 4, combo(1))
}

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
