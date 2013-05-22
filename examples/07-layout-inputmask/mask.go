// Package 07-layout-input-mask imitates the layout of a classical input mask.
//
// It demonstrates how to use GridGaps to define the spacing between
// components.
package main

import "code.google.com/p/ui2go/widget"

func main() {
	win := widget.NewWindow()

	gaps := widget.GridGaps{
		Top:            1,
		Begin:          1,
		End:            1,
		Bottom:         1,
		BetweenColumns: 2,
		BetweenRows:    1}.Unit(widget.Cm)
	win.SetGaps(gaps)

	win.Addf("%c growx 1 %c growx 2 wrap")
	win.Addf("%c         %c         wrap")
	win.Addf("%c         %c")

	win.Show()
	win.Run()
}
