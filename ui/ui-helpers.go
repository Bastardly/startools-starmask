package ui

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

var textColor = color.Black

func getText(text string) *canvas.Text {
	formattedText := canvas.NewText(text, textColor)
	formattedText.TextSize = 14
	// formattedText.Alignment = fyne.TextAlignCenter

	return formattedText
}

func getHelpfulText(text string) *canvas.Text {
	txt := getText(text)
	txt.TextStyle.Italic = true

	return txt
}

// func getButton(text string, icon fyne.Resource, action func()) *fyne.Container {
// 	button := widget.NewButtonWithIcon(text, icon, action)

// 	return container.New(layout.NewHBoxLayout(), layout.NewSpacer(), button, layout.NewSpacer())
// }

// Get input takes the label text, a hint text
func getInput(text, hint, initialValue string) (*widget.FormItem, *numericalEntry) {
	entry := numericEntry()
	entry.Text = initialValue

	input := widget.NewFormItem(text, entry)
	input.HintText = hint

	return input, entry
}

// Get string input takes the label text, a hint text
func getStringInput(text, hint, initialValue string) (*widget.FormItem, *widget.Entry) {
	entry := widget.NewEntry()
	entry.Text = initialValue

	input := widget.NewFormItem(text, entry)
	input.HintText = hint

	return input, entry
}
