package workflow

import (
	"sync"
)

// TODO: add error to the Action signature?
// A simple function definition with a single input and output.
type Action func(any) any

// Returns an action that does nothing and returns nil.
func NoOp() Action {
	return func(in any) any {
		return nil
	}
}

// Wraps provided actions so that "action" is called first and then "next" is called.
func Wrap(action Action, next Action) Action {
	if action == nil && next == nil {
		return NoOp()
	}
	if action == nil {
		return next
	}
	if next == nil {
		return action
	}
	return func(in any) any {
		out := action(in)
		return next(out)
	}
}

// Combines multiple actions into a single action that will execute based on the order the actions were passed.
func Combine(actions ...Action) Action {
	if len(actions) == 0 {
		return NoOp()
	}

	combined := actions[0]
	for i := 1; i < len(actions); i++ {
		combined = Wrap(combined, actions[i])
	}

	return combined
}

// Wrap a function with types into an action.
func Do[T1 any, T2 any](action func(T1) T2) Action {
	return func(in any) any {
		input := in.(T1)
		output := action(input)
		return output
	}
}

// Conditionally execute another action. Only one action will be executed.
func If[T any](condition func(in T) bool, ifTrue Action, ifFalse Action) Action {
	return func(in any) any {
		input := in.(T)
		var out any
		if condition(input) {
			out = ifTrue(in)
		} else {
			out = ifFalse(in)
		}
		return out
	}
}

// Execute multiple actions in parallel. The result function will combine all results into a single result.
func Parallel[T any](result func([]any) T, actions ...Action) Action {
	return func(in any) any {
		var outputs []any
		var lock sync.Mutex
		var wg sync.WaitGroup
		for _, v := range actions {
			wg.Add(1)
			go func(in any) {
				defer wg.Done()
				out := v(in)
				lock.Lock()
				defer lock.Unlock()
				outputs = append(outputs, out)
			}(in)
		}
		wg.Wait()
		return result(outputs)
	}
}
