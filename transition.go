package states


type Transition interface {
	Name() string
	SourceState() State
	TargetState() State
	EventType() Event
	EventHandler() EventHandler
}
