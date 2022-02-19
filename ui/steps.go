package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func GetStepOne(getFile func()) *fyne.Container {
	label := getText("STEP ONE: ")
	label.TextStyle.Monospace = true
	selectImageButton := widget.NewButtonWithIcon("Select a 16bit TIF file", theme.FileImageIcon(), getFile)
	// selectImageButton := getButton("Select file", theme.FileImageIcon(), getFile)
	helpfulText := getHelpfulText("Tip: You can always select a new image, and reuse the settings.")
	btnContainer := container.New(layout.NewHBoxLayout(), label, selectImageButton)

	return container.New(layout.NewVBoxLayout(), btnContainer, helpfulText)
}

// func GetStepTwo(removeStars, createAlpha func(b bool)) *fyne.Container {
// 	label := getText("STEP TWO: Select output- and star settings")
// 	label.TextStyle.Monospace = true
// 	areaText := getHelpfulText("Output settings")
// 	removeStarsComponent := widget.NewCheck("An image with removed stars", removeStars)
// 	createAlphaComponent := widget.NewCheck("A Star map - Great for post processing", createAlpha)

// 	cont := container.New(layout.NewVBoxLayout(), label, areaText, removeStarsComponent, createAlphaComponent)
// 	cont.Hidden = true

// 	return cont
// }

func GetStepTwoPartTwo(SubmitText string, submitFunc func(_, _, _, _ string)) *fyne.Container {
	label := getText("STEP TWO: Select star settings and run")
	label.TextStyle.Monospace = true
	areaText := getHelpfulText("Star settings")
	minStarSize, minStarSizeEntry := getInput("Min star size (px)", "Set the minimum star diameter (Measured in pixels) to be included", "0")
	maxStarSize, maxStarSizeEntry := getInput("Max star size (px)", "Set the maximum size of the star core without the glow (Measured in pixels)", "10")
	minContrastRatio, minContrastRatioEntry := getInput("Star clarity", "How clear the stars must be. Putting this setting too low, will often mask too much (1.7 is a good start)", "1.7")
	filePrefix, filePrefixEntry := getStringInput("File prefix", "Prefix for your output file.", "")

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			minStarSize, maxStarSize, minContrastRatio, filePrefix,
		},
		OnSubmit: func() {
			submitFunc(minStarSizeEntry.Text, maxStarSizeEntry.Text, minContrastRatioEntry.Text, filePrefixEntry.Text)
		},
	}
	form.SubmitText = SubmitText
	helpfulText := getHelpfulText("Tip: Your output files will be placed in the same folder as your input file.")

	cont := container.New(layout.NewVBoxLayout(), areaText, form, helpfulText)
	cont.Hidden = true

	return cont
}

// func GetPartThree(method func()) *fyne.Container {
// 	label := getText("STEP THREE: ")
// 	label.TextStyle.Monospace = true
// 	submitButton := widget.NewButtonWithIcon("RUN ST3K", theme.MediaPlayIcon(), method)
// 	helpfulText := getHelpfulText("Tip: Your output files will be placed in the same folder as your input file.")

// 	btnContainer := container.New(layout.NewHBoxLayout(), label, submitButton)

// 	cont := container.New(layout.NewVBoxLayout(), btnContainer, helpfulText)
// 	cont.Hidden = true

// 	return cont
// }

func GetSplashScreen(text string) *fyne.Container {
	label := getText(text)
	labelContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), label, layout.NewSpacer())

	return container.New(layout.NewVBoxLayout(), layout.NewSpacer(), labelContainer, layout.NewSpacer())
}
