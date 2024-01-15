package imagegenerator

import (
	"image"
	"math/rand"
)

type ImageGenerator interface {
	GenerateNextImage() image.Image
	GetPoints() []image.Point
}

func CreateNewImage(x, y int) (*image.NRGBA, []image.Point) {
	image := MakeNRGBASpace(x, y)
	return image, ImageToArray(image)
}

func ImageToArray(input image.Image) (result []image.Point) {
	for y := input.Bounds().Min.Y; y < input.Bounds().Max.Y; y++ {
		for x := input.Bounds().Min.X; x < input.Bounds().Max.X; x++ {
			result = append(result, image.Point{X: x, Y: y})
		}
	}
	return
}

func MakeNRGBASpace(x, y int) *image.NRGBA {
	return image.NewNRGBA(image.Rect(0, 0, x, y))
}

func RandomCoin() bool { return rand.Intn(2) == 0 }
