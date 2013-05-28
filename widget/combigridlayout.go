package widget

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strconv"
	"strings"
)

// Metrics of the working environment.
// XXX Dirk: Should be set automatically and adjust other length units.
var (
	// Screen resolution in dots per inch.
	dpi int = 90
	// Viewing distance in mm.
	distanceToScreen int = 600
)

// Length units in pixel units.
var (
	Mm   int = 10 * dpi / 254
	Cm   int = 100 * dpi / 254
	Inch int = dpi
	// Recommended height of capital letters (cap height) for standard font.
	Cap int = distanceToScreen * 7 * Mm / 1000
	// Recommended height of em (about the same as type height) for standard font.
	Rem int = distanceToScreen * Mm / 100
)

type GridGaps struct {
	Top            int
	Begin          int
	End            int
	Bottom         int
	BetweenRows    int
	BetweenColumns int
}

// Unit converts grid gaps in a defined length unit to pixel units.
//
// The unit parameter is a factor with which GridGaps are multiplied
// to get the length in pixel units.
//
// Predefined are: Mm, Cm, Inch.
func (g GridGaps) Unit(unit int) GridGaps {
	g.Top = g.Top * unit
	g.Begin = g.Begin * unit
	g.End = g.End * unit
	g.Bottom = g.Bottom * unit
	g.BetweenRows = g.BetweenRows * unit
	g.BetweenColumns = g.BetweenColumns * unit
	return g
}

// LayoutDef defines the layout of a component in CombiGridLayout.
type LayoutDef struct {
	GrowX int
	GrowY int
	SpanX int
	SpanY int
}

type CombiGridEntry struct {
	component Drawable
	layout    *LayoutDef
	// position in flexGrid
	gridPos image.Point
}

// CombiGridLayout implements Layout.
// It is a layout that supports laying out components
// in a grid and docking them together.
// CombiGridLayout was heavily inspired by MigLayout http://miglayout.com/.
//
// CombiGridLayout is a series of components (Drawables) arranged is a grid
// or docked together (which also results in a grid).
// Each component may be a Container consisting of other components.
type CombiGridLayout struct {
	gaps          GridGaps
	entries       []*CombiGridEntry
	grid          *flexGrid
	area          image.Rectangle
	screen        draw.Image
	layoutChanged bool
}

func NewCombiGridLayout() *CombiGridLayout {
	l := new(CombiGridLayout)
	l.entries = make([]*CombiGridEntry, 0, 20)
	l.grid = newFlexGrid()
	return l
}

func (l *CombiGridLayout) AddWithLayout(component Drawable, layout *LayoutDef) {
	l.layoutChanged = true
	// calculate layout constraints for flexGridLayout
	var flexGridLayout flexGridLayoutDef
	if component != nil {
		component.SetScreen(l.screen)
		flexGridLayout.dxPixel = component.MinSize().X + l.gaps.BetweenColumns
		flexGridLayout.dyPixel = component.MinSize().Y + l.gaps.BetweenRows
	}
	if layout != nil {
		flexGridLayout.growx = layout.GrowX
		flexGridLayout.growy = layout.GrowY
		flexGridLayout.spanx = layout.SpanX
		flexGridLayout.spany = layout.SpanY
	}
	gridPos := l.grid.AddComponent(flexGridLayout)

	// add new component to this layout
	if component != nil {
		entry := new(CombiGridEntry)
		entry.component = component
		entry.layout = layout
		entry.gridPos = gridPos
		l.entries = append(l.entries, entry)
	}
}

// Wrap jumps to the beginning of the next line.
func (l *CombiGridLayout) Wrap() {
	l.grid.NextRow()
}

func (l *CombiGridLayout) Addf(layoutDef string, components ...Drawable) {
	// parse layoutSpec
	definitions := strings.Fields(layoutDef)

	// Markers for initial and final definition for one component
	iDef := 0
	fDef := 0

	// Special case: no component at the beginning
	if len(definitions) > 0 && definitions[0] != "%c" {
		for i, _ := range definitions {
			if definitions[i] == "%c" {
				fDef = i
				break
			}
		}
		l.addLayoutCmds(definitions[iDef:fDef])
	}

	iComponent := 0
	iDef = 0
	for iDef < len(definitions) {
		if definitions[iDef] == "%c" {
			// init component
			var component Drawable
			if iComponent < len(components) {
				component = components[iComponent]
			} else {
				component = NewTestWidget()
			}
			for fDef = iDef + 1; fDef < len(definitions); fDef++ {
				if definitions[fDef] == "%c" {
					break
				}
			}
			if iDef < fDef {
				l.addComponentWithCmds(component, definitions[iDef+1:fDef])
			} else {
				l.addComponentWithCmds(component, make([]string, 0))
			}
			iDef = fDef
			iComponent++
		}
	}
}

func (l *CombiGridLayout) addLayoutCmds(layoutCmds []string) {
	for _, cmd := range layoutCmds {
		if cmd == "wrap" {
			l.Wrap()
		}
	}
}

func (l *CombiGridLayout) addComponentWithCmds(component Drawable, layoutCmds []string) {
	layout := &LayoutDef{SpanX: 1, SpanY: 1}
	shouldWrap := false
	for i, _ := range layoutCmds {
		switch layoutCmds[i] {
		case "growx":
			layout.GrowX = 1
			// check if next layout definition is an argument for growx
			if i+1 < len(layoutCmds) {
				value, err := strconv.Atoi(layoutCmds[i+1])
				if err == nil {
					layout.GrowX = value
				}
			}
		case "growy":
			layout.GrowY = 1
			// check if next layout definition is an argument for growy
			if i+1 < len(layoutCmds) {
				value, err := strconv.Atoi(layoutCmds[i+1])
				if err == nil {
					layout.GrowY = value
				}
			}
		case "growxy":
			layout.GrowX = 1
			layout.GrowY = 1
			// check if next layout definition is an argument for growxy
			if i+1 < len(layoutCmds) {
				value, err := strconv.Atoi(layoutCmds[i+1])
				if err == nil {
					layout.GrowX = value
					layout.GrowY = value
				}
			}
		case "spanx":
			value, err := strconv.Atoi(layoutCmds[i+1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			layout.SpanX = value
		case "spany":
			value, err := strconv.Atoi(layoutCmds[i+1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			layout.SpanY = value
		case "wrap":
			shouldWrap = true
		}
	}
	l.AddWithLayout(component, layout)
	if shouldWrap {
		l.Wrap()
	}
}

func (l *CombiGridLayout) layout() {
	// use flexGrid to calculate areas for every component
	l.adjustFlexGridSize()
	l.grid.CalculatePositions()

	// adjust layout area of every component
	for _, entry := range l.entries {
		// calculate layout size for a component
		var layoutArea image.Rectangle
		layoutArea.Min = l.grid.GetPositionInPixels(entry.gridPos)
		if entry.layout != nil {
			maxPos := entry.gridPos.Add(image.Point{entry.layout.SpanX, entry.layout.SpanY})
			layoutArea.Max = l.grid.GetPositionInPixels(maxPos)
		} else {
			maxPos := entry.gridPos.Add(image.Point{1, 1})
			layoutArea.Max = l.grid.GetPositionInPixels(maxPos)
		}
		// adjust offset and calculate gaps, borders, insets, growx, growy, ...

		// Offset by layout area position on the screen
		componentArea := layoutArea.Add(l.area.Min)

		// Offset by combigrid layout border
		componentArea = componentArea.Add(image.Point{X: l.gaps.Begin, Y: l.gaps.Top})

		// remove row and column gaps
		componentArea.Max = componentArea.Max.Sub(image.Point{X: l.gaps.BetweenColumns, Y: l.gaps.BetweenRows})

		entry.component.SetArea(componentArea)
		//entry.component.SetArea(layoutArea.Add(l.area.Min))
	}
	l.area.Max = l.area.Min.Add(l.grid.RealSizePixel)
	l.layoutChanged = false
}

func (l *CombiGridLayout) SetArea(drawRect image.Rectangle) {
	l.area = drawRect
	//l.grid.DefaultSizePixel.X = l.area.Dx()
	//l.grid.DefaultSizePixel.Y = l.area.Dy()
	l.layoutChanged = true
}

func (l *CombiGridLayout) adjustFlexGridSize() {
	l.grid.DefaultSizePixel.X = l.area.Dx() - l.gaps.Begin - l.gaps.End + l.gaps.BetweenColumns
	l.grid.DefaultSizePixel.Y = l.area.Dy() - l.gaps.Top - l.gaps.Bottom + l.gaps.BetweenRows
}

func (l *CombiGridLayout) Area() image.Rectangle {
	return l.area
}

func (l *CombiGridLayout) Draw() {
	//	backgroundImg := image.Uniform{C: color.RGBA{220, 220, 220, 255}}
	//	draw.Draw(l.screen, l.area, &backgroundImg, image.ZP, draw.Src)

	greyImg := image.Uniform{C: color.RGBA{220, 220, 220, 255}}
	redImg := image.Uniform{C: color.RGBA{255, 0, 0, 255}}
	innerArea := l.area.Inset(2)
	draw.Draw(l.screen, l.area, &redImg, image.ZP, draw.Src)
	draw.Draw(l.screen, innerArea, &greyImg, image.ZP, draw.Src)

	if l.layoutChanged {
		l.layout()
	}
	for _, entry := range l.entries {
		entry.component.Draw()
	}
}

func (l *CombiGridLayout) SetScreen(screen draw.Image) {
	l.screen = screen
	for _, entry := range l.entries {
		entry.component.SetScreen(screen)
	}
}

func (l *CombiGridLayout) Screen() draw.Image {
	return l.screen
}

func (l *CombiGridLayout) MinSize() image.Point {
	l.adjustFlexGridSize()
	flexSize := l.grid.MinSize()
	dx := flexSize.X + l.gaps.Begin + l.gaps.End
	dy := flexSize.Y + l.gaps.Top + l.gaps.Bottom
	return image.Point{X: dx, Y: dy}
}

func (l *CombiGridLayout) SetGaps(gaps GridGaps) {
	l.gaps = gaps
}

func (l *CombiGridLayout) Gaps() GridGaps {
	return l.gaps
}
