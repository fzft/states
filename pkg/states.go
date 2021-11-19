package main

type State struct {
	name string
}

func NewState(name string) State {
	return State{
		name: name,
	}
}

func (s State) Name() string {
	return s.name
}

func (s State) Equals(o interface{}) bool {
	if v, ok := o.(State); ok && v.name == s.name {
		return true
	}
	return false
}
