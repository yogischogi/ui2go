package widget

import (
	"github.com/ungerik/go-cairo"
	"image"
	"image/color"
)

// Style determines the drawing style like CSS.
type Style struct {
	FontSize float64

	Color      color.Color
	Background color.Color

	MarginTop    int
	MarginLeft   int
	MarginRight  int
	MarginBottom int

	PaddingTop    int
	PaddingLeft   int
	PaddingRight  int
	PaddingBottom int

	BorderTopColor    color.Color
	BorderLeftColor   color.Color
	BorderRightColor  color.Color
	BorderBottomColor color.Color

	BorderTopWidth    int
	BorderLeftWidth   int
	BorderRightWidth  int
	BorderBottomWidth int
}

// Size returns the size that is needed to draw the style without the inner content.
func (s *Style) Size() (dx, dy int) {
	dx = s.PaddingLeft + s.PaddingRight + s.BorderLeftWidth + s.BorderRightWidth +
		s.MarginLeft + s.MarginRight
	dy = s.PaddingTop + s.PaddingBottom + s.BorderTopWidth + s.BorderBottomWidth +
		s.MarginTop + s.MarginBottom
	return
}

// ContentPosition() returns the top left postion for content after
// the borders, paddings and margins are drawn.
func (s *Style) ContentPosition() (x, y int) {
	x = s.PaddingLeft + s.BorderLeftWidth + s.MarginLeft
	y = s.PaddingTop + s.BorderTopWidth + s.MarginTop
	return
}

func (s *Style) Draw(surface *cairo.Surface, area image.Rectangle) {
	// Outer border points
	borderXmin := float64(area.Min.X) + float64(s.MarginLeft)
	borderYmin := float64(area.Min.Y) + float64(s.MarginTop)
	borderXmax := float64(area.Max.X) - float64(s.MarginRight)
	borderYmax := float64(area.Max.Y) - float64(s.MarginBottom)

	// Inner border points
	innerXmin := borderXmin + float64(s.BorderLeftWidth)
	innerYmin := borderYmin + float64(s.BorderTopWidth)
	innerXmax := borderXmax - float64(s.BorderRightWidth)
	innerYmax := borderYmax - float64(s.BorderBottomWidth)

	// Draw top border
	surface.MoveTo(borderXmin, borderYmin)
	surface.LineTo(borderXmax, borderYmin)
	surface.LineTo(innerXmax, innerYmin)
	surface.LineTo(innerXmin, innerYmin)
	surface.ClosePath()
	surface.SetSourceRGBA(rgba(s.BorderTopColor))
	surface.Fill()

	// Draw right border
	surface.MoveTo(borderXmax, borderYmin)
	surface.LineTo(borderXmax, borderYmax)
	surface.LineTo(innerXmax, innerYmax)
	surface.LineTo(innerXmax, innerYmin)
	surface.ClosePath()
	surface.SetSourceRGBA(rgba(s.BorderRightColor))
	surface.Fill()

	// Draw bottom border
	surface.MoveTo(borderXmax, borderYmax)
	surface.LineTo(borderXmin, borderYmax)
	surface.LineTo(innerXmin, innerYmax)
	surface.LineTo(innerXmax, innerYmax)
	surface.ClosePath()
	surface.SetSourceRGBA(rgba(s.BorderBottomColor))
	surface.Fill()

	// Draw left border
	surface.MoveTo(borderXmin, borderYmax)
	surface.LineTo(borderXmin, borderYmin)
	surface.LineTo(innerXmin, innerYmin)
	surface.LineTo(innerXmin, innerYmax)
	surface.ClosePath()
	surface.SetSourceRGBA(rgba(s.BorderLeftColor))
	surface.Fill()

	// Draw background
	surface.MoveTo(innerXmin, innerYmin)
	surface.LineTo(innerXmax, innerYmin)
	surface.LineTo(innerXmax, innerYmax)
	surface.LineTo(innerXmin, innerYmax)
	surface.ClosePath()
	surface.SetSourceRGBA(rgba(s.Background))
	surface.Fill()
}
