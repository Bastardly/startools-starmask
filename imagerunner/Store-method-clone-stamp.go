package imagerunner

import (
	"math"
)

func (store Store) findCloneArea(centerRow, centerCol, startRow, startCol, endRow, endCol int, starRadius float64) (int, int) {
	factor := two64
	multiplier := starRadius * factor
	rowMod := centerRow - startRow
	colMod := centerCol - startCol
	rowLength := endRow - startRow
	colLength := endCol - startCol
	rowS := startRow
	colS := startCol
	limit := float64(60)
	lowestDiff := two64

	for radians := zero64; radians < limit; radians++ {
		if radians > 0 && int(radians)%20 == 0 {
			factor++ // We expand the search outwards
		}

		// We validate each coord, to see if it's out of frame.
		cloneRowCenter := centerRow + int(math.Sin(radians)*multiplier)
		cloneRowStart := cloneRowCenter - rowMod
		if cloneRowStart < 0 {
			continue
		}
		cloneRowEnd := cloneRowStart + rowLength
		if cloneRowEnd >= store.Height {
			continue
		}

		cloneColCenter := centerCol + int(math.Cos(radians)*multiplier)
		cloneColStart := cloneColCenter - colMod
		if cloneColStart < 0 {
			continue
		}
		cloneColEnd := cloneColStart + colLength
		if cloneColEnd >= store.Width {
			continue
		}

		hasStar := false
		for row := cloneRowStart; row < cloneRowEnd; row++ {
			for col := cloneColStart; col < cloneColEnd; col++ {
				if store.Pixels[row][col].IsStar || store.Pixels[row][col].glowStrength > 0.1 || store.Pixels[row][col].brightness > 0.8 {
					hasStar = true
				}
			}
		}
		diff := store.comparePixelBrightness(startRow, cloneRowStart, startCol, cloneColStart) + store.comparePixelBrightness(endRow, cloneRowEnd, endCol, cloneColEnd)
		if diff < lowestDiff && !hasStar {
			lowestDiff = diff
			rowS = cloneRowStart
			colS = cloneColStart
		}
	}

	return rowS, colS
}

func (store Store) cloneStampArea(centerRow, centerCol, startRow, startCol, endRow, endCol int, starRadius float64) {
	cloneRowStart, cloneColStart := store.findCloneArea(centerRow, centerCol, startRow, startCol, endRow, endCol, starRadius)
	rowLength := endRow - startRow
	colLength := endCol - startCol

	for row := 0; row < rowLength; row++ {
		for col := 0; col < colLength; col++ {
			cloneColor := store.getPixelColorFromCoords(cloneRowStart+row, cloneColStart+col)
			store.Pixels[startRow+row][startCol+col].setColor(cloneColor)
		}
	}

}
