package gameoflife

import (
	"image/color"

	"github.com/jlowell000/game-of-life/internal/cellularautomata"
	"github.com/jlowell000/game-of-life/internal/imagegenerator"
	"github.com/jlowell000/utils"
)

var (
	black = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	white = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

func NewGameOfLifeGenerator(xMax, yMax int) *cellularautomata.CellularAutomata[bool] {
	return cellularautomata.NewCellularAutomata(
		xMax, yMax,
		/* Rule Function */
		func(
			c *cellularautomata.Cell[bool],
			ca *cellularautomata.CellularAutomata[bool],
		) bool {
			readIndex, _ := ca.GetReadWriteIndexes()
			numberOfLiveNeighbors := 0
			utils.ForEach(
				c.Neighbors,
				func(nc *cellularautomata.Cell[bool]) {
					if nc.IsBound {
						newBoundCell := ca.Bounding(nc.X, nc.Y)
						nc.State = newBoundCell.State
					}
					if nc.State[readIndex] {
						numberOfLiveNeighbors++
					}
				},
			)

			if c.State[readIndex] {
				return numberOfLiveNeighbors >= 2 && numberOfLiveNeighbors <= 3
			} else {
				return numberOfLiveNeighbors == 3
			}
		},
		/* Color Function */
		func(s bool) color.NRGBA {
			if s {
				return white
			} else {
				return black
			}
		},
		/* Out of Bounds Function */
		cellularautomata.OverflowOOB,
		/* Initial Fill Function */
		func(_ *cellularautomata.Cell[bool], _ *cellularautomata.CellularAutomata[bool]) bool {
			return imagegenerator.RandomCoin() == 1
		},
	)
}
