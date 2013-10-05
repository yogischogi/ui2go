package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/skelterjohn/go.wde"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// Abs returns the absolute value of an integer.
func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// plainRectMask is an image to be used as a mask
// when drawing a rectangle.
type plainRectMask image.Rectangle

func (p plainRectMask) ColorModel() color.Model {
	return color.Alpha16Model
}

func (p plainRectMask) Bounds() image.Rectangle {
	return image.Rectangle(p)
}

func (p plainRectMask) At(x, y int) color.Color {
	return color.Opaque
}

// rectMask is a drawing mask for a rotated rectangle,
// used to draw thick lines.
type rectMask struct {
	// Offset is the difference between area.Min and the beginning of the line.
	Offset image.Point
	// Length of the line
	length float64
	// Thickness of the line
	width float64
	// Bounding box of the line
	area image.Rectangle
	// a1 is the "zero point" (top left) of the rotated rectangle, describing the line.
	a1 image.Point
	// sin(a)
	sina float64
	// cos(a)
	cosa float64
}

func newRectMask(p1, p2 image.Point, r int) *rectMask {
	if p1.X > p2.X {
		temp := p1
		p1 = p2
		p2 = temp
	}
	mask := new(rectMask)

	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	mask.length = math.Sqrt(dx*dx + dy*dy)
	mask.sina = dy / mask.length
	mask.cosa = dx / mask.length
	mask.width = float64(2 * r)

	offsetY := float64(r) * mask.cosa
	offsetX := math.Abs(float64(r) * mask.sina)

	if dy > 0 {
		mask.a1 = image.Point{X: int(2*offsetX + 0.5), Y: 0}
		mask.Offset = image.Point{X: int(offsetX + 0.5), Y: int(offsetY + 0.5)}
	} else {
		mask.Offset = image.Point{X: int(offsetX + 0.5), Y: int(-dy + offsetY + 0.5)}
		mask.a1 = image.Point{X: 0, Y: int(-dy + 0.5)}
	}
	dxx := dx + 2*offsetX
	dyy := math.Abs(dy) + 2*offsetY
	mask.area = image.Rect(0, 0, int(dxx+0.5), int(dyy+0.5))
	return mask
}

func (m rectMask) ColorModel() color.Model {
	return color.Alpha16Model
}

func (m rectMask) Bounds() image.Rectangle {
	return m.area
}

func (m rectMask) At(x, y int) color.Color {
	// Transform point to rectangle coordinates.
	v := image.Point{X: x, Y: y}.Sub(m.a1)
	vLen := math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
	sinb := float64(v.Y) / vLen
	cosb := float64(v.X) / vLen

	// dx, dy in rectangle coordinates
	dx := vLen * (cosb*m.cosa + sinb*m.sina)
	dy := vLen * (sinb*m.cosa - cosb*m.sina)
	result := color.Transparent
	if dx >= 0 && dx < m.length && dy >= 0 && dy < m.width {
		result = color.Opaque
	}
	return result
}

type Canvas struct {
	WidgetPrototype
	backgroundImg image.Uniform
	brushImg      image.Uniform
	brushRadius   int
	brushPosition image.Point
	brushMask     image.Image
}

func NewCanvas() *Canvas {
	c := new(Canvas)
	c.WidgetPrototype = *NewWidgetPrototype()
	c.backgroundImg = image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}
	c.brushImg = image.Uniform{C: color.RGBA{R: 0, G: 0, B: 0, A: 255}}
	c.brushRadius = 5
	c.minSize = image.Point{100, 100}
	c.brushMask = c.newCircleMask(c.brushRadius)
	return c
}

// SetBrushRadius sets the radius of the drawing brush
// in pixel units. For all drawing operations a circular
// brush is used.
func (c *Canvas) SetBrushRadius(radius int) {
	c.brushRadius = radius
	c.brushMask = c.newCircleMask(radius)
}

func (c *Canvas) SetBrushColor(color color.Color) {
	c.brushImg = image.Uniform{C: color}
}

func (c *Canvas) MoveTo(p image.Point) {
	c.brushPosition = p
}

func (c *Canvas) newCircleMask(r int) image.Image {
	size := image.Rect(0, 0, r*2, r*2)
	mask := image.NewAlpha16(size)
	draw.Draw(mask, size, &image.Uniform{C: color.Transparent}, image.ZP, draw.Src)
	for y := 0; y <= r; y++ {
		xy := int(math.Sqrt(float64(r*r-y*y)) + 0.5)
		for x := r - xy; x <= r+xy; x++ {
			mask.Set(x, r-y, color.Opaque)
			mask.Set(x, r+y, color.Opaque)
		}
	}
	return mask
}

func (c *Canvas) DrawCircle(p image.Point) {
	c.MoveTo(p)
	offset := image.Point{X: c.brushRadius, Y: c.brushRadius}
	target := c.brushMask.Bounds().Add(p).Sub(offset)
	draw.DrawMask(c.screen, target, &c.brushImg, image.ZP, c.brushMask, image.ZP, draw.Over)
	c.SendEvent(event.DisplayRequest{})
}

func (c *Canvas) drawRectangle(mask image.Image, p image.Point, offset image.Point) {
	target := mask.Bounds().Add(p).Sub(offset)
	draw.DrawMask(c.screen, target, &c.brushImg, image.ZP, mask, image.ZP, draw.Over)
	c.SendEvent(event.DisplayRequest{})
}

func (c *Canvas) LineTo(p image.Point) {
	fPos := p
	iPos := c.brushPosition
	dx := fPos.X - iPos.X
	dy := fPos.Y - iPos.Y
	c.DrawCircle(iPos)
	switch {
	case fPos == iPos:
		// Nothing to do.
	case fPos.X == iPos.X:
		min := image.Point{X: 0, Y: 0}
		max := image.Point{X: c.brushRadius * 2, Y: Abs(dy)}
		rectMask := plainRectMask{Min: min, Max: max}
		if dy > 0 {
			c.drawRectangle(rectMask, iPos, image.Point{X: c.brushRadius, Y: 0})
		} else {
			c.drawRectangle(rectMask, iPos, image.Point{X: c.brushRadius, Y: -dy})
		}
		c.DrawCircle(fPos)
	case fPos.Y == iPos.Y:
		min := image.Point{X: 0, Y: 0}
		max := image.Point{X: Abs(dx), Y: c.brushRadius * 2}
		rectMask := plainRectMask{Min: min, Max: max}
		if dx > 0 {
			c.drawRectangle(rectMask, iPos, image.Point{X: 0, Y: c.brushRadius})
		} else {
			c.drawRectangle(rectMask, iPos, image.Point{X: -dx, Y: c.brushRadius})
		}
		c.DrawCircle(fPos)
	default:
		// Rotated rectangle
		rectMask := newRectMask(iPos, fPos, c.brushRadius)
		if dx > 0 {
			c.drawRectangle(rectMask, iPos, rectMask.Offset)
		} else {
			offset := rectMask.Offset.Sub(image.Point{X: dx, Y: dy})
			c.drawRectangle(rectMask, iPos, offset)
		}
		c.DrawCircle(fPos)
	}
}

func (c *Canvas) SetBackgroundColor(color color.Color) {
	c.backgroundImg = image.Uniform{C: color}
	c.Draw()
}

func (c *Canvas) Clear() {
	draw.Draw(c.screen, c.area, &c.backgroundImg, image.ZP, draw.Src)
}

// Draw overwrites Draw in WidgetPrototype.
func (c *Canvas) Draw() {
	c.Clear()
	c.SendEvent(event.DisplayRequest{})
}

// scaleImage scales an image to the specified size.
// The aspect ratio is preserved.
func (c *Canvas) scaleImage(img image.Image, size image.Point) image.Image {
	src := img.Bounds()
	dst := image.Rect(0, 0, size.X, size.Y)
	// scale factor
	var s float32
	sx := float32(src.Dx()) / float32(dst.Dx())
	sy := float32(src.Dy()) / float32(dst.Dy())
	if sx > sy {
		s = sx
	} else {
		s = sy
	}
	// aspect ratio correction
	dx := int(float32(src.Dx()) / s)
	dy := int(float32(src.Dy()) / s)

	result := image.NewRGBA(dst)
	draw.Draw(c.screen, c.area, &c.backgroundImg, image.ZP, draw.Src)

	srcX := float32(src.Min.X)
	for x := 0; x < dx; x++ {
		srcY := float32(src.Min.Y)
		for y := 0; y < dy; y++ {
			result.Set(x, y, img.At(int(srcX), int(srcY)))
			srcY += s
		}
		srcX += s
	}
	return result
}

// scaleImage scales an image to the specified size.
// The aspect ratio is preserved.
func (c *Canvas) scaleImage1(img image.Image, size image.Point) image.Image {
	src := img.Bounds()
	dst := image.Rect(0, 0, size.X, size.Y)
	// scale factor
	var s float32
	sx := float32(src.Dx()) / float32(dst.Dx())
	sy := float32(src.Dy()) / float32(dst.Dy())
	if sx > sy {
		s = sx
	} else {
		s = sy
	}
	result := image.NewRGBA(dst)
	draw.Draw(c.screen, c.area, &c.backgroundImg, image.ZP, draw.Src)
	for x := dst.Min.X; x < dst.Max.X; x++ {
		for y := dst.Min.Y; y < dst.Max.Y; y++ {
			srcX := src.Min.X + int(float32(x-dst.Min.X)*s)
			srcY := src.Min.Y + int(float32(y-dst.Min.Y)*s)
			result.Set(x, y, img.At(srcX, srcY))
		}
	}
	return result
}

func (c *Canvas) LoadImage(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	img = c.scaleImage(img, c.area.Size())
	draw.Draw(c.screen, c.area, img, image.ZP, draw.Src)
	c.SendEvent(event.DisplayRequest{})
	return nil
}

// Draw Image draws an image onto the canvas.
func (c *Canvas) DrawImage(img image.Image) {
	draw.Draw(c.screen, c.area, img, image.ZP, draw.Src)
	c.SendEvent(event.DisplayRequest{})
}

// ReceiveEvent receives a single event.
// It overwrites event.Receiver.ReceiveEvent().
func (c *Canvas) ReceiveEvent(evt interface{}) {
	// Forward mouse events.
	switch evt := evt.(type) {
	case wde.MouseDownEvent, wde.MouseUpEvent, wde.MouseMovedEvent,
		wde.MouseDraggedEvent, wde.MouseEnteredEvent, wde.MouseExitedEvent:
		c.SendEvent(evt)
	}
}
