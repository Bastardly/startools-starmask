package imagerunner

import (
	"image"
	"starkiller/colortools"
)

const bitDivider = 257 // 16 bit to 8 bit divider

// Pixel struct example


func rgbaToPixel(R uint32, G uint32, B uint32, A uint32) Pixel {
	r := uint8(R / bitDivider)
	g := uint8(G / bitDivider)
	b := uint8(B / bitDivider)
	brightness := colortools.GetRGBBrightNess(colortools.ColorPixel{r, g, b})

	return Pixel{R, G, B, A, brightness, false, false, false}
}

// Get the bi-dimensional pixel array and image width and height
func getPixels(img image.Image) ([][]Pixel, int, int) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, width, height
}
