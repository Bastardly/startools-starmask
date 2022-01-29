package imagerunner

import "fmt"

func (p Pixel) IsValid() bool {
	return p.HasContrastChangeHorizontal && p.HasContrastChangeVertical
}

var count = 0

func findBrightestConnectedPixelPosition(pixels [][]Pixel, row, column int) (int, int) {
	bestRow := row
	bestCol := column
	currentRow := row
	currentCol := column
	brightest := pixels[bestRow][bestCol].brightness
	currentDirection := uint8(1) //  right (1), left(2), down/right(3)
	exploring := true
	pixels[currentRow][currentCol].HasBeenExplored = true

	goRight := func(isExploringDownwards bool) {
		if pixels[currentRow][currentCol-1].IsValid() {
			currentCol -= 1
		} else if !isExploringDownwards {
			currentDirection = 2 // We explore to the left
		} else {
			exploring = false // We stop looking for more stars.
		}
	}
	goLeft := func() {
		if pixels[currentRow][currentCol+1].IsValid() {
			currentCol += 1
		} else {
			currentDirection = 3 // We explore next row
		}
	}
	goDown := func() {
		if pixels[currentRow+1][currentCol].IsValid() {
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

		if pixels[currentRow][currentCol].HasBeenExplored {
			continue
		}

		pixels[currentRow][currentCol].HasBeenExplored = true

		if pixels[currentRow][currentCol].brightness > brightest {
			bestRow = currentRow
			bestCol = currentCol
			brightest = pixels[currentRow][currentCol].brightness
		}
	}

	return bestRow, bestCol
}

func addGlowToStarAlpha(pixels [][]Pixel, width, height int, settings Settings) {

	for row := 0; row < height; row++ {

		for column := 0; column < width; column++ {
			pixel := pixels[row][column]
			if !pixel.IsValid() || pixel.HasBeenExplored {
				continue
			}

			bestRow, bestCol := findBrightestConnectedPixelPosition(pixels, row, column)
			pixels[bestRow][bestCol].Rejected = true
			count++
		}
	}

	fmt.Println(count, "stars")
}
