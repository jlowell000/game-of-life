package cellularautomata

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getNeighborCells(t *testing.T) {
	const testSize int = 20
	type args struct {
		ca   CellularAutomata[int]
		cell Cell[int]
	}
	type test struct {
		name   string
		args   args
		result []*Cell[int]
	}

	makeCellularAutomata := func(x, y int) CellularAutomata[int] {
		table, list := makeCellTable[int](x, y)
		ca := CellularAutomata[int]{
			xMax: x, yMax: y,
			currentState: table,
			cells:        list,
			oobCellFunc:  func(x, y int, _, _, _, _ bool) *Cell[int] { return &Cell[int]{X: x, Y: y} },
		}
		return ca
	}

	makeResults := func(x, y int) (result []*Cell[int]) {
		for i := x - 1; i <= x+1; i++ {
			for j := y - 1; j <= y+1; j++ {
				if !(i == x && j == y) {
					result = append(result, NewCell[int](i, j))
				}
			}
		}
		return
	}

	makeTest := func(x, y int) test {
		return test{
			name: fmt.Sprintf("x %d y %d", x, y),
			args: args{
				ca:   makeCellularAutomata(testSize, testSize),
				cell: Cell[int]{X: x, Y: y},
			},
			result: makeResults(x, y),
		}
	}

	makeTests := func() (tests []test) {
		for i := 0; i < testSize; i++ {
			for j := 0; j < testSize; j++ {
				tests = append(tests, makeTest(i, j))
			}
		}
		return
	}

	tests := makeTests()
	tests = append(tests,
		test{
			name: "cell not in space then empty list",
			args: args{
				ca:   makeCellularAutomata(testSize, testSize),
				cell: Cell[int]{X: 1000, Y: 1000},
			},
			result: []*Cell[int]{},
		},
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.ca.getNeighborCells(&tt.args.cell)
			assert.Equal(t, len(got), len(tt.result), fmt.Sprintf("getNeighborCells() lengths got = %v, want %v", len(got), len(tt.result)))
			for _, v := range got {
				assert.True(t,
					slices.ContainsFunc(tt.result, func(c *Cell[int]) bool { return c.X == v.X && c.Y == v.Y }),
					fmt.Sprintf("getNeighborCells() got[v] = %v is not in %v", v, tt.result),
				)
			}
		})
	}
}

func Test_makeCellTable(t *testing.T) {
	type args struct {
		xMax int
		yMax int
	}

	type results struct {
		table    [][]*Cell[int]
		cellList []*Cell[int]
	}

	makeResults := func(x, y int) results {
		cellTable := make([][]*Cell[int], x)
		cellList := []*Cell[int]{}

		for i := range cellTable {
			cellTable[i] = make([]*Cell[int], y)
			for j := range cellTable[i] {
				c := NewCell[int](i, j)
				cellTable[i][j] = c
				cellList = append(cellList, c)
			}
		}
		return results{
			table:    cellTable,
			cellList: cellList,
		}
	}

	tests := []struct {
		name   string
		args   args
		result results
	}{
		{
			name:   "x 5 y 5",
			args:   args{xMax: 5, yMax: 5},
			result: makeResults(5, 5),
		},
		{
			name:   "x 100 y 1",
			args:   args{xMax: 100, yMax: 1},
			result: makeResults(100, 1),
		},
		{
			name:   "x 1 y 100",
			args:   args{xMax: 1, yMax: 100},
			result: makeResults(1, 100),
		},
		{
			name:   "x 0 y 0",
			args:   args{xMax: 0, yMax: 0},
			result: makeResults(0, 0),
		},
		{
			name:   "x 0 y 5",
			args:   args{xMax: 0, yMax: 5},
			result: makeResults(0, 5),
		},
		{
			name:   "x 5 y 0",
			args:   args{xMax: 5, yMax: 0},
			result: makeResults(5, 0),
		},
		{
			name:   "x 5 y -5->0",
			args:   args{xMax: 5, yMax: -5},
			result: makeResults(5, 0),
		},
		{
			name:   "x -5->0 y 5",
			args:   args{xMax: -5, yMax: 5},
			result: makeResults(0, -5),
		},
		{
			name:   "x -5->0 y -5->0",
			args:   args{xMax: -5, yMax: -5},
			result: makeResults(0, 0),
		},
		{
			name:   "x math.MaxInt y 1",
			args:   args{xMax: math.MaxInt, yMax: 1},
			result: makeResults(0, 1),
		},
		{
			name:   "x 1 y math.MaxInt",
			args:   args{xMax: 1, yMax: math.MaxInt},
			result: makeResults(1, 0),
		},
		{
			name:   "x math.MaxInt y math.MaxInt",
			args:   args{xMax: math.MaxInt, yMax: math.MaxInt},
			result: makeResults(0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := makeCellTable[int](tt.args.xMax, tt.args.yMax)
			assert.Equal(t, tt.result.table, got, fmt.Sprintf("makeCellTable() got = %v, want %v", got, tt.result.table))
			assert.Equal(t, tt.result.cellList, got1, "makeCellTable() got1 = %v, want %v", got1, tt.result.cellList)
		})
	}
}
