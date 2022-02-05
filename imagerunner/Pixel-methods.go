package imagerunner

import (
	"math"
	"starkiller/colortools"
)

func (p Pixel) IsValid() bool {
	return p.HasContrastChangeHorizontal && p.HasContrastChangeVertical
}

// calculateStarRadiusWithGlow finds the largest radius for the star, including halo
func (p *Pixel) calculateStarRadiusWithGlow(maxStarGlowInPx float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.starRadius = p.starCoreRadius + maxStarGlowInPx
}

func (p *Pixel) getSquareMapCoords(store Store, row, col int) (int, int, int, int) {
	radius := int(p.starRadius)

	startRow := validateCoord(row-radius, store.Height)
	startCol := validateCoord(col-radius, store.Width)
	endRow := validateCoord(row+radius, store.Height)
	endCol := validateCoord(col+radius, store.Width)

	return startRow, startCol, endRow, endCol
}

func (p *Pixel) markAsStarIfWithinRange(centerRow, centerCol, starRow, starCol int, starRadius, starCoreRadius float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	x := math.Abs(float64(centerCol - starCol))
	y := math.Abs(float64(centerRow - starRow))
	// Pythagoras
	a2 := math.Pow(x, 2)
	b2 := math.Pow(y, 2)
	c2 := math.Sqrt(a2 + b2)

	p.IsStar = p.IsStar || c2 <= starCoreRadius
	if c2 < starRadius {
		glowMaxDistance := starRadius - starCoreRadius
		glowDistance := c2 - starCoreRadius
		p.glowStrength = getRoundedFalloff(glowMaxDistance, glowDistance)
	}
}

func (p *Pixel) modifyColors(procentage float64, pxColor, srcColor Color) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.IsStar {
		procentage *= p.glowStrength
	}
	p.R = colortools.ChannelBlendByProcentage(procentage, pxColor.R, srcColor.R)
	p.G = colortools.ChannelBlendByProcentage(procentage, pxColor.G, srcColor.G)
	p.B = colortools.ChannelBlendByProcentage(procentage, pxColor.B, srcColor.B)
}

func (p *Pixel) getColor() Color {
	return Color{R: p.R, G: p.G, B: p.B}
}
