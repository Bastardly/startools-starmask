package imagerunner

import (
	"math"
	"starkiller/colortools"
	"sync"
)

func rgbaToPixel(R uint32, G uint32, B uint32, A uint32) Pixel {
	r := uint8(R / bitDivider)
	g := uint8(G / bitDivider)
	b := uint8(B / bitDivider)

	return Pixel{
		R:                           R,
		G:                           G,
		B:                           B,
		A:                           A,
		brightness:                  colortools.GetRGBBrightNess(r, g, b),
		glowStrength:                0,
		HasContrastChangeHorizontal: false,
		HasContrastChangeVertical:   false,
		HasBeenExplored:             false,
		IsStar:                      false,
		IsMapped:                    false,
		isStarCenter:                false,
		starRadiusStartHorizontal:   -1,
		starRadiusEndHorizontal:     -1,
		starRadiusStartVertical:     -1,
		starRadiusEndVertical:       -1,
		starCoreRadius:              -1,
		starRadius:                  -1,
		mu:                          sync.Mutex{},
	}
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

func getRoundedFalloff(max, position float64) float64 {
	return math.Sin((max - position) / max)
}
