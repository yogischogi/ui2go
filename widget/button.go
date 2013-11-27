package widget

import (
	"code.google.com/p/ui2go/event"
	"github.com/ungerik/go-cairo"
	"github.com/ungerik/go-cairo/extimage"
	"image"
	"image/draw"
)

// Button is a simple Button that is able to display an image.
type Button struct {
	WidgetPrototype
	Caption          string
	Command          string
	bgImage          *extimage.BGRA
	isLeftButtonDown bool
	isHighlighted    bool
}

func NewButton(caption string) *Button {
	b := new(Button)
	b.WidgetPrototype = *NewWidgetPrototype()
	b.Caption = caption
	b.Command = caption
	b.SetEvtHandler(func(evt interface{}) { b.onEvent(evt) })
	return b
}

func (b *Button) Draw() {
	if b.bgImage != nil {
		x, y, _, _ := RectSize(b.area)
		imgSf := cairo.NewSurfaceFromImage(b.bgImage)
		b.surface.SetSourceSurface(imgSf, x, y)
		b.surface.Paint()
		imgSf.Destroy()
	} else {
		drawDummyWidget(b.surface, b.area)
	}
}

func (b *Button) drawHighlighted() {
	if b.bgImage != nil {
		x, y, dx, dy := RectSize(b.area)
		b.surface.Rectangle(x, y, dx, dy)
		b.surface.SetSourceRGBA(0, 0, 0, 0.3)
		b.surface.Fill()
		b.surface.Flush()
	} else {
		drawDummyWidget(b.surface, b.area)
	}
}

func (b *Button) SetImage(img image.Image) {
	switch img := img.(type) {
	case *extimage.BGRA:
		b.bgImage = img
	default:
		b.bgImage = extimage.NewBGRA(img.Bounds())
		draw.Draw(b.bgImage, img.Bounds(), img, image.ZP, draw.Src)
	}
}

func (b *Button) LoadImage(filename string) error {
	img, err := LoadImage(filename)
	if err != nil {
		return err
	}
	b.SetImage(img)
	return nil
}

func (b *Button) MinSize() image.Point {
	var result image.Point
	if b.bgImage != nil {
		result = b.bgImage.Bounds().Size()
	} else {
		result = image.Point{80, 40}
	}
	return result
}

func (b *Button) onEvent(evt interface{}) {
	switch evt := evt.(type) {
	case event.PointerEvt:
		if evt.Type == event.PointerTouchEvt {
			b.drawHighlighted()
			b.SendEvent(event.DisplayRequest{})
			b.isHighlighted = true
		} else if evt.Type == event.PointerUntouchEvt && b.isHighlighted {
			b.Draw()
			b.isHighlighted = false
			b.SendEvent(event.DisplayRequest{})
			cmdEvent := event.Command{Command: b.Command}
			b.SendEvent(cmdEvent)
		}
	}
}
