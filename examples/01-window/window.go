// Package 01-window implements an empty window as the most
// basic example.
package main

import "code.google.com/p/ui2go/widget"

func main() {
	win := widget.NewWindow()
	win.Show()
	win.Run()
}
