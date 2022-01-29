package imagerunner

import (
	"starkiller/colortools"
)

func maskPixels(pixels [][]Pixel, starIndexStart, starIndexEnd, row, column int, isVertical bool) {
	radius := (starIndexEnd - starIndexStart) / 2

	for i := starIndexStart; i <= starIndexEnd; i++ {
		if isVertical {
			pixels[i][column].HasContrastChangeVertical = true
			pixels[i][column].starRadiusVertical = radius
		} else {
			pixels[row][i].HasContrastChangeHorizontal = true
			pixels[row][i].starRadiusHorizontal = radius
		}

	}
}

func mapAlphaAreasHorizontal(pixels [][]Pixel, width, height int, settings Settings) {

	// todo lave func med start slut som go routine
	for row := 0; row < height; row++ {
		// todo make this as a goRoutine + waitgroup
		starAreaActiveFrom := -1 // This is where we select the main star. Star and glow will form the mask.
		starAreaBrightnessLimit := float32(-1)
		previousBrightNess := pixels[row][0].brightness

		resetStarArea := func() {
			starAreaActiveFrom = -1
			starAreaBrightnessLimit = -1
		}

		for column := 0; column < width; column++ {
			pixel := pixels[row][column]

			// if start active
			if starAreaActiveFrom > -1 {
				if column-starAreaActiveFrom > settings.maxStarSizeInPx {
					// Star is too big and we ignore it.
					resetStarArea()
				} else if starAreaBrightnessLimit > pixel.brightness {
					// We look for a brightness lower than starAreaBrightnessLimit
					maskPixels(pixels, starAreaActiveFrom, column-1, row, column, false)
					resetStarArea()
				}
				// Start star area
			} else if pixel.brightness > previousBrightNess {
				contrast := colortools.GetContrastRatio(pixel.brightness, previousBrightNess)

				if contrast >= settings.wcagContrastMinimum {
					starAreaActiveFrom = column
					starAreaBrightnessLimit = pixel.brightness
				}
			}

			previousBrightNess = pixel.brightness
		}
	}
}

func mapAlphaAreasVertical(pixels [][]Pixel, width, height int, settings Settings) {

	// todo lave func med start slut som go routine
	for column := 0; column < width; column++ {
		// todo make this as a goRoutine + waitgroup
		starAreaActiveFrom := -1 // This is where we select the main star. Star and glow will form the mask.
		starAreaBrightnessLimit := float32(-1)
		previousBrightNess := pixels[0][column].brightness

		resetStarArea := func() {
			starAreaActiveFrom = -1
			starAreaBrightnessLimit = -1
		}

		for row := 0; row < height; row++ {
			pixel := pixels[row][column]

			// if start active
			if starAreaActiveFrom > -1 {
				if row-starAreaActiveFrom > settings.maxStarSizeInPx {
					// Star is too big and we ignore it.
					resetStarArea()
				} else if starAreaBrightnessLimit > pixel.brightness {
					// We look for a brightness lower than starAreaBrightnessLimit
					maskPixels(pixels, starAreaActiveFrom, row-1, row, column, true)
					resetStarArea()
				}
				// Start star area
			} else if pixel.brightness > previousBrightNess {
				contrast := colortools.GetContrastRatio(pixel.brightness, previousBrightNess)

				if contrast >= settings.wcagContrastMinimum {
					starAreaActiveFrom = row
					starAreaBrightnessLimit = pixel.brightness
				}
			}

			previousBrightNess = pixel.brightness
		}
	}
}
