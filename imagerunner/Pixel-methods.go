package imagerunner

import (
	"math"
	"starkiller/colortools"
	"sync"
)

func (p Pixel) IsValid() bool {
	return p.HasContrastChangeHorizontal && p.HasContrastChangeVertical
}

// calculateStarRadiusWithGlow finds the largest radius for the star, including halo
func (p *Pixel) calculateStarRadiusWithGlow(store Store, row, col int) {
	directions := []uint8{1, 2, 3, 4} // up, right, down, left

	var wg sync.WaitGroup

	for _, direction := range directions {
		wg.Add(1)

		go func(localDirection uint8) {
			defer wg.Done()
			switch localDirection {
			case 1: // up
				{

					coreRadius := row - p.starRadiusStartVertical
					starRadius := store.getHaloFalloffLengthVertical(row, col, -1, coreRadius)

					if starRadius > p.starRadius {
						p.starRadius = starRadius
					}
				}

			case 2: // right
				{
					coreRadius := p.starRadiusEndHorizontal - col
					starRadius := store.getHaloFalloffLengthHorisontal(row, col, 1, coreRadius)
					if starRadius > p.starRadius {
						p.starRadius = starRadius
					}
				}

			case 3: // down
				{
					coreRadius := p.starRadiusEndVertical - row
					starRadius := store.getHaloFalloffLengthVertical(row, col, 1, coreRadius)
					if starRadius > p.starRadius {
						p.starRadius = starRadius
					}
				}

			case 4: // left {
				coreRadius := col - p.starRadiusStartHorizontal
				starRadius := store.getHaloFalloffLengthHorisontal(row, col, -1, coreRadius)
				if starRadius > p.starRadius {
					p.starRadius = starRadius
				}
			}
		}(direction)
	}

	wg.Wait()
}

func (p *Pixel) getSquareMapCoords(store Store, row, col int) (int, int, int, int) {

	startRow := validateCoord(row-p.starRadius, store.Height)
	startCol := validateCoord(col-p.starRadius, store.Width)
	endRow := validateCoord(row+p.starRadius, store.Height)
	endCol := validateCoord(col+p.starRadius, store.Width)

	return startRow, startCol, endRow, endCol
}

func (p *Pixel) markAsStarIfWithinRange(centerRow, centerCol, starRadius, starRow, starCol int) {

	x := math.Abs(float64(centerCol - starCol))
	y := math.Abs(float64(centerRow - starRow))
	a2 := math.Pow(x, 2)
	b2 := math.Pow(y, 2)
	distance := math.Sqrt(a2 + b2)

	p.IsStar = float64(starRadius) < distance
}

func (p *Pixel) modifyColors(procentage float64, pxR, pxG, pxB, scR, scG, scB uint32) {

	p.R = colortools.ChannelBlendByProcentage(procentage, pxR, scR)
	p.G = colortools.ChannelBlendByProcentage(procentage, pxG, scG)
	p.B = colortools.ChannelBlendByProcentage(procentage, pxB, scB)
}
