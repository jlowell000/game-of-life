package cellularautomata

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_applyRules(t *testing.T) {
	const testSize int = 10
	tests := []struct {
		name       string
		args       CellularAutomata[int]
		expectFunc func(*Cell[int]) int
	}{
		{
			name: "initial readIndex: 0, writeIndex: 1",
			args: CellularAutomata[int]{
				xMax:       testSize,
				yMax:       testSize,
				readIndex:  0,
				writeIndex: 1,
				ruleFunc:   func(c *Cell[int], ca *CellularAutomata[int]) int { return c.X * c.Y },
			},
			expectFunc: func(c *Cell[int]) int { return c.X * c.Y },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.currentState, tt.args.cells = makeCellTable[int](tt.args.xMax, tt.args.yMax)
			tt.args.applyRules(tt.args.writeIndex, tt.args.ruleFunc)
			for _, c := range tt.args.cells {
				want := tt.expectFunc(c)
				got := c.State[tt.args.writeIndex]
				assert.Equal(
					t, want, got,
					fmt.Sprintf("applyRules() {x:%d, y:%d} got = %v, want %v", c.X, c.Y, got, want),
				)
			}

		})
	}
}

func Test_swapReadWrite(t *testing.T) {
	tests := []struct {
		name                                  string
		args                                  CellularAutomata[int]
		expectedReadIndex, expectedWriteIndex int
	}{
		{
			name:               "initial readIndex: 0, writeIndex: 1",
			args:               CellularAutomata[int]{readIndex: 0, writeIndex: 1},
			expectedReadIndex:  1,
			expectedWriteIndex: 0,
		},
		{
			name:               "initial readIndex: 1, writeIndex: 0",
			args:               CellularAutomata[int]{readIndex: 1, writeIndex: 0},
			expectedReadIndex:  0,
			expectedWriteIndex: 1,
		},
		{
			name:               "initial readIndex: nil (is treated as 0), writeIndex: 0",
			args:               CellularAutomata[int]{writeIndex: 0},
			expectedReadIndex:  1,
			expectedWriteIndex: 0,
		},
		{
			name:               "initial readIndex: 0, writeIndex: nil",
			args:               CellularAutomata[int]{readIndex: 0},
			expectedReadIndex:  1,
			expectedWriteIndex: 0,
		},
		{
			name:               "initial readIndex: nil (is treated as 0), writeIndex: 1",
			args:               CellularAutomata[int]{writeIndex: 1},
			expectedReadIndex:  1,
			expectedWriteIndex: 0,
		},
		{
			name:               "initial readIndex: 1, writeIndex: nil",
			args:               CellularAutomata[int]{readIndex: 1},
			expectedReadIndex:  0,
			expectedWriteIndex: 1,
		},
		{
			name:               "initial readIndex: nil, writeIndex: nil",
			args:               CellularAutomata[int]{},
			expectedReadIndex:  1,
			expectedWriteIndex: 0,
		},
		{
			name:               "initial readIndex: 123, writeIndex: 123",
			args:               CellularAutomata[int]{},
			expectedReadIndex:  1,
			expectedWriteIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.swapReadWrite()
			assert.Equal(
				t, tt.expectedReadIndex, tt.args.readIndex,
				fmt.Sprintf("swapReadWrite() readIndex got = %v, want %v", tt.args.readIndex, tt.expectedReadIndex),
			)
			assert.Equal(
				t, tt.expectedWriteIndex, tt.args.writeIndex,
				fmt.Sprintf("swapReadWrite() writeIndex got = %v, want %v", tt.args.writeIndex, tt.expectedWriteIndex),
			)
		})
	}
}

func Test_getNeighborCells(t *testing.T) {
	const testSize int = 10
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
			oobCellFunc:  func(_ *CellularAutomata[int], x, y int, _, _, _, _ bool) *Cell[int] { return &Cell[int]{X: x, Y: y} },
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
