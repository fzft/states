***

<div align="center">
    <b><em>Easy States</em></b><br>
    The simple, stupid States machine for golang, Enlightened by Java Project;
</div>

<div align="center">
</div>

## Tutorial

```go
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
		EventType(reflect.TypeOf(CoinEvent{})).
		EventHandler(NewUnLock()).
		TargetState(unlocked).
		Build()

	pushLocked := NewTransitionBuilder().
		Name("pushLocked").
		SourceState(locked).
		EventType(reflect.TypeOf(PushEvent{})).
		TargetState(locked).
		Build()

	lock := NewTransitionBuilder().
		Name("lock").
		SourceState(unlocked).
		EventType(reflect.TypeOf(PushEvent{})).
		EventHandler(NewLock()).
		TargetState(locked).
		Build()

	coinUnlocked := NewTransitionBuilder().
		Name("coinUnlocked").
		SourceState(unlocked).
		EventType(reflect.TypeOf(CoinEvent{})).
		TargetState(unlocked).
		Build()

	// 4. build FSM instance
	turnstileStateMachine := NewFiniteStateMachineBuilder(states, locked).
		registerTransitions(lock).
		registerTransitions(pushLocked).
		registerTransitions(unlock).
		registerTransitions(coinUnlocked).
		build()
```