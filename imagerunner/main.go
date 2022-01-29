package imagerunner

import (
	"image"
)

func Start(img image.Image) ([][]Pixel, int, int) {
	pixels, width, height := getPixels(img)

	// Todo, make this a parameter and pass these values from UI
	mockSettings := Settings{
		maxStarSizeInPx:     12,
		maxStarGlowInPx:     3,
		wcagContrastMinimum: 1.6,
	}

	// Todo go routines and waitGroups
	mapAlphaAreasHorizontal(pixels, width, height, mockSettings)
	mapAlphaAreasVertical(pixels, width, height, mockSettings)
	addGlowToStarAlpha(pixels, width, height, mockSettings)

	return pixels, width, height
}
