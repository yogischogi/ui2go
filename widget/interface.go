package widget

import (
	"code.google.com/p/ui2go/event"
	"image"
	"image/draw"
)

// Drawable is anything that could be drawn onto the screen.
type Drawable interface {
	Draw()
	SetArea(image.Rectangle)
	Area() image.Rectangle
	SetScreen(draw.Image)
	Screen() draw.Image

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
