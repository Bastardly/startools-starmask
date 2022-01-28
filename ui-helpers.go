package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var textColor = color.Black

func getText(text string) *canvas.Text {
	formattedText := canvas.NewText(text, textColor)
	formattedText.TextSize = 14
	formattedText.Alignment = fyne.TextAlignCenter

	return formattedText
}

func getButton(text string, icon fyne.Resource, action func()) *fyne.Container {
	button := widget.NewButtonWithIcon(text, icon, action)

	return container.New(layout.NewHBoxLayout(), layout.NewSpacer(), button, layout.NewSpacer())
}

