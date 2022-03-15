package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func GetHeader(resource fyne.Resource) *fyne.Container {
	title := getText("STAR TOOLS 3000 - STAR MASK")
	title.TextSize = 30

	image := canvas.NewImageFromResource(resource)
	image.SetMinSize(fyne.Size{Width: 50, Height: 50})

	spacer := GetBlock(30)

	// c := canvas.NewImageFromImage(image)
	// size := fyne.Size{
	// 	Width: 500, Height: 500,
	// }
	// image.Resize(size)

	// quote := getText("\"Brother Maynard – bring forth the holy hand grenade!\"")
	// quote.TextStyle.Italic = true

	// quoted := getText("- King Arthur")
	// quoted.TextStyle.Monospace = true
	// quoted.TextSize = 12

	cont := container.New(layout.NewHBoxLayout(), image, title)
	return container.New(layout.NewVBoxLayout(), cont, spacer)

}
