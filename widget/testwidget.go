package widget

import (
	"code.google.com/p/ui2go/event"
	"fmt"
	"github.com/ungerik/go-cairo"
	"image"
)

// TestWidget is a simple widget for testing purposes.
// It implements the Widget interface.
type TestWidget struct {
	event.Sender
	event.Receiver
	area    image.Rectangle
	minSize image.Point
	surface *cairo.Surface
}

func NewTestWidget() *TestWidget {
	widget := TestWidget{
		Sender:   event.NewSender(),
		Receiver: event.NewReceiver(),
		minSize:  image.Point{80, 40}}
	return &widget
}

func (w *TestWidget) Draw() {
	drawDummyWidget(w.surface, w.area)
}

func (w *TestWidget) SetArea(drawRect image.Rectangle) {
	w.area = drawRect
}

func (w *TestWidget) Area() image.Rectangle {
	return w.area
}

func (w *TestWidget) SetSurface(surface *cairo.Surface) {
	w.surface = surface
}

func (w *TestWidget) Surface() *cairo.Surface {
	return w.surface
}

func (w *TestWidget) MinSize() image.Point {
	return w.minSize
}

// ReceiveEvent receives a single event.
// It overwrites event.Receiver.ReceiveEvent().
func (w *TestWidget) ReceiveEvent(evt interface{}) {
	fmt.Printf("TestWidget.ReceiveEvent: %v\n", evt)
	// forward event
	//w.SendEvent(evt)
}
