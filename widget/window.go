// Package widget contains classes for widgets and the important
// main window of a program.
package widget

import (
	"code.google.com/p/ui2go/event"
	"fmt"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"os"
)

// Window represents a typical graphical window.
// It usually encapsulates a native system window.
// In this version of the code it is a wde window.
type Window struct {
	DefaultContainer
	osWindow            wde.Window
	receiverForEmbedded event.Receiver
}

func NewWindow() *Window {
	w := new(Window)
	win, err := wde.NewWindow(800, 600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.osWindow = win
	w.DefaultContainer = *NewDefaultContainer()
	w.receiverForEmbedded = event.NewReceiverFor(w.DefaultContainer)
	w.receiverForEmbedded.SetEvtChanHandler(func(ec <-chan interface{}) { w.ReceiveFromEmbeddedChan(ec) })
	return w
}

// Show draws the contents of the window and makes it visible on the screen.
func (w *Window) Show() {
	w.Draw()
	w.osWindow.Show()
}

// Draw completely draws the contents of a window.
// This includes a recalculation of the position of widgets.
func (w *Window) Draw() {
	w.DefaultContainer.SetScreen(w.osWindow.Screen())
	w.DefaultContainer.SetArea(w.osWindow.Screen().Bounds())
	w.DefaultContainer.Draw()
	w.osWindow.FlushImage(w.DefaultContainer.Screen().Bounds())
}

func (w *Window) Close() {
	err := w.osWindow.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	wde.Stop()
}

// readEvents forwards native window events to the embedded container.
func (w *Window) readEvents(osWinIn <-chan interface{}) {
	for evt := range osWinIn {
		w.DefaultContainer.ReceiveEvent(evt)
	}
}

// Run just starts the main event loop.
func (w *Window) Run() {
	go w.readEvents(w.osWindow.EventChan())
	wde.Run()
}

// ReceiveFromEmbeddedChan receives events from embedded components.
func (w *Window) ReceiveFromEmbeddedChan(ec <-chan interface{}) {
	for evt := range ec {
		switch ev := evt.(type) {
		case event.DisplayRequest:
			w.osWindow.FlushImage()
		case event.Command:
			if ev.Command == "Close" {
				w.Close()
			}
		case wde.CloseEvent:
			w.Close()
		case wde.ResizeEvent:
			w.Draw()
		}
	}
}
