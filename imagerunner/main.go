package imagerunner

import (
	"image"
	"sync"
)

func Start(img image.Image) ([][]Pixel, int, int) {
	mockSettings := Settings{
		maxStarSizeInPx:     4,
		maxStarGlowInPx:     1,
		wcagContrastMinimum: 1.6,
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

	store.findStarCenters()
	store.markStarRadiusAsStar()

	return store.Pixels, store.Width, store.Height
}
