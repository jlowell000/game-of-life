package cellularautomata

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/jlowell000/game-of-life/internal/imagegenerator"
	"github.com/jlowell000/utils"
)

type Cell[T any] struct {
	X, Y      int
	IsBound   bool
	State     [2]T
	Neighbors []*Cell[T]
}

func NewCell[T any](x, y int) *Cell[T] {
	c := &Cell[T]{X: x, Y: y, IsBound: false}
	return c
}

type CellularAutomata[T any] struct {
	xMax, yMax            int
	readIndex, writeIndex int
	currentImage          *image.NRGBA
	ruleFunc              func(*Cell[T], *CellularAutomata[T]) T
	colorFunc             func(T) color.NRGBA
	oobCellFunc           func(*CellularAutomata[T], int, int, bool, bool, bool, bool) *Cell[T]
	currentState          [][]*Cell[T]
	cells                 []*Cell[T]
}

func NewCellularAutomata[T any](
	xMax, yMax int,
	ruleFunc func(*Cell[T], *CellularAutomata[T]) T,
	colorFunc func(T) color.NRGBA,
	oobCellFunc func(*CellularAutomata[T], int, int, bool, bool, bool, bool) *Cell[T],
	initFillFunc func(*Cell[T], *CellularAutomata[T]) T,
) *CellularAutomata[T] {
	ca := CellularAutomata[T]{
		xMax:         xMax,
		yMax:         yMax,
		readIndex:    0,
		writeIndex:   1,
		currentImage: imagegenerator.CreateNewImage(xMax, yMax),
		ruleFunc:     ruleFunc,
		colorFunc:    colorFunc,
	}
	ca.oobCellFunc = oobCellFunc
	ca.currentState, ca.cells = makeCellTable[T](xMax, yMax)

	ca.applyRules(ca.readIndex, initFillFunc)
	utils.ForEachWG(ca.cells, func(c *Cell[T]) { c.Neighbors = ca.getNeighborCells(c) })

	return &ca
}

func (ca *CellularAutomata[T]) GetImageFromCurrentState() image.Image {
	if ca.currentImage == nil {
		ca.currentImage = imagegenerator.CreateNewImage(ca.xMax, ca.yMax)
	}
	utils.ForEachWG(
		ca.cells,
		func(c *Cell[T]) {
			ca.currentImage.Set(
				c.X, c.Y,
				ca.colorFunc(c.State[ca.readIndex]),
			)
		},
	)
	return ca.currentImage
}

func (ca *CellularAutomata[T]) GenerateNextState() {
	ca.applyRules(ca.writeIndex, ca.ruleFunc)
	ca.swapReadWrite()
}

func (ca *CellularAutomata[T]) GetCellAt(x, y int) (c *Cell[T], err error) {
	xMinOOB, xMaxOOB, yMinOOB, yMaxOOB := ca.oob(x, y)
	if xMinOOB || xMaxOOB || yMinOOB || yMaxOOB {
		err = fmt.Errorf("{x:%d, y:%d} out of range", x, y)
	} else {
		c = ca.currentState[x][y]
	}
	return
}

func (ca *CellularAutomata[T]) GetReadWriteIndexes() (readIndex int, writeIndex int) {
	return ca.readIndex, ca.writeIndex
}

func (ca *CellularAutomata[T]) Bounding(x, y int) *Cell[T] {
	xMinOOB, xMaxOOB, yMinOOB, yMaxOOB := ca.oob(x, y)
	if (xMinOOB || xMaxOOB) || (yMinOOB || yMaxOOB) {
		return ca.oobCellFunc(ca, x, y, xMinOOB, xMaxOOB, yMinOOB, yMaxOOB)
	}
	c, _ := ca.GetCellAt(x, y)
	return c
}

func (ca *CellularAutomata[T]) oob(x, y int) (xMinOOB, xMaxOOB, yMinOOB, yMaxOOB bool) {
	vIsOob := func(v, vMin, vMax int) (bool, bool) { return v < vMin, v > vMax }
	xMinOOB, xMaxOOB = vIsOob(x, 0, ca.xMax-1)
	yMinOOB, yMaxOOB = vIsOob(y, 0, ca.yMax-1)
	return
}

func (ca *CellularAutomata[T]) applyRules(index int, ruleFunc func(*Cell[T], *CellularAutomata[T]) T) {
	utils.ForEachWG(
		ca.cells,
		func(c *Cell[T]) { c.State[index] = ruleFunc(c, ca) },
	)
}

func (ca *CellularAutomata[T]) swapReadWrite() {
	if ca.readIndex == 0 {
		ca.readIndex = 1
		ca.writeIndex = 0
	} else {
		ca.readIndex = 0
		ca.writeIndex = 1
	}
}

func (ca *CellularAutomata[T]) getNeighborCells(c *Cell[T]) (result []*Cell[T]) {
	if c.X >= 0 && c.X < ca.xMax && c.Y >= 0 && c.Y < ca.yMax {
		for i := c.X - 1; i <= c.X+1; i++ {
			for j := c.Y - 1; j <= c.Y+1; j++ {
				if !(i == c.X && j == c.Y) {
					result = append(result, ca.Bounding(i, j))
				}
			}
		}
	}
	return
}

func makeCellTable[T any](xMax, yMax int) ([][]*Cell[T], []*Cell[T]) {
	if xMax < 0 || math.MaxInt <= xMax {
		xMax = 0
	}
	if yMax < 0 || math.MaxInt <= yMax {
		yMax = 0
	}
	cellList := []*Cell[T]{}
	cellTable := make([][]*Cell[T], xMax)
	for i := range cellTable {
		cellTable[i] = make([]*Cell[T], yMax)
		for j := range cellTable[i] {
			c := NewCell[T](i, j)
			cellTable[i][j] = c
			cellList = append(cellList, c)
		}
	}
	return cellTable, cellList
}
