// Package native represents a layer between the operating
// system and ui2go.
// It defines functions and interfaces, that are implemented
// operating system specific.
package native

import (
	"github.com/ungerik/go-cairo"
)

// Window represents a typical user interface window
// in an operating system independent way.
type Window interface {
	// Returns a cairo surface that can be used for drawing.
	Surface() *cairo.Surface
	// Flush flushes pending messages to the unterlying window system.
	Flush()
	// EventChan returns a channel of window events.
	EventChan() <-chan interface{}
	Close() error
}

// NewWindow is a function that returns a user interface window.
// It must be implemented individually by different window subsystems.
var NewWindow = func() Window {
	panic("Function NewWindow not defined." +
		"Possible cause: No operating system specific implementation.")
}
