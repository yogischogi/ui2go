// Package widget contains classes for widgets and the
// main window of a program.
package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/ungerik/go-cairo"
	"image"
	"strconv"
)

func init() {
	// Add constructors to ConstructorRegistry
	ConstructorRegistry["Button"] = NewButtonFromJson
	ConstructorRegistry["Canvas"] = NewCanvasFromJson
	ConstructorRegistry["Spacer"] = NewSpacerFromJson
}

// ComponentRegistry is a registry for all components that can
// be drawn onto the screen.
// The purpose is to access widgets through an id, for example
// Button("okButton").SetText("All Right"), where Button is a
// function that asks the registry for the appropriate widget.
// This is very useful when the widgets are to be identified
// by an id, like in web applications or a textual UI description.
// Newly created widgets should register at this registry and
// provide an access function like func Label(id string) *Label.
var ComponentRegistry = make(map[string]Drawable)

// ContructorRegistry maps component names to component constructors
// that use a JSON definition.
var ConstructorRegistry = make(map[string]func(json []byte) Drawable)

// idCounter is used to create unique IDs for different components.
var idCounter int = 0

// NewId creates a unique string as an identifier.
func NewId() string {
	idCounter++
	return strconv.Itoa(idCounter)
}

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
