package colortools

import "math"

type ColorPixel struct {
	R uint8
	G uint8
	B uint8
}

/**
 * Hvis vi arbejder med få farver, skal de første farver have en god kontrast.
 * AA er en ration på over 4,5 (Stor tekst)
 * AAA er en ratio på over 7 (Lille tekst)
 * https://webaim.org/resources/contrastchecker/
 * https://dev.to/alvaromontoro/building-your-own-color-contrast-checker-4j7o
 */

// Constants to calculate WCAG ratios
const comparer float32 = 0.03928
const divider float32 = 12.92
const bias float32 = 0.055
const biasDivider float32 = 1.055
const power float64 = 2.4
const rBrightnessFactor float32 = 0.2126
const gBrightnessFactor float32 = 0.7152
const bBrightnessFactor float32 = 0.0722
const luminanceBias float32 = 0.05

//  https://stackoverflow.com/questions/9733288/how-to-programmatically-calculate-the-contrast-ratio-between-two-colors
func getBrightnessFactor(v uint8) float32 {
	V := float32(v / 255)

	if V <= comparer {
		return V / divider
	}

	return float32(math.Pow(float64((V+bias)/biasDivider), power))
}

// GetRGBBrightNess Tager farvelysstyrken i en skala fra 0 - 1 ved hver farve, og returnerer WCAG ratio
func GetRGBBrightNess(color ColorPixel) float32 {
	r := getBrightnessFactor(color.R) * rBrightnessFactor
	g := getBrightnessFactor(color.G) * gBrightnessFactor
	b := getBrightnessFactor(color.B) * bBrightnessFactor

	return r + g + b
}

// GetContrastRatio tager farvelysstyrken i en skala fra 0 - 1 ved hver farve, og returnerer WCAG ratio
func GetContrastRatio(brightness1, brightness2 float32) float32 {
	b1 := brightness1 + luminanceBias
	b2 := brightness2 + luminanceBias
	if brightness1 > brightness2 {
		return b1 / b2
	}

	return b2 / b1
}
