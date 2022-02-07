package imagerunner

import (
	"math"
	"sync"
)

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
					if store.Pixels[goRow][goCol].IsStar || store.Pixels[goRow][goCol].glowStrength > 0 {
						sourceColorRow := colorList[colorIndex]
						pxColor := store.getPixelColorFromCoords(goRow, goCol)
						srcColor := store.getPixelColorFromCoords(sourceColorRow, goCol)
						distance := math.Abs(float64(sourceColorRow) - float64(goRow))
						procentage := getRoundedFalloff(numberOfRowsF64, distance)
						store.Pixels[goRow][goCol].modifyColors(procentage, pxColor, srcColor)

					} else { //} if store.Pixels[goRow][goCol].brightness < store.Pixels[targetRow][targetCol].brightness {
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
					if store.Pixels[goRow][goCol].IsStar || store.Pixels[goRow][goCol].glowStrength > 0 {
						sourceColorCol := colorList[colorIndex]
						pxColor := store.getPixelColorFromCoords(goRow, goCol)
						scrColor := store.getPixelColorFromCoords(goRow, sourceColorCol)
						distance := math.Abs(float64(sourceColorCol) - float64(goCol))
						procentage := getRoundedFalloff(numberOfColsF64, distance)
						store.Pixels[goRow][goCol].modifyColors(procentage, pxColor, scrColor)

					} else { // if store.Pixels[goRow][goCol].brightness < store.Pixels[targetRow][targetCol].brightness {
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

	wg.Wait()
}
