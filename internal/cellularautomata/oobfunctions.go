package cellularautomata

import "github.com/jlowell000/game-of-life/internal/imagegenerator"

/* wrapping oob */
func OverflowOOB[T any](ca *CellularAutomata[T]) func(int, int, bool, bool, bool, bool) *Cell[T] {
	return func(x, y int, xMinOOB, xMaxOOB, yMinOOB, yMaxOOB bool) *Cell[T] {
		f := func(v, vMin, vMax int, minOOB, maxOOB bool) int {
			if minOOB {
				v = vMax
			} else if maxOOB {
				v = vMin
			}
			return v
		}
		return ca.GetCellAt(
			f(x, 0, ca.xMax-1, xMinOOB, xMaxOOB),
			f(y, 0, ca.yMax-1, yMinOOB, yMaxOOB),
		)
	}
}

func NothingOOB[T any](ca *CellularAutomata[T]) func(int, int, bool, bool, bool, bool) *Cell[T] {
	return func(x, y int, _, _, _, _ bool) *Cell[T] {
		return &Cell[T]{X: x, Y: y, IsBound: true, State: [2]T{}}
	}
}

/* dead oob */
func DeadOOB(_ *CellularAutomata[int]) func(int, int, bool, bool, bool, bool) *Cell[int] {
	return func(x, y int, _, _, _, _ bool) *Cell[int] {
		return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{0, 0}}
	}
}

/* live oob */
func LiveOOB(_ *CellularAutomata[int]) func(int, int, bool, bool, bool, bool) *Cell[int] {
	return func(x, y int, _, _, _, _ bool) *Cell[int] {
		return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{1, 1}}
	}
}

/* random oob */
func RandomOOB(_ *CellularAutomata[int]) func(int, int, bool, bool, bool, bool) *Cell[int] {
	return func(x, y int, _, _, _, _ bool) *Cell[int] {
		r := imagegenerator.RandomCoin()
		return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{r, r}}
	}
}

/* random bottom oob */
func RandomBottomOOB(_ *CellularAutomata[int]) func(int, int, bool, bool, bool, bool) *Cell[int] {
	return func(x, y int, _, _, _, yMaxOOB bool) *Cell[int] {
		s := 0
		if yMaxOOB && imagegenerator.RandomCoin() == 1 {
			s = 1
		}
		return &Cell[int]{X: x, Y: y, IsBound: true, State: [2]int{s, s}}
	}
}
