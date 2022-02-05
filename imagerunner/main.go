package imagerunner

import (
	"fmt"
	"image"
	"sync"
)

func Start(img image.Image) ([][]Pixel, int, int) {
	mockSettings := Settings{
		starRadiusModifier:  2,
		maxStarSizeInPx:     5,
		maxStarGlowInPx:     2,
		wcagContrastMinimum: 1.3,
	}

	var store = Store{
		0, 0,
		mockSettings,
		[][]Pixel{},
	}

	store.fillStore(img)

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
