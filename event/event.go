// Package event provides classes and interfaces to implement
// event senders and receivers.
//
// There are also some types of events, but these only serve
// as examples and are not thought through.
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
