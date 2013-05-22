// Package 02-paint implements a simple example on how to use
// the mouse for painting on a canvas.
package main

import (
	"code.google.com/p/ui2go/event"
	"code.google.com/p/ui2go/widget"
	"code.google.com/p/x-go-binding/ui"
	"image"
)

// onEvent handles all events, that are sent from components
// embedded into the main window.
func onEvent(canvas *widget.Canvas, evt interface{}) {
	if ev, isMouseEvt := evt.(ui.MouseEvent); isMouseEvt {
		if ev.Buttons == 1 {
			canvas.DrawCircle(image.Point{X: ev.Loc.X, Y: ev.Loc.Y})
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
