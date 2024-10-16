package workflow

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func Test_Unit_Define(t *testing.T) {
	// arrange
	workflow := Define(Wrap(func(in int) int {
		return in * 5
	}))

	// act
	assertResult(t, 5, workflow.Start(1))
	assertResult(t, 10, workflow.Start(2))
	assertResult(t, 15, workflow.Start(3))
}

func Test_Unit_Next(t *testing.T) {
	// arrange
	workflow := Define(Wrap(func(in int) int {
		return in * 5
	})).Next(Wrap(func(in int) int {
		return in + 2
	}))

	// act
	assertResult(t, 7, workflow.Start(1))
	assertResult(t, 12, workflow.Start(2))
	assertResult(t, 17, workflow.Start(3))
}

func Test_Unit_MultipleNext(t *testing.T) {
	// arrange
	workflow := Define(Wrap(func(in int) int {
		return in * 5
	})).Next(Wrap(func(in int) int {
		return in + 2
	})).Next(Wrap(func(in int) int {
		return in - 1
	}))

	// act
	assertResult(t, 6, workflow.Start(1))
	assertResult(t, 11, workflow.Start(2))
	assertResult(t, 16, workflow.Start(3))
}

func Test_Unit_If(t *testing.T) {
	// arrange
	workflow := Define(Wrap(func(in int) int {
		return in * 5
	})).Next(Wrap(func(in int) int {
		return in + 2
	})).If(func(in []any) bool {
		v := in[0].(int)
		return v%2 == 1
	}, Define(Wrap(func(in int) int {
		return in + 4
	})), Define(Wrap(func(in int) int {
		return in - 4
	})))

	// act
	assertResult(t, 11, workflow.Start(1))
	assertResult(t, 8, workflow.Start(2))
	assertResult(t, 21, workflow.Start(3))
}

func assertResult[T1 any, T2 any](t *testing.T, expected T1, actual []T2) {
	assert.NotEmpty(t, actual)
	assert.Equal(t, expected, actual[0])
}
