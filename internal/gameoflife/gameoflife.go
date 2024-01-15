package gameoflife

import (
	"image"
	"image/color"
	"slices"

	"github.com/jlowell000/game-of-life/internal/imagegenerator"
	"github.com/jlowell000/utils"
)

var (
	black = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	white = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

type Cell struct {
	State bool
}

type GameOfLifeGenerator struct {
	XMax, YMax            int
	ReadIndex, WriteIndex int
	CurrentImage          *image.NRGBA
	CurrentState          [][][]*Cell
	Points                []*image.Point
}

func CreateInitalState(xMax, yMax int, initialLivePoints []image.Point) GameOfLifeGenerator {
	gen := GameOfLifeGenerator{
		XMax:       xMax,
		YMax:       yMax,
		ReadIndex:  0,
		WriteIndex: 1,
	}
	gen.CurrentState = gen.makeCellTables()
	gen.Points = gen.makePointsArray()
	gen.CurrentImage = imagegenerator.MakeNRGBASpace(gen.XMax, gen.YMax)

	if len(initialLivePoints) > 0 {
		utils.ForEach(
			gen.Points,
			func(p *image.Point) {
				gen.CurrentState[gen.ReadIndex][p.X][p.Y] = &Cell{
					State: slices.Contains(initialLivePoints, *p),
				}

			},
		)
	} else {
		utils.ForEach(
			gen.Points,
			func(p *image.Point) {
				gen.CurrentState[gen.ReadIndex][p.X][p.Y] = &Cell{
					State: imagegenerator.RandomCoin(),
				}
			},
		)
	}

	utils.ForEach(
		gen.Points,
		func(p *image.Point) {
			gen.CurrentState[gen.WriteIndex][p.X][p.Y] = &Cell{
				State: false,
			}
		},
	)

	return gen
}

func (r *GameOfLifeGenerator) GetImageFromState() image.Image {
	utils.ForEachWG(
		r.Points,
		func(p *image.Point) {
			r.CurrentImage.Set(p.X, p.Y, r.color(r.CurrentState[r.ReadIndex][p.X][p.Y].State))
		},
	)
	return r.CurrentImage
}

func (r *GameOfLifeGenerator) GenerateNextState() {
	utils.ForEachWG(
		r.Points,
		func(p *image.Point) {
			r.CurrentState[r.WriteIndex][p.X][p.Y].State = r.applyRules(p)
		},
	)
	r.swapReadWrite()
}

func (r *GameOfLifeGenerator) GetPoints() []*image.Point {
	return r.Points
}

func (r *GameOfLifeGenerator) applyRules(p *image.Point) bool {
	neighborCells := r.getNeighborCells(p)
	numberOfLiveNeighbors := len(
		utils.Filter(
			neighborCells,
			func(c *Cell) bool { return c.State },
		),
	)

	if r.CurrentState[r.ReadIndex][p.X][p.Y].State {
		if numberOfLiveNeighbors >= 2 && numberOfLiveNeighbors <= 3 {
			return true
		} else {
			return false
		}
	} else {
		if numberOfLiveNeighbors == 3 {
			return true
		} else {
			return false
		}
	}
}

func (r *GameOfLifeGenerator) getNeighborCells(p *image.Point) (result []*Cell) {
	for i := p.X - 1; i <= p.X+1; i++ {
		for j := p.Y - 1; j <= p.Y+1; j++ {
			if !(i == p.X && j == p.Y) {
				result = append(
					result,
					r.boundedCellAt(image.Point{X: i, Y: j}),
				)
			}
		}
	}
	return
}

func (r *GameOfLifeGenerator) boundedCellAt(p image.Point) *Cell {
	x, y := 0, 0
	if p.X < 0 {
		x = r.XMax - 1
	} else if p.X >= r.XMax {
		x = 0
	} else {
		x = p.X
	}
	if p.Y < 0 {
		y = r.YMax - 1
	} else if p.Y >= r.YMax {
		y = 0
	} else {
		y = p.Y
	}

	return r.CurrentState[r.ReadIndex][x][y]
}

func (r *GameOfLifeGenerator) makeCellTables() [][][]*Cell {
	cellTables := make([][][]*Cell, 2)
	for i := range cellTables {
		cellTables[i] = make([][]*Cell, r.XMax)
		for j := range cellTables[i] {
			cellTables[i][j] = make([]*Cell, r.YMax)
		}
	}
	return cellTables
}

func (r *GameOfLifeGenerator) makePointsArray() (result []*image.Point) {
	for y := 0; y < r.YMax; y++ {
		for x := 0; x < r.XMax; x++ {
			result = append(result, &image.Point{X: x, Y: y})
		}
	}
	return
}

func (r *GameOfLifeGenerator) swapReadWrite() {
	if r.ReadIndex == 0 {
		r.ReadIndex = 1
		r.WriteIndex = 0
	} else {
		r.ReadIndex = 0
		r.WriteIndex = 1
	}
}

func (r *GameOfLifeGenerator) color(s bool) color.NRGBA {
	if s {
		return white
	} else {
		return black
	}
}
