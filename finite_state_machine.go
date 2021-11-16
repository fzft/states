package states

import "sync"

type FiniteStateMachineBuilder struct {
}

type FiniteStateMachineImpl struct {
	mu             sync.Mutex
	currentState   State
	initialState   State
	finalStates    set
	states         set
	transitions    set
	lastEvent      Event
	lastTransition Transition
}

func NewFiniteStateMachineImpl(states set, initialState State) *FiniteStateMachineImpl {
	return &FiniteStateMachineImpl{
		states:       states,
		initialState: initialState,
		transitions:  set{},
		finalStates:  set{},
	}
}

func (m *FiniteStateMachineImpl) Fire(event Event) (State, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.finalStates) == 0 && m.finalStates.has(m.currentState) {
		return m.currentState, nil
	}

	for _, i := range m.transitions {
		t := i.(Transition)


	}
}
