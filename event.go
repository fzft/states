package states

import (
	"time"
)

type empty struct{}
type t interface{}
type Set map[t]empty

func (s Set) Has(item t) bool {
	_, exists := s[item]
	return exists
}

func (s Set) Insert(item t) {
	s[item] = empty{}
}

func (s Set) Delete(item t) {
	delete(s, item)
}


type Event interface {
	Name() string
	Timestamp() int64
}

type EventHandler interface {
	HandleEvent(e Event) error
}

type FiniteStateMachine interface {
	CurrentState() State
	InitialState() State
	FinalStates() Set
	States() Set
	Transitions() Set
	LastEvent() Event
	LastTransition() Transition
	Fire(e Event) (State, error)
}

type AbstractEvent struct {
	name      string
	timestamp int64
}

func NewAbstractEvent(name string) Event {
	return &AbstractEvent{
		name:      name,
		timestamp: time.Now().Unix(),
	}
}

func (e *AbstractEvent) Name() string {
	return e.name
}
func (e *AbstractEvent) Timestamp() int64 {
	return e.timestamp
}
