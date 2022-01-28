package createfile

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"starkiller/imagerunner"
)

func CreateAlpha(pixels [][]imagerunner.Pixel, uri string, width, height int) {
	dir := filepath.Dir(uri)
	fileName := filepath.Join(dir, "alpha.tif")
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[y][x]

			switch {
			case pixel.HasContrastChangeLeft:
				img.Set(x, y, color.White)
			default:
				img.Set(x, y, color.Black)
			}
		}
	}

	// Encode as Tiff.
	f, _ := os.Create(fileName)
	png.Encode(f, img)
}
