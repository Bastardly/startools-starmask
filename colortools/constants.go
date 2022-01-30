package colortools

/**
 * These are the values to calculate the WCAG contrast ratio.
 * This is normally used to calculate proper contrast for text. Which is at least 4,5 for headers, and 7 for text.
 * For stars, about 1.6 is a good mark.
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