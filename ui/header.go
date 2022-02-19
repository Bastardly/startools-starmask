package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func GetHeader() *fyne.Container {
	title := getText("STAR TOOLS 3000")
	title.TextSize = 30

	quote := getText("\"Brother Maynard – bring forth the holy hand grenade!\"")
	quote.TextStyle.Italic = true

	quoted := getText("- King Arthur")
	quoted.TextStyle.Monospace = true
	quoted.TextSize = 12

	return container.New(layout.NewVBoxLayout(), title, quote, quoted)

}
