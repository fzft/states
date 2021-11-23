package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
	"github.com/fzft/states"
)

type PushEvent struct {
	name      string
	timestamp int64
}

func NewPushEvent() Event {
	return &PushEvent{
		name:      "PushEvent",
		timestamp: time.Now().Unix(),
	}
}

func (e *PushEvent) Name() string {
	return e.name
}
func (e *PushEvent) Timestamp() int64 {
	return e.timestamp
}

type CoinEvent struct {
	name      string
	timestamp int64
}

func NewCoinEvent() Event {
	return &CoinEvent{
		name:      "CoinEvent",
		timestamp: time.Now().Unix(),
	}
}

func (e *CoinEvent) Name() string {
	return e.name
}
func (e *CoinEvent) Timestamp() int64 {
	return e.timestamp
}

// Lock an action
type Lock struct {
}

func NewLock() *Lock {
	return &Lock{}
}

func (l *Lock) HandleEvent(event Event) error {
	fmt.Printf("Notified about event %s triggered at %d\n", event.Name(), event.Timestamp())
	fmt.Println("Locking turnstile..")
	return nil
}

// UnLock an action
type UnLock struct {
}

func NewUnLock() *UnLock {
	return &UnLock{}
}

func (l *UnLock) HandleEvent(event Event) error {
	fmt.Printf("Notified about event %s triggered at %d\n", event.Name(), event.Timestamp())
	fmt.Println("Unlocking turnstile..")
	return nil
}

func main() {

	// 1. define states
	locked := NewState("locked")
	unlocked := NewState("unlocked")

	states := set{}
	states.insert(locked)
	states.insert(unlocked)

	// 2. define events

	// 3. transitions
	unlock := NewTransitionBuilder().
		Name("unlock").
		SourceState(locked).
		EventType(reflect.TypeOf(NewCoinEvent())).
		EventHandler(NewUnLock()).
		TargetState(unlocked).
		Build()

	pushLocked := NewTransitionBuilder().
		Name("pushLocked").
		SourceState(locked).
		EventType(reflect.TypeOf(NewPushEvent())).
		TargetState(locked).
		Build()

	lock := NewTransitionBuilder().
		Name("lock").
		SourceState(unlocked).
		EventType(reflect.TypeOf(NewPushEvent())).
		EventHandler(NewLock()).
		TargetState(locked).
		Build()

	coinUnlocked := NewTransitionBuilder().
		Name("coinUnlocked").
		SourceState(unlocked).
		EventType(reflect.TypeOf(NewCoinEvent())).
		TargetState(unlocked).
		Build()

	// 4. build FSM instance
	turnstileStateMachine := NewFiniteStateMachineBuilder(states, locked).
		registerTransitions(lock).
		registerTransitions(pushLocked).
		registerTransitions(unlock).
		registerTransitions(coinUnlocked).
		build()

	fmt.Printf("Turnstile initial state : %s\n", turnstileStateMachine.CurrentState().Name())

	fmt.Println("Which event do you want to fire")
	fmt.Println("1. Push [p]")
	fmt.Println("2. Coin [c]")
	fmt.Println("Press [q] to quit tutorial.")
	var input string
	for {
		fmt.Scan(&input)
		in := strings.TrimSpace(input)
		switch in {
		case "p":
			fmt.Printf("input = %s\n", in)
			fmt.Println("Firing push event ...")
			turnstileStateMachine.Fire(NewPushEvent())
			fmt.Printf("Turnstile state : %s\n", turnstileStateMachine.CurrentState().Name())
		case "c":
			fmt.Printf("input = %s\n", in)
			fmt.Println("Firing push event ...")
			turnstileStateMachine.Fire(NewCoinEvent())
			fmt.Printf("Turnstile state : %s\n", turnstileStateMachine.CurrentState().Name())
		case "q":
			fmt.Printf("input = %s\n", in)
			fmt.Println("bye!")
			os.Exit(0)
		}
	}

}
