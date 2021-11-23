package states

import (
	"time"
)

type empty struct{}
type t interface{}
type set map[t]empty

func (s set) has(item t) bool {
	_, exists := s[item]
	return exists
}

func (s set) insert(item t) {
	s[item] = empty{}
}

func (s set) delete(item t) {
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
	FinalStates() set
	States() set
	Transitions() set
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
