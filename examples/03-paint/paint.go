// Package 03-paint implements a slightly more sophisticated
// example for painting with the mouse.
//
// Programs like this can sometimes be found in department stores
// to delight small children. I hope you will be delighted too :-)
package main

import (
	"code.google.com/p/ui2go/event"
	"code.google.com/p/ui2go/toolbox"
	"code.google.com/p/ui2go/widget"
	"code.google.com/p/x-go-binding/ui"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

const resDir = "src/code.google.com/p/ui2go/examples/03-paint/resources"

func onCommand(evt interface{}, canvas *widget.Canvas) {
	if ev, isCommand := evt.(event.Command); isCommand {
		switch ev.Command {
		case "SmallBrush":
			canvas.SetBrushRadius(7)
		case "MediumBrush":
			canvas.SetBrushRadius(16)
		case "BigBrush":
			canvas.SetBrushRadius(27)
		case "RedBrush":
			canvas.SetBrushColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		case "GreenBrush":
			canvas.SetBrushColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		case "BlueBrush":
			canvas.SetBrushColor(color.RGBA{R: 0, G: 0, B: 255, A: 255})
		case "YellowBrush":
			canvas.SetBrushColor(color.RGBA{R: 255, G: 255, B: 0, A: 255})
		case "WhiteBrush":
			canvas.SetBrushColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		case "BlackBrush":
			canvas.SetBrushColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		case "BgRed":
			canvas.SetBackgroundColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		case "BgGreen":
			canvas.SetBackgroundColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		case "BgBlue":
			canvas.SetBackgroundColor(color.RGBA{R: 0, G: 0, B: 255, A: 255})
		case "BgYellow":
			canvas.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 0, A: 255})
		case "BgWhite":
			canvas.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		case "BgBlack":
			canvas.SetBackgroundColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		}
	}
}

func onMouseEventsFromCanvas(ec <-chan interface{}, canvas *widget.Canvas) {
	isDrawing := false
	for evt := range ec {
		if ev, isMouseEvt := evt.(ui.MouseEvent); isMouseEvt {
			location := image.Point{X: ev.Loc.X, Y: ev.Loc.Y}
			switch ev.Buttons {
			case 1:
				// left button pressed
				if isDrawing {
					canvas.LineTo(location)
				} else {
					canvas.DrawCircle(location)
					isDrawing = true
				}
			default:
				isDrawing = false
			}
		}
	}
}

func main() {
	resourcesDir, err := toolbox.FindResourcesDir(resDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	win := widget.NewWindow()

	// Button bar at the top
	topBar := widget.NewContainer()
	smallBrush := widget.NewButton("SmallBrush")
	smallBrush.LoadImage(filepath.Join(resourcesDir, "circle-small.png"))
	mediumBrush := widget.NewButton("MediumBrush")
	mediumBrush.LoadImage(filepath.Join(resourcesDir, "circle-med.png"))
	bigBrush := widget.NewButton("BigBrush")
	bigBrush.LoadImage(filepath.Join(resourcesDir, "circle-big.png"))
	whiteBrush := widget.NewButton("WhiteBrush")
	whiteBrush.LoadImage(filepath.Join(resourcesDir, "rectangle-small-white.png"))
	blackBrush := widget.NewButton("BlackBrush")
	blackBrush.LoadImage(filepath.Join(resourcesDir, "rectangle-small-black.png"))
	redBrush := widget.NewButton("RedBrush")
	redBrush.LoadImage(filepath.Join(resourcesDir, "rectangle-small-red.png"))
	greenBrush := widget.NewButton("GreenBrush")
	greenBrush.LoadImage(filepath.Join(resourcesDir, "rectangle-small-green.png"))
	blueBrush := widget.NewButton("BlueBrush")
	blueBrush.LoadImage(filepath.Join(resourcesDir, "rectangle-small-blue.png"))
	yellowBrush := widget.NewButton("YellowBrush")
	yellowBrush.LoadImage(filepath.Join(resourcesDir, "rectangle-small-yellow.png"))

	topBar.Addf("%c %c %c   ", smallBrush, mediumBrush, bigBrush)
	topBar.Addf("%c %c %c %c", redBrush, greenBrush, blueBrush, yellowBrush)
	topBar.Addf("%c %c      ", whiteBrush, blackBrush)

	// Button bar on the left
	leftBar := widget.NewContainer()
	bgWhite := widget.NewButton("BgWhite")
	bgWhite.LoadImage(filepath.Join(resourcesDir, "rectangle-big-white.png"))
	bgBlack := widget.NewButton("BgBlack")
	bgBlack.LoadImage(filepath.Join(resourcesDir, "rectangle-big-black.png"))
	bgRed := widget.NewButton("BgRed")
	bgRed.LoadImage(filepath.Join(resourcesDir, "rectangle-big-red.png"))
	bgGreen := widget.NewButton("BgGreen")
	bgGreen.LoadImage(filepath.Join(resourcesDir, "rectangle-big-green.png"))
	bgBlue := widget.NewButton("BgBlue")
	bgBlue.LoadImage(filepath.Join(resourcesDir, "rectangle-big-blue.png"))
	bgYellow := widget.NewButton("BgYellow")
	bgYellow.LoadImage(filepath.Join(resourcesDir, "rectangle-big-yellow.png"))

	leftBar.Addf("%c wrap %c wrap", bgRed, bgGreen)
	leftBar.Addf("%c wrap %c wrap", bgBlue, bgYellow)
	leftBar.Addf("%c wrap %c     ", bgWhite, bgBlack)

	canvas := widget.NewCanvas()

	// Layout components in the main window.
	// If you concentrate on the %c or the right parameters,
	// it is easy to spot the general layout idea.
	win.Addf("%c spanx 2   wrap", topBar)
	win.Addf("%c %c growx growy", leftBar, canvas)

	// Add event handlers.
	// Button commands and mouse events are handled separately.
	event.NewReceiverFor(win).SetEvtHandler(func(evt interface{}) { onCommand(evt, canvas) })
	event.NewReceiverFor(canvas).SetEvtChanHandler(func(ec <-chan interface{}) { onMouseEventsFromCanvas(ec, canvas) })

	win.Show()
	win.Run()
}
