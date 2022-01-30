package colortools

import (
	"math"
)

//  https://stackoverflow.com/questions/9733288/how-to-programmatically-calculate-the-contrast-ratio-between-two-colors
func getBrightnessFactor(v uint8) float32 {
	V := float32(v) / float32(255)

	if V <= comparer {

		return V / divider
	}

	return float32(math.Pow(float64((V+bias)/biasDivider), power))
}

// GetRGBBrightNess calculates the brightness at a scale from 0 - 1, and returns the sum of all RGB colors
func GetRGBBrightNess(R, G, B uint8) float32 {
	r := getBrightnessFactor(R) * rBrightnessFactor
	g := getBrightnessFactor(G) * gBrightnessFactor
	b := getBrightnessFactor(B) * bBrightnessFactor

	return r + g + b
}

// GetContrastRatio tager farvelysstyrken i en skala fra 0 - 1, and returns the WCAG contrast ratio (Scale for natual contrast for the naked eye)
func GetContrastRatio(brightness1, brightness2 float32) float32 {
	b1 := brightness1 + luminanceBias
	b2 := brightness2 + luminanceBias
	if brightness1 > brightness2 {
		return b1 / b2
	}

	return b2 / b1
}
