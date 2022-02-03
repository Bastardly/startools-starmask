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
			column := rgbaToPixel(R, G, B, A)
			row = append(row, column)
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

func (store Store) findStarCenters() {
	for row := 0; row < store.Height; row++ {

		for column := 0; column < store.Width; column++ {
			pixel := store.Pixels[row][column]
			if !pixel.IsValid() || pixel.HasBeenExplored {
				continue
			}

			bestRow, bestCol := store.findBrightestConnectedPixelPosition(row, column)
			store.Pixels[bestRow][bestCol].isStarCenter = true
		}
	}
}

func (store Store) getHaloFalloffLengthVertical(row, col, modifier, starCoreRadius int) int {
	prevBrightness := store.Pixels[row][col].brightness
	result := 0

	if store.settings.maxStarGlowInPx > 0 {
		for i := 1; i <= store.settings.maxStarGlowInPx; i++ {
			result = i
			nextRow := row + i*modifier
			if nextRow < 0 || nextRow >= store.Height {
				break
			}
			nextBrightness := store.Pixels[nextRow][col].brightness
			if nextBrightness >= prevBrightness {
				// it's not a deminishing halo
				break
			}
			prevBrightness = nextBrightness
		}
	}

	return result + starCoreRadius + extraPixel
}

func (store Store) getHaloFalloffLengthHorisontal(row, col, modifier, starCoreRadius int) int {
	prevBrightness := store.Pixels[row][col].brightness
	result := 0

	if store.settings.maxStarGlowInPx > 0 {
		for i := 1; i <= store.settings.maxStarGlowInPx; i++ {
			result = i
			nextCol := col + i*modifier
			if nextCol < 0 || nextCol >= store.Width {
				break
			}
			nextBrightness := store.Pixels[row][nextCol].brightness
			if nextBrightness >= prevBrightness {
				// it's not a deminishing halo
				break
			}
			prevBrightness = nextBrightness
		}
	}

	extra := starCoreRadius / 2

	return result + starCoreRadius + extra
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
					store.Pixels[goRow][goCol].calculateStarRadiusWithGlow(store, goRow, goCol)

					startRow, startCol, endRow, endCol := store.Pixels[goRow][goCol].getSquareMapCoords(store, goRow, goCol)
					store.markStarAreas(goRow, goCol, startRow, startCol, endRow, endCol)
				}
			}(row, col)
		}
	}

	wg.Wait()
}

// markStarAreas marks the area that will be affected by the star masking. Which will be initialised by this method since we have the selection ready.
func (store Store) markStarAreas(centerRow, centerCol, startRow, startCol, endRow, endCol int) {
	var wg sync.WaitGroup
	var starRadius = store.Pixels[centerRow][centerCol].starRadius

	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			wg.Add(1)

			go func(goRow, goCol int) {

				defer wg.Done()
				if goRow < 0 || goRow >= store.Height || goCol < 0 || goCol >= store.Width {
					fmt.Println(goCol, goRow)
					panic("Fuck arse!")
				}

				store.Pixels[goRow][goCol].markAsStarIfWithinRange(centerRow, centerCol, starRadius, goRow, goCol)
			}(row, col)
		}
	}
	wg.Wait()

	store.maskStarArea(startRow, startCol, endRow, endCol)
}

func (store Store) getPixelColorFromCoords(row, col int) Color {

	if row < 0 || col < 0 || row >= store.Height || col > store.Width {
		fmt.Println(row, store.Height, col, store.Width)
		panic("Fuuuuck")
	}

	p := store.Pixels[row][col]

	return Color{R: p.R, G: p.G, B: p.B}
}

// getFallBackColor return the darkest color of the four corners. Hopefulle this will suffice to avoid bright artefacts.
func (store Store) getFallBackColor(startRow, startCol, endRow, endCol int) Color {
	dark := store.Pixels[startRow][startCol]
	p2 := store.Pixels[startRow][endCol]
	p3 := store.Pixels[endRow][startCol]
	p4 := store.Pixels[endRow][endCol]

	if dark.brightness > p2.brightness {
		dark = p2
	}
	if dark.brightness > p3.brightness {
		dark = p3
	}
	if dark.brightness > p4.brightness {
		dark = p4
	}

	return Color{
		R: dark.R,
		G: dark.G,
		B: dark.B,
	}
}

func (store Store) setOppositeColorRow(row, col, modifier int, firstColor, fallbackColor Color) {
	oppositeColor := fallbackColor
	steps := 1
	nextRow := row

	for {
		nextRow = row + (steps * modifier)
		steps++

		if nextRow >= store.Height || nextRow < 0 {
			break
		}

		next := store.Pixels[nextRow][col]

		if !next.IsStar || next.IsMapped {
			oppositeColor = Color{
				R: next.R,
				G: next.G,
				B: next.B,
			}
			break
		}
	}

	store.Pixels[row][col].setMixedColor(firstColor, oppositeColor, steps)
}

func (store Store) setOppositeColorCol(row, col, modifier int, firstColor, fallbackColor Color) {
	oppositeColor := fallbackColor
	steps := 1
	nextCol := col

	// We search for the real opposite color
	for {
		nextCol = col + (steps * modifier)
		steps++

		if nextCol >= store.Width || nextCol < 0 {
			break
		}

		next := store.Pixels[row][nextCol]

		if !next.IsStar || next.IsMapped {
			oppositeColor = Color{
				R: next.R,
				G: next.G,
				B: next.B,
			}
			break
		}
	}

	store.Pixels[row][col].setMixedColor(firstColor, oppositeColor, steps)
}

func (store Store) setPixelMaskColorFromDirection(row, col int, direction string, fallbackColor Color) {
	firstColor := fallbackColor
	oppositeColor := fallbackColor
	switch direction {
	case "right":
		{
			// First color, if we go right, we look up
			if row-1 > 0 {
				firstColor = store.getPixelColorFromCoords(row-1, col)
			}
			store.setOppositeColorRow(row, col, 1, firstColor, oppositeColor)
		}

	case "left":
		{
			// First color, if we go left, we look down
			if row+1 < store.Height {
				firstColor = store.getPixelColorFromCoords(row+1, col)
			}
			store.setOppositeColorRow(row, col, -1, firstColor, oppositeColor)
		}

	case "up":
		{
			// First color, if we go up, we look left
			if col-1 > 0 {
				firstColor = store.getPixelColorFromCoords(row, col-1)
			}
			store.setOppositeColorCol(row, col, 1, firstColor, oppositeColor)
		}

	case "down":
		{
			// First color, if we go down, we look right (No fucking pun intended)
			if col+1 > store.Width {
				firstColor = store.getPixelColorFromCoords(row, col+1)
			}
			store.setOppositeColorCol(row, col, -1, firstColor, oppositeColor)
		}
	}
}

// func (store Store) maskStarArea(startRow, startCol, endRow, endCol int) {
// 	numberOfRows := endRow - startRow
// 	numberOfCols := endCol - startCol

// 	if numberOfCols < 1 || numberOfRows < 1 {
// 		return
// 	}

// 	fallbackColor := store.getFallBackColor(startRow, startCol, endRow, endCol)
// 	numberOfPixels := numberOfCols * numberOfRows
// 	direction := "right"
// 	row := startRow
// 	col := startCol

// 	for count := 0; count < numberOfPixels; count++ {
// 		p := store.Pixels[row][col]

// 		if p.IsStar {
// 			// We map the star with new color
// 			store.setPixelMaskColorFromDirection(row, col, direction, fallbackColor)
// 		}

// 		switch direction {
// 		case "right":
// 			{
// 				nextCol := col + 1
// 				if nextCol < endCol {
// 					col = nextCol
// 				} else {
// 					endCol -= 1
// 					direction = "down"
// 					row += 1
// 				}
// 			}

// 		case "left":
// 			{
// 				nextCol := col - 1
// 				if nextCol > startCol {
// 					col = nextCol
// 				} else {
// 					startCol += 1
// 					direction = "up"
// 					row -= 1
// 				}
// 			}

// 		case "up":
// 			{
// 				nextRow := row - 1
// 				if nextRow > startRow {
// 					row = nextRow
// 				} else {
// 					startRow += 1
// 					direction = "right"
// 					col += 1
// 				}
// 			}

// 		case "down":
// 			{
// 				nextRow := row + 1
// 				if nextRow < endRow {
// 					row = nextRow
// 				} else {
// 					endRow -= 1
// 					direction = "left"
// 					col -= 1
// 				}
// 			}
// 		}
// 	}
// }

//
//
//
//
//

func (store Store) maskStarArea(startRow, startCol, endRow, endCol int) {
	numberOfRows := endRow - startRow
	numberOfCols := endCol - startCol

	if numberOfCols < 1 || numberOfRows < 1 {
		return
	}

	var wg sync.WaitGroup
	numberOfRowsF64 := float64(numberOfRows)
	numberOfColsF64 := float64(numberOfCols)

	horizontal := func(reverse bool) {
		defer wg.Done()
		modifier, rowS, colS := getCorrectStartColsAndRowsIfReversed(startRow, startCol, endRow, endCol, reverse)
		// For each column, I need to find the row that contains the color I want.
		colorList := make([]int, numberOfCols+1)

		for row := 0; row <= numberOfRows; row++ {
			for col := 0; col <= numberOfCols; col++ {
				wg.Add(1)
				targetRow := rowS + (row * modifier)
				targetCol := colS + (col * modifier)
				if row == 0 {
					colorList[col] = targetRow
				}
				go func(goRow, goCol, colorIndex int) {
					defer wg.Done()
					if store.Pixels[goRow][goCol].IsStar {
						sourceColorRow := colorList[colorIndex]
						pxColor := store.getPixelColorFromCoords(goRow, goCol)
						srcColor := store.getPixelColorFromCoords(sourceColorRow, goCol)
						distance := math.Abs(float64(sourceColorRow) - float64(goRow))
						procentage := (numberOfRowsF64 - distance + 1) / numberOfRowsF64
						store.Pixels[goRow][goCol].modifyColors(procentage, pxColor, srcColor)

					} else { // Todo, add row limit 50% or run seriel if results are funky, but I doubt it will be a problem
						colorList[colorIndex] = goRow
					}

				}(targetRow, targetCol, col)
			}
		}
	}

	vertical := func(reverse bool) {
		defer wg.Done()
		modifier, rowS, colS := getCorrectStartColsAndRowsIfReversed(startRow, startCol, endRow, endCol, reverse)
		// For each column, I need to find the row that contains the color I want.
		colorList := make([]int, numberOfRows+1)

		for col := 0; col <= numberOfCols; col++ {
			for row := 0; row <= numberOfRows; row++ {
				wg.Add(1)
				targetRow := rowS + (row * modifier)
				targetCol := colS + (col * modifier)
				if col == 0 {
					colorList[row] = targetCol
				}
				go func(goRow, goCol, colorIndex int) {
					defer wg.Done()
					if store.Pixels[goRow][goCol].IsStar {
						sourceColorCol := colorList[colorIndex]
						pxColor := store.getPixelColorFromCoords(goRow, goCol)
						scrColor := store.getPixelColorFromCoords(goRow, sourceColorCol)
						distance := math.Abs(float64(sourceColorCol) - float64(goCol))
						procentage := (numberOfColsF64 - distance + 1) / numberOfColsF64
						store.Pixels[goRow][goCol].modifyColors(procentage, pxColor, scrColor)

					} else { // Todo, add row limit 50% or run seriel if results are funky, but I doubt it will be a problem
						colorList[colorIndex] = goCol
					}

				}(targetRow, targetCol, row)
			}
		}
	}

	wg.Add(4) // One for each direction!
	go horizontal(false)
	go horizontal(true)
	go vertical(false)
	go vertical(true)

	//          todo Vertical also!

	wg.Wait()
}
