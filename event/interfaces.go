package event

// Receiver is the basic interface that every event receiver
// must satisfy.
//
// Received events are handled by an event handler function,
// that gets the events directly through it's argument, or
// through a channel of events.
type Receiver interface {
	ListenTo(sender Sender)
	UnlistenTo(sender Sender)
	ReceiveEvent(evt interface{})
	SetEvtHandler(handler func(evt interface{}))
	SetEvtChanHandler(handler func(evtChan <-chan interface{}))
}

// Sender is the basic interface that every event sender
// must satisfy.
//
// An event sender should be able to send events to arbritrary numbers
// of event receivers.
type Sender interface {
	SendEvent(evt interface{})
	AddReceiver(receiver Receiver)
	RemoveReceiver(receiver Receiver)
}

// Event is the basic interface for all events.
type Event interface {
	// Sender of the event. This is implemented as a string,
	// because it's great for debugging and possible network
	// implementations.
	Sender() string

	// Time is important for many purposes as
	// double clicks, drawing operations, reaction
	// to events in a multitasking environments or
	// simple timeouts. It will be implemented later.
	// Time()
}
