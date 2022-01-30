package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()

	myWindow := a.NewWindow("Star Killer 3000")
	myWindow.Resize(fyne.NewSize(800, 600))

	r, _ := fyne.LoadResourceFromPath("Icon2.png")
	myWindow.SetIcon(r)

	title := getText("STAR KILLER 3000")
	title.TextSize = 30

	quote := getText("\"Now you see that evil will always triumph, because good is dumb.\"")
	quote.TextStyle.Italic = true

	quoted := getText("- Darth Helmet")
	quoted.TextStyle.Monospace = true
	quoted.TextSize = 12

	instructions := getText("STEP ONE: Select a TIF (or TIFF) file with tiny stars that we can annihilate!")
	instructions.TextStyle.Monospace = true

	selectImageButton := getButton("Select file", theme.FileImageIcon(), func() {
		selectFile()
	})
	firstPage := container.New(layout.NewVBoxLayout(), title, quote, quoted, layout.NewSpacer(), instructions, selectImageButton, layout.NewSpacer())

	myWindow.SetContent(firstPage)
	myWindow.ShowAndRun()
}
