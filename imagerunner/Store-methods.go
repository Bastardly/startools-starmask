package imagerunner

import (
	"image"
	"starkiller/colortools"
)

func (store *Store) fillStore(img image.Image) {
	bounds := img.Bounds()

	store.Width = bounds.Max.X
	store.Height = bounds.Max.Y

	for y := 0; y < store.Height; y++ {
		var row []Pixel
		for x := 0; x < store.Width; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			column := rgbaToPixel(R, G, B, A)
			row = append(row, column)
		}
		store.Pixels = append(store.Pixels, row)
	}
}

func (store Store) maskPixels(starIndexStart, starIndexEnd, row, column int, isVertical bool) {
	radius := (starIndexEnd - starIndexStart) / 2

	for i := starIndexStart; i <= starIndexEnd; i++ {
		if isVertical {
			store.Pixels[i][column].HasContrastChangeVertical = true
			store.Pixels[i][column].starRadiusVertical = radius
		} else {
			store.Pixels[row][i].HasContrastChangeHorizontal = true
			store.Pixels[row][i].starRadiusHorizontal = radius
		}
	}
}

func (store Store) mapAlphaAreasHorizontal() {

	// todo lave func med start slut som go routine
	for row := 0; row < store.Height; row++ {
		// todo make this as a goRoutine + waitgroup
		starAreaActiveFrom := -1 // This is where we select the main star. Star and glow will form the mask.
		starAreaBrightnessLimit := float32(-1)
		previousBrightNess := store.Pixels[row][0].brightness

		resetStarArea := func() {
			starAreaActiveFrom = -1
			starAreaBrightnessLimit = -1
		}

		for column := 0; column < store.Width; column++ {
			pixel := store.Pixels[row][column]

			// if start active
			if starAreaActiveFrom > -1 {
				if column-starAreaActiveFrom > store.settings.maxStarSizeInPx {
					// Star is too big and we ignore it.
					resetStarArea()
				} else if starAreaBrightnessLimit > pixel.brightness {
					// We look for a brightness lower than starAreaBrightnessLimit
					store.maskPixels(starAreaActiveFrom, column-1, row, column, false)
					resetStarArea()
				}
				// Start star area
			} else if pixel.brightness > previousBrightNess {
				contrast := colortools.GetContrastRatio(pixel.brightness, previousBrightNess)

				if contrast >= store.settings.wcagContrastMinimum {
					starAreaActiveFrom = column
					starAreaBrightnessLimit = pixel.brightness
				}
			}

			previousBrightNess = pixel.brightness
		}
	}
}

func (store Store) mapAlphaAreasVertical() {

	// todo lave func med start slut som go routine
	for column := 0; column < store.Width; column++ {
		// todo make this as a goRoutine + waitgroup
		starAreaActiveFrom := -1 // This is where we select the main star. Star and glow will form the mask.
		starAreaBrightnessLimit := float32(-1)
		previousBrightNess := store.Pixels[0][column].brightness

		resetStarArea := func() {
			starAreaActiveFrom = -1
			starAreaBrightnessLimit = -1
		}

		for row := 0; row < store.Height; row++ {
			pixel := store.Pixels[row][column]

			// if start active
			if starAreaActiveFrom > -1 {
				if row-starAreaActiveFrom > store.settings.maxStarSizeInPx {
					// Star is too big and we ignore it.
					resetStarArea()
				} else if starAreaBrightnessLimit > pixel.brightness {
					// We look for a brightness lower than starAreaBrightnessLimit
					store.maskPixels(starAreaActiveFrom, row-1, row, column, true)
					resetStarArea()
				}
				// Start star area
			} else if pixel.brightness > previousBrightNess {
				contrast := colortools.GetContrastRatio(pixel.brightness, previousBrightNess)

				if contrast >= store.settings.wcagContrastMinimum {
					starAreaActiveFrom = row
					starAreaBrightnessLimit = pixel.brightness
				}
			}

			previousBrightNess = pixel.brightness
		}
	}
}

func (store Store) findBrightestConnectedPixelPosition(row, column int) (int, int) {
	bestRow := row
	bestCol := column
	currentRow := row
	currentCol := column
	brightest := store.Pixels[bestRow][bestCol].brightness
	currentDirection := uint8(1) //  right (1), left(2), down/right(3)
	exploring := true
	store.Pixels[currentRow][currentCol].HasBeenExplored = true

	goRight := func(isExploringDownwards bool) {
		if store.Pixels[currentRow][currentCol-1].IsValid() {
			currentCol -= 1
		} else if !isExploringDownwards {
			currentDirection = 2 // We explore to the left
		} else {
			exploring = false // We stop looking for more stars.
		}
	}
	goLeft := func() {
		if store.Pixels[currentRow][currentCol+1].IsValid() {
			currentCol += 1
		} else {
			currentDirection = 3 // We explore next row
		}
	}
	goDown := func() {
		if store.Pixels[currentRow+1][currentCol].IsValid() {
			currentRow += 1
			currentDirection = 1 // We start exploring right again.
		} else {
			goRight(true)
		}
	}

	for exploring {
		if currentDirection == 1 {
			goRight(false)
		} else if currentDirection == 2 {
			goLeft()
		} else {
			goDown()
		}

		if store.Pixels[currentRow][currentCol].HasBeenExplored {
			continue
		}

		store.Pixels[currentRow][currentCol].HasBeenExplored = true

		if store.Pixels[currentRow][currentCol].brightness > brightest {
			bestRow = currentRow
			bestCol = currentCol
			brightest = store.Pixels[currentRow][currentCol].brightness
		}
	}

	return bestRow, bestCol
}

func (store Store) addGlowToStarAlpha() {

	for row := 0; row < store.Height; row++ {

		for column := 0; column < store.Width; column++ {
			pixel := store.Pixels[row][column]
			if !pixel.IsValid() || pixel.HasBeenExplored {
				continue
			}

			bestRow, bestCol := store.findBrightestConnectedPixelPosition(row, column)
			store.Pixels[bestRow][bestCol].Rejected = true

		}
	}
}
