package imagerunner

import (
	"fmt"
	"image"
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

func Start(img image.Image) ([][]Pixel, int, int) {
	initialSettings := Settings{
		starRadiusModifier:  3,
		maxStarSizeInPx:     10,
		minStarSizeInPx:     7,
		maxStarGlowInPx:     8,
		wcagContrastMinimum: 1.7,
		radialMaskStrength: 0.3,
		blendMode:           "cloneStamp",
	}
	initialSettingsSmall := Settings{
		starRadiusModifier:  2,
		maxStarSizeInPx:     7,
		minStarSizeInPx:     4,
		maxStarGlowInPx:     4,
		wcagContrastMinimum: 1.7,
		radialMaskStrength: 0.3,
		blendMode:           "cloneStamp",
	}
	initialSettingsSmall2 := Settings{
		starRadiusModifier:  1,
		maxStarSizeInPx:     4,
		minStarSizeInPx:     1,
		maxStarGlowInPx:     2,
		wcagContrastMinimum: 1.7,
		radialMaskStrength: 0.3,
		blendMode:           "cloneStamp",
	}

	var store = Store{
		0, 0,
		initialSettings,
		[][]Pixel{},
	}

	// todo - Validate settings. If stars are too small, no need to run it twice.

	store.fillStore(img)

	run(store)
	store.clearStars()
	store.settings = initialSettingsSmall
	run(store)
	store.clearStars()
	store.settings = initialSettingsSmall2
	run(store)
	store.clearStars()
	// Then we run it again, and remove larger files.
	return run(store)

}
