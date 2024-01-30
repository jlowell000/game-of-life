package cellularautomata

import "github.com/jlowell000/game-of-life/internal/imagegenerator"

/* wrapping oob */
func OverflowOOB[T any](ca *CellularAutomata[T], x, y int, xMinOOB, xMaxOOB, yMinOOB, yMaxOOB bool) *Cell[T] {
	f := func(v, vMin, vMax int, minOOB, maxOOB bool) int {
		if minOOB {
			v = vMax
		} else if maxOOB {
			v = vMin
		}
		return v
	}
	c, _ := ca.GetCellAt(
		f(x, 0, ca.xMax-1, xMinOOB, xMaxOOB),
		f(y, 0, ca.yMax-1, yMinOOB, yMaxOOB),
	)
	return c
}

func NothingOOB[T any](_ *CellularAutomata[T], x, y int, _, _, _, _ bool) *Cell[T] {
	return &Cell[T]{X: x, Y: y, IsBound: true, State: [2]T{}}
}

/* dead oob */
func DeadOOB(_ *CellularAutomata[int], x, y int, _, _, _, _ bool) *Cell[int] {
	return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{0, 0}}
}

/* live oob */
func LiveOOB(_ *CellularAutomata[int], x, y int, _, _, _, _ bool) *Cell[int] {
	return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{1, 1}}
}

/* random oob */
func RandomOOB(_ *CellularAutomata[int], x, y int, _, _, _, _ bool) *Cell[int] {
	r := imagegenerator.RandomCoin()
	return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{r, r}}
}

/* random bottom oob */
func RandomBottomOOB(_ *CellularAutomata[int], x, y int, _, _, _, yMaxOOB bool) *Cell[int] {
	s := 0
	if yMaxOOB && imagegenerator.RandomCoin() == 1 {
		s = 1
	}
	return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{s, s}}
}
