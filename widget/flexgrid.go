package widget

import (
	"image"
)

// flexGridLayoutDef defines the layout constraints
// of a component within a flexGrid.
type flexGridLayoutDef struct {
	dxPixel int
	dyPixel int
	growx   int
	growy   int
	spanx   int
	spany   int
}

// stripSize defines the width or height of a single column or row in flexGrid.
type stripSize struct {
	size int
	grow int
}

// flexGrid is a helper class used to calculate the layout of growable components.
type flexGrid struct {
	// Each cell holds just a marker if it is occupied.
	cells [][]bool

	// gridPoints hold the layout information.
	// Note, that gridPoints extends cells by 1 in each dimension.
	gridPoints [][]*flexGridLayoutDef

	// positions of the cells in pixel units
	xPos []int
	yPos []int

	// Minimum sizes of rows and columns.
	colMinWidths  []stripSize
	rowMinHeights []stripSize

	// size of the grid in cell units
	size image.Point

	// current location in cell units
	cursor image.Point

	// preferred size of the whole grid in pixel units
	DefaultSizePixel image.Point

	// size of the grid in pixel units after the layout calculation
	RealSizePixel image.Point
}

func newFlexGrid() *flexGrid {
	g := new(flexGrid)
	g.cells = make([][]bool, 1)
	g.cells[0] = make([]bool, 1)
	g.gridPoints = make([][]*flexGridLayoutDef, 1)
	g.gridPoints[0] = make([]*flexGridLayoutDef, 1)
	return g
}

// Adds a new cell to the grid.
// Returns the grid position of the newly added cell.
func (g *flexGrid) AddComponent(lay flexGridLayoutDef) image.Point {
	// look for next free position
	maxPoint := image.Point{X: g.cursor.X + lay.spanx, Y: g.cursor.Y + lay.spany}
	for {
		area := image.Rectangle{Min: g.cursor, Max: maxPoint}
		if g.isAreaOccupied(area) {
			g.nextColumn()
			maxPoint.X++
		} else {
			break
		}
	}

	// calculate actual size
	if g.size.X < maxPoint.X {
		g.size.X = maxPoint.X
	}
	if g.size.Y < maxPoint.Y {
		g.size.Y = maxPoint.Y
	}
	g.adjustMatrixSize(g.size.Y, g.size.X)

	// Add layout information to grid
	endPointX := g.cursor.X + lay.spanx
	endPointY := g.cursor.Y + lay.spany

	g.gridPoints[endPointY][endPointX] = &lay

	// Mark cells as occupied
	for row := g.cursor.Y; row < endPointY; row++ {
		for column := g.cursor.X; column < endPointX; column++ {
			g.cells[row][column] = true
		}
	}
	return g.cursor
}

func (g *flexGrid) nextColumn() {
	g.cursor.X++
}

func (g *flexGrid) NextRow() {
	g.cursor.Y++
	g.cursor.X = 0
}

// isAreaOccupied determines if a rectangular area of
// the layout grid is already occupied.
func (g *flexGrid) isAreaOccupied(area image.Rectangle) bool {
	result := false
	for x := area.Min.X; x < area.Max.X; x++ {
		for y := area.Min.Y; y < area.Max.Y; y++ {
			if g.isOccupied(x, y) {
				result = true
				break
			}
		}
	}
	return result
}

// isOccupied determines if a cell is already occupied.
func (g *flexGrid) isOccupied(x int, y int) bool {
	result := false
	if x < g.size.X && y < g.size.Y {
		if g.cells[y][x] == true {
			result = true
		}
	}
	return result
}

func (g *flexGrid) MinSize() image.Point {
	// positions of the grid points
	xPoints := make([]int, g.size.X+1)
	yPoints := make([]int, g.size.Y+1)

	// x-axis
	for x := 1; x <= g.size.X; x++ {
		for y := 1; y <= g.size.Y; y++ {
			if lay := g.gridPoints[y][x]; lay != nil {
				position := xPoints[x-lay.spanx] + lay.dxPixel
				if position > xPoints[x] {
					xPoints[x] = position
				}
			}
		}
	}
	// y-axis
	for y := 1; y <= g.size.Y; y++ {
		for x := 1; x <= g.size.X; x++ {
			if lay := g.gridPoints[y][x]; lay != nil {
				position := yPoints[y-lay.spany] + lay.dyPixel
				if position > yPoints[y] {
					yPoints[y] = position
				}
			}
		}
	}
	return image.Point{X: xPoints[g.size.X], Y: yPoints[g.size.Y]}
}

func (g *flexGrid) sum(values []stripSize, from int, to int) int {
	sum := 0
	for i := from; i < to; i++ {
		sum += values[i].size
	}
	return sum
}

func (g *flexGrid) CalculatePositions() {
	// calculate grid dimensions in pixel units
	minSize := g.MinSize()
	if g.DefaultSizePixel.X > minSize.X {
		g.RealSizePixel.X = g.DefaultSizePixel.X
	} else {
		g.RealSizePixel.X = minSize.X
	}
	if g.DefaultSizePixel.Y > minSize.Y {
		g.RealSizePixel.Y = g.DefaultSizePixel.Y
	} else {
		g.RealSizePixel.Y = minSize.Y
	}

	// calculate minumum widths and heights for all rows and columns.
	g.colMinWidths = make([]stripSize, g.size.X)
	g.rowMinHeights = make([]stripSize, g.size.Y)

	// x-axis
	iLastGrow := -1
	for x := 1; x <= g.size.X; x++ {
		for y := 1; y <= g.size.Y; y++ {
			if lay := g.gridPoints[y][x]; lay != nil {
				if x-lay.spanx > iLastGrow {
					dx := lay.dxPixel - g.sum(g.colMinWidths, x-lay.spanx, x-1)
					if dx > g.colMinWidths[x-1].size {
						g.colMinWidths[x-1].size = dx
					}
				}
				if lay.growx > g.colMinWidths[x-1].grow {
					g.colMinWidths[x-1].grow = lay.growx
				}
			}
		}
		if g.colMinWidths[x-1].grow > 0 {
			iLastGrow = x - 1
		}
	}
	// y-axis
	iLastGrow = -1
	for y := 1; y <= g.size.Y; y++ {
		for x := 1; x <= g.size.X; x++ {
			if lay := g.gridPoints[y][x]; lay != nil {
				if y-lay.spany > iLastGrow {
					dy := lay.dyPixel - g.sum(g.rowMinHeights, y-lay.spany, y-1)
					if dy > g.rowMinHeights[y-1].size {
						g.rowMinHeights[y-1].size = dy
					}
				}
				if lay.growy > g.rowMinHeights[y-1].grow {
					g.rowMinHeights[y-1].grow = lay.growy
				}
			}
		}
		if g.rowMinHeights[y-1].grow > 0 {
			iLastGrow = y - 1
		}
	}

	// Sum all weights for distribution of remaining space
	dxWeight := 0
	dyWeight := 0
	dxSize := 0
	dySize := 0
	for x := 0; x < g.size.X; x++ {
		dxWeight += g.colMinWidths[x].grow
		dxSize += g.colMinWidths[x].size
	}
	for y := 0; y < g.size.Y; y++ {
		dyWeight += g.rowMinHeights[y].grow
		dySize += g.rowMinHeights[y].size
	}

	// remaining space in both dimensions
	dxFree := g.RealSizePixel.X - dxSize
	dyFree := g.RealSizePixel.Y - dySize

	// Calculate positions by eating up remaining free space.
	g.xPos = make([]int, g.size.X+1)
	g.yPos = make([]int, g.size.Y+1)
	for i := 0; i < g.size.X; i++ {
		g.xPos[i+1] = g.xPos[i] + g.colMinWidths[i].size
		if dxWeight > 0 {
			// add some of the remaining free space
			g.xPos[i+1] += dxFree * g.colMinWidths[i].grow / dxWeight
		}
	}
	for i := 0; i < g.size.Y; i++ {
		g.yPos[i+1] = g.yPos[i] + g.rowMinHeights[i].size
		if dyWeight > 0 {
			// add some of the remaining free space
			g.yPos[i+1] += dyFree * g.rowMinHeights[i].grow / dyWeight
		}
	}
	return
}

// GetPositionInPixels returns the pixel position of point within
// a flexGrid. Before this method is called the first time, the
// positions must be calculated by calling CalculatePositions().
func (g *flexGrid) GetPositionInPixels(gridPoint image.Point) image.Point {
	return image.Point{X: g.xPos[gridPoint.X], Y: g.yPos[gridPoint.Y]}
}

func (g *flexGrid) adjustMatrixSize(rows int, columns int) {
	if rows > cap(g.cells) || columns > cap(g.cells[0]) {
		var newMatrix [][]bool
		newMatrix = make([][]bool, rows)
		for row := 0; row < rows; row++ {
			newMatrix[row] = make([]bool, columns)
		}
		for row := 0; row < len(g.cells); row++ {
			copy(newMatrix[row], g.cells[row])
		}
		g.cells = newMatrix
	}
	ny := rows + 1
	nx := columns + 1
	if ny > cap(g.gridPoints) || nx > cap(g.gridPoints[0]) {
		var newGridPoints [][]*flexGridLayoutDef
		newGridPoints = make([][]*flexGridLayoutDef, ny)
		for row := 0; row < ny; row++ {
			newGridPoints[row] = make([]*flexGridLayoutDef, nx)
		}
		for row := 0; row < len(g.gridPoints); row++ {
			copy(newGridPoints[row], g.gridPoints[row])
		}
		g.gridPoints = newGridPoints
	}
}
