// Package 04-imageviewer implements a simple imageviewer to
// browse images in a given directory.
//
// The images are scaled to fit into the window and retain
// their aspect ratio.
package main

import (
	"code.google.com/p/ui2go/event"
	"code.google.com/p/ui2go/widget"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const resDir = "src/code.google.com/p/ui2go/examples/04-imageviewer/resources"

var (
	imageNo    int
	imageDir   string
	imageNames []string
)

func onEvent(evt interface{}, canvas *widget.Canvas) {
	if ev, isCommand := evt.(event.Command); isCommand {
		switch ev.Command {
		case "PreviousImage":
			imageNo--
			if imageNo == -1 {
				imageNo = len(imageNames) - 1
			}
			fmt.Printf("%s\n", imageNames[imageNo])
			img, err := widget.LoadImage(filepath.Join(imageDir, imageNames[imageNo]))
			if err != nil {
				fmt.Printf("Error loading image: %s\n", err)
			}
			canvas.SetBackgroundImage(img)
		case "NextImage":
			imageNo++
			if imageNo == len(imageNames) {
				imageNo = 0
			}
			fmt.Printf("%s\n", imageNames[imageNo])
			img, err := widget.LoadImage(filepath.Join(imageDir, imageNames[imageNo]))
			if err != nil {
				fmt.Printf("Error loading image: %s\n", err)
			}
			canvas.SetBackgroundImage(img)
		}
	}
}

func readFilenames(dirName string) error {
	imageNames = make([]string, 0, 100)
	dir, err := os.Open(dirName)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	extensions := map[string]bool{".gif": true, ".png": true, ".jpg": true, ".jpeg": true}
	for _, name := range files {
		if extensions[strings.ToLower(filepath.Ext(name))] == true {
			imageNames = append(imageNames, name)
		}
	}
	if len(imageNames) == 0 {
		return errors.New("No image files found.")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <imagedirectory>\n", os.Args[0])
		return
	}
	imageDir = os.Args[1]
	resourcesDir, err := widget.FindResourcesDir(resDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = readFilenames(imageDir)
	if err != nil {
		fmt.Printf("Could not read files in directory: %s\n", err)
		return
	}

	win := widget.NewWindow()
	canvas := widget.NewCanvas()

	btnPrevious := widget.NewButton("PreviousImage")
	btnPrevious.LoadImage(filepath.Join(resourcesDir, "arrow-left.png"))

	btnNext := widget.NewButton("NextImage")
	btnNext.LoadImage(filepath.Join(resourcesDir, "arrow-right.png"))

	// 2 buttons and one filler
	win.Addf("%c %c %c growx         wrap", btnPrevious, btnNext)
	win.Addf("%c spanx 3 growx growy     ", canvas)

	event.NewReceiverFor(win).SetEvtHandler(func(evt interface{}) { onEvent(evt, canvas) })

	win.Show()
	win.Run()
}
