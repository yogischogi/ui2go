package event

import (
	"fmt"
)

// Command represents any kind of user command, for example
// a button click.
type Command struct {
	Command string
	Sender  string
}

func (c *Command) String() string {
	return fmt.Sprintf("event.Command Id: %s Cmd: %s ", c.Sender, c.Command)
}

// DisplayRequest is an event that is sended, whenever a widget
// wants to be displayed on the screen. In most cases this would
// be to update it's representation.
type DisplayRequest struct {
	Sender string
}

func (d *DisplayRequest) String() string {
	return fmt.Sprintf("event.DisplayRequest Sender: %s", d.Sender)
}
