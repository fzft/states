package states

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"reflect"
	"strings"
	"sync"
)

type FiniteStateMachineBuilder struct {
	fsm *FiniteStateMachineImpl
	mdv *FiniteStateMachineDefinitionValidator
	tdv *TransitionDefinitionValidator
}

func NewFiniteStateMachineBuilder(states Set, initialState State) *FiniteStateMachineBuilder {
	return &FiniteStateMachineBuilder{
		fsm: NewFiniteStateMachineImpl(states, initialState),
		mdv: NewMDV(),
		tdv: NewTDV(),
	}
}

func (b *FiniteStateMachineBuilder) RegisterTransitions(ts ...Transition) *FiniteStateMachineBuilder {
	group := new(errgroup.Group)
	for _, t := range ts {
		tv := t
		group.Go(func() error {
			return b.tdv.validateTransitionDefinition(tv, b.fsm)
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Printf("get error: %v\n", err)
		return b
	}
	b.fsm.RegisterTransition(ts...)
	return b
}

func (b *FiniteStateMachineBuilder) registerFinalStates(states ...State) *FiniteStateMachineBuilder {
	b.fsm.RegisterFinalState(states...)
	return b
}

func (b *FiniteStateMachineBuilder) Build() FiniteStateMachine {
	b.mdv.validateFiniteStateMachineDefinition(b.fsm)
	return b.fsm
}

type FiniteStateMachineImpl struct {
	mu             sync.Mutex
	currentState   State
	initialState   State
	finalStates    Set
	states         Set
	transitions    Set
	lastEvent      Event
	lastTransition Transition
}

func NewFiniteStateMachineImpl(states Set, initialState State) *FiniteStateMachineImpl {
	return &FiniteStateMachineImpl{
		states:       states,
		initialState: initialState,
		currentState: initialState,
		transitions:  Set{},
		finalStates:  Set{},
	}
}

func (m *FiniteStateMachineImpl) RegisterTransition(ts ...Transition) {
	for _, t := range ts {
		m.transitions.Insert(t)
	}
}

func (m *FiniteStateMachineImpl) RegisterFinalState(states ...State) {
	for _, s := range states {
		m.finalStates.Insert(s)
	}
}

//Fire ...
func (m *FiniteStateMachineImpl) Fire(event Event) (State, error) {
	defer recoverPanic()
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.finalStates) == 0 && m.finalStates.Has(m.currentState) {
		return m.currentState, nil
	}
	for t, _ := range m.transitions {
		if tv, ok := t.(Transition); ok {
			if m.currentState.Equals(tv.SourceState()) && tv.EventType() == reflect.TypeOf(event) && m.states.Has(tv.TargetState()) {
				if tv.EventHandler() != nil {
					tv.EventHandler().HandleEvent(event)
				}
				m.currentState = tv.TargetState()

				m.lastEvent = event
				m.lastTransition = tv
				break
			}
		}
	}
	return m.currentState, nil
}

//CurrentState ...
func (m *FiniteStateMachineImpl) CurrentState() State {
	return m.currentState
}

//InitialState ...
func (m *FiniteStateMachineImpl) InitialState() State {
	return m.initialState
}

//FinalStates ...
func (m *FiniteStateMachineImpl) FinalStates() Set {
	return m.finalStates
}

//States ...
func (m *FiniteStateMachineImpl) States() Set {
	return m.states
}

//Transitions ...
func (m *FiniteStateMachineImpl) Transitions() Set {
	return m.transitions
}

//LastEvent ...
func (m *FiniteStateMachineImpl) LastEvent() Event {
	return m.lastEvent
}

//LastTransition ...
func (m *FiniteStateMachineImpl) LastTransition() Transition {
	return m.lastTransition
}

type FiniteStateMachineDefinitionValidator struct {
}

func NewMDV() *FiniteStateMachineDefinitionValidator {
	return &FiniteStateMachineDefinitionValidator{}
}

func (v *FiniteStateMachineDefinitionValidator) validateFiniteStateMachineDefinition(m FiniteStateMachine) error {
	states := m.States()

	initialState := m.InitialState()
	if !states.Has(initialState) {
		return errors.New(fmt.Sprintf("Initial state %s must belong to FSM states: %s",
			initialState.Name(), dumpFSMStates(states)))
	}
	for s, _ := range m.FinalStates() {
		if !states.Has(s) {
			return errors.New(fmt.Sprintf("Final state %s must belong to FSM states: %s",
				initialState.Name(), dumpFSMStates(states)))
		}
	}
	return nil
}

func recoverPanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		fmt.Println(err)
	}
}

func dumpFSMStates(states Set) string {
	var b strings.Builder
	for s, _ := range states {
		if sv, ok := s.(State); ok {
			b.WriteString(sv.Name())
			b.WriteString(";")
		}
	}
	return b.String()
}
