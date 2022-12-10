package msgqueue

// Interface definition for events that are emitted using an EventEmitter
// Currently, the only requirement is that events are self-describing so that
// event emitter and listeners can infer an event's name.
type Event interface {
	EventName() string
}

// EventEmitter describes an interface for a class that emits events
type EventEmitter interface {
	Emit(e Event) error
}

type EventListener interface {
}
