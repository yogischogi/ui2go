// Package 11-json-layout-border creates a window with border layout from a JSON definition.
package main

import "code.google.com/p/ui2go/widget"

func main() {
	gui := "{\"Layout\": \"" +
		"TopButton spanx 3                          wrap " +
		"LeftButton MiddleButton growxy RightButton wrap " +
		"BottomButton spanx 3                            " +
		"\" }"
	win := widget.NewWindowFromJson([]byte(gui))
	win.Show()
	win.Run()
}
