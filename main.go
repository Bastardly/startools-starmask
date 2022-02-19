package main

import (
	"image"
	"starkiller/createfile"
	"starkiller/imagerunner"
	"starkiller/ui"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type AppData struct {
	imageFile        image.Image
	uri              string
	filePrefixEntry  string
	removeStars      bool
	createAlpha      bool
	minStarSize      int
	maxStarSize      int
	minContrastRatio float32
	firstPage        *fyne.Container
}

func startImagerun(data AppData) {
	imageRunnerData := imagerunner.IStart{
		ImageFile:        &data.imageFile,
		RemoveStars:      data.removeStars,
		CreateAlpha:      data.createAlpha,
		MinStarSize:      data.minStarSize,
		MaxStarSize:      data.maxStarSize,
		MinContrastRatio: data.minContrastRatio,
	}
	pixels, alphaPixels, width, height := imagerunner.Start(imageRunnerData)
	if data.createAlpha {
		createfile.CreateAlpha(alphaPixels, data.uri, data.filePrefixEntry, width, height)
	}
	if data.removeStars {
		createfile.CreatePreview(pixels, data.uri, data.filePrefixEntry, width, height)
	}
}

// Start app with air
func main() {
	var data AppData
	data.createAlpha = true
	a := app.New()

	myWindow := a.NewWindow("STAR TOOLS 3000")
	myWindow.Resize(fyne.NewSize(800, 600))

	r, _ := fyne.LoadResourceFromPath("Icon2.png")
	myWindow.SetIcon(r)

	header := ui.GetHeader()

	// stepTwo := ui.GetStepTwo(func(b bool) {
	// 	data.removeStars = b
	// }, func(b bool) {
	// 	data.createAlpha = b
	// })

	// this contains the submit
	stepTwoPartTwo := ui.GetStepTwoPartTwo("Create starmap", func(minStarSize, maxStarSize, minContrastRatio, filePrefixEntry string) {
		working := ui.GetSplashScreen("WORKING AT LUDICROUS SPEED!!")
		done := ui.GetSplashScreen("DONE")

		// Add form data to data
		minStarSizeValue, _ := strconv.Atoi(minStarSize)
		maxStarSizeValue, _ := strconv.Atoi(maxStarSize)
		contrastValue, _ := strconv.ParseFloat(minContrastRatio, 32)
		data.minStarSize = int(minStarSizeValue)
		data.maxStarSize = int(maxStarSizeValue)
		data.minContrastRatio = float32(contrastValue)
		data.filePrefixEntry = filePrefixEntry

		myWindow.SetContent(working)
		startImagerun(data)
		myWindow.SetContent(done)
		time.Sleep(1 * time.Second)
		myWindow.SetContent(data.firstPage)
	})

	// stepThree := ui.GetPartThree(func() {
	// 	working := ui.GetSplashScreen("WORKING AT LUDICROUS SPEED!!")
	// 	done := ui.GetSplashScreen("DONE")

	// 	myWindow.SetContent(working)
	// 	startImagerun(data)
	// 	myWindow.SetContent(done)
	// 	time.Sleep(1 * time.Second)
	// 	myWindow.SetContent(data.firstPage)
	// })

	stepOne := ui.GetStepOne(func() {
		file, uriStr, err := selectFile()
		if err == nil {
			data.imageFile = file
			data.uri = uriStr
			// stepTwo.Show()
			stepTwoPartTwo.Show()
			// stepThree.Show()
		}
	})

	data.firstPage = container.New(layout.NewVBoxLayout(), header, stepOne, layout.NewSpacer(), stepTwoPartTwo, layout.NewSpacer())

	myWindow.SetContent(data.firstPage)
	myWindow.ShowAndRun()
}
