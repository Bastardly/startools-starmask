package imagerunner

import (
	"image"
)

func Start(img image.Image) ([][]Pixel, int, int) {
	pixels, width, height := getPixels(img)

	// Todo, make this a parameter and pass these values from UI
	mockSettings := Settings{
		maxStarSizeInPx:     5,
		maxStarGlowInPx:     2,
		wcagContrastMinimum: 4,
	}

	// Todo go routines and waitGroups
	mapAlphaAreasForwardHorizontal(pixels, width, height, mockSettings)
	// mapContrastAreasBackwardHorizontal(pixels, width, height, mockSettings)

	return pixels, width, height

	// // Nope, brug alm "for-loop", for du skal hoppe mellem rækker og bruge index. Her skal du bruge height/width
	// for _, row := range pixels {
	// 	for _, pixel := range row {

	// 	}
	// }
}

// Run through rows x 4 Find contrast change... mark true until limit is reached, or that it changes natually
// (If the center from all axis is true, it does not matter since we will overwrite the center with surrounding star color.)
// Redo this from left, top, bottom, right.
// If all those values are true, that pixel is a star. Huzzah!
// OBS Make contrast detection fast! This will be the workhorse.

// Rest with old ideas...

// If we find contrast change, we start exploring in x, y
// If it is not currently hasContrastChange, and we notice contrast change in all directions, and  and within limit, it's a star
// If star - We then find the width of the star, and the height (both up and down, since we start unevenly) of the star.
// We then het height above entry point, and below, and width, to calculate center of star.
// From here we find the correct height and width of the star, and explore values within that rectangle or circle (circle + 1 should be best to avoid corners)
// We then mark all the bright pixes in that rectangle as hasContrastChange.
// We then overwrite the bright pixels with the surrounding pixel values (Perhaps mix them up a bit.)
// We continue.
// All hasContrastChange can be used to create an alpha map
