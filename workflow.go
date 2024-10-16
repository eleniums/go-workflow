package workflow

type Action func(any) any
type Condition func() bool

type Work struct {
	action Action
	first  *Work
	next   *Work
}

func Define[T1 any, T2 any](action func(T1) T2) *Work {
	work := &Work{
		action: Wrap(action),
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

func (w *Work) If(condition func() bool, ifTrue *Work, ifFalse *Work) *Work {
	// TODO: implement
	return nil
}

func (w *Work) Parallel(work1 Action, work2 Action) *Work {
	// TODO: implement
	return nil
}
