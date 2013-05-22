package widget

import (
	"code.google.com/p/ui2go/event" // Button is a simple Button that is able to display an image.
	"code.google.com/p/x-go-binding/ui"
	"image"
	"image/color"
	"image/draw"
	"os"
)

type Button struct {
	WidgetPrototype
	Caption          string
	Command          string
	bgImage          image.Image
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
		draw.Draw(b.screen, b.area, b.bgImage, image.ZP, draw.Src)
	} else {
		greenImg := image.Uniform{C: color.RGBA{0, 255, 0, 255}}
		redImg := image.Uniform{C: color.RGBA{255, 0, 0, 255}}
		innerArea := b.area.Inset(2)
		draw.Draw(b.screen, b.area, &greenImg, image.ZP, draw.Src)
		draw.Draw(b.screen, innerArea, &redImg, image.ZP, draw.Src)
	}
}

func (b *Button) drawHighlighted() {
	if b.bgImage != nil {
		top := image.Uniform{C: color.RGBA{0, 0, 0, 255}}
		mask := image.Uniform{C: color.RGBA{0, 0, 0, 100}}
		draw.DrawMask(b.screen, b.area, &top, image.ZP, &mask, image.ZP, draw.Over)
	} else {
		greenImg := image.Uniform{C: color.RGBA{0, 255, 0, 255}}
		yellowImg := image.Uniform{C: color.RGBA{255, 255, 0, 255}}
		innerArea := b.area.Inset(2)
		draw.Draw(b.screen, b.area, &greenImg, image.ZP, draw.Src)
		draw.Draw(b.screen, innerArea, &yellowImg, image.ZP, draw.Src)
	}
}

func (b *Button) SetImage(img image.Image) {
	b.bgImage = img
}

func (b *Button) LoadImage(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
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
	switch ev := evt.(type) {
	case ui.MouseEvent:
		if ev.Buttons == 1 {
			b.isLeftButtonDown = true
			if b.isHighlighted == false {
				b.drawHighlighted()
				b.SendEvent(event.DisplayRequest{})
				b.isHighlighted = true
			}
		}
		if ev.Buttons != 1 && b.isLeftButtonDown {
			b.isLeftButtonDown = false
			if b.isHighlighted {
				b.Draw()
				b.isHighlighted = false
				b.SendEvent(event.DisplayRequest{})
			}
			cmdEvent := event.Command{Command: b.Command}
			b.SendEvent(cmdEvent)
		}
	}
}
