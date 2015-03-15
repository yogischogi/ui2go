package widget

import (
	"fmt"
	"github.com/ungerik/go-cairo"
	"github.com/yogischogi/ui2go/event"
	"image"
)

// BlankWidget is a simple widget that is used as a default in Combigridlayout.
// It draws a frame and implements the Widget interface.
type BlankWidget struct {
	event.Sender
	event.Receiver
	area    image.Rectangle
	minSize image.Point
	surface *cairo.Surface
}

func NewBlankWidget() *BlankWidget {
	widget := BlankWidget{
		Sender:   event.NewSender(),
		Receiver: event.NewReceiver(),
		minSize:  image.Point{80, 40}}
	return &widget
}

func (w *BlankWidget) Draw() {
	drawDummyWidget(w.surface, w.area)
}

func (w *BlankWidget) SetArea(drawRect image.Rectangle) {
	w.area = drawRect
}

func (w *BlankWidget) Area() image.Rectangle {
	return w.area
}

func (w *BlankWidget) SetSurface(surface *cairo.Surface) {
	w.surface = surface
}

func (w *BlankWidget) Surface() *cairo.Surface {
	return w.surface
}

func (w *BlankWidget) MinSize() image.Point {
	return w.minSize
}

// ReceiveEvent receives a single event.
// It overwrites event.Receiver.ReceiveEvent().
func (w *BlankWidget) ReceiveEvent(evt interface{}) {
	fmt.Printf("BlankWidget.ReceiveEvent: %v\n", evt)
	// forward event
	//w.SendEvent(evt)
}
