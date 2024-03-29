package cellularautomata

import (
	"image"
	"image/color"

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
	ReadIndex, WriteIndex int
	currentImage          *image.NRGBA
	ruleFunc              func(*Cell[T], *CellularAutomata[T]) T
	colorFunc             func(T) color.NRGBA
	oobCellFunc           func(int, int, bool, bool, bool, bool) *Cell[T]
	currentState          [][]*Cell[T]
	cells                 []*Cell[T]
}

func NewCellularAutomata[T any](
	xMax, yMax int,
	ruleFunc func(*Cell[T], *CellularAutomata[T]) T,
	colorFunc func(T) color.NRGBA,
	oobCellFunc func(*CellularAutomata[T]) func(int, int, bool, bool, bool, bool) *Cell[T],
	initFillFunc func(*Cell[T], *CellularAutomata[T]) T,
) *CellularAutomata[T] {
	ca := CellularAutomata[T]{
		xMax:         xMax,
		yMax:         yMax,
		ReadIndex:    0,
		WriteIndex:   1,
		currentImage: imagegenerator.CreateNewImage(xMax, yMax),
		ruleFunc:     ruleFunc,
		colorFunc:    colorFunc,
	}
	ca.oobCellFunc = oobCellFunc(&ca)
	ca.currentState, ca.cells = makeCellTable[T](xMax, yMax)

	ca.applyRules(ca.ReadIndex, initFillFunc)
	utils.ForEachWG(ca.cells, func(c *Cell[T]) { c.Neighbors = ca.getNeighborCells(c) })

	return &ca
}

func (ca *CellularAutomata[T]) GetImageFromCurrentState() image.Image {
	utils.ForEachWG(
		ca.cells,
		func(c *Cell[T]) {
			ca.currentImage.Set(
				c.X, c.Y,
				ca.colorFunc(c.State[ca.ReadIndex]),
			)
		},
	)
	return ca.currentImage
}

func (ca *CellularAutomata[T]) GenerateNextState() {
	ca.applyRules(ca.WriteIndex, ca.ruleFunc)
	ca.swapReadWrite()
}

func (ca *CellularAutomata[T]) GetCellAt(x, y int) *Cell[T] {
	return ca.currentState[x][y]
}

func (ca *CellularAutomata[T]) Bounding(x, y int) *Cell[T] {
	vIsOob := func(v, vMin, vMax int) (bool, bool) { return v < vMin, v > vMax }
	xMinOOB, xMaxOOB := vIsOob(x, 0, ca.xMax-1)
	yMinOOB, yMaxOOB := vIsOob(y, 0, ca.yMax-1)
	if xMinOOB || xMaxOOB || yMinOOB || yMaxOOB {
		return ca.oobCellFunc(x, y, xMinOOB, xMaxOOB, yMinOOB, yMaxOOB)
	}
	return ca.GetCellAt(x, y)
}

func (ca *CellularAutomata[T]) applyRules(index int, ruleFunc func(*Cell[T], *CellularAutomata[T]) T) {
	utils.ForEachWG(
		ca.cells,
		func(c *Cell[T]) { c.State[index] = ruleFunc(c, ca) },
	)
}

func (ca *CellularAutomata[T]) swapReadWrite() {
	if ca.ReadIndex == 0 {
		ca.ReadIndex = 1
		ca.WriteIndex = 0
	} else {
		ca.ReadIndex = 0
		ca.WriteIndex = 1
	}
}

func (ca *CellularAutomata[T]) getNeighborCells(c *Cell[T]) (result []*Cell[T]) {
	for i := c.X - 1; i <= c.X+1; i++ {
		for j := c.Y - 1; j <= c.Y+1; j++ {
			if !(i == c.X && j == c.Y) {
				result = append(result, ca.Bounding(i, j))
			}
		}
	}
	return
}

func makeCellTable[T any](xMax, yMax int) ([][]*Cell[T], []*Cell[T]) {
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
