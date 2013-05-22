// Package widget contains classes for widgets and the important
// main window of a program.
package widget

import (
	"code.google.com/p/ui2go/event"
	"code.google.com/p/x-go-binding/ui"
	"code.google.com/p/x-go-binding/ui/x11"
	"fmt"
	"os"
)

// Window represents a typical graphical window.
// It encapsulates a native system window.
type Window struct {
	DefaultContainer
	osWindow            ui.Window
	receiverForEmbedded event.Receiver
}

func NewWindow() *Window {
	w := new(Window)
	win, err := x11.NewWindow()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.osWindow = win
	w.DefaultContainer = *NewDefaultContainer()
	w.DefaultContainer.SetScreen(w.osWindow.Screen())
	w.DefaultContainer.SetArea(w.osWindow.Screen().Bounds())
	w.receiverForEmbedded = event.NewReceiverFor(w.DefaultContainer)
	w.receiverForEmbedded.SetEvtChanHandler(func(ec <-chan interface{}) { w.ReceiveFromEmbeddedChan(ec) })
	return w
}

// Show draws the contents of the window and flushes them
// onto the screen.
func (w *Window) Show() {
	w.DefaultContainer.Draw()
	w.osWindow.FlushImage()
}

func (w *Window) Close() {
	w.osWindow.Close()
}

// readEvents forwards native window events to the embedded container.
func (w *Window) readEvents(x11in <-chan interface{}) {
	for evt := range x11in {
		w.DefaultContainer.ReceiveEvent(evt)
	}
}

// Run just starts the main event loop.
func (w *Window) Run() {
	w.readEvents(w.osWindow.EventChan())
}

// ReceiveFromEmbeddedChan receives events from embedded components.
func (w *Window) ReceiveFromEmbeddedChan(ec <-chan interface{}) {
	for evt := range ec {
		switch ev := evt.(type) {
		case event.DisplayRequest:
			w.osWindow.FlushImage()
		case event.Command:
			if ev.Command == "Close" {
				w.osWindow.Close()
			}
		}
	}
}
