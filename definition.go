package workflow

// Definition of a workflow.
type Definition struct {
	actions []Action
}

func Create() *Definition {
	return &Definition{}
}

func (d *Definition) Next(action Action) *Definition {
	d.actions = append(d.actions, action)
	return d
}

func (d *Definition) Run(in any) any {
	combined := d.Compile()
	return combined(in)
}

func (d *Definition) Compile() Action {
	return Combine(d.actions...)
}
