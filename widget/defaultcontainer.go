package widget

import (
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

// ReceiveEvent receives a single event.
func (c *DefaultContainer) ReceiveEvent(evt interface{}) {
	switch evt := evt.(type) {
	case event.PointerEvt:
		// Dispatch event
		for _, widget := range c.widgets {
			area := widget.Area()
			if evt.X >= area.Min.X &&
				evt.X <= area.Max.X &&
				evt.Y >= area.Min.Y &&
				evt.Y <= area.Max.Y {
				widget.ReceiveEvent(evt)
				break
			}
		}
	default:
		c.SendEvent(evt)
	}
}

// ReceiveFromEmbedded receives events from embedded components.
func (c *DefaultContainer) ReceiveFromEmbedded(evt interface{}) {
	// events from embedded widgets: forward them
	c.SendEvent(evt)
}
