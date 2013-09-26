package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/skelterjohn/go.wde"
	"image"
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

// dispatchEventToWidget forwards an event to the corresponding widget.
// The widget is determined by the location of the event.
func (c *DefaultContainer) dispatchEventToWidget(location image.Point, evt interface{}) {
	for _, widget := range c.widgets {
		area := widget.Area()
		if location.X >= area.Min.X &&
			location.X <= area.Max.X &&
			location.Y >= area.Min.Y &&
			location.Y <= area.Max.Y {
			widget.ReceiveEvent(evt)
			break
		}
	}
}

// ReceiveEvent receives a single event.
func (c *DefaultContainer) ReceiveEvent(evt interface{}) {
	switch evt := evt.(type) {
	case wde.MouseDownEvent:
		c.dispatchEventToWidget(evt.Where, evt)
	case wde.MouseUpEvent:
		c.dispatchEventToWidget(evt.Where, evt)
	case wde.MouseMovedEvent:
		c.dispatchEventToWidget(evt.Where, evt)
	case wde.MouseDraggedEvent:
		c.dispatchEventToWidget(evt.Where, evt)
	case wde.MouseEnteredEvent:
		c.dispatchEventToWidget(evt.Where, evt)
	case wde.MouseExitedEvent:
		c.dispatchEventToWidget(evt.Where, evt)
	default:
		c.SendEvent(evt)
	}
}

// ReceiveFromEmbedded receives events from embedded components.
func (c *DefaultContainer) ReceiveFromEmbedded(evt interface{}) {
	// events from embedded widgets: forward them
	c.SendEvent(evt)
}
