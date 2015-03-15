package widget

import (
	"github.com/yogischogi/ui2go/event"
	"github.com/yogischogi/ui2go/native"
	"image"
)

// Window represents a typical graphical window.
// It usually encapsulates a native system window.
// In this version of the code it is a wde window.
type Window struct {
	DefaultContainer
	nativeWin           native.Window
	receiverForEmbedded event.Receiver
}

func NewWindow() *Window {
	w := new(Window)
	win := native.NewWindow()
	w.nativeWin = win
	w.DefaultContainer = *NewDefaultContainer()
	w.DefaultContainer.SetSurface(w.nativeWin.Surface())
	w.receiverForEmbedded = event.NewReceiverFor(w.DefaultContainer)
	w.receiverForEmbedded.SetEvtChanHandler(func(ec <-chan interface{}) { w.ReceiveFromEmbeddedChan(ec) })
	return w
}

func NewWindowFromJson(jsonDef []byte) *Window {
	w := new(Window)
	win := native.NewWindow()
	w.nativeWin = win
	w.DefaultContainer = *NewDefaultContainerFromJson(jsonDef)
	w.DefaultContainer.SetSurface(w.nativeWin.Surface())
	w.receiverForEmbedded = event.NewReceiverFor(w.DefaultContainer)
	w.receiverForEmbedded.SetEvtChanHandler(func(ec <-chan interface{}) { w.ReceiveFromEmbeddedChan(ec) })
	return w
}

func GetWindow(name string) *Window {
	var result *Window
	drawable := ComponentRegistry[name]
	if drawable != nil {
		if win, ok := drawable.(*Window); ok {
			result = win
		} else {
			result = nil
		}
	}
	return result
}

// Show draws the contents of the window and makes it visible on the screen.
func (w *Window) Show() {
	w.Draw()
}

// Draw completely draws the contents of a window.
// This includes a recalculation of the position of widgets.
func (w *Window) Draw() {
	w.DefaultContainer.Draw()
	w.nativeWin.Flush()
}

func (w *Window) Close() {
	err := w.nativeWin.Close()
	if err != nil {
		panic("Could not close window.")
	}
}

// readEvents forwards native window events to the embedded container.
func (w *Window) readEvents(nativeWinIn <-chan interface{}) {
	for evt := range nativeWinIn {
		switch evt := evt.(type) {
		case event.ExposeEvt:
			w.DefaultContainer.SetArea(image.Rect(0, 0, evt.Dx, evt.Dy))
			w.Draw()
		case event.CloseEvt:
			w.Close()
		default:
			w.DefaultContainer.ReceiveEvent(evt)
		}
	}
}

// Run just starts the main event loop.
func (w *Window) Run() {
	w.readEvents(w.nativeWin.EventChan())
}

// ReceiveFromEmbeddedChan receives events from embedded components.
func (w *Window) ReceiveFromEmbeddedChan(ec <-chan interface{}) {
	for evt := range ec {
		switch ev := evt.(type) {
		case event.DisplayRequest:
			w.nativeWin.Flush()
		case event.Command:
			if ev.Command == "Close" {
				w.Close()
			}
		}
	}
}
