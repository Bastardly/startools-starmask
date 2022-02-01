package createfile

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"starkiller/imagerunner"
)

func CreateAlpha(pixels [][]imagerunner.Pixel, uri string, width, height int) {
	dir := filepath.Dir(uri)
	fileName := filepath.Join(dir, "alpha.png")
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[y][x]

			t := color.RGBA{255, 0, 0, 255}
			w := color.RGBA{255, 255, 255, 255}
			b := color.RGBA{0, 0, 0, 255}

			switch {
			case pixel.IsStar:
				img.Set(x, y, t)
			case pixel.IsValid():
				img.Set(x, y, w)
			default:
				img.Set(x, y, b)
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create(fileName)
	png.Encode(f, img)
	f.Close()
}

func CreatePreview(pixels [][]imagerunner.Pixel, uri string, width, height int) {
	dir := filepath.Dir(uri)
	fileName := filepath.Join(dir, "preview.png")
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[y][x]

			r := uint8(pixel.R / 257)
			g := uint8(pixel.G / 257)
			b := uint8(pixel.B / 257)
			a := uint8(pixel.A / 257)

			color := color.RGBA{r, g, b, a}

			img.Set(x, y, color)

		}
	}

	// Encode as PNG.
	f, _ := os.Create(fileName)
	png.Encode(f, img)
	f.Close()

	fmt.Println("DONE")
}
