package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func MockContent() *fyne.Container {
	text1 := canvas.NewText("Hello", textColor)
	text2 := canvas.NewText("There", textColor)
	text3 := canvas.NewText("(right)", textColor)
	text := container.New(layout.NewHBoxLayout(), text1, text2, layout.NewSpacer(), text3)

	header := canvas.NewText("This is the header", color.RGBA{20, 50, 50, 1})

	return container.New(layout.NewVBoxLayout(), header, text)
}
