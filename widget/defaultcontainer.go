package widget

import (
	"code.google.com/p/x-go-binding/ui"
	//	"fmt"
	"code.google.com/p/ui2go/event"
)

// DefaultContainer is a sample implementation that
// satisfies the Container interface.
type DefaultContainer struct {
	CombiGridLayout
	event.Sender
	event.Receiver
	receiverForEmbedded event.Receiver
	widgets             []Widget
}

func NewContainer() Container {
	return NewDefaultContainer()
}

func NewDefaultContainer() *DefaultContainer {
	c := DefaultContainer{
		Sender:              event.NewSender(),
		Receiver:            event.NewReceiver(),
		receiverForEmbedded: event.NewReceiver(),
		CombiGridLayout:     *NewCombiGridLayout(),
		widgets:             make([]Widget, 0)}
	c.SetEvtHandler(func(evt interface{}) { c.ReceiveEvent(evt) })
	c.receiverForEmbedded.SetEvtHandler(func(evt interface{}) { c.ReceiveFromEmbedded(evt) })
	return &c
}

// Addf overwrites Addf in Layout.
func (c *DefaultContainer) Addf(layoutDef string, components ...Drawable) {
	for _, component := range components {
		if cmp, isWidget := component.(Widget); isWidget {
			c.widgets = append(c.widgets, cmp)
			c.receiverForEmbedded.ListenTo(cmp)
		}
	}
	c.CombiGridLayout.Addf(layoutDef, components...)
}

// onEvent receives a single event.
func (c *DefaultContainer) ReceiveEvent(evt interface{}) {
	switch ev := evt.(type) {
	case ui.MouseEvent:
		// dispatch mouse events to embedded widgets
		for _, widget := range c.widgets {
			area := widget.Area()
			if ev.Loc.X >= area.Min.X &&
				ev.Loc.X <= area.Max.X &&
				ev.Loc.Y >= area.Min.Y &&
				ev.Loc.Y <= area.Max.Y {
				widget.ReceiveEvent(ev)
				break
			}
		}
	case ui.KeyEvent:
	case ui.ConfigEvent:
	case ui.ErrEvent:
	}
}

// ReceiveFromEmbedded receives events from embedded components.
func (c *DefaultContainer) ReceiveFromEmbedded(evt interface{}) {
	// events from embedded widgets: forward them
	c.SendEvent(evt)
}
