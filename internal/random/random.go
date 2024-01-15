package random

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/jlowell000/game-of-life/internal/imagegenerator"
	"github.com/jlowell000/utils"
)

type RandomImageGenerator struct {
	Xmax, Ymax int
	Points     []image.Point
}

func (r *RandomImageGenerator) GenerateNextImage() image.Image {
	newImage, points := imagegenerator.CreateNewImage(r.Xmax, r.Ymax)
	colors := utils.MapWG(points, func(_ image.Point) color.NRGBA { return getRandomNRGB() })

	for i, p := range points {
		newImage.Set(p.X, p.Y, colors[i])
	}
	return newImage
}

func (r *RandomImageGenerator) GetPoints() []image.Point {
	return r.Points
}

func randomV() uint8 { return uint8(rand.Intn(255)) }
func getRandomNRGB() color.NRGBA {
	return color.NRGBA{R: randomV(), G: randomV(), B: randomV(), A: 255}
}
