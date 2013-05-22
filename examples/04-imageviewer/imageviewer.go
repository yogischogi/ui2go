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
			err := canvas.LoadImage(filepath.Join(imageDir, imageNames[imageNo]))
			if err != nil {
				fmt.Printf("Error loading image: %s\n", err)
			}
		case "NextImage":
			imageNo++
			if imageNo == len(imageNames) {
				imageNo = 0
			}
			fmt.Printf("%s\n", imageNames[imageNo])
			err := canvas.LoadImage(filepath.Join(imageDir, imageNames[imageNo]))
			if err != nil {
				fmt.Printf("Error loading image: %s\n", err)
			}
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
	resourcesDir := filepath.Join(os.Getenv("GOPATH"), "src/code.google.com/p/ui2go/resources")

	err := readFilenames(imageDir)
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

	win.Addf("%c %c                  wrap", btnPrevious, btnNext)
	win.Addf("%c spanx 2 growx growy     ", canvas)

	event.NewReceiverFor(win).SetEvtHandler(func(evt interface{}) { onEvent(evt, canvas) })

	win.Show()
	win.Run()
}
