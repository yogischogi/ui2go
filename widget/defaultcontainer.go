package widget

import (
	"code.google.com/p/ui2go/event"
	"encoding/json"
	"strings"
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

func NewDefaultContainerFromJson(jsonDef []byte) *DefaultContainer {
	c := NewDefaultContainer()
	if jsonDef == nil {
		return c
	}

	// Interpret Json and populate container.
	var jsonStruct interface{}
	err := json.Unmarshal(jsonDef, &jsonStruct)
	if err != nil {
		panic("NewContainerFromJson: Invalid JSON format.")
	}
	jsonMap := jsonStruct.(map[string]interface{})
	for key, value := range jsonMap {
		// Look for layout string.
		if key == "Layout" {
			if value, isString := value.(string); isString {
				c.interpretLayout(value)
			}
		} else {
			// Must be a definition for a Widget/Drawable.
			c.getDrawable(key, value.([]byte))
		}
	}
	return c
}

// getDrawable returns a Drawable from a name.
// getDrawable first looks in the widget registry if a component
// of the specified name exists. If not it creates a new Drawable
// from the name and the JSON definition.
//
// For example when the name is blueButton a new Button named "blue"
// is created. The json parameter contains a JSON definition.
// If no JSON definition is provided this parameter may be nil.
// If the name does not fit any widget definition it is not
// possible to create a new Drawable. In such a case nil is returned.
func (c *DefaultContainer) getDrawable(name string, jsonDef []byte) Drawable {
	var result Drawable
	for widgetType, jsonConstructor := range ConstructorRegistry {
		// Parse key name for known widget type
		if strings.HasSuffix(name, widgetType) {
			// Extract widget name
			widgetName := strings.TrimSuffix(name, widgetType)
			if widgetName == "" {
				widgetName = NewId()
			}
			// Check registry for widget and create a new
			// one if it does not yet exist.
			if ComponentRegistry[widgetName] != nil {
				result = ComponentRegistry[widgetName]
				return result
			} else {
				result := jsonConstructor(jsonDef)
				ComponentRegistry[widgetName] = result
				return result
			}
		}
	}
	return result
}

// interpretLayout interpretes a layout definition from the JSON format.
// Example: "blueButton wrap redButton".
func (c *DefaultContainer) interpretLayout(jsonLayout string) {
	jsonFields := strings.Fields(jsonLayout)
	layoutCmds := make([]string, 0)
	drawables := make([]Drawable, 0)

	// Search for widget definition in the layout string.
	for _, entry := range jsonFields {
		drawable := c.getDrawable(entry, nil)
		if drawable == nil {
			layoutCmds = append(layoutCmds, entry)
		} else {
			drawables = append(drawables, drawable)
			layoutCmds = append(layoutCmds, "%c")
		}
	}

	// Create layout
	layoutDef := strings.Join(layoutCmds, " ")
	c.Addf(layoutDef, drawables...)
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
