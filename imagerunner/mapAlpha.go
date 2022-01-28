package imagerunner

import (
	"fmt"
	"starkiller/colortools"
)

func maskPixels(pixels [][]Pixel, startIndex, endIndex, row int, isReverse bool) {
	for i := startIndex; i <= endIndex; i++ {
		if isReverse {
			pixels[row][i].HasContrastChangeLeft = true
		} else {
			pixels[row][i].HasContrastChangeRight = true

		}
	}
}

func findGlowAndMaskPixels(pixels [][]Pixel, starIndexStart, starIndexEnd, maxStarGlowInPx, row, width int, isReverse bool) {
	starSize := starIndexEnd - starIndexStart
	if starSize < 0 {
		panic("Star start index is larger than end index. ")
	}

	alphaStartIndex := starIndexStart
	alphaEndIndex := starIndexEnd
	expandBy := int(starSize / 4) // We add a bit more than the glow depending on the star
	startActive := true
	endActive := true
	prevStartBrightness := pixels[row][starIndexStart].brightness
	prevEndBrightness := pixels[row][starIndexEnd].brightness

	for i := 1; i <= maxStarGlowInPx; i++ {
		if !startActive && !endActive {
			break
		}
		s := starIndexStart - i
		e := starIndexStart + i
		if startActive && s >= 0 {
			pixel := pixels[row][s]
			if pixel.brightness < prevStartBrightness {
				alphaStartIndex = s
			}
		}

		if endActive && e < width {
			pixel := pixels[row][e]
			if pixel.brightness < prevEndBrightness {
				alphaEndIndex = e
			}
		}
	}

	alphaStartIndex -= expandBy
	alphaEndIndex += expandBy

	if alphaStartIndex < 0 {
		alphaStartIndex = 0
	}
	if alphaEndIndex >= width {
		alphaEndIndex = width - 1
	}

	maskPixels(pixels, alphaStartIndex, alphaEndIndex, row, isReverse)
}

func mapAlphaAreasForwardHorizontal(pixels [][]Pixel, width, height int, settings Settings) {

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
					findGlowAndMaskPixels(pixels, starAreaActiveFrom, column-1, settings.maxStarGlowInPx, row, width, false)
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

// func mapContrastAreasBackwardHorizontal(pixels [][]Pixel, width, height int, settings Settings) {

// 	// todo lave func med start slut som go routine
// 	for row := height; row > 0; row-- {
// 		// todo make this as a goRoutine + waitgroup
// 		for column := width; column > 0; column-- {
// 			pixel := pixels[row][column]

// 		}
// 	}

// }
