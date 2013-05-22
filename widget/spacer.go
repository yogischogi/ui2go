package widget

import (
	"image"
	"image/draw"
)

// Spacer is a Drawable object that can be used to add some space
// in a layout. Apart from that, a spacer does nothing.
type Spacer struct {
	area    image.Rectangle
	minSize image.Point
	screen  draw.Image
}

func NewSpacer() *Spacer {
	spacer := Spacer{
		// Just add some size information so that the spacer is visible.
		minSize: image.Point{10, 10}}
	return &spacer
}

func (s *Spacer) Draw() {
	// do nothing
}

func (s *Spacer) SetArea(drawRect image.Rectangle) {
	s.area = drawRect
}

func (s *Spacer) Area() image.Rectangle {
	return s.area
}

func (s *Spacer) SetScreen(screen draw.Image) {
	s.screen = screen
}

func (s *Spacer) Screen() draw.Image {
	return s.screen
}

func (s *Spacer) MinSize() image.Point {
	return s.minSize
}
