package workflow

import (
	"sync"
)

type Action func(any) any

type Work struct {
	action Action
	first  *Work
	next   *Work
}

func Define(action func(in any) any) *Work {
	work := &Work{
		action: action,
	}
	work.first = work
	return work
}

func Wrap[T1 any, T2 any](action func(T1) T2) Action {
	return func(in interface{}) interface{} {
		input := in.(T1)
		output := action(input)
		return output
	}
}

func (w *Work) Start(in any) any {
	if w.first == nil {
		return nil
	}
	return w.first.Run(in)
}

func (w *Work) Run(in any) any {
	out := w.action(in)
	if w.next != nil {
		return w.next.Run(out)
	}
	return out
}

func (w *Work) Next(action Action) *Work {
	w.next = &Work{
		action: action,
		first:  w.first,
	}
	return w.next
}

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

func (w *Work) Parallel(results func([]any) any, work ...*Work) *Work {
	return w.Next(Wrap(func(in any) any {
		var lock sync.Mutex
		var outputs []any
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
		return results(outputs)
	}))
}
