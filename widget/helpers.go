package widget

import (
	"errors"
	"github.com/ungerik/go-cairo"
	"github.com/ungerik/go-cairo/extimage"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

// LoadImage loads an image from a file.
// For good drawing performance it returns an *extimage.BGRA.
// extimage.BGRA can be used directly for cairo image surfaces.
func LoadImage(filename string) (*extimage.BGRA, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	result := extimage.NewBGRA(img.Bounds())
	draw.Draw(result, img.Bounds(), img, image.ZP, draw.Src)
	return result, nil
}

// FindResourcesDir tries to locate a subdirectory that resides
// within the GOPATH environment variable.
//
// If the functions finds the specified subdirectory in one of the
// GOPATH directories it returns the full path of the directory.
func FindResourcesDir(resDir string) (dir string, err error) {
	goPaths := filepath.SplitList(os.Getenv("GOPATH"))
	for _, path := range goPaths {
		resourcesDir := filepath.Join(path, resDir)
		_, err := os.Stat(resourcesDir)
		if err == nil {
			return resourcesDir, err
		}
	}
	return "", errors.New("Could not locate resources directory.")
}

// drawDummyWidget draws a simple widget for testing purposes.
func drawDummyWidget(surface *cairo.Surface, area image.Rectangle) {
	// green image
	x, y, dx, dy := RectDimensions(area)
	surface.Rectangle(x, y, dx, dy)
	surface.SetSourceRGB(0, 255, 0)
	surface.Fill()

	// blue image
	innerArea := area.Inset(2)
	x, y, dx, dy = RectDimensions(innerArea)
	surface.Rectangle(x, y, dx, dy)
	surface.SetSourceRGB(0, 0, 255)
	surface.Fill()

	surface.Flush()
}

// RectDimensions returns the dimensions of a Rectangle
// as float64 values.
func RectDimensions(area image.Rectangle) (x, y, dx, dy float64) {
	x = float64(area.Min.X)
	y = float64(area.Min.Y)
	dx = float64(area.Dx())
	dy = float64(area.Dy())
	return
}

// rgba returns the red, gree, blue and alpha values
// for a color as float64 values.
func rgba(color color.Color) (r, g, b, a float64) {
	red, green, blue, alpha := color.RGBA()
	r = float64(red)
	g = float64(green)
	b = float64(blue)
	a = float64(alpha)
	return
}
