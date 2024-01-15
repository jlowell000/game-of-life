package gameoflife

import "image"

var Blinker = []image.Point{
	{X: 1, Y: 0},
	{X: 1, Y: 1},
	{X: 1, Y: 2},
}

var Glider = []image.Point{
	{X: 5, Y: 1},
	{X: 6, Y: 2},
	{X: 7, Y: 2},
	{X: 5, Y: 3},
	{X: 6, Y: 3},
}

var PentominoR = []image.Point{
	{X: 10, Y: 11},
	{X: 10, Y: 12},
	{X: 11, Y: 10},
	{X: 11, Y: 11},
	{X: 12, Y: 11},
}
