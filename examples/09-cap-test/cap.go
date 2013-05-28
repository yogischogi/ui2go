// Package 09-cap-test tests the Cap unit, the proposed height for capital letters.
//
// The gaps between the rows should be the height of a capital letter.
// Cap size depends on the viewing distance, which is defined
// in combigridlayout.go (distanceToScreen).
package main

import "code.google.com/p/ui2go/widget"

func main() {
	win := widget.NewWindow()

	gaps := widget.GridGaps{
		Top:            1,
		Bottom:         1,
		BetweenRows:    1}.Unit(widget.Cap)
	win.SetGaps(gaps)

	win.Addf("%c growx  wrap")
	win.Addf("%c        wrap")
	win.Addf("%c")

	win.Show()
	win.Run()
}
