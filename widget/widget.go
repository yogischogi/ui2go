// Package widget contains classes for widgets and the
// main window of a program.
package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/ungerik/go-cairo"
	"image"
)

// Drawable is anything that could be drawn onto the screen.
type Drawable interface {
	Draw()
	SetArea(image.Rectangle)
	Area() image.Rectangle
	SetSurface(*cairo.Surface)
	Surface() *cairo.Surface

	// MinSize is the minimum space the object needs to
	// to be displayed. It is calculated by the object itself.
	MinSize() image.Point
}

// Widget is something that can be drawn onto the screen
// and is capable to send and receive events.
type Widget interface {
	Drawable
	event.Receiver
	event.Sender
}

// Layout positions drawable objects.
type Layout interface {
	Drawable
	Addf(layoutDef string, components ...Drawable)
}

// Container is a container for other widgets that
// should be displayed.
// A typical container is a window.
type Container interface {
	Layout
	event.Receiver
	event.Sender
}

// Accessible is a simple interface for accessibility purposes.
type Accessible interface {
	Caption() string
	Tip() string
	Description() string
}
