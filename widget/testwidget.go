package widget

import (
	"code.google.com/p/ui2go/event"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// TestWidget is a simple widget for testing purposes.
// It implements the Widget interface.
type TestWidget struct {
	event.Sender
	event.Receiver
	area    image.Rectangle
	minSize image.Point
	screen  draw.Image
}

func NewTestWidget() *TestWidget {
	widget := TestWidget{
		Sender:   event.NewSender(),
		Receiver: event.NewReceiver(),
		minSize:  image.Point{80, 40}}
	return &widget
}

func (w *TestWidget) Draw() {
	greenImg := image.Uniform{C: color.RGBA{0, 255, 0, 255}}
	blueImg := image.Uniform{C: color.RGBA{0, 0, 255, 255}}
	innerArea := w.area.Inset(2)
	draw.Draw(w.screen, w.area, &greenImg, image.ZP, draw.Src)
	draw.Draw(w.screen, innerArea, &blueImg, image.ZP, draw.Src)
}

func (w *TestWidget) SetArea(drawRect image.Rectangle) {
	w.area = drawRect
}

func (w *TestWidget) Area() image.Rectangle {
	return w.area
}

func (w *TestWidget) SetScreen(screen draw.Image) {
	w.screen = screen
}

func (w *TestWidget) Screen() draw.Image {
	return w.screen
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
