package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/ungerik/go-cairo"
	"github.com/ungerik/go-cairo/extimage"
	"image"
	"image/color"
	"image/draw"
)

type Canvas struct {
	WidgetPrototype
	backgroundColor color.Color
	backgroundImage *extimage.BGRA
	brushColor      color.Color
	brushWidth      float64
	brushPosition   image.Point
}

func NewCanvas() *Canvas {
	return &Canvas{
		WidgetPrototype: *NewWidgetPrototype(),
		backgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		brushColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
		brushWidth:      10}
}

// SetBrushWidth sets the width of the drawing brush
// in pixel units. For all drawing operations a circular
// brush is used.
func (c *Canvas) SetBrushWidth(width int) {
	c.brushWidth = float64(width)
}

func (c *Canvas) SetBrushColor(color color.Color) {
	c.brushColor = color
}

func (c *Canvas) MoveTo(p image.Point) {
	c.brushPosition = p
}

func (c *Canvas) LineTo(p image.Point) {
	r, g, b, _ := rgba(c.brushColor)
	c.surface.SetSourceRGB(r, g, b)
	c.surface.SetLineJoin(cairo.LINE_JOIN_ROUND)
	c.surface.SetLineCap(cairo.LINE_CAP_ROUND)
	c.surface.SetLineWidth(c.brushWidth)
	c.surface.MoveTo(float64(c.brushPosition.X), float64(c.brushPosition.Y))
	c.surface.LineTo(float64(p.X), float64(p.Y))
	c.surface.Stroke()
	c.surface.Flush()
	c.MoveTo(p)
	c.SendEvent(event.DisplayRequest{})
}

func (c *Canvas) SetBackgroundColor(color color.Color) {
	c.backgroundColor = color
	c.Draw()
}

func (c *Canvas) SetBackgroundImage(img image.Image) {
	switch img := img.(type) {
	case *extimage.BGRA:
		c.backgroundImage = img
	default:
		c.backgroundImage = extimage.NewBGRA(img.Bounds())
		draw.Draw(c.backgroundImage, img.Bounds(), img, image.ZP, draw.Src)
	}
	c.Draw()
}

// Draw overwrites Draw in WidgetPrototype.
func (c *Canvas) Draw() {
	// Draw background color
	r, g, b, _ := rgba(c.backgroundColor)
	x, y, dx, dy := RectSize(c.area)
	c.surface.Rectangle(x, y, dx, dy)
	c.surface.SetSourceRGB(r, g, b)
	c.surface.Fill()

	// Draw scaled background image
	if c.backgroundImage != nil {
		img := c.scaleImage(c.backgroundImage, c.area)
		imgSf := cairo.NewSurfaceFromImage(img)
		c.surface.Save()
		c.surface.SetSourceSurface(imgSf, x, y)
		c.surface.Paint()
		c.surface.Flush()
		c.surface.Restore()
		imgSf.Destroy()
	}

	/*	// Draw scaled background image
		// XXX This is ridiculously slow and seems to use tons of memory.
		// Maybe cairo calculates some high quality image.
		// Maybe this code is just wrong.
		if c.backgroundImage != nil {
			img := c.backgroundImage
			imgSf := cairo.NewSurfaceFromImage(img)
			s := c.scaleFactor(img.Bounds(), c.area.Dx(), c.area.Dy())
			drawSf := cairo.NewSurface(cairo.FORMAT_ARGB32, img.Bounds().Dx()*int(s+1), img.Bounds().Dy()*int(s+1))
			drawSf.Scale(s, s)
			drawSf.SetSourceSurface(imgSf, 0, 0)
			drawSf.Paint()
			c.surface.SetSourceSurface(drawSf, x, y)
			c.surface.Paint()
			c.surface.Flush()
			drawSf.Destroy()
			imgSf.Destroy()
		}
	*/
	c.SendEvent(event.DisplayRequest{})
}

// scaleFactors calculates the factor necessary to scale an
// image of size rect to fit into the dimensions dx, dy.
// The aspect ratio is preserved.
func (c *Canvas) scaleFactor(rect image.Rectangle, dx, dy int) float64 {
	sx := float64(dx) / float64(rect.Dx())
	sy := float64(dy) / float64(rect.Dy())
	// Aspect ratio correction
	if sx < sy {
		return sx
	} else {
		return sy
	}
}

// scaleImage scales an image so that it fits into the specified rectangle.
// The aspect ratio is preserved.
func (c *Canvas) scaleImage(img image.Image, size image.Rectangle) image.Image {
	// Calculate size
	s := c.scaleFactor(img.Bounds(), size.Dx(), size.Dy())
	outDx := float64(img.Bounds().Dx()) * s
	outDy := float64(img.Bounds().Dy()) * s
	outSize := image.Rect(0, 0, int(outDx), int(outDy))
	result := extimage.NewBGRA(outSize)

	// Set result image values
	srcX := float64(img.Bounds().Min.X)
	step := 1 / s
	for x := 0; x < outSize.Dx(); x++ {
		srcY := float64(img.Bounds().Min.Y)
		for y := 0; y < outSize.Dy(); y++ {
			result.Set(x, y, img.At(int(srcX), int(srcY)))
			srcY += step
		}
		srcX += step
	}
	return result
}

// ReceiveEvent receives a single event.
// It overwrites event.Receiver.ReceiveEvent().
func (c *Canvas) ReceiveEvent(evt interface{}) {
	// Forward mouse events.
	switch evt := evt.(type) {
	case event.PointerEvt:
		c.SendEvent(evt)
	}
}
