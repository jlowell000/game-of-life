package cellularautomata

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"testing"

	"github.com/jlowell000/game-of-life/internal/imagegenerator"
	"github.com/jlowell000/utils"
	"github.com/stretchr/testify/assert"
)

func Test_GetImageFromCurrentState(t *testing.T) {
	const testSize int = 10
	type test struct {
		name string
		args CellularAutomata[int]
	}
	testColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	tests := []test{
		{
			name: "GetImageFromCurrentState no inital image",
			args: CellularAutomata[int]{
				xMax:       testSize,
				yMax:       testSize,
				readIndex:  0,
				writeIndex: 1,
				colorFunc: func(i int) color.NRGBA {
					return testColor
				},
			},
		},
		{
			name: "GetImageFromCurrentState inital image",
			args: CellularAutomata[int]{
				xMax:         testSize,
				yMax:         testSize,
				readIndex:    0,
				writeIndex:   1,
				currentImage: imagegenerator.CreateNewImage(testSize, testSize),
				colorFunc: func(i int) color.NRGBA {
					return testColor
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.currentState, tt.args.cells = makeCellTable[int](tt.args.xMax, tt.args.yMax)
			want := imagegenerator.CreateNewImage(tt.args.xMax, tt.args.yMax)
			for y := want.Bounds().Min.Y; y < want.Bounds().Max.Y; y++ {
				for x := want.Bounds().Min.X; x < want.Bounds().Max.X; x++ {
					want.Set(x, y, testColor)
				}
			}
			got := tt.args.GetImageFromCurrentState()
			assert.Equal(
				t, want, got,
				"GetImageFromCurrentState() generated image not equal",
			)
		})
	}
}

func Test_GenerateNextState(t *testing.T) {
	const testSize int = 10
	type test struct {
		name string
		args CellularAutomata[int]
	}
	ruleFuncCount := 0
	tests := []test{
		{
			name: "GenerateNextState",
			args: CellularAutomata[int]{
				xMax:       testSize,
				yMax:       testSize,
				readIndex:  0,
				writeIndex: 1,
				ruleFunc: func(c *Cell[int], ca *CellularAutomata[int]) int {
					ruleFuncCount++
					return 0
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.currentState, tt.args.cells = makeCellTable[int](tt.args.xMax, tt.args.yMax)
			tt.args.GenerateNextState()
			assert.True(
				t, ruleFuncCount > 0,
				fmt.Sprintf("GenerateNextState() ruleFuncCount got = %v, want greater than 0", ruleFuncCount),
			)
			assert.Equal(
				t, 1, tt.args.readIndex,
				fmt.Sprintf("GenerateNextState() readIndex got = %v, want %v", tt.args.readIndex, 1),
			)
			assert.Equal(
				t, 0, tt.args.writeIndex,
				fmt.Sprintf("GenerateNextState() writeIndex got = %v, want %v", tt.args.writeIndex, 0),
			)
			ruleFuncCount = 0
		})
	}
}

func Test_GetCellAt(t *testing.T) {
	const testSize int = 10
	type test struct {
		name string
		args CellularAutomata[int]
	}
	tests := []test{
		{
			name: "Bounding",
			args: CellularAutomata[int]{
				xMax:        testSize,
				yMax:        testSize,
				readIndex:   0,
				writeIndex:  1,
				oobCellFunc: LiveOOB,
			},
		},
	}

	assert := func(tt test, c *Cell[int]) {
		xMinOOB, xMaxOOB, yMinOOB, yMaxOOB := tt.args.oob(c.X, c.Y)
		got, gotErr := tt.args.GetCellAt(c.X, c.Y)
		if xMinOOB || xMaxOOB || yMinOOB || yMaxOOB {
			want := fmt.Errorf("{x:%d, y:%d} out of range", c.X, c.Y)
			assert.Equal(
				t, gotErr, want,
				fmt.Sprintf("GetCellAt() {x:%d, y:%d} got = %v, want %v", c.X, c.Y, got, want),
			)
		} else {
			want := tt.args.currentState[c.X][c.Y]
			assert.Equal(
				t, got, want,
				fmt.Sprintf("GetCellAt() {x:%d, y:%d} got = %v, want %v", c.X, c.Y, got, want),
			)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.currentState, tt.args.cells = makeCellTable[int](tt.args.xMax, tt.args.yMax)
			utils.ForEachWG(tt.args.cells, func(c *Cell[int]) { c.Neighbors = tt.args.getNeighborCells(c) })
			for _, c := range tt.args.cells {
				assert(tt, c)
				for _, n := range c.Neighbors {
					assert(tt, n)
				}
			}

		})
	}
}

func Test_GetReadWriteIndexes(t *testing.T) {
	type test struct {
		name string
		args CellularAutomata[int]
	}
	tests := []test{
		{
			name: "GetReadWriteIndexes",
			args: CellularAutomata[int]{
				readIndex:  0,
				writeIndex: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotW := tt.args.GetReadWriteIndexes()
			assert.Equal(
				t, tt.args.readIndex, gotR,
				fmt.Sprintf("GetReadWriteIndexes() readIndex got = %v, want %v", gotR, tt.args.readIndex),
			)
			assert.Equal(
				t, tt.args.writeIndex, gotW,
				fmt.Sprintf("GetReadWriteIndexes() writeIndex got = %v, want %v", gotW, tt.args.writeIndex),
			)
		})
	}
}

func Test_oob(t *testing.T) {
	const testSize int = 10
	type test struct {
		name string
		args CellularAutomata[int]
	}
	tests := []test{
		{
			name: "Bounding",
			args: CellularAutomata[int]{
				xMax:        testSize,
				yMax:        testSize,
				readIndex:   0,
				writeIndex:  1,
				oobCellFunc: LiveOOB,
			},
		},
	}

	assert := func(tt test, c *Cell[int]) {
		expectFunc := func(c *Cell[int], tt test) (xMinOOB, xMaxOOB, yMinOOB, yMaxOOB bool) {
			if c.X < 0 {
				xMinOOB = true
				xMaxOOB = false
			} else if c.X > tt.args.xMax-1 {
				xMinOOB = false
				xMaxOOB = true
			} else {
				xMinOOB = false
				xMaxOOB = false
			}
			if c.Y < 0 {
				yMinOOB = true
				yMaxOOB = false
			} else if c.Y > tt.args.yMax-1 {
				yMinOOB = false
				yMaxOOB = true
			} else {
				yMinOOB = false
				yMaxOOB = false
			}

			return
		}
		wantXMinOOB, wantXMaxOOB, wantYMinOOB, wantYMaxOOB := expectFunc(c, tt)
		gotXMinOOB, gotXMaxOOB, gotYMinOOB, gotYMaxOOB := tt.args.oob(c.X, c.Y)
		assert.Equal(
			t, gotXMinOOB, wantXMinOOB,
			fmt.Sprintf("oob() xMinOOB {x:%d, y:%d} got = %v, want %v", c.X, c.Y, gotXMinOOB, wantXMinOOB),
		)
		assert.Equal(
			t, wantXMaxOOB, gotXMaxOOB,
			fmt.Sprintf("oob() xMaxOOB {x:%d, y:%d} got = %v, want %v", c.X, c.Y, gotXMaxOOB, wantXMaxOOB),
		)
		assert.Equal(
			t, wantYMinOOB, gotYMinOOB,
			fmt.Sprintf("oob() yMinOOB {x:%d, y:%d} got = %v, want %v", c.X, c.Y, gotYMinOOB, wantYMinOOB),
		)
		assert.Equal(
			t, wantYMaxOOB, gotYMaxOOB,
			fmt.Sprintf("oob() yMaxOOB {x:%d, y:%d} got = %v, want %v", c.X, c.Y, gotYMaxOOB, wantYMaxOOB),
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.currentState, tt.args.cells = makeCellTable[int](tt.args.xMax, tt.args.yMax)
			utils.ForEachWG(tt.args.cells, func(c *Cell[int]) { c.Neighbors = tt.args.getNeighborCells(c) })
			for _, c := range tt.args.cells {
				assert(tt, c)
				for _, n := range c.Neighbors {
					assert(tt, n)
				}
			}

		})
	}
}

func Test_Bounding(t *testing.T) {
	const testSize int = 10
	type test struct {
		name       string
		args       CellularAutomata[int]
		expectFunc func(*Cell[int], test) int
	}
	tests := []test{
		{
			name: "Bounding",
			args: CellularAutomata[int]{
				xMax:        testSize,
				yMax:        testSize,
				readIndex:   0,
				writeIndex:  1,
				oobCellFunc: LiveOOB,
			},
			expectFunc: func(c *Cell[int], tt test) int {
				xMinOOB, xMaxOOB, yMinOOB, yMaxOOB := tt.args.oob(c.X, c.Y)
				if (xMinOOB || xMaxOOB) || (yMinOOB || yMaxOOB) {
					return 1
				}
				return 0
			},
		},
	}

	assert := func(tt test, c *Cell[int]) {
		want := tt.expectFunc(c, tt)
		got := tt.args.Bounding(c.X, c.Y).State[tt.args.readIndex]
		assert.Equal(
			t, want, got,
			fmt.Sprintf("Bounding() {x:%d, y:%d} got = %v, want %v", c.X, c.Y, got, want),
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.currentState, tt.args.cells = makeCellTable[int](tt.args.xMax, tt.args.yMax)
			utils.ForEachWG(tt.args.cells, func(c *Cell[int]) { c.Neighbors = tt.args.getNeighborCells(c) })
			for _, c := range tt.args.cells {
				assert(tt, c)
				for _, n := range c.Neighbors {
					assert(tt, n)
				}
			}

		})
	}
}

func Test_applyRules(t *testing.T) {
	const testSize int = 10
	tests := []struct {
		name       string
		args       CellularAutomata[int]
		expectFunc func(*Cell[int]) int
	}{
		{
			name: "func is called on all cells",
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
