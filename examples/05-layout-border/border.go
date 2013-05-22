// Package 05-layout-border contains an example for a classical border layout.
package main

import "code.google.com/p/ui2go/widget"

func main() {
	win := widget.NewWindow()

	win.Addf("%c spanx 3      wrap")
	win.Addf("%c %c growxy %c wrap")
	win.Addf("%c spanx 3          ")

	win.Show()
	win.Run()
}
