package workflow

import (
	"sync"
)

// A function that will be executed as part of a workflow.
type Action func(any) any

// Represents a block of work in a workflow.
type Work struct {
	// Action to be performed for this workflow.
	action Action

	// First block of work to be performed in the entire workflow.
	first *Work

	// Next block of work to be performed after this completes.
	next *Work
}

// Create the first block of work to be performed.
func Define(action func(in any) any) *Work {
	work := &Work{
		action: action,
	}
	work.first = work
	return work
}

// Convenience function to wrap a function with types into an Action function.
func Wrap[T1 any, T2 any](action func(T1) T2) Action {
	return func(in interface{}) interface{} {
		input := in.(T1)
		output := action(input)
		return output
	}
}

// Start the workflow from the very beginning.
func (w *Work) Start(in any) any {
	if w.first == nil {
		return nil
	}
	return w.first.run(in)
}

// Define the next block of work to be performed.
func (w *Work) Next(action Action) *Work {
	w.next = &Work{
		action: action,
		first:  w.first,
	}
	return w.next
}

// Conditionally execute another block of work. Only one path will be executed.
func (w *Work) If(condition func(in any) bool, ifTrue *Work, ifFalse *Work) *Work {
	return w.Next(Wrap(func(in any) any {
		var out any
		if condition(in) {
			out = ifTrue.Start(in)
		} else {
			out = ifFalse.Start(in)
		}
		return out
	}))
}

// Execute multiple blocks of work in parallel. The result function will combine all results into a single result.
func (w *Work) Parallel(result func([]any) any, work ...*Work) *Work {
	return w.Next(Wrap(func(in any) any {
		var outputs []any
		var lock sync.Mutex
		var wg sync.WaitGroup
		for _, w := range work {
			wg.Add(1)
			go func(in any) {
				defer wg.Done()
				out := w.Start(in)
				lock.Lock()
				defer lock.Unlock()
				outputs = append(outputs, out)
			}(in)
		}
		wg.Wait()
		return result(outputs)
	}))
}

// Recursive method to run a block of work and then run the next block of work.
// NOTE: This should remain internal because the recursive nature of this method
// could have unintended consequences to a consumer of this package.
func (w *Work) run(in any) any {
	out := w.action(in)
	if w.next != nil {
		return w.next.run(out)
	}
	return out
}
