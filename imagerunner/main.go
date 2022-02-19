package imagerunner

import (
	"fmt"
	"sync"
)

func run(store Store) ([][]Pixel, int, int) {
	var wg sync.WaitGroup

	wg.Add(2)
	go store.mapAlphaAreasHorizontal(&wg)
	go store.mapAlphaAreasVertical(&wg)

	wg.Wait()
	fmt.Println("mapAlphaAreas done")

	store.findStarCenters()
	fmt.Println("findStarCenters done")
	store.markStarRadiusAsStar()
	fmt.Println("markStarRadiusAsStar done")

	return store.Pixels, store.Width, store.Height
}

func Start(data IStart) ([][]Pixel, [][]Pixel, int, int) {
	fmt.Println("Starting imagerunner")
	// initialSettings := Settings{
	// 	starRadiusModifier:  3,
	// 	maxStarSizeInPx:     10,
	// 	minStarSizeInPx:     7,
	// 	maxStarGlowInPx:     8,
	// 	wcagContrastMinimum: 1.7,
	// 	radialMaskStrength:  0.3,
	// 	blendMode:           "cloneStamp",
	// }
	// initialSettingsSmall := Settings{
	// 	starRadiusModifier:  2,
	// 	maxStarSizeInPx:     7,
	// 	minStarSizeInPx:     4,
	// 	maxStarGlowInPx:     4,
	// 	wcagContrastMinimum: 1.7,
	// 	radialMaskStrength:  0.3,
	// 	blendMode:           "cloneStamp",
	// }
	// initialSettingsSmall2 := Settings{
	// 	starRadiusModifier:  1,
	// 	maxStarSizeInPx:     4,
	// 	minStarSizeInPx:     1,
	// 	maxStarGlowInPx:     2,
	// 	wcagContrastMinimum: 1.7,
	// 	radialMaskStrength:  0.3,
	// 	blendMode:           "cloneStamp",
	// }

	settings := getSettings(data)

	var store = Store{
		0, 0,
		settings[0],
		[][]Pixel{},
	}

	// todo - Validate settings. If stars are too small, no need to run it twice.

	store.fillStore(*data.ImageFile)
	alphaPixels := store.Pixels

	fmt.Println("Creating alpha", data)
	if data.CreateAlpha {
		store.settings.minStarSizeInPx = data.MinStarSize
		run(store)
		alphaPixels = store.Pixels
	}

	// run(store)
	// store.clearStars()
	// store.settings = initialSettingsSmall
	// run(store)
	// store.clearStars()
	// store.settings = initialSettingsSmall2
	// run(store)
	// store.clearStars()
	// Then we run it again, and remove larger files.
	// return run(store)

	if data.RemoveStars {
		for _, setting := range settings {
			store.settings = setting
			run(store)
			store.clearStars()
		}
	}

	return store.Pixels, alphaPixels, store.Width, store.Height

}
