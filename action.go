package workflow

import (
	"sync"
)

// A simple function definition with a single input and output.
type Action func(any) (any, error)

// Contains an action output and associated error, if any.
type Result struct {
	Out any
	Err error
}

// Encapsulate a function with types into an action.
func Do[T1 any, T2 any](action func(T1) (T2, error)) Action {
	return func(in any) (any, error) {
		input := in.(T1)
		return action(input)
	}
}

// Combines multiple actions into a single action that will execute based on the order the actions were passed.
func Sequential(actions ...Action) Action {
	if len(actions) == 0 {
		return NoOp()
	}

	sequential := actions[0]
	for i := 1; i < len(actions); i++ {
		sequential = wrap(sequential, actions[i])
	}

	return sequential
}

// Execute multiple actions in parallel. The reduce function should combine all parallel results into a single result.
func Parallel[T any](reduce func(in []Result) (T, error), actions ...Action) Action {
	return func(in any) (any, error) {
		var outputs []Result
		var lock sync.Mutex
		var wg sync.WaitGroup
		for _, v := range actions {
			wg.Add(1)
			go func(in any) {
				defer wg.Done()
				out, err := v(in)
				lock.Lock()
				defer lock.Unlock()
				outputs = append(outputs, Result{
					Out: out,
					Err: err,
				})
			}(in)
		}
		wg.Wait()
		return reduce(outputs)
	}
}

// Conditionally execute another action. Only one action will be executed.
func If[T any](condition func(in T) (bool, error), ifTrue Action, ifFalse Action) Action {
	// all functions must be valid
	if condition == nil || ifTrue == nil || ifFalse == nil {
		return NoOp()
	}

	return func(in any) (any, error) {
		input := in.(T)
		condition, err := condition(input)
		if err != nil {
			return nil, err
		}
		if condition {
			return ifTrue(in)
		} else {
			return ifFalse(in)
		}
	}
}

// Executes an action and calls the handle function if an error occurs.
func Catch(action Action, handle Action) Action {
	return func(in any) (any, error) {
		out, err := action(in)
		if err != nil {
			return handle(out)
		}
		return out, nil
	}
}

func Retry() Action {
	// TODO: implement retry function
	return nil
}

// Returns an action that does nothing and returns nil.
func NoOp() Action {
	return func(in any) (any, error) {
		return nil, nil
	}
}

// Wraps provided actions so that "action" is called first and then "next" is called.
func wrap(action Action, next Action) Action {
	if action == nil && next == nil {
		return NoOp()
	}
	if action == nil {
		return next
	}
	if next == nil {
		return action
	}
	return func(in any) (any, error) {
		out, err := action(in)
		if err != nil {
			return out, err
		}
		return next(out)
	}
}
