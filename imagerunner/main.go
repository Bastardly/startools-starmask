package imagerunner

import (
	"image"
)

func Start(img image.Image) ([][]Pixel, int, int) {
	mockSettings := Settings{
		maxStarSizeInPx:     12,
		maxStarGlowInPx:     3,
		wcagContrastMinimum: 1.6,
	}

	var store = Store{
		0, 0,
		mockSettings,
		[][]Pixel{},
	}

	store.fillStore(img)

	// Todo go routines and waitGroups
	store.mapAlphaAreasHorizontal()
	store.mapAlphaAreasVertical()
	store.addGlowToStarAlpha()

	return store.Pixels, store.Width, store.Height
}
