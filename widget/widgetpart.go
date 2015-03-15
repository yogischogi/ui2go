package widget

import (
	"github.com/ungerik/go-cairo"
	"github.com/yogischogi/ui2go/event"
	"image"
)

// WidgetPart implements most methods of the Widget interface.
// To fully implement the widget interface the Draw() method must
// be added.
// WidgetPart can be used to construct new widgets by mixing it into
// the new widget class.
// WidgetPart replaces the old WidgetPart to avoid
// overwriting methods in prototype based inheritence.
type WidgetPart struct {
	event.Sender
	event.Receiver
	area    image.Rectangle
	minSize image.Point
	surface *cairo.Surface
}

func NewWidgetPart() *WidgetPart {
	part := WidgetPart{
		Sender:   event.NewSender(),
		Receiver: event.NewReceiver(),
		minSize:  image.Point{10, 10}}
	return &part
}

func (w *WidgetPart) SetArea(drawRect image.Rectangle) {
	w.area = drawRect
}

func (w *WidgetPart) Area() image.Rectangle {
	return w.area
}

func (w *WidgetPart) SetSurface(surface *cairo.Surface) {
	w.surface = surface
}

func (w *WidgetPart) Surface() *cairo.Surface {
	return w.surface
}

func (w *WidgetPart) MinSize() image.Point {
	return w.minSize
}
