package imagegenerator

import (
	"image"
	"math/rand"
)

type ImageGenerator interface {
	GenerateNextState()
	GetImageFromCurrentState() image.Image
}

func CreateNewImageAndPoints(x, y int) (*image.NRGBA, []*image.Point) {
	newImage := CreateNewImage(x, y)
	return newImage, CreatePoints(newImage)
}

func CreateNewImage(x, y int) *image.NRGBA {
	return image.NewNRGBA(image.Rect(0, 0, x, y))
}

func CreatePoints(newImage *image.NRGBA) (points []*image.Point) {
	for y := newImage.Bounds().Min.Y; y < newImage.Bounds().Max.Y; y++ {
		for x := newImage.Bounds().Min.X; x < newImage.Bounds().Max.X; x++ {
			points = append(points, &image.Point{X: x, Y: y})
		}
	}
	return
}

func RandomCoin() int { return rand.Intn(2) }
