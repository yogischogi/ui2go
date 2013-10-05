// Package 10-cairo-test uses the go-cairo binding.
//
// In future versions of ui2go possibly all drawing operations
// will be done by using some sort of Cairo binding.
// Note that drawing operations in this example are very
// slow because they use the image interface and there are
// a lot of hidden image conversion routines at work.
//
// Also note that go-cairo requires some extra installation
// besides go get github.com/ungerik/go-cairo.
// Please visit https://github.com/ungerik/go-cairo for
// additional installation information and the Cairo homepage
// at http://cairographics.org/ for more documentation.
package main

import (
	"code.google.com/p/ui2go/widget"
	"github.com/ungerik/go-cairo"
)

func main() {
	win := widget.NewWindow()
	canvas := widget.NewCanvas()
	win.Addf("%c growxy", canvas)
	win.Show()

	img := canvas.Screen()
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, img.Bounds().Dx(), img.Bounds().Dy())
	surface.SelectFontFace("serif", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	surface.SetFontSize(64.0)
	surface.SetSourceRGB(1.0, 1.0, 1.0)
	surface.MoveTo(50, 100)
	surface.ShowText("Hello ui2go!")
	canvas.DrawImage(surface.GetImage())
	surface.Finish()

	win.Run()
}
