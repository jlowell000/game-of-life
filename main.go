package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
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
	"github.com/jlowell000/game-of-life/internal/gameoflife"
)

const imageSize int = 100

var (
	tag                = new(bool)
	framesToGen   uint = 0
	framesGened   uint = 1
	pressed       bool = false
	playing       bool = false
	gameoflifeGen gameoflife.GameOfLifeGenerator
	// gen           imagegenerator.ImageGenerator
)

func main() {
	// _, points := imagegenerator.CreateNewImage(imageSize, imageSize)
	// gen = &random.RandomImageGenerator{Xmax: imageSize, Ymax: imageSize, Points: points}
	gameoflifeGen = gameoflife.CreateInitalState(imageSize, imageSize, []image.Point{})
	// gameoflifeGen = gameoflife.CreateInitalState(imageSize, imageSize, gameoflife.Blinker)
	// gameoflifeGen = gameoflife.CreateInitalState(imageSize, imageSize, gameoflife.Glider)
	// gameoflifeGen = gameoflife.CreateInitalState(imageSize, imageSize, gameoflife.PentominoR)
	// gen = &gameoflifeGen

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
			if (playing || framesToGen > 0) && framesGened < 5000 {
				if framesToGen > 0 {
					framesToGen--
				}
				gameoflifeGen.GenerateNextState()
				framesGened++
			}
			createDebug(gtx)
			drawImage(&ops, gameoflifeGen.GetImageFromState())

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
	title := material.Body1(th, fmt.Sprintf("playing: %t; spacePressed:%t; framesToGen: %d; framesGened: %d", playing, pressed, framesToGen, framesGened))
	maroon := color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	title.Color = maroon
	title.Alignment = text.Middle
	return title.Layout(gtx)
}

func drawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Add(ops)
	scaleFactor := float32(1000 / imageSize)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scaleFactor, scaleFactor))).Add(ops)
	op.Offset(image.Pt(0, 5)).Add(ops)
	paint.PaintOp{}.Add(ops)
}
