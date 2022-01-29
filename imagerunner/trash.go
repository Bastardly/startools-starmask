package imagerunner

// func maskPixels(pixels [][]Pixel, glowStartIndex, glowEndIndex, starIndexStart, starIndexEnd, row int, isVertical bool) {
// 	for i := glowStartIndex; i <= glowEndIndex; i++ {
// 		pixels[row][i].IsGlow = true
// 		if i >= starIndexStart && i <= starIndexEnd {

// 			if isVertical {
// 				pixels[row][i].HasContrastChangeVertical = true
// 			} else {
// 				pixels[row][i].HasContrastChangeHorizontal = true

// 			}
// 		}
// 	}
// }

// func findGlowAndMaskPixels(pixels [][]Pixel, starIndexStart, starIndexEnd, maxStarGlowInPx, row, column, width int, isVertical bool) {
// 	starSize := starIndexEnd - starIndexStart
// 	if starSize < 0 {
// 		panic("Star start index is larger than end index. ")
// 	}

// 	alphaStartIndex := starIndexStart
// 	alphaEndIndex := starIndexEnd
// 	expandBy := int(starSize / 4) // We add a bit more than the glow depending on the star
// 	startActive := true
// 	endActive := true
// 	prevStartBrightness := pixels[row][starIndexStart].brightness
// 	prevEndBrightness := pixels[row][starIndexEnd].brightness

// 	for i := 1; i <= maxStarGlowInPx; i++ {
// 		if !startActive && !endActive {
// 			break
// 		}
// 		s := starIndexStart - i
// 		e := starIndexEnd + i
// 		if startActive && s >= 0 {
// 			pixel := pixels[row][s]
// 			if pixel.brightness < prevStartBrightness {
// 				alphaStartIndex = s
// 			} else {
// 				startActive = false
// 			}
// 		}

// 		if endActive && e < width {
// 			pixel := pixels[row][e]
// 			if pixel.brightness < prevEndBrightness {
// 				alphaEndIndex = e
// 			} else {
// 				endActive = false
// 			}
// 		}
// 	}

// 	alphaStartIndex -= expandBy
// 	alphaEndIndex += expandBy

// 	if alphaStartIndex < 0 {
// 		alphaStartIndex = 0
// 	}
// 	if alphaEndIndex >= width {
// 		alphaEndIndex = width - 1
// 	}

// 	maskPixels(pixels, alphaStartIndex, alphaEndIndex, starIndexStart, starIndexEnd, row, isVertical)
// }
