package imagerunner

import (
	"fmt"
	"image"
	"math"
	"starkiller/colortools"
	"sync"
)

func (store *Store) fillStore(img image.Image) {
	bounds := img.Bounds()

	store.Width = bounds.Max.X
	store.Height = bounds.Max.Y

	for y := 0; y < store.Height; y++ {
		var row []Pixel
		for x := 0; x < store.Width; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			row = append(row, rgbaToPixel(R, G, B, A))
		}
		store.Pixels = append(store.Pixels, row)
	}
}

func (store Store) maskPixels(starIndexStart, starIndexEnd, row, column int, isVertical bool) {

	for i := starIndexStart; i <= starIndexEnd; i++ {
		if isVertical {
			store.Pixels[i][column].HasContrastChangeVertical = true
			store.Pixels[i][column].starRadiusStartVertical = starIndexStart
			store.Pixels[i][column].starRadiusEndVertical = starIndexEnd
		} else {
			store.Pixels[row][i].HasContrastChangeHorizontal = true
			store.Pixels[row][i].starRadiusStartHorizontal = starIndexStart
			store.Pixels[row][i].starRadiusEndHorizontal = starIndexEnd
		}
	}
}

func (store Store) mapAlphaAreasHorizontal(wg *sync.WaitGroup) {
	defer wg.Done()
	for row := 0; row < store.Height; row++ {
		starAreaActiveFrom := -1 // This is where we select the main star. Star and glow will form the mask.
		starAreaBrightnessLimit := float32(-1)
		previousBrightNess := store.Pixels[row][0].brightness

		resetStarArea := func() {
			starAreaActiveFrom = -1
			starAreaBrightnessLimit = -1
		}

		for column := 0; column < store.Width; column++ {
			// if start active
			if starAreaActiveFrom > -1 {
				if column-starAreaActiveFrom > store.settings.maxStarSizeInPx {
					// Star is too big and we ignore it.
					resetStarArea()
				} else if starAreaBrightnessLimit > store.Pixels[row][column].brightness {
					// We look for a brightness lower than starAreaBrightnessLimit
					store.maskPixels(starAreaActiveFrom, column-1, row, column, false)
					resetStarArea()
				}
				// Start star area
			} else if store.Pixels[row][column].brightness > previousBrightNess {
				contrast := colortools.GetContrastRatio(store.Pixels[row][column].brightness, previousBrightNess)

				if contrast >= store.settings.wcagContrastMinimum {
					starAreaActiveFrom = column
					starAreaBrightnessLimit = store.Pixels[row][column].brightness
				}
			}

			previousBrightNess = store.Pixels[row][column].brightness
		}
	}
}

func (store Store) mapAlphaAreasVertical(wg *sync.WaitGroup) {
	defer wg.Done()

	for column := 0; column < store.Width; column++ {
		starAreaActiveFrom := -1 // This is where we select the main star. Star and glow will form the mask.
		starAreaBrightnessLimit := float32(-1)
		previousBrightNess := store.Pixels[0][column].brightness

		resetStarArea := func() {
			starAreaActiveFrom = -1
			starAreaBrightnessLimit = -1
		}

		for row := 0; row < store.Height; row++ {
			// if start active
			if starAreaActiveFrom > -1 {
				if row-starAreaActiveFrom > store.settings.maxStarSizeInPx {
					// Star is too big and we ignore it.
					resetStarArea()
				} else if starAreaBrightnessLimit > store.Pixels[row][column].brightness {
					// We look for a brightness lower than starAreaBrightnessLimit
					store.maskPixels(starAreaActiveFrom, row-1, row, column, true)
					resetStarArea()
				}
				// Start star area
			} else if store.Pixels[row][column].brightness > previousBrightNess {
				contrast := colortools.GetContrastRatio(store.Pixels[row][column].brightness, previousBrightNess)

				if contrast >= store.settings.wcagContrastMinimum {
					starAreaActiveFrom = row
					starAreaBrightnessLimit = store.Pixels[row][column].brightness
				}
			}

			previousBrightNess = store.Pixels[row][column].brightness
		}
	}
}

func (store Store) findStarCenterPixelPosition(row, column int) (int, int, float64) {
	startRow := row
	startCol := column
	endRow := row
	endCol := column

	currentRow := row
	currentCol := column
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

		if currentCol < startCol {
			startCol = currentCol
		}
		if currentCol > endCol {
			endCol = currentCol
		}
		if currentRow < startRow {
			// should not happen since we start top left, but just in case...
			startRow = currentRow
		}
		if currentRow > endRow {
			endRow = currentRow
		}
	}
	rowMod := (endRow - startRow) / 2
	colMod := (endCol - startCol) / 2
	bestRow := row + rowMod
	bestCol := column + colMod
	starCoreRadius := rowMod

	if colMod > starCoreRadius {
		starCoreRadius = colMod
	}

	if bestCol < 0 || bestCol >= store.Width || bestRow < 0 || bestRow >= store.Height {
		panic("Coordinates are out of range! Check findStarCenterPixelPosition")
	}

	return bestRow, bestCol, float64(starCoreRadius)
}

func (store Store) findStarCenters() {
	for row := 0; row < store.Height; row++ {

		for column := 0; column < store.Width; column++ {
			if !store.Pixels[row][column].IsValid() || store.Pixels[row][column].HasBeenExplored {
				continue
			}

			bestRow, bestCol, starCoreRadius := store.findStarCenterPixelPosition(row, column)
			store.Pixels[bestRow][bestCol].isStarCenter = true
			store.Pixels[bestRow][bestCol].starCoreRadius = starCoreRadius + store.settings.starRadiusModifier
		}
	}
}

// markStarRadiusAsStar finds the size of the star + halo, and calls mask star with the calculated area
func (store Store) markStarRadiusAsStar() {
	var wg sync.WaitGroup

	for row := 0; row < store.Height; row++ {
		for col := 0; col < store.Width; col++ {
			wg.Add(1)

			go func(goRow, goCol int) {
				defer wg.Done()

				if store.Pixels[goRow][goCol].isStarCenter {
					store.Pixels[goRow][goCol].calculateStarRadiusWithGlow(store.settings.maxStarGlowInPx)

					startRow, startCol, endRow, endCol := store.Pixels[goRow][goCol].getSquareMapCoords(store, goRow, goCol)
					store.markStarAreas(goRow, goCol, startRow, startCol, endRow, endCol)

				}
			}(row, col)
		}
	}

	wg.Wait()
	fmt.Println("finished markStarRadiusAsStar")
}

// markStarAreas marks the area that will be affected by the star masking. Which will be initialised by this method since we have the selection ready.
func (store Store) markStarAreas(centerRow, centerCol, startRow, startCol, endRow, endCol int) {
	var starRadius = store.Pixels[centerRow][centerCol].starRadius
	var starCoreRadius = store.Pixels[centerRow][centerCol].starCoreRadius
	var wg sync.WaitGroup

	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			wg.Add(1)

			go func(goRow, goCol int) {
				defer wg.Done()
				store.Pixels[goRow][goCol].markAsStarIfWithinRange(centerRow, centerCol, goRow, goCol, starRadius, starCoreRadius)
			}(row, col)
		}
	}

	wg.Wait()

	if store.settings.blendMode == "fast" {
		// Copy color from 4 directions
		store.maskStarArea(startRow, startCol, endRow, endCol)
	} else if store.settings.blendMode == "cloneStamp" {
		store.cloneStampArea(centerRow, centerCol, startRow, startCol, endRow, endCol, starRadius)
	}
}

func (store Store) getPixelColorFromCoords(row, col int) Color {

	if row < 0 || col < 0 || row >= store.Height || col > store.Width {
		fmt.Println(row, store.Height, col, store.Width)
		panic("Fuuuuck")
	}

	return store.Pixels[row][col].getColor()
}

func (store Store) comparePixelBrightness(row1, row2, col1, col2 int) float64 {
	brightness1 := float64(store.Pixels[row1][col1].brightness)
	brightness2 := float64(store.Pixels[row2][col2].brightness)

	return (math.Abs(brightness1-brightness2) / brightness1)
}

func (store Store) clearStars() {
			for row := 0; row < store.Height; row++ {
			for col := 0; col < store.Width; col++ {
				store.Pixels[row][col].reset()
			}
		}
}
