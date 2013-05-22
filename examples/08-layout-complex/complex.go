// Package 08-layout-complex demonstrates a slightly more complex real world example.
package main

import "code.google.com/p/ui2go/widget"

func main() {
	win := widget.NewWindow()

	topButtonBar := widget.NewDefaultContainer()
	topButtonBar.Addf("%c %c %c")

	statusArea := widget.NewDefaultContainer()
	statusArea.Addf("%c %c wrap")
	statusArea.Addf("%c %c wrap")
	statusArea.Addf("%c %c     ")

	rightButtonBar := widget.NewDefaultContainer()
	rightButtonBar.Addf("%c wrap")
	rightButtonBar.Addf("%c wrap")
	rightButtonBar.Addf("%c wrap")
	rightButtonBar.Addf("%c growy wrap", widget.NewSpacer())
	rightButtonBar.Addf("%c", statusArea)

	textFields := widget.NewDefaultContainer()
	textFields.Addf("%c growxy wrap")
	textFields.Addf("%c growy")

	win.Addf("%c spanx 2 wrap  ", topButtonBar)
	win.Addf("%c growxy      %c", textFields, rightButtonBar)

	win.Show()
	win.Run()
}
