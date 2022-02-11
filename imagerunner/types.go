package imagerunner

import (
	"starkiller/colortools"
	"sync"
)

type Pixel struct {
	R                           uint32
	G                           uint32
	B                           uint32
	A                           uint32
	brightness                  float32 // WCAG brightness made from 8bit estimate
	glowStrength                float64 // Procentage marked from 0 - 1 - We use this to calculate how hard it needs to be masked.
	radialGlowStrength          float64 // Procentage marked from 0 - 1 - We use this to calculate how hard it needs to be masked in radial blending mode.
	HasContrastChangeHorizontal bool
	HasContrastChangeVertical   bool
	HasBeenExplored             bool
	IsStar                      bool
	IsMapped                    bool // When we mask a star, the mask color is available. So we mark this star as mapped. So we can use the color
	isStarCenter                bool // Center of star, from where we will calculate it's final size with glow
	starRadiusStartHorizontal   int
	starRadiusEndHorizontal     int
	starRadiusStartVertical     int
	starRadiusEndVertical       int
	starCoreRadius              float64
	starRadius                  float64
	mu                          sync.Mutex
}

type Settings struct {
	maxStarSizeInPx     int
	minStarSizeInPx     int
	maxStarGlowInPx     float64
	radialMaskStrength  float64 // Procentage marked from 0 - 1 - We use this to calculate the maximum power of the radial mask.
	starRadiusModifier  float64
	wcagContrastMinimum float32 // The highter the star contrast is needed for detection.
	blendMode           string  // "fast" or "cloneStamp"
}

type Pixels = [][]Pixel

type Store struct {
	Width    int
	Height   int
	settings Settings
	Pixels   [][]Pixel
}

type ColorCoord struct {
	row int
	col int
}

type Color = colortools.Color

type StarFrame struct {
	startRow, startCol, endRow, endCol int
}
