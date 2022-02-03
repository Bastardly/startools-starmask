package imagerunner

import (
	"starkiller/colortools"
)

func rgbaToPixel(R uint32, G uint32, B uint32, A uint32) Pixel {
	r := uint8(R / bitDivider)
	g := uint8(G / bitDivider)
	b := uint8(B / bitDivider)
	brightness := colortools.GetRGBBrightNess(r, g, b)

	return Pixel{R, G, B, A, brightness, false, false, false, false, false, false, -1, -1, -1, -1, -1}
}

// validateCoord makes sure the calculated coordinate stays within table
func validateCoord(value, fallback int) int {
	if value < 0 {
		return 0
	}

	if value >= fallback {
		return fallback - 1
	}

	return value
}

func getCorrectStartColsAndRowsIfReversed(startRow, startCol, endRow, endCol int, reverse bool) (int, int, int) {
	modifier := 1
	if reverse {
		modifier = -1
	}
	rowS := startRow
	colS := startCol

	if reverse {
		rowS = endRow
		colS = endCol
	}

	return modifier, rowS, colS
}
