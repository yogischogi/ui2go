package widget

import (
	"code.google.com/p/ui2go/event"
	"image"
	"image/draw"
)

// WidgetPrototype is a simple widget that implements the Widget interface.
// It's main purpose is to be used as a parent class for other widgets.
type WidgetPrototype struct {
	event.Sender
	event.Receiver
	area    image.Rectangle
	minSize image.Point
	screen  draw.Image
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

func (w *WidgetPrototype) SetScreen(screen draw.Image) {
	w.screen = screen
}

func (w *WidgetPrototype) Screen() draw.Image {
	return w.screen
}

func (w *WidgetPrototype) MinSize() image.Point {
	return w.minSize
}
