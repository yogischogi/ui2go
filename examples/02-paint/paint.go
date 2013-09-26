// Package 02-paint implements a simple example on how to use
// the mouse for painting on a canvas.
package main

import (
	"code.google.com/p/ui2go/event"
	"code.google.com/p/ui2go/widget"
	"github.com/skelterjohn/go.wde"
	"image"
)

// onEvent handles all events, that are sent from components
// embedded into the main window.
func onEvent(canvas *widget.Canvas, evt interface{}) {
	switch evt := evt.(type) {
	case wde.MouseDownEvent:
		if evt.Which == wde.LeftButton {
			canvas.DrawCircle(image.Point{X: evt.Where.X, Y: evt.Where.Y})
		}
	case wde.MouseDraggedEvent:
		if evt.Which == wde.LeftButton {
			canvas.DrawCircle(image.Point{X: evt.Where.X, Y: evt.Where.Y})
		}
	}
}

func main() {
	win := widget.NewWindow()
	canvas := widget.NewCanvas()
	// Layout one component in the window.
	win.Addf("%c growxy", canvas)
	// Redirect all events from the main window to the onEvent method.
	event.NewReceiverFor(win).SetEvtHandler(func(evt interface{}) { onEvent(canvas, evt) })
	win.Show()
	win.Run()
}
