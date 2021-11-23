package states

import (
	"errors"
	"fmt"
	"reflect"
)

type Transition interface {
	Name() string
	SourceState() State
	TargetState() State
	EventType() reflect.Type
	EventHandler() EventHandler
}

type TransitionDefinitionValidator struct {
}

func NewTDV() *TransitionDefinitionValidator {
	return &TransitionDefinitionValidator{}
}

func (v *TransitionDefinitionValidator) validateTransitionDefinition(t Transition, m FiniteStateMachine) error {
	tn := t.Name()
	srcState := t.SourceState()
	tState := t.TargetState()

	if !m.States().has(srcState) {
		return errors.New(fmt.Sprintf("Source state %s is not registered in FSM states for transition %s", srcState.Name(), tn))
	}
	if !m.States().has(tState) {
		return errors.New(fmt.Sprintf("target state %s is not registered in FSM states for transition %s", srcState.Name(), tn))
	}
	return nil
}

type TransitionImpl struct {
	name         string
	sourceState  State
	targetState  State
	eventType    reflect.Type
	eventHandler EventHandler
}

func NewTransitionImpl() *TransitionImpl {
	return &TransitionImpl{}
}

func (t *TransitionImpl) Name() string {
	return t.name
}
func (t *TransitionImpl) SourceState() State {
	return t.sourceState
}
func (t *TransitionImpl) TargetState() State {
	return t.targetState
}
func (t *TransitionImpl) EventType() reflect.Type {
	return t.eventType
}
func (t *TransitionImpl) EventHandler() EventHandler {
	return t.eventHandler
}

func (t *TransitionImpl) Equals(o interface{}) bool {
	if tv, ok := o.(*TransitionImpl); ok {
		return tv.eventType == t.eventType && t.sourceState.Equals(tv.sourceState)
	}
	return false
}

type TransitionBuilder struct {
	transition *TransitionImpl
}

func NewTransitionBuilder() *TransitionBuilder {
	return &TransitionBuilder{
		transition: NewTransitionImpl(),
	}
}

func (b *TransitionBuilder) Name(name string) *TransitionBuilder {
	b.transition.name = name
	return b
}
func (b *TransitionBuilder) SourceState(sourceState State) *TransitionBuilder {
	b.transition.sourceState = sourceState
	return b
}
func (b *TransitionBuilder) TargetState(targetState State) *TransitionBuilder {
	b.transition.targetState = targetState
	return b
}
func (b *TransitionBuilder) EventType(eventType reflect.Type) *TransitionBuilder {
	b.transition.eventType = eventType
	return b
}
func (b *TransitionBuilder) EventHandler(eventHandler EventHandler) *TransitionBuilder {
	b.transition.eventHandler = eventHandler
	return b
}
func (b *TransitionBuilder) Build() Transition {
	return b.transition
}
