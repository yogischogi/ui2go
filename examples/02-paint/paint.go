// Package 02-paint implements a simple example on how to use
// the mouse for painting on a canvas.
package main

import (
	"github.com/yogischogi/ui2go/event"
	"github.com/yogischogi/ui2go/widget"
	"image"
)

// onEvent handles all events, that are sent from components
// embedded into the main window.
func onEvent(canvas *widget.Canvas, evt interface{}) {
	switch evt := evt.(type) {
	case event.PointerEvt:
		switch evt.Type {
		case event.PointerTouchEvt:
			canvas.MoveTo(image.Point{X: evt.X, Y: evt.Y})
		case event.PointerMoveEvt:
			if evt.State == event.PointerStateTouch {
				canvas.LineTo(image.Point{X: evt.X, Y: evt.Y})
			}
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
