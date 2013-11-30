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
	caption string
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
		caption:  "ui2go",
		style: Style{
			FontSize:          2 * float64(Rem),
			Color:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Background:        color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			MarginTop:         0,
			MarginLeft:        0,
			MarginRight:       0,
			MarginBottom:      0,
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
	// Draw style. Use minimal space.
	textExt := w.surface.TextExtents(w.caption)
	textSize := image.Point{X: int(textExt.Width + 0.5), Y: int(textExt.Height + 0.5)}
	size := textSize.Add(image.Pt(w.style.Size()))
	drawAreaMax := w.area.Min.Add(size)
	drawArea := image.Rectangle{Min: w.area.Min, Max: drawAreaMax}
	w.style.Draw(w.surface, drawArea)

	// Draw caption
	cx, cy := w.style.ContentPosition()
	x := float64(w.area.Min.X) + float64(cx) - float64(textExt.Xbearing)
	y := float64(w.area.Min.Y) + float64(cy) - float64(textExt.Ybearing)
	w.surface.MoveTo(x, y)
	w.surface.SetSourceRGB(0, 0, 0)
	w.surface.ShowText(w.caption)
}

func (w *TestWidget) SetArea(drawRect image.Rectangle) {
	w.area = drawRect
}

func (w *TestWidget) Area() image.Rectangle {
	return w.area
}

func (w *TestWidget) SetSurface(surface *cairo.Surface) {
	w.surface = surface
	w.surface.SetFontSize(w.style.FontSize)

	// XXX Add font face to Style
	w.surface.SelectFontFace("serif", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
}

func (w *TestWidget) Surface() *cairo.Surface {
	return w.surface
}

func (w *TestWidget) MinSize() image.Point {
	textExt := w.surface.TextExtents(w.caption)
	dx, dy := w.style.Size()
	dx += int(textExt.Width + 0.5)
	dy += int(textExt.Height + 0.5)
	return image.Point{X: dx, Y: dy}
}

// ReceiveEvent receives a single event.
// It overwrites event.Receiver.ReceiveEvent().
func (w *TestWidget) ReceiveEvent(evt interface{}) {
	fmt.Printf("TestWidget.ReceiveEvent: %v\n", evt)
	// forward event
	//w.SendEvent(evt)
}
