package widget

import (
	"code.google.com/p/ui2go/event"
	"fmt"
	"github.com/ungerik/go-cairo"
	"image"
	"image/color"
)

// TestWidget is a simple widget for testing purposes.
// It implements the Widget interface.
type TestWidget struct {
	event.Sender
	event.Receiver
	style   Style
	area    image.Rectangle
	minSize image.Point
	surface *cairo.Surface
}

func NewTestWidget() *TestWidget {
	return &TestWidget{
		Sender:   event.NewSender(),
		Receiver: event.NewReceiver(),
		minSize:  image.Point{80, 40},
		style: Style{
			Color:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Background:        color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			MarginTop:         10,
			MarginLeft:        10,
			MarginRight:       10,
			MarginBottom:      10,
			PaddingTop:        10,
			PaddingLeft:       10,
			PaddingRight:      10,
			PaddingBottom:     10,
			BorderTopColor:    color.NRGBA{R: 200, G: 0, B: 0, A: 255},
			BorderLeftColor:   color.NRGBA{R: 0, G: 200, B: 0, A: 255},
			BorderRightColor:  color.NRGBA{R: 0, G: 255, B: 255, A: 255},
			BorderBottomColor: color.NRGBA{R: 200, G: 200, B: 0, A: 255},
			BorderTopWidth:    10,
			BorderLeftWidth:   10,
			BorderRightWidth:  10,
			BorderBottomWidth: 10}}
}

func (w *TestWidget) Draw() {
	w.style.Draw(w.surface, w.area)
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
	dx, dy := w.style.Size()
	return image.Point{X: dx, Y: dy}
}

// ReceiveEvent receives a single event.
// It overwrites event.Receiver.ReceiveEvent().
func (w *TestWidget) ReceiveEvent(evt interface{}) {
	fmt.Printf("TestWidget.ReceiveEvent: %v\n", evt)
	// forward event
	//w.SendEvent(evt)
}
