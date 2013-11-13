package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/ungerik/go-cairo"
	"image"
)

// WidgetPrototype is a simple widget that implements the Widget interface.
// It's main purpose is to be used as a parent class for other widgets.
type WidgetPrototype struct {
	event.Sender
	event.Receiver
	area    image.Rectangle
	minSize image.Point
	surface *cairo.Surface
}

func NewWidgetPrototype() *WidgetPrototype {
	widget := WidgetPrototype{
		Sender:   event.NewSender(),
		Receiver: event.NewReceiver(),
		minSize:  image.Point{10, 10}}
	return &widget
}

func (w *WidgetPrototype) Draw() {
}

func (w *WidgetPrototype) SetArea(drawRect image.Rectangle) {
	w.area = drawRect
}

func (w *WidgetPrototype) Area() image.Rectangle {
	return w.area
}

func (w *WidgetPrototype) SetSurface(surface *cairo.Surface) {
	w.surface = surface
}

func (w *WidgetPrototype) Surface() *cairo.Surface {
	return w.surface
}

func (w *WidgetPrototype) MinSize() image.Point {
	return w.minSize
}
