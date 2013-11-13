// +build linux

package native

/*
#cgo pkg-config: xcb
#cgo pkg-config: cairo
#include <stdlib.h>
#include <xcb/xcb.h>
#include <cairo/cairo-xcb.h>

// getData0 is a helper function to extract the first
// 32 bit value out of msgData.
// msgData is a C union that is represented as an array
// of bytes in Go. So we use this function because
// access in C is much simpler.
xcb_atom_t getData0(xcb_client_message_data_t msgData){
	return msgData.data32[0];
}
*/
import "C"

import (
	"code.google.com/p/ui2go/event"
	"github.com/ungerik/go-cairo"
	"strconv"
	"unsafe"
)

// Window dimensions
const (
	winDx = 800
	winDy = 600
)

// init sets the NewWindow function to the window system
// specific implementation.
func init() {
	NewWindow = NewXWindow
}

// XWindow implements the Window interface for X-Windows.
// X-Window communication is done by using XCB directly.
//
// XCB resources:
//	https://en.wikipedia.org/wiki/XCB
//	http://static.usenix.org/publications/library/proceedings/als01/full_papers/massey/massey.pdf
//	http://xcb.freedesktop.org/
type XWindow struct {
	con     *C.xcb_connection_t
	screen  *C.xcb_screen_t
	win     C.xcb_window_t
	surface *C.cairo_surface_t
	context *C.cairo_t

	// Used to handle Close Window messages
	delReply *C.xcb_intern_atom_reply_t

	// Wrapper for surface and context
	cairoSurface *cairo.Surface
}

func NewXWindow() Window {
	x := new(XWindow)
	x.con = C.xcb_connect(nil, nil)
	x.screen = C.xcb_setup_roots_iterator(C.xcb_get_setup(x.con)).data
	x.win = (C.xcb_window_t)(C.xcb_generate_id(x.con))

	// Create graphics context
	foreground := (C.xcb_gcontext_t)(C.xcb_generate_id(x.con))
	drawable := (C.xcb_drawable_t)(x.screen.root)
	var evtMask C.uint32_t = C.XCB_GC_FOREGROUND | C.XCB_GC_GRAPHICS_EXPOSURES
	var values [2]C.uint32_t
	values[0] = x.screen.black_pixel
	values[1] = 0
	C.xcb_create_gc(x.con, foreground, drawable, evtMask, &values[0])

	// Create window
	evtMask = C.XCB_CW_BACK_PIXEL | C.XCB_CW_EVENT_MASK
	values[0] = x.screen.white_pixel
	values[1] = C.XCB_EVENT_MASK_EXPOSURE | C.XCB_EVENT_MASK_BUTTON_PRESS |
		C.XCB_EVENT_MASK_BUTTON_RELEASE | C.XCB_EVENT_MASK_POINTER_MOTION |
		C.XCB_EVENT_MASK_ENTER_WINDOW | C.XCB_EVENT_MASK_LEAVE_WINDOW |
		C.XCB_EVENT_MASK_KEY_PRESS | C.XCB_EVENT_MASK_KEY_RELEASE

	C.xcb_create_window(x.con, // Connection
		C.XCB_COPY_FROM_PARENT, // depth (same as root)
		x.win,         // window Id
		x.screen.root, // parent window
		0, 0,          // x, y
		winDx, winDy, // width, height
		0, // border_width
		C.XCB_WINDOW_CLASS_INPUT_OUTPUT, // class
		x.screen.root_visual,            // visual
		evtMask, &values[0])             // masks
	C.xcb_map_window(x.con, x.win)

	// Register event for window closing
	var protoCookie C.xcb_intern_atom_cookie_t
	var protoMsg = C.CString("WM_PROTOCOLS")
	protoCookie = C.xcb_intern_atom(x.con, 1, 12, protoMsg)
	var protoReply *C.xcb_intern_atom_reply_t
	protoReply = C.xcb_intern_atom_reply(x.con, protoCookie, nil)
	C.free(unsafe.Pointer(protoMsg))

	var delCookie C.xcb_intern_atom_cookie_t
	var delMsg = C.CString("WM_DELETE_WINDOW")
	delCookie = C.xcb_intern_atom(x.con, 0, 16, delMsg)
	x.delReply = C.xcb_intern_atom_reply(x.con, delCookie, nil)
	C.free(unsafe.Pointer(delMsg))

	C.xcb_change_property(x.con,
		C.XCB_PROP_MODE_REPLACE,
		x.win,
		protoReply.atom, 4, 32, 1,
		unsafe.Pointer(&(x.delReply.atom)))

	C.xcb_flush(x.con)
	return x
}

func readEvents(win *XWindow, out chan<- interface{}) {
	var evt *C.xcb_generic_event_t
	var et = newEventTranslator(win)
	for {
		evt = C.xcb_wait_for_event(win.con)
		if evt == nil {
			close(out)
			break
		}
		outEvt := et.translateEvent(evt)
		out <- outEvt
		C.free(unsafe.Pointer(evt))
		if _, ok := outEvt.(event.CloseEvt); ok {
			close(out)
			break
		}
	}
}

func (x *XWindow) EventChan() <-chan interface{} {
	evtChan := make(chan interface{})
	go readEvents(x, evtChan)
	return evtChan
}

func (x *XWindow) Flush() {
	C.xcb_flush(x.con)
}

func (x *XWindow) Close() error {
	C.cairo_device_finish(x.context)
	C.cairo_surface_finish(x.surface)
	x.cairoSurface = nil

	C.xcb_disconnect(x.con)
	C.xcb_flush(x.con)
	return nil
}

func (x *XWindow) Surface() *cairo.Surface {
	if x.surface == nil {
		x.surface, x.context = x.newSurface()
		x.cairoSurface = cairo.NewSurfaceFromC(x.surface, x.context)
	}
	return x.cairoSurface
}

func (x *XWindow) newSurface() (*C.cairo_surface_t, *C.cairo_t) {
	// Determine visual (color mapping)
	var visual *C.xcb_visualtype_t
	var depthIter C.xcb_depth_iterator_t
	depthIter = C.xcb_screen_allowed_depths_iterator(x.screen)
	for ; depthIter.rem != 0; C.xcb_depth_next(&depthIter) {
		var visualIter C.xcb_visualtype_iterator_t
		visualIter = C.xcb_depth_visuals_iterator(depthIter.data)
		for ; visualIter.rem != 0; C.xcb_visualtype_next(&visualIter) {
			if x.screen.root_visual == visualIter.data.visual_id {
				visual = visualIter.data
				break
			}
		}
	}
	if visual == nil {
		panic("Could not determine XCB visualtype.")
	}

	// Create surface for drawing
	var surface *C.cairo_surface_t
	surface = C.cairo_xcb_surface_create(x.con, (C.xcb_drawable_t)(x.win), visual, winDx, winDy)
	if surface == nil {
		panic("Could not create Cairo XCB surface.")
	}
	var context *C.cairo_t = C.cairo_create(surface)
	if C.cairo_status(context) != 0 {
		panic("Could not create Cairo context.")
	}
	return surface, context
}

// eventTranslator translates window events to ui2go events.
// It also does some low level window management, that should
// be hidden from the Window interface.
type eventTranslator struct {
	win      *XWindow
	pointerX int
	pointerY int
}

func newEventTranslator(win *XWindow) *eventTranslator {
	return &eventTranslator{win: win}
}

func (et *eventTranslator) translateEvent(evt *C.xcb_generic_event_t) event.Event {
	switch evt.response_type &^ 0x80 {

	case C.XCB_MOTION_NOTIFY:
		evtState := event.PointerStateNone
		evtType := event.PointerMoveEvt
		ev := (*C.xcb_motion_notify_event_t)(unsafe.Pointer(evt))
		if ev.state&C.XCB_BUTTON_MASK_1 > 0 {
			evtState = event.PointerStateTouch
		}
		et.pointerX = int(ev.event_x)
		et.pointerY = int(ev.event_y)
		return event.PointerEvt{Evt: event.Evt{SenderId: "Mouse"},
			Type:  evtType,
			State: evtState,
			X:     et.pointerX, Y: et.pointerY}

	case C.XCB_BUTTON_PRESS:
		// XXX currently same behaviour for all buttons
		//ev := (*C.xcb_button_press_event_t)(unsafe.Pointer(evt))
		return event.PointerEvt{Evt: event.Evt{SenderId: "Mouse"},
			Type:  event.PointerTouchEvt,
			State: event.PointerStateTouch,
			X:     et.pointerX,
			Y:     et.pointerY}

	case C.XCB_BUTTON_RELEASE:
		// XXX currently same behaviour for all buttons
		//ev := (*C.xcb_button_release_event_t)(unsafe.Pointer(evt))
		return event.PointerEvt{Evt: event.Evt{SenderId: "Mouse"},
			Type:  event.PointerUntouchEvt,
			State: event.PointerStateNone,
			X:     et.pointerX,
			Y:     et.pointerY}

	case C.XCB_EXPOSE:
		// Send every time when the window is exposed to the user.
		ev := (*C.xcb_expose_event_t)(unsafe.Pointer(evt))
		dx, dy := int(ev.width), int(ev.height)
		C.cairo_xcb_surface_set_size(et.win.surface, C.int(dx), C.int(dy))
		return event.ExposeEvt{Evt: event.Evt{SenderId: "Window"}, Dx: dx, Dy: dy}

	case C.XCB_CONFIGURE_NOTIFY:
		// Send when some attributes of the window change.
		ev := (*C.xcb_configure_notify_event_t)(unsafe.Pointer(evt))
		dx, dy := int(ev.width), int(ev.height)
		C.cairo_xcb_surface_set_size(et.win.surface, C.int(dx), C.int(dy))
		return event.ConfigEvt{Evt: event.Evt{SenderId: "Window"}, Dx: dx, Dy: dy}

	case C.XCB_CLIENT_MESSAGE:
		// Special messages, for example when the window is closed.
		ev := (*C.xcb_client_message_event_t)(unsafe.Pointer(evt))
		if C.getData0(ev.data) == et.win.delReply.atom {
			return event.CloseEvt{Evt: event.Evt{SenderId: "Window"}}
		} else {
			return event.Evt{SenderId: strconv.Itoa(int(evt.response_type))}
		}

	default:
		return event.Evt{SenderId: strconv.Itoa(int(evt.response_type))}
	}
}
