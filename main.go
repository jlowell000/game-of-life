package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/jlowell000/utils"
)

const imageSize int = 100

var (
	tag               = new(bool)
	framesToGen  uint = 0
	pressed      bool = false
	playing      bool = false
	currentImage image.Image
)

func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops

	nextStep := func() {
		if !playing {

			return
		}
		w.Invalidate()
	}

	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			keylistener(&ops, e.Queue)
			if playing || framesToGen > 0 {
				if framesToGen > 0 {
					framesToGen--
				}
				currentImage = generateImage(imageSize, imageSize)
			}
			createDebug(gtx)
			if currentImage != nil {
				drawImage(&ops, currentImage)
			}

			e.Frame(gtx.Ops)
			nextStep()
		}

	}
}

func keylistener(ops *op.Ops, q event.Queue) {
	for _, event := range q.Events(tag) {
		if x, ok := event.(key.Event); ok {
			switch x.State {
			case key.Press:
				if !pressed {
					pressed = true
					switch x.Name {
					case key.NameSpace:
						playing = !playing
						framesToGen = 0
					case key.NameRightArrow:
						framesToGen++
					}
				}
			case key.Release:
				pressed = false
			}
		}
	}

	key.InputOp{
		Tag: tag,
		Keys: "[" +
			key.NameSpace + "," +
			key.NameRightArrow +
			"]",
	}.Add(ops)

}

func createDebug(gtx layout.Context) layout.Dimensions {
	th := material.NewTheme()
	title := material.Body1(th, fmt.Sprintf("playing: %t; spacePressed:%t; framesToGen: %d", playing, pressed, framesToGen))
	maroon := color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	title.Color = maroon
	title.Alignment = text.Middle
	return title.Layout(gtx)
}

func drawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4))).Add(ops)
	op.Offset(image.Pt(0, 5)).Add(ops)
	paint.PaintOp{}.Add(ops)
}

func generateImage(x, y int) image.Image {
	newImage := makeNRGBASpace(x, y)
	points := imageToArray(newImage)
	colors := utils.MapWG(points, func(_ image.Point) color.NRGBA { return getRandomNRGB() })

	for i, p := range points {
		newImage.Set(p.X, p.Y, colors[i])
	}
	return newImage
}

func imageToArray(input image.Image) (result []image.Point) {
	for y := input.Bounds().Min.Y; y < input.Bounds().Max.Y; y++ {
		for x := input.Bounds().Min.X; x < input.Bounds().Max.X; x++ {
			result = append(result, image.Point{X: x, Y: y})
		}
	}
	return
}

func makeNRGBASpace(x, y int) *image.NRGBA {
	return image.NewNRGBA(image.Rect(0, 0, x, y))
}

func getRandomNRGB() color.NRGBA {
	f := func() uint8 { return uint8(rand.Intn(255)) }
	return color.NRGBA{R: f(), G: f(), B: f(), A: 255}
}
