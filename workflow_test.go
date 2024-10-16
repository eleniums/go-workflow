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
	assert.Equal(t, 5, workflow.Start(1))
	assert.Equal(t, 10, workflow.Start(2))
	assert.Equal(t, 15, workflow.Start(3))
}

func Test_Unit_Next(t *testing.T) {
	// arrange
	workflow := Define(Wrap(func(in int) int {
		return in * 5

	})).Next(Wrap(func(in int) int {
		return in + 2
	}))

	// act
	assert.Equal(t, 7, workflow.Start(1))
	assert.Equal(t, 12, workflow.Start(2))
	assert.Equal(t, 17, workflow.Start(3))
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
	assert.Equal(t, 6, workflow.Start(1))
	assert.Equal(t, 11, workflow.Start(2))
	assert.Equal(t, 16, workflow.Start(3))
}
