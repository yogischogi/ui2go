// Package 06-layout-grid shows a matrix of cells that grow with the size
// of the windows, just like the Java GridLayout.
package main

import "github.com/yogischogi/ui2go/widget"

func main() {
	win := widget.NewWindow()

	win.Addf("%c growxy %c growx %c growx wrap")
	win.Addf("%c growy  %c       %c wrap")
	win.Addf("%c growy  %c       %c     ")

	win.Show()
	win.Run()
}
