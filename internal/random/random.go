package random

import (
	"image/color"
	"math/rand"

	"github.com/jlowell000/game-of-life/internal/cellularautomata"
)

func NewRandomGenerator(xMax, yMax int) *cellularautomata.CellularAutomata[color.NRGBA] {
	rColor := func() color.NRGBA {
		r, g, b := rand.Intn(255), rand.Intn(255), rand.Intn(255)
		return color.NRGBA{R: uint8(r), G: uint8(b), B: uint8(g), A: 255}
	}

	return cellularautomata.NewCellularAutomata(
		xMax, yMax,
		/* Rule Function */
		func(
			_ *cellularautomata.Cell[color.NRGBA],
			_ *cellularautomata.CellularAutomata[color.NRGBA],
		) color.NRGBA {
			return rColor()
		},
		/* Color Function */
		func(s color.NRGBA) color.NRGBA {
			return s
		},
		/* Out of Bounds Function */
		cellularautomata.NothingOOB,
		/* Initial Fill Function */
		func(
			_ *cellularautomata.Cell[color.NRGBA],
			_ *cellularautomata.CellularAutomata[color.NRGBA],
		) color.NRGBA {
			return rColor()
		},
	)
}
