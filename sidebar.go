package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Sidebar() *fyne.Container {
	button := widget.NewButtonWithIcon("Select image", theme.FileImageIcon(), func() {
		log.Println("tapped home")
	})
	button2 := widget.NewButtonWithIcon("Very", theme.DocumentIcon(), func() {
		log.Println("tapped home")
	})
	button3 := widget.NewButtonWithIcon("Darkside!", theme.ColorChromaticIcon(), func() {
		log.Println("tapped home")
	})

	return container.New(layout.NewVBoxLayout(), button, button2, button3)
}
