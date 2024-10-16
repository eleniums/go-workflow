package workflow

type Action func(any) any

type Work struct {
	action Action
	next   *Work
}

func Define[T1 any, T2 any](action func(T1) T2) *Work {
	return &Work{
		action: Wrap(action),
	}
}

func Wrap[T1 any, T2 any](action func(T1) T2) Action {
	return func(in interface{}) interface{} {
		input := in.(T1)
		output := action(input)
		return output
	}
}

func (w *Work) Next(action Action) *Work {
	w.next = &Work{
		action: action,
	}
	return w.next
}

// func (w *Work[T1, T2]) Run(in T1) T2 {
// 	out := w.action(in)
// 	if w.next != nil {
// 		return w.next.Run(out)
// 	}
// 	return out
// }
