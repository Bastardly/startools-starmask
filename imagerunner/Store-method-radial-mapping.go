package imagerunner

import "math"

func (store Store) runRadialStarMasking(row, col int) {
	store.Pixels[row][col].calculateStarRadiusWithGlow(store.settings.maxStarGlowInPx)
	maskingStarRadius := store.Pixels[row][col].starRadius + 1 // We want to start just outside the star

	for maskingStarRadius > 0 {
		maskingStarRadius--
		lengthOfHalfCircleInPx := maskingStarRadius * math.Pi
		lengthOfFullCircle := lengthOfHalfCircleInPx * 2
		radialGlowStrength := getRoundedFalloff(store.Pixels[row][col].starRadius, maskingStarRadius)

		for step := zero64; step < lengthOfFullCircle; step++ {
			radians := step / lengthOfFullCircle * math.Pi * 2
			targetRow := int(math.Sin(radians) * maskingStarRadius)
			targetCol := int(math.Cos(radians) * maskingStarRadius)
			// We add a second row to avoid gaps
			targetRow2 := int(math.Sin(radians) * (maskingStarRadius - 1))
			targetCol2 := int(math.Cos(radians) * (maskingStarRadius - 1))

			// Top and bottom of semi circle - These are the coords we are masking
			topRow := store.validatedRow(row - targetRow)
			topCol := store.validatedColumn(col - targetCol)
			topRow2 := store.validatedRow(row - targetRow2)
			topCol2 := store.validatedColumn(col - targetCol2)

			// Next we get the coords from the px we want to clone
			goRowUp := radians > math.Pi
			goColUp := radians < halfPI && radians >= 0 || radians >= oneAndAHalfPI
			rowMask, colMask := store.getMaskPixel(topRow, topCol, int(maskingStarRadius), goRowUp, goColUp)

			color := store.Pixels[rowMask][colMask].getColor()

			// We add the glow strength for this radius
			store.Pixels[topRow][topCol].radialGlowStrength = radialGlowStrength
			store.Pixels[topRow][topCol].setRadialColor(color, store.settings.radialMaskStrength)
			store.Pixels[topRow2][topCol2].radialGlowStrength = radialGlowStrength
			store.Pixels[topRow2][topCol2].setRadialColor(color, store.settings.radialMaskStrength)
		}

	}
}
