// Package event provides classes and interfaces to implement
// event senders and receivers.
//
// There are also some types of events, but these only serve
// as examples and are not thought through.
package event

import (
	"fmt"
)

type PointerEvtType int

const (
	PointerMoveEvt PointerEvtType = iota
	PointerTouchEvt
	PointerUntouchEvt
)

type PointerState int

const (
	PointerStateNone PointerState = iota
	PointerStateTouch
)

// Evt is an implementation of the event interface.
// It's main purpose is to simplify the creation of other
// types of events.
type Evt struct {
	SenderId string
}

func (evt Evt) Sender() string {
	return evt.SenderId
}

// PointerEvt represents a typical pointer event.
// The idea is that a programmer should not worry about
// what kind of pointer device is being used.
// When thinking of pointers the mouse immediately
// comes to the mind of most people, but there
// are many other kinds of pointer devices:
// 	* Human finger
// 	* Foot mouse
// 	* Graphics tablet
//	* Trackball
// 	* Roll bar mouse
//	* Eye tracker
//	* Remote control for smart TV
//	* Electronic gun in combat simulation area.
// Unfortunately there are always some clever programmers
// around who try to utilize mouse events for special purposes.
// This often messes up user experience for non mouse users.
// The result may be some poor pilot shooting down a jumbo jet
// by accident.
type PointerEvt struct {
	Evt
	Type  PointerEvtType
	State PointerState
	X     int
	Y     int
}

type ConfigEvt struct {
	Evt
	Dx int
	Dy int
}

// An ExposeEvt is sent when the window is shown to
// the user (open window, resize, show after hide).
type ExposeEvt struct {
	Evt
	Dx int
	Dy int
}

// A CloseEvt is sent before the window is closed,
// so that a user may take some action before it is
// closed finally.
type CloseEvt struct {
	Evt
}

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
