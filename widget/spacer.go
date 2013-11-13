package widget

import (
	"github.com/ungerik/go-cairo"
	"image"
)

// Spacer is a Drawable object that can be used to add some space
// in a layout. Apart from that, a spacer does nothing.
type Spacer struct {
	area    image.Rectangle
	minSize image.Point
	surface *cairo.Surface
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

func (s *Spacer) SetSurface(surface *cairo.Surface) {
	s.surface = surface
}

func (s *Spacer) Surface() *cairo.Surface {
	return s.surface
}

func (s *Spacer) MinSize() image.Point {
	return s.minSize
}
