// Package 10-style demonstrates widget style properties.
package main

import "github.com/yogischogi/ui2go/widget"

func main() {
	win := widget.NewWindow()

	w1 := widget.NewTestWidget()
	w2 := widget.NewTestWidget()
	w3 := widget.NewTestWidget()
	w4 := widget.NewTestWidget()

	win.Addf("%c growxy %c growx wrap", w1, w2)
	win.Addf("%c growy  %c           ", w3, w4)

	win.Show()
	win.Run()
}
