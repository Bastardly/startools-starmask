package imagerunner

// getFallBackColor return the darkest color of the four corners. Hopefulle this will suffice to avoid bright artefacts.
// func (store Store) getFallBackColor(startRow, startCol, endRow, endCol int) Color {
// 	dark := store.Pixels[startRow][startCol]
// 	p2 := store.Pixels[startRow][endCol]
// 	p3 := store.Pixels[endRow][startCol]
// 	p4 := store.Pixels[endRow][endCol]

// 	if dark.brightness > p2.brightness {
// 		dark = p2
// 	}
// 	if dark.brightness > p3.brightness {
// 		dark = p3
// 	}
// 	if dark.brightness > p4.brightness {
// 		dark = p4
// 	}

// 	return Color{
// 		R: dark.R,
// 		G: dark.G,
// 		B: dark.B,
// 	}
// }

// func (store Store) setOppositeColorRow(row, col, modifier int, firstColor, fallbackColor Color) {
// 	oppositeColor := fallbackColor
// 	steps := 1
// 	nextRow := row

// 	for {
// 		nextRow = row + (steps * modifier)
// 		steps++

// 		if nextRow >= store.Height || nextRow < 0 {
// 			break
// 		}

// 		next := store.Pixels[nextRow][col]

// 		if !next.IsStar || next.IsMapped {
// 			oppositeColor = Color{
// 				R: next.R,
// 				G: next.G,
// 				B: next.B,
// 			}
// 			break
// 		}
// 	}

// 	store.Pixels[row][col].setMixedColor(firstColor, oppositeColor, steps)
// }

// func (store Store) setOppositeColorCol(row, col, modifier int, firstColor, fallbackColor Color) {
// 	oppositeColor := fallbackColor
// 	steps := 1
// 	nextCol := col

// 	// We search for the real opposite color
// 	for {
// 		nextCol = col + (steps * modifier)
// 		steps++

// 		if nextCol >= store.Width || nextCol < 0 {
// 			break
// 		}

// 		next := store.Pixels[row][nextCol]

// 		if !next.IsStar || next.IsMapped {
// 			oppositeColor = Color{
// 				R: next.R,
// 				G: next.G,
// 				B: next.B,
// 			}
// 			break
// 		}
// 	}

// 	store.Pixels[row][col].setMixedColor(firstColor, oppositeColor, steps)
// }

// func (store Store) setPixelMaskColorFromDirection(row, col int, direction string, fallbackColor Color) {
// 	firstColor := fallbackColor
// 	oppositeColor := fallbackColor
// 	switch direction {
// 	case "right":
// 		{
// 			// First color, if we go right, we look up
// 			if row-1 > 0 {
// 				firstColor = store.getPixelColorFromCoords(row-1, col)
// 			}
// 			store.setOppositeColorRow(row, col, 1, firstColor, oppositeColor)
// 		}

// 	case "left":
// 		{
// 			// First color, if we go left, we look down
// 			if row+1 < store.Height {
// 				firstColor = store.getPixelColorFromCoords(row+1, col)
// 			}
// 			store.setOppositeColorRow(row, col, -1, firstColor, oppositeColor)
// 		}

// 	case "up":
// 		{
// 			// First color, if we go up, we look left
// 			if col-1 > 0 {
// 				firstColor = store.getPixelColorFromCoords(row, col-1)
// 			}
// 			store.setOppositeColorCol(row, col, 1, firstColor, oppositeColor)
// 		}

// 	case "down":
// 		{
// 			// First color, if we go down, we look right (No fucking pun intended)
// 			if col+1 > store.Width {
// 				firstColor = store.getPixelColorFromCoords(row, col+1)
// 			}
// 			store.setOppositeColorCol(row, col, -1, firstColor, oppositeColor)
// 		}
// 	}
// }

// // Tried to run from the glow round and in. But this was too slow.
// func (store Store) maskStarArea2(startRow, startCol, endRow, endCol int) {
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

// //////// PX

// func (p *Pixel) setMixedColor(firstColor, oppositeColor Color, steps int) {
// 	color := colortools.MixColors(firstColor, oppositeColor, steps)
// 	p.R = color.R
// 	p.G = color.G
// 	p.B = color.B
// 	p.IsMapped = true
// }

// ////////////// store halo

// func (store Store) getHaloFalloffLengthVertical(row, col, modifier, starCoreRadius int) float64 {
// 	prevBrightness := store.Pixels[row][col].brightness
// 	result := 0

// 	if store.settings.maxStarGlowInPx > 0 {
// 		for i := 1; i <= store.settings.maxStarGlowInPx; i++ {
// 			result = i
// 			nextRow := row + i*modifier
// 			if nextRow < 0 || nextRow >= store.Height {
// 				break
// 			}
// 			nextBrightness := store.Pixels[nextRow][col].brightness
// 			if nextBrightness >= prevBrightness {
// 				// it's not a deminishing halo
// 				break
// 			}
// 			prevBrightness = nextBrightness
// 		}
// 	}

// 	return float64(result + starCoreRadius + extraPixel)
// }

// func (store Store) getHaloFalloffLengthHorisontal(row, col, modifier, starCoreRadius int) float64 {
// 	prevBrightness := store.Pixels[row][col].brightness
// 	result := 0

// 	if store.settings.maxStarGlowInPx > 0 {
// 		for i := 1; i <= store.settings.maxStarGlowInPx; i++ {
// 			result = i
// 			nextCol := col + i*modifier
// 			if nextCol < 0 || nextCol >= store.Width {
// 				break
// 			}
// 			nextBrightness := store.Pixels[row][nextCol].brightness
// 			if nextBrightness >= prevBrightness {
// 				// it's not a deminishing halo
// 				break
// 			}
// 			prevBrightness = nextBrightness
// 		}
// 	}

// 	extra := starCoreRadius / 2

// 	return float64(result + starCoreRadius + extra)
// }

// calculateStarRadiusWithGlow finds the largest radius for the star, including halo
// func (p *Pixel) calculateStarRadiusWithGlow(store Store, row, col int) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// directions := []uint8{1, 2, 3, 4} // up, right, down, left
// var wg sync.WaitGroup

// for _, direction := range directions {
// 	wg.Add(1)

// 	go func(localDirection uint8) {
// 		defer wg.Done()
// 		switch localDirection {
// 		case 1: // up
// 			{

// 				coreRadius := row - p.starRadiusStartVertical
// 				starRadius := store.getHaloFalloffLengthVertical(row, col, -1, coreRadius)

// 				if starRadius > p.starRadius {
// 					p.starRadius = starRadius
// 				}
// 			}

// case 2: // right
// 	{
// 		coreRadius := p.starRadiusEndHorizontal - col
// 		starRadius := store.getHaloFalloffLengthHorisontal(row, col, 1, coreRadius)
// 		if starRadius > p.starRadius {
// 			p.starRadius = starRadius
// 		}
// 	}

// case 3: // down
// 	{
// 		coreRadius := p.starRadiusEndVertical - row
// 		starRadius := store.getHaloFalloffLengthVertical(row, col, 1, coreRadius)
// 		if starRadius > p.starRadius {
// 			p.starRadius = starRadius
// 		}
// 	}

// 		case 4: // left {
// 			coreRadius := col - p.starRadiusStartHorizontal
// 			starRadius := store.getHaloFalloffLengthHorisontal(row, col, -1, coreRadius)
// 			if starRadius > p.starRadius {
// 				p.starRadius = starRadius
// 			}
// 		}
// 	}(direction)
// }

// wg.Wait()
// }
